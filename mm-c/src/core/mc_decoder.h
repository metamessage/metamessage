#ifndef MMC_CORE_DECODER_H
#define MMC_CORE_DECODER_H

#include <stdint.h>
#include <stddef.h>
#include "../ir/mc_ast.h"

typedef struct {
    const uint8_t* data;
    size_t size;
    size_t offset;
} mm_decoder_t;

mm_decoder_t* mm_decoder_new(const uint8_t* data, size_t size);
void mm_decoder_free(mm_decoder_t* d);
mm_node_t* mm_decoder_decode(mm_decoder_t* d);

#endif