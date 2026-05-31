package main

import (
	"fmt"
	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
	"os"
)

func main() {
	data, _ := os.ReadFile(os.Args[1])
	node, _ := core.ParseFromJSONC(string(data))
	_ = node
	
	// Walk the tree and print tag info
	walkNode(node, 0)
	_ = jsonc.ToJSONC(node)
}

func walkNode(n ir.Node, depth int) {
	if n == nil { return }
	indent := ""
	for i := 0; i < depth; i++ { indent += "  " }
	
	tag := n.GetTag()
	if tag != nil {
		fmt.Printf("%sTag: isInherit=%v ToString=%q Type=%v ChildType=%v Desc=%q ChildDesc=%q Nullable=%v ChildNullable=%v\n",
			indent, tag.IsInherit, tag.ToString(), tag.Type, tag.ChildType, tag.Desc, tag.ChildDesc, tag.Nullable, tag.ChildNullable)
	}
	
	switch v := n.(type) {
	case *ir.Object:
		for _, f := range v.Fields {
			fmt.Printf("%s%s:\n", indent, f.Key)
			walkNode(f.Value, depth+1)
		}
	case *ir.Array:
		fmt.Printf("%s[array] items=%d\n", indent, len(v.Items))
		for _, item := range v.Items {
			walkNode(item, depth+1)
		}
	case *ir.Value:
		fmt.Printf("%sValue: text=%q\n", indent, v.Text)
	}
}