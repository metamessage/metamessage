package core

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

// v0.5.0 ml>255 盲区根治回归测试。
//
// 历史缺陷链（v0.3.x 起潜伏、ml>255 盲区从未执行过）：
//   - 正向 tag 前缀读取写成 `l < 257`，对字节值恒真，[255,hi,lo] 三档为死代码；
//   - meta 串属性（desc/enums 家族）正向解码用 3-bit payload 比 `l <= 1<<8`，
//     payload==7 的 2 字节档同为死代码；
//   - ir/tag.go 串 writer 的 `l <= 1<<8` 允许 l=256 落入 1 字节槽，byte(256) 回绕；
//   - untailify 反推 inner 用旧 metaPrefixWidth 而非实际尾槽宽，ml=256 差一。
// 修复后：尾槽/前缀槽统一 canonical 三档（<254→1、254..255→2、256..65535→3），
// 反推式 inner = total - ml - slotW，死代码档全部激活。以下锁定两层不变式。

// TestFooterPrefixSlotPairing：全部合法 ml 区间上，尾槽与前缀槽宽度严格相等，
// 且 parseFooterSlot 与 footerSlotBytes 互逆。这是编解码边界一致性的数学不变式。
func TestFooterPrefixSlotPairing(t *testing.T) {
	for ml := 0; ml <= 65535; ml++ {
		slot, err := footerSlotBytes(ml)
		if err != nil {
			t.Fatalf("footerSlotBytes(%d): %v", ml, err)
		}
		prefix, err := metaPrefixBytes(ml)
		if err != nil {
			t.Fatalf("metaPrefixBytes(%d): %v", ml, err)
		}
		if len(prefix) != len(slot) {
			t.Fatalf("ml=%d: 前缀槽宽 %d ≠ 尾槽宽 %d", ml, len(prefix), len(slot))
		}
		got, w, err := parseFooterSlot(slot, len(slot))
		if err != nil || got != ml || w != len(slot) {
			t.Fatalf("ml=%d: parseFooterSlot 逆不还原 (%d,%d,%v)", ml, got, w, err)
		}
	}
	if _, err := footerSlotBytes(65536); err == nil {
		t.Fatal("ml=65536 超格式上限，必须报错")
	}
}

// TestTagMetaRoundTripAboveMax1Byte：覆盖原盲区及全部档位边界的 desc 内容级
// round-trip（编码 → 解码 → tag 字段逐字比对），杜绝旧 grep 式假阳性。
func TestTagMetaRoundTripAboveMax1Byte(t *testing.T) {
	// desc 属性 = key(1)+len 槽(2)+串长，65532 使 ml 恰达上限 65535。
	for _, l := range []int{1, 5, 6, 253, 254, 255, 256, 257, 261, 1000, 65532} {
		desc := strings.Repeat("x", l)
		src := "{\n  // mm: desc=" + desc + "\n  \"v\": 1\n}\n"
		wire, err := FromJSONC(src)
		if err != nil {
			t.Fatalf("len=%d: encode: %v", l, err)
		}
		node, err := Decode(wire)
		if err != nil {
			t.Fatalf("len=%d: decode: %v", l, err)
		}
		obj, ok := node.(*ir.NodeObject)
		if !ok || len(obj.Fields) != 1 {
			t.Fatalf("len=%d: 解码节点形态异常 %T", l, node)
		}
		scalar, ok := obj.Fields[0].Value.(*ir.NodeScalar)
		if !ok {
			t.Fatalf("len=%d: 字段值非标量 %T", l, obj.Fields[0].Value)
		}
		if got := scalar.GetTag().Desc; got != desc {
			t.Fatalf("len=%d: round-trip 后 desc 长度 %d，内容不符", l, len(got))
		}
	}
}
