package mm

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/lizongying/meta-message/internal/jsonc"
	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

// go test -v -run TestEncodeBytes
//
// go test ./internal/mm -v -run TestEncodeBytes/nil_pointer
//
// go test -bench=BenchmarkEncodeBytes -benchmem

type encodeBytesTestCase struct {
	name        string
	input       any
	expectedOut any
	expectedErr string
}

func TestEncodeBytes(t *testing.T) {
	x := []byte{1}
	var y []byte
	y = nil
	testCases := []encodeBytesTestCase{
		{
			name:        "nil_pointer",
			input:       []byte(nil),
			expectedOut: []byte(nil),
			expectedErr: "",
		},
		{
			name:        "pointer",
			input:       &x,
			expectedOut: []byte{1},
			expectedErr: "",
		},
		{
			name:        "pointer2",
			input:       &y,
			expectedOut: []byte{},
			expectedErr: "",
		},
		{
			name:        "Ordinary byte slice",
			input:       []byte("hello world"),
			expectedOut: []byte("hello world"),
			expectedErr: "",
		},
		{
			name:        "Empty slice ([]byte{})",
			input:       []byte{},
			expectedOut: []byte{},
			expectedErr: "",
		},
		{
			name:        "slice",
			input:       []byte{0, 0, 0, 0},
			expectedOut: []byte{0, 0, 0, 0},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs []byte
			bs, err := FromStruct(tc.input, "type=bytes")

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %s, actual error: %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			bs2, _ := Decode(bs)
			fmt.Println("decoded:", jsonc.Json(bs2), jsonc.ToString(bs2))
			if !reflect.DeepEqual(bs2.(*ast.Value).Data, tc.expectedOut) {
				t.Errorf("Expected output: %v %T, actual output: %v %T", tc.expectedOut, tc.expectedOut, bs2.(*ast.Value).Data, bs2.(*ast.Value).Data)
			}
		})
	}
}

func BenchmarkEncodeBytes(b *testing.B) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	data := []byte("benchmark test data 1234567890")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bf := &bytes.Buffer{}
		enc.Reset(bf)
		n, _ := jsonc.StructToJSONC(data, "")
		_, _ = enc.Encode(n)
		buf.Reset()
	}
}
