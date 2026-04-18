package mm

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/ast"
)

// go test ./internal/mm -v -run TestEncodeInteger
//
// go test ./internal/mm -v -run TestEncodeInteger/i8_nil_pointer
//
// go test ./internal/mm -bench=BenchmarkEncodeInteger -benchmem

func TestEncodeInteger(t *testing.T) {
	type testCase struct {
		name       string
		input      any
		wantErr    bool
		wantDecode any
	}

	x := 0
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
	xInt8 := int8(0)
	now := time.Now()
	testCases := []testCase{
		{
			name:       "time",
			input:      &now,
			wantErr:    false,
			wantDecode: int(0),
		},
		{
			name:       "i_nil_pointer",
			input:      i,
			wantErr:    false,
			wantDecode: int(0),
		},
		{
			name:       "i8_nil_pointer",
			input:      i8,
			wantErr:    false,
			wantDecode: int8(0),
		},
		{
			name:       "i16_nil_pointer",
			input:      i16,
			wantErr:    false,
			wantDecode: int16(0),
		},
		{
			name:       "i32_nil_pointer",
			input:      i32,
			wantErr:    false,
			wantDecode: int32(0),
		},
		{
			name:       "i64_nil_pointer",
			input:      i64,
			wantErr:    false,
			wantDecode: int64(0),
		},
		{
			name:       "u_nil_pointer",
			input:      u,
			wantErr:    false,
			wantDecode: uint(0),
		},
		{
			name:       "u8_nil_pointer",
			input:      u8,
			wantErr:    false,
			wantDecode: uint8(0),
		},
		{
			name:       "u16_nil_pointer",
			input:      u16,
			wantErr:    false,
			wantDecode: uint16(0),
		},
		{
			name:       "u32_nil_pointer",
			input:      u32,
			wantErr:    false,
			wantDecode: uint32(0),
		},
		{
			name:       "u64_nil_pointer",
			input:      u64,
			wantErr:    false,
			wantDecode: uint64(0),
		},
		{
			name:       "int zero pointer",
			input:      &x,
			wantErr:    false,
			wantDecode: int(0),
		},
		{
			name:       "int8 zero pointer",
			input:      &xInt8,
			wantErr:    false,
			wantDecode: int8(0),
		},

		{
			name:       "int zero",
			input:      0,
			wantErr:    false,
			wantDecode: int(0),
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
			name:       "int positive",
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
			name:       "uint zero",
			input:      uint(0),
			wantErr:    false,
			wantDecode: uint(0),
		},
		{
			name:       "uint positive",
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
			fmt.Println("res", bs)

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				fmt.Println("decoded:", jsonc.Json(gotVal), jsonc.ToString(gotVal))
				if !reflect.DeepEqual(gotVal.(*ast.Value).Data, tc.wantDecode) {
					t.Errorf("value mismatch: expected %v (%T), got %v (%T)",
						tc.wantDecode, tc.wantDecode, gotVal, gotVal)
				}
			}
		})
	}
}

func BenchmarkEncodeInteger(b *testing.B) {
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
