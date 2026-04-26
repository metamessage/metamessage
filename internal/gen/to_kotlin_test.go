package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc/ast"
)

func TestToKotlinGeneratesCode(t *testing.T) {
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

	out := ToKotlin(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "class User") {
		t.Fatalf("expected class declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "var id: Int = 0") {
		t.Fatalf("expected non-nullable Int field (var id: Int = 0), got:\n%s", out)
	}
	if !strings.Contains(out, "var name: String = \"\"") {
		t.Fatalf("expected non-nullable String field (var name: String = \"\"), got:\n%s", out)
	}
	if !strings.Contains(out, "var tags: List<String> = emptyList()") {
		t.Fatalf("expected non-nullable List<String> field, got:\n%s", out)
	}
	if !strings.Contains(out, "val userData = User()") {
		t.Fatalf("expected bound data instance declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.id = 1") {
		t.Fatalf("expected bound id assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.name = \"Alice\"") {
		t.Fatalf("expected bound name assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.tags = listOf(\n") {
		t.Fatalf("expected bound tags list assignment in output, got:\n%s", out)
	}
}
