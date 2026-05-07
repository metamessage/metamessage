package mm

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/ast"
	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test ./internal/mm -v -run TestEncodeBool
//
// go test ./internal/mm -v -run TestEncodeBool/bool_nil_pointer
//
// go test ./internal/mm -bench=BenchmarkEncodeBool -benchmem

func TestEncodeBool(t *testing.T) {
	type testCase struct {
		name       string
		input      any
		wantErr    bool
		wantDecode any
	}

	x := false
	xTrue := true
	var y *bool
	testCases := []testCase{
		{
			name:       "bool_nil_pointer",
			input:      y,
			wantErr:    false,
			wantDecode: bool(false),
		},
		{
			name:       "bool true",
			input:      true,
			wantErr:    false,
			wantDecode: bool(true),
		},
		{
			name:       "bool false",
			input:      x,
			wantErr:    false,
			wantDecode: bool(false),
		},
		{
			name:       "bool true pointer",
			input:      &xTrue,
			wantErr:    false,
			wantDecode: bool(true),
		},
		{
			name:       "bool false pointer",
			input:      &x,
			wantErr:    false,
			wantDecode: bool(false),
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
			bs, err := FromValue(tc.input, "")
			fmt.Println("res", bs)

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if !tc.wantErr {
				gotVal, decodeErr := Decode(bs)
				if decodeErr != nil {
					t.Fatalf("decode failed: %v", decodeErr)
				}

				fmt.Println("decoded:", Dump(gotVal), jsonc.ToJSONC(gotVal))
				if !reflect.DeepEqual(gotVal.(*ast.Value).Data, tc.wantDecode) {
					t.Errorf("value mismatch: expected %v (%T), got %v (%T)",
						tc.wantDecode, tc.wantDecode, gotVal, gotVal)
				}
			}
		})
	}
}

func BenchmarkEncodeBool(b *testing.B) {
	testInputs := []any{
		true,
		false,
	}
	var out bytes.Buffer
	enc := NewEncoder(&out)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var o bytes.Buffer
		enc.Reset(&o)
		n, _ := ValueToMM(testInputs[i%len(testInputs)], "")
		_, _ = enc.Encode(n)
	}
}

// TestEncodeBoolInStruct 测试结构体中的 bool 字段
func TestEncodeBoolInStruct(t *testing.T) {
	type testStruct struct {
		Active   bool `mm:"type=bool;desc=是否激活"`
		Verified bool `mm:"type=bool;desc=是否已验证"`
	}

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name:  "both_true",
			input: testStruct{Active: true, Verified: true},
		},
		{
			name:  "both_false",
			input: testStruct{Active: false, Verified: false},
		},
		{
			name:  "mixed",
			input: testStruct{Active: true, Verified: false},
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
			fields := obj.Fields
			if len(fields) < 2 {
				t.Fatalf("expected at least 2 fields, got %d", len(fields))
			}

			t.Logf("Encoded %s: %x", tc.name, bs)
		})
	}
}

// TestEncodeBoolArray 测试 bool 数组
func TestEncodeBoolArray(t *testing.T) {
	testCases := []struct {
		name  string
		input []bool
	}{
		{
			name:  "single_true",
			input: []bool{true},
		},
		{
			name:  "single_false",
			input: []bool{false},
		},
		{
			name:  "mixed_array",
			input: []bool{true, false, true, false, true},
		},
		{
			name:  "all_true",
			input: []bool{true, true, true, true},
		},
		{
			name:  "all_false",
			input: []bool{false, false, false, false},
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
				val := item.(*ast.Value).Data.(bool)
				if val != tc.input[i] {
					t.Errorf("array[%d] mismatch: expected %v, got %v", i, tc.input[i], val)
				}
			}

			t.Logf("Encoded %s (%d items): %x", tc.name, len(tc.input), bs)
		})
	}
}

// TestEncodeBoolNullable 测试可空的 bool 字段
func TestEncodeBoolNullable(t *testing.T) {
	type testStruct struct {
		Active    bool  `mm:"type=bool;desc=必需"`
		Suspended *bool `mm:"type=bool;nullable;desc=可空"`
	}

	falseVal := false
	trueVal := true

	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name:  "nullable_nil",
			input: testStruct{Active: true, Suspended: nil},
		},
		{
			name:  "nullable_true",
			input: testStruct{Active: true, Suspended: &trueVal},
		},
		{
			name:  "nullable_false",
			input: testStruct{Active: false, Suspended: &falseVal},
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

// TestEncodeBoolByteRepresentation 测试字节表示
func TestEncodeBoolByteRepresentation(t *testing.T) {
	testCases := []struct {
		name  string
		input bool
	}{
		{
			name:  "true_value",
			input: true,
		},
		{
			name:  "false_value",
			input: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := FromValue(tc.input, "")
			if err != nil {
				t.Fatalf("encode failed: %v", err)
			}

			// bool 至少应该有 1 字节的编码
			if len(bs) < 1 {
				t.Errorf("expected encoded bool to have at least 1 byte, got %d", len(bs))
			}

			// 验证能否解码回去
			decoded, err := Decode(bs)
			if err != nil {
				t.Fatalf("decode failed: %v", err)
			}

			decodedBool := decoded.(*ast.Value).Data.(bool)
			if decodedBool != tc.input {
				t.Errorf("value mismatch: expected %v, got %v", tc.input, decodedBool)
			}

			t.Logf("%s: %v -> %x -> %v", tc.name, tc.input, bs, decodedBool)
		})
	}
}
