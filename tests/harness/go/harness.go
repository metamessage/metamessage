package main

import (
	"fmt"
	"os"

	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/jsonc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: harness <file.jsonc>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}

	node, err := core.ParseFromJSONC(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	output := jsonc.ToJSONC(node)
	fmt.Print(output)
}