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

// Tag key prefixes (key = byte & 0xF8)
pub const TAG_IS_NULL: u8 = 0 << 3; // 0x00
pub const TAG_EXAMPLE: u8 = 1 << 3; // 0x08
pub const TAG_DESC: u8 = 2 << 3; // 0x10
pub const TAG_TYPE: u8 = 3 << 3; // 0x18
pub const TAG_RAW: u8 = 4 << 3; // 0x20
pub const TAG_NULLABLE: u8 = 5 << 3; // 0x28
pub const TAG_ALLOW_EMPTY: u8 = 6 << 3; // 0x30
pub const TAG_UNIQUE: u8 = 7 << 3; // 0x38
pub const TAG_DEFAULT: u8 = 8 << 3; // 0x40
pub const TAG_MIN: u8 = 9 << 3; // 0x48
pub const TAG_MAX: u8 = 10 << 3; // 0x50
pub const TAG_SIZE: u8 = 11 << 3; // 0x58
pub const TAG_ENUM: u8 = 12 << 3; // 0x60
pub const TAG_PATTERN: u8 = 13 << 3; // 0x68
pub const TAG_LOCATION: u8 = 14 << 3; // 0x70
pub const TAG_VERSION: u8 = 15 << 3; // 0x78
pub const TAG_MIME: u8 = 16 << 3; // 0x80

pub const TAG_CHILD_DESC: u8 = 17 << 3; // 0x88
pub const TAG_CHILD_TYPE: u8 = 18 << 3; // 0x90
pub const TAG_CHILD_RAW: u8 = 19 << 3; // 0x98
pub const TAG_CHILD_NULLABLE: u8 = 20 << 3; // 0xA0
pub const TAG_CHILD_ALLOW_EMPTY: u8 = 21 << 3; // 0xA8
pub const TAG_CHILD_UNIQUE: u8 = 22 << 3; // 0xB0
pub const TAG_CHILD_DEFAULT: u8 = 23 << 3; // 0xB8
pub const TAG_CHILD_MIN: u8 = 24 << 3; // 0xC0
pub const TAG_CHILD_MAX: u8 = 25 << 3; // 0xC8
pub const TAG_CHILD_SIZE: u8 = 26 << 3; // 0xD0
pub const TAG_CHILD_ENUM: u8 = 27 << 3; // 0xD8
pub const TAG_CHILD_PATTERN: u8 = 28 << 3; // 0xE0
pub const TAG_CHILD_LOCATION: u8 = 29 << 3; // 0xE8
pub const TAG_CHILD_VERSION: u8 = 30 << 3; // 0xF0
pub const TAG_CHILD_MIME: u8 = 31 << 3; // 0xF8

pub const TAG_KEY_MASK: u8 = 0xF8;
pub const TAG_PAYLOAD_MASK: u8 = 0x07;