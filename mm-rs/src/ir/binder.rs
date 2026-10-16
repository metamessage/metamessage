use crate::ir::ast::{NodeArray, Node, NodeObject, NodeScalar, ValueData};
use std::collections::HashMap;

/// 空容器注入的 example 哨兵项：只为保留元素类型，不参与绑定。
fn is_example(node: &Node) -> bool {
    node.get_tag().is_some_and(|t| t.example)
}

/// 节点是否带 is_null（编码侧真实值为 nil）。
fn is_null(node: &Node) -> bool {
    node.get_tag().is_some_and(|t| t.is_null)
}

pub fn bind(node: &Node, target: &mut dyn std::any::Any) -> Result<(), String> {
    if is_null(node) {
        // is_null：与 Go 对齐，绑定成 null
        if let Some(v) = target.downcast_mut::<serde_json::Value>() {
            *v = serde_json::Value::Null;
            return Ok(());
        }
    }

    match node {
        Node::Object(obj) => bind_object(obj, target),
        Node::Array(arr) => bind_array(arr, target),
        Node::Value(val) => bind_value(val, target),
    }
}

fn bind_object(obj: &NodeObject, target: &mut dyn std::any::Any) -> Result<(), String> {
    if let Some(map) = target.downcast_mut::<HashMap<String, serde_json::Value>>() {
        for field in &obj.fields {
            if is_example(&field.value) {
                continue; // 空 map 的 example 占位项：跳过后得到空 HashMap
            }
            let value = node_to_jsonvalue(&field.value)?;
            map.insert(field.key.clone(), value);
        }
        return Ok(());
    }

    Err("unsupported target type for object".to_string())
}

fn bind_array(arr: &NodeArray, target: &mut dyn std::any::Any) -> Result<(), String> {
    if let Some(vec) = target.downcast_mut::<Vec<serde_json::Value>>() {
        // 与 Go convertVec 对齐：绑定前先清空目标，避免遗留旧元素
        vec.clear();
        for item in &arr.items {
            // example 哨兵项不参与绑定；全为 example（空容器）时产出空 Vec。
            if is_example(item) {
                continue;
            }
            let value = node_to_jsonvalue(item)?;
            vec.push(value);
        }
        return Ok(());
    }

    Err("unsupported target type for array".to_string())
}

fn bind_value(val: &NodeScalar, target: &mut dyn std::any::Any) -> Result<(), String> {
    if val.tag.as_ref().is_some_and(|t| t.is_null) {
        // is_null：字符串/数值目标均绑定成零值，由调用方按 Option 处理
        if let Some(v) = target.downcast_mut::<serde_json::Value>() {
            *v = serde_json::Value::Null;
            return Ok(());
        }
    }

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
    if is_null(node) {
        return Ok(serde_json::Value::Null); // is_null：绑定成 null
    }

    match node {
        Node::Value(v) => value_to_jsonvalue(v),
        Node::Object(o) => {
            let mut map = serde_json::Map::new();
            for field in &o.fields {
                if is_example(&field.value) {
                    continue; // 空容器占位项不绑定
                }
                map.insert(field.key.clone(), node_to_jsonvalue(&field.value)?);
            }
            Ok(serde_json::Value::Object(map))
        }
        Node::Array(a) => {
            let mut vec = Vec::new();
            for item in &a.items {
                if is_example(item) {
                    continue; // 全为 example 时产出空数组
                }
                vec.push(node_to_jsonvalue(item)?);
            }
            Ok(serde_json::Value::Array(vec))
        }
    }
}

fn value_to_jsonvalue(val: &NodeScalar) -> Result<serde_json::Value, String> {
    match &val.data {
        ValueData::Bool(b) => Ok(serde_json::Value::Bool(*b)),
        ValueData::String(s) => Ok(serde_json::Value::String(s.clone())),
        ValueData::Int(i) => Ok(serde_json::Value::Number((*i).into())),
        ValueData::Uint(u) => Ok(serde_json::Value::Number((*u).into())),
        ValueData::Float(f) => Ok(serde_json::Number::from_f64(*f)
            .map(serde_json::Value::Number)
            .unwrap_or(serde_json::Value::Null)),
        ValueData::Bytes(b) => Ok(serde_json::Value::String(format!("{:?}", b))),
        ValueData::Null => Ok(serde_json::Value::Null),
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    use crate::ir::ast::Field;
    use crate::ir::tag::Tag;

    fn scalar(text: &str, example: bool, is_null: bool) -> Node {
        let mut tag = Tag::new();
        tag.example = example;
        tag.is_null = is_null;
        Node::Value(NodeScalar {
            data: ValueData::String(text.to_string()),
            text: text.to_string(),
            tag: Some(tag),
            path: String::new(),
        })
    }

    #[test]
    fn empty_vec_binds_to_empty_vec() {
        // 编码侧空切片：唯一 item 带 example，绑定后必须是空 Vec
        let arr = NodeArray {
            items: vec![scalar("", true, false)],
            tag: Some(Tag::new()),
            path: String::new(),
        };
        let mut target: Vec<serde_json::Value> = vec![serde_json::Value::Bool(true)];
        bind_array(&arr, &mut target).unwrap();
        assert!(target.is_empty());
    }

    #[test]
    fn empty_map_binds_to_empty_map() {
        let obj = NodeObject {
            fields: vec![Field { key: String::new(), value: scalar("", true, false) }],
            tag: Some(Tag::new()),
            path: String::new(),
        };
        let mut target: HashMap<String, serde_json::Value> = HashMap::new();
        target.insert("stale".to_string(), serde_json::Value::Null);
        bind_object(&obj, &mut target).unwrap();
        assert_eq!(target.len(), 1); // 占位项 key 不进入结果
        assert!(target.contains_key("stale"));
    }

    #[test]
    fn real_items_still_bind() {
        let arr = NodeArray {
            items: vec![scalar("a", false, false), scalar("b", false, false)],
            tag: Some(Tag::new()),
            path: String::new(),
        };
        let mut target: Vec<serde_json::Value> = Vec::new();
        bind_array(&arr, &mut target).unwrap();
        assert_eq!(target.len(), 2);
    }

    #[test]
    fn is_null_binds_to_null_json_value() {
        let node = scalar("", false, true);
        let mut target = serde_json::Value::Bool(true);
        bind(&node, &mut target).unwrap();
        assert_eq!(target, serde_json::Value::Null);
    }
}
