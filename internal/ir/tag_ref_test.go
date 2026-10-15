package ir

import (
	"bytes"
	"strings"
	"testing"
)

// KMore ref meta 的字节级黄金断言。
//
// wire 布局：[KMore | (payload总字节数-1)] [meta 位图] [每个出现值 1 字节]
// 两字段全出现：payload = 1 + 2 = 3 字节 → 低 3 bit = 2。
func TestTagBytesRefMetaAllFields(t *testing.T) {
	tag := &Tag{RefID: 1, RefTo: 2}
	got, err := tag.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}
	want := []byte{
		0xF8 | 2, // KMore | (3-1)
		0x03,     // meta = bit0|bit1
		0x01,     // ref_id
		0x02,     // ref_to
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected %x, got %x", want, got)
	}
}

func TestTagBytesRefMetaSingleField(t *testing.T) {
	// 单字段：payload = 1 + 1 = 2 字节 → 低 3 bit = 1
	tag := &Tag{RefTo: 44}
	got, err := tag.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}
	want := []byte{
		0xF8 | 1, // KMore | (2-1)
		0x02,     // meta = bit1
		0x2C,     // ref_to = 44
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected %x, got %x", want, got)
	}
}

// 边界值 MaxRefID 必须正常编码（u8 上限本身合法）。
func TestTagBytesRefMetaMaxValue(t *testing.T) {
	tag := &Tag{RefID: MaxRefID, RefTo: MaxRefID}
	got, err := tag.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}
	want := []byte{
		0xF8 | 2, // KMore | (3-1)
		0x03,     // meta = bit0|bit1
		0xFF,     // ref_id = 255
		0xFF,     // ref_to = 255
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected %x, got %x", want, got)
	}
}

// ref_id/ref_to 超过 MaxRefID(u8) 必须显式报错，绝不静默截断。
func TestTagBytesRefMetaOverflow(t *testing.T) {
	cases := []struct {
		name          string
		refID, refTo  int
	}{
		{"ref_to over", 0, MaxRefID + 1},
		{"ref_to large", 0, 300},
		{"ref_id over", MaxRefID + 1, 0},
		{"both over", 1000, 2000},
		{"negative", -1, 0},
	}
	for _, tc := range cases {
		tag := &Tag{RefID: tc.refID, RefTo: tc.refTo}
		if _, err := tag.Bytes(); err == nil {
			t.Fatalf("%s: expected overflow error for ref_id=%d ref_to=%d, got nil",
				tc.name, tc.refID, tc.refTo)
		}
	}
}

func TestTagBytesRefMetaZeroOmitted(t *testing.T) {
	// 两字段全 0：不写任何字节（零开销）
	tag := &Tag{}
	got, err := tag.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty bytes for zero ref fields, got %x", got)
	}
}

func TestTagBytesRefMetaMoreNotWritten(t *testing.T) {
	// More 已废弃：即使赋值也不编码（避免双 KMore 歧义）
	tag := &Tag{More: 5}
	got, err := tag.Bytes()
	if err != nil {
		t.Fatalf("Bytes() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected More to never be encoded, got %x", got)
	}
}

func TestTagToStringRefFields(t *testing.T) {
	tag := &Tag{RefID: 1, RefTo: 2}
	got := tag.ToString()
	want := "ref_id=1; ref_to=2"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTagToStringRefFieldsOmittedWhenZero(t *testing.T) {
	if got := (&Tag{}).ToString(); strings.Contains(got, "ref_") {
		t.Fatalf("zero ref fields must be omitted, got %q", got)
	}
}

func TestTagToStringParseMMTagRefRoundTrip(t *testing.T) {
	in := &Tag{RefID: 4, RefTo: 5}
	out, err := ParseMMTag(in.ToString())
	if err != nil {
		t.Fatalf("ParseMMTag failed: %v", err)
	}
	if out.RefID != in.RefID {
		t.Fatalf("ref_id: expected %d, got %d", in.RefID, out.RefID)
	}
	if out.RefTo != in.RefTo {
		t.Fatalf("ref_to: expected %d, got %d", in.RefTo, out.RefTo)
	}
}

func TestParseMMTagRefFieldsDirect(t *testing.T) {
	out, err := ParseMMTag("ref_id=7; ref_to=8")
	if err != nil {
		t.Fatalf("ParseMMTag failed: %v", err)
	}
	if out.RefID != 7 || out.RefTo != 8 {
		t.Fatalf("expected 7/8, got %d/%d", out.RefID, out.RefTo)
	}
}

// 文本侧 ref_id/ref_to 越界（0 或 > MaxRefID）必须报错，与二进制编码上限一致。
func TestParseMMTagRefFieldsOutOfRange(t *testing.T) {
	for _, s := range []string{"ref_id=0", "ref_id=256", "ref_to=0", "ref_to=300"} {
		if _, err := ParseMMTag(s); err == nil {
			t.Fatalf("expected out-of-range error for %q, got nil", s)
		}
	}
}

func TestMergeTagRefFields(t *testing.T) {
	dst := NewTag()
	src := &Tag{RefID: 1, RefTo: 2}
	MergeTag(dst, src)
	if dst.RefID != 1 || dst.RefTo != 2 {
		t.Fatalf("merge failed: %d/%d", dst.RefID, dst.RefTo)
	}

	// 零值不覆盖
	dst2 := NewTag()
	dst2.RefID = 9
	MergeTag(dst2, &Tag{RefID: 0, RefTo: 0})
	if dst2.RefID != 9 {
		t.Fatalf("zero src must not overwrite, got %d", dst2.RefID)
	}
}
