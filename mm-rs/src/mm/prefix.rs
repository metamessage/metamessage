use crate::mm::constants::*;

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
        let p = b & 0xF8;
        match p {
            0x00 => Some(Prefix::Simple),
            0x08 => Some(Prefix::PositiveInt),
            0x10 => Some(Prefix::NegativeInt),
            0x18 => Some(Prefix::PrefixFloat),
            0x20 => Some(Prefix::PrefixString),
            0x28 => Some(Prefix::PrefixBytes),
            0x30 => Some(Prefix::Container),
            0x38 => Some(Prefix::Tag),
            _ => None,
        }
    }

    pub fn is_array(prefix: u8) -> bool {
        (prefix & 0x10) != 0 && (prefix & 0xF8) == 0x30
    }
}

pub const PREFIX_SIMPLE: u8 = 0x00;
pub const PREFIX_POSITIVE_INT: u8 = 0x08;
pub const PREFIX_NEGATIVE_INT: u8 = 0x10;
pub const PREFIX_FLOAT: u8 = 0x18;
pub const PREFIX_STRING: u8 = 0x20;
pub const PREFIX_BYTES: u8 = 0x28;
pub const PREFIX_CONTAINER: u8 = 0x30;
pub const PREFIX_TAG: u8 = 0x38;

pub const FLOAT_POSITIVE_NEGATIVE_MASK: u8 = 0x10;
pub const FLOAT_LEN_MASK: u8 = 0x0F;
pub const FLOAT_LEN_1: u8 = 0x08;