import Foundation

public enum MMValueToNodeError: Error {
    case untypedNil
    case maxDepthExceeded
    case unsupportedType(String)
}

private let maxDepth = 32

public func valueToNode(_ value: Any?, tag tagStr: String = "") throws -> Node {
    var tag: Tag? = nil
    if !tagStr.isEmpty {
        tag = parseMMTag(tagStr)
    }

    return try valueToNode(value, tag: tag, depth: 0, path: "")
}

func valueToNode(_ value: Any?, tag: Tag?, depth: Int, path: String) throws -> Node {
    var effectiveTag = tag ?? Tag()
    let currentPath = path

    if depth > maxDepth {
        throw MMValueToNodeError.maxDepthExceeded
    }

    guard let v = value else {
        if effectiveTag.type == .unknown {
            throw MMValueToNodeError.untypedNil
        }
        if !effectiveTag.nullable {
            effectiveTag.nullable = true
        }
        effectiveTag.isNull = true
        return NodeScalar(data: nil, text: "null", tag: effectiveTag, path: currentPath)
    }

    return try convertScalar(v, tag: &effectiveTag, depth: depth, path: currentPath)
}

private func convertScalar(_ value: Any, tag: inout Tag, depth: Int, path: String) throws -> Node {
    if depth > maxDepth {
        throw MMValueToNodeError.maxDepthExceeded
    }
    let currentPath = path

    if let boolVal = value as? Bool {
        if tag.type == .unknown { tag.type = .bool }
        return NodeScalar(data: boolVal, text: boolVal ? "true" : "false", tag: tag, path: currentPath)
    }

    if let intVal = value as? Int {
        if tag.type == .unknown { tag.type = .i }
        return NodeScalar(data: intVal, text: String(intVal), tag: tag, path: currentPath)
    }

    if let intVal = value as? Int8 {
        if tag.type == .unknown { tag.type = .i8 }
        return NodeScalar(data: intVal, text: String(intVal), tag: tag, path: currentPath)
    }

    if let intVal = value as? Int16 {
        if tag.type == .unknown { tag.type = .i16 }
        return NodeScalar(data: intVal, text: String(intVal), tag: tag, path: currentPath)
    }

    if let intVal = value as? Int32 {
        if tag.type == .unknown { tag.type = .i32 }
        return NodeScalar(data: intVal, text: String(intVal), tag: tag, path: currentPath)
    }

    if let intVal = value as? Int64 {
        if tag.type == .unknown { tag.type = .i64 }
        return NodeScalar(data: intVal, text: String(intVal), tag: tag, path: currentPath)
    }

    if let uintVal = value as? UInt {
        if tag.type == .unknown { tag.type = .u }
        return NodeScalar(data: uintVal, text: String(uintVal), tag: tag, path: currentPath)
    }

    if let uintVal = value as? UInt8 {
        if tag.type == .unknown { tag.type = .u8 }
        return NodeScalar(data: uintVal, text: String(uintVal), tag: tag, path: currentPath)
    }

    if let uintVal = value as? UInt16 {
        if tag.type == .unknown { tag.type = .u16 }
        return NodeScalar(data: uintVal, text: String(uintVal), tag: tag, path: currentPath)
    }

    if let uintVal = value as? UInt32 {
        if tag.type == .unknown { tag.type = .u32 }
        return NodeScalar(data: uintVal, text: String(uintVal), tag: tag, path: currentPath)
    }

    if let uintVal = value as? UInt64 {
        if tag.type == .unknown { tag.type = .u64 }
        return NodeScalar(data: uintVal, text: String(uintVal), tag: tag, path: currentPath)
    }

    if let floatVal = value as? Float {
        if tag.type == .unknown { tag.type = .f32 }
        return NodeScalar(data: floatVal, text: String(floatVal), tag: tag, path: currentPath)
    }

    if let doubleVal = value as? Double {
        if tag.type == .unknown { tag.type = .f64 }
        return NodeScalar(data: doubleVal, text: String(doubleVal), tag: tag, path: currentPath)
    }

    if let stringVal = value as? String {
        if tag.type == .unknown { tag.type = .str }
        if tag.type == .enums {
            return NodeScalar(data: Int(stringVal), text: stringVal, tag: tag, path: currentPath)
        }
        return NodeScalar(data: stringVal, text: stringVal, tag: tag, path: currentPath)
    }

    if let dataVal = value as? Data {
        if tag.type == .unknown { tag.type = .bytes }
        let text = String(data: dataVal, encoding: .utf8) ?? ""
        return NodeScalar(data: dataVal, text: text, tag: tag, path: currentPath)
    }

    if let dateVal = value as? Date {
        if tag.type == .unknown { tag.type = .datetime }
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        let text = formatter.string(from: dateVal)
        return NodeScalar(data: dateVal, text: text, tag: tag, path: currentPath)
    }

    if let uuidVal = value as? UUID {
        if tag.type == .unknown { tag.type = .uuid }
        return NodeScalar(data: uuidVal, text: uuidVal.uuidString, tag: tag, path: currentPath)
    }

    let mirror = Mirror(reflecting: value)

    if mirror.displayStyle == .dictionary {
        return try dictToNode(value, tag: &tag, depth: depth, path: currentPath)
    }

    if mirror.displayStyle == .collection {
        return try arrayToNode(value, tag: &tag, depth: depth, path: currentPath)
    }

    if mirror.displayStyle == .class || mirror.displayStyle == .struct {
        return try structToNode(value, tag: &tag, depth: depth, path: currentPath)
    }

    if let optionalVal = value as? OptionalProtocol {
        if let wrapped = optionalVal.wrapped {
            return try convertScalar(wrapped, tag: &tag, depth: depth, path: currentPath)
        }
        tag.isNull = true
        if !tag.nullable {
            tag.nullable = true
        }
        return NodeScalar(data: nil, text: "null", tag: tag, path: currentPath)
    }

    throw MMValueToNodeError.unsupportedType(String(describing: type(of: value)))
}

private protocol OptionalProtocol {
    var wrapped: Any? { get }
}

extension Optional: OptionalProtocol {
    var wrapped: Any? {
        switch self {
        case .some(let value):
            return value
        case .none:
            return nil
        }
    }
}

private func structToNode(_ value: Any, tag: inout Tag, depth: Int, path: String) throws -> Node {
    let mirror = Mirror(reflecting: value)
    let newDepth = depth + 1

    let typeName = String(describing: type(of: value))
    if tag.name.isEmpty {
        tag.name = camelToSnake(typeName)
    }

    var currentPath = path
    if currentPath.isEmpty {
        currentPath = tag.name
    } else {
        currentPath = "\(currentPath).\(tag.name)"
    }

    tag.type = .obj

    let obj = NodeObject(tag: tag, path: currentPath)

    for child in mirror.children {
        guard let label = child.label else { continue }

        let fieldKey = camelToSnake(label)
        var childTag = Tag()
        childTag.name = fieldKey

        let subPath = "\(currentPath).\(fieldKey)"

        if let mmTagProvider = value as? MMTagProvider {
            if let fieldTagStr = mmTagProvider.tag(forField: label) {
                if let parsed = parseMMTag(fieldTagStr) {
                    childTag = parsed
                    if childTag.name.isEmpty {
                        childTag.name = fieldKey
                    }
                }
            }
        }

        let childNode = try convertScalar(child.value, tag: &childTag, depth: newDepth, path: subPath)
        let field = Field(key: fieldKey, value: childNode)
        obj.fields.append(field)
    }

    return obj
}

private func arrayToNode(_ value: Any, tag: inout Tag, depth: Int, path: String) throws -> Node {
    let mirror = Mirror(reflecting: value)
    let newDepth = depth + 1

    tag.type = .vec

    let arr = NodeArray(tag: tag, path: path)

    for child in mirror.children {
        var childTag = Tag()
        childTag.inherit(from: tag)

        let subPath = "\(path)[\(arr.items.count)]"
        let childNode = try convertScalar(child.value, tag: &childTag, depth: newDepth, path: subPath)
        arr.items.append(childNode)
    }

    return arr
}

private func dictToNode(_ value: Any, tag: inout Tag, depth: Int, path: String) throws -> Node {
    let mirror = Mirror(reflecting: value)
    let newDepth = depth + 1

    tag.type = .map

    let obj = NodeObject(tag: tag, path: path)

    for child in mirror.children {
        guard let label = child.label else { continue }
        let fieldKey = camelToSnake(label)

        var childTag = Tag()
        childTag.inherit(from: tag)
        childTag.name = fieldKey

        let subPath = "\(path)[\(fieldKey)]"
        let childNode = try convertScalar(child.value, tag: &childTag, depth: newDepth, path: subPath)
        let field = Field(key: fieldKey, value: childNode)
        obj.fields.append(field)
    }

    return obj
}

public protocol MMTagProvider {
    func tag(forField field: String) -> String?
}

public func camelToSnake(_ input: String) -> String {
    var result = ""
    for (index, char) in input.enumerated() {
        if char.isUppercase && index > 0 {
            result += "_"
        }
        result += String(char).lowercased()
    }
    return result
}