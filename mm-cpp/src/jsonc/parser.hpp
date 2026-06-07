#ifndef MMCPP_JSONC_PARSER_HPP
#define MMCPP_JSONC_PARSER_HPP

#include "../ir/ast.hpp"
#include "../ir/tag.hpp"
#include "token.hpp"
#include <cctype>
#include <cstdlib>
#include <memory>
#include <optional>
#include <stdexcept>
#include <string>
#include <vector>

namespace mmc {
namespace jsonc {

class Parser {
public:
  explicit Parser(const std::vector<Token> &tokens) : tokens_(tokens), pos_(0) {
    skipComments();
  }

  std::shared_ptr<ir::Node> parse() {
    if (pos_ >= tokens_.size())
      return nullptr;
    auto &tok = tokens_[pos_];
    switch (tok.type) {
    case TokenType::LBrace:
      return parseObject();
    case TokenType::LBracket:
      return parseArray();
    case TokenType::String:
    case TokenType::Number:
    case TokenType::True:
    case TokenType::False:
    case TokenType::Null:
      return parsePrimitive();
    default:
      return nullptr;
    }
  }

  void applyTags(std::shared_ptr<ir::Node> node) {
    if (!node)
      return;
    switch (node->getType()) {
    case ir::NodeType::NodeObject:
      applyTagsToObject(std::static_pointer_cast<ir::NodeObject>(node));
      break;
    case ir::NodeType::NodeArray:
      applyTagsToArray(std::static_pointer_cast<ir::NodeArray>(node));
      break;
    default:
      break;
    }
  }

private:
  const std::vector<Token> &tokens_;
  size_t pos_;
  std::optional<ir::Tag> pendingTag_;

  void skipComments() {
    size_t lastLine = 0;
    while (pos_ < tokens_.size()) {
      auto &tok = tokens_[pos_];
      if (tok.type == TokenType::Comment) {
        if (lastLine > 0 && tok.line - lastLine > 1) {
          pendingTag_.reset();
        }
        lastLine = tok.line;

        std::string comment = tok.literal;
        auto it = comment.find("mm:");
        if (it != std::string::npos) {
          std::string tagStr = comment.substr(it + 3);
          auto tag = ir::Tag::parse(tagStr);
          if (!pendingTag_.has_value()) {
            pendingTag_ = tag;
          } else {
            pendingTag_ = ir::mergeTag(
                pendingTag_.has_value() ? &pendingTag_.value() : nullptr, &tag);
          }
        }
        ++pos_;
      } else {
        break;
      }
    }
  }

  void applyTagToNode(std::shared_ptr<ir::Node> node) {
    if (pendingTag_.has_value()) {
      auto *tag = node->getTag();
      if (tag) {
        *tag = ir::mergeTag(tag, &pendingTag_.value());
      }
      pendingTag_.reset();
    }
  }

  void applyTagsToObject(std::shared_ptr<ir::NodeObject> obj) {
    for (auto &field : obj->fields) {
      applyTags(field.value);
    }
  }

  void applyTagsToArray(std::shared_ptr<ir::NodeArray> arr) {
    auto *arrTag = arr->getTag();
    for (auto &item : arr->items) {
      auto *childTag = item->getTag();
      if (childTag && arrTag) {
        childTag->isInherit = true;
        if (!arrTag->childDesc.empty())
          childTag->desc = arrTag->childDesc;
        if (arrTag->childType != ir::ValueType::Unknown) {
          childTag->type = arrTag->childType;
        }
        if (arrTag->childNullable)
          childTag->nullable = true;
        if (arrTag->childAllowEmpty)
          childTag->allowEmpty = true;
        if (arrTag->childUnique)
          childTag->unique = true;
        if (!arrTag->child_default_val.empty())
          childTag->default_val = arrTag->child_default_val;
        if (!arrTag->childMin.empty())
          childTag->min = arrTag->childMin;
        if (!arrTag->childMax.empty())
          childTag->max = arrTag->childMax;
        if (arrTag->childSize != 0)
          childTag->size = arrTag->childSize;
        if (!arrTag->child_enums.empty())
          childTag->enums = arrTag->child_enums;
        if (!arrTag->childPattern.empty())
          childTag->pattern = arrTag->childPattern;
        if (arrTag->childLocationOffset != ir::DefaultLocationOffset)
          childTag->locationOffset = arrTag->childLocationOffset;
        if (arrTag->childVersion != ir::DefaultVersion)
          childTag->version = arrTag->childVersion;
        if (!arrTag->childMime.empty())
          childTag->mime = arrTag->childMime;
      }
      applyTags(item);
    }
  }

  std::shared_ptr<ir::Node> parsePrimitive() {
    if (pos_ >= tokens_.size())
      return nullptr;
    auto &tok = tokens_[pos_];
    ++pos_;

    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();
    val->text = tok.literal;

    switch (tok.type) {
    case TokenType::String:
      tag->type = ir::ValueType::Str;
      break;
    case TokenType::Number:
      if (tok.literal.find('.') != std::string::npos ||
          tok.literal.find('e') != std::string::npos ||
          tok.literal.find('E') != std::string::npos) {
        tag->type = ir::ValueType::F64;
      } else {
        tag->type = ir::ValueType::I;
      }
      break;
    case TokenType::True:
    case TokenType::False:
      tag->type = ir::ValueType::Bool;
      break;
    case TokenType::Null:
      throw std::runtime_error("null is not supported");
    default:
      break;
    }

    applyTagToNode(val);
    return val;
  }

  std::shared_ptr<ir::Node> parseObject() {
    ++pos_;
    auto obj = ir::makeNodeObject();

    // Save tag that was set before { (object-level tag from outer scope)
    auto objTag = pendingTag_;
    pendingTag_.reset();

    skipComments(); // comments after { are for the first field, not the object

    while (pos_ < tokens_.size() && tokens_[pos_].type != TokenType::RBrace) {
      if (tokens_[pos_].type == TokenType::Comma) {
        ++pos_;
        skipComments();
        continue;
      }

      if (tokens_[pos_].type != TokenType::String) {
        ++pos_;
        skipComments();
        continue;
      }

      std::string key = tokens_[pos_].literal;
      ++pos_;
      skipComments();

      if (pos_ < tokens_.size() && tokens_[pos_].type == TokenType::Colon) {
        ++pos_;
        skipComments();
      }

      auto savedTag = pendingTag_;
      pendingTag_.reset();
      auto value = parseValue();
      if (value) {
        if (savedTag.has_value()) {
          if (pendingTag_.has_value()) {
            pendingTag_ = ir::mergeTag(&savedTag.value(), &pendingTag_.value());
          } else {
            pendingTag_ = savedTag;
          }
        }
        skipComments();
        applyTagToNode(value);
        obj->fields.emplace_back(key, value);
      }

      skipComments();
    }

    if (pos_ < tokens_.size() && tokens_[pos_].type == TokenType::RBrace) {
      ++pos_;
    }

    pendingTag_ = objTag;
    applyTagToNode(obj);

    return obj;
  }

  std::shared_ptr<ir::Node> parseArray() {
    ++pos_;
    auto arr = ir::makeNodeArray();
    skipComments();

    auto outerTag = pendingTag_;
    pendingTag_.reset();

    while (pos_ < tokens_.size() && tokens_[pos_].type != TokenType::RBracket) {
      if (tokens_[pos_].type == TokenType::Comma) {
        ++pos_;
        skipComments();
        continue;
      }

      auto value = parseValue();
      if (value) {
        arr->items.push_back(value);
      }

      skipComments();
    }

    if (pos_ < tokens_.size() && tokens_[pos_].type == TokenType::RBracket) {
      ++pos_;
    }

    pendingTag_ = outerTag;
    applyTagToNode(arr);
    return arr;
  }

  std::shared_ptr<ir::Node> parseValue() {
    skipComments();
    if (pos_ >= tokens_.size())
      return nullptr;

    switch (tokens_[pos_].type) {
    case TokenType::LBrace:
      return parseObject();
    case TokenType::LBracket:
      return parseArray();
    case TokenType::String:
    case TokenType::Number:
    case TokenType::True:
    case TokenType::False:
    case TokenType::Null:
      return parsePrimitive();
    default:
      ++pos_;
      return nullptr;
    }
  }
};

} // namespace jsonc
} // namespace mmc

#endif