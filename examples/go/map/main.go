package main

import (
	"fmt"
	"log"

	mm "github.com/metamessage/metamessage"
)

// go run examples/go/map/*.go
func main() {
	encoded, err := mm.EncodeFromObject(map[string]int{"test": 123456}, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Encoded: %v\n", encoded)

	// decode to JSONC
	resultJsonc, err := mm.DecodeToJSONC(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("decode result: %s", resultJsonc)

	// decode to map[string]int
	type T map[string]int
	var v T
	err = mm.Decode(encoded, &v)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// output:
	// // mm: type=map
	// {
	// 		"test": 123456,
	// }
	log.Printf("decode & bind: %+v", v)
}
