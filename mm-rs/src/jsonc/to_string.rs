use std::fmt::Write;
use crate::jsonc::ast::{Node, Object, Array, Value, Field, ValueData};

const INDENT_UNIT: &str = "  ";

pub fn to_string(node: &Node) -> String {
    let mut buf = String::new();
    write_node(&mut buf, node, 0);
    buf
}

fn write_node(buf: &mut String, node: &Node, indent: usize) {
    match node {
        Node::Value(v) => write_value(buf, v),
        Node::Object(o) => write_object(buf, o, indent),
        Node::Array(a) => write_array(buf, a, indent),
    }
}

fn write_value(buf: &mut String, val: &Value) {
    match &val.data {
        ValueData::Bool(b) => {
            buf.push_str(if *b { "true" } else { "false" });
        }
        ValueData::String(s) => {
            write_quoted_string(buf, s);
        }
        ValueData::Int(i) => {
            write!(buf, "{}", i).unwrap();
        }
        ValueData::Uint(u) => {
            write!(buf, "{}", u).unwrap();
        }
        ValueData::Float(f) => {
            write!(buf, "{}", f).unwrap();
        }
        ValueData::Bytes(b) => {
            write_quoted_string(buf, &format!("{:?}", b));
        }
        ValueData::Null => {
            buf.push_str("null");
        }
    }
}

fn write_quoted_string(buf: &mut String, s: &str) {
    buf.push('"');
    for c in s.chars() {
        match c {
            '"' => buf.push_str("\\\""),
            '\\' => buf.push_str("\\\\"),
            '\n' => buf.push_str("\\n"),
            '\r' => buf.push_str("\\r"),
            '\t' => buf.push_str("\\t"),
            _ => buf.push(c),
        }
    }
    buf.push('"');
}

fn write_object(buf: &mut String, obj: &Object, indent: usize) {
    buf.push_str("{\n");

    for (i, field) in obj.fields.iter().enumerate() {
        write_indent(buf, indent + 1);

        if let Some(tag) = field.value.get_tag() {
            let tag_str = format!("{:?}", tag);
            if !tag_str.is_empty() {
                write_indent(buf, indent + 1);
                buf.push_str("// ");
                buf.push_str(&tag_str);
                buf.push('\n');
                write_indent(buf, indent + 1);
            }
        }

        write_quoted_string(buf, &field.key);
        buf.push_str(": ");
        write_node(buf, &field.value, indent + 1);

        if i < obj.fields.len() - 1 {
            buf.push(',');
        }
        buf.push('\n');
    }

    write_indent(buf, indent);
    buf.push('}');
}

fn write_array(buf: &mut String, arr: &Array, indent: usize) {
    buf.push_str("[\n");

    for (i, item) in arr.items.iter().enumerate() {
        write_indent(buf, indent + 1);
        write_node(buf, item, indent + 1);

        if i < arr.items.len() - 1 {
            buf.push(',');
        }
        buf.push('\n');
    }

    write_indent(buf, indent);
    buf.push(']');
}

fn write_indent(buf: &mut String, indent: usize) {
    for _ in 0..indent {
        buf.push_str(INDENT_UNIT);
    }
}