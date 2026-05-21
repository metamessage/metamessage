/**
 * basic_usage.c
 *
 * MetaMessage C API 基础用法示例
 * 涵盖：创建值、构建对象/数组、二进制编码解码
 *
 * Compile:
 *   cd mm-c && gcc -I src -o ../examples/c/basic_usage ../examples/c/basic_usage.c src/*.c -lm
 *
 * 或用 cmake 构建后直接链接 libmmc.a
 */

#include "mm.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

int main(void) {
    printf("=== MetaMessage C API — 基础用法 ===\n\n");

    /* ---- 1. 创建各种类型的值 ---- */
    printf("-- 创建值 --\n");

    mm_node_t* age   = mm_int(30, .desc = "年龄", .min = 0, .max = 150);
    mm_node_t* name  = mm_str("Alice", .desc = "姓名", .min = 1, .max = 64);
    mm_node_t* score = mm_float(95.5, .desc = "分数");
    mm_node_t* flag  = mm_bool(true, .desc = "激活状态");

    printf("  age   = %s (type: %s)\n", age->data.value.text,
           mm_value_type_to_string(age->data.value.tag.type));
    printf("  name  = %s\n", name->data.value.text);
    printf("  score = %s\n", score->data.value.text);
    printf("  flag  = %s\n", flag->data.value.text);

    /* ---- 2. 构建对象 ---- */
    printf("\n-- 构建对象 --\n");

    mm_obj_t* person = mm_obj_new();
    mm_obj_set(person, "name",  name);
    mm_obj_set(person, "age",   age);
    mm_obj_set(person, "score", score);
    mm_obj_set(person, "active", flag);

    printf("  person 对象有 %zu 个字段\n", person->data.object.field_count);
    for (size_t i = 0; i < person->data.object.field_count; i++) {
        printf("    [%zu] %s => %s\n", i,
               person->data.object.fields[i].key,
               person->data.object.fields[i].value->data.value.text);
    }

    /* ---- 3. 构建数组 ---- */
    printf("\n-- 构建数组 --\n");

    mm_node_t* scores = mm_arr_new();
    mm_arr_add(scores, mm_int(95, .desc = "语文"));
    mm_arr_add(scores, mm_int(87, .desc = "数学"));
    mm_arr_add(scores, mm_int(92, .desc = "英语"));

    printf("  scores 数组有 %zu 个元素\n", scores->data.array.item_count);
    for (size_t i = 0; i < scores->data.array.item_count; i++) {
        printf("    [%zu] %s\n", i, scores->data.array.items[i]->data.value.text);
    }

    /* ---- 4. 对象嵌套数组 ---- */
    printf("\n-- 对象嵌套数组 --\n");

    mm_obj_t* student = mm_obj_new();
    mm_obj_set(student, "name",   mm_str("Bob", .desc = "姓名"));
    mm_obj_set(student, "scores", scores);

    printf("  student 有 %zu 个字段\n", student->data.object.field_count);

    /* ---- 5. 二进制编码 / 解码 ---- */
    printf("\n-- 二进制编码/解码 --\n");

    mm_buffer_t* enc = mm_encode(student);
    assert(enc != NULL);
    printf("  编码后: %zu 字节\n", enc->size);

    mm_node_t* dec = mm_decode(enc);
    assert(dec != NULL);
    assert(dec->type == MM_NODE_OBJECT);
    printf("  解码成功: 对象有 %zu 个字段\n", dec->data.object.field_count);

    /* ---- 6. JSONC 序列化 / 反序列化 ---- */
    printf("\n-- JSONC 序列化/反序列化 --\n");

    char* jsonc = mm_to_jsonc(student);
    assert(jsonc != NULL);
    printf("  JSONC 输出:\n%s\n", jsonc);

    mm_node_t* parsed = mm_from_jsonc(jsonc);
    assert(parsed != NULL);
    printf("  从 JSONC 解析: 对象有 %zu 个字段\n\n", parsed->data.object.field_count);

    /* ---- 清理 ---- */
    mm_node_free(student);
    mm_node_free(dec);
    mm_node_free(parsed);
    mm_buffer_free(enc);
    mm_string_free(jsonc);

    printf("=== 完成 ===\n");
    return 0;
}
