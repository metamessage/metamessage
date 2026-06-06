package parser

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc/token"
	"github.com/metamessage/metamessage/internal/utils"
)

const BitSize = 32 << (^uint(0) >> 63)
const maxDepth = 32

type Parser struct {
	toks []token.Token
	pos  int

	pending []token.Token

	depth int
}

func New(tokens []token.Token) *Parser {
	return &Parser{toks: tokens}
}

func (p *Parser) peek() token.Token {
	if p.pos >= len(p.toks) {
		return token.Token{Type: token.EOF}
	}
	return p.toks[p.pos]
}

func (p *Parser) next() token.Token {
	t := p.peek()
	p.pos++
	return t
}

func (p *Parser) consumeCommentsFor(anchorLine int) (*ir.Tag, error) {
	if len(p.pending) == 0 {
		return nil, nil
	}

	last := p.pending[len(p.pending)-1]
	if anchorLine-last.Line > 1 {
		p.pending = nil
		return nil, nil
	}

	var out *ir.Tag
	for _, ct := range p.pending {
		parsed, err := parseCommentsToTag(ct.Literal)
		if err != nil {
			return nil, err
		}
		if parsed == nil {
			continue
		}
		out = ir.MergeTag(out, parsed)
	}

	return out, nil
}

func (p *Parser) Parse() (val ir.Node, err error) {
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			if val == nil {
				val = &ir.NodeNull{}
			}
			return
		}

		if tok.Type == token.Comment {
			if len(p.pending) > 0 {
				last := p.pending[len(p.pending)-1]
				if tok.Line-last.Line > 1 {
					p.pending = nil
				}
			}
			p.pending = append(p.pending, tok)
			p.next()
			continue
		}

		var tag *ir.Tag
		if tag, err = p.consumeCommentsFor(tok.Line); err != nil {
			return nil, err
		}

		if tag == nil {
			tag = ir.NewTag()
		}

		if val, err = p.parse("", false, tag); err != nil {
			return
		}
	}
}

func (p *Parser) parse(path string, example bool, tag *ir.Tag) (val ir.Node, err error) {
	for {
		tok := p.next()
		var data any
		switch tok.Type {
		case token.EOF:
			return nil, nil
		case token.LBrace:
			return p.parseObject(tok.Line, path, tag)
		case token.LBracket:
			return p.parseArray(tok.Line, path, tag)
		case token.String:
			text := tok.Literal

			if tag.Type == ir.ValueTypeUnknown {
				tag.Type = ir.ValueTypeStr
			}

			switch text {
			case ir.SimpleCodeStr,
				ir.SimpleMessageStr,
				ir.SimpleDataStr,
				ir.SimpleSuccessStr,
				ir.SimpleErrorStr,
				ir.SimpleUnknownStr,
				ir.SimplePageStr,
				ir.SimpleLimitStr,
				ir.SimpleOffsetStr,
				ir.SimpleTotalStr,
				ir.SimpleIdStr,
				ir.SimpleNameStr,
				ir.SimpleDescriptionStr,
				ir.SimpleTypeStr,
				ir.SimpleVersionStr,
				ir.SimpleStatusStr,
				ir.SimpleUrlStr,
				ir.SimpleCreateTimeStr,
				ir.SimpleUpdateTimeStr,
				ir.SimpleDeleteTimeStr,
				ir.SimpleAccountStr,
				ir.SimpleTokenStr,
				ir.SimpleExpireTimeStr,
				ir.SimpleKeyStr,
				ir.SimpleValStr:
				tag.Type = ir.ValueTypeStr
				data, text, err = tag.ValidateStr(text, example || tag.Example)
				if err != nil {
					return nil, err
				}

			default:
				switch tag.Type {
				case ir.ValueTypeStr:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid string: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateStr(text, example || tag.Example)
					}

				case ir.ValueTypeBytes:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid bytes: %q, valid: %q", text, "")
						}

						data = []byte{}
					} else {
						var d []byte
						if d, err = base64.StdEncoding.DecodeString(text); err != nil {
							return nil, fmt.Errorf("invalid base64 bytes %q: %w", text, err)
						}

						data, text, err = tag.ValidateBytes(d, example || tag.Example)
					}

				case ir.ValueTypeDatetime:
					location := time.UTC
					if tag.Location != nil {
						location = tag.Location
					}

					if tag.IsNull {
						datetime := utils.DefaultTime.In(location).Format(time.DateTime)
						if text != datetime {
							return nil, fmt.Errorf("invalid datetime: %q, valid: %q", text, datetime)
						}

						data = utils.DefaultTime
					} else {
						var d time.Time
						d, err = time.ParseInLocation(time.DateTime, text, location)
						if err != nil {
							return nil, fmt.Errorf("invalid datetime %q: %w", text, err)
						}

						data, text, err = tag.ValidateDatetime(d, example || tag.Example)
					}

				case ir.ValueTypeDate:
					location := time.UTC
					if tag.Location != nil {
						location = tag.Location
					}

					if tag.IsNull {
						datetime := utils.DefaultTime.In(location).Format(time.DateOnly)
						if text != datetime {
							return nil, fmt.Errorf("invalid date: %q, valid: %q", text, datetime)
						}

						data = utils.DefaultTime
					} else {
						var d time.Time
						d, err = time.ParseInLocation(time.DateOnly, text, location)
						if err != nil {
							return nil, fmt.Errorf("invalid date %q: %w", text, err)
						}

						data, text, err = tag.ValidateDate(d, example || tag.Example)
					}

				case ir.ValueTypeTime:
					location := time.UTC
					if tag.Location != nil {
						location = tag.Location
					}

					if tag.IsNull {
						datetime := utils.DefaultTime.In(location).Format(time.TimeOnly)
						if text != datetime {
							return nil, fmt.Errorf("invalid time: %q, valid: %q", text, datetime)
						}

						data = utils.DefaultTime
					} else {
						var d time.Time
						d, err = time.ParseInLocation(time.TimeOnly, text, location)
						if err != nil {
							return nil, fmt.Errorf("invalid time %q: %w", text, err)
						}

						data, text, err = tag.ValidateTime(d, example || tag.Example)
					}

				case ir.ValueTypeUuid:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid uuid: %q, valid: %q", text, "")
						}

						data = [16]byte{}
					} else {
						data, text, err = tag.ValidateUuid(text, example || tag.Example)
					}

				case ir.ValueTypeDecimal:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid decimal: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateDecimal(text, example || tag.Example)
					}

				case ir.ValueTypeIp:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid ip: %q, valid: %q", text, "")
						}

						data = net.IP{}
					} else {
						ip := net.ParseIP(text)

						data, text, err = tag.ValidateIp(ip, example || tag.Example)
					}

				case ir.ValueTypeUrl:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid url: %q, valid: %q", text, "")
						}

						data = url.URL{}
					} else {
						var u *url.URL
						u, err = url.Parse(text)
						if err != nil {
							return nil, fmt.Errorf("invalid url %q: %w", text, err)
						}

						data, text, err = tag.ValidateUrl(*u, example || tag.Example)
					}

				case ir.ValueTypeEmail:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid email: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateEmail(text, example || tag.Example)
					}

				case ir.ValueTypeEnum:
					if tag.Enums == "" {
						err = fmt.Errorf("enum empty")
						return
					}

					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid enums: %q, valid: %q", text, "")
						}

						data = -1
					} else {
						data, text, err = tag.ValidateEnum(text, example || tag.Example)
					}

				case ir.ValueTypeMedia:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid media: %q, valid: %q", text, "")
						}

						data = []byte{}
					} else {
						var val []byte
						val, err = base64.StdEncoding.DecodeString(text)
						if err != nil {
							return nil, fmt.Errorf("invalid base64 media %q: %w", text, err)
						}

						data, text, err = tag.ValidateMedia(val, example || tag.Example)
					}

				default:
					return nil, fmt.Errorf("unsupported type %v for string literal", tag.Type)
				}
			}

			if err != nil {
				return nil, err
			}

			return &ir.Value{
				Data: data,
				Text: text,
				Tag:  tag,
				Path: path,
			}, nil

		case token.Number:
			text := tok.Literal

			if strings.Contains(text, ".") {
				if tag.Type == ir.ValueTypeUnknown {
					tag.Type = ir.ValueTypeF64
				}

				switch tag.Type {
				case ir.ValueTypeF32:
					if tag.IsNull {
						if text != "0.0" {
							return nil, fmt.Errorf("invalid float32: %v, valid: %v", text, "0.0")
						}

						data = float32(0.0)
					} else {
						var f64 float64
						if f64, err = strconv.ParseFloat(text, 32); err != nil {
							return nil, fmt.Errorf("invalid float32 %q: %w", text, err)
						}

						data, text, err = tag.ValidateF32(float32(f64), example || tag.Example)
					}

				case ir.ValueTypeF64:
					if tag.IsNull {
						if text != "0.0" {
							return nil, fmt.Errorf("invalid float64: %v, valid: %v", text, "0.0")
						}

						data = 0.0
					} else {
						var f64 float64
						if f64, err = strconv.ParseFloat(text, 64); err != nil {
							return nil, fmt.Errorf("invalid float64 %q: %w", text, err)
						}

						data, text, err = tag.ValidateF64(f64, example || tag.Example)
					}

				case ir.ValueTypeDecimal:
					if tag.IsNull {
						if text != "0.0" {
							return nil, fmt.Errorf("invalid decimal: %v, valid: %v", text, "0.0")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateDecimal(text, example || tag.Example)
					}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v for float literal", tag.Type)
				}
			} else if strings.HasPrefix(text, "-") {
				if tag.Type == ir.ValueTypeUnknown {
					tag.Type = ir.ValueTypeI
				}

				switch tag.Type {
				case ir.ValueTypeI:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int: %v, valid: %v", text, "0")
						}

						data = 0
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid int %q: %w", text, err)
						}

						data, text, err = tag.ValidateI(int(uv), example || tag.Example)
					}

				case ir.ValueTypeI8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int8: %v, valid: %v", text, "0")
						}

						data = int8(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid int8 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI8(int8(uv), example || tag.Example)
					}

				case ir.ValueTypeI16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int16: %v, valid: %v", text, "0")
						}

						data = int16(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid int16 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI16(int16(uv), example || tag.Example)
					}

				case ir.ValueTypeI32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int32: %v, valid: %v", text, "0")
						}

						data = int32(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid int32 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI32(int32(uv), example || tag.Example)
					}

				case ir.ValueTypeI64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int64: %v, valid: %v", text, "0")
						}

						data = int64(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid int64 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI64(uv, example || tag.Example)
					}

				case ir.ValueTypeBigint:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid bigint: %v, valid: %v", text, "0")
						}

						data = big.Int{}
					} else {
						bi, ok := new(big.Int).SetString(text, 10)
						if !ok {
							return nil, fmt.Errorf("invalid bigint %q", text)
						}

						data, text, err = tag.ValidateBigint(*bi, example || tag.Example)
					}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v for negative literal", tag.Type)
				}
			} else {
				if tag.Type == ir.ValueTypeUnknown {
					tag.Type = ir.ValueTypeI
				}

				switch tag.Type {
				case ir.ValueTypeI:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int: %v, valid: %v", text, "0")
						}

						data = 0
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid int %q: %w", text, err)
						}

						data, text, err = tag.ValidateI(int(uv), example || tag.Example)
					}

				case ir.ValueTypeI8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int8: %v, valid: %v", text, "0")
						}

						data = int8(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid int8 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI8(int8(uv), example || tag.Example)
					}

				case ir.ValueTypeI16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int16: %v, valid: %v", text, "0")
						}

						data = int16(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid int16 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI16(int16(uv), example || tag.Example)
					}

				case ir.ValueTypeI32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int32: %v, valid: %v", text, "0")
						}

						data = int32(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid int32 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI32(int32(uv), example || tag.Example)
					}

				case ir.ValueTypeI64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int64: %v, valid: %v", text, "0")
						}

						data = int64(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid int64 %q: %w", text, err)
						}

						data, text, err = tag.ValidateI64(uv, example || tag.Example)
					}

				case ir.ValueTypeU:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint: %v, valid: %v", text, "0")
						}

						data = uint(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid uint %q: %w", text, err)
						}

						data, text, err = tag.ValidateU(uint(uv), example || tag.Example)
					}

				case ir.ValueTypeU8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint8: %v, valid: %v", text, "0")
						}

						data = uint8(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid uint8 %q: %w", text, err)
						}

						data, text, err = tag.ValidateU8(uint8(uv), example || tag.Example)
					}

				case ir.ValueTypeU16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint16: %v, valid: %v", text, "0")
						}

						data = uint16(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid uint16 %q: %w", text, err)
						}

						data, text, err = tag.ValidateU16(uint16(uv), example || tag.Example)
					}

				case ir.ValueTypeU32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint32: %v, valid: %v", text, "0")
						}

						data = uint32(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid uint32 %q: %w", text, err)
						}

						data, text, err = tag.ValidateU32(uint32(uv), example || tag.Example)
					}

				case ir.ValueTypeU64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint64: %v, valid: %v", text, "0")
						}

						data = uint64(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid uint64 %q: %w", text, err)
						}

						data, text, err = tag.ValidateU64(uv, example || tag.Example)
					}

				case ir.ValueTypeBigint:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid bigint: %v, valid: %v", text, "0")
						}

						data = big.Int{}
					} else {
						bi, ok := new(big.Int).SetString(text, 10)
						if !ok {
							return nil, fmt.Errorf("invalid bigint %q", text)
						}

						data, text, err = tag.ValidateBigint(*bi, example || tag.Example)
					}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v", tag.Type)
				}
			}

			if err != nil {
				return nil, err
			}

			return &ir.Value{
				Data: data,
				Text: text,
				Tag:  tag,
				Path: path,
			}, nil

		case token.True:
			if tag.Type == ir.ValueTypeUnknown {
				tag.Type = ir.ValueTypeBool
			}

			switch tag.Type {
			case ir.ValueTypeBool:
				if tag.IsNull {
					return nil, fmt.Errorf("bool must false when bool is null")
				} else {
					if _, _, err = tag.ValidateBool(true, example || tag.Example); err != nil {
						return nil, err
					}
				}

			default:
				return nil, fmt.Errorf("unsupported type %v for boolean literal", tag.Type)
			}

			return &ir.Value{
				Data: true,
				Text: ir.True,
				Tag:  tag,
				Path: path,
			}, nil

		case token.False:
			if tag.Type == ir.ValueTypeUnknown {
				tag.Type = ir.ValueTypeBool
			}

			switch tag.Type {
			case ir.ValueTypeBool:
				if tag.IsNull {
				} else {
					if _, _, err = tag.ValidateBool(false, example || tag.Example); err != nil {
						return nil, err
					}
				}

			default:
				return nil, fmt.Errorf("unsupported type %v for boolean literal", tag.Type)
			}

			return &ir.Value{
				Data: false,
				Text: ir.False,
				Tag:  tag,
				Path: path,
			}, nil

		case token.Null:
			return nil, fmt.Errorf("null is not supported")

		default:
			return nil, fmt.Errorf("unexpected token %s", tok.Type)
		}
	}
}

func (p *Parser) parseObject(openLine int, path string, tag *ir.Tag) (*ir.Object, error) {
	p.depth++
	if p.depth > maxDepth {
		return nil, fmt.Errorf("max depth: %d", maxDepth)
	}

	if tag == nil {
		tag = ir.NewTag()
	}

	var err error

	if tag.Type == ir.ValueTypeUnknown {
		tag.Type = ir.ValueTypeObj
	}

	if tag.Name != "" {
		if path == "" {
			path = tag.Name
		} else {
			path = fmt.Sprintf("%s.%s", path, tag.Name)
		}
	}

	obj := &ir.Object{
		Tag:  tag,
		Path: path,
	}

	var val ir.Node
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			break
		}
		if tok.Type == token.RBrace {
			p.next()
			break
		}

		if tok.Type == token.Comment {
			if len(p.pending) > 0 {
				last := p.pending[len(p.pending)-1]
				if tok.Line-last.Line > 1 {
					p.pending = nil
				}
			}
			p.pending = append(p.pending, tok)
			p.next()
			continue
		}

		key := p.next()
		if key.Type != token.String {
			return nil, fmt.Errorf("expect string key")
		}
		keyStr := utils.CamelToSnake(key.Literal)

		p.next()

		var childTag *ir.Tag
		if childTag, err = p.consumeCommentsFor(tok.Line); err != nil {
			return nil, err
		}

		if childTag == nil {
			childTag = ir.NewTag()
		}

		pa := fmt.Sprintf("%s.%s", path, keyStr)
		if ir.ValueTypeMap == tag.Type {
			childTag.Inherit(tag)

			if childTag.Example {
				tag.IsEmpty = true
			}

			pa = fmt.Sprintf("%s[%s]", path, keyStr)
		}
		example := tag.Example || childTag.Example
		if val, err = p.parse(pa, example, childTag); err != nil {
			err = fmt.Errorf("%s parse err: %w", pa, err)
			return nil, err
		}
		if val == nil {
			continue
		}

		field := &ir.Field{
			Key:   keyStr,
			Value: val,
		}
		obj.Fields = append(obj.Fields, field)

		if p.peek().Type == token.Comma {
			p.next()
		}
	}

	if !tag.Example {
		switch tag.Type {
		case ir.ValueTypeMap:
			err = tag.ValidateMap()

		case ir.ValueTypeObj:
			err = tag.ValidateObj()
		}

		if err != nil {
			err = fmt.Errorf("validate failed: %w", err)
			return nil, err
		}
	}

	return obj, nil
}

func (p *Parser) parseArray(openLine int, path string, tag *ir.Tag) (*ir.Array, error) {
	p.depth++
	if p.depth > maxDepth {
		return nil, fmt.Errorf("max depth: %d", maxDepth)
	}

	if tag == nil {
		tag = ir.NewTag()
	}

	var err error

	if tag.Type == ir.ValueTypeUnknown {
		tag.Type = ir.ValueTypeVec
	}

	if tag.Name != "" {
		path = fmt.Sprintf("%s.%s", path, tag.Name)
	}
	arr := &ir.Array{
		Tag:  tag,
		Path: path,
	}

	var item ir.Node
	var i int
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			break
		}
		if tok.Type == token.RBracket {
			p.next()
			break
		}

		if tok.Type == token.Comment {
			if len(p.pending) > 0 {
				last := p.pending[len(p.pending)-1]
				if tok.Line-last.Line > 1 {
					p.pending = nil
				}
			}
			p.pending = append(p.pending, tok)
			p.next()
			continue
		}

		var childTag *ir.Tag
		if openLine != tok.Line {
			if childTag, err = p.consumeCommentsFor(tok.Line); err != nil {
				return nil, err
			}
		}

		if childTag == nil {
			childTag = ir.NewTag()
		}
		childTag.Inherit(tag)

		if childTag.Example {
			tag.IsEmpty = true
		}

		pa := fmt.Sprintf("%s[%d]", path, i)
		example := tag.Example || childTag.Example
		if item, err = p.parse(pa, example, childTag); err != nil {
			err = fmt.Errorf("%s parse err: %w", pa, err)
			return nil, err
		}
		if item == nil {
			continue
		}
		arr.Items = append(arr.Items, item)
		i++

		if p.peek().Type == token.Comma {
			p.next()
		}
	}

	if !tag.Example {
		switch tag.Type {
		case ir.ValueTypeArr:
			err = tag.ValidateArr(arr.Items)

		case ir.ValueTypeVec:
			err = tag.ValidateVec(arr.Items)
		}

		if err != nil {
			err = fmt.Errorf("validate failed: %w", err)
			return nil, err
		}
	}

	return arr, nil
}

func parseCommentsToTag(cs string) (*ir.Tag, error) {
	if after, ok := strings.CutPrefix(cs, "mm:"); ok {
		parsed, err := ir.ParseMMTag(after)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}
	return nil, nil
}
