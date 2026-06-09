package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc/parser"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
)

func main() {
	// Test Inherit behavior directly
	parentTag := ir.NewTag()
	parentTag.ChildDesc = "item name"
	parentTag.ChildType = ir.ValueTypeI8

	childTag := ir.NewTag()
	childTag.Inherit(parentTag)

	fmt.Printf("TEST Inherit: Child: IsInherit=%v, Desc=%q, ChildDesc=%q, Type=%v\n",
		childTag.IsInherit, childTag.Desc, childTag.ChildDesc, childTag.Type)
	fmt.Printf("TEST Inherit: Child tag bytes: %x\n\n", childTag.Bytes())

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file.jsonc>\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}

	toks := scanner.New(string(data)).ScanAll()
	p := parser.New(toks)
	node, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	obj, ok := node.(*ir.NodeObject)
	if !ok {
		fmt.Fprintf(os.Stderr, "expected object\n")
		os.Exit(1)
	}

	for _, field := range obj.Fields {
		arr, ok := field.Value.(*ir.NodeArray)
		if !ok {
			continue
		}
		fmt.Printf("Field: %s\n", field.Key)
		fmt.Printf("  Array tag JSON: %s\n", arr.Tag.Json())

		for i, item := range arr.Items {
			val, ok := item.(*ir.NodeScalar)
			if !ok {
				continue
			}
			b, _ := json.MarshalIndent(val.Tag, "", "    ")
			fmt.Printf("  Item[%d] tag JSON: %s\n", i, string(b))
			fmt.Printf("  Item[%d] tag bytes: %x\n", i, val.Tag.Bytes())
			if i >= 0 {
				break
			}
		}
	}
}
