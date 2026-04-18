package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func TestToRustGeneratesCode(t *testing.T) {
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

	out := ToRust(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "pub struct Sample") {
		t.Fatalf("expected top-level struct Sample, got:\n%s", out)
	}
	if !strings.Contains(out, "pub name: Option<String>") {
		t.Fatalf("expected name field with String type, got:\n%s", out)
	}
	if !strings.Contains(out, "pub data: Option<Vec<u8>>") {
		t.Fatalf("expected data field with Vec<u8> type, got:\n%s", out)
	}
	if !strings.Contains(out, "pub ids: Option<Vec<i32>>") {
		t.Fatalf("expected ids field with Vec<i32> type, got:\n%s", out)
	}
	if !strings.Contains(out, "pub struct Meta") {
		t.Fatalf("expected nested struct Meta, got:\n%s", out)
	}
	if !strings.Contains(out, "pub fn sample_data() -> Sample") {
		t.Fatalf("expected bound data function declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "name: Some(\"alice\".to_string())") {
		t.Fatalf("expected bound name value in output, got:\n%s", out)
	}
	if !strings.Contains(out, "data: Some(b\"abc\".to_vec())") {
		t.Fatalf("expected bound data bytes value in output, got:\n%s", out)
	}
	if !strings.Contains(out, "ids: Some(vec![1, 2])") {
		t.Fatalf("expected bound ids array value in output, got:\n%s", out)
	}
	if !strings.Contains(out, "meta: Some(Meta {\n") {
		t.Fatalf("expected bound meta initializer in output, got:\n%s", out)
	}
}
