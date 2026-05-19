package parser

import (
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
)

func scan(input string) *Parser {
	return New(scanner.New(input).ScanAll())
}

func TestParseObjectValNil_EOFAfterColon(t *testing.T) {
	p := scan(`{"key": `)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	obj, ok := node.(*ir.Object)
	if !ok {
		t.Fatalf("expected *ir.Object, got %T", node)
	}
	if len(obj.Fields) != 0 {
		t.Fatalf("expected 0 fields (val==nil path), got %d", len(obj.Fields))
	}
}

func TestParseObjectValNil_EOFAfterPartialObject(t *testing.T) {
	p := scan(`{"name": "Alice", "age": `)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	obj, ok := node.(*ir.Object)
	if !ok {
		t.Fatalf("expected *ir.Object, got %T", node)
	}
	if len(obj.Fields) != 1 {
		t.Fatalf("expected 1 field (name), got %d", len(obj.Fields))
	}
	if obj.Fields[0].Key != "name" {
		t.Fatalf("expected key 'name', got %q", obj.Fields[0].Key)
	}
}

func TestParseObjectValNil_EOFAfterComma(t *testing.T) {
	p := scan(`{"key": "val", `)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	obj, ok := node.(*ir.Object)
	if !ok {
		t.Fatalf("expected *ir.Object, got %T", node)
	}
	if len(obj.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(obj.Fields))
	}
}

func TestParseObjectValNil_ValidObject(t *testing.T) {
	p := scan(`{"name": "Alice", "age": 18}`)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	obj, ok := node.(*ir.Object)
	if !ok {
		t.Fatalf("expected *ir.Object, got %T", node)
	}
	if len(obj.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(obj.Fields))
	}
}

func TestParseObjectValNil_EmptyObject(t *testing.T) {
	p := scan(`{}`)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	obj, ok := node.(*ir.Object)
	if !ok {
		t.Fatalf("expected *ir.Object, got %T", node)
	}
	if len(obj.Fields) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(obj.Fields))
	}
}

func TestParseObjectValNil_RBraceAfterColon(t *testing.T) {
	p := scan(`{"key":}`)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for } after colon (parse sees unexpected token)")
	}
}

func TestParseArrayValNil_EOFAfterOpen(t *testing.T) {
	p := scan(`[`)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for empty array (ValidateVec: not allow empty)")
	}
}

func TestParseArrayValNil_EOFAfterItem(t *testing.T) {
	p := scan(`[1, `)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := node.(*ir.Array)
	if !ok {
		t.Fatalf("expected *ir.Array, got %T", node)
	}
	if len(arr.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(arr.Items))
	}
}

func TestParseValNil_DoesNotCrash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"EOF after open brace", `{`, false},
		{"EOF after open bracket", `[`, true},
		{"EOF after key colon", `{"a": `, false},
		{"EOF after comma in object", `{"a": 1, `, false},
		{"EOF after comma in array", `[1, `, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := scan(tt.input)
			_, err := p.Parse()
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error: %v, expectErr=%v", err, tt.expectErr)
			}
		})
	}
}

func TestParseNullLiteral_Unsupported(t *testing.T) {
	p := scan(`null`)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for null literal")
	}
	if err.Error() != "null literal is not supported" {
		t.Fatalf("unexpected error message: %v", err)
	}
}
