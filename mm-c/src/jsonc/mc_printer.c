#include "mc_printer.h"
#include "../ir/mc_tag.h"
#include "../ir/mc_value_type.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <ctype.h>

typedef struct {
    char* data;
    size_t size;
    size_t cap;
} strbuf_t;

static void sb_write(strbuf_t* sb, const char* s, size_t len)
{
    if (sb->size + len > sb->cap) {
        sb->cap = sb->cap ? sb->cap * 2 : 256;
        while (sb->size + len > sb->cap) {
            sb->cap *= 2;
        }
        sb->data = (char*)realloc(sb->data, sb->cap);
    }
    memcpy(sb->data + sb->size, s, len);
    sb->size += len;
}

static void sb_printf(strbuf_t* sb, const char* fmt, ...)
{
    va_list ap;
    va_start(ap, fmt);
    int n = vsnprintf(NULL, 0, fmt, ap);
    va_end(ap);
    if (n < 0) return;

    if (sb->size + (size_t)n + 1 > sb->cap) {
        sb->cap = sb->cap ? sb->cap * 2 : 256;
        while (sb->size + (size_t)n + 1 > sb->cap) {
            sb->cap *= 2;
        }
        sb->data = (char*)realloc(sb->data, sb->cap);
    }

    va_start(ap, fmt);
    vsnprintf(sb->data + sb->size, (size_t)n + 1, fmt, ap);
    va_end(ap);
    sb->size += (size_t)n;
}

static void sb_putc(strbuf_t* sb, char c)
{
    if (sb->size + 1 > sb->cap) {
        sb->cap = sb->cap ? sb->cap * 2 : 256;
        sb->data = (char*)realloc(sb->data, sb->cap);
    }
    sb->data[sb->size++] = c;
}

static void sb_puts(strbuf_t* sb, const char* s)
{
    if (s) {
        sb_write(sb, s, strlen(s));
    }
}

static void print_indent(strbuf_t* sb, int depth)
{
    for (int i = 0; i < depth; i++) {
        sb_putc(sb, '\t');
    }
}

static void print_string(strbuf_t* sb, const char* s)
{
    sb_putc(sb, '"');
    for (; *s; s++) {
        unsigned char c = (unsigned char)*s;
        switch (c) {
        case '"':
            sb_puts(sb, "\\\"");
            break;
        case '\\':
            sb_puts(sb, "\\\\");
            break;
        case '\n':
            sb_puts(sb, "\\n");
            break;
        case '\t':
            sb_puts(sb, "\\t");
            break;
        case '\r':
            sb_puts(sb, "\\r");
            break;
        default:
            if (c < 0x20) {
                char buf[8];
                snprintf(buf, sizeof(buf), "\\u%04x", c);
                sb_puts(sb, buf);
            } else {
                sb_putc(sb, c);
            }
            break;
        }
    }
    sb_putc(sb, '"');
}

static void write_leading_tag_comment(strbuf_t* sb, const mm_tag_t* tag, int depth)
{
    char* tag_str = mm_tag_to_string(tag);
    if (tag_str && tag_str[0] != '\0') {
        sb_putc(sb, '\n');
        print_indent(sb, depth);
        sb_printf(sb, "// mm: %s\n", tag_str);
    }
    free(tag_str);
}

static void print_value(strbuf_t* sb, const mm_value_t* value, int depth)
{
    (void)depth;
    if (!value->text) {
        sb_puts(sb, "null");
        return;
    }

    switch (value->tag.type) {
    case MM_VALUE_STR:
    case MM_VALUE_BYTES:
    case MM_VALUE_DATETIME:
    case MM_VALUE_DATE:
    case MM_VALUE_TIME:
    case MM_VALUE_UUID:
    case MM_VALUE_IP:
    case MM_VALUE_URL:
    case MM_VALUE_EMAIL:
    case MM_VALUE_ENUM:
        print_string(sb, value->text);
        break;
    default:
        sb_puts(sb, value->text);
        break;
    }
}

static void print_node(strbuf_t* sb, mm_node_t* node, int depth);

static void print_object(strbuf_t* sb, const mm_object_t* obj, int depth)
{
    sb_printf(sb, "{\n");

    for (size_t i = 0; i < obj->field_count; i++) {
        mm_field_t* field = &obj->fields[i];
        write_leading_tag_comment(sb, &field->value->tag, depth + 1);
        print_indent(sb, depth + 1);
        print_string(sb, field->key);
        sb_puts(sb, ": ");
        print_node(sb, field->value, depth + 1);
        sb_puts(sb, ",\n");
    }

    print_indent(sb, depth);
    sb_putc(sb, '}');
}

static void print_array(strbuf_t* sb, const mm_array_t* arr, int depth)
{
    sb_printf(sb, "[\n");

    for (size_t i = 0; i < arr->item_count; i++) {
        mm_node_t* item = arr->items[i];
        write_leading_tag_comment(sb, &item->tag, depth + 1);
        print_indent(sb, depth + 1);
        print_node(sb, item, depth + 1);
        sb_puts(sb, ",\n");
    }

    print_indent(sb, depth);
    sb_putc(sb, ']');
}

static void print_doc(strbuf_t* sb, const mm_doc_t* doc, int depth)
{
    sb_printf(sb, "{\n");

    for (size_t i = 0; i < doc->field_count; i++) {
        mm_field_t* field = &doc->fields[i];
        write_leading_tag_comment(sb, &field->value->tag, depth + 1);
        print_indent(sb, depth + 1);
        print_string(sb, field->key);
        sb_puts(sb, ": ");
        print_node(sb, field->value, depth + 1);
        sb_puts(sb, ",\n");
    }

    print_indent(sb, depth);
    sb_putc(sb, '}');
}

static void print_node(strbuf_t* sb, mm_node_t* node, int depth)
{
    if (!node) {
        sb_puts(sb, "null");
        return;
    }

    switch (node->type) {
    case MM_NODE_VALUE:
        print_value(sb, &node->data.value, depth);
        break;
    case MM_NODE_OBJECT:
        print_object(sb, &node->data.object, depth);
        break;
    case MM_NODE_ARRAY:
        print_array(sb, &node->data.array, depth);
        break;
    case MM_NODE_DOC:
        print_doc(sb, &node->data.doc, depth);
        break;
    default:
        sb_puts(sb, "null");
        break;
    }
}

char* mm_printer_to_jsonc(mm_node_t* node)
{
    if (!node) return strdup("");

    strbuf_t sb;
    memset(&sb, 0, sizeof(sb));

    write_leading_tag_comment(&sb, &node->tag, 0);
    print_node(&sb, node, 0);

    sb_putc(&sb, '\0');

    return sb.data;
}

void mm_printer_free(char* str)
{
    free(str);
}