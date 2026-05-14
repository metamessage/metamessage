import { JSONCScanner, TokenType, Token } from './scanner';
import { MMValue, MMObject, MMArray, MMDoc } from '../ast/ast';
import { ValidationResult } from '../ast/tag';
import { ValueType } from '../ast/value-type';
import { Tag, parseMMTag } from '../ast/tag';
import { Node } from '../ast/ast';
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
        break;
      }

      if (
        tok.type === TokenType.LINECOMMENT ||
        tok.type === TokenType.BLOCKCOMMENT
      ) {
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
    return new MMDoc(val || new MMObject());
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

  private parseValue(path: string): Node | null {
    this.depth++;
    if (this.depth > maxDepth) {
      throw new Error(`max depth: ${maxDepth}`);
    }

    while (true) {
      const tok = this.next();
      let data: any;
      let text: string;

      switch (tok.type) {
        case TokenType.EOF:
          this.depth--;
          return null;

        case TokenType.LCURLY:
          const objResult = this.parseObject(tok.line, path);
          this.depth--;
          return objResult;

        case TokenType.LBRACKET:
          const arrResult = this.parseArray(tok.line, path);
          this.depth--;
          return arrResult;

        case TokenType.STRING:
          let strTag = this.consumeCommentsFor(tok.line);
          text = tok.value;

          if (!strTag) {
            strTag = new Tag();
          }

          if (strTag.type === ValueType.Unknown) {
            strTag.type = ValueType.String;
          }

          switch (strTag.type) {
            case ValueType.String:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid string: "${text}", valid: ""`);
                }
                data = '';
              } else {
                const result = strTag.validateString(text);
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

            case ValueType.DateTime:
            case ValueType.Date:
            case ValueType.Time:
              if (strTag.isNull) {
                data = new Date(0);
              } else {
                try {
                  const dateValue = new Date(text);
                  let result: ValidationResult;
                  if (strTag.type === ValueType.Date) {
                    result = strTag.validateDate(dateValue);
                  } else if (strTag.type === ValueType.Time) {
                    result = strTag.validateTime(dateValue);
                  } else {
                    result = strTag.validateDateTime(dateValue);
                  }
                  if (!result.valid) {
                    throw new Error(
                      result.error || 'DateTime validation failed',
                    );
                  }
                  data = result.data;
                  text = result.text || text;
                } catch (e) {
                  throw new Error(
                    `invalid datetime "${text}": ${(e as Error).message}`,
                  );
                }
              }
              break;

            case ValueType.UUID:
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

            case ValueType.Enum:
              if (!strTag.enum) {
                throw new Error('enum empty');
              }
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid enum: "${text}", valid: ""`);
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

            case ValueType.Image:
              if (strTag.isNull) {
                if (text !== '') {
                  throw new Error(`invalid image: "${text}", valid: ""`);
                }
                data = new Uint8Array();
              } else {
                try {
                  const decoded = base64ToUint8(text);
                  const result = strTag.validateImage(decoded);
                  if (!result.valid) {
                    throw new Error(result.error || 'Image validation failed');
                  }
                  data = result.data;
                  text = result.text || text;
                } catch (e) {
                  throw new Error(
                    `invalid base64 image "${text}": ${(e as Error).message}`,
                  );
                }
              }
              break;

            case ValueType.URL:
            case ValueType.IP:
            case ValueType.Decimal:
            default:
              let defaultResult: ValidationResult;
              if (strTag.type === ValueType.URL) {
                defaultResult = strTag.validateURL(text);
              } else if (strTag.type === ValueType.IP) {
                defaultResult = strTag.validateIP(text);
              } else if (strTag.type === ValueType.Decimal) {
                defaultResult = strTag.validateDecimal(text);
              } else {
                defaultResult = strTag.validateString(text);
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

          const strValue = new MMValue(data, strTag);
          strValue.setPath(path);
          this.depth--;
          return strValue;

        case TokenType.NUMBER:
          let numTag = this.consumeCommentsFor(tok.line);
          text = tok.value;

          if (!numTag) {
            numTag = new Tag();
          }

          if (text.includes('.')) {
            if (numTag.type === ValueType.Unknown) {
              numTag.type = ValueType.Float64;
            }

            switch (numTag.type) {
              case ValueType.Float32:
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
                  const result = numTag.validateFloat32(f64);
                  if (!result.valid) {
                    throw new Error(
                      result.error || 'Float32 validation failed',
                    );
                  }
                  data = result.data;
                  text = result.text || text;
                }
                break;

              case ValueType.Float64:
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
                  const result = numTag.validateFloat64(f64);
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
                const floatDefaultResult = numTag.validateFloat64(
                  parseFloat(text),
                );
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
              numTag.type = ValueType.Int64;
            }

            switch (numTag.type) {
              case ValueType.Int8:
                data = this.parseAndValidateInt(text, numTag, 'validateInt8');
                break;
              case ValueType.Int16:
                data = this.parseAndValidateInt(text, numTag, 'validateInt16');
                break;
              case ValueType.Int32:
                data = this.parseAndValidateInt(text, numTag, 'validateInt32');
                break;
              case ValueType.Int64:
                data = this.parseAndValidateInt(text, numTag, 'validateInt64');
                break;
              case ValueType.BigInt:
                if (numTag.isNull) {
                  if (text !== '0') {
                    throw new Error(`invalid bigint: ${text}, valid: 0`);
                  }
                  data = BigInt(0);
                } else {
                  const bi = BigInt(text);
                  const result = numTag.validateBigInt(bi);
                  if (!result.valid) {
                    throw new Error(result.error || 'BigInt validation failed');
                  }
                  data = result.data;
                }
                break;
              default:
                const num = BigInt(text);
                let negResult = numTag.validateInt(num);
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
              numTag.type = ValueType.Int;
            }

            switch (numTag.type) {
              case ValueType.Int:
                data = this.parseAndValidateInt(text, numTag, 'validateInt');
                break;
              case ValueType.Int8:
                data = this.parseAndValidateInt(text, numTag, 'validateInt8');
                break;
              case ValueType.Int16:
                data = this.parseAndValidateInt(text, numTag, 'validateInt16');
                break;
              case ValueType.Int32:
                data = this.parseAndValidateInt(text, numTag, 'validateInt32');
                break;
              case ValueType.Int64:
                data = this.parseAndValidateInt(text, numTag, 'validateInt64');
                break;
              case ValueType.Uint:
                data = this.parseAndValidateUint(text, numTag, 'validateUint');
                break;
              case ValueType.Uint8:
                data = this.parseAndValidateUint(text, numTag, 'validateUint8');
                break;
              case ValueType.Uint16:
                data = this.parseAndValidateUint(
                  text,
                  numTag,
                  'validateUint16',
                );
                break;
              case ValueType.Uint32:
                data = this.parseAndValidateUint(
                  text,
                  numTag,
                  'validateUint32',
                );
                break;
              case ValueType.Uint64:
                data = this.parseAndValidateUint(
                  text,
                  numTag,
                  'validateUint64',
                );
                break;
              case ValueType.BigInt:
                if (numTag.isNull) {
                  if (text !== '0') {
                    throw new Error(`invalid bigint: ${text}, valid: 0`);
                  }
                  data = BigInt(0);
                } else {
                  const bi = BigInt(text);
                  const result = numTag.validateBigInt(bi);
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

          const numValue = new MMValue(data, numTag);
          numValue.setPath(path);
          this.depth--;
          return numValue;

        case TokenType.TRUE:
          let trueTag = this.consumeCommentsFor(tok.line);
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

          const trueValue = new MMValue(true, trueTag);
          trueValue.setPath(path);
          this.depth--;
          return trueValue;

        case TokenType.FALSE:
          let falseTag = this.consumeCommentsFor(tok.line);
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

          const falseValue = new MMValue(false, falseTag);
          falseValue.setPath(path);
          this.depth--;
          return falseValue;

        case TokenType.NULL:
          let nullTag = this.consumeCommentsFor(tok.line);
          if (!nullTag) {
            nullTag = new Tag();
          }
          nullTag.isNull = true;

          const nullValue = new MMValue(null, nullTag);
          nullValue.setPath(path);
          this.depth--;
          return nullValue;

        default:
          this.depth--;
          throw new Error(`unexpected token ${tok.type}`);
      }
    }
  }

  private parseAndValidateInt(
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

  private parseAndValidateUint(
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

  private parseObject(openLine: number, path: string): MMObject {
    let tag = new Tag();
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Object;
    }

    if (tag.name) {
      path = path ? `${path}.${tag.name}` : tag.name;
    }

    const obj = new MMObject();
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

      if (
        tok.type === TokenType.LINECOMMENT ||
        tok.type === TokenType.BLOCKCOMMENT
      ) {
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
      const keyStr = key.value;

      this.next();
      const pa =
        tag.type === ValueType.Map ? `${path}[${keyStr}]` : `${path}.${keyStr}`;
      val = this.parseValue(pa);
      if (!val) {
        this.next();
        continue;
      }

      const childTag = val.getTag();
      if (childTag && tag) {
        childTag.isInherit = true;
      }

      obj.setProperty(keyStr, val);

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    if (tag.type === ValueType.Map || tag.type === ValueType.Object) {
      const result = tag.validateStruct();
      if (!result.valid) {
        throw new Error(`validate failed: ${result.error}`);
      }
    }

    return obj;
  }

  private parseArray(openLine: number, path: string): MMArray {
    let tag = this.consumeCommentsFor(openLine);
    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = tag.size > 0 ? ValueType.Array : ValueType.Slice;
    }

    if (tag.name) {
      path = `${path}.${tag.name}`;
    }

    const arr = new MMArray();
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

      if (
        tok.type === TokenType.LINECOMMENT ||
        tok.type === TokenType.BLOCKCOMMENT
      ) {
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
      item = this.parseValue(pa);
      if (!item) {
        continue;
      }

      const childTag = item.getTag();
      if (childTag && tag) {
        childTag.isInherit = true;
      }

      arr.addElement(item);
      i++;

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    if (tag.type === ValueType.Array) {
      const result = tag.validateArray([]);
      if (!result.valid) {
        throw new Error(`validate failed: ${result.error}`);
      }
    }

    return arr;
  }

  private mergeNodeTag(n: Node, parsed: Tag): void {
    if (!n || !parsed) {
      return;
    }
    const existing = n.getTag();
    const merged = this.mergeTag(existing, parsed);
    n.setTag(merged);
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
    if (a.default) {
      b.default = a.default;
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
    if (a.enum) {
      b.enum = a.enum;
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
    let trimmed = cs.trim();

    if (trimmed.startsWith('//')) {
      trimmed = trimmed.substring(2).trim();
    } else if (trimmed.startsWith('/*')) {
      trimmed = trimmed.substring(2).trim();
      if (trimmed.endsWith('*/')) {
        trimmed = trimmed.slice(0, -2).trim();
      }
    }

    if (!trimmed.startsWith('mm:')) return null;
    const tagStr = trimmed.substring(3).trim();
    if (!tagStr) return null;
    return parseMMTag(tagStr);
  }
}

export function parseJSONC(input: string): MMDoc {
  const parser = new JSONCParser(input);
  return parser.parse();
}
