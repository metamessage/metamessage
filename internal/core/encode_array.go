package core

import "fmt"

func (e *encoder) encodeArray(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length > Max2Byte {
		err = fmt.Errorf("array too long, max length: %d, actual: %d", Max2Byte, length)
		return
	}

	sign := Container | ContainerArray
	if length < ContainerLen1Byte {
		sign |= Prefix(length)
		return e.writeBytesWithPrefix(bs, byte(sign))
	} else if length < Max1Byte {
		sign |= Prefix(ContainerLen1Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length))
	// 同 encode_string：最后一档必须 <=，否则 payload 正好 65535 字节时
	// 静默写 0 字节（数组内容整个丢失）。
	} else if length <= Max2Byte {
		sign |= Prefix(ContainerLen2Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length>>8), byte(length))
	}

	return
}
