#include "scanner.hpp"
#include <cctype>
#include <stdexcept>

namespace mmc {
namespace jsonc {

char Scanner::peek() const {
    if (position_ >= input_.size()) return '\0';
    return input_[position_];
}

void Scanner::advance(size_t count) {
    for (size_t i = 0; i < count; ++i) {
        if (position_ >= input_.size()) break;
        if (input_[position_] == '\n') {
            ++line_;
            column_ = 0;
        } else {
            ++column_;
        }
        ++position_;
    }
}

void Scanner::skipWhitespace() {
    while (position_ < input_.size()) {
        char ch = input_[position_];
        if (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') {
            if (ch == '\n') {
                ++line_;
                column_ = 0;
            }
            advance();
        } else {
            break;
        }
    }
}

bool Scanner::isLeadingComment() const {
    if (!lastToken_.has_value()) return true;
    switch (lastToken_.value().type) {
        case TokenType::Comma:
        case TokenType::Colon:
        case TokenType::LBrace:
        case TokenType::LBracket:
            return true;
        default:
            return false;
    }
}

Token Scanner::scanString() {
    advance();
    std::string sb;
    while (position_ < input_.size() && input_[position_] != '"') {
        if (input_[position_] == '\\' && position_ + 1 < input_.size()) {
            advance();
            char escaped = input_[position_];
            switch (escaped) {
                case 'n': sb.push_back('\n'); break;
                case 'r': sb.push_back('\r'); break;
                case 't': sb.push_back('\t'); break;
                case '"': sb.push_back('"'); break;
                case '\\': sb.push_back('\\'); break;
                default: sb.push_back(escaped); break;
            }
        } else {
            sb.push_back(input_[position_]);
        }
        advance();
    }
    if (position_ < input_.size() && input_[position_] == '"') advance();
    return createToken(TokenType::String, sb);
}

Token Scanner::scanComment() {
    advance();
    if (position_ >= input_.size())
        return createToken(TokenType::LeadingComment, "");

    if (input_[position_] == '/') {
        advance();
        std::string comment;
        while (position_ < input_.size() && input_[position_] != '\n') {
            comment.push_back(input_[position_]);
            advance();
        }
        TokenType type = isLeadingComment() ? TokenType::LeadingComment : TokenType::TrailingComment;
        return createToken(type, comment);
    } else if (input_[position_] == '*') {
        advance();
        std::string comment;
        while (position_ + 1 < input_.size()) {
            if (input_[position_] == '*' && input_[position_ + 1] == '/') {
                advance(2);
                break;
            }
            comment.push_back(input_[position_]);
            advance();
        }
        TokenType type = isLeadingComment() ? TokenType::LeadingComment : TokenType::TrailingComment;
        return createToken(type, comment);
    }
    return createToken(TokenType::Invalid, "");
}

Token Scanner::scanNumber() {
    std::string sb;
    if (input_[position_] == '-') {
        sb.push_back('-');
        advance();
    }
    while (position_ < input_.size()) {
        char ch = input_[position_];
        if (std::isdigit(static_cast<unsigned char>(ch)) || ch == '.' ||
            ch == 'e' || ch == 'E' || ch == '+' || ch == '_') {
            if (ch == '_') { advance(); continue; }
            sb.push_back(ch);
            advance();
        } else {
            break;
        }
    }
    return createToken(TokenType::Number, sb);
}

Token Scanner::scanIdentifier() {
    std::string sb;
    while (position_ < input_.size() &&
           (std::isalnum(static_cast<unsigned char>(input_[position_])) ||
            input_[position_] == '_')) {
        sb.push_back(input_[position_]);
        advance();
    }
    std::string lower = sb;
    for (auto& c : lower) c = static_cast<char>(std::tolower(static_cast<unsigned char>(c)));
    if (lower == "true") return createToken(TokenType::True, "true");
    if (lower == "false") return createToken(TokenType::False, "false");
    if (lower == "null") return createToken(TokenType::Null, "null");
    return createToken(TokenType::String, sb);
}

Token Scanner::nextToken() {
    skipWhitespace();
    if (position_ >= input_.size())
        return createToken(TokenType::EOF_, "");

    char ch = input_[position_];

    if (ch == '/') return scanComment();

    switch (ch) {
        case '{': advance(); return createToken(TokenType::LBrace);
        case '}': advance(); return createToken(TokenType::RBrace);
        case '[': advance(); return createToken(TokenType::LBracket);
        case ']': advance(); return createToken(TokenType::RBracket);
        case ':': advance(); return createToken(TokenType::Colon);
        case ',': advance(); return createToken(TokenType::Comma);
        case '"': return scanString();
        default:
            if (std::isdigit(static_cast<unsigned char>(ch)) || ch == '-')
                return scanNumber();
            else if (std::isalpha(static_cast<unsigned char>(ch)))
                return scanIdentifier();
            advance();
            return createToken(TokenType::Invalid, "");
    }
}

} // namespace jsonc
} // namespace mmc