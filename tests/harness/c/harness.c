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
        fprintf(stderr, "usage: harness <file.jsonc>\n");
        return 1;
    }

    char *input = read_file(argv[1]);
    if (!input) return 1;

    mm_node_t *node = mm_from_jsonc(input);
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