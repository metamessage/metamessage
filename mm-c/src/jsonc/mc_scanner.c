#include "mc_scanner.h"
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>

struct MMC_Scanner {
    const char* input;
    size_t input_len;
    size_t position;
    size_t line;
    size_t column;
    MMC_Token* last_token;
    MMC_Token* current_token;
};

static void advance(MMC_Scanner* scanner, size_t count) {
    for (size_t i = 0; i < count; i++) {
        if (scanner->position >= scanner->input_len) break;
        char ch = scanner->input[scanner->position];
        if (ch == '\n') {
            scanner->line++;
            scanner->column = 0;
        } else {
            scanner->column++;
        }
        scanner->position++;
    }
}

static void skip_whitespace(MMC_Scanner* scanner) {
    while (scanner->position < scanner->input_len) {
        char ch = scanner->input[scanner->position];
        if (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') {
            if (ch == '\n') {
                scanner->line++;
                scanner->column = 0;
            }
            scanner->position++;
            scanner->column++;
        } else {
            break;
        }
    }
}

static int is_leading_comment(MMC_Scanner* scanner) {
    if (scanner->last_token == NULL) return 1;
    switch (scanner->last_token->type) {
        case MMC_TOKEN_COMMA:
        case MMC_TOKEN_COLON:
        case MMC_TOKEN_LBRACE:
        case MMC_TOKEN_LBRACKET:
            return 1;
        default:
            return 0;
    }
}

static MMC_Token* create_token(MMC_Scanner* scanner, MMC_TokenType type, const char* literal) {
    MMC_Token* token = (MMC_Token*)malloc(sizeof(MMC_Token));
    if (token == NULL) return NULL;
    token->type = type;
    token->line = scanner->line;
    token->column = scanner->column;
    if (literal != NULL) {
        token->literal = (char*)malloc(strlen(literal) + 1);
        if (token->literal) {
            strcpy(token->literal, literal);
        }
    } else {
        token->literal = NULL;
    }
    return token;
}

static MMC_Token* scan_string(MMC_Scanner* scanner) {
    size_t start_line = scanner->line;
    size_t start_column = scanner->column;
    advance(scanner, 1);

    size_t start_pos = scanner->position;
    while (scanner->position < scanner->input_len && scanner->input[scanner->position] != '"') {
        if (scanner->input[scanner->position] == '\\' && scanner->position + 1 < scanner->input_len) {
            advance(scanner, 2);
        } else {
            if (scanner->input[scanner->position] == '\n') {
                scanner->line++;
                scanner->column = 0;
            }
            advance(scanner, 1);
        }
    }

    if (scanner->position < scanner->input_len) {
        advance(scanner, 1);
    }

    size_t len = scanner->position - start_pos - 1;
    char* literal = (char*)malloc(len + 1);
    if (literal == NULL) return NULL;
    memcpy(literal, scanner->input + start_pos, len);
    literal[len] = '\0';

    return create_token(scanner, MMC_TOKEN_STRING, literal);
}

static MMC_Token* scan_comment(MMC_Scanner* scanner) {
    size_t start_line = scanner->line;
    size_t start_column = scanner->column;
    size_t start_pos = scanner->position;
    advance(scanner, 1);

    if (scanner->position >= scanner->input_len) {
        return create_token(scanner, MMC_TOKEN_LEADING_COMMENT, "");
    }

    char next = scanner->input[scanner->position];
    if (next == '/') {
        advance(scanner, 1);
        while (scanner->position < scanner->input_len && scanner->input[scanner->position] != '\n') {
            advance(scanner, 1);
        }
        size_t len = scanner->position - start_pos;
        char* literal = (char*)malloc(len + 1);
        if (literal == NULL) return NULL;
        memcpy(literal, scanner->input + start_pos, len);
        literal[len] = '\0';
        MMC_TokenType type = is_leading_comment(scanner) ? MMC_TOKEN_LEADING_COMMENT : MMC_TOKEN_TRAILING_COMMENT;
        return create_token(scanner, type, literal);
    } else if (next == '*') {
        advance(scanner, 1);
        while (scanner->position + 1 < scanner->input_len) {
            if (scanner->input[scanner->position] == '*' && scanner->input[scanner->position + 1] == '/') {
                advance(scanner, 2);
                break;
            }
            if (scanner->input[scanner->position] == '\n') {
                scanner->line++;
                scanner->column = 0;
            }
            advance(scanner, 1);
        }
        size_t len = scanner->position - start_pos;
        char* literal = (char*)malloc(len + 1);
        if (literal == NULL) return NULL;
        memcpy(literal, scanner->input + start_pos, len);
        literal[len] = '\0';
        MMC_TokenType type = is_leading_comment(scanner) ? MMC_TOKEN_LEADING_COMMENT : MMC_TOKEN_TRAILING_COMMENT;
        return create_token(scanner, type, literal);
    }

    return create_token(scanner, MMC_TOKEN_INVALID, "");
}

static MMC_Token* scan_number(MMC_Scanner* scanner) {
    size_t start_pos = scanner->position;
    while (scanner->position < scanner->input_len) {
        char ch = scanner->input[scanner->position];
        if (isdigit(ch) || ch == '.' || ch == '-' || ch == 'e' || ch == 'E' || ch == '+' || ch == '_') {
            if (ch == '_') {
                advance(scanner, 1);
                continue;
            }
            advance(scanner, 1);
        } else {
            break;
        }
    }
    size_t len = scanner->position - start_pos;
    char* literal = (char*)malloc(len + 1);
    if (literal == NULL) return NULL;
    memcpy(literal, scanner->input + start_pos, len);
    literal[len] = '\0';
    return create_token(scanner, MMC_TOKEN_NUMBER, literal);
}

static MMC_Token* scan_identifier(MMC_Scanner* scanner) {
    size_t start_pos = scanner->position;
    while (scanner->position < scanner->input_len) {
        char ch = scanner->input[scanner->position];
        if (isalnum(ch) || ch == '_') {
            advance(scanner, 1);
        } else {
            break;
        }
    }
    size_t len = scanner->position - start_pos;
    char* literal = (char*)malloc(len + 1);
    if (literal == NULL) return NULL;
    memcpy(literal, scanner->input + start_pos, len);
    literal[len] = '\0';

    if (strcmp(literal, "true") == 0) {
        return create_token(scanner, MMC_TOKEN_TRUE, literal);
    } else if (strcmp(literal, "false") == 0) {
        return create_token(scanner, MMC_TOKEN_FALSE, literal);
    } else if (strcmp(literal, "null") == 0) {
        return create_token(scanner, MMC_TOKEN_NULL, literal);
    }
    return create_token(scanner, MMC_TOKEN_STRING, literal);
}

MMC_Scanner* mmc_scanner_new(const char* input) {
    MMC_Scanner* scanner = (MMC_Scanner*)malloc(sizeof(MMC_Scanner));
    if (scanner == NULL) return NULL;
    scanner->input = input;
    scanner->input_len = strlen(input);
    scanner->position = 0;
    scanner->line = 1;
    scanner->column = 1;
    scanner->last_token = NULL;
    scanner->current_token = NULL;
    return scanner;
}

void mmc_scanner_free(MMC_Scanner* scanner) {
    if (scanner == NULL) return;
    if (scanner->last_token != NULL) {
        mmc_token_free(scanner->last_token);
        free(scanner->last_token);
    }
    if (scanner->current_token != NULL) {
        mmc_token_free(scanner->current_token);
        free(scanner->current_token);
    }
    free(scanner);
}

MMC_Token* mmc_scanner_next_token(MMC_Scanner* scanner) {
    skip_whitespace(scanner);

    if (scanner->position >= scanner->input_len) {
        MMC_Token* token = create_token(scanner, MMC_TOKEN_EOF, "");
        if (scanner->last_token != NULL) {
            mmc_token_free(scanner->last_token);
            free(scanner->last_token);
        }
        scanner->last_token = scanner->current_token;
        scanner->current_token = token;
        return token;
    }

    char ch = scanner->input[scanner->position];

    if (ch == '/') {
        MMC_Token* token = scan_comment(scanner);
        if (scanner->last_token != NULL) {
            mmc_token_free(scanner->last_token);
            free(scanner->last_token);
        }
        scanner->last_token = scanner->current_token;
        scanner->current_token = token;
        return token;
    }

    MMC_Token* token = NULL;
    switch (ch) {
        case '{':
            token = create_token(scanner, MMC_TOKEN_LBRACE, NULL);
            advance(scanner, 1);
            break;
        case '}':
            token = create_token(scanner, MMC_TOKEN_RBRACE, NULL);
            advance(scanner, 1);
            break;
        case '[':
            token = create_token(scanner, MMC_TOKEN_LBRACKET, NULL);
            advance(scanner, 1);
            break;
        case ']':
            token = create_token(scanner, MMC_TOKEN_RBRACKET, NULL);
            advance(scanner, 1);
            break;
        case ':':
            token = create_token(scanner, MMC_TOKEN_COLON, NULL);
            advance(scanner, 1);
            break;
        case ',':
            token = create_token(scanner, MMC_TOKEN_COMMA, NULL);
            advance(scanner, 1);
            break;
        case '"':
            token = scan_string(scanner);
            break;
        default:
            if (isdigit(ch) || ch == '-') {
                token = scan_number(scanner);
            } else if (isalpha(ch)) {
                token = scan_identifier(scanner);
            } else {
                token = create_token(scanner, MMC_TOKEN_INVALID, "");
                advance(scanner, 1);
            }
            break;
    }

    if (scanner->last_token != NULL) {
        mmc_token_free(scanner->last_token);
        free(scanner->last_token);
    }
    scanner->last_token = scanner->current_token;
    scanner->current_token = token;
    return token;
}

const char* mmc_token_type_to_string(MMC_TokenType type) {
    switch (type) {
        case MMC_TOKEN_EOF: return "EOF";
        case MMC_TOKEN_LBRACE: return "LBRACE";
        case MMC_TOKEN_RBRACE: return "RBRACE";
        case MMC_TOKEN_LBRACKET: return "LBRACKET";
        case MMC_TOKEN_RBRACKET: return "RBRACKET";
        case MMC_TOKEN_COLON: return "COLON";
        case MMC_TOKEN_COMMA: return "COMMA";
        case MMC_TOKEN_STRING: return "STRING";
        case MMC_TOKEN_NUMBER: return "NUMBER";
        case MMC_TOKEN_TRUE: return "TRUE";
        case MMC_TOKEN_FALSE: return "FALSE";
        case MMC_TOKEN_NULL: return "NULL";
        case MMC_TOKEN_LEADING_COMMENT: return "LEADING_COMMENT";
        case MMC_TOKEN_TRAILING_COMMENT: return "TRAILING_COMMENT";
        case MMC_TOKEN_INVALID: return "INVALID";
        default: return "UNKNOWN";
    }
}
