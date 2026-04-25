import XCTest
@testable import MetaMessage

final class JSONCTagTests: XCTestCase {
    func testParseEmptyTag() {
        let tag = parseMMTag("")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .unknown)
    }

    func testParseIsNull() {
        let tag = parseMMTag("// mm:is_null")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.isNull, true)
        XCTAssertEqual(tag?.nullable, true)
    }

    func testParseType() {
        let tag = parseMMTag("// mm:type=str")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .string)
    }

    func testParseTypeInt() {
        let tag = parseMMTag("// mm:type=i")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .int)
    }

    func testParseTypeFloat() {
        let tag = parseMMTag("// mm:type=f64")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .float64)
    }

    func testParseTypeBool() {
        let tag = parseMMTag("// mm:type=bool")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .bool)
    }

    func testParseDesc() {
        let tag = parseMMTag("// mm:desc=\"test description\"")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.desc, "test description")
    }

    func testParseNullable() {
        let tag = parseMMTag("// mm:nullable")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.nullable, true)
    }

    func testParseDefault() {
        let tag = parseMMTag("// mm:default=value")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.defaultValue, "value")
    }

    func testParseMinMax() {
        let tag = parseMMTag("// mm:min=1;max=100")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.min, "1")
        XCTAssertEqual(tag?.max, "100")
    }

    func testParseSize() {
        let tag = parseMMTag("// mm:size=10")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.size, 10)
    }

    func testParseEnum() {
        let tag = parseMMTag("// mm:enum=a|b|c")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .enumValue)
        XCTAssertEqual(tag?.enumValues, "a|b|c")
    }

    func testParsePattern() {
        let tag = parseMMTag("// mm:pattern=^[a-z]+$")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.pattern, "^[a-z]+$")
    }

    func testParseLocation() {
        let tag = parseMMTag("// mm:location=8")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.locationOffset, 8)
    }

    func testParseVersion() {
        let tag = parseMMTag("// mm:version=4")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.version, 4)
    }

    func testParseChildType() {
        let tag = parseMMTag("// mm:child_type=str")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.childType, .string)
    }

    func testParseChildDesc() {
        let tag = parseMMTag("// mm:child_desc=\"child description\"")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.childDesc, "child description")
    }

    func testParseMultipleTags() {
        let tag = parseMMTag("// mm:type=str;desc=\"description\";nullable;size=100")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .string)
        XCTAssertEqual(tag?.desc, "description")
        XCTAssertEqual(tag?.nullable, true)
        XCTAssertEqual(tag?.size, 100)
    }

    func testParseComplexTag() {
        let tag = parseMMTag("// mm:type=array;size=5;child_type=i;child_nullable;desc=\"array of ints\"")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .array)
        XCTAssertEqual(tag?.size, 5)
        XCTAssertEqual(tag?.childType, .int)
        XCTAssertEqual(tag?.childNullable, true)
        XCTAssertEqual(tag?.desc, "array of ints")
    }

    func testTagStringValue() {
        let tag = JSONCTag()
        tag.type = .string
        tag.desc = "test"
        tag.nullable = true

        let str = tag.stringValue()
        XCTAssertTrue(str.contains("type=str"))
        XCTAssertTrue(str.contains("desc=\"test\""))
        XCTAssertTrue(str.contains("nullable"))
    }

    func testTagInherit() {
        let parent = JSONCTag()
        parent.childType = .int
        parent.childNullable = true
        parent.childDefault = "0"

        let child = JSONCTag()
        child.inherit(from: parent)

        XCTAssertEqual(child.type, .int)
        XCTAssertEqual(child.nullable, true)
        XCTAssertEqual(child.defaultValue, "0")
    }
}