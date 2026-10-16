package metamessage

import "testing"

// zzrepro_test.go — 空容器兩個復現 + 對照組：
//   1. 空 []string 編碼後 wire 帶 example 哨兵 ""，解回来變 [""]（len 1，不是 0）
//   2. 空 []int32 的 example 哨兵 0 類型在回讀時降級為 generic int，
//      綁進 int32 槽位硬錯 "target type must be int, got int32"
//   3. 對照組：真實元素 []int32{5} round-trip 正常；而真實零值元素
//      []int32{0} / []string{""} 在 encode 階段就被「not allow empty」拒絕——
//      只有 example 標記的零值哨兵能上 wire。

type zzStrBox struct {
	Names []string `mm:"allow_empty"`
}

type zzIntBox struct {
	Nums []int32 `mm:"allow_empty"`
}

// 問題 1：空 []string → 解成 [""]（len 1）
func TestZzEmptyStringArray(t *testing.T) {
	wire, err := EncodeFromValue(zzStrBox{Names: []string{}}, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	t.Logf("ENCODED len=%d: %q", len(wire), string(wire))
	if jsonc, err := DecodeToJsonc(wire); err == nil {
		t.Logf("JSONC:\n%s", jsonc)
	} else {
		t.Logf("DecodeToJsonc: %v", err)
	}
	var out zzStrBox
	if err = DecodeToValue(wire, &out); err != nil {
		t.Fatalf("decode err: %v", err)
	}
	t.Logf("DECODED Names len=%d %#v (input len 0)", len(out.Names), out.Names)
}

// 問題 2：空 []int32 → 硬解碼報錯
func TestZzEmptyIntArray(t *testing.T) {
	wire, err := EncodeFromValue(zzIntBox{Nums: []int32{}}, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	t.Logf("ENCODED len=%d: %q", len(wire), string(wire))
	if jsonc, err := DecodeToJsonc(wire); err == nil {
		t.Logf("JSONC:\n%s", jsonc)
	} else {
		t.Logf("DecodeToJsonc: %v", err)
	}
	var out zzIntBox
	if err = DecodeToValue(wire, &out); err != nil {
		t.Fatalf("decode err: %v", err) // 預期在這裏失敗：target type must be int, got int32
	}
	t.Logf("DECODED Nums len=%d %#v", len(out.Nums), out.Nums)
}

// 對照組：真實（非 example）元素
func TestZzRealElements(t *testing.T) {
	wire, err := EncodeFromValue(zzIntBox{Nums: []int32{5}}, "")
	if err != nil {
		t.Fatalf("encode [5]: %v", err)
	}
	var out zzIntBox
	if err = DecodeToValue(wire, &out); err != nil {
		t.Fatalf("decode [5]: %v", err)
	}
	t.Logf("real []int32{5} -> OK %#v（per-item 類型沒丟）", out.Nums)

	_, err = EncodeFromValue(zzIntBox{Nums: []int32{0}}, "")
	t.Logf(`encode []int32{0}（真實零值元素）-> err=%v`, err)

	wire2, err2 := EncodeFromValue(zzStrBox{Names: []string{""}}, "")
	t.Logf(`encode []string{""}（真實空字串元素）-> err=%v`, err2)
	if err2 == nil {
		var out2 zzStrBox
		t.Logf(`decode [""] -> err=%v out=%#v`, DecodeToValue(wire2, &out2), out2.Names)
	}
}

// 追加：example 注入的觸發面——map/vec 有、nil 指標也有、固定 arr 沒有。
type zzInner struct {
	A int
	B string
}

type zzNullBox struct {
	Ptr *zzInner           `mm:"allow_empty"` // nil 指標
	Arr [0]int             `mm:"allow_empty"` // 固定長度 0 的陣列
	M   map[string]zzInner `mm:"allow_empty"` // 空 map（對照，應有 example entry）
}

func TestZzExampleTriggerSurface(t *testing.T) {
	in := zzNullBox{Ptr: nil, M: map[string]zzInner{}}
	wire, err := EncodeFromValue(in, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	jsonc, err := DecodeToJsonc(wire)
	if err != nil {
		t.Fatalf("jsonc: %v", err)
	}
	t.Logf("JSONC:\n%s", jsonc)
	var out zzNullBox
	err = DecodeToValue(wire, &out)
	t.Logf("decode back: err=%v Ptr==nil?%v Arr(len=%d) M(len=%d)=%#v", err, out.Ptr == nil, len(out.Arr), len(out.M), out.M)
}
