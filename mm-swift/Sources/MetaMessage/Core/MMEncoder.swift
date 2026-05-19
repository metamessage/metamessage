import Foundation

public class MMEncoder {
    let buffer: MMBuffer

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

    public func reset() {
        buffer.reset()
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
        if value < 32 {
            return (0, Int(value))
        } else if value < UInt64(MMConstants.max1Byte) {
            return (1, 0)
        } else if value < UInt64(MMConstants.max2Byte) {
            return (2, 0)
        } else if value < UInt64(MMConstants.max3Byte) {
            return (3, 0)
        } else if value < UInt64(MMConstants.max4Byte) {
            return (4, 0)
        } else if value < UInt64(MMConstants.max5Byte) {
            return (5, 0)
        } else if value < UInt64(MMConstants.max6Byte) {
            return (6, 0)
        } else if value < UInt64(MMConstants.max7Byte) {
            return (7, 0)
        } else {
            return (8, 0)
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

        if len < 30 {
            buffer.write(MMPrefix.prefixString.rawValue | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.prefixString.rawValue | MMConstants.stringLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 30 {
            buffer.write(MMPrefix.prefixBytes.rawValue | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.prefixBytes.rawValue | MMConstants.bytesLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
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

// MARK: - TagKey constants (Go-compatible format)
// Upper 5 bits = key, lower 3 bits = length
// Matches Go's internal/ir/tag.go TagKey
enum TagKey {
    static let isNull: UInt8 = 0 << 3       // 0x00
    static let desc: UInt8 = 2 << 3         // 0x10
    static let type: UInt8 = 3 << 3         // 0x18
    static let raw: UInt8 = 4 << 3          // 0x20
    static let nullable: UInt8 = 5 << 3     // 0x28
    static let allowEmpty: UInt8 = 6 << 3   // 0x30
    static let unique: UInt8 = 7 << 3       // 0x38
    static let defaultValue: UInt8 = 8 << 3 // 0x40
    static let min: UInt8 = 9 << 3          // 0x48
    static let max: UInt8 = 10 << 3         // 0x50
    static let size: UInt8 = 11 << 3        // 0x58
    static let enumValues: UInt8 = 12 << 3  // 0x60
    static let pattern: UInt8 = 13 << 3     // 0x68
    static let location: UInt8 = 14 << 3    // 0x70
    static let version: UInt8 = 15 << 3     // 0x78
    static let mime: UInt8 = 16 << 3        // 0x80
}

// MARK: - JSONCNode encoding
extension MMEncoder {
    public func encodeNodeValue(_ node: JSONCValue) {
        guard let tag = node.getTag() else {
            encodeRawValue(node)
            return
        }

        let payloadStart = buffer.count
        encodeRawValue(node)
        let payloadBytes = Array(buffer.data[payloadStart..<buffer.count])

        if needsTagEncoding(tag) {
            let tagBuf = encodeTagToBytes(tag)
            let tagLenHeader = encodeTagBodyLength(tagBuf.count)
            let totalLen = tagLenHeader.count + tagBuf.count + payloadBytes.count
            let saved = Array(buffer.data[payloadStart..<buffer.count])
            buffer.seek(to: payloadStart)
            writeTagPrefix(totalLen)
            buffer.write(tagLenHeader)
            buffer.write(tagBuf)
            buffer.write(saved)
        }
    }

    public func encodeNodeArray(_ node: MMArray) {
        guard let tag = node.getTag() else {
            var valBuf = MMBuffer()
            for item in node.items {
                if let val = item as? JSONCValue {
                    let encoder = MMEncoder()
                    encoder.encodeNodeValue(val)
                    valBuf.write([UInt8](encoder.buffer.data))
                }
            }
            let payload = [UInt8](valBuf.data)
            let len = payload.count
            if len < 14 {
                buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
            } else if len < 256 {
                buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
                buffer.write(UInt8(len))
            } else {
                buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
                buffer.write(UInt8((len >> 8) & 0xFF))
                buffer.write(UInt8(len & 0xFF))
            }
            buffer.write(payload)
            return
        }

        let payloadStart = buffer.count
        var valBuf = MMBuffer()
        for item in node.items {
            if let val = item as? JSONCValue {
                let encoder = MMEncoder()
                encoder.encodeNodeValue(val)
                valBuf.write([UInt8](encoder.buffer.data))
            }
        }
        let payload = [UInt8](valBuf.data)
        let len = payload.count
        if len < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(len))
        } else if len < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            buffer.write(UInt8(len))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            buffer.write(UInt8((len >> 8) & 0xFF))
            buffer.write(UInt8(len & 0xFF))
        }
        buffer.write(payload)
        let payloadBytes = Array(buffer.data[payloadStart..<buffer.count])

        let tagBuf = encodeTagToBytes(tag)
        let tagLenHeader = encodeTagBodyLength(tagBuf.count)
        let totalLen = tagLenHeader.count + tagBuf.count + payloadBytes.count
        let saved = Array(buffer.data[payloadStart..<buffer.count])
        buffer.seek(to: payloadStart)
        writeTagPrefix(totalLen)
        buffer.write(tagLenHeader)
        buffer.write(tagBuf)
        buffer.write(saved)
    }

    public func encodeNodeObject(_ node: MMObject) {
        var keyBuf = MMBuffer()
        var valBuf = MMBuffer()

        for field in node.fields {
            let encoder = MMEncoder()
            encoder.encode(field.key)
            keyBuf.write([UInt8](encoder.buffer.data))

            let valEncoder = MMEncoder()
            if let val = field.value as? JSONCValue {
                valEncoder.encodeNodeValue(val)
            } else if let arr = field.value as? MMArray {
                valEncoder.encodeNodeArray(arr)
            } else if let obj = field.value as? MMObject {
                valEncoder.encodeNodeObject(obj)
            }
            valBuf.write([UInt8](valEncoder.buffer.data))
        }

        let keyBytes = [UInt8](keyBuf.data)
        let valBytes = [UInt8](valBuf.data)

        let keyBytesLen = keyBytes.count
        var keyArrayBuf = MMBuffer()
        if keyBytesLen < 14 {
            keyArrayBuf.write(MMPrefix.container.rawValue | MMConstants.containerArray | UInt8(keyBytesLen))
        } else if keyBytesLen < 256 {
            keyArrayBuf.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen1Byte)
            keyArrayBuf.write(UInt8(keyBytesLen))
        } else {
            keyArrayBuf.write(MMPrefix.container.rawValue | MMConstants.containerArray | MMConstants.containerLen2Byte)
            keyArrayBuf.write(UInt8((keyBytesLen >> 8) & 0xFF))
            keyArrayBuf.write(UInt8(keyBytesLen & 0xFF))
        }
        keyArrayBuf.write(keyBytes)
        let keyArrayBytes = [UInt8](keyArrayBuf.data)

        let allBytes = keyArrayBytes + valBytes
        let allLen = allBytes.count

        if allLen < 14 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerObject | UInt8(allLen))
        } else if allLen < 256 {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerObject | MMConstants.containerLen1Byte)
            buffer.write(UInt8(allLen))
        } else {
            buffer.write(MMPrefix.container.rawValue | MMConstants.containerObject | MMConstants.containerLen2Byte)
            buffer.write(UInt8((allLen >> 8) & 0xFF))
            buffer.write(UInt8(allLen & 0xFF))
        }
        buffer.write(keyArrayBytes)
        buffer.write(valBytes)

        if let tag = node.getTag() {
            let payloadBytes = Array(buffer.data)
            let tagBuf = encodeTagToBytes(tag)
            let tagLenHeader = encodeTagBodyLength(tagBuf.count)
            let totalLen = tagLenHeader.count + tagBuf.count + payloadBytes.count
            buffer.reset()
            writeTagPrefix(totalLen)
            buffer.write(tagLenHeader)
            buffer.write(tagBuf)
            buffer.write(payloadBytes)
        }
    }

    private func encodeRawValue(_ node: JSONCValue) {
        guard let data = node.data else {
            encodeNil()
            return
        }

        if let boolVal = data as? Bool {
            encode(boolVal)
        } else if let intVal = data as? Int {
            encode(intVal)
        } else if let intVal = data as? Int8 {
            encode(intVal)
        } else if let intVal = data as? Int16 {
            encode(intVal)
        } else if let intVal = data as? Int32 {
            encode(intVal)
        } else if let intVal = data as? Int64 {
            encode(intVal)
        } else if let uintVal = data as? UInt {
            encode(uintVal)
        } else if let uintVal = data as? UInt8 {
            encode(uintVal)
        } else if let uintVal = data as? UInt16 {
            encode(uintVal)
        } else if let uintVal = data as? UInt32 {
            encode(uintVal)
        } else if let uintVal = data as? UInt64 {
            encode(uintVal)
        } else if let floatVal = data as? Float {
            encode(floatVal)
        } else if let doubleVal = data as? Double {
            encode(doubleVal)
        } else if let stringVal = data as? String {
            encode(stringVal)
        } else if let dataVal = data as? Data {
            encode(dataVal)
        } else if data is Bool {
            let b = data as! Bool
            encode(b)
        } else {
            encodeNil()
        }
    }

    private func needsTagEncoding(_ tag: JSONCTag) -> Bool {
        return tag.type != .unknown || tag.name != "" || tag.desc != "" || tag.nullable || tag.isNull || tag.raw || tag.allowEmpty
    }

    private func encodeTagToBytes(_ tag: JSONCTag) -> [UInt8] {
        var bytes: [UInt8] = []

        if tag.isNull {
            bytes.append(TagKey.isNull | 1)
        }
        if !tag.desc.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.desc, value: tag.desc))
        }
        if tag.type != .unknown {
            bytes.append(TagKey.type)
            bytes.append(UInt8(tag.type.rawValue))
        }
        if tag.raw {
            bytes.append(TagKey.raw | 1)
        }
        if tag.nullable {
            bytes.append(TagKey.nullable | 1)
        }
        if tag.allowEmpty {
            bytes.append(TagKey.allowEmpty | 1)
        }
        if tag.unique {
            bytes.append(TagKey.unique | 1)
        }
        if !tag.defaultValue.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.defaultValue, value: tag.defaultValue))
        }
        if !tag.min.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.min, value: tag.min))
        }
        if !tag.max.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.max, value: tag.max))
        }
        if tag.size != 0 {
            bytes.append(contentsOf: encodeTagU64(TagKey.size, value: UInt64(tag.size)))
        }
        if !tag.enumValues.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.enumValues, value: tag.enumValues))
        }
        if !tag.pattern.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.pattern, value: tag.pattern))
        }
        if tag.locationOffset != 0 {
            let v = "\(tag.locationOffset)"
            bytes.append(contentsOf: encodeTagString(TagKey.location, value: v))
        }
        if tag.version != 0 {
            bytes.append(contentsOf: encodeTagU64(TagKey.version, value: UInt64(tag.version)))
        }
        if !tag.mime.isEmpty {
            bytes.append(contentsOf: encodeTagString(TagKey.mime, value: tag.mime))
        }
        return bytes
    }

    private func encodeTagString(_ key: UInt8, value: String) -> [UInt8] {
        guard let valData = value.data(using: .utf8) else { return [] }
        let valBytes = [UInt8](valData)
        let len = valBytes.count
        var bytes: [UInt8] = []
        if len <= 5 {
            bytes.append(key | UInt8(len))
        } else if len < 256 {
            bytes.append(key | 6)
            bytes.append(UInt8(len))
        } else {
            bytes.append(key | 7)
            bytes.append(UInt8((len >> 8) & 0xFF))
            bytes.append(UInt8(len & 0xFF))
        }
        bytes.append(contentsOf: valBytes)
        return bytes
    }

    private func encodeTagU64(_ key: UInt8, value: UInt64) -> [UInt8] {
        var bytes: [UInt8] = []
        switch value {
        case 0:
            bytes.append(key)
        case 1...255:
            bytes.append(key | 1)
            bytes.append(UInt8(value))
        case 256...65535:
            bytes.append(key | 2)
            bytes.append(UInt8((value >> 8) & 0xFF))
            bytes.append(UInt8(value & 0xFF))
        case 65536...16777215:
            bytes.append(key | 3)
            bytes.append(UInt8((value >> 16) & 0xFF))
            bytes.append(UInt8((value >> 8) & 0xFF))
            bytes.append(UInt8(value & 0xFF))
        case 16777216...4294967295:
            bytes.append(key | 4)
            bytes.append(UInt8((value >> 24) & 0xFF))
            bytes.append(UInt8((value >> 16) & 0xFF))
            bytes.append(UInt8((value >> 8) & 0xFF))
            bytes.append(UInt8(value & 0xFF))
        default:
            bytes.append(key | 5)
            for i in 0..<5 {
                bytes.append(UInt8((value >> (32 + i * 8)) & 0xFF))
            }
            bytes.append(UInt8((value >> 32) & 0xFF))
            bytes.append(UInt8((value >> 24) & 0xFF))
            bytes.append(UInt8((value >> 16) & 0xFF))
            bytes.append(UInt8((value >> 8) & 0xFF))
            bytes.append(UInt8(value & 0xFF))
        }
        return bytes
    }

    private func encodeTagBodyLength(_ length: Int) -> [UInt8] {
        if length < 254 {
            return [UInt8(length)]
        } else if length < 65536 {
            return [254, UInt8(length & 0xFF)]
        } else {
            return [255, UInt8((length >> 8) & 0xFF), UInt8(length & 0xFF)]
        }
    }

    private func writeTagPrefix(_ totalLen: Int) {
        if totalLen < 30 {
            buffer.write(MMPrefix.prefixTag.rawValue | UInt8(totalLen))
        } else if totalLen < 256 {
            buffer.write(MMPrefix.prefixTag.rawValue | MMConstants.tagLen1Byte)
            buffer.write(UInt8(totalLen))
        } else {
            buffer.write(MMPrefix.prefixTag.rawValue | MMConstants.tagLen2Byte)
            buffer.write(UInt8((totalLen >> 8) & 0xFF))
            buffer.write(UInt8(totalLen & 0xFF))
        }
    }
}