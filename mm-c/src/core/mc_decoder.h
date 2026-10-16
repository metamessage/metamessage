#ifndef MMC_CORE_DECODER_H
#define MMC_CORE_DECODER_H

#include <stdint.h>
#include <stddef.h>
#include "../ir/mc_ast.h"

typedef struct {
    const uint8_t* data;
    size_t size;
    size_t offset;
    int error;  /* v0.3.0：置位表示遇到截断/越界数据，顶层 decode 据此返回 NULL */
} mm_decoder_t;

mm_decoder_t* mm_decoder_new(const uint8_t* data, size_t size);
void mm_decoder_free(mm_decoder_t* d);
node_t* mm_decoder_decode(mm_decoder_t* d);

#endif