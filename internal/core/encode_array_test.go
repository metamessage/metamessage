package core

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// go test ./internal/core -v
//
// go test ./internal/core -v -run TestEncodeArray
//
// go test ./internal/core -v -run TestEncodeArray/nil_byte_arr
//
// go test ./internal/core -bench=BenchmarkEncodeArray -benchtime=1000000x

type encodeArrayTestCase struct {
	name        string
	input       any
	tag         string
	expectedBuf any
	expectedErr string
}

func TestEncodeArray(t *testing.T) {
	testCases := []encodeArrayTestCase{
		{
			name:        "child_location_1",
			input:       ([]byte)(nil),
			tag:         "location=1",
			expectedBuf: nil,
			expectedErr: "not support location",
		},
		{
			name:        "nil_byte_vec",
			input:       ([]byte)(nil),
			expectedBuf: nil,
			expectedErr: "not allow empty",
		},
		{
			name:        "nil_byte_arr",
			input:       [3]byte{},
			expectedBuf: nil,
			expectedErr: "not allow empty",
		},
		{
			name:        "nil_ordinary_slice[]int)",
			input:       ([]int)(nil),
			expectedBuf: nil,
			expectedErr: "not allow empty",
		},
		{
			name:        "empty_byte_slice_([]byte{})",
			input:       []byte{},
			expectedBuf: nil,
			expectedErr: "not allow empty",
		},
		{
			name:        "empty_ordinary_slice_([]int{})",
			input:       []int{},
			expectedBuf: nil,
			expectedErr: "not allow empty",
		},
		{
			name:        "Non - empty byte slice ([]byte{1,2,3})",
			input:       []byte{0x01, 0x02, 0x03},
			expectedBuf: []byte{0x01, 0x02, 0x03},
			expectedErr: "",
		},
		{
			name:        "Non - empty byte array ([3]byte{1,2,3})",
			input:       [3]byte{0x01, 0x02, 0x03},
			expectedBuf: [3]byte{0x01, 0x02, 0x03},
			expectedErr: "",
		},
		{
			name:        "Non - empty ordinary slice ([]int{10,20,30})",
			input:       []int{10, 20, 30},
			expectedBuf: []int{10, 20, 30},
			expectedErr: "",
		},
		{
			name:  "Non - empty ordinary array ([2]int{5,6})",
			input: [2]int{5, 6},
			expectedBuf: `
        // mm: size=2; child_type=i
        [
        	5,
        	6,
        ]
			`,
			expectedErr: "",
		},
		// {
		// 	name:        "Non - array/slice type (int)",
		// 	input:       123,
		// 	expectedBuf: nil,
		// 	expectedErr: "encodeArrayOrSlice: not array/slice type",
		// },
		{
			name:        "12345",
			input:       [5]int{1, 2, 3, 4, 5},
			expectedBuf: [5]int{1, 2, 3, 4, 5},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := FromValue(tc.input, tc.tag)
			if err != nil {
				fmt.Println("error", err)
			}
			// fmt.Println("encoded", encoded)
			if tc.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error: %s, Actual error: %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, decodeErr := Decode(encoded)
			if decodeErr != nil {
				t.Fatalf("decode failed: %v", decodeErr)
			}
			// fmt.Println("decoded:", Dump(gotVal), jsonc.ToJSONC(gotVal))

			// if !reflect.DeepEqual(str, tc.expectedBuf) {
			// 	t.Errorf("Expected buffer: %v, actual buffer: %v", tc.expectedBuf, str)
			// }
		})
	}
}

func BenchmarkEncodeArray(b *testing.B) {
	e := NewEncoder(nil)
	data := []byte("benchmark test data 1234567890abcdefghijklmnopqrstuvwxyz")
	val := reflect.ValueOf(data)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, _ := ValueToNode(val, "")
		_, _ = e.Encode(n)
	}
}
