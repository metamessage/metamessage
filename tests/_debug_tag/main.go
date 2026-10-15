package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/metamessage/internal/ir"
)

func main() {
	tag := ir.NewTag()
	tag.ChildDesc = "item name"
	tag.Type = ir.ValueTypeVec

	bytesVal, err := tag.Bytes()
	if err != nil {
		fmt.Println("tag.Bytes() error:", err)
		os.Exit(1)
	}
	fmt.Println("tag with child_desc=\"item name\"")
	fmt.Println("tag.Bytes() hex:", hex.EncodeToString(bytesVal))
	fmt.Println("tag.Bytes() length:", len(bytesVal))
	fmt.Println("tag.ToString():", tag.ToString())

	os.Exit(0)
}