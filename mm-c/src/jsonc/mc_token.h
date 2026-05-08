#ifndef MMC_JSONC_TOKEN_H
#define MMC_JSONC_TOKEN_H

#include <stddef.h>

typedef enum {
    MMC_TOKEN_EOF,
    MMC_TOKEN_LBRACE,
    MMC_TOKEN_RBRACE,
    MMC_TOKEN_LBRACKET,
    MMC_TOKEN_RBRACKET,
    MMC_TOKEN_COLON,
    MMC_TOKEN_COMMA,
    MMC_TOKEN_STRING,
    MMC_TOKEN_NUMBER,
    MMC_TOKEN_TRUE,
    MMC_TOKEN_FALSE,
    MMC_TOKEN_NULL,
    MMC_TOKEN_LEADING_COMMENT,
    MMC_TOKEN_TRAILING_COMMENT,
    MMC_TOKEN_INVALID
} MMC_TokenType;

typedef struct {
    MMC_TokenType type;
    char* literal;
    size_t line;
    size_t column;
} MMC_Token;

void mmc_token_free(MMC_Token* token);

#endif
