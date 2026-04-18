package mm

import "errors"

func (e *encoder) encodeBytes(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length > Max2Byte {
		err = errors.New("bytes too long")
		return
	}

	sign := PrefixBytes
	if length < BytesLen1Byte {
		sign |= Prefix(length)
		return e.writeBytesWithPrefix(bs, byte(sign))
	} else if length < Max1Byte {
		sign |= Prefix(BytesLen1Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length))
	} else if length < Max2Byte {
		sign |= Prefix(BytesLen2Byte)
		return e.writeBytesWithPrefix(bs, byte(sign), byte(length>>8), byte(length))
	}

	return
}
