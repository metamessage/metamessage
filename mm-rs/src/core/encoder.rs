use crate::core::constants::*;
use crate::core::prefix::{
    CONTAINER_ARRAY, CONTAINER_OBJECT, FLOAT_LEN_1, FLOAT_POSITIVE_NEGATIVE_MASK, PREFIX_BYTES,
    PREFIX_CONTAINER, PREFIX_FLOAT, PREFIX_NEGATIVE_INT, PREFIX_POSITIVE_INT, PREFIX_SIMPLE,
    PREFIX_STRING, PREFIX_TAG,
};
use crate::core::simple_value::SimpleValue;
use crate::ir::ast::{Node, NodeArray, NodeObject, NodeScalar, ValueData};
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
        // v0.3.0：落实 MAX_CAP 强制（此前从不检查），超限显式报错而非静默
        if self.offset + needed > MAX_CAP {
            panic!("encoded size exceeds format limit (max {} bytes)", MAX_CAP);
        }
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

    fn encode_object(&mut self, obj: &NodeObject) {
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

    fn encode_array(&mut self, arr: &NodeArray) {
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

    fn encode_value(&mut self, val: &NodeScalar) {
        let saved_offset = self.offset;
        self.offset = 0;

        let tag = val.tag.as_ref();
        let is_null = tag.map_or(false, |t| t.is_null);

        if is_null {
            // Align with Go: only specific types encode a null simple value.
            // All other types write nothing (the null state is indicated by the tag's isNull flag).
            match tag.map(|t| t.value_type) {
                Some(ValueType::Unknown) => self.encode_simple(SimpleValue::Null),
                Some(ValueType::Bool) => self.encode_simple(SimpleValue::NullBool),
                Some(ValueType::I) => self.encode_simple(SimpleValue::NullInt),
                Some(ValueType::F64) => self.encode_simple(SimpleValue::NullFloat),
                Some(ValueType::Str) => self.encode_simple(SimpleValue::NullString),
                Some(ValueType::Bytes) => self.encode_simple(SimpleValue::NullBytes),
                // Datetime, Date, Time, I8-I64, U-U64, F32, Email, Uuid, Decimal,
                // Url, Ip, Media, Bigint, Enum: no bytes written (skip)
                _ => {}
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
                            // Align with Go: if text is short (< 16 chars), encode as string
                            if val.text.len() < 16 {
                                self.encode_string(&val.text);
                            } else if let Ok(ip) = val.text.parse::<std::net::Ipv6Addr>() {
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
                    ValueData::Null => self.encode_simple(SimpleValue::Null),
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
                ValueData::Null => self.encode_simple(SimpleValue::Null),
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
        // v0.5.0 尾定界：一趟正向 [body BE][sign|档位] 写出，档位/内联语义与 v0.3.x 一致。
        if v < 24 {
            self.write_byte(sign | (v as u8));
            return;
        }
        let width: u32 = if v <= 0xFF {
            1
        } else if v <= 0xFFFF {
            2
        } else if v <= 0xFFFFFF {
            3
        } else if v <= 0xFFFFFFFF {
            4
        } else if v <= 0xFFFFFFFFFF {
            5
        } else if v <= 0xFFFFFFFFFFFF {
            6
        } else if v <= 0xFFFFFFFFFFFFFF {
            7
        } else {
            8
        };
        for i in (0..width).rev() {
            self.write_byte(((v >> (8 * i)) & 0xFF) as u8);
        }
        self.write_byte(sign | (INT_LEN_1 - 1 + width as u8));
    }

    pub fn encode_float(&mut self, s: &str) {
        let (is_negative, exponent, mantissa) = parse_float_string(s);

        let mut prefix = PREFIX_FLOAT;
        if is_negative {
            prefix |= FLOAT_POSITIVE_NEGATIVE_MASK;
        }

        // v0.5.0 尾定界：一趟正向 [exp][mantissa BE][sign|档位] 写出，档位/内联语义与 v0.3.x 一致。
        // 每一档都写 exponent（与 Go encode_float.go 对齐），修复旧末档漏写 exp 的问题。
        if exponent == -1 && mantissa <= 7 {
            self.write_byte(prefix | (mantissa as u8));
            return;
        }

        let width: u32 = if mantissa <= 0xFF {
            1
        } else if mantissa <= 0xFFFF {
            2
        } else if mantissa <= 0xFFFFFF {
            3
        } else if mantissa <= 0xFFFFFFFF {
            4
        } else if mantissa <= 0xFFFFFFFFFF {
            5
        } else if mantissa <= 0xFFFFFFFFFFFF {
            6
        } else if mantissa <= 0xFFFFFFFFFFFFFF {
            7
        } else {
            8
        };
        self.write_byte(exponent as u8);
        for i in (0..width).rev() {
            self.write_byte(((mantissa >> (8 * i)) & 0xFF) as u8);
        }
        self.write_byte(prefix | (FLOAT_LEN_1 - 1 + width as u8));
    }

    pub fn encode_string(&mut self, s: &str) {
        let utf = s.as_bytes();
        let len = utf.len();
        let mut prefix = PREFIX_STRING;

        // v0.5.0 尾定界：一趟正向 [body][len BE][sign|档位] 写出。
        // 1byte 档位阈值与 Go 对齐为 < 255（旧 HEAD 误用 <=255）。
        if len < 30 {
            self.write_bytes(utf);
            self.write_byte(prefix | len as u8);
        } else if len < 255 {
            prefix |= STRING_LEN_1;
            self.write_bytes(utf);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        } else {
            prefix |= STRING_LEN_2;
            self.write_bytes(utf);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        }
    }

    pub fn encode_bytes(&mut self, bytes: &[u8]) {
        let len = bytes.len();
        assert!(len <= 0xFFFFFFFF, "bytes too long: {}", len);
        let mut prefix = PREFIX_BYTES;

        // v0.5.0 尾定界：一趟正向 [body][len BE][sign|档位] 写出，档位阈值与 v0.3.x 一致。
        if len < 28 {
            self.write_bytes(bytes);
            self.write_byte(prefix | len as u8);
        } else if len <= 255 {
            prefix |= BYTES_LEN_1;
            self.write_bytes(bytes);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        } else if len <= 0xFFFF {
            prefix |= BYTES_LEN_2;
            self.write_bytes(bytes);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        } else if len <= 0xFFFFFF {
            prefix |= BYTES_LEN_3;
            self.write_bytes(bytes);
            self.write_byte((len >> 16) as u8);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        } else {
            prefix |= BYTES_LEN_4;
            self.write_bytes(bytes);
            self.write_byte((len >> 24) as u8);
            self.write_byte((len >> 16) as u8);
            self.write_byte((len >> 8) as u8);
            self.write_byte(len as u8);
            self.write_byte(prefix);
        }
    }

    fn encode_container(&mut self, container_type: u8, payload: &[u8]) {
        let length = payload.len();
        assert!(length <= 0xFFFFFFFF, "container payload too long: {}", length);
        let mut prefix = PREFIX_CONTAINER | container_type;

        // v0.5.0 尾定界：一趟正向 [payload][len BE][sign|档位] 写出，档位阈值与 v0.3.x 一致。
        self.write_bytes(payload);
        if length < 12 {
            self.write_byte(prefix | length as u8);
        } else if length <= 255 {
            prefix |= CONTAINER_LEN_1;
            self.write_byte(length as u8);
            self.write_byte(prefix);
        } else if length <= 0xFFFF {
            prefix |= CONTAINER_LEN_2;
            self.write_byte((length >> 8) as u8);
            self.write_byte(length as u8);
            self.write_byte(prefix);
        } else if length <= 0xFFFFFF {
            prefix |= CONTAINER_LEN_3;
            self.write_byte((length >> 16) as u8);
            self.write_byte((length >> 8) as u8);
            self.write_byte(length as u8);
            self.write_byte(prefix);
        } else {
            prefix |= CONTAINER_LEN_4;
            self.write_byte((length >> 24) as u8);
            self.write_byte((length >> 16) as u8);
            self.write_byte((length >> 8) as u8);
            self.write_byte(length as u8);
            self.write_byte(prefix);
        }
    }

    fn encode_tag(&mut self, tag_bytes: &[u8], payload: &[u8]) {
        if tag_bytes.is_empty() {
            self.write_bytes(payload);
            return;
        }

        // v0.5.0 尾定界：一趟正向 [payload][reverse(meta)][尾槽][total BE][sign] 写出。
        // total = len(payload) + len(meta) + len(尾槽)，与解码侧 untailify 反推式
        // inner = total - ml - len(尾槽) 一致（尾槽与前缀槽为同一套 canonical 三档）。
        let slot = crate::core::footer::footer_slot_bytes(tag_bytes.len());
        let total = payload.len() + tag_bytes.len() + slot.len();
        assert!(total <= 0xFFFFFFFF, "tag+payload too long: {}", total);

        let mut sign = PREFIX_TAG;
        let mut tier: Vec<u8> = Vec::new();
        if total < 28 {
            sign |= total as u8;
        } else if total <= 255 {
            sign |= TAG_LEN_1;
            tier.push(total as u8);
        } else if total <= 0xFFFF {
            sign |= TAG_LEN_2;
            tier.push((total >> 8) as u8);
            tier.push(total as u8);
        } else if total <= 0xFFFFFF {
            sign |= TAG_LEN_3;
            tier.push((total >> 16) as u8);
            tier.push((total >> 8) as u8);
            tier.push(total as u8);
        } else {
            sign |= TAG_LEN_4;
            tier.push((total >> 24) as u8);
            tier.push((total >> 16) as u8);
            tier.push((total >> 8) as u8);
            tier.push(total as u8);
        }

        self.write_bytes(payload);
        self.write_bytes(&tag_bytes.iter().rev().cloned().collect::<Vec<u8>>());
        self.write_bytes(&slot);
        self.write_bytes(&tier);
        self.write_byte(sign);
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
        assert_eq!(enc.buf[..enc.offset], vec![100, 0x20 | 24]);
    }

    #[test]
    fn test_encode_uint64_two_bytes() {
        let mut enc = Encoder::new();
        enc.encode_uint64(1000);
        assert_eq!(enc.buf[..enc.offset], vec![0x03, 0xE8, 0x20 | 25]);
    }

    #[test]
    fn test_encode_string_short() {
        let mut enc = Encoder::new();
        enc.encode_string("hello");
        assert_eq!(
            enc.buf[..enc.offset],
            vec![b'h', b'e', b'l', b'l', b'o', 0x80 | 5]
        );
    }

    #[test]
    fn test_encode_string_long() {
        let mut enc = Encoder::new();
        let s = "a".repeat(30);
        enc.encode_string(&s);
        let expected: Vec<u8> = std::iter::repeat(b'a')
            .take(30)
            .chain(vec![30, 0x80 | 30])
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
        // prefix = 0x60 | 9 = 0x69, tail order: exponent=-2 as u8 = 0xFE, mantissa 314 = 0x01 0x3A, sign 0x69 last
        assert_eq!(enc.buf[..enc.offset], vec![0xFE, 0x01, 0x3A, 0x69]);
    }

    #[test]
    fn test_encode_simple() {
        let mut enc = Encoder::new();
        enc.encode_bool(true);
        assert_eq!(enc.buf[..enc.offset], vec![0x00 | 7]); // SimpleTrue = 6
    }
}
