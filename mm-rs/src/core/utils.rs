use std::net::IpAddr;
use chrono::{DateTime, Utc};

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

struct BitReader<'a> {
    data: &'a [u8],
    pos: usize,
    total_bits: usize,
}

impl<'a> BitReader<'a> {
    fn new(data: &'a [u8], total_bits: usize) -> Self {
        BitReader { data, pos: 0, total_bits }
    }

    fn read_bit(&mut self) -> Result<i32, std::io::Error> {
        if self.pos >= self.total_bits {
            return Err(std::io::Error::new(std::io::ErrorKind::UnexpectedEof, "unexpected eof"));
        }
        let byte_idx = self.pos / 8;
        let bit_in_byte = 7 - (self.pos % 8);
        let bit = ((self.data[byte_idx] >> bit_in_byte) & 1) as i32;
        self.pos += 1;
        Ok(bit)
    }

    fn read_bits(&mut self, n: usize) -> Result<i32, std::io::Error> {
        let mut val = 0;
        for _ in 0..n {
            val = (val << 1) | self.read_bit()?;
        }
        Ok(val)
    }
}

/// Decode a bigint from bit-packed bytes.
/// The input bytes format: [digit_count (1 byte)] + [bit-packed data...]
/// Bit-packed format: sign(1 bit) + groups of 10/7/4 bits for 3/2/1 decimal digits
pub fn decode_big_int(data: &[u8]) -> Result<String, std::io::Error> {
    if data.is_empty() {
        return Ok("0".to_string());
    }
    let digit_len = data[0] as usize;
    if digit_len == 0 || data.len() <= 1 {
        return Ok("0".to_string());
    }

    let bit_data = &data[1..];
    let total_bits = bit_data.len() * 8;
    let mut reader = BitReader::new(bit_data, total_bits);

    // Read sign bit
    let sign = reader.read_bit()?;
    let neg = sign == 1;

    let mut result = String::with_capacity(digit_len);
    let mut remaining = digit_len;

    while remaining > 0 {
        if remaining >= 3 {
            let val = reader.read_bits(10)?;
            result.push_str(&format!("{:03}", val));
            remaining -= 3;
        } else if remaining == 2 {
            let val = reader.read_bits(7)?;
            result.push_str(&format!("{:02}", val));
            remaining -= 2;
        } else {
            let val = reader.read_bits(4)?;
            result.push_str(&format!("{:01}", val));
            remaining -= 1;
        }
    }

    // Trim leading zeros
    let trimmed = result.trim_start_matches('0');
    let final_str = if trimmed.is_empty() { "0" } else { trimmed };

    if neg && final_str != "0" {
        Ok(format!("-{}", final_str))
    } else {
        Ok(final_str.to_string())
    }
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