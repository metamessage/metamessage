#ifndef MMCPP_PARSER_HPP
#define MMCPP_PARSER_HPP

#include "ast.hpp"
#include "scanner.hpp"
#include <string>
#include <optional>
#include <stdexcept>

namespace mmc {

class Parser {
public:
    explicit Parser(const std::string& input);
    ~Parser();

    std::unique_ptr<Node> parse();
    std::optional<std::string> getError() const;

private:
    void nextToken();
    std::unique_ptr<Node> parseValue();
    std::unique_ptr<Node> parsePrimitive(const Token& token);
    std::unique_ptr<Node> parseObject();
    std::unique_ptr<Node> parseArray();

    Scanner scanner_;
    std::optional<Token> currentToken_;
    std::optional<std::string> error_;
};

}

#endif
