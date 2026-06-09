#ifndef MMCPP_MM_MM_HPP
#define MMCPP_MM_MM_HPP

#include "../core/decoder.hpp"
#include "../core/encoder.hpp"
#include "../ir/ast.hpp"
#include "../ir/tag.hpp"
#include "../ir/value_type.hpp"
#include "../jsonc/parser.hpp"
#include "../jsonc/printer.hpp"
#include <memory>
#include <string>
#include <vector>

namespace mmc {
namespace mm {

inline std::vector<uint8_t> fromJSONC(const std::string &s) {
  auto scanner = jsonc::Scanner(s);
  auto tokens = scanner.scanAll();
  auto parser = jsonc::Parser(tokens);
  auto node = parser.parse();
  parser.applyTags(node);

  core::Encoder encoder;
  return encoder.encode(node);
}

inline std::string toJSONC(const std::vector<uint8_t> &data) {
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

inline std::shared_ptr<ir::Node> toNode(const std::vector<uint8_t> &data) {
  core::Decoder decoder;
  return decoder.decode(data);
}

inline std::shared_ptr<ir::Node> parseJSONC(const std::string &s) {
  auto scanner = jsonc::Scanner(s);
  auto tokens = scanner.scanAll();
  auto parser = jsonc::Parser(tokens);
  auto node = parser.parse();
  parser.applyTags(node);
  return node;
}

inline std::shared_ptr<ir::Node> parseTaggedJSONC(const std::string &s) {
  auto scanner = jsonc::Scanner(s);
  auto tokens = scanner.scanAll();
  auto parser = jsonc::Parser(tokens);
  auto node = parser.parse();
  parser.applyTags(node);
  return node;
}

inline std::vector<uint8_t> encodeFromValue(std::shared_ptr<ir::Node> value,
                                            const std::string &tag = "") {
  if (!tag.empty()) {
    auto parsed_tag = ir::Tag::parse(tag);
    if (auto t = value->getTag()) {
      *t = ir::mergeTag(t, &parsed_tag);
    }
  }
  core::Encoder encoder;
  return encoder.encode(value);
}

inline std::vector<uint8_t> encodeFromJsonc(const std::string &jsonc) {
  return fromJSONC(jsonc);
}

inline std::shared_ptr<ir::Node> decodeToValue(const std::vector<uint8_t> &data) {
  return toNode(data);
}

inline std::string decodeToJsonc(const std::vector<uint8_t> &data) {
  return toJSONC(data);
}

inline std::string valueToJsonc(std::shared_ptr<ir::Node> value,
                                const std::string &tag = "") {
  if (!tag.empty()) {
    auto parsed_tag = ir::Tag::parse(tag);
    if (auto t = value->getTag()) {
      *t = ir::mergeTag(t, &parsed_tag);
    }
  }
  return toJSONCFromNode(value);
}

inline std::shared_ptr<ir::Node> jsoncToValue(const std::string &jsonc) {
  return parseJSONC(jsonc);
}

} // namespace mm
} // namespace mmc

#endif