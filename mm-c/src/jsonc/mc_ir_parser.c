#include "mc_ir_parser.h"
#include "../ir/mc_tag.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

typedef struct {
    const char* input;
    size_t pos;
    size_t len;
    char pending_tag[512];
    int has_pending_tag;
    int last_was_lbrace;
    int error;
} parse_ctx_t;

static void skip_ws(parse_ctx_t* ctx) {
    while (ctx->pos < ctx->len) {
        char c = ctx->input[ctx->pos];
        if (c == ' ' || c == '\t' || c == '\n' || c == '\r') {
            ctx->pos++;
        } else {
            break;
        }
    }
}

static int peek(parse_ctx_t* ctx) {
    skip_ws(ctx);
    if (ctx->pos >= ctx->len) return EOF;
    return (unsigned char)ctx->input[ctx->pos];
}

static int advance(parse_ctx_t* ctx) {
    if (ctx->pos >= ctx->len) return EOF;
    return (unsigned char)ctx->input[ctx->pos++];
}

static int match(parse_ctx_t* ctx, char c) {
    skip_ws(ctx);
    if (ctx->pos < ctx->len && ctx->input[ctx->pos] == c) {
        ctx->pos++;
        return 1;
    }
    return 0;
}

static void skip_line_comment(parse_ctx_t* ctx) {
    size_t start = ctx->pos;
    while (ctx->pos < ctx->len && ctx->input[ctx->pos] != '\n') {
        ctx->pos++;
    }
    size_t comment_len = ctx->pos - start;
    if (comment_len > 3 && ctx->input[start] == '/' && ctx->input[start+1] == '/') {
        const char* text = ctx->input + start;
        size_t text_len = comment_len;
        const char* mm_pos = text;
        while (mm_pos < text + text_len && *mm_pos != 'm') mm_pos++;
        if (mm_pos + 2 < text + text_len && mm_pos[0] == 'm' && mm_pos[1] == 'm' && mm_pos[2] == ':') {
            const char* tag_start = mm_pos + 3;
            while (tag_start < text + text_len && (*tag_start == ' ' || *tag_start == '\t')) tag_start++;
            size_t tag_len = text + text_len - tag_start;
            if (tag_len > 0 && tag_len < sizeof(ctx->pending_tag)) {
                memcpy(ctx->pending_tag, tag_start, tag_len);
                ctx->pending_tag[tag_len] = '\0';
                ctx->has_pending_tag = 1;
            }
        }
    }
}

static void skip_block_comment(parse_ctx_t* ctx) {
    while (ctx->pos + 1 < ctx->len) {
        if (ctx->input[ctx->pos] == '*' && ctx->input[ctx->pos+1] == '/') {
            ctx->pos += 2;
            return;
        }
        ctx->pos++;
    }
}

static void consume_comments(parse_ctx_t* ctx) {
    while (ctx->pos < ctx->len) {
        skip_ws(ctx);
        if (ctx->pos + 1 >= ctx->len) break;
        if (ctx->input[ctx->pos] == '/' && ctx->input[ctx->pos+1] == '/') {
            ctx->pos += 2;
            skip_line_comment(ctx);
        } else if (ctx->input[ctx->pos] == '/' && ctx->input[ctx->pos+1] == '*') {
            ctx->pos += 2;
            skip_block_comment(ctx);
        } else {
            break;
        }
    }
}

static char* parse_string(parse_ctx_t* ctx) {
    if (advance(ctx) != '"') {
        ctx->error = 1;
        return NULL;
    }
    size_t cap = 64, len = 0;
    char* result = malloc(cap);
    while (ctx->pos < ctx->len) {
        char c = ctx->input[ctx->pos];
        if (c == '"') {
            ctx->pos++;
            result[len] = '\0';
            return result;
        }
        if (c == '\\' && ctx->pos + 1 < ctx->len) {
            ctx->pos++;
            char esc = ctx->input[ctx->pos];
            switch (esc) {
                case '"': c = '"'; break;
                case '\\': c = '\\'; break;
                case '/': c = '/'; break;
                case 'n': c = '\n'; break;
                case 't': c = '\t'; break;
                case 'r': c = '\r'; break;
                default: c = esc; break;
            }
            ctx->pos++;
        } else {
            ctx->pos++;
        }
        if (len + 1 >= cap) {
            cap *= 2;
            result = realloc(result, cap);
        }
        result[len++] = c;
    }
    free(result);
    ctx->error = 1;
    return NULL;
}

static char* parse_raw_string(parse_ctx_t* ctx) {
    size_t start = ctx->pos;
    while (ctx->pos < ctx->len) {
        char c = ctx->input[ctx->pos];
        if (c == ',' || c == '}' || c == ']' || c == '/' || c == '\n' || c == '\r') {
            break;
        }
        ctx->pos++;
    }
    size_t len = ctx->pos - start;
    while (len > 0 && (ctx->input[start+len-1] == ' ' || ctx->input[start+len-1] == '\t')) {
        len--;
    }
    char* str = malloc(len + 1);
    memcpy(str, ctx->input + start, len);
    str[len] = '\0';
    return str;
}

static mm_node_t* parse_value(parse_ctx_t* ctx);

static mm_node_t* parse_object(parse_ctx_t* ctx) {
    mm_node_t* obj = mm_node_new_object();
    consume_comments(ctx);

    if (ctx->has_pending_tag) {
        mm_tag_t tag = mm_tag_parse(ctx->pending_tag);
        mm_tag_merge(&obj->data.object.tag, &tag);
        mm_tag_cleanup(&tag);
        ctx->has_pending_tag = 0;
    }

    if (ctx->pos >= ctx->len) { mm_node_free(obj); return NULL; }
    if (ctx->input[ctx->pos] != '{') { mm_node_free(obj); return NULL; }
    ctx->pos++;

    consume_comments(ctx);
    if (ctx->pos < ctx->len && ctx->input[ctx->pos] == '}') {
        ctx->pos++;
        return obj;
    }

    while (ctx->pos < ctx->len) {
        consume_comments(ctx);
        if (ctx->pos >= ctx->len) break;

        if (ctx->input[ctx->pos] == '}') {
            ctx->pos++;
            break;
        }

        char* key = parse_string(ctx);
        if (!key) break;

        consume_comments(ctx);
        if (!match(ctx, ':')) { free(key); break; }
        consume_comments(ctx);

        mm_node_t* val = parse_value(ctx);
        if (val) {
            if (ctx->has_pending_tag) {
                mm_tag_t tag = mm_tag_parse(ctx->pending_tag);
                mm_tag_merge(&val->data.value.tag, &tag);
                mm_tag_cleanup(&tag);
                ctx->has_pending_tag = 0;
            }
            mm_object_add_field(obj, key, val);
        }
        free(key);

        consume_comments(ctx);
        if (!match(ctx, ',')) {
            consume_comments(ctx);
            if (ctx->pos < ctx->len && ctx->input[ctx->pos] == '}') {
                ctx->pos++;
                break;
            }
        }
    }

    return obj;
}

static mm_node_t* parse_array(parse_ctx_t* ctx) {
    mm_node_t* arr = mm_node_new_array();
    if (ctx->pos >= ctx->len || ctx->input[ctx->pos] != '[') { mm_node_free(arr); return NULL; }
    ctx->pos++;

    consume_comments(ctx);
    if (ctx->pos < ctx->len && ctx->input[ctx->pos] == ']') {
        ctx->pos++;
        return arr;
    }

    while (ctx->pos < ctx->len) {
        consume_comments(ctx);
        if (ctx->input[ctx->pos] == ']') { ctx->pos++; break; }

        mm_node_t* item = parse_value(ctx);
        if (item) {
            if (ctx->has_pending_tag) {
                mm_tag_t tag = mm_tag_parse(ctx->pending_tag);
                mm_tag_merge(&item->data.value.tag, &tag);
                mm_tag_cleanup(&tag);
                ctx->has_pending_tag = 0;
            }
            mm_array_add_item(arr, item);
        } else {
            break;
        }

        consume_comments(ctx);
        if (!match(ctx, ',')) {
            if (ctx->pos < ctx->len && ctx->input[ctx->pos] == ']') {
                ctx->pos++;
                break;
            }
        }
    }

    return arr;
}

static mm_node_t* parse_value(parse_ctx_t* ctx) {
    consume_comments(ctx);

    if (ctx->pos >= ctx->len) return NULL;

    char c = ctx->input[ctx->pos];

    if (c == '{') return parse_object(ctx);
    if (c == '[') return parse_array(ctx);

    mm_node_t* node = mm_node_new_value();

    if (ctx->has_pending_tag && c != '{' && c != '[') {
        mm_tag_t tag = mm_tag_parse(ctx->pending_tag);
        mm_tag_merge(&node->data.value.tag, &tag);
        mm_tag_cleanup(&tag);
        ctx->has_pending_tag = 0;
    }

    if (c == '"') {
        char* str = parse_string(ctx);
        if (str) {
            node->data.value.text = str;
            node->data.value.tag.type = MM_VALUE_STR;
        } else {
            mm_node_free(node);
            return NULL;
        }
    } else {
        char* raw = parse_raw_string(ctx);
        if (raw) {
            node->data.value.text = raw;
            if (strcmp(raw, "true") == 0 || strcmp(raw, "false") == 0) {
                node->data.value.tag.type = MM_VALUE_BOOL;
            } else if (strcmp(raw, "null") == 0) {
                node->data.value.tag.type = MM_VALUE_UNKNOWN;
            } else if (strchr(raw, '.') || strchr(raw, 'e') || strchr(raw, 'E')) {
                node->data.value.tag.type = MM_VALUE_F64;
            } else {
                node->data.value.tag.type = MM_VALUE_I;
            }
        } else {
            mm_node_free(node);
            return NULL;
        }
    }

    return node;
}

mm_node_t* mm_jsonc_parse(const char* input) {
    if (!input) return NULL;

    parse_ctx_t ctx;
    ctx.input = input;
    ctx.pos = 0;
    ctx.len = strlen(input);
    ctx.has_pending_tag = 0;
    ctx.last_was_lbrace = 0;
    ctx.error = 0;
    memset(ctx.pending_tag, 0, sizeof(ctx.pending_tag));

    consume_comments(&ctx);

    if (ctx.pos >= ctx.len) return NULL;

    mm_node_t* root = NULL;
    if (ctx.input[ctx.pos] == '{') {
        root = parse_object(&ctx);
    } else {
        root = parse_value(&ctx);
    }

    if (ctx.error && root) {
        mm_node_free(root);
        return NULL;
    }

    return root;
}