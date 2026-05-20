package jsonc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metamessage/metamessage/internal/ir"
)

const indentUnit = "\t"

func writeIndent(b *strings.Builder, indent int) {
	b.WriteString(strings.Repeat(indentUnit, indent))
}

func writeValueJSONC(b *strings.Builder, v *ir.Value) {
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
		ir.ValueTypeEnum:
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

func writeArrayJSONC(b *strings.Builder, a *ir.Array, indent int) {
	b.WriteString("[\n")

	for _, item := range a.Items {
		writeLeadingComments(b, item.GetTag(), indent+1)

		writeIndent(b, indent+1)

		writeNodeJSONC(b, item, indent+1)

		b.WriteString(",\n")
	}

	writeIndent(b, indent)
	b.WriteString("]")
}

func writeObjectJSONC(b *strings.Builder, o *ir.Object, indent int) {
	b.WriteString("{\n")

	for _, f := range o.Fields {
		writeLeadingComments(b, f.Value.GetTag(), indent+1)

		writeIndent(b, indent+1)

		b.WriteString(strconv.Quote(f.Key))
		b.WriteString(": ")

		writeNodeJSONC(b, f.Value, indent+1)

		b.WriteString(",\n")
	}

	writeIndent(b, indent)
	b.WriteString("}")
}

func writeLeadingComments(b *strings.Builder, tag *ir.Tag, indent int) {
	tagStr := tag.ToString()
	if tagStr != "" {
		b.WriteString("\n")
		writeIndent(b, indent)
		fmt.Fprintf(b, "// mm: %s\n", tagStr)
	}
}

func writeNodeJSONC(b *strings.Builder, n ir.Node, indent int) {
	switch v := n.(type) {
	case *ir.Value:
		writeValueJSONC(b, v)
	case *ir.Object:
		writeObjectJSONC(b, v, indent)
	case *ir.Doc:
		writeObjectJSONC(b, &ir.Object{Fields: v.Fields, Tag: v.Tag, Path: v.Path}, indent)
	case *ir.Array:
		writeArrayJSONC(b, v, indent)
	default:
	}
}

func ToJSONC(n ir.Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	writeLeadingComments(&b, n.GetTag(), 0)
	writeNodeJSONC(&b, n, 0)
	return b.String()
}
