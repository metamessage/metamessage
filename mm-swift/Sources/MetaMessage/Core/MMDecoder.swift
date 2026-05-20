import Foundation

public class MMDecoder {
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
        case .nullBool, .nullInt, .nullFloat, .nullString, .nullBytes:
            return .null
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

        let tagFieldEnd = pos + tagBodyLen
        while pos < tagFieldEnd {
            guard pos < data.count else { break }
            let fieldByte = data[pos]
            let key = Int(fieldByte) & 0xF8
            let fieldLen = Int(fieldByte) & 0x07
            pos += 1

            if key == TagKey.type {
                guard pos < data.count else { break }
                pos += 1
                continue
            }

            let isBooleanFlag = key == TagKey.isNull || key == TagKey.raw ||
                                key == TagKey.nullable || key == TagKey.allowEmpty ||
                                key == TagKey.unique
            if isBooleanFlag {
                continue
            }

            if fieldLen <= 5 {
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

        pos = tagFieldEnd
        guard pos < data.count else {
            throw MMError.unexpectedEndOfData
        }

        let payloadData = Data(data[pos...])
        let innerDecoder = MMDecoder(data: payloadData)
        return try innerDecoder.decode()
    }
}