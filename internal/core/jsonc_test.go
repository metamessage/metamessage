package core

import (
	"fmt"
	"testing"

	"github.com/metamessage/metamessage/internal/jsonc"
)

// go test -v -run TestStrToJsonc
//
// go test internal/core/*.go -v -run TestStrToJsonc/user
func TestStrToJsonc(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name: "with_name",
			input: `{
  // mm: name=name
  "Name": "Alice",
}`,
		},
		{
			name: "invalid_email",
			input: `{
  // mm: type=email
  "Email": "Alice",
}`,
		},
		{
			name: "email",
			input: `{
  // mm: type=email
  "Email": "Alice@gmail.com",
}`,
		},
		{
			name: "invalid_url",
			input: `{
  // mm: type=url
  "Url": "/\/Alice",
}`,
		},
		{
			name: "url",
			input: `{
  // mm: type=url
  "Url": "https://Alice.com",
}`,
		},
		{
			name: "invalid_ip",
			input: `{
  // mm: type=ip
  "Url": "/\/Alice",
}`,
		},
		{
			name: "ip",
			input: `{
  // mm: type=ip
  "Url": "1.1.1.1",
}`,
		},
		{
			name: "ipv6",
			input: `{
  // mm: type=ip
  "Url": "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
}`,
		},
		{
			name: "invalid_decimal",
			input: `{
  // mm: type=decimal
  "Url": "/\/Alice",
}`,
		},
		{
			name: "decimal",
			input: `{
  // mm: type=decimal
  "Url": "123.0",
}`,
		},
		{
			name: "invalid_uuid",
			input: `{
  // mm: type=uuid
  "Url": "/\/Alice",
}`,
		},
		{
			name: "uuid",
			input: `{
  // mm: type=uuid
  "Url": "d46372a1-5b9e-4f8c-9a2d-7e1b3c5d7f9a",
}`,
		},
		{
			name:    "nil input",
			input:   "null",
			want:    "null",
			wantErr: false,
		},
		{
			name: "slice with empty values",
			input: `[
  "a",
  "",
  "c"
]`,
			want: `[
  "a",
  "",
  "c"
]`,
		},
		{
			name:    "unsupported type (channel)",
			input:   "",
			wantErr: true,
		},
		{
			name: "user",
			input: `
			
			// mm: desc=用户
    {
      // mm: type=i64; desc=用户ID
      "id": 666,
      // mm: desc=昵称
      "name": "abc",
	        // mm: type=u8
    "age": 20
    }
			
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseFromJSONC(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromJSONC error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(Dump(node))
			t.Log(jsonc.ToJSONC(node))
			if !tt.wantErr && jsonc.ToJSONC(node) != tt.want {
				t.Errorf("ToJSONC() = \n%v, want \n%v", node, tt.want)
			}

			encoder := getEncoder()
			defer putEncoder(encoder)
			wire, err := encoder.Encode(node)
			fmt.Println("wire", wire)
			_, _ = Decode(wire)

		})
	}
}
