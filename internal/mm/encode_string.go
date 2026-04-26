package mm

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
		return e.writeStringWithPrefix(s, byte(sign))
	} else if length < Max1Byte {
		sign |= Prefix(StringLen1Byte)
		return e.writeStringWithPrefix(s, byte(sign), byte(length))
	} else if length < Max2Byte {
		sign |= Prefix(StringLen2Byte)
		return e.writeStringWithPrefix(s, byte(sign), byte(length>>8), byte(length))
	}

	return
}
