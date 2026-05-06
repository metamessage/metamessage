package main

import (
	"log"

	mm "github.com/metamessage/metamessage"
)

// go run examples/go/int/*.go
func main() {
	var err error
	str := `123`

	// encode
	encoded, err := mm.EncodeFromJSONC(str)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("encoded len: %d, original len: %d", len(encoded), len(str))

	// decode to JSONC
	resultJsonc, err := mm.DecodeToJSONC(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("decode result: %s", resultJsonc)

	// decode to int
	var v int
	err = mm.Decode(encoded, &v)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// output:
	// 123
	log.Printf("decode & bind: %+v", v)
}
