"""
JSONC Round-Trip Example

Shows parsing JSONC (JSON with comments for schema tags),
encoding to MetaMessage binary, and decoding back.

JSONC format supports schema annotations like:
    // mm: type=str; desc=姓名
    "name": "Alice"
"""
from metamessage import parse_jsonc, to_jsonc, encode_from_value, decode_to_value, value_to_node


def main():
    # =============================================
    # Parse JSONC with schema annotations
    # =============================================
    
    jsonc_source = """{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 30,
    // mm: type=bool; desc=是否激活
    "active": true,
    // mm: type=slice; child_type=str; desc=标签
    "tags": ["admin", "user"],
    // mm: type=map; desc=地址
    "address": {
        // mm: type=str; desc=城市
        "city": "Beijing",
        // mm: type=str; desc=邮编
        "zip": "100000"
    }
}"""
    
    # Parse JSONC to Node tree
    node = parse_jsonc(jsonc_source)
    print("=== JSONC Parse ===")
    print(f"Parsed node type: {type(node).__name__}")
    print(f"Fields: {len(node.fields)}")
    for field in node.fields:
        print(f"  {field.key}: {field.value.tag.desc or '(no desc)'}")
    
    # Convert back to JSONC string
    jsonc_output = to_jsonc(node)
    print(f"\n=== To JSONC ===")
    print(jsonc_output)
    
    # =============================================
    # Encode parsed node to binary
    # =============================================
    
    # Encode a dict directly (encode_from_value handles value_to_node internally)
    binary = encode_from_value({
        "name": "Alice",
        "age": 30,
        "active": True,
        "tags": ["admin", "user"],
    })
    result = decode_to_value(binary)
    
    print(f"\n=== Encode/Decode Round-Trip ===")
    print(f"Encoded size: {len(binary)} bytes")
    print(f"Decoded: {result}")


if __name__ == "__main__":
    main()

