import Foundation

private protocol ObjectFieldsProvider {
    var fields: [Field] { get }
    var tag: Tag? { get }
}

extension NodeObject: ObjectFieldsProvider {}
extension MMDoc: ObjectFieldsProvider {}

public class JSONCPrinter {
    private var indentLevel: Int
    private let indentString: String
    private let useIndent: Bool

    public init(indentString: String = "\t", useIndent: Bool = true) {
        self.indentLevel = 0
        self.indentString = indentString
        self.useIndent = useIndent
    }

    public func print(_ node: Node?) -> String {
        guard let node = node else { return "" }

        var result = ""
        if let tag = node.getTag() {
            let tagStr = tag.stringValue()
            if !tagStr.isEmpty {
                result += "\n"
                result += "// mm: \(tagStr)\n"
            }
        }

        switch node.getType() {
        case .object:
            result += printObject(node as! NodeObject)
        case .array:
            result += printArray(node as! NodeArray)
        case .value:
            if node is NodeNull {
                result += "null"
            } else {
                result += printValue(node as! NodeScalar)
            }
        case .doc:
            result += printObject(node as! MMDoc)
        case .unknown:
            result += ""
        }
        return result
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
            return "{\n" + indent() + "}"
        }

        var result = "{\n"
        increaseIndent()

        for (_, field) in obj.fields.enumerated() {
            if let tag = field.value.getTag() {
                    let tagStr = tag.stringValue()
                    if !tagStr.isEmpty {
                        result += "\n"
                        result += indent()
                        result += "// mm: \(tagStr)\n"
                    }
                }

            result += indent()
            result += "\"\(field.key)\": "

            switch field.value.getType() {
            case .object:
                result += printObject(field.value as! NodeObject)
            case .array:
                result += printArray(field.value as! NodeArray)
            case .value:
                result += printValue(field.value as! NodeScalar)
            case .doc:
                result += printObject(field.value as! MMDoc)
            case .unknown:
                result += "null"
            }

            result += ","
            result += "\n"
        }

        decreaseIndent()
        result += indent() + "}"

        return result
    }

    private func printArray(_ arr: NodeArray) -> String {
        if arr.items.isEmpty {
            return "[\n" + indent() + "]"
        }

        var result = "[\n"
        increaseIndent()

        for (_, item) in arr.items.enumerated() {
            if let tag = item.getTag() {
                let tagStr = tag.stringValue()
                if !tagStr.isEmpty {
                    result += "\n"
                    result += indent()
                    result += "// mm: \(tagStr)\n"
                }
            }

            result += indent()

            switch item.getType() {
            case .object:
                result += printObject(item as! NodeObject)
            case .array:
                result += printArray(item as! NodeArray)
            case .value:
                result += printValue(item as! NodeScalar)
            case .doc:
                result += printObject(item as! MMDoc)
            case .unknown:
                result += "null"
            }

            result += ","
            result += "\n"
        }

        decreaseIndent()
        result += indent() + "]"

        return result
    }

    private func printValue(_ value: NodeScalar) -> String {
        if let tag = value.tag, tag.isNull {
            if let boolVal = value.data as? Bool {
                return boolVal ? "true" : "false"
            }
            if value.data is Int || value.data is Int64 || value.data is UInt || value.data is UInt64 {
                return "0"
            }
            if value.data is Float || value.data is Double {
                return "0.0"
            }
            if value.data is String {
                return "\"\""
            }
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
            // For encoded types (bytes, media, etc.), use the text representation (base64)
            if value.text.isEmpty {
                return "\"\(escapeString(String(data: dataVal, encoding: .utf8) ?? ""))\""
            }
            return "\"\(escapeString(value.text))\""
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

    public func printCompact(_ node: Node?) -> String {
        guard let node = node else { return "" }

        switch node.getType() {
        case .object:
            return printObjectCompact(node as! NodeObject)
        case .array:
            return printArrayCompact(node as! NodeArray)
        case .value:
            if node is NodeNull {
                return "null"
            }
            return printValueCompact(node as! NodeScalar)
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
                part += printObjectCompact(field.value as! NodeObject)
            case .array:
                part += printArrayCompact(field.value as! NodeArray)
            case .value:
                part += printValueCompact(field.value as! NodeScalar)
            case .doc:
                part += printObjectCompact(field.value as! MMDoc)
            case .unknown:
                part += "null"
            }

            parts.append(part)
        }

        return "{" + parts.joined(separator: ",") + "}"
    }

    private func printArrayCompact(_ arr: NodeArray) -> String {
        var parts: [String] = []

        for item in arr.items {
            var part = ""

            switch item.getType() {
            case .object:
                part += printObjectCompact(item as! NodeObject)
            case .array:
                part += printArrayCompact(item as! NodeArray)
            case .value:
                part += printValueCompact(item as! NodeScalar)
            case .doc:
                part += printObjectCompact(item as! MMDoc)
            case .unknown:
                part += "null"
            }

            parts.append(part)
        }

        return "[" + parts.joined(separator: ",") + "]"
    }

    private func printValueCompact(_ value: NodeScalar) -> String {
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
            // For encoded types (bytes, media, etc.), use the text representation (base64)
            if value.text.isEmpty {
                return "\"\(escapeString(String(data: dataVal, encoding: .utf8) ?? ""))\""
            }
            return "\"\(escapeString(value.text))\""
        }

        if value.data == nil {
            return "null"
        }

        return "\"\(escapeString(value.text))\""
    }
}