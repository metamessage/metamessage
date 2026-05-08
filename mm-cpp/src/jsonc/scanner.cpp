#include "scanner.hpp"
#include <cctype>
#include <stdexcept>

namespace mmc {

Scanner::Scanner(const std::string& input) : input_(input.begin(), input.end()) {}

char Scanner::peek() const {
    if (position_ >= input_.size()) {
        return '\0';
    }
    return input_[position_];
}

void Scanner::advance(size_t count) {
    for (size_t i = 0; i < count; ++i) {
        if (position_ >= input_.size()) break;
        char ch = input_[position_];
        if (ch == '\n') {
            line_++;
            column_ = 0;
        } else {
            column_++;
        }
        position_++;
    }
}

void Scanner::skipWhitespace() {
    while (position_ < input_.size()) {
        char ch = input_[position_];
        if (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') {
            if (ch == '\n') {
                line_++;
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

Token Scanner::createToken(TokenType type, const std::string& literal) {
    Token token(type, literal, line_, column_);
    lastToken_ = currentToken_;
    currentToken_ = token;
    return token;
}

Token Scanner::scanString() {
    size_t startLine = line_;
    size_t startColumn = column_;
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
                case 'b': sb.push_back('\x08'); break;
                case 'f': sb.push_back('\x0c'); break;
                case '"': sb.push_back('"'); break;
                case '\\': sb.push_back('\\'); break;
                case 'u':
                    if (position_ + 4 < input_.size()) {
                        advance();
                        std::string hex(input_.begin() + position_, input_.begin() + position_ + 4);
                        try {
                            unsigned int unicode = std::stoi(hex, nullptr, 16);
                            sb.push_back(static_cast<char>(unicode));
                            advance(3);
                        } catch (...) {
                            sb.push_back(escaped);
                        }
                    }
                    break;
                default: sb.push_back(escaped); break;
            }
        } else {
            if (input_[position_] == '\n') {
                line_++;
                column_ = 0;
            }
            sb.push_back(input_[position_]);
        }
        advance();
    }

    if (position_ < input_.size() && input_[position_] == '"') {
        advance();
    }

    return createToken(TokenType::String, sb);
}

Token Scanner::scanComment() {
    size_t startPos = position_;
    size_t startLine = line_;
    size_t startColumn = column_;
    advance();

    if (position_ >= input_.size()) {
        return createToken(TokenType::LeadingComment, "");
    }

    char next = input_[position_];
    if (next == '/') {
        advance();
        while (position_ < input_.size() && input_[position_] != '\n') {
            advance();
        }
        std::string literal(input_.begin() + startPos, input_.begin() + position_);
        TokenType type = isLeadingComment() ? TokenType::LeadingComment : TokenType::TrailingComment;
        return createToken(type, literal);
    } else if (next == '*') {
        advance();
        while (position_ + 1 < input_.size()) {
            if (input_[position_] == '*' && input_[position_ + 1] == '/') {
                advance(2);
                break;
            }
            if (input_[position_] == '\n') {
                line_++;
                column_ = 0;
            }
            advance();
        }
        std::string literal(input_.begin() + startPos, input_.begin() + position_);
        TokenType type = isLeadingComment() ? TokenType::LeadingComment : TokenType::TrailingComment;
        return createToken(type, literal);
    }

    return createToken(TokenType::Invalid, "");
}

Token Scanner::scanNumber() {
    size_t startPos = position_;
    std::string sb;

    if (input_[position_] == '-') {
        sb.push_back('-');
        advance();
    }

    while (position_ < input_.size()) {
        char ch = input_[position_];
        if (std::isdigit(ch) || ch == '.' || ch == 'e' || ch == 'E' || ch == '+' || ch == '_') {
            if (ch == '_') {
                advance();
                continue;
            }
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
    while (position_ < input_.size() && (std::isalnum(input_[position_]) || input_[position_] == '_')) {
        sb.push_back(input_[position_]);
        advance();
    }

    std::string lower = sb;
    for (auto& c : lower) c = std::tolower(c);

    if (lower == "true") {
        return createToken(TokenType::True, "true");
    } else if (lower == "false") {
        return createToken(TokenType::False, "false");
    } else if (lower == "null") {
        return createToken(TokenType::Null, "null");
    }
    return createToken(TokenType::String, sb);
}

Token Scanner::nextToken() {
    skipWhitespace();

    if (position_ >= input_.size()) {
        return createToken(TokenType::EOF_, "");
    }

    char ch = input_[position_];

    if (ch == '/') {
        return scanComment();
    }

    switch (ch) {
        case '{': return createToken(TokenType::LBrace);
        case '}': return createToken(TokenType::RBrace);
        case '[': return createToken(TokenType::LBracket);
        case ']': return createToken(TokenType::RBracket);
        case ':': return createToken(TokenType::Colon);
        case ',': return createToken(TokenType::Comma);
        case '"': return scanString();
        default:
            if (std::isdigit(ch) || ch == '-') {
                return scanNumber();
            } else if (std::isalpha(ch)) {
                return scanIdentifier();
            }
            advance();
            return createToken(TokenType::Invalid, "");
    }
}

}
