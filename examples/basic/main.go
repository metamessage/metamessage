package main

import (
	"fmt"

	"github.com/lizongying/meta-message/internal/jsonc"
	"github.com/lizongying/meta-message/internal/mm"
)

func main() {
	var s = `
 // mm: 0
	{
 // mm: 1
"Name": "Ed", // mm: 11
"some": -1,
"other": null,
"size": 21,
 // other: type=int 
 // mm: type=int 
"Age": null, // mm: type=int 
"Score": "23.0", // mm: type=int
"Text": "Knock knock."
/* 123
456
*/
// 2323233
"Group": [
1,// mm: a1
// mm: a2
2,
3, // mm: a
], // mm: b
}, // mm: c
`
	s = `
	// mm: type=map; nullable
	{
		//  "na": "232233",
		// "yes": true, // mm: type=int; desc=年纪 
	// "ok": false,
		// mm: name=name;
	//  "Name": "Ed",
	 		// mm: type=bi;
	//  "big": "1234",
	//  "size": 21, // mm: type=i8 
	//  "Score": 10.0233232323,

	// 這是不行的，因為間隔了空行
// mm: type=int

	// 這是可以的 size=32; 
	 // mm: child_type=i;
"Group": [
1,// mm: a1
// mm: a2
2,
3, // mm: a
], // mm: b
// "c":{
// "aa":123,
// "bb":false,
// "cc":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20]}
}
`
	s = `123`
	fmt.Println(s)
	fmt.Println("-----------------")
	r, e := mm.FromJSONC(s)
	if e != nil {
		fmt.Printf("Error: %s\n", e)
		return
	}
	fmt.Println("1-----------------")
	fmt.Println("bytes", len(r), "string", len(s))
	fmt.Println("2-----------------")
	r1, e := mm.ToJSONC(r)
	if e != nil {
		fmt.Printf("Error: %s\n", e)
	}
	fmt.Println("decode", r1)

	r2, e2 := mm.Decode(r)
	if e2 != nil {
		fmt.Printf("Error: %s\n", e)
	}
	jsonc.Json(r2)

	// type C struct {
	// 	Group [32]string
	// }

	// type C map[string][]int

	type C int
	var c int
	e3 := jsonc.Bind(r2, &c)
	if e3 != nil {
		fmt.Printf("Error: %v\n", e3)
	}
	fmt.Println("bind", c)

}
