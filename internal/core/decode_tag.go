package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/utils"
)

func (d *decoder) decodeTagBytes(tag *ir.Tag) (length int, err error) {
	if tag == nil {
		return 0, fmt.Errorf("tag nil")
	}

	var b byte
	if b, err = d.ReadByte(); err != nil {
		return
	}

	prefix := b & 0xF8

	l := int(b) & 0x07
	switch ir.TagKey(prefix) {
	case ir.KIsNull:
		tag.IsNull = l&0x01 == 1
		if tag.IsNull {
			tag.Nullable = true
		}
		length = 1

	case ir.KExample:
		tag.Example = l&0x01 == 1
		length = 1

	case ir.KDesc:
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

	case ir.KType:
		var b byte
		b, err = d.ReadByte()
		if err != nil {
			return
		}
		tag.Type = ir.ValueType(b)
		length = 1 + 1

	case ir.KRaw:
		tag.Raw = l&0x01 == 1
		length = 1

	case ir.KNullable:
		tag.Nullable = l&0x01 == 1
		length = 1

	case ir.KDefault:
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

	case ir.KMin:
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

	case ir.KMax:
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

	case ir.KSize:
		if l < 8 {
			for i := 0; i < l+1; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.Size = tag.Size<<8 | int(b)
			}
			length = 2 + l
			return
		}
		err = fmt.Errorf("size is too large")
		return

	case ir.KEnum:
		tag.Type = ir.ValueTypeEnum
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

	case ir.KPattern:
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

	case ir.KLocation:
		var bs []byte
		bs, err = d.ReadBytes(l)
		if err != nil {
			return
		}
		var num int
		num, err = strconv.Atoi(string(bs))
		tag.Location = utils.IntToLocation(num)
		length = 1 + l

	case ir.KVersion:
		if l < 8 {
			for i := 0; i < l+1; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.Version = tag.Version<<8 | int(b)
			}
			length = 2 + l
		}

	case ir.KMime:
		if l < 7 {
			tag.Mime = ir.MIME(uint8(l)).String()
			length = 1
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			tag.Mime = ir.MIME(uint8(l2)).String()
			length = 1 + 1
		}

	case ir.KChildDesc:
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

	case ir.KChildType:
		var b byte
		b, err = d.ReadByte()
		if err != nil {
			return
		}
		tag.ChildType = ir.ValueType(b)
		length = 1 + 1

	case ir.KChildRaw:
		tag.ChildRaw = l&0x01 == 1
		length = 1

	case ir.KChildNullable:
		tag.ChildNullable = l&0x01 == 1
		length = 1

	case ir.KChildDefault:
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

	case ir.KChildMin:
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

	case ir.KChildMax:
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

	case ir.KChildSize:
		if l < 8 {
			for i := 0; i < l+1; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.ChildSize = tag.ChildSize<<8 | int(b)
			}
			length = 2 + l
		}

	case ir.KChildEnum:
		tag.ChildType = ir.ValueTypeEnum
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

	case ir.KChildPattern:
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

	case ir.KChildLocation:
		var bs []byte
		bs, err = d.ReadBytes(l)
		if err != nil {
			return
		}
		var num int
		num, err = strconv.Atoi(string(bs))
		tag.ChildLocation = utils.IntToLocation(num)
		length = 1 + l

	case ir.KChildVersion:
		if l < 8 {
			for i := 0; i < l+1; i++ {
				var b byte
				b, err = d.ReadByte()
				if err != nil {
					return
				}
				tag.ChildVersion = tag.ChildVersion<<8 | int(b)
			}
			length = 2 + l
		}

	case ir.KChildMime:
		if l < 7 {
			tag.ChildMime = ir.MIME(uint8(l)).String()
			length = 1
		} else {
			var l2 byte
			l2, err = d.ReadByte()
			if err != nil {
				return
			}
			tag.ChildMime = ir.MIME(uint8(l2)).String()
			length = 1 + 1
		}

	default:
		err = errors.New("invalid data")
		return
	}

	return
}
