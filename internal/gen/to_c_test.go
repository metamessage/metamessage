package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToCGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "user"},
		Fields: []*ir.Field{
			{Key: "id", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI}, Text: "1"}},
			{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "Alice"}},
			{Key: "active", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeBool}, Text: "true"}},
			{Key: "score", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeF64}, Text: "95.5"}},
			{Key: "tags", Value: &ir.Array{Items: []ir.Node{
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "go"},
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "java"},
			}}},
			{Key: "profile", Value: &ir.Object{Fields: []*ir.Field{
				{Key: "age", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeU8}, Text: "30"}},
			}}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "#ifndef MM_GEN_USER_H") {
		t.Fatalf("expected header guard, got:\n%s", out)
	}
	if !strings.Contains(out, "typedef struct {") {
		t.Fatalf("expected struct definition, got:\n%s", out)
	}
	if !strings.Contains(out, "typedef struct {\n\tint64_t id;") {
		t.Fatalf("expected id field of type int64_t, got:\n%s", out)
	}
	if !strings.Contains(out, "char* name;") {
		t.Fatalf("expected name field of type char*, got:\n%s", out)
	}
	if !strings.Contains(out, "bool active;") {
		t.Fatalf("expected active field of type bool, got:\n%s", out)
	}
	if !strings.Contains(out, "double score;") {
		t.Fatalf("expected score field of type double, got:\n%s", out)
	}
	if !strings.Contains(out, "mm_arr_t tags;") {
		t.Fatalf("expected tags field of type mm_arr_t, got:\n%s", out)
	}
	if !strings.Contains(out, "#endif /* MM_GEN_USER_H */") {
		t.Fatalf("expected closing header guard, got:\n%s", out)
	}
}

func TestToCWithNoName(t *testing.T) {
	obj := &ir.Object{
		Fields: []*ir.Field{
			{Key: "x", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeF32}, Text: "1.5"}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "#ifndef MM_GEN_OBJ_H") {
		t.Fatalf("expected default header guard MM_GEN_OBJ_H, got:\n%s", out)
	}
	if !strings.Contains(out, "float x;") {
		t.Fatalf("expected float field, got:\n%s", out)
	}
}

func TestToCWithBytes(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "blob"},
		Fields: []*ir.Field{
			{Key: "data", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeBytes}, Text: "hello"}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "mm_bytes_t") {
		t.Fatalf("expected mm_bytes_t helper type, got:\n%s", out)
	}
	if !strings.Contains(out, "mm_bytes_t data;") {
		t.Fatalf("expected data field of type mm_bytes_t, got:\n%s", out)
	}
}

func TestToCWithNestedStruct(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "order"},
		Fields: []*ir.Field{
			{Key: "item", Value: &ir.Object{
				Tag: &ir.Tag{Name: "item"},
				Fields: []*ir.Field{
					{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "widget"}},
					{Key: "price", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeF64}, Text: "9.99"}},
				},
			}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "typedef struct {\n\tchar* name;\n\tdouble price;\n} Item") {
		t.Fatalf("expected nested Item struct, got:\n%s", out)
	}
}

func TestToCWithUnicodeFieldName(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "user_info"},
		Fields: []*ir.Field{
			{Key: "user name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "Alice"}},
			{Key: "123field", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI}, Text: "42"}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "char* user_name;") {
		t.Fatalf("expected snake_case field name 'user_name', got:\n%s", out)
	}
	if !strings.Contains(out, "int64_t f_123field;") {
		t.Fatalf("expected prefixed field name for digit-leading field, got:\n%s", out)
	}
}

func TestToCArrayNestedObject(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "classroom"},
		Fields: []*ir.Field{
			{Key: "students", Value: &ir.Array{Items: []ir.Node{
				&ir.Object{
					Fields: []*ir.Field{
						{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "Bob"}},
					},
				},
			}}},
		},
	}

	out := ToC(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "mm_arr_t") {
		t.Fatalf("expected array type, got:\n%s", out)
	}
}

func TestPrintCStruct(t *testing.T) {
	val := &ir.Value{
		Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeStr},
		Text: "hello",
	}
	PrintCStruct(val)
}

func TestToCNil(t *testing.T) {
	out := ToC(nil)
	if out != "" {
		t.Fatalf("expected empty string for nil input, got:\n%s", out)
	}
}
