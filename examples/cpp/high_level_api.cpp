/**
 * high_level_api.cpp
 *
 * MetaMessage C++ API — 高层便捷 API 示例
 * 展示 mm::fromJSONC / mm::toJSONC / mm::parseJSONC / mm::fromNode / mm::toNode 等
 *
 * Compile:
 *   cd mm-cpp && g++ -std=c++17 -I src -o ../examples/cpp/high_level_api \
 *     ../examples/cpp/high_level_api.cpp -lm
 */

#include "mm/mm.hpp"
#include <iostream>
#include <cassert>

int main() {
    using namespace mmc;

    std::cout << "=== MetaMessage C++ API — 高层便捷 API ===\n\n";

    /* ---- 1. JSONC → 二进制（mm::fromJSONC） ---- */
    std::cout << "-- 1. JSONC → 二进制 (mm::fromJSONC) --\n";

    std::string json = R"({"name":"Alice","age":30,"active":true})";
    auto bytes = mm::fromJSONC(json);

    std::cout << "  输入 JSONC:  " << json << "\n";
    std::cout << "  输出二进制:   " << bytes.size() << " 字节\n";

    /* ---- 2. 二进制 → JSONC（mm::toJSONC） ---- */
    std::cout << "\n-- 2. 二进制 → JSONC (mm::toJSONC) --\n";

    std::string jsonc = mm::toJSONC(bytes);
    std::cout << "  解码 JSONC:  " << jsonc << "\n";

    /* ---- 3. JSONC → AST 节点（mm::parseJSONC） ---- */
    std::cout << "\n-- 3. JSONC → AST 节点 (mm::parseJSONC) --\n";

    auto node = mm::parseJSONC(R"({"id":42,"score":95.5})");
    assert(node != nullptr);
    assert(node->getType() == ir::NodeType::Object);

    auto obj = std::static_pointer_cast<ir::Object>(node);
    for (const auto& f : obj->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        std::cout << "    " << f.key << " => " << (v ? v->text : "?") << "\n";
    }

    /* ---- 4. AST 节点 → 二进制（mm::fromNode） ---- */
    std::cout << "\n-- 4. AST 节点 → 二进制 (mm::fromNode) --\n";

    auto encoded = mm::fromNode(node);
    std::cout << "  编码后: " << encoded.size() << " 字节\n";

    /* ---- 5. 二进制 → AST 节点（mm::toNode） ---- */
    std::cout << "\n-- 5. 二进制 → AST 节点 (mm::toNode) --\n";

    auto decoded = mm::toNode(encoded);
    assert(decoded != nullptr);

    std::string jsonc2 = mm::toJSONCFromNode(decoded);
    std::cout << "  round-trip 结果: " << jsonc2 << "\n";

    /* ---- 6. 带 mm: 注释的 JSONC ---- */
    std::cout << "\n-- 6. 解析带 mm: 注释的 JSONC (mm::parseTaggedJSONC) --\n";

    auto tagged = mm::parseTaggedJSONC(R"({
        // mm: desc="student"
        "name": "Bob",  // mm: type=str; min=1; max=50
        "age": 22       // mm: type=u8; min=0; max=150
    })");

    auto obj2 = std::static_pointer_cast<ir::Object>(tagged);
    assert(obj2->tag.desc == "student");
    std::cout << "  根节点 desc: \"" << obj2->tag.desc << "\"\n";

    for (const auto& f : obj2->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        if (v) {
            std::cout << "    " << f.key << " = " << v->text
                      << "  type=" << ir::valueTypeToString(v->tag.type)
                      << " min=" << v->tag.min
                      << " max=" << v->tag.max << "\n";
        }
    }

    /* ---- 7. 完整 pipeline ---- */
    std::cout << "\n-- 7. 完整 pipeline --\n";

    std::string pipeline = R"({"msg":"Hello MetaMessage","count":5})";
    auto bin    = mm::fromJSONC(pipeline);
    auto ast    = mm::toNode(bin);
    auto text   = mm::toJSONCFromNode(ast);

    std::cout << "  JSONC  in:  " << pipeline << "\n";
    std::cout << "  Binary:     " << bin.size() << " bytes\n";
    std::cout << "  JSONC out:  " << text << "\n";
    assert(pipeline == text);

    std::cout << "\n=== 完成 ===\n";
    return 0;
}
