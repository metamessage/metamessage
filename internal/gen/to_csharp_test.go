package gen

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestToCSharpGeneratesCode(t *testing.T) {
	obj := &ir.Object{
		Tag: &ir.Tag{Name: "sample"},
		Fields: []*ir.Field{
			{Key: "when", Value: &ir.Value{Text: "2024-01-01T00:00:00Z", Tag: &ir.Tag{Type: ir.ValueTypeDatetime}}},
			{Key: "id", Value: &ir.Value{Text: "123", Tag: &ir.Tag{Type: ir.ValueTypeBigint}}},
			{Key: "data", Value: &ir.Array{
				Tag: &ir.Tag{Name: "data", Type: ir.ValueTypeArr, ChildType: ir.ValueTypeStr},
				Items: []ir.Node{
					&ir.Value{Text: "abc", Tag: &ir.Tag{Type: ir.ValueTypeStr}},
				},
			}},
			{Key: "nested", Value: &ir.Object{
				Tag:    &ir.Tag{Name: "nested"},
				Fields: []*ir.Field{{Key: "count", Value: &ir.Value{Tag: &ir.Tag{Type: ir.ValueTypeI}, Text: "5"}}},
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
