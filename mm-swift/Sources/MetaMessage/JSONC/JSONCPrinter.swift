import Foundation

private protocol ObjectFieldsProvider {
    var fields: [JSONCField] { get }
    var tag: JSONCTag? { get }
}

extension MMObject: ObjectFieldsProvider {}
extension MMDoc: ObjectFieldsProvider {}

public class JSONCPrinter {
    private var indentLevel: Int
    private let indentString: String
    private let useIndent: Bool

    public init(indentString: String = "  ", useIndent: Bool = true) {
        self.indentLevel = 0
        self.indentString = indentString
        self.useIndent = useIndent
    }

    public func print(_ node: JSONCNode?) -> String {
        guard let node = node else { return "" }

        switch node.getType() {
        case .object:
            return printObject(node as! MMObject)
        case .array:
            return printArray(node as! MMArray)
        case .value:
            return printValue(node as! JSONCValue)
        case .doc:
            return printObject(node as! MMDoc)
        case .unknown:
            return ""
        }
    }

    private func indent() -> String {
        if !useIndent {
            return ""
        }
        return String(repeating: indentString, count: indentLevel)
    }

    private func increaseIndent() {
        indentLevel += 1
    }

    private func decreaseIndent() {
        indentLevel = max(0, indentLevel - 1)
    }

    private func printObject(_ obj: ObjectFieldsProvider) -> String {
        if obj.fields.isEmpty {
            return "{}"
        }

        var result = "{\n"
        increaseIndent()

        for (index, field) in obj.fields.enumerated() {
            result += indent()

            if let tag = obj.tag, !useIndent {
            }

            result += "\"\(field.key)\": "

            switch field.value.getType() {
            case .object:
                result += printObject(field.value as! MMObject)
            case .array:
                result += printArray(field.value as! MMArray)
            case .value:
                result += printValue(field.value as! JSONCValue)
            case .doc:
                result += printObject(field.value as! MMDoc)
            case .unknown:
                result += "null"
            }

            if index < obj.fields.count - 1 {
                result += ","
            }
            result += "\n"
        }

        decreaseIndent()
        result += indent() + "}"

        return result
    }

    private func printArray(_ arr: MMArray) -> String {
        if arr.items.isEmpty {
            return "[]"
        }

        var result = "[\n"
        increaseIndent()

        for (index, item) in arr.items.enumerated() {
            result += indent()

            switch item.getType() {
            case .object:
                result += printObject(item as! MMObject)
            case .array:
                result += printArray(item as! MMArray)
            case .value:
                result += printValue(item as! JSONCValue)
            case .doc:
                result += printObject(item as! MMDoc)
            case .unknown:
                result += "null"
            }

            if index < arr.items.count - 1 {
                result += ","
            }
            result += "\n"
        }

        decreaseIndent()
        result += indent() + "]"

        return result
    }

    private func printValue(_ value: JSONCValue) -> String {
        if let tag = value.tag, tag.isNull {
            return "null"
        }

        if let boolVal = value.data as? Bool {
            return boolVal ? "true" : "false"
        }

        if let intVal = value.data as? Int {
            return String(intVal)
        }

        if let intVal = value.data as? Int64 {
            return String(intVal)
        }

        if let uintVal = value.data as? UInt {
            return String(uintVal)
        }

        if let uintVal = value.data as? UInt64 {
            return String(uintVal)
        }

        if let doubleVal = value.data as? Double {
            return formatDouble(doubleVal)
        }

        if let floatVal = value.data as? Float {
            return formatDouble(Double(floatVal))
        }

        if let stringVal = value.data as? String {
            let needsQuotes = value.tag?.type.needsQuotes ?? true
            if needsQuotes {
                return "\"\(escapeString(stringVal))\""
            } else {
                return escapeString(stringVal)
            }
        }

        if let dataVal = value.data as? Data {
            return "\"\(escapeString(String(data: dataVal, encoding: .utf8) ?? ""))\""
        }

        if value.data == nil {
            return "null"
        }

        let needsQuotes = value.tag?.type.needsQuotes ?? true
        if needsQuotes {
            return "\"\(escapeString(value.text))\""
        } else {
            return escapeString(value.text)
        }
    }

    private func formatDouble(_ value: Double) -> String {
        if value.isNaN {
            return "null"
        }
        if value.isInfinite {
            return value > 0 ? "Infinity" : "-Infinity"
        }

        if value == value.rounded() && abs(value) < Double(Int64.max) {
            return String(format: "%.1f", value)
        }

        return String(value)
    }

    private func escapeString(_ str: String) -> String {
        var result = ""
        for char in str {
            switch char {
            case "\"":
                result += "\\\""
            case "\\":
                result += "\\\\"
            case "\n":
                result += "\\n"
            case "\r":
                result += "\\r"
            case "\t":
                result += "\\t"
            case "\u{8}":
                result += "\\b"
            case "\u{C}":
                result += "\\f"
            default:
                if char.asciiValue != nil {
                    result += String(char)
                } else {
                    result += String(char)
                }
            }
        }
        return result
    }

    public func printCompact(_ node: JSONCNode?) -> String {
        guard let node = node else { return "" }

        switch node.getType() {
        case .object:
            return printObjectCompact(node as! MMObject)
        case .array:
            return printArrayCompact(node as! MMArray)
        case .value:
            return printValueCompact(node as! JSONCValue)
        case .doc:
            return printObjectCompact(node as! MMDoc)
        case .unknown:
            return ""
        }
    }

    private func printObjectCompact(_ obj: ObjectFieldsProvider) -> String {
        var parts: [String] = []

        for field in obj.fields {
            var part = "\"\(field.key)\":"

            switch field.value.getType() {
            case .object:
                part += printObjectCompact(field.value as! MMObject)
            case .array:
                part += printArrayCompact(field.value as! MMArray)
            case .value:
                part += printValueCompact(field.value as! JSONCValue)
            case .doc:
                part += printObjectCompact(field.value as! MMDoc)
            case .unknown:
                part += "null"
            }

            parts.append(part)
        }

        return "{" + parts.joined(separator: ",") + "}"
    }

    private func printArrayCompact(_ arr: MMArray) -> String {
        var parts: [String] = []

        for item in arr.items {
            var part = ""

            switch item.getType() {
            case .object:
                part += printObjectCompact(item as! MMObject)
            case .array:
                part += printArrayCompact(item as! MMArray)
            case .value:
                part += printValueCompact(item as! JSONCValue)
            case .doc:
                part += printObjectCompact(item as! MMDoc)
            case .unknown:
                part += "null"
            }

            parts.append(part)
        }

        return "[" + parts.joined(separator: ",") + "]"
    }

    private func printValueCompact(_ value: JSONCValue) -> String {
        if let tag = value.tag, tag.isNull {
            return "null"
        }

        if let boolVal = value.data as? Bool {
            return boolVal ? "true" : "false"
        }

        if let intVal = value.data as? Int {
            return String(intVal)
        }

        if let intVal = value.data as? Int64 {
            return String(intVal)
        }

        if let uintVal = value.data as? UInt {
            return String(uintVal)
        }

        if let uintVal = value.data as? UInt64 {
            return String(uintVal)
        }

        if let doubleVal = value.data as? Double {
            return formatDouble(doubleVal)
        }

        if let floatVal = value.data as? Float {
            return formatDouble(Double(floatVal))
        }

        if let stringVal = value.data as? String {
            return "\"\(escapeString(stringVal))\""
        }

        if let dataVal = value.data as? Data {
            return "\"\(escapeString(String(data: dataVal, encoding: .utf8) ?? ""))\""
        }

        if value.data == nil {
            return "null"
        }

        return "\"\(escapeString(value.text))\""
    }
}