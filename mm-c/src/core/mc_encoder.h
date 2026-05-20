#ifndef MMC_CORE_ENCODER_H
#define MMC_CORE_ENCODER_H

#include <stdint.h>
#include <stddef.h>
#include "../ir/mc_ast.h"

typedef struct {
    uint8_t* data;
    size_t size;
    size_t capacity;
} mm_encoder_buffer_t;

mm_encoder_buffer_t* mm_encoder_encode(mm_node_t* node);
void mm_encoder_buffer_free(mm_encoder_buffer_t* buf);

#endif