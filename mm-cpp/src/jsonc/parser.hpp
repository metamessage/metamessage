#ifndef MMCPP_JSONC_PARSER_HPP
#define MMCPP_JSONC_PARSER_HPP

#include "token.hpp"
#include "../ir/ast.hpp"
#include "../ir/tag.hpp"
#include <string>
#include <vector>
#include <memory>
#include <optional>
#include <stdexcept>
#include <cstdlib>
#include <cctype>

namespace mmc {
namespace jsonc {

class Parser {
public:
    explicit Parser(const std::vector<Token>& tokens)
        : tokens_(tokens), pos_(0) {
        skipComments();
    }

    std::shared_ptr<ir::Node> parse() {
        if (pos_ >= tokens_.size()) return nullptr;
        auto& tok = tokens_[pos_];
        switch (tok.type) {
            case TokenType::LBrace: return parseObject();
            case TokenType::LBracket: return parseArray();
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
        if (!node) return;
        switch (node->getType()) {
            case ir::NodeType::Object:
                applyTagsToObject(std::static_pointer_cast<ir::Object>(node));
                break;
            case ir::NodeType::Array:
                applyTagsToArray(std::static_pointer_cast<ir::Array>(node));
                break;
            default:
                break;
        }
    }

private:
    const std::vector<Token>& tokens_;
    size_t pos_;
    std::optional<ir::Tag> pendingTag_;

    void skipComments() {
        while (pos_ < tokens_.size()) {
            auto& tok = tokens_[pos_];
            if (tok.type == TokenType::LeadingComment) {
                std::string comment = tok.literal;
                auto it = comment.find("mm:");
                if (it != std::string::npos) {
                    std::string tagStr = comment.substr(it + 3);
                    pendingTag_ = ir::Tag::parse(tagStr);
                }
                ++pos_;
            } else if (tok.type == TokenType::TrailingComment) {
                std::string comment = tok.literal;
                auto it = comment.find("mm:");
                if (it != std::string::npos) {
                    std::string tagStr = comment.substr(it + 3);
                    auto tag = ir::Tag::parse(tagStr);
                    if (!pendingTag_.has_value()) {
                        pendingTag_ = tag;
                    } else {
                        pendingTag_ = ir::mergeTag(
                            pendingTag_.has_value() ? &pendingTag_.value() : nullptr,
                            &tag
                        );
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
            auto* tag = node->getTag();
            if (tag) {
                *tag = ir::mergeTag(tag, &pendingTag_.value());
            }
            pendingTag_.reset();
        }
    }

    void applyTagsToObject(std::shared_ptr<ir::Object> obj) {
        for (auto& field : obj->fields) {
            applyTags(field.value);
        }
    }

    void applyTagsToArray(std::shared_ptr<ir::Array> arr) {
        for (auto& item : arr->items) {
            applyTags(item);
        }
    }

    std::shared_ptr<ir::Node> parsePrimitive() {
        if (pos_ >= tokens_.size()) return nullptr;
        auto& tok = tokens_[pos_];
        ++pos_;

        auto val = ir::makeValue();
        auto* tag = val->getTag();
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
                tag->type = ir::ValueType::Unknown;
                tag->isNull = true;
                break;
            default:
                break;
        }

        applyTagToNode(val);
        return val;
    }

    std::shared_ptr<ir::Node> parseObject() {
        ++pos_;
        auto obj = ir::makeObject();
        skipComments();

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

            auto value = parseValue();
            if (value) {
                skipComments();
                applyTagToNode(value);
                obj->fields.emplace_back(key, value);
            }

            skipComments();
        }

        if (pos_ < tokens_.size() && tokens_[pos_].type == TokenType::RBrace) {
            ++pos_;
        }

        applyTagToNode(obj);

        return obj;
    }

    std::shared_ptr<ir::Node> parseArray() {
        ++pos_;
        auto arr = ir::makeArray();
        skipComments();

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

        applyTagToNode(arr);
        return arr;
    }

    std::shared_ptr<ir::Node> parseValue() {
        skipComments();
        if (pos_ >= tokens_.size()) return nullptr;

        switch (tokens_[pos_].type) {
            case TokenType::LBrace: return parseObject();
            case TokenType::LBracket: return parseArray();
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