package main

import (
	"fmt"

	"github.com/metamessage/metamessage"
)

func main() {
	data, err := metamessage.EncodeFromStruct(map[string]interface{}{"test": 123456}, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Encoded: %v\n", data)
	fmt.Printf("Hex: ")
	for _, b := range data {
		fmt.Printf("%02x ", b)
	}
	fmt.Println()
}
