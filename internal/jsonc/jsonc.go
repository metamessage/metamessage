package jsonc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metamessage/metamessage/internal/ast"
)

const indentUnit = "\t"

func writeIndent(b *strings.Builder, indent int) {
	b.WriteString(strings.Repeat(indentUnit, indent))
}

func writeValueJSONC(b *strings.Builder, v *ast.Value) {
	if v == nil {
		return
	}

	if v.Tag == nil {
		return
	}

	switch v.Tag.Type {
	case ast.ValueTypeString,
		ast.ValueTypeBytes,
		ast.ValueTypeDateTime,
		ast.ValueTypeDate,
		ast.ValueTypeTime,
		ast.ValueTypeUUID,
		ast.ValueTypeIP,
		ast.ValueTypeURL,
		ast.ValueTypeEmail,
		ast.ValueTypeEnum:
		b.WriteString(strconv.Quote(v.Text))

	case ast.ValueTypeInt, ast.ValueTypeInt8, ast.ValueTypeInt16, ast.ValueTypeInt32, ast.ValueTypeInt64,
		ast.ValueTypeUint, ast.ValueTypeUint8, ast.ValueTypeUint16, ast.ValueTypeUint32, ast.ValueTypeUint64,
		ast.ValueTypeBigInt,
		ast.ValueTypeDecimal,
		ast.ValueTypeBool:
		b.WriteString(v.Text)

	case ast.ValueTypeFloat32, ast.ValueTypeFloat64:
		b.WriteString(v.Text)

	default:
		b.WriteString(v.Text)
	}
}

func writeArrayJSONC(b *strings.Builder, a *ast.Array, indent int) {
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

func writeObjectJSONC(b *strings.Builder, o *ast.Object, indent int) {
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

func writeLeadingComments(b *strings.Builder, tag *ast.Tag, indent int) {
	tagStr := tag.String()
	if tagStr != "" {
		b.WriteString("\n")
		writeIndent(b, indent)
		fmt.Fprintf(b, "// mm: %s\n", tagStr)
	}
}

func writeNodeJSONC(b *strings.Builder, n ast.Node, indent int) {
	switch v := n.(type) {
	case *ast.Value:
		writeValueJSONC(b, v)
	case *ast.Object:
		writeObjectJSONC(b, v, indent)
	case *ast.Array:
		writeArrayJSONC(b, v, indent)
	default:
	}
}

func ToJSONC(n ast.Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	writeLeadingComments(&b, n.GetTag(), 0)
	writeNodeJSONC(&b, n, 0)
	return b.String()
}
