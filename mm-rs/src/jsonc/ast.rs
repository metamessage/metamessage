use std::collections::HashMap;
use crate::jsonc::tag::Tag;

#[derive(Debug, Clone)]
pub enum ValueData {
    Bool(bool),
    String(String),
    Int(i64),
    Uint(u64),
    Float(f64),
    Bytes(Vec<u8>),
    Null,
}

#[derive(Debug, Clone)]
pub struct Value {
    pub data: ValueData,
    pub text: String,
    pub tag: Option<Tag>,
}

#[derive(Debug, Clone)]
pub struct Field {
    pub key: String,
    pub value: Node,
}

#[derive(Debug, Clone)]
pub struct Object {
    pub fields: Vec<Field>,
    pub tag: Option<Tag>,
}

#[derive(Debug, Clone)]
pub struct Array {
    pub items: Vec<Node>,
    pub tag: Option<Tag>,
}

#[derive(Debug, Clone)]
pub enum Node {
    Value(Value),
    Object(Object),
    Array(Array),
}

impl Node {
    pub fn get_tag(&self) -> Option<&Tag> {
        match self {
            Node::Value(v) => v.tag.as_ref(),
            Node::Object(o) => o.tag.as_ref(),
            Node::Array(a) => a.tag.as_ref(),
        }
    }
}