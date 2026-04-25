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

print("Input JSONC:")
print(jsonc)

// 解析 JSONC
if let parsed = try? JSONC.parse(jsonc) {
    print("\nParsed:")
    print(parsed.description)
}

// 编码到 Wire 格式
// 注意: 需要先将 JSONC 转换为可编码的格式
print("\nEncoded Wire:")
