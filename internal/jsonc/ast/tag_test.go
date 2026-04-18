package ast

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"testing"
)

func TestParseMMTag_Basic(t *testing.T) {
	tag := "name=id; min=1; desc=用户ID; enums=active|pending|deleted"
	r, err := ParseMMTag(tag)
	fmt.Println("rrrr", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Name != "id" {
		t.Fatalf("expected name=id, got %s", r.Name)
	}
	if r.Desc != "用户ID" {
		t.Fatalf("expected desc=用户ID, got %s", r.Desc)
	}
	if r.Min != "1" {
		t.Fatalf("expected min=1, got %v", r.Min)
	}
}

func TestParseMMTag_Flags(t *testing.T) {
	tag := "nullable,default=abc,max=10,min=5"
	r, err := ParseMMTag(tag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Nullable {
		t.Fatalf("expected nullable=true")
	}
	if r.Default != "abc" {
		t.Fatalf("expected default=abc, got %s", r.Default)
	}
}

func TestParseMMTag_QuotedAndSemicolon(t *testing.T) {
	tag := `name="id"; desc="用户ID"; enum="active|pending"; pattern="^a,b$"; type=string; min=1; max=5; nullable; default="x"`
	r, err := ParseMMTag(tag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "id" {
		t.Fatalf("expected name=id, got %s", r.Name)
	}
	if r.Desc != "用户ID" {
		t.Fatalf("expected desc=用户ID, got %s", r.Desc)
	}
	if r.Enum != "active|pending" {
		t.Fatalf("expected enum=active|pending, got %s", r.Enum)
	}
	if r.Pattern != "^a,b$" {
		t.Fatalf("expected pattern=^a,b$, got %s", r.Pattern)
	}
	if r.Type != ValueTypeString {
		t.Fatalf("expected type=string, got %s", r.Type)
	}
	if r.Min != "1" || r.Max != "5" {
		t.Fatalf("expected min=1,max=5 got %v,%v", r.Min, r.Max)
	}
	if !r.Nullable {
		t.Fatalf("expected nullable flag present")
	}
	if r.Default != "x" {
		t.Fatalf("expected default=x got %s", r.Default)
	}
}

func TestTagValidate_NumericRange(t *testing.T) {
	tag := &Tag{Min: "1", Max: "10"}
	_, _, err := tag.ValidateInt(5)
	if err != nil {
		t.Fatalf("expected valid for 5, got %v", err)
	}
	_, _, err = tag.ValidateInt(0)
	if err == nil {
		t.Fatalf("expected invalid for 0")
	}
	_, _, err = tag.ValidateInt(11)
	if err == nil {
		t.Fatalf("expected invalid for 11")
	}
	// floats
	_, _, err = tag.ValidateFloat64(3.14)
	if err != nil {
		t.Fatalf("expected valid for 3.14, got  %v", err)
	}
}

func TestTagValidate_StringRuneLength(t *testing.T) {
	tag := &Tag{Min: "2", Max: "4"}
	// multibyte characters
	_, _, err := tag.ValidateString("你好") // 2 runes
	if err != nil {
		t.Fatalf("expected valid for 你好, got %v", err)
	}
	_, _, err = tag.ValidateString("😀😀😀😀😀") // 5 runes
	if err == nil {
		t.Fatalf("expected invalid for 5 runes")
	}
}

func TestTagValidate_PatternAndEnum(t *testing.T) {
	tag := &Tag{Pattern: `^\d{3}$`}
	_, _, err := tag.ValidateString("123")
	if err != nil {
		t.Fatalf("expected pattern match for 123, got  %v", err)
	}
	_, _, err = tag.ValidateString("12a")
	if err == nil {
		t.Fatalf("expected pattern mismatch for 12a")
	}

	etag := &Tag{Enum: "a|b|c"}
	_, _, err = etag.ValidateString("b")
	if err != nil {
		t.Fatalf("expected enum match for b")
	}
	_, _, err = etag.ValidateString("z")
	if err == nil {
		t.Fatalf("expected enum mismatch for z")
	}
}

func TestTagValidate_NilPointerNullable(t *testing.T) {
	tag := &Tag{Nullable: true}
	var p *int = nil
	_, _, err := tag.ValidateInt(*p)
	if err != nil {
		t.Fatalf("expected nil pointer accepted when nullable")
	}
	tag.Nullable = false
	_, _, err = tag.ValidateInt(*p)
	if err == nil {
		t.Fatalf("expected nil pointer rejected when not nullable")
	}
}

func TestTagValidate_TypeSpecificWrappers(t *testing.T) {
	intTag := &Tag{Min: "1", Max: "10"}
	_, _, err := intTag.ValidateInt(5)
	if err != nil {
		t.Fatalf("expected valid int value")
	}
	_, _, err = intTag.ValidateInt(0)
	if err == nil {
		t.Fatalf("expected invalid int value below min")
	}

	urlTag := &Tag{Type: ValueTypeURL, AllowEmpty: true}
	u, _ := url.Parse("https://example.com")
	_, _, err = urlTag.ValidateURL(*u)
	if err != nil {
		t.Fatalf("expected valid url value, got %v", err)
	}

	emailTag := &Tag{Type: ValueTypeEmail}
	_, _, err = emailTag.ValidateEmail("a@b.com")
	if err != nil {
		t.Fatalf("expected valid email value, got %v", err)
	}

	ipTag := &Tag{Type: ValueTypeIP, Version: 4}
	ip := net.ParseIP("127.0.0.1")
	_, _, err = ipTag.ValidateIP(ip)
	if err != nil {
		t.Fatalf("expected valid ipv4 value, got %v", err)
	}

	bigTag := &Tag{Type: ValueTypeBigInt}
	bi, _ := new(big.Int).SetString("12345678901234567890", 10)
	_, _, err = bigTag.ValidateBigInt(*bi)
	if err != nil {
		t.Fatalf("expected valid big.Int value, got %v", err)
	}
	b2 := big.NewInt(42)
	if _, _, err := bigTag.ValidateBigInt(*b2); err != nil {
		t.Fatalf("expected valid big.Int value object, got %v", err)
	}

	enumTag := &Tag{Enum: "a|b|c"}
	_, _, err = enumTag.ValidateEnum("b")
	if err != nil {
		t.Fatalf("expected valid enum value, got %v", err)
	}
	_, _, err = enumTag.ValidateEnum("d")
	if err == nil {
		t.Fatalf("expected invalid enum value")
	}
}
