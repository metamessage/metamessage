import Foundation

public class MMEncoder {
    private let buffer: MMBuffer

    public init() {
        self.buffer = MMBuffer()
    }

    public init(capacity: Int) {
        self.buffer = MMBuffer(capacity: capacity)
    }

    public func encode(_ value: Bool) {
        if value {
            buffer.write(MMSimpleValue.trueValue.rawValue)
        } else {
            buffer.write(MMSimpleValue.falseValue.rawValue)
        }
    }

    public func encodeNil() {
        buffer.write(MMSimpleValue.nullInt.rawValue)
    }

    public func encode(_ value: Int8) {
        encodeInt64(Int64(value))
    }

    public func encode(_ value: Int16) {
        encodeInt64(Int64(value))
    }

    public func encode(_ value: Int32) {
        encodeInt64(Int64(value))
    }

    public func encode(_ value: Int64) {
        encodeInt64(value)
    }

    public func encode(_ value: Int) {
        encodeInt64(Int64(value))
    }

    public func encode(_ value: UInt8) {
        encodeUInt64(UInt64(value))
    }

    public func encode(_ value: UInt16) {
        encodeUInt64(UInt64(value))
    }

    public func encode(_ value: UInt32) {
        encodeUInt64(UInt64(value))
    }

    public func encode(_ value: UInt64) {
        encodeUInt64(value)
    }

    public func encode(_ value: UInt) {
        encodeUInt64(UInt64(value))
    }

    private func encodeInt64(_ value: Int64) {
        if value >= 0 {
            encodeUInt64(UInt64(value))
        } else {
            let uv: UInt64
            if value == Int64.min {
                uv = 9223372036854775808
            } else {
                uv = UInt64(-value - 1)
            }
            encodeUInt64WithSign(MMPrefix.negativeInt.rawValue, uv)
        }
    }

    private func encodeUInt64(_ value: UInt64) {
        encodeUInt64WithSign(MMPrefix.positiveInt.rawValue, value)
    }

    private func encodeUInt64WithSign(_ sign: UInt8, _ value: UInt64) {
        let (extraBytes, _) = intExtraBytes(for: value)

        switch extraBytes {
        case 0:
            buffer.write(sign | UInt8(value))
        case 1:
            buffer.write(sign | MMConstants.intLen1Byte)
            buffer.write(UInt8(value & 0xFF))
        case 2:
            buffer.write(sign | MMConstants.intLen2Byte)
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 3:
            buffer.write(sign | MMConstants.intLen3Byte)
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 4:
            buffer.write(sign | MMConstants.intLen4Byte)
            buffer.write(UInt8((value >> 24) & 0xFF))
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 5:
            buffer.write(sign | MMConstants.intLen5Byte)
            buffer.write(UInt8((value >> 32) & 0xFF))
            buffer.write(UInt8((value >> 24) & 0xFF))
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 6:
            buffer.write(sign | MMConstants.intLen6Byte)
            buffer.write(UInt8((value >> 40) & 0xFF))
            buffer.write(UInt8((value >> 32) & 0xFF))
            buffer.write(UInt8((value >> 24) & 0xFF))
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 7:
            buffer.write(sign | MMConstants.intLen7Byte)
            buffer.write(UInt8((value >> 48) & 0xFF))
            buffer.write(UInt8((value >> 40) & 0xFF))
            buffer.write(UInt8((value >> 32) & 0xFF))
            buffer.write(UInt8((value >> 24) & 0xFF))
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        case 8:
            buffer.write(sign | MMConstants.intLen8Byte)
            buffer.write(UInt8((value >> 56) & 0xFF))
            buffer.write(UInt8((value >> 48) & 0xFF))
            buffer.write(UInt8((value >> 40) & 0xFF))
            buffer.write(UInt8((value >> 32) & 0xFF))
            buffer.write(UInt8((value >> 24) & 0xFF))
            buffer.write(UInt8((value >> 16) & 0xFF))
            buffer.write(UInt8((value >> 8) & 0xFF))
            buffer.write(UInt8(value & 0xFF))
        default:
            break
        }
    }

    private func intExtraBytes(for value: UInt64) -> (Int, Int) {
        if value < UInt64(MMConstants.max1Byte) {
            return (0, Int(value))
        } else if value < UInt64(MMConstants.max2Byte) {
            return (1, 0)
        } else if value < UInt64(MMConstants.max3Byte) {
            return (2, 0)
        } else if value < UInt64(MMConstants.max4Byte) {
            return (3, 0)
        } else if value < UInt64(MMConstants.max5Byte) {
            return (4, 0)
        } else if value < UInt64(MMConstants.max6Byte) {
            return (5, 0)
        } else if value < UInt64(MMConstants.max7Byte) {
            return (6, 0)
        } else {
            return (7, 0)
        }
    }

    public func encode(_ value: Float) {
        buffer.write(MMPrefix.prefixFloat.rawValue)
        buffer.writeFloat32(value)
    }

    public func encode(_ value: Double) {
        buffer.write(MMPrefix.prefixFloat.rawValue)
        buffer.writeFloat64(value)
    }

    public func encode(_ value: String) {
        guard let data = value.data(using: .utf8) else { return }
        let bytes = [UInt8](data)
        let len = bytes.count

        if len < 254 {
            buffer.write(MMPrefix.prefixString.rawValue | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.prefixString.rawValue | MMConstants.stringLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.prefixString.rawValue | MMConstants.stringLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(bytes)
    }

    public func encode(_ value: Data) {
        let bytes = [UInt8](value)
        let len = bytes.count

        if len < 254 {
            buffer.write(MMPrefix.prefixBytes.rawValue | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.prefixBytes.rawValue | MMConstants.bytesLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.prefixBytes.rawValue | MMConstants.bytesLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(bytes)
    }

    public func encodeArray(_ array: [Bool]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayStrings(_ array: [String]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayInt(_ array: [Int]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayUInt(_ array: [UInt]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayFloat(_ array: [Float]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayDouble(_ array: [Double]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public func encodeArrayData(_ array: [Data]) {
        var valBuf = MMBuffer()

        for element in array {
            let encoder = MMEncoder()
            encoder.encode(element)
            let data = encoder.buffer.data
            valBuf.write([UInt8](data))
        }

        let payload = [UInt8](valBuf.data)
        let len = payload.count

        if len < 254 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 65536 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len & 0xFF))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
    }

    public var result: Data {
        return buffer.data
    }
}