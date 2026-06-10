#ifndef MMC_IR_AST_H
#define MMC_IR_AST_H

#include "mc_tag.h"
#include <stddef.h>

typedef enum {
  MM_NODE_UNKNOWN = 0,
  MM_NODE_OBJECT,
  MM_NODE_ARRAY,
  MM_NODE_VALUE,
  MM_NODE_NULL,
  MM_NODE_DOC
} node_type_t;

typedef struct node node_t;
typedef struct mm_field mm_field_t;
typedef struct node_object node_object_t;
typedef struct node_array node_array_t;
typedef struct node_scalar node_scalar_t;
typedef struct node_doc node_doc_t;

struct mm_field {
  char *key;
  node_t *value;
};

struct node_object {
  mm_field_t *fields;
  size_t field_count;
  size_t capacity;
  mm_tag_t tag;
};

struct node_array {
  node_t **items;
  size_t item_count;
  size_t capacity;
  mm_tag_t tag;
};

struct node_scalar {
  char *text;
  mm_tag_t tag;
};

struct node_doc {
  mm_field_t *fields;
  size_t field_count;
  size_t capacity;
  mm_tag_t tag;
};

struct node {
  node_type_t type;
  mm_tag_t tag;
  char *path;
  union {
    node_object_t object;
    node_array_t array;
    node_scalar_t value;
    node_doc_t doc;
  } data;
};

node_t *node_new_object(void);
node_t *node_new_array(void);
node_t *node_new_scalar(void);
node_t *node_new_null(void);
node_t *node_new_doc(void);
void node_free(node_t *node);

void node_object_add_field(node_t *obj, const char *key, node_t *value);
void node_array_add_item(node_t *arr, node_t *item);
void node_doc_add_field(node_t *doc, const char *key, node_t *value);

#endif