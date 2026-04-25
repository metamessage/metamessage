import Foundation

public class MMBuffer {
    private var buffer: [UInt8]
    private var offset: Int

    public init(capacity: Int = 1024) {
        self.buffer = [UInt8](repeating: 0, count: capacity)
        self.offset = 0
    }

    public init(data: Data) {
        self.buffer = [UInt8](data)
        self.offset = 0
    }

    public var count: Int {
        return offset
    }

    public var data: Data {
        return Data(buffer[0..<offset])
    }

    public func reset() {
        offset = 0
    }

    public func seek(to position: Int) {
        offset = position
    }

    public func position() -> Int {
        return offset
    }

    private func ensureCapacity(_ needed: Int) {
        if offset + needed > buffer.count {
            let newCount = max(buffer.count * 2, offset + needed)
            buffer.append(contentsOf: [UInt8](repeating: 0, count: newCount - buffer.count))
        }
    }

    public func write(_ byte: UInt8) {
        ensureCapacity(1)
        buffer[offset] = byte
        offset += 1
    }

    public func write(_ bytes: [UInt8]) {
        ensureCapacity(bytes.count)
        buffer[offset..<offset + bytes.count].copy(from: bytes)
        offset += bytes.count
    }

    public func write<T: FixedWidthInteger>(_ value: T) {
        var v = value
        let bytes = withUnsafeBytes(of: &v) { Array($0) }
        write(bytes)
    }

    public func writeFloat32(_ value: Float) {
        var v = value
        let bytes = withUnsafeBytes(of: &v) { Array($0) }
        write(bytes)
    }

    public func writeFloat64(_ value: Double) {
        var v = value
        let bytes = withUnsafeBytes(of: &v) { Array($0) }
        write(bytes)
    }

    public func read() -> UInt8? {
        guard offset < buffer.count else { return nil }
        let byte = buffer[offset]
        offset += 1
        return byte
    }

    public func read(_ count: Int) -> [UInt8]? {
        guard offset + count <= buffer.count else { return nil }
        let bytes = Array(buffer[offset..<offset + count])
        offset += count
        return bytes
    }

    public func peek() -> UInt8? {
        guard offset < buffer.count else { return nil }
        return buffer[offset]
    }

    public func remaining() -> Int {
        return buffer.count - offset
    }

    public func readUInt8() throws -> UInt8 {
        guard let byte = read() else {
            throw MMError.unexpectedEndOfData
        }
        return byte
    }

    public func readInt8() throws -> Int8 {
        let byte = try readUInt8()
        return Int8(bitPattern: byte)
    }

    public func readUInt16() throws -> UInt16 {
        guard let bytes = read(2) else {
            throw MMError.unexpectedEndOfData
        }
        return UInt16(bytes[0]) << 8 | UInt16(bytes[1])
    }

    public func readInt16() throws -> Int16 {
        return Int16(bitPattern: try readUInt16())
    }

    public func readUInt32() throws -> UInt32 {
        guard let bytes = read(4) else {
            throw MMError.unexpectedEndOfData
        }
        return UInt32(bytes[0]) << 24 | UInt32(bytes[1]) << 16 | UInt32(bytes[2]) << 8 | UInt32(bytes[3])
    }

    public func readInt32() throws -> Int32 {
        return Int32(bitPattern: try readUInt32())
    }

    public func readUInt64() throws -> UInt64 {
        guard let bytes = read(8) else {
            throw MMError.unexpectedEndOfData
        }
        var result: UInt64 = 0
        for i in 0..<8 {
            result = result << 8 | UInt64(bytes[i])
        }
        return result
    }

    public func readInt64() throws -> Int64 {
        return Int64(bitPattern: try readUInt64())
    }

    public func readFloat32() throws -> Float {
        guard let bytes = read(4) else {
            throw MMError.unexpectedEndOfData
        }
        return Float(bitPattern: UInt32(bytes[0]) << 24 | UInt32(bytes[1]) << 16 | UInt32(bytes[2]) << 8 | UInt32(bytes[3]))
    }

    public func readFloat64() throws -> Double {
        guard let bytes = read(8) else {
            throw MMError.unexpectedEndOfData
        }
        var result: UInt64 = 0
        for i in 0..<8 {
            result = result << 8 | UInt64(bytes[i])
        }
        return Double(bitPattern: result)
    }
}

public enum MMError: Error {
    case unexpectedEndOfData
    case invalidData
    case invalidPrefix
    case invalidTag
    case typeMismatch
    case overflow
}