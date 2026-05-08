#ifndef MMCPP_SCANNER_HPP
#define MMCPP_SCANNER_HPP

#include "token.hpp"
#include <string>
#include <vector>
#include <optional>

namespace mmc {

class Scanner {
public:
    explicit Scanner(const std::string& input);

    Token nextToken();

private:
    void skipWhitespace();
    void advance(size_t count = 1);
    char peek() const;
    bool isLeadingComment() const;

    Token createToken(TokenType type, const std::string& literal = "");

    Token scanString();
    Token scanComment();
    Token scanNumber();
    Token scanIdentifier();

    const std::vector<char> input_;
    size_t position_ = 0;
    size_t line_ = 1;
    size_t column_ = 1;
    std::optional<Token> lastToken_;
    std::optional<Token> currentToken_;
};

}

#endif
