use crate::jsonc::ast::{Node, Object, Array, Value, Field, ValueData};
use crate::jsonc::tag::Tag;
use crate::jsonc::ValueType;
use crate::mm::prefix::Prefix;
use crate::mm::simple_value::SimpleValue;

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
            return Err(std::io::Error::new(std::io::ErrorKind::UnexpectedEof, "unexpected eof"));
        }
        let b = self.data[self.offset];
        self.offset += 1;
        Ok(b)
    }

    fn read_bytes(&mut self, n: usize) -> Result<&[u8], std::io::Error> {
        if self.offset + n > self.data.len() {
            return Err(std::io::Error::new(std::io::ErrorKind::UnexpectedEof, "unexpected eof"));
        }
        let start = self.offset;
        self.offset += n;
        Ok(&self.data[start..start + n])
    }

    pub fn decode(&mut self) -> Result<Node, std::io::Error> {
        let tag = Tag::new();
        self.decode_node(&tag)
    }

    fn decode_node(&mut self, tag: &Tag) -> Result<Node, std::io::Error> {
        let b = self.read_byte()?;
        let prefix = Prefix::from_byte(b).ok_or_else(|| {
            std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid prefix")
        })?;

        match prefix {
            Prefix::Tag => self.decode_tag(b),
            Prefix::Simple => self.decode_simple(b, tag),
            Prefix::PositiveInt => self.decode_positive_int(b, tag),
            Prefix::NegativeInt => self.decode_negative_int(b, tag),
            Prefix::PrefixFloat => self.decode_float(b, tag),
            Prefix::PrefixString => self.decode_string(b, tag),
            Prefix::PrefixBytes => self.decode_bytes(b, tag),
            Prefix::Container => self.decode_container(b, tag),
        }
    }

    fn decode_tag(&mut self, prefix: u8) -> Result<Node, std::io::Error> {
        let l1 = match prefix & 0x03 {
            0 => 0,
            1 => {
                let l = self.read_byte()?;
                l as usize
            }
            2 => {
                let l = self.read_bytes(2)?;
                ((l[0] as usize) << 8) | (l[1] as usize)
            }
            _ => return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid tag length")),
        };

        let mut tag = Tag::new();

        let b = self.read_byte()?;
        let mut l = b as usize;

        if l < 254 {
        } else if l < 257 {
            l = self.read_byte()? as usize;
        } else {
            let l3 = self.read_bytes(2)?;
            l = ((l3[0] as usize) << 8) | (l3[1] as usize);
        }

        for _ in 0..l {
            self.decode_tag_bytes(&mut tag)?;
        }

        if tag.is_null {
            let data = match tag.value_type {
                ValueType::Int => ValueData::Int(0),
                ValueType::Float64 => ValueData::Float(0.0),
                ValueType::Bool => ValueData::Bool(false),
                ValueType::String => ValueData::String(String::new()),
                ValueType::Bytes => ValueData::Bytes(vec![]),
                ValueType::DateTime => ValueData::String("1970-01-01T00:00:00Z".to_string()),
                _ => ValueData::Null,
            };
            Ok(Node::Value(Value {
                data,
                text: String::new(),
                tag: Some(tag),
            }))
        } else {
            self.decode_node(&tag)
        }
    }

    fn decode_tag_bytes(&mut self, _tag: &mut Tag) -> Result<usize, std::io::Error> {
        Ok(0)
    }

    fn decode_simple(&mut self, prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
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
            tag: Some(Tag::new()),
        }))
    }

    fn decode_positive_int(&mut self, prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
        let l1 = match prefix & 0x07 {
            0 => 0,
            1 => {
                let l = self.read_byte()?;
                l as usize
            }
            _ => return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid int length")),
        };

        let mut v = 0u64;
        for i in 0..l1 {
            let b = self.read_byte()?;
            v |= (b as u64) << (i * 8);
        }

        let data = ValueData::Int(v as i64);
        let text = v.to_string();

        Ok(Node::Value(Value {
            data,
            text,
            tag: Some(Tag::new()),
        }))
    }

    fn decode_negative_int(&mut self, prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
        let l1 = match prefix & 0x07 {
            0 => 0,
            1 => {
                let l = self.read_byte()?;
                l as usize
            }
            _ => return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid int length")),
        };

        let mut v = 0u64;
        for i in 0..l1 {
            let b = self.read_byte()?;
            v |= (b as u64) << (i * 8);
        }

        let data = ValueData::Int(-(v as i64));
        let text = format!("-{}", v);

        Ok(Node::Value(Value {
            data,
            text,
            tag: Some(Tag::new()),
        }))
    }

    fn decode_float(&mut self, _prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
        let s = String::from_utf8_lossy(self.read_bytes(8)?).to_string();
        let f = s.parse::<f64>().unwrap_or(0.0);
        Ok(Node::Value(Value {
            data: ValueData::Float(f),
            text: ryu::Buffer::new().format_finite(f).to_string(),
            tag: Some(Tag::new()),
        }))
    }

    fn decode_string(&mut self, prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
        let l1 = match prefix & 0x03 {
            0 => 0,
            1 => {
                let l = self.read_byte()?;
                l as usize
            }
            2 => {
                let l = self.read_bytes(2)?;
                ((l[0] as usize) << 8) | (l[1] as usize)
            }
            _ => return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid string length")),
        };

        let s = if l1 > 0 {
            String::from_utf8_lossy(self.read_bytes(l1)?).to_string()
        } else {
            String::new()
        };

        Ok(Node::Value(Value {
            data: ValueData::String(s.clone()),
            text: s,
            tag: Some(Tag::new()),
        }))
    }

    fn decode_bytes(&mut self, prefix: u8, _tag: &Tag) -> Result<Node, std::io::Error> {
        let l1 = match prefix & 0x03 {
            0 => 0,
            1 => {
                let l = self.read_byte()?;
                l as usize
            }
            2 => {
                let l = self.read_bytes(2)?;
                ((l[0] as usize) << 8) | (l[1] as usize)
            }
            _ => return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "invalid bytes length")),
        };

        let bytes = if l1 > 0 {
            self.read_bytes(l1)?.to_vec()
        } else {
            vec![]
        };

        Ok(Node::Value(Value {
            data: ValueData::Bytes(bytes.clone()),
            text: format!("{:?}", bytes),
            tag: Some(Tag::new()),
        }))
    }

    fn decode_container(&mut self, prefix: u8, tag: &Tag) -> Result<Node, std::io::Error> {
        let is_array = (prefix & 0x01) != 0;
        if is_array {
            self.decode_array(tag)
        } else {
            self.decode_object(tag)
        }
    }

    fn decode_array(&mut self, tag: &Tag) -> Result<Node, std::io::Error> {
        let l = self.read_byte()?;
        let len = l as usize;

        let mut items = Vec::new();
        for _ in 0..len {
            let item = self.decode_node(&Tag::new())?;
            items.push(item);
        }

        Ok(Node::Array(Array {
            items,
            tag: Some(Tag::new()),
        }))
    }

    fn decode_object(&mut self, tag: &Tag) -> Result<Node, std::io::Error> {
        let mut fields = Vec::new();

        let l = self.read_byte()?;
        let _len = l as usize;

        loop {
            if self.offset >= self.data.len() {
                break;
            }
            let b = self.data[self.offset];
            if b == 0x00 || b == 0x01 {
                break;
            }

            let key_node = self.decode_node(&Tag::new())?;
            if let Node::Value(v) = key_node {
                let key = v.text;
                let value = self.decode_node(&Tag::new())?;
                fields.push(Field { key, value });
            } else {
                break;
            }
        }

        Ok(Node::Object(Object {
            fields,
            tag: Some(Tag::new()),
        }))
    }
}