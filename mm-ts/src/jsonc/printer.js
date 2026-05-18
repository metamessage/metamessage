"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.JSONCPrinter = void 0;
exports.toJSONC = toJSONC;
exports.printJSONCCompact = printJSONCCompact;
exports.uint8ToBase64 = uint8ToBase64;
const ast_1 = require("../ast/ast");
const value_type_1 = require("../ast/value-type");
class JSONCPrinter {
    constructor(indent = '  ') {
        this.indent = indent;
        this.indentLevel = 0;
    }
    print(node) {
        this.indentLevel = 0;
        const tag = node.getTag();
        let result = '';
        if (tag.toString() !== '') {
            result += `// mm: ${tag.toString()}\n`;
        }
        return (result += this.printNode(node));
    }
    printCompact(node) {
        this.indentLevel = 0;
        return this.printNodeCompact(node);
    }
    printNode(node) {
        if (node instanceof ast_1.MMValue) {
            return this.printValue(node);
        }
        else if (node instanceof ast_1.MMObject) {
            return this.printObject(node);
        }
        else if (node instanceof ast_1.MMArray) {
            return this.printArray(node);
        }
        else if (node instanceof ast_1.MMDoc) {
            return this.printNode(node.getRoot());
        }
        return '';
    }
    printNodeCompact(node) {
        if (node instanceof ast_1.MMValue) {
            return this.printValueCompact(node);
        }
        else if (node instanceof ast_1.MMObject) {
            return this.printObjectCompact(node);
        }
        else if (node instanceof ast_1.MMArray) {
            return this.printArrayCompact(node);
        }
        else if (node instanceof ast_1.MMDoc) {
            return this.printNodeCompact(node.getRoot());
        }
        return '';
    }
    printValue(value) {
        return `${this.valueToStringOnly(value)}`;
    }
    printValueCompact(value) {
        return this.valueToStringOnly(value);
    }
    valueToStringOnly(value) {
        const tag = value.getTag();
        const type = tag.type;
        const val = value.getValue();
        switch (type) {
            case value_type_1.ValueType.Unknown:
                return 'null';
            case value_type_1.ValueType.BigInt:
                return val.toString();
            case value_type_1.ValueType.Bool:
                return val ? 'true' : 'false';
            case value_type_1.ValueType.String:
                return `"${val}"`;
            case value_type_1.ValueType.Bytes:
                return `"${uint8ToBase64(val)}"`;
            case value_type_1.ValueType.DateTime:
                return val.toString();
            case value_type_1.ValueType.Date:
                return val.toString();
            case value_type_1.ValueType.Time:
                return val.toString();
            case value_type_1.ValueType.UUID:
                return `"${val}"`;
            case value_type_1.ValueType.IP:
                return val.toString();
            case value_type_1.ValueType.URL:
                return val.toString();
            case value_type_1.ValueType.Email:
                return `"${val}"`;
            case value_type_1.ValueType.Enum:
                return val.toString();
            case value_type_1.ValueType.Int:
                return val.toString();
            case value_type_1.ValueType.Int8:
                return val.toString();
            case value_type_1.ValueType.Int16:
                return val.toString();
            case value_type_1.ValueType.Int32:
                return val.toString();
            case value_type_1.ValueType.Int64:
                return val.toString();
            case value_type_1.ValueType.Uint:
                return val.toString();
            case value_type_1.ValueType.Uint8:
                return val.toString();
            case value_type_1.ValueType.Uint16:
                return val.toString();
            case value_type_1.ValueType.Uint32:
                return val.toString();
            case value_type_1.ValueType.Uint64:
                return val.toString();
            case value_type_1.ValueType.Float32:
                return val.toString();
            case value_type_1.ValueType.Float64:
                return val.toString();
            default:
                return val.toString();
        }
    }
    printObject(obj) {
        const properties = obj.getProperties();
        if (Object.keys(properties).length === 0) {
            return '{}';
        }
        this.indentLevel++;
        const indent = this.getIndent();
        const entries = [];
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
    printObjectCompact(obj) {
        const properties = obj.getProperties();
        if (Object.keys(properties).length === 0) {
            return '{}';
        }
        const entries = [];
        for (const [key, value] of Object.entries(properties)) {
            entries.push(`${JSON.stringify(key)}:${this.printNodeCompact(value)}`);
        }
        return `{${entries.join(',')}}`;
    }
    printArray(array) {
        const elements = array.getElements();
        if (elements.length === 0) {
            return '[]';
        }
        this.indentLevel++;
        const indent = this.getIndent();
        const entries = [];
        for (const element of elements) {
            entries.push(`${indent}${this.printNode(element)}`);
        }
        this.indentLevel--;
        const closingIndent = this.getIndent();
        return `[\n${entries.join(',\n')}\n${closingIndent}]`;
    }
    printArrayCompact(array) {
        const elements = array.getElements();
        if (elements.length === 0) {
            return '[]';
        }
        const entries = [];
        for (const element of elements) {
            entries.push(this.printNodeCompact(element));
        }
        return `[${entries.join(',')}]`;
    }
    getIndent() {
        return this.indent.repeat(this.indentLevel);
    }
}
exports.JSONCPrinter = JSONCPrinter;
function toJSONC(node) {
    const printer = new JSONCPrinter();
    return printer.print(node);
}
function printJSONCCompact(node) {
    const printer = new JSONCPrinter();
    return printer.printCompact(node);
}

function uint8ToBase64(bytes) {
    return btoa(Array.from(bytes, (c) => String.fromCharCode(c)).join(''));
}
