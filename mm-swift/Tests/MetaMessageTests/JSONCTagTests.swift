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
        XCTAssertEqual(tag?.type, .str)
    }

    func testParseTypeInt() {
        let tag = parseMMTag("// mm:type=i")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .i)
    }

    func testParseTypeFloat() {
        let tag = parseMMTag("// mm:type=f64")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .f64)
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
        let tag = parseMMTag("// mm:default_val=value")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.defaultVal, "value")
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
        let tag = parseMMTag("// mm:enums=a|b|c")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .enums)
        XCTAssertEqual(tag?.enums, "a|b|c")
    }

    func testParsePattern() {
        let tag = parseMMTag("// mm:pattern=^[a-z]+$")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.pattern, "^[a-z]+$")
    }

    func testParseLocation() {
        let tag = parseMMTag("// mm:location=8")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.location, 8)
    }

    func testParseVersion() {
        let tag = parseMMTag("// mm:version=4")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.version, 4)
    }

    func testParseChildType() {
        let tag = parseMMTag("// mm:child_type=str")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.childType, .str)
    }

    func testParseChildDesc() {
        let tag = parseMMTag("// mm:child_desc=\"child description\"")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.childDesc, "child description")
    }

    func testParseMultipleTags() {
        let tag = parseMMTag("// mm:type=str;desc=\"description\";nullable;size=100")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .str)
        XCTAssertEqual(tag?.desc, "description")
        XCTAssertEqual(tag?.nullable, true)
        XCTAssertEqual(tag?.size, 100)
    }

    func testParseComplexTag() {
        let tag = parseMMTag("// mm:type=array;size=5;child_type=i;child_nullable;desc=\"array of ints\"")
        XCTAssertNotNil(tag)
        XCTAssertEqual(tag?.type, .arr)
        XCTAssertEqual(tag?.size, 5)
        XCTAssertEqual(tag?.childType, .i)
        XCTAssertEqual(tag?.childNullable, true)
        XCTAssertEqual(tag?.desc, "array of ints")
    }

    func testTagStringValue() {
        let tag = Tag()
        tag.type = .str
        tag.desc = "test"
        tag.nullable = true

        let str = tag.stringValue()
        XCTAssertTrue(str.contains("type=str"))
        XCTAssertTrue(str.contains("desc=\"test\""))
        XCTAssertTrue(str.contains("nullable"))
    }

    func testTagInherit() {
        let parent = Tag()
        parent.childType = .i
        parent.childNullable = true
        parent.childDefaultVal = "0"

        let child = Tag()
        child.inherit(from: parent)

        XCTAssertEqual(child.type, .i)
        XCTAssertEqual(child.nullable, true)
        XCTAssertEqual(child.defaultVal, "0")
    }
}