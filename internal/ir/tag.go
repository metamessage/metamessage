package ir

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
	TIsNull     = "is_null"
	TExample    = "example"
	TDeprecated = "deprecated"

	TName       = "name"
	TDesc       = "desc"
	TType       = "type"
	TNullable   = "nullable"
	TAllowEmpty = "allow_empty"
	TUnique     = "unique"
	TDefaultVal = "default_val"
	TMin        = "min"
	TMax        = "max"
	TSize       = "size"
	TEnum       = "enums"
	TPattern    = "pattern"
	TLocation   = "location"
	TVersion    = "version"
	TMime       = "mime"

	TChildDesc       = "child_desc"
	TChildType       = "child_type"
	TChildNullable   = "child_nullable"
	TChildAllowEmpty = "child_allow_empty"
	TChildUnique     = "child_unique"
	TChildDefaultVal = "child_default_val"
	TChildMin        = "child_min"
	TChildMax        = "child_max"
	TChildSize       = "child_size"
	TChildEnums      = "child_enums"
	TChildPattern    = "child_pattern"
	TChildLocation   = "child_location"
	TChildVersion    = "child_version"
	TChildMime       = "child_mime"
)

type TagKey uint8

const (
	KIsNull     TagKey = 0 << 3
	KExample           = 1 << 3
	KDeprecated        = 2 << 3

	KDesc       = 3 << 3
	KType       = 4 << 3
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
	KChildNullable   = 19 << 3
	KChildAllowEmpty = 20 << 3
	KChildUnique     = 21 << 3
	KChildDefaultVal = 22 << 3
	KChildMin        = 23 << 3
	KChildMax        = 24 << 3
	KChildSize       = 25 << 3
	KChildEnums      = 26 << 3
	KChildPattern    = 27 << 3
	KChildLocation   = 28 << 3
	KChildVersion    = 29 << 3
	KChildMime       = 30 << 3

	KMore = 31 << 3
)

type Tag struct {
	Name string // name=... For parsing only

	IsNull  bool // is_null
	Example bool // example

	Desc       string         // desc=...
	Type       ValueType      // type=...
	Deprecated bool           // deprecated
	Nullable   bool           // nullable
	AllowEmpty bool           // allow_empty
	Unique     bool           // unique
	DefaultVal string         // default_val=...
	Min        string         // min=...
	Max        string         // max=...
	Size       int            // size=... default 0
	Enums      string         // enums=...|...
	Pattern    string         // pattern=...
	Location   *time.Location // location=0  for time.Time [-12, +14]
	Version    int            // version=0 for uuid/ip
	Mime       string         // mime=...
	More       int            // more

	ChildDesc       string         // child_desc=...
	ChildType       ValueType      // child_type=...
	ChildNullable   bool           // child_nullable
	ChildAllowEmpty bool           // child_allow_empty
	ChildUnique     bool           // child_unique
	ChildDefaultVal string         // child_default_val=...
	ChildMin        string         // child_min=...
	ChildMax        string         // child_max=...
	ChildSize       int            // child_size=... default 0
	ChildEnums      string         // child_enums=...|...
	ChildPattern    string         // child_pattern=...
	ChildLocation   *time.Location // child_location=0  for time.Time [-12, +14]
	ChildVersion    int            // child_version=0 for uuid/ip
	ChildMime       string         // child_mime=...

	IsInherit bool
}

const (
	DefaultVersion int = 0
)

var DefaultLocation *time.Location = time.UTC

func NewTag() *Tag {
	return &Tag{
		Version:       DefaultVersion,
		ChildVersion:  DefaultVersion,
		Location:      DefaultLocation,
		ChildLocation: DefaultLocation,
	}
}

func (t *Tag) Inherit(tag *Tag) {
	t.IsInherit = true

	if tag.ChildDesc != "" {
		t.Desc = tag.ChildDesc
	}

	if tag.ChildType != ValueTypeUnknown {
		t.Type = tag.ChildType
	}

	if tag.ChildNullable {
		t.Nullable = tag.ChildNullable
	}

	if tag.ChildAllowEmpty {
		t.AllowEmpty = tag.ChildAllowEmpty
	}

	if tag.ChildUnique {
		t.Unique = tag.ChildUnique
	}

	if tag.ChildDefaultVal != "" {
		t.DefaultVal = tag.ChildDefaultVal
	}

	if tag.ChildMin != "" {
		t.Min = tag.ChildMin
	}

	if tag.ChildMax != "" {
		t.Max = tag.ChildMax
	}

	if tag.ChildSize != 0 {
		t.Size = tag.ChildSize
	}

	if tag.ChildEnums != "" {
		t.Enums = tag.ChildEnums
	}

	if tag.ChildPattern != "" {
		t.Pattern = tag.ChildPattern
	}

	if utils.GetLocationOffsetHour(tag.ChildLocation) != 0 {
		t.Location = tag.ChildLocation
	}

	if tag.ChildVersion != DefaultVersion {
		t.Version = tag.ChildVersion
	}

	if tag.ChildMime != "" {
		t.Mime = tag.ChildMime
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

func (t *Tag) ToString() string {
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

	if t.Type != ValueTypeUnknown && !t.IsInherit {
		if t.Type == ValueTypeStr ||
			t.Type == ValueTypeI ||
			t.Type == ValueTypeF64 ||
			t.Type == ValueTypeBool ||
			t.Type == ValueTypeObj ||
			t.Type == ValueTypeVec {
		} else {
			if t.Type == ValueTypeArr && t.Size > 0 ||
				t.Type == ValueTypeEnum && t.Enums != "" {

			} else {
				add(TType + "=" + t.Type.String())
			}
		}
	}

	if t.Example {
		add(TExample)
	}

	if t.IsNull {
		add(TIsNull)
	}

	if t.Nullable && !t.IsInherit {
		if !t.IsNull {
			add(TNullable)
		}
	}

	if t.Desc != "" && !t.IsInherit {
		add(TDesc + "=" + strconv.Quote(t.Desc))
	}

	if t.Deprecated && !t.IsInherit {
		add("deprecated")
	}

	if t.AllowEmpty && !t.IsInherit {
		add(TAllowEmpty)
	}

	if t.Unique && !t.IsInherit {
		add(TUnique)
	}

	if t.DefaultVal != "" && !t.IsInherit {
		add(TDefaultVal + "=" + t.DefaultVal)
	}

	if t.Min != "" && !t.IsInherit {
		add(TMin + "=" + t.Min)
	}

	if t.Max != "" && !t.IsInherit {
		add(TMax + "=" + t.Max)
	}

	if t.Size != 0 && !t.IsInherit {
		add(TSize + "=" + strconv.Itoa(t.Size))
	}

	if t.Enums != "" && !t.IsInherit {
		add(TEnum + "=" + t.Enums)
	}

	if t.Pattern != "" && !t.IsInherit {
		add(TPattern + "=" + t.Pattern)
	}

	locationOffsetHour := utils.GetLocationOffsetHour(t.Location)
	if locationOffsetHour != 0 && !t.IsInherit {
		add(TLocation + "=" + strconv.Itoa(locationOffsetHour))
	}

	if t.Version != DefaultVersion && !t.IsInherit {
		add(TVersion + "=" + strconv.Itoa(t.Version))
	}

	if t.Mime != "" && !t.IsInherit {
		add(TMime + "=" + t.Mime)
	}

	if t.ChildDesc != "" {
		add(TChildDesc + "=" + strconv.Quote(t.ChildDesc))
	}

	if t.ChildType != ValueTypeUnknown {
		if t.ChildType == ValueTypeStr ||
			t.ChildType == ValueTypeI ||
			t.ChildType == ValueTypeF64 ||
			t.ChildType == ValueTypeBool ||
			t.ChildType == ValueTypeObj ||
			t.ChildType == ValueTypeVec {
		} else {
			if t.ChildType == ValueTypeArr && t.ChildSize > 0 ||
				t.ChildType == ValueTypeEnum && t.ChildEnums != "" {

			} else {
				add(TChildType + "=" + t.ChildType.String())
			}
		}
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

	if t.ChildDefaultVal != "" {
		add(TChildDefaultVal + "=" + t.ChildDefaultVal)
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

	if t.ChildEnums != "" {
		add(TChildEnums + "=" + t.ChildEnums)
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
	if t.Example {
		bs.WriteByte(byte(KExample | 1))
	}

	if t.IsNull {
		bs.WriteByte(byte(KIsNull | 1))
	}

	if t.Nullable && !t.IsInherit {
		if !t.IsNull {
			bs.WriteByte(byte(KNullable | 1))
		}
	}

	if t.Desc != "" && !t.IsInherit {
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

	if t.Type != ValueTypeUnknown && !t.IsInherit {
		if t.Type == ValueTypeStr ||
			t.Type == ValueTypeBytes ||
			t.Type == ValueTypeI ||
			t.Type == ValueTypeF64 ||
			t.Type == ValueTypeBool ||
			t.Type == ValueTypeObj ||
			t.Type == ValueTypeVec {
		} else {
			if t.Type == ValueTypeArr && t.Size > 0 ||
				t.Type == ValueTypeEnum && t.Enums != "" {

			} else {
				bs.WriteByte(byte(KType))
				bs.WriteByte(byte(t.Type))
			}
		}
	}

	if t.Deprecated && !t.IsInherit {
		bs.WriteByte(byte(KDeprecated | 1))
	}

	if t.AllowEmpty && !t.IsInherit {
		bs.WriteByte(byte(KAllowEmpty | 1))
	}

	if t.Unique && !t.IsInherit {
		bs.WriteByte(byte(KUnique | 1))
	}

	if t.DefaultVal != "" && !t.IsInherit {
		l := len(t.DefaultVal)
		if l < 7 {
			bs.WriteByte(byte(KDefault) | byte(l))
			bs.WriteString(t.DefaultVal)
		} else {
			bs.WriteByte(byte(KDefault) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.DefaultVal)
		}
	}

	if t.Min != "" && !t.IsInherit {
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

	if t.Max != "" && !t.IsInherit {
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

	if t.Size != 0 && !t.IsInherit {
		encodeU64(&bs, KSize, uint64(t.Size))
	}

	if t.Enums != "" && !t.IsInherit {
		l := len(t.Enums)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KEnum) | byte(l))
			bs.WriteString(t.Enums)
		case l <= 1<<8:
			bs.WriteByte(byte(KEnum) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Enums)
		case l <= 1<<16:
			bs.WriteByte(byte(KEnum) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.Enums)
		default:
			// err = fmt.Errorf("enums too long")
			// return
		}
	}

	if t.Pattern != "" && !t.IsInherit {
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
	if locationOffsetHour != 0 && !t.IsInherit {
		v := strconv.Itoa(locationOffsetHour)
		bs.WriteByte(byte(KLocation) | byte(len(v)))
		bs.WriteString(v)
	}

	if t.Version != DefaultVersion && !t.IsInherit {
		encodeU64(&bs, KVersion, uint64(t.Version))
	}

	if t.Mime != "" && !t.IsInherit {
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
		if t.ChildType == ValueTypeStr ||
			t.ChildType == ValueTypeI ||
			t.ChildType == ValueTypeF64 ||
			t.ChildType == ValueTypeBool ||
			t.ChildType == ValueTypeObj ||
			t.ChildType == ValueTypeVec {
		} else {
			if t.ChildType == ValueTypeArr && t.ChildSize > 0 ||
				t.ChildType == ValueTypeEnum && t.ChildEnums != "" {

			} else {
				bs.WriteByte(byte(KChildType))
				bs.WriteByte(byte(t.ChildType))
			}
		}
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

	if t.ChildDefaultVal != "" {
		l := len(t.ChildDefaultVal)
		if l < 7 {
			bs.WriteByte(byte(KChildDefaultVal) | byte(l))
			bs.WriteString(t.ChildDefaultVal)
		} else {
			bs.WriteByte(byte(KChildDefaultVal) | byte(7))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildDefaultVal)
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
		encodeU64(&bs, KChildSize, uint64(t.ChildSize))
	}

	if t.ChildEnums != "" {
		l := len(t.ChildEnums)
		switch {
		case l <= 5:
			bs.WriteByte(byte(KChildEnums) | byte(l))
			bs.WriteString(t.ChildEnums)
		case l <= 1<<8:
			bs.WriteByte(byte(KChildEnums) | byte(6))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildEnums)
		case l <= 1<<16:
			bs.WriteByte(byte(KChildEnums) | byte(7))
			bs.WriteByte(byte(l >> 8))
			bs.WriteByte(byte(l))
			bs.WriteString(t.ChildEnums)
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
		encodeU64(&bs, KChildVersion, uint64(t.ChildVersion))
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

	if t.More != 0 {
		encodeU64(&bs, KMore, uint64(t.More))
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

func encodeU64(buf *bytes.Buffer, sign TagKey, uv uint64) {
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

	if src.Deprecated {
		dst.Deprecated = true
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

	if src.DefaultVal != "" {
		dst.DefaultVal = src.DefaultVal
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

	if src.Enums != "" {
		dst.Enums = src.Enums
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

	if src.ChildNullable {
		dst.ChildNullable = true
	}

	if src.ChildAllowEmpty {
		dst.ChildAllowEmpty = true
	}

	if src.ChildUnique {
		dst.ChildUnique = true
	}

	if src.ChildDefaultVal != "" {
		dst.ChildDefaultVal = src.ChildDefaultVal
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

	if src.ChildEnums != "" {
		dst.ChildEnums = src.ChildEnums
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

		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			if s, err := strconv.Unquote(v); err == nil {
				v = s
			}
		}

		lower := strings.ToLower(k)
		switch lower {
		case TName:
			r.Name = v

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

		case TDeprecated:
			r.Deprecated = true

		case TNullable:
			r.Nullable = true

		case TAllowEmpty:
			r.AllowEmpty = true

		case TUnique:
			r.Unique = true

		case TDefaultVal:
			r.DefaultVal = v

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
			r.Enums = v

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

		case TChildNullable:
			r.ChildNullable = true

		case TChildAllowEmpty:
			r.ChildAllowEmpty = true

		case TChildUnique:
			r.ChildUnique = true

		case TChildDefaultVal:
			r.ChildDefaultVal = v

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

		case TChildEnums:
			r.ChildEnums = v
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
