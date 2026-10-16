package core

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/utils"
)

func newBoolTagReader(get func(*ir.Tag) bool) func(ir.Node) bool {
	return func(n ir.Node) bool {
		if n == nil {
			return false
		}
		t := n.GetTag()
		return t != nil && get(t)
	}
}

// isExampleNode 报告节点是否为空容器注入的 example 哨兵项。
// example 项只是承载容器元素类型信息的占位符，不进入绑定结果。
var isExampleNode = newBoolTagReader(func(t *ir.Tag) bool { return t.Example })

// isNullNode 报告节点是否带 is_null（编码侧真实值为 nil）。
var isNullNode = newBoolTagReader(func(t *ir.Tag) bool { return t.IsNull })

// bindNil 把目标置为零值：指针/切片/映射置 nil，其余类型置零值。
func bindNil(outVal reflect.Value) error {
	if !outVal.CanSet() {
		return fmt.Errorf("target %s is not settable", outVal.Type())
	}
	outVal.Set(reflect.Zero(outVal.Type()))
	return nil
}

func Bind(node ir.Node, out any) error {
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Pointer {
		return fmt.Errorf("out must be a pointer")
	}

	if outVal.IsNil() {
		if node.GetType() == ir.NodeTypeNull {
			return nil
		} else {
			return fmt.Errorf("out pointer must not be nil for non-null node")
		}
	}

	outVal = outVal.Elem()

	// When the target is an interface{}, create concrete values dynamically
	if outVal.Kind() == reflect.Interface {
		if isNullNode(node) {
			outVal.Set(reflect.Zero(outVal.Type()))
			return nil
		}

		switch n := node.(type) {
		case *ir.NodeObject:
			m := make(map[string]any)
			for _, field := range n.Fields {
				if isExampleNode(field.Value) {
					continue // 空容器占位项不绑定
				}
				var v any
				if err := Bind(field.Value, &v); err != nil {
					return fmt.Errorf("failed to bind field %s: %w", field.Key, err)
				}
				m[field.Key] = v
			}
			outVal.Set(reflect.ValueOf(m))
			return nil

		case *ir.NodeArray:
			s := make([]any, 0, len(n.Items))
			for i, item := range n.Items {
				if isExampleNode(item) {
					continue // 空容器占位项不绑定
				}
				var v any
				if err := Bind(item, &v); err != nil {
					return fmt.Errorf("failed to bind array item %d: %w", i, err)
				}
				s = append(s, v)
			}
			outVal.Set(reflect.ValueOf(s))
			return nil

		case *ir.NodeScalar:
			if n.Data != nil {
				outVal.Set(reflect.ValueOf(n.Data))
			}
			return nil

		case *ir.NodeNull:
			return nil

		default:
			return fmt.Errorf("unsupported node type: %T", node)
		}
	}

	switch n := node.(type) {
	case *ir.NodeObject:
		if n.Tag.Type == ir.ValueTypeObj {
			return convertObj(n, outVal)
		} else {
			return convertMap(n, outVal)
		}

	case *ir.NodeArray:
		if n.Tag.Type == ir.ValueTypeArr {
			return convertArr(n, outVal)
		} else {
			return convertVec(n, outVal)
		}

	case *ir.NodeScalar:
		return convertScalar(n, outVal)

	case *ir.NodeNull:
		return nil

	default:
		return fmt.Errorf("unsupported node type: %T", node)
	}
}

func getFieldMap(t reflect.Type, val reflect.Value) map[string]reflect.Value {
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldMap[utils.CamelToSnake(field.Name)] = val.Field(i)
	}

	return fieldMap
}

func derefNullable(path string, tag *ir.Tag, outVal reflect.Value) (reflect.Value, bool, error) {
	if tag == nil || !tag.Nullable {
		return outVal, false, nil
	}
	if outVal.Kind() != reflect.Pointer {
		return outVal, false, fmt.Errorf("%s requires pointer type, got %s", path, outVal.Kind())
	}
	if tag.IsNull {
		return outVal, true, nil
	}
	if outVal.IsNil() {
		outVal.Set(reflect.New(outVal.Type().Elem()))
	}
	return outVal.Elem(), false, nil
}

func convertObj(obj *ir.NodeObject, outVal reflect.Value) error {
	target, isNull, err := derefNullable(obj.Path, obj.Tag, outVal)
	if err != nil {
		return err
	}
	if isNull {
		return bindNil(outVal) // is_null：绑定成 nil 指针
	}
	outVal = target

	if outVal.Kind() != reflect.Struct {
		return fmt.Errorf("%s requires struct type, got %s", obj.Path, outVal.Kind())
	}

	fieldMap := getFieldMap(outVal.Type(), outVal)

	for _, field := range obj.Fields {
		fieldKey := field.Key
		fieldVal, found := fieldMap[fieldKey]

		if !found {
			return fmt.Errorf("%s has no field '%s'", obj.Path, fieldKey)
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

func convertMap(obj *ir.NodeObject, outVal reflect.Value) error {
	target, isNull, err := derefNullable(obj.Path, obj.Tag, outVal)
	if err != nil {
		return err
	}
	if isNull {
		return bindNil(outVal) // is_null：绑定成 nil 映射
	}
	outVal = target

	if outVal.Kind() != reflect.Map {
		return fmt.Errorf("%s requires map type, got %s", obj.Path, outVal.Kind())
	}

	if outVal.IsNil() {
		outVal.Set(reflect.MakeMap(outVal.Type()))
	}

	if outVal.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("convertMap requires map with string key, got %s", outVal.Type().Key().Kind())
	}

	if outVal.Len() > 0 {
		return fmt.Errorf("target map is not empty (length: %d), clear it before binding", outVal.Len())
	}

	mapValType := outVal.Type().Elem()
	for _, field := range obj.Fields {
		if isExampleNode(field.Value) {
			continue // 空 map 的 example 占位项：跳过后得到非 nil 空映射
		}

		key := reflect.ValueOf(field.Key)
		val := reflect.New(mapValType).Elem()
		if err := Bind(field.Value, val.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert field %s: %w", field.Key, err)
		}

		outVal.SetMapIndex(key, val)
	}
	return nil
}

func convertArr(arr *ir.NodeArray, outVal reflect.Value) error {
	target, isNull, err := derefNullable(arr.Path, arr.Tag, outVal)
	if err != nil {
		return err
	}
	if isNull {
		return bindNil(outVal) // is_null：绑定成 nil 指针
	}
	outVal = target

	if outVal.Kind() != reflect.Array {
		return fmt.Errorf("%s requires array type, got %s", arr.Path, outVal.Kind())
	}

	arrayLen := outVal.Len()
	size := arr.Tag.Size
	if size != arrayLen {
		return fmt.Errorf("array length mismatch: target array length %d, got %d items", arrayLen, size)
	}

	idx := 0
	for _, item := range arr.Items {
		if isExampleNode(item) {
			continue // 固定数组长度由 tag.size 决定，example 占位项不占下标
		}
		if idx >= outVal.Len() {
			return fmt.Errorf("%s: real items (%d) exceed target array length %d", arr.Path, idx+1, outVal.Len())
		}
		elem := outVal.Index(idx)
		if err := Bind(item, elem.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert array item %d: %w", idx, err)
		}
		idx++
	}

	return nil
}

func convertVec(arr *ir.NodeArray, outVal reflect.Value) error {
	target, isNull, err := derefNullable(arr.Path, arr.Tag, outVal)
	if err != nil {
		return err
	}
	if isNull {
		return bindNil(outVal) // is_null：绑定成 nil 切片
	}
	outVal = target

	if outVal.Kind() != reflect.Slice {
		return fmt.Errorf("%s requires slice type, got %s", arr.Path, outVal.Kind())
	}

	// example 哨兵项只为保留元素类型，不参与绑定；
	// 全部为 example（编码侧为空容器）时产出长度 0 的非 nil 切片。
	items := make([]ir.Node, 0, len(arr.Items))
	for _, item := range arr.Items {
		if isExampleNode(item) {
			continue
		}
		items = append(items, item)
	}

	size := len(items)
	slice := reflect.MakeSlice(outVal.Type(), size, size)

	for i, item := range items {
		elem := slice.Index(i)
		if err := Bind(item, elem.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to convert slice item %d: %w", i, err)
		}
	}

	outVal.Set(slice)
	return nil
}

func convertScalar(val *ir.NodeScalar, outVal reflect.Value) error {
	target, isNull, err := derefNullable(val.Path, val.Tag, outVal)
	if err != nil {
		return err
	}
	if isNull {
		return bindNil(outVal) // is_null：绑定成 nil 指针
	}
	outVal = target

	data := val.Data
	text := val.Text

	tag := val.GetTag()
	if tag == nil {
		return fmt.Errorf("%s: tag is nil", val.Path)
	}

	switch tag.Type {
	case ir.ValueTypeTime, ir.ValueTypeDate, ir.ValueTypeDatetime:
		targetType := reflect.TypeFor[time.Time]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be time.Time, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case time.Time:
			outVal.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("unsupported type for time conversion: %T (expected time.Time)", d)
		}

	case ir.ValueTypeBigint:
		targetType := reflect.TypeFor[big.Int]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be big.Int, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case big.Int:
			outVal.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("unsupported type for big.Int conversion: %T", d)
		}

	case ir.ValueTypeUuid:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case string:
			outVal.SetString(d)
		default:
			return fmt.Errorf("unsupported type for uuid conversion: %T", d)
		}

	case ir.ValueTypeDecimal:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case string:
			outVal.SetString(d)
		default:
			return fmt.Errorf("unsupported type for decimal conversion: %T", d)
		}

	case ir.ValueTypeEmail:
		targetType := reflect.TypeFor[string]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case string:
			outVal.SetString(d)
		default:
			return fmt.Errorf("unsupported type for email conversion: %T", d)
		}

	case ir.ValueTypeIp:
		targetType := reflect.TypeFor[net.IP]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be net.IP, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case net.IP:
			outVal.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("unsupported type for ip conversion: %T", d)
		}

	case ir.ValueTypeUrl:
		targetType := reflect.TypeFor[url.URL]()
		if outVal.Type() != targetType {
			return fmt.Errorf("target type must be url.URL, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case url.URL:
			outVal.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("unsupported type for url conversion: %T", d)
		}

	case ir.ValueTypeEnum:
		if outVal.Type() != reflect.TypeFor[string]() {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case string:
			outVal.SetString(text)
		default:
			return fmt.Errorf("unsupported type for enum conversion: %T", d)
		}

	case ir.ValueTypeDoc:
		return fmt.Errorf("unsupported type: %s", tag.Type)

	case ir.ValueTypeMedia:
		return fmt.Errorf("unsupported type: %s", tag.Type)

	case ir.ValueTypeI:
		if outVal.Type() != reflect.TypeFor[int]() {
			return fmt.Errorf("target type must be int, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case int:
			outVal.SetInt(int64(d))
		default:
			return fmt.Errorf("unsupported type for int conversion: %T", d)
		}

	case ir.ValueTypeI8:
		if outVal.Type() != reflect.TypeFor[int8]() {
			return fmt.Errorf("target type must be int8, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case int8:
			outVal.SetInt(int64(d))
		default:
			return fmt.Errorf("unsupported type for int8 conversion: %T", d)
		}

	case ir.ValueTypeI16:
		if outVal.Type() != reflect.TypeFor[int16]() {
			return fmt.Errorf("target type must be int16, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case int16:
			outVal.SetInt(int64(d))
		default:
			return fmt.Errorf("unsupported type for int16 conversion: %T (expected int16)", d)
		}

	case ir.ValueTypeI32:
		if outVal.Type() != reflect.TypeFor[int32]() {
			return fmt.Errorf("target type must be int32, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case int32:
			outVal.SetInt(int64(d))
		default:
			return fmt.Errorf("unsupported type for int32 conversion: %T (expected int32)", d)
		}

	case ir.ValueTypeI64:
		if outVal.Type() != reflect.TypeFor[int64]() {
			return fmt.Errorf("target type must be int64, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case int64:
			outVal.SetInt(d)
		default:
			return fmt.Errorf("unsupported type for int64 conversion: %T (expected int64)", d)
		}

	case ir.ValueTypeU:
		if outVal.Type() != reflect.TypeFor[uint]() {
			return fmt.Errorf("target type must be uint, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case uint:
			outVal.SetUint(uint64(d))
		default:
			return fmt.Errorf("unsupported type for uint conversion: %T (expected uint)", d)
		}

	case ir.ValueTypeU8:
		if outVal.Type() != reflect.TypeFor[uint8]() {
			return fmt.Errorf("target type must be uint8, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case uint8:
			outVal.SetUint(uint64(d))
		default:
			return fmt.Errorf("unsupported type for uint8 conversion: %T (expected uint8)", d)
		}

	case ir.ValueTypeU16:
		if outVal.Type() != reflect.TypeFor[uint16]() {
			return fmt.Errorf("target type must be uint16, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case uint16:
			outVal.SetUint(uint64(d))
		default:
			return fmt.Errorf("unsupported type for uint16 conversion: %T (expected uint16)", d)
		}

	case ir.ValueTypeU32:
		if outVal.Type() != reflect.TypeFor[uint32]() {
			return fmt.Errorf("target type must be uint32, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case uint32:
			outVal.SetUint(uint64(d))
		default:
			return fmt.Errorf("unsupported type for uint32 conversion: %T (expected uint32)", d)
		}

	case ir.ValueTypeU64:
		if outVal.Type() != reflect.TypeFor[uint64]() {
			return fmt.Errorf("target type must be uint64, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case uint64:
			outVal.SetUint(d)
		default:
			return fmt.Errorf("unsupported type for uint64 conversion: %T (expected uint64)", d)
		}

	case ir.ValueTypeF32:
		if outVal.Type() != reflect.TypeFor[float32]() {
			return fmt.Errorf("target type must be float32, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case float32:
			outVal.SetFloat(float64(d))
		default:
			return fmt.Errorf("unsupported type for float32 conversion: %T (expected float32)", d)
		}

	case ir.ValueTypeF64:
		if outVal.Type() != reflect.TypeFor[float64]() {
			return fmt.Errorf("target type must be float64, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case float64:
			outVal.SetFloat(d)
		default:
			return fmt.Errorf("unsupported type for float64 conversion: %T (expected float64)", d)
		}

	case ir.ValueTypeStr:
		if outVal.Type() != reflect.TypeFor[string]() {
			return fmt.Errorf("target type must be string, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case string:
			outVal.SetString(text)
		default:
			return fmt.Errorf("unsupported type for string conversion: %T", d)
		}

	case ir.ValueTypeBool:
		if outVal.Type() != reflect.TypeFor[bool]() {
			return fmt.Errorf("target type must be bool, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case bool:
			outVal.SetBool(d)
		default:
			return fmt.Errorf("unsupported type for bool conversion: %T (expected bool)", d)
		}

	case ir.ValueTypeBytes:
		if outVal.Type() != reflect.TypeFor[[]byte]() {
			return fmt.Errorf("target type must be []byte, got %s", outVal.Type())
		}

		switch d := data.(type) {
		case []byte:
			outVal.SetBytes(d)
		default:
			return fmt.Errorf("unsupported type for bytes conversion: %T (expected []byte)", d)
		}

	case ir.ValueTypeMap:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ir.ValueTypeObj:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ir.ValueTypeArr:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ir.ValueTypeVec:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	case ir.ValueTypeUnknown:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	default:
		return fmt.Errorf("unsupported type: %s", tag.Type)
	}

	return nil
}
