package mm

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/ast"
)

func TestEncodeInt(t *testing.T) {
	type testCase struct {
		name       string
		input      any
		wantErr    bool
		wantDecode any
	}

	var i *int
	var i8 *int8
	var i16 *int16
	var i32 *int32
	var i64 *int64
	var u *uint
	var u8 *uint8
	var u16 *uint16
	var u32 *uint32
	var u64 *uint64
	xInt8 := int8(1)
	xInt := 1

	testCases := []testCase{
		{
			name:       "i_nil_pointer",
			input:      i,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "i8_nil_pointer",
			input:      i8,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "i16_nil_pointer",
			input:      i16,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "i32_nil_pointer",
			input:      i32,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "i64_nil_pointer",
			input:      i64,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "u_nil_pointer",
			input:      u,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "u8_nil_pointer",
			input:      u8,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "u16_nil_pointer",
			input:      u16,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "u32_nil_pointer",
			input:      u32,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "u64_nil_pointer",
			input:      u64,
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "int one pointer",
			input:      &xInt,
			wantErr:    false,
			wantDecode: int(1),
		},
		{
			name:       "int8 one pointer",
			input:      &xInt8,
			wantErr:    false,
			wantDecode: int8(1),
		},
		{
			name:       "int positive",
			input:      1,
			wantErr:    false,
			wantDecode: int(1),
		},
		{
			name:       "int 23",
			input:      23,
			wantErr:    false,
			wantDecode: int(23),
		},
		{
			name:       "int 24",
			input:      24,
			wantErr:    false,
			wantDecode: int(24),
		},
		{
			name:       "int large",
			input:      123456,
			wantErr:    false,
			wantDecode: int(123456),
		},
		{
			name:       "int negative",
			input:      -7890,
			wantErr:    false,
			wantDecode: int(-7890),
		},
		{
			name:       "int8 min",
			input:      int8(-128),
			wantErr:    false,
			wantDecode: int8(-128),
		},
		{
			name:       "int8 max",
			input:      int8(127),
			wantErr:    false,
			wantDecode: int8(127),
		},
		{
			name:       "int16 min",
			input:      int16(-32768),
			wantErr:    false,
			wantDecode: int16(-32768),
		},
		{
			name:       "int16 max",
			input:      int16(32767),
			wantErr:    false,
			wantDecode: int16(32767),
		},
		{
			name:       "int32 min",
			input:      int32(-2147483648),
			wantErr:    false,
			wantDecode: int32(-2147483648),
		},
		{
			name:       "int32 max",
			input:      int32(2147483647),
			wantErr:    false,
			wantDecode: int32(2147483647),
		},
		{
			name:       "int64 min",
			input:      int64(-9223372036854775808),
			wantErr:    false,
			wantDecode: int64(-9223372036854775808),
		},
		{
			name:       "int64 max",
			input:      int64(9223372036854775807),
			wantErr:    false,
			wantDecode: int64(9223372036854775807),
		},
		{
			name:       "uint positive",
			input:      uint(1),
			wantErr:    false,
			wantDecode: uint(1),
		},
		{
			name:       "uint large",
			input:      uint(987654),
			wantErr:    false,
			wantDecode: uint(987654),
		},
		{
			name:       "uint8 max",
			input:      uint8(255),
			wantErr:    false,
			wantDecode: uint8(255),
		},
		{
			name:       "uint16 max",
			input:      uint16(65535),
			wantErr:    false,
			wantDecode: uint16(65535),
		},
		{
			name:       "uint32 max",
			input:      uint32(4294967295),
			wantErr:    false,
			wantDecode: uint32(4294967295),
		},
		{
			name:       "uint64 max",
			input:      uint64(18446744073709551615),
			wantErr:    false,
			wantDecode: uint64(18446744073709551615),
		},
		{
			name:       "uintptr",
			input:      uintptr(123456789),
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "nil input",
			input:      nil,
			wantErr:    true,
			wantDecode: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !reflect.DeepEqual(gotVal.(*ast.Value).Data, tc.wantDecode) {
					t.Errorf("value mismatch: expected %v (%T), got %v (%T)",
						tc.wantDecode, tc.wantDecode, gotVal, gotVal)
				}
			}
		})
	}
}

func TestEncodeIntInStruct(t *testing.T) {
	type testStruct struct {
		IntField    int
		Int8Field   int8
		Int16Field  int16
		Int32Field  int32
		Int64Field  int64
		UintField   uint
		Uint8Field  uint8
		Uint16Field uint16
		Uint32Field uint32
		Uint64Field uint64
	}

	type testCase struct {
		name      string
		input     testStruct
		wantErr   bool
		wantCheck func(any) bool
	}

	testCases := []testCase{
		{
			name: "all positive values",
			input: testStruct{
				IntField:    123,
				Int8Field:   127,
				Int16Field:  32767,
				Int32Field:  2147483647,
				Int64Field:  9223372036854775807,
				UintField:   123,
				Uint8Field:  255,
				Uint16Field: 65535,
				Uint32Field: 4294967295,
				Uint64Field: 18446744073709551615,
			},
			wantErr: false,
			wantCheck: func(val any) bool {
				decodedData := val.(*ast.Value).Data.(map[string]any)
				return decodedData["IntField"] == int(123) &&
					decodedData["Int8Field"] == int8(127)
			},
		},
		{
			name: "negative integers",
			input: testStruct{
				IntField:    -456,
				Int8Field:   -128,
				Int16Field:  -32768,
				Int32Field:  -2147483648,
				Int64Field:  -9223372036854775808,
				UintField:   100,
				Uint8Field:  100,
				Uint16Field: 100,
				Uint32Field: 100,
				Uint64Field: 100,
			},
			wantErr: false,
			wantCheck: func(val any) bool {
				decodedData := val.(*ast.Value).Data.(map[string]any)
				return decodedData["IntField"] == int(-456) &&
					decodedData["Int8Field"] == int8(-128)
			},
		},
		{
			name: "mixed values",
			input: testStruct{
				IntField:    789,
				Int8Field:   64,
				Int16Field:  16000,
				Int32Field:  1000000000,
				Int64Field:  5000000000000,
				UintField:   999,
				Uint8Field:  200,
				Uint16Field: 50000,
				Uint32Field: 2000000000,
				Uint64Field: 9000000000000000000,
			},
			wantErr: false,
			wantCheck: func(val any) bool {
				decodedData := val.(*ast.Value).Data.(map[string]any)
				return decodedData["IntField"] == int(789)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !tc.wantCheck(gotVal) {
					t.Errorf("value check failed for decoded value: %v", gotVal)
				}
			}
		})
	}
}

func TestEncodeIntArray(t *testing.T) {
	type testCase struct {
		name      string
		input     any
		wantErr   bool
		wantCheck func(any) bool
	}

	testCases := []testCase{
		{
			name:    "int array",
			input:   []int{1, 2, 3, 100, -50},
			wantErr: false,
			wantCheck: func(val any) bool {
				arr := val.(*ast.Value).Data.([]any)
				return len(arr) == 5 && arr[0] == int(1)
			},
		},
		{
			name:    "int8 array",
			input:   []int8{10, 20, 127, -128, 1},
			wantErr: false,
			wantCheck: func(val any) bool {
				arr := val.(*ast.Value).Data.([]any)
				return len(arr) == 5
			},
		},
		{
			name:    "uint64 array",
			input:   []uint64{1000, 2000, 3000, 9223372036854775807},
			wantErr: false,
			wantCheck: func(val any) bool {
				arr := val.(*ast.Value).Data.([]any)
				return len(arr) == 4
			},
		},
		{
			name:    "mixed int sizes array",
			input:   []int32{100, -200, 32767, -32768, 1},
			wantErr: false,
			wantCheck: func(val any) bool {
				arr := val.(*ast.Value).Data.([]any)
				return len(arr) == 5
			},
		},
		{
			name:    "uint array with large values",
			input:   []uint{1, 100, 1000, 10000},
			wantErr: false,
			wantCheck: func(val any) bool {
				arr := val.(*ast.Value).Data.([]any)
				return len(arr) == 4
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !tc.wantCheck(gotVal) {
					t.Errorf("array check failed for decoded value: %v", gotVal)
				}
			}
		})
	}
}

func TestEncodeIntNullable(t *testing.T) {
	type testCase struct {
		name      string
		input     any
		wantErr   bool
		wantCheck func(any) bool
	}

	val1 := int(1)
	val8 := int8(10)
	val64 := int64(999999)

	testCases := []testCase{
		{
			name:    "nullable int pointer",
			input:   &val1,
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int(1)
			},
		},
		{
			name:    "nullable int8 pointer",
			input:   &val8,
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int8(10)
			},
		},
		{
			name:    "nullable int64 pointer",
			input:   &val64,
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int64(999999)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !tc.wantCheck(gotVal) {
					t.Errorf("value mismatch: %v", gotVal)
				}
			}
		})
	}
}

func TestEncodeIntBoundary(t *testing.T) {
	type testCase struct {
		name      string
		input     any
		wantErr   bool
		wantCheck func(any) bool
	}

	testCases := []testCase{
		{
			name:    "int8 min boundary",
			input:   int8(-128),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int8(-128)
			},
		},
		{
			name:    "int8 max boundary",
			input:   int8(127),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int8(127)
			},
		},
		{
			name:    "int16 min boundary",
			input:   int16(-32768),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int16(-32768)
			},
		},
		{
			name:    "int16 max boundary",
			input:   int16(32767),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int16(32767)
			},
		},
		{
			name:    "int32 min boundary",
			input:   int32(-2147483648),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int32(-2147483648)
			},
		},
		{
			name:    "int64 max boundary",
			input:   int64(9223372036854775807),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int64(9223372036854775807)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !tc.wantCheck(gotVal) {
					t.Errorf("boundary check failed: %v", gotVal)
				}
			}
		})
	}
}

func TestEncodeIntRepresentation(t *testing.T) {
	type testCase struct {
		name      string
		input     any
		wantErr   bool
		wantCheck func(any) bool
	}

	testCases := []testCase{
		{
			name:    "int64 large positive",
			input:   int64(9000000000000000000),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int64(9000000000000000000)
			},
		},
		{
			name:    "int64 large negative",
			input:   int64(-9000000000000000000),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int64(-9000000000000000000)
			},
		},
		{
			name:    "uint64 very large",
			input:   uint64(18000000000000000000),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == uint64(18000000000000000000)
			},
		},
		{
			name:    "int32 boundary positive",
			input:   int32(2147483647),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == int32(2147483647)
			},
		},
		{
			name:    "uint32 boundary max",
			input:   uint32(4294967295),
			wantErr: false,
			wantCheck: func(v any) bool {
				return v.(*ast.Value).Data == uint32(4294967295)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromStruct(tc.input, "")

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				if !tc.wantCheck(gotVal) {
					t.Errorf("representation check failed: %v", gotVal)
				}
			}
		})
	}
}

func BenchmarkEncodeInt(b *testing.B) {
	testInputs := []any{
		int64(123456),
		uint64(987654),
		int8(-128),
		uint32(4294967295),
	}
	var out bytes.Buffer
	enc := NewEncoder(&out)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var o bytes.Buffer
		enc.Reset(&o)
		n, _ := jsonc.StructToJSONC(testInputs[i%len(testInputs)], "")
		_, _ = enc.Encode(n)
	}
}
