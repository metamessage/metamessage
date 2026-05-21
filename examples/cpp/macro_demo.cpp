/**
 * macro_demo.cpp
 *
 * MetaMessage C++ API — MM_OBJECT / MM_FIELD 宏声明式示例
 * 展示如何使用声明式宏自动生成 Schema 定义和序列化逻辑。
 *
 * Compile:
 *   cd mm-cpp && g++ -std=c++17 -I src -o ../examples/cpp/macro_demo \
 *     ../examples/cpp/macro_demo.cpp -lm
 */

#include "mm/mm.hpp"
#include "mm/macro.hpp"
#include <iostream>
#include <string>
#include <cassert>

/* ========== 声明式 Schema 定义 ========== */

struct Person {
    std::string name;
    uint8_t age;
    std::string email;
    bool isActive;
};

MM_OBJECT(Person,
    MM_FIELD(name,     str,   .desc="姓名",    .min=1, .max=64),
    MM_FIELD(age,      u8,    .desc="年龄",    .min=0, .max=150),
    MM_FIELD(email,    email, .desc="电子邮箱"),
    MM_FIELD(isActive, bool,  .desc="是否激活")
);

struct Product {
    int64_t id;
    std::string title;
    double price;
};

MM_OBJECT(Product,
    MM_FIELD(id,    i64,  .desc="商品ID",   .min=1),
    MM_FIELD(title, str,  .desc="商品名称",  .min=1, .max=200),
    MM_FIELD(price, f64,  .desc="价格",     .min=0)
);

/* ======================================= */

int main() {
    using namespace mmc;

    std::cout << "=== MetaMessage C++ API — MM_OBJECT 宏声明式用法 ===\n\n";

    /* ---- 1. 检查自动生成的字段元信息 ---- */
    std::cout << "-- 1. 字段元信息 --\n";

    std::cout << "  Person 有 " << _mm_field_count_Person << " 个字段:\n";
    for (size_t i = 0; i < _mm_field_count_Person; i++) {
        const auto& fd = _mm_fields_Person[i];
        std::cout << "    [" << i << "] " << fd.name
                  << " (type=" << fd.type
                  << ", desc=" << fd.tag.desc
                  << ", min=" << fd.tag.min
                  << ", max=" << fd.tag.max << ")\n";
    }

    std::cout << "\n  Product 有 " << _mm_field_count_Product << " 个字段:\n";
    for (size_t i = 0; i < _mm_field_count_Product; i++) {
        const auto& fd = _mm_fields_Product[i];
        std::cout << "    [" << i << "] " << fd.name
                  << " (type=" << fd.type
                  << ", desc=" << fd.tag.desc << ")\n";
    }

    /* ---- 2. 从 FieldDescriptor 构建 Tag ---- */
    std::cout << "\n-- 2. 构建 Tag --\n";

    auto tag = _mm_build_field_tag_Person(_mm_fields_Person[0]);  // name
    assert(tag.type == ir::ValueType::Str);
    std::cout << "  name tag: type=" << ir::valueTypeToString(tag.type)
              << " desc=\"" << tag.desc << "\" min=" << tag.min
              << " max=" << tag.max << "\n";

    auto tag2 = _mm_build_field_tag_Product(_mm_fields_Product[2]);  // price
    assert(tag2.type == ir::ValueType::F64);
    std::cout << "  price tag: type=" << ir::valueTypeToString(tag2.type)
              << " desc=\"" << tag2.desc << "\" min=" << tag2.min << "\n";

    /* ---- 3. 用 _mm_to_node 将 struct 转为 AST 节点 ---- */
    std::cout << "\n-- 3. struct → AST 节点（_mm_to_node_XXX） --\n";

    Person alice;
    alice.name = "Alice";
    alice.age = 28;
    alice.email = "alice@example.com";
    alice.isActive = true;

    auto personNode = _mm_to_node_Person(alice);
    assert(personNode != nullptr);
    assert(personNode->getType() == ir::NodeType::Object);

    std::cout << "  Person 节点: " << personNode->fields.size() << " 个字段\n";
    for (const auto& f : personNode->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        if (v) {
            std::cout << "    " << f.key << " = " << v->text
                      << "  (type=" << ir::valueTypeToString(v->tag.type)
                      << ", desc=\"" << v->tag.desc << "\")\n";
        }
    }

    /* ---- 4. 编码/解码 ---- */
    std::cout << "\n-- 4. 二进制编码/解码 --\n";

    core::Encoder enc;
    auto encoded = enc.encode(personNode);
    std::cout << "  编码后: " << encoded.size() << " 字节\n";

    core::Decoder dec;
    auto decoded = dec.decode(encoded);
    assert(decoded != nullptr);

    std::string jsonc = mm::toJSONCFromNode(decoded);
    std::cout << "  解码后 JSONC:\n" << jsonc << "\n";

    /* ---- 5. 从 AST 节点恢复 struct ---- */
    std::cout << "\n-- 5. 节点层级字段遍历 --\n";

    auto obj = std::static_pointer_cast<ir::Object>(decoded);
    for (const auto& f : obj->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        if (v) {
            std::cout << "    " << f.key << " => " << v->text << "\n";
        }
    }

    /* ---- 6. 给根节点添加 desc ---- */
    personNode->tag.desc = "人员信息";
    std::cout << "\n  根节点 desc: \"" << personNode->tag.desc << "\"\n";

    std::cout << "\n=== 完成 ===\n";
    return 0;
}
