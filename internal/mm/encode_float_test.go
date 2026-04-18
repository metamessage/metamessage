package mm

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/ast"
)

// go test ./internal/mm -v -run TestEncodeFloat
//
// go test ./internal/mm -v -run TestEncodeFloat/float32_nil_pointer
//
// go test ./internal/mm -v -run TestEncodeFloat
//
// go test ./internal/mm -v -run "TestEncodeFloat.*"

// go test ./internal/mm -v -run TestEncode_Float_Others

func TestEncodeFloat(t *testing.T) {
	type testCase struct {
		name       string
		input      any
		wantErr    bool
		wantDecode any
	}

	var f32 *float32
	var f64 *float64
	testCases := []testCase{
		{
			name:       "float32_inf",
			input:      math.Inf(1),
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "float32_nil_pointer",
			input:      f32,
			wantErr:    false,
			wantDecode: float32(0),
		},
		{
			name:       "float64_nil_pointer",
			input:      f64,
			wantErr:    false,
			wantDecode: float64(0),
		},
		{
			name:       "float32_00",
			input:      float32(0),
			wantErr:    false,
			wantDecode: float32(0),
		},
		{
			name:       "float64_00",
			input:      float64(0),
			wantErr:    false,
			wantDecode: float64(0),
		},
		{
			name:       "float32_positive_basic",
			input:      float32(0.123456),
			wantErr:    false,
			wantDecode: float32(0.123456),
		},
		{
			name:       "float32_negative_boundary",
			input:      float32(-math.MaxFloat32),
			wantErr:    false,
			wantDecode: float32(-math.MaxFloat32),
		},
		{
			name:       "float32_negative_inf",
			input:      float32(math.Inf(-1)),
			wantErr:    false,
			wantDecode: float32(math.Inf(-1)),
		},
		{
			name:       "float32_nan",
			input:      float32(math.NaN()),
			wantErr:    true,
			wantDecode: nil,
		},
		{
			name:       "float64_nan",
			input:      float64(math.NaN()),
			wantErr:    true,
			wantDecode: nil,
		},

		{
			name:       "float64_positive_basic",
			input:      float64(9876.54321),
			wantErr:    false,
			wantDecode: float64(9876.54321),
		},
		{
			name:       "float64_zero",
			input:      float64(0.0),
			wantErr:    false,
			wantDecode: float64(0.0),
		},

		{
			name:       "float32_smallest_positive",
			input:      float32(math.SmallestNonzeroFloat32),
			wantErr:    false,
			wantDecode: float32(math.SmallestNonzeroFloat32),
		},
		{
			name:       "float64_smallest_positive",
			input:      float64(math.SmallestNonzeroFloat64),
			wantErr:    true,
			wantDecode: float64(math.SmallestNonzeroFloat64),
		},
		{
			name:       "float32_epsilon",
			input:      float32(math.Nextafter32(1.0, 2.0) - 1.0),
			wantErr:    false,
			wantDecode: float32(math.Nextafter32(1.0, 2.0) - 1.0),
		},
		{
			name:       "float64_epsilon",
			input:      float64(math.Nextafter(1.0, 2.0) - 1.0),
			wantErr:    false,
			wantDecode: float64(math.Nextafter(1.0, 2.0) - 1.0),
		},

		// {
		// 	name:       "input_non_float",
		// 	input:      "not a float",
		// 	wantErr:    false,
		// 	wantDecode: "not a float",
		// },
		// {
		// 	name:       "input_nil",
		// 	input:      nil,
		// 	wantErr:    true,
		// 	wantDecode: nil,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			var bs []byte
			bs, err = FromStruct(tc.input, "")
			fmt.Printf("bs %08b\n", bs)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Mismatch in error status: expected = %t, actual = %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}
				fmt.Println("decoded:", jsonc.Json(gotVal), jsonc.ToString(gotVal))
				if !reflect.DeepEqual(gotVal.(*ast.Value).Data, tc.wantDecode) {
					t.Errorf("value mismatch: expected %v (%T), got %v (%T)",
						tc.wantDecode, tc.wantDecode, gotVal.(*ast.Value).Data, gotVal.(*ast.Value).Data)
				}
			}
		})
	}
}

func TestEncodeFloat_Others(t *testing.T) {
	// var zero float32 = 0.0
	// var negZero float32 = -0.0

	// // 1. 比较结果
	// fmt.Println(zero == negZero) // ✅ 输出：true （数值相等）

	// // 2. 底层二进制完全不同！
	// fmt.Printf("0.0  : %x\n", math.Float32bits(zero))    // 00000000
	// fmt.Printf("-0.0 : %x\n", math.Float32bits(negZero)) // 80000000

	negZero := math.Copysign(0.0, -1.0)

	fmt.Println(negZero) // 输出？
	fmt.Printf("%f\n", negZero)
	fmt.Printf("%g\n", negZero)
	fmt.Printf("%+f\n", negZero)
}
