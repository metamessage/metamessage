package main

import (
	"bytes"
	"fmt"

	"github.com/metamessage/metamessage/internal/core"
)

// go run examples/encode_stream/*.go
func main() {
	var buf bytes.Buffer
	enc := core.NewEncoder(&buf)

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
