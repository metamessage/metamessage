use crate::ir::ast::{Array, Node, Object, Value, ValueData};
use crate::ir::value_type::ValueType;
use std::fmt::Write;

const INDENT_UNIT: &str = "\t";

pub fn to_string(node: &Node) -> String {
    let mut buf = String::new();
    write_node(&mut buf, node, 0);
    buf
}

pub fn to_compact_string(node: &Node) -> String {
    let mut buf = String::new();
    write_node_compact(&mut buf, node);
    buf
}

fn write_node(buf: &mut String, node: &Node, indent: usize) {
    match node {
        Node::Value(v) => write_value(buf, v),
        Node::Object(o) => write_object(buf, o, indent),
        Node::Array(a) => write_array(buf, a, indent),
    }
}

fn write_node_compact(buf: &mut String, node: &Node) {
    match node {
        Node::Value(v) => write_value(buf, v),
        Node::Object(o) => write_object_compact(buf, o),
        Node::Array(a) => write_array_compact(buf, a),
    }
}

fn write_value(buf: &mut String, val: &Value) {
    let value_type = val.tag.as_ref().map(|t| t.value_type);
    match &val.data {
        ValueData::Bool(b) => {
            buf.push_str(if *b { "true" } else { "false" });
        }
        ValueData::String(s) => {
            let needs_quotes = val
                .tag
                .as_ref()
                .map(|t| t.value_type.needs_quotes())
                .unwrap_or(true);
            if needs_quotes {
                write_quoted_string(buf, s);
            } else {
                buf.push_str(s);
            }
        }
        ValueData::Int(i) => {
            if value_type == Some(ValueType::Datetime) {
                let naive = chrono::DateTime::from_timestamp(*i, 0)
                    .map(|dt| dt.naive_utc())
                    .unwrap_or_default();
                let s = naive.format("%Y-%m-%d %H:%M:%S").to_string();
                write_quoted_string(buf, &s);
            } else {
                write!(buf, "{}", i).unwrap();
            }
        }
        ValueData::Uint(u) => {
            write!(buf, "{}", u).unwrap();
        }
        ValueData::Float(f) => {
            write!(buf, "{}", f).unwrap();
        }
        ValueData::Bytes(b) => {
            if value_type == Some(ValueType::Uuid) {
                let hex = b.iter().map(|x| format!("{:02x}", x)).collect::<String>();
                let uuid_str = format!(
                    "{}-{}-{}-{}-{}",
                    &hex[0..8],
                    &hex[8..12],
                    &hex[12..16],
                    &hex[16..20],
                    &hex[20..32]
                );
                write_quoted_string(buf, &uuid_str);
            } else {
                write_quoted_string(buf, &format!("{:?}", b));
            }
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

    for (_i, field) in obj.fields.iter().enumerate() {
        if let Some(tag) = field.value.get_tag() {
            let tag_str = tag.to_string();
            if !tag_str.is_empty() {
                write_indent(buf, indent + 1);
                buf.push_str("// mm: ");
                buf.push_str(&tag_str);
                buf.push('\n');
            }
        }

        write_indent(buf, indent + 1);
        write_quoted_string(buf, &field.key);
        buf.push_str(": ");
        write_node(buf, &field.value, indent + 1);

        buf.push(',');
        buf.push('\n');
    }

    write_indent(buf, indent);
    buf.push('}');
}

fn write_object_compact(buf: &mut String, obj: &Object) {
    buf.push('{');

    for (_i, field) in obj.fields.iter().enumerate() {
        if _i > 0 {
            buf.push(',');
        }
        write_quoted_string(buf, &field.key);
        buf.push(':');
        write_node_compact(buf, &field.value);
    }

    buf.push('}');
}

fn write_array(buf: &mut String, arr: &Array, indent: usize) {
    buf.push_str("[\n");

    for (_i, item) in arr.items.iter().enumerate() {
        write_indent(buf, indent + 1);
        write_node(buf, item, indent + 1);

        buf.push(',');
        buf.push('\n');
    }

    write_indent(buf, indent);
    buf.push(']');
}

fn write_array_compact(buf: &mut String, arr: &Array) {
    buf.push('[');

    for (_i, item) in arr.items.iter().enumerate() {
        if _i > 0 {
            buf.push(',');
        }
        write_node_compact(buf, item);
    }

    buf.push(']');
}

fn write_indent(buf: &mut String, indent: usize) {
    for _ in 0..indent {
        buf.push_str(INDENT_UNIT);
    }
}