package mm

import (
	"errors"
)

func (e *encoder) encodeString(s string) (n uint32, err error) {
	// defer func() {
	// 	fmt.Printf("encodeString: %q\n", s)
	// }()

	length := len(s)
	if length > Max2Byte {
		err = errors.New("string too long")
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
