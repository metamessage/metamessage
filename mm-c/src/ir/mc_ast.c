#include "mc_ast.h"
#include <stdlib.h>
#include <string.h>

mm_node_t* mm_node_new_object(void)
{
    mm_node_t* node = (mm_node_t*)malloc(sizeof(mm_node_t));
    if (!node) return NULL;

    node->type = MM_NODE_OBJECT;
    node->data.object.fields = NULL;
    node->data.object.field_count = 0;
    node->data.object.capacity = 0;
    mm_tag_init(&node->tag);
    node->path = NULL;

    return node;
}

mm_node_t* mm_node_new_array(void)
{
    mm_node_t* node = (mm_node_t*)malloc(sizeof(mm_node_t));
    if (!node) return NULL;

    node->type = MM_NODE_ARRAY;
    node->data.array.items = NULL;
    node->data.array.item_count = 0;
    node->data.array.capacity = 0;
    mm_tag_init(&node->tag);
    node->path = NULL;

    return node;
}

mm_node_t* mm_node_new_value(void)
{
    mm_node_t* node = (mm_node_t*)malloc(sizeof(mm_node_t));
    if (!node) return NULL;

    node->type = MM_NODE_VALUE;
    node->data.value.text = NULL;
    mm_tag_init(&node->tag);
    node->path = NULL;

    return node;
}

mm_node_t* mm_node_new_doc(void)
{
    mm_node_t* node = (mm_node_t*)malloc(sizeof(mm_node_t));
    if (!node) return NULL;

    node->type = MM_NODE_DOC;
    node->data.doc.fields = NULL;
    node->data.doc.field_count = 0;
    node->data.doc.capacity = 0;
    mm_tag_init(&node->tag);
    node->path = NULL;

    return node;
}

void mm_node_free(mm_node_t* node)
{
    if (!node) return;

    switch (node->type) {
    case MM_NODE_OBJECT:
        for (size_t i = 0; i < node->data.object.field_count; i++) {
            free(node->data.object.fields[i].key);
            mm_node_free(node->data.object.fields[i].value);
        }
        free(node->data.object.fields);
        break;
    case MM_NODE_ARRAY:
        for (size_t i = 0; i < node->data.array.item_count; i++) {
            mm_node_free(node->data.array.items[i]);
        }
        free(node->data.array.items);
        break;
    case MM_NODE_VALUE:
        free(node->data.value.text);
        break;
    case MM_NODE_DOC:
        for (size_t i = 0; i < node->data.doc.field_count; i++) {
            free(node->data.doc.fields[i].key);
            mm_node_free(node->data.doc.fields[i].value);
        }
        free(node->data.doc.fields);
        break;
    default:
        break;
    }

    mm_tag_cleanup(&node->tag);
    free(node->path);
    free(node);
}

void mm_object_add_field(mm_node_t* obj, const char* key, mm_node_t* value)
{
    mm_field_t* new_fields = (mm_field_t*)realloc(
        obj->data.object.fields,
        (obj->data.object.field_count + 1) * sizeof(mm_field_t));

    if (!new_fields) return;

    obj->data.object.fields = new_fields;
    obj->data.object.fields[obj->data.object.field_count].key = strdup(key);
    obj->data.object.fields[obj->data.object.field_count].value = value;
    obj->data.object.field_count++;
}

void mm_array_add_item(mm_node_t* arr, mm_node_t* item)
{
    mm_node_t** new_items = (mm_node_t**)realloc(
        arr->data.array.items,
        (arr->data.array.item_count + 1) * sizeof(mm_node_t*));

    if (!new_items) return;

    arr->data.array.items = new_items;
    arr->data.array.items[arr->data.array.item_count] = item;
    arr->data.array.item_count++;
}

void mm_doc_add_field(mm_node_t* doc, const char* key, mm_node_t* value)
{
    mm_field_t* new_fields = (mm_field_t*)realloc(
        doc->data.doc.fields,
        (doc->data.doc.field_count + 1) * sizeof(mm_field_t));

    if (!new_fields) return;

    doc->data.doc.fields = new_fields;
    doc->data.doc.fields[doc->data.doc.field_count].key = strdup(key);
    doc->data.doc.fields[doc->data.doc.field_count].value = value;
    doc->data.doc.field_count++;
}