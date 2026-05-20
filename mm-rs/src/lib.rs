pub mod core;
pub mod ir;
pub mod jsonc;

use core::{Decoder, Encoder};
use ir::ast::Node;
use jsonc::Parser;

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
    ir::to_string(node)
}

#[cfg(test)]
mod tests {
    use super::*;
    use ir::{Field, Node, Object, Tag, Value, ValueData, ValueType};

    #[test]
    fn test_tag_parse_direct() {
        let result = Tag::parse("// mm:type=str;desc=用户名");
        assert!(result.is_some(), "Tag::parse should succeed");
        let tag = result.unwrap();
        assert_eq!(tag.value_type, ValueType::String);
        assert_eq!(tag.desc, Some("用户名".to_string()));
    }

    #[test]
    fn test_simple_encode_decode() {
        let obj = Node::Object(Object {
            fields: vec![Field {
                key: "name".to_string(),
                value: Node::Value(Value {
                    data: ValueData::String("test".to_string()),
                    text: "test".to_string(),
                    path: String::new(),
                    tag: None,
                }),
            }],
            tag: Some(Tag::new()),
            path: String::new(),
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
                let tag = tag.unwrap();
                assert_eq!(tag.value_type, ValueType::String);
                assert_eq!(tag.desc, Some("用户名".to_string()));
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

    // ===== Tag parsing tests (aligned with Go's tag_test.go) =====

    #[test]
    fn test_tag_parse_basic() {
        let tag =
            Tag::parse("// mm:name=id;min=1;desc=用户ID;enum=active|pending|deleted").unwrap();
        assert_eq!(tag.value_type, ValueType::Enum);
        assert_eq!(tag.desc, Some("用户ID".to_string()));
        assert_eq!(tag.min, Some("1".to_string()));
        assert_eq!(tag.enum_values, Some("active|pending|deleted".to_string()));
    }

    #[test]
    fn test_tag_parse_flags() {
        let tag = Tag::parse("// mm:nullable;default=abc;max=10;min=5").unwrap();
        assert!(tag.nullable);
        assert_eq!(tag.default, Some("abc".to_string()));
        assert_eq!(tag.max, Some("10".to_string()));
        assert_eq!(tag.min, Some("5".to_string()));
    }

    #[test]
    fn test_tag_parse_quoted_and_semicolon() {
        let tag = Tag::parse(r#"// mm:name="id";desc="用户ID";enum="active|pending";pattern="^a,b$";type=str;min=1;max=5;nullable;default="x""#).unwrap();
        assert_eq!(tag.desc, Some("用户ID".to_string()));
        assert_eq!(tag.enum_values, Some("active|pending".to_string()));
        assert_eq!(tag.pattern, Some("^a,b$".to_string()));
        assert_eq!(tag.value_type, ValueType::String);
        assert_eq!(tag.min, Some("1".to_string()));
        assert_eq!(tag.max, Some("5".to_string()));
        assert!(tag.nullable);
        assert_eq!(tag.default, Some("x".to_string()));
    }

    #[test]
    fn test_tag_parse_block_comment() {
        let tag = Tag::parse("/* mm:type=i64;desc=用户ID */").unwrap();
        assert_eq!(tag.value_type, ValueType::Int64);
        assert_eq!(tag.desc, Some("用户ID".to_string()));
    }

    #[test]
    fn test_tag_parse_non_mm_comment() {
        let result = Tag::parse("// this is a regular comment");
        assert!(result.is_none());
    }

    #[test]
    fn test_tag_parse_empty_mm() {
        let tag = Tag::parse("// mm:").unwrap();
        assert_eq!(tag.value_type, ValueType::Unknown);
    }

    #[test]
    fn test_tag_parse_type_aliases() {
        let tag = Tag::parse("// mm:type=slice").unwrap();
        assert_eq!(tag.value_type, ValueType::Slice);
        let tag = Tag::parse("// mm:type=array").unwrap();
        assert_eq!(tag.value_type, ValueType::Array);
        let tag = Tag::parse("// mm:type=vec").unwrap();
        assert_eq!(tag.value_type, ValueType::Slice);
        let tag = Tag::parse("// mm:type=arr").unwrap();
        assert_eq!(tag.value_type, ValueType::Array);
    }

    // ===== Parser edge case tests (aligned with Go's parser_test.go) =====

    #[test]
    fn test_parse_object_empty() {
        let node = parse_jsonc(r#"{}"#).unwrap();
        match node {
            Node::Object(o) => assert_eq!(o.fields.len(), 0),
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_valid() {
        let node = parse_jsonc(r#"{"name": "Alice", "age": 18}"#).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 2);
                assert_eq!(o.fields[0].key, "name");
                assert_eq!(o.fields[1].key, "age");
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_eof_after_colon() {
        let node = parse_jsonc(r#"{"key": "#).unwrap();
        match node {
            Node::Object(o) => assert_eq!(o.fields.len(), 0),
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_eof_after_partial() {
        let node = parse_jsonc(r#"{"name": "Alice", "age": "#).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 1);
                assert_eq!(o.fields[0].key, "name");
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_eof_after_comma() {
        let node = parse_jsonc(r#"{"key": "val", "#).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 1);
                assert_eq!(o.fields[0].key, "key");
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_rbrace_after_colon() {
        let result = parse_jsonc(r#"{"key":}"#);
        assert!(result.is_err());
    }

    #[test]
    fn test_parse_array_eof_after_open() {
        let node = parse_jsonc(r#"["#).unwrap();
        match node {
            Node::Array(a) => assert_eq!(a.items.len(), 0),
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_parse_array_eof_after_item() {
        let node = parse_jsonc(r#"[1, "#).unwrap();
        match node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 1);
            }
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_parse_array_empty() {
        let node = parse_jsonc(r#"[]"#).unwrap();
        match node {
            Node::Array(a) => assert_eq!(a.items.len(), 0),
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_parse_array_with_items() {
        let node = parse_jsonc(r#"[1, "two", true, null, 3.14]"#).unwrap();
        match node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 5);
            }
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_parse_nested_objects() {
        let node = parse_jsonc(r#"{"outer": {"inner": "value"}}"#).unwrap();
        match node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 1);
                match &o.fields[0].value {
                    Node::Object(inner) => {
                        assert_eq!(inner.fields.len(), 1);
                        assert_eq!(inner.fields[0].key, "inner");
                    }
                    _ => panic!("expected nested object"),
                }
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_nested_arrays() {
        let node = parse_jsonc(r#"[[1, 2], [3, 4]]"#).unwrap();
        match node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 2);
            }
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_parse_negative_number() {
        let node = parse_jsonc(r#"{"val": -42}"#).unwrap();
        match node {
            Node::Object(o) => {
                let val = &o.fields[0].value;
                let tag = val.get_tag().unwrap();
                assert_eq!(tag.value_type, ValueType::Int);
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_float_number() {
        let node = parse_jsonc(r#"{"val": 3.14}"#).unwrap();
        match node {
            Node::Object(o) => {
                let val = &o.fields[0].value;
                let tag = val.get_tag().unwrap();
                assert_eq!(tag.value_type, ValueType::Float64);
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_boolean() {
        let node = parse_jsonc(r#"{"a": true, "b": false}"#).unwrap();
        match node {
            Node::Object(o) => {
                let tag_a = o.fields[0].value.get_tag().unwrap();
                let tag_b = o.fields[1].value.get_tag().unwrap();
                assert_eq!(tag_a.value_type, ValueType::Bool);
                assert_eq!(tag_b.value_type, ValueType::Bool);
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_val_nil_does_not_crash() {
        let cases = vec![
            ("EOF after open brace", "{", true),
            ("EOF after open bracket", "[", true),
            ("EOF after key colon", r#"{"a": "#, true),
            ("EOF after comma in object", r#"{"a": 1, "#, true),
            ("EOF after comma in array", "[1, ", true),
        ];
        for (name, input, expect_ok) in cases {
            let result = parse_jsonc(input);
            assert_eq!(result.is_ok(), expect_ok, "test case '{}' failed", name);
        }
    }

    // ===== Tag with inheritance tests =====

    #[test]
    fn test_tag_inheritance() {
        let mut child = Tag::new();
        let parent = Tag {
            child_desc: Some("用户名".to_string()),
            child_type: ValueType::String,
            child_nullable: true,
            ..Tag::new()
        };
        child.inherit(&parent);
        assert_eq!(child.desc, Some("用户名".to_string()));
        assert_eq!(child.value_type, ValueType::String);
        assert!(child.nullable);
        assert!(child.is_inherit);
    }

    // ===== Comment edge case tests =====

    #[test]
    fn test_parse_comment_blank_line_separator() {
        let input = "{\n  // mm:type=str;desc=用户名\n\n  \"name\": \"test\"\n}";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                let tag = o.fields[0].value.get_tag().unwrap();
                assert_eq!(
                    tag.desc, None,
                    "comment separated by blank line should be dropped"
                );
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_multiple_leading_comments() {
        let input = "{\n  // mm:type=str\n  // mm:desc=用户名\n  \"name\": \"test\"\n}";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                let tag = o.fields[0].value.get_tag().unwrap();
                assert_eq!(tag.value_type, ValueType::String);
                assert_eq!(tag.desc, Some("用户名".to_string()));
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_trailing_comment_merges_tag() {
        let input = "{\n  \"name\": \"test\" // mm:desc=用户名\n}";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                let tag = o.fields[0].value.get_tag().unwrap();
                assert_eq!(tag.desc, Some("用户名".to_string()));
                assert_eq!(tag.value_type, ValueType::String);
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_object_comment_before_brace() {
        let input = "// mm:desc=用户对象\n{\n  \"name\": \"test\"\n}";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                let tag = o.tag.as_ref().unwrap();
                assert_eq!(tag.desc, Some("用户对象".to_string()));
            }
            _ => panic!("expected object"),
        }
    }

    #[test]
    fn test_parse_array_size_based_type() {
        let input = "// mm:size=3\n[1, 2, 3]";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Array(a) => {
                let tag = a.tag.as_ref().unwrap();
                assert_eq!(tag.value_type, ValueType::Array);
                assert_eq!(tag.size, Some(3));
            }
            _ => panic!("expected array"),
        }
    }

    #[test]
    fn test_comment_with_mm_only_no_tag() {
        let input = "{\n  // this is a normal comment\n  \"name\": \"test\"\n}";
        let node = parse_jsonc(input).unwrap();
        match node {
            Node::Object(o) => {
                let tag = o.fields[0].value.get_tag();
                assert!(tag.is_some());
                let tag = tag.unwrap();
                assert_eq!(tag.desc, None);
                assert_eq!(tag.value_type, ValueType::String);
            }
            _ => panic!("expected object"),
        }
    }
}
