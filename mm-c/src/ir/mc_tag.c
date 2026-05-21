#include "mc_tag.h"
#include "mc_value_type.h"
#include <ctype.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static size_t quoted_string_len(const char* s)
{
    size_t n = 2;
    for (; *s; s++) {
        unsigned char c = (unsigned char)*s;
        if (c == '"' || c == '\\') {
            n += 2;
        } else if (c == '\n') {
            n += 2;
        } else if (c == '\r') {
            n += 2;
        } else if (c == '\t') {
            n += 2;
        } else if (c < 0x20) {
            n += 4;
        } else {
            n++;
        }
    }
    return n;
}

static char* quote_string(const char* s)
{
    if (!s) return NULL;
    size_t n = quoted_string_len(s);
    char* out = (char*)malloc(n + 1);
    if (!out) return NULL;
    char* p = out;
    *p++ = '"';
    for (; *s; s++) {
        unsigned char c = (unsigned char)*s;
        if (c == '"' || c == '\\') {
            *p++ = '\\';
            *p++ = (char)c;
        } else if (c == '\n') {
            *p++ = '\\';
            *p++ = 'n';
        } else if (c == '\r') {
            *p++ = '\\';
            *p++ = 'r';
        } else if (c == '\t') {
            *p++ = '\\';
            *p++ = 't';
        } else if (c < 0x20) {
            p += sprintf(p, "\\x%02x", c);
        } else {
            *p++ = (char)c;
        }
    }
    *p++ = '"';
    *p = '\0';
    return out;
}

static char* unquote_string(const char* s, size_t len)
{
    if (len < 2 || s[0] != '"' || s[len - 1] != '"') return NULL;
    size_t inner_len = len - 2;
    const char* inner = s + 1;
    char* out = (char*)malloc(inner_len + 1);
    if (!out) return NULL;
    size_t o = 0;
    for (size_t i = 0; i < inner_len; i++) {
        if (inner[i] == '\\' && i + 1 < inner_len) {
            i++;
            switch (inner[i]) {
            case '"':  out[o++] = '"'; break;
            case '\\': out[o++] = '\\'; break;
            case '/':  out[o++] = '/'; break;
            case 'n':  out[o++] = '\n'; break;
            case 'r':  out[o++] = '\r'; break;
            case 't':  out[o++] = '\t'; break;
            default:   out[o++] = inner[i]; break;
            }
        } else {
            out[o++] = inner[i];
        }
    }
    out[o] = '\0';
    return out;
}

static void append_str(char** buf, size_t* cap, size_t* len, const char* s)
{
    size_t slen = strlen(s);
    if (*len + slen + 1 > *cap) {
        *cap = *cap ? *cap * 2 : 64;
        while (*len + slen + 1 > *cap) *cap *= 2;
        *buf = (char*)realloc(*buf, *cap);
    }
    memcpy(*buf + *len, s, slen);
    *len += slen;
    (*buf)[*len] = '\0';
}

static void append_fmt(char** buf, size_t* cap, size_t* len, const char* fmt, ...)
{
    va_list args;
    va_start(args, fmt);
    int needed = vsnprintf(NULL, 0, fmt, args);
    va_end(args);
    if (needed < 0) return;
    size_t needed_s = (size_t)needed + 1;
    if (*len + needed_s > *cap) {
        *cap = *cap ? *cap * 2 : 64;
        while (*len + needed_s > *cap) *cap *= 2;
        *buf = (char*)realloc(*buf, *cap);
    }
    va_start(args, fmt);
    vsnprintf(*buf + *len, needed_s, fmt, args);
    va_end(args);
    *len += (size_t)needed;
}

static bool is_simple_type(mm_value_type_t vt)
{
    return vt == MM_VALUE_STR || vt == MM_VALUE_I || vt == MM_VALUE_F64 ||
           vt == MM_VALUE_BOOL || vt == MM_VALUE_OBJ || vt == MM_VALUE_VEC;
}

static bool is_bytes_simple_type(mm_value_type_t vt)
{
    return vt == MM_VALUE_STR || vt == MM_VALUE_BYTES || vt == MM_VALUE_I ||
           vt == MM_VALUE_F64 || vt == MM_VALUE_BOOL || vt == MM_VALUE_OBJ ||
           vt == MM_VALUE_VEC;
}

static void encode_string(uint8_t** buf, size_t* cap, size_t* len, int key, const char* s)
{
    size_t slen = strlen(s);
    size_t needed;
    if (slen <= 5) {
        needed = 1 + slen;
    } else if (slen <= 255) {
        needed = 2 + slen;
    } else if (slen <= 65535) {
        needed = 3 + slen;
    } else {
        return;
    }
    if (*len + needed > *cap) {
        *cap = *cap ? *cap * 2 : 128;
        while (*len + needed > *cap) *cap *= 2;
        *buf = (uint8_t*)realloc(*buf, *cap);
    }
    if (slen <= 5) {
        (*buf)[(*len)++] = (uint8_t)(key | (int)slen);
        memcpy(*buf + *len, s, slen);
        *len += slen;
    } else if (slen <= 255) {
        (*buf)[(*len)++] = (uint8_t)(key | 6);
        (*buf)[(*len)++] = (uint8_t)(slen & 0xFF);
        memcpy(*buf + *len, s, slen);
        *len += slen;
    } else {
        (*buf)[(*len)++] = (uint8_t)(key | 7);
        (*buf)[(*len)++] = (uint8_t)((slen >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(slen & 0xFF);
        memcpy(*buf + *len, s, slen);
        *len += slen;
    }
}

static void encode_string_simple(uint8_t** buf, size_t* cap, size_t* len, int key, const char* s)
{
    size_t slen = strlen(s);
    size_t needed;
    if (slen < 7) {
        needed = 1 + slen;
    } else {
        needed = 2 + slen;
    }
    if (*len + needed > *cap) {
        *cap = *cap ? *cap * 2 : 128;
        while (*len + needed > *cap) *cap *= 2;
        *buf = (uint8_t*)realloc(*buf, *cap);
    }
    if (slen < 7) {
        (*buf)[(*len)++] = (uint8_t)(key | (int)slen);
    } else {
        (*buf)[(*len)++] = (uint8_t)(key | 7);
        (*buf)[(*len)++] = (uint8_t)(slen & 0xFF);
    }
    memcpy(*buf + *len, s, slen);
    *len += slen;
}

static void encode_u64(uint8_t** buf, size_t* cap, size_t* len, int key, uint64_t uv)
{
    size_t needed;
    if (uv <= 0xFFULL) {
        needed = 2;
    } else if (uv <= 0xFFFFULL) {
        needed = 3;
    } else if (uv <= 0xFFFFFFULL) {
        needed = 4;
    } else if (uv <= 0xFFFFFFFFULL) {
        needed = 5;
    } else if (uv <= 0xFFFFFFFFFFULL) {
        needed = 6;
    } else if (uv <= 0xFFFFFFFFFFFFULL) {
        needed = 7;
    } else if (uv <= 0xFFFFFFFFFFFFFFULL) {
        needed = 8;
    } else {
        needed = 9;
    }
    if (*len + needed > *cap) {
        *cap = *cap ? *cap * 2 : 128;
        while (*len + needed > *cap) *cap *= 2;
        *buf = (uint8_t*)realloc(*buf, *cap);
    }
    if (uv <= 0xFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 0);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 1);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 2);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 3);
        (*buf)[(*len)++] = (uint8_t)((uv >> 24) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFFFFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 4);
        (*buf)[(*len)++] = (uint8_t)((uv >> 32) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 24) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFFFFFFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 5);
        (*buf)[(*len)++] = (uint8_t)((uv >> 40) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 32) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 24) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else if (uv <= 0xFFFFFFFFFFFFFFULL) {
        (*buf)[(*len)++] = (uint8_t)(key | 6);
        (*buf)[(*len)++] = (uint8_t)((uv >> 48) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 40) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 32) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 24) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    } else {
        (*buf)[(*len)++] = (uint8_t)(key | 7);
        (*buf)[(*len)++] = (uint8_t)((uv >> 56) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 48) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 40) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 32) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 24) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 16) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)((uv >> 8) & 0xFF);
        (*buf)[(*len)++] = (uint8_t)(uv & 0xFF);
    }
}

static uint8_t parse_mime(const char* s)
{
    if (!s) return 0;
    size_t slen = strlen(s);
    if (slen >= 64) return 0;
    char lower[64];
    for (size_t i = 0; i < slen; i++) {
        char c = s[i];
        if (c >= 'A' && c <= 'Z') c += 32;
        lower[i] = c;
    }
    lower[slen] = '\0';

    if (strcmp(lower, "image/jpeg") == 0) return 1;
    if (strcmp(lower, "image/jpg") == 0) return 1;
    if (strcmp(lower, "image/png") == 0) return 2;
    if (strcmp(lower, "image/gif") == 0) return 3;
    if (strcmp(lower, "image/webp") == 0) return 4;
    if (strcmp(lower, "image/svg+xml") == 0) return 5;
    if (strcmp(lower, "image/avif") == 0) return 6;
    if (strcmp(lower, "image/bmp") == 0) return 7;
    if (strcmp(lower, "image/x-icon") == 0) return 8;
    if (strcmp(lower, "image/tiff") == 0) return 9;
    if (strcmp(lower, "image/heic") == 0) return 10;
    if (strcmp(lower, "image/heif") == 0) return 11;
    if (strcmp(lower, "text/plain") == 0) return 12;
    if (strcmp(lower, "text/html") == 0) return 13;
    if (strcmp(lower, "text/css") == 0) return 14;
    if (strcmp(lower, "text/javascript") == 0) return 15;
    if (strcmp(lower, "application/javascript") == 0) return 15;
    if (strcmp(lower, "application/json") == 0) return 16;
    if (strcmp(lower, "text/csv") == 0) return 17;
    if (strcmp(lower, "text/markdown") == 0) return 18;
    if (strcmp(lower, "application/pdf") == 0) return 19;
    if (strcmp(lower, "application/zip") == 0) return 20;
    if (strcmp(lower, "application/gzip") == 0) return 21;
    if (strcmp(lower, "application/x-tar") == 0) return 22;
    if (strcmp(lower, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") == 0) return 23;
    if (strcmp(lower, "application/vnd.openxmlformats-officedocument.wordprocessingml.document") == 0) return 24;
    if (strcmp(lower, "application/octet-stream") == 0) return 25;
    if (strcmp(lower, "video/mp4") == 0) return 26;
    if (strcmp(lower, "video/webm") == 0) return 27;
    if (strcmp(lower, "video/mov") == 0) return 28;
    if (strcmp(lower, "audio/mpeg") == 0) return 29;
    if (strcmp(lower, "audio/wav") == 0) return 30;
    if (strcmp(lower, "audio/flac") == 0) return 31;
    if (strcmp(lower, "font/woff2") == 0) return 32;
    if (strcmp(lower, "font/ttf") == 0) return 33;

    return 0;
}

void mm_tag_init(mm_tag_t* tag)
{
    if (!tag) return;
    memset(tag, 0, sizeof(*tag));
    tag->version = MM_TAG_DEFAULT_VERSION;
    tag->child_version = MM_TAG_DEFAULT_VERSION;
}

void mm_tag_cleanup(mm_tag_t* tag)
{
    if (!tag) return;
    free(tag->name);
    free(tag->desc);
    free(tag->default_val);
    free(tag->min);
    free(tag->max);
    free(tag->enums);
    free(tag->pattern);
    free(tag->mime);
    free(tag->child_desc);
    free(tag->child_default);
    free(tag->child_min);
    free(tag->child_max);
    free(tag->child_enum);
    free(tag->child_pattern);
    free(tag->child_mime);
    memset(tag, 0, sizeof(*tag));
}

void mm_tag_inherit(mm_tag_t* tag, const mm_tag_t* parent)
{
    if (!tag || !parent) return;
    tag->is_inherit = true;

    if (parent->child_desc) {
        free(tag->desc);
        tag->desc = strdup(parent->child_desc);
    }

    if (parent->child_type != MM_VALUE_UNKNOWN) {
        tag->type = parent->child_type;
    }

    if (parent->child_raw) {
        tag->raw = parent->child_raw;
    }

    if (parent->child_nullable) {
        tag->nullable = parent->child_nullable;
    }

    if (parent->child_allow_empty) {
        tag->allow_empty = parent->child_allow_empty;
    }

    if (parent->child_unique) {
        tag->unique = parent->child_unique;
    }

    if (parent->child_default) {
        free(tag->default_val);
        tag->default_val = strdup(parent->child_default);
    }

    if (parent->child_min) {
        free(tag->min);
        tag->min = strdup(parent->child_min);
    }

    if (parent->child_max) {
        free(tag->max);
        tag->max = strdup(parent->child_max);
    }

    if (parent->child_size != 0) {
        tag->size = parent->child_size;
    }

    if (parent->child_enum) {
        free(tag->enums);
        tag->enums = strdup(parent->child_enum);
    }

    if (parent->child_pattern) {
        free(tag->pattern);
        tag->pattern = strdup(parent->child_pattern);
    }

    if (parent->child_location_offset != 0) {
        tag->location_offset = parent->child_location_offset;
    }

    if (parent->child_version != MM_TAG_DEFAULT_VERSION) {
        tag->version = parent->child_version;
    }

    if (parent->child_mime) {
        free(tag->mime);
        tag->mime = strdup(parent->child_mime);
    }
}

char* mm_tag_to_string(const mm_tag_t* tag)
{
    if (!tag) return strdup("");

    char* buf = NULL;
    size_t cap = 0;
    size_t len = 0;
    bool first = true;

    if (tag->type != MM_VALUE_UNKNOWN && !tag->is_inherit) {
        if (is_simple_type(tag->type)) {
        } else {
            if ((tag->type == MM_VALUE_ARR && tag->size > 0) ||
                (tag->type == MM_VALUE_ENUM && tag->enums)) {
            } else {
                if (!first) append_str(&buf, &cap, &len, "; ");
                append_str(&buf, &cap, &len, "type=");
                append_str(&buf, &cap, &len, mm_value_type_to_string(tag->type));
                first = false;
            }
        }
    }

    if (tag->example) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "example");
        first = false;
    }

    if (tag->is_null) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "is_null");
        first = false;
    }

    if (tag->nullable && !tag->is_inherit) {
        if (!tag->is_null) {
            if (!first) append_str(&buf, &cap, &len, "; ");
            append_str(&buf, &cap, &len, "nullable");
            first = false;
        }
    }

    if (tag->desc && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "desc=");
        char* q = quote_string(tag->desc);
        if (q) {
            append_str(&buf, &cap, &len, q);
            free(q);
        }
        first = false;
    }

    if (tag->raw && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "raw");
        first = false;
    }

    if (tag->allow_empty && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "allow_empty");
        first = false;
    }

    if (tag->unique && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "unique");
        first = false;
    }

    if (tag->default_val && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "default=");
        append_str(&buf, &cap, &len, tag->default_val);
        first = false;
    }

    if (tag->min && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "min=");
        append_str(&buf, &cap, &len, tag->min);
        first = false;
    }

    if (tag->max && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "max=");
        append_str(&buf, &cap, &len, tag->max);
        first = false;
    }

    if (tag->size != 0 && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "size=%d", tag->size);
        first = false;
    }

    if (tag->enums && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "enum=");
        append_str(&buf, &cap, &len, tag->enums);
        first = false;
    }

    if (tag->pattern && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "pattern=");
        append_str(&buf, &cap, &len, tag->pattern);
        first = false;
    }

    if (tag->location_offset != 0 && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "location=%d", tag->location_offset);
        first = false;
    }

    if (tag->version != MM_TAG_DEFAULT_VERSION && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "version=%d", tag->version);
        first = false;
    }

    if (tag->mime && !tag->is_inherit) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "mime=");
        append_str(&buf, &cap, &len, tag->mime);
        first = false;
    }

    if (tag->child_desc) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_desc=");
        char* q = quote_string(tag->child_desc);
        if (q) {
            append_str(&buf, &cap, &len, q);
            free(q);
        }
        first = false;
    }

    if (tag->child_type != MM_VALUE_UNKNOWN) {
        if (is_simple_type(tag->child_type)) {
        } else {
            if ((tag->child_type == MM_VALUE_ARR && tag->child_size > 0) ||
                (tag->child_type == MM_VALUE_ENUM && tag->child_enum)) {
            } else {
                if (!first) append_str(&buf, &cap, &len, "; ");
                append_str(&buf, &cap, &len, "child_type=");
                append_str(&buf, &cap, &len, mm_value_type_to_string(tag->child_type));
                first = false;
            }
        }
    }

    if (tag->child_raw) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_raw");
        first = false;
    }

    if (tag->child_nullable) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_nullable");
        first = false;
    }

    if (tag->child_allow_empty) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_allow_empty");
        first = false;
    }

    if (tag->child_unique) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_unique");
        first = false;
    }

    if (tag->child_default) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_default=");
        append_str(&buf, &cap, &len, tag->child_default);
        first = false;
    }

    if (tag->child_min) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_min=");
        append_str(&buf, &cap, &len, tag->child_min);
        first = false;
    }

    if (tag->child_max) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_max=");
        append_str(&buf, &cap, &len, tag->child_max);
        first = false;
    }

    if (tag->child_size != 0) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "child_size=%d", tag->child_size);
        first = false;
    }

    if (tag->child_enum) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_enum=");
        append_str(&buf, &cap, &len, tag->child_enum);
        first = false;
    }

    if (tag->child_pattern) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_pattern=");
        append_str(&buf, &cap, &len, tag->child_pattern);
        first = false;
    }

    if (tag->child_location_offset != 0) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "child_location=%d", tag->child_location_offset);
        first = false;
    }

    if (tag->child_version != MM_TAG_DEFAULT_VERSION) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_fmt(&buf, &cap, &len, "child_version=%d", tag->child_version);
        first = false;
    }

    if (tag->child_mime) {
        if (!first) append_str(&buf, &cap, &len, "; ");
        append_str(&buf, &cap, &len, "child_mime=");
        append_str(&buf, &cap, &len, tag->child_mime);
        first = false;
    }

    if (!buf) return strdup("");
    return buf;
}

// Split by ';' respecting double-quoted strings
static char** split_tag_respecting_quotes(const char* tag_str, int* count)
{
    *count = 0;
    if (!tag_str || !*tag_str) return NULL;

    size_t len = strlen(tag_str);
    int capacity = 16;
    char** parts = (char**)malloc(sizeof(char*) * (size_t)capacity);
    if (!parts) return NULL;

    bool in_quote = false;
    size_t start = 0;
    for (size_t i = 0; i <= len; i++) {
        if (i < len && tag_str[i] == '"') {
            in_quote = !in_quote;
        }
        if (!in_quote && (i == len || tag_str[i] == ';')) {
            size_t part_len = i - start;
            while (part_len > 0 && tag_str[start] == ' ') { start++; part_len--; }
            while (part_len > 0 && tag_str[start + part_len - 1] == ' ') part_len--;
            if (part_len > 0) {
                if (*count >= capacity) {
                    capacity *= 2;
                    parts = (char**)realloc(parts, sizeof(char*) * (size_t)capacity);
                }
                parts[*count] = (char*)malloc(part_len + 1);
                memcpy(parts[*count], tag_str + start, part_len);
                parts[*count][part_len] = '\0';
                (*count)++;
            }
            start = i + 1;
        }
    }

    if (*count == 0) {
        free(parts);
        return NULL;
    }
    return parts;
}

static void free_split_parts(char** parts, int count)
{
    if (!parts) return;
    for (int i = 0; i < count; i++) {
        free(parts[i]);
    }
    free(parts);
}

mm_tag_t mm_tag_parse(const char* tag_str)
{
    mm_tag_t r;
    mm_tag_init(&r);

    if (!tag_str) return r;

    const char* p = tag_str;
    while (*p == ' ' || *p == '\t' || *p == '\n' || *p == '\r') p++;

    if (p[0] == '/' && p[1] == '/') p += 2;

    while (*p == ' ' || *p == '\t' || *p == '\n' || *p == '\r') p++;

    if (p[0] == 'm' && p[1] == 'm' && p[2] == ':') p += 3;

    while (*p == ' ' || *p == '\t' || *p == '\n' || *p == '\r') p++;

    if (!*p) return r;

    int part_count = 0;
    char** parts = split_tag_respecting_quotes(p, &part_count);

    for (int i = 0; i < part_count; i++) {
        char* eq = strchr(parts[i], '=');
        char* k = NULL;
        char* v = NULL;
        if (eq) {
            size_t k_len = (size_t)(eq - parts[i]);
            while (k_len > 0 && parts[i][k_len - 1] == ' ') k_len--;
            k = (char*)malloc(k_len + 1);
            memcpy(k, parts[i], k_len);
            k[k_len] = '\0';

            const char* vstart = eq + 1;
            while (*vstart == ' ') vstart++;
            v = strdup(vstart);
        } else {
            k = strdup(parts[i]);
            v = NULL;
        }

        // Lowercase key
        for (char* kp = k; *kp; kp++) {
            if (*kp >= 'A' && *kp <= 'Z') *kp += 32;
        }

        // Unquote value if quoted
        char* unquoted_v = NULL;
        if (v) {
            size_t vlen = strlen(v);
            if (vlen >= 2 && v[0] == '"' && v[vlen - 1] == '"') {
                unquoted_v = unquote_string(v, vlen);
            }
        }
        const char* val = unquoted_v ? unquoted_v : (v ? v : "");

        if (strcmp(k, "name") == 0) {
            free(r.name);
            r.name = strdup(val);
        } else if (strcmp(k, "is_null") == 0) {
            r.is_null = true;
            r.nullable = true;
        } else if (strcmp(k, "example") == 0) {
            r.example = true;
        } else if (strcmp(k, "desc") == 0) {
            free(r.desc);
            r.desc = strdup(val);
        } else if (strcmp(k, "type") == 0) {
            r.type = mm_value_type_parse(val);
        } else if (strcmp(k, "raw") == 0) {
            r.raw = true;
        } else if (strcmp(k, "nullable") == 0) {
            r.nullable = true;
        } else if (strcmp(k, "allow_empty") == 0) {
            r.allow_empty = true;
        } else if (strcmp(k, "unique") == 0) {
            r.unique = true;
        } else if (strcmp(k, "default") == 0) {
            free(r.default_val);
            r.default_val = strdup(val);
        } else if (strcmp(k, "min") == 0) {
            free(r.min);
            r.min = strdup(val);
        } else if (strcmp(k, "max") == 0) {
            free(r.max);
            r.max = strdup(val);
        } else if (strcmp(k, "size") == 0) {
            char* end;
            long u = strtol(val, &end, 10);
            if (*end == '\0' && u >= 0) {
                r.size = (int)u;
            }
        } else if (strcmp(k, "enum") == 0) {
            r.type = MM_VALUE_ENUM;
            free(r.enums);
            r.enums = strdup(val);
        } else if (strcmp(k, "pattern") == 0) {
            free(r.pattern);
            r.pattern = strdup(val);
        } else if (strcmp(k, "location") == 0) {
            char* end;
            long d = strtol(val, &end, 10);
            if (*end == '\0' && d >= -12 && d <= 14) {
                r.location_offset = (int)d;
            }
        } else if (strcmp(k, "version") == 0) {
            char* end;
            long d = strtol(val, &end, 10);
            if (*end == '\0' && d >= 1 && d <= 10) {
                r.version = (int)d;
            }
        } else if (strcmp(k, "mime") == 0) {
            free(r.mime);
            r.mime = strdup(val);
        } else if (strcmp(k, "child_desc") == 0) {
            free(r.child_desc);
            r.child_desc = strdup(val);
        } else if (strcmp(k, "child_type") == 0) {
            r.child_type = mm_value_type_parse(val);
        } else if (strcmp(k, "child_raw") == 0) {
            r.child_raw = true;
        } else if (strcmp(k, "child_nullable") == 0) {
            r.child_nullable = true;
        } else if (strcmp(k, "child_allow_empty") == 0) {
            r.child_allow_empty = true;
        } else if (strcmp(k, "child_unique") == 0) {
            r.child_unique = true;
        } else if (strcmp(k, "child_default") == 0) {
            free(r.child_default);
            r.child_default = strdup(val);
        } else if (strcmp(k, "child_min") == 0) {
            free(r.child_min);
            r.child_min = strdup(val);
        } else if (strcmp(k, "child_max") == 0) {
            free(r.child_max);
            r.child_max = strdup(val);
        } else if (strcmp(k, "child_size") == 0) {
            char* end;
            long u = strtol(val, &end, 10);
            if (*end == '\0' && u >= 0) {
                r.child_size = (int)u;
            }
        } else if (strcmp(k, "child_enum") == 0) {
            free(r.child_enum);
            r.child_enum = strdup(val);
            r.child_type = MM_VALUE_ENUM;
        } else if (strcmp(k, "child_pattern") == 0) {
            free(r.child_pattern);
            r.child_pattern = strdup(val);
        } else if (strcmp(k, "child_location") == 0) {
            char* end;
            long d = strtol(val, &end, 10);
            if (*end == '\0' && d >= -12 && d <= 14) {
                r.child_location_offset = (int)d;
            }
        } else if (strcmp(k, "child_version") == 0) {
            char* end;
            long d = strtol(val, &end, 10);
            if (*end == '\0' && d >= 1 && d <= 10) {
                r.child_version = (int)d;
            }
        } else if (strcmp(k, "child_mime") == 0) {
            free(r.child_mime);
            r.child_mime = strdup(val);
        }

        free(k);
        free(v);
        free(unquoted_v);
    }

    free_split_parts(parts, part_count);
    return r;
}

uint8_t* mm_tag_bytes(const mm_tag_t* tag, size_t* out_len)
{
    uint8_t* buf = NULL;
    size_t cap = 0;
    size_t len = 0;

    if (tag->example) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KEXAMPLE | 1);
    }

    if (tag->is_null) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KISNULL | 1);
    }

    if (tag->nullable && !tag->is_inherit) {
        if (!tag->is_null) {
            if (len + 1 > cap) {
                cap = cap ? cap * 2 : 128;
                buf = (uint8_t*)realloc(buf, cap);
            }
            buf[len++] = (uint8_t)(MM_TAG_KNULLABLE | 1);
        }
    }

    if (tag->desc && !tag->is_inherit) {
        encode_string(&buf, &cap, &len, MM_TAG_KDESC, tag->desc);
    }

    if (tag->type != MM_VALUE_UNKNOWN && !tag->is_inherit) {
        if (is_bytes_simple_type(tag->type)) {
        } else {
            if ((tag->type == MM_VALUE_ARR && tag->size > 0) ||
                (tag->type == MM_VALUE_ENUM && tag->enums)) {
            } else {
                if (len + 2 > cap) {
                    cap = cap ? cap * 2 : 128;
                    buf = (uint8_t*)realloc(buf, cap);
                }
                buf[len++] = (uint8_t)(MM_TAG_KTYPE);
                buf[len++] = (uint8_t)(tag->type);
            }
        }
    }

    if (tag->raw && !tag->is_inherit) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KRAW | 1);
    }

    if (tag->allow_empty && !tag->is_inherit) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KALLOWEMPTY | 1);
    }

    if (tag->unique && !tag->is_inherit) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KUNIQUE | 1);
    }

    if (tag->default_val && !tag->is_inherit) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KDEFAULT, tag->default_val);
    }

    if (tag->min && !tag->is_inherit) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KMIN, tag->min);
    }

    if (tag->max && !tag->is_inherit) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KMAX, tag->max);
    }

    if (tag->size != 0 && !tag->is_inherit) {
        encode_u64(&buf, &cap, &len, MM_TAG_KSIZE, (uint64_t)tag->size);
    }

    if (tag->enums && !tag->is_inherit) {
        encode_string(&buf, &cap, &len, MM_TAG_KENUM, tag->enums);
    }

    if (tag->pattern && !tag->is_inherit) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KPATTERN, tag->pattern);
    }

    if (tag->location_offset != 0 && !tag->is_inherit) {
        char loc_buf[16];
        int loc_len = snprintf(loc_buf, sizeof(loc_buf), "%d", tag->location_offset);
        size_t l = (size_t)loc_len;
        if (len + 1 + l > cap) {
            cap = cap ? cap * 2 : 128;
            while (len + 1 + l > cap) cap *= 2;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KLOCATION | (int)l);
        memcpy(buf + len, loc_buf, l);
        len += l;
    }

    if (tag->version != MM_TAG_DEFAULT_VERSION && !tag->is_inherit) {
        encode_u64(&buf, &cap, &len, MM_TAG_KVERSION, (uint64_t)tag->version);
    }

    if (tag->mime && !tag->is_inherit) {
        uint8_t l = parse_mime(tag->mime);
        if (l < 7) {
            if (len + 1 > cap) {
                cap = cap ? cap * 2 : 128;
                buf = (uint8_t*)realloc(buf, cap);
            }
            buf[len++] = (uint8_t)(MM_TAG_KMIME | l);
        } else {
            if (len + 2 > cap) {
                cap = cap ? cap * 2 : 128;
                buf = (uint8_t*)realloc(buf, cap);
            }
            buf[len++] = (uint8_t)(MM_TAG_KMIME | 7);
            buf[len++] = l;
        }
    }

    if (tag->child_desc) {
        encode_string(&buf, &cap, &len, MM_TAG_KCHILDDESC, tag->child_desc);
    }

    if (tag->child_type != MM_VALUE_UNKNOWN) {
        if (is_simple_type(tag->child_type)) {
        } else {
            if ((tag->child_type == MM_VALUE_ARR && tag->child_size > 0) ||
                (tag->child_type == MM_VALUE_ENUM && tag->child_enum)) {
            } else {
                if (len + 2 > cap) {
                    cap = cap ? cap * 2 : 128;
                    buf = (uint8_t*)realloc(buf, cap);
                }
                buf[len++] = (uint8_t)(MM_TAG_KCHILDTYPE);
                buf[len++] = (uint8_t)(tag->child_type);
            }
        }
    }

    if (tag->child_raw) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KCHILDRAW | 1);
    }

    if (tag->child_nullable) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KCHILDNULLABLE | 1);
    }

    if (tag->child_allow_empty) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KCHILDALLOWEMPTY | 1);
    }

    if (tag->child_unique) {
        if (len + 1 > cap) {
            cap = cap ? cap * 2 : 128;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KCHILDUNIQUE | 1);
    }

    if (tag->child_default) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KCHILDDEFAULT, tag->child_default);
    }

    if (tag->child_min) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KCHILDMIN, tag->child_min);
    }

    if (tag->child_max) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KCHILDMAX, tag->child_max);
    }

    if (tag->child_size != 0) {
        encode_u64(&buf, &cap, &len, MM_TAG_KCHILDSIZE, (uint64_t)tag->child_size);
    }

    if (tag->child_enum) {
        encode_string(&buf, &cap, &len, MM_TAG_KCHILDENUM, tag->child_enum);
    }

    if (tag->child_pattern) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KCHILDPATTERN, tag->child_pattern);
    }

    if (tag->child_location_offset != 0) {
        char loc_buf[16];
        int loc_len = snprintf(loc_buf, sizeof(loc_buf), "%d", tag->child_location_offset);
        size_t l = (size_t)loc_len;
        if (len + 1 + l > cap) {
            cap = cap ? cap * 2 : 128;
            while (len + 1 + l > cap) cap *= 2;
            buf = (uint8_t*)realloc(buf, cap);
        }
        buf[len++] = (uint8_t)(MM_TAG_KCHILDLOCATION | (int)l);
        memcpy(buf + len, loc_buf, l);
        len += l;
    }

    if (tag->child_version != MM_TAG_DEFAULT_VERSION) {
        encode_u64(&buf, &cap, &len, MM_TAG_KCHILDVERSION, (uint64_t)tag->child_version);
    }

    if (tag->child_mime) {
        encode_string_simple(&buf, &cap, &len, MM_TAG_KCHILDMIME, tag->child_mime);
    }

    *out_len = len;
    return buf;
}

void mm_tag_merge(mm_tag_t* dst, const mm_tag_t* src)
{
    if (!dst || !src) return;

    if (src->is_null) {
        dst->is_null = src->is_null;
    }

    if (src->example) {
        dst->example = src->example;
    }

    if (src->desc) {
        free(dst->desc);
        dst->desc = strdup(src->desc);
    }

    if (src->type != MM_VALUE_UNKNOWN) {
        dst->type = src->type;
    }

    if (src->raw) {
        dst->raw = true;
    }

    if (src->nullable) {
        dst->nullable = true;
    }

    if (src->allow_empty) {
        dst->allow_empty = true;
    }

    if (src->unique) {
        dst->unique = true;
    }

    if (src->default_val) {
        free(dst->default_val);
        dst->default_val = strdup(src->default_val);
    }

    if (src->min) {
        free(dst->min);
        dst->min = strdup(src->min);
    }

    if (src->max) {
        free(dst->max);
        dst->max = strdup(src->max);
    }

    if (src->size != 0) {
        dst->size = src->size;
    }

    if (src->enums) {
        free(dst->enums);
        dst->enums = strdup(src->enums);
    }

    if (src->pattern) {
        free(dst->pattern);
        dst->pattern = strdup(src->pattern);
    }

    if (src->location_offset != 0) {
        dst->location_offset = src->location_offset;
    }

    if (src->version != MM_TAG_DEFAULT_VERSION) {
        dst->version = src->version;
    }

    if (src->mime) {
        free(dst->mime);
        dst->mime = strdup(src->mime);
    }

    if (src->child_desc) {
        free(dst->child_desc);
        dst->child_desc = strdup(src->child_desc);
    }

    if (src->child_type != MM_VALUE_UNKNOWN) {
        dst->child_type = src->child_type;
    }

    if (src->child_raw) {
        dst->child_raw = true;
    }

    if (src->child_nullable) {
        dst->child_nullable = true;
    }

    if (src->child_allow_empty) {
        dst->child_allow_empty = true;
    }

    if (src->child_unique) {
        dst->child_unique = true;
    }

    if (src->child_default) {
        free(dst->child_default);
        dst->child_default = strdup(src->child_default);
    }

    if (src->child_min) {
        free(dst->child_min);
        dst->child_min = strdup(src->child_min);
    }

    if (src->child_max) {
        free(dst->child_max);
        dst->child_max = strdup(src->child_max);
    }

    if (src->child_size != 0) {
        dst->child_size = src->child_size;
    }

    if (src->child_enum) {
        free(dst->child_enum);
        dst->child_enum = strdup(src->child_enum);
    }

    if (src->child_pattern) {
        free(dst->child_pattern);
        dst->child_pattern = strdup(src->child_pattern);
    }

    if (src->child_location_offset != 0) {
        dst->child_location_offset = src->child_location_offset;
    }

    if (src->child_version != MM_TAG_DEFAULT_VERSION) {
        dst->child_version = src->child_version;
    }

    if (src->child_mime) {
        free(dst->child_mime);
        dst->child_mime = strdup(src->child_mime);
    }
}