package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/metamessage/metamessage/internal/utils"
)

const (
	TIsNull  = "is_null"
	TExample = "example"

	TDesc       = "desc"
	TType       = "type"
	TRaw        = "raw"
	TNullable   = "nullable"
	TAllowEmpty = "allow_empty"
	TUnique     = "unique"
	TDefault    = "default"
	TMin        = "min"
	TMax        = "max"
	TSize       = "size"
	TEnum       = "enum"
	TPattern    = "pattern"
	TLocation   = "location"
	TVersion    = "version"
	TMime       = "mime"

	TChildDesc       = "child_desc"
	TChildType       = "child_type"
	TChildRaw        = "child_raw"
	TChildNullable   = "child_nullable"
	TChildAllowEmpty = "child_allow_empty"
	TChildUnique     = "child_unique"
	TChildDefault    = "child_default"
	TChildMin        = "child_min"
	TChildMax        = "child_max"
	TChildSize       = "child_size"
	TChildEnum       = "child_enum"
	TChildPattern    = "child_pattern"
	TChildLocation   = "child_location"
	TChildVersion    = "child_version"
	TChildMime       = "child_mime"
)

type TagKey uint8

const (
	KIsNull  TagKey = 0 << 3
	KExample        = 1 << 3

	KDesc       = 2 << 3
	KType       = 3 << 3
	KRaw        = 4 << 3
	KNullable   = 5 << 3
	KAllowEmpty = 6 << 3
	KUnique     = 7 << 3
	KDefault    = 8 << 3
	KMin        = 9 << 3
	KMax        = 10 << 3
	KSize       = 11 << 3
	KEnum       = 12 << 3
	KPattern    = 13 << 3
	KLocation   = 14 << 3
	KVersion    = 15 << 3
	KMime       = 16 << 3

	KChildDesc       = 17 << 3
	KChildType       = 18 << 3
	KChildRaw        = 19 << 3
	KChildNullable   = 20 << 3
	KChildAllowEmpty = 21 << 3
	KChildUnique     = 22 << 3
	KChildDefault    = 23 << 3
	KChildMin        = 24 << 3
	KChildMax        = 25 << 3
	KChildSize       = 26 << 3
	KChildEnum       = 27 << 3
	KChildPattern    = 28 << 3
	KChildLocation   = 29 << 3
	KChildVersion    = 30 << 3
	KChildMime       = 31 << 3
)

type Tag struct {
	Name string // name=... For parsing only

	IsNull  bool // is_null
	Example bool // example

	Desc       string         // desc=...
	Type       ValueType      // type=...
	Raw        bool           // raw=...
	Nullable   bool           // nullable
	AllowEmpty bool           // allow_empty=...
	Unique     bool           // unique
	Default    string         // default=...
	Min        string         // min=...
	Max        string         // max=...
	Size       int            // size=... default 0
	Enum       string         // enum=...|...
	Pattern    string         // pattern=...
	Location   *time.Location // location=0  for time.Time
	Version    int            // version=4
	Mime       string         // mime=...

	ChildDesc       string         // child_desc=...
	ChildType       ValueType      // child_type=...
	ChildRaw        bool           // child_raw=...
	ChildNullable   bool           // child_nullable
	ChildAllowEmpty bool           // child_allow_empty
	ChildUnique     bool           // child_unique
	ChildDefault    string         // child_default=...
	ChildMin        string         // child_min=...
	ChildMax        string         // child_max=...
	ChildSize       int            // child_size=... default 0
	ChildEnum       string         // child_enum=...|...
	ChildPattern    string         // child_pattern=...
	ChildLocation   *time.Location // child_location=0  for time.Time [-12, +14]
	ChildVersion    int            // child_version=0 for uuid/ip
	ChildMime       string         // child_mime=...
	ChildIsNull     bool
	ChildExample    bool

	ParentDesc       string
	ParentType       ValueType
	ParentRaw        bool
	ParentNullable   bool
	ParentAllowEmpty bool
	ParentUnique     bool
	ParentDefault    string
	ParentMin        string
	ParentMax        string
	ParentSize       int
	ParentEnum       string
	ParentPattern    string
	ParentLocation   *time.Location
	ParentVersion    int
	ParentMime       string
	ParentIsNull     bool
	ParentExample    bool
}

const (
	DefaultVersion int = 0
)

var DefaultLocation *time.Location = time.UTC

func NewTag() *Tag {
	return &Tag{
		Version:        DefaultVersion,
		ChildVersion:   DefaultVersion,
		ParentVersion:  DefaultVersion,
		Location:       DefaultLocation,
		ChildLocation:  DefaultLocation,
		ParentLocation: DefaultLocation,
	}
}

func (t *Tag) GetPattern() (*regexp.Regexp, error) {
	if t.Pattern == "" {
		return nil, nil
	}
	return regexp.Compile(t.Pattern)
}

func (t *Tag) Json() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func (t *Tag) String() string {
	if t == nil {
		return ""
	}

	var b strings.Builder
	first := true
	add := func(s string) {
		if !first {
			b.WriteString("; ")
		}
		b.WriteString(s)
		first = false
	}

	if t.Type != ValueTypeUnknown && t.Type != t.ParentType {
		if t.Type == ValueTypeString ||
			t.Type == ValueTypeInt ||
			t.Type == ValueTypeFloat64 ||
			t.Type == ValueTypeBool ||
			t.Type == ValueTypeStruct ||
			t.Type == ValueTypeSlice {
		} else {
			if t.Type == ValueTypeArray && t.Size > 0 ||
				t.Type == ValueTypeEnum && t.Enum != "" {

			} else {
				add(TType + "=" + t.Type.String())
			}
		}
	}

	if t.Example && !t.ParentExample {
		add(TExample)
	}

	if t.Nullable && t.Nullable != t.ParentNullable {
		if t.IsNull {
			add(TIsNull)
		} else {
			add(TNullable)
		}
	}

	if t.Desc != "" && t.Desc != t.ParentDesc {
		add(TDesc + "=" + strconv.Quote(t.Desc))
	}

	if t.Raw && t.Raw != t.ParentRaw {
		add(TRaw)
	}

	if t.AllowEmpty && t.AllowEmpty != t.ParentAllowEmpty {
		add(TAllowEmpty)
	}

	if t.Unique && t.Unique != t.ParentUnique {
		add(TUnique)
	}

	if t.Default != "" && t.Default != t.ParentDefault {
		add(TDefault + "=" + t.Default)
	}

	if t.Min != "" && t.Min != t.ParentMin {
		add(TMin + "=" + t.Min)
	}

	if t.Max != "" && t.Max != t.ParentMax {
		add(TMax + "=" + t.Max)
	}

	if t.Size != 0 && t.Size != t.ParentSize {
		add(TSize + "=" + strconv.Itoa(t.Size))
	}

	if t.Enum != "" && t.Enum != t.ParentEnum {
		add(TEnum + "=" + t.Enum)
	}

	if t.Pattern != "" && t.Pattern != t.ParentPattern {
		add(TPattern + "=" + t.Pattern)
	}

	locationOffsetHour := utils.GetLocationOffsetHour(t.Location)
	if locationOffsetHour != 0 && t.Location != t.ParentLocation {
		add(TLocation + "=" + strconv.Itoa(locationOffsetHour))
	}

	if t.Version != DefaultVersion && t.Version != t.ParentVersion {
		add(TVersion + "=" + strconv.Itoa(t.Version))
	}

	if t.Mime != "" && t.Mime != t.ParentMime {
		add(TMime + "=" + t.Mime)
	}

	if t.ChildDesc != "" {
		add(TChildDesc + "=" + strconv.Quote(t.ChildDesc))
	}

	if t.ChildType != ValueTypeUnknown {
		if t.ChildType == ValueTypeString ||
			t.ChildType == ValueTypeInt ||
			t.ChildType == ValueTypeFloat64 ||
			t.ChildType == ValueTypeBool ||
			t.ChildType == ValueTypeStruct ||
			t.ChildType == ValueTypeSlice {
		} else {
			if t.ChildType == ValueTypeArray && t.ChildSize > 0 ||
				t.ChildType == ValueTypeEnum && t.ChildEnum != "" {

			} else {
				add(TChildType + "=" + t.ChildType.String())
			}
		}
	}

	if t.ChildRaw {
		add(TChildRaw)
	}

	if t.ChildNullable {
		add(TChildNullable)
	}

	if t.ChildAllowEmpty {
		add(TChildAllowEmpty)
	}

	if t.ChildUnique {
		add(TChildUnique)
	}

	if t.ChildDefault != "" {
		add(TChildDefault + "=" + t.ChildDefault)
	}

	if t.ChildMin != "" {
		add(TChildMin + "=" + t.ChildMin)
	}

	if t.ChildMax != "" {
		add(TChildMax + "=" + t.ChildMax)
	}

	if t.ChildSize != 0 {
		add(TChildSize + "=" + strconv.Itoa(t.ChildSize))
	}

	if t.ChildEnum != "" {
		add(TChildEnum + "=" + t.ChildEnum)
	}

	if t.ChildPattern != "" {
		add(TChildPattern + "=" + t.ChildPattern)
	}

	childLocationOffsetHour := utils.GetLocationOffsetHour(t.ChildLocation)
	if childLocationOffsetHour != 0 && t.ChildLocation != DefaultLocation {
		add(TChildLocation + "=" + strconv.Itoa(childLocationOffsetHour))
	}

	if t.ChildVersion != DefaultVersion {
		add(TChildVersion + "=" + strconv.Itoa(t.ChildVersion))
	}

	if t.ChildMime != "" {
		add(TChildMime + "=" + t.ChildMime)
	}

	return b.String()
}

func (t *Tag) Bytes() []byte {
	var bs bytes.Buffer
	if t.Example && !t.ParentExample {
		b := KExample | 1
		bs.WriteByte(byte(b))
	}

	if t.Nullable && t.Nullable != t.ParentNullable {
		if t.IsNull {
			b := KIsNull | 1
			bs.WriteByte(byte(b))
		} else {
			b := KNullable | 1
			bs.WriteByte(byte(b))
		}
	}

	if t.Desc != "" && t.Desc != t.ParentDesc {
		l := len(t.Desc)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KDesc) | byte(l))
			bs.WriteString(t.Desc)
		case l <= 1<<8:
			bs.WriteByte(byte(KDesc) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Desc)
		case l <= 1<<16:
			bs.WriteByte(byte(KDesc) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Desc)
		default:
			// err = fmt.Errorf("desc too long")
			// return
		}
	}

	if t.Type != ValueTypeUnknown && t.Type != t.ParentType {
		if t.Type == ValueTypeString ||
			t.Type == ValueTypeBytes ||
			t.Type == ValueTypeInt ||
			t.Type == ValueTypeFloat64 ||
			t.Type == ValueTypeBool ||
			t.Type == ValueTypeStruct ||
			t.Type == ValueTypeSlice {
		} else {
			if t.Type == ValueTypeArray && t.Size > 0 ||
				t.Type == ValueTypeEnum && t.Enum != "" {

			} else {
				bs.WriteByte(byte(KType))
				bs.WriteByte(byte(t.Type))
			}
		}
	}

	if t.Raw && t.Raw != t.ParentRaw {
		bs.WriteByte(byte(KRaw | 1))
	}

	if t.AllowEmpty && t.AllowEmpty != t.ParentAllowEmpty {
		bs.WriteByte(byte(KAllowEmpty | 1))
	}

	if t.Unique && t.Unique != t.ParentUnique {
		bs.WriteByte(byte(KUnique | 1))
	}

	if t.Default != "" && t.Default != t.ParentDefault {
		l := len(t.Default)
		if l < 7 {
			bs.WriteByte(byte(KDefault) | byte(l))
			bs.WriteString(t.Default)
		} else {
			bs.WriteByte(byte(KDefault) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Default)
		}
	}

	if t.Min != "" && t.Min != t.ParentMin {
		l := len(t.Min)
		if l < 7 {
			bs.WriteByte(byte(KMin) | byte(l))
			bs.WriteString(t.Min)
		} else {
			bs.WriteByte(byte(KMin) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Min)
		}
	}

	if t.Max != "" && t.Max != t.ParentMax {
		l := len(t.Max)
		if l < 7 {
			bs.WriteByte(byte(KMax) | byte(l))
			bs.WriteString(t.Max)
		} else {
			bs.WriteByte(byte(KMax) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Max)
		}
	}

	if t.Size != 0 && t.Size != t.ParentSize {
		encodeUint64(&bs, KSize, uint64(t.Size))
	}

	if t.Enum != "" && t.Enum != t.ParentEnum {
		l := len(t.Enum)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KEnum) | byte(l))
			bs.WriteString(t.Enum)
		case l <= 1<<8:
			bs.WriteByte(byte(KEnum) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Enum)
		case l <= 1<<16:
			bs.WriteByte(byte(KEnum) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Enum)
		default:
			// err = fmt.Errorf("enum too long")
			// return
		}
	}

	if t.Pattern != "" && t.Pattern != t.ParentPattern {
		l := len(t.Pattern)
		if l < 7 {
			bs.WriteByte(byte(KPattern) | byte(l))
			bs.WriteString(t.Pattern)
		} else {
			bs.WriteByte(byte(KPattern) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Pattern)
		}
	}

	locationOffsetHour := utils.GetLocationOffsetHour(t.Location)
	if locationOffsetHour != 0 && t.Location != t.ParentLocation {
		v := strconv.Itoa(locationOffsetHour)
		bs.WriteByte(byte(KLocation) | byte(len(v)))
		bs.WriteString(v)
	}

	if t.Version != DefaultVersion && t.Version != t.ParentVersion {
		encodeUint64(&bs, KVersion, uint64(t.Version))
	}

	if t.Mime != "" {
		l, _ := ParseMIME(t.Mime)
		if l < 7 {
			bs.WriteByte(byte(KMime) | byte(l))
		} else {
			bs.WriteByte(byte(KMime) | byte(7))
			bs.WriteByte(byte(l))
		}
	}

	if t.ChildDesc != "" {
		l := len(t.ChildDesc)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KChildDesc) | byte(l))
			bs.WriteString(t.ChildDesc)
		case l <= 1<<8:
			bs.WriteByte(byte(KChildDesc) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildDesc)
		case l <= 1<<16:
			bs.WriteByte(byte(KChildDesc) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildDesc)
		default:
			// err = fmt.Errorf("child desc too long")
			// return
		}
	}
	if t.ChildType != ValueTypeUnknown {
		if t.ChildType == ValueTypeString ||
			t.ChildType == ValueTypeInt ||
			t.ChildType == ValueTypeFloat64 ||
			t.ChildType == ValueTypeBool ||
			t.ChildType == ValueTypeStruct ||
			t.ChildType == ValueTypeSlice {
		} else {
			if t.ChildType == ValueTypeArray && t.ChildSize > 0 ||
				t.ChildType == ValueTypeEnum && t.ChildEnum != "" {

			} else {
				bs.WriteByte(byte(KChildType))
				bs.WriteByte(byte(t.ChildType))
			}
		}
	}

	if t.ChildRaw {
		bs.WriteByte(byte(KChildRaw | 1))
	}

	if t.ChildNullable {
		bs.WriteByte(byte(KChildNullable | 1))
	}

	if t.ChildAllowEmpty {
		bs.WriteByte(byte(KChildAllowEmpty | 1))
	}

	if t.ChildUnique {
		bs.WriteByte(byte(KChildUnique | 1))
	}

	if t.ChildDefault != "" {
		l := len(t.ChildDefault)
		if l < 7 {
			bs.WriteByte(byte(KChildDefault) | byte(l))
			bs.WriteString(t.ChildDefault)
		} else {
			bs.WriteByte(byte(KChildDefault) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildDefault)
		}
	}

	if t.ChildMin != "" {
		l := len(t.ChildMin)
		if l < 7 {
			bs.WriteByte(byte(KChildMin) | byte(l))
			bs.WriteString(t.ChildMin)
		} else {
			bs.WriteByte(byte(KChildMin) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildMin)
		}
	}

	if t.ChildMax != "" {
		l := len(t.ChildMax)
		if l < 7 {
			bs.WriteByte(byte(KChildMax) | byte(l))
			bs.WriteString(t.ChildMax)
		} else {
			bs.WriteByte(byte(KChildMax) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildMax)
		}
	}

	if t.ChildSize != 0 {
		encodeUint64(&bs, KChildSize, uint64(t.ChildSize))
	}

	if t.ChildEnum != "" {
		l := len(t.ChildEnum)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KChildEnum) | byte(l))
			bs.WriteString(t.ChildEnum)
		case l <= 1<<8:
			bs.WriteByte(byte(KChildEnum) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildEnum)
		case l <= 1<<16:
			bs.WriteByte(byte(KChildEnum) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildEnum)
		default:
			// err = fmt.Errorf("child enum too long")
			// return
		}
	}

	if t.ChildPattern != "" {
		l := len(t.ChildPattern)
		if l < 7 {
			bs.WriteByte(byte(KChildPattern) | byte(l))
			bs.WriteString(t.ChildPattern)
		} else {
			bs.WriteByte(byte(KChildPattern) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildPattern)
		}
	}

	childLocationOffsetHour := utils.GetLocationOffsetHour(t.ChildLocation)
	if childLocationOffsetHour != 0 {
		v := strconv.Itoa(childLocationOffsetHour)
		bs.WriteByte(byte(KChildLocation) | byte(len(v)))
		bs.WriteString(v)
	}

	if t.ChildVersion != DefaultVersion {
		encodeUint64(&bs, KChildVersion, uint64(t.ChildVersion))
	}

	if t.ChildMime != "" {
		l := len(t.ChildMime)
		if l < 7 {
			bs.WriteByte(byte(KChildMime) | byte(l))
			bs.WriteString(t.ChildMime)
		} else {
			bs.WriteByte(byte(KChildMime) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildMime)
		}
	}

	return bs.Bytes()
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

func encodeUint64(buf *bytes.Buffer, sign TagKey, uv uint64) {
	switch {
	case uv <= Max1Byte:
		sign |= 0
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv))
	case uv <= Max2Byte:
		sign |= 1
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max3Byte:
		sign |= 2
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max4Byte:
		sign |= 3
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 24))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max5Byte:
		sign |= 4
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 32))
		buf.WriteByte(byte(uv >> 24))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max6Byte:
		sign |= 5
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 40))
		buf.WriteByte(byte(uv >> 32))
		buf.WriteByte(byte(uv >> 24))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max7Byte:
		sign |= 6
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 48))
		buf.WriteByte(byte(uv >> 40))
		buf.WriteByte(byte(uv >> 32))
		buf.WriteByte(byte(uv >> 24))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	case uv <= Max8Byte:
		sign |= 7
		buf.WriteByte(byte(sign))
		buf.WriteByte(byte(uv >> 56))
		buf.WriteByte(byte(uv >> 48))
		buf.WriteByte(byte(uv >> 40))
		buf.WriteByte(byte(uv >> 32))
		buf.WriteByte(byte(uv >> 24))
		buf.WriteByte(byte(uv >> 16))
		buf.WriteByte(byte(uv >> 8))
		buf.WriteByte(byte(uv))
	default:
	}
}

func MergeTag(dst *Tag, src *Tag) *Tag {
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	if src.IsNull {
		dst.IsNull = src.IsNull
	}

	if src.Example {
		dst.Example = src.Example
	}

	if src.Desc != "" {
		dst.Desc = src.Desc
	}

	if src.Type != ValueTypeUnknown {
		dst.Type = src.Type
	}

	if src.Raw {
		dst.Raw = true
	}

	if src.Nullable {
		dst.Nullable = true
	}

	if src.AllowEmpty {
		dst.AllowEmpty = true
	}

	if src.Unique {
		dst.Unique = true
	}

	if src.Default != "" {
		dst.Default = src.Default
	}

	if src.Min != "" {
		dst.Min = src.Min
	}

	if src.Max != "" {
		dst.Max = src.Max
	}

	if src.Size != 0 {
		dst.Size = src.Size
	}

	if src.Enum != "" {
		dst.Enum = src.Enum
	}

	if src.Pattern != "" {
		dst.Pattern = src.Pattern
	}

	if utils.GetLocationOffsetHour(src.Location) != 0 {
		dst.Location = src.Location
	}

	if src.Version != DefaultVersion {
		dst.Version = src.Version
	}

	if src.Mime != "" {
		dst.Mime = src.Mime
	}

	if src.ChildDesc != "" {
		dst.ChildDesc = src.ChildDesc
	}

	if src.ChildType != ValueTypeUnknown {
		dst.ChildType = src.ChildType
	}

	if src.ChildRaw {
		dst.ChildRaw = true
	}

	if src.ChildNullable {
		dst.ChildNullable = true
	}

	if src.ChildAllowEmpty {
		dst.ChildAllowEmpty = true
	}

	if src.ChildUnique {
		dst.ChildUnique = true
	}

	if src.ChildDefault != "" {
		dst.ChildDefault = src.ChildDefault
	}

	if src.ChildMin != "" {
		dst.ChildMin = src.ChildMin
	}

	if src.ChildMax != "" {
		dst.ChildMax = src.ChildMax
	}

	if src.ChildSize != 0 {
		dst.ChildSize = src.ChildSize
	}

	if src.ChildEnum != "" {
		dst.ChildEnum = src.ChildEnum
	}

	if src.ChildPattern != "" {
		dst.ChildPattern = src.ChildPattern
	}

	if utils.GetLocationOffsetHour(src.ChildLocation) != 0 {
		dst.ChildLocation = src.ChildLocation
	}

	if src.ChildVersion != DefaultVersion {
		dst.ChildVersion = src.ChildVersion
	}

	if src.ChildMime != "" {
		dst.ChildMime = src.ChildMime
	}

	return dst
}

func ParseMMTag(tag string) (*Tag, error) {
	r := NewTag()
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "//")
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "mm:")
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return r, nil
	}

	parts := splitTag(tag)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var k, v string
		if strings.Contains(p, "=") {
			kv := strings.SplitN(p, "=", 2)
			k = strings.TrimSpace(kv[0])
			v = strings.TrimSpace(kv[1])
		} else {
			k = strings.TrimSpace(p)
			v = ""
		}

		lower := strings.ToLower(k)
		switch lower {
		case TIsNull:
			r.IsNull = true
			r.Nullable = true

		case TExample:
			r.Example = true

		case TDesc:
			r.Desc = v

		case TType:
			t, err := ParseValueType(v)
			if err != nil {
				return nil, fmt.Errorf("parsing failed %v %w", v, err)
			}
			r.Type = t

		case TRaw:
			r.Raw = true

		case TNullable:
			r.Nullable = true

		case TAllowEmpty:
			r.AllowEmpty = true

		case TUnique:
			r.Unique = true

		case TDefault:
			r.Default = v

		case TPattern:
			r.Pattern = v

		case TMin:
			r.Min = v

		case TMax:
			r.Max = v

		case TSize:
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse size value '%v' to uint64. %w", v, err)
			}
			intMax := uint64(math.MaxInt)
			if u > intMax {
				return nil, fmt.Errorf(
					"uint64 value %d exceeds Go's int type max value %d (Go language limit: array/slice length must be int, not uint64)",
					u, intMax,
				)
			}
			r.Size = int(u)

		case TEnum:
			r.Type = ValueTypeEnum
			r.Enum = v

		case TLocation:
			d, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse location %v %w", v, err)
			}
			if d < -12 || d > 14 {
				return nil, fmt.Errorf("location offset hours must be between -12 and +14, got %d", d)
			}
			r.Location = utils.IntToLocation(d)

		case TVersion:
			d, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version %v %w", v, err)
			}
			if d < 1 || d > 10 {
				return nil, fmt.Errorf("version must be between 1 and 10, got %d", d)
			}
			r.Version = d

		case TMime:
			r.Mime = v

		case TChildDesc:
			r.ChildDesc = v

		case TChildType:
			t, err := ParseValueType(v)
			if err != nil {
				return nil, fmt.Errorf("parsing failed %v %w", v, err)
			}
			r.ChildType = t

		case TChildRaw:
			r.ChildRaw = true

		case TChildNullable:
			r.ChildNullable = true

		case TChildAllowEmpty:
			r.ChildAllowEmpty = true

		case TChildUnique:
			r.ChildUnique = true

		case TChildDefault:
			r.ChildDefault = v

		case TChildPattern:
			r.ChildPattern = v

		case TChildMin:
			r.ChildMin = v

		case TChildMax:
			r.ChildMax = v

		case TChildSize:
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse size value '%v' to uint64. %w", v, err)
			}
			intMax := uint64(math.MaxInt)
			if u > intMax {
				return nil, fmt.Errorf(
					"uint64 value %d exceeds Go's int type max value %d (Go language limit: array/slice length must be int, not uint64)",
					u, intMax,
				)
			}
			r.ChildSize = int(u)

		case TChildEnum:
			r.ChildEnum = v
			r.ChildType = ValueTypeEnum

		case TChildLocation:
			d, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse location %v %w", v, err)
			}
			if d < -12 || d > 14 {
				return nil, fmt.Errorf("location offset hours must be between -12 and +14, got %d", d)
			}
			r.ChildLocation = utils.IntToLocation(d)

		case TChildVersion:
			d, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version %v %w", v, err)
			}
			if d < 1 || d > 10 {
				return nil, fmt.Errorf("version must be between 1 and 10, got %d", d)
			}
			r.ChildVersion = d

		case TChildMime:
			r.ChildMime = v

		default:

		}
	}
	return r, nil
}

func splitTag(tag string) []string {
	if tag == "" {
		return nil
	}

	parts := strings.Split(tag, ";")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
