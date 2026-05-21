#include "mm.h"
#include "core/mc_encoder.h"
#include "core/mc_decoder.h"
#include "jsonc/mc_printer.h"
#include "jsonc/mc_ir_parser.h"
#include "ir/mc_tag.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <limits.h>

mm_node_t* mm_value_create_str(const char* text, mm_value_type_t type, mm_field_attr_t attr)
{
    mm_node_t* node = mm_node_new_value();
    if (!node) return NULL;

    node->data.value.text = strdup(text);
    node->data.value.tag.type = type;

    if (attr.desc && strlen(attr.desc) > 0) {
        node->data.value.tag.desc = strdup(attr.desc);
    }

    if (attr.default_val && strlen(attr.default_val) > 0) {
        node->data.value.tag.default_val = strdup(attr.default_val);
    }

    if (attr.min != INT64_MIN) {
        char buf[32];
        snprintf(buf, sizeof(buf), "%lld", (long long)attr.min);
        node->data.value.tag.min = strdup(buf);
    }

    if (attr.max != INT64_MIN) {
        char buf[32];
        snprintf(buf, sizeof(buf), "%lld", (long long)attr.max);
        node->data.value.tag.max = strdup(buf);
    }

    if (attr.size > 0) {
        node->data.value.tag.size = attr.size;
    }

    if (attr.enums && strlen(attr.enums) > 0) {
        node->data.value.tag.enums = strdup(attr.enums);
    }

    if (attr.pattern && strlen(attr.pattern) > 0) {
        node->data.value.tag.pattern = strdup(attr.pattern);
    }

    if (attr.location != INT_MIN) {
        node->data.value.tag.location_offset = attr.location;
    }

    node->data.value.tag.nullable = attr.nullable;
    node->data.value.tag.raw = attr.raw;
    node->data.value.tag.allow_empty = attr.allow_empty;
    node->data.value.tag.unique = attr.unique;

    if (attr.version > 0) {
        node->data.value.tag.version = attr.version;
    }

    if (attr.mime && strlen(attr.mime) > 0) {
        node->data.value.tag.mime = strdup(attr.mime);
    }

    if (attr.child_desc && strlen(attr.child_desc) > 0) {
        node->data.value.tag.child_desc = strdup(attr.child_desc);
    }

    if (attr.child_min != INT64_MIN) {
        char buf[32];
        snprintf(buf, sizeof(buf), "%lld", (long long)attr.child_min);
        node->data.value.tag.child_min = strdup(buf);
    }

    if (attr.child_max != INT64_MIN) {
        char buf[32];
        snprintf(buf, sizeof(buf), "%lld", (long long)attr.child_max);
        node->data.value.tag.child_max = strdup(buf);
    }

    if (attr.child_size > 0) {
        node->data.value.tag.child_size = attr.child_size;
    }

    node->data.value.tag.child_nullable = attr.child_nullable;
    node->data.value.tag.child_raw = attr.child_raw;
    node->data.value.tag.child_allow_empty = attr.child_allow_empty;
    node->data.value.tag.child_unique = attr.child_unique;

    if (attr.child_version > 0) {
        node->data.value.tag.child_version = attr.child_version;
    }

    if (attr.child_mime && strlen(attr.child_mime) > 0) {
        node->data.value.tag.child_mime = strdup(attr.child_mime);
    }

    if (attr.child_type && strlen(attr.child_type) > 0) {
        node->data.value.tag.child_type = mm_value_type_parse(attr.child_type);
    }

    if (attr.child_default && strlen(attr.child_default) > 0) {
        node->data.value.tag.child_default = strdup(attr.child_default);
    }

    if (attr.child_enum && strlen(attr.child_enum) > 0) {
        node->data.value.tag.child_enum = strdup(attr.child_enum);
    }

    if (attr.child_pattern && strlen(attr.child_pattern) > 0) {
        node->data.value.tag.child_pattern = strdup(attr.child_pattern);
    }

    if (attr.child_location != INT_MIN) {
        node->data.value.tag.child_location_offset = attr.child_location;
    }

    return node;
}

mm_node_t* mm_int_create(int64_t val, mm_field_attr_t attr)
{
    char buf[32];
    snprintf(buf, sizeof(buf), "%lld", (long long)val);
    return mm_value_create_str(buf, MM_VALUE_I, attr);
}

mm_node_t* mm_str_create(const char* val, mm_field_attr_t attr)
{
    return mm_value_create_str(val, MM_VALUE_STR, attr);
}

mm_node_t* mm_bool_create(bool val, mm_field_attr_t attr)
{
    const char* text = val ? "true" : "false";
    return mm_value_create_str(text, MM_VALUE_BOOL, attr);
}

mm_node_t* mm_float_create(double val, mm_field_attr_t attr)
{
    char buf[64];
    snprintf(buf, sizeof(buf), "%g", val);
    return mm_value_create_str(buf, MM_VALUE_F64, attr);
}

mm_obj_t* mm_obj_new(void)
{
    return mm_node_new_object();
}

void mm_obj_set(mm_obj_t* obj, const char* key, mm_node_t* value)
{
    mm_object_add_field(obj, key, value);
}

void mm_obj_free(mm_obj_t* obj)
{
    mm_node_free(obj);
}

mm_node_t* mm_arr_new(void)
{
    return mm_node_new_array();
}

void mm_arr_add(mm_node_t* arr, mm_node_t* item)
{
    mm_array_add_item(arr, item);
}

mm_buffer_t* mm_encode(mm_node_t* node)
{
    mm_encoder_buffer_t* enc = mm_encoder_encode(node);
    if (!enc) return NULL;

    mm_buffer_t* buf = (mm_buffer_t*)malloc(sizeof(mm_buffer_t));
    if (!buf) {
        mm_encoder_buffer_free(enc);
        return NULL;
    }

    buf->data = (uint8_t*)malloc(enc->capacity);
    if (!buf->data) {
        free(buf);
        mm_encoder_buffer_free(enc);
        return NULL;
    }

    memcpy(buf->data, enc->data, enc->size);
    buf->size = enc->size;
    buf->capacity = enc->capacity;

    mm_encoder_buffer_free(enc);

    return buf;
}

mm_node_t* mm_decode(const mm_buffer_t* buf)
{
    mm_decoder_t* d = mm_decoder_new(buf->data, buf->size);
    if (!d) return NULL;

    mm_node_t* node = mm_decoder_decode(d);
    mm_decoder_free(d);

    return node;
}

char* mm_to_jsonc(mm_node_t* node)
{
    return mm_printer_to_jsonc(node);
}

mm_node_t* mm_from_jsonc(const char* jsonc_str)
{
    return mm_jsonc_parse(jsonc_str);
}

void mm_buffer_free(mm_buffer_t* buf)
{
    if (buf) {
        free(buf->data);
        free(buf);
    }
}

void mm_string_free(char* str)
{
    free(str);
}