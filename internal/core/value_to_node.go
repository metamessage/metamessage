package core

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/utils"
)

const maxDepth = 32

func NilToNode(valueType ir.ValueType) (*ir.Value, error) {
	tag := ir.NewTag()
	tag.Type = valueType
	return &ir.Value{
		Data: nil,
		Text: ir.Null,
		Tag:  tag,
	}, nil
}

func resolveDefaultValue(tag *ir.Tag) (data any, text string, ok bool) {
	if tag.DefaultVal == "" {
		return nil, "", false
	}
	text = tag.DefaultVal
	switch tag.Type {
	case ir.ValueTypeBytes, ir.ValueTypeMedia:
		data = []byte(tag.DefaultVal)
	case ir.ValueTypeBool:
		if tag.DefaultVal == "true" || tag.DefaultVal == "1" {
			data = true
			text = ir.True
		} else {
			data = false
			text = ir.False
		}
	case ir.ValueTypeI:
		if v, err := strconv.Atoi(tag.DefaultVal); err == nil {
			data = v
		}
	case ir.ValueTypeI8:
		if v, err := strconv.ParseInt(tag.DefaultVal, 10, 8); err == nil {
			data = int8(v)
		}
	case ir.ValueTypeI16:
		if v, err := strconv.ParseInt(tag.DefaultVal, 10, 16); err == nil {
			data = int16(v)
		}
	case ir.ValueTypeI32:
		if v, err := strconv.ParseInt(tag.DefaultVal, 10, 32); err == nil {
			data = int32(v)
		}
	case ir.ValueTypeI64:
		if v, err := strconv.ParseInt(tag.DefaultVal, 10, 64); err == nil {
			data = v
		}
	case ir.ValueTypeU:
		if v, err := strconv.ParseUint(tag.DefaultVal, 10, 64); err == nil {
			data = uint(v)
		}
	case ir.ValueTypeU8:
		if v, err := strconv.ParseUint(tag.DefaultVal, 10, 8); err == nil {
			data = uint8(v)
		}
	case ir.ValueTypeU16:
		if v, err := strconv.ParseUint(tag.DefaultVal, 10, 16); err == nil {
			data = uint16(v)
		}
	case ir.ValueTypeU32:
		if v, err := strconv.ParseUint(tag.DefaultVal, 10, 32); err == nil {
			data = uint32(v)
		}
	case ir.ValueTypeU64:
		if v, err := strconv.ParseUint(tag.DefaultVal, 10, 64); err == nil {
			data = v
		}
	case ir.ValueTypeF32:
		if v, err := strconv.ParseFloat(tag.DefaultVal, 32); err == nil {
			data = float32(v)
		}
	case ir.ValueTypeF64:
		if v, err := strconv.ParseFloat(tag.DefaultVal, 64); err == nil {
			data = v
		}
	case ir.ValueTypeStr, ir.ValueTypeDecimal, ir.ValueTypeEmail, ir.ValueTypeUrl, ir.ValueTypeIp:
		data = tag.DefaultVal
	case ir.ValueTypeEnum:
		data = -1
	case ir.ValueTypeUuid:
		data = [16]byte{}
	case ir.ValueTypeBigint:
		if v, ok := new(big.Int).SetString(tag.DefaultVal, 10); ok {
			data = *v
		}
	case ir.ValueTypeDatetime:
		if t, err := time.Parse(time.DateTime, tag.DefaultVal); err == nil {
			data = t
		}
	case ir.ValueTypeDate:
		if t, err := time.Parse(time.DateOnly, tag.DefaultVal); err == nil {
			data = t
		}
	case ir.ValueTypeTime:
		if t, err := time.Parse(time.TimeOnly, tag.DefaultVal); err == nil {
			data = t
		}
	default:
		data = tag.DefaultVal
	}
	if data == nil {
		data = tag.DefaultVal
	}
	return data, text, true
}

func ValueToNode(v any, tagStr string) (node ir.Node, err error) {
	if v == nil {
		return &ir.NodeNull{}, nil
	}

	var tag *ir.Tag
	if tagStr != "" {
		if tag, err = ir.ParseMMTag(tagStr); err != nil {
			return nil, fmt.Errorf("parse tag failed: %w", err)
		}
	}

	if tag == nil {
		tag = ir.NewTag()
	}

	example := tag.Example

	return valueToNode(v, tag, 0, "", example)
}

func valueToNode(v any, tag *ir.Tag, depth int, path string, example bool) (node ir.Node, err error) {
	if tag == nil {
		tag = ir.NewTag()
	}

	var data any
	var text string
	text = ir.Null
	switch val := v.(type) {
	case nil:
		// var val interface{} = nil
		// val.IsValid()

		if tag.Type == ir.ValueTypeUnknown {
			return nil, fmt.Errorf("invalid input: v is untyped nil (no concrete type/value)")
		}
		tag.IsNull = true

	case []byte:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeVec
		}

		switch tag.Type {
		case ir.ValueTypeBytes:
			data, text, err = tag.ValidateBytes(val, example)

		case ir.ValueTypeMedia:
			data, text, err = tag.ValidateMedia(val, example)

		case ir.ValueTypeVec:
			return anyToJSONC(v, tag, depth, path, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *[]byte:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeVec
		}

		switch tag.Type {
		case ir.ValueTypeBytes:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = []byte{}
					text = ""
				}
			} else {
				data, text, err = tag.ValidateBytes(*val, example)
			}

		case ir.ValueTypeMedia:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = []byte{}
					text = ""
				}
			} else {
				data, text, err = tag.ValidateMedia(*val, example)
			}

		case ir.ValueTypeVec:
			return anyToJSONC(v, tag, depth, path, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case bool:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBool
		}

		switch tag.Type {
		case ir.ValueTypeBool:
			data, text, err = tag.ValidateBool(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *bool:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBool
		}

		switch tag.Type {
		case ir.ValueTypeBool:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = false
					text = ir.False
				}
			} else {
				data, text, err = tag.ValidateBool(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI
		}

		switch tag.Type {
		case ir.ValueTypeI:
			data, text, err = tag.ValidateI(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI
		}

		switch tag.Type {
		case ir.ValueTypeI:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = 0
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateI(*val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int8:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI8
		}

		switch tag.Type {
		case ir.ValueTypeI8:
			data, text, err = tag.ValidateI8(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int8:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI8
		}

		switch tag.Type {
		case ir.ValueTypeI8:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = int8(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateI8(*val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int16:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI16
		}

		switch tag.Type {
		case ir.ValueTypeI16:
			data, text, err = tag.ValidateI16(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int16:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI16
		}

		switch tag.Type {
		case ir.ValueTypeI16:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = int16(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateI16(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI32
		}

		switch tag.Type {
		case ir.ValueTypeI32:
			data, text, err = tag.ValidateI32(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI32
		}

		switch tag.Type {
		case ir.ValueTypeI32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = int32(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateI32(*val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI64
		}

		switch tag.Type {
		case ir.ValueTypeI64:
			data, text, err = tag.ValidateI64(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeI64
		}

		switch tag.Type {
		case ir.ValueTypeI64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = int64(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateI64(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU
		}

		switch tag.Type {
		case ir.ValueTypeU:
			data, text, err = tag.ValidateU(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU
		}

		switch tag.Type {
		case ir.ValueTypeU:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = uint(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateU(*val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint8:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU8
		}

		switch tag.Type {
		case ir.ValueTypeU8:
			data, text, err = tag.ValidateU8(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint8:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU8
		}

		switch tag.Type {
		case ir.ValueTypeU8:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = uint8(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateU8(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint16:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU16
		}

		switch tag.Type {
		case ir.ValueTypeU16:
			data, text, err = tag.ValidateU16(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint16:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU16
		}

		switch tag.Type {
		case ir.ValueTypeU16:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = uint16(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateU16(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU32
		}
		switch tag.Type {
		case ir.ValueTypeU32:
			data, text, err = tag.ValidateU32(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU32
		}

		switch tag.Type {
		case ir.ValueTypeU32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = uint32(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateU32(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU64
		}

		switch tag.Type {
		case ir.ValueTypeU64:
			data, text, err = tag.ValidateU64(val, example)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeU64
		}

		switch tag.Type {
		case ir.ValueTypeU64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = uint64(0)
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateU64(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case float32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeF32
		}

		switch tag.Type {
		case ir.ValueTypeF32:
			if math.IsInf(float64(val), +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(float64(val), -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(float64(val)) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateF32(val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *float32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeF32
		}

		switch tag.Type {
		case ir.ValueTypeF32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = float32(0)
					text = "0.0"
				}
			} else if math.IsInf(float64(*val), +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(float64(*val), -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(float64(*val)) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateF32(*val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case float64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeF64
		}

		switch tag.Type {
		case ir.ValueTypeF64:
			if math.IsInf(val, +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(val, -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(val) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateF64(val, example)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *float64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeF64
		}

		switch tag.Type {
		case ir.ValueTypeF64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = float64(0)
					text = "0.0"
				}
			} else if math.IsInf(*val, +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(*val, -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(*val) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateF64(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case string:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeStr
		}

		switch tag.Type {
		case ir.ValueTypeStr:
			data, text, err = tag.ValidateStr(val, example)

		case ir.ValueTypeDecimal:
			data, text, err = tag.ValidateDecimal(val, example)

		case ir.ValueTypeEmail:
			data, text, err = tag.ValidateEmail(val, example)

		case ir.ValueTypeEnum:
			data, text, err = tag.ValidateEnum(val, example)

		case ir.ValueTypeUuid:
			data, text, err = tag.ValidateUuid(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *string:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeStr
		}

		switch tag.Type {
		case ir.ValueTypeStr:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = ""
					text = ""
				}
			} else {
				data, text, err = tag.ValidateStr(*val, example)
			}

		case ir.ValueTypeDecimal:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = ""
					text = ""
				}
			} else {
				data, text, err = tag.ValidateDecimal(*val, example)
			}

		case ir.ValueTypeEmail:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = ""
					text = ""
				}
			} else {
				data, text, err = tag.ValidateEmail(*val, example)
			}

		case ir.ValueTypeEnum:
			if tag.Enums == "" {
				err = fmt.Errorf("enum empty")
				return
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = -1
					text = ""
				}
			} else {
				data, text, err = tag.ValidateEnum(*val, example)
			}

		case ir.ValueTypeUuid:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = [16]byte{}
					text = ""
				}
			} else {
				data, text, err = tag.ValidateUuid(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case big.Int:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBigint
		}

		switch tag.Type {
		case ir.ValueTypeBigint:
			data, text, err = tag.ValidateBigint(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *big.Int:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBigint
		}

		switch tag.Type {
		case ir.ValueTypeBigint:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = big.Int{}
					text = "0"
				}
			} else {
				data, text, err = tag.ValidateBigint(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case url.URL:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUrl
		}

		switch tag.Type {
		case ir.ValueTypeUrl:
			data, text, err = tag.ValidateUrl(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *url.URL:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUrl
		}

		switch tag.Type {
		case ir.ValueTypeUrl:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = url.URL{}
					text = ""
				}
			} else {
				data, text, err = tag.ValidateUrl(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case net.IP:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeIp
		}

		switch tag.Type {
		case ir.ValueTypeIp:
			data, text, err = tag.ValidateIp(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *net.IP:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeIp
		}

		switch tag.Type {
		case ir.ValueTypeIp:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = net.IP{}
					text = ""
				}
			} else {
				data, text, err = tag.ValidateIp(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case time.Time:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeDatetime
		}

		switch tag.Type {
		case ir.ValueTypeDatetime:
			data, text, err = tag.ValidateDatetime(val, example)

		case ir.ValueTypeDate:
			data, text, err = tag.ValidateDate(val, example)

		case ir.ValueTypeTime:
			data, text, err = tag.ValidateTime(val, example)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *time.Time:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeDatetime
		}

		switch tag.Type {
		case ir.ValueTypeDatetime:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = utils.DefaultTime
					text = utils.DefaultTime.In(location).Format(time.DateTime)
				}
			} else {
				data, text, err = tag.ValidateTime(*val, example)
			}

		case ir.ValueTypeDate:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = utils.DefaultTime
					text = utils.DefaultTime.In(location).Format(time.DateOnly)
				}
			} else {
				data, text, err = tag.ValidateTime(*val, example)
			}

		case ir.ValueTypeTime:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("%s: value is nil and not nullable", path)
					return
				}

				if !example {
					tag.IsNull = true
				}

				if d, t, ok := resolveDefaultValue(tag); ok {
					data = d
					text = t
				} else {
					data = utils.DefaultTime
					text = utils.DefaultTime.In(location).Format(time.TimeOnly)
				}
			} else {
				data, text, err = tag.ValidateTime(*val, example)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	default:
		return anyToJSONC(v, tag, depth, path, example)
	}

	if err != nil {
		// err = fmt.Errorf("validate failed: %w", err)
		return
	}

	return &ir.Value{
		Data: data,
		Text: text,
		Tag:  tag,
		Path: path,
	}, nil
}

// Create sample values according to the element Type (special for empty slices)
func createExampleValue(elemType reflect.Type) (any, error) {
	if elemType.Kind() == reflect.Pointer {
		base := elemType.Elem()
		baseVal, err := createExampleValue(base)
		if err != nil {
			return nil, err
		}
		ptr := reflect.New(base)
		ptr.Elem().Set(reflect.ValueOf(baseVal))
		return ptr.Interface(), nil
	}

	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return reflect.Zero(elemType).Interface(), nil

	case reflect.String:
		return "", nil

	case reflect.Bool:
		return false, nil

	case reflect.Slice:
		// Recursively create sample values for other slices (e.g., []int → []int{0})
		sliceVal := reflect.MakeSlice(elemType, 1, 1)
		elemExample, err := createExampleValue(elemType.Elem())
		if err != nil {
			return nil, err
		}
		sliceVal.Index(0).Set(reflect.ValueOf(elemExample))
		return sliceVal.Interface(), nil

	case reflect.Array:
		arrayVal := reflect.New(elemType).Elem()
		if elemType.Len() > 0 {
			elemExample, err := createExampleValue(elemType.Elem())
			if err != nil {
				return nil, err
			}
			arrayVal.Index(0).Set(reflect.ValueOf(elemExample))
		}
		return arrayVal.Interface(), nil

	case reflect.Struct:
		switch elemType {
		case reflect.TypeFor[time.Time]():
			return utils.DefaultTime, nil

		case reflect.TypeFor[big.Int]():
			return big.NewInt(0), nil

		case reflect.TypeFor[net.IP]():
			return net.IP{}, nil

		case reflect.TypeFor[url.URL]():
			return url.URL{}, nil

		default:
			structVal := reflect.New(elemType).Elem()
			for i := 0; i < elemType.NumField(); i++ {
				field := elemType.Field(i)
				if !field.IsExported() {
					continue
				}

				fieldVal, err := createExampleValue(field.Type)
				if err != nil {
					return nil, fmt.Errorf("struct %s field %s: %w", elemType.Name(), field.Name, err)
				}

				val := reflect.ValueOf(fieldVal)
				if structVal.Field(i).CanSet() && val.IsValid() {
					structVal.Field(i).Set(val)
				}
			}
			return structVal.Interface(), nil
		}

	case reflect.Map:
		mapVal := reflect.MakeMap(elemType)
		keyExample, err := createExampleValue(elemType.Key())
		if err != nil {
			return nil, fmt.Errorf("create map key example: %w", err)
		}
		valExample, err := createExampleValue(elemType.Elem())
		if err != nil {
			return nil, fmt.Errorf("create map val example: %w", err)
		}
		mapVal.SetMapIndex(reflect.ValueOf(keyExample), reflect.ValueOf(valExample))
		return mapVal.Interface(), nil

	default:
		return nil, fmt.Errorf("createExampleValue unsupported example type: %s", elemType.Kind())
	}
}

func anyToJSONC(obj any, tag *ir.Tag, depth int, path string, example bool) (ir.Node, error) {
	depth++
	if depth > maxDepth {
		return nil, fmt.Errorf("max depth: %d", maxDepth)
	}

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if val.Kind() == reflect.Interface {
		return nil, fmt.Errorf("unsupported type: interface{} (cannot reconstruct concrete value from abstract interface)")
	}

	if tag == nil {
		tag = ir.NewTag()
	}

	if val.Kind() == reflect.Pointer {
		tag.Nullable = true
		if val.IsNil() {
			if !example {
				tag.IsNull = true
			}
			typ = typ.Elem()
			example = true
			elemVal, err := createExampleValue(typ)
			if err != nil {
				return nil, fmt.Errorf("%s create element %s err: %w", path, typ, err)
			}
			val = reflect.ValueOf(elemVal)
		} else {
			val = val.Elem()
			typ = typ.Elem()
		}
		if val.Kind() == reflect.Interface {
			return nil, fmt.Errorf("%s unsupported interface{} (cannot reconstruct concrete value from abstract interface)", path)
		}
		if val.Kind() == reflect.Pointer {
			return nil, fmt.Errorf("%s unsupported multi-level pointer (%s) (only single-level pointer is allowed)", path, typ)
		}
	}

	var err error

	if tag.ToString() == "" {
		mmMethod, hasMM := typ.MethodByName("MM")
		if hasMM && mmMethod.Type.NumIn() == 1 && mmMethod.Type.NumOut() == 1 {
			ret := mmMethod.Func.Call([]reflect.Value{val})
			if len(ret) > 0 {
				mmTag := ret[0].String()
				if mmTag != "" {
					var tagNode *ir.Tag
					if tagNode, err = ir.ParseMMTag(mmTag); err != nil {
						return nil, fmt.Errorf("parse mm tag for struct %s: %w", tag.Name, err)
					} else {
						tag = ir.MergeTag(tag, tagNode)
					}
				}
			}
		}
	}

	tag.Type = ir.ValueTypeObj
	tag.Name = utils.CamelToSnake(typ.Name())
	if tag.Name != "" {
		if path == "" {
			path = tag.Name
		} else {
			path = fmt.Sprintf("%s.%s", path, tag.Name)
		}
	}

	switch val.Kind() {
	case reflect.Struct:
		switch val.Type() {
		default:
			tag.Type = ir.ValueTypeObj

			objNode := &ir.Object{
				Tag:  tag,
				Path: path,
			}

			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				fieldVal := val.Field(i)

				if field.PkgPath != "" {
					continue
				}

				fieldKey := utils.CamelToSnake(field.Name)
				mmTagStr := field.Tag.Get("mm")
				var tagField *ir.Tag
				if mmTagStr != "" {
					if tagField, err = ir.ParseMMTag(mmTagStr); err != nil {
						return nil, fmt.Errorf("parse mm tag for field %s: %w", fieldKey, err)
					} else {
						if tagField != nil {
							if tagField.Name != "" {
								if tagField.Name == "-" {
									continue
								}
								fieldKey = tagField.Name
							}
						}
					}
				}

				if tagField == nil {
					tagField = ir.NewTag()
				}

				if tagField.Name == "" {
					tagField.Name = fieldKey
				}

				p := fmt.Sprintf("%s.%s", path, fieldKey)
				fieldNode, err := valueToNode(fieldVal.Interface(), tagField, depth, p, example)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", p, err)
				}

				objNode.Fields = append(objNode.Fields, &ir.Field{
					Key:   fieldKey,
					Value: fieldNode,
				})
			}

			if !example {
				err = tag.ValidateObj()
				if err != nil {
					err = fmt.Errorf("validate failed: %w", err)
					return nil, err
				}
			}

			return objNode, nil
		}

	case reflect.Map:
		mapValueType := val.Type().Elem()

		if mapValueType.Kind() == reflect.Interface {
			return nil, fmt.Errorf("path %q: map value type cannot be any/interface{}, please use concrete type (string/int/bool/struct/map etc.)", path)
		}

		tag.Type = ir.ValueTypeMap

		node := &ir.Object{
			Tag:  tag,
			Path: path,
		}

		setTag := false
		for _, key := range val.MapKeys() {
			keyStr, ok := key.Interface().(string)
			if !ok {
				return nil, fmt.Errorf("map key must be string, got %T", key.Interface())
			}

			keyStr = utils.CamelToSnake(keyStr)

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Name = keyStr

			p := fmt.Sprintf("%s[%s]", path, keyStr)
			valNode, err := valueToNode(val.MapIndex(key).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = valNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefaultVal = tagItem.DefaultVal
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnums = tagItem.Enums
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Fields = append(node.Fields, &ir.Field{
				Key:   keyStr,
				Value: valNode,
			})
		}

		if val.Len() == 0 {
			tag.IsEmpty = true
			keyType := typ.Key()
			if keyType.Kind() != reflect.String {
				return nil, fmt.Errorf("map key must be string, got %T", keyType.Kind())
			}

			var exampleVal any
			elemType := typ.Elem()
			exampleVal, err = createExampleValue(elemType)
			if err != nil {
				return nil, fmt.Errorf("create example value for empty slice: %w", err)
			}

			keyStr := ""

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Name = keyStr

			tagItem.Example = true

			p := fmt.Sprintf("%s[%s]", path, keyStr)
			valNode, err := valueToNode(exampleVal, tagItem, depth, p, true)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = valNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefaultVal = tagItem.DefaultVal
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnums = tagItem.Enums
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Fields = append(node.Fields, &ir.Field{
				Key:   keyStr,
				Value: valNode,
			})
		}

		if !example {
			err = tag.ValidateMap()
			if err != nil {
				err = fmt.Errorf("validate failed: %w", err)
				return nil, err
			}
		}

		return node, nil

	case reflect.Slice:
		tag.Type = ir.ValueTypeVec

		node := &ir.Array{
			Tag:  tag,
			Path: path,
		}

		size := val.Len()
		setTag := false
		for i := range size {
			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			p := fmt.Sprintf("%s[%d]", path, i)
			itemNode, err := valueToNode(val.Index(i).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefaultVal = tagItem.DefaultVal
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnums = tagItem.Enums
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if size == 0 {
			tag.IsEmpty = true
			var exampleVal any
			exampleVal, err = createExampleValue(typ.Elem())
			if err != nil {
				return nil, fmt.Errorf("create example value for empty slice: %w", err)
			}

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Example = true

			p := fmt.Sprintf("%s[%d]", path, 0)
			itemNode, err := valueToNode(exampleVal, tagItem, depth, p, true)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefaultVal = tagItem.DefaultVal
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnums = tagItem.Enums
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if !example {
			err = tag.ValidateVec(node.Items)
			if err != nil {
				err = fmt.Errorf("validate failed: %w", err)
				return nil, err
			}
		}

		return node, nil

	case reflect.Array:
		tag.Type = ir.ValueTypeArr

		tag.Size = val.Len()
		node := &ir.Array{
			Tag:  tag,
			Path: path,
		}

		setTag := false
		for i := 0; i < tag.Size; i++ {
			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			p := fmt.Sprintf("%s[%d]", path, i)
			itemNode, err := valueToNode(val.Index(i).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefaultVal = tagItem.DefaultVal
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnums = tagItem.Enums
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if tag.Size == 0 {
			tag.IsEmpty = true
		}

		if !example {
			err = tag.ValidateArr(node.Items)
			if err != nil {
				err = fmt.Errorf("validate failed: %w", err)
				return nil, err
			}
		}

		return node, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", val.Kind())
	}
}
