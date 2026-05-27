#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ValueType {
    Unknown,
    Doc,
    Vec,
    Arr,
    Obj,
    Map,
    Str,
    Bytes,
    Bool,
    I,
    I8,
    I16,
    I32,
    I64,
    U,
    U8,
    U16,
    U32,
    U64,
    F32,
    F64,
    Bigint,
    Datetime,
    Date,
    Time,
    Uuid,
    Decimal,
    Ip,
    Url,
    Email,
    Enum,
    Image,
    Video,
}

impl ValueType {
    #[allow(clippy::should_implement_trait)]
    pub fn from_str(s: &str) -> Self {
        match s.to_lowercase().as_str() {
            "unknown" => ValueType::Unknown,
            "doc" => ValueType::Doc,
            "vec" => ValueType::Vec,
            "arr" => ValueType::Arr,
            "obj" => ValueType::Obj,
            "map" => ValueType::Map,
            "str" => ValueType::Str,
            "bytes" => ValueType::Bytes,
            "bool" => ValueType::Bool,
            "i" => ValueType::I,
            "i8" => ValueType::I8,
            "i16" => ValueType::I16,
            "i32" => ValueType::I32,
            "i64" => ValueType::I64,
            "u" => ValueType::U,
            "u8" => ValueType::U8,
            "u16" => ValueType::U16,
            "u32" => ValueType::U32,
            "u64" => ValueType::U64,
            "f32" => ValueType::F32,
            "f64" => ValueType::F64,
            "bigint" => ValueType::Bigint,
            "datetime" => ValueType::Datetime,
            "date" => ValueType::Date,
            "time" => ValueType::Time,
            "uuid" => ValueType::Uuid,
            "decimal" => ValueType::Decimal,
            "ip" => ValueType::Ip,
            "url" => ValueType::Url,
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
            ValueType::Vec => "vec",
            ValueType::Arr => "arr",
            ValueType::Obj => "obj",
            ValueType::Map => "map",
            ValueType::Str => "str",
            ValueType::Bytes => "bytes",
            ValueType::Bool => "bool",
            ValueType::I => "i",
            ValueType::I8 => "i8",
            ValueType::I16 => "i16",
            ValueType::I32 => "i32",
            ValueType::I64 => "i64",
            ValueType::U => "u",
            ValueType::U8 => "u8",
            ValueType::U16 => "u16",
            ValueType::U32 => "u32",
            ValueType::U64 => "u64",
            ValueType::F32 => "f32",
            ValueType::F64 => "f64",
            ValueType::Bigint => "bigint",
            ValueType::Datetime => "datetime",
            ValueType::Date => "date",
            ValueType::Time => "time",
            ValueType::Uuid => "uuid",
            ValueType::Decimal => "decimal",
            ValueType::Ip => "ip",
            ValueType::Url => "url",
            ValueType::Email => "email",
            ValueType::Enum => "enum",
            ValueType::Image => "image",
            ValueType::Video => "video",
        }
    }

    pub fn needs_quotes(&self) -> bool {
        matches!(
            self,
            ValueType::Unknown
                | ValueType::Str
                | ValueType::Bytes
                | ValueType::Datetime
                | ValueType::Date
                | ValueType::Time
                | ValueType::Uuid
                | ValueType::Ip
                | ValueType::Url
                | ValueType::Email
                | ValueType::Enum
        )
    }

    pub fn from_code(code: u8) -> Self {
        match code {
            0 => ValueType::Unknown,
            1 => ValueType::Doc,
            2 => ValueType::Vec,
            3 => ValueType::Arr,
            4 => ValueType::Obj,
            5 => ValueType::Map,
            6 => ValueType::Str,
            7 => ValueType::Bytes,
            8 => ValueType::Bool,
            9 => ValueType::I,
            10 => ValueType::I8,
            11 => ValueType::I16,
            12 => ValueType::I32,
            13 => ValueType::I64,
            14 => ValueType::U,
            15 => ValueType::U8,
            16 => ValueType::U16,
            17 => ValueType::U32,
            18 => ValueType::U64,
            19 => ValueType::F32,
            20 => ValueType::F64,
            21 => ValueType::Bigint,
            22 => ValueType::Datetime,
            23 => ValueType::Date,
            24 => ValueType::Time,
            25 => ValueType::Uuid,
            26 => ValueType::Decimal,
            27 => ValueType::Ip,
            28 => ValueType::Url,
            29 => ValueType::Email,
            30 => ValueType::Enum,
            31 => ValueType::Image,
            32 => ValueType::Video,
            _ => ValueType::Unknown,
        }
    }
}
