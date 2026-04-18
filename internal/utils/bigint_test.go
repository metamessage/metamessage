package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecodeBigInt_RoundTrip(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "zero", input: "0", want: "0"},
		{name: "single digit", input: "7", want: "7"},
		{name: "two digits", input: "12", want: "12"},
		{name: "three digits", input: "123", want: "123"},
		{name: "four digits", input: "1234", want: "1234"},
		{name: "five digits", input: "12345", want: "12345"},
		{name: "six digits", input: "123456", want: "123456"},
		{name: "leading zeros", input: "000", want: "000"},
		{name: "negative small", input: "-1", want: "-1"},
		{name: "negative multi", input: "-123456789", want: "-123456789"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			EncodeBigInt(buf, tc.input)
			encoded := buf.Bytes()
			got, err := DecodeBigInt(encoded, len(strings.TrimPrefix(tc.input, "-")))
			if err != nil {
				t.Fatalf("DecodeBigInt(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("DecodeBigInt(EncodeBigInt(%q)) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDecodeBigInt_Empty(t *testing.T) {
	got, err := DecodeBigInt(nil, 0)
	if err != nil {
		t.Fatalf("DecodeBigInt(nil) returned error: %v", err)
	}
	if got != "" {
		t.Fatalf("DecodeBigInt(nil) = %q, want empty string", got)
	}
}

func TestToBitsFromBits(t *testing.T) {
	cases := []struct {
		value int
		size  int
		bits  []int
	}{
		{0, 4, []int{0, 0, 0, 0}},
		{1, 4, []int{0, 0, 0, 1}},
		{5, 4, []int{0, 1, 0, 1}},
		{15, 4, []int{1, 1, 1, 1}},
		{127, 7, []int{1, 1, 1, 1, 1, 1, 1}},
		{512, 10, []int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for _, c := range cases {
		got := toBits(c.value, c.size)
		if !equalBits(got, c.bits) {
			t.Fatalf("toBits(%d, %d) = %v, want %v", c.value, c.size, got, c.bits)
		}

		back := fromBits(got)
		if back != c.value {
			t.Fatalf("fromBits(toBits(%d, %d)) = %d, want %d", c.value, c.size, back, c.value)
		}
	}
}

func equalBits(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestWriteAndBytesToBits(t *testing.T) {
	bits := []int{1, 0, 1, 1, 1, 0, 0, 1, 1}
	buf := bytes.NewBuffer(nil)
	writeBits(buf, bits)
	encoded := buf.Bytes()

	recovered := bytesToBits(encoded)
	if len(recovered) < len(bits) {
		t.Fatalf("bytesToBits returned too few bits: got %d want >= %d", len(recovered), len(bits))
	}

	for i, b := range bits {
		if recovered[i] != b {
			t.Fatalf("bytesToBits mismatch at position %d: got %d want %d", i, recovered[i], b)
		}
	}
}
