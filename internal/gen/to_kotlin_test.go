package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToKotlinGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "user"},
		Fields: []*ir.Field{
			{Key: "id", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI}, Text: "1"}},
			{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "Alice"}},
			{Key: "tags", Value: &ir.Array{Items: []ir.Node{
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "go"},
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "java"},
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
