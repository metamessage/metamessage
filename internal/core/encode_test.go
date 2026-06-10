package core

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test ./internal/core -v -run TestEncode
//
// go test ./internal/core -v -run "^TestEncode/nil$"
func TestEncode(t *testing.T) {
	testCases := []encodeStrTestCase{
		{
			name:        "nil",
			input:       nil,
			tag:         "",
			expectedOut: nil,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs []byte
			bs, err := FromValue(tc.input, tc.tag)
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Printf("encode len: %d\n", len(bs))

			if tc.expectedErr != "" {
				if err == nil {
					t.Fatalf("Expected error: %s, actual error: nil", tc.expectedErr)
				}

				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error contains: %q\nActual error: %q", tc.expectedErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			bs2, err := Decode(bs)
			if err != nil {
				fmt.Println("decoded err", err)
			}
			// fmt.Println("decoded:", Dump(bs2))
			fmt.Println("jsonc:", jsonc.ToJSONC(bs2))
			v, ok := bs2.(*ir.NodeScalar)
			if ok {
				if dataTime, isTime := v.Data.(time.Time); isTime {
					if expectedTime, ok := tc.expectedOut.(time.Time); ok {
						if !dataTime.Equal(expectedTime) {
							t.Errorf("Expected output: %v, actual output: %v", tc.expectedOut, v.Data)
						}
						return
					}
				}
				if !reflect.DeepEqual(v.Data, tc.expectedOut) {
					t.Errorf("Expected output: %v %T, actual output: %v %T", tc.expectedOut, tc.expectedOut, bs2.(*ir.NodeScalar).Data, bs2.(*ir.NodeScalar).Data)
				}
			}
		})
	}
}
