package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToPHPGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "sample"},
		Fields: []*ir.Field{
			{Key: "name", Value: &ir.Value{Text: "alice", Tag: &ir.Tag{Type: ir.ValueTypeStr}}},
			{Key: "data", Value: &ir.Value{Text: "abc", Tag: &ir.Tag{Type: ir.ValueTypeBytes}}},
			{Key: "ids", Value: &ir.Array{Tag: &ir.Tag{Type: ir.ValueTypeArr, ChildType: ir.ValueTypeI}, Items: []ir.Node{
				&ir.Value{Text: "1", Tag: &ir.Tag{Type: ir.ValueTypeI}},
			}}},
			{Key: "meta", Value: &ir.Object{Tag: &ir.Tag{Name: "meta"}, Fields: []*ir.Field{
				{Key: "count", Value: &ir.Value{Text: "5", Tag: &ir.Tag{Type: ir.ValueTypeI}}},
			}}},
		},
	}

	out := ToPHP(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "<?php") {
		t.Fatalf("expected php opening tag, got:\n%s", out)
	}
	if !strings.Contains(out, "class Sample") {
		t.Fatalf("expected class Sample, got:\n%s", out)
	}
	if !strings.Contains(out, "public string $name;") {
		t.Fatalf("expected non-nullable name property (public string $name;), got:\n%s", out)
	}
	if !strings.Contains(out, "public string $data;") {
		t.Fatalf("expected non-nullable data property (public string $data;), got:\n%s", out)
	}
	if !strings.Contains(out, "public array $ids;") {
		t.Fatalf("expected non-nullable ids property (public array $ids;), got:\n%s", out)
	}
	if !strings.Contains(out, "class Meta") {
		t.Fatalf("expected nested class Meta, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData = new Sample();") {
		t.Fatalf("expected bound php instance declaration, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData->name = \"alice\";") {
		t.Fatalf("expected bound name assignment, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData->data = \"abc\";") {
		t.Fatalf("expected bound data assignment, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData->ids = [1];") {
		t.Fatalf("expected bound ids assignment, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData->meta = new Meta();") {
		t.Fatalf("expected bound meta instantiation, got:\n%s", out)
	}
	if !strings.Contains(out, "$sampleData->meta->count = 5;") {
		t.Fatalf("expected bound meta count assignment, got:\n%s", out)
	}
}
