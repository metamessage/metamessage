#include "parser.hpp"
#include <cstdlib>
#include <cmath>

namespace mmc {

std::unique_ptr<Node> createValueNode() {
    auto node = std::make_unique<Node>(Node::Type::Value);
    return node;
}

std::unique_ptr<Node> createObjectNode() {
    auto node = std::make_unique<Node>(Node::Type::Object);
    return node;
}

std::unique_ptr<Node> createArrayNode() {
    auto node = std::make_unique<Node>(Node::Type::Array);
    return node;
}

Parser::Parser(const std::string& input)
    : scanner_(input), currentToken_(std::nullopt) {
    nextToken();
}

Parser::~Parser() = default;

void Parser::nextToken() {
    currentToken_ = scanner_.nextToken();
}

std::unique_ptr<Node> Parser::parseValue() {
    if (!currentToken_.has_value()) {
        error_ = "No token available";
        return nullptr;
    }

    const Token& tok = currentToken_.value();

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
            return parsePrimitive(tok);
        case TokenType::EOF_:
            error_ = "Unexpected EOF";
            return nullptr;
        default:
            error_ = "Unexpected token";
            return nullptr;
    }
}

std::unique_ptr<Node> Parser::parsePrimitive(const Token& token) {
    auto node = createValueNode();
    auto& value = std::get<Value>(node->data);

    value.text = token.literal;

    switch (token.type) {
        case TokenType::String:
            value.data = token.literal;
            value.tag.type = ValueType::String;
            break;
        case TokenType::Number:
            if (token.literal.find('.') != std::string::npos) {
                value.data = std::stod(token.literal);
                value.tag.type = ValueType::Float64;
            } else if (!token.literal.empty() && token.literal[0] == '-') {
                value.data = std::stoll(token.literal);
                value.tag.type = ValueType::Int64;
            } else {
                try {
                    value.data = std::stoll(token.literal);
                    value.tag.type = ValueType::Int64;
                } catch (...) {
                    value.data = std::stoull(token.literal);
                    value.tag.type = ValueType::Uint64;
                }
            }
            break;
        case TokenType::True:
            value.data = true;
            value.tag.type = ValueType::Bool;
            break;
        case TokenType::False:
            value.data = false;
            value.tag.type = ValueType::Bool;
            break;
        case TokenType::Null:
            value.data = std::monostate{};
            value.tag.type = ValueType::Unknown;
            break;
        default:
            error_ = "Invalid primitive type";
            return nullptr;
    }

    nextToken();
    return node;
}

std::unique_ptr<Node> Parser::parseObject() {
    auto node = createObjectNode();
    auto& obj = std::get<Object>(node->data);
    obj.tag.type = ValueType::Object;

    nextToken();

    while (currentToken_.has_value() && currentToken_.value().type != TokenType::RBrace) {
        if (currentToken_.value().type == TokenType::Comma) {
            nextToken();
            continue;
        }

        if (currentToken_.value().type != TokenType::String) {
            if (currentToken_.value().type != TokenType::EOF_) {
                nextToken();
            }
            continue;
        }

        std::string key = currentToken_.value().literal;
        nextToken();

        if (!currentToken_.has_value() || currentToken_.value().type != TokenType::Colon) {
            continue;
        }
        nextToken();

        auto value = parseValue();
        if (value == nullptr) {
            continue;
        }

        obj.fields.emplace_back(std::move(key), std::move(value));
    }

    if (currentToken_.has_value() && currentToken_.value().type == TokenType::RBrace) {
        nextToken();
    }

    return node;
}

std::unique_ptr<Node> Parser::parseArray() {
    auto node = createArrayNode();
    auto& arr = std::get<Array>(node->data);
    arr.tag.type = ValueType::Array;

    nextToken();

    while (currentToken_.has_value() && currentToken_.value().type != TokenType::RBracket) {
        if (currentToken_.value().type == TokenType::Comma) {
            nextToken();
            continue;
        }

        auto value = parseValue();
        if (value != nullptr) {
            arr.items.push_back(std::move(value));
        } else {
            break;
        }
    }

    if (currentToken_.has_value() && currentToken_.value().type == TokenType::RBracket) {
        nextToken();
    }

    return node;
}

std::unique_ptr<Node> Parser::parse() {
    if (currentToken_.has_value() &&
        (currentToken_.value().type == TokenType::LBrace ||
         currentToken_.value().type == TokenType::LBracket ||
         currentToken_.value().type == TokenType::String ||
         currentToken_.value().type == TokenType::Number ||
         currentToken_.value().type == TokenType::True ||
         currentToken_.value().type == TokenType::False ||
         currentToken_.value().type == TokenType::Null)) {
        return parseValue();
    }
    return nullptr;
}

std::optional<std::string> Parser::getError() const {
    return error_;
}

}
