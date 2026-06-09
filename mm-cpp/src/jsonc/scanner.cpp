#include "scanner.hpp"
#include <cctype>
#include <stdexcept>

namespace mmc {
namespace jsonc {

char Scanner::peek() const {
  if (position_ >= input_.size())
    return '\0';
  return input_[position_];
}

void Scanner::advance(size_t count) {
  for (size_t i = 0; i < count; ++i) {
    if (position_ >= input_.size())
      break;
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
      advance();
    } else {
      break;
    }
  }
}

Token Scanner::scanString() {
  size_t startLine = line_;
  size_t startCol = column_;
  advance();
  std::string sb;
  while (position_ < input_.size() && input_[position_] != '"') {
    if (input_[position_] == '\\' && position_ + 1 < input_.size()) {
      sb.push_back(input_[position_]);
      advance();
      sb.push_back(input_[position_]);
    } else {
      sb.push_back(input_[position_]);
    }
    advance();
  }
  if (position_ < input_.size() && input_[position_] == '"')
    advance();
  return Token(TokenType::String, sb, startLine, startCol);
}

Token Scanner::scanComment() {
  size_t startLine = line_;
  size_t startCol = column_;
  advance();
  if (position_ >= input_.size())
    return Token(TokenType::Comment, "", startLine, startCol);

  if (input_[position_] == '/') {
    advance();
    std::string comment;
    while (position_ < input_.size() && input_[position_] != '\n') {
      comment.push_back(input_[position_]);
      advance();
    }
    // Trim whitespace (matching Go's strings.TrimSpace)
    size_t start = 0;
    while (start < comment.size() &&
           (comment[start] == ' ' || comment[start] == '\t' ||
            comment[start] == '\r' || comment[start] == '\n'))
      ++start;
    size_t end = comment.size();
    while (end > start &&
           (comment[end - 1] == ' ' || comment[end - 1] == '\t' ||
            comment[end - 1] == '\r' || comment[end - 1] == '\n'))
      --end;
    comment = comment.substr(start, end - start);
    return Token(TokenType::Comment, comment, startLine, startCol);
  }
  return Token(TokenType::Comment, "", startLine, startCol);
}

Token Scanner::scanNumber() {
  size_t startLine = line_;
  size_t startCol = column_;
  std::string sb;
  if (input_[position_] == '-') {
    sb.push_back('-');
    advance();
  }
  while (position_ < input_.size()) {
    char ch = input_[position_];
    if (std::isdigit(static_cast<unsigned char>(ch)) || ch == '.' ||
        ch == 'e' || ch == 'E' || ch == '+' || ch == '_') {
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
  return Token(TokenType::Number, sb, startLine, startCol);
}

Token Scanner::scanIdentifier() {
  size_t startLine = line_;
  size_t startCol = column_;
  std::string sb;
  while (position_ < input_.size() &&
         (std::isalnum(static_cast<unsigned char>(input_[position_])) ||
          input_[position_] == '_')) {
    sb.push_back(input_[position_]);
    advance();
  }
  if (sb == "true")
    return Token(TokenType::True, sb, startLine, startCol);
  if (sb == "false")
    return Token(TokenType::False, sb, startLine, startCol);
  if (sb == "null")
    return Token(TokenType::Null, sb, startLine, startCol);
  return Token(TokenType::String, sb, startLine, startCol);
}

Token Scanner::nextToken() {
  skipWhitespace();
  if (position_ >= input_.size())
    return createToken(TokenType::EOF_, "");

  char ch = input_[position_];
  size_t startLine = line_;
  size_t startCol = column_;

  if (ch == '/')
    return scanComment();

  switch (ch) {
  case '{':
    advance();
    return Token(TokenType::LBrace, "", startLine, startCol);
  case '}':
    advance();
    return Token(TokenType::RBrace, "", startLine, startCol);
  case '[':
    advance();
    return Token(TokenType::LBracket, "", startLine, startCol);
  case ']':
    advance();
    return Token(TokenType::RBracket, "", startLine, startCol);
  case ':':
    advance();
    return Token(TokenType::Colon, "", startLine, startCol);
  case ',':
    advance();
    return Token(TokenType::Comma, "", startLine, startCol);
  case '"':
    return scanString();
  default:
    if (std::isdigit(static_cast<unsigned char>(ch)) || ch == '-')
      return scanNumber();
    else if (std::isalpha(static_cast<unsigned char>(ch)))
      return scanIdentifier();
    advance();
    return Token(TokenType::Comment, "", startLine, startCol);
  }
}

} // namespace jsonc
} // namespace mmc