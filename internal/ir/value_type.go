package ir

import (
	"fmt"
	"strings"
)

type ValueType uint8

const (
	ValueTypeUnknown ValueType = iota

	ValueTypeDoc
	ValueTypeVec
	ValueTypeArr
	ValueTypeObj
	ValueTypeMap

	ValueTypeStr
	ValueTypeBytes
	ValueTypeBool

	ValueTypeI
	ValueTypeI8
	ValueTypeI16
	ValueTypeI32
	ValueTypeI64
	ValueTypeU
	ValueTypeU8
	ValueTypeU16
	ValueTypeU32
	ValueTypeU64

	ValueTypeF32
	ValueTypeF64

	ValueTypeBigint
	ValueTypeDatetime
	ValueTypeDate
	ValueTypeTime

	ValueTypeUuid
	ValueTypeDecimal
	ValueTypeIp
	ValueTypeUrl
	ValueTypeEmail

	ValueTypeEnum

	ValueTypeImage
	ValueTypeVideo

	vtUnknownStr = "unknown"

	vtDocStr = "doc"
	vtArrStr = "arr"
	vtVecStr = "vec"
	vtObjStr = "obj"
	vtMapStr = "map"

	vtStrStr   = "str"
	vtBytesStr = "bytes"
	vtBoolStr  = "bool"

	vtIStr   = "i"
	vtI8Str  = "i8"
	vtI16Str = "i16"
	vtI32Str = "i32"
	vtI64Str = "i64"
	vtUStr   = "u"
	vtU8Str  = "u8"
	vtU16Str = "u16"
	vtU32Str = "u32"
	vtU64Str = "u64"

	vtF32Str = "f32"
	vtF64Str = "f64"

	vtBigintStr   = "bigint"
	vtDatetimeStr = "datetime"
	vtDateStr     = "date"
	vtTimeStr     = "time"
	vtUuidStr     = "uuid"
	vtDecimalStr  = "decimal"
	vtIpStr       = "ip"
	vtUrlStr      = "url"
	vtEmailStr    = "email"

	vtEnumStr = "enum"

	vtImageStr = "image"
	vtVideoStr = "video"
)

func (v ValueType) String() string {
	switch v {
	case ValueTypeUnknown:
		return vtUnknownStr
	case ValueTypeDoc:
		return vtDocStr
	case ValueTypeArr:
		return vtArrStr
	case ValueTypeVec:
		return vtVecStr
	case ValueTypeObj:
		return vtObjStr
	case ValueTypeMap:
		return vtMapStr
	case ValueTypeStr:
		return vtStrStr
	case ValueTypeBytes:
		return vtBytesStr
	case ValueTypeBool:
		return vtBoolStr
	case ValueTypeI:
		return vtIStr
	case ValueTypeI8:
		return vtI8Str
	case ValueTypeI16:
		return vtI16Str
	case ValueTypeI32:
		return vtI32Str
	case ValueTypeI64:
		return vtI64Str
	case ValueTypeU:
		return vtUStr
	case ValueTypeU8:
		return vtU8Str
	case ValueTypeU16:
		return vtU16Str
	case ValueTypeU32:
		return vtU32Str
	case ValueTypeU64:
		return vtU64Str
	case ValueTypeF32:
		return vtF32Str
	case ValueTypeF64:
		return vtF64Str
	case ValueTypeBigint:
		return vtBigintStr
	case ValueTypeDatetime:
		return vtDatetimeStr
	case ValueTypeDate:
		return vtDateStr
	case ValueTypeTime:
		return vtTimeStr
	case ValueTypeUuid:
		return vtUuidStr
	case ValueTypeDecimal:
		return vtDecimalStr
	case ValueTypeIp:
		return vtIpStr
	case ValueTypeUrl:
		return vtUrlStr
	case ValueTypeEmail:
		return vtEmailStr
	case ValueTypeEnum:
		return vtEnumStr
	case ValueTypeImage:
		return vtImageStr
	case ValueTypeVideo:
		return vtVideoStr
	default:
		return fmt.Sprintf("ValueType(%d)", v)
	}
}

var strToValueType = map[string]ValueType{
	vtUnknownStr:  ValueTypeUnknown,
	vtDocStr:      ValueTypeDoc,
	vtArrStr:      ValueTypeArr,
	vtVecStr:      ValueTypeVec,
	vtObjStr:      ValueTypeObj,
	vtMapStr:      ValueTypeMap,
	vtStrStr:      ValueTypeStr,
	vtBytesStr:    ValueTypeBytes,
	vtBoolStr:     ValueTypeBool,
	vtIStr:        ValueTypeI,
	vtI8Str:       ValueTypeI8,
	vtI16Str:      ValueTypeI16,
	vtI32Str:      ValueTypeI32,
	vtI64Str:      ValueTypeI64,
	vtUStr:        ValueTypeU,
	vtU8Str:       ValueTypeU8,
	vtU16Str:      ValueTypeU16,
	vtU32Str:      ValueTypeU32,
	vtU64Str:      ValueTypeU64,
	vtF32Str:      ValueTypeF32,
	vtF64Str:      ValueTypeF64,
	vtBigintStr:   ValueTypeBigint,
	vtDatetimeStr: ValueTypeDatetime,
	vtDateStr:     ValueTypeDate,
	vtTimeStr:     ValueTypeTime,
	vtUuidStr:     ValueTypeUuid,
	vtDecimalStr:  ValueTypeDecimal,
	vtIpStr:       ValueTypeIp,
	vtUrlStr:      ValueTypeUrl,
	vtEmailStr:    ValueTypeEmail,
	vtEnumStr:     ValueTypeEnum,
	vtImageStr:    ValueTypeImage,
	vtVideoStr:    ValueTypeVideo,
}

func ParseValueType(s string) (ValueType, error) {
	s = strings.ToLower(s)

	vt, ok := strToValueType[s]
	if !ok {
		return ValueTypeUnknown, fmt.Errorf("Invalid ValueType string: %s", s)
	}
	return vt, nil
}
