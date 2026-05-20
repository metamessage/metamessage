#ifndef MMC_JSONC_PRINTER_H
#define MMC_JSONC_PRINTER_H

#include "../ir/mc_ast.h"

char* mm_printer_to_jsonc(mm_node_t* node);
void mm_printer_free(char* str);

#endif