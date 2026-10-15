package main

import (
"fmt"
"github.com/metamessage/metamessage/internal/ir"
)

func init() {
parentTag := ir.NewTag()
parentTag.ChildDesc = "item name"
parentTag.ChildType = ir.ValueTypeI8

fmt.Printf("Parent: ChildDesc=%q, ChildType=%v\n", parentTag.ChildDesc, parentTag.ChildType)

childTag := ir.NewTag()
childTag.Inherit(parentTag)

fmt.Printf("Child AFTER Inherit: IsInherit=%v, Desc=%q, ChildDesc=%q, Type=%v\n",
childTag.IsInherit, childTag.Desc, childTag.ChildDesc, childTag.Type)
tb, err := childTag.Bytes()
if err != nil {
fmt.Printf("Child tag bytes error: %v\n", err)
} else {
fmt.Printf("Child tag bytes: %x\n", tb)
}
}
