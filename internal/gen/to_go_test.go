package gen

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

// gofmt -w . && go test ./internal/jsonc/ast -run TestGenerateGoStruct -v
// go test ./internal/jsonc -v -run TestGenerateGoStruct
func TestGenerateGoStruct(t *testing.T) {
	obj := &ast.Object{
		Fields: []*ast.Field{
			{Key: "name", Value: &ast.Value{Tag: &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20", Max: "30"}}},
			{Key: "data", Value: &ast.Value{
				Data: "werwerwe",
				Text: "werwerwe",
				Tag:  &ast.Tag{Name: "data", Type: ast.ValueTypeBytes}}},
			{Key: "sss", Value: &ast.Object{
				Fields: []*ast.Field{
					{Key: "name", Value: &ast.Value{Tag: &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"}}},
					{Key: "data", Value: &ast.Value{Tag: &ast.Tag{Name: "data", Type: ast.ValueTypeBytes}}},
					{Key: "sss", Value: &ast.Value{Tag: &ast.Tag{Name: "data", Type: ast.ValueTypeBytes}}},
				},
				Tag: &ast.Tag{
					Name: "obj",
				},
			}},
			{Key: "arr", Value: &ast.Array{
				Items: []ast.Node{
					&ast.Value{
						Data: "121212",
						Tag:  &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"}},
					&ast.Value{
						Data: "121212",
						Tag:  &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"}},
					&ast.Value{
						Data: "44334",
						Tag:  &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"}},
				},
				Tag: &ast.Tag{Name: "data", Type: ast.ValueTypeBytes}},
			},

			{Key: "String", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}}},
		},
		Tag: &ast.Tag{
			Name: "obj",
		},
	}

	out := ToGo(obj)
	fmt.Printf("Generated struct: \n%s", out)
}

func TestPrintGoStruct(t *testing.T) {

	val := &ast.Value{
		Data: "abc",
		Text: "abc",
		Tag:  &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"},
	}

	PrintGoStruct(val)

	arr := &ast.Array{
		Items: []ast.Node{&ast.Value{
			Data: "abc",
			Text: "abc",
			Tag:  &ast.Tag{Name: "name", Type: ast.ValueTypeString, Min: "20"},
		},
		},
		Tag: &ast.Tag{Name: "arr", Type: ast.ValueTypeArray},
	}

	// var arr []string = []string{"abc", "def"}
	// var arr []any = []any{"abc", "def"}
	// var arr = []any{"abc", "def"}
	// var arr []any = []any{"abc", "def"}

	PrintGoStruct(arr)
}

func TestToGoGeneratesValidGo(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "sample"},
		Fields: []*ast.Field{
			{Key: "when", Value: &ast.Value{Text: "2024-01-01T00:00:00Z", Tag: &ast.Tag{Type: ast.ValueTypeDateTime}}},
			{Key: "ip", Value: &ast.Value{Text: "127.0.0.1", Tag: &ast.Tag{Type: ast.ValueTypeIP}}},
			{Key: "site", Value: &ast.Value{Text: "https://example.com", Tag: &ast.Tag{Type: ast.ValueTypeURL}}},
			{Key: "id", Value: &ast.Value{Text: "123", Tag: &ast.Tag{Type: ast.ValueTypeBigInt}}},
			{Key: "data", Value: &ast.Value{Text: "abc", Tag: &ast.Tag{Type: ast.ValueTypeBytes}}},
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
	valueNode := &ast.Value{
		Tag: &ast.Tag{
			Name: "username",
			Type: ast.ValueTypeString,
		},
		Text: "zhangsan",
	}
	fmt.Println("=== Value 生成结果 ===")
	fmt.Println(ToGo(valueNode))

	// 测试 Array 类型
	arrayNode := &ast.Array{
		Tag: &ast.Tag{
			Name: "ages",
			Type: ast.ValueTypeInt8,
		},
		Items: []ast.Node{
			&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
			&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "20"},
		},
	}
	fmt.Println("\n=== Array 生成结果 ===")
	fmt.Println(ToGo(arrayNode))

	// 测试 Object 类型
	objectNode := &ast.Object{
		Tag: &ast.Tag{
			Name: "user_info",
			Type: ast.ValueTypeStruct,
		},
		Fields: []*ast.Field{
			{
				Key:   "user_name",
				Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "zhangsan"},
			},
			{
				Key:   "age",
				Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
			},
		},
	}
	fmt.Println("\n=== Object 生成结果 ===")
	fmt.Println(ToGo(objectNode))

	objectNode1 := &ast.Object{
		Tag: &ast.Tag{
			Name: "user_info",
			Type: ast.ValueTypeStruct,
		},
		Fields: []*ast.Field{
			{
				Key: "user_name",
				Value: &ast.Array{
					Tag: &ast.Tag{
						Name: "ages",
						Type: ast.ValueTypeInt8,
					},
					Items: []ast.Node{
						&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
						&ast.Object{
							Tag: &ast.Tag{
								Name: "user_info",
								Type: ast.ValueTypeStruct,
							},
							Fields: []*ast.Field{
								{
									Key:   "user_name",
									Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "zhangsan"},
								},
								{
									Key:   "age",
									Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
								},
							},
						},
					},
				},
			},
			{
				Key:   "age",
				Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt8}, Text: "18"},
			},
		},
	}
	fmt.Println("\n=== Object 生成结果 ===")
	fmt.Println(ToGo(objectNode1))
}
