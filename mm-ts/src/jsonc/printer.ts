import { JSONCValue, JSONCObject, JSONCArray, JSONCDoc, JSONCTag, JSONCValueType } from './ast';

export interface JSONCNode {
  getType(): string;
  getTag(): any;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: any): void;
}

const QUOTED_TYPES: JSONCValueType[] = [
  JSONCValueType.String,
  JSONCValueType.Bytes,
  JSONCValueType.DateTime,
  JSONCValueType.Date,
  JSONCValueType.Time,
  JSONCValueType.UUID,
  JSONCValueType.IP,
  JSONCValueType.URL,
  JSONCValueType.Email,
  JSONCValueType.Enum,
];

export class JSONCPrinter {
  private indent: string;
  private indentLevel: number;

  constructor(indent: string = '  ') {
    this.indent = indent;
    this.indentLevel = 0;
  }

  print(node: JSONCNode): string {
    this.indentLevel = 0;
    return this.printNode(node);
  }

  printCompact(node: JSONCNode): string {
    this.indentLevel = 0;
    return this.printNodeCompact(node);
  }

  private printNode(node: JSONCNode): string {
    if (node instanceof JSONCValue) {
      return this.printValue(node);
    } else if (node instanceof JSONCObject) {
      return this.printObject(node);
    } else if (node instanceof JSONCArray) {
      return this.printArray(node);
    } else if (node instanceof JSONCDoc) {
      return this.printNode(node.getRoot());
    }
    return '';
  }

  private printNodeCompact(node: JSONCNode): string {
    if (node instanceof JSONCValue) {
      return this.printValueCompact(node);
    } else if (node instanceof JSONCObject) {
      return this.printObjectCompact(node);
    } else if (node instanceof JSONCArray) {
      return this.printArrayCompact(node);
    } else if (node instanceof JSONCDoc) {
      return this.printNodeCompact(node.getRoot());
    }
    return '';
  }

  private printValue(value: JSONCValue): string {
    const tag = value.getTag();
    let result = '';

    if (tag && tag.desc) {
      result += `// mm: ${tag.toString()}\n${this.getIndent()}`;
    }

    result += this.valueToStringOnly(value);
    return result;
  }

  private printValueCompact(value: JSONCValue): string {
    return this.valueToStringOnly(value);
  }

  private valueToStringOnly(value: JSONCValue): string {
    const tag = value.getTag();
    const type = tag?.type || JSONCValueType.Unknown;
    const needsQuotes = QUOTED_TYPES.includes(type);

    const val = value.getValue();
    if (val === null) {
      return 'null';
    } else if (typeof val === 'boolean') {
      return val ? 'true' : 'false';
    } else if (typeof val === 'number') {
      return val.toString();
    } else if (typeof val === 'string') {
      return needsQuotes ? `"${val}"` : JSON.stringify(val);
    }
    return 'null';
  }

  private printObject(obj: JSONCObject): string {
    const properties = obj.getProperties();
    if (properties.size === 0) {
      return '{}';
    }

    this.indentLevel++;
    const indent = this.getIndent();
    const entries: string[] = [];

    for (const [key, value] of properties.entries()) {
      const tag = value.getTag();
      let entry = '';
      if (tag && tag.desc && value instanceof JSONCValue) {
        entry += `${indent}// mm: ${tag.toString()}\n${indent}`;
      }
      entry += `${JSON.stringify(key)}: ${this.printNode(value)}`;
      entries.push(entry);
    }

    this.indentLevel--;
    const closingIndent = this.getIndent();

    return `{\n${entries.join(',\n')}\n${closingIndent}}`;
  }

  private printObjectCompact(obj: JSONCObject): string {
    const properties = obj.getProperties();
    if (properties.size === 0) {
      return '{}';
    }

    const entries: string[] = [];
    for (const [key, value] of properties.entries()) {
      entries.push(`${JSON.stringify(key)}:${this.printNodeCompact(value)}`);
    }

    return `{${entries.join(',')}}`;
  }

  private printArray(array: JSONCArray): string {
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

  private printArrayCompact(array: JSONCArray): string {
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

export function printJSONC(node: JSONCNode): string {
  const printer = new JSONCPrinter();
  return printer.print(node);
}

export function printJSONCCompact(node: JSONCNode): string {
  const printer = new JSONCPrinter();
  return printer.printCompact(node);
}