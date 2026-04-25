export interface JSONCNode {
  getType(): string;
  getTag(): JSONCTag | null;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: JSONCTag): void;
}

export class JSONCValue implements JSONCNode {
  private value: any;
  private tag: JSONCTag | null;
  private path: string;

  constructor(value: any, tag: JSONCTag | null = null) {
    this.value = value;
    this.tag = tag;
    this.path = '';
  }

  getType(): string {
    if (this.value === null) return 'null';
    if (typeof this.value === 'boolean') return 'boolean';
    if (typeof this.value === 'number') return 'number';
    if (typeof this.value === 'string') return 'string';
    return 'unknown';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getValue(): any {
    return this.value;
  }

  setValue(value: any): void {
    this.value = value;
  }
}

export class JSONCObject implements JSONCNode {
  private properties: Map<string, JSONCNode>;
  private tag: JSONCTag | null;
  private path: string;

  constructor() {
    this.properties = new Map();
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'object';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getProperties(): Map<string, JSONCNode> {
    return this.properties;
  }

  setProperty(key: string, value: JSONCNode): void {
    this.properties.set(key, value);
  }

  getProperty(key: string): JSONCNode | undefined {
    return this.properties.get(key);
  }

  hasProperty(key: string): boolean {
    return this.properties.has(key);
  }

  size(): number {
    return this.properties.size;
  }
}

export class JSONCArray implements JSONCNode {
  private elements: JSONCNode[];
  private tag: JSONCTag | null;
  private path: string;

  constructor() {
    this.elements = [];
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'array';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getElements(): JSONCNode[] {
    return this.elements;
  }

  addElement(element: JSONCNode): void {
    this.elements.push(element);
  }

  getElement(index: number): JSONCNode | undefined {
    return this.elements[index];
  }

  size(): number {
    return this.elements.length;
  }
}

export class JSONCDoc implements JSONCNode {
  private root: JSONCNode;
  private tag: JSONCTag | null;
  private path: string;

  constructor(root: JSONCNode) {
    this.root = root;
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'document';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getRoot(): JSONCNode {
    return this.root;
  }

  setRoot(root: JSONCNode): void {
    this.root = root;
  }
}

export enum JSONCValueType {
  Unknown = 'unknown',
  String = 'string',
  Int = 'int',
  Int8 = 'int8',
  Int16 = 'int16',
  Int32 = 'int32',
  Int64 = 'int64',
  Uint = 'uint',
  Uint8 = 'uint8',
  Uint16 = 'uint16',
  Uint32 = 'uint32',
  Uint64 = 'uint64',
  Float32 = 'float32',
  Float64 = 'float64',
  Bool = 'bool',
  Bytes = 'bytes',
  BigInt = 'bigint',
  DateTime = 'datetime',
  Date = 'date',
  Time = 'time',
  UUID = 'uuid',
  Decimal = 'decimal',
  IP = 'ip',
  URL = 'url',
  Email = 'email',
  Enum = 'enum',
  Array = 'array',
  Struct = 'struct',
  Null = 'null',
}

export class JSONCTag {
  type: JSONCValueType = JSONCValueType.Unknown;
  desc: string = '';
  nullable: boolean = false;
  isNull: boolean = false;
  raw: boolean = false;
  allowEmpty: boolean = false;
  unique: boolean = false;
  defaultValue: string = '';
  min: string = '';
  max: string = '';
  size: number = 0;
  enum: string = '';
  pattern: string = '';
  location: string = '';
  version: number = 0;
  mime: string = '';

  toString(): string {
    const parts: string[] = [];
    if (this.type !== JSONCValueType.Unknown) {
      parts.push(`type=${typeToString(this.type)}`);
    }
    if (this.desc) {
      parts.push(`desc=${this.desc}`);
    }
    if (this.nullable) {
      parts.push('nullable');
    }
    if (this.isNull) {
      parts.push('is_null');
    }
    if (this.raw) {
      parts.push('raw');
    }
    if (this.defaultValue) {
      parts.push(`default=${this.defaultValue}`);
    }
    if (this.enum) {
      parts.push(`enum=${this.enum}`);
    }
    return parts.join('; ');
  }
}

export function typeToString(type: JSONCValueType): string {
  switch (type) {
    case JSONCValueType.String: return 'str';
    case JSONCValueType.Int: return 'i';
    case JSONCValueType.Int8: return 'i8';
    case JSONCValueType.Int16: return 'i16';
    case JSONCValueType.Int32: return 'i32';
    case JSONCValueType.Int64: return 'i64';
    case JSONCValueType.Uint: return 'u';
    case JSONCValueType.Uint8: return 'u8';
    case JSONCValueType.Uint16: return 'u16';
    case JSONCValueType.Uint32: return 'u32';
    case JSONCValueType.Uint64: return 'u64';
    case JSONCValueType.Float32: return 'f32';
    case JSONCValueType.Float64: return 'f64';
    case JSONCValueType.Bool: return 'bool';
    case JSONCValueType.Bytes: return 'bytes';
    case JSONCValueType.BigInt: return 'bi';
    case JSONCValueType.DateTime: return 'datetime';
    case JSONCValueType.Date: return 'date';
    case JSONCValueType.Time: return 'time';
    case JSONCValueType.UUID: return 'uuid';
    case JSONCValueType.Decimal: return 'decimal';
    case JSONCValueType.IP: return 'ip';
    case JSONCValueType.URL: return 'url';
    case JSONCValueType.Email: return 'email';
    case JSONCValueType.Enum: return 'enum';
    case JSONCValueType.Array: return 'arr';
    case JSONCValueType.Struct: return 'struct';
    case JSONCValueType.Null: return 'null';
    default: return 'unknown';
  }
}

export function parseMMTag(tagStr: string): JSONCTag {
  const tag = new JSONCTag();
  const parts = tagStr.split(';');

  for (const part of parts) {
    const trimmed = part.trim();
    if (!trimmed) continue;
    const kv = trimmed.split('=', 2);
    if (kv.length === 0 || !kv[0]) continue;
    const key = kv[0].toLowerCase();
    const value = (kv.length > 1 && kv[1]) ? kv[1].trim() : '';

    switch (key) {
      case 'is_null':
        tag.isNull = true;
        tag.nullable = true;
        break;
      case 'example':
        break;
      case 'desc':
        tag.desc = value;
        break;
      case 'type':
        tag.type = stringToType(value);
        break;
      case 'nullable':
        tag.nullable = true;
        break;
      case 'raw':
        tag.raw = true;
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
        tag.type = JSONCValueType.Enum;
        tag.enum = value;
        break;
      case 'pattern':
        tag.pattern = value;
        break;
      case 'location':
        tag.location = value;
        break;
      case 'version':
        tag.version = parseInt(value, 10) || 0;
        break;
      case 'mime':
        tag.mime = value;
        break;
    }
  }

  return tag;
}

function stringToType(value: string): JSONCValueType {
  switch (value) {
    case 'str': return JSONCValueType.String;
    case 'i': return JSONCValueType.Int;
    case 'i8': return JSONCValueType.Int8;
    case 'i16': return JSONCValueType.Int16;
    case 'i32': return JSONCValueType.Int32;
    case 'i64': return JSONCValueType.Int64;
    case 'u': return JSONCValueType.Uint;
    case 'u8': return JSONCValueType.Uint8;
    case 'u16': return JSONCValueType.Uint16;
    case 'u32': return JSONCValueType.Uint32;
    case 'u64': return JSONCValueType.Uint64;
    case 'f32': return JSONCValueType.Float32;
    case 'f64': return JSONCValueType.Float64;
    case 'bool': return JSONCValueType.Bool;
    case 'bytes': return JSONCValueType.Bytes;
    case 'bi': return JSONCValueType.BigInt;
    case 'datetime': return JSONCValueType.DateTime;
    case 'date': return JSONCValueType.Date;
    case 'time': return JSONCValueType.Time;
    case 'uuid': return JSONCValueType.UUID;
    case 'decimal': return JSONCValueType.Decimal;
    case 'ip': return JSONCValueType.IP;
    case 'url': return JSONCValueType.URL;
    case 'email': return JSONCValueType.Email;
    case 'enum': return JSONCValueType.Enum;
    case 'arr': return JSONCValueType.Array;
    case 'struct': return JSONCValueType.Struct;
    default: return JSONCValueType.Unknown;
  }
}