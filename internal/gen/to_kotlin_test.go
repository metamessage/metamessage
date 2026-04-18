package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
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
	if !strings.Contains(out, "var id: Int? = null") {
		t.Fatalf("expected Int field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "var name: String? = null") {
		t.Fatalf("expected String field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "var tags: List<String>? = null") {
		t.Fatalf("expected List<String> field declaration in output, got:\n%s", out)
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
