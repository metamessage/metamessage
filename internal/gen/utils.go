package gen

import (
	"strings"
	"unicode"

	"github.com/metamessage/metamessage/internal/ir"
)

const indentUnit = "\t"

func WriteIndent(b *strings.Builder, indent int) {
	b.WriteString(strings.Repeat(indentUnit, indent))
}

func ExportName(s string) string {
	if s == "" {
		return "Field"
	}

	var cleaned []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, r)
		} else {
			cleaned = append(cleaned, ' ')
		}
	}

	parts := strings.Fields(string(cleaned))
	if len(parts) == 0 {
		return "Field"
	}

	var sb strings.Builder
	for _, p := range parts {
		sb.WriteRune(unicode.ToUpper(rune(p[0])))
		sb.WriteString(p[1:])
	}
	name := sb.String()

	for len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		name = name[1:]
	}
	if name == "" {
		return "Field"
	}
	return name
}

func toTitle(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func findFirstObjectInArray(a *ir.Array) *ir.Object {
	if a == nil {
		return nil
	}
	for _, item := range a.Items {
		if obj, ok := item.(*ir.Object); ok {
			return obj
		}
	}
	return nil
}
