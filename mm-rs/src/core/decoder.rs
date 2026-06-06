use crate::core::constants::{
    BYTES_LEN_1, CONTAINER_LEN_1, CONTAINER_LEN_MASK, INT_LEN_1, INT_LEN_MASK, STRING_LEN_1,
    TAG_ALLOW_EMPTY, TAG_CHILD_ALLOW_EMPTY, TAG_CHILD_DEFAULT_VAL, TAG_CHILD_DESC, TAG_CHILD_ENUMS,
    TAG_CHILD_LOCATION, TAG_CHILD_MAX, TAG_CHILD_MIME, TAG_CHILD_MIN, TAG_CHILD_NULLABLE,
    TAG_CHILD_PATTERN, TAG_CHILD_SIZE, TAG_CHILD_TYPE, TAG_CHILD_UNIQUE, TAG_CHILD_VERSION,
    TAG_DEFAULT_VAL, TAG_DEPRECATED, TAG_DESC, TAG_ENUMS, TAG_EXAMPLE, TAG_IS_NULL, TAG_KEY_MASK,
    TAG_LEN_1, TAG_LOCATION, TAG_MAX, TAG_MIME, TAG_MIN, TAG_MORE, TAG_NULLABLE, TAG_PATTERN,
    TAG_PAYLOAD_MASK, TAG_SIZE, TAG_TYPE, TAG_UNIQUE, TAG_VERSION,
};
use crate::core::prefix::{
    CONTAINER_TYPE_MASK, FLOAT_LEN_1, FLOAT_LEN_MASK, FLOAT_POSITIVE_NEGATIVE_MASK, PREFIX_BYTES,
    PREFIX_CONTAINER, PREFIX_FLOAT, PREFIX_MASK, PREFIX_NEGATIVE_INT, PREFIX_POSITIVE_INT,
    PREFIX_SIMPLE, PREFIX_STRING, PREFIX_TAG,
};
use crate::core::simple_value::SimpleValue;
use crate::ir::ast::{Array, Field, Node, Object, Value, ValueData};
use crate::ir::mime::mime_to_str;
use crate::ir::tag::Tag;
use crate::ir::ValueType;
use chrono::TimeZone;

pub struct Decoder {
    data: Vec<u8>,
    offset: usize,
}

impl Decoder {
    pub fn new(data: Vec<u8>) -> Self {
        Self { data, offset: 0 }
    }

    fn read_byte(&mut self) -> Result<u8, std::io::Error> {
        if self.offset >= self.data.len() {
            return Err(std::io::Error::new(
                std::io::ErrorKind::UnexpectedEof,
                "unexpected eof",
            ));
        }
        let b = self.data[self.offset];
        self.offset += 1;
        Ok(b)
    }

    fn read_bytes(&mut self, n: usize) -> Result<&[u8], std::io::Error> {
        if self.offset + n > self.data.len() {
            return Err(std::io::Error::new(
                std::io::ErrorKind::UnexpectedEof,
                "unexpected eof",
            ));
        }
        let start = self.offset;
        self.offset += n;
        Ok(&self.data[start..start + n])
    }

    pub fn decode(&mut self) -> Result<Node, std::io::Error> {
        let tag = Tag::new();
        self.decode_node(&tag, "")
    }

    fn decode_node(&mut self, tag: &Tag, _path: &str) -> Result<Node, std::io::Error> {
        let b = self.read_byte()?;
        match b & PREFIX_MASK {
            PREFIX_TAG => self.decode_tag(b),
            PREFIX_SIMPLE => self.decode_simple(b, tag),
            PREFIX_POSITIVE_INT => self.decode_positive_int(b, tag),
            PREFIX_NEGATIVE_INT => self.decode_negative_int(b, tag),
            PREFIX_FLOAT => self.decode_float(b, tag),
            PREFIX_STRING => self.decode_string(b, tag),
            PREFIX_BYTES => self.decode_bytes(b, tag),
            PREFIX_CONTAINER => self.decode_container(b, tag),
            _ => Err(std::io::Error::new(
                std::io::ErrorKind::InvalidData,
                "invalid prefix",
            )),
        }
    }

    fn decode_tag(&mut self, prefix: u8) -> Result<Node, std::io::Error> {
        let _total_len = {
            let l = prefix & INT_LEN_MASK;
            if l < TAG_LEN_1 {
                l as usize
            } else if l == TAG_LEN_1 {
                self.read_byte()? as usize
            } else {
                let b = self.read_bytes(2)?;
                ((b[0] as usize) << 8) | (b[1] as usize)
            }
        };

        let mut tag = Tag::new();

        let tag_len_byte = self.read_byte()?;
        let tag_bytes_len = if tag_len_byte < 254 {
            tag_len_byte as usize
        } else if tag_len_byte == 254 {
            self.read_byte()? as usize
        } else {
            let b = self.read_bytes(2)?;
            ((b[0] as usize) << 8) | (b[1] as usize)
        };

        let mut remaining = tag_bytes_len;
        while remaining > 0 {
            let n = self.decode_tag_bytes(&mut tag)?;
            if n == 0 || n > remaining {
                break;
            }
            remaining -= n;
        }

        if tag.is_null {
            let inner = self.decode_node(&Tag::new(), "")?;
            let (data, text) = match &inner {
                Node::Value(v) => (v.data.clone(), v.text.clone()),
                _ => (ValueData::Null, String::new()),
            };
            Ok(Node::Value(Value {
                data,
                text,
                path: String::new(),
                tag: Some(tag),
            }))
        } else {
            let mut inner_tag = Tag::new();
            inner_tag.inherit(&tag);
            if tag.child_type == ValueType::Unknown {
                inner_tag.value_type = tag.value_type;
                if tag.enums.is_some() {
                    inner_tag.enums = tag.enums.clone();
                }
                if tag.mime.is_some() {
                    inner_tag.mime = tag.mime.clone();
                }
                if tag.version.is_some() {
                    inner_tag.version = tag.version;
                }
            }
            if inner_tag.location.is_none() {
                inner_tag.location = tag.location;
            }
            if inner_tag.version.is_none() {
                inner_tag.version = tag.version;
            }
            self.decode_node(&inner_tag, "")
        }
    }

    fn decode_tag_bytes(&mut self, tag: &mut Tag) -> Result<usize, std::io::Error> {
        let b = self.read_byte()?;
        let key = b & TAG_KEY_MASK;
        let payload = (b & TAG_PAYLOAD_MASK) as usize;

        match key {
            TAG_IS_NULL => {
                tag.is_null = (payload & 1) == 1;
                if tag.is_null {
                    tag.nullable = true;
                }
                Ok(1)
            }
            TAG_EXAMPLE => {
                tag.example = (payload & 1) == 1;
                Ok(1)
            }
            TAG_DESC => {
                let s = self.read_tag_str(payload)?;
                tag.desc = Some(s.clone());
                let n = if payload <= 5 {
                    1 + payload
                } else if payload == 6 {
                    2 + s.len()
                } else {
                    3 + s.len()
                };
                Ok(n)
            }
            TAG_TYPE => {
                let tb = self.read_byte()?;
                tag.value_type = ValueType::from_code(tb);
                Ok(2)
            }
            TAG_DEPRECATED => {
                tag.deprecated = (payload & 1) == 1;
                Ok(1)
            }
            TAG_NULLABLE => {
                tag.nullable = (payload & 1) == 1;
                Ok(1)
            }
            TAG_ALLOW_EMPTY => {
                tag.allow_empty = (payload & 1) == 1;
                Ok(1)
            }
            TAG_UNIQUE => {
                tag.unique = (payload & 1) == 1;
                Ok(1)
            }
            TAG_DEFAULT_VAL => {
                let s = self.read_tag_short_str(payload)?;
                tag.default_val = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_MIN => {
                let s = self.read_tag_short_str(payload)?;
                tag.min = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_MAX => {
                let s = self.read_tag_short_str(payload)?;
                tag.max = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_SIZE => {
                let v = self.read_tag_uint(payload)?;
                tag.size = Some(v);
                Ok(2 + payload)
            }
            TAG_ENUMS => {
                tag.value_type = ValueType::Enum;
                let s = self.read_tag_str(payload)?;
                tag.enums = Some(s.clone());
                let n = if payload <= 5 {
                    1 + payload
                } else if payload == 6 {
                    2 + s.len()
                } else {
                    3 + s.len()
                };
                Ok(n)
            }
            TAG_PATTERN => {
                let s = self.read_tag_short_str(payload)?;
                tag.pattern = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_LOCATION => {
                let s = self.read_tag_ascii(payload)?;
                if let Ok(n) = s.parse::<i32>() {
                    tag.location = Some(n);
                }
                Ok(1 + payload)
            }
            TAG_VERSION => {
                let v = self.read_tag_uint(payload)?;
                tag.version = Some(v as i32);
                Ok(2 + payload)
            }
            TAG_MIME => {
                tag.value_type = ValueType::Media;
                let mime_id = if payload < 7 {
                    payload as u8
                } else {
                    self.read_byte()?
                };
                tag.mime = Some(mime_to_str(mime_id).to_string());
                let extra = if payload < 7 { 0 } else { 1 };
                Ok(1 + extra)
            }
            TAG_CHILD_DESC => {
                let s = self.read_tag_str(payload)?;
                tag.child_desc = Some(s.clone());
                let n = if payload <= 5 {
                    1 + payload
                } else if payload == 6 {
                    2 + s.len()
                } else {
                    3 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_TYPE => {
                let tb = self.read_byte()?;
                tag.child_type = ValueType::from_code(tb);
                Ok(2)
            }
            TAG_MORE => {
                tag.more = (payload & 0x07) as u8;
                Ok(1)
            }
            TAG_CHILD_NULLABLE => {
                tag.child_nullable = (payload & 1) == 1;
                Ok(1)
            }
            TAG_CHILD_ALLOW_EMPTY => {
                tag.child_allow_empty = (payload & 1) == 1;
                Ok(1)
            }
            TAG_CHILD_UNIQUE => {
                tag.child_unique = (payload & 1) == 1;
                Ok(1)
            }
            TAG_CHILD_DEFAULT_VAL => {
                let s = self.read_tag_short_str(payload)?;
                tag.child_default_val = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_MIN => {
                let s = self.read_tag_short_str(payload)?;
                tag.child_min = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_MAX => {
                let s = self.read_tag_short_str(payload)?;
                tag.child_max = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_SIZE => {
                let v = self.read_tag_uint(payload)?;
                tag.child_size = Some(v);
                Ok(2 + payload)
            }
            TAG_CHILD_ENUMS => {
                tag.child_type = ValueType::Enum;
                let s = self.read_tag_str(payload)?;
                tag.child_enums = Some(s.clone());
                let n = if payload <= 5 {
                    1 + payload
                } else if payload == 6 {
                    2 + s.len()
                } else {
                    3 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_PATTERN => {
                let s = self.read_tag_short_str(payload)?;
                tag.child_pattern = Some(s.clone());
                let n = if payload < 7 {
                    1 + payload
                } else {
                    2 + s.len()
                };
                Ok(n)
            }
            TAG_CHILD_LOCATION => {
                let s = self.read_tag_ascii(payload)?;
                if let Ok(n) = s.parse::<i32>() {
                    tag.child_location = Some(n);
                }
                Ok(1 + payload)
            }
            TAG_CHILD_VERSION => {
                let v = self.read_tag_uint(payload)?;
                tag.child_version = Some(v as i32);
                Ok(2 + payload)
            }
            TAG_CHILD_MIME => {
                tag.child_type = ValueType::Media;
                let mime_id = if payload < 7 {
                    payload as u8
                } else {
                    self.read_byte()?
                };
                tag.child_mime = Some(mime_to_str(mime_id).to_string());
                let extra = if payload < 7 { 0 } else { 1 };
                Ok(1 + extra)
            }
            _ => Ok(1),
        }
    }

    fn read_tag_str(&mut self, payload: usize) -> Result<String, std::io::Error> {
        let len = if payload <= 5 {
            payload
        } else if payload == 6 {
            self.read_byte()? as usize
        } else {
            let hi = self.read_byte()? as usize;
            let lo = self.read_byte()? as usize;
            (hi << 8) | lo
        };
        let bytes = self.read_bytes(len)?.to_vec();
        Ok(String::from_utf8_lossy(&bytes).to_string())
    }

    fn read_tag_short_str(&mut self, payload: usize) -> Result<String, std::io::Error> {
        let len = if payload < 7 {
            payload
        } else {
            self.read_byte()? as usize
        };
        let bytes = self.read_bytes(len)?.to_vec();
        Ok(String::from_utf8_lossy(&bytes).to_string())
    }

    fn read_tag_ascii(&mut self, payload: usize) -> Result<String, std::io::Error> {
        let bytes = self.read_bytes(payload)?.to_vec();
        Ok(String::from_utf8_lossy(&bytes).to_string())
    }

    fn read_tag_uint(&mut self, payload: usize) -> Result<u64, std::io::Error> {
        let nbytes = payload + 1;
        let mut v: u64 = 0;
        for _ in 0..nbytes {
            let b = self.read_byte()?;
            v = (v << 8) | (b as u64);
        }
        Ok(v)
    }

    fn decode_simple(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let value = SimpleValue::from_byte(prefix & 0x1F).ok_or_else(|| {
            std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid simple value")
        })?;

        let (data, text) = match value {
            SimpleValue::False => (ValueData::Bool(false), "false".to_string()),
            SimpleValue::True => (ValueData::Bool(true), "true".to_string()),
            SimpleValue::NullBool => (ValueData::Bool(false), "false".to_string()),
            SimpleValue::NullFloat => (ValueData::Float(0.0), "0.0".to_string()),
            SimpleValue::NullInt => (ValueData::Int(0), "0".to_string()),
            SimpleValue::NullString => (ValueData::String(String::new()), String::new()),
            SimpleValue::NullBytes => (ValueData::Bytes(vec![]), String::new()),
            _ => (ValueData::String(String::new()), String::new()),
        };

        Ok(Node::Value(Value {
            data,
            text,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_positive_int(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = prefix & INT_LEN_MASK;
        let v = if l < INT_LEN_1 {
            l as u64
        } else {
            let extra_bytes = (l - INT_LEN_1 + 1) as usize;
            let mut v = 0u64;
            for _ in 0..extra_bytes {
                let b = self.read_byte()?;
                v = (v << 8) | (b as u64);
            }
            v
        };

        let (data, text): (ValueData, String) = match tag.value_type {
            ValueType::Datetime => {
                let naive = chrono::DateTime::from_timestamp(v as i64, 0)
                    .map(|dt| dt.naive_utc())
                    .unwrap_or(chrono::NaiveDateTime::default());
                let offset = if let Some(offset_hours) = tag.location {
                    chrono::FixedOffset::east_opt(offset_hours * 3600)
                        .unwrap_or_else(|| chrono::FixedOffset::east_opt(0).unwrap())
                } else {
                    chrono::FixedOffset::east_opt(0).unwrap()
                };
                let dt = offset.from_utc_datetime(&naive);
                let text = dt.format("%Y-%m-%d %H:%M:%S").to_string();
                (ValueData::Int(v as i64), text)
            }
            ValueType::Date => {
                let epoch = chrono::NaiveDate::from_ymd_opt(1970, 1, 1).unwrap();
                let date = epoch + chrono::Duration::days(v as i64);
                let text = if let Some(offset_hours) = tag.location {
                    let dt = date.and_hms_opt(0, 0, 0).unwrap();
                    let offset = chrono::FixedOffset::east_opt(offset_hours * 3600)
                        .unwrap_or_else(|| chrono::FixedOffset::east_opt(0).unwrap());
                    let localized = offset.from_utc_datetime(&dt);
                    localized.format("%Y-%m-%d").to_string()
                } else {
                    date.format("%Y-%m-%d").to_string()
                };
                (ValueData::Int(v as i64), text)
            }
            ValueType::Time => {
                let hours = (v / 3600) % 24;
                let minutes = (v % 3600) / 60;
                let seconds = v % 60;
                let text = if let Some(offset_hours) = tag.location {
                    let base = chrono::NaiveTime::from_hms_opt(0, 0, 0).unwrap();
                    let time = chrono::NaiveTime::from_hms_opt(
                        hours as u32,
                        minutes as u32,
                        seconds as u32,
                    )
                    .unwrap_or(base);
                    let base_dt = chrono::NaiveDate::from_ymd_opt(1970, 1, 1)
                        .unwrap()
                        .and_time(time);
                    let offset = chrono::FixedOffset::east_opt(offset_hours * 3600)
                        .unwrap_or_else(|| chrono::FixedOffset::east_opt(0).unwrap());
                    let localized = offset.from_utc_datetime(&base_dt);
                    localized.format("%H:%M:%S").to_string()
                } else {
                    format!("{:02}:{:02}:{:02}", hours, minutes, seconds)
                };
                (ValueData::Int(v as i64), text)
            }
            ValueType::Enum => {
                if let Some(ref enums_str) = tag.enums {
                    let enums: Vec<&str> = enums_str.split('|').collect();
                    let idx = v as usize;
                    if idx < enums.len() {
                        let text = enums[idx].trim().to_string();
                        (ValueData::Int(v as i64), text)
                    } else {
                        let text = v.to_string();
                        (ValueData::Int(v as i64), text)
                    }
                } else {
                    let text = v.to_string();
                    (ValueData::Int(v as i64), text)
                }
            }
            _ => {
                let text = v.to_string();
                let data = ValueData::Int(v as i64);
                (data, text)
            }
        };

        Ok(Node::Value(Value {
            data,
            text,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_negative_int(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = prefix & INT_LEN_MASK;
        let v = if l < INT_LEN_1 {
            l as u64
        } else {
            let extra_bytes = (l - INT_LEN_1 + 1) as usize;
            let mut v = 0u64;
            for _ in 0..extra_bytes {
                let b = self.read_byte()?;
                v = (v << 8) | (b as u64);
            }
            v
        };

        let data = ValueData::Int(-(v as i64));
        let text = format!("-{}", v);

        Ok(Node::Value(Value {
            data,
            text,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_float(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = prefix & FLOAT_LEN_MASK;
        let v: f64;

        if l < FLOAT_LEN_1 {
            v = l as f64 / 10.0;
        } else {
            let exp = self.read_byte()? as i8;
            let extra_bytes = (l - FLOAT_LEN_1 + 1) as usize;

            let mantissa: u64 = if extra_bytes == 0 {
                0
            } else {
                let mut m = 0u64;
                for _ in 0..extra_bytes {
                    let b = self.read_byte()?;
                    m = (m << 8) | (b as u64);
                }
                m
            };

            let dec = mantissa_to_decimal(mantissa, exp);
            v = dec.parse().unwrap_or(0.0);
        }

        let v = if (prefix & FLOAT_POSITIVE_NEGATIVE_MASK) != 0 {
            -v
        } else {
            v
        };

        Ok(Node::Value(Value {
            data: ValueData::Float(v),
            text: ryu::Buffer::new().format_finite(v).to_string(),
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_string(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = prefix & INT_LEN_MASK;
        let len = if l < STRING_LEN_1 {
            l as usize
        } else if l == STRING_LEN_1 {
            self.read_byte()? as usize
        } else {
            let b = self.read_bytes(2)?;
            ((b[0] as usize) << 8) | (b[1] as usize)
        };

        let s = if len > 0 {
            String::from_utf8_lossy(self.read_bytes(len)?).to_string()
        } else {
            String::new()
        };

        Ok(Node::Value(Value {
            data: ValueData::String(s.clone()),
            text: s,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_bytes(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = prefix & INT_LEN_MASK;
        let len = if l < BYTES_LEN_1 {
            l as usize
        } else if l == BYTES_LEN_1 {
            self.read_byte()? as usize
        } else {
            let b = self.read_bytes(2)?;
            ((b[0] as usize) << 8) | (b[1] as usize)
        };

        let bytes = if len > 0 {
            self.read_bytes(len)?.to_vec()
        } else {
            vec![]
        };

        let text = match tag.value_type {
            ValueType::Uuid => {
                if bytes.len() == 16 {
                    let mut uuid_bytes = [0u8; 16];
                    uuid_bytes.copy_from_slice(&bytes);
                    uuid::Uuid::from_bytes(uuid_bytes).to_string()
                } else {
                    use base64::{engine::general_purpose, Engine as _};
                    general_purpose::STANDARD.encode(&bytes)
                }
            }
            ValueType::Bigint => {
                let decoded =
                    crate::core::utils::decode_big_int(&bytes).unwrap_or_else(|_| "0".to_string());
                decoded
            }
            _ => {
                use base64::{engine::general_purpose, Engine as _};
                general_purpose::STANDARD.encode(&bytes)
            }
        };

        Ok(Node::Value(Value {
            data: ValueData::Bytes(bytes.clone()),
            text,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_container(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        if (prefix & CONTAINER_TYPE_MASK) != 0 {
            self.decode_array(prefix, tag)
        } else {
            self.decode_object(prefix, tag)
        }
    }

    fn decode_container_len(&mut self, prefix: u8) -> Result<usize, std::io::Error> {
        let l = prefix & CONTAINER_LEN_MASK;
        if l < CONTAINER_LEN_1 {
            Ok(l as usize)
        } else if l == CONTAINER_LEN_1 {
            Ok(self.read_byte()? as usize)
        } else {
            let b = self.read_bytes(2)?;
            Ok(((b[0] as usize) << 8) | (b[1] as usize))
        }
    }

    fn decode_array(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let total_bytes = self.decode_container_len(prefix)?;

        let mut items = Vec::new();
        let mut index = 0;
        while index < total_bytes {
            let mut child_tag = Tag::new();
            child_tag.inherit(tag);
            if tag.child_type == ValueType::Unknown {
                if child_tag.value_type == ValueType::Unknown {
                    child_tag.value_type = tag.value_type;
                }
                if child_tag.enums.is_none() {
                    child_tag.enums = tag.enums.clone();
                }
                if child_tag.mime.is_none() {
                    child_tag.mime = tag.mime.clone();
                }
                if child_tag.version.is_none() {
                    child_tag.version = tag.version;
                }
                if child_tag.location.is_none() && tag.location.is_some() {
                    child_tag.location = tag.location;
                }
            }
            let before = self.offset;
            let item = self.decode_node(&child_tag, "")?;
            let consumed = self.offset - before;
            index += consumed;
            items.push(item);
        }

        Ok(Node::Array(Array {
            items,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }

    fn decode_object(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let _total_bytes = self.decode_container_len(prefix)?;

        let keys_node = self.decode_node(&Tag::new(), "")?;
        let keys = if let Node::Array(arr) = keys_node {
            arr.items
        } else {
            vec![]
        };

        let mut fields = Vec::with_capacity(keys.len());
        for key_node in &keys {
            let key = if let Node::Value(v) = key_node {
                v.text.clone()
            } else {
                continue;
            };

            let mut child_tag = Tag::new();
            child_tag.inherit(tag);
            if tag.child_type == ValueType::Unknown {
                if child_tag.value_type == ValueType::Unknown {
                    child_tag.value_type = tag.value_type;
                }
                if child_tag.enums.is_none() {
                    child_tag.enums = tag.enums.clone();
                }
                if child_tag.mime.is_none() {
                    child_tag.mime = tag.mime.clone();
                }
                if child_tag.version.is_none() {
                    child_tag.version = tag.version;
                }
                if child_tag.location.is_none() && tag.location.is_some() {
                    child_tag.location = tag.location;
                }
            }
            let value = self.decode_node(&child_tag, "")?;
            fields.push(Field { key, value });
        }

        Ok(Node::Object(Object {
            fields,
            path: String::new(),
            tag: Some(tag.clone()),
        }))
    }
}

fn mantissa_to_decimal(mantissa: u64, exp: i8) -> String {
    let num_str = mantissa.to_string();
    let decimal_pos = num_str.len() as i32 + (exp as i32);

    if decimal_pos <= 0 {
        format!("0.{}{}", "0".repeat((-decimal_pos) as usize), num_str)
    } else if (decimal_pos as usize) < num_str.len() {
        format!(
            "{}.{}",
            &num_str[..decimal_pos as usize],
            &num_str[decimal_pos as usize..]
        )
    } else {
        format!(
            "{}{}",
            num_str,
            "0".repeat((decimal_pos as usize).saturating_sub(num_str.len()))
        )
    }
}
