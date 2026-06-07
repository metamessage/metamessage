import { NodeScalar, NodeObject, NodeArray, MMDoc } from '../ir/ast';
import { ValueType } from '../ir/value-type';
import { Node } from '../ir/ast';

export class JSONCPrinter {
  private indent: string;
  private indentLevel: number;

  constructor(indent: string = '\t') {
    this.indent = indent;
    this.indentLevel = 0;
  }

  print(node: Node): string {
    this.indentLevel = 0;
    const tag = node.getTag();
    let result = '';

    if (tag.toString() !== '') {
      result += `\n// mm: ${tag.toString()}\n`;
    }

    return (result += this.printNode(node));
  }

  printCompact(node: Node): string {
    this.indentLevel = 0;
    return this.printNodeCompact(node);
  }

  private printNode(node: Node): string {
    if (node instanceof NodeScalar) {
      return this.printValue(node);
    } else if (node instanceof NodeObject) {
      return this.printObject(node);
    } else if (node instanceof NodeArray) {
      return this.printArray(node);
    } else if (node instanceof MMDoc) {
      return this.printNode(node.getRoot());
    }
    return '';
  }

  private printNodeCompact(node: Node): string {
    if (node instanceof NodeScalar) {
      return this.printValueCompact(node);
    } else if (node instanceof NodeObject) {
      return this.printObjectCompact(node);
    } else if (node instanceof NodeArray) {
      return this.printArrayCompact(node);
    } else if (node instanceof MMDoc) {
      return this.printNodeCompact(node.getRoot());
    }
    return '';
  }

  private printValue(value: NodeScalar): string {
    return `${this.valueToStringOnly(value)}`;
  }

  private printValueCompact(value: NodeScalar): string {
    return this.valueToStringOnly(value);
  }

  private valueToStringOnly(value: NodeScalar): string {
    const tag = value.getTag();
    if (tag.isNull) {
      switch (tag.type) {
        case ValueType.I:
        case ValueType.I8:
        case ValueType.I16:
        case ValueType.I32:
        case ValueType.I64:
        case ValueType.U:
        case ValueType.U8:
        case ValueType.U16:
        case ValueType.U32:
        case ValueType.U64:
        case ValueType.Bigint:
          return '0';
        case ValueType.F32:
        case ValueType.F64:
          return '0.0';
        case ValueType.Bool:
          return 'false';
        default:
          return '""';
      }
    }
    const type = tag.type;
    const text = value.getText();

    switch (type) {
      case ValueType.Str:
      case ValueType.Bytes:
      case ValueType.Datetime:
      case ValueType.Date:
      case ValueType.Time:
      case ValueType.Uuid:
      case ValueType.Ip:
      case ValueType.Url:
      case ValueType.Email:
      case ValueType.Enums:
      case ValueType.Media:
        return `"${text}"`;
      default:
        return text;
    }
  }

  private printObject(obj: NodeObject): string {
    const properties = obj.getProperties();
    if (Object.keys(properties).length === 0) {
      const indent = this.getIndent();
      return `{\n${indent}}`;
    }

    this.indentLevel++;
    const indent = this.getIndent();
    const entries: string[] = [];

    for (const [key, value] of Object.entries(properties)) {
      const tag = value.getTag();
      let entry = '';
      if (tag.toString() !== '') {
        entry += `\n${indent}// mm: ${tag.toString()}\n${indent}`;
      } else {
        entry += `${indent}`;
      }
      entry += `${JSON.stringify(key)}: ${this.printNode(value)},`;
      entries.push(entry);
    }

    this.indentLevel--;
    const closingIndent = this.getIndent();

    return `{\n${entries.join('\n')}\n${closingIndent}}`;
  }

  private printObjectCompact(obj: NodeObject): string {
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

  private printArray(array: NodeArray): string {
    const elements = array.getElements();
    if (elements.length === 0) {
      const indent = this.getIndent();
      return `[\n${indent}]`;
    }

    this.indentLevel++;
    const indent = this.getIndent();
    const entries: string[] = [];

    for (const element of elements) {
      const tag = element.getTag();
      if (tag.toString() !== '') {
        entries.push(
          `\n${indent}// mm: ${tag.toString()}\n${indent}${this.printNode(element)},`,
        );
      } else {
        entries.push(`${indent}${this.printNode(element)},`);
      }
    }

    this.indentLevel--;
    const closingIndent = this.getIndent();

    return `[\n${entries.join('\n')}\n${closingIndent}]`;
  }

  private printArrayCompact(array: NodeArray): string {
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
