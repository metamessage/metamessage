package main

import (
	"fmt"
	"log"

	mm "github.com/metamessage/metamessage"
)

// go run examples/go/object/*.go
func main() {
	type T struct {
		Name string
		Age  int
	}

	t := T{
		Name: "Alice",
		Age:  30,
	}
	encoded, err := mm.EncodeFromValue(t, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Printf("encoded len: %d", len(encoded))

	// decode to JSONC
	resultJsonc, err := mm.DecodeToJSONC(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("decode result: %s", resultJsonc)

	// decode to object
	var v T
	err = mm.DecodeToValue(encoded, &v)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// output:
	// {
	//         "name": "Alice",
	//         "age": 30,
	// }
	log.Printf("decode & bind: %+v", v)
}
