package mm

import (
	"fmt"
)

func (e *encoder) encodeTag(bs []byte, tag []byte) (n uint32, err error) {
	if len(tag) == 0 {
		return
	}

	length := len(bs) + len(tag)
	if length > Max2Byte {
		err = fmt.Errorf("tag+payload too long, max length: %d, actual: %d", Max2Byte, length)
		return
	}

	sign := PrefixTag

	var ns uint32
	var nb uint32
	if length < TagLen1Byte {
		sign |= Prefix(length)
		if ns, err = e.writeBytesWithPrefix(tag, byte(sign)); err != nil {
			return
		}
		nb, err = e.writeBytes(bs)
	} else if length < Max1Byte {
		sign |= Prefix(TagLen1Byte)
		if ns, err = e.writeBytesWithPrefix(tag, byte(sign), byte(length)); err != nil {
			return
		}
		nb, err = e.writeBytes(bs)
	} else if length < Max2Byte {
		sign |= Prefix(TagLen2Byte)
		if ns, err = e.writeBytesWithPrefix(tag, byte(sign), byte(length>>8), byte(length)); err != nil {
			return
		}
		nb, err = e.writeBytes(bs)
	}
	if err != nil {
		return
	}

	n = nb + ns
	return
}

func (e *encoder) encodeT(bs []byte) (n uint32, err error) {
	length := len(bs)
	if length == 0 {
		return
	}

	if length > Max2Byte {
		err = fmt.Errorf("tag too long, max length: %d, actual: %d", Max2Byte, length)
		return
	}

	if length < 254 {
		return e.writeBytesWithPrefix(bs, byte(length))
	}

	if length < 257 {
		return e.writeBytesWithPrefix(bs, 254, byte(length))
	}

	return e.writeBytesWithPrefix(bs, 255, byte(length>>8), byte(length))
}
