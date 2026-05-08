#ifndef MMC_JSONC_PARSER_H
#define MMC_JSONC_PARSER_H

#include "mc_ast.h"
#include "mc_scanner.h"

typedef struct MMC_Parser {
    MMC_Scanner* scanner;
    MMC_Token* current_token;
    char* error_message;
} MMC_Parser;

MMC_Parser* mmc_parser_new(const char* input);
void mmc_parser_free(MMC_Parser* parser);
MMC_Node* mmc_parser_parse(MMC_Parser* parser);
const char* mmc_parser_get_error(MMC_Parser* parser);

#endif
