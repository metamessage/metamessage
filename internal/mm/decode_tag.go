package mm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
	"github.com/lizongying/meta-message/internal/utils"
)

func (d *decoder) decodeTagBytes(tag *ast.Tag) (length int, err error) {
	if tag == nil {
		return 0, fmt.Errorf("tag nil")
	}

	var b byte
	if b, err = d.ReadByte(); err != nil {
		return
	}

	prefix := b & 0xF8

	l := int(b) & 0x07
	switch ast.TagKey(prefix) {
	case ast.KIsNull:
		tag.IsNull = l&0x01 == 1
		if tag.IsNull {
			tag.Nullable = true
		}
		length = 1

	case ast.KExample:
		tag.Example = l&0x01 == 1
		length = 1

	case ast.KDesc:
		switch {
		case l <= 5:
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}

			tag.Desc = string(bs)
			length = 1 + l
		case l <= 1<<8:
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}

			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Desc = string(bs)
			length = 1 + 1 + l
		case l <= 1<<16:
			var l2 []byte
			l2, err = d.ReadBytes(2)
			if err != nil {
				return
			}

			l = int(l2[0])<<8 | int(l2[1])
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Desc = string(bs)
			length = 1 + 2 + l
		default:
			// err = fmt.Errorf("desc too long")
			// return
		}

	case ast.KType:
		var b byte
		b, err = d.ReadByte()
		if err != nil {
			return
		}
		tag.Type = ast.ValueType(b)
		length = 1 + 1

	case ast.KRaw:
		tag.Raw = l&0x01 == 1
		length = 1

	case ast.KNullable:
		tag.Nullable = l&0x01 == 1
		length = 1

	case ast.KDefault:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Default = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Default = string(bs)
			length = 1 + 1 + l
		}

	case ast.KMin:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Min = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Min = string(bs)
			length = 1 + 1 + l
		}

	case ast.KMax:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Max = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Max = string(bs)
			length = 1 + 1 + l
		}

	case ast.KSize:
		if l < 8 {
			for i := 0; i < l; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.Size = tag.Size<<8 | int(b)
			}
			length = 1 + l
		}

	case ast.KEnum:
		tag.Type = ast.ValueTypeEnum
		switch {
		case l <= 5:
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}

			tag.Enum = string(bs)
			length = 1 + l
		case l <= 1<<8:
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}

			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Enum = string(bs)
			length = 1 + 1 + l
		case l <= 1<<16:
			var l2 []byte
			l2, err = d.ReadBytes(2)
			if err != nil {
				return
			}

			l = int(l2[0])<<8 | int(l2[1])
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Enum = string(bs)
			length = 1 + 2 + l
		default:
			// err = fmt.Errorf("enum too long")
			// return
		}

	case ast.KPattern:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Pattern = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.Pattern = string(bs)
			length = 1 + 1 + l
		}

	case ast.KLocation:
		var bs []byte
		bs, err = d.ReadBytes(l)
		if err != nil {
			return
		}
		var num int
		num, err = strconv.Atoi(string(bs))
		tag.Location = utils.IntToLocation(num)
		length = 1 + l

	case ast.KVersion:
		if l < 8 {
			for i := 0; i < l; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.Version = tag.Version<<8 | int(b)
			}
			length = 1 + l
		}

	case ast.KMime:
		if l < 7 {
			tag.Mime = ast.MIME(uint8(l)).String()
			length = 1
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			tag.Mime = ast.MIME(uint8(l2)).String()
			length = 1 + 1
		}

	case ast.KChildDesc:
		switch {
		case l <= 5:
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}

			tag.ChildDesc = string(bs)
			length = 1 + l
		case l <= 1<<8:
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}

			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildDesc = string(bs)
			length = 1 + 1 + l
		case l <= 1<<16:
			var l2 []byte
			l2, err = d.ReadBytes(2)
			if err != nil {
				return
			}

			l = int(l2[0])<<8 | int(l2[1])
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildDesc = string(bs)
			length = 1 + 2 + l
		default:
			// err = fmt.Errorf("child desc too long")
			// return
		}

	case ast.KChildType:
		var b byte
		b, err = d.ReadByte()
		if err != nil {
			return
		}
		tag.ChildType = ast.ValueType(b)
		length = 1 + 1

	case ast.KChildRaw:
		tag.ChildRaw = l&0x01 == 1
		length = 1

	case ast.KChildNullable:
		tag.ChildNullable = l&0x01 == 1
		length = 1

	case ast.KChildDefault:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildDefault = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildDefault = string(bs)
			length = 1 + 1 + l
		}

	case ast.KChildMin:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildMin = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildMin = string(bs)
			length = 1 + 1 + l
		}

	case ast.KChildMax:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildMax = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildMax = string(bs)
			length = 1 + 1 + l
		}

	case ast.KChildSize:
		if l < 8 {
			for i := 0; i < l; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.ChildSize = tag.ChildSize<<8 | int(b)
			}
			length = 1 + l
		}

	case ast.KChildEnum:
		tag.ChildType = ast.ValueTypeEnum
		switch {
		case l <= 5:
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}

			tag.ChildEnum = string(bs)
			length = 1 + l
		case l <= 1<<8:
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}

			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildEnum = string(bs)
			length = 1 + 1 + l
		case l <= 1<<16:
			var l2 []byte
			l2, err = d.ReadBytes(2)
			if err != nil {
				return
			}

			l = int(l2[0])<<8 | int(l2[1])
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildEnum = string(bs)
			length = 1 + 2 + l
		default:
			// err = fmt.Errorf("child enum too long")
			// return
		}

	case ast.KChildPattern:
		if l < 7 {
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildPattern = string(bs)
			length = 1 + l
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			l = int(l2)
			var bs []byte
			bs, err = d.ReadBytes(l)
			if err != nil {
				return
			}
			tag.ChildPattern = string(bs)
			length = 1 + 1 + l
		}

	case ast.KChildLocation:
		var bs []byte
		bs, err = d.ReadBytes(l)
		if err != nil {
			return
		}
		var num int
		num, err = strconv.Atoi(string(bs))
		tag.ChildLocation = utils.IntToLocation(num)
		length = 1 + l

	case ast.KChildVersion:
		if l < 8 {
			for i := 0; i < l; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.ChildVersion = tag.ChildVersion<<8 | int(b)
			}
			length = 1 + l
		}

	case ast.KChildMime:
		if l < 7 {
			tag.ChildMime = ast.MIME(uint8(l)).String()
			length = 1
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			tag.ChildMime = ast.MIME(uint8(l2)).String()
			length = 1 + 1
		}

	default:
		err = errors.New("invalid data")
		return
	}

	return
}
