/* MetaMessage C test harness - parse JSONC file and re-print to JSONC. */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "mm.h"

static char *read_file(const char *path) {
    FILE *f = fopen(path, "rb");
    if (!f) {
        fprintf(stderr, "read error: cannot open %s\n", path);
        return NULL;
    }
    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    fseek(f, 0, SEEK_SET);
    char *buf = malloc(size + 1);
    if (!buf) {
        fclose(f);
        return NULL;
    }
    fread(buf, 1, size, f);
    buf[size] = '\0';
    fclose(f);
    return buf;
}

int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: harness [--encode|--decode] <file.jsonc>\n");
        return 1;
    }

    if (strcmp(argv[1], "--encode") == 0) {
        if (argc < 3) {
            fprintf(stderr, "usage: harness --encode <file.jsonc>\n");
            return 1;
        }
        char *input = read_file(argv[2]);
        if (!input) return 1;

        mm_buffer_t *buf = mm_encode_from_jsonc(input);
        free(input);
        if (!buf) {
            fprintf(stderr, "encode error\n");
            return 1;
        }

        for (size_t i = 0; i < buf->size; i++) {
            printf("%02x", buf->data[i]);
        }
        mm_buffer_free(buf);
        return 0;
    }

    if (strcmp(argv[1], "--decode") == 0) {
        char hex_input[65536];
        size_t hex_len = 0;
        int ch;
        while ((ch = getchar()) != EOF && hex_len < sizeof(hex_input) - 1) {
            if (ch != '\n' && ch != ' ') {
                hex_input[hex_len++] = ch;
            }
        }
        hex_input[hex_len] = '\0';

        if (hex_len % 2 != 0) {
            fprintf(stderr, "hex decode error: odd length\n");
            return 1;
        }

        size_t wire_len = hex_len / 2;
        uint8_t *wire = malloc(wire_len);
        for (size_t i = 0; i < wire_len; i++) {
            unsigned int byte;
            sscanf(&hex_input[i * 2], "%2x", &byte);
            wire[i] = (uint8_t)byte;
        }

        mm_buffer_t buf;
        buf.data = wire;
        buf.size = wire_len;
        buf.capacity = wire_len;

        char *output = mm_decode_to_jsonc(&buf);
        free(wire);
        if (!output) {
            fprintf(stderr, "decode error\n");
            return 1;
        }

        printf("%s", output);
        mm_string_free(output);
        return 0;
    }

    // Existing behavior
    char *input = read_file(argv[1]);
    if (!input) return 1;

    node_t *node = mm_from_jsonc(input);
    free(input);
    if (!node) {
        fprintf(stderr, "parse error\n");
        return 1;
    }

    char *output = mm_to_jsonc(node);
    if (!output) {
        fprintf(stderr, "print error\n");
        return 1;
    }

    printf("%s", output);
    mm_string_free(output);
    return 0;
}