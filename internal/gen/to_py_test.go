package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToPyGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "user"},
		Fields: []*ir.Field{
			{Key: "id", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeInt}, Text: "1"}},
			{Key: "name", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "Alice"}},
			{Key: "tags", Value: &ir.Array{Items: []ir.Node{
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "go"},
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeString}, Text: "java"},
			}}},
			{Key: "profile", Value: &ir.Object{Fields: []*ir.Field{
				{Key: "age", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeInt}, Text: "30"}},
			}}},
		},
	}

	out := ToPy(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "class User:") {
		t.Fatalf("expected top-level class declaration, got:\n%s", out)
	}
	if !strings.Contains(out, "id: int") {
		t.Fatalf("expected non-nullable id field (id: int), got:\n%s", out)
	}
	if !strings.Contains(out, "name: str") {
		t.Fatalf("expected non-nullable name field (name: str), got:\n%s", out)
	}
	if !strings.Contains(out, "tags: List[str]") {
		t.Fatalf("expected non-nullable tags field (tags: List[str]), got:\n%s", out)
	}
	if !strings.Contains(out, "class Profile:") {
		t.Fatalf("expected nested Profile class, got:\n%s", out)
	}
	if !strings.Contains(out, "user_data = User(\n") {
		t.Fatalf("expected multiline bound data output, got:\n%s", out)
	}
	if !strings.Contains(out, "\tid=1,\n") {
		t.Fatalf("expected formatted id field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tname=\"Alice\",\n") {
		t.Fatalf("expected formatted name field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\ttags=[\n") {
		t.Fatalf("expected formatted tags list in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tprofile=Profile(\n") {
		t.Fatalf("expected formatted profile instance in bound data, got:\n%s", out)
	}
}
