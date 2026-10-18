package core

import (
	"fmt"
)

func (e *encoder) encodeString(s string) (n uint32, err error) {
	length := len(s)
	if length > Max2Byte {
		err = fmt.Errorf("string too long, max length: %d, actual: %d", Max2Byte, length)
		return
	}

	sign := PrefixString
	if length < StringLen1Byte {
		sign |= Prefix(length)
		return e.writeStringWithSuffix(s, byte(sign))
	} else if length < Max1Byte {
		sign |= Prefix(StringLen1Byte)
		return e.writeStringWithSuffix(s, byte(length), byte(sign))
	// 最后一档必须用 <=：Max2Byte(65535) 本身是可表达的（2 字节长度字段）。
	// 写成 < 会让 length == 65535 三个分支全落空，静默写 0 字节（数据丢失）。
	} else if length <= Max2Byte {
		sign |= Prefix(StringLen2Byte)
		return e.writeStringWithSuffix(s, byte(length>>8), byte(length), byte(sign))
	}

	return
}
