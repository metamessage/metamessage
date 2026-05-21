/**
 * nested_structures.cpp
 *
 * MetaMessage C++ API — 嵌套结构与高级用法示例
 * 展示：嵌套对象、数组嵌套对象、Tag 解析与合并
 *
 * Compile:
 *   cd mm-cpp && g++ -std=c++17 -I src -o ../examples/cpp/nested_structures \
 *     ../examples/cpp/nested_structures.cpp -lm
 */

#include "mm/mm.hpp"
#include "ir/tag.hpp"
#include <iostream>
#include <cassert>

int main() {
    using namespace mmc;

    std::cout << "=== MetaMessage C++ API — 嵌套结构与高级用法 ===\n\n";

    /* ---- 1. 嵌套对象 ---- */
    std::cout << "-- 1. 嵌套对象 --\n";

    /* 地址 */
    auto address = ir::makeObject();
    address->tag.desc = "地址信息";

    auto street = ir::makeValue();
    street->text = "123 Main St";
    street->tag.type = ir::ValueType::Str;
    street->tag.desc = "街道";
    street->tag.min = "1";
    street->tag.max = "200";

    auto city = ir::makeValue();
    city->text = "Beijing";
    city->tag.type = ir::ValueType::Str;
    city->tag.desc = "城市";

    address->fields.emplace_back("street", street);
    address->fields.emplace_back("city", city);

    /* 人员 */
    auto person = ir::makeObject();
    person->tag.desc = "人员信息";

    auto name = ir::makeValue();
    name->text = "Alice";
    name->tag.type = ir::ValueType::Str;
    name->tag.desc = "姓名";

    person->fields.emplace_back("name", name);
    person->fields.emplace_back("address", address);

    /* 编码/解码 */
    core::Encoder enc1;
    auto encoded = enc1.encode(person);
    core::Decoder dec1;
    auto decoded = dec1.decode(encoded);
    auto obj1 = std::static_pointer_cast<ir::Object>(decoded);

    std::cout << "  人员对象: " << obj1->fields.size() << " 个字段\n";
    for (const auto& f : obj1->fields) {
        if (f.value->getType() == ir::NodeType::Object) {
            auto inner = std::static_pointer_cast<ir::Object>(f.value);
            std::cout << "    " << f.key << " (嵌套对象, "
                      << inner->fields.size() << " 字段):\n";
            for (const auto& innerF : inner->fields) {
                auto v = std::static_pointer_cast<ir::Value>(innerF.value);
                if (v) std::cout << "      " << innerF.key << " = " << v->text << "\n";
            }
        } else {
            auto v = std::static_pointer_cast<ir::Value>(f.value);
            if (v) std::cout << "    " << f.key << " = " << v->text << "\n";
        }
    }

    /* ---- 2. 对象数组 ---- */
    std::cout << "\n-- 2. 对象数组 --\n";

    auto people = ir::makeArray();
    people->tag.desc = "人员列表";

    auto makePerson = [](const std::string& n, int a) -> std::shared_ptr<ir::Object> {
        auto p = ir::makeObject();
        auto vn = ir::makeValue(); vn->text = n; vn->tag.type = ir::ValueType::Str;
        auto va = ir::makeValue(); va->text = std::to_string(a); va->tag.type = ir::ValueType::U8;
        p->fields.emplace_back("name", vn);
        p->fields.emplace_back("age", va);
        return p;
    };

    people->items.push_back(makePerson("Alice", 30));
    people->items.push_back(makePerson("Bob", 25));
    people->items.push_back(makePerson("Charlie", 35));

    std::cout << "  数组有 " << people->items.size() << " 个元素\n";
    for (size_t i = 0; i < people->items.size(); i++) {
        auto p = std::static_pointer_cast<ir::Object>(people->items[i]);
        auto vn = std::static_pointer_cast<ir::Value>(p->fields[0].value);
        auto va = std::static_pointer_cast<ir::Value>(p->fields[1].value);
        std::cout << "    [" << i << "] " << vn->text << ", age=" << va->text << "\n";
    }

    /* ---- 3. Tag 解析与合并 ---- */
    std::cout << "\n-- 3. Tag 解析与合并 --\n";

    /* 从 mm: 注释字符串解析 Tag */
    auto parsedTag = ir::Tag::parse("// mm: desc=\"person\"; min=0; max=200; nullable");
    std::cout << "  解析 Tag:\n";
    std::cout << "    desc=\"" << parsedTag.desc << "\"\n";
    std::cout << "    min=" << parsedTag.min << "\n";
    std::cout << "    max=" << parsedTag.max << "\n";
    std::cout << "    nullable=" << (parsedTag.nullable ? "true" : "false") << "\n";

    /* Tag 合并 */
    ir::Tag base, override;
    base.desc = "base";
    base.nullable = false;
    base.min = "0";

    override.desc = "override";
    override.nullable = true;

    ir::Tag merged = ir::mergeTag(&base, &override);
    std::cout << "\n  Tag 合并结果:\n";
    std::cout << "    desc=\"" << merged.desc << "\" (覆盖)\n";
    std::cout << "    nullable=" << (merged.nullable ? "true" : "false") << " (覆盖)\n";
    std::cout << "    min=" << merged.min << " (继承)\n";

    /* Tag 序列化 */
    ir::Tag tag = ir::Tag::create();
    tag.type = ir::ValueType::U8;
    tag.desc = "年龄";
    tag.min = "0";
    tag.max = "150";
    tag.nullable = false;

    std::string tagStr = tag.toString();
    std::cout << "\n  Tag toString: " << tagStr << "\n";

    /* Tag 的 child 属性继承 */
    std::cout << "\n-- 4. Child 属性继承 --\n";

    ir::Tag arrayTag = ir::Tag::create();
    arrayTag.childDesc = "分数";
    arrayTag.childType = ir::ValueType::I;
    arrayTag.childMin = "0";
    arrayTag.childMax = "100";
    arrayTag.childNullable = true;

    ir::Tag fieldTag = ir::Tag::create();
    fieldTag.inherit(arrayTag);

    std::cout << "  child 属性继承后:\n";
    std::cout << "    desc=\"" << fieldTag.desc << "\"\n";
    std::cout << "    type=" << ir::valueTypeToString(fieldTag.type) << "\n";
    std::cout << "    min=" << fieldTag.min << "\n";
    std::cout << "    max=" << fieldTag.max << "\n";
    std::cout << "    nullable=" << (fieldTag.nullable ? "true" : "false") << "\n";

    /* ---- 5. Tag 二进制编码 ---- */
    std::cout << "\n-- 5. Tag 二进制编码 --\n";

    auto tagBytes = tag.bytes();
    std::cout << "  tag 编码后: " << tagBytes.size() << " 字节\n";
    for (size_t i = 0; i < tagBytes.size() && i < 20; i++) {
        printf("    %02x\n", tagBytes[i]);
    }

    std::cout << "\n=== 完成 ===\n";
    return 0;
}
