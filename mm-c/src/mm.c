#include "mm.h"
#include "core/mc_decoder.h"
#include "core/mc_encoder.h"
#include "ir/mc_tag.h"
#include "jsonc/mc_ir_parser.h"
#include "jsonc/mc_printer.h"
#include <limits.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

node_t *mm_value_create_str(const char *text, mm_value_type_t type,
                               mm_field_attr_t attr) {
  node_t *node = node_new_scalar();
  if (!node)
    return NULL;

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
  node->data.value.tag.deprecated = attr.deprecated;
  node->data.value.tag.allow_empty = attr.allow_empty;
  node->data.value.tag.unique = attr.unique;

  if (attr.version > 0) {
    node->data.value.tag.version = attr.version;
  }

  if (attr.mime && strlen(attr.mime) > 0) {
    node->data.value.tag.mime = strdup(attr.mime);
  }

  return node;
}

void mm_container_apply_attr(node_t *container, mm_container_attr_t attr) {
  if (!container)
    return;

  mm_tag_t *tag = NULL;
  if (container->type == MM_NODE_ARRAY) {
    tag = &container->data.array.tag;
  } else if (container->type == MM_NODE_OBJECT) {
    // obj does not have child_* attrs, only apply base fields
    tag = &container->data.object.tag;
  }

  if (!tag)
    return;

  if (attr.desc && strlen(attr.desc) > 0) {
    tag->desc = strdup(attr.desc);
  }
  if (attr.default_val && strlen(attr.default_val) > 0) {
    tag->default_val = strdup(attr.default_val);
  }
  if (attr.min != INT64_MIN) {
    char buf[32];
    snprintf(buf, sizeof(buf), "%lld", (long long)attr.min);
    tag->min = strdup(buf);
  }
  if (attr.max != INT64_MIN) {
    char buf[32];
    snprintf(buf, sizeof(buf), "%lld", (long long)attr.max);
    tag->max = strdup(buf);
  }
  if (attr.size > 0) {
    tag->size = attr.size;
  }
  if (attr.enums && strlen(attr.enums) > 0) {
    tag->enums = strdup(attr.enums);
  }
  if (attr.pattern && strlen(attr.pattern) > 0) {
    tag->pattern = strdup(attr.pattern);
  }
  if (attr.location != INT_MIN) {
    tag->location_offset = attr.location;
  }
  tag->nullable = attr.nullable;
  tag->deprecated = attr.deprecated;
  tag->allow_empty = attr.allow_empty;
  tag->unique = attr.unique;
  if (attr.version > 0) {
    tag->version = attr.version;
  }
  if (attr.mime && strlen(attr.mime) > 0) {
    tag->mime = strdup(attr.mime);
  }

  // child_* only for arr/vec/map (not obj)
  if (container->type == MM_NODE_OBJECT)
    return;

  if (attr.child_desc && strlen(attr.child_desc) > 0) {
    tag->child_desc = strdup(attr.child_desc);
  }
  if (attr.child_min != INT64_MIN) {
    char buf[32];
    snprintf(buf, sizeof(buf), "%lld", (long long)attr.child_min);
    tag->child_min = strdup(buf);
  }
  if (attr.child_max != INT64_MIN) {
    char buf[32];
    snprintf(buf, sizeof(buf), "%lld", (long long)attr.child_max);
    tag->child_max = strdup(buf);
  }
  if (attr.child_size > 0) {
    tag->child_size = attr.child_size;
  }
  tag->child_nullable = attr.child_nullable;
  tag->child_allow_empty = attr.child_allow_empty;
  tag->child_unique = attr.child_unique;
  if (attr.child_version > 0) {
    tag->child_version = attr.child_version;
  }
  if (attr.child_mime && strlen(attr.child_mime) > 0) {
    tag->child_mime = strdup(attr.child_mime);
  }
  if (attr.child_type && strlen(attr.child_type) > 0) {
    tag->child_type = mm_value_type_parse(attr.child_type);
  }
  if (attr.child_default_val && strlen(attr.child_default_val) > 0) {
    tag->child_default_val = strdup(attr.child_default_val);
  }
  if (attr.child_enums && strlen(attr.child_enums) > 0) {
    tag->child_enums = strdup(attr.child_enums);
  }
  if (attr.child_pattern && strlen(attr.child_pattern) > 0) {
    tag->child_pattern = strdup(attr.child_pattern);
  }
  if (attr.child_location != INT_MIN) {
    tag->child_location_offset = attr.child_location;
  }
}

node_t *mm_int_create(int64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%lld", (long long)val);
  return mm_value_create_str(buf, MM_VALUE_I, attr);
}

node_t *mm_str_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_STR, attr);
}

node_t *mm_bool_create(bool val, mm_field_attr_t attr) {
  const char *text = val ? "true" : "false";
  return mm_value_create_str(text, MM_VALUE_BOOL, attr);
}

node_t *mm_float_create(double val, mm_field_attr_t attr) {
  char buf[64];
  snprintf(buf, sizeof(buf), "%g", val);
  return mm_value_create_str(buf, MM_VALUE_F64, attr);
}

node_t *mm_i8_create(int64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%lld", (long long)val);
  return mm_value_create_str(buf, MM_VALUE_I8, attr);
}

node_t *mm_i16_create(int64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%lld", (long long)val);
  return mm_value_create_str(buf, MM_VALUE_I16, attr);
}

node_t *mm_i32_create(int64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%lld", (long long)val);
  return mm_value_create_str(buf, MM_VALUE_I32, attr);
}

node_t *mm_i64_create(int64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%lld", (long long)val);
  return mm_value_create_str(buf, MM_VALUE_I64, attr);
}

node_t *mm_uint_create(uint64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%llu", (unsigned long long)val);
  return mm_value_create_str(buf, MM_VALUE_U, attr);
}

node_t *mm_u8_create(uint64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%llu", (unsigned long long)val);
  return mm_value_create_str(buf, MM_VALUE_U8, attr);
}

node_t *mm_u16_create(uint64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%llu", (unsigned long long)val);
  return mm_value_create_str(buf, MM_VALUE_U16, attr);
}

node_t *mm_u32_create(uint64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%llu", (unsigned long long)val);
  return mm_value_create_str(buf, MM_VALUE_U32, attr);
}

node_t *mm_u64_create(uint64_t val, mm_field_attr_t attr) {
  char buf[32];
  snprintf(buf, sizeof(buf), "%llu", (unsigned long long)val);
  return mm_value_create_str(buf, MM_VALUE_U64, attr);
}

node_t *mm_f32_create(double val, mm_field_attr_t attr) {
  char buf[64];
  snprintf(buf, sizeof(buf), "%g", val);
  return mm_value_create_str(buf, MM_VALUE_F32, attr);
}

node_t *mm_bytes_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_BYTES, attr);
}

node_t *mm_bigint_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_BIGINT, attr);
}

node_t *mm_datetime_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_DATETIME, attr);
}

node_t *mm_date_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_DATE, attr);
}

node_t *mm_time_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_TIME, attr);
}

node_t *mm_uuid_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_UUID, attr);
}

node_t *mm_decimal_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_DECIMAL, attr);
}

node_t *mm_ip_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_IP, attr);
}

node_t *mm_url_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_URL, attr);
}

node_t *mm_email_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_EMAIL, attr);
}

node_t *mm_enum_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_ENUMS, attr);
}

node_t *mm_image_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_IMAGE, attr);
}

node_t *mm_video_create(const char *val, mm_field_attr_t attr) {
  return mm_value_create_str(val, MM_VALUE_VIDEO, attr);
}

mm_obj_t *mm_obj_new(void) { return node_new_object(); }

void mm_obj_set(mm_obj_t *obj, const char *key, node_t *value) {
  node_object_add_field(obj, key, value);
}

void mm_obj_free(mm_obj_t *obj) { node_free(obj); }

node_t *mm_arr_new(void) { return node_new_array(); }

void mm_arr_add(node_t *arr, node_t *item) {
  node_array_add_item(arr, item);
}

mm_buffer_t *mm_encode(node_t *node) {
  mm_encoder_buffer_t *enc = mm_encoder_encode(node);
  if (!enc)
    return NULL;

  mm_buffer_t *buf = (mm_buffer_t *)malloc(sizeof(mm_buffer_t));
  if (!buf) {
    mm_encoder_buffer_free(enc);
    return NULL;
  }

  buf->data = (uint8_t *)malloc(enc->capacity);
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

node_t *mm_decode(const mm_buffer_t *buf) {
  mm_decoder_t *d = mm_decoder_new(buf->data, buf->size);
  if (!d)
    return NULL;

  node_t *node = mm_decoder_decode(d);
  mm_decoder_free(d);

  return node;
}

char *mm_to_jsonc(node_t *node) { return mm_printer_to_jsonc(node); }

node_t *mm_from_jsonc(const char *jsonc_str) {
  return mm_jsonc_parse(jsonc_str);
}

void mm_buffer_free(mm_buffer_t *buf) {
  if (buf) {
    free(buf->data);
    free(buf);
  }
}

void mm_string_free(char *str) { free(str); }

mm_buffer_t *mm_encode_from_jsonc(const char *jsonc_str) {
  node_t *node = mm_from_jsonc(jsonc_str);
  if (!node)
    return NULL;

  mm_buffer_t *buf = mm_encode(node);
  node_free(node);
  return buf;
}

char *mm_decode_to_jsonc(const mm_buffer_t *buf) {
  node_t *node = mm_decode(buf);
  if (!node)
    return NULL;

  char *jsonc = mm_to_jsonc(node);
  node_free(node);
  return jsonc;
}