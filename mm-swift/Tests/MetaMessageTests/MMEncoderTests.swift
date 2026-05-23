import XCTest
@testable import MetaMessage

final class EncoderTests: XCTestCase {
    func testEncodeBool() {
        let encoder = Encoder()

        encoder.encode(true)
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)
        XCTAssertEqual(data[0], MMSimpleValue.trueValue.rawValue)

        encoder.reset()
        encoder.encode(false)
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)
        XCTAssertEqual(data[0], MMSimpleValue.falseValue.rawValue)
    }

    func testEncodeNil() {
        let encoder = Encoder()
        encoder.encodeNil()
        let data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)
        XCTAssertEqual(data[0], MMSimpleValue.nullInt.rawValue)
    }

    func testEncodeInt() {
        let encoder = Encoder()

        encoder.encode(Int(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(Int(23))
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)
        XCTAssertEqual(data[0], 0b001_10111)

        encoder.reset()
        encoder.encode(Int(24))
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 2)
        XCTAssertEqual(data[0], 0b001_11000)

        encoder.reset()
        encoder.encode(Int(123456))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)

        encoder.reset()
        encoder.encode(Int(-7890))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeInt8() {
        let encoder = Encoder()

        encoder.encode(Int8(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(Int8(-128))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)

        encoder.reset()
        encoder.encode(Int8(127))
        data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)
    }

    func testEncodeInt16() {
        let encoder = Encoder()

        encoder.encode(Int16(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(Int16(-32768))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)

        encoder.reset()
        encoder.encode(Int16(32767))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeInt32() {
        let encoder = Encoder()

        encoder.encode(Int32(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(Int32(-2147483648))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)

        encoder.reset()
        encoder.encode(Int32(2147483647))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeInt64() {
        let encoder = Encoder()

        encoder.encode(Int64(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(Int64(-9223372036854775808))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)

        encoder.reset()
        encoder.encode(Int64(9223372036854775807))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt() {
        let encoder = Encoder()

        encoder.encode(UInt(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt(123456))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)

        encoder.reset()
        encoder.encode(UInt(987654))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt8() {
        let encoder = Encoder()

        encoder.encode(UInt8(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt8(255))
        let data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)
    }

    func testEncodeUInt16() {
        let encoder = Encoder()

        encoder.encode(UInt16(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt16(65535))
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt32() {
        let encoder = Encoder()

        encoder.encode(UInt32(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt32(4294967295))
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt64() {
        let encoder = Encoder()

        encoder.encode(UInt64(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt64(18446744073709551615))
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeFloat() {
        let encoder = Encoder()

        encoder.encode(Float(0.0))
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)

        encoder.reset()
        encoder.encode(Float(3.14))
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 9)
    }

    func testEncodeDouble() {
        let encoder = Encoder()

        encoder.encode(Double(0.0))
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)

        encoder.reset()
        encoder.encode(Double(3.14159265359))
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 7)
    }

    func testEncodeString() {
        let encoder = Encoder()

        encoder.encode("")
        var data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)

        encoder.reset()
        encoder.encode("hello")
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 5)

        encoder.reset()
        encoder.encode("hello world")
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 11)

        let longString = String(repeating: "a", count: 300)
        encoder.reset()
        encoder.encode(longString)
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 300)
    }

    func testEncodeData() {
        let encoder = Encoder()

        encoder.encode(Data())
        var data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)

        encoder.reset()
        encoder.encode(Data([0x01, 0x02, 0x03]))
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)

        let longData = Data(repeating: 0xAB, count: 300)
        encoder.reset()
        encoder.encode(longData)
        data = encoder.buffer.data
        XCTAssertTrue(data.count > 300)
    }

    func testEncodeArrayBool() {
        let encoder = Encoder()
        encoder.encodeArray([true, false, true])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayStrings() {
        let encoder = Encoder()
        encoder.encodeArrayStrings(["a", "b", "c"])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayInt() {
        let encoder = Encoder()
        encoder.encodeArrayInt([10, 20, 30])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayUInt() {
        let encoder = Encoder()
        encoder.encodeArrayUInt([100, 200, 300])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayFloat() {
        let encoder = Encoder()
        encoder.encodeArrayFloat([1.0, 2.0, 3.0])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayDouble() {
        let encoder = Encoder()
        encoder.encodeArrayDouble([1.0, 2.0, 3.0])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayData() {
        let encoder = Encoder()
        encoder.encodeArrayData([Data([0x01]), Data([0x02]), Data([0x03])])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeEmptyArray() {
        let encoder = Encoder()
        encoder.encodeArray([Bool]())
        let data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)
    }

    func testEncodeLargeArray() {
        let encoder = Encoder()
        let largeArray = [Bool](repeating: true, count: 1000)
        encoder.encodeArray(largeArray)
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 1000)
    }
}