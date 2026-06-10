import MetaMessage

// JSONC 字符串
let jsonc = """
{

    // mm: desc=姓名
    "name": "Alice",

    // mm: desc=年龄
    "age": 25,

    // mm: desc=是否激活
    "active": true,
    
    // mm: desc=分数
    "scores": [95, 87, 92]
}
"""

print("Input JSONC:")
print(jsonc)

// 从 JSONC 编码到 Wire 格式
if let wire = try? fromJSONC(jsonc) {
    print("\nEncoded to Wire:")
    print(bytesToHex(wire))

    // 使用 valueToJSONC 转换
    let user = User(name: "Alice", age: 25, active: true, scores: [95, 87, 92])
    if let jsoncOut = try? valueToJSONC(user) {
        print("\nBound to object:")
        print(jsoncOut)
    }
}

func bytesToHex(_ data: Data) -> String {
    return data.map { String(format: "%02x", $0) }.joined()
}
