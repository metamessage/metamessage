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

const maxDepth = 256

// copyChildTag copies child-relevant attributes from child tag to parent tag's Child* fields.
func copyChildTag(parent *ir.Tag, child *ir.Tag) {
	parent.ChildDesc = child.Desc
	parent.ChildType = child.Type
	parent.ChildNullable = child.Nullable
	parent.ChildAllowEmpty = child.AllowEmpty
	parent.ChildUnique = child.Unique
	parent.ChildDefaultVal = child.DefaultVal
	parent.ChildMin = child.Min
	parent.ChildMax = child.Max
	parent.ChildSize = child.Size
	parent.ChildEnums = child.Enums
	parent.ChildPattern = child.Pattern
	parent.ChildLocation = child.Location
	parent.ChildVersion = child.Version
	parent.ChildMime = child.Mime
}

// refCtx 承载循环引用检测与 ref_id 分配状态，仅在单次 ValueToNode 调用内有效。
//
// visiting 即「当前递归路径栈」（D3）：容器类型进入 anyToJSONC / createExampleValue
// 时写入、返回时删除，因此同一类型出现在两个平行分支时各自独立编号，不会被误判为环。
type refCtx struct {
	visiting map[reflect.Type]*refFrame
	nextID   int // 每次编码调用从 1 开始自增
}

// refFrame 是路径栈上某一容器类型的帧。
// tag 指向该类型当前正在构建的节点的 tag（createExampleValue 内部登记的帧无 tag）。
// 环命中瞬间惰性分配 id 并直接回写 tag.RefID：tag 为指针，Bytes()/ToString()
// 在建树之后才调用，提前写安全；id 只在真被引用时分配，无空洞、无后扫描。
type refFrame struct {
	tag *ir.Tag
	id  int // 0 = 尚未被引用
}

// refMarker 是 createExampleValue 的环命中哨兵：命中帧带 tag（anyToJSONC 帧），
// id 已写入其 tag.RefID，调用方据 refTo 产出引用空壳节点。
type refMarker struct {
	refTo int
}

// cycleStop 是 createExampleValue 内部成环的哨兵（命中帧无 tag）：
// 停止展开该分支并置零，引用语义由 anyToJSONC 二次编码阶段重新检测产生。
type cycleStop struct{}

func newRefCtx() *refCtx {
	return &refCtx{
		visiting: make(map[reflect.Type]*refFrame),
		nextID:   1,
	}
}

// isContainerKind 报告该 kind 是否可能参与循环（需要登记路径栈）。
// Pointer 仅用于 createExampleValue 的自嵌套防护（如 type A *A）；
// anyToJSONC 在检测前已解引用指针，实际不会以 Pointer kind 登记。
func isContainerKind(k reflect.Kind) bool {
	switch k {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Pointer:
		return true
	}
	return false
}

// hitCycle 报告 typ 是否已在路径栈上；命中时对目标帧惰性分配 ref_id 并回写其 tag。
func (ctx *refCtx) hitCycle(typ reflect.Type) (id int, frame *refFrame, ok bool) {
	frame, onStack := ctx.visiting[typ]
	if !onStack {
		return 0, nil, false
	}
	if frame.id == 0 {
		frame.id = ctx.nextID
		ctx.nextID++
		if frame.tag != nil {
			frame.tag.RefID = frame.id
		}
	}
	return frame.id, frame, true
}

// enter 登记容器类型入栈；调用方须在返回前执行 leave 出栈。
func (ctx *refCtx) enter(typ reflect.Type, tag *ir.Tag) {
	ctx.visiting[typ] = &refFrame{tag: tag}
}

func (ctx *refCtx) leave(typ reflect.Type) {
	delete(ctx.visiting, typ)
}

func NilToNode(valueType ir.ValueType) (*ir.NodeScalar, error) {
	tag := ir.NewTag()
	tag.Type = valueType
	return &ir.NodeScalar{
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
	var tag *ir.Tag
	if tagStr != "" {
		if tag, err = ir.ParseMMTag(tagStr); err != nil {
			return nil, fmt.Errorf("parse tag failed: %w", err)
		}
	}

	if tag == nil {
		tag = ir.NewTag()
	}

	if v == nil {
		return &ir.NodeNull{
			Tag: tag,
		}, nil
	}

	example := tag.Example

	return valueToNode(newRefCtx(), v, tag, 0, "", example)
}

func valueToNode(ctx *refCtx, v any, tag *ir.Tag, depth int, path string, example bool) (node ir.Node, err error) {
	if tag == nil {
		tag = ir.NewTag()
	}

	var data any
	var text string
	text = ir.Null
	switch val := v.(type) {
	case nil:

		// var val any = nil
		return &ir.NodeNull{}, nil

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
			return anyToJSONC(ctx, v, tag, depth, path, example)
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
			return anyToJSONC(ctx, v, tag, depth, path, example)
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
				data, text, err = tag.ValidateDatetime(*val, example)
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
				data, text, err = tag.ValidateDate(*val, example)
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
		return anyToJSONC(ctx, v, tag, depth, path, example)
	}

	if err != nil {
		// err = fmt.Errorf("validate failed: %w", err)
		return
	}

	return &ir.NodeScalar{
		Data: data,
		Text: text,
		Tag:  tag,
		Path: path,
	}, nil
}

// createExampleValue 按元素类型构造示例值（空容器专用）。
//
// 循环防护：容器类型进入时登记路径栈（帧不带 tag）。
//   - 命中带 tag 的帧（anyToJSONC 登记）→ 返回 refMarker：
//     id 已写入被引用节点的 tag.RefID，调用方产出引用空壳；
//   - 命中无 tag 的帧（本函数内部登记）→ 返回 cycleStop：
//     分支置零停止展开，引用语义由 anyToJSONC 二次编码阶段重新检测产生。
func createExampleValue(ctx *refCtx, elemType reflect.Type) (any, error) {
	if isContainerKind(elemType.Kind()) {
		id, frame, hit := ctx.hitCycle(elemType)
		if hit {
			if frame.tag != nil {
				return &refMarker{refTo: id}, nil
			}
			return &cycleStop{}, nil
		}
		ctx.enter(elemType, nil)
		defer ctx.leave(elemType)
	}
	return createExampleValueInner(ctx, elemType)
}

// createExampleValueInner 构造示例值本体。所有哨兵只在分支边界处理，
// 不会嵌套进返回的真实值内部（struct 分支对哨兵字段置零跳过）。
func createExampleValueInner(ctx *refCtx, elemType reflect.Type) (any, error) {
	if elemType.Kind() == reflect.Pointer {
		base := elemType.Elem()
		baseVal, err := createExampleValue(ctx,base)
		if err != nil {
			return nil, err
		}
		switch baseVal.(type) {
		case *refMarker:
			// 指针指向环内类型：整体退化为引用
			return baseVal, nil
		case *cycleStop:
			// 内部成环：置为 nil 指针
			return reflect.Zero(elemType).Interface(), nil
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
		elemExample, err := createExampleValue(ctx,elemType.Elem())
		if err != nil {
			return nil, err
		}
		switch elemExample.(type) {
		case *refMarker:
			return elemExample, nil
		case *cycleStop:
			return reflect.MakeSlice(elemType, 0, 0).Interface(), nil
		}
		sliceVal := reflect.MakeSlice(elemType, 1, 1)
		sliceVal.Index(0).Set(reflect.ValueOf(elemExample))
		return sliceVal.Interface(), nil

	case reflect.Array:
		arrayVal := reflect.New(elemType).Elem()
		if elemType.Len() > 0 {
			elemExample, err := createExampleValue(ctx,elemType.Elem())
			if err != nil {
				return nil, err
			}
			switch elemExample.(type) {
			case *refMarker:
				return elemExample, nil
			case *cycleStop:
				return arrayVal.Interface(), nil
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

				fieldVal, err := createExampleValue(ctx,field.Type)
				if err != nil {
					return nil, fmt.Errorf("struct %s field %s: %w", elemType.Name(), field.Name, err)
				}

				// 环哨兵：该字段保持零值。refMarker 的引用语义、cycleStop 的
				// 展开语义均由 anyToJSONC 二次编码阶段重新检测并正确落地。
				switch fieldVal.(type) {
				case *refMarker, *cycleStop:
					continue
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
		keyExample, err := createExampleValue(ctx,elemType.Key())
		if err != nil {
			return nil, fmt.Errorf("create map key example: %w", err)
		}
		valExample, err := createExampleValue(ctx,elemType.Elem())
		if err != nil {
			return nil, fmt.Errorf("create map val example: %w", err)
		}
		switch valExample.(type) {
		case *refMarker:
			return valExample, nil
		case *cycleStop:
			return mapVal.Interface(), nil
		}
		switch keyExample.(type) {
		case *refMarker, *cycleStop:
			return mapVal.Interface(), nil
		}
		mapVal.SetMapIndex(reflect.ValueOf(keyExample), reflect.ValueOf(valExample))
		return mapVal.Interface(), nil

	default:
		return nil, fmt.Errorf("createExampleValue unsupported example type: %s", elemType.Kind())
	}
}

func anyToJSONC(ctx *refCtx, obj any, tag *ir.Tag, depth int, path string, example bool) (ir.Node, error) {
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
			elemVal, err := createExampleValue(ctx, typ)
			if err != nil {
				return nil, fmt.Errorf("%s create element %s err: %w", path, typ, err)
			}
			// 环内类型：产出引用空壳（D2：被引用类型的空形式，不展开字段）
			if r, isRef := elemVal.(*refMarker); isRef {
				tag.RefTo = r.refTo
				tag.Type = ir.ValueTypeObj
				return &ir.NodeObject{Tag: tag, Path: path}, nil
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

	// —— 循环引用检测（D3：按当前递归路径栈）——
	// 此处 val 已解引用指针，kind 只可能是 struct/map/slice/array。
	if isContainerKind(val.Kind()) {
		if id, _, hit := ctx.hitCycle(typ); hit {
			if example {
				// D2：引用节点 = 被引用类型的空形式，不展开字段
				tag.RefTo = id
				tag.Type = ir.ValueTypeObj
				return &ir.NodeObject{Tag: tag, Path: path}, nil
			}
			// 非 example（真实数据环）：不产壳、不重复登记，
			// 落入原逻辑由 maxDepth 兜底报错，保持既有行为不变。
		} else {
			ctx.enter(typ, tag)
			defer ctx.leave(typ)
		}
	}

	switch val.Kind() {
	case reflect.Struct:
		switch val.Type() {
		default:
			tag.Type = ir.ValueTypeObj

			objNode := &ir.NodeObject{
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
				fieldNode, err := valueToNode(ctx, fieldVal.Interface(), tagField, depth, p, example)
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

		node := &ir.NodeObject{
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
			valNode, err := valueToNode(ctx, val.MapIndex(key).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = valNode.GetTag()

			if !setTag {
				copyChildTag(node.Tag, tagItem)
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
			exampleVal, err = createExampleValue(ctx, elemType)
			if err != nil {
				return nil, fmt.Errorf("create example value for empty slice: %w", err)
			}

			keyStr := ""

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Name = keyStr

			tagItem.Example = true

			p := fmt.Sprintf("%s[%s]", path, keyStr)

			var valNode ir.Node
			if r, isRef := exampleVal.(*refMarker); isRef {
				// 环内类型：产出引用空壳（D2，ref_to 标在空壳自身 tag 上）
				tagItem.RefTo = r.refTo
				tagItem.Type = ir.ValueTypeObj
				valNode = &ir.NodeObject{Tag: tagItem, Path: p}
			} else {
				valNode, err = valueToNode(ctx, exampleVal, tagItem, depth, p, true)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", p, err)
				}
			}

			tagItem = valNode.GetTag()

			if !setTag {
				copyChildTag(node.Tag, tagItem)
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

		node := &ir.NodeArray{
			Tag:  tag,
			Path: path,
		}

		size := val.Len()
		setTag := false
		for i := range size {
			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			p := fmt.Sprintf("%s[%d]", path, i)
			itemNode, err := valueToNode(ctx, val.Index(i).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				copyChildTag(node.Tag, tagItem)
				setTag = true
			}

			node.Items = append(node.Items, itemNode)
		}

		if size == 0 {
			tag.IsEmpty = true
			var exampleVal any
			exampleVal, err = createExampleValue(ctx, typ.Elem())
			if err != nil {
				return nil, fmt.Errorf("create example value for empty slice: %w", err)
			}

			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			tagItem.Example = true

			p := fmt.Sprintf("%s[%d]", path, 0)

			var itemNode ir.Node
			if r, isRef := exampleVal.(*refMarker); isRef {
				// 环内类型：产出引用空壳（D2，ref_to 标在空壳自身 tag 上）
				tagItem.RefTo = r.refTo
				tagItem.Type = ir.ValueTypeObj
				itemNode = &ir.NodeObject{Tag: tagItem, Path: p}
			} else {
				itemNode, err = valueToNode(ctx, exampleVal, tagItem, depth, p, true)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", p, err)
				}
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				copyChildTag(node.Tag, tagItem)
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
		node := &ir.NodeArray{
			Tag:  tag,
			Path: path,
		}

		setTag := false
		for i := 0; i < tag.Size; i++ {
			tagItem := ir.NewTag()
			tagItem.Inherit(tag)

			p := fmt.Sprintf("%s[%d]", path, i)
			itemNode, err := valueToNode(ctx, val.Index(i).Interface(), tagItem, depth, p, example)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}

			tagItem = itemNode.GetTag()

			if !setTag {
				copyChildTag(node.Tag, tagItem)
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
