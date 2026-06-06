package core

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBigintArrEncoding(t *testing.T) {
	v1 := big.NewInt(0)
	v1.SetString("12345678901234567890", 10)
	v2 := big.NewInt(0)
	v2.SetString("98765432109876543210", 10)

	input := []any{*v1, *v2}
	bs, err := FromValue(input, "")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Go hex:", fmt.Sprintf("%x", bs))
	fmt.Println("Go len:", len(bs))
	fmt.Printf("Go bytes: %v\n", bs)
}