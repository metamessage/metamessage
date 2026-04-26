package main

import (
	"log"

	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/mm"
)

// go run examples/basic/*.go
func main() {
	var err error
	str := `123`
	encoded, err := mm.FromJSONC(str)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("encoded len: %d, original len: %d", len(encoded), len(str))

	resultJsonc, err := mm.ToJSONC(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Printf("decode result: %q", resultJsonc)

	result, err := mm.Decode(encoded)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	jsonc.Json(result)

	type V int
	var v V
	err = jsonc.Bind(result, &v)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	log.Printf("bind: %d", v)
}
