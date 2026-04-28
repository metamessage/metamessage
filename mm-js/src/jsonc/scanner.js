import { MmValidator } from '../mm/validator.js';

export const TokenType = {
  EOF: 'EOF',
  LCURLY: 'LCURLY',
  RCURLY: 'RCURLY',
  LBRACKET: 'LBRACKET',
  RBRACKET: 'RBRACKET',
  COLON: 'COLON',
  COMMA: 'COMMA',
  STRING: 'STRING',
  NUMBER: 'NUMBER',
  TRUE: 'TRUE',
  FALSE: 'FALSE',
  NULL: 'NULL',
  LEADING_COMMENT: 'LEADING_COMMENT',
  TRAILING_COMMENT: 'TRAILING_COMMENT'
};

export class Token {
  constructor(type, literal, line, column) {
    this.type = type;
    this.literal = literal;
    this.line = line;
    this.column = column;
  }

  toString() {
    return `Token(${this.type}, "${this.literal}", line:${this.line}, col:${this.column})`;
  }
}

export class JSONCScanner {
  constructor(input) {
    this.src = Array.from(input);
    this.pos = 0;
    this.line = 1;
    this.col = 1;
    this.newLine = true;
  }

  peek() {
    if (this.pos >= this.src.length) {
      return null;
    }
    return this.src[this.pos];
  }

  next() {
    if (this.pos >= this.src.length) {
      return null;
    }
    const ch = this.src[this.pos];
    this.pos++;

    if (ch === '\n') {
      this.newLine = true;
      this.line++;
      this.col = 1;
    } else {
      this.col++;
    }
    return ch;
  }

  skipWhitespace() {
    while (true) {
      const ch = this.peek();
      if (ch === null) break;
      if (ch === ' ' || ch === '\t' || ch === '\r' || ch === '\n') {
        this.next();
      } else {
        break;
      }
    }
  }

  nextToken() {
    this.skipWhitespace();

    const ch = this.peek();
    if (ch === null) {
      return new Token(TokenType.EOF, '', this.line, this.col);
    }

    const startLine = this.line;
    const startCol = this.col;

    switch (ch) {
      case '{':
        this.next();
        return new Token(TokenType.LCURLY, '{', startLine, startCol);
      case '}':
        this.next();
        return new Token(TokenType.RCURLY, '}', startLine, startCol);
      case '[':
        this.next();
        return new Token(TokenType.LBRACKET, '[', startLine, startCol);
      case ']':
        this.next();
        return new Token(TokenType.RBRACKET, ']', startLine, startCol);
      case ':':
        this.next();
        this.newLine = false;
        return new Token(TokenType.COLON, ':', startLine, startCol);
      case ',':
        this.next();
        this.newLine = false;
        return new Token(TokenType.COMMA, ',', startLine, startCol);
      case '"':
        return this.scanString(startLine, startCol);
      case '/':
        return this.scanComment(startLine, startCol);
      default:
        return this.scanLiteral(startLine, startCol);
    }
  }

  scanString(startLine, startCol) {
    this.next();

    let buf = '';
    while (true) {
      const ch = this.peek();
      if (ch === null || ch === '\n') {
        break;
      }
      const c = this.next();
      if (c === '"') {
        break;
      }
      if (c === '\\') {
        buf += c;
        const escaped = this.peek();
        if (escaped !== null) {
          buf += this.next();
        }
        continue;
      }
      buf += c;
    }

    return new Token(TokenType.STRING, buf, startLine, startCol);
  }

  scanComment(startLine, startCol) {
    this.next();

    if (this.peek() === '/') {
      const commentType = this.newLine ? TokenType.LEADING_COMMENT : TokenType.TRAILING_COMMENT;
      this.next();

      let buf = '';
      while (true) {
        const ch = this.peek();
        if (ch === null || ch === '\n') {
          break;
        }
        buf += this.next();
      }
      const trimmed = buf.trim();

      return new Token(commentType, trimmed, startLine, startCol);
    }

    if (this.peek() === '*') {
      const commentType = this.newLine ? TokenType.LEADING_COMMENT : TokenType.TRAILING_COMMENT;
      this.next();

      let buf = '';
      while (true) {
        const ch = this.peek();
        if (ch === null) {
          break;
        }
        if (ch === '*' && this.pos + 1 < this.src.length && this.src[this.pos + 1] === '/') {
          this.next();
          this.next();
          break;
        }
        buf += this.next();
      }
      const trimmed = buf.trim();

      return new Token(commentType, trimmed, startLine, startCol);
    }

    return new Token(TokenType.EOF, '', startLine, startCol);
  }

  scanLiteral(startLine, startCol) {
    let buf = '';

    while (true) {
      const ch = this.peek();
      if (ch === null) break;
      if (ch === ' ' || ch === '\t' || ch === '\r' || ch === '\n' || ch === '{' || ch === '}' || ch === '[' || ch === ']' || ch === ':' || ch === ',') {
        break;
      }
      buf += this.next();
    }

    switch (buf) {
      case 'true':
        return new Token(TokenType.TRUE, buf, startLine, startCol);
      case 'false':
        return new Token(TokenType.FALSE, buf, startLine, startCol);
      case 'null':
        return new Token(TokenType.NULL, buf, startLine, startCol);
      default:
        return new Token(TokenType.NUMBER, buf, startLine, startCol);
    }
  }
}

export function parseJSONC(input) {
  const scanner = new JSONCScanner(input);
  const tokens = [];

  while (true) {
    const token = scanner.nextToken();
    tokens.push(token);
    if (token.type === TokenType.EOF) {
      break;
    }
  }

  const parser = new JSONCParser(tokens);
  return parser.parse();
}

class JSONCParser {
  constructor(tokens) {
    this.tokens = tokens;
    this.pos = 0;
    this.pendingComments = [];
    this.depth = 0;
    this.maxDepth = 32;
  }

  peek() {
    if (this.pos >= this.tokens.length) {
      return new Token(TokenType.EOF, '', 0, 0);
    }
    return this.tokens[this.pos];
  }

  next() {
    const token = this.peek();
    this.pos++;
    return token;
  }

  parse() {
    let result = null;

    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        return result;
      }

      if (tok.type === TokenType.LEADING_COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      if (tok.type === TokenType.TRAILING_COMMENT) {
        if (result) {
          const parsed = parseCommentToTag(tok.literal);
          if (parsed) {
            mergeNodeTag(result, parsed);
          }
        }
        this.next();
        continue;
      }

      result = this.parseNode('');
    }
  }

  parseNode(path) {
    const tok = this.next();

    switch (tok.type) {
      case TokenType.EOF:
        return null;

      case TokenType.LCURLY:
        return this.parseObject(tok.line, path);

      case TokenType.LBRACKET:
        return this.parseArray(tok.line, path);

      case TokenType.STRING: {
        let tag = this.consumeCommentsFor(tok.line) || new JSONCTag();
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.String;
        }
        const text = tok.literal;
        
        // 验证值
        const convertedTag = convertJSONCTagToTag(tag);
        const result = MmValidator.validate(text, convertedTag);
        if (!result.valid) {
          throw new Error(result.error || 'String validation failed');
        }
        
        return new JSONCValue(text, text, tag, path);
      }

      case TokenType.NUMBER: {
        let tag = this.consumeCommentsFor(tok.line) || new JSONCTag();
        if (tag.type === ValueType.Unknown) {
          if (tok.literal.includes('.')) {
            tag.type = ValueType.Float64;
          } else {
            tag.type = ValueType.Int;
          }
        }
        let data;
        if (tok.literal.includes('.')) {
          data = parseFloat(tok.literal);
        } else if (tok.literal.startsWith('-')) {
          data = BigInt(tok.literal);
        } else {
          data = BigInt(tok.literal);
        }
        
        // 验证值
        const convertedTag = convertJSONCTagToTag(tag);
        // 将 BigInt 转换为数字进行验证
        const validationData = typeof data === 'bigint' ? Number(data) : data;
        const result = MmValidator.validate(validationData, convertedTag);
        if (!result.valid) {
          throw new Error(result.error || 'Number validation failed');
        }
        
        return new JSONCValue(data, tok.literal, tag, path);
      }

      case TokenType.TRUE: {
        let tag = this.consumeCommentsFor(tok.line) || new JSONCTag();
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Bool;
        }
        
        // 验证值
        const convertedTag = convertJSONCTagToTag(tag);
        const result = MmValidator.validate(true, convertedTag);
        if (!result.valid) {
          throw new Error(result.error || 'Boolean validation failed');
        }
        
        return new JSONCValue(true, 'true', tag, path);
      }

      case TokenType.FALSE: {
        let tag = this.consumeCommentsFor(tok.line) || new JSONCTag();
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Bool;
        }
        
        // 验证值
        const convertedTag = convertJSONCTagToTag(tag);
        const result = MmValidator.validate(false, convertedTag);
        if (!result.valid) {
          throw new Error(result.error || 'Boolean validation failed');
        }
        
        return new JSONCValue(false, 'false', tag, path);
      }

      case TokenType.NULL: {
        let tag = this.consumeCommentsFor(tok.line) || new JSONCTag();
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Unknown;
        }
        tag.isNull = true;
        return new JSONCValue(null, 'null', tag, path);
      }

      default:
        throw new Error(`Unexpected token: ${tok.type}`);
    }
  }

  parseObject(openLine, path) {
    this.depth++;
    if (this.depth > this.maxDepth) {
      throw new Error('Max depth exceeded');
    }

    let tag = this.consumeCommentsFor(openLine) || new JSONCTag();
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Struct;
    }

    // 验证结构体 tag
    const convertedTag = convertJSONCTagToTag(tag);
    const result = MmValidator.validateStruct(convertedTag);
    if (!result.valid) {
      throw new Error(result.error || 'Struct validation failed');
    }

    const obj = new JSONCObject(tag, path);

    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        break;
      }
      if (tok.type === TokenType.RCURLY) {
        this.next();
        break;
      }

      if (tok.type === TokenType.LEADING_COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      if (tok.type === TokenType.TRAILING_COMMENT) {
        if (obj.fields.length > 0) {
          const lastField = obj.fields[obj.fields.length - 1];
          const parsed = parseCommentToTag(tok.literal);
          if (parsed) {
            mergeNodeTag(lastField.value, parsed);
          }
        }
        this.next();
        continue;
      }

      const keyTok = this.next();
      if (keyTok.type !== TokenType.STRING) {
        throw new Error('Expected string key');
      }

      const key = keyTok.literal;
      this.next();

      const childPath = `${path}.${key}`;
      const val = this.parseNode(childPath);
      if (val) {
        const childTag = val.tag;
        if (childTag && tag) {
          childTag.inherit(tag);
        }
        obj.fields.push(new JSONCField(key, val));
      }

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    this.depth--;
    return obj;
  }

  parseArray(openLine, path) {
    this.depth++;
    if (this.depth > this.maxDepth) {
      throw new Error('Max depth exceeded');
    }

    let tag = this.consumeCommentsFor(openLine) || new JSONCTag();
    if (tag.type === ValueType.Unknown) {
      if (tag.size > 0) {
        tag.type = ValueType.Array;
      } else {
        tag.type = ValueType.Slice;
      }
    }

    // 验证数组 tag
    const convertedTag = convertJSONCTagToTag(tag);
    const result = MmValidator.validateStruct(convertedTag);
    if (!result.valid) {
      throw new Error(result.error || 'Array validation failed');
    }

    const arr = new JSONCArray(tag, path);
    let index = 0;

    while (true) {
      const tok = this.peek();
      if (tok.type === TokenType.EOF) {
        break;
      }
      if (tok.type === TokenType.RBRACKET) {
        this.next();
        break;
      }

      if (tok.type === TokenType.LEADING_COMMENT) {
        if (this.pendingComments.length > 0) {
          const last = this.pendingComments[this.pendingComments.length - 1];
          if (tok.line - last.line > 1) {
            this.pendingComments = [];
          }
        }
        this.pendingComments.push(tok);
        this.next();
        continue;
      }

      if (tok.type === TokenType.TRAILING_COMMENT) {
        if (arr.items.length > 0) {
          const lastItem = arr.items[arr.items.length - 1];
          const parsed = parseCommentToTag(tok.literal);
          if (parsed) {
            mergeNodeTag(lastItem, parsed);
          }
        }
        this.next();
        continue;
      }

      const itemPath = `${path}[${index}]`;
      const item = this.parseNode(itemPath);
      if (item) {
        const childTag = item.tag;
        if (childTag && tag) {
          childTag.inherit(tag);
        }
        arr.items.push(item);
        index++;
      }

      if (this.peek().type === TokenType.COMMA) {
        this.next();
      }
    }

    this.depth--;
    return arr;
  }

  consumeCommentsFor(anchorLine) {
    if (this.pendingComments.length === 0) {
      return null;
    }

    const last = this.pendingComments[this.pendingComments.length - 1];
    if (anchorLine - last.line > 1) {
      this.pendingComments = [];
      return null;
    }

    let result = null;
    for (const comment of this.pendingComments) {
      const parsed = parseCommentToTag(comment.literal);
      if (parsed) {
        result = mergeTag(result, parsed);
      }
    }

    this.pendingComments = [];
    return result;
  }
}

class JSONCTag {
  constructor() {
    this.name = '';
    this.isNull = false;
    this.example = false;
    this.desc = '';
    this.type = ValueType.Unknown;
    this.raw = false;
    this.nullable = false;
    this.allowEmpty = false;
    this.unique = false;
    this.defaultValue = '';
    this.min = '';
    this.max = '';
    this.size = 0;
    this.enumValues = '';
    this.pattern = '';
    this.locationOffset = 0;
    this.version = 0;
    this.mime = '';
    this.childDesc = '';
    this.childType = ValueType.Unknown;
    this.childRaw = false;
    this.childNullable = false;
    this.childAllowEmpty = false;
    this.childUnique = false;
    this.childDefault = '';
    this.childMin = '';
    this.childMax = '';
    this.childSize = 0;
    this.childEnum = '';
    this.childPattern = '';
    this.childLocationOffset = 0;
    this.childVersion = 0;
    this.childMime = '';
    this.isInherit = false;
  }

  inherit(parent) {
    this.desc = parent.childDesc;
    this.type = parent.childType;
    this.raw = parent.childRaw;
    this.nullable = parent.childNullable;
    this.allowEmpty = parent.childAllowEmpty;
    this.unique = parent.childUnique;
    this.defaultValue = parent.childDefault;
    this.min = parent.childMin;
    this.max = parent.childMax;
    this.size = parent.childSize;
    this.enumValues = parent.childEnum;
    this.pattern = parent.childPattern;
    this.locationOffset = parent.childLocationOffset;
    this.version = parent.childVersion;
    this.mime = parent.childMime;
  }
}

function mergeTag(dst, src) {
  if (!dst) return src;

  const merged = new JSONCTag();
  merged.name = src.name || dst.name;
  if (src.isNull) merged.isNull = src.isNull;
  if (src.example) merged.example = src.example;
  if (src.desc) merged.desc = src.desc;
  if (src.type !== ValueType.Unknown) merged.type = src.type;
  if (src.raw) merged.raw = src.raw;
  if (src.nullable) merged.nullable = src.nullable;
  if (src.allowEmpty) merged.allowEmpty = src.allowEmpty;
  if (src.unique) merged.unique = src.unique;
  if (src.defaultValue) merged.defaultValue = src.defaultValue;
  if (src.min) merged.min = src.min;
  if (src.max) merged.max = src.max;
  if (src.size) merged.size = src.size;
  if (src.enumValues) merged.enumValues = src.enumValues;
  if (src.pattern) merged.pattern = src.pattern;
  if (src.locationOffset) merged.locationOffset = src.locationOffset;
  if (src.version) merged.version = src.version;
  if (src.mime) merged.mime = src.mime;

  if (src.childDesc) merged.childDesc = src.childDesc;
  if (src.childType !== ValueType.Unknown) merged.childType = src.childType;
  if (src.childRaw) merged.childRaw = src.childRaw;
  if (src.childNullable) merged.childNullable = src.childNullable;
  if (src.childAllowEmpty) merged.childAllowEmpty = src.childAllowEmpty;
  if (src.childUnique) merged.childUnique = src.childUnique;
  if (src.childDefault) merged.childDefault = src.childDefault;
  if (src.childMin) merged.childMin = src.childMin;
  if (src.childMax) merged.childMax = src.childMax;
  if (src.childSize) merged.childSize = src.childSize;
  if (src.childEnum) merged.childEnum = src.childEnum;
  if (src.childPattern) merged.childPattern = src.childPattern;
  if (src.childLocationOffset) merged.childLocationOffset = src.childLocationOffset;
  if (src.childVersion) merged.childVersion = src.childVersion;
  if (src.childMime) merged.childMime = src.childMime;

  return merged;
}

function parseCommentToTag(literal) {
  if (literal.startsWith('mm:')) {
    return parseMMTag(literal.slice(3));
  }
  return null;
}

function parseMMTag(tagStr) {
  const tag = new JSONCTag();
  const trimmed = tagStr.trim();

  if (!trimmed) {
    return tag;
  }

  const parts = trimmed.split(';').map(p => p.trim());

  for (const part of parts) {
    if (!part) continue;

    let key, value;
    const equalIndex = part.indexOf('=');
    if (equalIndex >= 0) {
      key = part.slice(0, equalIndex).trim().toLowerCase();
      value = part.slice(equalIndex + 1).trim();
    } else {
      key = part.trim().toLowerCase();
      value = '';
    }

    switch (key) {
      case 'is_null':
        tag.isNull = true;
        tag.nullable = true;
        break;
      case 'example':
        tag.example = true;
        break;
      case 'desc':
        tag.desc = value.replace(/^["']|["']$/g, '');
        break;
      case 'type':
        tag.type = parseValueType(value);
        break;
      case 'raw':
        tag.raw = true;
        break;
      case 'nullable':
        tag.nullable = true;
        break;
      case 'allow_empty':
        tag.allowEmpty = true;
        break;
      case 'unique':
        tag.unique = true;
        break;
      case 'default':
        tag.defaultValue = value;
        break;
      case 'min':
        tag.min = value;
        break;
      case 'max':
        tag.max = value;
        break;
      case 'size':
        tag.size = parseInt(value, 10) || 0;
        break;
      case 'enum':
        tag.type = ValueType.Enum;
        tag.enumValues = value;
        break;
      case 'pattern':
        tag.pattern = value;
        break;
      case 'location':
        tag.locationOffset = parseInt(value, 10) || 0;
        break;
      case 'version':
        tag.version = parseInt(value, 10) || 0;
        break;
      case 'mime':
        tag.mime = value;
        break;
      case 'child_desc':
        tag.childDesc = value.replace(/^["']|["']$/g, '');
        break;
      case 'child_type':
        tag.childType = parseValueType(value);
        break;
      case 'child_raw':
        tag.childRaw = true;
        break;
      case 'child_nullable':
        tag.childNullable = true;
        break;
      case 'child_allow_empty':
        tag.childAllowEmpty = true;
        break;
      case 'child_unique':
        tag.childUnique = true;
        break;
      case 'child_default':
        tag.childDefault = value;
        break;
      case 'child_min':
        tag.childMin = value;
        break;
      case 'child_max':
        tag.childMax = value;
        break;
      case 'child_size':
        tag.childSize = parseInt(value, 10) || 0;
        break;
      case 'child_enum':
        tag.childEnum = value;
        tag.childType = ValueType.Enum;
        break;
      case 'child_location':
        tag.childLocationOffset = parseInt(value, 10) || 0;
        break;
      case 'child_version':
        tag.childVersion = parseInt(value, 10) || 0;
        break;
      case 'child_mime':
        tag.childMime = value;
        break;
    }
  }

  return tag;
}

function parseValueType(s) {
  switch (s.toLowerCase()) {
    case 'unknown': return ValueType.Unknown;
    case 'doc': return ValueType.Doc;
    case 'slice': return ValueType.Slice;
    case 'arr':
    case 'array': return ValueType.Array;
    case 'struct': return ValueType.Struct;
    case 'map': return ValueType.Map;
    case 'str':
    case 'string': return ValueType.String;
    case 'bytes': return ValueType.Bytes;
    case 'bool': return ValueType.Bool;
    case 'i':
    case 'int': return ValueType.Int;
    case 'i8': return ValueType.Int8;
    case 'i16': return ValueType.Int16;
    case 'i32': return ValueType.Int32;
    case 'i64': return ValueType.Int64;
    case 'u':
    case 'uint': return ValueType.Uint;
    case 'u8': return ValueType.Uint8;
    case 'u16': return ValueType.Uint16;
    case 'u32': return ValueType.Uint32;
    case 'u64': return ValueType.Uint64;
    case 'f32': return ValueType.Float32;
    case 'f64':
    case 'float':
    case 'double': return ValueType.Float64;
    case 'bi':
    case 'bigint': return ValueType.BigInt;
    case 'datetime': return ValueType.DateTime;
    case 'date': return ValueType.Date;
    case 'time': return ValueType.Time;
    case 'uuid': return ValueType.Uuid;
    case 'decimal': return ValueType.Decimal;
    case 'ip': return ValueType.Ip;
    case 'url': return ValueType.Url;
    case 'email': return ValueType.Email;
    case 'enum': return ValueType.Enum;
    case 'image': return ValueType.Image;
    case 'video': return ValueType.Video;
    default: return ValueType.Unknown;
  }
}

class JSONCNode {
  getTag() { return null; }
  getType() { return 'unknown'; }
  getPath() { return ''; }
}

class JSONCField {
  constructor(key, value) {
    this.key = key;
    this.value = value;
  }
}

class JSONCObject extends JSONCNode {
  constructor(tag, path) {
    super();
    this.fields = [];
    this.tag = tag;
    this.path = path;
  }

  getTag() { return this.tag; }
  getType() { return 'object'; }
  getPath() { return this.path; }
}

class JSONCArray extends JSONCNode {
  constructor(tag, path) {
    super();
    this.items = [];
    this.tag = tag;
    this.path = path;
  }

  getTag() { return this.tag; }
  getType() { return 'array'; }
  getPath() { return this.path; }
}

class JSONCValue extends JSONCNode {
  constructor(data, text, tag, path) {
    super();
    this.data = data;
    this.text = text;
    this.tag = tag;
    this.path = path;
  }

  getTag() { return this.tag; }
  getType() { return 'value'; }
  getPath() { return this.path; }
}

function mergeNodeTag(node, tag) {
  if (!node || !tag) return;

  const existing = node.getTag();
  if (existing) {
    const merged = mergeTag(existing, tag);
    if (node instanceof JSONCValue) {
      node.tag = merged;
    } else if (node instanceof JSONCObject) {
      node.tag = merged;
    } else if (node instanceof JSONCArray) {
      node.tag = merged;
    }
  }
}

// 转换 JSONCTag 为 Tag 类型
function convertJSONCTagToTag(jsoncTag) {
  return {
    name: jsoncTag.name,
    isNull: jsoncTag.isNull,
    example: jsoncTag.example,
    desc: jsoncTag.desc,
    type: jsoncTag.type,
    raw: jsoncTag.raw,
    nullable: jsoncTag.nullable,
    allowEmpty: jsoncTag.allowEmpty,
    unique: jsoncTag.unique,
    default: jsoncTag.defaultValue,
    min: jsoncTag.min,
    max: jsoncTag.max,
    size: jsoncTag.size,
    enum: jsoncTag.enumValues,
    pattern: jsoncTag.pattern,
    location: jsoncTag.locationOffset,
    version: jsoncTag.version,
    mime: jsoncTag.mime,
    childDesc: jsoncTag.childDesc,
    childType: jsoncTag.childType,
    childRaw: jsoncTag.childRaw,
    childNullable: jsoncTag.childNullable,
    childAllowEmpty: jsoncTag.childAllowEmpty,
    childUnique: jsoncTag.childUnique,
    childDefault: jsoncTag.childDefault,
    childMin: jsoncTag.childMin,
    childMax: jsoncTag.childMax,
    childSize: jsoncTag.childSize,
    childEnum: jsoncTag.childEnum,
    childPattern: jsoncTag.childPattern,
    childLocation: jsoncTag.childLocationOffset,
    childVersion: jsoncTag.childVersion,
    childMime: jsoncTag.childMime,
    isInherit: jsoncTag.isInherit
  };
}

export { JSONCParser, JSONCTag, JSONCNode, JSONCField, JSONCObject, JSONCArray, JSONCValue };

const ValueType = {
  Unknown: 'unknown',
  Doc: 'doc',
  Slice: 'slice',
  Array: 'array',
  Struct: 'struct',
  Map: 'map',
  String: 'str',
  Bytes: 'bytes',
  Bool: 'bool',
  Int: 'int',
  Int8: 'i8',
  Int16: 'i16',
  Int32: 'i32',
  Int64: 'i64',
  Uint: 'uint',
  Uint8: 'u8',
  Uint16: 'u16',
  Uint32: 'u32',
  Uint64: 'u64',
  Float32: 'f32',
  Float64: 'f64',
  BigInt: 'bi',
  DateTime: 'datetime',
  Date: 'date',
  Time: 'time',
  Uuid: 'uuid',
  Decimal: 'decimal',
  Ip: 'ip',
  Url: 'url',
  Email: 'email',
  Enum: 'enum',
  Image: 'image',
  Video: 'video'
};