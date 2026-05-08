#ifndef MMCPP_TOKEN_HPP
#define MMCPP_TOKEN_HPP

#include <string>
#include <variant>

namespace mmc {

enum class TokenType {
    EOF_,
    LBrace,
    RBrace,
    LBracket,
    RBracket,
    Colon,
    Comma,
    String,
    Number,
    True,
    False,
    Null,
    LeadingComment,
    TrailingComment,
    Invalid
};

struct Token {
    TokenType type;
    std::string literal;
    size_t line;
    size_t column;

    Token(TokenType t, std::string lit = "", size_t l = 1, size_t c = 1)
        : type(t), literal(std::move(lit)), line(l), column(c) {}
};

}

#endif
