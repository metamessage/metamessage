use crate::ir::ast::{Array, Node, Object, Value, ValueData};
use crate::ir::ValueType;
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
    let tag_type = val.tag.as_ref().map(|t| t.value_type);

    match &val.data {
        ValueData::Bool(b) => {
            buf.push_str(if *b { "true" } else { "false" });
        }
        ValueData::String(_s) => {
            let needs_quotes = tag_type.map(|t| t.needs_quotes()).unwrap_or(true);
            if needs_quotes {
                write_quoted_string(buf, &val.text);
            } else {
                buf.push_str(&val.text);
            }
        }
        ValueData::Int(_i) => {
            let should_quote = tag_type
                .map(|t| {
                    matches!(
                        t,
                        ValueType::Datetime | ValueType::Date | ValueType::Time | ValueType::Enum
                    )
                })
                .unwrap_or(false);
            if should_quote {
                write_quoted_string(buf, &val.text);
            } else {
                buf.push_str(&val.text);
            }
        }
        ValueData::Uint(u) => {
            write!(buf, "{}", u).unwrap();
        }
        ValueData::Float(f) => {
            write!(buf, "{}", f).unwrap();
        }
        ValueData::Bytes(b) => {
            let is_uuid = tag_type.map(|t| t == ValueType::Uuid).unwrap_or(false);
            let is_bigint = tag_type.map(|t| t == ValueType::Bigint).unwrap_or(false);
            if is_uuid {
                write_quoted_string(buf, &val.text);
            } else if is_bigint {
                buf.push_str(&val.text);
            } else {
                use base64::{engine::general_purpose, Engine as _};
                write_quoted_string(buf, &general_purpose::STANDARD.encode(b));
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
                buf.push('\n');
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
        if let Some(tag) = item.get_tag() {
            let tag_str = tag.to_string();
            if !tag_str.is_empty() {
                buf.push('\n');
                write_indent(buf, indent + 1);
                buf.push_str("// mm: ");
                buf.push_str(&tag_str);
                buf.push('\n');
            }
        }

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
