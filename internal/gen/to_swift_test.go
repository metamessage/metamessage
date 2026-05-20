package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToSwiftGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "sample"},
		Fields: []*ir.Field{
			{Key: "name", Value: &ir.Value{Text: "alice", Tag: &ir.Tag{Type: ir.ValueTypeStr}}},
			{Key: "data", Value: &ir.Value{Text: "abc", Tag: &ir.Tag{Type: ir.ValueTypeBytes}}},
			{Key: "ids", Value: &ir.Array{Tag: &ir.Tag{Type: ir.ValueTypeArr, ChildType: ir.ValueTypeI}, Items: []ir.Node{
				&ir.Value{Text: "1", Tag: &ir.Tag{Type: ir.ValueTypeI}},
				&ir.Value{Text: "2", Tag: &ir.Tag{Type: ir.ValueTypeI}},
			}}},
			{Key: "meta", Value: &ir.Object{Tag: &ir.Tag{Name: "meta"}, Fields: []*ir.Field{
				{Key: "count", Value: &ir.Value{Text: "5", Tag: &ir.Tag{Type: ir.ValueTypeI}}},
			}}},
		},
	}

	out := ToSwift(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "struct Sample") {
		t.Fatalf("expected top-level struct Sample, got:\n%s", out)
	}
	if !strings.Contains(out, "var name: String?") {
		t.Fatalf("expected name field with String type, got:\n%s", out)
	}
	if !strings.Contains(out, "var data: Data?") {
		t.Fatalf("expected data field with Data type, got:\n%s", out)
	}
	if !strings.Contains(out, "var ids: [Int]?") {
		t.Fatalf("expected ids field with [Int] type, got:\n%s", out)
	}
	if !strings.Contains(out, "struct Meta") {
		t.Fatalf("expected nested struct Meta, got:\n%s", out)
	}
	if !strings.Contains(out, "let sampleData: Sample = Sample(\n") {
		t.Fatalf("expected multiline bound data output, got:\n%s", out)
	}
	if !strings.Contains(out, "\tname: \"alice\",\n") {
		t.Fatalf("expected formatted name field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tdata: Data(),\n") {
		t.Fatalf("expected formatted data field in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tids: [\n") {
		t.Fatalf("expected formatted ids array in bound data, got:\n%s", out)
	}
	if !strings.Contains(out, "\tmeta: Meta(\n") {
		t.Fatalf("expected formatted meta initializer in bound data, got:\n%s", out)
	}
}
