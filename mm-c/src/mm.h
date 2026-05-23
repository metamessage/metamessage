#ifndef MM_H
#define MM_H

#include "ir/mc_ast.h"
#include "ir/mc_value_type.h"
#include <limits.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

typedef mm_node_t mm_obj_t;

typedef struct {
  uint8_t *data;
  size_t size;
  size_t capacity;
} mm_buffer_t;

typedef struct {
  const char *desc;
  const char *default_val;
  int64_t min;
  int64_t max;
  int size;
  const char *enums;
  const char *pattern;
  bool nullable;
  bool raw;
  bool allow_empty;
  bool unique;
  int location;
  int version;
  const char *mime;
} mm_field_attr_t;

typedef struct {
  const char *desc;
  const char *default_val;
  int64_t min;
  int64_t max;
  int size;
  const char *enums;
  const char *pattern;
  bool nullable;
  bool raw;
  bool allow_empty;
  bool unique;
  int location;
  int version;
  const char *mime;
  const char *child_desc;
  int64_t child_min;
  int64_t child_max;
  int child_size;
  bool child_nullable;
  bool child_raw;
  bool child_allow_empty;
  bool child_unique;
  const char *child_default_val;
  const char *child_enums;
  const char *child_pattern;
  int child_location;
  int child_version;
  const char *child_mime;
  const char *child_type;
} mm_container_attr_t;

mm_node_t *mm_value_create_str(const char *text, mm_value_type_t type,
                               mm_field_attr_t attr);
void mm_container_apply_attr(mm_node_t *container, mm_container_attr_t attr);

mm_node_t *mm_int_create(int64_t val, mm_field_attr_t attr);
mm_node_t *mm_i8_create(int64_t val, mm_field_attr_t attr);
mm_node_t *mm_i16_create(int64_t val, mm_field_attr_t attr);
mm_node_t *mm_i32_create(int64_t val, mm_field_attr_t attr);
mm_node_t *mm_i64_create(int64_t val, mm_field_attr_t attr);
mm_node_t *mm_uint_create(uint64_t val, mm_field_attr_t attr);
mm_node_t *mm_u8_create(uint64_t val, mm_field_attr_t attr);
mm_node_t *mm_u16_create(uint64_t val, mm_field_attr_t attr);
mm_node_t *mm_u32_create(uint64_t val, mm_field_attr_t attr);
mm_node_t *mm_u64_create(uint64_t val, mm_field_attr_t attr);
mm_node_t *mm_str_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_bool_create(bool val, mm_field_attr_t attr);
mm_node_t *mm_float_create(double val, mm_field_attr_t attr);
mm_node_t *mm_f32_create(double val, mm_field_attr_t attr);
mm_node_t *mm_bytes_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_bigint_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_datetime_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_date_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_time_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_uuid_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_decimal_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_ip_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_url_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_email_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_enum_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_image_create(const char *val, mm_field_attr_t attr);
mm_node_t *mm_video_create(const char *val, mm_field_attr_t attr);

#define mm_i(val, ...)                                                         \
  mm_int_create((int64_t)(val), (mm_field_attr_t){.min = INT64_MIN,            \
                                                  .max = INT64_MIN,            \
                                                  .size = -1,                  \
                                                  .version = -1,               \
                                                  .location = INT_MIN,         \
                                                  .desc = NULL,                \
                                                  .default_val = NULL,         \
                                                  .enums = NULL,            \
                                                  .pattern = NULL,             \
                                                  .nullable = false,           \
                                                  .raw = false,                \
                                                  .allow_empty = false,        \
                                                  .unique = false,             \
                                                  .mime = NULL,                \
                                                  __VA_ARGS__})
#define mm_str(val, ...)                                                       \
  mm_str_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,       \
                                                       .max = INT64_MIN,       \
                                                       .size = -1,             \
                                                       .version = -1,          \
                                                       .location = INT_MIN,    \
                                                       .desc = NULL,           \
                                                       .default_val = NULL,    \
                                                       .enums = NULL,       \
                                                       .pattern = NULL,        \
                                                       .nullable = false,      \
                                                       .raw = false,           \
                                                       .allow_empty = false,   \
                                                       .unique = false,        \
                                                       .mime = NULL,           \
                                                       __VA_ARGS__})
#define mm_bool(val, ...)                                                      \
  mm_bool_create((bool)(val), (mm_field_attr_t){.min = INT64_MIN,              \
                                                .max = INT64_MIN,              \
                                                .size = -1,                    \
                                                .version = -1,                 \
                                                .location = INT_MIN,           \
                                                .desc = NULL,                  \
                                                .default_val = NULL,           \
                                                .enums = NULL,              \
                                                .pattern = NULL,               \
                                                .nullable = false,             \
                                                .raw = false,                  \
                                                .allow_empty = false,          \
                                                .unique = false,               \
                                                .mime = NULL,                  \
                                                __VA_ARGS__})
#define mm_f64(val, ...)                                                       \
  mm_float_create((double)(val), (mm_field_attr_t){.min = INT64_MIN,           \
                                                   .max = INT64_MIN,           \
                                                   .size = -1,                 \
                                                   .version = -1,              \
                                                   .location = INT_MIN,        \
                                                   .desc = NULL,               \
                                                   .default_val = NULL,        \
                                                   .enums = NULL,           \
                                                   .pattern = NULL,            \
                                                   .nullable = false,          \
                                                   .raw = false,               \
                                                   .allow_empty = false,       \
                                                   .unique = false,            \
                                                   .mime = NULL,               \
                                                   __VA_ARGS__})
#define mm_i8(val, ...)                                                        \
  mm_i8_create((int64_t)(val), (mm_field_attr_t){.min = INT64_MIN,             \
                                                 .max = INT64_MIN,             \
                                                 .size = -1,                   \
                                                 .version = -1,                \
                                                 .location = INT_MIN,          \
                                                 .desc = NULL,                 \
                                                 .default_val = NULL,          \
                                                 .enums = NULL,             \
                                                 .pattern = NULL,              \
                                                 .nullable = false,            \
                                                 .raw = false,                 \
                                                 .allow_empty = false,         \
                                                 .unique = false,              \
                                                 .mime = NULL,                 \
                                                 __VA_ARGS__})
#define mm_i16(val, ...)                                                       \
  mm_i16_create((int64_t)(val), (mm_field_attr_t){.min = INT64_MIN,            \
                                                  .max = INT64_MIN,            \
                                                  .size = -1,                  \
                                                  .version = -1,               \
                                                  .location = INT_MIN,         \
                                                  .desc = NULL,                \
                                                  .default_val = NULL,         \
                                                  .enums = NULL,            \
                                                  .pattern = NULL,             \
                                                  .nullable = false,           \
                                                  .raw = false,                \
                                                  .allow_empty = false,        \
                                                  .unique = false,             \
                                                  .mime = NULL,                \
                                                  __VA_ARGS__})
#define mm_i32(val, ...)                                                       \
  mm_i32_create((int64_t)(val), (mm_field_attr_t){.min = INT64_MIN,            \
                                                  .max = INT64_MIN,            \
                                                  .size = -1,                  \
                                                  .version = -1,               \
                                                  .location = INT_MIN,         \
                                                  .desc = NULL,                \
                                                  .default_val = NULL,         \
                                                  .enums = NULL,            \
                                                  .pattern = NULL,             \
                                                  .nullable = false,           \
                                                  .raw = false,                \
                                                  .allow_empty = false,        \
                                                  .unique = false,             \
                                                  .mime = NULL,                \
                                                  __VA_ARGS__})
#define mm_i64(val, ...)                                                       \
  mm_i64_create((int64_t)(val), (mm_field_attr_t){.min = INT64_MIN,            \
                                                  .max = INT64_MIN,            \
                                                  .size = -1,                  \
                                                  .version = -1,               \
                                                  .location = INT_MIN,         \
                                                  .desc = NULL,                \
                                                  .default_val = NULL,         \
                                                  .enums = NULL,            \
                                                  .pattern = NULL,             \
                                                  .nullable = false,           \
                                                  .raw = false,                \
                                                  .allow_empty = false,        \
                                                  .unique = false,             \
                                                  .mime = NULL,                \
                                                  __VA_ARGS__})
#define mm_u(val, ...)                                                         \
  mm_uint_create((uint64_t)(val), (mm_field_attr_t){.min = INT64_MIN,          \
                                                    .max = INT64_MIN,          \
                                                    .size = -1,                \
                                                    .version = -1,             \
                                                    .location = INT_MIN,       \
                                                    .desc = NULL,              \
                                                    .default_val = NULL,       \
                                                    .enums = NULL,          \
                                                    .pattern = NULL,           \
                                                    .nullable = false,         \
                                                    .raw = false,              \
                                                    .allow_empty = false,      \
                                                    .unique = false,           \
                                                    .mime = NULL,              \
                                                    __VA_ARGS__})
#define mm_u8(val, ...)                                                        \
  mm_u8_create((uint64_t)(val), (mm_field_attr_t){.min = INT64_MIN,            \
                                                  .max = INT64_MIN,            \
                                                  .size = -1,                  \
                                                  .version = -1,               \
                                                  .location = INT_MIN,         \
                                                  .desc = NULL,                \
                                                  .default_val = NULL,         \
                                                  .enums = NULL,            \
                                                  .pattern = NULL,             \
                                                  .nullable = false,           \
                                                  .raw = false,                \
                                                  .allow_empty = false,        \
                                                  .unique = false,             \
                                                  .mime = NULL,                \
                                                  __VA_ARGS__})
#define mm_u16(val, ...)                                                       \
  mm_u16_create((uint64_t)(val), (mm_field_attr_t){.min = INT64_MIN,           \
                                                   .max = INT64_MIN,           \
                                                   .size = -1,                 \
                                                   .version = -1,              \
                                                   .location = INT_MIN,        \
                                                   .desc = NULL,               \
                                                   .default_val = NULL,        \
                                                   .enums = NULL,           \
                                                   .pattern = NULL,            \
                                                   .nullable = false,          \
                                                   .raw = false,               \
                                                   .allow_empty = false,       \
                                                   .unique = false,            \
                                                   .mime = NULL,               \
                                                   __VA_ARGS__})
#define mm_u32(val, ...)                                                       \
  mm_u32_create((uint64_t)(val), (mm_field_attr_t){.min = INT64_MIN,           \
                                                   .max = INT64_MIN,           \
                                                   .size = -1,                 \
                                                   .version = -1,              \
                                                   .location = INT_MIN,        \
                                                   .desc = NULL,               \
                                                   .default_val = NULL,        \
                                                   .enums = NULL,           \
                                                   .pattern = NULL,            \
                                                   .nullable = false,          \
                                                   .raw = false,               \
                                                   .allow_empty = false,       \
                                                   .unique = false,            \
                                                   .mime = NULL,               \
                                                   __VA_ARGS__})
#define mm_u64(val, ...)                                                       \
  mm_u64_create((uint64_t)(val), (mm_field_attr_t){.min = INT64_MIN,           \
                                                   .max = INT64_MIN,           \
                                                   .size = -1,                 \
                                                   .version = -1,              \
                                                   .location = INT_MIN,        \
                                                   .desc = NULL,               \
                                                   .default_val = NULL,        \
                                                   .enums = NULL,           \
                                                   .pattern = NULL,            \
                                                   .nullable = false,          \
                                                   .raw = false,               \
                                                   .allow_empty = false,       \
                                                   .unique = false,            \
                                                   .mime = NULL,               \
                                                   __VA_ARGS__})
#define mm_f32(val, ...)                                                       \
  mm_f32_create((double)(val), (mm_field_attr_t){.min = INT64_MIN,             \
                                                 .max = INT64_MIN,             \
                                                 .size = -1,                   \
                                                 .version = -1,                \
                                                 .location = INT_MIN,          \
                                                 .desc = NULL,                 \
                                                 .default_val = NULL,          \
                                                 .enums = NULL,             \
                                                 .pattern = NULL,              \
                                                 .nullable = false,            \
                                                 .raw = false,                 \
                                                 .allow_empty = false,         \
                                                 .unique = false,              \
                                                 .mime = NULL,                 \
                                                 __VA_ARGS__})
#define mm_bytes(val, ...)                                                     \
  mm_bytes_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,     \
                                                         .max = INT64_MIN,     \
                                                         .size = -1,           \
                                                         .version = -1,        \
                                                         .location = INT_MIN,  \
                                                         .desc = NULL,         \
                                                         .default_val = NULL,  \
                                                         .enums = NULL,     \
                                                         .pattern = NULL,      \
                                                         .nullable = false,    \
                                                         .raw = false,         \
                                                         .allow_empty = false, \
                                                         .unique = false,      \
                                                         .mime = NULL,         \
                                                         __VA_ARGS__})
#define mm_bigint(val, ...)                                                    \
  mm_bigint_create((const char *)(val),                                        \
                   (mm_field_attr_t){.min = INT64_MIN,                         \
                                     .max = INT64_MIN,                         \
                                     .size = -1,                               \
                                     .version = -1,                            \
                                     .location = INT_MIN,                      \
                                     .desc = NULL,                             \
                                     .default_val = NULL,                      \
                                     .enums = NULL,                         \
                                     .pattern = NULL,                          \
                                     .nullable = false,                        \
                                     .raw = false,                             \
                                     .allow_empty = false,                     \
                                     .unique = false,                          \
                                     .mime = NULL,                             \
                                     __VA_ARGS__})
#define mm_datetime(val, ...)                                                  \
  mm_datetime_create((const char *)(val),                                      \
                     (mm_field_attr_t){.min = INT64_MIN,                       \
                                       .max = INT64_MIN,                       \
                                       .size = -1,                             \
                                       .version = -1,                          \
                                       .location = INT_MIN,                    \
                                       .desc = NULL,                           \
                                       .default_val = NULL,                    \
                                       .enums = NULL,                       \
                                       .pattern = NULL,                        \
                                       .nullable = false,                      \
                                       .raw = false,                           \
                                       .allow_empty = false,                   \
                                       .unique = false,                        \
                                       .mime = NULL,                           \
                                       __VA_ARGS__})
#define mm_date(val, ...)                                                      \
  mm_date_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,      \
                                                        .max = INT64_MIN,      \
                                                        .size = -1,            \
                                                        .version = -1,         \
                                                        .location = INT_MIN,   \
                                                        .desc = NULL,          \
                                                        .default_val = NULL,   \
                                                        .enums = NULL,      \
                                                        .pattern = NULL,       \
                                                        .nullable = false,     \
                                                        .raw = false,          \
                                                        .allow_empty = false,  \
                                                        .unique = false,       \
                                                        .mime = NULL,          \
                                                        __VA_ARGS__})
#define mm_time(val, ...)                                                      \
  mm_time_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,      \
                                                        .max = INT64_MIN,      \
                                                        .size = -1,            \
                                                        .version = -1,         \
                                                        .location = INT_MIN,   \
                                                        .desc = NULL,          \
                                                        .default_val = NULL,   \
                                                        .enums = NULL,      \
                                                        .pattern = NULL,       \
                                                        .nullable = false,     \
                                                        .raw = false,          \
                                                        .allow_empty = false,  \
                                                        .unique = false,       \
                                                        .mime = NULL,          \
                                                        __VA_ARGS__})
#define mm_uuid(val, ...)                                                      \
  mm_uuid_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,      \
                                                        .max = INT64_MIN,      \
                                                        .size = -1,            \
                                                        .version = -1,         \
                                                        .location = INT_MIN,   \
                                                        .desc = NULL,          \
                                                        .default_val = NULL,   \
                                                        .enums = NULL,      \
                                                        .pattern = NULL,       \
                                                        .nullable = false,     \
                                                        .raw = false,          \
                                                        .allow_empty = false,  \
                                                        .unique = false,       \
                                                        .mime = NULL,          \
                                                        __VA_ARGS__})
#define mm_decimal(val, ...)                                                   \
  mm_decimal_create((const char *)(val),                                       \
                    (mm_field_attr_t){.min = INT64_MIN,                        \
                                      .max = INT64_MIN,                        \
                                      .size = -1,                              \
                                      .version = -1,                           \
                                      .location = INT_MIN,                     \
                                      .desc = NULL,                            \
                                      .default_val = NULL,                     \
                                      .enums = NULL,                        \
                                      .pattern = NULL,                         \
                                      .nullable = false,                       \
                                      .raw = false,                            \
                                      .allow_empty = false,                    \
                                      .unique = false,                         \
                                      .mime = NULL,                            \
                                      __VA_ARGS__})
#define mm_ip(val, ...)                                                        \
  mm_ip_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,        \
                                                      .max = INT64_MIN,        \
                                                      .size = -1,              \
                                                      .version = -1,           \
                                                      .location = INT_MIN,     \
                                                      .desc = NULL,            \
                                                      .default_val = NULL,     \
                                                      .enums = NULL,        \
                                                      .pattern = NULL,         \
                                                      .nullable = false,       \
                                                      .raw = false,            \
                                                      .allow_empty = false,    \
                                                      .unique = false,         \
                                                      .mime = NULL,            \
                                                      __VA_ARGS__})
#define mm_url(val, ...)                                                       \
  mm_url_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,       \
                                                       .max = INT64_MIN,       \
                                                       .size = -1,             \
                                                       .version = -1,          \
                                                       .location = INT_MIN,    \
                                                       .desc = NULL,           \
                                                       .default_val = NULL,    \
                                                       .enums = NULL,       \
                                                       .pattern = NULL,        \
                                                       .nullable = false,      \
                                                       .raw = false,           \
                                                       .allow_empty = false,   \
                                                       .unique = false,        \
                                                       .mime = NULL,           \
                                                       __VA_ARGS__})
#define mm_email(val, ...)                                                     \
  mm_email_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,     \
                                                         .max = INT64_MIN,     \
                                                         .size = -1,           \
                                                         .version = -1,        \
                                                         .location = INT_MIN,  \
                                                         .desc = NULL,         \
                                                         .default_val = NULL,  \
                                                         .enums = NULL,     \
                                                         .pattern = NULL,      \
                                                         .nullable = false,    \
                                                         .raw = false,         \
                                                         .allow_empty = false, \
                                                         .unique = false,      \
                                                         .mime = NULL,         \
                                                         __VA_ARGS__})
#define mm_enum(val, ...)                                                      \
  mm_enum_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,      \
                                                        .max = INT64_MIN,      \
                                                        .size = -1,            \
                                                        .version = -1,         \
                                                        .location = INT_MIN,   \
                                                        .desc = NULL,          \
                                                        .default_val = NULL,   \
                                                        .enums = NULL,      \
                                                        .pattern = NULL,       \
                                                        .nullable = false,     \
                                                        .raw = false,          \
                                                        .allow_empty = false,  \
                                                        .unique = false,       \
                                                        .mime = NULL,          \
                                                        __VA_ARGS__})
#define mm_image(val, ...)                                                     \
  mm_image_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,     \
                                                         .max = INT64_MIN,     \
                                                         .size = -1,           \
                                                         .version = -1,        \
                                                         .location = INT_MIN,  \
                                                         .desc = NULL,         \
                                                         .default_val = NULL,  \
                                                         .enums = NULL,     \
                                                         .pattern = NULL,      \
                                                         .nullable = false,    \
                                                         .raw = false,         \
                                                         .allow_empty = false, \
                                                         .unique = false,      \
                                                         .mime = NULL,         \
                                                         __VA_ARGS__})
#define mm_video(val, ...)                                                     \
  mm_video_create((const char *)(val), (mm_field_attr_t){.min = INT64_MIN,     \
                                                         .max = INT64_MIN,     \
                                                         .size = -1,           \
                                                         .version = -1,        \
                                                         .location = INT_MIN,  \
                                                         .desc = NULL,         \
                                                         .default_val = NULL,  \
                                                         .enums = NULL,     \
                                                         .pattern = NULL,      \
                                                         .nullable = false,    \
                                                         .raw = false,         \
                                                         .allow_empty = false, \
                                                         .unique = false,      \
                                                         .mime = NULL,         \
                                                         __VA_ARGS__})

mm_obj_t *mm_obj_new(void);
void mm_obj_set(mm_obj_t *obj, const char *key, mm_node_t *value);
void mm_obj_free(mm_obj_t *obj);

mm_node_t *mm_arr_new(void);
void mm_arr_add(mm_node_t *arr, mm_node_t *item);

mm_buffer_t *mm_encode(mm_node_t *node);
mm_node_t *mm_decode(const mm_buffer_t *buf);
void mm_buffer_free(mm_buffer_t *buf);

char *mm_to_jsonc(mm_node_t *node);
mm_node_t *mm_from_jsonc(const char *jsonc_str);
void mm_string_free(char *str);

#endif