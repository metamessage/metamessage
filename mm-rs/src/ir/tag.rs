use crate::ir::value_type::ValueType;

#[derive(Debug, Clone)]
pub struct Tag {
    pub value_type: ValueType,

    pub name: bool,
    pub is_null: bool,
    pub example: bool,
    pub nullable: bool,
    pub raw: bool,
    pub allow_empty: bool,
    pub unique: bool,
    pub is_inherit: bool,

    pub desc: Option<String>,
    pub default: Option<String>,
    pub min: Option<String>,
    pub max: Option<String>,
    pub size: Option<u64>,
    pub enum_values: Option<String>,
    pub pattern: Option<String>,
    pub location: Option<i32>,
    pub version: Option<i32>,
    pub mime: Option<String>,

    pub child_desc: Option<String>,
    pub child_type: ValueType,
    pub child_nullable: bool,
    pub child_raw: bool,
    pub child_allow_empty: bool,
    pub child_unique: bool,
    pub child_default: Option<String>,
    pub child_min: Option<String>,
    pub child_max: Option<String>,
    pub child_size: Option<u64>,
    pub child_enum: Option<String>,
    pub child_pattern: Option<String>,
    pub child_location: Option<i32>,
    pub child_version: Option<i32>,
    pub child_mime: Option<String>,
}

impl Tag {
    pub fn new() -> Self {
        Self {
            value_type: ValueType::Unknown,
            name: false,
            is_null: false,
            example: false,
            nullable: false,
            raw: false,
            allow_empty: false,
            unique: false,
            is_inherit: false,
            desc: None,
            default: None,
            min: None,
            max: None,
            size: None,
            enum_values: None,
            pattern: None,
            location: None,
            version: None,
            mime: None,
            child_desc: None,
            child_type: ValueType::Unknown,
            child_nullable: false,
            child_raw: false,
            child_allow_empty: false,
            child_unique: false,
            child_default: None,
            child_min: None,
            child_max: None,
            child_size: None,
            child_enum: None,
            child_pattern: None,
            child_location: None,
            child_version: None,
            child_mime: None,
        }
    }

    pub fn inherit(&mut self, parent: &Tag) {
        self.is_inherit = true;

        if parent.child_type != ValueType::Unknown {
            self.value_type = parent.child_type;
        }

        if let Some(ref v) = parent.child_desc {
            self.desc = Some(v.clone());
        }

        if parent.child_nullable {
            self.nullable = true;
        }

        if parent.child_raw {
            self.raw = true;
        }

        if parent.child_allow_empty {
            self.allow_empty = true;
        }

        if parent.child_unique {
            self.unique = true;
        }

        if let Some(ref v) = parent.child_default {
            self.default = Some(v.clone());
        }

        if let Some(ref v) = parent.child_min {
            self.min = Some(v.clone());
        }

        if let Some(ref v) = parent.child_max {
            self.max = Some(v.clone());
        }

        if let Some(v) = parent.child_size {
            self.size = Some(v);
        }

        if let Some(ref v) = parent.child_enum {
            self.enum_values = Some(v.clone());
            self.value_type = ValueType::Enum;
        }

        if let Some(ref v) = parent.child_pattern {
            self.pattern = Some(v.clone());
        }

        if let Some(v) = parent.child_location {
            self.location = Some(v);
        }

        if let Some(v) = parent.child_version {
            self.version = Some(v);
        }

        if let Some(ref v) = parent.child_mime {
            self.mime = Some(v.clone());
        }
    }

    pub fn merge(dst: Option<Tag>, src: Tag) -> Tag {
        let mut dst = dst.unwrap_or_default();

        if src.is_null {
            dst.is_null = true;
        }

        if src.example {
            dst.example = true;
        }

        if let Some(ref v) = src.desc {
            dst.desc = Some(v.clone());
        }

        if src.nullable {
            dst.nullable = true;
        }

        if src.raw {
            dst.raw = true;
        }

        if src.allow_empty {
            dst.allow_empty = true;
        }

        if src.unique {
            dst.unique = true;
        }

        if src.value_type != ValueType::Unknown {
            dst.value_type = src.value_type;
        }

        if let Some(ref v) = src.default {
            dst.default = Some(v.clone());
        }

        if let Some(ref v) = src.min {
            dst.min = Some(v.clone());
        }

        if let Some(ref v) = src.max {
            dst.max = Some(v.clone());
        }

        if let Some(v) = src.size {
            dst.size = Some(v);
        }

        if let Some(ref v) = src.enum_values {
            dst.enum_values = Some(v.clone());
        }

        if let Some(ref v) = src.pattern {
            dst.pattern = Some(v.clone());
        }

        if let Some(v) = src.location {
            dst.location = Some(v);
        }

        if let Some(v) = src.version {
            dst.version = Some(v);
        }

        if let Some(ref v) = src.mime {
            dst.mime = Some(v.clone());
        }

        dst
    }

    pub fn parse(s: &str) -> Option<Tag> {
        let s = s.trim();

        let s = if let Some(rest) = s.strip_prefix("//") {
            rest
        } else if let Some(rest) = s.strip_prefix("/*") {
            rest.strip_suffix("*/")?
        } else {
            return None;
        };

        let s = s.trim();

        let s = s.strip_prefix("mm:")?;

        let s = s.trim();
        if s.is_empty() {
            return Some(Tag::new());
        }

        let mut tag = Tag::new();

        for part in s.split(';') {
            let part = part.trim();
            if part.is_empty() {
                continue;
            }

            let (key, value) = if let Some(idx) = part.find('=') {
                (
                    part[..idx].trim(),
                    Some(part[idx + 1..].trim().to_string()),
                )
            } else {
                (part.trim(), None)
            };

            match key {
                "name" => {
                    tag.name = true;
                    let _v = value;
                }
                "type" => {
                    if let Some(ref v) = value {
                        tag.value_type = ValueType::from_str(v);
                    }
                }
                "desc" => {
                    if let Some(ref v) = value {
                        tag.desc = Some(v.trim_matches('"').to_string());
                    }
                }
                "default" => {
                    if let Some(ref v) = value {
                        tag.default = Some(v.trim_matches('"').to_string());
                    }
                }
                "min" => {
                    if let Some(ref v) = value {
                        tag.min = Some(v.trim_matches('"').to_string());
                    }
                }
                "max" => {
                    if let Some(ref v) = value {
                        tag.max = Some(v.trim_matches('"').to_string());
                    }
                }
                "size" => {
                    if let Some(ref v) = value {
                        tag.size = v.parse::<u64>().ok();
                    }
                }
                "enum" => {
                    if let Some(ref v) = value {
                        tag.value_type = ValueType::Enum;
                        tag.enum_values = Some(v.trim_matches('"').to_string());
                    }
                }
                "location" => {
                    if let Some(ref v) = value {
                        if let Ok(offset) = v.parse::<i32>() {
                            if (-12..=14).contains(&offset) {
                                tag.location = Some(offset);
                            }
                        }
                    }
                }
                "version" => {
                    if let Some(ref v) = value {
                        if let Ok(ver) = v.parse::<i32>() {
                            if (1..=10).contains(&ver) {
                                tag.version = Some(ver);
                            }
                        }
                    }
                }
                "mime" => {
                    if let Some(ref v) = value {
                        tag.mime = Some(v.clone());
                    }
                }
                "nullable" => {
                    if let Some(ref v) = value {
                        tag.nullable = v == "true";
                    } else {
                        tag.nullable = true;
                    }
                }
                "allow_empty" => {
                    if let Some(ref v) = value {
                        tag.allow_empty = v == "true";
                    } else {
                        tag.allow_empty = true;
                    }
                }
                "raw" => {
                    if let Some(ref v) = value {
                        tag.raw = v == "true";
                    } else {
                        tag.raw = true;
                    }
                }
                "unique" => {
                    if let Some(ref v) = value {
                        tag.unique = v == "true";
                    } else {
                        tag.unique = true;
                    }
                }
                "pattern" => {
                    if let Some(ref v) = value {
                        tag.pattern = Some(v.trim_matches('"').to_string());
                    }
                }
                "child_desc" => {
                    if let Some(ref v) = value {
                        tag.child_desc = Some(v.trim_matches('"').to_string());
                    }
                }
                "child_type" => {
                    if let Some(ref v) = value {
                        tag.child_type = ValueType::from_str(v);
                    }
                }
                "child_nullable" => {
                    if let Some(ref v) = value {
                        tag.child_nullable = v == "true";
                    } else {
                        tag.child_nullable = true;
                    }
                }
                "child_raw" => {
                    if let Some(ref v) = value {
                        tag.child_raw = v == "true";
                    } else {
                        tag.child_raw = true;
                    }
                }
                "child_allow_empty" => {
                    if let Some(ref v) = value {
                        tag.child_allow_empty = v == "true";
                    } else {
                        tag.child_allow_empty = true;
                    }
                }
                "child_unique" => {
                    if let Some(ref v) = value {
                        tag.child_unique = v == "true";
                    } else {
                        tag.child_unique = true;
                    }
                }
                "child_default" => {
                    if let Some(ref v) = value {
                        tag.child_default = Some(v.trim_matches('"').to_string());
                    }
                }
                "child_min" => {
                    if let Some(ref v) = value {
                        tag.child_min = Some(v.trim_matches('"').to_string());
                    }
                }
                "child_max" => {
                    if let Some(ref v) = value {
                        tag.child_max = Some(v.trim_matches('"').to_string());
                    }
                }
                "child_size" => {
                    if let Some(ref v) = value {
                        tag.child_size = v.parse::<u64>().ok();
                    }
                }
                "child_enum" => {
                    if let Some(ref v) = value {
                        tag.child_enum = Some(v.trim_matches('"').to_string());
                        tag.child_type = ValueType::Enum;
                    }
                }
                "child_location" => {
                    if let Some(ref v) = value {
                        if let Ok(offset) = v.parse::<i32>() {
                            if (-12..=14).contains(&offset) {
                                tag.child_location = Some(offset);
                            }
                        }
                    }
                }
                "child_version" => {
                    if let Some(ref v) = value {
                        if let Ok(ver) = v.parse::<i32>() {
                            if (1..=10).contains(&ver) {
                                tag.child_version = Some(ver);
                            }
                        }
                    }
                }
                "child_mime" => {
                    if let Some(ref v) = value {
                        tag.child_mime = Some(v.clone());
                    }
                }
                _ => {}
            }
        }

        Some(tag)
    }

    #[allow(clippy::inherent_to_string)]
    pub fn to_string(&self) -> String {
        let mut parts = Vec::new();

        if self.value_type != ValueType::Unknown && !self.is_inherit {
            match self.value_type {
                ValueType::String
                | ValueType::Int
                | ValueType::Float64
                | ValueType::Bool
                | ValueType::Struct
                | ValueType::Slice => {}
                ValueType::Array => if self.size.is_none() || self.size.unwrap_or(0) == 0 {},
                ValueType::Enum => if self.enum_values.is_some() {},
                _ => {
                    parts.push(format!("type={}", self.value_type.to_str()));
                }
            }
        }

        if self.example {
            parts.push("example".to_string());
        }

        if self.is_null {
            parts.push("is_null".to_string());
        }

        if self.nullable && !self.is_inherit && !self.is_null {
            parts.push("nullable".to_string());
        }

        if let Some(ref desc) = self.desc {
            if !self.is_inherit {
                parts.push(format!("desc=\"{}\"", desc));
            }
        }

        if self.raw && !self.is_inherit {
            parts.push("raw".to_string());
        }

        if self.allow_empty && !self.is_inherit {
            parts.push("allow_empty".to_string());
        }

        if self.unique && !self.is_inherit {
            parts.push("unique".to_string());
        }

        if let Some(ref v) = self.default {
            if !self.is_inherit {
                parts.push(format!("default={}", v));
            }
        }

        if let Some(ref v) = self.min {
            if !self.is_inherit {
                parts.push(format!("min={}", v));
            }
        }

        if let Some(ref v) = self.max {
            if !self.is_inherit {
                parts.push(format!("max={}", v));
            }
        }

        if let Some(v) = self.size {
            if !self.is_inherit {
                parts.push(format!("size={}", v));
            }
        }

        if let Some(ref v) = self.enum_values {
            if !self.is_inherit {
                parts.push(format!("enum={}", v));
            }
        }

        if let Some(ref v) = self.pattern {
            if !self.is_inherit {
                parts.push(format!("pattern={}", v));
            }
        }

        if let Some(v) = self.location {
            if !self.is_inherit {
                parts.push(format!("location={}", v));
            }
        }

        if let Some(v) = self.version {
            if !self.is_inherit {
                parts.push(format!("version={}", v));
            }
        }

        if let Some(ref v) = self.mime {
            if !self.is_inherit {
                parts.push(format!("mime={}", v));
            }
        }

        if let Some(ref v) = self.child_desc {
            parts.push(format!("child_desc=\"{}\"", v));
        }

        if self.child_type != ValueType::Unknown {
            match self.child_type {
                ValueType::String
                | ValueType::Int
                | ValueType::Float64
                | ValueType::Bool
                | ValueType::Struct
                | ValueType::Slice => {}
                ValueType::Array => {
                    if self.child_size.is_none() || self.child_size.unwrap_or(0) == 0 {}
                }
                ValueType::Enum => if self.child_enum.is_some() {},
                _ => {
                    parts.push(format!("child_type={}", self.child_type.to_str()));
                }
            }
        }

        if self.child_raw {
            parts.push("child_raw".to_string());
        }

        if self.child_nullable {
            parts.push("child_nullable".to_string());
        }

        if self.child_allow_empty {
            parts.push("child_allow_empty".to_string());
        }

        if self.child_unique {
            parts.push("child_unique".to_string());
        }

        if let Some(ref v) = self.child_default {
            parts.push(format!("child_default={}", v));
        }

        if let Some(ref v) = self.child_min {
            parts.push(format!("child_min={}", v));
        }

        if let Some(ref v) = self.child_max {
            parts.push(format!("child_max={}", v));
        }

        if let Some(v) = self.child_size {
            parts.push(format!("child_size={}", v));
        }

        if let Some(ref v) = self.child_enum {
            parts.push(format!("child_enum={}", v));
        }

        if let Some(ref v) = self.child_pattern {
            parts.push(format!("child_pattern={}", v));
        }

        if let Some(v) = self.child_location {
            parts.push(format!("child_location={}", v));
        }

        if let Some(v) = self.child_version {
            parts.push(format!("child_version={}", v));
        }

        if let Some(ref v) = self.child_mime {
            parts.push(format!("child_mime={}", v));
        }

        if parts.is_empty() {
            String::new()
        } else {
            parts.join("; ")
        }
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

        if let Some(ref v) = self.default {
            parts.push(format!("default={}", v));
        }

        if let Some(ref v) = self.min {
            parts.push(format!("min={}", v));
        }

        if let Some(ref v) = self.max {
            parts.push(format!("max={}", v));
        }

        if let Some(v) = self.size {
            parts.push(format!("size={}", v));
        }

        if let Some(ref v) = self.enum_values {
            parts.push(format!("enum={}", v));
        }

        if let Some(ref v) = self.pattern {
            parts.push(format!("pattern={}", v));
        }

        if let Some(v) = self.location {
            parts.push(format!("location={}", v));
        }

        if let Some(v) = self.version {
            parts.push(format!("version={}", v));
        }

        if let Some(ref v) = self.mime {
            parts.push(format!("mime={}", v));
        }

        if let Some(ref v) = self.child_desc {
            parts.push(format!("child_desc={}", v));
        }

        if self.child_type != ValueType::Unknown {
            parts.push(format!("child_type={}", self.child_type.to_str()));
        }

        if self.child_raw {
            parts.push("child_raw=true".to_string());
        }

        if self.child_nullable {
            parts.push("child_nullable=true".to_string());
        }

        if let Some(ref v) = self.child_default {
            parts.push(format!("child_default={}", v));
        }

        if let Some(ref v) = self.child_min {
            parts.push(format!("child_min={}", v));
        }

        if let Some(ref v) = self.child_max {
            parts.push(format!("child_max={}", v));
        }

        if let Some(v) = self.child_size {
            parts.push(format!("child_size={}", v));
        }

        if let Some(ref v) = self.child_enum {
            parts.push(format!("child_enum={}", v));
        }

        if let Some(ref v) = self.child_pattern {
            parts.push(format!("child_pattern={}", v));
        }

        if let Some(v) = self.child_location {
            parts.push(format!("child_location={}", v));
        }

        if let Some(v) = self.child_version {
            parts.push(format!("child_version={}", v));
        }

        if let Some(ref v) = self.child_mime {
            parts.push(format!("child_mime={}", v));
        }

        parts.join(";").as_bytes().to_vec()
    }
}

impl Default for Tag {
    fn default() -> Self {
        Self::new()
    }
}