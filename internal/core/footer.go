package core

import "fmt"

// v0.5.0 尾定界：节点的 meta 块不在头部，而是整体逆序后放到节点尾部，
// 并在其前追加一个"反向可读"的自描述长度槽（尾部槽）。
//
// 节点布局（带 meta 时才有 PrefixTag 封装，meta 为空则裸编码）：
//
//	[inner wire][reverse(meta 正序字节)][尾部槽][total_len 大端][PrefixTag|档位 sign]
//
// total_len = len(inner) + len(meta) + len(尾部槽)。
// 尾部槽紧贴 total_len 左侧，从节点末端向前反向读取：槽末字节即档位标记或直接
// 长度，反向读到的字节流恰好等于 meta 的正序内容，因此解码器"从尾部开始、按反向
// 读取流程"即可还原 meta。
//
// 尾部槽三档（canonical 最小宽度编码）：
//   - metaLen < 254：[metaLen]                        槽宽 1
//   - metaLen <= 255：[metaLen][254]                  槽宽 2
//   - metaLen <= Max2Byte：[hi][lo][255]（u16 大端）   槽宽 3
const (
	FooterSlotDirectMax = 253 // 一档可直写的 meta 长度上限（尾字节 < 254）
	FooterSlotMarker1   = 254 // 二档标记：倒数第 2 字节为 u8 长度
	FooterSlotMarker2   = 255 // 三档标记：倒数第 3..2 字节为 u16 大端长度
)

// footerSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
// metaLen 为 meta 正序字节长度（即逆序区长度，两者等长）。
func footerSlotBytes(metaLen int) ([]byte, error) {
	switch {
	case metaLen <= FooterSlotDirectMax:
		return []byte{byte(metaLen)}, nil
	case metaLen <= Max1Byte:
		// 逆序区末字节之后依次是 [len(u8)][254]，从尾部反向读到的顺序为 254 → len。
		return []byte{byte(metaLen), FooterSlotMarker1}, nil
	case metaLen <= Max2Byte:
		return []byte{byte(metaLen >> 8), byte(metaLen), FooterSlotMarker2}, nil
	default:
		return nil, fmt.Errorf("tag too long, max length: %d, actual: %d", Max2Byte, metaLen)
	}
}

// parseFooterSlot 从节点尾端反向解析长度槽。
// data 为完整 wire 缓冲，nodeEnd 为节点（含封装）的结束偏移（不含）。
// 返回 meta 长度与槽宽。
func parseFooterSlot(data []byte, nodeEnd int) (metaLen int, slotW int, err error) {
	if nodeEnd < 1 {
		return 0, 0, fmt.Errorf("footer slot out of range")
	}
	marker := data[nodeEnd-1]
	switch marker {
	case FooterSlotMarker1:
		if nodeEnd < 2 {
			return 0, 0, fmt.Errorf("footer slot truncated")
		}
		return int(data[nodeEnd-2]), 2, nil
	case FooterSlotMarker2:
		if nodeEnd < 3 {
			return 0, 0, fmt.Errorf("footer slot truncated")
		}
		return int(data[nodeEnd-3])<<8 | int(data[nodeEnd-2]), 3, nil
	default:
		return int(marker), 1, nil
	}
}

// reverseBytes 返回 bs 的字节逆序拷贝。
func reverseBytes(bs []byte) []byte {
	out := make([]byte, len(bs))
	for i, b := range bs {
		out[len(bs)-1-i] = b
	}
	return out
}

// v0.5.0 尾定界布局中值节点不另设 payload_len 槽：各类型 sign 低 5 位的
// 档位/内联语义（v0.3.0 原样）同时承担定界信息，换位/定界算法见 tailframe.go。
