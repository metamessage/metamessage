/**
 * jsonc_roundtrip.c
 *
 * MetaMessage C API — JSONC 往返示例
 * 展示 JSONC（带注释的 JSON）序列化与反序列化，
 * 以及 `mm_from_jsonc` / `mm_to_jsonc` 的完整生命周期。
 *
 * Compile:
 *   cd mm-c && gcc -I src -o ../examples/c/jsonc_roundtrip \
 *     ../examples/c/jsonc_roundtrip.c src/*.c -lm
 */

#include "mm.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

int main(void) {
    printf("=== MetaMessage C API — JSONC 往返 ===\n\n");

    /* ---- 1. 构建一个带标签的对象 ---- */
    printf("-- 1. 构建带元数据的对象 --\n");

    mm_obj_t* user = mm_obj_new();

    mm_obj_set(user, "id",
        mm_value_create_str("42", MM_VALUE_U64,
            (mm_field_attr_t){ .desc = "用户ID", .min = 1, .max = 999999 }));

    mm_obj_set(user, "name",
        mm_str("Alice", .desc = "用户名", .min = 1, .max = 50));

    mm_obj_set(user, "email",
        mm_value_create_str("alice@example.com", MM_VALUE_EMAIL,
            (mm_field_attr_t){ .desc = "电子邮箱" }));

    mm_obj_set(user, "age",
        mm_int(28, .desc = "年龄", .min = 0, .max = 150));

    mm_obj_set(user, "active",
        mm_bool(true, .desc = "是否激活"));

    printf("  user 对象有 %zu 个字段\n", user->data.object.field_count);

    /* ---- 2. 对象 → JSONC ---- */
    printf("\n-- 2. 对象 → JSONC --\n");

    char* jsonc_str = mm_to_jsonc(user);
    assert(jsonc_str != NULL);
    printf("%s\n", jsonc_str);

    /* ---- 3. JSONC → 对象（反序列化） ---- */
    printf("-- 3. JSONC → 对象 --\n");

    mm_node_t* parsed_user = mm_from_jsonc(jsonc_str);
    assert(parsed_user != NULL);
    assert(parsed_user->type == MM_NODE_OBJECT);
    printf("  解析成功: 对象有 %zu 个字段\n", parsed_user->data.object.field_count);

    for (size_t i = 0; i < parsed_user->data.object.field_count; i++) {
        mm_node_t* val = parsed_user->data.object.fields[i].value;
        printf("  [%zu] %s = %s  (type=%s)\n", i,
               parsed_user->data.object.fields[i].key,
               val->data.value.text,
               mm_value_type_to_string(val->data.value.tag.type));
    }

    /* ---- 4. 对象 → 二进制 → JSONC（完整往返） ---- */
    printf("\n-- 4. 对象 → 二进制 → JSONC（完整往返） --\n");

    mm_buffer_t* bin = mm_encode(user);
    assert(bin != NULL);
    printf("  二进制编码: %zu 字节\n", bin->size);

    mm_node_t* decoded = mm_decode(bin);
    assert(decoded != NULL);

    char* jsonc2 = mm_to_jsonc(decoded);
    assert(jsonc2 != NULL);
    printf("  解码后 JSONC:\n%s\n", jsonc2);

    /* ---- 5. 验证 round-trip 一致性 ---- */
    printf("-- 5. 验证一致性 --\n");

    int same = (strcmp(jsonc_str, jsonc2) == 0);
    printf("  JSONC round-trip 结果 %s\n", same ? "一致 ✓" : "不一致 ✗");

    /* ---- 清理 ---- */
    mm_node_free(user);
    mm_node_free(parsed_user);
    mm_node_free(decoded);
    mm_buffer_free(bin);
    mm_string_free(jsonc_str);
    mm_string_free(jsonc2);

    printf("\n=== 完成 ===\n");
    return 0;
}
