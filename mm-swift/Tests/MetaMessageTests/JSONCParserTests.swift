import XCTest
@testable import MetaMessage

final class JSONCParserTests: XCTestCase {
    func testParseEmptyObject() throws {
        let result = try parseJSONC("{}")
        XCTAssertNotNil(result)
        XCTAssertEqual(result?.getType(), .object)
    }

    func testParseEmptyArray() throws {
        let result = try parseJSONC("[]")
        XCTAssertNotNil(result)
        XCTAssertEqual(result?.getType(), .array)
    }

    func testParseSimpleObject() throws {
        let result = try parseJSONC("{\"name\": \"test\"}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        XCTAssertEqual(obj.fields.count, 1)
        XCTAssertEqual(obj.fields[0].key, "name")

        guard let valueNode = obj.fields[0].value as? JSONCValue else {
            XCTFail("Expected JSONCValue")
            return
        }

        XCTAssertEqual(valueNode.data as? String, "test")
    }

    func testParseSimpleArray() throws {
        let result = try parseJSONC("[1, 2, 3]")
        XCTAssertNotNil(result)

        guard let arr = result as? JSONCArray else {
            XCTFail("Expected JSONCArray")
            return
        }

        XCTAssertEqual(arr.items.count, 3)
    }

    func testParseNumbers() throws {
        let result = try parseJSONC("{\"int\": 123, \"float\": 3.14, \"neg\": -456}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        let intField = obj.fields.first { $0.key == "int" }
        XCTAssertNotNil(intField)

        let floatField = obj.fields.first { $0.key == "float" }
        XCTAssertNotNil(floatField)

        let negField = obj.fields.first { $0.key == "neg" }
        XCTAssertNotNil(negField)
    }

    func testParseBooleans() throws {
        let result = try parseJSONC("{\"true\": true, \"false\": false}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        let trueField = obj.fields.first { $0.key == "true" }
        XCTAssertNotNil(trueField)

        let falseField = obj.fields.first { $0.key == "false" }
        XCTAssertNotNil(falseField)
    }

    func testParseNull() throws {
        let result = try parseJSONC("{\"null\": null}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        let nullField = obj.fields.first { $0.key == "null" }
        XCTAssertNotNil(nullField)
    }

    func testParseNestedObject() throws {
        let result = try parseJSONC("{\"outer\": {\"inner\": \"value\"}}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        let outerField = obj.fields.first { $0.key == "outer" }
        XCTAssertNotNil(outerField)

        guard let innerObj = outerField?.value as? JSONCObject else {
            XCTFail("Expected inner JSONCObject")
            return
        }

        let innerField = innerObj.fields.first { $0.key == "inner" }
        XCTAssertNotNil(innerField)
    }

    func testParseNestedArray() throws {
        let result = try parseJSONC("{\"matrix\": [[1, 2], [3, 4]]}")
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        let matrixField = obj.fields.first { $0.key == "matrix" }
        XCTAssertNotNil(matrixField)
    }

    func testParseWithComment() throws {
        let result = try parseJSONC("""
        {
            // mm:type=str
            "name": "test"
        }
        """)
        XCTAssertNotNil(result)
    }

    func testParseWithBlockComment() throws {
        let result = try parseJSONC("""
        {
            /* mm:type=str */
            "name": "test"
        }
        """)
        XCTAssertNotNil(result)
    }

    func testParseComplex() throws {
        let json = """
        {
            "id": 1,
            "name": "test",
            "active": true,
            "items": [1, 2, 3],
            "metadata": {
                "created": "2024-01-01",
                "tags": ["a", "b"]
            }
        }
        """
        let result = try parseJSONC(json)
        XCTAssertNotNil(result)

        guard let obj = result as? JSONCObject else {
            XCTFail("Expected JSONCObject")
            return
        }

        XCTAssertEqual(obj.fields.count, 5)
    }
}