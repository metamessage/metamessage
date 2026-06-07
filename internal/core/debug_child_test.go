package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

func TestDebugChildTags(t *testing.T) {
	hexStr := "cf03f6def3856e616d65738673636f72657385636f756e74896f6c645f6669656c6485695f6172728669385f617272876931365f617272876933325f617272876936345f61727285755f6172728675385f617272877531365f617272877533325f617272877536345f617272876633325f617272876636345f617272877374725f61727288626f6f6c5f6172728962797465735f6172728a626967696e745f6172728c6461746574696d655f61727288646174655f6172728874696d655f61727288757569645f6172728b646563696d616c5f6172728669705f6172728775726c5f61727289656d61696c5f61727288656e756d5f617272fe320b8e096974656d206e616d65de24f20b8e096974656d206e616d6585616c696365f00b8e096974656d206e616d6583626f62fe2b0799b931c3313030de21ea0799b931c33130303850ea0799b931c3313030385aea0799b931c3313030385fe603423432382ae90111866c6567616379d3212223f402900ade0fe402900a21e402900a22e402900a23f402900bde0fe402900b21e402900b22e402900b23f402900cde0fe402900c21e402900c22e402900c23f402900dde0fe402900d21e402900d22e402900d23f402900ede0fe402900e21e402900e22e402900e23f402900fde0fe402900f21e402900f22e402900f23f4029010de0fe402901021e402901022e402901023f4029011de0fe402901121e402901122e402901123f4029012de0fe402901221e402901222e402901223fa029013de15e602901368ff0fe602901368ff19e602901368ff23d968ff0f68ff1968ff23d6816181628163d3060506f9029007de14e9029007a568656c6c6fe9029007a5776f726c64fe23029015de1eee029015aa140f6e462a062b3535a0ee029015aa147b747282315fad80a0f7029016de12e80290163b65920080e80290163b666d8948f3029017de0ee6029017394d0be6029017394db1f4029018de0fe602901839afc8e70290183a01517ffe2f029019de2af4029019b0550e8400e29b41d4a716446655440000f4029019b06ba7b8109dad11d180b400c04fd430c8f702901ade12e802901a6afb04cb2fe802901a6afb0425d4fe2202901bde1def02901b8b3139322e3136382e312e31ec02901b8831302e302e302e31fe3202901cde2df702901c9368747470733a2f2f6578616d706c652e636f6df402901c9068747470733a2f2f746573742e6f7267fe2d02901dde28f402901d9075736572406578616d706c652e636f6df202901d8e61646d696e40746573742e6f7267fe4c10d60e7265647c677265656e7c626c7565de39f210d60e7265647c677265656e7c626c756520f210d60e7265647c677265656e7c626c756521f210d60e7265647c677265656e7c626c756522"
	wire, _ := hex.DecodeString(hexStr)
	node, err := Decode(wire)
	if err != nil {
		t.Fatal("decode error:", err)
	}

	obj := node.(*ir.NodeObject)
	for _, f := range obj.Fields {
		if f.Key == "datetime_arr" {
			arr := f.Value.(*ir.NodeArray)
			fmt.Printf("datetime_arr: tag type=%v, child_type=%v, isInherit=%v\n", arr.Tag.Type, arr.Tag.ChildType, arr.Tag.IsInherit)
			for i, item := range arr.Items {
				v := item.(*ir.NodeScalar)
				fmt.Printf("  child[%d]: type=%v, data=%T(%v), text=%q, isInherit=%v\n",
					i, v.Tag.Type, v.Data, v.Data, v.Text, v.Tag.IsInherit)
			}
		}
	}

	// Test encode path
	fmt.Println("\n--- Testing encode path ---")
	node2, _ := ParseFromJSONC(`{
  // mm: child_type=datetime
  "datetime_arr": ["2024-01-01 00:00:00", "2024-06-15 12:30:00"]
}`)
	obj2 := node2.(*ir.NodeObject)
	for _, f := range obj2.Fields {
		if f.Key == "datetime_arr" {
			arr := f.Value.(*ir.NodeArray)
			fmt.Printf("Parsed datetime_arr: tag type=%v, child_type=%v, isInherit=%v\n", arr.Tag.Type, arr.Tag.ChildType, arr.Tag.IsInherit)
			for i, item := range arr.Items {
				v := item.(*ir.NodeScalar)
				fmt.Printf("  Parsed child[%d]: type=%v, child_type=%v, data=%T(%v), text=%q, isInherit=%v\n",
					i, v.Tag.Type, v.Tag.ChildType, v.Data, v.Data, v.Text, v.Tag.IsInherit)
				fmt.Printf("  Child tag Bytes() = %v\n", v.Tag.Bytes())
				fmt.Printf("  Child tag ToString() = %q\n", v.Tag.ToString())
				fmt.Printf("  Full tag dump: %+v\n", v.Tag)
			}
		}
	}
	enc := getEncoder()
	bs, _ := enc.Encode(node2)
	fmt.Printf("Encoded hex: %x\n", bs)

	decoded, _ := Decode(bs)
	obj3 := decoded.(*ir.NodeObject)
	for _, f := range obj3.Fields {
		if f.Key == "datetime_arr" {
			arr := f.Value.(*ir.NodeArray)
			fmt.Printf("Decoded datetime_arr: tag type=%v, child_type=%v, isInherit=%v\n", arr.Tag.Type, arr.Tag.ChildType, arr.Tag.IsInherit)
			for i, item := range arr.Items {
				v := item.(*ir.NodeScalar)
				fmt.Printf("  Decoded child[%d]: type=%v, child_type=%v, data=%T(%v), text=%q, isInherit=%v\n",
					i, v.Tag.Type, v.Tag.ChildType, v.Data, v.Data, v.Text, v.Tag.IsInherit)
			}
		}
	}
}

func TestDebugMimeTag(t *testing.T) {
	hexStr := "ce11d887636f6e74656e74e7028001a3313233"
	wire, _ := hex.DecodeString(hexStr)
	node, err := Decode(wire)
	if err != nil {
		t.Fatal("decode error:", err)
	}

	obj := node.(*ir.NodeObject)
	for _, f := range obj.Fields {
		v := f.Value.(*ir.NodeScalar)
		fmt.Printf("%s: type=%v, data=%T(%v), text=%q\n", f.Key, v.Tag.Type, v.Data, v.Data, v.Text)
	}
}