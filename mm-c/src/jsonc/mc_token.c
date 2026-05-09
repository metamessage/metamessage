#include "mc_token.h"
#include <stdlib.h>
#include <string.h>

void mmc_token_free(MMC_Token* token) {
    if (token == NULL) return;
    if (token->literal != NULL) {
        free(token->literal);
        token->literal = NULL;
    }
}
