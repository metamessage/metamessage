package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToRustGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "sample"},
		Fields: []*ir.Field{
			{Key: "name", Value: &ir.Value{Text: "alice", Tag: &ir.Tag{Type: ir.ValueTypeString}}},
			{Key: "data", Value: &ir.Value{Text: "abc", Tag: &ir.Tag{Type: ir.ValueTypeBytes}}},
			{Key: "ids", Value: &ir.Array{Tag: &ir.Tag{Type: ir.ValueTypeArray, ChildType: ir.ValueTypeInt}, Items: []ir.Node{
				&ir.Value{Text: "1", Tag: &ir.Tag{Type: ir.ValueTypeInt}},
				&ir.Value{Text: "2", Tag: &ir.Tag{Type: ir.ValueTypeInt}},
			}}},
			{Key: "meta", Value: &ir.Object{Tag: &ir.Tag{Name: "meta"}, Fields: []*ir.Field{
				{Key: "count", Value: &ir.Value{Text: "5", Tag: &ir.Tag{Type: ir.ValueTypeInt}}},
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
