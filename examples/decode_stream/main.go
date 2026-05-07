package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/metamessage/metamessage/internal/mm"
)

// go run examples/decode_stream/*.go
func main() {
	type Message struct {
		Arr  []int
		Text string
	}

	message := Message{
		Arr:  []int{},
		Text: "abc",
	}

	bs, err := mm.FromValue(message, "")
	if err != nil {
		log.Fatalf("encode error: %v\n", err)
	}

	dec := mm.NewDecoder(bytes.NewReader(bs))
	for {
		var m Message
		_, err := dec.DecodeStream(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("decode: %v", err)
		}

		fmt.Printf("result %+v\n", m)
	}
}
