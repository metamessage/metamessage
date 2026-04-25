import XCTest
@testable import MetaMessage

final class MetaMessageIntegrationTests: XCTestCase {
    func testEncodeDecodeBool() throws {
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

    func testEncodeDecodeInt() throws {
        let encoder = MMEncoder()
        encoder.encode(Int(123456))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, 123456)
    }

    func testEncodeDecodeNegativeInt() throws {
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

    func testEncodeDecodeFloat() throws {
        let encoder = MMEncoder()
        encoder.encode(Float(3.14159))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .float(let f) = value else {
            XCTFail("Expected float")
            return
        }
        XCTAssertEqual(Float(f), 3.14159, accuracy: 0.001)
    }

    func testEncodeDecodeDouble() throws {
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

    func testEncodeDecodeString() throws {
        let encoder = MMEncoder()
        encoder.encode("hello world")
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .string(let s) = value else {
            XCTFail("Expected string")
            return
        }
        XCTAssertEqual(s, "hello world")
    }

    func testEncodeDecodeData() throws {
        let encoder = MMEncoder()
        encoder.encode(Data([0xDE, 0xAD, 0xBE, 0xEF]))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .data(let d) = value else {
            XCTFail("Expected data")
            return
        }
        XCTAssertEqual(d, Data([0xDE, 0xAD, 0xBE, 0xEF]))
    }

    func testEncodeDecodeArray() throws {
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

    func testEncodeDecodeIntArray() throws {
        let encoder = MMEncoder()
        encoder.encodeArrayInt([10, 20, 30, 40, 50])
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .array(let arr) = value else {
            XCTFail("Expected array")
            return
        }
        XCTAssertEqual(arr.count, 5)
    }

    func testMetaMessageStaticEncode() {
        let data = MetaMessage.encode(Int(42))
        XCTAssertTrue(data.count > 0)

        let data2 = MetaMessage.encode("test")
        XCTAssertTrue(data2.count > 0)

        let data3 = MetaMessage.encode(true)
        XCTAssertTrue(data3.count > 0)
    }

    func testMetaMessageStaticDecode() throws {
        let encoder = MMEncoder()
        encoder.encode(Int(42))
        let data = encoder.buffer.data

        let value = try MetaMessage.decode(data)
        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, 42)
    }

    func testJSONCParserAndPrinter() throws {
        let json = """
        {
            "name": "test",
            "value": 123,
            "active": true
        }
        """

        let node = try parseJSONC(json)
        XCTAssertNotNil(node)

        let printer = JSONCPrinter()
        let output = printer.print(node)
        XCTAssertTrue(output.contains("name"))
        XCTAssertTrue(output.contains("test"))
    }

    func testJSONCCompactPrinter() throws {
        let json = """
        {
            "name": "test",
            "value": 123
        }
        """

        let node = try parseJSONC(json)
        XCTAssertNotNil(node)

        let printer = JSONCPrinter()
        let output = printer.printCompact(node)
        XCTAssertFalse(output.contains("\n"))
        XCTAssertTrue(output.contains("name"))
        XCTAssertTrue(output.contains("test"))
    }

    func testJSONCBinder() throws {
        let json = """
        {
            "name": "test",
            "value": 123
        }
        """

        let node = try parseJSONC(json)
        XCTAssertNotNil(node)

        let binder = JSONCBinder()
        let result = try binder.bind(node, to: TestStruct.self)
        XCTAssertEqual(result.name, "test")
    }

    func testValueTypeParsing() {
        XCTAssertEqual(ValueType.parse("str"), .string)
        XCTAssertEqual(ValueType.parse("i"), .int)
        XCTAssertEqual(ValueType.parse("i64"), .int64)
        XCTAssertEqual(ValueType.parse("u"), .uint)
        XCTAssertEqual(ValueType.parse("f64"), .float64)
        XCTAssertEqual(ValueType.parse("bool"), .bool)
        XCTAssertEqual(ValueType.parse("arr"), .array)
        XCTAssertEqual(ValueType.parse("unknown"), .unknown)
    }

    func testEncodeLargeString() throws {
        let longString = String(repeating: "a", count: 1000)
        let encoder = MMEncoder()
        encoder.encode(longString)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .string(let s) = value else {
            XCTFail("Expected string")
            return
        }
        XCTAssertEqual(s.count, 1000)
    }

    func testEncodeLargeData() throws {
        let longData = Data(repeating: 0xAB, count: 1000)
        let encoder = MMEncoder()
        encoder.encode(longData)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .data(let d) = value else {
            XCTFail("Expected data")
            return
        }
        XCTAssertEqual(d.count, 1000)
    }

    func testEncodeDecodeUInt64Max() throws {
        let encoder = MMEncoder()
        encoder.encode(UInt64(UInt64.max))
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .uint(let u) = value else {
            XCTFail("Expected uint")
            return
        }
        XCTAssertEqual(u, UInt64.max)
    }

    func testEncodeDecodeInt64Min() throws {
        let encoder = MMEncoder()
        encoder.encode(Int64.min)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, Int64.min)
    }

    func testEncodeDecodeInt64Max() throws {
        let encoder = MMEncoder()
        encoder.encode(Int64.max)
        let data = encoder.buffer.data

        let decoder = MMDecoder(data: data)
        let value = try decoder.decode()

        guard case .int(let i) = value else {
            XCTFail("Expected int")
            return
        }
        XCTAssertEqual(i, Int64.max)
    }
}

struct TestStruct: Codable {
    let name: String
    let value: Int
}