package mm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test ./internal/mm -v
//
// go test ./internal/mm -v -run TestEncodeArray
//
// go test ./internal/mm -v -run TestEncodeArray/child_location_1
//
// go test ./internal/mm -bench=BenchmarkEncodeArray -benchtime=1000000x

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
			expectedBuf: ([]byte)(nil),
			expectedErr: "",
		},
		{
			name:        "nil byte slice",
			input:       ([]byte)(nil),
			expectedBuf: ([]byte)(nil),
			expectedErr: "",
		},
		{
			name:        "nil ordinary slice[]int)",
			input:       ([]int)(nil),
			expectedBuf: ([]int)(nil),
			expectedErr: "",
		},
		{
			name:        "Empty byte slice ([]byte{})",
			input:       []byte{},
			expectedBuf: []byte{},
			expectedErr: "",
		},
		{
			name:        "Empty ordinary slice ([]int{})",
			input:       []int{},
			expectedBuf: []int{},
			expectedErr: "",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := FromValue(tc.input, tc.tag)

			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %s, Actual error: %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			gotVal, decodeErr := Decode(encoded)
			if decodeErr != nil {
				t.Fatalf("decode failed: %v", decodeErr)
			}
			fmt.Println("decoded:", Dump(gotVal), jsonc.ToJSONC(gotVal))

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
		n, _ := ValueToMM(val, "")
		_, _ = e.Encode(n)
	}
}
