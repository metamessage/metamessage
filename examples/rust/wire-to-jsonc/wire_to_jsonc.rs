use metamessage::{parse_jsonc, to_jsonc_string, encode, decode};

fn main() {
    // JSONC 字符串
    let jsonc = r#"{
    // mm: type=datetime; desc=创建时间
    "create_time": "2026-01-01 00:00:00",
    // mm: type=str; desc=用户名称
    "user_name": "Alice",
    // mm: type=bool; desc=是否激活
    "is_active": true,
    // mm: type=array; child_type=i
    "scores": [95, 87, 92]
}"#;

    println!("Original JSONC:");
    println!("{}", jsonc);

    // 解析并编码
    let node = parse_jsonc(jsonc).unwrap();
    let wire = encode(&node);

    println!("\nEncoded Wire:");
    println!("{}", bytes_to_hex(&wire));

    // 从 Wire 解码
    let decoded = decode(&wire).unwrap();
    println!("\nDecoded Node: {:?}", decoded);

    // 转回 JSONC 字符串
    let jsonc_out = to_jsonc_string(&decoded);
    println!("\nDecoded to JSONC:");
    println!("{}", jsonc_out);
}

fn bytes_to_hex(bytes: &[u8]) -> String {
    bytes.iter().map(|b| format!("{:02x}", b)).collect()
}
