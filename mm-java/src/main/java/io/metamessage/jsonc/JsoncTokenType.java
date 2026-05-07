package io.github.metamessage.jsonc;

public enum JsoncTokenType {
    EOF,
    LBRACE,
    RBRACE,
    LBRACKET,
    RBRACKET,
    COLON,
    COMMA,
    STRING,
    NUMBER,
    TRUE,
    FALSE,
    NULL,
    LEADING_COMMENT,
    TRAILING_COMMENT
}
