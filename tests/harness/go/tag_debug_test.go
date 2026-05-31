package main

import (
	"fmt"
	"github.com/metamessage/metamessage/internal/ir"
)

func debugTagInherit() {
	parent := ir.NewTag()
	parent.ChildType = ir.ValueTypeI8
	parent.ChildDesc = "test desc"
	parent.ChildNullable = true
	parent.ChildMin = "1"
	parent.ChildMax = "100"

	child := ir.NewTag()
	child.Inherit(parent)

	fmt.Printf("After Inherit:\n")
	fmt.Printf("  isInherit: %v\n", child.IsInherit)
	fmt.Printf("  Type: %v\n", child.Type)
	fmt.Printf("  ChildType: %v\n", child.ChildType)
	fmt.Printf("  Desc: %q\n", child.Desc)
	fmt.Printf("  ChildDesc: %q\n", child.ChildDesc)
	fmt.Printf("  Nullable: %v\n", child.Nullable)
	fmt.Printf("  Min: %q\n", child.Min)
	fmt.Printf("  Max: %q\n", child.Max)
	fmt.Printf("  ToString(): %q\n", child.ToString())
}