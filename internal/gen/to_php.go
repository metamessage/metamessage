package gen

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/metamessage/metamessage/internal/ir"
)

var phpTypeMap = map[ir.ValueType]string{
	ir.ValueTypeUnknown:  "mixed",
	ir.ValueTypeString:   "string",
	ir.ValueTypeBytes:    "string",
	ir.ValueTypeBool:     "bool",
	ir.ValueTypeArray:    "array",
	ir.ValueTypeSlice:    "array",
	ir.ValueTypeMap:      "array",
	ir.ValueTypeInt:      "int",
	ir.ValueTypeInt8:     "int",
	ir.ValueTypeInt16:    "int",
	ir.ValueTypeInt32:    "int",
	ir.ValueTypeInt64:    "int",
	ir.ValueTypeUint:     "int",
	ir.ValueTypeUint8:    "int",
	ir.ValueTypeUint16:   "int",
	ir.ValueTypeUint32:   "int",
	ir.ValueTypeUint64:   "int",
	ir.ValueTypeFloat32:  "float",
	ir.ValueTypeFloat64:  "float",
	ir.ValueTypeBigInt:   "string",
	ir.ValueTypeDateTime: "string",
	ir.ValueTypeDate:     "string",
	ir.ValueTypeTime:     "string",
	ir.ValueTypeUUID:     "string",
	ir.ValueTypeDecimal:  "string",
	ir.ValueTypeEmail:    "string",
	ir.ValueTypeIP:       "string",
	ir.ValueTypeURL:      "string",
	ir.ValueTypeEnum:     "string",
	ir.ValueTypeImage:    "string",
}

func ToPHP(n ir.Node) string {
	if n == nil {
		return ""
	}

	topName := "Obj"
	if obj, ok := n.(*ir.Object); ok && obj.Tag != nil && obj.Tag.Name != "" {
		topName = exportPhpClassName(obj.Tag.Name)
	}

	var sb strings.Builder
	sb.WriteString("<?php\n")
	sb.WriteString("declare(strict_types=1);\n\n")

	classNames := make(map[string]struct{})
	genPHPClass(&sb, topName, n, classNames)
	genPHPNestedClasses(&sb, n, classNames)

	sb.WriteString("$" + exportPhpInstanceName(topName) + " = new " + topName + "();\n")
	genPHPObjectAssignments(&sb, "$"+exportPhpInstanceName(topName), n, 0)

	return sb.String()
}

func genPHPClass(b *strings.Builder, className string, n ir.Node, generated map[string]struct{}) {
	if _, ok := generated[className]; ok {
		return
	}
	generated[className] = struct{}{}

	b.WriteString("class ")
	b.WriteString(className)
	b.WriteString("\n")
	b.WriteString("{\n")
	genPHPFields(b, n, 1)
	b.WriteString("}\n\n")
}

func genPHPFields(b *strings.Builder, n ir.Node, indent int) {
	obj, ok := n.(*ir.Object)
	if !ok {
		return
	}

	for _, f := range obj.Fields {
		if f == nil {
			continue
		}
		WriteIndent(b, indent)
		typeName := getPhpTypeForField(f)
		b.WriteString("public ")

		// Check if field is nullable
		isNullable := false
		if tag := f.Value.GetTag(); tag != nil {
			isNullable = tag.Nullable
		}

		if isNullable {
			b.WriteString("?")
		}
		b.WriteString(typeName)
		b.WriteString(" $")
		b.WriteString(exportPhpFieldName(f.Key))
		if isNullable {
			b.WriteString(" = null;\n")
		} else {
			b.WriteString(";\n")
		}
	}
}

func getPhpTypeForField(f *ir.Field) string {
	switch v := f.Value.(type) {
	case *ir.Value:
		return getPhpType(v)
	case *ir.Object:
		return getPhpObjectType(f.Key, v)
	case *ir.Array:
		return getPhpArrayType(f.Key, v)
	default:
		return "mixed"
	}
}

func getPhpType(v *ir.Value) string {
	if v != nil && v.Tag != nil {
		if t, ok := phpTypeMap[v.Tag.Type]; ok {
			return t
		}
	}
	return "mixed"
}

func getPhpObjectType(fieldKey string, obj *ir.Object) string {
	if obj != nil && obj.Tag != nil && obj.Tag.Name != "" {
		return exportPhpClassName(obj.Tag.Name)
	}
	return exportPhpClassName(fieldKey)
}

func getPhpArrayType(fieldKey string, a *ir.Array) string {
	if a == nil {
		return "array"
	}
	if a.Tag != nil && a.Tag.ChildType != ir.ValueTypeUnknown {
		if _, ok := phpTypeMap[a.Tag.ChildType]; ok {
			return "array"
		}
	}
	if len(a.Items) > 0 {
		switch item := a.Items[0].(type) {
		case *ir.Object:
			return "array"
		case *ir.Value:
			if item.Tag != nil {
				if _, ok := phpTypeMap[item.Tag.Type]; ok {
					return "array"
				}
			}
		}
	}
	return "array"
}

func genPHPNestedClasses(b *strings.Builder, n ir.Node, generated map[string]struct{}) {
	obj, ok := n.(*ir.Object)
	if !ok {
		return
	}

	for _, f := range obj.Fields {
		if f == nil {
			continue
		}

		switch v := f.Value.(type) {
		case *ir.Object:
			className := getPhpObjectType(f.Key, v)
			genPHPClass(b, className, v, generated)
			genPHPNestedClasses(b, v, generated)
		case *ir.Array:
			if nestedObj := findFirstObjectInArrayPHP(v); nestedObj != nil {
				className := getPhpObjectType(f.Key, nestedObj)
				genPHPClass(b, className, nestedObj, generated)
				genPHPNestedClasses(b, nestedObj, generated)
			}
		}
	}
}

func findFirstObjectInArrayPHP(a *ir.Array) *ir.Object {
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

func exportPhpClassName(s string) string {
	name := ExportName(s)
	if len(name) == 0 {
		return "Obj"
	}
	return name
}

func exportPhpFieldName(s string) string {
	if s == "" {
		return "field"
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	if len(parts) == 0 {
		return "field"
	}
	name := parts[0]
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		name += strings.Title(part)
	}
	if unicode.IsDigit(rune(name[0])) {
		name = "f_" + name
	}
	return name
}

func exportPhpInstanceName(name string) string {
	if name == "" {
		return "objData"
	}
	if len(name) == 1 {
		return strings.ToLower(name) + "Data"
	}
	return strings.ToLower(string(name[0])) + name[1:] + "Data"
}

func genPHPObjectAssignments(b *strings.Builder, varName string, n ir.Node, indent int) {
	obj, ok := n.(*ir.Object)
	if !ok {
		return
	}

	for _, f := range obj.Fields {
		if f == nil {
			continue
		}
		prop := varName + "->" + exportPhpFieldName(f.Key)
		switch v := f.Value.(type) {
		case *ir.Value:
			WriteIndent(b, indent)
			b.WriteString(prop + " = " + formatPhpValueLiteral(v) + ";\n")
		case *ir.Object:
			className := getPhpObjectType(f.Key, v)
			WriteIndent(b, indent)
			b.WriteString(prop + " = new " + className + "();\n")
			genPHPObjectAssignments(b, prop, v, indent)
		case *ir.Array:
			if nestedObj := findFirstObjectInArrayPHP(v); nestedObj != nil {
				WriteIndent(b, indent)
				b.WriteString(prop + " = [\n")
				for i, item := range v.Items {
					WriteIndent(b, indent+1)
					switch iv := item.(type) {
					case *ir.Object:
						b.WriteString(genPHPObjectLiteral(iv, f.Key, indent+1))
					default:
						b.WriteString(formatPhpValueLiteral(iv.(*ir.Value)))
					}
					if i < len(v.Items)-1 {
						b.WriteString(",\n")
					}
				}
				b.WriteString("\n")
				WriteIndent(b, indent)
				b.WriteString("];\n")
			} else {
				WriteIndent(b, indent)
				b.WriteString(prop + " = " + genPHPArrayLiteral(v) + ";\n")
			}
		default:
			WriteIndent(b, indent)
			b.WriteString(prop + " = null;\n")
		}
	}
}

func genPHPArrayLiteral(a *ir.Array) string {
	if a == nil {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[")
	for i, item := range a.Items {
		if i > 0 {
			sb.WriteString(", ")
		}
		switch iv := item.(type) {
		case *ir.Value:
			sb.WriteString(formatPhpValueLiteral(iv))
		case *ir.Object:
			sb.WriteString(genPHPObjectLiteral(iv, "", 0))
		default:
			sb.WriteString("null")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func genPHPObjectLiteral(obj *ir.Object, fieldKey string, indent int) string {
	var sb strings.Builder
	className := getPhpObjectType(fieldKey, obj)
	sb.WriteString("(function() {\n")
	WriteIndent(&sb, indent+1)
	sb.WriteString("$obj = new " + className + "();\n")
	for _, f := range obj.Fields {
		if f == nil {
			continue
		}
		WriteIndent(&sb, indent+1)
		sb.WriteString("$obj->" + exportPhpFieldName(f.Key) + " = ")
		switch iv := f.Value.(type) {
		case *ir.Value:
			sb.WriteString(formatPhpValueLiteral(iv))
		case *ir.Object:
			sb.WriteString(genPHPObjectLiteral(iv, f.Key, indent+1))
		case *ir.Array:
			sb.WriteString(genPHPArrayLiteral(iv))
		default:
			sb.WriteString("null")
		}
		sb.WriteString(";\n")
	}
	WriteIndent(&sb, indent+1)
	sb.WriteString("return $obj;\n")
	WriteIndent(&sb, indent)
	sb.WriteString("})()")
	return sb.String()
}

func formatPhpValueLiteral(v *ir.Value) string {
	if v == nil {
		return "null"
	}

	if v.Tag != nil {
		switch v.Tag.Type {
		case ir.ValueTypeString, ir.ValueTypeBytes, ir.ValueTypeDateTime, ir.ValueTypeDate, ir.ValueTypeTime, ir.ValueTypeUUID, ir.ValueTypeDecimal, ir.ValueTypeEmail, ir.ValueTypeIP, ir.ValueTypeURL, ir.ValueTypeEnum, ir.ValueTypeImage:
			return "\"" + strings.ReplaceAll(v.Text, "\"", "\\\"") + "\""
		case ir.ValueTypeBool:
			if strings.EqualFold(v.Text, "true") {
				return "true"
			}
			return "false"
		case ir.ValueTypeInt, ir.ValueTypeInt8, ir.ValueTypeInt16, ir.ValueTypeInt32, ir.ValueTypeInt64, ir.ValueTypeUint, ir.ValueTypeUint8, ir.ValueTypeUint16, ir.ValueTypeUint32, ir.ValueTypeUint64, ir.ValueTypeFloat32, ir.ValueTypeFloat64:
			if v.Text == "" {
				return "0"
			}
			return v.Text
		}
	}
	return "\"" + strings.ReplaceAll(v.Text, "\"", "\\\"") + "\""
}

func PrintPHPStruct(n ir.Node) {
	fmt.Println(ToPHP(n))
}
