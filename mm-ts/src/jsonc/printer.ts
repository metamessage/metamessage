import { MMValue, MMObject, MMArray, MMDoc } from '../ir/ast';
import { ValueType } from '../ir/value-type';
import { Node } from '../ir/ast';

export class JSONCPrinter {
  private indent: string;
  private indentLevel: number;

  constructor(indent: string = '  ') {
    this.indent = indent;
    this.indentLevel = 0;
  }

  print(node: Node): string {
    this.indentLevel = 0;
    const tag = node.getTag();
    let result = '';

    if (tag.toString() !== '') {
      result += `// mm: ${tag.toString()}\n`;
    }

    return (result += this.printNode(node));
  }

  printCompact(node: Node): string {
    this.indentLevel = 0;
    return this.printNodeCompact(node);
  }

  private printNode(node: Node): string {
    if (node instanceof MMValue) {
      return this.printValue(node);
    } else if (node instanceof MMObject) {
      return this.printObject(node);
    } else if (node instanceof MMArray) {
      return this.printArray(node);
    } else if (node instanceof MMDoc) {
      return this.printNode(node.getRoot());
    }
    return '';
  }

  private printNodeCompact(node: Node): string {
    if (node instanceof MMValue) {
      return this.printValueCompact(node);
    } else if (node instanceof MMObject) {
      return this.printObjectCompact(node);
    } else if (node instanceof MMArray) {
      return this.printArrayCompact(node);
    } else if (node instanceof MMDoc) {
      return this.printNodeCompact(node.getRoot());
    }
    return '';
  }

  private printValue(value: MMValue): string {
    return `${this.valueToStringOnly(value)}`;
  }

  private printValueCompact(value: MMValue): string {
    return this.valueToStringOnly(value);
  }

  private valueToStringOnly(value: MMValue): string {
    const tag = value.getTag();
    const type = tag.type;
    const text = value.getText();
    const val = value.getValue();

    if (text) {
      switch (type) {
        case ValueType.Str:
        case ValueType.Bytes:
        case ValueType.DateTime:
        case ValueType.Date:
        case ValueType.Time:
        case ValueType.UUID:
        case ValueType.IP:
        case ValueType.URL:
        case ValueType.Email:
        case ValueType.Enum:
          return `"${text}"`;
        default:
          return text;
      }
    }

    switch (type) {
      case ValueType.Unknown:
        return 'null';
      case ValueType.Str:
      case ValueType.UUID:
      case ValueType.Email:
        return `"${val}"`;
      case ValueType.Bytes:
        return `"${uint8ToBase64(val)}"`;
      case ValueType.DateTime:
      case ValueType.Date:
      case ValueType.Time:
        return `"${this.dateToText(val)}"`;
      case ValueType.IP:
      case ValueType.URL:
      case ValueType.Enum:
        return `"${val}"`;
      case ValueType.Bool:
        return val ? 'true' : 'false';
      case ValueType.BigInt:
      case ValueType.Int:
      case ValueType.Int8:
      case ValueType.Int16:
      case ValueType.Int32:
      case ValueType.Int64:
      case ValueType.Uint:
      case ValueType.Uint8:
      case ValueType.Uint16:
      case ValueType.Uint32:
      case ValueType.Uint64:
      case ValueType.Float32:
      case ValueType.Float64:
      default:
        return String(val);
    }
  }

  private dateToText(val: any): string {
    if (val instanceof Date) {
      const pad = (n: number) => String(n).padStart(2, '0');
      return `${val.getFullYear()}-${pad(val.getMonth() + 1)}-${pad(val.getDate())} ${pad(val.getHours())}:${pad(val.getMinutes())}:${pad(val.getSeconds())}`;
    }
    return String(val);
  }

  private printObject(obj: MMObject): string {
    const properties = obj.getProperties();
    if (Object.keys(properties).length === 0) {
      return '{}';
    }

    this.indentLevel++;
    const indent = this.getIndent();
    const entries: string[] = [];

    for (const [key, value] of Object.entries(properties)) {
      const tag = value.getTag();
      let entry = '';
      if (tag.toString() !== '') {
        entry += `${indent}// mm: ${tag.toString()}\n${indent}`;
      }
      entry += `${JSON.stringify(key)}: ${this.printNode(value)}`;
      entries.push(entry);
    }

    this.indentLevel--;
    const closingIndent = this.getIndent();

    return `{\n${entries.join(',\n\n')}\n${closingIndent}}`;
  }

  private printObjectCompact(obj: MMObject): string {
    const properties = obj.getProperties();
    if (Object.keys(properties).length === 0) {
      return '{}';
    }

    const entries: string[] = [];
    for (const [key, value] of Object.entries(properties)) {
      entries.push(`${JSON.stringify(key)}:${this.printNodeCompact(value)}`);
    }

    return `{${entries.join(',')}}`;
  }

  private printArray(array: MMArray): string {
    const elements = array.getElements();
    if (elements.length === 0) {
      return '[]';
    }

    this.indentLevel++;
    const indent = this.getIndent();
    const entries: string[] = [];

    for (const element of elements) {
      entries.push(`${indent}${this.printNode(element)}`);
    }

    this.indentLevel--;
    const closingIndent = this.getIndent();

    return `[\n${entries.join(',\n')}\n${closingIndent}]`;
  }

  private printArrayCompact(array: MMArray): string {
    const elements = array.getElements();
    if (elements.length === 0) {
      return '[]';
    }

    const entries: string[] = [];
    for (const element of elements) {
      entries.push(this.printNodeCompact(element));
    }

    return `[${entries.join(',')}]`;
  }

  private getIndent(): string {
    return this.indent.repeat(this.indentLevel);
  }
}

export function toJSONC(node: Node): string {
  const printer = new JSONCPrinter();
  return printer.print(node);
}

export function printJSONCCompact(node: Node): string {
  const printer = new JSONCPrinter();
  return printer.printCompact(node);
}

export function uint8ToBase64(bytes: Uint8Array): string {
  return btoa(Array.from(bytes, (c) => String.fromCharCode(c)).join(''));
}

export function base64ToUint8(base64: string): Uint8Array {
  const binaryString = atob(base64);
  const length = binaryString.length;
  const bytes = new Uint8Array(length);

  for (let i = 0; i < length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return bytes;
}
