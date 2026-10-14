package parser

import (
	"fmt"
	"strings"
	"testing"
)

// nestArrays 构造 n 层嵌套数组，最内层是 inner。
// 顶层数组算第 1 层，所以 n == maxDepth 时应被接受、n == maxDepth+1 时应被拒。
func nestArrays(n int, inner string) string {
	return strings.Repeat("[", n) + inner + strings.Repeat("]", n)
}

// TestParseArraySiblingsDoNotAccumulateDepth 回归：数组的同级元素不得把前面
// 所有元素累计的深度继承过去。
//
// 修复前 parseArray/parseObject 只 p.depth++ 从不 p.depth--，depth 退化成
// 「至今累计见过多少层」的全局计数器，于是 [{"a":1}, ...] 放到第 32 个元素就
// 假报 max depth——而它的真实嵌套深度只有 2。
func TestParseArraySiblingsDoNotAccumulateDepth(t *testing.T) {
	const n = 5000
	items := make([]string, n)
	for i := range items {
		items[i] = fmt.Sprintf(`{"a":%d}`, i+1)
	}
	p := scan("[" + strings.Join(items, ",") + "]")

	node, err := p.Parse()
	if err != nil {
		t.Fatalf("5000 sibling objects (real depth 2) must parse, got: %v", err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	// 解析成功後 depth 必须归零：这才是「深度」而不是「计数器」。
	if p.depth != 0 {
		t.Fatalf("depth not restored after successful parse: %d, want 0", p.depth)
	}
}

// TestParseObjectSiblingsDoNotAccumulateDepth 同上，但对象有大量「值是对象」的字段。
func TestParseObjectSiblingsDoNotAccumulateDepth(t *testing.T) {
	const n = 5000
	fields := make([]string, n)
	for i := range fields {
		fields[i] = fmt.Sprintf(`"f%d":{"a":%d}`, i, i+1)
	}
	p := scan("{" + strings.Join(fields, ",") + "}")

	node, err := p.Parse()
	if err != nil {
		t.Fatalf("object with %d object-valued fields (real depth 2) must parse, got: %v", n, err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if p.depth != 0 {
		t.Fatalf("depth not restored after successful parse: %d, want 0", p.depth)
	}
}

// TestParseMaxDepthStillEnforced 深度上限必须仍然生效——防止「把检查删掉」被
// 当成修复。真实嵌套超过 maxDepth 依旧要报错。
func TestParseMaxDepthStillEnforced(t *testing.T) {
	if _, err := scan(nestArrays(maxDepth, "1")).Parse(); err != nil {
		t.Fatalf("%d nested arrays (== maxDepth) must parse, got: %v", maxDepth, err)
	}
	_, err := scan(nestArrays(maxDepth+1, "1")).Parse()
	if err == nil {
		t.Fatalf("%d nested arrays must exceed maxDepth=%d, but parsed fine", maxDepth+1, maxDepth)
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("max depth: %d", maxDepth)) {
		t.Fatalf("error = %q, want it to mention max depth %d", err, maxDepth)
	}
}

// TestParseDepthRestoredAcrossSiblings 混合形状：宽（兄弟多）+ 深（嵌套到上限附近）
// 同时出现。修复前只要前面出现过足够多的宽数组，后面的合法深层嵌套就会被误杀。
func TestParseDepthRestoredAcrossSiblings(t *testing.T) {
	const siblings = 200
	items := make([]string, siblings)
	for i := range items {
		// 每个元素自身嵌套到接近上限（maxDepth-1 层，加外层数组正好 maxDepth）
		items[i] = nestArrays(maxDepth-1, fmt.Sprintf("%d", i+1))
	}
	// 顶层数组 + 每个元素 maxDepth-1 层 ⇒ 峰值正好 maxDepth，必须全部通过
	if _, err := scan("[" + strings.Join(items, ",") + "]").Parse(); err != nil {
		t.Fatalf("%d siblings each nested %d deep (peak == maxDepth) must parse, got: %v",
			siblings, maxDepth-1, err)
	}

	// 再多一层就越界，必须报错
	deep := make([]string, siblings)
	for i := range deep {
		deep[i] = nestArrays(maxDepth, fmt.Sprintf("%d", i+1))
	}
	if _, err := scan("[" + strings.Join(deep, ",") + "]").Parse(); err == nil {
		t.Fatal("siblings nested one level past maxDepth must be rejected")
	}
}
