package core

import "fmt"

func (e *encoder) encodeArray(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length > Max4Byte {
		err = fmt.Errorf("array too long, max length: %d, actual: %d", Max4Byte, length)
		return
	}

	sign := Container | ContainerArray
	if length < ContainerLen1Byte {
		sign |= Prefix(length)
		return e.writeBytesWithSuffix(bs, byte(sign))
	} else if length <= Max1Byte {
		sign |= Prefix(ContainerLen1Byte)
		return e.writeBytesWithSuffix(bs, byte(length), byte(sign))
	} else if length <= Max2Byte {
		sign |= Prefix(ContainerLen2Byte)
		return e.writeBytesWithSuffix(bs, byte(length>>8), byte(length), byte(sign))
	} else if length <= Max3Byte {
		sign |= Prefix(ContainerLen3Byte)
		return e.writeBytesWithSuffix(bs, byte(length>>16), byte(length>>8), byte(length), byte(sign))
	// 同 encode_string：最后一档必须 <=，否则 payload 正好 Max4Byte 时
	// 静默写 0 字节（数组内容整个丢失）。
	} else if length <= Max4Byte {
		sign |= Prefix(ContainerLen4Byte)
		return e.writeBytesWithSuffix(bs,
			byte(length>>24), byte(length>>16), byte(length>>8), byte(length), byte(sign))
	}

	return
}
