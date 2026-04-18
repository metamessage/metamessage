package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func TestToJsGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "user"},
		Fields: []*ast.Field{
			{Key: "id", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt}, Text: "1"}},
			{Key: "name", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "Alice"}},
			{Key: "tags", Value: &ast.Array{Items: []ast.Node{
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "go"},
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "java"},
			}}},
			{Key: "profile", Value: &ast.Object{Fields: []*ast.Field{
				{Key: "age", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt}, Text: "30"}},
			}}},
		},
	}

	out := ToJS(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "export class User") {
		t.Fatalf("expected top-level class declaration, got:\n%s", out)
	}
	if !strings.Contains(out, "this.id = null;") {
		t.Fatalf("expected id property initialization, got:\n%s", out)
	}
	if !strings.Contains(out, "this.tags = [];") {
		t.Fatalf("expected tags array initialization, got:\n%s", out)
	}
	if !strings.Contains(out, "export class Profile") {
		t.Fatalf("expected nested Profile class, got:\n%s", out)
	}
	if !strings.Contains(out, "export const userData = new User();") {
		t.Fatalf("expected bound data instance declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.id = 1;") {
		t.Fatalf("expected bound id assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.name = \"Alice\";") {
		t.Fatalf("expected bound name assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.tags = [") {
		t.Fatalf("expected bound tags assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "userData.profile = new Profile();") {
		t.Fatalf("expected bound profile object creation in output, got:\n%s", out)
	}
}
