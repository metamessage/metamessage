use crate::ir::{Tag, Node, Value, ValueData, Field, Object, Array, ValueType};
use crate::core::utils::camel_to_snake;
use std::collections::HashMap;

const MAX_DEPTH: usize = 32;

pub trait ToNode {
    fn to_node(&self, tag: Option<Tag>) -> Node;
}

pub fn nil_to_node(value_type: ValueType) -> Value {
    let mut tag = Tag::new();
    tag.value_type = value_type;
    Value {
        data: ValueData::Null,
        text: "null".to_string(),
        tag: Some(tag),
        path: String::new(),
    }
}

pub fn value_to_node<T: ToNode + ?Sized>(v: &T, tag: Option<Tag>) -> Node {
    v.to_node(tag)
}

pub fn value_to_node_internal<T: ?Sized>(v: &T, tag: Option<Tag>, depth: usize, path: &str) -> Node
where
    T: ToNode,
{
    if depth > MAX_DEPTH {
        return Value {
            data: ValueData::Null,
            text: "null".to_string(),
            tag,
            path: path.to_string(),
        }.into();
    }
    v.to_node(tag)
}

macro_rules! impl_to_node_int {
    ($t:ty, $vt:ident, $f:ident) => {
        impl ToNode for $t {
            fn to_node(&self, tag: Option<Tag>) -> Node {
                let mut tag = tag.unwrap_or_default();
                if tag.value_type == ValueType::Unknown {
                    tag.value_type = ValueType::$vt;
                }
                let text = self.to_string();
                Node::Value(Value {
                    data: ValueData::$f(*self as i64),
                    text,
                    tag: Some(tag),
                    path: String::new(),
                })
            }
        }
    };
}

macro_rules! impl_to_node_uint {
    ($t:ty, $vt:ident, $f:ident) => {
        impl ToNode for $t {
            fn to_node(&self, tag: Option<Tag>) -> Node {
                let mut tag = tag.unwrap_or_default();
                if tag.value_type == ValueType::Unknown {
                    tag.value_type = ValueType::$vt;
                }
                let text = self.to_string();
                Node::Value(Value {
                    data: ValueData::$f(*self as u64),
                    text,
                    tag: Some(tag),
                    path: String::new(),
                })
            }
        }
    };
}

impl_to_node_int!(i8, I8, Int);
impl_to_node_int!(i16, I16, Int);
impl_to_node_int!(i32, I32, Int);
impl_to_node_int!(i64, I64, Int);
impl_to_node_int!(isize, I, Int);
impl_to_node_uint!(u8, U8, Uint);
impl_to_node_uint!(u16, U16, Uint);
impl_to_node_uint!(u32, U32, Uint);
impl_to_node_uint!(u64, U64, Uint);
impl_to_node_uint!(usize, U, Uint);

impl ToNode for f32 {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::F32;
        }
        let text = crate::core::utils::format_float32(*self);
        Node::Value(Value {
            data: ValueData::Float(*self as f64),
            text,
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl ToNode for f64 {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::F64;
        }
        let text = crate::core::utils::format_float64(*self);
        Node::Value(Value {
            data: ValueData::Float(*self),
            text,
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl ToNode for bool {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::Bool;
        }
        let text = if *self { "true".to_string() } else { "false".to_string() };
        Node::Value(Value {
            data: ValueData::Bool(*self),
            text,
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl ToNode for String {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::Str;
        }
        Node::Value(Value {
            data: ValueData::String(self.clone()),
            text: self.clone(),
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl ToNode for &str {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::Str;
        }
        Node::Value(Value {
            data: ValueData::String(self.to_string()),
            text: self.to_string(),
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl<T: ToNode> ToNode for Vec<T> {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        tag.value_type = ValueType::Vec;

        let items: Vec<Node> = self
            .iter()
            .enumerate()
            .map(|(i, item)| {
                let mut item_tag = Tag::new();
                item_tag.inherit(&tag);
                let path = format!("[{}]", i);
                item.to_node_value_with_depth(item_tag, 0, &path)
            })
            .collect();

        let mut set_tag = false;
        for item in &items {
            if let Some(item_tag) = item.get_tag() {
                if !set_tag {
                    let t = &mut tag;
                    t.child_desc.clone_from(&item_tag.desc);
                    t.child_type = item_tag.value_type;
                    t.child_nullable = item_tag.nullable;
                    t.child_raw = item_tag.raw;
                    t.child_allow_empty = item_tag.allow_empty;
                    t.child_unique = item_tag.unique;
                    t.child_default_val.clone_from(&item_tag.default_val);
                    t.child_min.clone_from(&item_tag.min);
                    t.child_max.clone_from(&item_tag.max);
                    t.child_size = item_tag.size;
                    t.child_enums.clone_from(&item_tag.enums);
                    t.child_pattern.clone_from(&item_tag.pattern);
                    t.child_location = item_tag.location;
                    t.child_version = item_tag.version;
                    t.child_mime.clone_from(&item_tag.mime);
                    set_tag = true;
                }
            }
        }

        if items.is_empty() {
            let mut example_tag = Tag::new();
            example_tag.inherit(&tag);
            example_tag.example = true;
            let empty_node: Node = Node::Value(Value {
                data: ValueData::String(String::new()),
                text: String::new(),
                tag: Some(example_tag),
                path: "[0]".to_string(),
            });
            let items = vec![empty_node];
            return Node::Array(Array { items, tag: Some(tag), path: String::new() });
        }

        Node::Array(Array {
            items,
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl<T: ToNode> ToNode for HashMap<String, T> {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        tag.value_type = ValueType::Map;

        let mut fields = Vec::new();
        let mut set_tag = false;

        for (key, val) in self {
            let key_str = camel_to_snake(key);
            let mut item_tag = Tag::new();
            item_tag.inherit(&tag);
            item_tag.name = Some(key_str.clone());
            let path = format!("[{}]", key_str);
            let field_node = val.to_node_value_with_depth(item_tag, 0, &path);
            let child_tag = field_node.get_tag().cloned();

            if !set_tag {
                if let Some(ref ct) = child_tag {
                    tag.child_desc.clone_from(&ct.desc);
                    tag.child_type = ct.value_type;
                    tag.child_nullable = ct.nullable;
                    tag.child_raw = ct.raw;
                    tag.child_allow_empty = ct.allow_empty;
                    tag.child_unique = ct.unique;
                    tag.child_default_val.clone_from(&ct.default_val);
                    tag.child_min.clone_from(&ct.min);
                    tag.child_max.clone_from(&ct.max);
                    tag.child_size = ct.size;
                    tag.child_enums.clone_from(&ct.enums);
                    tag.child_pattern.clone_from(&ct.pattern);
                    tag.child_location = ct.location;
                    tag.child_version = ct.version;
                    tag.child_mime.clone_from(&ct.mime);
                }
                set_tag = true;
            }

            fields.push(Field {
                key: key_str,
                value: field_node,
            });
        }

        if fields.is_empty() {
            let mut example_tag = Tag::new();
            example_tag.inherit(&tag);
            example_tag.example = true;
            fields.push(Field {
                key: String::new(),
                value: Node::Value(Value {
                    data: ValueData::String(String::new()),
                    text: String::new(),
                    tag: Some(example_tag),
                    path: "[]".to_string(),
                }),
            });
        }

        Node::Object(Object {
            fields,
            tag: Some(tag),
            path: String::new(),
        })
    }
}

impl<T: ToNode> ToNode for Option<T> {
    fn to_node(&self, tag: Option<Tag>) -> Node {
        let mut tag = tag.unwrap_or_default();
        tag.nullable = true;

        match self {
            None => {
                tag.is_null = true;
                Node::Value(Value {
                    data: ValueData::Null,
                    text: "null".to_string(),
                    tag: Some(tag),
                    path: String::new(),
                })
            }
            Some(val) => val.to_node(Some(tag)),
        }
    }
}

impl From<Value> for Node {
    fn from(v: Value) -> Self {
        Node::Value(v)
    }
}

impl From<Object> for Node {
    fn from(o: Object) -> Self {
        Node::Object(o)
    }
}

impl From<Array> for Node {
    fn from(a: Array) -> Self {
        Node::Array(a)
    }
}

trait ToNodeExt: ToNode {
    fn to_node_value_with_depth(&self, tag: Tag, depth: usize, path: &str) -> Node {
        if depth > MAX_DEPTH {
            return Node::Value(Value {
                data: ValueData::Null,
                text: "null".to_string(),
                tag: Some(tag),
                path: path.to_string(),
            });
        }
        let mut node = self.to_node(Some(tag));
        match &mut node {
            Node::Value(v) => v.path = path.to_string(),
            Node::Object(o) => o.path = path.to_string(),
            Node::Array(a) => a.path = path.to_string(),
        }
        node
    }
}

impl<T: ToNode> ToNodeExt for T {}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ir::{Tag, ValueType};

    #[test]
    fn test_i32_to_node() {
        let val: i32 = 42;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "42");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::I32);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_i64_to_node() {
        let val: i64 = -100;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "-100");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::I64);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_u8_to_node() {
        let val: u8 = 255;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "255");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::U8);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_f64_to_node() {
        let val: f64 = 3.14;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert!(v.text.contains("3.14"));
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::F64);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_bool_to_node() {
        let val = true;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "true");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::Bool);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_string_to_node() {
        let val = "hello".to_string();
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "hello");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::Str);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_str_to_node() {
        let val: &str = "world";
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "world");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::Str);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_vec_u8_to_node() {
        let val: Vec<u8> = vec![1, 2, 3];
        let node = val.to_node(None);
        match &node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 3);
                assert_eq!(a.tag.as_ref().unwrap().value_type, ValueType::Vec);
            }
            _ => panic!("expected Array node"),
        }
    }

    #[test]
    fn test_option_some_to_node() {
        let val: Option<i32> = Some(42);
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "42");
                assert!(v.tag.as_ref().unwrap().nullable);
                assert!(!v.tag.as_ref().unwrap().is_null);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_option_none_to_node() {
        let val: Option<String> = None;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert!(v.tag.as_ref().unwrap().is_null);
                assert!(v.tag.as_ref().unwrap().nullable);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_vec_to_node() {
        let val: Vec<i32> = vec![1, 2, 3];
        let node = val.to_node(None);
        match &node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 3);
                assert_eq!(a.tag.as_ref().unwrap().value_type, ValueType::Vec);
            }
            _ => panic!("expected Array node"),
        }
    }

    #[test]
    fn test_empty_vec_to_node() {
        let val: Vec<i32> = vec![];
        let node = val.to_node(None);
        match &node {
            Node::Array(a) => {
                assert_eq!(a.items.len(), 1);
                assert!(a.items[0].get_tag().unwrap().example);
            }
            _ => panic!("expected Array node"),
        }
    }

    #[test]
    fn test_hashmap_to_node() {
        let mut val: HashMap<String, i32> = HashMap::new();
        val.insert("key1".to_string(), 100);
        val.insert("key2".to_string(), 200);
        let node = val.to_node(None);
        match &node {
            Node::Object(o) => {
                assert_eq!(o.fields.len(), 2);
                assert_eq!(o.tag.as_ref().unwrap().value_type, ValueType::Map);
            }
            _ => panic!("expected Object node"),
        }
    }

    #[test]
    fn test_nil_to_node() {
        let node = nil_to_node(ValueType::Str);
        assert_eq!(node.text, "null");
        assert_eq!(node.tag.unwrap().value_type, ValueType::Str);
    }

    #[test]
    fn test_with_explicit_tag() {
        let val: i64 = 42;
        let mut tag = Tag::new();
        tag.desc = Some("用户ID".to_string());
        let node = val.to_node(Some(tag));
        match &node {
            Node::Value(v) => {
                assert_eq!(v.tag.as_ref().unwrap().desc, Some("用户ID".to_string()));
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::I64);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_string_with_email_tag() {
        let val = "test@example.com".to_string();
        let mut tag = Tag::new();
        tag.value_type = ValueType::Email;
        let node = val.to_node(Some(tag));
        match &node {
            Node::Value(v) => {
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::Email);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_nested_vec_to_node() {
        let val: Vec<Vec<i32>> = vec![vec![1, 2], vec![3, 4]];
        let node = val.to_node(None);
        match &node {
            Node::Array(outer) => {
                assert_eq!(outer.items.len(), 2);
                assert_eq!(outer.tag.as_ref().unwrap().value_type, ValueType::Vec);
            }
            _ => panic!("expected Array node"),
        }
    }

    #[test]
    fn test_isize_to_node() {
        let val: isize = -1;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "-1");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::I);
            }
            _ => panic!("expected Value node"),
        }
    }

    #[test]
    fn test_usize_to_node() {
        let val: usize = 100;
        let node = val.to_node(None);
        match &node {
            Node::Value(v) => {
                assert_eq!(v.text, "100");
                assert_eq!(v.tag.as_ref().unwrap().value_type, ValueType::U);
            }
            _ => panic!("expected Value node"),
        }
    }
}