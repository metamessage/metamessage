package main

import (
	"bytes"
	"fmt"

	"github.com/lizongying/meta-message/internal/mm"
)

func main() {
	var w bytes.Buffer
	encoder := mm.NewEncoder(&w)
	for range 10 {
		n, e := encoder.EncodeStream(1)
		fmt.Println(999, n, e)
	}
	fmt.Println(w.String())
}
