import Foundation

public class MetaMessageEncoder {
    private var encoder: MMEncoder

    public init() {
        self.encoder = MMEncoder()
    }

    public init(capacity: Int) {
        self.encoder = MMEncoder(capacity: capacity)
    }

    public func encode(_ value: Bool) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Int8) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Int16) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Int32) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Int64) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Int) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: UInt8) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: UInt16) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: UInt32) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: UInt64) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: UInt) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Float) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Double) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: String) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encode(_ value: Data) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    public func encodeNil() -> Data {
        encoder.encodeNil()
        return encoder.result
    }

    public func encodeArray(_ array: [Bool]) -> Data {
        encoder.encodeArray(array)
        return encoder.result
    }

    public func encodeArrayStrings(_ array: [String]) -> Data {
        encoder.encodeArrayStrings(array)
        return encoder.result
    }

    public func encodeArrayInt(_ array: [Int]) -> Data {
        encoder.encodeArrayInt(array)
        return encoder.result
    }

    public func encodeArrayUInt(_ array: [UInt]) -> Data {
        encoder.encodeArrayUInt(array)
        return encoder.result
    }

    public func encodeArrayFloat(_ array: [Float]) -> Data {
        encoder.encodeArrayFloat(array)
        return encoder.result
    }

    public func encodeArrayDouble(_ array: [Double]) -> Data {
        encoder.encodeArrayDouble(array)
        return encoder.result
    }

    public func encodeArrayData(_ array: [Data]) -> Data {
        encoder.encodeArrayData(array)
        return encoder.result
    }
}

public class MetaMessageDecoder {
    private var decoder: MMDecoder

    public init(data: Data) {
        self.decoder = MMDecoder(data: data)
    }

    public func decode() throws -> MMDecoder.DecodedValue {
        return try decoder.decode()
    }

    public func decodeToBool() throws -> Bool {
        let value = try decode()
        guard case .bool(let b) = value else {
            throw MMError.typeMismatch
        }
        return b
    }

    public func decodeToInt() throws -> Int {
        let value = try decode()
        guard case .int(let i) = value else {
            throw MMError.typeMismatch
        }
        return Int(i)
    }

    public func decodeToInt64() throws -> Int64 {
        let value = try decode()
        guard case .int(let i) = value else {
            throw MMError.typeMismatch
        }
        return i
    }

    public func decodeToUInt() throws -> UInt {
        let value = try decode()
        guard case .uint(let u) = value else {
            throw MMError.typeMismatch
        }
        return u
    }

    public func decodeToUInt64() throws -> UInt64 {
        let value = try decode()
        guard case .uint(let u) = value else {
            throw MMError.typeMismatch
        }
        return u
    }

    public func decodeToFloat() throws -> Float {
        let value = try decode()
        guard case .float(let f) = value else {
            throw MMError.typeMismatch
        }
        return Float(f)
    }

    public func decodeToDouble() throws -> Double {
        let value = try decode()
        guard case .float(let f) = value else {
            throw MMError.typeMismatch
        }
        return f
    }

    public func decodeToString() throws -> String {
        let value = try decode()
        guard case .string(let s) = value else {
            throw MMError.typeMismatch
        }
        return s
    }

    public func decodeToData() throws -> Data {
        let value = try decode()
        guard case .data(let d) = value else {
            throw MMError.typeMismatch
        }
        return d
    }

    public func decodeToArray() throws -> [MMDecoder.DecodedValue] {
        let value = try decode()
        guard case .array(let arr) = value else {
            throw MMError.typeMismatch
        }
        return arr
    }

    public func decodeToObject() throws -> [String: MMDecoder.DecodedValue] {
        let value = try decode()
        guard case .object(let obj) = value else {
            throw MMError.typeMismatch
        }
        return obj
    }
}

public enum MetaMessage {
    public static func encode(_ value: Bool) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Int) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Int8) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Int16) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Int32) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Int64) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: UInt) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: UInt8) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: UInt16) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: UInt32) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: UInt64) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Float) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Double) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: String) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func encode(_ value: Data) -> Data {
        return MetaMessageEncoder().encode(value)
    }

    public static func decode(_ data: Data) throws -> MMDecoder.DecodedValue {
        return try MetaMessageDecoder(data: data).decode()
    }

    public static func toJSONC(_ data: Data) throws -> String {
        let node = try decode(data)
        return nodeToString(node)
    }

    private static func nodeToString(_ node: MMDecoder.DecodedValue) -> String {
        let printer = JSONCPrinter()
        switch node {
        case .bool(let b):
            return b ? "true" : "false"
        case .int(let i):
            return String(i)
        case .uint(let u):
            return String(u)
        case .float(let f):
            return String(f)
        case .string(let s):
            return "\"\(s)\""
        case .data(let d):
            return "\"\(String(data: d, encoding: .utf8) ?? "")\""
        case .array(let arr):
            let items = arr.map { nodeToString($0) }
            return "[" + items.joined(separator: ", ") + "]"
        case .object(let obj):
            let items = obj.map { "\"\($0.key)\": \(nodeToString($0.value))" }
            return "{" + items.joined(separator: ", ") + "}"
        case .null:
            return "null"
        }
    }
}

public func structToJSONC<T>(_ value: T, tag: String = "") throws -> JSONCNode? {
    let mirror = Mirror(reflecting: value)
    return reflectToNode(mirror, tag: tag)
}

private func reflectToNode(_ mirror: Mirror, tag: String) -> JSONCNode? {
    if mirror.displayStyle == .class || mirror.displayStyle == .struct {
        var fields: [JSONCField] = []
        var nodeTag: JSONCTag? = nil

        if !tag.isEmpty {
            nodeTag = parseMMTag(tag)
        }

        for child in mirror.children {
            guard let label = child.label else { continue }

            let snakeKey = camelToSnake(label)

            var childTag: JSONCTag? = nil
            if let tagString = getTagString(from: child) {
                childTag = parseMMTag(tagString)
            }

            if let childNode = reflectValueToNode(child.value) {
                let field = JSONCField(key: snakeKey, value: childNode)
                fields.append(field)
            }
        }

        return JSONCObject(fields: fields, tag: nodeTag)
    } else if mirror.displayStyle == .array || mirror.displayStyle == .collection {
        var items: [JSONCNode] = []

        for child in mirror.children {
            if let childNode = reflectValueToNode(child.value) {
                items.append(childNode)
            }
        }

        return JSONCArray(items: items)
    } else if mirror.displayStyle == .dictionary {
        return nil
    }

    return nil
}

private func reflectValueToNode(_ value: Any) -> JSONCNode? {
    let mirror = Mirror(reflecting: value)

    if let boolVal = value as? Bool {
        return JSONCValue(data: boolVal, text: boolVal ? "true" : "false")
    }

    if let intVal = value as? Int {
        return JSONCValue(data: intVal, text: String(intVal))
    }

    if let intVal = value as? Int64 {
        return JSONCValue(data: intVal, text: String(intVal))
    }

    if let intVal = value as? Int32 {
        return JSONCValue(data: intVal, text: String(intVal))
    }

    if let intVal = value as? Int16 {
        return JSONCValue(data: intVal, text: String(intVal))
    }

    if let intVal = value as? Int8 {
        return JSONCValue(data: intVal, text: String(intVal))
    }

    if let uintVal = value as? UInt {
        return JSONCValue(data: uintVal, text: String(uintVal))
    }

    if let uintVal = value as? UInt64 {
        return JSONCValue(data: uintVal, text: String(uintVal))
    }

    if let uintVal = value as? UInt32 {
        return JSONCValue(data: uintVal, text: String(uintVal))
    }

    if let uintVal = value as? UInt16 {
        return JSONCValue(data: uintVal, text: String(uintVal))
    }

    if let uintVal = value as? UInt8 {
        return JSONCValue(data: uintVal, text: String(uintVal))
    }

    if let doubleVal = value as? Double {
        return JSONCValue(data: doubleVal, text: String(doubleVal))
    }

    if let floatVal = value as? Float {
        return JSONCValue(data: floatVal, text: String(floatVal))
    }

    if let stringVal = value as? String {
        return JSONCValue(data: stringVal, text: stringVal)
    }

    if let dataVal = value as? Data {
        return JSONCValue(data: dataVal, text: String(data: dataVal, encoding: .utf8) ?? "")
    }

    if mirror.displayStyle == .class || mirror.displayStyle == .struct {
        return reflectToNode(mirror, tag: "")
    }

    if mirror.displayStyle == .array || mirror.displayStyle == .collection {
        return reflectToNode(mirror, tag: "")
    }

    return JSONCValue(data: nil, text: "null")
}

private func getTagString(from child: Mirror.Child) -> String? {
    return nil
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