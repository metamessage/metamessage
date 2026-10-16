import XCTest

@testable import MetaMessage

/// 空容器 / is_null 绑定回归用例（对齐 Go zzrepro_test.go）：
/// Go 侧为空容器编码时会注入 example 哨兵项、为 nil 值注入 is_null，
/// 绑定侧必须跳过哨兵项并产出空的非 nil 集合 / nil。
final class EmptyContainerBindTests: XCTestCase {

    final class StrBox: NSObject {
        @objc var names: [Any] = ["stale"]
    }

    final class MapBox: NSObject {
        @objc var m: [String: Any] = ["stale": 1]
    }

    final class PtrBox: NSObject {
        @objc var ptr: Any? = "stale"
    }

    /// 构造一个带 example 哨兵项的 vec 节点（等价于 Go 编码空 []string 的 wire 结构）
    private func exampleVecNode() -> NodeObject {
        let itemTag = Tag()
        itemTag.example = true
        let item = NodeScalar(data: "", text: "", tag: itemTag, path: "names[0]")

        let vecTag = Tag()
        vecTag.type = .vec
        let vec = NodeArray(items: [item], tag: vecTag, path: "names")

        let rootTag = Tag()
        rootTag.type = .obj
        return NodeObject(fields: [Field(key: "names", value: vec)], tag: rootTag, path: "")
    }

    func testEmptyVecBindsToEmptyArray() throws {
        let box = StrBox()
        try bindNode(exampleVecNode(), to: box)
        XCTAssertEqual(box.names.count, 0, "example 哨兵项不得进入绑定结果")
    }

    func testRealItemsStillBind() throws {
        let realItem = NodeScalar(data: "a", text: "a", tag: Tag(), path: "names[0]")
        let vecTag = Tag()
        vecTag.type = .vec
        let vec = NodeArray(items: [realItem], tag: vecTag, path: "names")
        let rootTag = Tag()
        rootTag.type = .obj
        let root = NodeObject(fields: [Field(key: "names", value: vec)], tag: rootTag, path: "")

        let box = StrBox()
        try bindNode(root, to: box)
        XCTAssertEqual(box.names.count, 1)
    }

    func testEmptyMapBindsToEmptyDictionary() throws {
        let entryTag = Tag()
        entryTag.example = true
        let entry = NodeScalar(data: "", text: "", tag: entryTag, path: "m[]")

        let mapTag = Tag()
        mapTag.type = .map
        let map = NodeObject(fields: [Field(key: "", value: entry)], tag: mapTag, path: "m")

        let rootTag = Tag()
        rootTag.type = .obj
        let root = NodeObject(fields: [Field(key: "m", value: map)], tag: rootTag, path: "")

        let box = MapBox()
        try bindNode(root, to: box)
        XCTAssertEqual(box.m.count, 0, "空 map 的 example 占位项不得进入绑定结果")
    }

    func testIsNullBindsToNil() throws {
        let nullTag = Tag()
        nullTag.isNull = true
        nullTag.type = .obj
        let inner = NodeScalar(data: 0, text: "0", tag: Tag(), path: "ptr.a")
        let ptr = NodeObject(fields: [Field(key: "a", value: inner)], tag: nullTag, path: "ptr")

        let rootTag = Tag()
        rootTag.type = .obj
        let root = NodeObject(fields: [Field(key: "ptr", value: ptr)], tag: rootTag, path: "")

        let box = PtrBox()
        try bindNode(root, to: box)
        XCTAssertNil(box.ptr, "is_null 必须绑定成 nil")
    }
}
