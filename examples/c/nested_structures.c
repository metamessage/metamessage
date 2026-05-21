/**
 * nested_structures.c
 *
 * MetaMessage C API — 嵌套结构示例
 * 展示：嵌套对象、数组嵌套对象、混合类型数组
 *
 * Compile:
 *   cd mm-c && gcc -I src -o ../examples/c/nested_structures \
 *     ../examples/c/nested_structures.c src/*.c -lm
 */

#include "mm.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

int main(void) {
    printf("=== MetaMessage C API — 嵌套结构 ===\n\n");

    /* ---- 1. 嵌套对象 ---- */
    printf("-- 1. 嵌套对象 --\n");

    /* 地址对象 */
    mm_obj_t* address = mm_obj_new();
    mm_obj_set(address, "street", mm_str("123 Main St", .desc = "街道", .min = 1, .max = 200));
    mm_obj_set(address, "city",   mm_str("Beijing", .desc = "城市"));
    mm_obj_set(address, "zip",    mm_value_create_str("100000", MM_VALUE_STR,
        (mm_field_attr_t){ .desc = "邮编", .pattern = "^[0-9]{6}$" }));

    /* 人员对象包含地址 */
    mm_obj_t* person = mm_obj_new();
    mm_obj_set(person, "name", mm_str("Alice", .desc = "姓名"));
    mm_obj_set(person, "address", address);

    /* 编码/解码 */
    mm_buffer_t* enc = mm_encode(person);
    assert(enc != NULL);
    mm_node_t* dec = mm_decode(enc);
    assert(dec != NULL && dec->type == MM_NODE_OBJECT);

    printf("  人员对象有 %zu 个字段\n", dec->data.object.field_count);

    for (size_t i = 0; i < dec->data.object.field_count; i++) {
        const char* key = dec->data.object.fields[i].key;
        mm_node_t* val = dec->data.object.fields[i].value;
        if (val->type == MM_NODE_OBJECT) {
            printf("    %s: (嵌套对象, %zu 字段)\n", key, val->data.object.field_count);
            for (size_t j = 0; j < val->data.object.field_count; j++) {
                printf("      [%zu] %s = %s\n", j,
                       val->data.object.fields[j].key,
                       val->data.object.fields[j].value->data.value.text);
            }
        } else {
            printf("    %s = %s\n", key, val->data.value.text);
        }
    }

    mm_node_free(person);
    mm_node_free(dec);
    mm_buffer_free(enc);

    /* ---- 2. 数组嵌套对象（对象列表） ---- */
    printf("\n-- 2. 数组嵌套对象（对象列表） --\n");

    mm_node_t* people = mm_arr_new();

    mm_obj_t* p1 = mm_obj_new();
    mm_obj_set(p1, "name", mm_str("Alice"));
    mm_obj_set(p1, "age",  mm_int(30));
    mm_arr_add(people, p1);

    mm_obj_t* p2 = mm_obj_new();
    mm_obj_set(p2, "name", mm_str("Bob"));
    mm_obj_set(p2, "age",  mm_int(25));
    mm_arr_add(people, p2);

    mm_obj_t* p3 = mm_obj_new();
    mm_obj_set(p3, "name", mm_str("Charlie"));
    mm_obj_set(p3, "age",  mm_int(35));
    mm_arr_add(people, p3);

    printf("  people 数组有 %zu 个元素\n", people->data.array.item_count);
    for (size_t i = 0; i < people->data.array.item_count; i++) {
        mm_node_t* obj = people->data.array.items[i];
        printf("    [%zu] 对象 (%zu 字段)\n", i, obj->data.object.field_count);
        for (size_t j = 0; j < obj->data.object.field_count; j++) {
            printf("      %s = %s\n",
                   obj->data.object.fields[j].key,
                   obj->data.object.fields[j].value->data.value.text);
        }
    }

    /* 编码后 JSONC */
    char* jsonc = mm_to_jsonc(people);
    assert(jsonc != NULL);
    printf("  JSONC:\n%s\n", jsonc);
    mm_string_free(jsonc);
    mm_node_free(people);

    /* ---- 3. 混合类型数组 ---- */
    printf("\n-- 3. 混合类型数组 --\n");

    mm_node_t* mixed = mm_arr_new();
    mm_arr_add(mixed, mm_int(42));
    mm_arr_add(mixed, mm_str("hello"));
    mm_arr_add(mixed, mm_bool(true));
    mm_arr_add(mixed, mm_float(3.14));

    mm_buffer_t* mix_enc = mm_encode(mixed);
    assert(mix_enc != NULL);
    mm_node_t* mix_dec = mm_decode(mix_enc);
    assert(mix_dec != NULL && mix_dec->type == MM_NODE_ARRAY);

    printf("  混合数组有 %zu 个元素:\n", mix_dec->data.array.item_count);
    for (size_t i = 0; i < mix_dec->data.array.item_count; i++) {
        mm_node_t* item = mix_dec->data.array.items[i];
        printf("    [%zu] %s  (type=%s)\n", i, item->data.value.text,
               mm_value_type_to_string(item->data.value.tag.type));
    }

    mm_node_free(mixed);
    mm_node_free(mix_dec);
    mm_buffer_free(mix_enc);

    printf("\n=== 完成 ===\n");
    return 0;
}
