import Foundation

// MARK: - Internal Encoder

class MetaMessageEncoder {
    private var encoder: Encoder

    init() {
        self.encoder = Encoder()
    }

    init(capacity: Int) {
        self.encoder = Encoder(capacity: capacity)
    }

    func encode(_ value: Bool) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Int8) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Int16) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Int32) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Int64) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Int) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: UInt8) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: UInt16) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: UInt32) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: UInt64) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: UInt) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Float) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Double) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: String) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encode(_ value: Data) -> Data {
        encoder.encode(value)
        return encoder.result
    }

    func encodeNil() -> Data {
        encoder.encodeNil()
        return encoder.result
    }

    func encodeArray(_ array: [Bool]) -> Data {
        encoder.encodeArray(array)
        return encoder.result
    }

    func encodeArrayStrings(_ array: [String]) -> Data {
        encoder.encodeArrayStrings(array)
        return encoder.result
    }

    func encodeArrayInt(_ array: [Int]) -> Data {
        encoder.encodeArrayInt(array)
        return encoder.result
    }

    func encodeArrayUInt(_ array: [UInt]) -> Data {
        encoder.encodeArrayUInt(array)
        return encoder.result
    }

    func encodeArrayFloat(_ array: [Float]) -> Data {
        encoder.encodeArrayFloat(array)
        return encoder.result
    }

    func encodeArrayDouble(_ array: [Double]) -> Data {
        encoder.encodeArrayDouble(array)
        return encoder.result
    }

    func encodeArrayData(_ array: [Data]) -> Data {
        encoder.encodeArrayData(array)
        return encoder.result
    }
}

// MARK: - Internal Decoder

class MetaMessageDecoder {
    private var decoder: Decoder

    init(data: Data) {
        self.decoder = Decoder(data: data)
    }

    func decode() throws -> Decoder.DecodedValue {
        return try decoder.decode()
    }

    func decodeToBool() throws -> Bool {
        let value = try decode()
        guard case .bool(let b) = value else {
            throw MMError.typeMismatch
        }
        return b
    }

    func decodeToInt() throws -> Int {
        let value = try decode()
        guard case .int(let i) = value else {
            throw MMError.typeMismatch
        }
        return Int(i)
    }

    func decodeToInt64() throws -> Int64 {
        let value = try decode()
        guard case .int(let i) = value else {
            throw MMError.typeMismatch
        }
        return i
    }

    func decodeToUInt() throws -> UInt {
        let value = try decode()
        guard case .uint(let u) = value else {
            throw MMError.typeMismatch
        }
        return UInt(u)
    }

    func decodeToUInt64() throws -> UInt64 {
        let value = try decode()
        guard case .uint(let u) = value else {
            throw MMError.typeMismatch
        }
        return u
    }

    func decodeToFloat() throws -> Float {
        let value = try decode()
        guard case .float(let f) = value else {
            throw MMError.typeMismatch
        }
        return Float(f)
    }

    func decodeToDouble() throws -> Double {
        let value = try decode()
        guard case .float(let f) = value else {
            throw MMError.typeMismatch
        }
        return f
    }

    func decodeToString() throws -> String {
        let value = try decode()
        guard case .string(let s) = value else {
            throw MMError.typeMismatch
        }
        return s
    }

    func decodeToData() throws -> Data {
        let value = try decode()
        guard case .data(let d) = value else {
            throw MMError.typeMismatch
        }
        return d
    }

    func decodeToArray() throws -> [Decoder.DecodedValue] {
        let value = try decode()
        guard case .array(let arr) = value else {
            throw MMError.typeMismatch
        }
        return arr
    }

    func decodeToObject() throws -> [String: Decoder.DecodedValue] {
        let value = try decode()
        guard case .object(let obj) = value else {
            throw MMError.typeMismatch
        }
        return obj
    }
}

// MARK: - Private Helpers

func encodeNode(_ node: Node, with encoder: Encoder) throws -> Data {
    switch node {
    case let obj as NodeObject:
        encoder.encodeNodeObject(obj)
    case let arr as NodeArray:
        encoder.encodeNodeArray(arr)
    case let val as NodeScalar:
        encoder.encodeNodeValue(val)
    case let nullNode as NodeNull:
        encoder.encodeNodeNull(nullNode)
    default:
        throw MMError.unsupportedType
    }
    return encoder.result
}

func nodeToString(_ node: Decoder.DecodedValue) -> String {
    switch node {
    case .bool(let b):
        return b ? "true" : "false"
    case .int(let i):
        return String(i)
    case .uint(let u):
        return String(u)
    case .float(let f):
        if f == f.rounded() && f != 0 {
            return String(format: "%.1f", f)
        }
        return String(f)
    case .string(let s):
        let escaped = s
            .replacingOccurrences(of: "\\", with: "\\\\")
            .replacingOccurrences(of: "\"", with: "\\\"")
            .replacingOccurrences(of: "\n", with: "\\n")
            .replacingOccurrences(of: "\r", with: "\\r")
            .replacingOccurrences(of: "\t", with: "\\t")
        return "\"\(escaped)\""
    case .bigint(let s):
        return s
    case .data(let d):
        let str = String(data: d, encoding: .utf8) ?? d.map { String(format: "\\x%02x", $0) }.joined()
        let escaped = str
            .replacingOccurrences(of: "\\", with: "\\\\")
            .replacingOccurrences(of: "\"", with: "\\\"")
            .replacingOccurrences(of: "\n", with: "\\n")
            .replacingOccurrences(of: "\r", with: "\\r")
            .replacingOccurrences(of: "\t", with: "\\t")
        return "\"\(escaped)\""
    case .array(let arr):
        let items = arr.map { nodeToString($0) }
        return "[" + items.joined(separator: ",") + "]"
    case .object(let obj):
        let sortedKeys = obj.keys.sorted()
        let items = sortedKeys.map { "\"\($0)\":\(nodeToString(obj[$0]!))" }
        return "{" + items.joined(separator: ",") + "}"
    case .null:
        return "null"
    }
}