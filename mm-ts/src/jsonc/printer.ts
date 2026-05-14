import { MMValue, MMObject, MMArray, MMDoc } from '../ast/ast';
import { ValueType } from '../ast/value-type';
import { Node } from '../ast/ast';

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
    const val = value.getValue();
    switch (type) {
      case ValueType.Unknown:
        return 'null';
      case ValueType.BigInt:
        return val.toString();

      case ValueType.Bool:
        return val ? 'true' : 'false';
      case ValueType.String:
        return `"${val}"`;
      case ValueType.Bytes:
        return `"${uint8ToBase64(val)}"`;
      case ValueType.DateTime:
        return val.toString();
      case ValueType.Date:
        return val.toString();
      case ValueType.Time:
        return val.toString();
      case ValueType.UUID:
        return `"${val}"`;
      case ValueType.IP:
        return val.toString();
      case ValueType.URL:
        return val.toString();
      case ValueType.Email:
        return `"${val}"`;
      case ValueType.Enum:
        return val.toString();
      case ValueType.Int:
        return val.toString();
      case ValueType.Int8:
        return val.toString();
      case ValueType.Int16:
        return val.toString();
      case ValueType.Int32:
        return val.toString();
      case ValueType.Int64:
        return val.toString();
      case ValueType.Uint:
        return val.toString();
      case ValueType.Uint8:
        return val.toString();
      case ValueType.Uint16:
        return val.toString();
      case ValueType.Uint32:
        return val.toString();
      case ValueType.Uint64:
        return val.toString();
      case ValueType.Float32:
        return val.toString();
      case ValueType.Float64:
        return val.toString();
      default:
        return val.toString();
    }
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

function uint8ToBase64(bytes: Uint8Array): string {
  return btoa(Array.from(bytes, (c) => String.fromCharCode(c)).join(''));
}
