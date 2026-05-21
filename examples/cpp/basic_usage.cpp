/**
 * basic_usage.cpp
 *
 * MetaMessage C++ API — 基础用法示例
 * 涵盖：创建值、构建对象/数组、二进制编码解码
 *
 * Compile:
 *   cd mm-cpp && g++ -std=c++17 -I src -o ../examples/cpp/basic_usage \
 *     ../examples/cpp/basic_usage.cpp -lm
 */

#include "mm/mm.hpp"
#include <iostream>
#include <cassert>

int main() {
    using namespace mmc;

    std::cout << "=== MetaMessage C++ API — 基础用法 ===\n\n";

    /* ---- 1. 创建值 ---- */
    std::cout << "-- 创建值 --\n";

    auto nameVal = ir::makeValue();
    nameVal->text = "Alice";
    nameVal->tag.type = ir::ValueType::Str;
    nameVal->tag.desc = "姓名";
    nameVal->tag.min = "1";
    nameVal->tag.max = "64";

    auto ageVal = ir::makeValue();
    ageVal->text = "30";
    ageVal->tag.type = ir::ValueType::U8;
    ageVal->tag.desc = "年龄";
    ageVal->tag.min = "0";
    ageVal->tag.max = "150";

    auto scoreVal = ir::makeValue();
    scoreVal->text = "95.5";
    scoreVal->tag.type = ir::ValueType::F64;
    scoreVal->tag.desc = "分数";

    auto flagVal = ir::makeValue();
    flagVal->text = "true";
    flagVal->tag.type = ir::ValueType::Bool;
    flagVal->tag.desc = "激活状态";

    std::cout << "  name  = " << nameVal->text << " (type: "
              << ir::valueTypeToString(nameVal->tag.type) << ")\n";
    std::cout << "  age   = " << ageVal->text << "\n";
    std::cout << "  score = " << scoreVal->text << "\n";
    std::cout << "  flag  = " << flagVal->text << "\n";

    /* ---- 2. 构建对象 ---- */
    std::cout << "\n-- 构建对象 --\n";

    auto person = ir::makeObject();
    person->fields.emplace_back("name", nameVal);
    person->fields.emplace_back("age", ageVal);
    person->fields.emplace_back("score", scoreVal);
    person->fields.emplace_back("active", flagVal);

    std::cout << "  person 对象有 " << person->fields.size() << " 个字段\n";
    for (const auto& f : person->fields) {
        auto v = std::static_pointer_cast<ir::Value>(f.value);
        std::cout << "    " << f.key << " => " << (v ? v->text : "?") << "\n";
    }

    /* ---- 3. 构建数组 ---- */
    std::cout << "\n-- 构建数组 --\n";

    auto scores = ir::makeArray();
    auto s1 = ir::makeValue(); s1->text = "95"; s1->tag.type = ir::ValueType::I; s1->tag.desc = "语文";
    auto s2 = ir::makeValue(); s2->text = "87"; s2->tag.type = ir::ValueType::I; s2->tag.desc = "数学";
    auto s3 = ir::makeValue(); s3->text = "92"; s3->tag.type = ir::ValueType::I; s3->tag.desc = "英语";
    scores->items.push_back(s1);
    scores->items.push_back(s2);
    scores->items.push_back(s3);

    std::cout << "  scores 数组有 " << scores->items.size() << " 个元素\n";
    for (size_t i = 0; i < scores->items.size(); i++) {
        auto v = std::static_pointer_cast<ir::Value>(scores->items[i]);
        std::cout << "    [" << i << "] " << (v ? v->text : "?") << "\n";
    }

    /* ---- 4. 二进制编码 / 解码 ---- */
    std::cout << "\n-- 二进制编码/解码 --\n";

    core::Encoder encoder;
    auto encoded = encoder.encode(person);
    std::cout << "  编码后: " << encoded.size() << " 字节\n";

    core::Decoder decoder;
    auto decoded = decoder.decode(encoded);
    assert(decoded != nullptr);
    assert(decoded->getType() == ir::NodeType::Object);
    auto decodedObj = std::static_pointer_cast<ir::Object>(decoded);
    std::cout << "  解码成功: 对象有 " << decodedObj->fields.size() << " 个字段\n";

    /* ---- 5. JSONC 序列化 ---- */
    std::cout << "\n-- JSONC 序列化 --\n";

    std::string jsonc = mm::toJSONCFromNode(person);
    std::cout << jsonc << "\n";

    std::cout << "\n=== 完成 ===\n";
    return 0;
}
