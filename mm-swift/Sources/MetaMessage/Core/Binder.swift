import Foundation

public enum MMBindError: Error {
    case notAnObject
    case fieldNotFound(String)
    case typeMismatch(expected: String, actual: String)
    case unsupportedType(String)
    case notAClass
}

public func bindJSONCNode(_ node: JSONCNode, to out: AnyObject) throws {
    guard let obj = node as? MMObject else {
        throw MMBindError.notAnObject
    }
    try bindObject(obj, to: out)
}

public func bindObject(_ obj: MMObject, to out: AnyObject) throws {
    let mirror = Mirror(reflecting: out)

    for field in obj.fields {
        let fieldKey = field.key
        let camelKey = snakeToCamel(fieldKey)

        let propertyName = findPropertyName(in: mirror, matching: camelKey) ?? camelKey

        let childValue = field.value

        if let valueNode = childValue as? JSONCValue {
            let boundValue = try extractValue(from: valueNode)
            setProperty(propertyName, value: boundValue, on: out)
        } else if let childObj = childValue as? MMObject {
            if childObj.tag?.type == .structType || childObj.tag?.type == .unknown || childObj.tag == nil {
                if let nestedObject = createNestedInstance(for: propertyName, on: out, mirror: mirror) {
                    try bindObject(childObj, to: nestedObject)
                    setProperty(propertyName, value: nestedObject, on: out)
                    continue
                }
            }
            let dict = try extractDict(from: childObj)
            setProperty(propertyName, value: dict, on: out)
        } else if let childArr = childValue as? MMArray {
            let arr = try extractArray(from: childArr)
            setProperty(propertyName, value: arr, on: out)
        }
    }
}

private func findPropertyName(in mirror: Mirror, matching snakeKey: String) -> String? {
    let loweredSnake = snakeKey.lowercased()
    for child in mirror.children {
        guard let label = child.label else { continue }
        if label.lowercased() == loweredSnake {
            return label
        }
    }
    return nil
}

private func extractValue(from node: JSONCValue) throws -> Any? {
    let tag = node.getTag()
    if tag?.isNull == true {
        return nil
    }
    return node.data
}

private func extractDict(from obj: MMObject) throws -> [String: Any] {
    var result: [String: Any] = [:]
    for field in obj.fields {
        if let valueNode = field.value as? JSONCValue {
            result[field.key] = valueNode.data ?? NSNull()
        } else if let childObj = field.value as? MMObject {
            result[field.key] = try extractDict(from: childObj)
        } else if let childArr = field.value as? MMArray {
            result[field.key] = try extractArray(from: childArr)
        }
    }
    return result
}

private func extractArray(from arr: MMArray) throws -> [Any] {
    var result: [Any] = []
    for item in arr.items {
        if let valueNode = item as? JSONCValue {
            let tag = valueNode.getTag()
            if tag?.isNull == true {
                result.append(NSNull())
            } else {
                result.append(valueNode.data ?? NSNull())
            }
        } else if let childObj = item as? MMObject {
            result.append(try extractDict(from: childObj))
        } else if let childArr = item as? MMArray {
            result.append(try extractArray(from: childArr))
        }
    }
    return result
}

private func createNestedInstance(for propertyName: String, on out: AnyObject, mirror: Mirror) -> AnyObject? {
    for child in mirror.children {
        guard let label = child.label else { continue }
        if label == propertyName {
            let value = child.value
            if let obj = value as? NSObject {
                return type(of: obj).init()
            }
        }
    }
    return nil
}

private func setProperty(_ name: String, value: Any?, on object: AnyObject) {
    let sel = NSSelectorFromString("set\(name.prefix(1).uppercased())\(name.dropFirst())):")
    if object.responds(to: sel) {
        object.setValue(value, forKey: name)
        return
    }
    object.setValue(value, forKey: name)
}

public func snakeToCamel(_ input: String) -> String {
    let parts = input.split(separator: "_", omittingEmptySubsequences: true)
    var result = ""
    for (i, part) in parts.enumerated() {
        if i == 0 {
            result += part.lowercased()
        } else {
            result += part.capitalized
        }
    }
    return result
}