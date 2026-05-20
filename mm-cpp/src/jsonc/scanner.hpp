#ifndef MMCPP_JSONC_SCANNER_HPP
#define MMCPP_JSONC_SCANNER_HPP

#include "token.hpp"
#include <string>
#include <vector>
#include <optional>

namespace mmc {
namespace jsonc {

class Scanner {
public:
    explicit Scanner(const std::string& input)
        : input_(input.begin(), input.end()), position_(0), line_(1), column_(1) {}

    Token nextToken();
    std::vector<Token> scanAll() {
        std::vector<Token> tokens;
        while (true) {
            Token t = nextToken();
            tokens.push_back(t);
            if (t.type == TokenType::EOF_) break;
        }
        return tokens;
    }

private:
    void skipWhitespace();
    char peek() const;
    void advance(size_t count = 1);

    Token scanString();
    Token scanComment();
    Token scanNumber();
    Token scanIdentifier();

    bool isLeadingComment() const;

    Token createToken(TokenType type, const std::string& literal = "") {
        Token token(type, literal, line_, column_);
        lastToken_ = currentToken_;
        currentToken_ = token;
        return token;
    }

    std::vector<char> input_;
    size_t position_;
    size_t line_;
    size_t column_;
    std::optional<Token> lastToken_;
    std::optional<Token> currentToken_;
};

} // namespace jsonc
} // namespace mmc

#endif