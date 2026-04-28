package ast

import (
	"fmt"
	"strings"
)

type ValueType uint8

const (
	ValueTypeUnknown ValueType = iota

	ValueTypeDoc
	ValueTypeSlice
	ValueTypeArray
	ValueTypeObject
	ValueTypeMap

	ValueTypeString
	ValueTypeBytes
	ValueTypeBool

	ValueTypeInt
	ValueTypeInt8
	ValueTypeInt16
	ValueTypeInt32
	ValueTypeInt64
	ValueTypeUint
	ValueTypeUint8
	ValueTypeUint16
	ValueTypeUint32
	ValueTypeUint64

	ValueTypeFloat32
	ValueTypeFloat64

	ValueTypeBigInt
	ValueTypeDateTime
	ValueTypeDate
	ValueTypeTime

	ValueTypeUUID
	ValueTypeDecimal
	ValueTypeIP
	ValueTypeURL
	ValueTypeEmail

	ValueTypeEnum

	ValueTypeImage
	ValueTypeVideo

	vtUnknownStr = "unknown"

	vtDocStr    = "doc"
	vtArrayStr  = "arr"
	vtSliceStr  = "slice"
	vtObjectStr = "obj"
	vtMapStr    = "map"

	vtStringStr = "str"
	vtBytesStr  = "bytes"
	vtBoolStr   = "bool"

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

	vtBiStr       = "bi"
	vtDateTimeStr = "datetime"
	vtDateStr     = "date"
	vtTimeStr     = "time"
	vtUUIDStr     = "uuid"
	vtDecimalStr  = "decimal"
	vtIPStr       = "ip"
	vtURLStr      = "url"
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
	case ValueTypeArray:
		return vtArrayStr
	case ValueTypeSlice:
		return vtSliceStr
	case ValueTypeObject:
		return vtObjectStr
	case ValueTypeMap:
		return vtMapStr
	case ValueTypeString:
		return vtStringStr
	case ValueTypeBytes:
		return vtBytesStr
	case ValueTypeBool:
		return vtBoolStr
	case ValueTypeInt:
		return vtIStr
	case ValueTypeInt8:
		return vtI8Str
	case ValueTypeInt16:
		return vtI16Str
	case ValueTypeInt32:
		return vtI32Str
	case ValueTypeInt64:
		return vtI64Str
	case ValueTypeUint:
		return vtUStr
	case ValueTypeUint8:
		return vtU8Str
	case ValueTypeUint16:
		return vtU16Str
	case ValueTypeUint32:
		return vtU32Str
	case ValueTypeUint64:
		return vtU64Str
	case ValueTypeFloat32:
		return vtF32Str
	case ValueTypeFloat64:
		return vtF64Str
	case ValueTypeBigInt:
		return vtBiStr
	case ValueTypeDateTime:
		return vtDateTimeStr
	case ValueTypeDate:
		return vtDateStr
	case ValueTypeTime:
		return vtTimeStr
	case ValueTypeUUID:
		return vtUUIDStr
	case ValueTypeDecimal:
		return vtDecimalStr
	case ValueTypeIP:
		return vtIPStr
	case ValueTypeURL:
		return vtURLStr
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
	vtArrayStr:    ValueTypeArray,
	vtSliceStr:    ValueTypeSlice,
	vtObjectStr:   ValueTypeObject,
	vtMapStr:      ValueTypeMap,
	vtStringStr:   ValueTypeString,
	vtBytesStr:    ValueTypeBytes,
	vtBoolStr:     ValueTypeBool,
	vtIStr:        ValueTypeInt,
	vtI8Str:       ValueTypeInt8,
	vtI16Str:      ValueTypeInt16,
	vtI32Str:      ValueTypeInt32,
	vtI64Str:      ValueTypeInt64,
	vtUStr:        ValueTypeUint,
	vtU8Str:       ValueTypeUint8,
	vtU16Str:      ValueTypeUint16,
	vtU32Str:      ValueTypeUint32,
	vtU64Str:      ValueTypeUint64,
	vtF32Str:      ValueTypeFloat32,
	vtF64Str:      ValueTypeFloat64,
	vtBiStr:       ValueTypeBigInt,
	vtDateTimeStr: ValueTypeDateTime,
	vtDateStr:     ValueTypeDate,
	vtTimeStr:     ValueTypeTime,
	vtUUIDStr:     ValueTypeUUID,
	vtDecimalStr:  ValueTypeDecimal,
	vtIPStr:       ValueTypeIP,
	vtURLStr:      ValueTypeURL,
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
