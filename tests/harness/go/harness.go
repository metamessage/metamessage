package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/jsonc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: harness [--encode|--decode] <file.jsonc>")
		os.Exit(1)
	}

	if os.Args[1] == "--encode" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: harness --encode <file.jsonc>")
			os.Exit(1)
		}
		data, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			os.Exit(1)
		}
		wire, err := core.FromJSONC(string(data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(hex.EncodeToString(wire))
		return
	}

	if os.Args[1] == "--decode" {
		hexBytes, _ := ioutil.ReadAll(os.Stdin)
		hexStr := string(hexBytes)
		wire, err := hex.DecodeString(hexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hex decode error: %v\n", err)
			os.Exit(1)
		}
		node, err := core.Decode(wire)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
			os.Exit(1)
		}
		output := jsonc.ToJSONC(node)
		fmt.Print(output)
		return
	}

	// Existing behavior
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