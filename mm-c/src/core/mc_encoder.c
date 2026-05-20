#include "mc_encoder.h"
#include "../ir/mc_tag.h"
#include "../ir/mc_value_type.h"
#include "mc_constants.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

typedef struct {
    uint8_t* buf;
    size_t size;
    size_t cap;
} encoder_t;

static void enc_encode_node_object(encoder_t* e, mm_object_t* obj);
static void enc_encode_node_array(encoder_t* e, mm_array_t* arr);
static void enc_encode_node_value(encoder_t* e, mm_value_t* val);
static void enc_encode_node_doc(encoder_t* e, mm_doc_t* doc);

static void enc_write_byte(encoder_t* e, uint8_t b) {
    if (e->size >= e->cap) {
        e->cap = e->cap == 0 ? 64 : e->cap * 2;
        e->buf = realloc(e->buf, e->cap);
    }
    e->buf[e->size++] = b;
}

static void enc_write_bytes(encoder_t* e, const uint8_t* data, size_t len) {
    if (len == 0) return;
    while (e->size + len > e->cap) {
        e->cap = e->cap == 0 ? 64 : e->cap * 2;
        e->buf = realloc(e->buf, e->cap);
    }
    memcpy(e->buf + e->size, data, len);
    e->size += len;
}

static void enc_encode_int(encoder_t* e, uint8_t sign, const char* text) {
    uint64_t uv;
    if (sign == MM_PREFIX_NEGATIVEINT) {
        int64_t sv = (int64_t)strtoll(text, NULL, 10);
        if (sv < 0) {
            uv = (uint64_t)(-sv);
        } else {
            uv = (uint64_t)sv;
            sign = MM_PREFIX_POSITIVEINT;
        }
    } else {
        uv = strtoull(text, NULL, 10);
    }

    if (uv <= 23) {
        enc_write_byte(e, (uint8_t)(sign | uv));
    } else if (uv <= 0xFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_INTLEN1BYTE));
        enc_write_byte(e, (uint8_t)uv);
    } else if (uv <= 0xFFFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_INTLEN2BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    } else if (uv <= 0xFFFFFFFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_INTLEN4BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 24));
        enc_write_byte(e, (uint8_t)(uv >> 16));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    } else {
        enc_write_byte(e, (uint8_t)(sign | MM_INTLEN8BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 56));
        enc_write_byte(e, (uint8_t)(uv >> 48));
        enc_write_byte(e, (uint8_t)(uv >> 40));
        enc_write_byte(e, (uint8_t)(uv >> 32));
        enc_write_byte(e, (uint8_t)(uv >> 24));
        enc_write_byte(e, (uint8_t)(uv >> 16));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    }
}

static void enc_encode_float(encoder_t* e, const char* text) {
    double d = strtod(text, NULL);
    uint64_t uv;
    memcpy(&uv, &d, sizeof(uv));

    uint8_t sign = MM_PREFIX_FLOAT;

    if (uv <= 7) {
        enc_write_byte(e, (uint8_t)(sign | uv));
    } else if (uv <= 0xFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_FLOATLEN1BYTE));
        enc_write_byte(e, (uint8_t)uv);
    } else if (uv <= 0xFFFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_FLOATLEN2BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    } else if (uv <= 0xFFFFFFFF) {
        enc_write_byte(e, (uint8_t)(sign | MM_FLOATLEN4BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 24));
        enc_write_byte(e, (uint8_t)(uv >> 16));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    } else {
        enc_write_byte(e, (uint8_t)(sign | MM_FLOATLEN8BYTE));
        enc_write_byte(e, (uint8_t)(uv >> 56));
        enc_write_byte(e, (uint8_t)(uv >> 48));
        enc_write_byte(e, (uint8_t)(uv >> 40));
        enc_write_byte(e, (uint8_t)(uv >> 32));
        enc_write_byte(e, (uint8_t)(uv >> 24));
        enc_write_byte(e, (uint8_t)(uv >> 16));
        enc_write_byte(e, (uint8_t)(uv >> 8));
        enc_write_byte(e, (uint8_t)uv);
    }
}

static void enc_encode_string(encoder_t* e, const char* s) {
    size_t len = strlen(s);
    uint8_t sign = MM_PREFIX_STRING;

    if (len <= 29) {
        enc_write_byte(e, (uint8_t)(sign | len));
    } else if (len <= 255) {
        enc_write_byte(e, (uint8_t)(sign | MM_STRINGLEN1BYTE));
        enc_write_byte(e, (uint8_t)len);
    } else if (len <= 65535) {
        enc_write_byte(e, (uint8_t)(sign | MM_STRINGLEN2BYTE));
        enc_write_byte(e, (uint8_t)(len >> 8));
        enc_write_byte(e, (uint8_t)(len & 0xFF));
    }

    enc_write_bytes(e, (const uint8_t*)s, len);
}

static void enc_encode_bytes(encoder_t* e, const uint8_t* data, size_t len) {
    uint8_t sign = MM_PREFIX_BYTES;

    if (len <= 29) {
        enc_write_byte(e, (uint8_t)(sign | len));
    } else if (len <= 255) {
        enc_write_byte(e, (uint8_t)(sign | MM_BYTESLEN1BYTE));
        enc_write_byte(e, (uint8_t)len);
    } else if (len <= 65535) {
        enc_write_byte(e, (uint8_t)(sign | MM_BYTESLEN2BYTE));
        enc_write_byte(e, (uint8_t)(len >> 8));
        enc_write_byte(e, (uint8_t)(len & 0xFF));
    }

    enc_write_bytes(e, data, len);
}

static void enc_encode_simple(encoder_t* e, uint8_t value) {
    enc_write_byte(e, (uint8_t)(MM_PREFIX_SIMPLE | value));
}

static void enc_encode_bool(encoder_t* e, const char* text) {
    if (strcmp(text, "true") == 0) {
        enc_encode_simple(e, MM_SIMPLE_TRUE);
    } else {
        enc_encode_simple(e, MM_SIMPLE_FALSE);
    }
}

static void enc_encode_array(encoder_t* e, const uint8_t* data, size_t len) {
    uint8_t sign = (uint8_t)(MM_PREFIX_CONTAINER | MM_CONTAINER_ARRAY);

    if (len <= 13) {
        enc_write_byte(e, (uint8_t)(sign | len));
    } else if (len <= 255) {
        enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN1BYTE));
        enc_write_byte(e, (uint8_t)len);
    } else if (len <= 65535) {
        enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN2BYTE));
        enc_write_byte(e, (uint8_t)(len >> 8));
        enc_write_byte(e, (uint8_t)(len & 0xFF));
    }

    enc_write_bytes(e, data, len);
}

static void enc_encode_container(encoder_t* e, uint8_t container_type, const uint8_t* data, size_t len) {
    uint8_t sign = (uint8_t)(MM_PREFIX_CONTAINER | container_type);

    if (len <= 13) {
        enc_write_byte(e, (uint8_t)(sign | len));
    } else if (len <= 255) {
        enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN1BYTE));
        enc_write_byte(e, (uint8_t)len);
    } else if (len <= 65535) {
        enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN2BYTE));
        enc_write_byte(e, (uint8_t)(len >> 8));
        enc_write_byte(e, (uint8_t)(len & 0xFF));
    }

    enc_write_bytes(e, data, len);
}

static void enc_encode_tag(encoder_t* e, mm_tag_t* tag, const uint8_t* inner_data, size_t inner_len) {
    size_t tag_len;
    uint8_t* tag_bytes = mm_tag_bytes(tag, &tag_len);

    if (tag_bytes == NULL || tag_len == 0) {
        free(tag_bytes);
        enc_write_bytes(e, inner_data, inner_len);
        return;
    }

    encoder_t tag_enc = {0};

    if (tag_len < 254) {
        enc_write_byte(&tag_enc, (uint8_t)tag_len);
    } else if (tag_len < 257) {
        enc_write_byte(&tag_enc, 254);
        enc_write_byte(&tag_enc, (uint8_t)tag_len);
    } else {
        enc_write_byte(&tag_enc, 255);
        enc_write_byte(&tag_enc, (uint8_t)(tag_len >> 8));
        enc_write_byte(&tag_enc, (uint8_t)(tag_len & 0xFF));
    }

    enc_write_bytes(&tag_enc, tag_bytes, tag_len);
    free(tag_bytes);

    size_t combined_len = tag_enc.size + inner_len;

    if (combined_len < MM_TAGLEN1BYTE) {
        enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | combined_len));
    } else if (combined_len < 256) {
        enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | MM_TAGLEN1BYTE));
        enc_write_byte(e, (uint8_t)combined_len);
    } else if (combined_len < 65536) {
        enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | MM_TAGLEN2BYTE));
        enc_write_byte(e, (uint8_t)(combined_len >> 8));
        enc_write_byte(e, (uint8_t)(combined_len & 0xFF));
    }

    enc_write_bytes(e, tag_enc.buf, tag_enc.size);
    enc_write_bytes(e, inner_data, inner_len);

    free(tag_enc.buf);
}

static void enc_encode_node_value(encoder_t* e, mm_value_t* val) {
    encoder_t tmp = {0};

    switch (val->tag.type) {
        case MM_VALUE_STR:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
            } else {
                enc_encode_string(&tmp, val->text);
            }
            break;

        case MM_VALUE_BOOL:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLBOOL);
            } else {
                enc_encode_bool(&tmp, val->text);
            }
            break;

        case MM_VALUE_I:
        case MM_VALUE_I8:
        case MM_VALUE_I16:
        case MM_VALUE_I32:
        case MM_VALUE_I64: {
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
            } else {
                int64_t sv = (int64_t)strtoll(val->text, NULL, 10);
                if (sv < 0) {
                    enc_encode_int(&tmp, MM_PREFIX_NEGATIVEINT, val->text);
                } else {
                    enc_encode_int(&tmp, MM_PREFIX_POSITIVEINT, val->text);
                }
            }
            break;
        }

        case MM_VALUE_U:
        case MM_VALUE_U8:
        case MM_VALUE_U16:
        case MM_VALUE_U32:
        case MM_VALUE_U64:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
            } else {
                enc_encode_int(&tmp, MM_PREFIX_POSITIVEINT, val->text);
            }
            break;

        case MM_VALUE_F32:
        case MM_VALUE_F64:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLFLOAT);
            } else {
                enc_encode_float(&tmp, val->text);
            }
            break;

        case MM_VALUE_BYTES:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLBYTES);
            } else {
                enc_encode_bytes(&tmp, (const uint8_t*)val->text, strlen(val->text));
            }
            break;

        case MM_VALUE_EMAIL:
        case MM_VALUE_URL:
        case MM_VALUE_DECIMAL:
        case MM_VALUE_BIGINT:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
            } else {
                enc_encode_string(&tmp, val->text);
            }
            break;

        case MM_VALUE_UUID:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLBYTES);
            } else {
                enc_encode_string(&tmp, val->text);
            }
            break;

        case MM_VALUE_DATE:
        case MM_VALUE_TIME:
        case MM_VALUE_DATETIME:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
            } else {
                enc_encode_string(&tmp, val->text);
            }
            break;

        case MM_VALUE_IP:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
            } else {
                enc_encode_string(&tmp, val->text);
            }
            break;

        case MM_VALUE_ENUM:
            if (val->tag.is_null) {
                enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
            } else {
                enc_encode_int(&tmp, MM_PREFIX_POSITIVEINT, val->text);
            }
            break;

        default:
            enc_encode_string(&tmp, val->text);
            break;
    }

    enc_encode_tag(e, &val->tag, tmp.buf, tmp.size);
    free(tmp.buf);
}

static void enc_encode_node_array(encoder_t* e, mm_array_t* arr) {
    encoder_t items = {0};

    for (size_t i = 0; i < arr->item_count; i++) {
        encoder_t tmp = {0};
        mm_node_t* item = arr->items[i];

        switch (item->type) {
            case MM_NODE_OBJECT:
                enc_encode_node_object(&tmp, &item->data.object);
                break;
            case MM_NODE_ARRAY:
                enc_encode_node_array(&tmp, &item->data.array);
                break;
            case MM_NODE_VALUE:
                enc_encode_node_value(&tmp, &item->data.value);
                break;
            case MM_NODE_DOC:
                enc_encode_node_doc(&tmp, &item->data.doc);
                break;
            default:
                break;
        }

        enc_write_bytes(&items, tmp.buf, tmp.size);
        free(tmp.buf);
    }

    encoder_t container_enc = {0};
    enc_encode_array(&container_enc, items.buf, items.size);
    free(items.buf);

    enc_encode_tag(e, &arr->tag, container_enc.buf, container_enc.size);
    free(container_enc.buf);
}

static void enc_encode_node_object(encoder_t* e, mm_object_t* obj) {
    encoder_t keys = {0};
    encoder_t vals = {0};

    for (size_t i = 0; i < obj->field_count; i++) {
        mm_field_t* field = &obj->fields[i];

        encoder_t tmp_val = {0};
        switch (field->value->type) {
            case MM_NODE_OBJECT:
                enc_encode_node_object(&tmp_val, &field->value->data.object);
                break;
            case MM_NODE_ARRAY:
                enc_encode_node_array(&tmp_val, &field->value->data.array);
                break;
            case MM_NODE_VALUE:
                enc_encode_node_value(&tmp_val, &field->value->data.value);
                break;
            case MM_NODE_DOC:
                enc_encode_node_doc(&tmp_val, &field->value->data.doc);
                break;
            default:
                break;
        }
        enc_write_bytes(&vals, tmp_val.buf, tmp_val.size);
        free(tmp_val.buf);

        encoder_t tmp_key = {0};
        enc_encode_string(&tmp_key, field->key);
        enc_write_bytes(&keys, tmp_key.buf, tmp_key.size);
        free(tmp_key.buf);
    }

    encoder_t arr_enc = {0};
    enc_encode_array(&arr_enc, keys.buf, keys.size);
    free(keys.buf);

    encoder_t combined = {0};
    enc_write_bytes(&combined, arr_enc.buf, arr_enc.size);
    enc_write_bytes(&combined, vals.buf, vals.size);
    free(arr_enc.buf);
    free(vals.buf);

    encoder_t container_enc = {0};
    enc_encode_container(&container_enc, MM_CONTAINER_OBJECT, combined.buf, combined.size);
    free(combined.buf);

    enc_encode_tag(e, &obj->tag, container_enc.buf, container_enc.size);
    free(container_enc.buf);
}

static void enc_encode_node_doc(encoder_t* e, mm_doc_t* doc) {
    encoder_t keys = {0};
    encoder_t vals = {0};

    for (size_t i = 0; i < doc->field_count; i++) {
        mm_field_t* field = &doc->fields[i];

        encoder_t tmp_val = {0};
        switch (field->value->type) {
            case MM_NODE_OBJECT:
                enc_encode_node_object(&tmp_val, &field->value->data.object);
                break;
            case MM_NODE_ARRAY:
                enc_encode_node_array(&tmp_val, &field->value->data.array);
                break;
            case MM_NODE_VALUE:
                enc_encode_node_value(&tmp_val, &field->value->data.value);
                break;
            case MM_NODE_DOC:
                enc_encode_node_doc(&tmp_val, &field->value->data.doc);
                break;
            default:
                break;
        }
        enc_write_bytes(&vals, tmp_val.buf, tmp_val.size);
        free(tmp_val.buf);

        encoder_t tmp_key = {0};
        enc_encode_string(&tmp_key, field->key);
        enc_write_bytes(&keys, tmp_key.buf, tmp_key.size);
        free(tmp_key.buf);
    }

    encoder_t arr_enc = {0};
    enc_encode_array(&arr_enc, keys.buf, keys.size);
    free(keys.buf);

    encoder_t combined = {0};
    enc_write_bytes(&combined, arr_enc.buf, arr_enc.size);
    enc_write_bytes(&combined, vals.buf, vals.size);
    free(arr_enc.buf);
    free(vals.buf);

    encoder_t container_enc = {0};
    enc_encode_container(&container_enc, MM_CONTAINER_OBJECT, combined.buf, combined.size);
    free(combined.buf);

    enc_encode_tag(e, &doc->tag, container_enc.buf, container_enc.size);
    free(container_enc.buf);
}

static void enc_encode(encoder_t* e, mm_node_t* node) {
    if (node == NULL) return;

    switch (node->type) {
        case MM_NODE_OBJECT:
            enc_encode_node_object(e, &node->data.object);
            break;
        case MM_NODE_ARRAY:
            enc_encode_node_array(e, &node->data.array);
            break;
        case MM_NODE_VALUE:
            enc_encode_node_value(e, &node->data.value);
            break;
        case MM_NODE_DOC:
            enc_encode_node_doc(e, &node->data.doc);
            break;
        default:
            break;
    }
}

mm_encoder_buffer_t* mm_encoder_encode(mm_node_t* node) {
    encoder_t e = {0};
    enc_encode(&e, node);
    mm_encoder_buffer_t* buf = malloc(sizeof(mm_encoder_buffer_t));
    buf->data = e.buf;
    buf->size = e.size;
    buf->capacity = e.cap;
    return buf;
}

void mm_encoder_buffer_free(mm_encoder_buffer_t* buf) {
    if (buf) {
        free(buf->data);
        free(buf);
    }
}