/**
 * tag_attributes.c
 *
 * MetaMessage C API — 标签属性（Tag Attributes）示例
 * 涵盖：desc, min, max, nullable, allow_empty, unique, raw,
 *       version, enums, pattern, default, mime, child 属性
 *
 * Compile:
 *   cd mm-c && gcc -I src -o ../examples/c/tag_attributes \
 *     ../examples/c/tag_attributes.c src/*.c -lm
 */

#include "mm.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

static void print_tag(const mm_tag_t* tag) {
    printf("    type=%s", mm_value_type_to_string(tag->type));
    if (tag->desc)        printf(" desc=\"%s\"", tag->desc);
    if (tag->min)         printf(" min=%s", tag->min);
    if (tag->max)         printf(" max=%s", tag->max);
    if (tag->nullable)    printf(" nullable");
    if (tag->allow_empty) printf(" allow_empty");
    if (tag->unique)      printf(" unique");
    if (tag->raw)         printf(" raw");
    if (tag->version > 0) printf(" version=%d", tag->version);
    if (tag->enums)       printf(" enum=%s", tag->enums);
    if (tag->pattern)     printf(" pattern=\"%s\"", tag->pattern);
    if (tag->default_val) printf(" default=%s", tag->default_val);
    if (tag->mime)        printf(" mime=%s", tag->mime);
    if (tag->location_offset != INT_MIN) printf(" location=%d", tag->location_offset);
    printf("\n");
}

int main(void) {
    printf("=== MetaMessage C API — 标签属性 ===\n\n");

    /* ---- 1. 基础属性：desc, min, max ---- */
    printf("-- 基础属性: desc, min, max --\n");
    mm_node_t* v1 = mm_int(25, .desc = "年龄", .min = 0, .max = 150);
    printf("  int with desc/min/max:\n");
    print_tag(&v1->data.value.tag);
    mm_node_free(v1);

    /* ---- 2. 约束属性：nullable, allow_empty, unique, raw ---- */
    printf("\n-- 约束属性: nullable, allow_empty, unique, raw --\n");

    mm_node_t* v2 = mm_str(NULL, .nullable = true, .desc = "可选名称");
    printf("  nullable str:\n");
    print_tag(&v2->data.value.tag);
    mm_node_free(v2);

    mm_node_t* v3 = mm_str("", .allow_empty = true, .desc = "允许为空");
    printf("  allow_empty str:\n");
    print_tag(&v3->data.value.tag);
    mm_node_free(v3);

    mm_node_t* v4 = mm_str("unique_key", .unique = true, .desc = "唯一键");
    printf("  unique str:\n");
    print_tag(&v4->data.value.tag);
    mm_node_free(v4);

    mm_node_t* v5 = mm_str("*raw*", .raw = true, .desc = "原始数据");
    printf("  raw str:\n");
    print_tag(&v5->data.value.tag);
    mm_node_free(v5);

    /* ---- 3. 版本与枚举 ---- */
    printf("\n-- 版本与枚举 --\n");

    mm_node_t* v6 = mm_str("1.0.0", .version = 2, .desc = "版本号");
    printf("  versioned:\n");
    print_tag(&v6->data.value.tag);
    mm_node_free(v6);

    mm_node_t* v7 = mm_int(0, .enums = "red|green|blue", .desc = "颜色枚举");
    printf("  enum:\n");
    print_tag(&v7->data.value.tag);
    mm_node_free(v7);

    /* ---- 4. pattern 与 default ---- */
    printf("\n-- pattern 与 default --\n");

    mm_node_t* v8 = mm_str("abc123", .pattern = "^[a-z0-9]+$", .desc = "字母数字");
    printf("  with pattern:\n");
    print_tag(&v8->data.value.tag);
    mm_node_free(v8);

    mm_node_t* v9 = mm_int(100, .default_val = "50", .desc = "默认50");
    printf("  with default:\n");
    print_tag(&v9->data.value.tag);
    mm_node_free(v9);

    /* ---- 5. MIME 类型 ---- */
    printf("\n-- MIME 类型 --\n");

    mm_node_t* v10 = mm_value_create_str(
        "image.png", MM_VALUE_STR,
        (mm_field_attr_t){ .desc = "头像", .mime = "image/png" });
    printf("  with mime:\n");
    print_tag(&v10->data.value.tag);
    mm_node_free(v10);

    /* ---- 6. 非默认类型 —— 用 mm_value_create_str 指定 ---- */
    printf("\n-- 非默认类型 --\n");

    mm_node_t* v11 = mm_value_create_str(
        "a1b2c3d4-e5f6-7890-abcd-ef1234567890", MM_VALUE_UUID,
        (mm_field_attr_t){ .desc = "UUID" });
    printf("  uuid:\n");
    print_tag(&v11->data.value.tag);
    mm_node_free(v11);

    mm_node_t* v12 = mm_value_create_str(
        "user@example.com", MM_VALUE_EMAIL,
        (mm_field_attr_t){ .desc = "电子邮箱" });
    printf("  email:\n");
    print_tag(&v12->data.value.tag);
    mm_node_free(v12);

    mm_node_t* v13 = mm_value_create_str(
        "https://example.com", MM_VALUE_URL,
        (mm_field_attr_t){ .desc = "网址" });
    printf("  url:\n");
    print_tag(&v13->data.value.tag);
    mm_node_free(v13);

    mm_node_t* v14 = mm_value_create_str(
        "192.168.1.1", MM_VALUE_IP,
        (mm_field_attr_t){ .desc = "IP 地址" });
    printf("  ip:\n");
    print_tag(&v14->data.value.tag);
    mm_node_free(v14);

    /* ---- 7. child 属性（数组/对象的子元素约束） ---- */
    printf("\n-- child 属性 --\n");

    /* 使用底层的 tag 操作 */
    mm_node_t* arr = mm_arr_new();
    arr->data.array.tag.child_desc = strdup("分数列表");
    arr->data.array.tag.child_type = MM_VALUE_I;
    arr->data.array.tag.child_min   = strdup("0");
    arr->data.array.tag.child_max   = strdup("100");

    mm_arr_add(arr, mm_int(85));
    mm_arr_add(arr, mm_int(90));

    printf("  array child attributes:\n");
    printf("    child_desc=%s\n", arr->data.array.tag.child_desc);
    printf("    child_type=%s\n", mm_value_type_to_string(arr->data.array.tag.child_type));
    printf("    child_min=%s  child_max=%s\n",
           arr->data.array.tag.child_min, arr->data.array.tag.child_max);

    /* ---- 8. 完整 round-trip 验证 ---- */
    printf("\n-- Round-trip 验证 --\n");

    mm_buffer_t* enc = mm_encode(arr);
    assert(enc != NULL);
    mm_node_t* dec = mm_decode(enc);
    assert(dec != NULL && dec->type == MM_NODE_ARRAY);

    printf("  decoded array: %zu items\n", dec->data.array.item_count);
    printf("  child_desc=%s\n", dec->data.array.tag.child_desc);
    printf("  child_min=%s  child_max=%s\n",
           dec->data.array.tag.child_min, dec->data.array.tag.child_max);

    mm_node_free(arr);
    mm_node_free(dec);
    mm_buffer_free(enc);

    printf("\n=== 完成 ===\n");
    return 0;
}
