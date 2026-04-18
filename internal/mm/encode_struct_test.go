package mm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test ./internal/mm -v -run TestEncodeStruct
//
// go test ./internal/mm -v -run TestEncodeStruct/invalid_uuid
//
// go test ./internal/mm -test.fullpath=true -benchmem -run=^$ -bench ^BenchmarkEncodeStruct_MM$ -benchtime=1000000x
// go test ./internal/mm -test.fullpath=true -benchmem -run=^$ -bench ^BenchmarkEncodeStruct_JSON$ -benchtime=1000000x

type User struct {
	ID      int    `mm:"min=2"`
	Name    string `mm:"allow_empty"`
	Email   string `mm:"allow_empty"`
	UUID    string `mm:"type=uuid; allow_empty"`
	Decimal string `mm:"type=decimal"`
	Age     uint8  `mm:"allow_empty"`
	Profile struct {
		Address string `mm:"allow_empty"`
		Phone   []byte `mm:"type=bytes; allow_empty"`
	} ``
	Nullable *string ``
}

type encodeTestCase struct {
	name    string
	input   User
	tag     string
	wantErr bool
}

func TestEncodeStruct(t *testing.T) {
	var nilStr *string

	testCases := []encodeTestCase{
		{
			name: "invalid_uuid",
			input: User{
				ID:   1,
				UUID: "121",
			},
			wantErr: true,
		},
		{
			name: "uuid",
			input: User{
				ID:      2,
				UUID:    "d46372a1-5b9e-4f8c-9a2d-7e1b3c5d7f9a",
				Decimal: "0.0",
			},
			wantErr: false,
		},
		{
			name: "invalid_decimal",
			input: User{
				ID:      1,
				UUID:    "d46372a1-5b9e-4f8c-9a2d-7e1b3c5d7f9a",
				Decimal: "aaa12123223",
			},
			wantErr: true,
		},
		{
			name: "ok_decimal",
			input: User{
				UUID:    "d46372a1-5b9e-4f8c-9a2d-7e1b3c5d7f9a",
				Decimal: "12123223",
			},
			wantErr: true,
		},
		{
			name: "bytes",
			input: User{
				ID:      2,
				Name:    "",
				Age:     30,
				Decimal: "0.0",
				Profile: struct {
					Address string `mm:"allow_empty"`
					Phone   []byte `mm:"type=bytes; allow_empty"`
				}{
					Address: "",
					Phone:   []byte{},
				},
				Nullable: nilStr,
			},
			wantErr: false,
		},
		{
			name: "email",
			input: User{
				ID:      3,
				Name:    "",
				Email:   "lisi@test.com",
				Decimal: "0.0",
				Age:     28,
			},
			wantErr: false,
		},
		{
			name: "id0",
			input: User{
				ID:      2,
				Name:    "",
				Decimal: "0.0",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, err := FromStruct(tc.input, tc.tag)

			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch: expected err=%t, got err=%v", tc.wantErr, err)
			}

			if err != nil {
				fmt.Println("output err:", err)
				return
			}

			bs2, _ := Decode(gotBytes)
			// fmt.Println("decoded node:", jsonc.Json(bs2))
			fmt.Println("decoded jsonc:", jsonc.ToString(bs2))
		})
	}
}

func BenchmarkEncodeStruct_MM(b *testing.B) {
	e := NewEncoder(nil)
	data := User{
		ID:   100,
		Name: "",
		Age:  25,
		Profile: struct {
			Address string `mm:"allow_empty"`
			Phone   []byte `mm:"type=bytes; allow_empty"`
		}{Address: ""},
		Decimal: "0.0",
	}
	n, err := jsonc.StructToJSONC(data, "")
	out, _ := e.Encode(n)
	b.Log("out", len(out), err)
	b.ResetTimer()

	for b.Loop() {
		n, _ := jsonc.StructToJSONC(data, "")
		_, _ = e.Encode(n)
	}
}

func BenchmarkEncodeStruct_JSON(b *testing.B) {
	data := User{
		ID:   100,
		Name: "",
		Age:  25,
		Profile: struct {
			Address string `mm:"allow_empty"`
			Phone   []byte `mm:"type=bytes; allow_empty"`
		}{Address: ""},
		Decimal: "0.0",
	}

	out, _ := json.Marshal(data)
	b.Log("out", len(out))
	b.ResetTimer()

	for b.Loop() {
		_, _ = json.Marshal(data)
	}
}
