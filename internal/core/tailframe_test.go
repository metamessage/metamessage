package core

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
	"github.com/metamessage/metamessage/internal/jsonc"
)

// v0.5.0 尾定界测试：编码器一趟正向直接产出尾定界 wire（档位/内联不丢），
// 解码入口 Untailify 从尾反向剥节点还原头部形态后走原正向解码器。

func hx(s string) []byte {
	bs, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bs
}

// TestTailFloatInline 内联档 float 保持单字节 wire（v0.3.0 的 0x71 回归）。
func TestTailFloatInline(t *testing.T) {
	bs, err := FromValue(-0.1, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(bs) != 1 || bs[0] != 0x71 {
		t.Fatalf("want single 0x71, got %x", bs)
	}
	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v := node.(*ir.NodeScalar).Data.(float64); v != -0.1 {
		t.Fatalf("decode value = %v", v)
	}
}

// TestTailIntInline 内联档 int 保持单字节 wire。
func TestTailIntInline(t *testing.T) {
	bs, err := FromValue(5, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	// 无 tag 时正文本 [0x05]；内联 int sign = PositiveInt|5
	if len(bs) != 1 || bs[0] != byte(PositiveInt|Prefix(5)) {
		t.Fatalf("want single 0x25, got %x", bs)
	}
	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v := node.(*ir.NodeScalar).Data.(int); v != 5 {
		t.Fatalf("decode value = %v", v)
	}
}

// TestTailIntTier int 带体节点：[body 原序][sign]。
func TestTailIntTier(t *testing.T) {
	bs, err := FromValue(300, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	// head: [sign=PositiveInt|IntLen2Byte][0x01][0x2C] → wire: [0x01][0x2C][sign]
	want := byte(PositiveInt | Prefix(IntLen2Byte))
	if len(bs) != 3 || bs[2] != want || bs[0] != 0x01 || bs[1] != 0x2C {
		t.Fatalf("wire = %x, want 01 2c %02x", bs, want)
	}
	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v := node.(*ir.NodeScalar).Data.(int); v != 300 {
		t.Fatalf("decode value = %v", v)
	}
}

// TestTailStringTier 字符串带档位：[body][长度 BE][sign]。
func TestTailStringTier(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" // 52 字节 > 29
	bs, err := FromValue(long, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	sign := byte(PrefixString | Prefix(StringLen1Byte))
	if len(bs) != len(long)+2 || bs[len(bs)-1] != sign || bs[len(bs)-2] != byte(len(long)) {
		t.Fatalf("wire tail wrong: len=%d last=%02x prev=%02x", len(bs), bs[len(bs)-1], bs[len(bs)-2])
	}
	node, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s := node.(*ir.NodeScalar).Text; s != long {
		t.Fatalf("decode text mismatch")
	}
}

// TestTailContainerChildrenForward 容器子节点 wire 形态正向拼接在 payload 区。
func TestTailContainerChildrenForward(t *testing.T) {
	src := `{"a":1,"b":2}`
	bs, err := FromJSONC(src)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	back, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode %x: %v", bs, err)
	}
	obj := back.(*ir.NodeObject)
	if len(obj.Fields) != 2 || obj.Fields[0].Key != "a" || obj.Fields[1].Key != "b" {
		t.Fatalf("fields order wrong: %+v", obj.Fields)
	}
}

func encodeNode(node ir.Node) ([]byte, error) {
	e := NewEncoder(nil)
	return e.Encode(node)
}

// TestTailTagWrapper tag 包装段换位：[inner][reverse(meta)][slot][tier BE][sign]。
func TestTailTagWrapper(t *testing.T) {
	src := "// mm: desc=hello\n{\"v\": \"x\"}"
	bs, err := FromJSONC(src)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	// 根节点带 meta，尾字节必须是 PrefixTag sign
	if len(bs) == 0 || GetPrefix(bs[len(bs)-1]) != PrefixTag {
		t.Fatalf("want tag-wrapped wire, got %x", bs)
	}
	back, err := Decode(bs)
	if err != nil {
		t.Fatalf("decode %x: %v", bs, err)
	}
	out := jsonc.ToJSONC(back)
	if !strings.Contains(out, "hello") {
		t.Fatalf("tag lost in roundtrip: %s", out)
	}
}

// TestTailUntailBoundary 篡改防御：坏边界必须报错而不是静默错位。
func TestTailUntailBoundary(t *testing.T) {
	bs, err := FromValue(300, "")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	truncated := bs[:len(bs)-1] // 砍掉尾 sign
	if _, err := Decode(truncated); err == nil {
		t.Fatalf("want error on truncated wire")
	}
}
