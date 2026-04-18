package mm

import "errors"

func (e *encoder) encodeInt(sign Prefix, uv uint64) (n uint32, err error) {
	switch {
	case uv < IntLen1Byte:
		sign |= Prefix(uv)
		return e.writeByte(byte(sign))
	case uv <= Max1Byte:
		sign |= IntLen1Byte
		return e.writeByte(byte(sign), byte(uv))
	case uv <= Max2Byte:
		sign |= IntLen2Byte
		return e.writeByte(byte(sign), byte(uv>>8), byte(uv))
	case uv <= Max3Byte:
		sign |= IntLen3Byte
		return e.writeByte(byte(sign), byte(uv>>16), byte(uv>>8), byte(uv))
	case uv <= Max4Byte:
		sign |= IntLen4Byte
		return e.writeByte(byte(sign), byte(uv>>24), byte(uv>>16), byte(uv>>8), byte(uv))
	case uv <= Max5Byte:
		sign |= IntLen5Byte
		return e.writeByte(byte(sign), byte(uv>>32), byte(uv>>24), byte(uv>>16), byte(uv>>8), byte(uv))
	case uv <= Max6Byte:
		sign |= IntLen6Byte
		return e.writeByte(byte(sign), byte(uv>>40), byte(uv>>32), byte(uv>>24), byte(uv>>16), byte(uv>>8), byte(uv))
	case uv <= Max7Byte:
		sign |= IntLen7Byte
		return e.writeByte(byte(sign), byte(uv>>48), byte(uv>>40), byte(uv>>32), byte(uv>>24), byte(uv>>16), byte(uv>>8), byte(uv))
	case uv <= Max8Byte:
		sign |= IntLen8Byte
		return e.writeByte(byte(sign), byte(uv>>56), byte(uv>>48), byte(uv>>40), byte(uv>>32), byte(uv>>24), byte(uv>>16), byte(uv>>8), byte(uv))
	default:
		err = errors.New("invalid byte size")
		return
	}
}
