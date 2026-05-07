package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ast"
)

func TestToJavaGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "user"},
		Fields: []*ast.Field{
			{Key: "id", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt}, Text: "1"}},
			{Key: "name", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "Alice"}},
			{Key: "tags", Value: &ast.Array{Tag: &ast.Tag{Type: ast.ValueTypeArray}, Items: []ast.Node{
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "go"},
				&ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeString}, Text: "java"},
			}}},
		},
	}

	out := ToJava(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "public class User") {
		t.Fatalf("expected class declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "public int id;") {
		t.Fatalf("expected int field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "public String name;") {
		t.Fatalf("expected String field declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "public class UserData") {
		t.Fatalf("expected data class declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "public static final User userData = createUser();") {
		t.Fatalf("expected bound data field in output, got:\n%s", out)
	}
	if !strings.Contains(out, "obj.id = 1;") {
		t.Fatalf("expected data binding assignment in output, got:\n%s", out)
	}
	if !strings.Contains(out, "obj.tags = List.of(\"go\", \"java\");") {
		t.Fatalf("expected list data binding in output, got:\n%s", out)
	}
	if !strings.Contains(out, "import java.util.List;") {
		t.Fatalf("expected java.util.List import in output, got:\n%s", out)
	}
}
