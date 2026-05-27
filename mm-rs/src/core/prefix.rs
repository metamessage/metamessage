#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Prefix {
    Simple,
    PositiveInt,
    NegativeInt,
    PrefixFloat,
    PrefixString,
    PrefixBytes,
    Container,
    Tag,
}

impl Prefix {
    pub fn from_byte(b: u8) -> Option<Self> {
        match b & PREFIX_MASK {
            0x00 => Some(Prefix::Simple),
            0x20 => Some(Prefix::PositiveInt),
            0x40 => Some(Prefix::NegativeInt),
            0x60 => Some(Prefix::PrefixFloat),
            0x80 => Some(Prefix::PrefixString),
            0xA0 => Some(Prefix::PrefixBytes),
            0xC0 => Some(Prefix::Container),
            0xE0 => Some(Prefix::Tag),
            _ => None,
        }
    }

    pub fn is_array(b: u8) -> bool {
        (b & CONTAINER_TYPE_MASK) != 0 && (b & PREFIX_MASK) == PREFIX_CONTAINER
    }
}

pub const PREFIX_MASK: u8 = 0xE0;
pub const SUFFIX_MASK: u8 = 0x1F;

pub const PREFIX_SIMPLE: u8 = 0x00;
pub const PREFIX_POSITIVE_INT: u8 = 0x20;
pub const PREFIX_NEGATIVE_INT: u8 = 0x40;
pub const PREFIX_FLOAT: u8 = 0x60;
pub const PREFIX_STRING: u8 = 0x80;
pub const PREFIX_BYTES: u8 = 0xA0;
pub const PREFIX_CONTAINER: u8 = 0xC0;
pub const PREFIX_TAG: u8 = 0xE0;

pub const CONTAINER_TYPE_MASK: u8 = 0x10;
pub const CONTAINER_OBJECT: u8 = 0x00;
pub const CONTAINER_ARRAY: u8 = 0x10;

pub const FLOAT_POSITIVE_NEGATIVE_MASK: u8 = 0x10;
pub const FLOAT_LEN_MASK: u8 = 0x0F;
pub const FLOAT_LEN_1: u8 = 0x08;