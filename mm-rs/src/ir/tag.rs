use crate::ir::mime::parse_mime;
use crate::ir::value_type::ValueType;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TagKey {
    IsNull = 0 << 3,
    Example = 1 << 3,
    Deprecated = 2 << 3,
    Desc = 3 << 3,
    Type = 4 << 3,
    Nullable = 5 << 3,
    AllowEmpty = 6 << 3,
    Unique = 7 << 3,
    Default = 8 << 3,
    Min = 9 << 3,
    Max = 10 << 3,
    Size = 11 << 3,
    Enum = 12 << 3,
    Pattern = 13 << 3,
    Location = 14 << 3,
    Version = 15 << 3,
    Mime = 16 << 3,
    ChildDesc = 17 << 3,
    ChildType = 18 << 3,
    ChildNullable = 19 << 3,
    ChildAllowEmpty = 20 << 3,
    ChildUnique = 21 << 3,
    ChildDefaultVal = 22 << 3,
    ChildMin = 23 << 3,
    ChildMax = 24 << 3,
    ChildSize = 25 << 3,
    ChildEnums = 26 << 3,
    ChildPattern = 27 << 3,
    ChildLocation = 28 << 3,
    ChildVersion = 29 << 3,
    ChildMime = 30 << 3,
    More = 31 << 3,
}

#[derive(Debug, Clone)]
pub struct Tag {
    pub value_type: ValueType,

    pub name: Option<String>,
    pub is_null: bool,
    pub example: bool,
    pub nullable: bool,
    pub deprecated: bool,
    pub allow_empty: bool,
    pub unique: bool,
    pub is_inherit: bool,

    pub desc: Option<String>,
    pub default_val: Option<String>,
    pub min: Option<String>,
    pub max: Option<String>,
    pub size: Option<u64>,
    pub enums: Option<String>,
    pub pattern: Option<String>,
    pub location: Option<i32>,
    pub version: Option<i32>,
    pub mime: Option<String>,

    pub child_desc: Option<String>,
    pub child_type: ValueType,
    pub child_nullable: bool,
    pub child_allow_empty: bool,
    pub child_unique: bool,
    pub child_default_val: Option<String>,
    pub child_min: Option<String>,
    pub child_max: Option<String>,
    pub child_size: Option<u64>,
    pub child_enums: Option<String>,
    pub child_pattern: Option<String>,
    pub child_location: Option<i32>,
    pub child_version: Option<i32>,
    pub child_mime: Option<String>,

    pub more: u8,
}

impl Tag {
    pub fn new() -> Self {
        Self {
            value_type: ValueType::Unknown,
            name: None,
            is_null: false,
            example: false,
            nullable: false,
            deprecated: false,
            allow_empty: false,
            unique: false,
            is_inherit: false,
            desc: None,
            default_val: None,
            min: None,
            max: None,
            size: None,
            enums: None,
            pattern: None,
            location: None,
            version: None,
            mime: None,
            child_desc: None,
            child_type: ValueType::Unknown,
            child_nullable: false,
            child_allow_empty: false,
            child_unique: false,
            child_default_val: None,
            child_min: None,
            child_max: None,
            child_size: None,
            child_enums: None,
            child_pattern: None,
            child_location: None,
            child_version: None,
            child_mime: None,
            more: 0,
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

        if parent.child_allow_empty {
            self.allow_empty = true;
        }

        if parent.child_unique {
            self.unique = true;
        }

        if let Some(ref v) = parent.child_default_val {
            self.default_val = Some(v.clone());
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

        if let Some(ref v) = parent.child_enums {
            self.enums = Some(v.clone());
            self.value_type = ValueType::Enum;
        }

        if let Some(ref v) = parent.child_pattern {
            self.pattern = Some(v.clone());
        }

        if let Some(v) = parent.child_location {
            if v != 0 {
                self.location = Some(v);
            }
        }

        if let Some(v) = parent.child_version {
            self.version = Some(v);
        }

        if let Some(ref v) = parent.child_mime {
            self.mime = Some(v.clone());
            self.value_type = ValueType::Media;
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

        if src.deprecated {
            dst.deprecated = true;
        }

        if src.more != 0 {
            dst.more = src.more;
        }

        if let Some(ref v) = src.desc {
            dst.desc = Some(v.clone());
        }

        if src.nullable {
            dst.nullable = true;
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

        if let Some(ref v) = src.default_val {
            dst.default_val = Some(v.clone());
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

        if let Some(ref v) = src.enums {
            dst.enums = Some(v.clone());
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

        if let Some(ref v) = src.child_desc {
            dst.child_desc = Some(v.clone());
        }

        if src.child_type != ValueType::Unknown {
            dst.child_type = src.child_type;
        }

        if src.child_nullable {
            dst.child_nullable = true;
        }

        if src.child_allow_empty {
            dst.child_allow_empty = true;
        }

        if src.child_unique {
            dst.child_unique = true;
        }

        if let Some(ref v) = src.child_default_val {
            dst.child_default_val = Some(v.clone());
        }

        if let Some(ref v) = src.child_min {
            dst.child_min = Some(v.clone());
        }

        if let Some(ref v) = src.child_max {
            dst.child_max = Some(v.clone());
        }

        if let Some(v) = src.child_size {
            dst.child_size = Some(v);
        }

        if let Some(ref v) = src.child_enums {
            dst.child_enums = Some(v.clone());
        }

        if let Some(ref v) = src.child_pattern {
            dst.child_pattern = Some(v.clone());
        }

        if let Some(v) = src.child_location {
            dst.child_location = Some(v);
        }

        if let Some(v) = src.child_version {
            dst.child_version = Some(v);
        }

        if let Some(ref v) = src.child_mime {
            dst.child_mime = Some(v.clone());
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
            s
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
                (part[..idx].trim(), Some(part[idx + 1..].trim().to_string()))
            } else {
                (part.trim(), None)
            };

            match key {
                "name" => {
                    tag.name = value.map(|v| v.trim_matches('"').to_string());
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
                "default_val" => {
                    if let Some(ref v) = value {
                        tag.default_val = Some(v.trim_matches('"').to_string());
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
                "enums" => {
                    if let Some(ref v) = value {
                        tag.value_type = ValueType::Enum;
                        tag.enums = Some(v.trim_matches('"').to_string());
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
                        tag.value_type = ValueType::Media;
                    }
                }
                "example" => {
                    if let Some(ref v) = value {
                        tag.example = v == "true";
                    } else {
                        tag.example = true;
                    }
                }
                "is_null" => {
                    if let Some(ref v) = value {
                        tag.is_null = v == "true";
                    } else {
                        tag.is_null = true;
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
                "deprecated" => {
                    if let Some(ref v) = value {
                        tag.deprecated = v == "true";
                    } else {
                        tag.deprecated = true;
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
                    // backward compatibility, no-op
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
                "child_default_val" => {
                    if let Some(ref v) = value {
                        tag.child_default_val = Some(v.trim_matches('"').to_string());
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
                "child_enums" => {
                    if let Some(ref v) = value {
                        tag.child_enums = Some(v.trim_matches('"').to_string());
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
                        tag.child_type = ValueType::Media;
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
                ValueType::Str
                | ValueType::I
                | ValueType::F64
                | ValueType::Bool
                | ValueType::Obj
                | ValueType::Vec => {}
                ValueType::Arr => if self.size.is_none() || self.size.unwrap_or(0) == 0 {},
                ValueType::Enum => if self.enums.is_some() {},
                ValueType::Media => if self.mime.is_some() {},
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

        if self.deprecated && !self.is_inherit {
            parts.push("deprecated".to_string());
        }

        if self.allow_empty && !self.is_inherit {
            parts.push("allow_empty".to_string());
        }

        if self.unique && !self.is_inherit {
            parts.push("unique".to_string());
        }

        if let Some(ref v) = self.default_val {
            if !self.is_inherit {
                parts.push(format!("default_val={}", v));
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

        if let Some(ref v) = self.enums {
            if !self.is_inherit {
                parts.push(format!("enums={}", v));
            }
        }

        if let Some(ref v) = self.pattern {
            if !self.is_inherit {
                parts.push(format!("pattern={}", v));
            }
        }

        if let Some(v) = self.location {
            if !self.is_inherit && v != 0 {
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
                ValueType::Str
                | ValueType::I
                | ValueType::F64
                | ValueType::Bool
                | ValueType::Obj
                | ValueType::Vec => {}
                ValueType::Arr => {
                    if self.child_size.is_none() || self.child_size.unwrap_or(0) == 0 {}
                }
                ValueType::Enum => if self.child_enums.is_some() {},
                ValueType::Media => if self.child_mime.is_some() {},
                _ => {
                    parts.push(format!("child_type={}", self.child_type.to_str()));
                }
            }
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

        if let Some(ref v) = self.child_default_val {
            parts.push(format!("child_default_val={}", v));
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

        if let Some(ref v) = self.child_enums {
            parts.push(format!("child_enums={}", v));
        }

        if let Some(ref v) = self.child_pattern {
            parts.push(format!("child_pattern={}", v));
        }

        if let Some(v) = self.child_location {
            if v != 0 {
                parts.push(format!("child_location={}", v));
            }
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
        let mut bs = Vec::new();

        if self.example {
            bs.push(TagKey::Example as u8 | 1);
        }

        if self.is_null {
            bs.push(TagKey::IsNull as u8 | 1);
        }

        if self.nullable && !self.is_inherit {
            if !self.is_null {
                bs.push(TagKey::Nullable as u8 | 1);
            }
        }

        if self.deprecated && !self.is_inherit {
            bs.push(TagKey::Deprecated as u8 | 1);
        }

        if let Some(ref desc) = self.desc {
            if !self.is_inherit {
                let l = desc.len();
                if l <= 5 {
                    bs.push(TagKey::Desc as u8 | l as u8);
                    bs.extend_from_slice(desc.as_bytes());
                } else if l <= 255 {
                    bs.push(TagKey::Desc as u8 | 6);
                    bs.push(l as u8);
                    bs.extend_from_slice(desc.as_bytes());
                } else if l <= 65535 {
                    bs.push(TagKey::Desc as u8 | 7);
                    bs.push((l >> 8) as u8);
                    bs.push(l as u8);
                    bs.extend_from_slice(desc.as_bytes());
                }
            }
        }

        if self.value_type != ValueType::Unknown && !self.is_inherit {
            match self.value_type {
                ValueType::Str
                | ValueType::Bytes
                | ValueType::I
                | ValueType::F64
                | ValueType::Bool
                | ValueType::Obj
                | ValueType::Vec => {}
                ValueType::Arr => if self.size.is_some() && self.size.unwrap_or(0) > 0 {},
                ValueType::Enum => if self.enums.is_some() {},
                ValueType::Media => if self.mime.is_some() {},
                _ => {
                    bs.push(TagKey::Type as u8);
                    bs.push(self.value_type as u8);
                }
            }
        }

        if self.allow_empty && !self.is_inherit {
            bs.push(TagKey::AllowEmpty as u8 | 1);
        }

        if self.unique && !self.is_inherit {
            bs.push(TagKey::Unique as u8 | 1);
        }

        if let Some(ref v) = self.default_val {
            if !self.is_inherit {
                let l = v.len();
                if l < 7 {
                    bs.push(TagKey::Default as u8 | l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else {
                    bs.push(TagKey::Default as u8 | 7);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                }
            }
        }

        if let Some(ref v) = self.min {
            if !self.is_inherit {
                let l = v.len();
                if l < 7 {
                    bs.push(TagKey::Min as u8 | l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else {
                    bs.push(TagKey::Min as u8 | 7);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                }
            }
        }

        if let Some(ref v) = self.max {
            if !self.is_inherit {
                let l = v.len();
                if l < 7 {
                    bs.push(TagKey::Max as u8 | l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else {
                    bs.push(TagKey::Max as u8 | 7);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                }
            }
        }

        if let Some(v) = self.size {
            if !self.is_inherit {
                Self::encode_u64(&mut bs, TagKey::Size as u8, v);
            }
        }

        if let Some(ref v) = self.enums {
            if !self.is_inherit {
                let l = v.len();
                if l <= 5 {
                    bs.push(TagKey::Enum as u8 | l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else if l <= 255 {
                    bs.push(TagKey::Enum as u8 | 6);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else if l <= 65535 {
                    bs.push(TagKey::Enum as u8 | 7);
                    bs.push((l >> 8) as u8);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                }
            }
        }

        if let Some(ref v) = self.pattern {
            if !self.is_inherit {
                let l = v.len();
                if l < 7 {
                    bs.push(TagKey::Pattern as u8 | l as u8);
                    bs.extend_from_slice(v.as_bytes());
                } else {
                    bs.push(TagKey::Pattern as u8 | 7);
                    bs.push(l as u8);
                    bs.extend_from_slice(v.as_bytes());
                }
            }
        }

        if let Some(v) = self.location {
            if !self.is_inherit && v != 0 {
                let s = v.to_string();
                bs.push(TagKey::Location as u8 | s.len() as u8);
                bs.extend_from_slice(s.as_bytes());
            }
        }

        if let Some(v) = self.version {
            if !self.is_inherit {
                Self::encode_u64(&mut bs, TagKey::Version as u8, v as u64);
            }
        }

        if let Some(ref v) = self.mime {
            if !self.is_inherit {
                let mime_id = parse_mime(v);
                Self::encode_u64(&mut bs, TagKey::Mime as u8, mime_id as u64);
            }
        }

        if let Some(ref v) = self.child_desc {
            let l = v.len();
            if l <= 5 {
                bs.push(TagKey::ChildDesc as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else if l <= 255 {
                bs.push(TagKey::ChildDesc as u8 | 6);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else if l <= 65535 {
                bs.push(TagKey::ChildDesc as u8 | 7);
                bs.push((l >> 8) as u8);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if self.child_type != ValueType::Unknown {
            match self.child_type {
                ValueType::Str
                | ValueType::I
                | ValueType::F64
                | ValueType::Bool
                | ValueType::Obj
                | ValueType::Vec => {}
                ValueType::Arr => {
                    if self.child_size.is_some() && self.child_size.unwrap_or(0) > 0 {}
                }
                ValueType::Enum => if self.child_enums.is_some() {},
                ValueType::Media => if self.child_mime.is_some() {},
                _ => {
                    bs.push(TagKey::ChildType as u8);
                    bs.push(self.child_type as u8);
                }
            }
        }

        if self.child_nullable {
            bs.push(TagKey::ChildNullable as u8 | 1);
        }

        if self.child_allow_empty {
            bs.push(TagKey::ChildAllowEmpty as u8 | 1);
        }

        if self.child_unique {
            bs.push(TagKey::ChildUnique as u8 | 1);
        }

        if let Some(ref v) = self.child_default_val {
            let l = v.len();
            if l < 7 {
                bs.push(TagKey::ChildDefaultVal as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else {
                bs.push(TagKey::ChildDefaultVal as u8 | 7);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if let Some(ref v) = self.child_min {
            let l = v.len();
            if l < 7 {
                bs.push(TagKey::ChildMin as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else {
                bs.push(TagKey::ChildMin as u8 | 7);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if let Some(ref v) = self.child_max {
            let l = v.len();
            if l < 7 {
                bs.push(TagKey::ChildMax as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else {
                bs.push(TagKey::ChildMax as u8 | 7);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if let Some(v) = self.child_size {
            Self::encode_u64(&mut bs, TagKey::ChildSize as u8, v);
        }

        if let Some(ref v) = self.child_enums {
            let l = v.len();
            if l <= 5 {
                bs.push(TagKey::ChildEnums as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else if l <= 255 {
                bs.push(TagKey::ChildEnums as u8 | 6);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else if l <= 65535 {
                bs.push(TagKey::ChildEnums as u8 | 7);
                bs.push((l >> 8) as u8);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if let Some(ref v) = self.child_pattern {
            let l = v.len();
            if l < 7 {
                bs.push(TagKey::ChildPattern as u8 | l as u8);
                bs.extend_from_slice(v.as_bytes());
            } else {
                bs.push(TagKey::ChildPattern as u8 | 7);
                bs.push(l as u8);
                bs.extend_from_slice(v.as_bytes());
            }
        }

        if let Some(v) = self.child_location {
            if v != 0 {
                let s = v.to_string();
                bs.push(TagKey::ChildLocation as u8 | s.len() as u8);
                bs.extend_from_slice(s.as_bytes());
            }
        }

        if let Some(v) = self.child_version {
            Self::encode_u64(&mut bs, TagKey::ChildVersion as u8, v as u64);
        }

        if let Some(ref v) = self.child_mime {
            let mime_id = parse_mime(v);
            Self::encode_u64(&mut bs, TagKey::ChildMime as u8, mime_id as u64);
        }

        if self.more != 0 {
            Self::encode_u64(&mut bs, TagKey::More as u8, self.more as u64);
        }

        bs
    }

    fn encode_u64(buf: &mut Vec<u8>, sign: u8, uv: u64) {
        match uv {
            0..=0xFF => {
                buf.push(sign | 0);
                buf.push(uv as u8);
            }
            0x100..=0xFFFF => {
                buf.push(sign | 1);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            0x10000..=0xFFFFFF => {
                buf.push(sign | 2);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            0x1000000..=0xFFFFFFFF => {
                buf.push(sign | 3);
                buf.push((uv >> 24) as u8);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            0x100000000..=0xFFFFFFFFFF => {
                buf.push(sign | 4);
                buf.push((uv >> 32) as u8);
                buf.push((uv >> 24) as u8);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            0x10000000000..=0xFFFFFFFFFFFF => {
                buf.push(sign | 5);
                buf.push((uv >> 40) as u8);
                buf.push((uv >> 32) as u8);
                buf.push((uv >> 24) as u8);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            0x1000000000000..=0xFFFFFFFFFFFFFF => {
                buf.push(sign | 6);
                buf.push((uv >> 48) as u8);
                buf.push((uv >> 40) as u8);
                buf.push((uv >> 32) as u8);
                buf.push((uv >> 24) as u8);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
            _ => {
                buf.push(sign | 7);
                buf.push((uv >> 56) as u8);
                buf.push((uv >> 48) as u8);
                buf.push((uv >> 40) as u8);
                buf.push((uv >> 32) as u8);
                buf.push((uv >> 24) as u8);
                buf.push((uv >> 16) as u8);
                buf.push((uv >> 8) as u8);
                buf.push(uv as u8);
            }
        }
    }
}

impl Default for Tag {
    fn default() -> Self {
        Self::new()
    }
}
