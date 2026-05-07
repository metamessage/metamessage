package mm

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/ast"
	"github.com/metamessage/metamessage/internal/jsonc"
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
			name:       "float32_positive_basic",
			input:      float32(0.123456),
			wantErr:    false,
			wantDecode: float32(0.123456),
		},
		{
			name:       "float64_positive_basic",
			input:      float64(9876.54321),
			wantErr:    false,
			wantDecode: float64(9876.54321),
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
			name:       "float32_smallest_positive",
			input:      float32(math.SmallestNonzeroFloat32),
			wantErr:    false,
			wantDecode: float32(math.SmallestNonzeroFloat32),
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
		{
			name:       "float32_positive_value",
			input:      float32(123.456),
			wantErr:    false,
			wantDecode: float32(123.456),
		},
		{
			name:       "float64_positive_value",
			input:      float64(654.321),
			wantErr:    false,
			wantDecode: float64(654.321),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			var bs []byte
			bs, err = FromValue(tc.input, "")
			fmt.Printf("bs %08b\n", bs)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Mismatch in error status: expected = %t, actual = %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}
				fmt.Println("decoded:", Dump(gotVal), jsonc.ToJSONC(gotVal))
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

// TestEncodeFloatInStruct 测试结构体中的 float 字段
func TestEncodeFloatInStruct(t *testing.T) {
	type testStruct struct {
		Price    float64 `mm:"type=f64;desc=价格"`
		Discount float32 `mm:"type=f32;desc=折扣"`
	}

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name: "basic_values",
			input: testStruct{
				Price:    99.99,
				Discount: 0.15,
			},
		},
		{
			name: "large_values",
			input: testStruct{
				Price:    1234567.89,
				Discount: 999999.99,
			},
		},
		{
			name: "small_values",
			input: testStruct{
				Price:    0.00001,
				Discount: 0.0001,
			},
		},
		{
			name: "positive_decimals",
			input: testStruct{
				Price:    123.45,
				Discount: 0.5,
			},
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

			obj := decoded.(*ast.Object)
			if len(obj.Fields) < 2 {
				t.Fatalf("expected at least 2 fields, got %d", len(obj.Fields))
			}

			t.Logf("Encoded %s: %x", tc.name, bs)
		})
	}
}

// TestEncodeFloatArray 测试 float 数组
func TestEncodeFloatArray(t *testing.T) {
	testCases := []struct {
		name  string
		input []float64
	}{
		{
			name:  "single_value",
			input: []float64{3.14159},
		},
		{
			name:  "multiple_values",
			input: []float64{1.1, 2.2, 3.3, 4.4, 5.5},
		},
		{
			name:  "negative_values",
			input: []float64{-1.5, -2.5, -3.5},
		},
		{
			name:  "positive_and_negative",
			input: []float64{-1.0, 1.0, -2.5, 2.5},
		},
		{
			name:  "scientific_notation",
			input: []float64{1e-5, 1e-3, 1e3, 1e5},
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

			arr := decoded.(*ast.Array)
			if len(arr.Items) != len(tc.input) {
				t.Errorf("array length mismatch: expected %d, got %d", len(tc.input), len(arr.Items))
			}

			for i, item := range arr.Items {
				val := item.(*ast.Value).Data.(float64)
				if !almostEqual(val, tc.input[i], 1e-10) {
					t.Errorf("array[%d] mismatch: expected %v, got %v", i, tc.input[i], val)
				}
			}

			t.Logf("Encoded %s (%d items): %x", tc.name, len(tc.input), bs)
		})
	}
}

// TestEncodeFloatNullable 测试可空的 float 字段
func TestEncodeFloatNullable(t *testing.T) {
	type testStruct struct {
		Required float64 `mm:"type=f64;desc=必需"`
		Optional float64 `mm:"type=f64;nullable;desc=可空"`
	}

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name: "both_set",
			input: testStruct{
				Required: 123.45,
				Optional: 67.89,
			},
		},
		{
			name: "small_precision",
			input: testStruct{
				Required: 0.123456789,
				Optional: 0.987654321,
			},
		},
		{
			name: "large_values",
			input: testStruct{
				Required: 1000000.001,
				Optional: 999999.999,
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

// TestEncodeFloatBoundary 测试 float 边界值
func TestEncodeFloatBoundary(t *testing.T) {
	testCases := []struct {
		name  string
		input float64
	}{
		{
			name:  "one",
			input: 1.0,
		},
		{
			name:  "pi",
			input: math.Pi,
		},
		{
			name:  "e",
			input: math.E,
		},
		{
			name:  "sqrt2",
			input: math.Sqrt2,
		},
		{
			name:  "negative_pi",
			input: -math.Pi,
		},
		{
			name:  "negative_e",
			input: -math.E,
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

			decodedVal := decoded.(*ast.Value).Data.(float64)
			if !almostEqual(decodedVal, tc.input, 1e-10) {
				t.Errorf("value mismatch: expected %v, got %v", tc.input, decodedVal)
			}

			t.Logf("%s: %v -> %x -> %v", tc.name, tc.input, bs, decodedVal)
		})
	}
}

// TestEncodeFloat32Precision 测试 float32 精度
func TestEncodeFloat32Precision(t *testing.T) {
	testCases := []struct {
		name  string
		input float32
	}{
		{
			name:  "quarter",
			input: 0.25,
		},
		{
			name:  "third",
			input: 1.0 / 3.0,
		},
		{
			name:  "pi_f32",
			input: float32(math.Pi),
		},
		{
			name:  "e_f32",
			input: float32(math.E),
		},
		{
			name:  "small_f32",
			input: float32(0.00123),
		},
		{
			name:  "large_f32",
			input: float32(1234567),
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

			decodedVal := decoded.(*ast.Value).Data.(float32)
			if !almostEqualF32(decodedVal, tc.input, 1e-5) {
				t.Errorf("value mismatch: expected %v, got %v", tc.input, decodedVal)
			}

			t.Logf("%s: %v -> %x -> %v", tc.name, tc.input, bs, decodedVal)
		})
	}
}

// TestEncodeFloatRepresentation 测试 float 表示
func TestEncodeFloatRepresentation(t *testing.T) {
	testCases := []struct {
		name  string
		input float64
	}{
		{
			name:  "integer_value",
			input: 42.5,
		},
		{
			name:  "simple_decimal",
			input: 3.5,
		},
		{
			name:  "complex_decimal",
			input: 0.123456789,
		},
		{
			name:  "very_small",
			input: 1e-50,
		},
		{
			name:  "large_value",
			input: 1e15,
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

			decodedVal := decoded.(*ast.Value).Data.(float64)
			if !almostEqual(decodedVal, tc.input, 1e-10) {
				t.Errorf("value mismatch: expected %v, got %v", tc.input, decodedVal)
			}

			t.Logf("%s: %v -> encoded -> %v", tc.name, tc.input, decodedVal)
		})
	}
}

// 辅助函数：比较浮点数是否在误差范围内相等
func almostEqual(a, b, epsilon float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 0) && math.IsInf(b, 0) {
		return math.Signbit(a) == math.Signbit(b)
	}
	return math.Abs(a-b) < epsilon
}

// 辅助函数：比较 float32 是否在误差范围内相等
func almostEqualF32(a, b float32, epsilon float32) bool {
	if math.IsNaN(float64(a)) && math.IsNaN(float64(b)) {
		return true
	}
	return math.Abs(float64(a-b)) < float64(epsilon)
}
