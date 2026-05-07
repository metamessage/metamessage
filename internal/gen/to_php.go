package gen

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/metamessage/metamessage/internal/ast"
)

var phpTypeMap = map[ast.ValueType]string{
	ast.ValueTypeUnknown:  "mixed",
	ast.ValueTypeString:   "string",
	ast.ValueTypeBytes:    "string",
	ast.ValueTypeBool:     "bool",
	ast.ValueTypeArray:    "array",
	ast.ValueTypeSlice:    "array",
	ast.ValueTypeMap:      "array",
	ast.ValueTypeInt:      "int",
	ast.ValueTypeInt8:     "int",
	ast.ValueTypeInt16:    "int",
	ast.ValueTypeInt32:    "int",
	ast.ValueTypeInt64:    "int",
	ast.ValueTypeUint:     "int",
	ast.ValueTypeUint8:    "int",
	ast.ValueTypeUint16:   "int",
	ast.ValueTypeUint32:   "int",
	ast.ValueTypeUint64:   "int",
	ast.ValueTypeFloat32:  "float",
	ast.ValueTypeFloat64:  "float",
	ast.ValueTypeBigInt:   "string",
	ast.ValueTypeDateTime: "string",
	ast.ValueTypeDate:     "string",
	ast.ValueTypeTime:     "string",
	ast.ValueTypeUUID:     "string",
	ast.ValueTypeDecimal:  "string",
	ast.ValueTypeEmail:    "string",
	ast.ValueTypeIP:       "string",
	ast.ValueTypeURL:      "string",
	ast.ValueTypeEnum:     "string",
	ast.ValueTypeImage:    "string",
}

func ToPHP(n ast.Node) string {
	if n == nil {
		return ""
	}

	topName := "Obj"
	if obj, ok := n.(*ast.Object); ok && obj.Tag != nil && obj.Tag.Name != "" {
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

func genPHPClass(b *strings.Builder, className string, n ast.Node, generated map[string]struct{}) {
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

func genPHPFields(b *strings.Builder, n ast.Node, indent int) {
	obj, ok := n.(*ast.Object)
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

func getPhpTypeForField(f *ast.Field) string {
	switch v := f.Value.(type) {
	case *ast.Value:
		return getPhpType(v)
	case *ast.Object:
		return getPhpObjectType(f.Key, v)
	case *ast.Array:
		return getPhpArrayType(f.Key, v)
	default:
		return "mixed"
	}
}

func getPhpType(v *ast.Value) string {
	if v != nil && v.Tag != nil {
		if t, ok := phpTypeMap[v.Tag.Type]; ok {
			return t
		}
	}
	return "mixed"
}

func getPhpObjectType(fieldKey string, obj *ast.Object) string {
	if obj != nil && obj.Tag != nil && obj.Tag.Name != "" {
		return exportPhpClassName(obj.Tag.Name)
	}
	return exportPhpClassName(fieldKey)
}

func getPhpArrayType(fieldKey string, a *ast.Array) string {
	if a == nil {
		return "array"
	}
	if a.Tag != nil && a.Tag.ChildType != ast.ValueTypeUnknown {
		if _, ok := phpTypeMap[a.Tag.ChildType]; ok {
			return "array"
		}
	}
	if len(a.Items) > 0 {
		switch item := a.Items[0].(type) {
		case *ast.Object:
			return "array"
		case *ast.Value:
			if item.Tag != nil {
				if _, ok := phpTypeMap[item.Tag.Type]; ok {
					return "array"
				}
			}
		}
	}
	return "array"
}

func genPHPNestedClasses(b *strings.Builder, n ast.Node, generated map[string]struct{}) {
	obj, ok := n.(*ast.Object)
	if !ok {
		return
	}

	for _, f := range obj.Fields {
		if f == nil {
			continue
		}

		switch v := f.Value.(type) {
		case *ast.Object:
			className := getPhpObjectType(f.Key, v)
			genPHPClass(b, className, v, generated)
			genPHPNestedClasses(b, v, generated)
		case *ast.Array:
			if nestedObj := findFirstObjectInArrayPHP(v); nestedObj != nil {
				className := getPhpObjectType(f.Key, nestedObj)
				genPHPClass(b, className, nestedObj, generated)
				genPHPNestedClasses(b, nestedObj, generated)
			}
		}
	}
}

func findFirstObjectInArrayPHP(a *ast.Array) *ast.Object {
	if a == nil {
		return nil
	}
	for _, item := range a.Items {
		if obj, ok := item.(*ast.Object); ok {
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

func genPHPObjectAssignments(b *strings.Builder, varName string, n ast.Node, indent int) {
	obj, ok := n.(*ast.Object)
	if !ok {
		return
	}

	for _, f := range obj.Fields {
		if f == nil {
			continue
		}
		prop := varName + "->" + exportPhpFieldName(f.Key)
		switch v := f.Value.(type) {
		case *ast.Value:
			WriteIndent(b, indent)
			b.WriteString(prop + " = " + formatPhpValueLiteral(v) + ";\n")
		case *ast.Object:
			className := getPhpObjectType(f.Key, v)
			WriteIndent(b, indent)
			b.WriteString(prop + " = new " + className + "();\n")
			genPHPObjectAssignments(b, prop, v, indent)
		case *ast.Array:
			if nestedObj := findFirstObjectInArrayPHP(v); nestedObj != nil {
				WriteIndent(b, indent)
				b.WriteString(prop + " = [\n")
				for i, item := range v.Items {
					WriteIndent(b, indent+1)
					switch iv := item.(type) {
					case *ast.Object:
						b.WriteString(genPHPObjectLiteral(iv, f.Key, indent+1))
					default:
						b.WriteString(formatPhpValueLiteral(iv.(*ast.Value)))
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

func genPHPArrayLiteral(a *ast.Array) string {
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
		case *ast.Value:
			sb.WriteString(formatPhpValueLiteral(iv))
		case *ast.Object:
			sb.WriteString(genPHPObjectLiteral(iv, "", 0))
		default:
			sb.WriteString("null")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func genPHPObjectLiteral(obj *ast.Object, fieldKey string, indent int) string {
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
		case *ast.Value:
			sb.WriteString(formatPhpValueLiteral(iv))
		case *ast.Object:
			sb.WriteString(genPHPObjectLiteral(iv, f.Key, indent+1))
		case *ast.Array:
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

func formatPhpValueLiteral(v *ast.Value) string {
	if v == nil {
		return "null"
	}

	if v.Tag != nil {
		switch v.Tag.Type {
		case ast.ValueTypeString, ast.ValueTypeBytes, ast.ValueTypeDateTime, ast.ValueTypeDate, ast.ValueTypeTime, ast.ValueTypeUUID, ast.ValueTypeDecimal, ast.ValueTypeEmail, ast.ValueTypeIP, ast.ValueTypeURL, ast.ValueTypeEnum, ast.ValueTypeImage:
			return "\"" + strings.ReplaceAll(v.Text, "\"", "\\\"") + "\""
		case ast.ValueTypeBool:
			if strings.EqualFold(v.Text, "true") {
				return "true"
			}
			return "false"
		case ast.ValueTypeInt, ast.ValueTypeInt8, ast.ValueTypeInt16, ast.ValueTypeInt32, ast.ValueTypeInt64, ast.ValueTypeUint, ast.ValueTypeUint8, ast.ValueTypeUint16, ast.ValueTypeUint32, ast.ValueTypeUint64, ast.ValueTypeFloat32, ast.ValueTypeFloat64:
			if v.Text == "" {
				return "0"
			}
			return v.Text
		}
	}
	return "\"" + strings.ReplaceAll(v.Text, "\"", "\\\"") + "\""
}

func PrintPHPStruct(n ast.Node) {
	fmt.Println(ToPHP(n))
}
