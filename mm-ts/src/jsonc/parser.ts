import { JSONCScanner, TokenType, Token } from './scanner';
import { JSONCValue, JSONCObject, JSONCArray, JSONCDoc, JSONCTag, parseMMTag, JSONCValueType } from './ast';
import { MmValidator } from '../mm/validator';
import { ValueType, Tag } from '../mm/types';

// 转换 JSONCTag 为 Tag 类型
function convertJSONCTagToTag(jsoncTag: JSONCTag): Tag {
  const tag: Tag = {
    isNull: jsoncTag.isNull,
    desc: jsoncTag.desc,
    type: jsoncTag.type as unknown as ValueType,
    raw: jsoncTag.raw,
    nullable: jsoncTag.nullable,
    allowEmpty: jsoncTag.allowEmpty,
    unique: jsoncTag.unique,
    default: jsoncTag.defaultValue,
    min: jsoncTag.min,
    max: jsoncTag.max,
    size: jsoncTag.size,
    enum: jsoncTag.enum,
    pattern: jsoncTag.pattern,
    version: jsoncTag.version,
    mime: jsoncTag.mime
  };
  
  if (jsoncTag.location) {
    tag.location = parseInt(jsoncTag.location, 10);
  }
  
  return tag;
}

export interface JSONCNode {
  getType(): string;
  getTag(): any;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: any): void;
}

export class JSONCParser {
  private scanner: JSONCScanner;
  private currentToken: Token;
  private pendingComments: Token[] = [];

  constructor(input: string) {
    this.scanner = new JSONCScanner(input);
    this.currentToken = this.scanner.nextToken();
  }

  parse(): JSONCDoc {
    const root = this.parseValue();
    this.expect(TokenType.EOF);
    return new JSONCDoc(root);
  }

  private parseValue(): JSONCNode {
    this.skipLeadingComments();

    switch (this.currentToken.type) {
      case TokenType.LCURLY:
        return this.parseObject();
      case TokenType.LBRACKET:
        return this.parseArray();
      case TokenType.STRING:
        return this.parseString();
      case TokenType.NUMBER:
        return this.parseNumber();
      case TokenType.TRUE:
        return this.parseBoolean(true);
      case TokenType.FALSE:
        return this.parseBoolean(false);
      case TokenType.NULL:
        return this.parseNull();
      default:
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
    }
  }

  private parseObject(): JSONCObject {
    this.expect(TokenType.LCURLY);
    const tag = this.consumeComments();
    const obj = new JSONCObject();
    if (tag) {
      obj.setTag(tag);
      // 验证结构体 tag
      const convertedTag = convertJSONCTagToTag(tag);
      const result = MmValidator.validateStruct(convertedTag);
      if (!result.valid) {
        throw new Error(result.error || 'Struct validation failed');
      }
    }

    this.skipLeadingComments();

    while (this.currentToken.type !== (TokenType.RCURLY as any)) {
      this.skipLeadingComments();

      if (this.currentToken.type === TokenType.LINECOMMENT || this.currentToken.type === TokenType.BLOCKCOMMENT) {
        this.pendingComments.push(this.currentToken);
        this.consumeToken();
        continue;
      }

      const keyToken = this.currentToken;
      this.expect(TokenType.STRING);
      const key = keyToken.value;

      this.skipLeadingComments();
      this.expect(TokenType.COLON);

      this.skipLeadingComments();
      const value = this.parseValue();

      this.skipTrailingComments();

      obj.setProperty(key, value);

      if (this.currentToken.type === (TokenType.COMMA as any)) {
        this.expect(TokenType.COMMA);
        this.skipLeadingComments();
      } else if (this.currentToken.type !== (TokenType.RCURLY as any)) {
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
      }
    }

    this.expect(TokenType.RCURLY);
    return obj;
  }

  private parseArray(): JSONCArray {
    this.expect(TokenType.LBRACKET);
    const tag = this.consumeComments();
    const arr = new JSONCArray();
    if (tag) {
      arr.setTag(tag);
      // 验证数组 tag
      const convertedTag = convertJSONCTagToTag(tag);
      const result = MmValidator.validateStruct(convertedTag);
      if (!result.valid) {
        throw new Error(result.error || 'Array validation failed');
      }
    }

    this.skipLeadingComments();

    let index = 0;
    while (this.currentToken.type !== (TokenType.RBRACKET as any)) {
      this.skipLeadingComments();

      if (this.currentToken.type === TokenType.LINECOMMENT || this.currentToken.type === TokenType.BLOCKCOMMENT) {
        this.pendingComments.push(this.currentToken);
        this.consumeToken();
        continue;
      }

      const value = this.parseValue();
      arr.addElement(value);

      this.skipTrailingComments();

      if (this.currentToken.type === (TokenType.COMMA as any)) {
        this.expect(TokenType.COMMA);
        this.skipLeadingComments();
      } else if (this.currentToken.type !== (TokenType.RBRACKET as any)) {
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
      }
      index++;
    }

    this.expect(TokenType.RBRACKET);
    return arr;
  }

  private parseString(): JSONCValue {
    const token = this.currentToken;
    this.expect(TokenType.STRING);
    const tag = this.consumeComments() || new JSONCTag();
    if (tag.type === JSONCValueType.Unknown) {
      tag.type = JSONCValueType.String;
    }
    
    // 验证值
    const convertedTag = convertJSONCTagToTag(tag);
    let value: any = token.value;
    
    // 处理日期时间类型
    if (tag.type === JSONCValueType.DateTime || tag.type === JSONCValueType.Date || tag.type === JSONCValueType.Time) {
      value = new Date(token.value);
    }
    
    const result = MmValidator.validate(value, convertedTag);
    if (!result.valid) {
      throw new Error(result.error || 'String validation failed');
    }
    
    return new JSONCValue(token.value, tag);
  }

  private parseNumber(): JSONCValue {
    const token = this.currentToken;
    this.expect(TokenType.NUMBER);
    const tag = this.consumeComments() || new JSONCTag();
    if (tag.type === JSONCValueType.Unknown) {
      tag.type = token.value.includes('.') ? JSONCValueType.Float64 : JSONCValueType.Int;
    }
    const value = parseFloat(token.value);
    
    // 验证值
    const convertedTag = convertJSONCTagToTag(tag);
    const result = MmValidator.validate(value, convertedTag);
    if (!result.valid) {
      throw new Error(result.error || 'Number validation failed');
    }
    
    return new JSONCValue(value, tag);
  }

  private parseBoolean(value: boolean): JSONCValue {
    this.expect(value ? TokenType.TRUE : TokenType.FALSE);
    const tag = this.consumeComments() || new JSONCTag();
    tag.type = JSONCValueType.Bool;
    
    // 验证值
    const convertedTag = convertJSONCTagToTag(tag);
    const result = MmValidator.validate(value, convertedTag);
    if (!result.valid) {
      throw new Error(result.error || 'Boolean validation failed');
    }
    
    return new JSONCValue(value, tag);
  }

  private parseNull(): JSONCValue {
    this.expect(TokenType.NULL);
    const tag = this.consumeComments() || new JSONCTag();
    tag.isNull = true;
    return new JSONCValue(null, tag);
  }

  private skipLeadingComments(): void {
    while (this.currentToken.type === TokenType.LINECOMMENT || this.currentToken.type === TokenType.BLOCKCOMMENT) {
      this.pendingComments.push(this.currentToken);
      this.consumeToken();
    }
  }

  private skipTrailingComments(): void {
    while (this.currentToken.type === TokenType.LINECOMMENT || this.currentToken.type === TokenType.BLOCKCOMMENT) {
      this.pendingComments.push(this.currentToken);
      this.consumeToken();
    }
  }

  private consumeComments(): JSONCTag | null {
    if (this.pendingComments.length === 0) return null;

    const comments = this.pendingComments.slice();
    this.pendingComments = [];

    let merged = new JSONCTag();
    for (const c of comments) {
      const t = tagFromComment(c.value);
      if (t) {
        merged = mergeTag(merged, t);
      }
    }
    return merged;
  }

  private expect(type: TokenType): void {
    if (this.currentToken.type !== type) {
      throw new Error(`Expected ${type}, got ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
    }
    this.consumeToken();
  }

  private consumeToken(): void {
    this.currentToken = this.scanner.nextToken();
  }
}

function tagFromComment(comment: string): JSONCTag | null {
  let trimmed = comment.trim();

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

function mergeTag(a: JSONCTag, b: JSONCTag): JSONCTag {
  if (a.type !== JSONCValueType.Unknown) b.type = a.type;
  if (a.desc) b.desc = a.desc;
  if (a.nullable) b.nullable = true;
  if (a.isNull) b.isNull = true;
  if (a.defaultValue) b.defaultValue = a.defaultValue;
  if (a.min) b.min = a.min;
  if (a.max) b.max = a.max;
  if (a.size !== 0) b.size = a.size;
  if (a.enum) b.enum = a.enum;
  if (a.pattern) b.pattern = a.pattern;
  if (a.location) b.location = a.location;
  if (a.version !== 0) b.version = a.version;
  if (a.mime) b.mime = a.mime;
  return b;
}

export function parseJSONC(input: string): JSONCDoc {
  const parser = new JSONCParser(input);
  return parser.parse();
}