package main

import (
	"log"

	mm "github.com/metamessage/metamessage"
)

// go run examples/go/some/*.go
func main() {
	var err error
	type User struct {
		ID       int64  `mm:"desc=用戶ID"`
		Name     string `mm:"desc=用戶名稱;min=1;max=50"`
		Email    string `mm:"type=email;desc=電子郵箱"`
		Age      uint8  `mm:"desc=年齡;min=0;max=150"`
		IsActive bool   `mm:"desc=是否激活"`
	}

	var users = []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 25, IsActive: true},
		// {ID: 2, Name: "Bob", Email: "bob@example.com", Age: 30, IsActive: true},
		// {ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 35, IsActive: false},
	}

	type ListUsersResponse struct {
		Total int64  `mm:"desc=總數"`
		Users []User `mm:"desc=用戶列表"`
	}

	data := ListUsersResponse{
		Total: int64(len(users)),
		Users: users,
	}

	encoded, err := mm.EncodeFromValue(data, "")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// decode to JSONC
	resultJsonc, err := mm.DecodeToJsonc(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("decode result: %s", resultJsonc)

	// decode to int
	var v ListUsersResponse
	err = mm.DecodeToValue(encoded, &v)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// output:
	// 123
	log.Printf("decode & bind: %+v", v)
}
