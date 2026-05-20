#ifndef MM_H
#define MM_H

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>
#include <limits.h>
#include "ir/mc_ast.h"
#include "ir/mc_value_type.h"

typedef mm_node_t mm_obj_t;

typedef struct {
    uint8_t* data;
    size_t size;
    size_t capacity;
} mm_buffer_t;

typedef struct {
    const char* desc;
    const char* default_val;
    int64_t min;
    int64_t max;
    int size;
    const char* enum_val;
    const char* pattern;
    bool nullable;
    bool raw;
    bool allow_empty;
    bool unique;
    int version;
    const char* mime;
    const char* child_desc;
    int64_t child_min;
    int64_t child_max;
    int child_size;
    bool child_nullable;
    bool child_raw;
    bool child_allow_empty;
    bool child_unique;
    int child_version;
    const char* child_mime;
} mm_field_attr_t;

mm_node_t* mm_value_create_str(const char* text, mm_value_type_t type, mm_field_attr_t attr);

mm_node_t* mm_int_create(int64_t val, mm_field_attr_t attr);
mm_node_t* mm_str_create(const char* val, mm_field_attr_t attr);
mm_node_t* mm_bool_create(bool val, mm_field_attr_t attr);
mm_node_t* mm_float_create(double val, mm_field_attr_t attr);

#define mm_int(val, ...) mm_int_create((int64_t)(val), (mm_field_attr_t){ .min = INT64_MIN, .max = INT64_MIN, .size = -1, .version = -1, __VA_ARGS__ })
#define mm_str(val, ...) mm_str_create((const char*)(val), (mm_field_attr_t){ .min = INT64_MIN, .max = INT64_MIN, .size = -1, .version = -1, __VA_ARGS__ })
#define mm_bool(val, ...) mm_bool_create((bool)(val), (mm_field_attr_t){ .min = INT64_MIN, .max = INT64_MIN, .size = -1, .version = -1, __VA_ARGS__ })
#define mm_float(val, ...) mm_float_create((double)(val), (mm_field_attr_t){ .min = INT64_MIN, .max = INT64_MIN, .size = -1, .version = -1, __VA_ARGS__ })

mm_obj_t* mm_obj_new(void);
void mm_obj_set(mm_obj_t* obj, const char* key, mm_node_t* value);
void mm_obj_free(mm_obj_t* obj);

mm_node_t* mm_arr_new(void);
void mm_arr_add(mm_node_t* arr, mm_node_t* item);

mm_buffer_t* mm_encode(mm_node_t* node);
mm_node_t* mm_decode(const mm_buffer_t* buf);
void mm_buffer_free(mm_buffer_t* buf);

char* mm_to_jsonc(mm_node_t* node);
mm_node_t* mm_from_jsonc(const char* jsonc_str);
void mm_string_free(char* str);

#endif