import { JSONCScanner, TokenType, Token } from './scanner';
import { NodeScalar, NodeObject, NodeArray, MMDoc } from '../ir/ast';
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
                  const dateValue = new Date(text.replace(' ', 'T') + 'Z');
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

          const trueValue = new NodeScalar(true, trueTag);
          trueValue.setPath(path);
          trueValue.setText('true');
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

          const falseValue = new NodeScalar(false, falseTag);
          falseValue.setPath(path);
          falseValue.setText('false');
          this.depth--;
          return falseValue;

        case TokenType.NULL:
          this.depth--;
          throw new Error(`null is not supported`);

        default:
          this.depth--;
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

  private parseObject(openLine: number, path: string): NodeObject {
    let tag = new Tag();
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
      if (childTag && tag && this.hasChildFields(tag)) {
        const origType = childTag.type;
        childTag.inherit(tag);
        if (val instanceof NodeScalar && origType !== childTag.type) {
          this.revalidateValue(val, childTag, origType);
        }
      }

      obj.setProperty(keyStr, val);

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    if (tag.type === ValueType.Map || tag.type === ValueType.Obj) {
      const result = tag.validateObj();
      if (!result.valid) {
        throw new Error(`validate failed: ${result.error}`);
      }
    }

    return obj;
  }

  private parseArray(openLine: number, path: string): NodeArray {
    let tag = this.consumeCommentsFor(openLine);
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
      item = this.parseValue(pa);
      if (!item) {
        continue;
      }

      const childTag = item.getTag();
      if (childTag && tag && this.hasChildFields(tag)) {
        const origType = childTag.type;
        childTag.inherit(tag);
        if (item instanceof NodeScalar && origType !== childTag.type) {
          this.revalidateValue(item, childTag, origType);
        }
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

    if (tag.type === ValueType.Arr) {
      const result = tag.validateArr(arr.getElements());
      if (!result.valid) {
        throw new Error(`validate failed: ${result.error}`);
      }
    }

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

  private hasChildFields(tag: Tag): boolean {
    return (
      tag.childDesc !== '' ||
      tag.childType !== ValueType.Unknown ||
      tag.childNullable ||
      tag.childAllowEmpty ||
      tag.childUnique ||
      tag.childDefaultVal !== '' ||
      tag.childMin !== '' ||
      tag.childMax !== '' ||
      tag.childSize !== 0n ||
      tag.childEnums !== '' ||
      tag.childPattern !== '' ||
      tag.childLocation !== 0 ||
      tag.childVersion !== 0 ||
      tag.childMime !== ''
    );
  }

  private revalidateValue(
    item: NodeScalar,
    tag: Tag,
    origType: ValueType,
  ): void {
    const text = item.getText() || String(item.getValue());

    switch (tag.type) {
      case ValueType.Bytes: {
        try {
          const decoded = base64ToUint8(text);
          const result = tag.validateBytes(decoded);
          if (!result.valid) {
            throw new Error(result.error || 'Bytes validation failed');
          }
          item.setValue(result.data);
          if (result.text) item.setText(result.text);
        } catch (e) {
          throw new Error(
            `invalid base64 bytes "${text}": ${(e as Error).message}`,
          );
        }
        break;
      }

      case ValueType.Datetime:
      case ValueType.Date:
      case ValueType.Time: {
        let dateValue: Date;
        if (tag.type === ValueType.Time) {
          dateValue = new Date(`1970-01-01T${text}Z`);
        } else {
          dateValue = new Date(text.replace(' ', 'T') + 'Z');
        }
        let result: ValidationResult;
        if (tag.type === ValueType.Date) {
          result = tag.validateDate(dateValue);
        } else if (tag.type === ValueType.Time) {
          result = tag.validateTime(dateValue);
        } else {
          result = tag.validateDatetime(dateValue);
        }
        if (!result.valid) {
          throw new Error(result.error || 'Datetime validation failed');
        }
        item.setValue(result.data);
        break;
      }

      case ValueType.Uuid: {
        const result = tag.validateUUID(text);
        if (!result.valid) {
          throw new Error(result.error || 'UUID validation failed');
        }
        item.setValue(result.data);
        if (result.text) item.setText(result.text);
        break;
      }

      case ValueType.Email: {
        const result = tag.validateEmail(text);
        if (!result.valid) {
          throw new Error(result.error || 'Email validation failed');
        }
        item.setValue(result.data);
        if (result.text) item.setText(result.text);
        break;
      }

      case ValueType.Url: {
        const result = tag.validateURL(text);
        if (!result.valid) {
          throw new Error(result.error || 'URL validation failed');
        }
        item.setValue(result.data);
        if (result.text) item.setText(result.text);
        break;
      }

      case ValueType.Ip: {
        const result = tag.validateIP(text);
        if (!result.valid) {
          throw new Error(result.error || 'IP validation failed');
        }
        item.setValue(result.data);
        if (result.text) item.setText(result.text);
        break;
      }

      case ValueType.Decimal: {
        const result = tag.validateDecimal(text);
        if (!result.valid) {
          throw new Error(result.error || 'Decimal validation failed');
        }
        item.setValue(result.data);
        if (result.text) item.setText(result.text);
        break;
      }

      case ValueType.Bigint: {
        try {
          const bi = BigInt(text);
          const result = tag.validateBigint(bi);
          if (!result.valid) {
            throw new Error(result.error || 'BigInt validation failed');
          }
          item.setValue(result.data);
          if (result.text) item.setText(result.text);
        } catch (e) {
          throw new Error(`invalid bigint "${text}": ${(e as Error).message}`);
        }
        break;
      }

      case ValueType.Enums: {
        if (tag.enums) {
          const result = tag.validateEnum(text);
          if (!result.valid) {
            throw new Error(result.error || 'Enum validation failed');
          }
          item.setValue(result.data);
          if (result.text) item.setText(result.text);
        }
        break;
      }

      default:
        break;
    }
  }
}

export function parseJSONC(input: string): MMDoc {
  const parser = new JSONCParser(input);
  return parser.parse();
}
