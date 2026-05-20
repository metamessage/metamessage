#ifndef MMCPP_MM_MM_HPP
#define MMCPP_MM_MM_HPP

#include "../ir/ast.hpp"
#include "../ir/value_type.hpp"
#include "../ir/tag.hpp"
#include "../core/encoder.hpp"
#include "../core/decoder.hpp"
#include "../jsonc/printer.hpp"
#include "../jsonc/parser.hpp"
#include <string>
#include <vector>
#include <memory>

namespace mmc {
namespace mm {

inline std::vector<uint8_t> fromJSONC(const std::string& s) {
    auto scanner = jsonc::Scanner(s);
    auto tokens = scanner.scanAll();
    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();

    core::Encoder encoder;
    return encoder.encode(node);
}

inline std::string toJSONC(const std::vector<uint8_t>& data) {
    core::Decoder decoder;
    auto node = decoder.decode(data);
    return jsonc::toJSONC(node);
}

inline std::string toJSONCFromNode(std::shared_ptr<ir::Node> node) {
    return jsonc::toJSONC(node);
}

inline std::vector<uint8_t> fromNode(std::shared_ptr<ir::Node> node) {
    core::Encoder encoder;
    return encoder.encode(node);
}

inline std::shared_ptr<ir::Node> toNode(const std::vector<uint8_t>& data) {
    core::Decoder decoder;
    return decoder.decode(data);
}

inline std::shared_ptr<ir::Node> parseJSONC(const std::string& s) {
    auto scanner = jsonc::Scanner(s);
    auto tokens = scanner.scanAll();
    auto parser = jsonc::Parser(tokens);
    return parser.parse();
}

inline std::shared_ptr<ir::Node> parseTaggedJSONC(const std::string& s) {
    auto scanner = jsonc::Scanner(s);
    auto tokens = scanner.scanAll();
    auto parser = jsonc::Parser(tokens);
    auto node = parser.parse();
    parser.applyTags(node);
    return node;
}

} // namespace mm
} // namespace mmc

#endif