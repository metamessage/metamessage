use std::collections::HashMap;
use crate::jsonc::value_type::ValueType;

#[derive(Debug, Clone)]
pub struct Tag {
    pub value_type: ValueType,
    pub desc: Option<String>,
    pub nullable: bool,
    pub is_null: bool,
    pub default: Option<String>,
    pub min: Option<String>,
    pub max: Option<String>,
    pub size: Option<usize>,
    pub enum_values: Option<Vec<String>>,
    pub pattern: Option<String>,
    pub location: Option<String>,
    pub version: Option<String>,
    pub mime: Option<String>,
    pub child_type: Option<ValueType>,
    pub child_desc: Option<String>,
    pub key_desc: Option<String>,
    pub value_desc: Option<String>,
    pub ele_desc: Option<String>,
    pub name: Option<String>,
    pub example: bool,
}

impl Tag {
    pub fn new() -> Self {
        Self {
            value_type: ValueType::Unknown,
            desc: None,
            nullable: false,
            is_null: false,
            default: None,
            min: None,
            max: None,
            size: None,
            enum_values: None,
            pattern: None,
            location: None,
            version: None,
            mime: None,
            child_type: None,
            child_desc: None,
            key_desc: None,
            value_desc: None,
            ele_desc: None,
            name: None,
            example: false,
        }
    }

    pub fn parse(s: &str) -> Option<Self> {
        let s = s.trim();
        if !s.starts_with("//") && !s.starts_with("/*") {
            return None;
        }

        let mut s = s;
        if s.starts_with("//") {
            s = &s[2..];
        } else if s.starts_with("/*") {
            s = &s[2..];
            if s.ends_with("*/") {
                s = &s[..s.len() - 2];
            }
        }

        if !s.trim().starts_with("mm:") {
            return None;
        }

        s = &s.trim()[3..];

        let mut tag = Tag::new();

        let pairs: Vec<&str> = s.split(';').collect();
        for pair in pairs {
            let pair = pair.trim();
            if pair.is_empty() {
                continue;
            }

            if let Some(idx) = pair.find('=') {
                let key = pair[..idx].trim().to_lowercase();
                let value = pair[idx + 1..].trim();

                match key.as_str() {
                    "type" => {
                        tag.value_type = ValueType::from_str(value);
                    }
                    "desc" => {
                        tag.desc = Some(value.to_string());
                    }
                    "nullable" => {
                        tag.nullable = value.eq_ignore_ascii_case("true");
                    }
                    "default" => {
                        tag.default = Some(value.to_string());
                    }
                    "min" => {
                        tag.min = Some(value.to_string());
                    }
                    "max" => {
                        tag.max = Some(value.to_string());
                    }
                    "size" => {
                        tag.size = value.parse().ok();
                    }
                    "enum" => {
                        tag.enum_values = Some(value.split('|').map(|s| s.trim().to_string()).collect());
                    }
                    "pattern" => {
                        tag.pattern = Some(value.to_string());
                    }
                    "location" => {
                        tag.location = Some(value.to_string());
                    }
                    "version" => {
                        tag.version = Some(value.to_string());
                    }
                    "mime" => {
                        tag.mime = Some(value.to_string());
                    }
                    "child_type" => {
                        tag.child_type = Some(ValueType::from_str(value));
                    }
                    "child_desc" => {
                        tag.child_desc = Some(value.to_string());
                    }
                    "key_desc" => {
                        tag.key_desc = Some(value.to_string());
                    }
                    "value_desc" => {
                        tag.value_desc = Some(value.to_string());
                    }
                    "ele_desc" => {
                        tag.ele_desc = Some(value.to_string());
                    }
                    _ => {}
                }
            }
        }

        Some(tag)
    }

    pub fn to_bytes(&self) -> Vec<u8> {
        let mut parts = Vec::new();

        parts.push(format!("type={}", self.value_type.to_str()));

        if let Some(ref desc) = self.desc {
            parts.push(format!("desc={}", desc));
        }

        if self.nullable {
            parts.push("nullable=true".to_string());
        }

        if let Some(ref default) = self.default {
            parts.push(format!("default={}", default));
        }

        if let Some(ref min) = self.min {
            parts.push(format!("min={}", min));
        }

        if let Some(ref max) = self.max {
            parts.push(format!("max={}", max));
        }

        if let Some(size) = self.size {
            parts.push(format!("size={}", size));
        }

        if let Some(ref enum_values) = self.enum_values {
            parts.push(format!("enum={}", enum_values.join("|")));
        }

        if let Some(ref pattern) = self.pattern {
            parts.push(format!("pattern={}", pattern));
        }

        if let Some(ref location) = self.location {
            parts.push(format!("location={}", location));
        }

        if let Some(ref version) = self.version {
            parts.push(format!("version={}", version));
        }

        if let Some(ref mime) = self.mime {
            parts.push(format!("mime={}", mime));
        }

        if let Some(child_type) = self.child_type {
            parts.push(format!("child_type={}", child_type.to_str()));
        }

        if let Some(ref child_desc) = self.child_desc {
            parts.push(format!("child_desc={}", child_desc));
        }

        if let Some(ref key_desc) = self.key_desc {
            parts.push(format!("key_desc={}", key_desc));
        }

        if let Some(ref value_desc) = self.value_desc {
            parts.push(format!("value_desc={}", value_desc));
        }

        if let Some(ref ele_desc) = self.ele_desc {
            parts.push(format!("ele_desc={}", ele_desc));
        }

        parts.join(";").as_bytes().to_vec()
    }
}

impl Default for Tag {
    fn default() -> Self {
        Self::new()
    }
}