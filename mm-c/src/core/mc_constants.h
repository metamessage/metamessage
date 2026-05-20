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
    MM_SIMPLE_NULLBOOL   = 0,
    MM_SIMPLE_NULLINT    = 1,
    MM_SIMPLE_NULLFLOAT  = 2,
    MM_SIMPLE_NULLSTRING = 3,
    MM_SIMPLE_NULLBYTES  = 4,
    MM_SIMPLE_FALSE = 5,
    MM_SIMPLE_TRUE  = 6
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
    MM_FLOATLEN1BYTE = 8,
    MM_FLOATLEN2BYTE = 9,
    MM_FLOATLEN3BYTE = 10,
    MM_FLOATLEN4BYTE = 11,
    MM_FLOATLEN5BYTE = 12,
    MM_FLOATLEN6BYTE = 13,
    MM_FLOATLEN7BYTE = 14,
    MM_FLOATLEN8BYTE = 15
};

enum {
    MM_STRINGLEN1BYTE = 30,
    MM_STRINGLEN2BYTE = 31
};

enum {
    MM_BYTESLEN1BYTE = 30,
    MM_BYTESLEN2BYTE = 31
};

enum {
    MM_CONTAINERLEN1BYTE = 14,
    MM_CONTAINERLEN2BYTE = 15
};

enum {
    MM_TAGLEN1BYTE = 30,
    MM_TAGLEN2BYTE = 31
};

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
    if (l < MM_BYTESLEN1BYTE) return 0;
    if (l == MM_BYTESLEN1BYTE) return 1;
    return 2;
}

static inline int mm_bytes_inline_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_BYTESLEN1BYTE) return l;
    return 0;
}

static inline int mm_container_extra_len(uint8_t b) {
    int l = b & 0b01111;
    if (l < MM_CONTAINERLEN1BYTE) return 0;
    if (l == MM_CONTAINERLEN1BYTE) return 1;
    return 2;
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
    if (l < MM_TAGLEN1BYTE) return 0;
    if (l == MM_TAGLEN1BYTE) return 1;
    return 2;
}

static inline int mm_tag_inline_len(uint8_t b) {
    int l = b & 0b11111;
    if (l < MM_TAGLEN1BYTE) return l;
    return 0;
}

#endif