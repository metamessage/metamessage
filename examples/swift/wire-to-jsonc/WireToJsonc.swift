import MetaMessage

// JSONC 字符串
let jsonc = """
{
    // mm: type=datetime; desc=创建时间
    "create_time": "2026-01-01 00:00:00",
    // mm: type=str; desc=用户名称
    "user_name": "Alice",
    // mm: type=bool; desc=是否激活
    "is_active": true,
    // mm: type=array; child_type=i
    "scores": [95, 87, 92]
}
"""

print("Original JSONC:")
print(jsonc)

// 编码到 Wire
let wire = MetaMessage.encode(Person(name: "test", age: 0))

print("\nEncoded Wire:")
print(bytesToHex(wire))

// 从 Wire 解码到 JSONC
if let jsoncOut = try? MetaMessage.toJSONC(wire) {
    print("\nDecoded to JSONC:")
    print(jsoncOut)
}

func bytesToHex(_ data: Data) -> String {
    return data.map { String(format: "%02x", $0) }.joined()
}
