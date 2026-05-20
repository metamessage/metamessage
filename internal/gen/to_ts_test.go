package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToTSGeneratesCode(t *testing.T) {
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

	out := ToTS(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "export interface User") {
		t.Fatalf("expected interface declaration in output, got:\n%s", out)
	}
	if !strings.Contains(out, "id: number;") {
		t.Fatalf("expected number field declaration (non-nullable), got:\n%s", out)
	}
	if !strings.Contains(out, "name: string;") {
		t.Fatalf("expected string field declaration (non-nullable), got:\n%s", out)
	}
	if !strings.Contains(out, "tags: Array<string>;") {
		t.Fatalf("expected Array<string> field declaration (non-nullable), got:\n%s", out)
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

// TestToTSNullableFields verifies that nullable and non-nullable fields are generated correctly
func TestToTSNullableFields(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "profile"},
		Fields: []*ir.Field{
			// Non-nullable (default)
			{Key: "id", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI, Nullable: false}, Text: "1"}},
			// Nullable
			{Key: "nickname", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr, Nullable: true}, Text: ""}},
			// Non-nullable array
			{Key: "tags", Value: &ir.Array{Tag: &ir.Tag{Nullable: false}, Items: []ir.Node{
				&ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeStr}, Text: "dev"},
			}}},
		},
	}

	out := ToTS(obj)
	t.Log("\n" + out)

	// Non-nullable fields should NOT have ?
	if !strings.Contains(out, "id: number;") {
		t.Fatalf("expected non-nullable id field (id: number;), got:\n%s", out)
	}
	if !strings.Contains(out, "tags: Array<string>;") {
		t.Fatalf("expected non-nullable tags field (tags: Array<string>;), got:\n%s", out)
	}

	// Nullable fields SHOULD have ?
	if !strings.Contains(out, "nickname?: string;") {
		t.Fatalf("expected nullable nickname field (nickname?: string;), got:\n%s", out)
	}
}
