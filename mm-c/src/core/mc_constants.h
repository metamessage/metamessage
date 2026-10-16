#ifndef MMC_CORE_CONSTANTS_H
#define MMC_CORE_CONSTANTS_H

#include <stdint.h>
#include <stdbool.h>

enum {
    MM_PREFIX_SIMPLE      = 0b000 << 5,
    MM_PREFIX_POSITIVEINT = 0b001 << 5,
    MM_PREFIX_NEGATIVEINT = 0b010 << 5,
    MM_PREFIX_FLOAT       = 0b011 << 5,
    MM_PREFIX_STRING      = 0b100 << 5,
    MM_PREFIX_BYTES       = 0b101 << 5,
    MM_PREFIX_CONTAINER   = 0b110 << 5,
    MM_PREFIX_TAG         = 0b111 << 5
};

enum {
    MM_SIMPLE_NULL       = 0,
    MM_SIMPLE_NULLBOOL   = 1,
    MM_SIMPLE_NULLINT    = 2,
    MM_SIMPLE_NULLFLOAT  = 3,
    MM_SIMPLE_NULLSTRING = 4,
    MM_SIMPLE_NULLBYTES  = 5,
    MM_SIMPLE_FALSE = 6,
    MM_SIMPLE_TRUE  = 7
};

enum {
    MM_CONTAINER_OBJECT = 0x00,
    MM_CONTAINER_ARRAY  = 0x10
};

enum {
    MM_INTLEN1BYTE = 24,
    MM_INTLEN2BYTE = 25,
    MM_INTLEN3BYTE = 26,
    MM_INTLEN4BYTE = 27,
    MM_INTLEN5BYTE = 28,
    MM_INTLEN6BYTE = 29,
    MM_INTLEN7BYTE = 30,
    MM_INTLEN8BYTE = 31
};

enum {
    MM_FLOAT_NEG_MASK   = 0x10,
    MM_FLOATLEN1BYTE = 8,
    MM_FLOATLEN2BYTE = 9,
    MM_FLOATLEN3BYTE = 10,
    MM_FLOATLEN4BYTE = 11,
    MM_FLOATLEN5BYTE = 12,
    MM_FLOATLEN6BYTE = 13,
    MM_FLOATLEN7BYTE = 14,
    MM_FLOATLEN8BYTE = 15
};

/* 长度上限（v0.3.0）。用宏而非 enum：0xFFFFFFFF 超出 int 表示范围。 */
#define MM_MAX1BYTE 0xFFu
#define MM_MAX2BYTE 0xFFFFu
#define MM_MAX3BYTE 0xFFFFFFu
#define MM_MAX4BYTE 0xFFFFFFFFu

/*
 * v0.3.0 长度档位分配（仅 Bytes / Container / Tag，String 不在本次范围内）：
 *
 *   Bytes / Tag（5-bit suffix）：
 *     0..27      = 内联长度（length 直接等于 suffix）
 *     28/29/30/31 = 1/2/3/4 字节大端长度档位标记
 *   Container（suffix 第 4 位是 Array/Object 标志，低 4 位是长度）：
 *     0..11      = 内联长度
 *     12/13/14/15 = 1/2/3/4 字节大端长度档位标记
 *
 * 最大可表达长度 = MM_MAX4BYTE(0xFFFFFFFF) ≈ 4.29GB。
 *
 * 这是不兼容变更：v0.2.x 用 30/31 作 1/2 字节标记、内联范围 0..29（Bytes/Tag）
 * 与 0..13（Container），现在内联收窄以腾出四个档位标记。
 *
 * String 不参与本次扩展：维持 v0.2.x 语义（内联 0..29，30/31 = 1/2 字节，
 * 上限 MM_MAX2BYTE = 65535），字节级与 v0.2.x 完全一致。
 */

/* String：v0.2.x 语义，未改动 */
enum {
    MM_STRINGLEN1BYTE = 30,
    MM_STRINGLEN2BYTE = 31
};

enum {
    MM_BYTESLEN1BYTE = 28,
    MM_BYTESLEN2BYTE = 29,
    MM_BYTESLEN3BYTE = 30,
    MM_BYTESLEN4BYTE = 31
};

enum {
    MM_CONTAINERLEN1BYTE = 12,
    MM_CONTAINERLEN2BYTE = 13,
    MM_CONTAINERLEN3BYTE = 14,
    MM_CONTAINERLEN4BYTE = 15
};

enum {
    MM_TAGLEN1BYTE = 28,
    MM_TAGLEN2BYTE = 29,
    MM_TAGLEN3BYTE = 30,
    MM_TAGLEN4BYTE = 31
};

/* TagPayload（encodeT 内部方案）：v0.2.x 语义，未改动 */
enum {
    MM_TAGPAYLOAD1BYTE = 30,
    MM_TAGPAYLOAD2BYTE = 31
};

static inline int mm_prefix_of(uint8_t b) { return b & 0b11100000; }
static inline int mm_suffix_of(uint8_t b) { return b & 0b00011111; }

static inline int mm_int_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_INTLEN1BYTE) return 0;
    return l - MM_INTLEN1BYTE + 1;
}

static inline int mm_float_len(uint8_t b) {
    int l = b & 0b01111;
    if (l < MM_FLOATLEN1BYTE) return 0;
    return l - MM_FLOATLEN1BYTE + 1;
}

static inline int mm_string_extra_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_STRINGLEN1BYTE) return 0;
    if (l == MM_STRINGLEN1BYTE) return 1;
    return 2;
}

static inline int mm_string_inline_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_STRINGLEN1BYTE) return l;
    return 0;
}

static inline int mm_bytes_extra_len(uint8_t b) {
    int l = b & 0b11111;
    switch (l) {
    case MM_BYTESLEN1BYTE: return 1;
    case MM_BYTESLEN2BYTE: return 2;
    case MM_BYTESLEN3BYTE: return 3;
    case MM_BYTESLEN4BYTE: return 4;
    default: return 0;
    }
}

static inline int mm_bytes_inline_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_BYTESLEN1BYTE) return l;
    return 0;
}

static inline int mm_container_extra_len(uint8_t b) {
    int l = b & 0b01111;
    switch (l) {
    case MM_CONTAINERLEN1BYTE: return 1;
    case MM_CONTAINERLEN2BYTE: return 2;
    case MM_CONTAINERLEN3BYTE: return 3;
    case MM_CONTAINERLEN4BYTE: return 4;
    default: return 0;
    }
}

static inline int mm_container_inline_len(uint8_t b) {
    int l = b & 0b01111;
    if (l < MM_CONTAINERLEN1BYTE) return l;
    return 0;
}

static inline bool mm_is_array_container(uint8_t b) {
    return (b & 0b10000) == MM_CONTAINER_ARRAY;
}

static inline int mm_tag_extra_len(uint8_t b) {
    int l = b & 0b11111;
    switch (l) {
    case MM_TAGLEN1BYTE: return 1;
    case MM_TAGLEN2BYTE: return 2;
    case MM_TAGLEN3BYTE: return 3;
    case MM_TAGLEN4BYTE: return 4;
    default: return 0;
    }
}

static inline int mm_tag_inline_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_TAGLEN1BYTE) return l;
    return 0;
}

#endif