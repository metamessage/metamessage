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

	bytes := tag.Bytes()
	fmt.Println("tag with child_desc=\"item name\"")
	fmt.Println("tag.Bytes() hex:", hex.EncodeToString(bytes))
	fmt.Println("tag.Bytes() length:", len(bytes))
	fmt.Println("tag.ToString():", tag.ToString())

	os.Exit(0)
}