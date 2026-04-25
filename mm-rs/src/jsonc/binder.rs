use std::collections::HashMap;
use crate::jsonc::ast::{Node, Object, Array, Value, Field, ValueData};
use crate::jsonc::tag::Tag;

pub fn bind(node: &Node, target: &mut dyn std::any::Any) -> Result<(), String> {
    match node {
        Node::Object(obj) => bind_object(obj, target),
        Node::Array(arr) => bind_array(arr, target),
        Node::Value(val) => bind_value(val, target),
    }
}

fn bind_object(obj: &Object, target: &mut dyn std::any::Any) -> Result<(), String> {
    if let Some(map) = target.downcast_mut::<HashMap<String, serde_json::Value>>() {
        for field in &obj.fields {
            let value = node_to_jsonvalue(&field.value)?;
            map.insert(field.key.clone(), value);
        }
        return Ok(());
    }

    Err("unsupported target type for object".to_string())
}

fn bind_array(arr: &Array, target: &mut dyn std::any::Any) -> Result<(), String> {
    if let Some(vec) = target.downcast_mut::<Vec<serde_json::Value>>() {
        for item in &arr.items {
            let value = node_to_jsonvalue(item)?;
            vec.push(value);
        }
        return Ok(());
    }

    Err("unsupported target type for array".to_string())
}

fn bind_value(val: &Value, target: &mut dyn std::any::Any) -> Result<(), String> {
    if let Some(s) = target.downcast_mut::<String>() {
        *s = val.text.clone();
        return Ok(());
    }

    if let Some(b) = target.downcast_mut::<bool>() {
        if let ValueData::Bool(v) = &val.data {
            *b = *v;
            return Ok(());
        }
    }

    if let Some(i) = target.downcast_mut::<i64>() {
        if let ValueData::Int(v) = &val.data {
            *i = *v;
            return Ok(());
        }
    }

    if let Some(f) = target.downcast_mut::<f64>() {
        if let ValueData::Float(v) = &val.data {
            *f = *v;
            return Ok(());
        }
    }

    Err("unsupported target type for value".to_string())
}

fn node_to_jsonvalue(node: &Node) -> Result<serde_json::Value, String> {
    match node {
        Node::Value(v) => value_to_jsonvalue(v),
        Node::Object(o) => {
            let mut map = serde_json::Map::new();
            for field in &o.fields {
                map.insert(field.key.clone(), node_to_jsonvalue(&field.value)?);
            }
            Ok(serde_json::Value::Object(map))
        }
        Node::Array(a) => {
            let mut vec = Vec::new();
            for item in &a.items {
                vec.push(node_to_jsonvalue(item)?);
            }
            Ok(serde_json::Value::Array(vec))
        }
    }
}

fn value_to_jsonvalue(val: &Value) -> Result<serde_json::Value, String> {
    match &val.data {
        ValueData::Bool(b) => Ok(serde_json::Value::Bool(*b)),
        ValueData::String(s) => Ok(serde_json::Value::String(s.clone())),
        ValueData::Int(i) => Ok(serde_json::Value::Number((*i).into())),
        ValueData::Uint(u) => Ok(serde_json::Value::Number((*u).into())),
        ValueData::Float(f) => {
            Ok(serde_json::Number::from_f64(*f)
                .map(serde_json::Value::Number)
                .unwrap_or(serde_json::Value::Null))
        }
        ValueData::Bytes(b) => {
            Ok(serde_json::Value::String(format!("{:?}", b)))
        }
        ValueData::Null => Ok(serde_json::Value::Null),
    }
}