import Foundation

public enum JSONCBinderError: Error {
    case typeMismatch(expected: String, actual: String)
    case invalidValue(String)
    case missingRequiredField(String)
    case conversionFailed(String)
}

public class JSONCBinder {
    public init() {}

    public func bind<T: Decodable>(_ node: JSONCNode, to type: T.Type) throws -> T {
        let decoder = JSONCDecoder(node: node)
        return try T(from: decoder)
    }

    public func bindToObject(_ node: JSONCNode, to target: Any) throws {
        guard let obj = node as? JSONCObject else {
            throw JSONCBinderError.typeMismatch(expected: "Object", actual: "\(node.getType())")
        }

        let mirror = Mirror(reflecting: target)
        try bindObject(obj, to: target, mirror: mirror)
    }

    private func bindObject(_ obj: JSONCObject, to target: Any, mirror: Mirror) throws {
        for child in mirror.children {
            guard let label = child.label else { continue }

            let snakeKey = camelToSnake(label)

            if let field = obj.fields.first(where: { $0.key == snakeKey }) {
                try bindValue(field.value, to: child.value, label: label)
            }
        }
    }

    private func bindValue(_ node: JSONCNode, to target: Any, label: String) throws {
        switch target {
        case is Bool:
            guard let valueNode = node as? JSONCValue, let boolVal = valueNode.data as? Bool else {
                throw JSONCBinderError.typeMismatch(expected: "Bool", actual: "\(type(of: node))")
            }
            if let field = target as? Bool {
                _ = field
            }

        case is Int:
            guard let valueNode = node as? JSONCValue else {
                throw JSONCBinderError.typeMismatch(expected: "Int", actual: "\(type(of: node))")
            }
            if let intVal = valueNode.data as? Int {
                _ = intVal
            } else if let int64Val = valueNode.data as? Int64 {
                _ = Int(int64Val)
            }

        case is Int8:
            guard let valueNode = node as? JSONCValue, let intVal = valueNode.data as? Int else {
                throw JSONCBinderError.typeMismatch(expected: "Int8", actual: "\(type(of: node))")
            }
            _ = Int8(intVal)

        case is Int16:
            guard let valueNode = node as? JSONCValue, let intVal = valueNode.data as? Int else {
                throw JSONCBinderError.typeMismatch(expected: "Int16", actual: "\(type(of: node))")
            }
            _ = Int16(intVal)

        case is Int32:
            guard let valueNode = node as? JSONCValue, let intVal = valueNode.data as? Int else {
                throw JSONCBinderError.typeMismatch(expected: "Int32", actual: "\(type(of: node))")
            }
            _ = Int32(intVal)

        case is Int64:
            guard let valueNode = node as? JSONCValue, let intVal = valueNode.data as? Int64 else {
                throw JSONCBinderError.typeMismatch(expected: "Int64", actual: "\(type(of: node))")
            }
            _ = intVal

        case is UInt:
            guard let valueNode = node as? JSONCValue, let uintVal = valueNode.data as? UInt else {
                throw JSONCBinderError.typeMismatch(expected: "UInt", actual: "\(type(of: node))")
            }
            _ = uintVal

        case is UInt8:
            guard let valueNode = node as? JSONCValue, let uintVal = valueNode.data as? UInt else {
                throw JSONCBinderError.typeMismatch(expected: "UInt8", actual: "\(type(of: node))")
            }
            _ = UInt8(uintVal)

        case is UInt16:
            guard let valueNode = node as? JSONCValue, let uintVal = valueNode.data as? UInt else {
                throw JSONCBinderError.typeMismatch(expected: "UInt16", actual: "\(type(of: node))")
            }
            _ = UInt16(uintVal)

        case is UInt32:
            guard let valueNode = node as? JSONCValue, let uintVal = valueNode.data as? UInt else {
                throw JSONCBinderError.typeMismatch(expected: "UInt32", actual: "\(type(of: node))")
            }
            _ = UInt32(uintVal)

        case is UInt64:
            guard let valueNode = node as? JSONCValue, let uintVal = valueNode.data as? UInt64 else {
                throw JSONCBinderError.typeMismatch(expected: "UInt64", actual: "\(type(of: node))")
            }
            _ = uintVal

        case is Float:
            guard let valueNode = node as? JSONCValue, let doubleVal = valueNode.data as? Double else {
                throw JSONCBinderError.typeMismatch(expected: "Float", actual: "\(type(of: node))")
            }
            _ = Float(doubleVal)

        case is Double:
            guard let valueNode = node as? JSONCValue, let doubleVal = valueNode.data as? Double else {
                throw JSONCBinderError.typeMismatch(expected: "Double", actual: "\(type(of: node))")
            }
            _ = doubleVal

        case is String:
            guard let valueNode = node as? JSONCValue, let stringVal = valueNode.data as? String else {
                throw JSONCBinderError.typeMismatch(expected: "String", actual: "\(type(of: node))")
            }
            _ = stringVal

        case is Data:
            guard let valueNode = node as? JSONCValue, let dataVal = valueNode.data as? Data else {
                throw JSONCBinderError.typeMismatch(expected: "Data", actual: "\(type(of: node))")
            }
            _ = dataVal

        case is [String: Any]:
            if let obj = node as? JSONCObject {
                var dict: [String: Any] = [:]
                for field in obj.fields {
                    dict[field.key] = field.value
                }
                _ = dict
            }

        case is [Any]:
            if let arr = node as? JSONCArray {
                var result: [Any] = []
                for item in arr.items {
                    result.append(item)
                }
                _ = result
            }

        default:
            if let obj = node as? JSONCObject {
                let targetMirror = Mirror(reflecting: target)
                try bindObject(obj, to: target, mirror: targetMirror)
            }
        }
    }

    private func camelToSnake(_ input: String) -> String {
        var result = ""
        for (index, char) in input.enumerated() {
            if char.isUppercase && index > 0 {
                result += "_"
            }
            result += String(char).lowercased()
        }
        return result
    }
}

private class JSONCDecoder: Decoder {
    let node: JSONCNode

    init(node: JSONCNode) {
        self.node = node
    }

    var codingPath: [CodingKey] = []

    var userInfo: [CodingUserInfoKey: Any] = [:]

    func container<Key>(keyedBy type: Key.Type) throws -> KeyedDecodingContainer<Key> where Key: CodingKey {
        guard let obj = node as? JSONCObject else {
            throw JSONCBinderError.typeMismatch(expected: "Object", actual: "\(type(of: node))")
        }
        return KeyedDecodingContainer(JSONCKeyedDecodingContainer(node: obj, codingPath: codingPath))
    }

    func unkeyedContainer() throws -> UnkeyedDecodingContainer {
        guard let arr = node as? JSONCArray else {
            throw JSONCBinderError.typeMismatch(expected: "Array", actual: "\(type(of: node))")
        }
        return JSONCUnkeyedDecodingContainer(node: arr, codingPath: codingPath)
    }

    func singleValueContainer() throws -> SingleValueDecodingContainer {
        guard let value = node as? JSONCValue else {
            throw JSONCBinderError.typeMismatch(expected: "Value", actual: "\(type(of: node))")
        }
        return JSONCSingleValueDecodingContainer(node: value, codingPath: codingPath)
    }
}

private struct JSONCKeyedDecodingContainer<Key: CodingKey>: KeyedDecodingContainerProtocol {
    let node: JSONCObject
    let codingPath: [CodingKey]

    var allKeys: [Key] {
        return node.fields.compactMap { Key(stringValue: $0.key) }
    }

    func contains(_ key: Key) -> Bool {
        return node.fields.contains(where: { $0.key == key.stringValue })
    }

    func decode(_ type: Bool.Type, forKey key: Key) throws -> Bool {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let boolVal = valueNode.data as? Bool else {
            throw JSONCBinderError.invalidValue("Bool value for key \(key.stringValue)")
        }
        return boolVal
    }

    func decode(_ type: Int.Type, forKey key: Key) throws -> Int {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int value for key \(key.stringValue)")
        }
        return intVal
    }

    func decode(_ type: Int8.Type, forKey key: Key) throws -> Int8 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int8 value for key \(key.stringValue)")
        }
        return Int8(intVal)
    }

    func decode(_ type: Int16.Type, forKey key: Key) throws -> Int16 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int16 value for key \(key.stringValue)")
        }
        return Int16(intVal)
    }

    func decode(_ type: Int32.Type, forKey key: Key) throws -> Int32 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int32 value for key \(key.stringValue)")
        }
        return Int32(intVal)
    }

    func decode(_ type: Int64.Type, forKey key: Key) throws -> Int64 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int64 else {
            throw JSONCBinderError.invalidValue("Int64 value for key \(key.stringValue)")
        }
        return intVal
    }

    func decode(_ type: UInt.Type, forKey key: Key) throws -> UInt {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt value for key \(key.stringValue)")
        }
        return uintVal
    }

    func decode(_ type: UInt8.Type, forKey key: Key) throws -> UInt8 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt8 value for key \(key.stringValue)")
        }
        return UInt8(uintVal)
    }

    func decode(_ type: UInt16.Type, forKey key: Key) throws -> UInt16 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt16 value for key \(key.stringValue)")
        }
        return UInt16(uintVal)
    }

    func decode(_ type: UInt32.Type, forKey key: Key) throws -> UInt32 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt32 value for key \(key.stringValue)")
        }
        return UInt32(uintVal)
    }

    func decode(_ type: UInt64.Type, forKey key: Key) throws -> UInt64 {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let uintVal = valueNode.data as? UInt64 else {
            throw JSONCBinderError.invalidValue("UInt64 value for key \(key.stringValue)")
        }
        return uintVal
    }

    func decode(_ type: Float.Type, forKey key: Key) throws -> Float {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let doubleVal = valueNode.data as? Double else {
            throw JSONCBinderError.invalidValue("Float value for key \(key.stringValue)")
        }
        return Float(doubleVal)
    }

    func decode(_ type: Double.Type, forKey key: Key) throws -> Double {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let doubleVal = valueNode.data as? Double else {
            throw JSONCBinderError.invalidValue("Double value for key \(key.stringValue)")
        }
        return doubleVal
    }

    func decode(_ type: String.Type, forKey key: Key) throws -> String {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let stringVal = valueNode.data as? String else {
            throw JSONCBinderError.invalidValue("String value for key \(key.stringValue)")
        }
        return stringVal
    }

    func decode(_ type: Data.Type, forKey key: Key) throws -> Data {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let dataVal = valueNode.data as? Data else {
            throw JSONCBinderError.invalidValue("Data value for key \(key.stringValue)")
        }
        return dataVal
    }

    func decode<T>(_ type: T.Type, forKey key: Key) throws -> T where T: Decodable {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }) else {
            throw JSONCBinderError.missingRequiredField(key.stringValue)
        }

        let decoder = JSONCDecoder(node: field.value)
        return try T(from: decoder)
    }

    func decodeIfPresent(_ type: Bool.Type, forKey key: Key) throws -> Bool? {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let boolVal = valueNode.data as? Bool else {
            return nil
        }
        return boolVal
    }

    func decodeIfPresent(_ type: Int.Type, forKey key: Key) throws -> Int? {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            return nil
        }
        return intVal
    }

    func decodeIfPresent(_ type: String.Type, forKey key: Key) throws -> String? {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let valueNode = field.value as? JSONCValue,
              let stringVal = valueNode.data as? String else {
            return nil
        }
        return stringVal
    }

    func decodeIfPresent<T>(_ type: T.Type, forKey key: Key) throws -> T? where T: Decodable {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }) else {
            return nil
        }

        let decoder = JSONCDecoder(node: field.value)
        return try T(from: decoder)
    }

    func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type, forKey key: Key) throws -> KeyedDecodingContainer<NestedKey> where NestedKey: CodingKey {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let obj = field.value as? JSONCObject else {
            throw JSONCBinderError.invalidValue("Nested object for key \(key.stringValue)")
        }
        return KeyedDecodingContainer(JSONCKeyedDecodingContainer<NestedKey>(node: obj, codingPath: codingPath))
    }

    func nestedUnkeyedContainer(forKey key: Key) throws -> UnkeyedDecodingContainer {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }),
              let arr = field.value as? JSONCArray else {
            throw JSONCBinderError.invalidValue("Nested array for key \(key.stringValue)")
        }
        return JSONCUnkeyedDecodingContainer(node: arr, codingPath: codingPath)
    }

    func superDecoder() throws -> Decoder {
        return JSONCDecoder(node: node)
    }

    func superDecoder(forKey key: Key) throws -> Decoder {
        guard let field = node.fields.first(where: { $0.key == key.stringValue }) else {
            throw JSONCBinderError.missingRequiredField(key.stringValue)
        }
        return JSONCDecoder(node: field.value)
    }
}

private struct JSONCUnkeyedDecodingContainer: UnkeyedDecodingContainer {
    let node: JSONCArray
    let codingPath: [CodingKey]

    var count: Int? {
        return node.items.count
    }

    var isAtEnd: Bool {
        return currentIndex >= node.items.count
    }

    var currentIndex: Int = 0

    init(node: JSONCArray, codingPath: [CodingKey]) {
        self.node = node
        self.codingPath = codingPath
    }

    mutating func decode(_ type: Bool.Type) throws -> Bool {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let boolVal = valueNode.data as? Bool else {
            throw JSONCBinderError.invalidValue("Bool at index \(currentIndex)")
        }
        currentIndex += 1
        return boolVal
    }

    mutating func decode(_ type: Int.Type) throws -> Int {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int at index \(currentIndex)")
        }
        currentIndex += 1
        return intVal
    }

    mutating func decode(_ type: Int8.Type) throws -> Int8 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int8 at index \(currentIndex)")
        }
        currentIndex += 1
        return Int8(intVal)
    }

    mutating func decode(_ type: Int16.Type) throws -> Int16 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int16 at index \(currentIndex)")
        }
        currentIndex += 1
        return Int16(intVal)
    }

    mutating func decode(_ type: Int32.Type) throws -> Int32 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            throw JSONCBinderError.invalidValue("Int32 at index \(currentIndex)")
        }
        currentIndex += 1
        return Int32(intVal)
    }

    mutating func decode(_ type: Int64.Type) throws -> Int64 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int64 else {
            throw JSONCBinderError.invalidValue("Int64 at index \(currentIndex)")
        }
        currentIndex += 1
        return intVal
    }

    mutating func decode(_ type: UInt.Type) throws -> UInt {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt at index \(currentIndex)")
        }
        currentIndex += 1
        return uintVal
    }

    mutating func decode(_ type: UInt8.Type) throws -> UInt8 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt8 at index \(currentIndex)")
        }
        currentIndex += 1
        return UInt8(uintVal)
    }

    mutating func decode(_ type: UInt16.Type) throws -> UInt16 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt16 at index \(currentIndex)")
        }
        currentIndex += 1
        return UInt16(uintVal)
    }

    mutating func decode(_ type: UInt32.Type) throws -> UInt32 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let uintVal = valueNode.data as? UInt else {
            throw JSONCBinderError.invalidValue("UInt32 at index \(currentIndex)")
        }
        currentIndex += 1
        return UInt32(uintVal)
    }

    mutating func decode(_ type: UInt64.Type) throws -> UInt64 {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let uintVal = valueNode.data as? UInt64 else {
            throw JSONCBinderError.invalidValue("UInt64 at index \(currentIndex)")
        }
        currentIndex += 1
        return uintVal
    }

    mutating func decode(_ type: Float.Type) throws -> Float {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let doubleVal = valueNode.data as? Double else {
            throw JSONCBinderError.invalidValue("Float at index \(currentIndex)")
        }
        currentIndex += 1
        return Float(doubleVal)
    }

    mutating func decode(_ type: Double.Type) throws -> Double {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let doubleVal = valueNode.data as? Double else {
            throw JSONCBinderError.invalidValue("Double at index \(currentIndex)")
        }
        currentIndex += 1
        return doubleVal
    }

    mutating func decode(_ type: String.Type) throws -> String {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let stringVal = valueNode.data as? String else {
            throw JSONCBinderError.invalidValue("String at index \(currentIndex)")
        }
        currentIndex += 1
        return stringVal
    }

    mutating func decode(_ type: Data.Type) throws -> Data {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let dataVal = valueNode.data as? Data else {
            throw JSONCBinderError.invalidValue("Data at index \(currentIndex)")
        }
        currentIndex += 1
        return dataVal
    }

    mutating func decode<T>(_ type: T.Type) throws -> T where T: Decodable {
        guard currentIndex < node.items.count else {
            throw JSONCBinderError.invalidValue("Value at index \(currentIndex)")
        }

        let decoder = JSONCDecoder(node: node.items[currentIndex])
        currentIndex += 1
        return try T(from: decoder)
    }

    mutating func decodeIfPresent(_ type: Bool.Type) throws -> Bool? {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let boolVal = valueNode.data as? Bool else {
            return nil
        }
        currentIndex += 1
        return boolVal
    }

    mutating func decodeIfPresent(_ type: Int.Type) throws -> Int? {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let intVal = valueNode.data as? Int else {
            return nil
        }
        currentIndex += 1
        return intVal
    }

    mutating func decodeIfPresent(_ type: String.Type) throws -> String? {
        guard currentIndex < node.items.count,
              let valueNode = node.items[currentIndex] as? JSONCValue,
              let stringVal = valueNode.data as? String else {
            return nil
        }
        currentIndex += 1
        return stringVal
    }

    mutating func decodeIfPresent<T>(_ type: T.Type) throws -> T? where T: Decodable {
        guard currentIndex < node.items.count else {
            return nil
        }

        let decoder = JSONCDecoder(node: node.items[currentIndex])
        currentIndex += 1
        return try T(from: decoder)
    }

    mutating func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type) throws -> KeyedDecodingContainer<NestedKey> where NestedKey: CodingKey {
        guard currentIndex < node.items.count,
              let obj = node.items[currentIndex] as? JSONCObject else {
            throw JSONCBinderError.invalidValue("Nested object at index \(currentIndex)")
        }
        currentIndex += 1
        return KeyedDecodingContainer(JSONCKeyedDecodingContainer<NestedKey>(node: obj, codingPath: codingPath))
    }

    mutating func nestedUnkeyedContainer() throws -> UnkeyedDecodingContainer {
        guard currentIndex < node.items.count,
              let arr = node.items[currentIndex] as? JSONCArray else {
            throw JSONCBinderError.invalidValue("Nested array at index \(currentIndex)")
        }
        currentIndex += 1
        return JSONCUnkeyedDecodingContainer(node: arr, codingPath: codingPath)
    }

    mutating func superDecoder() throws -> Decoder {
        guard currentIndex < node.items.count else {
            throw JSONCBinderError.invalidValue("Value at index \(currentIndex)")
        }
        let decoder = JSONCDecoder(node: node.items[currentIndex])
        currentIndex += 1
        return decoder
    }
}

private struct JSONCSingleValueDecodingContainer: SingleValueDecodingContainer {
    let node: JSONCValue
    let codingPath: [CodingKey]

    init(node: JSONCValue, codingPath: [CodingKey]) {
        self.node = node
        self.codingPath = codingPath
    }

    func decode(_ type: Bool.Type) throws -> Bool {
        guard let boolVal = node.data as? Bool else {
            throw JSONCBinderError.invalidValue("Bool")
        }
        return boolVal
    }

    func decode(_ type: Int.Type) throws -> Int {
        if let intVal = node.data as? Int {
            return intVal
        }
        if let int64Val = node.data as? Int64 {
            return Int(int64Val)
        }
        throw JSONCBinderError.invalidValue("Int")
    }

    func decode(_ type: Int8.Type) throws -> Int8 {
        if let intVal = node.data as? Int {
            return Int8(intVal)
        }
        throw JSONCBinderError.invalidValue("Int8")
    }

    func decode(_ type: Int16.Type) throws -> Int16 {
        if let intVal = node.data as? Int {
            return Int16(intVal)
        }
        throw JSONCBinderError.invalidValue("Int16")
    }

    func decode(_ type: Int32.Type) throws -> Int32 {
        if let intVal = node.data as? Int {
            return Int32(intVal)
        }
        throw JSONCBinderError.invalidValue("Int32")
    }

    func decode(_ type: Int64.Type) throws -> Int64 {
        if let intVal = node.data as? Int64 {
            return intVal
        }
        if let intVal = node.data as? Int {
            return Int64(intVal)
        }
        throw JSONCBinderError.invalidValue("Int64")
    }

    func decode(_ type: UInt.Type) throws -> UInt {
        if let uintVal = node.data as? UInt {
            return uintVal
        }
        throw JSONCBinderError.invalidValue("UInt")
    }

    func decode(_ type: UInt8.Type) throws -> UInt8 {
        if let uintVal = node.data as? UInt {
            return UInt8(uintVal)
        }
        throw JSONCBinderError.invalidValue("UInt8")
    }

    func decode(_ type: UInt16.Type) throws -> UInt16 {
        if let uintVal = node.data as? UInt {
            return UInt16(uintVal)
        }
        throw JSONCBinderError.invalidValue("UInt16")
    }

    func decode(_ type: UInt32.Type) throws -> UInt32 {
        if let uintVal = node.data as? UInt {
            return UInt32(uintVal)
        }
        throw JSONCBinderError.invalidValue("UInt32")
    }

    func decode(_ type: UInt64.Type) throws -> UInt64 {
        if let uintVal = node.data as? UInt64 {
            return uintVal
        }
        throw JSONCBinderError.invalidValue("UInt64")
    }

    func decode(_ type: Float.Type) throws -> Float {
        if let doubleVal = node.data as? Double {
            return Float(doubleVal)
        }
        throw JSONCBinderError.invalidValue("Float")
    }

    func decode(_ type: Double.Type) throws -> Double {
        if let doubleVal = node.data as? Double {
            return doubleVal
        }
        throw JSONCBinderError.invalidValue("Double")
    }

    func decode(_ type: String.Type) throws -> String {
        if let stringVal = node.data as? String {
            return stringVal
        }
        throw JSONCBinderError.invalidValue("String")
    }

    func decode(_ type: Data.Type) throws -> Data {
        if let dataVal = node.data as? Data {
            return dataVal
        }
        throw JSONCBinderError.invalidValue("Data")
    }

    func decode<T>(_ type: T.Type) throws -> T where T: Decodable {
        let decoder = JSONCDecoder(node: node)
        return try T(from: decoder)
    }

    func decodeNil() -> Bool {
        return node.data == nil || node.tag?.isNull == true
    }
}