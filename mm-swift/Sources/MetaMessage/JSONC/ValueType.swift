import Foundation

public enum ValueType: UInt8, Codable {
    case unknown = 0

    case doc = 1
    case slice = 2
    case array = 3
    case structType = 4
    case map = 5

    case string = 6
    case bytes = 7
    case bool = 8

    case int = 9
    case int8 = 10
    case int16 = 11
    case int32 = 12
    case int64 = 13
    case uint = 14
    case uint8 = 15
    case uint16 = 16
    case uint32 = 17
    case uint64 = 18

    case float32 = 19
    case float64 = 20

    case bigInt = 21
    case dateTime = 22
    case date = 23
    case time = 24

    case uuid = 25
    case decimal = 26
    case ip = 27
    case url = 28
    case email = 29

    case enumValue = 30

    case image = 31
    case video = 32

    public var stringValue: String {
        switch self {
        case .unknown: return "unknown"
        case .doc: return "doc"
        case .slice: return "slice"
        case .array: return "arr"
        case .structType: return "struct"
        case .map: return "map"
        case .string: return "str"
        case .bytes: return "bytes"
        case .bool: return "bool"
        case .int: return "i"
        case .int8: return "i8"
        case .int16: return "i16"
        case .int32: return "i32"
        case .int64: return "i64"
        case .uint: return "u"
        case .uint8: return "u8"
        case .uint16: return "u16"
        case .uint32: return "u32"
        case .uint64: return "u64"
        case .float32: return "f32"
        case .float64: return "f64"
        case .bigInt: return "bi"
        case .dateTime: return "datetime"
        case .date: return "date"
        case .time: return "time"
        case .uuid: return "uuid"
        case .decimal: return "decimal"
        case .ip: return "ip"
        case .url: return "url"
        case .email: return "email"
        case .enumValue: return "enum"
        case .image: return "image"
        case .video: return "video"
        }
    }

    public static func parse(_ s: String) -> ValueType? {
        switch s.lowercased() {
        case "unknown": return .unknown
        case "doc": return .doc
        case "slice": return .slice
        case "arr", "array": return .array
        case "struct": return .structType
        case "map": return .map
        case "str", "string": return .string
        case "bytes": return .bytes
        case "bool": return .bool
        case "i", "int": return .int
        case "i8", "int8": return .int8
        case "i16", "int16": return .int16
        case "i32", "int32": return .int32
        case "i64", "int64": return .int64
        case "u", "uint": return .uint
        case "u8", "uint8": return .uint8
        case "u16", "uint16": return .uint16
        case "u32", "uint32": return .uint32
        case "u64", "uint64": return .uint64
        case "f32", "float32": return .float32
        case "f64", "float64", "float": return .float64
        case "bi", "bigint": return .bigInt
        case "datetime": return .dateTime
        case "date": return .date
        case "time": return .time
        case "uuid": return .uuid
        case "decimal": return .decimal
        case "ip": return .ip
        case "url": return .url
        case "email": return .email
        case "enum": return .enumValue
        case "image": return .image
        case "video": return .video
        default: return nil
        }
    }

    public var needsQuotes: Bool {
        switch self {
        case .string, .bytes, .dateTime, .date, .time, .uuid, .ip, .url, .email, .enumValue:
            return true
        default:
            return false
        }
    }
}