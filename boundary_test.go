package metamessage

import (
	"bytes"
	"strings"
	"testing"
)

// 边界回归（公开 API 层）：长度正好 Max2Byte(65535) 时必须无损往返。
//
// 修复前 5 个编码器（string/bytes/array/object/tag）的最后一档写成
// `length < Max2Byte`，length == 65535 时静默写 0 字节——编码不报错，但内容
// 整个消失（实测：65535 字节的字符串编码后 wire 只有 4 字节，解回来 0 个元素）。
// 这里用公开 API 做往返，确保「不报错」真的等价于「没丢数据」。
//
// 注意数组/对象限的是**序列化后的 payload 字节数**，不是元素个数；下面的尺寸
// 是实测标定出来的（含 key、tag 等协议开销），不要凭感觉改。

// wantWire = 1 字节 sign + 2 字节长度前缀 + 65535 字节负载。
const wantWire = 65535 + 3

func TestBoundaryStringRoundTrip(t *testing.T) {
	in := strings.Repeat("x", 65535)
	wire, err := EncodeFromValue(in, "")
	if err != nil {
		t.Fatalf("65535-byte string must encode, got: %v", err)
	}
	if len(wire) != wantWire {
		t.Fatalf("wire = %d bytes, want %d", len(wire), wantWire)
	}
	var out string
	if err := DecodeToValue(wire, &out); err != nil {
		t.Fatalf("DecodeToValue: %v", err)
	}
	if out != in {
		t.Fatalf("string corrupted: got %d bytes, want %d", len(out), len(in))
	}
}

func TestBoundaryBytesRoundTrip(t *testing.T) {
	// []byte 必须显式 type=bytes，否则会被当成「每个元素一个字节」的数组
	// （那条路 2 字节/元素，永远撞不到奇数 65535）。
	// 对象 payload = wire("b") + [1 sign + 2 len + 65529 payload + 1 tag] = 65535。
	type bytesTyped struct {
		B []byte `mm:"type=bytes"`
	}
	in := bytesTyped{B: bytes.Repeat([]byte("y"), 65529)}
	wire, err := EncodeFromValue(in, "")
	if err != nil {
		t.Fatalf("[]byte field with object payload exactly %d must encode, got: %v", 65535, err)
	}
	if len(wire) != wantWire {
		t.Fatalf("wire = %d bytes, want %d", len(wire), wantWire)
	}
	var out bytesTyped
	if err := DecodeToValue(wire, &out); err != nil {
		t.Fatalf("DecodeToValue: %v", err)
	}
	if !bytes.Equal(out.B, in.B) {
		t.Fatalf("[]byte corrupted: got %d bytes, want %d", len(out.B), len(in.B))
	}
}

func TestBoundaryArrayPayloadRoundTrip(t *testing.T) {
	// 每个元素走 2 字节长度档：wire = 1(sign) + 2(len) + 768 = 771
	// 85 × 771 = 65535，正好顶到数组 payload 上限。
	// 顶层 tag 必须留空——tag 字节也计入 payload（"allow_empty" 会多加 5 字节）。
	const n, each = 85, 768
	in := make([]string, n)
	for i := range in {
		in[i] = strings.Repeat("z", each)
	}
	wire, err := EncodeFromValue(in, "")
	if err != nil {
		t.Fatalf("array with payload exactly %d must encode, got: %v", 65535, err)
	}
	if len(wire) != wantWire {
		t.Fatalf("wire = %d bytes, want %d", len(wire), wantWire)
	}
	var out []string
	if err := DecodeToValue(wire, &out); err != nil {
		t.Fatalf("DecodeToValue: %v", err)
	}
	if len(out) != n {
		t.Fatalf("decoded %d items, want %d", len(out), n)
	}
	for i, s := range out {
		if len(s) != each {
			t.Fatalf("item[%d] = %d bytes, want %d", i, len(s), each)
		}
	}
}

func TestBoundaryObjectPayloadRoundTrip(t *testing.T) {
	// 单字段结构体，对象 payload = wire("a") + wire(65529 B) = 2 + 65533 = 65535。
	type oneField struct {
		A string
	}
	in := oneField{A: strings.Repeat("w", 65529)}
	wire, err := EncodeFromValue(in, "")
	if err != nil {
		t.Fatalf("object with payload exactly %d must encode, got: %v", 65535, err)
	}
	if len(wire) != wantWire {
		t.Fatalf("wire = %d bytes, want %d", len(wire), wantWire)
	}
	var out oneField
	if err := DecodeToValue(wire, &out); err != nil {
		t.Fatalf("DecodeToValue: %v", err)
	}
	if out.A != in.A {
		t.Fatalf("field corrupted: got %d bytes, want %d", len(out.A), len(in.A))
	}
}

// TestBoundaryPlusOneStillRejected 超过 65535 照旧报错，且错误信息可读——
// 防止修复时把守卫一起放宽。
func TestBoundaryPlusOneStillRejected(t *testing.T) {
	if _, err := EncodeFromValue(strings.Repeat("x", 65536), ""); err == nil {
		t.Fatal("string of 65536 bytes must be rejected")
	} else if !strings.Contains(err.Error(), "string too long") {
		t.Fatalf("string: error = %q, want \"string too long\"", err)
	}

	// 数组 payload 超限：86×768 = 66288 > 65535
	in := make([]string, 86)
	for i := range in {
		in[i] = strings.Repeat("z", 768)
	}
	if _, err := EncodeFromValue(in, ""); err == nil {
		t.Fatal("array with payload 66288 must be rejected")
	} else if !strings.Contains(err.Error(), "array too long") {
		t.Fatalf("array: error = %q, want \"array too long\"", err)
	}
}
