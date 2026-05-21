package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToCppGeneratesCode(t *testing.T) {
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

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "#pragma once") {
		t.Fatalf("expected #pragma once, got:\n%s", out)
	}
	if !strings.Contains(out, "#include <cstdint>") {
		t.Fatalf("expected #include <cstdint>, got:\n%s", out)
	}
	if !strings.Contains(out, "#include <string>") {
		t.Fatalf("expected #include <string>, got:\n%s", out)
	}
	if !strings.Contains(out, "#include <vector>") {
		t.Fatalf("expected #include <vector>, got:\n%s", out)
	}

	if !strings.Contains(out, "struct User {") {
		t.Fatalf("expected struct User definition, got:\n%s", out)
	}
	if !strings.Contains(out, "int64_t id;") {
		t.Fatalf("expected int64_t id field, got:\n%s", out)
	}
	if !strings.Contains(out, "std::string name;") {
		t.Fatalf("expected std::string name field, got:\n%s", out)
	}
	if !strings.Contains(out, "bool active;") {
		t.Fatalf("expected bool active field, got:\n%s", out)
	}
	if !strings.Contains(out, "double score;") {
		t.Fatalf("expected double score field, got:\n%s", out)
	}
	if !strings.Contains(out, "std::vector<std::string> tags;") {
		t.Fatalf("expected std::vector<std::string> tags field, got:\n%s", out)
	}
	if !strings.Contains(out, "inline const User") {
		t.Fatalf("expected inline instance data, got:\n%s", out)
	}
}

func TestToCppWithNoName(t *testing.T) {
	obj := &ir.Object{
		Fields: []*ir.Field{
			{Key: "x", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeF32}, Text: "1.5"}},
		},
	}

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "struct Obj {") {
		t.Fatalf("expected default struct Obj, got:\n%s", out)
	}
	if !strings.Contains(out, "float x;") {
		t.Fatalf("expected float field, got:\n%s", out)
	}
}

func TestToCppWithBytes(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "blob"},
		Fields: []*ir.Field{
			{Key: "data", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeBytes}, Text: "hello"}},
		},
	}

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "std::vector<uint8_t>") {
		t.Fatalf("expected std::vector<uint8_t> type, got:\n%s", out)
	}
}

func TestToCppWithNested(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "order"},
		Fields: []*ir.Field{
			{Key: "item", Value: &ir.Object{
				Tag: &ir.Tag{Name: "item"},
				Fields: []*ir.Field{
					{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "widget"}},
				},
			}},
		},
	}

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "struct Order {") {
		t.Fatalf("expected Order struct, got:\n%s", out)
	}
	if !strings.Contains(out, "Item item;") {
		t.Fatalf("expected Item item field, got:\n%s", out)
	}
}

func TestToCppWithMapType(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "cfg"},
		Fields: []*ir.Field{
			{Key: "props", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeMap}, Text: ""}},
		},
	}

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "#include <map>") {
		t.Fatalf("expected #include <map> for map type, got:\n%s", out)
	}
	if !strings.Contains(out, "std::map<std::string, std::any> props;") {
		t.Fatalf("expected map field, got:\n%s", out)
	}
}

func TestToCppArrayNestedObject(t *testing.T) {
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

	out := ToCpp(obj)
	t.Log("\n" + out)

	if !strings.Contains(out, "std::vector<Students> students;") {
		t.Fatalf("expected vector field for array of objects, got:\n%s", out)
	}
}

func TestPrintCppStruct(t *testing.T) {
	val := &ir.Value{
		Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeStr},
		Text: "hello",
	}
	PrintCppStruct(val)
}

func TestToCppNil(t *testing.T) {
	out := ToCpp(nil)
	if out != "" {
		t.Fatalf("expected empty string for nil input, got:\n%s", out)
	}
}
