package metamessage

import (
	"fmt"
	"strings"
	"testing"
)

// TestEncodeFromJsoncSiblingsDoNotAccumulateDepth 回归：JSONC 文本入口不得把
// 「元素多」误判成「嵌套深」。
//
// 修复前 parseArray/parseObject 只 p.depth++ 从不 p.depth--，depth 退化成全局
// 累计计数器，于是 `[{"a":1}, ...]` 放到第 32 个元素就报 max depth: 32，
// 而真实嵌套深度只有 2。Go 值路径（EncodeFromValue）从未受影响。
func TestEncodeFromJsoncSiblingsDoNotAccumulateDepth(t *testing.T) {
	const n = 5000
	items := make([]string, n)
	for i := range items {
		items[i] = fmt.Sprintf(`{"a":%d}`, i+1)
	}
	wire, err := EncodeFromJsonc("[" + strings.Join(items, ",") + "]")
	if err != nil {
		t.Fatalf("%d sibling objects (real depth 2) must encode, got: %v", n, err)
	}
	if len(wire) == 0 {
		t.Fatal("expected non-empty wire")
	}

	// 解回来校验一遍，确认不只是"没报错"
	var v []struct {
		A int
	}
	if err := DecodeToValue(wire, &v); err != nil {
		t.Fatalf("DecodeToValue: %v", err)
	}
	if len(v) != n {
		t.Fatalf("decoded %d items, want %d", len(v), n)
	}
	for i, it := range v {
		if it.A != i+1 {
			t.Fatalf("item[%d].A = %d, want %d", i, it.A, i+1)
		}
	}
}

// TestEncodeFromJsoncDeepNestingStillRejected 深度上限必须仍然生效：
// 真实嵌套超过 256 层依旧报错，防止"把检查删掉"被当成修复。
func TestEncodeFromJsoncDeepNestingStillRejected(t *testing.T) {
	deep := strings.Repeat("[", 257) + "1" + strings.Repeat("]", 257)
	if _, err := EncodeFromJsonc(deep); err == nil {
		t.Fatal("257-level nesting must be rejected")
	} else if !strings.Contains(err.Error(), "max depth: 256") {
		t.Fatalf("error = %q, want it to mention max depth: 256", err)
	}
}
