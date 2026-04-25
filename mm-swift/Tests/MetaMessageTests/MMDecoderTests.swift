import XCTest
@testable import MetaMessage

final class MMDecoderTests: XCTestCase {
    func testDecodeBool() throws {
        let encoder = MMEncoder()
        encoder.encode(true)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .bool(let b) = value else {
            XCTFail("Expected bool")
            return
        }
        XCTAssertEqual(b, true)
    }

    func testDecodeFalse() throws {
        let encoder = MMEncoder()
        encoder.encode(false)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .bool(let b) = value else {
            XCTFail("Expected bool")
            return
        }
        XCTAssertEqual(b, false)
    }

    func testDecodeNil() throws {
        let encoder = MMEncoder()
        encoder.encodeNil()
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .null = value else {
            XCTFail("Expected null")
            return
        }
    }

    func testDecodeInt() throws {
        let encoder = MMEncoder()
        encoder.encode(Int(23))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, 23)
    }

    func testDecodeNegativeInt() throws {
        let encoder = MMEncoder()
        encoder.encode(Int(-7890))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, -7890)
    }

    func testDecodeInt8() throws {
        let encoder = MMEncoder()
        encoder.encode(Int8(127))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(Int8(i), 127)
    }

    func testDecodeInt16() throws {
        let encoder = MMEncoder()
        encoder.encode(Int16(32767))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(Int16(i), 32767)
    }

    func testDecodeInt32() throws {
        let encoder = MMEncoder()
        encoder.encode(Int32(2147483647))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(Int32(i), 2147483647)
    }

    func testDecodeInt64() throws {
        let encoder = MMEncoder()
        encoder.encode(Int64(9223372036854775807))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, 9223372036854775807)
    }

    func testDecodeUInt() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt(987654))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(u, 987654)
    }

    func testDecodeUInt8() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt8(255))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(UInt8(u), 255)
    }

    func testDecodeUInt16() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt16(65535))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(UInt16(u), 65535)
    }

    func testDecodeUInt32() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt32(4294967295))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(UInt32(u), 4294967295)
    }

    func testDecodeUInt64() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt64(18446744073709551615))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(u, 18446744073709551615)
    }

    func testDecodeFloat() throws {
        let encoder = MMEncoder()
        encoder.encode(Float(3.14))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .float(let f) = value else {
            XCTFail("Expected float")
            return
        }
        XCTAssertEqual(Float(f), 3.14, accuracy: 0.001)
    }

    func testDecodeDouble() throws {
        let encoder = MMEncoder()
        encoder.encode(Double(3.14159265359))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .float(let f) = value else {
            XCTFail("Expected float")
            return
        }
        XCTAssertEqual(f, 3.14159265359, accuracy: 0.0001)
    }

    func testDecodeString() throws {
        let encoder = MMEncoder()
        encoder.encode("hello")
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .string(let s) = value else {
            XCTFail("Expected string")
            return
        }
        XCTAssertEqual(s, "hello")
    }

    func testDecodeEmptyString() throws {
        let encoder = MMEncoder()
        encoder.encode("")
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .string(let s) = value else {
            XCTFail("Expected string")
            return
        }
        XCTAssertEqual(s, "")
    }

    func testDecodeData() throws {
        let encoder = MMEncoder()
        encoder.encode(Data([0x01, 0x02, 0x03]))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .data(let d) = value else {
            XCTFail("Expected data")
            return
        }
        XCTAssertEqual(d, Data([0x01, 0x02, 0x03]))
    }

    func testDecodeArrayBool() throws {
        let encoder = MMEncoder()
        encoder.encodeArray([true, false, true])
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .array(let arr) = value else {
            XCTFail("Expected array")
            return
        }
        XCTAssertEqual(arr.count, 3)
    }

    func testDecodeArrayStrings() throws {
        let encoder = MMEncoder()
        encoder.encodeArrayStrings(["a", "b", "c"])
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .array(let arr) = value else {
            XCTFail("Expected array")
            return
        }
        XCTAssertEqual(arr.count, 3)
    }

    func testDecodeArrayInt() throws {
        let encoder = MMEncoder()
        encoder.encodeArrayInt([10, 20, 30])
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .array(let arr) = value else {
            XCTFail("Expected array")
            return
        }
        XCTAssertEqual(arr.count, 3)
    }

    func testDecodeEmptyArray() throws {
        let encoder = MMEncoder()
        encoder.encodeArray([Bool]())
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .array(let arr) = value else {
            XCTFail("Expected array")
            return
        }
        XCTAssertEqual(arr.count, 0)
    }

    func testDecodeUnexpectedEndOfData() {
        let decoder = MMDecoder(data: Data())
        XCTAssertThrowsError(try decoder.decode()) { error in
            if case MMError.unexpectedEndOfData = error {
            } else {
                XCTFail("Expected unexpectedEndOfData")
            }
        }
    }
}