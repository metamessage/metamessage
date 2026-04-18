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

	"github.com/lizongying/meta-message/internal/jsonc/ast"
	"github.com/lizongying/meta-message/internal/jsonc/token"
	"github.com/lizongying/meta-message/internal/utils"
)

const BitSize = 32 << (^uint(0) >> 63)

type Parser struct {
	toks []token.Token
	pos  int

	pending []token.Token
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

func (p *Parser) consumeCommentsFor(anchorLine int) (*ast.Tag, error) {
	if len(p.pending) == 0 {
		return nil, nil
	}

	// if the last pending comment is separated from the anchor by a blank line,
	// drop all pending comments
	last := p.pending[len(p.pending)-1]
	if anchorLine-last.Line > 1 {
		p.pending = nil
		return nil, nil
	}

	var out *ast.Tag
	for _, ct := range p.pending {
		parsed, err := parseCommentsToTag(ct.Literal)
		if err != nil {
			return nil, err
		}
		if parsed == nil {
			continue
		}
		out = ast.MergeTag(out, parsed)
	}

	return out, nil
}

func (p *Parser) Parse() (val ast.Node, err error) {
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			return
		}

		if tok.Type == token.LeadingComment {
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

		if tok.Type == token.TrailingComment {
			if val != nil {
				var parsed *ast.Tag
				parsed, err = parseCommentsToTag(tok.Literal)
				if err != nil {
					return
				} else if parsed != nil {
					mergeNodeTag(val, parsed)
				}
			}
			p.next()
			continue
		}

		if val, err = p.parse(); err != nil {
			return
		}
	}
}

func (p *Parser) parse() (val ast.Node, err error) {
	for {
		tok := p.next()
		var data any
		switch tok.Type {
		case token.EOF:
			return nil, nil
		case token.LBrace:
			return p.parseObject(tok.Line)
		case token.LBracket:
			return p.parseArray(tok.Line)
		case token.String:
			var tag *ast.Tag
			if tag, err = p.consumeCommentsFor(tok.Line); err != nil {
				return nil, err
			}
			if tag == nil {
				tag = ast.NewTag()
			}
			text := tok.Literal
			t := tag.Type

			if t == ast.ValueTypeUnknown {
				tag.Type = ast.ValueTypeString
			}

			switch text {
			case ast.SimpleCodeStr,
				ast.SimpleMessageStr,
				ast.SimpleDataStr,
				ast.SimpleSuccessStr,
				ast.SimpleErrorStr,
				ast.SimpleUnknownStr,
				ast.SimplePageStr,
				ast.SimpleLimitStr,
				ast.SimpleOffsetStr,
				ast.SimpleTotalStr,
				ast.SimpleIdStr,
				ast.SimpleNameStr,
				ast.SimpleDescriptionStr,
				ast.SimpleTypeStr,
				ast.SimpleVersionStr,
				ast.SimpleStatusStr,
				ast.SimpleUrlStr,
				ast.SimpleCreateTimeStr,
				ast.SimpleUpdateTimeStr,
				ast.SimpleDeleteTimeStr,
				ast.SimpleAccountStr,
				ast.SimpleTokenStr,
				ast.SimpleExpireTimeStr,
				ast.SimpleKeyStr,
				ast.SimpleValStr:
				tag.Type = ast.ValueTypeString

			default:
				switch t {
				case ast.ValueTypeString:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid string: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateString(text)
					}

				case ast.ValueTypeBytes:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid bytes: %q, valid: %q", text, "")
						}

						data = []byte{}
					} else {
						var d []byte
						if d, err = base64.StdEncoding.DecodeString(text); err != nil {
							return nil, fmt.Errorf("invalid base64 bytes literal: %w", err)
						}

						data, text, err = tag.ValidateBytes(d)
					}

				case ast.ValueTypeDateTime:
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
							return nil, fmt.Errorf("invalid datatime literal: %w", err)
						}

						data, text, err = tag.ValidateDateTime(d)
					}

				case ast.ValueTypeDate:
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
							return nil, fmt.Errorf("invalid data literal: %w", err)
						}

						data, text, err = tag.ValidateDate(d)
					}

				case ast.ValueTypeTime:
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
							return nil, fmt.Errorf("invalid time literal: %w", err)
						}

						data, text, err = tag.ValidateTime(d)
					}

				case ast.ValueTypeUUID:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid uuid: %q, valid: %q", text, "")
						}

						data = [16]byte{}
					} else {
						data, text, err = tag.ValidateUUID(text)
					}

				case ast.ValueTypeDecimal:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid decimal: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateDecimal(text)
					}

				case ast.ValueTypeIP:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid email: %q, valid: %q", text, "")
						}

						data = net.IP{}
					} else {
						ip := net.ParseIP(text)

						data, text, err = tag.ValidateIP(ip)
					}

				case ast.ValueTypeURL:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid url: %q, valid: %q", text, "")
						}

						data = url.URL{}
					} else {
						var u *url.URL
						u, err = url.Parse(text)
						if err != nil {
							return nil, fmt.Errorf("invalid url: %w", err)
						}

						data, text, err = tag.ValidateURL(*u)
					}

				case ast.ValueTypeEmail:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid email: %q, valid: %q", text, "")
						}

						data = ""
					} else {
						data, text, err = tag.ValidateEmail(text)
					}

				case ast.ValueTypeEnum:
					if tag.Enum == "" {
						err = fmt.Errorf("enum empty")
						return
					}

					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid enum: %q, valid: %q", text, "")
						}

						data = -1
					} else {
						data, text, err = tag.ValidateEnum(text)
					}

				case ast.ValueTypeImage:
					if tag.IsNull {
						if text != "" {
							return nil, fmt.Errorf("invalid image: %q, valid: %q", text, "")
						}

						data = []byte{}
					} else {
						var val []byte
						val, err = base64.StdEncoding.DecodeString(text)
						if err != nil {
							return nil, fmt.Errorf("invalid base64 image literal: %w", err)
						}

						data, text, err = tag.ValidateImage(val)
					}

				default:
					return nil, fmt.Errorf("unsupported type %v for string literal", t)
				}
			}

			if err != nil {
				return nil, err
			}

			return &ast.Value{
				Data: data,
				Text: text,
				Tag:  tag,
			}, nil

		case token.Number:
			var tag *ast.Tag
			if tag, err = p.consumeCommentsFor(tok.Line); err != nil {
				return nil, err
			}
			if tag == nil {
				tag = ast.NewTag()
			}
			text := tok.Literal
			t := tag.Type

			if strings.Contains(text, ".") {
				if t == ast.ValueTypeUnknown {
					tag.Type = ast.ValueTypeFloat64
				}

				switch t {
				case ast.ValueTypeFloat32:
					if tag.IsNull {
						if text != "0.0" {
							return nil, fmt.Errorf("invalid float32: %v, valid: %v", text, "0.0")
						}

						data = float32(0.0)
					} else {
						var f64 float64
						if f64, err = strconv.ParseFloat(text, 32); err != nil {
							return nil, fmt.Errorf("invalid float32 literal: %w", err)
						}

						data, text, err = tag.ValidateFloat32(float32(f64))
					}

				case ast.ValueTypeFloat64:
					if tag.IsNull {
						if text != "0.0" {
							return nil, fmt.Errorf("invalid float64: %v, valid: %v", text, "0.0")
						}

						data = 0.0
					} else {
						var f64 float64
						if f64, err = strconv.ParseFloat(text, 64); err != nil {
							return nil, fmt.Errorf("invalid float64 literal: %w", err)
						}

						data, text, err = tag.ValidateFloat64(f64)
					}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v for float literal", t)
				}
			} else if strings.HasPrefix(text, "-") {
				if t == ast.ValueTypeUnknown {
					tag.Type = ast.ValueTypeInt
				}

				switch t {
				case ast.ValueTypeInt:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int: %v, valid: %v", text, "0")
						}

						data = 0
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid uint literal: %w", err)
						}

						data, text, err = tag.ValidateInt(int(uv))
					}

				case ast.ValueTypeInt8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int8: %v, valid: %v", text, "0")
						}

						data = int8(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid int8 literal: %w", err)
						}

						data, text, err = tag.ValidateInt8(int8(uv))
					}

				case ast.ValueTypeInt16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int16: %v, valid: %v", text, "0")
						}

						data = int16(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid int16 literal: %w", err)
						}

						data, text, err = tag.ValidateInt16(int16(uv))
					}

				case ast.ValueTypeInt32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int32: %v, valid: %v", text, "0")
						}

						data = int32(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid int32 literal: %w", err)
						}

						data, text, err = tag.ValidateInt32(int32(uv))
					}

				case ast.ValueTypeInt64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int64: %v, valid: %v", text, "0")
						}

						data = int64(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid int64 literal: %w", err)
						}

						data, text, err = tag.ValidateInt64(uv)
					}

				case ast.ValueTypeBigInt:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid bigint: %v, valid: %v", text, "0")
						}

						data = big.Int{}
					} else {
						bi, ok := new(big.Int).SetString(text, 10)
						if !ok {
							return nil, fmt.Errorf("invalid bigint literal: %v", text)
						}

						data, text, err = tag.ValidateBigInt(*bi)
					}

				// // It is recommended that floating - point numbers must include a decimal point, e.g., 1.0.
				// case ast.ValueTypeFloat32:
				// 	if tag.IsNull {
				// 		data = float32(0.0)
				// 		text = "0.0"
				// 	} else {
				// 		var f64 float64
				// 		if f64, err = strconv.ParseFloat(text, 32); err != nil {
				// 			return nil, fmt.Errorf("invalid float32 literal: %w", err)
				// 		}

				// 		data, text, err = tag.ValidateFloat32(float32(f64))
				// 	}

				// // It is recommended that floating - point numbers must include a decimal point, e.g., 1.0.
				// case ast.ValueTypeFloat64:
				// 	if tag.IsNull {
				// 		data = 0.0
				// 		text = "0.0"
				// 	} else {
				// 		var f64 float64
				// 		if f64, err = strconv.ParseFloat(text, 64); err != nil {
				// 			return nil, fmt.Errorf("invalid float64 literal: %w", err)
				// 		}

				// 		data, text, err = tag.ValidateFloat64(f64)
				// 	}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v for negative literal", t)
				}
			} else {
				if t == ast.ValueTypeUnknown {
					tag.Type = ast.ValueTypeInt
				}

				switch t {
				case ast.ValueTypeInt:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int: %v, valid: %v", text, "0")
						}

						data = 0
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid uint literal: %w", err)
						}

						data, text, err = tag.ValidateInt(int(uv))
					}

				case ast.ValueTypeInt8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int8: %v, valid: %v", text, "0")
						}

						data = int8(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid int8 literal: %w", err)
						}

						data, text, err = tag.ValidateInt8(int8(uv))
					}

				case ast.ValueTypeInt16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int16: %v, valid: %v", text, "0")
						}

						data = int16(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid int16 literal: %w", err)
						}

						data, text, err = tag.ValidateInt16(int16(uv))
					}

				case ast.ValueTypeInt32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int32: %v, valid: %v", text, "0")
						}

						data = int32(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid int32 literal: %w", err)
						}

						data, text, err = tag.ValidateInt32(int32(uv))
					}

				case ast.ValueTypeInt64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid int64: %v, valid: %v", text, "0")
						}

						data = int64(0)
					} else {
						var uv int64
						if uv, err = strconv.ParseInt(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid int64 literal: %w", err)
						}

						data, text, err = tag.ValidateInt64(uv)
					}

				case ast.ValueTypeUint:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint: %v, valid: %v", text, "0")
						}

						data = uint(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, BitSize); err != nil {
							return nil, fmt.Errorf("invalid uint literal: %w", err)
						}

						data, text, err = tag.ValidateUint(uint(uv))
					}

				case ast.ValueTypeUint8:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint8: %v, valid: %v", text, "0")
						}

						data = uint8(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 8); err != nil {
							return nil, fmt.Errorf("invalid uint8 literal: %w", err)
						}

						data, text, err = tag.ValidateUint8(uint8(uv))
					}

				case ast.ValueTypeUint16:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint16: %v, valid: %v", text, "0")
						}

						data = uint16(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 16); err != nil {
							return nil, fmt.Errorf("invalid uint16 literal: %w", err)
						}

						data, text, err = tag.ValidateUint16(uint16(uv))
					}

				case ast.ValueTypeUint32:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint32: %v, valid: %v", text, "0")
						}

						data = uint32(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 32); err != nil {
							return nil, fmt.Errorf("invalid uint32 literal: %w", err)
						}

						data, text, err = tag.ValidateUint32(uint32(uv))
					}

				case ast.ValueTypeUint64:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid uint64: %v, valid: %v", text, "0")
						}

						data = uint64(0)
					} else {
						var uv uint64
						if uv, err = strconv.ParseUint(text, 10, 64); err != nil {
							return nil, fmt.Errorf("invalid uint64 literal: %w", err)
						}

						data, text, err = tag.ValidateUint64(uv)
					}

				case ast.ValueTypeBigInt:
					if tag.IsNull {
						if text != "0" {
							return nil, fmt.Errorf("invalid bigint: %v, valid: %v", text, "0")
						}

						data = big.Int{}
					} else {
						bi, ok := new(big.Int).SetString(text, 10)
						if !ok {
							return nil, fmt.Errorf("invalid bigint literal: %v", text)
						}

						data, text, err = tag.ValidateBigInt(*bi)
					}

				// // It is recommended that floating - point numbers must include a decimal point, e.g., 1.0.
				// case ast.ValueTypeFloat32:
				// 	if tag.IsNull {
				// 		data = float32(0.0)
				// 		text = "0.0"
				// 	} else {
				// 		var f64 float64
				// 		if f64, err = strconv.ParseFloat(text, 32); err != nil {
				// 			return nil, fmt.Errorf("invalid float32 literal: %w", err)
				// 		}

				// 		data, text, err = tag.ValidateFloat32(float32(f64))
				// 	}

				// // It is recommended that floating - point numbers must include a decimal point, e.g., 1.0.
				// case ast.ValueTypeFloat64:
				// 	if tag.IsNull {
				// 		data = 0.0
				// 		text = "0.0"
				// 	} else {
				// 		var f64 float64
				// 		if f64, err = strconv.ParseFloat(text, 64); err != nil {
				// 			return nil, fmt.Errorf("invalid float64 literal: %w", err)
				// 		}

				// 		data, text, err = tag.ValidateFloat64(f64)
				// 	}

				default:
					return nil, fmt.Errorf("unsupported numeric type %v", t)
				}
			}

			if err != nil {
				return nil, err
			}

			return &ast.Value{
				Data: data,
				Text: text,
				Tag:  tag,
			}, nil

		case token.True:
			var tag *ast.Tag
			tag, err = p.consumeCommentsFor(tok.Line)
			if err != nil {
				return nil, err
			}

			if tag == nil {
				tag = ast.NewTag()
			}
			if tag.Type == ast.ValueTypeUnknown {
				tag.Type = ast.ValueTypeBool
			}

			switch tag.Type {
			case ast.ValueTypeBool:
				if tag.IsNull {
					return nil, fmt.Errorf("bool must false when bool is null")
				} else {
					if _, _, err = tag.ValidateBool(true); err != nil {
						return nil, err
					}
				}

			default:
				return nil, fmt.Errorf("unsupported type %v for boolean literal", tag.Type)
			}

			return &ast.Value{
				Data: true,
				Text: ast.True,
				Tag:  tag,
			}, nil

		case token.False:
			var tag *ast.Tag
			tag, err = p.consumeCommentsFor(tok.Line)
			if err != nil {
				return nil, err
			}

			if tag == nil {
				tag = ast.NewTag()
			}
			if tag.Type == ast.ValueTypeUnknown {
				tag.Type = ast.ValueTypeBool
			}

			switch tag.Type {
			case ast.ValueTypeBool:
				if tag.IsNull {
				} else {
					if _, _, err = tag.ValidateBool(false); err != nil {
						return nil, err
					}
				}

			default:
				return nil, fmt.Errorf("unsupported type %v for boolean literal", tag.Type)
			}

			return &ast.Value{
				Data: false,
				Text: ast.False,
				Tag:  tag,
			}, nil

		// case token.Null:
		// 	var tag *ast.Tag
		// 	if tag, err = p.consumeCommentsFor(tok.Line); err != nil {
		// 		return nil, err
		// 	}
		// 	if tag == nil {
		// 		tag = ast.NewTag()
		// 		tag.Type = ast.ValueTypeNull
		// 	} else {
		// 		switch tag.Type {
		// 		case ast.ValueTypeUnknown:
		// 			tag.Type = ast.ValueTypeNull
		// 		case ast.ValueTypeNull:
		// 		default:
		// 			return nil, fmt.Errorf("unsupported type %v for null literal", tag.Type)
		// 		}
		// 	}
		// 	return &ast.Value{
		// 		Data: nil,
		// 		Text: ast.Null,
		// 		Tag:  tag,
		// 	}, nil

		default:
			return nil, fmt.Errorf("unexpected token %s", tok.Type)
		}
	}
}

func (p *Parser) parseObject(openLine int) (*ast.Object, error) {
	tag, err := p.consumeCommentsFor(openLine)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		tag = ast.NewTag()
	}
	if tag.Type == ast.ValueTypeUnknown {
		tag.Type = ast.ValueTypeStruct
	}

	obj := &ast.Object{
		Tag: tag,
	}

	var val ast.Node
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			return obj, nil
		}
		if tok.Type == token.RBrace {
			p.next()
			return obj, nil
		}

		if tok.Type == token.LeadingComment {
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

		if tok.Type == token.TrailingComment {
			if val != nil {
				parsed, err := parseCommentsToTag(tok.Literal)
				if err != nil {
					return nil, err
				} else if parsed != nil {
					mergeNodeTag(val, parsed)
				}
			}
			p.next()
			continue
		}

		key := p.next()
		if key.Type != token.String {
			return nil, fmt.Errorf("expect string key")
		}

		p.next()
		if val, err = p.parse(); err != nil {
			return nil, err
		}
		if val == nil {
			p.next()
			continue
		}
		childTag := val.GetTag()
		if childTag != nil && tag != nil {
			childTag.ParentDesc = tag.ChildDesc
			childTag.ParentType = tag.ChildType
			childTag.ParentRaw = tag.ChildRaw
			childTag.ParentNullable = tag.ChildNullable
			childTag.ParentDefault = tag.ChildDefault
			childTag.ParentMin = tag.ChildMin
			childTag.ParentMax = tag.ChildMax
			childTag.ParentSize = tag.ChildSize
			childTag.ParentEnum = tag.ChildEnum
			childTag.ParentPattern = tag.ChildPattern
			childTag.ParentLocation = tag.ChildLocation
			childTag.ParentVersion = tag.ChildVersion
		}
		field := &ast.Field{
			Key:   utils.CamelToSnake(key.Literal),
			Value: val,
		}
		obj.Fields = append(obj.Fields, field)

		if p.peek().Type == token.Comma {
			p.next()
		}
	}
}

func (p *Parser) parseArray(openLine int) (*ast.Array, error) {
	tag, err := p.consumeCommentsFor(openLine)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		tag = ast.NewTag()
	}
	if tag.Type == ast.ValueTypeUnknown {
		if tag.Size > 0 {
			tag.Type = ast.ValueTypeArray
		} else {
			tag.Type = ast.ValueTypeSlice
		}
	}

	arr := &ast.Array{
		Tag: tag,
	}

	var item ast.Node
	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			return arr, nil
		}
		if tok.Type == token.RBracket {
			p.next()
			return arr, nil
		}

		if tok.Type == token.LeadingComment {
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

		if tok.Type == token.TrailingComment {
			if item != nil {
				parsed, err := parseCommentsToTag(tok.Literal)
				if err != nil {
					return nil, err
				} else if parsed != nil {
					mergeNodeTag(item, parsed)
				}
			}
			p.next()
			continue
		}

		if item, err = p.parse(); err != nil {
			return nil, err
		}
		if item == nil {
			continue
		}
		childTag := item.GetTag()
		if childTag != nil && tag != nil {
			childTag.ParentDesc = tag.ChildDesc
			childTag.ParentType = tag.ChildType
			childTag.ParentRaw = tag.ChildRaw
			childTag.ParentNullable = tag.ChildNullable
			childTag.ParentDefault = tag.ChildDefault
			childTag.ParentMin = tag.ChildMin
			childTag.ParentMax = tag.ChildMax
			childTag.ParentSize = tag.ChildSize
			childTag.ParentEnum = tag.ChildEnum
			childTag.ParentPattern = tag.ChildPattern
			childTag.ParentLocation = tag.ChildLocation
			childTag.ParentVersion = tag.ChildVersion
		}
		arr.Items = append(arr.Items, item)

		if p.peek().Type == token.Comma {
			p.next()
		}
	}
}

func mergeNodeTag(n ast.Node, parsed *ast.Tag) {
	if n == nil || parsed == nil {
		return
	}
	existing := n.GetTag()
	merged := ast.MergeTag(existing, parsed)
	switch t := n.(type) {
	case *ast.Value:
		t.Tag = merged
	case *ast.Object:
		t.Tag = merged
	case *ast.Array:
		t.Tag = merged
	default:
	}
}

func parseCommentsToTag(cs string) (*ast.Tag, error) {
	if after, ok := strings.CutPrefix(cs, "mm:"); ok {
		parsed, err := ast.ParseMMTag(after)
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}
	return nil, nil
}
