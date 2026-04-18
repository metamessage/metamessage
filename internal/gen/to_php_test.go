package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func TestToPHPGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "sample"},
		Fields: []*ast.Field{
			{Key: "name", Value: &ast.Value{Text: "alice", Tag: &ast.Tag{Type: ast.ValueTypeString}}},
			{Key: "data", Value: &ast.Value{Text: "abc", Tag: &ast.Tag{Type: ast.ValueTypeBytes}}},
			{Key: "ids", Value: &ast.Array{Tag: &ast.Tag{Type: ast.ValueTypeArray, ChildType: ast.ValueTypeInt}, Items: []ast.Node{
				&ast.Value{Text: "1", Tag: &ast.Tag{Type: ast.ValueTypeInt}},
			}}},
			{Key: "meta", Value: &ast.Object{Tag: &ast.Tag{Name: "meta"}, Fields: []*ast.Field{
				{Key: "count", Value: &ast.Value{Text: "5", Tag: &ast.Tag{Type: ast.ValueTypeInt}}},
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
	if !strings.Contains(out, "public ?string $name = null;") {
		t.Fatalf("expected name property with string type, got:\n%s", out)
	}
	if !strings.Contains(out, "public ?string $data = null;") {
		t.Fatalf("expected data property with string type, got:\n%s", out)
	}
	if !strings.Contains(out, "public ?array $ids = null;") {
		t.Fatalf("expected ids property with array type, got:\n%s", out)
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
