package jsonc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metamessage/metamessage/internal/ir"
)

const indentUnit = '\t'

type JSONCConfig struct {
	Pretty bool
	Indent byte
}

func writeIndent(b *strings.Builder, indent int, config *JSONCConfig) {
	if config.Pretty {
		for range indent {
			b.WriteByte(config.Indent)
		}
	}
}

func writeNullJSONC(b *strings.Builder, v *ir.NodeNull) {
	if v == nil {
		return
	}

	b.WriteString(ir.Null)
}

func writeValueJSONC(b *strings.Builder, v *ir.NodeScalar) {
	if v == nil {
		return
	}

	if v.Tag == nil {
		return
	}

	switch v.Tag.Type {
	case ir.ValueTypeStr,
		ir.ValueTypeBytes,
		ir.ValueTypeDatetime,
		ir.ValueTypeDate,
		ir.ValueTypeTime,
		ir.ValueTypeUuid,
		ir.ValueTypeIp,
		ir.ValueTypeUrl,
		ir.ValueTypeEmail,
		ir.ValueTypeEnum,
		ir.ValueTypeMedia:
		b.WriteString(strconv.Quote(v.Text))

	case ir.ValueTypeI, ir.ValueTypeI8, ir.ValueTypeI16, ir.ValueTypeI32, ir.ValueTypeI64,
		ir.ValueTypeU, ir.ValueTypeU8, ir.ValueTypeU16, ir.ValueTypeU32, ir.ValueTypeU64,
		ir.ValueTypeBigint,
		ir.ValueTypeDecimal,
		ir.ValueTypeBool:
		b.WriteString(v.Text)

	case ir.ValueTypeF32, ir.ValueTypeF64:
		b.WriteString(v.Text)

	default:
		b.WriteString(v.Text)
	}
}

func writeArrayJSONC(b *strings.Builder, a *ir.NodeArray, indent int, config *JSONCConfig) {
	b.WriteString("[\n")

	for _, item := range a.Items {
		writeLeadingComments(b, item.GetTag(), indent+1, config)

		writeIndent(b, indent+1, config)

		writeNodeJSONC(b, item, indent+1, config)

		b.WriteString(",\n")
	}

	writeIndent(b, indent, config)
	b.WriteString("]")
}

func writeObjectJSONC(b *strings.Builder, o *ir.NodeObject, indent int, config *JSONCConfig) {
	b.WriteString("{\n")

	for _, f := range o.Fields {
		writeLeadingComments(b, f.Value.GetTag(), indent+1, config)

		writeIndent(b, indent+1, config)

		b.WriteString(strconv.Quote(f.Key))
		b.WriteString(": ")

		writeNodeJSONC(b, f.Value, indent+1, config)

		b.WriteString(",\n")
	}

	writeIndent(b, indent, config)
	b.WriteString("}")
}

func writeLeadingComments(b *strings.Builder, tag *ir.Tag, indent int, config *JSONCConfig) {
	tagStr := tag.ToString()
	if tagStr != "" {
		b.WriteString("\n")
		writeIndent(b, indent, config)
		fmt.Fprintf(b, "// mm: %s\n", tagStr)
	}
}

func writeNodeJSONC(b *strings.Builder, n ir.Node, indent int, config *JSONCConfig) {
	switch v := n.(type) {
	case *ir.Doc:
		writeObjectJSONC(b, &ir.NodeObject{Fields: v.Fields, Tag: v.Tag, Path: v.Path}, indent, config)

	case *ir.NodeObject:
		writeObjectJSONC(b, v, indent, config)

	case *ir.NodeArray:
		writeArrayJSONC(b, v, indent, config)

	case *ir.NodeScalar:
		writeValueJSONC(b, v)

	case *ir.NodeNull:
		writeNullJSONC(b, v)

	default:
	}
}

func ToJSONC(n ir.Node) string {
	return ToJSONCWithConfig(n, &JSONCConfig{Pretty: true, Indent: indentUnit})
}

func ToJSONCCompact(n ir.Node) string {
	return ToJSONCWithConfig(n, &JSONCConfig{})
}

func ToJSONCWithConfig(n ir.Node, config *JSONCConfig) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	writeLeadingComments(&b, n.GetTag(), 0, config)
	writeNodeJSONC(&b, n, 0, config)
	return b.String()
}
