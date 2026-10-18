package core

import (
	"fmt"
)

// encodeTag 尾定界写出 tag 包装段（一趟正向，零回写）：
//
//	[inner 已编码 wire][reverse(meta 正序)][meta_len 尾部槽][total_len 档位大端][tag sign]
//
// bs 为已完成尾定界编码的 inner 字节；meta 为 Tag.Bytes() 正序字节。
// total = len(bs)+len(meta)+len(尾槽)，与解码侧（tailframe.go）反推式
// inner = total - ml - slotW 严格互逆；尾槽与前缀槽为同一套 canonical 三档。
// meta_len 尾槽三档见 footer.go footerSlotBytes。
func (e *encoder) encodeTag(bs []byte, meta []byte) (n uint32, err error) {
	if len(meta) == 0 {
		return
	}

	slot, err := footerSlotBytes(len(meta))
	if err != nil {
		return
	}

	length := len(bs) + len(meta) + len(slot)
	if length > Max4Byte {
		err = fmt.Errorf("tag+payload too long, max length: %d, actual: %d", Max4Byte, length)
		return
	}

	sign := PrefixTag
	var tier []byte
	if length < TagLen1Byte {
		sign |= Prefix(length)
	} else if length <= Max1Byte {
		sign |= Prefix(TagLen1Byte)
		tier = []byte{byte(length)}
	} else if length <= Max2Byte {
		sign |= Prefix(TagLen2Byte)
		tier = []byte{byte(length >> 8), byte(length)}
	} else if length <= Max3Byte {
		sign |= Prefix(TagLen3Byte)
		tier = []byte{byte(length >> 16), byte(length >> 8), byte(length)}
	// 同 encode_string：最后一档必须 <=，否则 tag+payload 正好 Max4Byte 时
	// 静默不写档位字节（节点连 tag 带内容整个丢失）。
	} else {
		sign |= Prefix(TagLen4Byte)
		tier = []byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}
	}

	rmq := reverseBytes(meta)
	if _, err = e.writeBytes(bs); err != nil {
		return
	}
	if _, err = e.writeBytes(rmq); err != nil {
		return
	}
	if _, err = e.writeBytes(slot); err != nil {
		return
	}
	if _, err = e.writeBytes(tier); err != nil {
		return
	}
	if _, err = e.writeByte(byte(sign)); err != nil {
		return
	}

	return uint32(len(bs) + len(rmq) + len(slot) + len(tier) + 1), nil
}
