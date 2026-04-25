pub mod mm;
pub mod jsonc;

use mm::{Encoder, Decoder};
use jsonc::{Parser, Node};

pub fn encode(node: &Node) -> Vec<u8> {
    let mut encoder = Encoder::new();
    encoder.encode(node)
}

pub fn decode(data: &[u8]) -> Result<Node, std::io::Error> {
    let mut decoder = Decoder::new(data.to_vec());
    decoder.decode()
}

pub fn parse_jsonc(input: &str) -> Result<Node, String> {
    let mut parser = Parser::new(input);
    parser.parse()
}

pub fn to_jsonc_string(node: &Node) -> String {
    jsonc::to_string(node)
}

#[cfg(test)]
mod tests {
    use super::*;
    use jsonc::{Node, Object, Value, Field, ValueData, Tag, ValueType};

    #[test]
    fn test_simple_encode_decode() {
        let obj = Node::Object(Object {
            fields: vec![
                Field {
                    key: "name".to_string(),
                    value: Node::Value(Value {
                        data: ValueData::String("test".to_string()),
                        text: "test".to_string(),
                        tag: None,
                    }),
                },
            ],
            tag: Some(Tag::new()),
        });

        let encoded = encode(&obj);
        assert!(!encoded.is_empty());
    }

    #[test]
    fn test_jsonc_parse() {
        let input = r#"{"name": "test", "value": 123}"#;
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 2);
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_jsonc_with_comment() {
        let input = r#"{
            // mm:type=str;desc=用户名
            "name": "test"
        }"#;
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 1);
                let tag = o.fields[0].value.get_tag();
                assert!(tag.is_some());
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_jsonc_roundtrip() {
        let input = r#"{"name": "test", "age": 25}"#;
        let node = parse_jsonc(input).unwrap();
        let output = to_jsonc_string(&node);
        assert!(output.contains("\"name\""));
        assert!(output.contains("test"));
    }
}