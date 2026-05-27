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

pub const INT_LEN_1: u8 = 24;
pub const INT_LEN_2: u8 = 25;
pub const INT_LEN_3: u8 = 26;
pub const INT_LEN_4: u8 = 27;
pub const INT_LEN_5: u8 = 28;
pub const INT_LEN_6: u8 = 29;
pub const INT_LEN_7: u8 = 30;
pub const INT_LEN_8: u8 = 31;
pub const INT_LEN_MASK: u8 = 0x1F;

pub const STRING_LEN_1: u8 = 30;
pub const STRING_LEN_2: u8 = 31;

pub const BYTES_LEN_1: u8 = 30;
pub const BYTES_LEN_2: u8 = 31;

pub const CONTAINER_LEN_1: u8 = 14;
pub const CONTAINER_LEN_2: u8 = 15;
pub const CONTAINER_LEN_MASK: u8 = 0x0F;

pub const TAG_LEN_1: u8 = 30;
pub const TAG_LEN_2: u8 = 31;

// Tag key prefixes (key = byte & 0xF8)
pub const TAG_IS_NULL: u8 = 0 << 3;
pub const TAG_EXAMPLE: u8 = 1 << 3;
pub const TAG_DEPRECATED: u8 = 2 << 3;
pub const TAG_DESC: u8 = 3 << 3;
pub const TAG_TYPE: u8 = 4 << 3;
pub const TAG_NULLABLE: u8 = 5 << 3;
pub const TAG_ALLOW_EMPTY: u8 = 6 << 3;
pub const TAG_UNIQUE: u8 = 7 << 3;
pub const TAG_DEFAULT_VAL: u8 = 8 << 3;
pub const TAG_MIN: u8 = 9 << 3;
pub const TAG_MAX: u8 = 10 << 3;
pub const TAG_SIZE: u8 = 11 << 3;
pub const TAG_ENUMS: u8 = 12 << 3;
pub const TAG_PATTERN: u8 = 13 << 3;
pub const TAG_LOCATION: u8 = 14 << 3;
pub const TAG_VERSION: u8 = 15 << 3;
pub const TAG_MIME: u8 = 16 << 3;

pub const TAG_CHILD_DESC: u8 = 17 << 3;
pub const TAG_CHILD_TYPE: u8 = 18 << 3;
pub const TAG_CHILD_NULLABLE: u8 = 19 << 3;
pub const TAG_CHILD_ALLOW_EMPTY: u8 = 20 << 3;
pub const TAG_CHILD_UNIQUE: u8 = 21 << 3;
pub const TAG_CHILD_DEFAULT_VAL: u8 = 22 << 3;
pub const TAG_CHILD_MIN: u8 = 23 << 3;
pub const TAG_CHILD_MAX: u8 = 24 << 3;
pub const TAG_CHILD_SIZE: u8 = 25 << 3;
pub const TAG_CHILD_ENUMS: u8 = 26 << 3;
pub const TAG_CHILD_PATTERN: u8 = 27 << 3;
pub const TAG_CHILD_LOCATION: u8 = 28 << 3;
pub const TAG_CHILD_VERSION: u8 = 29 << 3;
pub const TAG_CHILD_MIME: u8 = 30 << 3;
pub const TAG_MORE: u8 = 31 << 3;

pub const TAG_KEY_MASK: u8 = 0xF8;
pub const TAG_PAYLOAD_MASK: u8 = 0x07;