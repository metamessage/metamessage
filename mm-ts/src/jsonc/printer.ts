import { JSONCValue, JSONCObject, JSONCArray, JSONCDoc } from './ast';

export interface JSONCNode {
  getType(): string;
  getTag(): any;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: any): void;
}

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
      return this.printValue(node);
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
    const val = value.getValue();
    if (val === null) {
      return 'null';
    } else if (typeof val === 'boolean') {
      return val ? 'true' : 'false';
    } else if (typeof val === 'number') {
      return val.toString();
    } else if (typeof val === 'string') {
      return JSON.stringify(val);
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
      entries.push(`${indent}${JSON.stringify(key)}: ${this.printNode(value)}`);
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