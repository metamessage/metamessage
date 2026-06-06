import Foundation

public enum ValueType: UInt8, Codable {
    case unknown = 0

    case doc = 1
    case vec = 2
    case arr = 3
    case obj = 4
    case map = 5

    case str = 6
    case bytes = 7
    case bool = 8

    case i = 9
    case i8 = 10
    case i16 = 11
    case i32 = 12
    case i64 = 13
    case u = 14
    case u8 = 15
    case u16 = 16
    case u32 = 17
    case u64 = 18

    case f32 = 19
    case f64 = 20

    case bigint = 21
    case datetime = 22
    case date = 23
    case time = 24

    case uuid = 25
    case decimal = 26
    case ip = 27
    case url = 28
    case email = 29

    case enums = 30
    case media = 31

    public var stringValue: String {
        switch self {
        case .unknown: return "unknown"
        case .doc: return "doc"
        case .vec: return "vec"
        case .arr: return "arr"
        case .obj: return "obj"
        case .map: return "map"
        case .str: return "str"
        case .bytes: return "bytes"
        case .bool: return "bool"
        case .i: return "i"
        case .i8: return "i8"
        case .i16: return "i16"
        case .i32: return "i32"
        case .i64: return "i64"
        case .u: return "u"
        case .u8: return "u8"
        case .u16: return "u16"
        case .u32: return "u32"
        case .u64: return "u64"
        case .f32: return "f32"
        case .f64: return "f64"
        case .bigint: return "bigint"
        case .datetime: return "datetime"
        case .date: return "date"
        case .time: return "time"
        case .uuid: return "uuid"
        case .decimal: return "decimal"
        case .ip: return "ip"
        case .url: return "url"
        case .email: return "email"
        case .enums: return "enums"
        case .media: return "media"
        }
    }

    public static func parse(_ s: String) -> ValueType? {
        switch s.lowercased() {
        case "unknown": return .unknown
        case "doc": return .doc
        case "vec": return .vec
        case "arr", "array": return .arr
        case "obj": return .obj
        case "map": return .map
        case "str": return .str
        case "bytes": return .bytes
        case "bool": return .bool
        case "i": return .i
        case "i8": return .i8
        case "i16": return .i16
        case "i32": return .i32
        case "i64": return .i64
        case "u": return .u
        case "u8": return .u8
        case "u16": return .u16
        case "u32": return .u32
        case "u64": return .u64
        case "f32": return .f32
        case "f64": return .f64
        case "bigint": return .bigint
        case "datetime": return .datetime
        case "date": return .date
        case "time": return .time
        case "uuid": return .uuid
        case "decimal": return .decimal
        case "ip": return .ip
        case "url": return .url
        case "email": return .email
        case "enums": return .enums
        case "media": return .media
        default: return nil
        }
    }

    public var needsQuotes: Bool {
        switch self {
        case .str, .bytes, .datetime, .date, .time, .uuid, .ip, .url, .email, .enums, .unknown:
            return true
        default:
            return false
        }
    }
}