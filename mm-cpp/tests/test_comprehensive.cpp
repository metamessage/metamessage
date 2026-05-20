#include <iostream>
#include <string>
#include <cassert>
#include <sstream>

#include "ir/value_type.hpp"
#include "ir/tag.hpp"
#include "ir/ast.hpp"
#include "core/constants.hpp"
#include "core/encoder.hpp"
#include "core/decoder.hpp"
#include "jsonc/scanner.hpp"
#include "jsonc/parser.hpp"
#include "jsonc/printer.hpp"
#include "mm/mm.hpp"
#include "mm/macro.hpp"

using namespace mmc;

struct Person {
    std::string name;
    uint8_t age;
};

MM_OBJECT(Person,
    MM_FIELD(name, str, .desc="姓名", .min=1, .max=64),
    MM_FIELD(age, u8, .desc="年龄", .min=0, .max=150)
);

static int testsPassed = 0;
static int testsFailed = 0;

#define TEST(name, expr) do { \
    if (!(expr)) { \
        std::cerr << "FAIL: " << name << std::endl; \
        ++testsFailed; \
    } else { \
        std::cout << "PASS: " << name << std::endl; \
        ++testsPassed; \
    } \
} while(0)

void testValueType() {
    std::cout << "\n=== ValueType Tests ===\n";

    TEST("ValueType::Str to string",
         ir::valueTypeToString(ir::ValueType::Str) == "str");
    TEST("ValueType::I to string",
         ir::valueTypeToString(ir::ValueType::I) == "i");
    TEST("ValueType::F64 to string",
         ir::valueTypeToString(ir::ValueType::F64) == "f64");
    TEST("ValueType::Bool to string",
         ir::valueTypeToString(ir::ValueType::Bool) == "bool");
    TEST("ValueType::Uuid to string",
         ir::valueTypeToString(ir::ValueType::Uuid) == "uuid");
    TEST("ValueType::Datetime to string",
         ir::valueTypeToString(ir::ValueType::Datetime) == "datetime");
    TEST("ValueType::Unknown to string",
         ir::valueTypeToString(ir::ValueType::Unknown) == "unknown");

    TEST("Parse str", ir::parseValueType("str") == ir::ValueType::Str);
    TEST("Parse i32", ir::parseValueType("i32") == ir::ValueType::I32);
    TEST("Parse uuid", ir::parseValueType("uuid") == ir::ValueType::Uuid);
    TEST("Parse UPPERCASE", ir::parseValueType("BOOL") == ir::ValueType::Bool);
    TEST("Parse unknown", ir::parseValueType("invalid") == ir::ValueType::Unknown);
}

void testTag() {
    std::cout << "\n=== Tag Tests ===\n";

    ir::Tag tag = ir::Tag::create();
    tag.type = ir::ValueType::Str;
    tag.desc = "姓名";
    tag.nullable = true;
    tag.min = "1";
    tag.max = "64";

    std::string ts = tag.toString();
    TEST("Tag toString contains type",
         ts.find("type=") == std::string::npos); // Str is default, not printed
    TEST("Tag toString contains desc",
         ts.find("desc") != std::string::npos);
    TEST("Tag toString contains nullable",
         ts.find("nullable") != std::string::npos);
    TEST("Tag toString contains min",
         ts.find("min=1") != std::string::npos);
    TEST("Tag toString contains max",
         ts.find("max=64") != std::string::npos);

    tag.type = ir::ValueType::U8;
    ts = tag.toString();
    TEST("Non-default type is printed",
         ts.find("type=u8") != std::string::npos);

    auto tag2 = ir::Tag::parse(" // mm: desc=\"姓名\" ; min=1 ; max=64");
    TEST("Parse desc", tag2.desc == "姓名");
    TEST("Parse min", tag2.min == "1");
    TEST("Parse max", tag2.max == "64");

    auto tag3 = ir::Tag::parse("mm: nullable ; desc=\"test\"");
    TEST("Parse nullable", tag3.nullable == true);
    TEST("Parse desc from mm: prefix", tag3.desc == "test");

    ir::Tag childTag = ir::Tag::create();
    childTag.childDesc = "child field";
    childTag.childType = ir::ValueType::U32;
    childTag.childMin = "0";
    childTag.childMax = "100";

    ir::Tag fieldTag = ir::Tag::create();
    fieldTag.inherit(childTag);
    TEST("Inherit child desc", fieldTag.desc == "child field");
    TEST("Inherit child type", fieldTag.type == ir::ValueType::U32);
    TEST("Inherit child min", fieldTag.min == "0");
    TEST("Inherit child max", fieldTag.max == "100");

    ir::Tag t1, t2;
    t1.desc = "original";
    t2.desc = "override";
    t1.nullable = false;
    t2.nullable = true;

    ir::Tag merged = ir::mergeTag(&t1, &t2);
    TEST("Merge desc (override)", merged.desc == "override"); // src overrides dst
    TEST("Merge nullable (true)", merged.nullable == true);

    auto encoded = tag.bytes();
    TEST("Tag bytes non-empty", !encoded.empty());
}

void testAST() {
    std::cout << "\n=== AST Tests ===\n";

    auto obj = ir::makeObject();
    TEST("Object type", obj->getType() == ir::NodeType::Object);
    TEST("Object tag not null", obj->getTag() != nullptr);

    auto arr = ir::makeArray();
    TEST("Array type", arr->getType() == ir::NodeType::Array);

    auto val = ir::makeValue();
    TEST("Value type", val->getType() == ir::NodeType::Value);
    val->text = "hello";
    TEST("Value text", val->text == "hello");

    auto doc = ir::makeDoc();
    TEST("Doc type", doc->getType() == ir::NodeType::Doc);

    obj->fields.emplace_back("name", val);
    TEST("Object has 1 field", obj->fields.size() == 1);
    TEST("Field key", obj->fields[0].key == "name");

    ir::Field field("age", ir::makeValue());
    TEST("Field key from constructor", field.key == "age");

    obj->setPath("/root");
    TEST("Path setting", obj->getPath() == "/root");

    val->setPath("/root/name");
    TEST("Value path", val->getPath() == "/root/name");
}

void testJSONCScannerParser() {
    std::cout << "\n=== JSONC Scanner/Parser Tests ===\n";

    std::string input = "{\"name\": \"alice\", \"age\": 30}";
    auto scanner = jsonc::Scanner(input);
    auto tokens = scanner.scanAll();

    TEST("Scanner produces tokens", tokens.size() > 0);
    TEST("First token is LBrace", tokens[0].type == jsonc::TokenType::LBrace);

    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();
    TEST("Parse produces node", node != nullptr);
    TEST("Parsed node is Object", node->getType() == ir::NodeType::Object);
    auto obj = std::static_pointer_cast<ir::Object>(node);
    TEST("Object has 2 fields", obj->fields.size() == 2);

    std::string taggedInput = R"({
        // mm: desc="person object"
        "name": "alice",
        "age": 30  // mm: type=u8; min=0; max=150
    })";
    auto scanner2 = jsonc::Scanner(taggedInput);
    auto tokens2 = scanner2.scanAll();
    auto parser2 = jsonc::Parser(tokens2);
    auto node2 = parser2.parse();
    parser2.applyTags(node2);
    auto obj2 = std::static_pointer_cast<ir::Object>(node2);

    if (obj2 && obj2->fields.size() >= 2) {
        auto ageVal = std::static_pointer_cast<ir::Value>(obj2->fields[1].value);
        TEST("Age field has type from mm tag",
             ageVal->getTag()->type == ir::ValueType::U8);
        TEST("Age field has min from mm tag",
             ageVal->getTag()->min == "0");
        TEST("Age field has max from mm tag",
             ageVal->getTag()->max == "150");
    } else {
        TEST("Tagged parse produces valid object", false);
    }

    std::string arrInput = "[1, 2, 3]";
    auto scanner3 = jsonc::Scanner(arrInput);
    auto tokens3 = scanner3.scanAll();
    auto parser3 = jsonc::Parser(tokens3);
    auto node3 = parser3.parse();
    TEST("Array parse", node3->getType() == ir::NodeType::Array);
    auto arr3 = std::static_pointer_cast<ir::Array>(node3);
    TEST("Array has 3 items", arr3->items.size() == 3);
}

void testJSONCPrinter() {
    std::cout << "\n=== JSONC Printer Tests ===\n";

    std::string input = R"({"name":"alice","age":30})";
    auto scanner = jsonc::Scanner(input);
    auto tokens = scanner.scanAll();
    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();

    std::string output = jsonc::toJSONC(node);
    TEST("Round-trip JSONC", !output.empty());
    TEST("Printer output valid json", output.find("alice") != std::string::npos);
}

void testEncoderDecoder() {
    std::cout << "\n=== Encoder/Decoder Tests ===\n";

    auto obj = ir::makeObject();
    obj->tag.type = ir::ValueType::Obj;

    auto nameVal = ir::makeValue();
    nameVal->text = "alice";
    nameVal->tag.type = ir::ValueType::Str;
    obj->fields.emplace_back("name", nameVal);

    auto ageVal = ir::makeValue();
    ageVal->text = "30";
    ageVal->tag.type = ir::ValueType::U8;
    obj->fields.emplace_back("age", ageVal);

    core::Encoder encoder;
    auto encoded = encoder.encode(obj);
    TEST("Encoder produces bytes", !encoded.empty());
    TEST("Encoded size > 2", encoded.size() > 2);

    core::Decoder decoder;
    auto decoded = decoder.decode(encoded);
    TEST("Decoded node not null", decoded != nullptr);
    TEST("Decoded is Object", decoded->getType() == ir::NodeType::Object);
    auto decodedObj = std::static_pointer_cast<ir::Object>(decoded);
    TEST("Decoded has 2 fields", decodedObj->fields.size() == 2);
    if (decodedObj->fields.size() >= 1) {
        auto nameField = std::static_pointer_cast<ir::Value>(decodedObj->fields[0].value);
        if (nameField) {
            TEST("Decoded name field exists",
                 decodedObj->fields[0].key == "name");
        }
    }
}

void testMacroTagSystem() {
    std::cout << "\n=== Macro Tag System Tests ===\n";

    size_t count = _mm_field_count_Person;
    TEST("Person has 2 fields", count == 2);

    const auto& fields = _mm_fields_Person;
    TEST("First field name is name", std::string(fields[0].name) == "name");
    TEST("First field type is str", std::string(fields[0].type) == "str");
    TEST("First field tag has desc",
         std::string(fields[0].tag.desc) == "姓名");
    TEST("First field tag min",
         fields[0].tag.min == 1);
    TEST("First field tag max",
         fields[0].tag.max == 64);

    TEST("Second field name is age", std::string(fields[1].name) == "age");
    TEST("Second field type is u8", std::string(fields[1].type) == "u8");
    TEST("Second field tag desc",
         std::string(fields[1].tag.desc) == "年龄");
    TEST("Second field tag min",
         fields[1].tag.min == 0);
    TEST("Second field tag max",
         fields[1].tag.max == 150);

    auto tag = _mm_build_field_tag_Person(fields[0]);
    TEST("Field tag has correct type", tag.type == ir::ValueType::Str);
    TEST("Field tag has desc", tag.desc == "姓名");
    TEST("Field tag has min", tag.min == "1");
    TEST("Field tag has max", tag.max == "64");

    auto tag2 = _mm_build_field_tag_Person(fields[1]);
    TEST("Age field tag type U8", tag2.type == ir::ValueType::U8);
    TEST("Age field tag min 0", tag2.min == "0");
    TEST("Age field tag max 150", tag2.max == "150");
}

void testFullRoundTrip() {
    std::cout << "\n=== Full Round Trip Tests ===\n";

    std::string input = R"({"name":"bob","age":25,"active":true})";
    auto scanner = jsonc::Scanner(input);
    auto tokens = scanner.scanAll();
    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();

    core::Encoder encoder;
    auto encoded = encoder.encode(node);

    core::Decoder decoder;
    auto decoded = decoder.decode(encoded);

    std::string output = jsonc::toJSONC(decoded);

    TEST("Round trip produces output", !output.empty());
    TEST("Round trip contains name", output.find("bob") != std::string::npos);
    TEST("Round trip contains age", output.find("25") != std::string::npos);
    TEST("Round trip contains active", output.find("true") != std::string::npos);

    auto encoded2 = mm::fromJSONC(input);
    TEST("mm::fromJSONC produces bytes", !encoded2.empty());
    auto decoded2 = mm::toNode(encoded2);
    TEST("mm::toNode returns node", decoded2 != nullptr);
}

void testValueTypeEncodeDecode() {
    std::cout << "\n=== Value Type Encode/Decode Tests ===\n";

    auto boolVal = ir::makeValue();
    boolVal->text = "true";
    boolVal->tag.type = ir::ValueType::Bool;
    core::Encoder enc;
    auto encoded = enc.encode(boolVal);
    core::Decoder dec;
    auto decoded = dec.decode(encoded);
    TEST("Bool round-trip", decoded != nullptr);
    if (decoded) {
        auto v = std::static_pointer_cast<ir::Value>(decoded);
        if (v) TEST("Bool value true", v->text == "true");
    }

    auto intVal = ir::makeValue();
    intVal->text = "42";
    intVal->tag.type = ir::ValueType::I;
    core::Encoder enc2;
    auto encoded2 = enc2.encode(intVal);
    core::Decoder dec2;
    auto decoded2 = dec2.decode(encoded2);
    TEST("Int round-trip", decoded2 != nullptr);
    if (decoded2) {
        auto v = std::static_pointer_cast<ir::Value>(decoded2);
        if (v) TEST("Int value 42", v->text == "42");
    }

    auto strVal = ir::makeValue();
    strVal->text = "hello world";
    strVal->tag.type = ir::ValueType::Str;
    core::Encoder enc3;
    auto encoded3 = enc3.encode(strVal);
    core::Decoder dec3;
    auto decoded3 = dec3.decode(encoded3);
    TEST("String round-trip", decoded3 != nullptr);
    if (decoded3) {
        auto v = std::static_pointer_cast<ir::Value>(decoded3);
        if (v) TEST("String value hello world", v->text == "hello world");
    }
}

void testNestedObject() {
    std::cout << "\n=== Nested Object Tests ===\n";

    auto outer = ir::makeObject();
    auto inner = ir::makeObject();
    inner->tag.type = ir::ValueType::Obj;

    auto val = ir::makeValue();
    val->text = "inner_value";
    val->tag.type = ir::ValueType::Str;
    inner->fields.emplace_back("inner_field", val);

    outer->fields.emplace_back("nested", inner);

    auto val2 = ir::makeValue();
    val2->text = "outer_value";
    val2->tag.type = ir::ValueType::Str;
    outer->fields.emplace_back("outer_field", val2);

    core::Encoder enc;
    auto encoded = enc.encode(outer);
    TEST("Nested encoder produces bytes", !encoded.empty());

    core::Decoder dec;
    auto decoded = dec.decode(encoded);
    TEST("Nested decode produces node", decoded != nullptr);
    auto obj = std::static_pointer_cast<ir::Object>(decoded);
    TEST("Nested has 2 fields", obj->fields.size() == 2);

    if (obj->fields.size() >= 2) {
        auto nestedField = obj->fields[0].value;
        TEST("First field is nested object",
             nestedField->getType() == ir::NodeType::Object);
        if (nestedField->getType() == ir::NodeType::Object) {
            auto nestedObj = std::static_pointer_cast<ir::Object>(nestedField);
            TEST("Nested object has 1 field", nestedObj->fields.size() == 1);
        }
    }
}

void testEncodingConstants() {
    std::cout << "\n=== Encoding Constants Tests ===\n";

    TEST("Prefix Simple", static_cast<uint8_t>(core::Prefix::Simple) == 0b00000000);
    TEST("Prefix PositiveInt", static_cast<uint8_t>(core::Prefix::PositiveInt) == 0b00100000);
    TEST("Prefix NegativeInt", static_cast<uint8_t>(core::Prefix::NegativeInt) == 0b01000000);
    TEST("Prefix Float", static_cast<uint8_t>(core::Prefix::Float) == 0b01100000);
    TEST("Prefix String", static_cast<uint8_t>(core::Prefix::String) == 0b10000000);
    TEST("Prefix Bytes", static_cast<uint8_t>(core::Prefix::Bytes) == 0b10100000);
    TEST("Prefix Container", static_cast<uint8_t>(core::Prefix::Container) == 0b11000000);
    TEST("Prefix Tag", static_cast<uint8_t>(core::Prefix::Tag) == 0b11100000);

    TEST("getSuffix extracts lower bits",
         core::getSuffix(0b11111111) == 0b00011111);
    TEST("getPrefix extracts upper bits",
         core::getPrefix(0b11111111) == core::Prefix::Tag);

    TEST("stringInlineLen for small values",
         core::stringInlineLen(0b10000011) == 3);
    TEST("stringExtraLen for len 31",
         core::stringExtraLen(0b10011110) == 1);

    TEST("isArrayContainer true",
         core::isArrayContainer(0b11010000) == true);
    TEST("isArrayContainer false for object",
         core::isArrayContainer(0b11000000) == false);
}

int main() {
    std::cout << "MM-C++ Comprehensive Tests\n";
    std::cout << "==========================\n";

    testEncodingConstants();
    testValueType();
    testTag();
    testAST();
    testJSONCScannerParser();
    testJSONCPrinter();
    testEncoderDecoder();
    testMacroTagSystem();
    testFullRoundTrip();
    testValueTypeEncodeDecode();
    testNestedObject();

    std::cout << "\n==========================\n";
    std::cout << "Tests passed: " << testsPassed << std::endl;
    std::cout << "Tests failed: " << testsFailed << std::endl;

    if (testsFailed > 0) {
        std::cerr << "\n*** SOME TESTS FAILED ***\n";
        return 1;
    }
    std::cout << "\n*** ALL TESTS PASSED ***\n";
    return 0;
}