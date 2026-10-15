package core

import "fmt"

func (e *encoder) encodeBytes(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length > Max4Byte {
		err = fmt.Errorf("bytes too long, max length: %d, actual: %d", Max4Byte, length)
		return
	}

	sign := PrefixBytes
	if length < BytesLen1Byte {
		sign |= Prefix(length)
		return e.writeBytesWithPrefix(bs, byte(sign))
	} else if length <= Max1Byte {
		sign |= Prefix(BytesLen1Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length))
	} else if length <= Max2Byte {
		sign |= Prefix(BytesLen2Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length>>8), byte(length))
	} else if length <= Max3Byte {
		sign |= Prefix(BytesLen3Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length>>16), byte(length>>8), byte(length))
	// 同 encode_string：最后一档必须 <=，否则 length == Max4Byte 静默写 0 字节。
	} else if length <= Max4Byte {
		sign |= Prefix(BytesLen4Byte)
		return e.writeBytesWithPrefix(bs, byte(sign),
			byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	}

	return
}
