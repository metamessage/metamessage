#include "mc_parser.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <inttypes.h>

static void parser_error(MMC_Parser* parser, const char* message) {
    if (parser->error_message != NULL) {
        free(parser->error_message);
    }
    parser->error_message = (char*)malloc(256);
    if (parser->error_message != NULL) {
        snprintf(parser->error_message, 256, "Error at line %zu, column %zu: %s",
                 parser->current_token ? parser->current_token->line : 0,
                 parser->current_token ? parser->current_token->column : 0,
                 message);
    }
}

static MMC_Node* parse_value(MMC_Parser* parser);
static MMC_Node* parse_primitive(MMC_Parser* parser, MMC_Token* token);
static MMC_Node* parse_object(MMC_Parser* parser);
static MMC_Node* parse_array(MMC_Parser* parser);
static void next_token(MMC_Parser* parser);

static void next_token(MMC_Parser* parser) {
    if (parser->current_token != NULL) {
        mmc_token_free(parser->current_token);
        free(parser->current_token);
    }
    parser->current_token = mmc_scanner_next_token(parser->scanner);
}

static MMC_Node* parse_value(MMC_Parser* parser) {
    if (parser->current_token == NULL) {
        parser_error(parser, "No token available");
        return NULL;
    }

    switch (parser->current_token->type) {
        case MMC_TOKEN_LBRACE:
            return parse_object(parser);
        case MMC_TOKEN_LBRACKET:
            return parse_array(parser);
        case MMC_TOKEN_STRING:
        case MMC_TOKEN_NUMBER:
        case MMC_TOKEN_TRUE:
        case MMC_TOKEN_FALSE:
        case MMC_TOKEN_NULL:
            return parse_primitive(parser, parser->current_token);
        case MMC_TOKEN_EOF:
            parser_error(parser, "Unexpected EOF");
            return NULL;
        default:
            parser_error(parser, "Unexpected token");
            return NULL;
    }
}

static MMC_Node* parse_primitive(MMC_Parser* parser, MMC_Token* token) {
    MMC_Node* node = mmc_node_new_value();
    if (node == NULL) return NULL;

    node->data.value.text = token->literal ? strdup(token->literal) : NULL;

    switch (token->type) {
        case MMC_TOKEN_STRING:
            node->data.value.data_type = MMC_DATA_STRING;
            node->data.value.data.string_value = token->literal ? strdup(token->literal) : NULL;
            node->data.value.tag.type = MMC_VALUE_TYPE_STRING;
            break;
        case MMC_TOKEN_NUMBER:
            if (strchr(token->literal, '.') != NULL) {
                node->data.value.data_type = MMC_DATA_FLOAT;
                node->data.value.data.float_value = atof(token->literal);
                node->data.value.tag.type = MMC_VALUE_TYPE_FLOAT64;
            } else if (token->literal[0] == '-') {
                node->data.value.data_type = MMC_DATA_INT;
                node->data.value.data.int_value = atoll(token->literal);
                node->data.value.tag.type = MMC_VALUE_TYPE_INT64;
            } else {
                node->data.value.data_type = MMC_DATA_UINT;
                node->data.value.data.uint_value = strtoull(token->literal, NULL, 10);
                node->data.value.tag.type = MMC_VALUE_TYPE_UINT64;
            }
            break;
        case MMC_TOKEN_TRUE:
            node->data.value.data_type = MMC_DATA_BOOL;
            node->data.value.data.bool_value = 1;
            node->data.value.tag.type = MMC_VALUE_TYPE_BOOL;
            break;
        case MMC_TOKEN_FALSE:
            node->data.value.data_type = MMC_DATA_BOOL;
            node->data.value.data.bool_value = 0;
            node->data.value.tag.type = MMC_VALUE_TYPE_BOOL;
            break;
        case MMC_TOKEN_NULL:
            node->data.value.data_type = MMC_DATA_NULL;
            node->data.value.tag.type = MMC_VALUE_TYPE_UNKNOWN;
            break;
        default:
            mmc_node_free(node);
            parser_error(parser, "Invalid primitive type");
            return NULL;
    }

    next_token(parser);
    return node;
}

static MMC_Node* parse_object(MMC_Parser* parser) {
    MMC_Node* node = mmc_node_new_object();
    if (node == NULL) return NULL;

    node->data.object.tag.type = MMC_VALUE_TYPE_OBJECT;

    next_token(parser);

    size_t capacity = 4;
    node->data.object.fields = (MMC_Field*)malloc(capacity * sizeof(MMC_Field));
    if (node->data.object.fields == NULL) {
        mmc_node_free(node);
        return NULL;
    }
    node->data.object.field_count = 0;

    while (parser->current_token != NULL && parser->current_token->type != MMC_TOKEN_RBRACE) {
        if (parser->current_token->type == MMC_TOKEN_COMMA) {
            next_token(parser);
            continue;
        }

        if (parser->current_token->type != MMC_TOKEN_STRING) {
            if (parser->current_token->type != MMC_TOKEN_EOF) {
                next_token(parser);
            }
            continue;
        }

        char* key = strdup(parser->current_token->literal);
        next_token(parser);

        if (parser->current_token == NULL || parser->current_token->type != MMC_TOKEN_COLON) {
            free(key);
            continue;
        }
        next_token(parser);

        MMC_Node* value = parse_value(parser);
        if (value == NULL) {
            free(key);
            continue;
        }

        if (node->data.object.field_count >= capacity) {
            capacity *= 2;
            MMC_Field* new_fields = (MMC_Field*)realloc(node->data.object.fields, capacity * sizeof(MMC_Field));
            if (new_fields == NULL) {
                free(key);
                mmc_node_free(value);
                continue;
            }
            node->data.object.fields = new_fields;
        }

        node->data.object.fields[node->data.object.field_count].key = key;
        node->data.object.fields[node->data.object.field_count].value = value;
        node->data.object.field_count++;
    }

    if (parser->current_token != NULL && parser->current_token->type == MMC_TOKEN_RBRACE) {
        next_token(parser);
    }

    return node;
}

static MMC_Node* parse_array(MMC_Parser* parser) {
    MMC_Node* node = mmc_node_new_array();
    if (node == NULL) return NULL;

    node->data.array.tag.type = MMC_VALUE_TYPE_ARRAY;

    next_token(parser);

    size_t capacity = 4;
    node->data.array.items = (MMC_Node**)malloc(capacity * sizeof(MMC_Node*));
    if (node->data.array.items == NULL) {
        mmc_node_free(node);
        return NULL;
    }
    node->data.array.item_count = 0;

    while (parser->current_token != NULL && parser->current_token->type != MMC_TOKEN_RBRACKET) {
        if (parser->current_token->type == MMC_TOKEN_COMMA) {
            next_token(parser);
            continue;
        }

        MMC_Node* value = parse_value(parser);
        if (value != NULL) {
            if (node->data.array.item_count >= capacity) {
                capacity *= 2;
                MMC_Node** new_items = (MMC_Node**)realloc(node->data.array.items, capacity * sizeof(MMC_Node*));
                if (new_items == NULL) {
                    mmc_node_free(value);
                    continue;
                }
                node->data.array.items = new_items;
            }
            node->data.array.items[node->data.array.item_count] = value;
            node->data.array.item_count++;
        } else {
            break;
        }
    }

    if (parser->current_token != NULL && parser->current_token->type == MMC_TOKEN_RBRACKET) {
        next_token(parser);
    }

    return node;
}

MMC_Parser* mmc_parser_new(const char* input) {
    MMC_Parser* parser = (MMC_Parser*)malloc(sizeof(MMC_Parser));
    if (parser == NULL) return NULL;

    parser->scanner = mmc_scanner_new(input);
    if (parser->scanner == NULL) {
        free(parser);
        return NULL;
    }

    parser->current_token = NULL;
    parser->error_message = NULL;

    next_token(parser);
    return parser;
}

void mmc_parser_free(MMC_Parser* parser) {
    if (parser == NULL) return;
    if (parser->scanner != NULL) {
        mmc_scanner_free(parser->scanner);
    }
    if (parser->current_token != NULL) {
        mmc_token_free(parser->current_token);
        free(parser->current_token);
    }
    if (parser->error_message != NULL) {
        free(parser->error_message);
    }
    free(parser);
}

MMC_Node* mmc_parser_parse(MMC_Parser* parser) {
    if (parser == NULL || parser->current_token == NULL) {
        return NULL;
    }

    MMC_Node* result = parse_value(parser);

    if (parser->error_message != NULL) {
        mmc_node_free(result);
        return NULL;
    }

    return result;
}

const char* mmc_parser_get_error(MMC_Parser* parser) {
    if (parser == NULL) return "Parser is NULL";
    return parser->error_message ? parser->error_message : NULL;
}
