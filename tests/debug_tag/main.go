package main

import (
	"fmt"
	"github.com/metamessage/metamessage/internal/ir"
)

func main() {
	// Create array tag with child_desc
	arrTag := ir.NewTag()
	arrTag.ChildDesc = "item name"
	arrTag.ChildType = ir.ValueTypeI8

	// Create child tag and inherit
	childTag := ir.NewTag()
	childTag.Inherit(arrTag)

	fmt.Printf("After Inherit:\n")
	fmt.Printf("  IsInherit: %v\n", childTag.IsInherit)
	fmt.Printf("  Desc: %q\n", childTag.Desc)
	fmt.Printf("  ChildDesc: %q\n", childTag.ChildDesc)
	fmt.Printf("  Type: %v\n", childTag.Type)
	fmt.Printf("  ChildType: %v\n", childTag.ChildType)
	fmt.Printf("  Bytes: %v\n", childTag.Bytes())
	fmt.Printf("  Bytes hex: %x\n", childTag.Bytes())

	// Now check what Bytes() produces
	bytesVal := childTag.Bytes()
	fmt.Printf("  Bytes len: %d\n", len(bytesVal))
	for i, b := range bytesVal {
		fmt.Printf("    [%d] = 0x%02x (%08b)\n", i, b, b)
	}
}