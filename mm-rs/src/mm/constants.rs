pub const DEFAULT_BUF_SIZE: usize = 1024;
pub const MAX_CAP: usize = 1024 * 1024 * 1024;

pub const MAX_1: u8 = 0xFF;
pub const MAX_2: u16 = 0xFFFF;
pub const MAX_3: u32 = 0xFFFFFF;
pub const MAX_4: u32 = 0xFFFFFFFF;
pub const MAX_5: u64 = 0xFFFFFFFFFF;
pub const MAX_6: u64 = 0xFFFFFFFFFFFF;
pub const MAX_7: u64 = 0xFFFFFFFFFFFFFF;
pub const MAX_8: u64 = 0xFFFFFFFFFFFFFFFF;

pub const INT_LEN_1: u8 = 0x80;
pub const INT_LEN_2: u8 = 0x81;
pub const INT_LEN_3: u8 = 0x82;
pub const INT_LEN_4: u8 = 0x83;
pub const INT_LEN_5: u8 = 0x84;
pub const INT_LEN_6: u8 = 0x85;
pub const INT_LEN_7: u8 = 0x86;
pub const INT_LEN_8: u8 = 0x87;
pub const INT_LEN_MASK: u8 = 0x07;

pub const STRING_LEN_1: u8 = 0x80;
pub const STRING_LEN_2: u8 = 0x81;

pub const BYTES_LEN_1: u8 = 0x80;
pub const BYTES_LEN_2: u8 = 0x81;

pub const CONTAINER_LEN_1: u8 = 0x08;
pub const CONTAINER_LEN_2: u8 = 0x09;
pub const CONTAINER_LEN_MASK: u8 = 0x07;

pub const TAG_LEN_1: u8 = 0x80;
pub const TAG_LEN_2: u8 = 0x81;

pub const FLOAT_POSITIVE_NEGATIVE_MASK: u8 = 0x10;
pub const FLOAT_LEN_MASK: u8 = 0x0F;
pub const FLOAT_LEN_1: u8 = 0x08;

pub const CONTAINER_ARRAY: u8 = 0x10;
pub const CONTAINER_MAP: u8 = 0x00;
pub const CONTAINER_TYPE_MASK: u8 = 0x10;