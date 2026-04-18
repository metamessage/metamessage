package jsonc

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"
	"unicode"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func Bind(node ast.Node, out any) error {
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Pointer || outVal.IsNil() {
		return fmt.Errorf("out must be a non-nil pointer")
	}

	outVal = outVal.Elem()

	switch n := node.(type) {
	case *ast.Object:
		if n.Tag.Type == ast.ValueTypeStruct {
			return convertStruct(n, outVal)
		} else {
			return convertMap(n, outVal)
		}

	case *ast.Array:
		if n.Tag.Type == ast.ValueTypeArray {
			return convertArray(n, outVal)
		} else {
			return convertSlice(n, outVal)
		}

	case *ast.Value:
		return convertValue(n, outVal)

	default:
		return fmt.Errorf("unsupported node type: %T", node)
	}
}

func convertStruct(obj *ast.Object, outVal reflect.Value) error {
	if obj.Tag.Nullable {
		if outVal.Kind() != reflect.Pointer {
			return fmt.Errorf("convertStruct requires pointer type, got %s", outVal.Kind())
		}
		if outVal.IsNil() {
			return fmt.Errorf("convertStruct: nullable pointer must point to struct, got nil")
		}
		outVal = outVal.Elem()
	}
	if outVal.Kind() != reflect.Struct {
		return fmt.Errorf("convertStruct requires struct type, got %s", outVal.Kind())
	}

	for _, field := range obj.Fields {
		fieldKey := field.Key
		var fieldVal reflect.Value
		found := false
		runes := []rune(fieldKey)
		runes[0] = unicode.ToUpper(runes[0])
		name := string(runes)
		for i := 0; i < outVal.NumField(); i++ {
			structField := outVal.Type().Field(i)
			if structField.Name == name {
				fieldVal = outVal.Field(i)
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("struct has no field '%s'", fieldKey)
		}

		if !fieldVal.CanSet() {
			return fmt.Errorf("struct field '%s' is unexported (cannot set)", fieldKey)
		}

		if err := Bind(field.Value, fieldVal.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to bind struct field %s: %w", field.Key, err)
		}
	}

	return nil
}

func convertMap(obj *ast.Object, outVal reflect.Value) error {
	if obj.Tag.Nullable {
		if outVal.Kind() != reflect.Pointer {
			return fmt.Errorf("convertMap requires pointer type, got %s", outVal.Kind())
		}
		if outVal.IsNil() {
			return fmt.Errorf("convertMap: nullable pointer must point to map, got nil")
		}
		outVal = outVal.Elem()
	}
	if outVal.Kind() != reflect.Map {
		return fmt.Errorf("convertMap requires map type, got %s", outVal.Kind())
	}

	if outVal.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("convertMap requires map with string key, got %s", outVal.Type().Key().Kind())
	}

	if outVal.Len() > 0 {
		return fmt.Errorf("target map is not empty (length: %d), clear it before binding", outVal.Len())
	}

	mapValType := outVal.Type().Elem()
	for _, field := range obj.Fields {
		key := reflect.ValueOf(field.Key)
		val := reflect.New(mapValType).Elem()
		if err := Bind(field.Value, val.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert field %s: %w", field.Key, err)
		}

		outVal.SetMapIndex(key, val)
	}
	return nil
}

func convertArray(arr *ast.Array, outVal reflect.Value) error {
	if arr.Tag.Nullable {
		if outVal.Kind() != reflect.Pointer {
			return fmt.Errorf("convertArray requires pointer type, got %s", outVal.Kind())
		}
		if outVal.IsNil() {
			return fmt.Errorf("convertArray: nullable pointer must point to array, got nil")
		}
		outVal = outVal.Elem()
	}
	if outVal.Kind() != reflect.Array {
		return fmt.Errorf("it can only be converted to array, got %s", outVal.Kind())
	}

	arrayLen := outVal.Len()
	size := arr.Tag.Size
	if size != arrayLen {
		return fmt.Errorf("array length mismatch: target array length %d, got %d items", arrayLen, size)
	}

	for i, item := range arr.Items {
		elem := outVal.Index(i)
		if err := Bind(item, elem.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert array item %d: %w", i, err)
		}
	}

	return nil
}

func convertSlice(arr *ast.Array, outVal reflect.Value) error {
	if arr.Tag.Nullable {
		if outVal.Kind() != reflect.Pointer {
			return fmt.Errorf("convertArray requires pointer type, got %s", outVal.Kind())
		}
		if outVal.IsNil() {
			return fmt.Errorf("convertArray: nullable pointer must point to slice, got nil")
		}
		outVal = outVal.Elem()
	}
	if outVal.Kind() != reflect.Slice {
		return fmt.Errorf("it can only be converted to slice, got %s", outVal.Kind())
	}

	size := len(arr.Items)
	slice := reflect.MakeSlice(outVal.Type(), size, size)

	for i, item := range arr.Items {
		elem := slice.Index(i)
		if err := Bind(item, elem.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert array item %d: %w", i, err)
		}
	}

	outVal.Set(slice)
	return nil
}

func convertValue(val *ast.Value, outVal reflect.Value) error {
	if val.Tag.Nullable {
		if outVal.Kind() != reflect.Pointer {
			return fmt.Errorf("convertValue requires pointer type, got %s", outVal.Kind())
		}
		if outVal.IsNil() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}
		outVal = outVal.Elem()
	}

	data := val.Data
	text := val.Text

	tag := val.GetTag()
	if tag == nil {
		return fmt.Errorf("")
	}

	switch tag.Type {
	case ast.ValueTypeTime, ast.ValueTypeDate, ast.ValueTypeDateTime:
		targetType := reflect.TypeFor[time.Time]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be time.Time, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target time.Time value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case time.Time:
			*(*time.Time)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for time conversion: %T (expected time.Time)", d)
		}

	case ast.ValueTypeBigInt:
		targetType := reflect.TypeFor[big.Int]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be big.Int, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *big.Int pointer is nil")
		}

		switch d := data.(type) {
		case big.Int:
			*(*big.Int)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for big.Int conversion: %T", d)
		}

	case ast.ValueTypeUUID:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *string pointer is nil")
		}

		switch d := data.(type) {
		case string:
			*(*string)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uuid conversion: %T", d)
		}

	case ast.ValueTypeDecimal:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *string pointer is nil")
		}

		switch d := data.(type) {
		case string:
			*(*string)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for decimal conversion: %T", d)
		}

	case ast.ValueTypeEmail:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *string pointer is nil")
		}

		switch d := data.(type) {
		case string:
			*(*string)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for email conversion: %T", d)
		}

	case ast.ValueTypeIP:
		targetType := reflect.TypeFor[net.IP]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be net.IP, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *net.IP pointer is nil")
		}

		switch d := data.(type) {
		case net.IP:
			*(*net.IP)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for ip conversion: %T", d)
		}

	case ast.ValueTypeURL:
		targetType := reflect.TypeFor[url.URL]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be url.URL, got %s", outVal.Type())
		}
		if outVal.IsNil() {
			return fmt.Errorf("target *url.URL pointer is nil")
		}

		switch d := data.(type) {
		case url.URL:
			*(*url.URL)(outVal.UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for url conversion: %T", d)
		}

	case ast.ValueTypeInt:
		if outVal.Type() != reflect.TypeFor[int]() {
			return fmt.Errorf("target type must be int, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target int value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case int:
			*(*int)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for int conversion: %T", d)
		}

	case ast.ValueTypeInt8:
		if outVal.Type() != reflect.TypeFor[int8]() {
			return fmt.Errorf("target type must be int8, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target int8 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case int8:
			*(*int8)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for int8 conversion: %T", d)
		}

	case ast.ValueTypeInt16:
		if outVal.Type() != reflect.TypeFor[int16]() {
			return fmt.Errorf("target type must be int16, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target int16 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case int16:
			*(*int16)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for int16 conversion: %T (expected int16)", d)
		}

	case ast.ValueTypeInt32:
		if outVal.Type() != reflect.TypeFor[int32]() {
			return fmt.Errorf("target type must be int32, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target int32 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case int32:
			*(*int32)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for int32 conversion: %T (expected int32)", d)
		}

	case ast.ValueTypeInt64:
		if outVal.Type() != reflect.TypeFor[int64]() {
			return fmt.Errorf("target type must be int64, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target int64 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case int64:
			*(*int64)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for int64 conversion: %T (expected int64)", d)
		}

	case ast.ValueTypeUint:
		if outVal.Type() != reflect.TypeFor[uint]() {
			return fmt.Errorf("target type must be uint, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target uint value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case uint:
			*(*uint)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uint conversion: %T (expected uint)", d)
		}

	case ast.ValueTypeUint8:
		if outVal.Type() != reflect.TypeFor[uint8]() {
			return fmt.Errorf("target type must be uint8, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target uint8 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case uint8:
			*(*uint8)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uint8 conversion: %T (expected uint8)", d)
		}

	case ast.ValueTypeUint16:
		if outVal.Type() != reflect.TypeFor[uint16]() {
			return fmt.Errorf("target type must be uint16, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target uint16 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case uint16:
			*(*uint16)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uint16 conversion: %T (expected uint16)", d)
		}

	case ast.ValueTypeUint32:
		if outVal.Type() != reflect.TypeFor[uint32]() {
			return fmt.Errorf("target type must be uint32, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target uint32 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case uint32:
			*(*uint32)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uint32 conversion: %T (expected uint32)", d)
		}

	case ast.ValueTypeUint64:
		if outVal.Type() != reflect.TypeFor[uint64]() {
			return fmt.Errorf("target type must be uint64, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target uint64 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case uint64:
			*(*uint64)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for uint64 conversion: %T (expected uint64)", d)
		}

	case ast.ValueTypeFloat32:
		if outVal.Type() != reflect.TypeFor[float32]() {
			return fmt.Errorf("target type must be float32, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target float32 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case float32:
			*(*float32)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for float32 conversion: %T (expected float32)", d)
		}

	case ast.ValueTypeFloat64:
		if outVal.Type() != reflect.TypeFor[float64]() {
			return fmt.Errorf("target type must be float64, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target float64 value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case float64:
			*(*float64)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for float64 conversion: %T (expected float64)", d)
		}

	case ast.ValueTypeString:
		if outVal.Type() != reflect.TypeFor[string]() {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target string value is unaddressable (bug)")
		}
		*(*string)(outVal.Addr().UnsafePointer()) = text

	case ast.ValueTypeBool:
		if outVal.Type() != reflect.TypeFor[bool]() {
			return fmt.Errorf("target type must be bool, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target bool value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case bool:
			*(*bool)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for bool conversion: %T (expected bool)", d)
		}

	case ast.ValueTypeBytes:
		if outVal.Type() != reflect.TypeFor[[]byte]() {
			return fmt.Errorf("target type must be []byte, got %s", outVal.Type())
		}
		if !outVal.CanAddr() {
			return fmt.Errorf("unexpected: target []byte value is unaddressable (bug)")
		}
		switch d := data.(type) {
		case []byte:
			*(*[]byte)(outVal.Addr().UnsafePointer()) = d
		default:
			return fmt.Errorf("unsupported type for bytes conversion: %T (expected []byte)", d)
		}

	case ast.ValueTypeMap:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ast.ValueTypeStruct:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ast.ValueTypeArray:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ast.ValueTypeSlice:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ast.ValueTypeUnknown:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	default:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	}

	return nil
}
