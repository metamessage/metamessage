#ifndef MMCPP_JSONC_TOKEN_HPP
#define MMCPP_JSONC_TOKEN_HPP

#include <string>

namespace mmc {
namespace jsonc {

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

    Token(TokenType t = TokenType::Invalid, std::string lit = "",
          size_t l = 1, size_t c = 1)
        : type(t), literal(std::move(lit)), line(l), column(c) {}
};

} // namespace jsonc
} // namespace mmc

#endif