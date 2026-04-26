package main

import (
	"bytes"
	"fmt"

	"github.com/metamessage/metamessage/internal/mm"
)

// go run examples/encode_stream/*.go
func main() {
	var buf bytes.Buffer
	enc := mm.NewEncoder(&buf)

	for range 10 {
		_, err := enc.EncodeStream(1)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
	}

	// output:
	// result: !!!!!!!!!!
	fmt.Println("result:", buf.String())
}
