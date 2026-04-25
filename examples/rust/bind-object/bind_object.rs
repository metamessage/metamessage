use metamessage::{parse_jsonc, to_jsonc_string, jsonc::{Node, Binder}};

#[derive(Debug)]
struct User {
    name: String,
    age: i32,
    active: bool,
}

fn main() {
    // JSONC 字符串
    let jsonc = r#"{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25,
    // mm: type=bool; desc=是否激活
    "active": true
}"#;

    println!("Input JSONC:");
    println!("{}", jsonc);

    // 解析 JSONC
    let node = parse_jsonc(jsonc).unwrap();

    // 绑定到 User
    // 注意: Rust 需要手动实现 Binder trait 或使用反射
    println!("\nBound to object:");
    println!("Node: {:?}", node);
}
