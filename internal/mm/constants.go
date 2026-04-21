package mm

import (
	"fmt"
)

type Prefix byte

const (
	Simple       Prefix = 0b000 << 5
	PositiveInt  Prefix = 0b001 << 5
	NegativeInt  Prefix = 0b010 << 5
	PrefixFloat  Prefix = 0b011 << 5
	PrefixString Prefix = 0b100 << 5
	PrefixBytes  Prefix = 0b101 << 5
	Container    Prefix = 0b110 << 5
	PrefixTag    Prefix = 0b111 << 5
)

func (p Prefix) String() string {
	switch p {
	case Simple:
		return "Simple"
	case PositiveInt:
		return "PositiveInt"
	case NegativeInt:
		return "NegativeInt"
	case PrefixFloat:
		return "Float"
	case PrefixString:
		return "String"
	case PrefixBytes:
		return "Bytes"
	case Container:
		return "Container"
	case PrefixTag:
		return "Tag"
	default:
		return fmt.Sprintf("Prefix(0x%02x)", byte(p))
	}
}

// Simple
type SimpleValue byte

const (
	SimpleNullBool SimpleValue = iota
	SimpleNullInt
	SimpleNullFloat
	SimpleNullString
	SimpleNullBytes

	SimpleFalse
	SimpleTrue

	SimpleCode
	SimpleMessage
	SimpleData
	SimpleSuccess
	SimpleError
	SimpleUnknown

	SimplePage
	SimpleLimit
	SimpleOffset
	SimpleTotal
	SimpleId
	SimpleName
	SimpleDescription
	SimpleType
	SimpleVersion
	SimpleStatus
	SimpleUrl
	SimpleCreateTime
	SimpleUpdateTime
	SimpleDeleteTime
	SimpleAccount
	SimpleToken
	SimpleExpireTime
	SimpleKey
	SimpleVal
)

func (s SimpleValue) String() string {
	switch s {
	case SimpleNullBool:
		return "null_bool"
	case SimpleNullInt:
		return "null_int"
	case SimpleNullFloat:
		return "null_float"
	case SimpleNullString:
		return "null_string"
	case SimpleNullBytes:
		return "null_bytes"
	case SimpleFalse:
		return "false"
	case SimpleTrue:
		return "true"
	case SimpleCode:
		return "code"
	case SimpleMessage:
		return "message"
	case SimpleData:
		return "data"
	case SimpleSuccess:
		return "success"
	case SimpleError:
		return "error"
	case SimpleUnknown:
		return "unknown"
	case SimplePage:
		return "page"
	case SimpleLimit:
		return "limit"
	case SimpleOffset:
		return "offset"
	case SimpleTotal:
		return "total"
	case SimpleId:
		return "id"
	case SimpleName:
		return "name"
	case SimpleDescription:
		return "description"
	case SimpleType:
		return "type"
	case SimpleVersion:
		return "version"
	case SimpleStatus:
		return "status"
	case SimpleUrl:
		return "url"
	case SimpleCreateTime:
		return "create_time"
	case SimpleUpdateTime:
		return "update_time"
	case SimpleDeleteTime:
		return "delete_time"
	case SimpleAccount:
		return "account"
	case SimpleToken:
		return "token"
	case SimpleExpireTime:
		return "expire_time"
	case SimpleKey:
		return "key"
	case SimpleVal:
		return "value"

	default:
		return fmt.Sprintf("SimpleValue(%d)", uint8(s))
	}
}

func (s SimpleValue) IsValid() bool {
	return s < 32
}

// Int
const (
	Max1Byte = 0xFF
	Max2Byte = 0xFFFF
	Max3Byte = 0xFFFFFF
	Max4Byte = 0xFFFFFFFF
	Max5Byte = 0xFFFFFFFFFF
	Max6Byte = 0xFFFFFFFFFFFF
	Max7Byte = 0xFFFFFFFFFFFFFF
	Max8Byte = 0xFFFFFFFFFFFFFFFF

	IntLenMask  = 0b11111
	IntLen1Byte = IntLenMask - 7
	IntLen2Byte = IntLenMask - 6
	IntLen3Byte = IntLenMask - 5
	IntLen4Byte = IntLenMask - 4
	IntLen5Byte = IntLenMask - 3
	IntLen6Byte = IntLenMask - 2
	IntLen7Byte = IntLenMask - 1
	IntLen8Byte = IntLenMask
)

// Float
const (
	FloatPositiveNegativeMask = 0b10000

	FloatLenMask  = 0b01111
	FloatLen1Byte = FloatLenMask - 7
	FloatLen2Byte = FloatLenMask - 6
	FloatLen3Byte = FloatLenMask - 5
	FloatLen4Byte = FloatLenMask - 4
	FloatLen5Byte = FloatLenMask - 3
	FloatLen6Byte = FloatLenMask - 2
	FloatLen7Byte = FloatLenMask - 1
	FloatLen8Byte = FloatLenMask
)

// String
const (
	StringLenMask  = 0b11111
	StringLen1Byte = StringLenMask - 1
	StringLen2Byte = StringLenMask
)

// Bytes
const (
	BytesLenMask  = 0b11111
	BytesLen1Byte = BytesLenMask - 1
	BytesLen2Byte = BytesLenMask
)

// Container
const (
	ContainerMask  = 0b10000
	ContainerMap   = 0b00000
	ContainerArray = 0b10000

	ContainerLenMask  = 0b01111
	ContainerLen1Byte = ContainerLenMask - 1
	ContainerLen2Byte = ContainerLenMask
)

// Tag
const (
	TagLenMask  = 0b11111
	TagLen1Byte = TagLenMask - 1
	TagLen2Byte = TagLenMask
)

// TagPayload
const (
	TagPayload1Byte = BytesLenMask - 1
	TagPayload2Byte = BytesLenMask
)

const (
	PrefixMask = 0b11100000
	SuffixMask = 0b00011111
)

func GetPrefix(b byte) Prefix {
	return Prefix(b & PrefixMask)
}

func GetSuffix(b byte) byte {
	return b & SuffixMask
}

func TagLen(b byte) (int, int) {
	l := int(b & TagLenMask)
	switch l {
	case TagLen1Byte:
		return 1, 0
	case TagLen2Byte:
		return 2, 0
	default:
		return 0, l
	}
}

func ContainerLen(b byte) (int, int) {
	l := int(b & ContainerLenMask)
	switch l {
	case ContainerLen1Byte:
		return 1, 0
	case ContainerLen2Byte:
		return 2, 0
	default:
		return 0, l
	}
}

func IsArray(b byte) bool {
	return b&ContainerMask == ContainerArray
}

func StringLen(b byte) (int, int) {
	l := int(b & StringLenMask)
	if l < StringLen1Byte {
		return 0, l
	} else if l == StringLen1Byte {
		return 1, l
	} else {
		return 2, l
	}
}

func BytesLen(b byte) (int, int) {
	l := int(b & BytesLenMask)
	if l < BytesLen1Byte {
		return 0, l
	} else if l == BytesLen1Byte {
		return 1, l
	} else {
		return 2, l
	}
}

func IntLen(b byte) (int, int) {
	l := int(b & IntLenMask)
	if l < IntLen1Byte {
		return 0, l
	} else {
		return l - IntLen1Byte + 1, 0
	}
}

func FloatLen(b byte) (int, int) {
	l := int(b & FloatLenMask)
	if l < FloatLen1Byte {
		return 0, l
	} else {
		return l - FloatLen1Byte + 1, 0
	}
}
