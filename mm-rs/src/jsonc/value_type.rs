#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ValueType {
    Unknown,
    Doc,
    Slice,
    Array,
    Struct,
    Map,
    String,
    Bytes,
    Bool,
    Int,
    Int8,
    Int16,
    Int32,
    Int64,
    Uint,
    Uint8,
    Uint16,
    Uint32,
    Uint64,
    Float32,
    Float64,
    BigInt,
    DateTime,
    Date,
    Time,
    UUID,
    Decimal,
    IP,
    URL,
    Email,
    Enum,
    Image,
    Video,
}

impl ValueType {
    pub fn from_str(s: &str) -> Self {
        match s.to_lowercase().as_str() {
            "unknown" => ValueType::Unknown,
            "doc" => ValueType::Doc,
            "slice" => ValueType::Slice,
            "arr" => ValueType::Array,
            "obj" => ValueType::Struct,
            "map" => ValueType::Map,
            "str" => ValueType::String,
            "bytes" => ValueType::Bytes,
            "bool" => ValueType::Bool,
            "i" => ValueType::Int,
            "i8" => ValueType::Int8,
            "i16" => ValueType::Int16,
            "i32" => ValueType::Int32,
            "i64" => ValueType::Int64,
            "u" => ValueType::Uint,
            "u8" => ValueType::Uint8,
            "u16" => ValueType::Uint16,
            "u32" => ValueType::Uint32,
            "u64" => ValueType::Uint64,
            "f32" => ValueType::Float32,
            "f64" => ValueType::Float64,
            "bi" => ValueType::BigInt,
            "datetime" => ValueType::DateTime,
            "date" => ValueType::Date,
            "time" => ValueType::Time,
            "uuid" => ValueType::UUID,
            "decimal" => ValueType::Decimal,
            "ip" => ValueType::IP,
            "url" => ValueType::URL,
            "email" => ValueType::Email,
            "enum" => ValueType::Enum,
            "image" => ValueType::Image,
            "video" => ValueType::Video,
            _ => ValueType::Unknown,
        }
    }

    pub fn to_str(&self) -> &str {
        match self {
            ValueType::Unknown => "unknown",
            ValueType::Doc => "doc",
            ValueType::Slice => "slice",
            ValueType::Array => "arr",
            ValueType::Struct => "obj",
            ValueType::Map => "map",
            ValueType::String => "str",
            ValueType::Bytes => "bytes",
            ValueType::Bool => "bool",
            ValueType::Int => "i",
            ValueType::Int8 => "i8",
            ValueType::Int16 => "i16",
            ValueType::Int32 => "i32",
            ValueType::Int64 => "i64",
            ValueType::Uint => "u",
            ValueType::Uint8 => "u8",
            ValueType::Uint16 => "u16",
            ValueType::Uint32 => "u32",
            ValueType::Uint64 => "u64",
            ValueType::Float32 => "f32",
            ValueType::Float64 => "f64",
            ValueType::BigInt => "bi",
            ValueType::DateTime => "datetime",
            ValueType::Date => "date",
            ValueType::Time => "time",
            ValueType::UUID => "uuid",
            ValueType::Decimal => "decimal",
            ValueType::IP => "ip",
            ValueType::URL => "url",
            ValueType::Email => "email",
            ValueType::Enum => "enum",
            ValueType::Image => "image",
            ValueType::Video => "video",
        }
    }

    pub fn needs_quotes(&self) -> bool {
        matches!(
            self,
            ValueType::String
                | ValueType::Bytes
                | ValueType::DateTime
                | ValueType::Date
                | ValueType::Time
                | ValueType::UUID
                | ValueType::IP
                | ValueType::URL
                | ValueType::Email
                | ValueType::Enum
        )
    }
}