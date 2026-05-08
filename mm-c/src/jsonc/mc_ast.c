#include "mc_ast.h"
#include <stdlib.h>
#include <string.h>

MMC_Node* mmc_node_new_value(void) {
    MMC_Node* node = (MMC_Node*)malloc(sizeof(MMC_Node));
    if (node == NULL) return NULL;
    node->type = MMC_NODE_VALUE;
    memset(&node->data.value, 0, sizeof(MMC_Value));
    return node;
}

MMC_Node* mmc_node_new_object(void) {
    MMC_Node* node = (MMC_Node*)malloc(sizeof(MMC_Node));
    if (node == NULL) return NULL;
    node->type = MMC_NODE_OBJECT;
    node->data.object.fields = NULL;
    node->data.object.field_count = 0;
    memset(&node->data.object.tag, 0, sizeof(MMC_Tag));
    return node;
}

MMC_Node* mmc_node_new_array(void) {
    MMC_Node* node = (MMC_Node*)malloc(sizeof(MMC_Node));
    if (node == NULL) return NULL;
    node->type = MMC_NODE_ARRAY;
    node->data.array.items = NULL;
    node->data.array.item_count = 0;
    memset(&node->data.array.tag, 0, sizeof(MMC_Tag));
    return node;
}

static void free_value(MMC_Value* val) {
    if (val == NULL) return;
    if (val->text != NULL) {
        free(val->text);
    }
    if (val->data_type == MMC_DATA_STRING && val->data.string_value != NULL) {
        free(val->data.string_value);
    }
}

void mmc_node_free(MMC_Node* node) {
    if (node == NULL) return;
    switch (node->type) {
        case MMC_NODE_VALUE:
            free_value(&node->data.value);
            break;
        case MMC_NODE_OBJECT:
            for (size_t i = 0; i < node->data.object.field_count; i++) {
                if (node->data.object.fields[i].key != NULL) {
                    free(node->data.object.fields[i].key);
                }
                if (node->data.object.fields[i].value != NULL) {
                    mmc_node_free(node->data.object.fields[i].value);
                }
            }
            if (node->data.object.fields != NULL) {
                free(node->data.object.fields);
            }
            break;
        case MMC_NODE_ARRAY:
            for (size_t i = 0; i < node->data.array.item_count; i++) {
                if (node->data.array.items[i] != NULL) {
                    mmc_node_free(node->data.array.items[i]);
                }
            }
            if (node->data.array.items != NULL) {
                free(node->data.array.items);
            }
            break;
    }
    free(node);
}

MMC_Tag* mmc_tag_new(void) {
    MMC_Tag* tag = (MMC_Tag*)malloc(sizeof(MMC_Tag));
    if (tag == NULL) return NULL;
    tag->type = MMC_VALUE_TYPE_UNKNOWN;
    tag->size = -1;
    tag->required = -1;
    tag->min_value = 0;
    tag->max_value = 0;
    tag->min_length = 0;
    tag->max_length = 0;
    tag->deprecated = 0;
    return tag;
}

void mmc_tag_free(MMC_Tag* tag) {
    if (tag != NULL) {
        free(tag);
    }
}
