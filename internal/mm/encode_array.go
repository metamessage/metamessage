package mm

import (
	"errors"
)

func (e *encoder) encodeArray(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length > Max2Byte {
		err = errors.New("array too long")
		return
	}

	sign := Container | ContainerArray
	if length < ContainerLen1Byte {
		sign |= Prefix(length)
		return e.writeBytesWithPrefix(bs, byte(sign))
	} else if length < Max1Byte {
		sign |= Prefix(ContainerLen1Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length))
	} else if length < Max2Byte {
		sign |= Prefix(ContainerLen2Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length>>8), byte(length))
	}

	return
}
