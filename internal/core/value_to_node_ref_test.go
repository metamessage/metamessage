package core

import (
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
)

// ProbeNode 是触发原 stack overflow 的最小自引用结构：
// 编码零值时 nil Ref 走 nil 指针示例分支、空 Kids 走空 slice 示例分支，
// createExampleValue 在 struct↔slice 之间无限递归直至爆栈。
type ProbeNode struct {
	Name string
	Kids []ProbeNode
	Ref  *ProbeNode
}

// 修复回归：自引用结构编码不再 stack overflow，且按 D1/D2/D5 正确标注
// ref_id（被引用节点）、ref_to（引用节点）。
// example 模式（schema/示例生成）是本特性服务的场景。
func TestValueToNode_SelfReference_NoStackOverflow(t *testing.T) {
	node, err := ValueToNode(ProbeNode{}, "example")
	if err != nil {
		t.Fatalf("ProbeNode must encode without stack overflow, got err: %v", err)
	}

	root, ok := node.(*ir.NodeObject)
	if !ok {
		t.Fatalf("expected root NodeObject, got %T", node)
	}
	if root.Tag.RefID != 1 {
		t.Fatalf("expected root ref_id=1, got %d", root.Tag.RefID)
	}
	if root.Tag.RefTo != 0 {
		t.Fatalf("root must not carry ref_to, got %d", root.Tag.RefTo)
	}

	var kids *ir.NodeArray
	var refNode *ir.NodeObject
	for _, f := range root.Fields {
		switch f.Key {
		case "kids":
			arr, ok := f.Value.(*ir.NodeArray)
			if !ok {
				t.Fatalf("kids should be NodeArray, got %T", f.Value)
			}
			kids = arr
		case "ref":
			obj, ok := f.Value.(*ir.NodeObject)
			if !ok {
				t.Fatalf("ref should be NodeObject shell, got %T", f.Value)
			}
			refNode = obj
		}
	}
	if kids == nil {
		t.Fatalf("kids field missing")
	}
	if len(kids.Items) != 1 {
		t.Fatalf("expected 1 example item in kids, got %d", len(kids.Items))
	}
	item, ok := kids.Items[0].(*ir.NodeObject)
	if !ok {
		t.Fatalf("kids[0] should be NodeObject shell, got %T", kids.Items[0])
	}
	if item.Tag.RefTo != 1 {
		t.Fatalf("expected kids[0] ref_to=1, got %d", item.Tag.RefTo)
	}
	if len(item.Fields) != 0 {
		t.Fatalf("kids[0] shell must have no fields, got %d", len(item.Fields))
	}

	if refNode == nil {
		t.Fatalf("ref field missing")
	}
	if refNode.Tag.RefTo != 1 {
		t.Fatalf("expected ref ref_to=1, got %d", refNode.Tag.RefTo)
	}
	if len(refNode.Fields) != 0 {
		t.Fatalf("ref shell must have no fields, got %d", len(refNode.Fields))
	}
}

// jsonc 文本表示（D5）：只有引用节点和被引用节点标注 ref 标签。
func TestValueToNode_SelfReference_JSONCShape(t *testing.T) {
	node, err := ValueToNode(ProbeNode{}, "example")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	out := jsonc.ToJSONC(node)
	for _, want := range []string{"ref_id=1", "ref_to=1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("jsonc output missing %q:\n%s", want, out)
		}
	}
}

// 平行分支不误判（D3）：同一非环类型出现在两个兄弟字段，均完整展开、零 ref 标记。
func TestValueToNode_ParallelBranchNoFalsePositive(t *testing.T) {
	type simple struct {
		N int
	}
	type twoSiblings struct {
		A simple
		B simple
	}

	node, err := ValueToNode(twoSiblings{}, "example")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	assertNoRefTags(t, node)
}

// 深层无环结构：零 ref 标记（既有空容器行为不变的守卫）。
func TestValueToNode_AcyclicNoRefTags(t *testing.T) {
	cases := []struct {
		name string
		in   any
	}{
		{"empty int slice", []int{}},
		{"nested empty slices", [][]string{}},
		{"triple nested slices", [][][]int{}},
		{"nested map slices", [][]map[string]string{}},
		{"fixed array of slices", [2][]*uint8{}},
		{"map of slices", map[string][]int{}},
		{"struct with empty slice", struct {
			Arr []int
		}{}},
		{"struct with nil ptr slice", struct {
			Ptrs []*string
		}{}},
		{"deep acyclic nesting", struct {
			Inner struct {
				List []struct {
					Vals []int
				}
			}
		}{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := ValueToNode(tc.in, "")
			if err != nil {
				// 既有行为：部分空容器在非 example 模式下校验报错，保持不变
				t.Logf("expected-behavior err: %v", err)
				return
			}
			assertNoRefTags(t, node)
		})
	}
}

// 真实数据环（非 example 路径）保持既有行为：maxDepth 报错，而非引用壳。
func TestValueToNode_RealCycleStillErrors(t *testing.T) {
	type Recursive struct {
		Self *Recursive
	}
	r := &Recursive{}
	r.Self = r

	_, err := ValueToNode(r, "")
	if err == nil {
		t.Fatal("expected error for real (non-example) cycle, got nil")
	}
	if !strings.Contains(err.Error(), "max depth") {
		t.Fatalf("expected max depth error, got: %v", err)
	}
}

// 自嵌套 slice 类型（不经 anyToJSONC 栈的自嵌套防护）。
func TestValueToNode_InfSliceType(t *testing.T) {
	type InfSlice []InfSlice
	node, err := ValueToNode(InfSlice{}, "example")
	if err != nil {
		t.Fatalf("InfSlice must encode without overflow, got err: %v", err)
	}
	arr, ok := node.(*ir.NodeArray)
	if !ok {
		t.Fatalf("expected NodeArray, got %T", node)
	}
	if len(arr.Items) != 1 {
		t.Fatalf("expected 1 example item, got %d", len(arr.Items))
	}
	item, ok := arr.Items[0].(*ir.NodeObject)
	if !ok {
		t.Fatalf("expected shell NodeObject item, got %T", arr.Items[0])
	}
	if item.Tag.RefTo != 1 {
		t.Fatalf("expected item ref_to=1, got %d", item.Tag.RefTo)
	}
}

// 自嵌套指针类型（type A *A）：本质是多级指针，允许报错，但绝不能爆栈。
func TestValueToNode_InfPointerType(t *testing.T) {
	type InfPtr *InfPtr
	_, err := ValueToNode(InfPtr(nil), "example")
	if err != nil {
		t.Logf("InfPtr rejected with error (acceptable, no overflow): %v", err)
	}
}

// map 值自引用。
func TestValueToNode_MapSelfReference(t *testing.T) {
	type Tree struct {
		Children map[string]Tree
	}
	node, err := ValueToNode(Tree{}, "example")
	if err != nil {
		t.Fatalf("Tree must encode without overflow, got err: %v", err)
	}
	root := node.(*ir.NodeObject)
	if root.Tag.RefID != 1 {
		t.Fatalf("expected root ref_id=1, got %d", root.Tag.RefID)
	}
	for _, f := range root.Fields {
		if f.Key == "children" {
			child := f.Value.(*ir.NodeObject)
			if len(child.Fields) != 1 {
				t.Fatalf("expected 1 example field, got %d", len(child.Fields))
			}
			leaf := child.Fields[0].Value.(*ir.NodeObject)
			if leaf.Tag.RefTo != 1 {
				t.Fatalf("expected leaf ref_to=1, got %d", leaf.Tag.RefTo)
			}
		}
	}
}

// assertNoRefTags 递归断言整棵树不含任何 ref 标记。
func assertNoRefTags(t *testing.T, n ir.Node) {
	t.Helper()
	if n == nil {
		return
	}
	tag := n.GetTag()
	if tag != nil && (tag.RefID != 0 || tag.RefTo != 0) {
		t.Fatalf("unexpected ref tag on %T (path=%s): ref_id=%d ref_to=%d",
			n, n.GetPath(), tag.RefID, tag.RefTo)
	}
	switch v := n.(type) {
	case *ir.NodeObject:
		for _, f := range v.Fields {
			assertNoRefTags(t, f.Value)
		}
	case *ir.NodeArray:
		for _, item := range v.Items {
			assertNoRefTags(t, item)
		}
	}
}
