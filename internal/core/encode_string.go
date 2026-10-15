package core

import (
	"fmt"
)

func (e *encoder) encodeString(s string) (n uint32, err error) {
	length := len(s)
	if length > Max4Byte {
		err = fmt.Errorf("string too long, max length: %d, actual: %d", Max4Byte, length)
		return
	}

	sign := PrefixString
	if length < StringLen1Byte {
		sign |= Prefix(length)
		return e.writeStringWithPrefix(s, byte(sign))
	} else if length <= Max1Byte {
		sign |= Prefix(StringLen1Byte)
		return e.writeStringWithPrefix(s, byte(sign), byte(length))
	} else if length <= Max2Byte {
		sign |= Prefix(StringLen2Byte)
		return e.writeStringWithPrefix(s, byte(sign), byte(length>>8), byte(length))
	} else if length <= Max3Byte {
		sign |= Prefix(StringLen3Byte)
		return e.writeStringWithPrefix(s, byte(sign), byte(length>>16), byte(length>>8), byte(length))
	// 最后一档必须用 <=：Max4Byte(4294967295) 本身是可表达的（4 字节长度字段）。
	// 写成 < 会让 length == Max4Byte 时所有分支全落空，静默写 0 字节（数据丢失）。
	} else if length <= Max4Byte {
		sign |= Prefix(StringLen4Byte)
		return e.writeStringWithPrefix(s, byte(sign),
			byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	}

	return
}
