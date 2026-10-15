package core

import (
	"testing"

	"github.com/metamessage/metamessage/internal/ir"
)

// decodeTagBytes 对 encodeRefMeta 的逆运算：逐字节解析 ref meta。
// decoder 从 d.data 读取（见 Decode），故直接设置 data 后调用。
func decodeTagBytesForTest(t *testing.T, bs []byte) (*ir.Tag, int, error) {
	t.Helper()
	d := NewDecoder(nil)
	d.data = bs
	d.offset = 0
	tag := ir.NewTag()
	length, err := d.decodeTagBytes(tag)
	return tag, length, err
}

func TestDecodeTagBytesRefMetaAllFields(t *testing.T) {
	src := &ir.Tag{RefID: 1, RefTo: 2}
	bs, err := src.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}

	got, length, err := decodeTagBytesForTest(t, bs)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if length != len(bs) {
		t.Fatalf("expected length %d, got %d", len(bs), length)
	}
	if got.RefID != 1 || got.RefTo != 2 {
		t.Fatalf("expected ref_id=1 ref_to=2, got ref_id=%d ref_to=%d", got.RefID, got.RefTo)
	}
}

// 单字段边界值 MaxRefID：1 字节编码往返。
func TestDecodeTagBytesRefMetaSingleMax(t *testing.T) {
	src := &ir.Tag{RefTo: ir.MaxRefID}
	bs, err := src.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}

	got, _, err := decodeTagBytesForTest(t, bs)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if got.RefTo != ir.MaxRefID {
		t.Fatalf("expected ref_to=%d, got %d", ir.MaxRefID, got.RefTo)
	}
}

// 超过 MaxRefID：编码侧显式报错，绝不静默截断。
func TestEncodeRefMetaOverflow(t *testing.T) {
	src := &ir.Tag{RefTo: 300}
	if _, err := src.Bytes(); err == nil {
		t.Fatal("expected overflow error for ref_to=300, got nil")
	}
}

// 保留位（bit2 起）非零 → corrupt error。
func TestDecodeTagBytesRefMetaReservedBits(t *testing.T) {
	// [KMore|2][0x08][0x00 0x01]：meta=0x08 触碰保留位
	bs := []byte{0xF8 | 2, 0x08, 0x00, 0x01}
	if _, _, err := decodeTagBytesForTest(t, bs); err == nil {
		t.Fatal("expected error for reserved meta bits, got nil")
	}
}

// payload 长度与位图不符 → corrupt error。
func TestDecodeTagBytesRefMetaLengthMismatch(t *testing.T) {
	// meta=bit0|bit1（k=2，期望 payload 1+2=3 字节）但 sign 低 3 bit = 1（payload 2 字节）
	bs := []byte{0xF8 | 1, 0x03, 0x01}
	if _, _, err := decodeTagBytesForTest(t, bs); err == nil {
		t.Fatal("expected error for length mismatch, got nil")
	}
}

// 位图声明字段存在但取值为 0 → corrupt error（0 = 未设置，不可同时置位）。
func TestDecodeTagBytesRefMetaZeroValue(t *testing.T) {
	// [KMore|1][0x01][0x00]：meta 声明有 ref_id，但值字节为 0
	bs := []byte{0xF8 | 1, 0x01, 0x00}
	if _, _, err := decodeTagBytesForTest(t, bs); err == nil {
		t.Fatal("expected error for zero ref value, got nil")
	}
}

// 空 payload 防御。
func TestDecodeTagBytesRefMetaEmptyPayload(t *testing.T) {
	bs := []byte{0xF8 | 0}
	if _, _, err := decodeTagBytesForTest(t, bs); err == nil {
		t.Fatal("expected error for empty ref meta payload, got nil")
	}
}

// wire 端到端：ProbeNode 编码 → 解码，ref 标记完整往返。
func TestWire_RefTagRoundTrip(t *testing.T) {
	node, err := ValueToNode(ProbeNode{}, "example")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	encoder := getEncoder()
	defer putEncoder(encoder)

	wire, err := encoder.Encode(node)
	if err != nil {
		t.Fatalf("wire encode failed: %v", err)
	}

	decoded, err := Decode(wire)
	if err != nil {
		t.Fatalf("wire decode failed: %v", err)
	}

	root, ok := decoded.(*ir.NodeObject)
	if !ok {
		t.Fatalf("expected root NodeObject, got %T", decoded)
	}
	if root.Tag.RefID != 1 {
		t.Fatalf("expected decoded root ref_id=1, got %d", root.Tag.RefID)
	}

	for _, f := range root.Fields {
		switch v := f.Value.(type) {
		case *ir.NodeArray:
		if f.Key == "kids" {
			if len(v.Items) != 1 {
					t.Fatalf("expected 1 decoded item, got %d", len(v.Items))
				}
				item, ok := v.Items[0].(*ir.NodeObject)
				if !ok {
					t.Fatalf("expected decoded item NodeObject, got %T", v.Items[0])
				}
				if item.Tag.RefTo != 1 {
					t.Fatalf("expected decoded item ref_to=1, got %d", item.Tag.RefTo)
				}
			}
		case *ir.NodeObject:
			if f.Key == "ref" && v.Tag.RefTo != 1 {
				t.Fatalf("expected decoded ref ref_to=1, got %d", v.Tag.RefTo)
			}
		}
	}
}

// 无环数据的 wire 往返不携带 ref 标记。
func TestWire_AcyclicNoRefMeta(t *testing.T) {
	type plain struct {
		Name string
		Arr  []int
	}
	node, err := ValueToNode(plain{Name: "x", Arr: []int{1}}, "")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	encoder := getEncoder()
	defer putEncoder(encoder)

	wire, err := encoder.Encode(node)
	if err != nil {
		t.Fatalf("wire encode failed: %v", err)
	}

	// wire 中不应出现 KMore 头字节后跟合法 meta 的模式；
	// 直接解码验证 tag 字段为零。
	decoded, err := Decode(wire)
	if err != nil {
		t.Fatalf("wire decode failed: %v", err)
	}
	assertNoRefTags(t, decoded)
}
