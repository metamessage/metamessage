import XCTest
@testable import MetaMessage

final class MMEncoderTests: XCTestCase {
    func testEncodeBool() {
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()
        encoder.encodeNil()
        let data = encoder.buffer.data
        XCTAssertEqual(data.count, 1)
        XCTAssertEqual(data[0], MMSimpleValue.nullInt.rawValue)
    }

    func testEncodeInt() {
        let encoder = MMEncoder()

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
        XCTAssertEqual(data.count, 1)
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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

        encoder.encode(UInt8(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt8(255))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)
    }

    func testEncodeUInt16() {
        let encoder = MMEncoder()

        encoder.encode(UInt16(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt16(65535))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt32() {
        let encoder = MMEncoder()

        encoder.encode(UInt32(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt32(4294967295))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeUInt64() {
        let encoder = MMEncoder()

        encoder.encode(UInt64(0))
        XCTAssertEqual(encoder.buffer.data.count, 1)

        encoder.reset()
        encoder.encode(UInt64(18446744073709551615))
        var data = encoder.buffer.data
        XCTAssertTrue(data.count > 1)
    }

    func testEncodeFloat() {
        let encoder = MMEncoder()

        encoder.encode(Float(0.0))
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 9)

        encoder.reset()
        encoder.encode(Float(3.14))
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 9)
    }

    func testEncodeDouble() {
        let encoder = MMEncoder()

        encoder.encode(Double(0.0))
        var data = encoder.buffer.data
        XCTAssertEqual(data.count, 9)

        encoder.reset()
        encoder.encode(Double(3.14159265359))
        data = encoder.buffer.data
        XCTAssertEqual(data.count, 9)
    }

    func testEncodeString() {
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()

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
        let encoder = MMEncoder()
        encoder.encodeArray([true, false, true])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayStrings() {
        let encoder = MMEncoder()
        encoder.encodeArrayStrings(["a", "b", "c"])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayInt() {
        let encoder = MMEncoder()
        encoder.encodeArrayInt([10, 20, 30])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayUInt() {
        let encoder = MMEncoder()
        encoder.encodeArrayUInt([100, 200, 300])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayFloat() {
        let encoder = MMEncoder()
        encoder.encodeArrayFloat([1.0, 2.0, 3.0])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayDouble() {
        let encoder = MMEncoder()
        encoder.encodeArrayDouble([1.0, 2.0, 3.0])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeArrayData() {
        let encoder = MMEncoder()
        encoder.encodeArrayData([Data([0x01]), Data([0x02]), Data([0x03])])
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 3)
    }

    func testEncodeEmptyArray() {
        let encoder = MMEncoder()
        encoder.encodeArray([Bool]())
        let data = encoder.buffer.data
        XCTAssertTrue(data.count >= 1)
    }

    func testEncodeLargeArray() {
        let encoder = MMEncoder()
        let largeArray = [Bool](repeating: true, count: 1000)
        encoder.encodeArray(largeArray)
        let data = encoder.buffer.data
        XCTAssertTrue(data.count > 1000)
    }
}