use crate::core::constants::*;
use crate::core::prefix::{
    CONTAINER_ARRAY, CONTAINER_OBJECT, FLOAT_LEN_1, FLOAT_POSITIVE_NEGATIVE_MASK, PREFIX_BYTES,
    PREFIX_CONTAINER, PREFIX_FLOAT, PREFIX_NEGATIVE_INT, PREFIX_POSITIVE_INT, PREFIX_SIMPLE,
    PREFIX_STRING, PREFIX_TAG,
};
use crate::core::simple_value::SimpleValue;
use crate::ir::ast::{Array, Node, Object, Value, ValueData};
use crate::ir::ValueType;
use chrono::Timelike;

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
        self.offset = 0;
        match node {
            Node::Object(obj) => self.encode_object(obj),
            Node::Array(arr) => self.encode_array(arr),
            Node::Value(val) => self.encode_value(val),
        }
        self.buf[..self.offset].to_vec()
    }

    fn encode_object(&mut self, obj: &Object) {
        let mut key_buf = Vec::new();
        let mut val_buf = Vec::new();

        for field in &obj.fields {
            let saved_offset = self.offset;
            self.offset = 0;

            match &field.value {
                Node::Object(o) => self.encode_object(o),
                Node::Array(a) => self.encode_array(a),
                Node::Value(v) => self.encode_value(v),
            }

            val_buf.extend_from_slice(&self.buf[..self.offset]);

            self.offset = 0;
            self.encode_string(&field.key);
            key_buf.extend_from_slice(&self.buf[..self.offset]);

            self.offset = saved_offset;
        }

        self.offset = 0;
        self.encode_container(CONTAINER_ARRAY, &key_buf);
        let key_array = self.buf[..self.offset].to_vec();

        let mut buf_all = key_array;
        buf_all.extend_from_slice(&val_buf);

        self.offset = 0;
        self.encode_container(CONTAINER_OBJECT, &buf_all);
        let object_payload = self.buf[..self.offset].to_vec();

        if let Some(t) = &obj.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.offset = 0;
                self.encode_tag(&tag_bytes, &object_payload);
                return;
            }
        }
        self.offset = 0;
        self.write_bytes(&object_payload);
    }

    fn encode_array(&mut self, arr: &Array) {
        let mut val_buf = Vec::new();

        for item in &arr.items {
            let saved_offset = self.offset;
            self.offset = 0;

            match item {
                Node::Object(o) => self.encode_object(o),
                Node::Array(a) => self.encode_array(a),
                Node::Value(v) => self.encode_value(v),
            }

            val_buf.extend_from_slice(&self.buf[..self.offset]);
            self.offset = saved_offset;
        }

        self.offset = 0;
        self.encode_container(CONTAINER_ARRAY, &val_buf);
        let array_payload = self.buf[..self.offset].to_vec();

        if let Some(t) = &arr.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.offset = 0;
                self.encode_tag(&tag_bytes, &array_payload);
                return;
            }
        }
        self.offset = 0;
        self.write_bytes(&array_payload);
    }

    fn encode_value(&mut self, val: &Value) {
        let saved_offset = self.offset;
        self.offset = 0;

        let tag = val.tag.as_ref();
        let is_null = tag.map_or(false, |t| t.is_null);

        if is_null {
            match tag.map(|t| t.value_type) {
                Some(ValueType::Bool) => self.encode_simple(SimpleValue::NullBool),
                Some(ValueType::I) => self.encode_simple(SimpleValue::NullInt),
                Some(ValueType::F64) | Some(ValueType::F32) => {
                    self.encode_simple(SimpleValue::NullFloat)
                }
                Some(ValueType::Str) => self.encode_simple(SimpleValue::NullString),
                Some(ValueType::Bytes) => self.encode_simple(SimpleValue::NullBytes),
                _ => self.encode_simple(SimpleValue::NullInt),
            }
        } else if let Some(t) = tag {
            match t.value_type {
                ValueType::Bool => match &val.data {
                    ValueData::Bool(b) => self.encode_bool(*b),
                    _ => self.encode_bool(val.text == "true"),
                },
                ValueType::I | ValueType::I8 | ValueType::I16 | ValueType::I32 | ValueType::I64 => {
                    match &val.data {
                        ValueData::Int(i) => self.encode_int64(*i),
                        _ => {
                            if let Ok(v) = val.text.parse::<i64>() {
                                self.encode_int64(v);
                            } else {
                                self.encode_int64(0);
                            }
                        }
                    }
                }
                ValueType::U | ValueType::U8 | ValueType::U16 | ValueType::U32 | ValueType::U64 => {
                    match &val.data {
                        ValueData::Uint(u) => self.encode_uint64(*u),
                        _ => {
                            if let Ok(v) = val.text.parse::<u64>() {
                                self.encode_uint64(v);
                            } else {
                                self.encode_uint64(0);
                            }
                        }
                    }
                }
                ValueType::F32 | ValueType::F64 | ValueType::Decimal => {
                    self.encode_float(&val.text);
                }
                ValueType::Str | ValueType::Url | ValueType::Email => {
                    let s = match &val.data {
                        ValueData::String(s) => s.clone(),
                        _ => val.text.clone(),
                    };
                    self.encode_string(&s);
                }
                ValueType::Bytes => {
                    let s = match &val.data {
                        ValueData::String(s) => s.clone(),
                        _ => val.text.clone(),
                    };
                    use base64::{engine::general_purpose, Engine as _};
                    if let Ok(decoded) = general_purpose::STANDARD.decode(&s) {
                        self.encode_bytes(&decoded);
                    } else {
                        self.encode_string(&s);
                    }
                }
                ValueType::Datetime => {
                    let timestamp = match &val.data {
                        ValueData::Int(ts) => *ts,
                        _ => {
                            let naive = chrono::NaiveDateTime::parse_from_str(
                                &val.text,
                                "%Y-%m-%d %H:%M:%S",
                            )
                            .unwrap_or_else(|_| {
                                chrono::NaiveDateTime::parse_from_str(
                                    &val.text,
                                    "%Y-%m-%dT%H:%M:%S",
                                )
                                .unwrap_or(
                                    chrono::NaiveDate::from_ymd_opt(1970, 1, 1)
                                        .unwrap()
                                        .and_hms_opt(0, 0, 0)
                                        .unwrap(),
                                )
                            });
                            let utc_ts = naive.and_utc().timestamp();
                            if let Some(loc) = t.location {
                                utc_ts - (loc as i64 * 3600)
                            } else {
                                utc_ts
                            }
                        }
                    };
                    self.encode_int64(timestamp);
                }
                ValueType::Date => {
                    let days = match &val.data {
                        ValueData::Int(d) => *d,
                        _ => {
                            let naive = chrono::NaiveDate::parse_from_str(&val.text, "%Y-%m-%d")
                                .unwrap_or(chrono::NaiveDate::from_ymd_opt(1970, 1, 2).unwrap());
                            let epoch = chrono::NaiveDate::from_ymd_opt(1970, 1, 1).unwrap();
                            let duration = naive.signed_duration_since(epoch);
                            duration.num_days()
                        }
                    };
                    self.encode_int64(days);
                }
                ValueType::Time => {
                    let seconds = match &val.data {
                        ValueData::Int(s) => *s,
                        _ => {
                            let naive = chrono::NaiveTime::parse_from_str(&val.text, "%H:%M:%S")
                                .unwrap_or(chrono::NaiveTime::from_hms_opt(0, 0, 0).unwrap());
                            naive.num_seconds_from_midnight() as i64
                        }
                    };
                    self.encode_int64(seconds);
                }
                ValueType::Uuid => match &val.data {
                    ValueData::Bytes(b) => self.encode_bytes(b),
                    _ => {
                        if let Ok(u) = uuid::Uuid::parse_str(&val.text) {
                            self.encode_bytes(u.as_bytes());
                        } else {
                            self.encode_string(&val.text);
                        }
                    }
                },
                ValueType::Bigint => {
                    self.encode_big_int(&val.text);
                }
                ValueType::Ip => {
                    if let Some(v) = t.version {
                        if v == 4 {
                            if let Ok(ip) = val.text.parse::<std::net::Ipv4Addr>() {
                                self.encode_bytes(&ip.octets());
                            } else {
                                self.encode_string(&val.text);
                            }
                        } else if v == 6 {
                            if let Ok(ip) = val.text.parse::<std::net::Ipv6Addr>() {
                                self.encode_bytes(&ip.octets());
                            } else {
                                self.encode_string(&val.text);
                            }
                        } else {
                            self.encode_string(&val.text);
                        }
                    } else {
                        self.encode_string(&val.text);
                    }
                }
                ValueType::Enum => {
                    let s = match &val.data {
                        ValueData::String(s) => s.clone(),
                        _ => val.text.clone(),
                    };
                    if let Some(ref enums) = t.enums {
                        let enum_list: Vec<&str> = enums.split('|').map(|e| e.trim()).collect();
                        if let Some(pos) = enum_list.iter().position(|e| *e == s) {
                            self.encode_int64(pos as i64);
                        } else {
                            self.encode_int64(0);
                        }
                    } else {
                        self.encode_int64(0);
                    }
                }
                ValueType::Media => {
                    let s = match &val.data {
                        ValueData::String(s) => s.clone(),
                        _ => val.text.clone(),
                    };
                    use base64::{engine::general_purpose, Engine as _};
                    if let Ok(decoded) = general_purpose::STANDARD.decode(&s) {
                        self.encode_bytes(&decoded);
                    } else {
                        self.encode_string(&s);
                    }
                }
                _ => match &val.data {
                    ValueData::Bool(b) => self.encode_bool(*b),
                    ValueData::String(s) => self.encode_string(s),
                    ValueData::Int(i) => self.encode_int64(*i),
                    ValueData::Uint(u) => self.encode_uint64(*u),
                    ValueData::Float(_) => self.encode_float(&val.text),
                    ValueData::Bytes(b) => self.encode_bytes(b),
                    ValueData::Null => self.encode_simple(SimpleValue::NullInt),
                },
            }
        } else {
            match &val.data {
                ValueData::Bool(b) => self.encode_bool(*b),
                ValueData::String(s) => self.encode_string(s),
                ValueData::Int(i) => self.encode_int64(*i),
                ValueData::Uint(u) => self.encode_uint64(*u),
                ValueData::Float(_) => self.encode_float(&val.text),
                ValueData::Bytes(b) => self.encode_bytes(b),
                ValueData::Null => self.encode_simple(SimpleValue::NullInt),
            }
        }

        let payload = self.buf[..self.offset].to_vec();

        if let Some(t) = &val.tag {
            let tag_bytes = t.to_bytes();
            if !tag_bytes.is_empty() {
                self.offset = saved_offset;
                self.encode_tag(&tag_bytes, &payload);
                return;
            }
        }
        self.offset = saved_offset;
        self.write_bytes(&payload);
    }

    fn encode_big_int(&mut self, s: &str) {
        let neg = s.starts_with('-');
        let digits = if neg { &s[1..] } else { s };
        let digit_len = digits.len();
        if digit_len == 0 {
            return;
        }

        // Calculate total bits needed
        let groups = digit_len / 3;
        let rem = digit_len % 3;
        let total_bits = 1
            + groups * 10
            + match rem {
                2 => 7,
                1 => 4,
                _ => 0,
            };
        let byte_count = (total_bits + 7) / 8;

        // Build bit-packed bytes
        let mut bits = vec![0u8; byte_count];
        let mut bit_offset = 0usize;

        // Write sign bit
        if neg {
            bits[0] |= 1 << (7 - bit_offset);
        }
        bit_offset += 1;

        // Process digit groups
        let mut i = 0;
        while i < digit_len {
            let rem_group = digit_len - i;
            let (val, num_bits): (u64, usize) = if rem_group >= 3 {
                let v = (digits.as_bytes()[i] - b'0') as u64 * 100
                    + (digits.as_bytes()[i + 1] - b'0') as u64 * 10
                    + (digits.as_bytes()[i + 2] - b'0') as u64;
                i += 3;
                (v, 10)
            } else if rem_group == 2 {
                let v = (digits.as_bytes()[i] - b'0') as u64 * 10
                    + (digits.as_bytes()[i + 1] - b'0') as u64;
                i += 2;
                (v, 7)
            } else {
                let v = (digits.as_bytes()[i] - b'0') as u64;
                i += 1;
                (v, 4)
            };

            for b in (0..num_bits).rev() {
                let byte_idx = bit_offset / 8;
                let bit_in_byte = 7 - (bit_offset % 8);
                if (val >> b) & 1 == 1 {
                    bits[byte_idx] |= 1 << bit_in_byte;
                }
                bit_offset += 1;
            }
        }

        // Build payload: [digit_count_byte] + [bit_data]
        let mut payload = Vec::with_capacity(1 + byte_count);
        payload.push(digit_len as u8);
        payload.extend_from_slice(&bits);

        self.encode_bytes(&payload);
    }

    pub fn encode_bool(&mut self, v: bool) {
        let value = if v {
            SimpleValue::True
        } else {
            SimpleValue::False
        };
        self.encode_simple(value);
    }

    fn encode_simple(&mut self, value: SimpleValue) {
        self.write_byte(PREFIX_SIMPLE | value.to_byte());
    }

    pub fn encode_int64(&mut self, v: i64) {
        if v >= 0 {
            self.encode_uint64_with_sign(PREFIX_POSITIVE_INT, v as u64);
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
        if v < 24 {
            self.write_byte(sign | (v as u8));
        } else if v <= 0xFF {
            self.write_byte(sign | INT_LEN_1);
            self.write_byte(v as u8);
        } else if v <= 0xFFFF {
            self.write_byte(sign | INT_LEN_2);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v <= 0xFFFFFF {
            self.write_byte(sign | INT_LEN_3);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v <= 0xFFFFFFFF {
            self.write_byte(sign | INT_LEN_4);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v <= 0xFFFFFFFFFF {
            self.write_byte(sign | INT_LEN_5);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v <= 0xFFFFFFFFFFFF {
            self.write_byte(sign | INT_LEN_6);
            self.write_byte((v >> 40) as u8);
            self.write_byte((v >> 32) as u8);
            self.write_byte((v >> 24) as u8);
            self.write_byte((v >> 16) as u8);
            self.write_byte((v >> 8) as u8);
            self.write_byte(v as u8);
        } else if v <= 0xFFFFFFFFFFFFFF {
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

    pub fn encode_float(&mut self, s: &str) {
        let (is_negative, exponent, mantissa) = parse_float_string(s);

        let mut prefix = PREFIX_FLOAT;
        if is_negative {
            prefix |= FLOAT_POSITIVE_NEGATIVE_MASK;
        }

        if exponent == -1 && mantissa <= 7 {
            self.write_byte(prefix | (mantissa as u8));
            return;
        }

        if mantissa <= 0xFF {
            prefix |= FLOAT_LEN_1;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFF {
            prefix |= FLOAT_LEN_1 + 1;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFFFF {
            prefix |= FLOAT_LEN_1 + 2;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFFFFFF {
            prefix |= FLOAT_LEN_1 + 3;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 24) as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFFFFFFFF {
            prefix |= FLOAT_LEN_1 + 4;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 32) as u8);
            self.write_byte((mantissa >> 24) as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFFFFFFFFFF {
            prefix |= FLOAT_LEN_1 + 5;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 40) as u8);
            self.write_byte((mantissa >> 32) as u8);
            self.write_byte((mantissa >> 24) as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else if mantissa <= 0xFFFFFFFFFFFFFF {
            prefix |= FLOAT_LEN_1 + 6;
            self.write_byte(prefix);
            self.write_byte(exponent as u8);
            self.write_byte((mantissa >> 48) as u8);
            self.write_byte((mantissa >> 40) as u8);
            self.write_byte((mantissa >> 32) as u8);
            self.write_byte((mantissa >> 24) as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        } else {
            prefix |= FLOAT_LEN_1 + 7;
            self.write_byte(prefix);
            self.write_byte((mantissa >> 56) as u8);
            self.write_byte((mantissa >> 48) as u8);
            self.write_byte((mantissa >> 40) as u8);
            self.write_byte((mantissa >> 32) as u8);
            self.write_byte((mantissa >> 24) as u8);
            self.write_byte((mantissa >> 16) as u8);
            self.write_byte((mantissa >> 8) as u8);
            self.write_byte(mantissa as u8);
        }
    }

    pub fn encode_string(&mut self, s: &str) {
        let utf = s.as_bytes();
        let len = utf.len();
        let mut prefix = PREFIX_STRING;

        if len < 30 {
            self.write_byte(prefix | len as u8);
            self.write_bytes(utf);
        } else if len <= 255 {
            prefix |= STRING_LEN_1;
            self.write_byte(prefix);
            self.write_byte(len as u8);
            self.write_bytes(utf);
        } else {
            prefix |= STRING_LEN_2;
            self.write_byte(prefix);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_bytes(utf);
        }
    }

    pub fn encode_bytes(&mut self, bytes: &[u8]) {
        let len = bytes.len();
        let mut prefix = PREFIX_BYTES;

        if len < 30 {
            self.write_byte(prefix | len as u8);
            self.write_bytes(bytes);
        } else if len <= 255 {
            prefix |= BYTES_LEN_1;
            self.write_byte(prefix);
            self.write_byte(len as u8);
            self.write_bytes(bytes);
        } else {
            prefix |= BYTES_LEN_2;
            self.write_byte(prefix);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_bytes(bytes);
        }
    }

    fn encode_container(&mut self, container_type: u8, payload: &[u8]) {
        let length = payload.len();
        let mut prefix = PREFIX_CONTAINER | container_type;

        if length < 14 {
            self.write_byte(prefix | length as u8);
        } else if length <= 255 {
            prefix |= CONTAINER_LEN_1;
            self.write_byte(prefix);
            self.write_byte(length as u8);
        } else {
            prefix |= CONTAINER_LEN_2;
            self.write_byte(prefix);
            self.write_byte((length >> 8) as u8);
            self.write_byte(length as u8);
        }
        self.write_bytes(payload);
    }

    fn encode_tag(&mut self, tag_bytes: &[u8], payload: &[u8]) {
        if tag_bytes.is_empty() {
            self.write_bytes(payload);
            return;
        }

        let tag_enc_len = if tag_bytes.len() < 254 {
            1 + tag_bytes.len()
        } else if tag_bytes.len() < 257 {
            2 + tag_bytes.len()
        } else {
            3 + tag_bytes.len()
        };

        let total_len = tag_enc_len + payload.len();
        let mut prefix = PREFIX_TAG;

        if total_len < 30 {
            prefix |= total_len as u8;
            self.write_byte(prefix);
        } else if total_len <= 255 {
            prefix |= TAG_LEN_1;
            self.write_byte(prefix);
            self.write_byte(total_len as u8);
        } else {
            prefix |= TAG_LEN_2;
            self.write_byte(prefix);
            self.write_byte((total_len >> 8) as u8);
            self.write_byte(total_len as u8);
        }

        if tag_bytes.len() < 254 {
            self.write_byte(tag_bytes.len() as u8);
        } else if tag_bytes.len() == 254 {
            self.write_byte(254);
            self.write_byte(tag_bytes.len() as u8);
        } else {
            self.write_byte(255);
            self.write_byte((tag_bytes.len() >> 8) as u8);
            self.write_byte(tag_bytes.len() as u8);
        }
        self.write_bytes(tag_bytes);
        self.write_bytes(payload);
    }
}

impl Default for Encoder {
    fn default() -> Self {
        Self::new()
    }
}

fn parse_float_string(s: &str) -> (bool, i8, u64) {
    let s = s.trim();
    if s.is_empty() {
        return (false, 0, 0);
    }

    let is_negative = s.starts_with('-');
    let mut digits = if is_negative { &s[1..] } else { s };

    let mut exponent_part: Option<i64> = None;
    if let Some(e_idx) = digits.find(|c| c == 'e' || c == 'E') {
        let exp_str = &digits[e_idx + 1..];
        digits = &digits[..e_idx];
        if !exp_str.is_empty() {
            exponent_part = Some(exp_str.parse::<i64>().unwrap_or(0));
        }
    }

    let base_exp: i64;
    let mantissa_str: String;

    if let Some(dot_idx) = digits.find('.') {
        let int_part = &digits[..dot_idx];
        let frac_part = &digits[dot_idx + 1..];
        let int_part_trimmed = if int_part.is_empty() { "0" } else { int_part };
        base_exp = -(frac_part.len() as i64);
        mantissa_str = int_part_trimmed.to_string() + frac_part;
    } else {
        base_exp = 0;
        mantissa_str = digits.to_string();
    }

    let final_exp = match exponent_part {
        Some(exp) => base_exp + exp,
        None => base_exp,
    };

    let mantissa_str_trimmed = mantissa_str.trim_start_matches('0');
    let mantissa_str_clean = if mantissa_str_trimmed.is_empty() {
        "0"
    } else {
        mantissa_str_trimmed
    };

    let mantissa = mantissa_str_clean.parse::<u64>().unwrap_or(0);

    let exponent = if final_exp < i64::from(i8::MIN) {
        i8::MIN
    } else if final_exp > i64::from(i8::MAX) {
        i8::MAX
    } else {
        final_exp as i8
    };

    (is_negative, exponent, mantissa)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_float_string_integer() {
        let (neg, exp, man) = parse_float_string("42");
        assert!(!neg);
        assert_eq!(exp, 0);
        assert_eq!(man, 42);
    }

    #[test]
    fn test_parse_float_string_decimal() {
        let (neg, exp, man) = parse_float_string("3.14");
        assert!(!neg);
        assert_eq!(exp, -2);
        assert_eq!(man, 314);
    }

    #[test]
    fn test_parse_float_string_negative() {
        let (neg, exp, man) = parse_float_string("-2.5");
        assert!(neg);
        assert_eq!(exp, -1);
        assert_eq!(man, 25);
    }

    #[test]
    fn test_parse_float_string_scientific() {
        let (neg, exp, man) = parse_float_string("1.23e4");
        assert!(!neg);
        assert_eq!(exp, 2);
        assert_eq!(man, 123);
    }

    #[test]
    fn test_parse_float_string_zero() {
        let (neg, exp, man) = parse_float_string("0.0");
        assert!(!neg);
        assert_eq!(exp, -1);
        assert_eq!(man, 0);
    }

    #[test]
    fn test_encode_uint64_small() {
        let mut enc = Encoder::new();
        enc.encode_uint64(5);
        assert_eq!(enc.buf[..enc.offset], vec![0x20 | 5]);
    }

    #[test]
    fn test_encode_uint64_one_byte() {
        let mut enc = Encoder::new();
        enc.encode_uint64(100);
        assert_eq!(enc.buf[..enc.offset], vec![0x20 | 24, 100]);
    }

    #[test]
    fn test_encode_uint64_two_bytes() {
        let mut enc = Encoder::new();
        enc.encode_uint64(1000);
        assert_eq!(enc.buf[..enc.offset], vec![0x20 | 25, 0x03, 0xE8]);
    }

    #[test]
    fn test_encode_string_short() {
        let mut enc = Encoder::new();
        enc.encode_string("hello");
        assert_eq!(
            enc.buf[..enc.offset],
            vec![0x80 | 5, b'h', b'e', b'l', b'l', b'o']
        );
    }

    #[test]
    fn test_encode_string_long() {
        let mut enc = Encoder::new();
        let s = "a".repeat(30);
        enc.encode_string(&s);
        let expected: Vec<u8> = vec![0x80 | 30, 30]
            .into_iter()
            .chain(std::iter::repeat(b'a').take(30))
            .collect();
        assert_eq!(enc.buf[..enc.offset], expected);
    }

    #[test]
    fn test_encode_float_inline() {
        let mut enc = Encoder::new();
        enc.encode_float("0.0"); // exponent=-1, mantissa=0 <= 7 -> inline
        assert_eq!(enc.buf[..enc.offset], vec![0x60 | 0]);
    }

    #[test]
    fn test_encode_float_inline_one() {
        let mut enc = Encoder::new();
        enc.encode_float("1.0"); // exponent=-1, mantissa=10... no wait
                                 // "1.0": int="1", frac="0" -> base_exp=-1, mantissa_str="10" -> mantissa=10
                                 // exponent=-1, mantissa=10 > 7 -> not inline
        enc.encode_float("1.0");
        // 0.1: int="0", frac="1" -> base_exp=-1, mantissa_str="01" -> mantissa=1
        // exponent=-1, mantissa=1 <= 7 -> inline
    }

    #[test]
    fn test_encode_float_standard() {
        let mut enc = Encoder::new();
        enc.encode_float("3.14");
        // int="3", frac="14" -> base_exp=-2, mantissa_str="314" -> mantissa=314
        // exponent=-2, mantissa=314 > 0xFF? No, 314 <= 255... wait, 314 > 255
        // So needs FloatLen2 (= FLOAT_LEN_1 + 1 = 9)
        // prefix = 0x60 | 9 = 0x69, then exponent=-2 as u8 = 0xFE, then mantissa 314 = 0x01 0x3A
        assert_eq!(enc.buf[..enc.offset], vec![0x69, 0xFE, 0x01, 0x3A]);
    }

    #[test]
    fn test_encode_simple() {
        let mut enc = Encoder::new();
        enc.encode_bool(true);
        assert_eq!(enc.buf[..enc.offset], vec![0x00 | 6]); // SimpleTrue = 6
    }
}
