import MetaMessage

// JSONC 字符串
let jsonc = """
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25,
    // mm: type=bool; desc=是否激活
    "active": true,
    // mm: type=array; child_type=i; desc=分数
    "scores": [95, 87, 92]
}
"""

print("Input JSONC:")
print(jsonc)

// 解析并绑定到对象
if let node = try? JSONC.parse(jsonc) {
    print("\nParsed JSONC:")
    print(node.description)
}

// 使用 valueToJSONC 转换
let user = User(name: "Alice", age: 25, active: true, scores: [95, 87, 92])
if let jsoncNode = try? valueToJSONC(user) {
    print("\nBound to object:")
    print(jsoncNode.description)
}
