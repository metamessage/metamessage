#ifndef MMC_JSONC_AST_H
#define MMC_JSONC_AST_H

#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

typedef enum {
    MMC_VALUE_TYPE_UNKNOWN,
    MMC_VALUE_TYPE_STRING,
    MMC_VALUE_TYPE_BOOL,
    MMC_VALUE_TYPE_INT,
    MMC_VALUE_TYPE_INT8,
    MMC_VALUE_TYPE_INT16,
    MMC_VALUE_TYPE_INT32,
    MMC_VALUE_TYPE_INT64,
    MMC_VALUE_TYPE_UINT,
    MMC_VALUE_TYPE_UINT8,
    MMC_VALUE_TYPE_UINT16,
    MMC_VALUE_TYPE_UINT32,
    MMC_VALUE_TYPE_UINT64,
    MMC_VALUE_TYPE_FLOAT32,
    MMC_VALUE_TYPE_FLOAT64,
    MMC_VALUE_TYPE_BYTES,
    MMC_VALUE_TYPE_BIG_INT,
    MMC_VALUE_TYPE_ARRAY,
    MMC_VALUE_TYPE_SLICE,
    MMC_VALUE_TYPE_OBJECT
} MMC_ValueType;

typedef struct MMC_Node MMC_Node;
typedef struct MMC_Field MMC_Field;
typedef struct MMC_Object MMC_Object;
typedef struct MMC_Array MMC_Array;
typedef struct MMC_Value MMC_Value;

typedef struct MMC_Tag {
    MMC_ValueType type;
    int size;
    int required;
    int min_value;
    int max_value;
    int min_length;
    int max_length;
    int deprecated;
} MMC_Tag;

typedef enum {
    MMC_DATA_BOOL,
    MMC_DATA_STRING,
    MMC_DATA_INT,
    MMC_DATA_UINT,
    MMC_DATA_FLOAT,
    MMC_DATA_BYTES,
    MMC_DATA_NULL
} MMC_DataType;

typedef struct MMC_Value {
    MMC_DataType data_type;
    union {
        int bool_value;
        char* string_value;
        int64_t int_value;
        uint64_t uint_value;
        double float_value;
        struct {
            uint8_t* bytes;
            size_t len;
        } bytes_value;
    } data;
    char* text;
    MMC_Tag tag;
} MMC_Value;

struct MMC_Field {
    char* key;
    MMC_Node* value;
};

struct MMC_Object {
    MMC_Field* fields;
    size_t field_count;
    MMC_Tag tag;
};

struct MMC_Array {
    MMC_Node** items;
    size_t item_count;
    MMC_Tag tag;
};

enum MMC_NodeType {
    MMC_NODE_VALUE,
    MMC_NODE_OBJECT,
    MMC_NODE_ARRAY
};

struct MMC_Node {
    enum MMC_NodeType type;
    union {
        MMC_Value value;
        MMC_Object object;
        MMC_Array array;
    } data;
};

MMC_Node* mmc_node_new_value(void);
MMC_Node* mmc_node_new_object(void);
MMC_Node* mmc_node_new_array(void);
void mmc_node_free(MMC_Node* node);
MMC_Tag* mmc_tag_new(void);
void mmc_tag_free(MMC_Tag* tag);

#endif
