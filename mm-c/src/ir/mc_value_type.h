#ifndef MMC_IR_VALUE_TYPE_H
#define MMC_IR_VALUE_TYPE_H

#include <stdint.h>
#include <string.h>

typedef enum {
    MM_VALUE_UNKNOWN = 0,
    MM_VALUE_DOC = 1,
    MM_VALUE_VEC = 2,
    MM_VALUE_ARR = 3,
    MM_VALUE_OBJ = 4,
    MM_VALUE_MAP = 5,
    MM_VALUE_STR = 6,
    MM_VALUE_BYTES = 7,
    MM_VALUE_BOOL = 8,
    MM_VALUE_I = 9,
    MM_VALUE_I8 = 10,
    MM_VALUE_I16 = 11,
    MM_VALUE_I32 = 12,
    MM_VALUE_I64 = 13,
    MM_VALUE_U = 14,
    MM_VALUE_U8 = 15,
    MM_VALUE_U16 = 16,
    MM_VALUE_U32 = 17,
    MM_VALUE_U64 = 18,
    MM_VALUE_F32 = 19,
    MM_VALUE_F64 = 20,
    MM_VALUE_BIGINT = 21,
    MM_VALUE_DATETIME = 22,
    MM_VALUE_DATE = 23,
    MM_VALUE_TIME = 24,
    MM_VALUE_UUID = 25,
    MM_VALUE_DECIMAL = 26,
    MM_VALUE_IP = 27,
    MM_VALUE_URL = 28,
    MM_VALUE_EMAIL = 29,
    MM_VALUE_ENUM = 30,
    MM_VALUE_IMAGE = 31,
    MM_VALUE_VIDEO = 32
} mm_value_type_t;

static inline const char* mm_value_type_to_string(mm_value_type_t vt) {
    switch (vt) {
        case MM_VALUE_UNKNOWN: return "unknown";
        case MM_VALUE_DOC: return "doc";
        case MM_VALUE_VEC: return "vec";
        case MM_VALUE_ARR: return "arr";
        case MM_VALUE_OBJ: return "obj";
        case MM_VALUE_MAP: return "map";
        case MM_VALUE_STR: return "str";
        case MM_VALUE_BYTES: return "bytes";
        case MM_VALUE_BOOL: return "bool";
        case MM_VALUE_I: return "i";
        case MM_VALUE_I8: return "i8";
        case MM_VALUE_I16: return "i16";
        case MM_VALUE_I32: return "i32";
        case MM_VALUE_I64: return "i64";
        case MM_VALUE_U: return "u";
        case MM_VALUE_U8: return "u8";
        case MM_VALUE_U16: return "u16";
        case MM_VALUE_U32: return "u32";
        case MM_VALUE_U64: return "u64";
        case MM_VALUE_F32: return "f32";
        case MM_VALUE_F64: return "f64";
        case MM_VALUE_BIGINT: return "bigint";
        case MM_VALUE_DATETIME: return "datetime";
        case MM_VALUE_DATE: return "date";
        case MM_VALUE_TIME: return "time";
        case MM_VALUE_UUID: return "uuid";
        case MM_VALUE_DECIMAL: return "decimal";
        case MM_VALUE_IP: return "ip";
        case MM_VALUE_URL: return "url";
        case MM_VALUE_EMAIL: return "email";
        case MM_VALUE_ENUM: return "enum";
        case MM_VALUE_IMAGE: return "image";
        case MM_VALUE_VIDEO: return "video";
        default: return "unknown";
    }
}

static inline mm_value_type_t mm_value_type_parse(const char* s) {
    if (!s) return MM_VALUE_UNKNOWN;
    size_t len = strlen(s);
    char lower[32];
    if (len >= sizeof(lower)) return MM_VALUE_UNKNOWN;
    for (size_t i = 0; i < len; ++i) {
        char c = s[i];
        if (c >= 'A' && c <= 'Z') c += 32;
        lower[i] = c;
    }
    lower[len] = '\0';

    if (strcmp(lower, "doc") == 0) return MM_VALUE_DOC;
    if (strcmp(lower, "vec") == 0) return MM_VALUE_VEC;
    if (strcmp(lower, "arr") == 0) return MM_VALUE_ARR;
    if (strcmp(lower, "obj") == 0) return MM_VALUE_OBJ;
    if (strcmp(lower, "map") == 0) return MM_VALUE_MAP;
    if (strcmp(lower, "str") == 0) return MM_VALUE_STR;
    if (strcmp(lower, "string") == 0) return MM_VALUE_STR;
    if (strcmp(lower, "bytes") == 0) return MM_VALUE_BYTES;
    if (strcmp(lower, "bool") == 0) return MM_VALUE_BOOL;
    if (strcmp(lower, "boolean") == 0) return MM_VALUE_BOOL;
    if (strcmp(lower, "i") == 0) return MM_VALUE_I;
    if (strcmp(lower, "int") == 0) return MM_VALUE_I;
    if (strcmp(lower, "i8") == 0) return MM_VALUE_I8;
    if (strcmp(lower, "int8") == 0) return MM_VALUE_I8;
    if (strcmp(lower, "i16") == 0) return MM_VALUE_I16;
    if (strcmp(lower, "int16") == 0) return MM_VALUE_I16;
    if (strcmp(lower, "i32") == 0) return MM_VALUE_I32;
    if (strcmp(lower, "int32") == 0) return MM_VALUE_I32;
    if (strcmp(lower, "i64") == 0) return MM_VALUE_I64;
    if (strcmp(lower, "int64") == 0) return MM_VALUE_I64;
    if (strcmp(lower, "u") == 0) return MM_VALUE_U;
    if (strcmp(lower, "uint") == 0) return MM_VALUE_U;
    if (strcmp(lower, "u8") == 0) return MM_VALUE_U8;
    if (strcmp(lower, "uint8") == 0) return MM_VALUE_U8;
    if (strcmp(lower, "u16") == 0) return MM_VALUE_U16;
    if (strcmp(lower, "uint16") == 0) return MM_VALUE_U16;
    if (strcmp(lower, "u32") == 0) return MM_VALUE_U32;
    if (strcmp(lower, "uint32") == 0) return MM_VALUE_U32;
    if (strcmp(lower, "u64") == 0) return MM_VALUE_U64;
    if (strcmp(lower, "uint64") == 0) return MM_VALUE_U64;
    if (strcmp(lower, "f32") == 0) return MM_VALUE_F32;
    if (strcmp(lower, "float32") == 0) return MM_VALUE_F32;
    if (strcmp(lower, "f64") == 0) return MM_VALUE_F64;
    if (strcmp(lower, "float64") == 0) return MM_VALUE_F64;
    if (strcmp(lower, "float") == 0) return MM_VALUE_F64;
    if (strcmp(lower, "bigint") == 0) return MM_VALUE_BIGINT;
    if (strcmp(lower, "datetime") == 0) return MM_VALUE_DATETIME;
    if (strcmp(lower, "date") == 0) return MM_VALUE_DATE;
    if (strcmp(lower, "time") == 0) return MM_VALUE_TIME;
    if (strcmp(lower, "uuid") == 0) return MM_VALUE_UUID;
    if (strcmp(lower, "decimal") == 0) return MM_VALUE_DECIMAL;
    if (strcmp(lower, "ip") == 0) return MM_VALUE_IP;
    if (strcmp(lower, "url") == 0) return MM_VALUE_URL;
    if (strcmp(lower, "email") == 0) return MM_VALUE_EMAIL;
    if (strcmp(lower, "enum") == 0) return MM_VALUE_ENUM;
    if (strcmp(lower, "image") == 0) return MM_VALUE_IMAGE;
    if (strcmp(lower, "video") == 0) return MM_VALUE_VIDEO;
    return MM_VALUE_UNKNOWN;
}

#endif