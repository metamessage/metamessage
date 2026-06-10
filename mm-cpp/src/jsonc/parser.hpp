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
  static constexpr size_t kMaxDepth = 32;
  const std::vector<Token> &tokens_;
  size_t pos_;
  size_t depth_ = 0;
  std::optional<ir::Tag> pendingTag_;

  static std::string camelToSnake(const std::string &s) {
    if (s.empty())
      return {};
    std::string result;
    result.reserve(s.size() + 4);
    for (size_t i = 0; i < s.size(); ++i) {
      char c = s[i];
      if (std::isupper(static_cast<unsigned char>(c))) {
        if (i > 0) {
          bool prevUpper = std::isupper(static_cast<unsigned char>(s[i - 1]));
          bool nextUpper = i + 1 < s.size() &&
                           std::isupper(static_cast<unsigned char>(s[i + 1]));
          if (!prevUpper || !nextUpper) {
            result.push_back('_');
          }
        }
        result.push_back(
            static_cast<char>(std::tolower(static_cast<unsigned char>(c))));
      } else {
        result.push_back(c);
      }
    }
    return result;
  }

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
        // Trim whitespace
        size_t start = 0;
        while (start < comment.size() &&
               (comment[start] == ' ' || comment[start] == '\t'))
          ++start;
        if (start < comment.size() && comment.compare(start, 3, "mm:") == 0) {
          std::string tagStr = comment.substr(start + 3);
          // Trim leading whitespace from tagStr
          size_t ts = 0;
          while (ts < tagStr.size() &&
                 (tagStr[ts] == ' ' || tagStr[ts] == '\t'))
            ++ts;
          tagStr = tagStr.substr(ts);
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

    // Handle null early - return NodeNull without creating a scalar
    if (tok.type == TokenType::Null) {
      auto nullVal = ir::makeNodeNull();
      applyTagToNode(nullVal);
      // Check that the applied tag doesn't conflict with null
      if (nullVal->getTag()->type != ir::ValueType::Unknown) {
        throw std::runtime_error(
            "null is not supported for type " +
            std::to_string(static_cast<int>(nullVal->getTag()->type)));
      }
      return nullVal;
    }

    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();
    val->text = tok.literal;

    switch (tok.type) {
    case TokenType::String:
      if (tag->type == ir::ValueType::Unknown)
        tag->type = ir::ValueType::Str;
      break;
    case TokenType::Number:
      if (tag->type == ir::ValueType::Unknown) {
        if (tok.literal.find('.') != std::string::npos ||
            tok.literal.find('e') != std::string::npos ||
            tok.literal.find('E') != std::string::npos) {
          tag->type = ir::ValueType::F64;
        } else {
          tag->type = ir::ValueType::I;
        }
      }
      break;
    case TokenType::True:
    case TokenType::False:
      if (tag->type == ir::ValueType::Unknown)
        tag->type = ir::ValueType::Bool;
      break;
    default:
      break;
    }

    applyTagToNode(val);
    return val;
  }

  std::shared_ptr<ir::Node> parseObject() {
    ++pos_;
    auto obj = ir::makeNodeObject();

    depth_++;
    if (depth_ > kMaxDepth) {
      throw std::runtime_error("max depth: 32");
    }

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

      std::string key = camelToSnake(tokens_[pos_].literal);
      ++pos_;
      skipComments();

      if (pos_ < tokens_.size() && tokens_[pos_].type == TokenType::Colon) {
        ++pos_;
        skipComments();
      }

      auto savedTag = pendingTag_;
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
    // Set default type if not set via tag
    if (obj->tag.type == ir::ValueType::Unknown) {
      obj->tag.type = ir::ValueType::Obj;
    }
    applyTagToNode(obj);

    depth_--;
    return obj;
  }

  std::shared_ptr<ir::Node> parseArray() {
    ++pos_;
    auto arr = ir::makeNodeArray();

    depth_++;
    if (depth_ > kMaxDepth) {
      throw std::runtime_error("max depth: 32");
    }

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
    // Set default type if not set via tag
    if (arr->tag.type == ir::ValueType::Unknown) {
      arr->tag.type = ir::ValueType::Vec;
    }
    applyTagToNode(arr);

    // Run validation
    if (!arr->tag.example) {
      if (arr->tag.type == ir::ValueType::Vec) {
        auto vr = arr->tag.validateVec(arr->items.size());
        if (!vr.isValid) {
          throw std::runtime_error(vr.error);
        }
      } else if (arr->tag.type == ir::ValueType::Arr) {
        auto vr = arr->tag.validateArr(arr->items.size());
        if (!vr.isValid) {
          throw std::runtime_error(vr.error);
        }
      }
    }

    depth_--;
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