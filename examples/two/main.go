package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/lizongying/meta-message/internal/mm"
)

func main() {
	type Message struct {
		Name []int
		Text string
	}
	s := Message{
		Name: []int{},
		// Text: "2333333",
	}

	bs, e := mm.FromStruct(s, "")
	if e != nil {
		fmt.Printf("Error: %s\n", e)
		return
	}

	dec := mm.NewDecoder(bytes.NewReader(bs))
	for {
		var m Message
		if _, err := dec.DecodeStream(&m); err == io.EOF {
			fmt.Printf("00000 %s %s\n", m.Name, m.Text)
			log.Fatal(err)
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("11111 %v: %s\n", m.Name, m.Text)
	}
}
