package mm

import (
	"fmt"

	"github.com/lizongying/meta-message/internal/utils"
)

func (e *encoder) encodeFloat(s string) (n uint32, err error) {
	isNegative, exponent, mantissa, err := utils.ParseStringToUint64(s)
	if err != nil {
		return
	}

	sign := PrefixFloat
	if isNegative {
		sign |= FloatPositiveNegativeMask
	}

	// 0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9
	// 0.00, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09
	// 0, 1, 2, 3
	switch {
	case exponent == -1 && mantissa <= 7:
		sign |= Prefix(mantissa)
		return e.writeByte(byte(sign))
	// case exponent == -2 && mantissa <= 9:
	// 	sign |= Prefix(10 + mantissa)
	// 	return e.writeByte(byte(sign))
	// case exponent == 0 && mantissa <= 3:
	// 	sign |= Prefix(20 + mantissa)
	// 	return e.writeByte(byte(sign))
	case mantissa <= Max1Byte:
		sign |= FloatLen1Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa))
	case mantissa <= Max2Byte:
		sign |= FloatLen2Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max3Byte:
		sign |= FloatLen3Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max4Byte:
		sign |= FloatLen4Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>24), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max5Byte:
		sign |= FloatLen5Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>32), byte(mantissa>>24), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max6Byte:
		sign |= FloatLen6Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>40), byte(mantissa>>32), byte(mantissa>>24), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max7Byte:
		sign |= FloatLen7Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>48), byte(mantissa>>40), byte(mantissa>>32), byte(mantissa>>24), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	case mantissa <= Max8Byte:
		sign |= FloatLen8Byte
		return e.writeByte(byte(sign), byte(exponent), byte(mantissa>>56), byte(mantissa>>48), byte(mantissa>>40), byte(mantissa>>32), byte(mantissa>>24), byte(mantissa>>16), byte(mantissa>>8), byte(mantissa))
	default:
		err = fmt.Errorf("unsupported mantissa sign length: %d", sign)
		return
	}
}
