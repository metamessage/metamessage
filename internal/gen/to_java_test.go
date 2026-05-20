package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToJavaGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "user"},
		Fields: []*ir.Field{
			{Key: "id", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI}, Text: "1"}},
			{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "Alice"}},
			{Key: "tags", Value: &ir.Array{Tag: &ir.Tag{Type: ir.ValueTypeArr}, Items: []ir.Node{
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "go"},
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "java"},
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
