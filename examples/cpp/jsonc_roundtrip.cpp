/**
 * jsonc_roundtrip.cpp
 *
 * MetaMessage C++ API — JSONC 往返示例
 * 展示 JSONC Scanner / Parser / Printer 的完整生命周期，
 * 以及带 mm: 注释标签的解析。
 *
 * Compile:
 *   cd mm-cpp && g++ -std=c++17 -I src -o ../examples/cpp/jsonc_roundtrip \
 *     ../examples/cpp/jsonc_roundtrip.cpp src/jsonc/scanner.cpp -lm
 */

#include "mm/mm.hpp"
#include "jsonc/scanner.hpp"
#include "jsonc/parser.hpp"
#include "jsonc/printer.hpp"
#include <iostream>
#include <cassert>

int main() {
    using namespace mmc;

    std::cout << "=== MetaMessage C++ API — JSONC 往返 ===\n\n";

    /* ---- 1. 扫描 JSONC 字符串 ---- */
    std::cout << "-- 1. 扫描 JSON --\n";

    std::string json = R"({"name":"Alice","age":30,"active":true})";
    auto scanner = jsonc::Scanner(json);
    auto tokens = scanner.scanAll();

    std::cout << "  扫描得到 " << tokens.size() << " 个 token\n";
    for (size_t i = 0; i < tokens.size(); i++) {
        std::cout << "    [" << i << "] type=" << static_cast<int>(tokens[i].type)
                  << " value=\"" << tokens[i].value << "\"\n";
    }

    /* ---- 2. 解析为 AST ---- */
    std::cout << "\n-- 2. 解析为 AST --\n";

    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();
    assert(node != nullptr);
    assert(node->getType() == ir::NodeType::Object);

    auto obj = std::static_pointer_cast<ir::Object>(node);
    std::cout << "  解析对象: " << obj->fields.size() << " 个字段\n";
    for (const auto& f : obj->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        std::cout << "    " << f.key << " => " << (v ? v->text : "?") << "\n";
    }

    /* ---- 3. 打印回 JSONC ---- */
    std::cout << "\n-- 3. 打印回 JSONC --\n";

    std::string output = jsonc::toJSONC(node);
    std::cout << output << "\n";

    /* ---- 4. 带 mm: 注释标签的 JSONC ---- */
    std::cout << "\n-- 4. 带 mm: 注释标签的 JSONC --\n";

    std::string taggedJsonc = R"({
        // mm: desc="用户信息"
        "name": "Bob",
        "age": 25,  // mm: type=u8; min=0; max=150
        "email": "bob@example.com"  // mm: type=email
    })";

    auto scanner2 = jsonc::Scanner(taggedJsonc);
    auto tokens2 = scanner2.scanAll();
    auto parser2 = jsonc::Parser(tokens2);
    auto node2 = parser2.parse();
    parser2.applyTags(node2);  // 将 mm: 注释应用到节点 tag

    auto obj2 = std::static_pointer_cast<ir::Object>(node2);
    std::cout << "  解析带标签的 JSONC: " << obj2->fields.size() << " 个字段\n";

    for (const auto& f : obj2->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        if (v) {
            std::cout << "    " << f.key << " = " << v->text
                      << "  desc=\"" << v->tag.desc
                      << "\" type=" << ir::valueTypeToString(v->tag.type);
            if (!v->tag.min.empty()) std::cout << " min=" << v->tag.min;
            if (!v->tag.max.empty()) std::cout << " max=" << v->tag.max;
            std::cout << "\n";
        }
    }

    /* ---- 5. 使用高层 API（mm::parseTaggedJSONC） ---- */
    std::cout << "\n-- 5. 高层 API mm::parseTaggedJSONC --\n";

    auto tagged = mm::parseTaggedJSONC(R"({
        // mm: desc="product"
        "id": 1001,  // mm: type=u64; min=1
        "price": 29.99  // mm: type=f64; desc="价格"
    })");

    auto obj3 = std::static_pointer_cast<ir::Object>(tagged);
    for (const auto& f : obj3->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        if (v) {
            std::cout << "    " << f.key << " = " << v->text
                      << "  type=" << ir::valueTypeToString(v->tag.type)
                      << " desc=\"" << v->tag.desc << "\"\n";
        }
    }

    std::cout << "\n=== 完成 ===\n";
    return 0;
}
