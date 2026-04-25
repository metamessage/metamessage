import { TokenType, JSONCScanner } from './scanner.js';

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
      return JSON.stringify(value.data);
    }

    if (value.data === null || value.data === undefined) {
      return 'null';
    }

    return JSON.stringify(value.text);
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
      return JSON.stringify(value.data);
    }

    if (value.data === null || value.data === undefined) {
      return 'null';
    }

    return JSON.stringify(value.text);
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