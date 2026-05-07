package ast

import "testing"

func TestParseMIME(t *testing.T) {
	cases := []struct {
		input string
		want  MIME
	}{
		{"application/json", MIMEJson},
		{"text/html", MIMEHtml},
		{"image/png", MIMEPng},
		{"audio/wav", MIMEWav},
		{"application/javascript", MIMEJavaScript},
	}

	for _, c := range cases {
		got, err := ParseMIME(c.input)
		if err != nil {
			t.Fatalf("ParseMIME(%q) unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Fatalf("ParseMIME(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestParseMIME_Invalid(t *testing.T) {
	if _, err := ParseMIME("not/a/mime"); err == nil {
		t.Fatal("expected error for invalid MIME string")
	}
}
