package core

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net"
	"time"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/utils"
)

type Encoder interface {
	Reset(io.Writer)
	EncodeStream(in any) (n int, err error)
	Encode(node ir.Node) (out []byte, err error)
}

type encoder struct {
	w      io.Writer
	buf    []byte
	offset uint32
	maxCap uint32
}

const defaultBufSize = 1024       // 1KB
const maxCap = 1024 * 1024 * 1024 // 1GB

func NewEncoder(w io.Writer) *encoder {
	e := &encoder{
		buf:    make([]byte, defaultBufSize),
		offset: 0,
		maxCap: maxCap,
	}
	e.Reset(w)
	return e
}

func (e *encoder) Reset(w io.Writer) {
	if w == nil {
		return
	}

	e.w = w
}

func (e *encoder) encodeNodeObject(obj *ir.Object) (n uint32, err error) {
	var bufKey bytes.Buffer
	var buf bytes.Buffer
	tag := obj.GetTag()

	for _, field := range obj.Fields {
		switch val := field.Value.(type) {
		case *ir.Object:
			if n, err = e.encodeNodeObject(val); err != nil {
				return
			}

		case *ir.Array:
			if n, err = e.encodeNodeArray(val); err != nil {
				return
			}

		case *ir.Value:
			if n, err = e.encodeNodeValue(val); err != nil {
				err = fmt.Errorf("%s: %w", val.GetPath(), err)
				return
			}

		default:
			return 0, fmt.Errorf("unsupported type %T", val)
		}

		encodedSub := e.buf[e.offset-n : e.offset]
		_, err = buf.Write(encodedSub)
		if err != nil {
			return
		}

		// maybe special encode. it must be a string
		var ns uint32
		if ns, err = e.encodeString(field.Key); err != nil {
			return
		}
		encodedKey := e.buf[e.offset-ns : e.offset]
		_, err = bufKey.Write(encodedKey)
		if err != nil {
			return
		}
	}

	var nk uint32
	if nk, err = e.encodeArray(bufKey.Bytes()); err != nil {
		return
	}

	encodedKeyArray := e.buf[e.offset-nk : e.offset]
	bufAll := make([]byte, len(encodedKeyArray)+buf.Len())
	copy(bufAll, encodedKeyArray)
	copy(bufAll[len(encodedKeyArray):], buf.Bytes())

	if n, err = e.encodeObject(bufAll); err != nil {
		return
	}

	var n1 uint32
	n1, err = e.encodeComment(e.buf[e.offset-n:e.offset], tag)
	if err != nil {
		return
	}

	if n1 == 0 {
		return
	}

	n = n1

	return
}

func (e *encoder) encodeComment(payload []byte, tag *ir.Tag) (n uint32, err error) {
	var ns uint32
	if ns, err = e.encodeT(tag.Bytes()); err != nil {
		return
	}

	if ns == 0 {
		return
	}

	return e.encodeTag(payload, e.buf[e.offset-ns:e.offset])
}

func (e *encoder) encodeNodeArray(arr *ir.Array) (n uint32, err error) {
	var buf bytes.Buffer

	tag := arr.GetTag()

	for _, item := range arr.Items {
		switch val := item.(type) {
		case *ir.Object:
			if n, err = e.encodeNodeObject(val); err != nil {
				return
			}
		case *ir.Array:
			if n, err = e.encodeNodeArray(val); err != nil {
				return
			}
		case *ir.Value:
			if n, err = e.encodeNodeValue(val); err != nil {
				return
			}
		default:
			return 0, fmt.Errorf("unsupported type %T", val)
		}

		encodedSub := e.buf[e.offset-n : e.offset]
		_, err = buf.Write(encodedSub)
		if err != nil {
			return
		}
	}

	if n, err = e.encodeArray(buf.Bytes()); err != nil {
		return
	}

	var n1 uint32
	n1, err = e.encodeComment(e.buf[e.offset-n:e.offset], tag)
	if err != nil {
		return
	}

	if n1 == 0 {
		return
	}

	n = n1

	return
}

func (e *encoder) encodeNodeValue(val *ir.Value) (n uint32, err error) {
	// defer func() {
	// 	fmt.Println("encodeNodeValue:", val.Data, val.Text, val.Tag.Type)
	// }()

	tag := val.GetTag()

	switch val.Tag.Type {
	case ir.ValueTypeDateTime:
		if tag.IsNull {
		} else {
			n, err = e.encodeDateTime(val.Data.(time.Time))
		}

	case ir.ValueTypeDate:
		if tag.IsNull {
		} else {
			n, err = e.encodeDate(val.Data.(time.Time))
		}

	case ir.ValueTypeTime:
		if tag.IsNull {
		} else {
			n, err = e.encodeTime(val.Data.(time.Time))
		}

	case ir.ValueTypeInt:
		if tag.IsNull {
			n, err = e.encodeSimple(SimpleNullInt)
		} else {
			n, err = e.encodeInt64(int64(val.Data.(int)))
		}

	case ir.ValueTypeInt8:
		if tag.IsNull {
		} else {
			n, err = e.encodeInt64(int64(val.Data.(int8)))
		}

	case ir.ValueTypeInt16:
		if tag.IsNull {
		} else {
			n, err = e.encodeInt64(int64(val.Data.(int16)))
		}

	case ir.ValueTypeInt32:
		if tag.IsNull {
		} else {
			n, err = e.encodeInt64(int64(val.Data.(int32)))
		}

	case ir.ValueTypeInt64:
		if tag.IsNull {
		} else {
			n, err = e.encodeInt64(val.Data.(int64))
		}

	case ir.ValueTypeUint:
		if tag.IsNull {
		} else {
			n, err = e.encodeUint64(uint64(val.Data.(uint)))
		}

	case ir.ValueTypeUint8:
		if tag.IsNull {
		} else {
			n, err = e.encodeUint64(uint64(val.Data.(uint8)))
		}

	case ir.ValueTypeUint16:
		if tag.IsNull {
		} else {
			n, err = e.encodeUint64(uint64(val.Data.(uint16)))
		}

	case ir.ValueTypeUint32:
		if tag.IsNull {
		} else {
			n, err = e.encodeUint64(uint64(val.Data.(uint32)))
		}

	case ir.ValueTypeUint64:
		if tag.IsNull {
		} else {
			n, err = e.encodeUint64(val.Data.(uint64))
		}

	case ir.ValueTypeFloat32:
		if tag.IsNull {
		} else {
			n, err = e.encodeFloat(val.Text)
		}

	case ir.ValueTypeFloat64:
		if tag.IsNull {
			n, err = e.encodeSimple(SimpleNullFloat)
		} else {
			n, err = e.encodeFloat(val.Text)
		}

	case ir.ValueTypeString:
		if tag.IsNull {
			n, err = e.encodeSimple(SimpleNullString)
		} else {
			n, err = e.encodeString(val.Text)
		}

	case ir.ValueTypeEmail:
		if tag.IsNull {
		} else {
			n, err = e.encodeString(val.Text)
		}

	case ir.ValueTypeUUID:
		if tag.IsNull {
		} else {
			arr, ok := val.Data.([16]byte)
			if !ok {
				return 0, fmt.Errorf("invalid UUID string: %q", val.Text)
			}
			n, err = e.encodeBytes(arr[:])
		}

	case ir.ValueTypeDecimal:
		if tag.IsNull {
		} else {
			n, err = e.encodeFloat(val.Text)
		}

	case ir.ValueTypeURL:
		if tag.IsNull {
		} else {
			n, err = e.encodeString(val.Text)
		}

	case ir.ValueTypeIP:
		if tag.IsNull {
		} else {
			ip := val.Data.(net.IP)
			switch tag.Version {
			case 0:
				n, err = e.encodeString(val.Text)
			case 4:
				n, err = e.encodeBytes(ip)
			case 6:
				if len(val.Text) < 16 {
					n, err = e.encodeString(val.Text)
				} else {
					n, err = e.encodeBytes(ip)
				}
			default:
				return 0, fmt.Errorf("unsupported IP version: %d", tag.Version)
			}
		}

	case ir.ValueTypeBytes:
		if tag.IsNull {
			n, err = e.encodeSimple(SimpleNullBytes)
		} else {
			n, err = e.encodeBytes(val.Data.([]byte))
		}

	case ir.ValueTypeBigInt:
		if tag.IsNull {
		} else {
			n, err = e.encodeBigInt(val.Text)
		}

	case ir.ValueTypeBool:
		if tag.IsNull {
			n, err = e.encodeSimple(SimpleNullBool)
		} else {
			n, err = e.encodeBool(val.Data.(bool))
		}

	case ir.ValueTypeEnum:
		if tag.IsNull {
		} else {
			n, err = e.encodeInt64(int64(val.Data.(int)))
		}

	default:
		return 0, fmt.Errorf("type error: unsupported type: %v, value: %v", val.Tag.Type, val.Data)
	}

	if err != nil {
		return
	}

	var n1 uint32
	n1, err = e.encodeComment(e.buf[e.offset-n:e.offset], tag)
	if err != nil {
		return
	}

	if n1 == 0 {
		return
	}

	n = n1

	return
}

func (e *encoder) Encode(node ir.Node) (out []byte, err error) {
	var n uint32
	switch val := node.(type) {
	case *ir.Object:
		n, err = e.encodeNodeObject(val)
	case *ir.Array:
		n, err = e.encodeNodeArray(val)
	case *ir.Value:
		n, err = e.encodeNodeValue(val)
	default:
		err = fmt.Errorf("unsupported type %T", val)
	}

	if err != nil {
		err = fmt.Errorf("encode error: %w", err)
		return
	}

	out = e.buf[e.offset-n : e.offset]
	e.offset = 0
	return
}

func (e *encoder) EncodeStream(in any) (n int, err error) {
	if e.w == nil {
		return 0, fmt.Errorf("writer cannot be nil")
	}

	var node ir.Node
	node, err = ValueToNode(in, "")
	if err != nil {
		return
	}

	var bs []byte
	bs, err = e.Encode(node)
	if err != nil {
		return
	}

	n, err = e.w.Write(bs)
	if err != nil {
		return n, fmt.Errorf("write to writer failed: %w", err)
	}

	l := len(bs)
	if n != l {
		return n, fmt.Errorf("short write: wrote %d bytes, expected %d", n, l)
	}

	return n, nil
}

func (e *encoder) encodeSimple(value SimpleValue) (n uint32, err error) {
	sign := Simple
	sign |= Prefix(value)
	return e.writeByte(byte(sign))
}

func (e *encoder) encodeBool(v bool) (n uint32, err error) {
	value := SimpleTrue
	if !v {
		value = SimpleFalse
	}
	return e.encodeSimple(value)
}

func (e *encoder) encodeUint64(uv uint64) (n uint32, err error) {
	return e.encodeInt(PositiveInt, uv)
}

func (e *encoder) encodeInt64(v int64) (n uint32, err error) {
	var sign Prefix
	var uv uint64

	if v >= 0 {
		sign = PositiveInt
		uv = uint64(v)
	} else {
		sign = NegativeInt
		if v == math.MinInt64 {
			uv = 9223372036854775808
		} else {
			uv = uint64(-v)
		}
	}

	return e.encodeInt(sign, uv)
}

func printBytes(buf []byte) {
	for i, b := range buf {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%02x", b)
	}
	fmt.Println()
}

func (e *encoder) encodeBigInt(s string) (n uint32, err error) {
	e.writeByte(byte(len(s)))
	n = utils.EncodeBigInt(e, s)
	// fmt.Println("n", n)
	return e.encodeBytes(e.buf[e.offset-n-1 : e.offset])
}

func (e *encoder) encodeDateTime(t time.Time) (n uint32, err error) {
	v := t.Unix()

	var sign Prefix
	var uv uint64

	if v >= 0 {
		sign = PositiveInt
		uv = uint64(v)
	} else {
		sign = NegativeInt
		if v == math.MinInt64 {
			uv = 9223372036854775808
		} else {
			uv = uint64(-v)
		}
	}

	return e.encodeInt(sign, uv)
}

func (e *encoder) encodeDate(t time.Time) (n uint32, err error) {
	v1 := t.UTC().Truncate(24 * time.Hour)
	v := int64(v1.Sub(utils.DefaultTime).Hours() / 24)

	var sign Prefix
	var uv uint64

	if v >= 0 {
		sign = PositiveInt
		uv = uint64(v)
	} else {
		sign = NegativeInt
		if v == math.MinInt64 {
			uv = 9223372036854775808
		} else {
			uv = uint64(-v)
		}
	}

	return e.encodeInt(sign, uv)
}

func (e *encoder) encodeTime(t time.Time) (n uint32, err error) {
	v1 := t.UTC()
	v := v1.Hour()*3600 + v1.Minute()*60 + v1.Second()
	return e.encodeInt(PositiveInt, uint64(v))
}

func (e *encoder) getEncodedBytes(written uint32) []byte {
	return e.buf[e.offset-written : e.offset]
}
