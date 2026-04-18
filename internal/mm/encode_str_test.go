package mm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/lizongying/meta-message/internal/jsonc"
	"github.com/lizongying/meta-message/internal/jsonc/ast"
)

// go test ./internal/mm -v -run TestEncodeStr
//
// go test ./internal/mm -v -run TestEncodeStr/location
//
// go test ./internal/mm -bench=BenchmarkEncodeStr -benchmem
// go test ./internal/mm -test.fullpath=true -benchmem -run=^$ -bench ^BenchmarkEncodeStr_MM$ -benchtime=1000000x
// go test ./internal/mm -test.fullpath=true -benchmem -run=^$ -bench ^BenchmarkEncodeStr_JSON$ -benchtime=1000000x

type encodeStrTestCase struct {
	name        string
	input       any
	tag         string
	expectedOut any
	expectedErr string
}

func TestEncodeStr(t *testing.T) {
	x := "hello world"
	var y string
	y = ""
	var s *string
	var bi *big.Int

	var ip *net.IP
	var Url *url.URL

	testCases := []encodeStrTestCase{
		{
			name:        "str_location__10",
			input:       "",
			tag:         "location=-10",
			expectedOut: nil,
			expectedErr: "",
		},
		{
			name:        "location__10",
			input:       time.Now(),
			tag:         "location=-10",
			expectedOut: nil,
			expectedErr: "",
		},
		{
			name:        "location__12",
			input:       time.Now(),
			tag:         "location=-12",
			expectedOut: nil,
			expectedErr: "",
		},
		{
			name:        "location__13",
			input:       time.Now(),
			tag:         "location=-13",
			expectedOut: nil,
			expectedErr: "parse tag failed",
		},
		{
			name:        "location_14",
			input:       time.Now(),
			tag:         "location=14",
			expectedOut: nil,
			expectedErr: "",
		},
		{
			name:        "location_15",
			input:       time.Now(),
			tag:         "location=15",
			expectedOut: nil,
			expectedErr: "parse tag failed",
		},

		{
			name:        "ip_nil_pointer",
			input:       ip,
			expectedOut: net.IP{},
			expectedErr: "",
		},
		{
			name:        "ip_v4_zero",
			input:       net.IPv4zero,
			expectedOut: net.IPv4zero,
			expectedErr: "",
		},
		{
			name:        "ip_v4_allrouter",
			input:       net.IPv4allrouter,
			expectedOut: net.IPv4allrouter,
			expectedErr: "",
		},
		{
			name:        "ip_v6_zero",
			input:       net.IPv6zero,
			expectedOut: net.IPv6zero,
			expectedErr: "",
		},
		{
			name:        "ip_v6_linklocalallnodes",
			input:       net.IPv6linklocalallnodes,
			expectedOut: net.IPv6linklocalallnodes,
			expectedErr: "",
		},
		{
			name:        "ip_v6_all",
			input:       net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expectedOut: net.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			expectedErr: "",
		},
		{
			name:        "url_nil_pointer",
			input:       Url,
			expectedOut: url.URL{},
			expectedErr: "",
		},
		{
			name:        "url_zero",
			input:       url.URL{},
			expectedOut: url.URL{},
			expectedErr: "",
		},
		{
			name:        "nil_pointer",
			input:       s,
			expectedOut: "",
			expectedErr: "",
		},
		{
			name:        "bi_nil_pointer",
			input:       bi,
			expectedOut: big.Int{},
			expectedErr: "",
		},
		{
			name:        "bi_long",
			input:       big.NewInt(234343434343423232),
			expectedOut: big.NewInt(234343434343423232),
			expectedErr: "",
		},
		{
			name:        "enum",
			input:       "SECOND",
			tag:         "enum=FIRST|SECOND|THIRD",
			expectedOut: "SECOND",
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
			bs, err := FromStruct(tc.input, tc.tag)
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
			fmt.Println("decoded:", jsonc.Json(bs2))
			fmt.Println("jsonc:", jsonc.ToString(bs2))
			v, ok := bs2.(*ast.Value)
			if ok {
				if !reflect.DeepEqual(v.Data, tc.expectedOut) {
					t.Errorf("Expected output: %v %T, actual output: %v %T", tc.expectedOut, tc.expectedOut, bs2.(*ast.Value).Data, bs2.(*ast.Value).Data)
				}
			}
		})
	}
}

func BenchmarkEncodeStr_MM(b *testing.B) {
	e := NewEncoder(nil)
	data := "benchmark test data 1234567890"
	n, err := jsonc.StructToJSONC(data, "")
	out, _ := e.Encode(n)
	b.Log("out", len(out), err)
	b.ResetTimer()

	for b.Loop() {
		n, _ := jsonc.StructToJSONC(data, "")
		_, _ = e.Encode(n)
	}
}

func BenchmarkEncodeStr_JSON(b *testing.B) {
	data := "benchmark test data 1234567890"
	out, _ := json.Marshal(data)
	b.Log("out", len(out))
	b.ResetTimer()

	for b.Loop() {
		_, _ = json.Marshal(data)
	}
}
