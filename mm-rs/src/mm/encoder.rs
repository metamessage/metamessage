use crate::jsonc::ast::{Node, Object, Array, Value};
use crate::mm::constants::*;
use crate::mm::prefix::{PREFIX_SIMPLE, PREFIX_POSITIVE_INT, PREFIX_NEGATIVE_INT, PREFIX_FLOAT, PREFIX_STRING, PREFIX_BYTES, PREFIX_CONTAINER, PREFIX_TAG, FLOAT_POSITIVE_NEGATIVE_MASK};
use crate::mm::simple_value::SimpleValue;

pub struct Encoder {
    buf: Vec<u8>,
    offset: usize,
}

impl Encoder {
    pub fn new() -> Self {
        Self {
            buf: vec![0u8; DEFAULT_BUF_SIZE],
            offset: 0,
        }
    }

    pub fn with_capacity(capacity: usize) -> Self {
        Self {
            buf: vec![0u8; capacity],
            offset: 0,
        }
    }

    fn ensure_capacity(&mut self, needed: usize) {
        while self.offset + needed > self.buf.len() {
            self.buf.resize(self.buf.len() * 2, 0);
        }
    }

    fn write_byte(&mut self, b: u8) {
        self.ensure_capacity(1);
        self.buf[self.offset] = b;
        self.offset += 1;
    }

    fn write_bytes(&mut self, bytes: &[u8]) {
        self.ensure_capacity(bytes.len());
        self.buf[self.offset..self.offset + bytes.len()].copy_from_slice(bytes);
        self.offset += bytes.len();
    }

    pub fn encode(&mut self, node: &Node) -> Vec<u8> {
        match node {
            Node::Object(obj) => {
                self.encode_object(obj);
            }
            Node::Array(arr) => {
                self.encode_array(arr);
            }
            Node::Value(val) => {
                self.encode_value(val);
            }
        }
        self.buf[..self.offset].to_vec()
    }

    fn encode_object(&mut self, obj: &Object) {
        let mut key_buf = Vec::new();
        let mut val_buf = Vec::new();

        for field in &obj.fields {
            let field_offset = self.offset;
            self.offset = 0;

            match &field.value {
                Node::Object(o) => self.encode_object(o),
                Node::Array(a) => self.encode_array(a),
                Node::Value(v) => self.encode_value(v),
            }

            let encoded = self.buf[..self.offset].to_vec();
            val_buf.extend_from_slice(&encoded);

            self.offset = 0;
            self.encode_string(&field.key);
            let encoded_key = self.buf[..self.offset].to_vec();
            key_buf.extend_from_slice(&encoded_key);

            self.offset = field_offset;
        }

        let total_len = key_buf.len() + val_buf.len();

        let mut prefix = PREFIX_CONTAINER | CONTAINER_MAP;

        if total_len < 254 {
            self.write_byte(prefix);
            self.write_bytes(&val_buf);
            self.write_bytes(&key_buf);
        } else if total_len < 65536 {
            prefix |= CONTAINER_LEN_1;
            self.write_byte(prefix);
            self.write_byte((total_len & 0xFF) as u8);
            self.write_bytes(&val_buf);
            self.write_bytes(&key_buf);
        } else {
            prefix |= CONTAINER_LEN_2;
            self.write_byte(prefix);
            self.write_byte((total_len >> 8 & 0xFF) as u8);
            self.write_byte((total_len & 0xFF) as u8);
            self.write_bytes(&val_buf);
            self.write_bytes(&key_buf);
        }

        if let Some(t) = &obj.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.encode_tag(&tag_bytes);
            }
        }
    }

    fn encode_array(&mut self, arr: &Array) {
        let mut val_buf = Vec::new();

        for item in &arr.items {
            let field_offset = self.offset;
            self.offset = 0;

            match item {
                Node::Object(o) => self.encode_object(o),
                Node::Array(a) => self.encode_array(a),
                Node::Value(v) => self.encode_value(v),
            }

            let encoded = self.buf[..self.offset].to_vec();
            val_buf.extend_from_slice(&encoded);

            self.offset = field_offset;
        }

        let total_len = val_buf.len();

        let mut prefix = PREFIX_CONTAINER | CONTAINER_ARRAY;

        if total_len < 254 {
            self.write_byte(prefix);
            self.write_bytes(&val_buf);
        } else if total_len < 65536 {
            prefix |= CONTAINER_LEN_1;
            self.write_byte(prefix);
            self.write_byte((total_len & 0xFF) as u8);
            self.write_bytes(&val_buf);
        } else {
            prefix |= CONTAINER_LEN_2;
            self.write_byte(prefix);
            self.write_byte((total_len >> 8 & 0xFF) as u8);
            self.write_byte((total_len & 0xFF) as u8);
            self.write_bytes(&val_buf);
        }

        if let Some(t) = &arr.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.encode_tag(&tag_bytes);
            }
        }
    }

    fn encode_value(&mut self, val: &Value) {
        match &val.data {
            crate::jsonc::ast::ValueData::Bool(b) => {
                self.encode_bool(*b);
            }
            crate::jsonc::ast::ValueData::String(s) => {
                self.encode_string(s);
            }
            crate::jsonc::ast::ValueData::Int(i) => {
                self.encode_int64(*i);
            }
            crate::jsonc::ast::ValueData::Float(f) => {
                self.encode_float(*f);
            }
            crate::jsonc::ast::ValueData::Bytes(b) => {
                self.encode_bytes(b);
            }
            crate::jsonc::ast::ValueData::Null => {
                self.encode_simple(SimpleValue::NullInt);
            }
            _ => {}
        }

        if let Some(t) = &val.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.encode_tag(&tag_bytes);
            }
        }
    }

    pub fn encode_bool(&mut self, v: bool) {
        let value = if v { SimpleValue::True } else { SimpleValue::False };
        self.encode_simple(value);
    }

    fn encode_simple(&mut self, value: SimpleValue) {
        self.write_byte(PREFIX_SIMPLE | value.to_byte());
    }

    pub fn encode_int64(&mut self, v: i64) {
        if v >= 0 {
            self.encode_uint64(v as u64);
        } else {
            let uv = if v == i64::MIN {
                9223372036854775808u64
            } else {
                (-v) as u64
            };
            self.encode_uint64_with_sign(PREFIX_NEGATIVE_INT, uv);
        }
    }

    pub fn encode_uint64(&mut self, v: u64) {
        self.encode_uint64_with_sign(PREFIX_POSITIVE_INT, v);
    }

    fn encode_uint64_with_sign(&mut self, sign: u8, v: u64) {
        if v < INT_LEN_1 as u64 {
            self.write_byte(sign | (v as u8));
        } else if v < MAX_1 as u64 {
            self.write_byte(sign | INT_LEN_1);
            self.write_byte(v as u8);
        } else if v < MAX_2 as u64 {
            self.write_byte(sign | INT_LEN_2);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v < MAX_3 as u64 {
            self.write_byte(sign | INT_LEN_3);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v < MAX_4 as u64 {
            self.write_byte(sign | INT_LEN_4);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v < MAX_5 as u64 {
            self.write_byte(sign | INT_LEN_5);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v < MAX_6 as u64 {
            self.write_byte(sign | INT_LEN_6);
            self.write_byte((v >> 40) as u8);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v < MAX_7 as u64 {
            self.write_byte(sign | INT_LEN_7);
            self.write_byte((v >> 48) as u8);
            self.write_byte((v >> 40) as u8);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else {
            self.write_byte(sign | INT_LEN_8);
            self.write_byte((v >> 56) as u8);
            self.write_byte((v >> 48) as u8);
            self.write_byte((v >> 40) as u8);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        }
    }

    pub fn encode_float(&mut self, f: f64) {
        let mut buffer = ryu::Buffer::new();
        let bytes = buffer.format(f).as_bytes();
        self.write_byte(PREFIX_FLOAT | FLOAT_POSITIVE_NEGATIVE_MASK);
        self.write_bytes(bytes);
    }

    pub fn encode_string(&mut self, s: &str) {
        let utf = s.as_bytes();
        let len = utf.len();
        let mut prefix = PREFIX_STRING;

        if len < 254 {
            self.write_byte(prefix | len as u8);
            self.write_bytes(utf);
        } else if len < 65536 {
            prefix |= STRING_LEN_1;
            self.write_byte(prefix);
            self.write_byte((len & 0xFF) as u8);
            self.write_bytes(utf);
        } else {
            prefix |= STRING_LEN_2;
            self.write_byte(prefix);
            self.write_byte((len >> 8 & 0xFF) as u8);
            self.write_byte((len & 0xFF) as u8);
            self.write_bytes(utf);
        }
    }

    pub fn encode_bytes(&mut self, bytes: &[u8]) {
        let len = bytes.len();
        let mut prefix = PREFIX_BYTES;

        if len < 254 {
            self.write_byte(prefix | len as u8);
            self.write_bytes(bytes);
        } else if len < 65536 {
            prefix |= BYTES_LEN_1;
            self.write_byte(prefix);
            self.write_byte((len & 0xFF) as u8);
            self.write_bytes(bytes);
        } else {
            prefix |= BYTES_LEN_2;
            self.write_byte(prefix);
            self.write_byte((len >> 8 & 0xFF) as u8);
            self.write_byte((len & 0xFF) as u8);
            self.write_bytes(bytes);
        }
    }

    fn encode_tag(&mut self, tag_bytes: &[u8]) {
        let len = tag_bytes.len();
        let mut prefix = PREFIX_TAG;

        if len < 254 {
            self.write_byte(prefix);
            self.write_byte(len as u8);
            self.write_bytes(tag_bytes);
        } else if len < 65536 {
            prefix |= TAG_LEN_1;
            self.write_byte(prefix);
            self.write_byte(len as u8);
            self.write_bytes(tag_bytes);
        } else {
            prefix |= TAG_LEN_2;
            self.write_byte(prefix);
            self.write_byte((len >> 8 & 0xFF) as u8);
            self.write_byte((len & 0xFF) as u8);
            self.write_bytes(tag_bytes);
        }
    }
}

impl Default for Encoder {
    fn default() -> Self {
        Self::new()
    }
}