#ifndef MMC_IR_AST_H
#define MMC_IR_AST_H

#include "mc_tag.h"
#include <stddef.h>

typedef enum {
    MM_NODE_UNKNOWN = 0,
    MM_NODE_OBJECT,
    MM_NODE_ARRAY,
    MM_NODE_VALUE,
    MM_NODE_DOC
} mm_node_type_t;

typedef struct mm_node mm_node_t;
typedef struct mm_field mm_field_t;
typedef struct mm_object mm_object_t;
typedef struct mm_array mm_array_t;
typedef struct mm_value mm_value_t;
typedef struct mm_doc mm_doc_t;

struct mm_field {
    char* key;
    mm_node_t* value;
};

struct mm_object {
    mm_field_t* fields;
    size_t field_count;
    size_t capacity;
    mm_tag_t tag;
};

struct mm_array {
    mm_node_t** items;
    size_t item_count;
    size_t capacity;
    mm_tag_t tag;
};

struct mm_value {
    char* text;
    mm_tag_t tag;
};

struct mm_doc {
    mm_field_t* fields;
    size_t field_count;
    size_t capacity;
    mm_tag_t tag;
};

struct mm_node {
    mm_node_type_t type;
    mm_tag_t tag;
    char* path;
    union {
        mm_object_t object;
        mm_array_t array;
        mm_value_t value;
        mm_doc_t doc;
    } data;
};

mm_node_t* mm_node_new_object(void);
mm_node_t* mm_node_new_array(void);
mm_node_t* mm_node_new_value(void);
mm_node_t* mm_node_new_doc(void);
void mm_node_free(mm_node_t* node);

void mm_object_add_field(mm_node_t* obj, const char* key, mm_node_t* value);
void mm_array_add_item(mm_node_t* arr, mm_node_t* item);
void mm_doc_add_field(mm_node_t* doc, const char* key, mm_node_t* value);

#endif