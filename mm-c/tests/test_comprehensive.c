#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include "mm.h"
#include "ir/mc_value_type.h"
#include "ir/mc_tag.h"
#include "ir/mc_ast.h"

static int tests_passed = 0;
static int tests_failed = 0;

#define TEST(name, expr) do { \
    if (!(expr)) { \
        printf("  FAIL: %s\n", name); \
        tests_failed++; \
    } else { \
        printf("  PASS: %s\n", name); \
        tests_passed++; \
    } \
} while(0)

#define ASSERT_EQ(actual, expected) do { \
    if ((actual) != (expected)) { \
        printf("  FAIL: expected %lld but got %lld\n", (long long)(expected), (long long)(actual)); \
        tests_failed++; \
    } else { \
        printf("  PASS\n"); \
        tests_passed++; \
    } \
} while(0)

#define ASSERT_STR_EQ(actual, expected) do { \
    if (strcmp((actual), (expected)) != 0) { \
        printf("  FAIL: expected \"%s\" but got \"%s\"\n", (expected), (actual)); \
        tests_failed++; \
    } else { \
        printf("  PASS\n"); \
        tests_passed++; \
    } \
} while(0)

static void test_single_value_roundtrip(void) {
    printf("\n=== Single Value Roundtrip Tests ===\n");

    printf("  int roundtrip: ");
    {
        mm_node_t* v = mm_int(42, .desc="answer");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("type is value", dec->type == MM_NODE_VALUE);
        TEST("text is 42", dec->data.value.text && strcmp(dec->data.value.text, "42") == 0);
        TEST("tag type I", dec->data.value.tag.type == MM_VALUE_I);
        TEST("desc answer", dec->data.value.tag.desc && strcmp(dec->data.value.tag.desc, "answer") == 0);
        mm_node_free(v);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  str roundtrip: ");
    {
        mm_node_t* v = mm_str("hello", .desc="greeting");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("type is value", dec->type == MM_NODE_VALUE);
        TEST("text is hello", dec->data.value.text && strcmp(dec->data.value.text, "hello") == 0);
        TEST("tag type STR", dec->data.value.tag.type == MM_VALUE_STR);
        TEST("desc greeting", dec->data.value.tag.desc && strcmp(dec->data.value.tag.desc, "greeting") == 0);
        mm_node_free(v);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  bool true roundtrip: ");
    {
        mm_node_t* v = mm_bool(true, .desc="flag");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("type is value", dec->type == MM_NODE_VALUE);
        TEST("text is true", dec->data.value.text && strcmp(dec->data.value.text, "true") == 0);
        TEST("tag type BOOL", dec->data.value.tag.type == MM_VALUE_BOOL);
        mm_node_free(v);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  bool false roundtrip: ");
    {
        mm_node_t* v = mm_bool(false, .desc="inactive");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("text is false", dec->data.value.text && strcmp(dec->data.value.text, "false") == 0);
        mm_node_free(v);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  float roundtrip: ");
    {
        mm_node_t* v = mm_float(3.14, .desc="pi");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("type is F64", dec->data.value.tag.type == MM_VALUE_F64);
        TEST("desc pi", dec->data.value.tag.desc && strcmp(dec->data.value.tag.desc, "pi") == 0);
        mm_node_free(v);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }
}

static void test_object_roundtrip(void) {
    printf("\n=== Object Roundtrip Tests ===\n");

    printf("  person with tag attrs: ");
    {
        mm_node_t* age = mm_int(30, .min=0, .max=150, .desc="年龄");
        mm_node_t* name = mm_str("Alice", .min=1, .max=64, .desc="姓名");
        mm_obj_t* person = mm_obj_new();
        mm_obj_set(person, "name", name);
        mm_obj_set(person, "age", age);
        TEST("person has 2 fields", person->data.object.field_count == 2);

        mm_buffer_t* enc = mm_encode(person);
        TEST("encode not null", enc != NULL);
        TEST("encode size > 0", enc->size > 0);

        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec) {
            TEST("decoded is object", dec->type == MM_NODE_OBJECT);
            if (dec->type == MM_NODE_OBJECT) {
                TEST("decoded has 2 fields", dec->data.object.field_count == 2);
                if (dec->data.object.field_count == 2) {
                    TEST("field 0 key is name", strcmp(dec->data.object.fields[0].key, "name") == 0);
                    TEST("field 1 key is age", strcmp(dec->data.object.fields[1].key, "age") == 0);

                    mm_node_t* name_val = dec->data.object.fields[0].value;
                    mm_node_t* age_val = dec->data.object.fields[1].value;

                    TEST("name text Alice", name_val->data.value.text && strcmp(name_val->data.value.text, "Alice") == 0);
                    TEST("age text 30", age_val->data.value.text && strcmp(age_val->data.value.text, "30") == 0);

                    TEST("name tag desc 姓名", name_val->data.value.tag.desc && strcmp(name_val->data.value.tag.desc, "姓名") == 0);
                    TEST("name tag min 1", name_val->data.value.tag.min && strcmp(name_val->data.value.tag.min, "1") == 0);
                    TEST("name tag max 64", name_val->data.value.tag.max && strcmp(name_val->data.value.tag.max, "64") == 0);

                    TEST("age tag desc 年龄", age_val->data.value.tag.desc && strcmp(age_val->data.value.tag.desc, "年龄") == 0);
                    TEST("age tag min 0", age_val->data.value.tag.min && strcmp(age_val->data.value.tag.min, "0") == 0);
                    TEST("age tag max 150", age_val->data.value.tag.max && strcmp(age_val->data.value.tag.max, "150") == 0);
                }
            }
        }
        mm_node_free(person);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  empty object: ");
    {
        mm_obj_t* obj = mm_obj_new();
        mm_buffer_t* enc = mm_encode(obj);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("decoded is object", dec->type == MM_NODE_OBJECT);
        TEST("decoded has 0 fields", dec->data.object.field_count == 0);
        mm_node_free(obj);
        mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  nested object: ");
    {
        mm_node_t* street = mm_str("123 Main St", .desc="street");
        mm_obj_t* addr = mm_obj_new();
        mm_obj_set(addr, "street", street);

        mm_node_t* name = mm_str("Bob");
        mm_obj_t* person = mm_obj_new();
        mm_obj_set(person, "name", name);
        mm_obj_set(person, "address", addr);

        mm_buffer_t* enc = mm_encode(person);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec && dec->type == MM_NODE_OBJECT) {
            TEST("decoded has 2 fields", dec->data.object.field_count == 2);
            if (dec->data.object.field_count == 2) {
                mm_node_t* addr_val = NULL;
                for (size_t i = 0; i < dec->data.object.field_count; i++) {
                    if (strcmp(dec->data.object.fields[i].key, "address") == 0) {
                        addr_val = dec->data.object.fields[i].value;
                        break;
                    }
                }
                TEST("address field exists", addr_val != NULL);
                TEST("address is object", addr_val != NULL && addr_val->type == MM_NODE_OBJECT);
                if (addr_val && addr_val->type == MM_NODE_OBJECT) {
                    TEST("address has 1 field", addr_val->data.object.field_count == 1);
                    if (addr_val->data.object.field_count == 1) {
                        TEST("street text", strcmp(addr_val->data.object.fields[0].value->data.value.text, "123 Main St") == 0);
                    }
                }
            }
        }
        mm_node_free(person);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }
}

static void test_array_roundtrip(void) {
    printf("\n=== Array Roundtrip Tests ===\n");

    printf("  simple int array: ");
    {
        mm_node_t* arr = mm_arr_new();
        mm_node_t* v1 = mm_int(1);
        mm_node_t* v2 = mm_int(2);
        mm_node_t* v3 = mm_int(3);
        mm_arr_add(arr, v1);
        mm_arr_add(arr, v2);
        mm_arr_add(arr, v3);
        TEST("array has 3 items", arr->data.array.item_count == 3);

        mm_buffer_t* enc = mm_encode(arr);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("decoded is array", dec->type == MM_NODE_ARRAY);
        if (dec && dec->type == MM_NODE_ARRAY) {
            TEST("decoded has 3 items", dec->data.array.item_count == 3);
            if (dec->data.array.item_count == 3) {
                TEST("item 0 is 1", strcmp(dec->data.array.items[0]->data.value.text, "1") == 0);
                TEST("item 1 is 2", strcmp(dec->data.array.items[1]->data.value.text, "2") == 0);
                TEST("item 2 is 3", strcmp(dec->data.array.items[2]->data.value.text, "3") == 0);
            }
        }
        mm_node_free(arr);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  mixed array: ");
    {
        mm_node_t* arr = mm_arr_new();
        mm_arr_add(arr, mm_str("hello"));
        mm_arr_add(arr, mm_int(42));
        mm_arr_add(arr, mm_bool(true));
        TEST("array has 3 items", arr->data.array.item_count == 3);

        mm_buffer_t* enc = mm_encode(arr);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec && dec->type == MM_NODE_ARRAY) {
            TEST("decoded has 3 items", dec->data.array.item_count == 3);
            if (dec->data.array.item_count == 3) {
                TEST("item 0 str", dec->data.array.items[0]->data.value.tag.type == MM_VALUE_STR);
                TEST("item 1 int", dec->data.array.items[1]->data.value.tag.type == MM_VALUE_I);
                TEST("item 2 bool", dec->data.array.items[2]->data.value.tag.type == MM_VALUE_BOOL);
            }
        }
        mm_node_free(arr);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  array in object: ");
    {
        mm_node_t* arr = mm_arr_new();
        mm_arr_add(arr, mm_int(10));
        mm_arr_add(arr, mm_int(20));
        mm_arr_add(arr, mm_int(30));

        mm_obj_t* obj = mm_obj_new();
        mm_obj_set(obj, "scores", arr);
        mm_obj_set(obj, "name", mm_str("test"));

        mm_buffer_t* enc = mm_encode(obj);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec && dec->type == MM_NODE_OBJECT) {
            TEST("decoded has 2 fields", dec->data.object.field_count == 2);
            if (dec->data.object.field_count == 2) {
                mm_node_t* scores = NULL;
                for (size_t i = 0; i < dec->data.object.field_count; i++) {
                    if (strcmp(dec->data.object.fields[i].key, "scores") == 0) {
                        scores = dec->data.object.fields[i].value;
                        break;
                    }
                }
                TEST("scores field exists", scores != NULL);
                TEST("scores is array", scores != NULL && scores->type == MM_NODE_ARRAY);
                if (scores && scores->type == MM_NODE_ARRAY) {
                    TEST("scores has 3 items", scores->data.array.item_count == 3);
                }
            }
        }
        mm_node_free(obj);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }
}

static void test_tag_attributes_roundtrip(void) {
    printf("\n=== Tag Attributes Roundtrip Tests ===\n");

    printf("  nullable with desc: ");
    {
        mm_node_t* v = mm_str("", .nullable=true, .desc="optional name");
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec) {
            TEST("desc preserved", dec->data.value.tag.desc && strcmp(dec->data.value.tag.desc, "optional name") == 0);
            TEST("nullable preserved", dec->data.value.tag.nullable == true);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  allow_empty: ");
    {
        mm_node_t* v = mm_str("", .allow_empty=true);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("allow_empty preserved", dec->data.value.tag.allow_empty == true);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  unique: ");
    {
        mm_node_t* v = mm_str("unique", .unique=true);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("unique preserved", dec->data.value.tag.unique == true);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  raw: ");
    {
        mm_node_t* v = mm_str("rawdata", .raw=true);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("raw preserved", dec->data.value.tag.raw == true);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  version: ");
    {
        mm_node_t* v = mm_str("1.0.0", .version=1);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("version preserved", dec->data.value.tag.version == 1);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  enum: ");
    {
        mm_node_t* v = mm_int(0, .enums="a|b|c");
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("enum preserved", dec->data.value.tag.enums && strcmp(dec->data.value.tag.enums, "a|b|c") == 0);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  pattern: ");
    {
        mm_node_t* v = mm_str("abc", .pattern="^[a-z]+$");
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("pattern preserved", dec->data.value.tag.pattern && strcmp(dec->data.value.tag.pattern, "^[a-z]+$") == 0);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  default value: ");
    {
        mm_node_t* v = mm_int(0, .default_val="100");
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        if (dec) {
            TEST("default preserved", dec->data.value.tag.default_val && strcmp(dec->data.value.tag.default_val, "100") == 0);
        }
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  child attributes: ");
    {
        mm_node_t* arr = mm_arr_new();
        arr->data.array.tag.child_desc = strdup("child items");
        arr->data.array.tag.child_min = strdup("1");
        arr->data.array.tag.child_max = strdup("10");
        arr->data.array.tag.child_type = MM_VALUE_I;
        mm_arr_add(arr, mm_int(5));

        mm_buffer_t* enc = mm_encode(arr);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        if (dec && dec->type == MM_NODE_ARRAY) {
            TEST("child_desc preserved", dec->data.array.tag.child_desc && strcmp(dec->data.array.tag.child_desc, "child items") == 0);
            TEST("child_min preserved", dec->data.array.tag.child_min && strcmp(dec->data.array.tag.child_min, "1") == 0);
            TEST("child_max preserved", dec->data.array.tag.child_max && strcmp(dec->data.array.tag.child_max, "10") == 0);
            if (dec->data.array.item_count > 0) {
                TEST("child type inherited", dec->data.array.items[0]->data.value.tag.type == MM_VALUE_I);
            }
        }
        mm_node_free(arr);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }
}

static void test_negative_int_roundtrip(void) {
    printf("\n=== Negative Int Roundtrip Tests ===\n");

    printf("  negative int: ");
    {
        mm_node_t* v = mm_int(-42);
        mm_buffer_t* enc = mm_encode(v);
        TEST("encode not null", enc != NULL);
        mm_node_t* dec = mm_decode(enc);
        TEST("decode not null", dec != NULL);
        TEST("text is -42", dec && dec->data.value.text && strcmp(dec->data.value.text, "-42") == 0);
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  zero: ");
    {
        mm_node_t* v = mm_int(0);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        TEST("text is 0", dec && dec->data.value.text && strcmp(dec->data.value.text, "0") == 0);
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  large positive: ");
    {
        mm_node_t* v = mm_int(100000);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        TEST("text is 100000", dec && dec->data.value.text && strcmp(dec->data.value.text, "100000") == 0);
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }

    printf("  large negative: ");
    {
        mm_node_t* v = mm_int(-100000);
        mm_buffer_t* enc = mm_encode(v);
        mm_node_t* dec = mm_decode(enc);
        TEST("text is -100000", dec && dec->data.value.text && strcmp(dec->data.value.text, "-100000") == 0);
        mm_node_free(v);
        if (dec) mm_node_free(dec);
        mm_buffer_free(enc);
    }
}

static void test_value_type(void) {
    printf("\n=== ValueType Tests ===\n");
    TEST("Str to string", strcmp(mm_value_type_to_string(MM_VALUE_STR), "str") == 0);
    TEST("U8 to string", strcmp(mm_value_type_to_string(MM_VALUE_U8), "u8") == 0);
    TEST("Parse str", mm_value_type_parse("str") == MM_VALUE_STR);
    TEST("Parse u8", mm_value_type_parse("u8") == MM_VALUE_U8);
    TEST("Parse int", mm_value_type_parse("int") == MM_VALUE_I);
    TEST("Parse uint64", mm_value_type_parse("uint64") == MM_VALUE_U64);
    TEST("Parse I32", mm_value_type_parse("I32") == MM_VALUE_I32);
    TEST("Parse unknown", mm_value_type_parse("invalid") == MM_VALUE_UNKNOWN);
    TEST("Parse NULL", mm_value_type_parse(NULL) == MM_VALUE_UNKNOWN);
    TEST("Unknown to string", strcmp(mm_value_type_to_string(MM_VALUE_UNKNOWN), "unknown") == 0);
}

static void test_tag(void) {
    printf("\n=== Tag Tests ===\n");
    mm_tag_t tag;
    mm_tag_init(&tag);
    TEST("Tag init type unknown", tag.type == MM_VALUE_UNKNOWN);
    TEST("Tag init not null", tag.is_null == 0);

    tag.desc = strdup("test description");
    tag.type = MM_VALUE_STR;

    char* str = mm_tag_to_string(&tag);
    TEST("Tag to string contains desc", strstr(str, "test description") != NULL);
    free(str);

    mm_tag_cleanup(&tag);

    mm_tag_t parsed = mm_tag_parse("// mm: desc=\"姓名\"; min=1; max=64; type=str");
    TEST("Parsed desc", parsed.desc && strcmp(parsed.desc, "姓名") == 0);
    TEST("Parsed min", parsed.min && strcmp(parsed.min, "1") == 0);
    TEST("Parsed max", parsed.max && strcmp(parsed.max, "64") == 0);
    TEST("Parsed type", parsed.type == MM_VALUE_STR);
    mm_tag_cleanup(&parsed);

    mm_tag_t t1, t2;
    mm_tag_init(&t1);
    mm_tag_init(&t2);
    t1.desc = strdup("original");
    t1.nullable = 0;
    t2.desc = strdup("override");
    t2.nullable = 1;
    mm_tag_merge(&t1, &t2);
    TEST("Merge desc override", t1.desc && strcmp(t1.desc, "override") == 0);
    TEST("Merge nullable", t1.nullable == 1);
    mm_tag_cleanup(&t1);
    mm_tag_cleanup(&t2);

    mm_tag_t child_parent;
    mm_tag_init(&child_parent);
    child_parent.child_desc = strdup("inherited");
    child_parent.child_type = MM_VALUE_STR;
    child_parent.child_min = strdup("5");
    child_parent.child_nullable = true;

    mm_tag_t child;
    mm_tag_init(&child);
    mm_tag_inherit(&child, &child_parent);
    TEST("Inherit desc", child.desc && strcmp(child.desc, "inherited") == 0);
    TEST("Inherit type", child.type == MM_VALUE_STR);
    TEST("Inherit min", child.min && strcmp(child.min, "5") == 0);
    TEST("Inherit nullable", child.nullable == true);
    mm_tag_cleanup(&child);
    mm_tag_cleanup(&child_parent);

    mm_tag_t parsed_child = mm_tag_parse("// mm: child_desc=\"items\"; child_min=1; child_max=10; child_type=i");
    TEST("Parsed child_desc", parsed_child.child_desc && strcmp(parsed_child.child_desc, "items") == 0);
    TEST("Parsed child_min", parsed_child.child_min && strcmp(parsed_child.child_min, "1") == 0);
    TEST("Parsed child_max", parsed_child.child_max && strcmp(parsed_child.child_max, "10") == 0);
    TEST("Parsed child_type", parsed_child.child_type == MM_VALUE_I);
    mm_tag_cleanup(&parsed_child);
}

static void test_ast(void) {
    printf("\n=== AST Tests ===\n");
    mm_node_t* obj = mm_node_new_object();
    TEST("Object created", obj != NULL);
    TEST("Object type", obj->type == MM_NODE_OBJECT);
    TEST("Object no fields", obj->data.object.field_count == 0);
    mm_node_free(obj);

    mm_node_t* arr = mm_node_new_array();
    TEST("Array created", arr != NULL);
    TEST("Array type", arr->type == MM_NODE_ARRAY);
    mm_node_free(arr);

    mm_node_t* val = mm_node_new_value();
    TEST("Value created", val != NULL);
    TEST("Value type", val->type == MM_NODE_VALUE);
    TEST("Value text NULL", val->data.value.text == NULL);
    mm_node_free(val);

    mm_node_t* obj2 = mm_node_new_object();
    mm_node_t* v1 = mm_node_new_value();
    v1->data.value.text = strdup("hello");
    mm_object_add_field(obj2, "msg", v1);
    TEST("Object has 1 field", obj2->data.object.field_count == 1);
    TEST("Field key correct", strcmp(obj2->data.object.fields[0].key, "msg") == 0);
    mm_node_free(obj2);

    mm_node_t* arr2 = mm_node_new_array();
    mm_node_t* v2 = mm_node_new_value();
    v2->data.value.text = strdup("item1");
    mm_array_add_item(arr2, v2);
    TEST("Array has 1 item", arr2->data.array.item_count == 1);
    mm_node_free(arr2);
}

static void test_byte_compatibility(void) {
    printf("\n=== Byte Compatibility Tests ===\n");

    printf("  tag bytes with string desc: ");
    {
        mm_tag_t tag;
        mm_tag_init(&tag);
        tag.desc = strdup("hello");
        size_t len;
        uint8_t* bytes = mm_tag_bytes(&tag, &len);
        TEST("bytes not null", bytes != NULL);
        TEST("length > 0", len > 0);
        if (bytes && len > 0) {
            uint8_t key = bytes[0] & 0xF8;
            TEST("first key is KDESC", key == MM_TAG_KDESC);
            uint8_t payload_len = bytes[0] & 0x07;
            TEST("desc length 5", payload_len == 5);
            if (len >= 6) {
                char desc_str[16] = {0};
                memcpy(desc_str, bytes + 1, 5);
                TEST("desc content hello", strcmp(desc_str, "hello") == 0);
            }
        }
        free(bytes);
        mm_tag_cleanup(&tag);
    }

    printf("  tag bytes with min and max: ");
    {
        mm_tag_t tag;
        mm_tag_init(&tag);
        tag.min = strdup("0");
        tag.max = strdup("150");
        size_t len;
        uint8_t* bytes = mm_tag_bytes(&tag, &len);
        TEST("bytes not null", bytes != NULL);
        TEST("contains min and max", len >= 4);
        free(bytes);
        mm_tag_cleanup(&tag);
    }
}

int main(void) {
    printf("=== MMC Comprehensive Tests ===\n");

    test_value_type();
    test_tag();
    test_ast();
    test_byte_compatibility();
    test_single_value_roundtrip();
    test_object_roundtrip();
    test_array_roundtrip();
    test_tag_attributes_roundtrip();
    test_negative_int_roundtrip();

    printf("\n=== Summary ===\n");
    printf("Tests passed: %d\n", tests_passed);
    printf("Tests failed: %d\n", tests_failed);

    return tests_failed > 0 ? 1 : 0;
}