package mm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/metamessage/metamessage/internal/ast"
	"github.com/metamessage/metamessage/internal/utils"
)

type Decoder interface {
	Reset(r io.Reader)
	Decode(encoded []byte) (node ast.Node, err error)
	DecodeStream(out any) (n int, err error)
}

type decoder struct {
	r      io.Reader
	buf    []byte
	data   []byte
	offset uint64
}

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:    r,
		buf:  make([]byte, defaultBufSize),
		data: make([]byte, 0, defaultBufSize),
	}
	return d
}

func (d *decoder) Reset(r io.Reader) {
	if r == nil {
		return
	}

	d.r = r
}

// ReadByte Implement io.ByteReader to read a byte from the internal buffer at the current offset and move the offset forward.
func (d *decoder) ReadByte() (byte, error) {
	if int(d.offset) >= len(d.data) {
		return 0, io.EOF
	}
	b := d.data[d.offset]
	d.offset++
	return b, nil
}

// ReadBytes Read n bytes from the internal buffer at the current offset and move the offset forward by n.
func (d *decoder) ReadBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("invalid length")
	}
	if int(d.offset)+n > len(d.data) {
		return nil, io.EOF
	}
	start := int(d.offset)
	bs := d.data[start : start+n]
	d.offset += uint64(n)
	return bs, nil
}

// Read Implement io.Reader to read data from the current offset of the internal buffer.
func (d *decoder) Read(p []byte) (int, error) {
	if d == nil || p == nil || len(p) == 0 {
		return 0, nil
	}
	if int(d.offset) >= len(d.data) {
		return 0, io.EOF
	}
	n := copy(p, d.data[int(d.offset):])
	d.offset += uint64(n)
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (d *decoder) readAllWithDynamicBuf() (int, error) {
	all := 0
	for {
		n, err := d.r.Read(d.buf)
		if n > 0 {
			d.data = append(d.data, d.buf[:n]...)
			all += n
		}
		if err == io.EOF {
			return all, io.EOF
		}
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, errors.New("zero read without error")
		}
	}
}

func (d *decoder) Decode(encoded []byte) (node ast.Node, err error) {
	d.data = encoded
	d.offset = 0
	node, _, err = d.decode(nil, "")
	return
}

func (d *decoder) DecodeStream(out any) (n int, err error) {
	n, err = d.readAllWithDynamicBuf()

	eof := false
	if err != nil {
		if err != io.EOF {
			err = errors.New("eof decode: " + err.Error())
			return
		} else {
			eof = true
		}
	}

	d.offset = 0
	var node ast.Node
	node, _, err = d.decode(nil, "")
	if err != nil {
		return
	}

	// fmt.Println("node", jsonc.Json(node))

	err = Bind(node, out)
	if err != nil {
		return
	}

	if eof {
		err = io.EOF
	}
	return
}

func (d *decoder) decode(tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	var b byte
	b, err = d.ReadByte()
	prefix := GetPrefix(b)
	// fmt.Printf("decode prefix: %s %b\n", prefix.String(), b)

	switch prefix {
	case PrefixTag:
		return d.decodeTag(b, path)
	case Simple:
		return d.decodeSimple(b, tag, path)
	case PositiveInt:
		return d.decodePositiveInt(b, tag, path)
	case NegativeInt:
		return d.decodeNegativeInt(b, tag, path)
	case PrefixFloat:
		return d.decodeFloat(b, tag, path)
	case PrefixString:
		return d.decodeString(b, tag, path)
	case PrefixBytes:
		return d.decodeBytes(b, tag, path)
	case Container:
		return d.decodeContainer(b, tag, path)
	default:
		err = errors.New("invalid prefix")
		return
	}
}

func (d *decoder) decodeTag(prefix byte, path string) (node ast.Node, length int, err error) {
	l1, l2 := TagLen(prefix)

	switch l1 {
	case 0:
	case 1:
		var l byte
		if l, err = d.ReadByte(); err != nil {
			return
		}
		l2 = int(l)
	case 2:
		var l []byte
		if l, err = d.ReadBytes(2); err != nil {
			return
		}
		l2 = int(l[0])<<8 | int(l[1])
	default:
		err = errors.New("invalid data")
		return
	}

	tag := ast.NewTag()

	var b byte
	b, err = d.ReadByte()
	if err != nil {
		return
	}

	l := int(b)
	if l < 254 {
	} else if l < 257 {
		b, err = d.ReadByte()
		if err != nil {
			return
		}
		l = int(b)
	} else {
		var l3 []byte
		if l3, err = d.ReadBytes(2); err != nil {
			return
		}
		l = int(l3[0])<<8 | int(l3[1])
	}

	for l > 0 {
		var n int
		if n, err = d.decodeTagBytes(tag); err != nil {
			return
		}

		if n == 0 {
			err = errors.New("tag error")
			return
		}

		if n > l {
			err = errors.New("tag overflow")
			return
		}

		l -= n
	}

	// fmt.Println("decode tag", tag.String())
	if tag.IsNull {
		switch tag.Type {
		case ast.ValueTypeDateTime:
			var text string
			if tag.Location == nil {
				text = utils.DefaultTime.Format(time.DateTime)
			} else {
				text = utils.DefaultTime.In(tag.Location).Format(time.DateTime)
			}
			node = &ast.Value{
				Data: utils.DefaultTime,
				Text: text,
				Tag:  tag,
			}
		case ast.ValueTypeDate:
			var text string
			if tag.Location == nil {
				text = utils.DefaultTime.Format(time.DateOnly)
			} else {
				text = utils.DefaultTime.In(tag.Location).Format(time.DateOnly)
			}
			node = &ast.Value{
				Data: utils.DefaultTime,
				Text: text,
				Tag:  tag,
			}
		case ast.ValueTypeTime:
			var text string
			if tag.Location == nil {
				text = utils.DefaultTime.Format(time.TimeOnly)
			} else {
				text = utils.DefaultTime.In(tag.Location).Format(time.TimeOnly)
			}
			node = &ast.Value{
				Data: utils.DefaultTime,
				Text: text,
				Tag:  tag,
			}

		case ast.ValueTypeInt8:
			node = &ast.Value{
				Data: int8(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeInt16:
			node = &ast.Value{
				Data: int16(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeInt32:
			node = &ast.Value{
				Data: int32(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeInt64:
			node = &ast.Value{
				Data: int64(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeUint:
			node = &ast.Value{
				Data: uint(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeUint8:
			node = &ast.Value{
				Data: uint8(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeUint16:
			node = &ast.Value{
				Data: uint16(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeUint32:
			node = &ast.Value{
				Data: uint32(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeUint64:
			node = &ast.Value{
				Data: uint64(0),
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeFloat32:
			node = &ast.Value{
				Data: float32(0.0),
				Text: "0.0",
				Tag:  tag,
			}
		case ast.ValueTypeEmail, ast.ValueTypeUUID, ast.ValueTypeDecimal:
			node = &ast.Value{
				Data: "",
				Text: "",
				Tag:  tag,
			}
		case ast.ValueTypeBigInt:
			node = &ast.Value{
				Data: big.Int{},
				Text: "0",
				Tag:  tag,
			}
		case ast.ValueTypeURL:
			node = &ast.Value{
				Data: url.URL{},
				Text: "",
				Tag:  tag,
			}
		case ast.ValueTypeIP:
			var text string
			switch tag.Version {
			case 0:
				text = ""
			case 4:
				text = "0.0.0.0"
			case 6:
				text = "::"
			default:
				err = fmt.Errorf("unsupported IP version: %v", tag.Version)
				return
			}
			node = &ast.Value{
				Data: net.IP{},
				Text: text,
				Tag:  tag,
			}
		default:
			node, length, err = d.decode(tag, path)
			if err != nil {
				return
			}
		}
	} else {
		node, length, err = d.decode(tag, path)
		if err != nil {
			return
		}
	}

	length = l1 + 1 + l2
	return
}

func (d *decoder) decodeString(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	l1, l2 := StringLen(prefix)

	switch l1 {
	case 0:
	case 1:
		var l byte
		if l, err = d.ReadByte(); err != nil {
			return
		}
		l2 = int(l)
	case 2:
		var l []byte
		if l, err = d.ReadBytes(2); err != nil {
			return
		}
		l2 = int(l[0])<<8 | int(l[1])
	default:
		err = errors.New("invalid data")
		return
	}

	var bs []byte
	if l2 > 0 {
		bs, err = d.ReadBytes(l2)
		if err != nil {
			return
		}
	}

	var data any
	text := string(bs)
	if tag == nil {
		tag = ast.NewTag()
	}

	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeString
	}

	switch tag.Type {
	case ast.ValueTypeEmail:
		data = text

	case ast.ValueTypeURL:
		var u *url.URL
		u, err = url.Parse(text)
		if err != nil {
			return
		}
		data = *u

	case ast.ValueTypeIP:
		d := net.ParseIP(text)
		data = d
		text = d.String()

	case ast.ValueTypeString:
		data = text

	default:
		err = fmt.Errorf("unsupported string type: %v", tag.Type)
		return
	}

	node = &ast.Value{
		Data: data,
		Text: text,
		Tag:  tag,
	}

	length = l1 + 1 + l2
	return
}

func (d *decoder) decodeBytes(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	l1, l2 := BytesLen(prefix)

	switch l1 {
	case 0:
	case 1:
		var l byte
		if l, err = d.ReadByte(); err != nil {
			return
		}
		l2 = int(l)
	case 2:
		var l []byte
		if l, err = d.ReadBytes(2); err != nil {
			return
		}
		l2 = int(l[0])<<8 | int(l[1])
	default:
		err = errors.New("invalid data")
		return
	}

	var bs []byte
	if l2 > 0 {
		bs, err = d.ReadBytes(l2)
		if err != nil {
			return
		}
	} else {
		bs = []byte{}
	}

	var data any
	text := ""
	if tag == nil {
		tag = ast.NewTag()
	}

	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeBytes
	}

	switch tag.Type {
	case ast.ValueTypeBigInt:
		text, err = utils.DecodeBigInt(bs[1:], int(bs[0]))
		if err != nil {
			return
		}
		var ok bool
		data, ok = new(big.Int).SetString(text, 10)
		if !ok {
			err = fmt.Errorf("Failed to parse big integer: Invalid decimal string: '%s'", text)
			return
		}

	case ast.ValueTypeBytes:
		data = bs
		text = base64.StdEncoding.EncodeToString(bs)

	case ast.ValueTypeUUID:
		d := [16]byte(bs)
		data = d
		text = utils.BytesToUUIDString(d)

	case ast.ValueTypeIP:
		data = net.IP(bs)
		text = data.(net.IP).String()

	default:
		err = fmt.Errorf("unsupported string types: %v", tag.Type)
		return
	}

	node = &ast.Value{
		Data: data,
		Text: text,
		Tag:  tag,
	}

	length = l1 + 1 + l2
	return
}

func (d *decoder) decodePositiveInt(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	l1, l2 := IntLen(prefix)

	var v uint64
	var l []byte
	switch l1 {
	case 0:
		v = uint64(l2)
	case 1:
		if l, err = d.ReadBytes(1); err != nil {
			return
		}
		v = uint64(l[0])
	case 2:
		if l, err = d.ReadBytes(2); err != nil {
			return
		}
		v = uint64(l[0])<<8 | uint64(l[1])
	case 3:
		if l, err = d.ReadBytes(3); err != nil {
			return
		}
		v = uint64(l[0])<<16 | uint64(l[1])<<8 | uint64(l[2])
	case 4:
		if l, err = d.ReadBytes(4); err != nil {
			return
		}
		v = uint64(l[0])<<24 | uint64(l[1])<<16 | uint64(l[2])<<8 | uint64(l[3])
	case 5:
		if l, err = d.ReadBytes(5); err != nil {
			return
		}
		v = uint64(l[0])<<32 | uint64(l[1])<<24 | uint64(l[2])<<16 | uint64(l[3])<<8 | uint64(l[4])
	case 6:
		if l, err = d.ReadBytes(6); err != nil {
			return
		}
		v = uint64(l[0])<<40 | uint64(l[1])<<32 | uint64(l[2])<<24 | uint64(l[3])<<16 | uint64(l[4])<<8 | uint64(l[5])
	case 7:
		if l, err = d.ReadBytes(7); err != nil {
			return
		}
		v = uint64(l[0])<<48 | uint64(l[1])<<40 | uint64(l[2])<<32 | uint64(l[3])<<24 | uint64(l[4])<<16 | uint64(l[5])<<8 | uint64(l[6])
	case 8:
		if l, err = d.ReadBytes(8); err != nil {
			return
		}
		v = uint64(l[0])<<56 | uint64(l[1])<<48 | uint64(l[2])<<40 | uint64(l[3])<<32 | uint64(l[4])<<24 | uint64(l[5])<<16 | uint64(l[6])<<8 | uint64(l[7])
	default:
		err = fmt.Errorf("unsupported numerical length: %v", l2)
		return
	}

	var data any
	text := strconv.FormatUint(v, 10)
	if tag == nil {
		tag = ast.NewTag()
	}

	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeInt
	}

	switch tag.Type {
	case ast.ValueTypeInt:
		data = int(v)
	case ast.ValueTypeInt8:
		data = int8(v)
	case ast.ValueTypeInt16:
		data = int16(v)
	case ast.ValueTypeInt32:
		data = int32(v)
	case ast.ValueTypeInt64:
		data = int64(v)
	case ast.ValueTypeUint:
		data = uint(v)
	case ast.ValueTypeUint8:
		data = uint8(v)
	case ast.ValueTypeUint16:
		data = uint16(v)
	case ast.ValueTypeUint32:
		data = uint32(v)
	case ast.ValueTypeUint64:
		data = uint64(v)
	case ast.ValueTypeDateTime:
		if tag.IsNull {
			data = nil
			text = ""
		} else {
			if v > math.MaxInt64 {
				err = errors.New("decodeDateTime: time value out of math.MaxInt64")
				return
			}
			d := time.Unix(int64(v), 0)
			if tag.Location != nil {
				d = d.In(tag.Location)
			} else {
				d = d.UTC()
			}
			data = d
			text = d.Format(time.DateTime)
		}

	case ast.ValueTypeDate:
		if tag.IsNull {
			data = nil
			text = ""
		} else {
			if v > math.MaxInt {
				err = errors.New("decodeDate: time value out of math.MaxInt")
				return
			}
			d := utils.DefaultTime.AddDate(0, 0, int(v)).Truncate(24 * time.Hour)
			if tag.Location != nil {
				d = d.In(tag.Location)
			} else {
				d = d.UTC()
			}
			data = d
			text = d.Format(time.DateOnly)
		}

	case ast.ValueTypeTime:
		if tag.IsNull {
			data = nil
			text = ""
		} else {
			if v > 86399 {
				err = errors.New("decodeTime: time value out of range (0-86399)")
				return
			}

			hour := v / 3600
			minute := (v % 3600) / 60
			second := v % 60

			d := time.Date(1970, 1, 1, int(hour), int(minute), int(second), 0, time.UTC)
			if tag.Location != nil {
				d = d.In(tag.Location)
			} else {
				d = d.UTC()
			}
			data = d
			text = d.Format(time.TimeOnly)
		}

	case ast.ValueTypeEnum:
		if tag.IsNull {
			data = -1
			text = ""
		} else {
			if tag.Enum != "" {
				enums := strings.Split(tag.Enum, "|")
				d := int(v)
				if d >= len(enums) {
					err = fmt.Errorf("enum index out of range")
					return
				}
				data = d
				text = strings.TrimSpace(enums[d])
			} else {
				err = fmt.Errorf("only enum are supported")
				return
			}
		}

	default:
		err = fmt.Errorf("unsupported value types: %v", tag.Type)
		return
	}

	node = &ast.Value{
		Data: data,
		Text: text,
		Tag:  tag,
	}

	length = l1 + 1
	return
}

func (d *decoder) decodeNegativeInt(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	l1, l2 := IntLen(prefix)

	var v uint64
	var l []byte
	switch l1 {
	case 0:
		v = uint64(l2)
	case 1:
		if l, err = d.ReadBytes(1); err != nil {
			return
		}
		v = uint64(l[0])
	case 2:
		if l, err = d.ReadBytes(2); err != nil {
			return
		}
		v = uint64(l[0])<<8 | uint64(l[1])
	case 3:
		if l, err = d.ReadBytes(3); err != nil {
			return
		}
		v = uint64(l[0])<<16 | uint64(l[1])<<8 | uint64(l[2])
	case 4:
		if l, err = d.ReadBytes(4); err != nil {
			return
		}
		v = uint64(l[0])<<24 | uint64(l[1])<<16 | uint64(l[2])<<8 | uint64(l[3])
	case 5:
		if l, err = d.ReadBytes(5); err != nil {
			return
		}
		v = uint64(l[0])<<32 | uint64(l[1])<<24 | uint64(l[2])<<16 | uint64(l[3])<<8 | uint64(l[4])
	case 6:
		if l, err = d.ReadBytes(6); err != nil {
			return
		}
		v = uint64(l[0])<<40 | uint64(l[1])<<32 | uint64(l[2])<<24 | uint64(l[3])<<16 | uint64(l[4])<<8 | uint64(l[5])
	case 7:
		if l, err = d.ReadBytes(7); err != nil {
			return
		}
		v = uint64(l[0])<<48 | uint64(l[1])<<40 | uint64(l[2])<<32 | uint64(l[3])<<24 | uint64(l[4])<<16 | uint64(l[5])<<8 | uint64(l[6])
	case 8:
		if l, err = d.ReadBytes(8); err != nil {
			return
		}
		v = uint64(l[0])<<56 | uint64(l[1])<<48 | uint64(l[2])<<40 | uint64(l[3])<<32 | uint64(l[4])<<24 | uint64(l[5])<<16 | uint64(l[6])<<8 | uint64(l[7])
	default:
		err = fmt.Errorf("unsupported numerical length: %v", l2)
		return
	}

	var data any
	text := "-" + strconv.FormatUint(v, 10)
	if tag == nil {
		tag = ast.NewTag()
	}

	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeInt
	}

	if tag != nil {
		switch tag.Type {
		case ast.ValueTypeInt:
			data = -int(v)
		case ast.ValueTypeInt8:
			data = -int8(v)
		case ast.ValueTypeInt16:
			data = -int16(v)
		case ast.ValueTypeInt32:
			data = -int32(v)
		case ast.ValueTypeInt64:
			data = -int64(v)
		case ast.ValueTypeDateTime:
			if tag.IsNull {
				data = nil
				text = ""
			} else {
				// TODO
				// if v < math.MinInt64 {
				// 	err = errors.New("decodeDateTime: time value out of math.MinInt64")
				// 	return
				// }
				d := time.Unix(-int64(v), 0)
				if tag.Location != nil {
					d = d.In(tag.Location)
				} else {
					d = d.UTC()
				}
				if tag.Nullable {
					data = &d
				} else {
					data = d
				}
				text = d.Format(time.DateTime)
			}

		case ast.ValueTypeDate:
			if tag.IsNull {
				data = nil
				text = ""
			} else {
				// TODO
				// if v < math.MinInt {
				// 	err = errors.New("decodeDate: time value out of math.MinInt")
				// 	return
				// }
				d := utils.DefaultTime.AddDate(0, 0, -int(v)).Truncate(24 * time.Hour)
				if tag.Location != nil {
					d = d.In(tag.Location)
				} else {
					d = d.UTC()
				}
				if tag.Nullable {
					data = &d
				} else {
					data = d
				}
				text = d.Format(time.DateOnly)
			}

		default:
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
	}

	node = &ast.Value{
		Data: data,
		Text: text,
		Tag:  tag,
	}

	length = l1 + 1
	return
}

func mantissaToDecimal(mantissa uint64, exp int8) string {
	numStr := strconv.FormatUint(mantissa, 10)

	decimalPos := len(numStr) + int(exp)

	var result string
	switch {
	case decimalPos <= 0:
		result = "0." + strings.Repeat("0", -decimalPos) + numStr

	case decimalPos > 0 && decimalPos < len(numStr):
		result = numStr[:decimalPos] + "." + numStr[decimalPos:]

	default:
		trailingZeros := decimalPos - len(numStr)
		result = numStr + strings.Repeat("0", trailingZeros)
	}

	// result = cleanTrailingZeros(result)

	return result
}

func cleanTrailingZeros(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}

	// for strings.HasSuffix(s, "0") {
	// 	s = s[:len(s)-1]
	// }

	s = strings.TrimSuffix(s, ".")

	return s
}

func (d *decoder) decodeFloat(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	l1, l2 := FloatLen(prefix)

	p := Prefix(prefix)

	v := 0.0

	switch {
	// 0.x
	case p >= PrefixFloat && p <= PrefixFloat+7:
		v = float64(p&0xF) / 10
		length = 1

	// // 0.0x
	// case p >= PositiveFloat+10 && p <= PositiveFloat+19:
	// 	v = float64(p&0xF-10) / 100
	// 	length = 1

	// // 0–3
	// case p >= PositiveFloat+20 && p <= PositiveFloat+23:
	// 	v = float64(p&0xF - 20)
	// 	length = 1

	default:
		var b byte
		if b, err = d.ReadByte(); err != nil {
			return
		}

		exp := int8(b)
		var mantissa uint64
		var l []byte
		switch l1 {
		case 0:
			mantissa = uint64(l2)
		case 1:
			if l, err = d.ReadBytes(1); err != nil {
				return
			}
			mantissa = uint64(l[0])
		case 2:
			if l, err = d.ReadBytes(2); err != nil {
				return
			}
			mantissa = uint64(l[0])<<8 | uint64(l[1])
		case 3:
			if l, err = d.ReadBytes(3); err != nil {
				return
			}
			mantissa = uint64(l[0])<<16 | uint64(l[1])<<8 | uint64(l[2])
		case 4:
			if l, err = d.ReadBytes(4); err != nil {
				return
			}
			mantissa = uint64(l[0])<<24 | uint64(l[1])<<16 | uint64(l[2])<<8 | uint64(l[3])
		case 5:
			if l, err = d.ReadBytes(5); err != nil {
				return
			}
			mantissa = uint64(l[0])<<32 | uint64(l[1])<<24 | uint64(l[2])<<16 | uint64(l[3])<<8 | uint64(l[4])
		case 6:
			if l, err = d.ReadBytes(6); err != nil {
				return
			}
			mantissa = uint64(l[0])<<40 | uint64(l[1])<<32 | uint64(l[2])<<24 | uint64(l[3])<<16 | uint64(l[4])<<8 | uint64(l[5])
		case 7:
			if l, err = d.ReadBytes(7); err != nil {
				return
			}
			mantissa = uint64(l[0])<<48 | uint64(l[1])<<40 | uint64(l[2])<<32 | uint64(l[3])<<24 | uint64(l[4])<<16 | uint64(l[5])<<8 | uint64(l[6])
		case 8:
			if l, err = d.ReadBytes(8); err != nil {
				return
			}
			mantissa = uint64(l[0])<<56 | uint64(l[1])<<48 | uint64(l[2])<<40 | uint64(l[3])<<32 | uint64(l[4])<<24 | uint64(l[5])<<16 | uint64(l[6])<<8 | uint64(l[7])
		default:
			err = fmt.Errorf("unsupported numerical length: %v", l2)
			return
		}

		// TODO mantissa 1<<53 - 1
		// v = float64(mantissa) * math.Pow10(int(exp))

		v, err = strconv.ParseFloat(mantissaToDecimal(mantissa, exp), 10)
		if err != nil {
			err = fmt.Errorf("Failed to convert to float64 | Mantissa = %d | Exponent = %d | Original error: %w", mantissa, exp, err)
			return
		}
		length = l1 + 2
	}

	if p&FloatPositiveNegativeMask != 0 {
		v = -v
	}

	var data any
	text := ""
	if tag == nil {
		tag = ast.NewTag()
	}
	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeFloat64
	}

	switch tag.Type {
	case ast.ValueTypeFloat32:
		data = float32(v)
		text = utils.FormatFloat32(data.(float32))
	case ast.ValueTypeFloat64:
		data = v
		text = utils.FormatFloat64(data.(float64))
	case ast.ValueTypeDecimal:
		data = v
		text = utils.FormatFloat64(data.(float64))
	default:
		err = fmt.Errorf("unsupported value types: %v", tag.Type)
		return
	}

	node = &ast.Value{
		Data: data,
		Text: text,
		Tag:  tag,
	}

	return
}

func (d *decoder) decodeSimple(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	if tag == nil {
		tag = ast.NewTag()
	}

	switch SimpleValue(prefix) {
	case SimpleFalse:
		tag.Type = ast.ValueTypeBool
		node = &ast.Value{
			Data: false,
			Text: ast.False,
			Tag:  tag,
		}

	case SimpleTrue:
		tag.Type = ast.ValueTypeBool
		node = &ast.Value{
			Data: true,
			Text: ast.True,
			Tag:  tag,
		}

	case SimpleNullBool:
		if tag.Type == ast.ValueTypeUnknown {
			tag.Type = ast.ValueTypeBool
		} else if tag.Type != ast.ValueTypeBool {
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
		node = &ast.Value{
			Data: false,
			Text: ast.False,
			Tag:  tag,
		}

	case SimpleNullFloat:
		if tag.Type == ast.ValueTypeUnknown {
			tag.Type = ast.ValueTypeFloat64
		} else if tag.Type != ast.ValueTypeFloat64 && tag.Type != ast.ValueTypeFloat32 {
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
		switch tag.Type {
		case ast.ValueTypeFloat32:
			node = &ast.Value{
				Data: float32(0.0),
				Text: "0.0",
				Tag:  tag,
			}
		case ast.ValueTypeFloat64:
			node = &ast.Value{
				Data: 0.0,
				Text: "0.0",
				Tag:  tag,
			}
		default:
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}

	case SimpleNullInt:
		if tag.Type == ast.ValueTypeUnknown {
			tag.Type = ast.ValueTypeInt
		} else if tag.Type != ast.ValueTypeInt {
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
		switch tag.Type {
		case ast.ValueTypeInt:
			node = &ast.Value{
				Data: 0,
				Text: "0",
				Tag:  tag,
			}

		default:
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}

	case SimpleNullString:
		if tag.Type == ast.ValueTypeUnknown {
			tag.Type = ast.ValueTypeString
		} else if tag.Type != ast.ValueTypeString {
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
		switch tag.Type {
		case ast.ValueTypeString:
			node = &ast.Value{
				Data: "",
				Text: "",
				Tag:  tag,
			}

		default:
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}

	case SimpleNullBytes:
		if tag.Type == ast.ValueTypeUnknown {
			tag.Type = ast.ValueTypeBytes
		} else if tag.Type != ast.ValueTypeBytes {
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}
		switch tag.Type {
		case ast.ValueTypeBytes:
			node = &ast.Value{
				Data: []byte{},
				Text: "",
				Tag:  tag,
			}

		default:
			err = fmt.Errorf("unsupported value types: %v", tag.Type)
			return
		}

	case SimpleCode:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleCodeStr,
			Tag:  tag,
		}
	case SimpleMessage:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleMessageStr,
			Tag:  tag,
		}
	case SimpleData:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleDataStr,
			Tag:  tag,
		}
	case SimpleSuccess:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleSuccessStr,
			Tag:  tag,
		}
	case SimpleError:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleErrorStr,
			Tag:  tag,
		}
	case SimpleUnknown:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleUnknownStr,
			Tag:  tag,
		}
	case SimplePage:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimplePageStr,
			Tag:  tag,
		}
	case SimpleLimit:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleLimitStr,
			Tag:  tag,
		}
	case SimpleOffset:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleOffsetStr,
			Tag:  tag,
		}
	case SimpleTotal:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleTotalStr,
			Tag:  tag,
		}
	case SimpleId:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleIdStr,
			Tag:  tag,
		}
	case SimpleName:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleNameStr,
			Tag:  tag,
		}
	case SimpleDescription:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleDescriptionStr,
			Tag:  tag,
		}
	case SimpleType:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleTypeStr,
			Tag:  tag,
		}
	case SimpleVersion:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleVersionStr,
			Tag:  tag,
		}
	case SimpleStatus:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleStatusStr,
			Tag:  tag,
		}
	case SimpleUrl:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleUrlStr,
			Tag:  tag,
		}
	case SimpleCreateTime:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleCreateTimeStr,
			Tag:  tag,
		}
	case SimpleUpdateTime:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleUpdateTimeStr,
			Tag:  tag,
		}
	case SimpleDeleteTime:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleDeleteTimeStr,
			Tag:  tag,
		}
	case SimpleAccount:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleAccountStr,
			Tag:  tag,
		}
	case SimpleToken:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleTokenStr,
			Tag:  tag,
		}
	case SimpleExpireTime:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleExpireTimeStr,
			Tag:  tag,
		}
	case SimpleKey:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleKeyStr,
			Tag:  tag,
		}
	case SimpleVal:
		tag.Type = ast.ValueTypeString
		node = &ast.Value{
			Data: nil,
			Text: ast.SimpleValStr,
			Tag:  tag,
		}
	default:
		err = fmt.Errorf("unsupported value: %v", prefix)
		return
	}

	length = 1
	return
}

func (d *decoder) decodeContainer(prefix byte, tag *ast.Tag, path string) (node ast.Node, length int, err error) {
	if IsArray(prefix) {
		return d.decodeArray(prefix, tag, path)
	}

	return d.decodeObject(prefix, tag, path)
}

func (d *decoder) decodeArray(prefix byte, tag *ast.Tag, path string) (node *ast.Array, length int, err error) {
	if tag == nil {
		tag = ast.NewTag()
		tag.Type = ast.ValueTypeSlice
	}
	if tag.Type == ast.ValueTypeUnknown {
		if tag.Size > 0 {
			tag.Type = ast.ValueTypeArray
		} else {
			tag.Type = ast.ValueTypeSlice
		}
	}
	l1, l2 := ContainerLen(prefix)

	switch l1 {
	case 0:
	case 1:
		if len(d.data) < 2 {
			err = fmt.Errorf("%s: invalid data", path)
			return
		}
		var l byte
		if l, err = d.ReadByte(); err != nil {
			err = fmt.Errorf("%s: %w", path, err)
			return
		}
		l2 = int(l)
	case 2:
		if len(d.data) < 3 {
			err = fmt.Errorf("%s: invalid data", path)
			return
		}
		var l []byte
		if l, err = d.ReadBytes(2); err != nil {
			err = fmt.Errorf("%s: %w", path, err)
			return
		}
		l2 = int(l[0])<<8 | int(l[1])
	default:
		err = fmt.Errorf("%s: invalid data", path)
		return
	}

	node = &ast.Array{
		Tag:  tag,
		Path: path,
	}

	index := 0
	for index < l2 {
		tagValue := ast.NewTag()
		tagValue.Inherit(tag)

		p := fmt.Sprintf("%s.%s", path, strconv.Itoa(index))
		n, l, e := d.decode(tagValue, p)
		if e != nil || l <= 0 {
			err = fmt.Errorf("%s: %w", p, e)
			return
		}

		node.Items = append(node.Items, n)
		index += l
	}

	length = l1 + 1 + l2
	return
}

func (d *decoder) decodeObject(prefix byte, tag *ast.Tag, path string) (node *ast.Object, length int, err error) {
	if tag == nil {
		tag = ast.NewTag()
		tag.Type = ast.ValueTypeObject
	}
	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeObject
	}
	l1, l2 := ContainerLen(prefix)

	switch l1 {
	case 0:
	case 1:
		if len(d.data) < 2 {
			err = fmt.Errorf("%s: invalid data", path)
			return
		}
		var l byte
		if l, err = d.ReadByte(); err != nil {
			err = fmt.Errorf("%s: %w", path, err)
			return
		}
		l2 = int(l)
	case 2:
		if len(d.data) < 3 {
			err = fmt.Errorf("%s: invalid data", path)
			return
		}
		var l []byte
		if l, err = d.ReadBytes(2); err != nil {
			err = fmt.Errorf("%s: %w", path, err)
			return
		}
		l2 = int(l[0])<<8 | int(l[1])
	default:
		err = fmt.Errorf("%s: invalid data", path)
		return
	}

	node = &ast.Object{
		Tag:  tag,
		Path: path,
	}

	var lArray byte
	if lArray, err = d.ReadByte(); err != nil {
		err = fmt.Errorf("%s: %w", path, err)
		return
	}

	nKeys, lKeys, eKeys := d.decodeArray(lArray, tag, path)
	if eKeys != nil {
		err = fmt.Errorf("%s: %w", path, eKeys)
		return
	}

	index := lKeys
	i := 0
	for index < l2 {
		tagValue := ast.NewTag()
		tagValue.Inherit(tag)

		key := nKeys.Items[i].(*ast.Value).Text
		p := fmt.Sprintf("%s.%s", path, key)
		n, l, e := d.decode(tagValue, p)
		if e != nil || l <= 0 {
			err = fmt.Errorf("%s: %w", p, e)
			return
		}

		node.Fields = append(node.Fields, &ast.Field{
			Key:   key,
			Value: n,
		})
		index += l
		i++
	}

	length = l1 + 1 + l2
	return
}
