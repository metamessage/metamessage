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
            value |= UInt64(b) << (i * 8)
        }

        let suffix = Int(byte & 0x1F)
        value |= UInt64(suffix)

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
            value |= UInt64(b) << (i * 8)
        }

        let suffix = Int(byte & 0x1F)
        value |= UInt64(suffix)

        return .int(-Int64(value))
    }

    private func decodeFloat() throws -> DecodedValue {
        _ = buffer.read()
        let value = try buffer.readFloat64()
        return .float(value)
    }

    private func decodeString() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = stringLen(byte)

        var totalLen = len
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            for (i, b) in bytes.enumerated() {
                totalLen |= Int(b) << (i * 8)
            }
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

        var totalLen = len
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            for (i, b) in bytes.enumerated() {
                totalLen |= Int(b) << (i * 8)
            }
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

        var totalLen = len
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            for (i, b) in bytes.enumerated() {
                totalLen |= Int(b) << (i * 8)
            }
        }

        if isArrayContainer {
            var elements: [DecodedValue] = []
            for _ in 0..<totalLen {
                let element = try decode()
                elements.append(element)
            }
            return .array(elements)
        } else {
            var dict: [String: DecodedValue] = [:]
            for _ in 0..<(totalLen / 2) {
                let keyValue = try decode()
                let value = try decode()
                if case .string(let key) = keyValue {
                    dict[key] = value
                }
            }
            return .object(dict)
        }
    }

    private func decodeTag() throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = tagLen(byte)

        if extraBytes > 0 {
            guard buffer.read(extraBytes) != nil else {
                throw MMError.unexpectedEndOfData
            }
        }

        return try decode()
    }
}