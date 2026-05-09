#ifndef MMC_JSONC_SCANNER_H
#define MMC_JSONC_SCANNER_H

#include "mc_token.h"

typedef struct MMC_Scanner MMC_Scanner;

MMC_Scanner* mmc_scanner_new(const char* input);
void mmc_scanner_free(MMC_Scanner* scanner);
MMC_Token* mmc_scanner_next_token(MMC_Scanner* scanner);
const char* mmc_token_type_to_string(MMC_TokenType type);

#endif
