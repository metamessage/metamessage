import Foundation

public class Decoder {
    private let buffer: MMBuffer

    public init(data: Data) {
        self.buffer = MMBuffer(data: data)
    }

    public enum DecodedValue {
        case bool(Bool)
        case int(Int64)
        case uint(UInt64)
        case float(Double)
        case string(String)
        case data(Data)
        case array([DecodedValue])
        case object([String: DecodedValue])
        case null
    }

    public func decode() throws -> DecodedValue {
        guard let byte = buffer.peek() else {
            throw MMError.unexpectedEndOfData
        }

        guard let prefix = getPrefix(byte) else {
            throw MMError.invalidPrefix
        }

        switch prefix {
        case .simple:
            return try decodeSimple()
        case .positiveInt:
            return try decodePositiveInt()
        case .negativeInt:
            return try decodeNegativeInt()
        case .prefixFloat:
            return try decodeFloat()
        case .prefixString:
            return try decodeString()
        case .prefixBytes:
            return try decodeBytes()
        case .container:
            return try decodeContainer()
        case .prefixTag:
            return try decodeTag()
        }
    }

    private func decodeSimple() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        guard let simpleValue = MMSimpleValue(rawValue: byte) else {
            throw MMError.invalidData
        }

        switch simpleValue {
        case .trueValue:
            return .bool(true)
        case .falseValue:
            return .bool(false)
        case .nullBool:
            return .bool(false)
        case .nullInt:
            return .int(0)
        case .nullFloat:
            return .float(0.0)
        case .nullString:
            return .string("")
        case .nullBytes:
            return .data(Data())
        default:
            return .null
        }
    }

    private func decodePositiveInt() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for i in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        if value > UInt64(Int64.max) {
            return .uint(value)
        }
        return .int(Int64(value))
    }

    private func decodeNegativeInt() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for i in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        if value > UInt64(Int64.max) {
            return .int(Int64.min)
        }
        return .int(-Int64(value))
    }

    private func decodeFloat() throws -> DecodedValue {
        guard let prefix = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let prefixLower = prefix & 0x0F

        if prefixLower <= 6 {
            let val = Double(prefixLower) / 10.0
            if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
                return .float(-val)
            }
            return .float(val)
        }

        guard let expByte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }
        let exponent = Int8(bitPattern: expByte)

        let (extraBytes, _) = floatLen(prefix)
        var mantissa: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            mantissa = (mantissa << 8) | UInt64(b)
        }

        let decimalStr = mantissaToDecimal(mantissa, exponent)
        guard let value = Double(decimalStr) else {
            throw MMError.invalidData
        }

        if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
            return .float(-value)
        }
        return .float(value)
    }

    private func mantissaToDecimal(_ mantissa: UInt64, _ exp: Int8) -> String {
        let numStr = String(mantissa)
        let decimalPos = numStr.count + Int(exp)

        if decimalPos <= 0 {
            return "0." + String(repeating: "0", count: -decimalPos) + numStr
        } else if decimalPos > 0 && decimalPos < numStr.count {
            let idx = numStr.index(numStr.startIndex, offsetBy: decimalPos)
            return String(numStr[..<idx]) + "." + String(numStr[idx...])
        } else {
            let trailingZeros = decimalPos - numStr.count
            return numStr + String(repeating: "0", count: trailingZeros)
        }
    }

    private func decodeString() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = stringLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let bytes = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        guard let string = String(bytes: bytes, encoding: .utf8) else {
            throw MMError.invalidData
        }

        return .string(string)
    }

    private func decodeBytes() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = bytesLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let result = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        return .data(Data(result))
    }

    private func decodeContainer() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let isArrayContainer = isArray(byte)

        let (extraBytes, len) = containerLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        if isArrayContainer {
            var elements: [DecodedValue] = []
            let startPos = buffer.position()
            let endPos = startPos + totalLen
            while buffer.position() < endPos {
                let element = try decode()
                elements.append(element)
            }
            return .array(elements)
        } else {
            var dict: [String: DecodedValue] = [:]

            let keyArrayValue = try decode()
            guard case .array(let keyItems) = keyArrayValue else {
                throw MMError.invalidData
            }

            for item in keyItems {
                guard case .string(let key) = item else {
                    continue
                }
                let value = try decode()
                dict[key] = value
            }

            return .object(dict)
        }
    }

    private func decodeTag() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = tagLen(byte)
        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let data = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        var pos = 0
        guard pos < data.count else {
            throw MMError.unexpectedEndOfData
        }

        let tagBodyLenByte = data[pos]
        pos += 1
        var tagBodyLen = Int(tagBodyLenByte)
        if tagBodyLenByte >= 254 {
            if tagBodyLenByte == 254 {
                guard pos < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = Int(data[pos])
                pos += 1
            } else {
                guard pos + 1 < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                pos += 2
            }
        }

        let tag = Tag()
        let tagFieldEnd = pos + tagBodyLen
        while pos < tagFieldEnd {
            guard pos < data.count else { break }
            let fieldByte = data[pos]
            let key = fieldByte & 0xF8
            let fieldLen = Int(fieldByte) & 0x07
            pos += 1

            switch key {
            case TagKey.type:
                guard pos < data.count else { break }
                tag.type = ValueType(rawValue: data[pos]) ?? .unknown
                pos += 1

            case TagKey.isNull:
                tag.isNull = (fieldLen & 0x01) == 1
                if tag.isNull {
                    tag.nullable = true
                }

            case TagKey.deprecated:
                tag.deprecated = (fieldLen & 0x01) == 1

            case TagKey.nullable:
                tag.nullable = (fieldLen & 0x01) == 1

            case TagKey.allowEmpty:
                tag.allowEmpty = (fieldLen & 0x01) == 1

            case TagKey.unique:
                tag.unique = (fieldLen & 0x01) == 1

            case TagKey.example:
                break

            default:
                let isU64Field = key == TagKey.size || key == TagKey.version ||
                                 key == TagKey.more || key == TagKey.childSize ||
                                 key == TagKey.childVersion
                if isU64Field {
                    pos += fieldLen + 1
                } else if fieldLen <= 5 {
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1 + strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2 + strLen
                }
            }
        }

        pos = tagFieldEnd

        if tag.isNull {
            if tag.type != .unknown {
                return tagNullValue(tag.type)
            }
        }

        guard pos < data.count else {
            if tag.isNull {
                return .null
            }
            throw MMError.unexpectedEndOfData
        }

        let payloadData = Data(data[pos...])
        let innerDecoder = Decoder(data: payloadData)
        let value = try innerDecoder.decode()

        if tag.type != .unknown {
            return applyTagType(value, tag: tag)
        }
        return value
    }

    private func tagNullValue(_ type: ValueType) -> DecodedValue {
        switch type {
        case .i, .i8, .i16, .i32, .i64:
            return .int(0)
        case .u, .u8, .u16, .u32, .u64:
            return .uint(0)
        case .f32, .f64:
            return .float(0.0)
        case .bool:
            return .bool(false)
        case .str, .email, .url, .datetime, .date, .time, .enums:
            return .string("")
        case .bytes, .uuid, .ip:
            return .data(Data())
        default:
            return .null
        }
    }

    private func applyTagType(_ value: DecodedValue, tag: Tag) -> DecodedValue {
        switch tag.type {
        case .datetime:
            if case .int(let ts) = value {
                let date = Date(timeIntervalSince1970: TimeInterval(ts))
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
            if case .uint(let ts) = value {
                let date = Date(timeIntervalSince1970: TimeInterval(ts))
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }

        case .date:
            if case .int(let days) = value {
                let ref = Date(timeIntervalSince1970: 0)
                let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: Int(days), to: ref) ?? ref
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
            if case .uint(let days) = value {
                let ref = Date(timeIntervalSince1970: 0)
                let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: Int(days), to: ref) ?? ref
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }

        case .time:
            if case .int(let secs) = value {
                let hours = secs / 3600
                let mins = (secs % 3600) / 60
                let sec = secs % 60
                return .string(String(format: "%02d:%02d:%02d", hours, mins, sec))
            }
            if case .uint(let secs) = value {
                let hours = secs / 3600
                let mins = (secs % 3600) / 60
                let sec = secs % 60
                return .string(String(format: "%02d:%02d:%02d", hours, mins, sec))
            }

        case .uuid:
            if case .data(let d) = value, d.count == 16 {
                let bytes = [UInt8](d)
                let hexStr = String(format: "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
                    bytes[0], bytes[1], bytes[2], bytes[3],
                    bytes[4], bytes[5], bytes[6], bytes[7],
                    bytes[8], bytes[9], bytes[10], bytes[11],
                    bytes[12], bytes[13], bytes[14], bytes[15])
                return .string(hexStr)
            }

        case .enums:
            if !tag.enums.isEmpty {
                let enumValues = tag.enums.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
                if case .int(let idx) = value, idx >= 0, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
                if case .uint(let idx) = value, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
            }

        default:
            break
        }

        return value
    }
}