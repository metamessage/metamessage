package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func TestToTSGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "user"},
		Fields: []*ast.Field{
			{Key: "id", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt}, Text: "1"}},
			{Key: "name", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "Alice"}},
			{Key: "tags", Value: &ast.Array{Items: []ast.Node{
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "go"},
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "java"},
			}}},
		},
	}

	out := ToTS(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "export interface User") {
		t.Fatalf("expected interface declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "id?: number;") {
		t.Fatalf("expected number field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "name?: string;") {
		t.Fatalf("expected string field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "tags?: Array<string>;") {
		t.Fatalf("expected Array<string> field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export const userData: User = {") {
		t.Fatalf("expected bound data object declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "id: 1,") {
		t.Fatalf("expected bound data id value in output, got:\n%s", out)
	}
	if !strings.Contains(out, "name: \"Alice\",") {
		t.Fatalf("expected bound data name value in output, got:\n%s", out)
	}
	if !strings.Contains(out, "tags: [") {
		t.Fatalf("expected bound data array value in output, got:\n%s", out)
	}
}
