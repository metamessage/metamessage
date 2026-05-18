package core

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"net/url"
	"reflect"
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

func ValueToNode(v any, tagStr string) (node ir.Node, err error) {
	var tag *ir.Tag
	if tagStr != "" {
		if tag, err = ir.ParseMMTag(tagStr); err != nil {
			return nil, fmt.Errorf("parse tag failed: %w", err)
		}
	}

	return valueToNode(v, tag, 0, "")
}

func valueToNode(v any, tag *ir.Tag, depth int, path string) (node ir.Node, err error) {
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
			tag.Type = ir.ValueTypeSlice
		}

		switch tag.Type {
		case ir.ValueTypeBytes:
			data, text, err = tag.ValidateBytes(val)

		case ir.ValueTypeImage:
			data, text, err = tag.ValidateImage(val)

		case ir.ValueTypeSlice:
			return anyToJSONC(v, tag, depth, path)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *[]byte:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeSlice
		}

		switch tag.Type {
		case ir.ValueTypeBytes:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = []byte{}
				text = ""
			} else {
				data, text, err = tag.ValidateBytes(*val)
			}

		case ir.ValueTypeImage:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = []byte{}
				text = ""
			} else {
				data, text, err = tag.ValidateImage(*val)
			}

		case ir.ValueTypeSlice:
			return anyToJSONC(v, tag, depth, path)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case bool:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBool
		}

		switch tag.Type {
		case ir.ValueTypeBool:
			data, text, err = tag.ValidateBool(val)

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
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = false
				text = ir.False
			} else {
				data, text, err = tag.ValidateBool(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt
		}

		switch tag.Type {
		case ir.ValueTypeInt:
			data, text, err = tag.ValidateInt(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt
		}

		switch tag.Type {
		case ir.ValueTypeInt:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = 0
				text = "0"
			} else {
				data, text, err = tag.ValidateInt(*val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int8:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt8
		}

		switch tag.Type {
		case ir.ValueTypeInt8:
			data, text, err = tag.ValidateInt8(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int8:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt8
		}

		switch tag.Type {
		case ir.ValueTypeInt8:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = int8(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateInt8(*val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int16:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt16
		}

		switch tag.Type {
		case ir.ValueTypeInt16:
			data, text, err = tag.ValidateInt16(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int16:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt16
		}

		switch tag.Type {
		case ir.ValueTypeInt16:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = int16(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateInt16(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt32
		}

		switch tag.Type {
		case ir.ValueTypeInt32:
			data, text, err = tag.ValidateInt32(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt32
		}

		switch tag.Type {
		case ir.ValueTypeInt32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = int32(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateInt32(*val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case int64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt64
		}

		switch tag.Type {
		case ir.ValueTypeInt64:
			data, text, err = tag.ValidateInt64(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *int64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeInt64
		}

		switch tag.Type {
		case ir.ValueTypeInt64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = int64(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateInt64(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint
		}

		switch tag.Type {
		case ir.ValueTypeUint:
			data, text, err = tag.ValidateUint(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint
		}

		switch tag.Type {
		case ir.ValueTypeUint:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = uint(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateUint(*val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint8:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint8
		}

		switch tag.Type {
		case ir.ValueTypeUint8:
			data, text, err = tag.ValidateUint8(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint8:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint8
		}

		switch tag.Type {
		case ir.ValueTypeUint8:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = uint8(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateUint8(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint16:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint16
		}

		switch tag.Type {
		case ir.ValueTypeUint16:
			data, text, err = tag.ValidateUint16(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint16:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint16
		}

		switch tag.Type {
		case ir.ValueTypeUint16:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = uint16(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateUint16(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint32
		}
		switch tag.Type {
		case ir.ValueTypeUint32:
			data, text, err = tag.ValidateUint32(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint32
		}

		switch tag.Type {
		case ir.ValueTypeUint32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = uint32(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateUint32(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case uint64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint64
		}

		switch tag.Type {
		case ir.ValueTypeUint64:
			data, text, err = tag.ValidateUint64(val)
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *uint64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeUint64
		}

		switch tag.Type {
		case ir.ValueTypeUint64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = uint64(0)
				text = "0"
			} else {
				data, text, err = tag.ValidateUint64(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case float32:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeFloat32
		}

		switch tag.Type {
		case ir.ValueTypeFloat32:
			if math.IsInf(float64(val), +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(float64(val), -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(float64(val)) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateFloat32(val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *float32:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeFloat32
		}

		switch tag.Type {
		case ir.ValueTypeFloat32:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = float32(0)
				text = "0.0"
			} else if math.IsInf(float64(*val), +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(float64(*val), -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(float64(*val)) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateFloat32(*val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case float64:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeFloat64
		}

		switch tag.Type {
		case ir.ValueTypeFloat64:
			if math.IsInf(val, +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(val, -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(val) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateFloat64(val)
			}
		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *float64:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeFloat64
		}

		switch tag.Type {
		case ir.ValueTypeFloat64:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = float64(0)
				text = "0.0"
			} else if math.IsInf(*val, +1) {
				return nil, fmt.Errorf("%s unsupported value: +Inf", tag.Type.String())
			} else if math.IsInf(*val, -1) {
				return nil, fmt.Errorf("%s unsupported value: -Inf", tag.Type.String())
			} else if math.IsNaN(*val) {
				return nil, fmt.Errorf("%s unsupported value: NaN", tag.Type.String())
			} else {
				data, text, err = tag.ValidateFloat64(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case string:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeString
		}

		switch tag.Type {
		case ir.ValueTypeString:
			data, text, err = tag.ValidateString(val)

		case ir.ValueTypeDecimal:
			data, text, err = tag.ValidateDecimal(val)

		case ir.ValueTypeEmail:
			data, text, err = tag.ValidateEmail(val)

		case ir.ValueTypeEnum:
			data, text, err = tag.ValidateEnum(val)

		case ir.ValueTypeUUID:
			data, text, err = tag.ValidateUUID(val)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *string:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeString
		}

		switch tag.Type {
		case ir.ValueTypeString:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = ""
				text = ""
			} else {
				data, text, err = tag.ValidateString(*val)
			}

		case ir.ValueTypeDecimal:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = ""
				text = ""
			} else {
				data, text, err = tag.ValidateDecimal(*val)
			}

		case ir.ValueTypeEmail:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = ""
				text = ""
			} else {
				data, text, err = tag.ValidateEmail(*val)
			}

		case ir.ValueTypeEnum:
			if tag.Enum == "" {
				err = fmt.Errorf("enum empty")
				return
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = -1
				text = ""
			} else {
				data, text, err = tag.ValidateEnum(*val)
			}

		case ir.ValueTypeUUID:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = [16]byte{}
				text = ""
			} else {
				data, text, err = tag.ValidateUUID(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case big.Int:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBigInt
		}

		switch tag.Type {
		case ir.ValueTypeBigInt:
			data, text, err = tag.ValidateBigInt(val)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *big.Int:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeBigInt
		}

		switch tag.Type {
		case ir.ValueTypeBigInt:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = big.Int{}
				text = "0"
			} else {
				data, text, err = tag.ValidateBigInt(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case url.URL:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeURL
		}

		switch tag.Type {
		case ir.ValueTypeURL:
			data, text, err = tag.ValidateURL(val)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *url.URL:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeURL
		}

		switch tag.Type {
		case ir.ValueTypeURL:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = url.URL{}
				text = ""
			} else {
				data, text, err = tag.ValidateURL(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case net.IP:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeIP
		}

		switch tag.Type {
		case ir.ValueTypeIP:
			data, text, err = tag.ValidateIP(val)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *net.IP:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeIP
		}

		switch tag.Type {
		case ir.ValueTypeIP:
			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = net.IP{}
				text = ""
			} else {
				data, text, err = tag.ValidateIP(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case time.Time:
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeDateTime
		}

		switch tag.Type {
		case ir.ValueTypeDateTime:
			data, text, err = tag.ValidateDateTime(val)

		case ir.ValueTypeDate:
			data, text, err = tag.ValidateDate(val)

		case ir.ValueTypeTime:
			data, text, err = tag.ValidateTime(val)

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	case *time.Time:
		tag.Nullable = true
		if tag.Type == ir.ValueTypeUnknown {
			tag.Type = ir.ValueTypeDateTime
		}

		switch tag.Type {
		case ir.ValueTypeDateTime:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = utils.DefaultTime
				text = utils.DefaultTime.In(location).Format(time.DateTime)
			} else {
				data, text, err = tag.ValidateTime(*val)
			}

		case ir.ValueTypeDate:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = utils.DefaultTime
				text = utils.DefaultTime.In(location).Format(time.DateOnly)
			} else {
				data, text, err = tag.ValidateTime(*val)
			}

		case ir.ValueTypeTime:
			location := time.UTC
			if tag.Location != nil {
				location = tag.Location
			}

			if val == nil {
				if !tag.Nullable {
					err = fmt.Errorf("value is nil and not nullable")
					return
				}
				tag.IsNull = true
				data = utils.DefaultTime
				text = utils.DefaultTime.In(location).Format(time.TimeOnly)
			} else {
				data, text, err = tag.ValidateTime(*val)
			}

		default:
			return nil, fmt.Errorf("%s unsupported type: %T", tag.Type.String(), val)
		}

	default:
		return anyToJSONC(v, tag, depth, path)
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

func anyToJSONC(obj any, tag *ir.Tag, depth int, path string) (ir.Node, error) {
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
			tag.IsNull = true
			typ = typ.Elem()
			elemVal, err := createExampleValue(typ)
			if err != nil {
				return nil, fmt.Errorf("create element %s: %w", typ, err)
			}
			val = reflect.ValueOf(elemVal)
		} else {
			val = val.Elem()
			typ = typ.Elem()
		}
		if val.Kind() == reflect.Interface {
			return nil, fmt.Errorf("unsupported type: interface{} (cannot reconstruct concrete value from abstract interface)")
		}
		if val.Kind() == reflect.Pointer {
			return nil, fmt.Errorf("unsupported type: multi-level pointer (%s) (only single-level pointer is allowed)", typ)
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

	tag.Type = ir.ValueTypeObject
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
			tag.Type = ir.ValueTypeObject

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
				fieldNode, err := valueToNode(fieldVal.Interface(), tagField, depth, p)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", p, err)
				}

				objNode.Fields = append(objNode.Fields, &ir.Field{
					Key:   fieldKey,
					Value: fieldNode,
				})
			}

			err = tag.ValidateStruct()
			if err != nil {
				err = fmt.Errorf("validate failed: %w", err)
				return nil, err
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
			valNode, err := valueToNode(val.MapIndex(key).Interface(), tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = valNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
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
			valNode, err := valueToNode(exampleVal, tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = valNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
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

		err = tag.ValidateMap()
		if err != nil {
			err = fmt.Errorf("validate failed: %w", err)
			return nil, err
		}
		return node, nil

	case reflect.Slice:
		tag.Type = ir.ValueTypeSlice

		node := &ir.Array{
			Tag:  tag,
			Path: path,
		}

		setTag := false
		for i := 0; i < val.Len(); i++ {
			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			p := fmt.Sprintf("%s[%d]", path, i)
			itemNode, err := valueToNode(val.Index(i).Interface(), tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if val.Len() == 0 {
			var exampleVal any
			exampleVal, err = createExampleValue(typ.Elem())
			if err != nil {
				return nil, fmt.Errorf("create example value for empty slice: %w", err)
			}

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Example = true

			p := fmt.Sprintf("%s[%d]", path, 0)
			itemNode, err := valueToNode(exampleVal, tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		err = tag.ValidateSlice(node.Items)
		if err != nil {
			err = fmt.Errorf("validate failed: %w", err)
			return nil, err
		}

		return node, nil

	case reflect.Array:
		tag.Type = ir.ValueTypeArray

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
			itemNode, err := valueToNode(val.Index(i).Interface(), tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if tag.Size == 0 {
			var exampleVal any
			exampleVal, err = createExampleValue(typ.Elem())
			if err != nil {
				return nil, fmt.Errorf("create example value for empty array: %w", err)
			}

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Example = true

			p := fmt.Sprintf("%s[%d]", path, 0)
			itemNode, err := valueToNode(exampleVal, tagItem, depth, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				node.Tag.ChildDesc = tagItem.Desc
				node.Tag.ChildType = tagItem.Type
				node.Tag.ChildRaw = tagItem.Raw
				node.Tag.ChildNullable = tagItem.Nullable
				node.Tag.ChildAllowEmpty = tagItem.AllowEmpty
				node.Tag.ChildUnique = tagItem.Unique
				node.Tag.ChildDefault = tagItem.Default
				node.Tag.ChildMin = tagItem.Min
				node.Tag.ChildMax = tagItem.Max
				node.Tag.ChildSize = tagItem.Size
				node.Tag.ChildEnum = tagItem.Enum
				node.Tag.ChildPattern = tagItem.Pattern
				node.Tag.ChildLocation = tagItem.Location
				node.Tag.ChildVersion = tagItem.Version
				node.Tag.ChildMime = tagItem.Mime
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		err = tag.ValidateArray(node.Items)
		if err != nil {
			err = fmt.Errorf("validate failed: %w", err)
			return nil, err
		}

		return node, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", val.Kind())
	}
}
