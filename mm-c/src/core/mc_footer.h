#ifndef MC_CORE_FOOTER_H
#define MC_CORE_FOOTER_H

#include <stdint.h>
#include <stddef.h>

/*
 * v0.5.0 尾定界辅助（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
 *
 *   1. mc_footer_slot_bytes / mc_reverse_bytes：编码器写出 meta 逆序区 + 尾部长度槽。
 *   2. mc_untailify：解码入口，从节点尾 sign 反向定界，把尾定界 wire 还原成
 *      「sign 在前」的头部形态，再交给原有正向解码器。
 *
 * 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本文件不含编码侧换位逻辑。
 */

/* footer_slot_bytes：编码尾部长度槽写入 out（out 至少 3 字节），返回槽宽；meta_len 超 65535 返回 0。 */
size_t mc_footer_slot_bytes(size_t meta_len, uint8_t *out);

/* reverse_bytes：把 src 的 len 字节逆序拷贝到 dst（dst 至少 len 字节）。 */
void mc_reverse_bytes(const uint8_t *src, uint8_t *dst, size_t len);

/* untailify：整帧尾定界 wire → 头部形态。成功返回 malloc 缓冲区（*out_len 设长度），失败返回 NULL。 */
uint8_t *mc_untailify(const uint8_t *wire, size_t len, size_t *out_len);

#endif
