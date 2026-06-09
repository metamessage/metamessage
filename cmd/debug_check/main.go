package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc/parser"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
	"github.com/metamessage/metamessage/internal/jsonc/token"
)

func main() {
	data, _ := os.ReadFile("tests/fixtures/03_tags/child_tags.jsonc")
	sc := scanner.New(string(data))
	var toks []token.Token
	for {
		t := sc.NextToken()
		toks = append(toks, t)
		if t.Type == token.EOF {
			break
		}
	}

	p := parser.New(toks)
	node, err := p.Parse()
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	dumpNode(node, "")
}

func dumpNode(node ir.Node, indent string) {
	switch n := node.(type) {
	case *ir.NodeObject:
		for _, f := range n.Fields {
			fmt.Printf("%sField %q:\n", indent, f.Key)
			dumpNode(f.Value, indent+"  ")
		}
	case *ir.NodeArray:
		fmt.Printf("%sArray (items: %d):\n", indent, len(n.Items))
		if n.Tag != nil {
			dumpTag(indent+"  ArrayTag: ", n.Tag)
			tb := n.Tag.Bytes()
			fmt.Printf("%s  Bytes(): %s (len=%d)\n", indent, hex.EncodeToString(tb), len(tb))
		}
		for i, item := range n.Items {
			if v, ok := item.(*ir.NodeScalar); ok {
				fmt.Printf("%s  [%d] Value data=%v text=%q:\n", indent, i, v.Data, v.Text)
				if v.Tag != nil {
					dumpTag(indent+"    Tag: ", v.Tag)
					tb := v.Tag.Bytes()
					fmt.Printf("%s    Bytes(): %s (len=%d)\n", indent, hex.EncodeToString(tb), len(tb))
				} else {
					fmt.Printf("%s    Tag: nil\n", indent)
				}
			} else {
				fmt.Printf("%s  [%d]: (non-Value)\n", indent, i)
				dumpNode(item, indent+"    ")
			}
		}
	case *ir.NodeScalar:
		fmt.Printf("%sValue: data=%v text=%q\n", indent, n.Data, n.Text)
	}
}

func dumpTag(prefix string, tag *ir.Tag) {
	fmt.Printf("%sIsInherit=%v Type=%v Desc=%q ChildDesc=%q\n",
		prefix, tag.IsInherit, tag.Type, tag.Desc, tag.ChildDesc)
	fmt.Printf("%s  ChildType=%v Location=%v Enums=%q ChildEnums=%q\n",
		prefix, tag.ChildType, tag.Location, tag.Enums, tag.ChildEnums)
}
