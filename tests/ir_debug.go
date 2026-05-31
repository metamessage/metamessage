package main

import (
	"fmt"
	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
)

func main() {
	input := `{
// mm: child_desc=item name
"names": ["alice", "bob"]
}`
	node, err := core.ParseFromJSONC(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	obj := node.(*ir.Object)
	for _, f := range obj.Fields {
		fmt.Printf("Field: %s\n", f.Key)
		fmt.Printf("  Field tag ToString: %q\n", f.Value.GetTag().ToString())
		
		if arr, ok := f.Value.(*ir.Array); ok {
			fmt.Printf("  Array tag ToString: %q\n", arr.Tag.ToString())
			fmt.Printf("  Array tag ChildDesc: %q\n", arr.Tag.ChildDesc)
			for i, item := range arr.Items {
				t := item.GetTag()
				fmt.Printf("  Item[%d] tag ToString: %q\n", i, t.ToString())
				fmt.Printf("    Type: %v, IsInherit: %v, Desc: %q, ChildDesc: %q\n", 
					t.Type, t.IsInherit, t.Desc, t.ChildDesc)
			}
		}
	}

	fmt.Println("\n=== Full JSONC Output ===")
	fmt.Print(jsonc.ToJSONC(node))
}