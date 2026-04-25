import { TokenType, JSONCScanner } from './scanner.js';

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

export class JSONCPrinter {
  constructor(indentString = '  ', useIndent = true) {
    this.indentLevel = 0;
    this.indentString = indentString;
    this.useIndent = useIndent;
  }

  indent() {
    if (!this.useIndent) {
      return '';
    }
    return this.indentString.repeat(this.indentLevel);
  }

  increaseIndent() {
    this.indentLevel++;
  }

  decreaseIndent() {
    this.indentLevel = Math.max(0, this.indentLevel - 1);
  }

  print(node) {
    if (!node) return '';

    switch (node.getType()) {
      case 'object':
        return this.printObject(node);
      case 'array':
        return this.printArray(node);
      case 'value':
        return this.printValue(node);
      default:
        return '';
    }
  }

  printObject(obj) {
    if (obj.fields.length === 0) {
      return '{}';
    }

    let result = '{\n';
    this.increaseIndent();

    for (let i = 0; i < obj.fields.length; i++) {
      const field = obj.fields[i];
      result += this.indent();
      result += `"${field.key}": `;

      switch (field.value.getType()) {
        case 'object':
          result += this.printObject(field.value);
          break;
        case 'array':
          result += this.printArray(field.value);
          break;
        case 'value':
          result += this.printValue(field.value);
          break;
        default:
          result += 'null';
      }

      if (i < obj.fields.length - 1) {
        result += ',';
      }
      result += '\n';
    }

    this.decreaseIndent();
    result += this.indent() + '}';

    return result;
  }

  printArray(arr) {
    if (arr.items.length === 0) {
      return '[]';
    }

    let result = '[\n';
    this.increaseIndent();

    for (let i = 0; i < arr.items.length; i++) {
      const item = arr.items[i];
      result += this.indent();

      switch (item.getType()) {
        case 'object':
          result += this.printObject(item);
          break;
        case 'array':
          result += this.printArray(item);
          break;
        case 'value':
          result += this.printValue(item);
          break;
        default:
          result += 'null';
      }

      if (i < arr.items.length - 1) {
        result += ',';
      }
      result += '\n';
    }

    this.decreaseIndent();
    result += this.indent() + ']';

    return result;
  }

  printValue(value) {
    if (value.tag && value.tag.isNull) {
      return 'null';
    }

    if (value.data === true || value.data === false) {
      return value.data ? 'true' : 'false';
    }

    if (typeof value.data === 'bigint') {
      return String(value.data);
    }

    if (typeof value.data === 'number') {
      return this.formatNumber(value.data);
    }

    if (typeof value.data === 'string') {
      return this.valueToStringOnly(value);
    }

    if (value.data === null || value.data === undefined) {
      return 'null';
    }

    return this.valueToStringOnly(value);
  }

  valueToStringOnly(value) {
    const tag = value.tag;
    const type = tag?.type || ValueType.Unknown;

    const needsQuotes = this.typeNeedsQuotes(type);

    if (needsQuotes) {
      return `"${value.text}"`;
    }
    return value.text;
  }

  typeNeedsQuotes(type) {
    return type === ValueType.String ||
           type === ValueType.Bytes ||
           type === ValueType.DateTime ||
           type === ValueType.Date ||
           type === ValueType.Time ||
           type === ValueType.Uuid ||
           type === ValueType.Ip ||
           type === ValueType.Url ||
           type === ValueType.Email ||
           type === ValueType.Enum;
  }

  tagToString(tag) {
    const parts = [];

    if (tag.type !== ValueType.Unknown) {
      parts.push(`type=${this.typeToString(tag.type)}`);
    }
    if (tag.desc) {
      parts.push(`desc=${tag.desc}`);
    }
    if (tag.nullable) {
      parts.push('nullable');
    }
    if (tag.isNull) {
      parts.push('is_null');
    }
    if (tag.raw) {
      parts.push('raw');
    }
    if (tag.defaultValue) {
      parts.push(`default=${tag.defaultValue}`);
    }
    if (tag.enumValues) {
      parts.push(`enum=${tag.enumValues}`);
    }

    return parts.join('; ');
  }

  typeToString(type) {
    switch (type) {
      case ValueType.String: return 'str';
      case ValueType.Int: return 'i';
      case ValueType.Int8: return 'i8';
      case ValueType.Int16: return 'i16';
      case ValueType.Int32: return 'i32';
      case ValueType.Int64: return 'i64';
      case ValueType.Uint: return 'u';
      case ValueType.Uint8: return 'u8';
      case ValueType.Uint16: return 'u16';
      case ValueType.Uint32: return 'u32';
      case ValueType.Uint64: return 'u64';
      case ValueType.Float32: return 'f32';
      case ValueType.Float64: return 'f64';
      case ValueType.Bool: return 'bool';
      case ValueType.Bytes: return 'bytes';
      case ValueType.BigInt: return 'bi';
      case ValueType.DateTime: return 'datetime';
      case ValueType.Date: return 'date';
      case ValueType.Time: return 'time';
      case ValueType.Uuid: return 'uuid';
      case ValueType.Decimal: return 'decimal';
      case ValueType.Ip: return 'ip';
      case ValueType.Url: return 'url';
      case ValueType.Email: return 'email';
      case ValueType.Enum: return 'enum';
      case ValueType.Array: return 'arr';
      case ValueType.Struct: return 'struct';
      default: return 'unknown';
    }
  }

  formatNumber(value) {
    if (Number.isNaN(value)) {
      return 'null';
    }
    if (!Number.isFinite(value)) {
      return value > 0 ? 'Infinity' : '-Infinity';
    }
    if (value === Math.floor(value) && Math.abs(value) < Number.MAX_SAFE_INTEGER) {
      return String(value) + '.0';
    }
    return String(value);
  }

  printCompact(node) {
    if (!node) return '';

    switch (node.getType()) {
      case 'object':
        return this.printObjectCompact(node);
      case 'array':
        return this.printArrayCompact(node);
      case 'value':
        return this.printValueCompact(node);
      default:
        return '';
    }
  }

  printObjectCompact(obj) {
    const parts = [];

    for (const field of obj.fields) {
      let part = `"${field.key}":`;

      switch (field.value.getType()) {
        case 'object':
          part += this.printObjectCompact(field.value);
          break;
        case 'array':
          part += this.printArrayCompact(field.value);
          break;
        case 'value':
          part += this.printValueCompact(field.value);
          break;
        default:
          part += 'null';
      }

      parts.push(part);
    }

    return '{' + parts.join(',') + '}';
  }

  printArrayCompact(arr) {
    const parts = [];

    for (const item of arr.items) {
      let part = '';

      switch (item.getType()) {
        case 'object':
          part += this.printObjectCompact(item);
          break;
        case 'array':
          part += this.printArrayCompact(item);
          break;
        case 'value':
          part += this.printValueCompact(item);
          break;
        default:
          part += 'null';
      }

      parts.push(part);
    }

    return '[' + parts.join(',') + ']';
  }

  printValueCompact(value) {
    if (value.tag && value.tag.isNull) {
      return 'null';
    }

    if (value.data === true || value.data === false) {
      return value.data ? 'true' : 'false';
    }

    if (typeof value.data === 'bigint') {
      return String(value.data);
    }

    if (typeof value.data === 'number') {
      return this.formatNumber(value.data);
    }

    if (typeof value.data === 'string') {
      return this.valueToStringOnly(value);
    }

    if (value.data === null || value.data === undefined) {
      return 'null';
    }

    return this.valueToStringOnly(value);
  }
}

export function printJSONC(node) {
  const printer = new JSONCPrinter();
  return printer.print(node);
}

export function printJSONCCompact(node) {
  const printer = new JSONCPrinter();
  return printer.printCompact(node);
}