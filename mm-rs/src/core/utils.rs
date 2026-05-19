use std::net::IpAddr;
use chrono::{DateTime, TimeZone, Utc};

pub const DEFAULT_TIME: DateTime<Utc> = DateTime::from_timestamp(0, 0).unwrap();

pub fn camel_to_snake(s: &str) -> String {
    let mut result = String::new();
    for (i, c) in s.chars().enumerate() {
        if c.is_uppercase() && i > 0 {
            result.push('_');
        }
        result.push(c.to_ascii_lowercase());
    }
    result
}

pub fn snake_to_camel(s: &str) -> String {
    let mut result = String::new();
    let mut capitalize_next = false;
    for c in s.chars() {
        if c == '_' {
            capitalize_next = true;
        } else if capitalize_next {
            result.push(c.to_ascii_uppercase());
            capitalize_next = false;
        } else {
            result.push(c);
        }
    }
    result
}

pub fn format_float32(f: f32) -> String {
    ryu::Buffer::new().format_finite(f).to_string()
}

pub fn format_float64(f: f64) -> String {
    ryu::Buffer::new().format_finite(f).to_string()
}

pub fn encode_big_int(buf: &mut Vec<u8>, s: &str) -> usize {
    let bytes = s.as_bytes();
    let len = bytes.len();
    buf.extend_from_slice(bytes);
    len
}

pub fn decode_big_int(data: &[u8], len: usize) -> Result<String, std::io::Error> {
    if len > data.len() {
        return Err(std::io::Error::new(std::io::ErrorKind::UnexpectedEof, "unexpected eof"));
    }
    let s = String::from_utf8_lossy(&data[..len]).to_string();
    Ok(s)
}

pub fn bytes_to_uuid_string(bytes: [u8; 16]) -> String {
    format!(
        "{:02x}{:02x}{:02x}{:02x}-{:02x}{:02x}-{:02x}{:02x}-{:02x}{:02x}-{:02x}{:02x}{:02x}{:02x}{:02x}{:02x}",
        bytes[0], bytes[1], bytes[2], bytes[3],
        bytes[4], bytes[5],
        bytes[6], bytes[7],
        bytes[8], bytes[9],
        bytes[10], bytes[11], bytes[12], bytes[13], bytes[14], bytes[15]
    )
}

pub fn parse_ip(s: &str) -> Option<IpAddr> {
    s.parse().ok()
}