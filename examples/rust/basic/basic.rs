use metamessage::{encode, decode, parse_jsonc, to_jsonc_string, jsonc::{Node, Object, Value, Field, ValueData, Tag}};

struct Person {
    name: String,
    age: i32,
}

fn main() {
    // 创建 Person 对象
    let person = Person {
        name: "Ed".to_string(),
        age: 30,
    };

    println!("Original: Name={}, Age={}", person.name, person.age);

    // 创建 Node 并编码
    let node = Node::Object(Object {
        fields: vec![
            Field {
                key: "name".to_string(),
                value: Node::Value(Value {
                    data: ValueData::String(person.name.clone()),
                    text: person.name.clone(),
                    tag: None,
                }),
            },
            Field {
                key: "age".to_string(),
                value: Node::Value(Value {
                    data: ValueData::Int(person.age),
                    text: person.age.to_string(),
                    tag: None,
                }),
            },
        ],
        tag: Some(Tag::new()),
    });

    // 编码到 Wire 格式
    let wire = encode(&node);
    println!("Encoded: {:?}", bytes_to_hex(&wire));

    // 从 Wire 解码
    let decoded = decode(&wire).unwrap();
    println!("Decoded: {:?}", decoded);
}

fn bytes_to_hex(bytes: &[u8]) -> String {
    bytes.iter().map(|b| format!("{:02x}", b)).collect()
}
