package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ast"
)

func TestToSwiftGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "sample"},
		Fields: []*ast.Field{
			{Key: "name", Value: &ast.Value{Text: "alice", Tag: &ast.Tag{Type: ast.ValueTypeString}}},
			{Key: "data", Value: &ast.Value{Text: "abc", Tag: &ast.Tag{Type: ast.ValueTypeBytes}}},
			{Key: "ids", Value: &ast.Array{Tag: &ast.Tag{Type: ast.ValueTypeArray, ChildType: ast.ValueTypeInt}, Items: []ast.Node{
				&ast.Value{Text: "1", Tag: &ast.Tag{Type: ast.ValueTypeInt}},
				&ast.Value{Text: "2", Tag: &ast.Tag{Type: ast.ValueTypeInt}},
			}}},
			{Key: "meta", Value: &ast.Object{Tag: &ast.Tag{Name: "meta"}, Fields: []*ast.Field{
				{Key: "count", Value: &ast.Value{Text: "5", Tag: &ast.Tag{Type: ast.ValueTypeInt}}},
			}}},
		},
	}

	out := ToSwift(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "struct Sample") {
		t.Fatalf("expected top-level struct Sample, got:\n%s", out)
	}
	if !strings.Contains(out, "var name: String?") {
		t.Fatalf("expected name field with String type, got:\n%s", out)
	}
	if !strings.Contains(out, "var data: Data?") {
		t.Fatalf("expected data field with Data type, got:\n%s", out)
	}
	if !strings.Contains(out, "var ids: [Int]?") {
		t.Fatalf("expected ids field with [Int] type, got:\n%s", out)
	}
	if !strings.Contains(out, "struct Meta") {
		t.Fatalf("expected nested struct Meta, got:\n%s", out)
	}
	if !strings.Contains(out, "let sampleData: Sample = Sample(\n") {
		t.Fatalf("expected multiline bound data output, got:\n%s", out)
	}
	if !strings.Contains(out, "\tname: \"alice\",\n") {
		t.Fatalf("expected formatted name field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tdata: Data(),\n") {
		t.Fatalf("expected formatted data field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tids: [\n") {
		t.Fatalf("expected formatted ids array in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tmeta: Meta(\n") {
		t.Fatalf("expected formatted meta initializer in bound data, got:\n%s", out)
	}
}
