package mm

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/ast"
	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test -v -run TestEncodeBytes
//
// go test ./internal/mm -v -run TestEncodeBytes/nil_pointer
// go test ./internal/mm -v -run TestEncodesBytesRepresentation/medium_data
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
	xData := []byte("hello")

	testCases := []encodeBytesTestCase{
		{
			name:        "pointer_with_data",
			input:       &x,
			expectedOut: []byte{1},
			expectedErr: "",
		},
		{
			name:        "string_bytes",
			input:       []byte("hello world"),
			expectedOut: []byte("hello world"),
			expectedErr: "",
		},
		{
			name:        "zero_values",
			input:       []byte{0, 0, 0, 0},
			expectedOut: []byte{0, 0, 0, 0},
			expectedErr: "",
		},
		{
			name:        "binary_data",
			input:       []byte{255, 254, 253, 1, 2, 3},
			expectedOut: []byte{255, 254, 253, 1, 2, 3},
			expectedErr: "",
		},
		{
			name:        "pointer_with_content",
			input:       &xData,
			expectedOut: []byte("hello"),
			expectedErr: "",
		},
		{
			name:        "large_data",
			input:       bytes.Repeat([]byte{0xAB}, 256),
			expectedOut: bytes.Repeat([]byte{0xAB}, 256),
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs []byte
			bs, err := FromValue(tc.input, "type=bytes")

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
			fmt.Println("decoded:", Dump(bs2), jsonc.ToJSONC(bs2))
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
		n, _ := ValueToMM(data, "")
		_, _ = enc.Encode(n)
		buf.Reset()
	}
}

// TestEncodeBytesInStruct 测试结构体中的 bytes 字段
func TestEncodeBytesInStruct(t *testing.T) {
	type testStruct struct {
		Data     []byte `mm:"type=bytes;desc=数据"`
		Checksum []byte `mm:"type=bytes;desc=校验和"`
	}

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name: "with_data",
			input: testStruct{
				Data:     []byte("hello world"),
				Checksum: []byte{0x12, 0x34, 0x56, 0x78},
			},
		},
		{
			name: "binary_data",
			input: testStruct{
				Data:     []byte{255, 254, 253, 0, 1, 2, 3},
				Checksum: []byte{0xFF, 0x00},
			},
		},
		{
			name: "medium_data",
			input: testStruct{
				Data:     bytes.Repeat([]byte{0x42}, 256),
				Checksum: []byte{0xAB, 0xCD, 0xEF},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			_, err = Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			t.Logf("Encoded %s: %x", tc.name, bs)
		})
	}
}

// TestEncodeBytesArray 测试 bytes 数组
func TestEncodeBytesArray(t *testing.T) {
	testCases := []struct {
		name  string
		input [][]byte
	}{
		{
			name: "single_data",
			input: [][]byte{
				[]byte("test"),
			},
		},
		{
			name: "multiple_items",
			input: [][]byte{
				[]byte("first"),
				[]byte("second"),
				[]byte("third"),
			},
		},
		{
			name: "binary_array",
			input: [][]byte{
				{0xFF, 0xFE},
				{0x01, 0x02},
				{1, 2, 3, 4, 5},
			},
		},
		{
			name: "varied_sizes",
			input: [][]byte{
				[]byte{1},
				[]byte{1, 2, 3, 4, 5, 6, 7, 8},
				[]byte("hello"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			_, err = Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			t.Logf("Encoded %s (%d items): %x", tc.name, len(tc.input), bs)
		})
	}
}

// TestEncodesBytesNullable 测试可空的 bytes 字段
func TestEncodesBytesNullable(t *testing.T) {
	type testStruct struct {
		Required []byte `mm:"type=bytes;desc=必需"`
		Optional []byte `mm:"type=bytes;nullable;desc=可空"`
	}

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name: "with_data",
			input: testStruct{
				Required: []byte("required"),
				Optional: []byte("optional"),
			},
		},
		{
			name: "different_sizes",
			input: testStruct{
				Required: []byte("r"),
				Optional: []byte("optional data"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			_, err = Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			t.Logf("Encoded %s: %x", tc.name, bs)
		})
	}
}

// TestEncodesBytesRepresentation 测试字节数据表示
func TestEncodesBytesRepresentation(t *testing.T) {
	testCases := []struct {
		name           string
		input          []byte
		minEncodedSize int
	}{
		{
			name:           "single_byte",
			input:          []byte{0x42},
			minEncodedSize: 1,
		},
		{
			name:           "ascii_string",
			input:          []byte("hello"),
			minEncodedSize: 5,
		},
		{
			name:           "binary_data",
			input:          []byte{255, 254, 253, 252, 251},
			minEncodedSize: 5,
		},
		{
			name:           "medium_data",
			input:          bytes.Repeat([]byte{0x55}, 128),
			minEncodedSize: 128,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			if len(bs) < tc.minEncodedSize {
				t.Errorf("expected encoded size >= %d, got %d", tc.minEncodedSize, len(bs))
			}

			decoded, err := Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			value, ok := decoded.(*ast.Value)

			if !ok {
				t.Fatalf("expected *ast.Value, got %T", decoded)
			}

			decodedBytes := value.Data.([]byte)
			if !reflect.DeepEqual(decodedBytes, tc.input) {
				t.Errorf("value mismatch: expected %v, got %v", tc.input, decodedBytes)
			}

			t.Logf("%s: %d bytes -> encoded %x -> %d bytes", tc.name, len(tc.input), bs, len(decodedBytes))
		})
	}
}

// TestEncodesLargeBytes 测试大型字节数据
func TestEncodesLargeBytes(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{
			name:  "1KB",
			input: make([]byte, 1024),
		},
		{
			name:  "10KB",
			input: make([]byte, 10*1024),
		},
		{
			name:  "100KB",
			input: make([]byte, 100*1024),
		},
		{
			name:  "1MB",
			input: make([]byte, 1024*1024),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fill with some pattern
			for i := range tc.input {
				tc.input[i] = byte(i % 256)
			}

			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			decoded, err := Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			decodedBytes := decoded.(*ast.Value).Data.([]byte)
			if !reflect.DeepEqual(decodedBytes, tc.input) {
				t.Errorf("value mismatch for %s: size check failed", tc.name)
			}

			t.Logf("Encoded %s (%d bytes): %d bytes encoded", tc.name, len(tc.input), len(bs))
		})
	}
}

// TestEncodesUTF8Bytes 测试 UTF-8 编码的字节数据
func TestEncodesUTF8Bytes(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{
			name:  "ascii",
			input: []byte("hello world"),
		},
		{
			name:  "utf8_chinese",
			input: []byte("你好世界"),
		},
		{
			name:  "utf8_emoji",
			input: []byte("Hello 😀 World 🎉"),
		},
		{
			name:  "special_chars",
			input: []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
		},
		{
			name:  "mixed_content",
			input: []byte("Hello-世界-مرحبا-🌍"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			decoded, err := Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			value, ok := decoded.(*ast.Value)
			if !ok {
				t.Fatalf("expected *ast.Value, got %T", decoded)
			}

			decodedBytes := value.Data.([]byte)
			if !reflect.DeepEqual(decodedBytes, tc.input) {
				t.Errorf("value mismatch: expected %s, got %s", string(tc.input), string(decodedBytes))
			}

			t.Logf("Encoded %s: %s -> %x", tc.name, string(tc.input), bs)
		})
	}
}
