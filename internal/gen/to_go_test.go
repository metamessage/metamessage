package gen

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

// gofmt -w . && go test ./internal/jsonc/ast -run TestGenerateGoStruct -v
// go test ./internal/jsonc -v -run TestGenerateGoStruct
func TestGenerateGoStruct(t *testing.T) {
	obj := &ir.Object{
		Path: "",
		Fields: []*ir.Field{
			{Key: "name", Value: &ir.Value{Path: "name", Tag: &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20", Max: "30"}}},
			{Key: "data", Value: &ir.Value{
				Path: "data",
				Data: "werwerwe",
				Text: "werwerwe",
				Tag:  &ir.Tag{Name: "data", Type: ir.ValueTypeBytes}}},
			{Key: "sss", Value: &ir.Object{
				Path: "sss",
				Fields: []*ir.Field{
					{Key: "name", Value: &ir.Value{Path: "sss.name", Tag: &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"}}},
					{Key: "data", Value: &ir.Value{Path: "sss.data", Tag: &ir.Tag{Name: "data", Type: ir.ValueTypeBytes}}},
					{Key: "sss", Value: &ir.Value{Path: "sss.sss", Tag: &ir.Tag{Name: "data", Type: ir.ValueTypeBytes}}},
				},
				Tag: &ir.Tag{
					Name: "obj",
				},
			}},
			{Key: "arr", Value: &ir.Array{
				Path: "arr",
				Items: []ir.Node{
					&ir.Value{
						Path: "arr.0",
						Data: "121212",
						Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"}},
					&ir.Value{
						Path: "arr.1",
						Data: "121212",
						Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"}},
					&ir.Value{
						Path: "arr.2",
						Data: "44334",
						Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"}},
				},
				Tag: &ir.Tag{Name: "data", Type: ir.ValueTypeBytes}},
			},

			{Key: "String", Value: &ir.Value{Path: "String", Tag: &ir.Tag{Type: ir.ValueTypeString}}},
		},
		Tag: &ir.Tag{
			Name: "obj",
		},
	}

	out := ToGo(obj)
	fmt.Printf("Generated struct: \n%s", out)
}

func TestPrintGoStruct(t *testing.T) {

	val := &ir.Value{
		Path: "name",
		Data: "abc",
		Text: "abc",
		Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"},
	}

	PrintGoStruct(val)

	arr := &ir.Array{
		Path: "arr",
		Items: []ir.Node{&ir.Value{
			Path: "arr.0",
			Data: "abc",
			Text: "abc",
			Tag:  &ir.Tag{Name: "name", Type: ir.ValueTypeString, Min: "20"},
		},
		},
		Tag: &ir.Tag{Name: "arr", Type: ir.ValueTypeArray},
	}

	// var arr []string = []string{"abc", "def"}
	// var arr []any = []any{"abc", "def"}
	// var arr = []any{"abc", "def"}
	// var arr []any = []any{"abc", "def"}

	PrintGoStruct(arr)
}

func TestToGoGeneratesValidGo(t *testing.T) {
	obj := &ir.Object{
		Path: "",
		Tag:  &ir.Tag{Name: "sample"},
		Fields: []*ir.Field{
			{Key: "when", Value: &ir.Value{Path: "when", Text: "2024-01-01T00:00:00Z", Tag: &ir.Tag{Type: ir.ValueTypeDateTime}}},
			{Key: "ip", Value: &ir.Value{Path: "ip", Text: "127.0.0.1", Tag: &ir.Tag{Type: ir.ValueTypeIP}}},
			{Key: "site", Value: &ir.Value{Path: "site", Text: "https://example.com", Tag: &ir.Tag{Type: ir.ValueTypeURL}}},
			{Key: "id", Value: &ir.Value{Path: "id", Text: "123", Tag: &ir.Tag{Type: ir.ValueTypeBigInt}}},
			{Key: "data", Value: &ir.Value{Path: "data", Text: "abc", Tag: &ir.Tag{Type: ir.ValueTypeBytes}}},
		},
	}

	out := ToGo(obj)
	if !strings.Contains(out, "import (") {
		t.Fatalf("expected import block, got:\n%s", out)
	}
	for _, pkg := range []string{"math/big", "net", "net/url", "time"} {
		if !strings.Contains(out, pkg) {
			t.Fatalf("expected import %q in generated code, got:\n%s", pkg, out)
		}
	}
	if !strings.Contains(out, "[]byte(") {
		t.Fatalf("expected []byte literal in generated code, got:\n%s", out)
	}

	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("generated Go code is invalid: %v\n%s", err, out)
	}
}

func TestGenerateGoStruct2(t *testing.T) {
	// 测试 Value 类型
	valueNode := &ir.Value{
		Path: "username",
		Tag: &ir.Tag{
			Name: "username",
			Type: ir.ValueTypeString,
		},
		Text: "zhangsan",
	}
	fmt.Println("=== Value 生成结果 ===")
	fmt.Println(ToGo(valueNode))

	// 测试 Array 类型
	arrayNode := &ir.Array{
		Path: "ages",
		Tag: &ir.Tag{
			Name: "ages",
			Type: ir.ValueTypeInt8,
		},
		Items: []ir.Node{
			&ir.Value{Path: "ages.0", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "18"},
			&ir.Value{Path: "ages.1", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "20"},
		},
	}
	fmt.Println("\n=== Array 生成结果 ===")
	fmt.Println(ToGo(arrayNode))

	// 测试 Object 类型
	objectNode := &ir.Object{
		Path: "",
		Tag: &ir.Tag{
			Name: "user_info",
			Type: ir.ValueTypeObject,
		},
		Fields: []*ir.Field{
			{
				Key:   "user_name",
				Value: &ir.Value{Path: "user_name", Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "zhangsan"},
			},
			{
				Key:   "age",
				Value: &ir.Value{Path: "age", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "18"},
			},
		},
	}
	fmt.Println("\n=== Object 生成结果 ===")
	fmt.Println(ToGo(objectNode))

	objectNode1 := &ir.Object{
		Path: "",
		Tag: &ir.Tag{
			Name: "user_info",
			Type: ir.ValueTypeObject,
		},
		Fields: []*ir.Field{
			{
				Key: "user_name",
				Value: &ir.Array{
					Path: "user_name",
					Tag: &ir.Tag{
						Name: "ages",
						Type: ir.ValueTypeInt8,
					},
					Items: []ir.Node{
						&ir.Value{Path: "user_name.0", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "18"},
						&ir.Object{
							Path: "user_name.1",
							Tag: &ir.Tag{
								Name: "user_info",
								Type: ir.ValueTypeObject,
							},
							Fields: []*ir.Field{
								{
									Key:   "user_name",
									Value: &ir.Value{Path: "user_name.1.user_name", Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "zhangsan"},
								},
								{
									Key:   "age",
									Value: &ir.Value{Path: "user_name.1.age", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "18"},
								},
							},
						},
					},
				},
			},
			{
				Key:   "age",
				Value: &ir.Value{Path: "age", Tag: &ir.Tag{Type: ir.ValueTypeInt8}, Text: "18"},
			},
		},
	}
	fmt.Println("\n=== Object 生成结果 ===")
	fmt.Println(ToGo(objectNode1))
}

// TestDuplicateStructNamesWithPath demonstrates how Path prevents duplicate struct names
func TestDuplicateStructNamesWithPath(t *testing.T) {
	// Two structs with the same name but different paths should generate different struct names
	obj := &ir.Object{
		Path: "",
		Tag:  &ir.Tag{Name: "root"},
		Fields: []*ir.Field{
			{
				Key: "user",
				Value: &ir.Object{
					Path: "user",
					Tag:  &ir.Tag{Name: "profile"},
					Fields: []*ir.Field{
						{Key: "name", Value: &ir.Value{Path: "user.name", Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "alice"}},
					},
				},
			},
			{
				Key: "admin",
				Value: &ir.Object{
					Path: "admin",
					Tag:  &ir.Tag{Name: "profile"},
					Fields: []*ir.Field{
						{Key: "name", Value: &ir.Value{Path: "admin.name", Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "bob"}},
					},
				},
			},
		},
	}

	out := ToGo(obj)
	fmt.Printf("Generated code with different paths for same struct name: \n%s", out)

	// Verify both profiles are generated as separate structs
	if !strings.Contains(out, "type AdminProfile struct") {
		t.Fatalf("expected AdminProfile struct, got:\n%s", out)
	}
	if !strings.Contains(out, "type UserProfile struct") {
		t.Fatalf("expected UserProfile struct, got:\n%s", out)
	}
}
