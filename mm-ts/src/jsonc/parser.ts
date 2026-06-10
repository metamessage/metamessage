import { JSONCScanner, TokenType, Token } from './scanner';
import { NodeScalar, NodeObject, NodeArray, NodeNull, MMDoc } from '../ir/ast';
import { ValidationResult } from '../ir/tag';
import { ValueType } from '../ir/value-type';
import { Tag, parseMMTag } from '../ir/tag';
import { Node } from '../ir/ast';
import { base64ToUint8 } from './printer';

const maxDepth = 32;

export class JSONCParser {
  private scanner: JSONCScanner;
  private currentToken: Token;
  private pendingComments: Token[] = [];
  private depth = 0;

  constructor(input: string) {
    this.scanner = new JSONCScanner(input);
    this.currentToken = this.scanner.nextToken();
  }

  parse(): MMDoc {
    let val: Node | null = null;
    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        if (val === null) {
          throw new Error('no value parsed');
        }
        break;
      }

      if (tok.type === TokenType.COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (last && tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      val = this.parseValue('');
    }
    return new MMDoc(val || new NodeObject());
  }

  private peek(): Token {
    return this.currentToken;
  }

  private next(): Token {
    const t = this.currentToken;
    this.currentToken = this.scanner.nextToken();
    return t;
  }

  private consumeCommentsFor(anchorLine: number): Tag | null {
    if (this.pendingComments.length === 0) {
      return null;
    }

    const last = this.pendingComments[this.pendingComments.length - 1];
    if (last && anchorLine - last.line > 1) {
      this.pendingComments = [];
      return null;
    }

    let merged: Tag | null = null;
    for (const ct of this.pendingComments) {
      const parsed = this.parseCommentsToTag(ct.value);
      if (parsed) {
        merged = this.mergeTag(merged, parsed);
      }
    }
    this.pendingComments = [];
    return merged;
  }

  private parseValue(path: string, existingTag?: Tag): Node | null {
    while (true) {
      const tok = this.next();
      let data: any;
      let text: string;

      switch (tok.type) {
        case TokenType.EOF:
          return null;

        case TokenType.LCURLY:
          return this.parseObject(tok.line, path, existingTag);

        case TokenType.LBRACKET:
          return this.parseArray(tok.line, path, existingTag);

        case TokenType.STRING:
          let strTag = existingTag || this.consumeCommentsFor(tok.line);
          text = tok.value;

          if (!strTag) {
            strTag = new Tag();
          }

          if (strTag.type === ValueType.Unknown) {
            strTag.type = ValueType.Str;
          }

          switch (strTag.type) {
            case ValueType.Str:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid string: "${text}", valid: ""`);
                }
                data = '';
              } else {
                const result = strTag.validateStr(text);
                if (!result.valid) {
                  throw new Error(result.error || 'String validation failed');
                }
                data = result.data;
                text = result.text || text;
              }
              break;

            case ValueType.Bytes:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid bytes: "${text}", valid: ""`);
                }
                data = new Uint8Array();
              } else {
                try {
                  const decoded = base64ToUint8(text);
                  const result = strTag.validateBytes(decoded);
                  if (!result.valid) {
                    throw new Error(result.error || 'Bytes validation failed');
                  }
                  data = result.data;
                  text = result.text || text;
                } catch (e) {
                  throw new Error(
                    `invalid base64 bytes "${text}": ${(e as Error).message}`,
                  );
                }
              }
              break;

            case ValueType.Datetime:
            case ValueType.Date:
            case ValueType.Time:
              if (strTag.isNull) {
                data = new Date(0);
              } else {
                try {
                  let dateValue: Date;
                  if (strTag.type === ValueType.Time) {
                    dateValue = new Date(`1970-01-01T${text}Z`);
                  } else {
                    dateValue = new Date(text.replace(' ', 'T') + 'Z');
                  }
                  let result: ValidationResult;
                  if (strTag.type === ValueType.Date) {
                    result = strTag.validateDate(dateValue);
                  } else if (strTag.type === ValueType.Time) {
                    result = strTag.validateTime(dateValue);
                  } else {
                    result = strTag.validateDatetime(dateValue);
                  }
                  if (!result.valid) {
                    throw new Error(
                      result.error || 'DateTime validation failed',
                    );
                  }
                  data = result.data;
                  // Keep original text from JSONC input, not UTC-formatted result.text
                } catch (e) {
                  throw new Error(
                    `invalid datetime "${text}": ${(e as Error).message}`,
                  );
                }
              }
              break;

            case ValueType.Uuid:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid uuid: "${text}", valid: ""`);
                }
                data = '';
              } else {
                const result = strTag.validateUUID(text);
                if (!result.valid) {
                  throw new Error(result.error || 'UUID validation failed');
                }
                data = text;
                text = result.text || text;
              }
              break;

            case ValueType.Email:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid email: "${text}", valid: ""`);
                }
                data = '';
              } else {
                const result = strTag.validateEmail(text);
                if (!result.valid) {
                  throw new Error(result.error || 'Email validation failed');
                }
                data = result.data;
                text = result.text || text;
              }
              break;

            case ValueType.Enums:
              if (!strTag.enums) {
                throw new Error('enum empty');
              }
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid enums: "${text}", valid: ""`);
                }
                data = -1;
              } else {
                const result = strTag.validateEnum(text);
                if (!result.valid) {
                  throw new Error(result.error || 'Enum validation failed');
                }
                data = result.data;
                text = result.text || text;
              }
              break;

            case ValueType.Media:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid media: "${text}", valid: ""`);
                }
                data = new Uint8Array();
              } else {
                try {
                  const decoded = base64ToUint8(text);
                  const result = strTag.validateMedia(decoded);
                  if (!result.valid) {
                    throw new Error(result.error || 'Media validation failed');
                  }
                  data = result.data;
                  text = result.text || text;
                } catch (e) {
                  throw new Error(
                    `invalid base64 media "${text}": ${(e as Error).message}`,
                  );
                }
              }
              break;

            case ValueType.Url:
            case ValueType.Ip:
            case ValueType.Decimal:
            default:
              let defaultResult: ValidationResult;
              if (strTag.type === ValueType.Url) {
                defaultResult = strTag.validateURL(text);
              } else if (strTag.type === ValueType.Ip) {
                defaultResult = strTag.validateIP(text);
              } else if (strTag.type === ValueType.Decimal) {
                defaultResult = strTag.validateDecimal(text);
              } else {
                defaultResult = strTag.validateStr(text);
              }
              if (!defaultResult.valid) {
                throw new Error(
                  defaultResult.error ||
                    `Validation failed for type ${strTag.type}`,
                );
              }
              data = defaultResult.data;
              text = defaultResult.text || text;
              break;
          }

          const strValue = new NodeScalar(data, strTag);
          strValue.setPath(path);
          strValue.setText(text);
          return strValue;

        case TokenType.NUMBER:
          let numTag = existingTag || this.consumeCommentsFor(tok.line);
          text = tok.value;

          if (!numTag) {
            numTag = new Tag();
          }

          if (text.includes('.')) {
            if (numTag.type === ValueType.Unknown) {
              numTag.type = ValueType.F64;
            }

            switch (numTag.type) {
              case ValueType.F32:
                if (numTag.isNull) {
                  if (text !== '0.0') {
                    throw new Error(`invalid float32: ${text}, valid: 0.0`);
                  }
                  data = 0.0;
                } else {
                  const f64 = parseFloat(text);
                  if (isNaN(f64)) {
                    throw new Error(`invalid float32 "${text}"`);
                  }
                  const result = numTag.validateF32(f64);
                  if (!result.valid) {
                    throw new Error(
                      result.error || 'Float32 validation failed',
                    );
                  }
                  data = result.data;
                  text = result.text || text;
                }
                break;

              case ValueType.F64:
                if (numTag.isNull) {
                  if (text !== '0.0') {
                    throw new Error(`invalid float64: ${text}, valid: 0.0`);
                  }
                  data = 0.0;
                } else {
                  const f64 = parseFloat(text);
                  if (isNaN(f64)) {
                    throw new Error(`invalid float64 "${text}"`);
                  }
                  const result = numTag.validateF64(f64);
                  if (!result.valid) {
                    throw new Error(
                      result.error || 'Float64 validation failed',
                    );
                  }
                  data = result.data;
                  text = result.text || text;
                }
                break;

              default:
                const floatDefaultResult = numTag.validateF64(parseFloat(text));
                if (!floatDefaultResult.valid) {
                  throw new Error(
                    floatDefaultResult.error ||
                      `Validation failed for numeric type ${numTag.type}`,
                  );
                }
                data = floatDefaultResult.data;
                text = floatDefaultResult.text || text;
                break;
            }
          } else if (text.startsWith('-')) {
            if (numTag.type === ValueType.Unknown) {
              numTag.type = ValueType.I;
            }

            switch (numTag.type) {
              case ValueType.I8:
                data = this.parseAndValidateI(text, numTag, 'validateI8');
                break;
              case ValueType.I16:
                data = this.parseAndValidateI(text, numTag, 'validateI16');
                break;
              case ValueType.I32:
                data = this.parseAndValidateI(text, numTag, 'validateI32');
                break;
              case ValueType.I64:
                data = this.parseAndValidateI(text, numTag, 'validateI64');
                break;
              case ValueType.Bigint:
                if (numTag.isNull) {
                  if (text !== '0') {
                    throw new Error(`invalid bigint: ${text}, valid: 0`);
                  }
                  data = BigInt(0);
                } else {
                  const bi = BigInt(text);
                  const result = numTag.validateBigint(bi);
                  if (!result.valid) {
                    throw new Error(result.error || 'BigInt validation failed');
                  }
                  data = result.data;
                }
                break;
              default:
                const num = BigInt(text);
                let negResult = numTag.validateI(num);
                if (!negResult.valid) {
                  throw new Error(
                    negResult.error ||
                      `Validation failed for type ${numTag.type}`,
                  );
                }
                data = negResult.data;
                break;
            }
          } else {
            if (numTag.type === ValueType.Unknown) {
              numTag.type = ValueType.I;
            }

            switch (numTag.type) {
              case ValueType.I:
                data = this.parseAndValidateI(text, numTag, 'validateI');
                break;
              case ValueType.I8:
                data = this.parseAndValidateI(text, numTag, 'validateI8');
                break;
              case ValueType.I16:
                data = this.parseAndValidateI(text, numTag, 'validateI16');
                break;
              case ValueType.I32:
                data = this.parseAndValidateI(text, numTag, 'validateI32');
                break;
              case ValueType.I64:
                data = this.parseAndValidateI(text, numTag, 'validateI64');
                break;
              case ValueType.U:
                data = this.parseAndValidateU(text, numTag, 'validateU');
                break;
              case ValueType.U8:
                data = this.parseAndValidateU(text, numTag, 'validateU8');
                break;
              case ValueType.U16:
                data = this.parseAndValidateU(text, numTag, 'validateU16');
                break;
              case ValueType.U32:
                data = this.parseAndValidateU(text, numTag, 'validateU32');
                break;
              case ValueType.U64:
                data = this.parseAndValidateU(text, numTag, 'validateU64');
                break;
              case ValueType.Bigint:
                if (numTag.isNull) {
                  if (text !== '0') {
                    throw new Error(`invalid bigint: ${text}, valid: 0`);
                  }
                  data = BigInt(0);
                } else {
                  const bi = BigInt(text);
                  const result = numTag.validateBigint(bi);
                  if (!result.valid) {
                    throw new Error(result.error || 'BigInt validation failed');
                  }
                  data = result.data;
                }
                break;
              default:
                throw new Error(
                  `unsupported type ${numTag.type} for numeric literal`,
                );
            }
          }

          const numValue = new NodeScalar(data, numTag);
          numValue.setPath(path);
          numValue.setText(text);
          return numValue;

        case TokenType.TRUE:
          let trueTag = existingTag || this.consumeCommentsFor(tok.line);
          if (!trueTag) {
            trueTag = new Tag();
          }
          if (trueTag.type === ValueType.Unknown) {
            trueTag.type = ValueType.Bool;
          }

          if (trueTag.type === ValueType.Bool) {
            if (trueTag.isNull) {
              throw new Error('bool must false when bool is null');
            } else {
              const result = trueTag.validateBool(true);
              if (!result.valid) {
                throw new Error(result.error || 'Boolean validation failed');
              }
            }
          } else {
            throw new Error(
              `unsupported type ${trueTag.type} for boolean literal`,
            );
          }

          const trueValue = new NodeScalar(true, trueTag);
          trueValue.setPath(path);
          trueValue.setText('true');
          return trueValue;

        case TokenType.FALSE:
          let falseTag = existingTag || this.consumeCommentsFor(tok.line);
          if (!falseTag) {
            falseTag = new Tag();
          }
          if (falseTag.type === ValueType.Unknown) {
            falseTag.type = ValueType.Bool;
          }

          if (falseTag.type === ValueType.Bool) {
            if (!falseTag.isNull) {
              const result = falseTag.validateBool(false);
              if (!result.valid) {
                throw new Error(result.error || 'Boolean validation failed');
              }
            }
          } else {
            throw new Error(
              `unsupported type ${falseTag.type} for boolean literal`,
            );
          }

          const falseValue = new NodeScalar(false, falseTag);
          falseValue.setPath(path);
          falseValue.setText('false');
          return falseValue;

        case TokenType.NULL:
          let nullTag = existingTag || this.consumeCommentsFor(tok.line);
          if (!nullTag) {
            nullTag = new Tag();
          }
          if (nullTag.type !== ValueType.Unknown) {
            throw new Error(`null is not supported for type ${nullTag.type}`);
          }
          const nullNode = new NodeNull(nullTag);
          nullNode.setPath(path);
          return nullNode;

        default:
          throw new Error(`unexpected token ${tok.type}`);
      }
    }
  }

  private parseAndValidateI(
    text: string,
    tag: Tag,
    methodName: string,
  ): number {
    if (tag.isNull) {
      if (text !== '0') {
        throw new Error(`invalid ${tag.type}: ${text}, valid: 0`);
      }
      return 0;
    }

    const num = parseInt(text, 10);
    if (isNaN(num)) {
      throw new Error(`invalid ${tag.type} "${text}"`);
    }

    const result = (tag as any)[methodName](num);
    if (!result.valid) {
      throw new Error(result.error || `${tag.type} validation failed`);
    }
    return result.data;
  }

  private parseAndValidateU(
    text: string,
    tag: Tag,
    methodName: string,
  ): number {
    if (tag.isNull) {
      if (text !== '0') {
        throw new Error(`invalid ${tag.type}: ${text}, valid: 0`);
      }
      return 0;
    }

    const num = parseInt(text, 10);
    if (isNaN(num)) {
      throw new Error(`invalid ${tag.type} "${text}"`);
    }

    const result = (tag as any)[methodName](num);
    if (!result.valid) {
      throw new Error(result.error || `${tag.type} validation failed`);
    }
    return result.data;
  }

  private parseObject(
    openLine: number,
    path: string,
    parentTag?: Tag,
  ): NodeObject {
    this.depth++;
    if (this.depth > maxDepth) {
      throw new Error(`max depth: ${maxDepth}`);
    }

    let tag = parentTag || new Tag();
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Obj;
    }

    if (tag.name) {
      path = path ? `${path}.${tag.name}` : tag.name;
    }

    const obj = new NodeObject();
    obj.setTag(tag);
    obj.setPath(path);

    let val: Node | null = null;
    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        break;
      }
      if (tok.type === TokenType.RCURLY) {
        this.next();
        break;
      }

      if (tok.type === TokenType.COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (last && tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      const key = this.next();
      if (key.type !== TokenType.STRING) {
        throw new Error('expect string key');
      }
      const keyStr = camelToSnake(key.value);

      this.next();
      const pa =
        tag.type === ValueType.Map ? `${path}[${keyStr}]` : `${path}.${keyStr}`;

      let childTag: Tag | null = this.consumeCommentsFor(key.line);
      if (!childTag) {
        childTag = new Tag();
      }
      if (tag.type === ValueType.Map) {
        childTag.inherit(tag);
      }

      val = this.parseValue(pa, childTag);
      if (!val) {
        this.next();
        continue;
      }

      obj.setProperty(keyStr, val);

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    switch (tag.type) {
      case ValueType.Map: {
        const result = tag.validateMap();
        if (!result.valid) {
          throw new Error(`validate failed: ${result.error}`);
        }
        break;
      }
      case ValueType.Obj: {
        const result = tag.validateObj();
        if (!result.valid) {
          throw new Error(`validate failed: ${result.error}`);
        }
        break;
      }
    }

    this.depth--;
    return obj;
  }

  private parseArray(
    openLine: number,
    path: string,
    parentTag?: Tag,
  ): NodeArray {
    this.depth++;
    if (this.depth > maxDepth) {
      throw new Error(`max depth: ${maxDepth}`);
    }

    let tag = parentTag || this.consumeCommentsFor(openLine);
    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Vec;
    }

    if (tag.name) {
      path = `${path}.${tag.name}`;
    }

    const arr = new NodeArray();
    arr.setTag(tag);
    arr.setPath(path);

    let item: Node | null = null;
    let i = 0;
    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        break;
      }
      if (tok.type === TokenType.RBRACKET) {
        this.next();
        break;
      }

      if (tok.type === TokenType.COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (last && tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      const pa = `${path}[${i}]`;

      let childTag: Tag | null = null;
      if (openLine !== tok.line) {
        childTag = this.consumeCommentsFor(tok.line);
      }
      if (!childTag) {
        childTag = new Tag();
      }
      childTag.inherit(tag);

      item = this.parseValue(pa, childTag);
      if (!item) {
        continue;
      }

      if (item instanceof NodeScalar) {
        const itemTag = item.getTag();
        if (
          itemTag.type === ValueType.Datetime ||
          itemTag.type === ValueType.Date ||
          itemTag.type === ValueType.Time
        ) {
          if (typeof item.getValue() === 'string') {
            const text = item.getText();
            let dateValue: Date;
            if (itemTag.type === ValueType.Time) {
              dateValue = new Date(`1970-01-01T${text}Z`);
            } else {
              dateValue = new Date(text.replace(' ', 'T') + 'Z');
            }
            if (!isNaN(dateValue.getTime())) {
              item.setValue(dateValue);
            }
          }
        }
      }

      arr.addElement(item);
      i++;

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    switch (tag.type) {
      case ValueType.Arr: {
        const result = tag.validateArr(arr.getElements());
        if (!result.valid && !tag.example) {
          throw new Error(`validate failed: ${result.error}`);
        }
        break;
      }
      case ValueType.Vec: {
        const result = tag.validateVec(arr.getElements());
        if (!result.valid && !tag.example) {
          throw new Error(`validate failed: ${result.error}`);
        }
        break;
      }
    }

    this.depth--;
    return arr;
  }

  private mergeTag(a: Tag | null, b: Tag): Tag {
    if (!a) {
      return b;
    }
    if (a.type !== ValueType.Unknown) {
      b.type = a.type;
    }
    if (a.desc) {
      b.desc = a.desc;
    }
    if (a.nullable) {
      b.nullable = true;
    }
    if (a.isNull) {
      b.isNull = true;
    }
    if (a.default_val) {
      b.default_val = a.default_val;
    }
    if (a.min) {
      b.min = a.min;
    }
    if (a.max) {
      b.max = a.max;
    }
    if (a.size !== 0n) {
      b.size = a.size;
    }
    if (a.enums) {
      b.enums = a.enums;
    }
    if (a.pattern) {
      b.pattern = a.pattern;
    }
    if (a.location) {
      b.location = a.location;
    }
    if (a.version !== 0) {
      b.version = a.version;
    }
    if (a.mime) {
      b.mime = a.mime;
    }
    return b;
  }

  private parseCommentsToTag(cs: string): Tag | null {
    if (!cs.startsWith('mm:')) return null;
    const tagStr = cs.substring(3).trim();
    if (!tagStr) return null;
    return parseMMTag(tagStr);
  }
}

function camelToSnake(s: string): string {
  if (s === '') return '';
  let result = '';
  for (let i = 0; i < s.length; i++) {
    const ch = s[i] as string;
    if (ch >= 'A' && ch <= 'Z') {
      if (
        i > 0 &&
        (!((s[i - 1] as string) >= 'A' && (s[i - 1] as string) <= 'Z') ||
          (i + 1 < s.length &&
            !((s[i + 1] as string) >= 'A' && (s[i + 1] as string) <= 'Z')))
      ) {
        result += '_';
      }
      result += ch.toLowerCase();
    } else {
      result += ch;
    }
  }
  return result;
}

export function parseJSONC(input: string): MMDoc {
  const parser = new JSONCParser(input);
  return parser.parse();
}
