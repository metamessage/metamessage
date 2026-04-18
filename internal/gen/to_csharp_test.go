package gen

import (
	"strings"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

func TestToCSharpGeneratesCode(t *testing.T) {
	obj := &ast.Object{
		Tag: &ast.Tag{Name: "sample"},
		Fields: []*ast.Field{
			{Key: "when", Value: &ast.Value{Text: "2024-01-01T00:00:00Z", Tag: &ast.Tag{Type: ast.ValueTypeDateTime}}},
			{Key: "id", Value: &ast.Value{Text: "123", Tag: &ast.Tag{Type: ast.ValueTypeBigInt}}},
			{Key: "data", Value: &ast.Array{
				Tag: &ast.Tag{Name: "data", Type: ast.ValueTypeArray, ChildType: ast.ValueTypeString},
				Items: []ast.Node{
					&ast.Value{Text: "abc", Tag: &ast.Tag{Type: ast.ValueTypeString}},
				},
			}},
			{Key: "nested", Value: &ast.Object{
				Tag:    &ast.Tag{Name: "nested"},
				Fields: []*ast.Field{{Key: "count", Value: &ast.Value{Tag: &ast.Tag{Type: ast.ValueTypeInt}, Text: "5"}}},
			}},
		},
	}

	out := ToCSharp(obj)
	t.Log("\n" + out)
	if !strings.Contains(out, "using System;") {
		t.Fatalf("expected System import, got:\n%s", out)
	}
	if !strings.Contains(out, "using System.Collections.Generic;") {
		t.Fatalf("expected System.Collections.Generic import, got:\n%s", out)
	}
	if !strings.Contains(out, "public class Sample") {
		t.Fatalf("expected top-level class Sample, got:\n%s", out)
	}
	if !strings.Contains(out, "public List<string> Data") {
		t.Fatalf("expected List<string> property for data, got:\n%s", out)
	}
	if !strings.Contains(out, "public class Nested") {
		t.Fatalf("expected nested class Nested, got:\n%s", out)
	}
	if !strings.Contains(out, "public static class SampleData") {
		t.Fatalf("expected data helper class SampleData, got:\n%s", out)
	}
	if !strings.Contains(out, "public static readonly Sample Instance = CreateSample();") {
		t.Fatalf("expected bound sample instance declaration, got:\n%s", out)
	}
	if !strings.Contains(out, "Data = new List<string> {\n") {
		t.Fatalf("expected data list initializer, got:\n%s", out)
	}
	if !strings.Contains(out, "Nested = new Nested {") {
		t.Fatalf("expected nested object initializer, got:\n%s", out)
	}
}
