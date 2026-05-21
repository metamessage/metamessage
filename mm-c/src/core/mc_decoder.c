#include "mc_decoder.h"
#include "../ir/mc_tag.h"
#include "../ir/mc_value_type.h"
#include "mc_constants.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>

static uint8_t dec_read_byte(mm_decoder_t* d)
{
    if (d->offset >= d->size) {
        errno = EINVAL;
        return 0;
    }
    return d->data[d->offset++];
}

static const uint8_t* dec_read_bytes(mm_decoder_t* d, size_t len)
{
    if (d->offset + len > d->size) {
        errno = EINVAL;
        return NULL;
    }
    const uint8_t* p = d->data + d->offset;
    d->offset += len;
    return p;
}

static uint64_t dec_read_uint64_be(mm_decoder_t* d, int byte_len)
{
    uint64_t v = 0;
    for (int i = 0; i < byte_len; i++) {
        v = (v << 8) | dec_read_byte(d);
    }
    return v;
}

static const char* mime_id_to_string(uint8_t id)
{
    switch (id) {
        case 1:  return "image/jpeg";
        case 2:  return "image/png";
        case 3:  return "image/gif";
        case 4:  return "image/webp";
        case 5:  return "image/svg+xml";
        case 6:  return "image/avif";
        case 7:  return "image/bmp";
        case 8:  return "image/x-icon";
        case 9:  return "image/tiff";
        case 10: return "image/heic";
        case 11: return "image/heif";
        case 12: return "text/plain";
        case 13: return "text/html";
        case 14: return "text/css";
        case 15: return "text/javascript";
        case 16: return "application/json";
        case 17: return "text/csv";
        case 18: return "text/markdown";
        case 19: return "application/pdf";
        case 20: return "application/zip";
        case 21: return "application/gzip";
        case 22: return "application/x-tar";
        case 23: return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";
        case 24: return "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
        case 25: return "application/octet-stream";
        case 26: return "video/mp4";
        case 27: return "video/webm";
        case 28: return "video/mov";
        case 29: return "audio/mpeg";
        case 30: return "audio/wav";
        case 31: return "audio/flac";
        case 32: return "font/woff2";
        case 33: return "font/ttf";
        default: return "";
    }
}

static mm_node_t* dec_decode_simple(mm_decoder_t* d, uint8_t b)
{
    (void)d;
    int suffix = mm_suffix_of(b);
    mm_node_t* node = mm_node_new_value();
    mm_value_t* val = &node->data.value;
    mm_tag_init(&val->tag);

    switch (suffix) {
        case MM_SIMPLE_NULLBOOL:
            val->tag.type = MM_VALUE_BOOL;
            val->text = strdup("false");
            break;
        case MM_SIMPLE_NULLINT:
            val->tag.type = MM_VALUE_I;
            val->text = strdup("0");
            break;
        case MM_SIMPLE_NULLFLOAT:
            val->tag.type = MM_VALUE_F64;
            val->text = strdup("0.0");
            break;
        case MM_SIMPLE_NULLSTRING:
            val->tag.type = MM_VALUE_STR;
            val->text = strdup("");
            break;
        case MM_SIMPLE_NULLBYTES:
            val->tag.type = MM_VALUE_BYTES;
            val->text = strdup("");
            break;
        case MM_SIMPLE_FALSE:
            val->tag.type = MM_VALUE_BOOL;
            val->text = strdup("false");
            break;
        case MM_SIMPLE_TRUE:
            val->tag.type = MM_VALUE_BOOL;
            val->text = strdup("true");
            break;
        default:
            val->tag.type = MM_VALUE_STR;
            val->text = strdup("");
            break;
    }
    return node;
}

static mm_node_t* dec_decode_int(mm_decoder_t* d, uint8_t b, int is_positive)
{
    int byte_len = mm_int_len(b);
    uint64_t uv;

    if (byte_len == 0) {
        uv = (uint64_t)mm_suffix_of(b);
    } else {
        uv = dec_read_uint64_be(d, byte_len);
    }

    mm_node_t* node = mm_node_new_value();
    mm_value_t* val = &node->data.value;
    mm_tag_init(&val->tag);

    if (is_positive) {
        val->tag.type = MM_VALUE_I;
        char buf[24];
        snprintf(buf, sizeof(buf), "%llu", (unsigned long long)uv);
        val->text = strdup(buf);
    } else {
        val->tag.type = MM_VALUE_I;
        char buf[24];
        snprintf(buf, sizeof(buf), "-%llu", (unsigned long long)uv);
        val->text = strdup(buf);
    }

    return node;
}

static mm_node_t* dec_decode_float(mm_decoder_t* d, uint8_t b)
{
    int byte_len = mm_float_len(b);
    int suffix = mm_suffix_of(b);
    int is_negative = (b & 0x10) != 0;

    double v;

    if (byte_len == 0) {
        v = (double)suffix / 10.0;
        if (is_negative) {
            v = -v;
        }
    } else {
        int8_t exp = (int8_t)dec_read_byte(d);
        uint64_t mantissa;
        if (byte_len >= 1 && byte_len <= 8) {
            mantissa = dec_read_uint64_be(d, byte_len);
        } else {
            mantissa = 0;
        }
        char buf[128];
        char mantissa_str[32];
        snprintf(mantissa_str, sizeof(mantissa_str), "%llu", (unsigned long long)mantissa);
        int decimal_pos = (int)strlen(mantissa_str) + (int)exp;
        if (decimal_pos <= 0) {
            char pad[64];
            int pad_len = -decimal_pos;
            memset(pad, '0', (size_t)pad_len);
            pad[pad_len] = '\0';
            snprintf(buf, sizeof(buf), "0.%s%s", pad, mantissa_str);
        } else if (decimal_pos > 0 && decimal_pos < (int)strlen(mantissa_str)) {
            size_t m_len = strlen(mantissa_str);
            size_t int_part = (size_t)decimal_pos;
            char int_buf[64], frac_buf[64];
            memcpy(int_buf, mantissa_str, int_part);
            int_buf[int_part] = '\0';
            memcpy(frac_buf, mantissa_str + int_part, m_len - int_part + 1);
            snprintf(buf, sizeof(buf), "%s.%s", int_buf, frac_buf);
        } else {
            int trailing = decimal_pos - (int)strlen(mantissa_str);
            char pad[64];
            memset(pad, '0', (size_t)trailing);
            pad[trailing] = '\0';
            snprintf(buf, sizeof(buf), "%s%s", mantissa_str, pad);
        }
        if (is_negative) {
            char neg_buf[132];
            snprintf(neg_buf, sizeof(neg_buf), "-%s", buf);
            v = -strtod(buf, NULL);
        } else {
            v = strtod(buf, NULL);
        }
    }

    mm_node_t* node = mm_node_new_value();
    mm_value_t* val = &node->data.value;
    mm_tag_init(&val->tag);
    val->tag.type = MM_VALUE_F64;
    char text_buf[64];
    snprintf(text_buf, sizeof(text_buf), "%g", v);
    val->text = strdup(text_buf);

    return node;
}

static mm_node_t* dec_decode_string(mm_decoder_t* d, uint8_t b)
{
    int extra_len = mm_string_extra_len(b);
    int inline_len = mm_string_inline_len(b);
    size_t str_len;

    if (extra_len == 0) {
        str_len = (size_t)inline_len;
    } else if (extra_len == 1) {
        str_len = (size_t)dec_read_byte(d);
    } else {
        uint8_t hi = dec_read_byte(d);
        uint8_t lo = dec_read_byte(d);
        str_len = ((size_t)hi << 8) | (size_t)lo;
    }

    char* text = NULL;
    if (str_len > 0) {
        const uint8_t* raw = dec_read_bytes(d, str_len);
        if (!raw) {
            text = strdup("");
        } else {
            text = (char*)malloc(str_len + 1);
            memcpy(text, raw, str_len);
            text[str_len] = '\0';
        }
    } else {
        text = strdup("");
    }

    mm_node_t* node = mm_node_new_value();
    mm_value_t* val = &node->data.value;
    mm_tag_init(&val->tag);
    val->tag.type = MM_VALUE_STR;
    val->text = text;

    return node;
}

static mm_node_t* dec_decode_bytes(mm_decoder_t* d, uint8_t b)
{
    int extra_len = mm_bytes_extra_len(b);
    int inline_len = mm_bytes_inline_len(b);
    size_t bytes_len;

    if (extra_len == 0) {
        bytes_len = (size_t)inline_len;
    } else if (extra_len == 1) {
        bytes_len = (size_t)dec_read_byte(d);
    } else {
        uint8_t hi = dec_read_byte(d);
        uint8_t lo = dec_read_byte(d);
        bytes_len = ((size_t)hi << 8) | (size_t)lo;
    }

    mm_node_t* node = mm_node_new_value();
    mm_value_t* val = &node->data.value;
    mm_tag_init(&val->tag);
    val->tag.type = MM_VALUE_BYTES;
    val->text = strdup("");

    if (bytes_len > 0) {
        const uint8_t* raw = dec_read_bytes(d, bytes_len);
        (void)raw;
    }

    return node;
}

static mm_node_t* dec_decode_array(mm_decoder_t* d, size_t total_len, mm_tag_t* parent_tag)
{
    mm_node_t* node = mm_node_new_array();
    mm_array_t* arr = &node->data.array;
    if (parent_tag) {
        mm_tag_inherit(&arr->tag, parent_tag);
    }
    if (arr->tag.type == MM_VALUE_UNKNOWN) {
        arr->tag.type = MM_VALUE_VEC;
    }

    size_t end_offset = d->offset + total_len;
    while (d->offset < end_offset) {
        mm_node_t* item = mm_decoder_decode(d);
        if (!item) break;
        mm_array_add_item(node, item);
    }

    if (d->offset < end_offset) {
        d->offset = end_offset;
    }

    return node;
}

static mm_node_t* dec_decode_object(mm_decoder_t* d, size_t total_len, mm_tag_t* parent_tag)
{
    mm_node_t* node = mm_node_new_object();
    mm_object_t* obj = &node->data.object;
    if (parent_tag) {
        mm_tag_inherit(&obj->tag, parent_tag);
    }
    if (obj->tag.type == MM_VALUE_UNKNOWN) {
        obj->tag.type = MM_VALUE_OBJ;
    }

    size_t start_offset = d->offset;
    size_t end_offset = start_offset + total_len;

    mm_node_t* key_array = mm_decoder_decode(d);
    if (!key_array || key_array->type != MM_NODE_ARRAY) {
        if (key_array) mm_node_free(key_array);
        d->offset = end_offset;
        return node;
    }

    for (size_t i = 0; i < key_array->data.array.item_count && d->offset < end_offset; i++) {
        mm_node_t* key_item = key_array->data.array.items[i];
        if (key_item->type != MM_NODE_VALUE || !key_item->data.value.text) {
            mm_decoder_decode(d);
            continue;
        }
        mm_node_t* val_node = mm_decoder_decode(d);
        if (val_node) {
            mm_object_add_field(node, key_item->data.value.text, val_node);
        }
    }

    mm_node_free(key_array);
    d->offset = end_offset;

    return node;
}

static mm_node_t* dec_decode_container(mm_decoder_t* d, uint8_t b, mm_tag_t* parent_tag)
{
    int is_array = mm_is_array_container(b);
    int extra_len = mm_container_extra_len(b);
    int inline_len = mm_container_inline_len(b);
    size_t container_len;

    if (extra_len == 0) {
        container_len = (size_t)inline_len;
    } else if (extra_len == 1) {
        container_len = (size_t)dec_read_byte(d);
    } else {
        uint8_t hi = dec_read_byte(d);
        uint8_t lo = dec_read_byte(d);
        container_len = ((size_t)hi << 8) | (size_t)lo;
    }

    if (is_array) {
        return dec_decode_array(d, container_len, parent_tag);
    } else {
        return dec_decode_object(d, container_len, parent_tag);
    }
}

static bool dec_parse_one_tag_entry(mm_decoder_t* d, mm_tag_t* tag)
{
    uint8_t b = dec_read_byte(d);
    if (errno) return false;

    int key = b & 0xF8;
    int payload = b & 0x07;

    switch (key) {
        case MM_TAG_KISNULL:
            tag->is_null = (payload & 0x01) != 0;
            if (tag->is_null) {
                tag->nullable = true;
            }
            return true;

        case MM_TAG_KEXAMPLE:
            tag->example = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KDESC: {
            size_t str_len;
            if (payload <= 5) {
                str_len = (size_t)payload;
            } else if (payload == 6) {
                str_len = (size_t)dec_read_byte(d);
            } else {
                uint8_t hi = dec_read_byte(d);
                uint8_t lo = dec_read_byte(d);
                str_len = ((size_t)hi << 8) | (size_t)lo;
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->desc);
                tag->desc = (char*)malloc(str_len + 1);
                memcpy(tag->desc, raw, str_len);
                tag->desc[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KTYPE:
            tag->type = (mm_value_type_t)dec_read_byte(d);
            return true;

        case MM_TAG_KRAW:
            tag->raw = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KNULLABLE:
            tag->nullable = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KALLOWEMPTY:
            tag->allow_empty = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KUNIQUE:
            tag->unique = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KDEFAULT: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->default_val);
                tag->default_val = (char*)malloc(str_len + 1);
                memcpy(tag->default_val, raw, str_len);
                tag->default_val[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KMIN: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->min);
                tag->min = (char*)malloc(str_len + 1);
                memcpy(tag->min, raw, str_len);
                tag->min[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KMAX: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->max);
                tag->max = (char*)malloc(str_len + 1);
                memcpy(tag->max, raw, str_len);
                tag->max[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KSIZE: {
            int nbytes = payload + 1;
            uint64_t val = 0;
            for (int i = 0; i < nbytes; i++) {
                val = (val << 8) | dec_read_byte(d);
            }
            tag->size = (int)val;
            return true;
        }

        case MM_TAG_KENUM: {
            tag->type = MM_VALUE_ENUM;
            size_t str_len;
            if (payload <= 5) {
                str_len = (size_t)payload;
            } else if (payload == 6) {
                str_len = (size_t)dec_read_byte(d);
            } else {
                uint8_t hi = dec_read_byte(d);
                uint8_t lo = dec_read_byte(d);
                str_len = ((size_t)hi << 8) | (size_t)lo;
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->enums);
                tag->enums = (char*)malloc(str_len + 1);
                memcpy(tag->enums, raw, str_len);
                tag->enums[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KPATTERN: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->pattern);
                tag->pattern = (char*)malloc(str_len + 1);
                memcpy(tag->pattern, raw, str_len);
                tag->pattern[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KLOCATION: {
            size_t str_len = (size_t)payload;
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                char buf[16];
                memcpy(buf, raw, str_len);
                buf[str_len] = '\0';
                tag->location_offset = (int)strtol(buf, NULL, 10);
            }
            return true;
        }

        case MM_TAG_KVERSION: {
            int nbytes = payload + 1;
            uint64_t val = 0;
            for (int i = 0; i < nbytes; i++) {
                val = (val << 8) | dec_read_byte(d);
            }
            tag->version = (int)val;
            return true;
        }

        case MM_TAG_KMIME: {
            uint8_t mime_id;
            if (payload < 7) {
                mime_id = (uint8_t)payload;
            } else {
                mime_id = dec_read_byte(d);
            }
            free(tag->mime);
            tag->mime = strdup(mime_id_to_string(mime_id));
            return true;
        }

        case MM_TAG_KCHILDDESC: {
            size_t str_len;
            if (payload <= 5) {
                str_len = (size_t)payload;
            } else if (payload == 6) {
                str_len = (size_t)dec_read_byte(d);
            } else {
                uint8_t hi = dec_read_byte(d);
                uint8_t lo = dec_read_byte(d);
                str_len = ((size_t)hi << 8) | (size_t)lo;
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_desc);
                tag->child_desc = (char*)malloc(str_len + 1);
                memcpy(tag->child_desc, raw, str_len);
                tag->child_desc[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDTYPE:
            tag->child_type = (mm_value_type_t)dec_read_byte(d);
            return true;

        case MM_TAG_KCHILDRAW:
            tag->child_raw = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KCHILDNULLABLE:
            tag->child_nullable = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KCHILDALLOWEMPTY:
            tag->child_allow_empty = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KCHILDUNIQUE:
            tag->child_unique = (payload & 0x01) != 0;
            return true;

        case MM_TAG_KCHILDDEFAULT: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_default);
                tag->child_default = (char*)malloc(str_len + 1);
                memcpy(tag->child_default, raw, str_len);
                tag->child_default[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDMIN: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_min);
                tag->child_min = (char*)malloc(str_len + 1);
                memcpy(tag->child_min, raw, str_len);
                tag->child_min[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDMAX: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_max);
                tag->child_max = (char*)malloc(str_len + 1);
                memcpy(tag->child_max, raw, str_len);
                tag->child_max[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDSIZE: {
            int nbytes = payload + 1;
            uint64_t val = 0;
            for (int i = 0; i < nbytes; i++) {
                val = (val << 8) | dec_read_byte(d);
            }
            tag->child_size = (int)val;
            return true;
        }

        case MM_TAG_KCHILDENUM: {
            tag->child_type = MM_VALUE_ENUM;
            size_t str_len;
            if (payload <= 5) {
                str_len = (size_t)payload;
            } else if (payload == 6) {
                str_len = (size_t)dec_read_byte(d);
            } else {
                uint8_t hi = dec_read_byte(d);
                uint8_t lo = dec_read_byte(d);
                str_len = ((size_t)hi << 8) | (size_t)lo;
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_enum);
                tag->child_enum = (char*)malloc(str_len + 1);
                memcpy(tag->child_enum, raw, str_len);
                tag->child_enum[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDPATTERN: {
            size_t str_len;
            if (payload < 7) {
                str_len = (size_t)payload;
            } else {
                str_len = (size_t)dec_read_byte(d);
            }
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                free(tag->child_pattern);
                tag->child_pattern = (char*)malloc(str_len + 1);
                memcpy(tag->child_pattern, raw, str_len);
                tag->child_pattern[str_len] = '\0';
            }
            return true;
        }

        case MM_TAG_KCHILDLOCATION: {
            size_t str_len = (size_t)payload;
            const uint8_t* raw = dec_read_bytes(d, str_len);
            if (raw) {
                char buf[16];
                memcpy(buf, raw, str_len);
                buf[str_len] = '\0';
                tag->child_location_offset = (int)strtol(buf, NULL, 10);
            }
            return true;
        }

        case MM_TAG_KCHILDVERSION: {
            int nbytes = payload + 1;
            uint64_t val = 0;
            for (int i = 0; i < nbytes; i++) {
                val = (val << 8) | dec_read_byte(d);
            }
            tag->child_version = (int)val;
            return true;
        }

        case MM_TAG_KCHILDMIME: {
            uint8_t mime_id;
            if (payload < 7) {
                mime_id = (uint8_t)payload;
            } else {
                mime_id = dec_read_byte(d);
            }
            free(tag->child_mime);
            tag->child_mime = strdup(mime_id_to_string(mime_id));
            return true;
        }

        default:
            return true;
    }
}

static mm_tag_t dec_read_tag_bytes(mm_decoder_t* d)
{
    mm_tag_t tag;
    mm_tag_init(&tag);

    uint8_t first = dec_read_byte(d);
    if (errno) return tag;

    size_t tag_bytes_len;
    if (first < 254) {
        tag_bytes_len = (size_t)first;
    } else if (first == 254) {
        tag_bytes_len = (size_t)dec_read_byte(d);
    } else {
        uint8_t hi = dec_read_byte(d);
        uint8_t lo = dec_read_byte(d);
        tag_bytes_len = ((size_t)hi << 8) | (size_t)lo;
    }

    size_t end_offset = d->offset + tag_bytes_len;

    while (d->offset < end_offset) {
        if (!dec_parse_one_tag_entry(d, &tag)) {
            break;
        }
        if (errno) break;
    }

    if (d->offset < end_offset) {
        d->offset = end_offset;
    }

    return tag;
}

static mm_node_t* dec_decode_tag(mm_decoder_t* d, uint8_t b, mm_tag_t* parent_tag)
{
    int extra_len = mm_tag_extra_len(b);
    int inline_len = mm_tag_inline_len(b);
    size_t total_len;

    if (extra_len == 0) {
        total_len = (size_t)inline_len;
    } else if (extra_len == 1) {
        total_len = (size_t)dec_read_byte(d);
    } else {
        uint8_t hi = dec_read_byte(d);
        uint8_t lo = dec_read_byte(d);
        total_len = ((size_t)hi << 8) | (size_t)lo;
    }

    size_t tag_data_start = d->offset;

    mm_tag_t tag = dec_read_tag_bytes(d);
    if (parent_tag) {
        mm_tag_inherit(&tag, parent_tag);
    }

    size_t tag_bytes_consumed = d->offset - tag_data_start;

    mm_node_t* inner = NULL;

    if (tag.is_null) {
        inner = mm_node_new_value();
        mm_value_t* val = &inner->data.value;
        mm_tag_cleanup(&val->tag);
        val->tag = tag;

        switch (tag.type) {
            case MM_VALUE_I:
            case MM_VALUE_I8:
            case MM_VALUE_I16:
            case MM_VALUE_I32:
            case MM_VALUE_I64:
            case MM_VALUE_U:
            case MM_VALUE_U8:
            case MM_VALUE_U16:
            case MM_VALUE_U32:
            case MM_VALUE_U64:
                val->text = strdup("0");
                break;
            case MM_VALUE_F32:
            case MM_VALUE_F64:
                val->text = strdup("0.0");
                break;
            case MM_VALUE_STR:
            case MM_VALUE_BYTES:
            case MM_VALUE_EMAIL:
            case MM_VALUE_UUID:
            case MM_VALUE_DECIMAL:
            case MM_VALUE_URL:
            case MM_VALUE_BIGINT:
                val->text = strdup("");
                break;
            case MM_VALUE_DATETIME:
            case MM_VALUE_DATE:
            case MM_VALUE_TIME:
                val->text = strdup("");
                break;
            case MM_VALUE_IP:
                val->text = strdup("");
                break;
            default:
                val->text = strdup("");
                break;
        }

        size_t remaining = total_len - tag_bytes_consumed;
        if (d->offset + remaining <= d->size) {
            d->offset += remaining;
        }
    } else {
        inner = mm_decoder_decode(d);
        if (inner) {
            if (inner->type == MM_NODE_VALUE) {
                mm_tag_merge(&inner->data.value.tag, &tag);
                mm_tag_cleanup(&tag);
                if (inner->data.value.tag.type == MM_VALUE_UNKNOWN) {
                    inner->data.value.tag.type = MM_VALUE_STR;
                }
            } else if (inner->type == MM_NODE_ARRAY) {
                mm_tag_merge(&inner->data.array.tag, &tag);
                mm_tag_cleanup(&tag);
                if (inner->data.array.tag.type == MM_VALUE_UNKNOWN) {
                    inner->data.array.tag.type = MM_VALUE_VEC;
                }
            } else if (inner->type == MM_NODE_OBJECT) {
                mm_tag_merge(&inner->data.object.tag, &tag);
                mm_tag_cleanup(&tag);
                if (inner->data.object.tag.type == MM_VALUE_UNKNOWN) {
                    inner->data.object.tag.type = MM_VALUE_OBJ;
                }
            }
        }
    }

    return inner;
}

mm_decoder_t* mm_decoder_new(const uint8_t* data, size_t size)
{
    mm_decoder_t* d = (mm_decoder_t*)malloc(sizeof(mm_decoder_t));
    if (!d) return NULL;
    d->data = data;
    d->size = size;
    d->offset = 0;
    return d;
}

void mm_decoder_free(mm_decoder_t* d)
{
    free(d);
}

mm_node_t* mm_decoder_decode(mm_decoder_t* d)
{
    if (!d || d->offset >= d->size) {
        return NULL;
    }

    errno = 0;
    uint8_t b = dec_read_byte(d);
    if (errno) return NULL;

    int prefix = mm_prefix_of(b);

    switch (prefix) {
        case MM_PREFIX_SIMPLE:
            return dec_decode_simple(d, b);
        case MM_PREFIX_POSITIVEINT:
            return dec_decode_int(d, b, 1);
        case MM_PREFIX_NEGATIVEINT:
            return dec_decode_int(d, b, 0);
        case MM_PREFIX_FLOAT:
            return dec_decode_float(d, b);
        case MM_PREFIX_STRING:
            return dec_decode_string(d, b);
        case MM_PREFIX_BYTES:
            return dec_decode_bytes(d, b);
        case MM_PREFIX_CONTAINER:
            return dec_decode_container(d, b, NULL);
        case MM_PREFIX_TAG:
            return dec_decode_tag(d, b, NULL);
        default:
            return NULL;
    }
}