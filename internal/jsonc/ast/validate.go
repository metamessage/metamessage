package ast

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lizongying/meta-message/internal/utils"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var decimalRegex = regexp.MustCompile(`^-?\d+\.\d+$`)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (t *Tag) ValidateString(val string) (data any, text string, err error) {
	if val == "" {
		if t.AllowEmpty {
			data = val
			text = val
			return
		}
		err = fmt.Errorf("type string not allow empty value %q", val)
		return
	}

	if t.Pattern != "" {
		var re *regexp.Regexp
		re, err = t.GetPattern()
		if err != nil {
			return
		}
		if re != nil {
			if !re.MatchString(val) {
				err = fmt.Errorf("value '%s' does not match pattern %s", val, t.Pattern)
				return
			}
		} else {
			err = fmt.Errorf("pattern error")
			return
		}
	}

	l := utf8.RuneCountInString(val)

	if t.Min != "" {
		var mini int
		mini, err = strconv.Atoi(t.Min)
		if err == nil {
			return
		}
		if l < mini {
			err = fmt.Errorf("string length %d < min %d", l, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int
		maxi, err = strconv.Atoi(t.Max)
		if err == nil {
			return
		}
		if l > maxi {
			err = fmt.Errorf("string length %d > max %d", l, maxi)
			return
		}
	}

	if t.Size != 0 {
		if l != t.Size {
			err = fmt.Errorf("string length %d != size %d", l, t.Size)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type string not support location UTC%d", location)
		return
	}

	data = val
	text = val

	return
}

func (t *Tag) ValidateBytes(val []byte) (data any, text string, err error) {
	l := len(val)

	if l == 0 {
		if t.AllowEmpty {
			data = val
			text = ""
			return
		}
		err = fmt.Errorf("type []byte not allow empty value []byte{}")
		return
	}

	if t.Min != "" {
		var mini int
		mini, err = strconv.Atoi(t.Min)
		if err != nil {
			return
		}
		if l < mini {
			err = fmt.Errorf("[]byte length %d < min %d", l, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int
		maxi, err = strconv.Atoi(t.Max)
		if err != nil {
			return
		}
		if l > maxi {
			err = fmt.Errorf("[]byte length %d > max %d", l, maxi)
			return
		}
	}

	if t.Size != 0 {
		if l != t.Size {
			err = fmt.Errorf("[]byte length %d != size %d", l, t.Size)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type []byte not support location UTC%d", location)
		return
	}

	data = val
	text = base64.StdEncoding.EncodeToString(val)

	return
}

func (t *Tag) ValidateBool(val bool) (data any, text string, err error) {
	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	if t.AllowEmpty {
		err = fmt.Errorf("type bool not support allow empty")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type bool not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatBool(val)

	return
}

func (t *Tag) ValidateSlice(value *[]any) (data any, text string, err error) {
	if utils.GetLocationOffsetHour(t.Location) != 0 {
		err = fmt.Errorf("location offset hour not zero: type slice not supported")
		return
	}
	return
}

func (t *Tag) ValidateArray(value []any) (ok bool, err error) {
	// if t.Size != 0 {
	// 	if l != t.Size {
	// 		return false, fmt.Errorf("string length %d != size %d", l, t.Size)
	// 	}
	// }

	return
}

func (t *Tag) ValidateStruct(value *any) (data any, text string, err error) {
	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type struct not support location UTC%d", location)
		return
	}
	return
}

func (t *Tag) ValidateMap(value *map[string]any) (data any, text string, err error) {
	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type map not support location UTC%d", location)
		return
	}
	return
}

func (t *Tag) ValidateInt(val int) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type int not allow empty value %v", val)
		return
	}

	val64 := int64(val)

	if t.Min != "" {
		var mini int64
		mini, err = strconv.ParseInt(t.Min, 10, BitSize)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as int: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val64, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int64
		maxi, err = strconv.ParseInt(t.Max, 10, BitSize)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as int: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val64, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type int not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.Itoa(val)

	return
}

func (t *Tag) ValidateInt8(val int8) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type int8 not allow empty value %v", val)
		return
	}

	val64 := int64(val)
	if t.Min != "" {
		var mini int64
		mini, err = strconv.ParseInt(t.Min, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as int8: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val64, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int64
		maxi, err = strconv.ParseInt(t.Max, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as int8: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val64, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type int8 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatInt(int64(val), 10)

	return
}

func (t *Tag) ValidateInt16(val int16) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type int16 not allow empty value %v", val)
		return
	}

	val64 := int64(val)

	if t.Min != "" {
		var mini int64
		mini, err = strconv.ParseInt(t.Min, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as int16: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val64, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int64
		maxi, err = strconv.ParseInt(t.Max, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as int16: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val64, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type int16 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatInt(int64(val), 10)

	return
}

func (t *Tag) ValidateInt32(val int32) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type int32 not allow empty value %v", val)
		return
	}

	val64 := int64(val)

	if t.Min != "" {
		var mini int64
		mini, err = strconv.ParseInt(t.Min, 10, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as int32: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val64, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int64
		maxi, err = strconv.ParseInt(t.Max, 10, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as int32: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val64, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type int32 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatInt(int64(val), 10)

	return
}

func (t *Tag) ValidateInt64(val int64) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type int64 not allow empty value %v", val)
		return
	}

	if t.Min != "" {
		var mini int64
		mini, err = strconv.ParseInt(t.Min, 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as int64: %w", err)
			return
		}
		if val < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int64
		maxi, err = strconv.ParseInt(t.Max, 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as int64: %w", err)
			return
		}
		if val > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type int64 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatInt(val, 10)

	return
}

func (t *Tag) ValidateUint(val uint) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type uint not allow empty value %v", val)
		return
	}

	val64 := uint64(val)

	if t.Min != "" {
		var mini uint64
		mini, err = strconv.ParseUint(t.Min, 10, BitSize)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as uint: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi uint64
		maxi, err = strconv.ParseUint(t.Max, 10, BitSize)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as uint: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uint not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatUint(uint64(val), 10)

	return
}

func (t *Tag) ValidateUint8(val uint8) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type uint8 not allow empty value %v", val)
		return
	}

	val64 := uint64(val)

	if t.Min != "" {
		var mini uint64
		mini, err = strconv.ParseUint(t.Min, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as uint8: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi uint64
		maxi, err = strconv.ParseUint(t.Max, 10, 8)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as uint8: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uint8 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatUint(uint64(val), 10)

	return
}

func (t *Tag) ValidateUint16(val uint16) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type uint16 not allow empty value %v", val)
		return
	}

	val64 := uint64(val)

	if t.Min != "" {
		var mini uint64
		mini, err = strconv.ParseUint(t.Min, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as uint16: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi uint64
		maxi, err = strconv.ParseUint(t.Max, 10, 16)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as uint16: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uint16 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatUint(uint64(val), 10)

	return
}

func (t *Tag) ValidateUint32(val uint32) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type uint32 not allow empty value %v", val)
		return
	}

	val64 := uint64(val)

	if t.Min != "" {
		var mini uint64
		mini, err = strconv.ParseUint(t.Min, 10, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as uint32: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi uint64
		maxi, err = strconv.ParseUint(t.Max, 10, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as uint32: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uint32 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatUint(uint64(val), 10)

	return
}

func (t *Tag) ValidateUint64(val uint64) (data any, text string, err error) {
	if val == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type uint64 not allow empty value %v", val)
		return
	}

	if t.Min != "" {
		var mini uint64
		mini, err = strconv.ParseUint(t.Min, 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as uint64: %w", err)
			return
		}
		if val < mini {
			err = fmt.Errorf("value %v is less than the minimum limit %v", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi uint64
		maxi, err = strconv.ParseUint(t.Max, 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as uint64: %w", err)
			return
		}
		if val > maxi {
			err = fmt.Errorf("value %v exceeds the maximum limit %v", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uint64 not support location UTC%d", location)
		return
	}

	data = val
	text = strconv.FormatUint(val, 10)

	return
}

func (t *Tag) ValidateFloat32(val float32) (data any, text string, err error) {
	if val == 0.0 {
		if t.AllowEmpty {
			data = val
			text = "0.0"
			return
		}
		err = fmt.Errorf("type float32 not allow empty value 0.0")
		return
	}

	val64 := float64(val)

	if t.Min != "" {
		var mini float64
		mini, err = strconv.ParseFloat(t.Min, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as float32: %w", err)
			return
		}
		if val64 < mini {
			err = fmt.Errorf("%f < min %f", val64, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi float64
		maxi, err = strconv.ParseFloat(t.Max, 32)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as float32: %w", err)
			return
		}
		if val64 > maxi {
			err = fmt.Errorf("%f > max %f", val64, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type float32 not support location UTC%d", location)
		return
	}

	data = val
	text = utils.FormatFloat32(val)

	return
}

func (t *Tag) ValidateFloat64(val float64) (data any, text string, err error) {
	if val == 0.0 {
		if t.AllowEmpty {
			data = val
			text = "0.0"
			return
		}
		err = fmt.Errorf("type float64not allow empty value 0.0")
		return
	}

	if t.Min != "" {
		var mini float64
		mini, err = strconv.ParseFloat(t.Min, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Min as float64: %w", err)
			return
		}
		if val < mini {
			err = fmt.Errorf("%f < min %f", val, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi float64
		maxi, err = strconv.ParseFloat(t.Max, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse t.Max as float64: %w", err)
			return
		}
		if val > maxi {
			err = fmt.Errorf("%f > max %f", val, maxi)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type float64 not support location UTC%d", location)
		return
	}

	data = val
	text = utils.FormatFloat64(val)

	return
}

func (t *Tag) ValidateBigInt(val big.Int) (data any, text string, err error) {
	if val.Sign() == 0 {
		if t.AllowEmpty {
			data = val
			text = "0"
			return
		}
		err = fmt.Errorf("type big.Int not allow empty value 0")
		return
	}

	if t.Min != "" {
		mini, ok := new(big.Int).SetString(t.Min, 10)
		if !ok {
			err = fmt.Errorf("invalid min %q for big.Int", t.Min)
			return
		}
		if val.Cmp(mini) == -1 {
			err = fmt.Errorf("big.Int length %s < min %s", val.String(), mini.String())
			return
		}
	}

	if t.Max != "" {
		maxi, ok := new(big.Int).SetString(t.Max, 10)
		if !ok {
			err = fmt.Errorf("invalid max %q for big.Int", t.Min)
			return
		}
		if val.Cmp(maxi) == 1 {
			err = fmt.Errorf("big.Int length %s > max %s", val.String(), maxi.String())
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type big.Int not support location UTC%d", location)
		return
	}

	data = val
	text = val.String()

	return
}

func (t *Tag) ValidateDateTime(val time.Time) (data any, text string, err error) {
	location := time.UTC
	if t.Location != nil {
		location = t.Location
	}

	val = val.Truncate(time.Second)
	format := val.In(location).Format(time.DateTime)
	if val.Unix() == 0 {
		if t.AllowEmpty {
			data = val
			text = format
			return
		}
		err = fmt.Errorf("type datetime not allow empty %v", format)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	data = val
	text = format

	return
}

func (t *Tag) ValidateDate(val time.Time) (data any, text string, err error) {
	location := time.UTC
	if t.Location != nil {
		location = t.Location
	}

	val = val.Truncate(time.Second)
	format := val.In(location).Format(time.DateOnly)
	if val.Unix() == 0 {
		if t.AllowEmpty {
			data = val
			text = format
			return
		}
		err = fmt.Errorf("type date not allow empty %v", format)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	data = val
	text = format

	return
}

func (t *Tag) ValidateTime(val time.Time) (data any, text string, err error) {
	location := time.UTC
	if t.Location != nil {
		location = t.Location
	}

	val = val.Truncate(time.Second)
	format := val.In(location).Format(time.TimeOnly)
	if val.Unix() == 0 {
		if t.AllowEmpty {
			data = val
			text = format
			return
		}
		err = fmt.Errorf("type time not allow empty %v", format)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	data = val
	text = format

	return
}

func (t *Tag) ValidateUUID(val string) (data any, text string, err error) {
	if val == "" {
		if t.AllowEmpty {
			data = [16]byte{}
			text = val
			return
		}
		err = fmt.Errorf("type uuid not allow empty value \"\"")
		return
	}

	if !uuidRegex.MatchString(val) {
		err = fmt.Errorf("value '%s' does not match UUID pattern", val)
		return
	}

	var uuid [16]byte
	uuid, err = utils.UUIDStringToBytes(val)
	if err != nil {
		err = fmt.Errorf("invalid uuid: %w", err)
		return
	}

	if t.Version != 0 && t.Version != int((uuid[6]>>4)&0x0F) {
		err = fmt.Errorf("invalid uuid")
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type uuid not support location UTC%d", location)
		return
	}

	data = uuid
	text = val

	return
}

func (t *Tag) ValidateDecimal(val string) (data any, text string, err error) {
	if val == "" {
		if t.AllowEmpty {
			data = val
			text = val
			return
		}
		err = fmt.Errorf("type decimal not allow empty value \"\"")
		return
	}

	if !decimalRegex.MatchString(val) {
		err = fmt.Errorf("invalid decimal %q, must be like \"0.0\"", val)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type decimal not support location UTC%d", location)
		return
	}

	data = val
	text = val

	return
}

func (t *Tag) ValidateIP(val net.IP) (data any, text string, err error) {
	if val.String() == "" {
		if t.AllowEmpty {
			data = val
			text = ""
			return
		}
		err = fmt.Errorf("type ip not allow empty value \"\"")
		return
	}

	if t.Version == 4 {
		if val.To4() == nil {
			err = fmt.Errorf("invalid ip: %s", val.String())
			return
		}
	}

	if t.Version == 6 {
		if val.To4() != nil {
			err = fmt.Errorf("invalid ip: %s", val.String())
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type ip not support location UTC%d", location)
		return
	}

	data = val
	text = val.String()

	return
}

func (t *Tag) ValidateURL(val url.URL) (data any, text string, err error) {
	if val.String() == "" {
		if t.AllowEmpty {
			data = val
			text = ""
			return
		}
		err = fmt.Errorf("type url not allow empty value \"\"")
		return
	}

	if val.Scheme != "http" && val.Scheme != "https" {
		err = fmt.Errorf("invalid url: %s", val.String())
		return
	}

	if val.Host == "" {
		err = fmt.Errorf("invalid url: %s", val.String())
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type url not support location UTC%d", location)
		return
	}

	data = val
	text = val.String()

	return
}

func (t *Tag) ValidateEmail(val string) (data any, text string, err error) {
	if val == "" {
		if t.AllowEmpty {
			data = val
			text = val
			return
		}
		err = fmt.Errorf("type email not allow empty value \"\"")
		return
	}

	if !emailRegex.MatchString(val) {
		err = fmt.Errorf("value '%s' does not match email pattern", val)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type email not support location UTC%d", location)
		return
	}

	data = val
	text = val

	return
}

func (t *Tag) ValidateEnum(val string) (data any, text string, err error) {
	if val == "" {
		if t.AllowEmpty {
			data = -1
			text = val
			return
		}
		err = fmt.Errorf("type enum not allow empty value \"\"")
		return
	}

	enums := strings.Split(t.Enum, "|")
	idx := -1
	for i, s := range enums {
		if strings.TrimSpace(s) == val {
			idx = i
			break
		}
	}

	if idx == -1 {
		err = fmt.Errorf("value '%s' not found in enum: %v", val, enums)
		return
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type enum not support location UTC%d", location)
		return
	}

	data = idx
	text = val

	return
}

func (t *Tag) ValidateImage(val []byte) (data any, text string, err error) {
	l := len(val)

	if l == 0 {
		if t.AllowEmpty {
			data = val
			text = ""
			return
		}
		err = fmt.Errorf("type image not allow empty value []byte{}")
		return
	}

	if t.Min != "" {
		var mini int
		mini, err = strconv.Atoi(t.Min)
		if err != nil {
			return
		}
		if l < mini {
			err = fmt.Errorf("[]byte length %d < min %d", l, mini)
			return
		}
	}

	if t.Max != "" {
		var maxi int
		maxi, err = strconv.Atoi(t.Max)
		if err != nil {
			return
		}
		if l > maxi {
			err = fmt.Errorf("[]byte length %d > max %d", l, maxi)
			return
		}
	}

	if t.Size != 0 {
		if l != t.Size {
			err = fmt.Errorf("[]byte length %d != size %d", l, t.Size)
			return
		}
	}

	if len(t.Desc) > 65535 {
		err = fmt.Errorf("desc length exceeds 65535 bytes")
		return
	}

	location := utils.GetLocationOffsetHour(t.Location)
	if location != 0 {
		err = fmt.Errorf("type image not support location UTC%d", location)
		return
	}

	data = val
	text = base64.StdEncoding.EncodeToString(val)

	return
}
