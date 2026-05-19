import XCTest
@testable import MetaMessage

final class ValueToNodeBinderTests: XCTestCase {

    // MARK: - Basic ValueToNode Tests

    func testValueToNodeBool() throws {
        let node = try valueToNode(true, tag: "bool")
        XCTAssertEqual(node.getType(), .value)
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Bool, true)
        XCTAssertEqual(val.text, "true")
    }

    func testValueToNodeInt() throws {
        let node = try valueToNode(42, tag: "int")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Int, 42)
        XCTAssertEqual(val.text, "42")
    }

    func testValueToNodeString() throws {
        let node = try valueToNode("hello", tag: "str")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? String, "hello")
    }

    func testValueToNodeDouble() throws {
        let node = try valueToNode(3.14, tag: "f64")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Double, 3.14)
    }

    func testValueToNodeData() throws {
        let data = Data([0x01, 0x02, 0x03])
        let node = try valueToNode(data, tag: "bytes")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Data, data)
    }

    func testValueToNodeDate() throws {
        let date = Date()
        let node = try valueToNode(date, tag: "datetime")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertNotNil(val.data as? Date)
    }

    func testValueToNodeUUID() throws {
        let uuid = UUID()
        let node = try valueToNode(uuid, tag: "uuid")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? UUID, uuid)
    }

    func testValueToNodeNil() throws {
        let node = try valueToNode(nil, tag: "type=str;nullable")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertNil(val.data)
        XCTAssertEqual(val.text, "null")
        XCTAssertEqual(val.getTag()?.isNull, true)
    }

    // MARK: - Struct Tests

    struct SimpleUser {
        var id: Int = 1001
        var name: String = "张三"
        var age: Int = 28
        var isActive: Bool = true
    }

    func testValueToNodeSimpleStruct() throws {
        let user = SimpleUser()
        let node = try valueToNode(user, tag: "user")
        XCTAssertEqual(node.getType(), .object)
        guard let obj = node as? MMObject else {
            XCTFail("Expected MMObject"); return
        }
        XCTAssertGreaterThan(obj.fields.count, 0)

        let idField = obj.fields.first(where: { $0.key == "id" })
        XCTAssertNotNil(idField)
        if let val = idField?.value as? JSONCValue {
            XCTAssertEqual(val.data as? Int, 1001)
        }

        let nameField = obj.fields.first(where: { $0.key == "name" })
        XCTAssertNotNil(nameField)
        if let val = nameField?.value as? JSONCValue {
            XCTAssertEqual(val.data as? String, "张三")
        }

        let activeField = obj.fields.first(where: { $0.key == "is_active" })
        XCTAssertNotNil(activeField)
        if let val = activeField?.value as? JSONCValue {
            XCTAssertEqual(val.data as? Bool, true)
        }
    }

    struct Address {
        var province: String = "北京市"
        var city: String = "朝阳区"
        var zipCode: String = "100000"
    }

    struct UserWithAddress {
        var id: Int64 = 1001
        var name: String = "张三"
        var addr: Address = Address()
    }

    func testValueToNodeNestedStruct() throws {
        let user = UserWithAddress()
        let node = try valueToNode(user, tag: "user")
        guard let obj = node as? MMObject else {
            XCTFail("Expected MMObject"); return
        }

        let addrField = obj.fields.first(where: { $0.key == "addr" })
        XCTAssertNotNil(addrField)
        XCTAssertEqual(addrField?.value.getType(), .object)

        if let addrObj = addrField?.value as? MMObject {
            let provinceField = addrObj.fields.first(where: { $0.key == "province" })
            XCTAssertNotNil(provinceField)
            if let val = provinceField?.value as? JSONCValue {
                XCTAssertEqual(val.data as? String, "北京市")
            }

            let cityField = addrObj.fields.first(where: { $0.key == "city" })
            XCTAssertNotNil(cityField)
        }
    }

    // MARK: - Array Tests

    struct UserWithTags {
        var tags: [String] = ["swift", "ast", "jsonc"]
    }

    func testValueToNodeArray() throws {
        let user = UserWithTags()
        let node = try valueToNode(user, tag: "user")
        guard let obj = node as? MMObject else {
            XCTFail("Expected MMObject"); return
        }

        let tagsField = obj.fields.first(where: { $0.key == "tags" })
        XCTAssertNotNil(tagsField)
        XCTAssertEqual(tagsField?.value.getType(), .array)

        if let arr = tagsField?.value as? MMArray {
            XCTAssertEqual(arr.items.count, 3)
        }
    }

    func testValueToNodeIntArray() throws {
        let node = try valueToNode([1, 2, 3], tag: "arr_int")
        guard let arr = node as? MMArray else {
            XCTFail("Expected MMArray"); return
        }
        XCTAssertEqual(arr.items.count, 3)
    }

    // MARK: - Round-trip Tests

    func testRoundTripSimpleStruct() throws {
        let user = SimpleUser()
        let data = try MetaMessage.fromValue(user, tag: "simple_user")

        let decoded = try MetaMessage.decode(data)
        switch decoded {
        case .object(let obj):
            XCTAssertNotNil(obj["id"])
            XCTAssertNotNil(obj["name"])
            XCTAssertNotNil(obj["is_active"])
        default:
            XCTFail("Expected object")
        }
    }

    func testRoundTripInt() throws {
        let data = try MetaMessage.fromValue(42, tag: "int")
        let decoded = try MetaMessage.decode(data)
        guard case .int(let i) = decoded else {
            XCTFail("Expected int"); return
        }
        XCTAssertEqual(i, 42)
    }

    func testRoundTripString() throws {
        let data = try MetaMessage.fromValue("hello", tag: "str")
        let decoded = try MetaMessage.decode(data)
        guard case .string(let s) = decoded else {
            XCTFail("Expected string"); return
        }
        XCTAssertEqual(s, "hello")
    }

    func testRoundTripBool() throws {
        let data = try MetaMessage.fromValue(true, tag: "bool")
        let decoded = try MetaMessage.decode(data)
        guard case .bool(let b) = decoded else {
            XCTFail("Expected bool"); return
        }
        XCTAssertEqual(b, true)
    }

    func testRoundTripDouble() throws {
        let data = try MetaMessage.fromValue(3.14159, tag: "f64")
        let decoded = try MetaMessage.decode(data)
        guard case .float(let f) = decoded else {
            XCTFail("Expected float"); return
        }
        XCTAssertEqual(f, 3.14159, accuracy: 0.001)
    }

    // MARK: - ValueToJSONC Tests

    func testValueToJSONCSimple() throws {
        let jsonc = try MetaMessage.valueToJSONC(SimpleUser(), name: "user")
        XCTAssertTrue(jsonc.contains("id"))
        XCTAssertTrue(jsonc.contains("name"))
        XCTAssertTrue(jsonc.contains("is_active"))
    }

    func testValueToJSONCNested() throws {
        let user = UserWithAddress()
        let jsonc = try MetaMessage.valueToJSONC(user, name: "user")
        XCTAssertTrue(jsonc.contains("addr"))
        XCTAssertTrue(jsonc.contains("province"))
        XCTAssertTrue(jsonc.contains("city"))
        XCTAssertTrue(jsonc.contains("北京市"))
    }

    // MARK: - ValueToNode Error Tests

    func testValueToNodeUntypedNil() {
        XCTAssertThrowsError(try valueToNode(nil, tag: "")) { error in
            XCTAssertTrue(error is MMValueToNodeError)
        }
    }

    func testValueToNodeMaxDepth() throws {
        class Recursive {
            var child: Recursive?
            init() {}
        }
        let r = Recursive()
        r.child = r
        XCTAssertThrowsError(try valueToNode(r, tag: "recursive")) { error in
            XCTAssertTrue(error is MMValueToNodeError)
        }
    }

    // MARK: - CamelSnake Conversion

    func testCamelToSnake() {
        XCTAssertEqual(camelToSnake("camelCase"), "camel_case")
        XCTAssertEqual(camelToSnake("simpleUser"), "simple_user")
        XCTAssertEqual(camelToSnake("id"), "id")
        XCTAssertEqual(camelToSnake("XMLParser"), "x_m_l_parser")
    }

    func testSnakeToCamel() {
        XCTAssertEqual(snakeToCamel("snake_case"), "snakeCase")
        XCTAssertEqual(snakeToCamel("simple_user"), "simpleUser")
        XCTAssertEqual(snakeToCamel("id"), "id")
        XCTAssertEqual(snakeToCamel("alreadyCamel"), "alreadycamel")
    }

    // MARK: - Optional Handling

    func testValueToNodeOptional() throws {
        let value: Int? = 42
        let node = try valueToNode(value, tag: "int")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Int, 42)
    }

    func testValueToNodeNilOptionalWithTag() throws {
        let value: String? = nil
        let node = try valueToNode(value, tag: "type=str;nullable")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertNil(val.data)
        XCTAssertEqual(val.getTag()?.isNull, true)
    }

    // MARK: - UInt Types

    func testValueToNodeUInt() throws {
        let node = try valueToNode(UInt(255), tag: "u")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? UInt, 255)
    }

    func testValueToNodeUInt64() throws {
        let node = try valueToNode(UInt64.max, tag: "u64")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? UInt64, UInt64.max)
    }

    // MARK: - Float Types

    func testValueToNodeFloat() throws {
        let node = try valueToNode(Float(3.14), tag: "f32")
        guard let val = node as? JSONCValue else {
            XCTFail("Expected JSONCValue"); return
        }
        XCTAssertEqual(val.data as? Float, Float(3.14))
    }

    // MARK: - Binder Tests (NSObject subclasses)

    class TestBindable: NSObject {
        @objc dynamic var name: String = ""
        @objc dynamic var value: Int = 0
    }

    func testBindSimpleObject() throws {
        let jsonc = """
        {
            "name": "test_name",
            "value": 42
        }
        """

        let obj = TestBindable()
        try MetaMessage.bindFromJSONC(jsonc, to: obj)
        XCTAssertEqual(obj.name, "test_name")
        XCTAssertEqual(obj.value, 42)
    }

    class NestedBindable: NSObject {
        @objc dynamic var title: String = ""
        @objc dynamic var count: Int = 0
    }

    class ParentBindable: NSObject {
        @objc dynamic var id: Int = 0
        @objc dynamic var nested: NestedBindable = NestedBindable()
    }

    func testBindNestedObject() throws {
        let jsonc = """
        {
            "id": 123,
            "nested": {
                "title": "nested_title",
                "count": 99
            }
        }
        """

        let obj = ParentBindable()
        try MetaMessage.bindFromJSONC(jsonc, to: obj)
        XCTAssertEqual(obj.id, 123)
        XCTAssertEqual(obj.nested.title, "nested_title")
        XCTAssertEqual(obj.nested.count, 99)
    }

    // MARK: - MMTagProvider

    struct TaggedStruct: MMTagProvider {
        var id: Int = 1
        var name: String = "test"

        func tag(forField field: String) -> String? {
            switch field {
            case "id": return "min=1;desc=ID"
            case "name": return "required;max_len=20;desc=名称"
            default: return nil
            }
        }
    }

    func testValueToNodeWithTagProvider() throws {
        let tagged = TaggedStruct()
        let node = try valueToNode(tagged, tag: "tagged")
        guard let obj = node as? MMObject else {
            XCTFail("Expected MMObject"); return
        }

        let nameField = obj.fields.first(where: { $0.key == "name" })
        XCTAssertNotNil(nameField)
        let nameTag = nameField?.value.getTag()
        XCTAssertNotNil(nameTag)
        XCTAssertEqual(nameTag?.desc, "名称")
    }

    // MARK: - Nested Array Roundtrip

    struct IntArrayContainer {
        var numbers: [Int] = [10, 20, 30]
    }

    func testRoundTripIntArray() throws {
        let container = IntArrayContainer()
        let data = try MetaMessage.fromValue(container, tag: "container")

        let decoded = try MetaMessage.decode(data)
        switch decoded {
        case .object(let obj):
            XCTAssertNotNil(obj["numbers"])
        default:
            XCTFail("Expected object")
        }
    }

    // MARK: - ValueToNode fromJSONC encoding

    func testFromJSONC() throws {
        let jsonc = """
        {"name": "test", "value": 42}
        """
        let data = try MetaMessage.fromJSONC(jsonc)
        XCTAssertGreaterThan(data.count, 0)

        let decoded = try MetaMessage.decode(data)
        switch decoded {
        case .object(let obj):
            XCTAssertNotNil(obj["name"])
            XCTAssertNotNil(obj["value"])
        default:
            XCTFail("Expected object")
        }
    }

    // MARK: - Nullable handling

    struct NullableFields {
        var name: String? = "hello"
        var optionalNil: String? = nil
    }

    func testValueToNodeNullable() throws {
        let nullable = NullableFields()
        let node = try valueToNode(nullable, tag: "nullable")
        guard let obj = node as? MMObject else {
            XCTFail("Expected MMObject"); return
        }

        let nameField = obj.fields.first(where: { $0.key == "name" })
        XCTAssertNotNil(nameField)
        if let val = nameField?.value as? JSONCValue {
            XCTAssertEqual(val.data as? String, "hello")
        }

        let nilField = obj.fields.first(where: { $0.key == "optional_nil" })
        XCTAssertNotNil(nilField)
        XCTAssertEqual(nilField?.value.getTag()?.isNull, true)
    }
}