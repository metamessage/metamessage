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
        return create_token(scanner, MMC_TOKEN_COMMENT, "");
    }

    char next = scanner->input[scanner->position];
    if (next != '/') {
        return create_token(scanner, MMC_TOKEN_INVALID, "");
    }

    advance(scanner, 1);
    while (scanner->position < scanner->input_len && scanner->input[scanner->position] != '\n') {
        advance(scanner, 1);
    }

    size_t content_start = start_pos + 2;
    size_t len = scanner->position - content_start;
    const char* raw = scanner->input + content_start;

    while (len > 0 && (raw[0] == ' ' || raw[0] == '\t' || raw[0] == '\r')) {
          raw++;
          len--;
      }
    while (len > 0 && (raw[len - 1] == ' ' || raw[len - 1] == '\t' || raw[len - 1] == '\r')) {
        len--;
    }

    char* literal = (char*)malloc(len + 1);
    if (literal == NULL) return NULL;
    memcpy(literal, raw, len);
    literal[len] = '\0';

    return create_token(scanner, MMC_TOKEN_COMMENT, literal);
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
    return scanner;
}

void mmc_scanner_free(MMC_Scanner* scanner) {
    if (scanner == NULL) return;
}

MMC_Token* mmc_scanner_next_token(MMC_Scanner* scanner) {
    skip_whitespace(scanner);

    if (scanner->position >= scanner->input_len) {
        return create_token(scanner, MMC_TOKEN_EOF, "");
    }

    char ch = scanner->input[scanner->position];

    if (ch == '/') {
        return scan_comment(scanner);
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
        case MMC_TOKEN_COMMENT: return "COMMENT";
        case MMC_TOKEN_INVALID: return "INVALID";
        default: return "UNKNOWN";
    }
}
