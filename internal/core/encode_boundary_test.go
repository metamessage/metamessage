package core

import (
	"bytes"
	"strings"
	"testing"
)

// 边界回归：长度正好 Max2Byte(65535) 时必须能编码。
//
// 曾经 encodeString / encodeBytes / encodeArray / encodeObject / encodeTag 的
// 最后一档都写成 `length < Max2Byte`，而守卫是 `length > Max2Byte`。于是
// length == 65535 时：守卫不报错，三个分支又全部落空 ⇒ 函数静默返回 (0, nil)，
// 一个字节都不写。调用方拿不到错误，数据无声丢失（实测：65535 字节的字符串
// 编码后 wire 只有 4 字节，解回来 0 个元素）。
//
// 解码端其实能表达 65535——2 字节长度字段 `int(l[0])<<8 | int(l[1])` 最大就是
// 65535（见 decode.go 的 decodeString/decodeBytes）。所以这是纯编码端 off-by-one，
// 修法是把最后一档的 `<` 改成 `<=`。
//
// 每个用例除了「写够了字节」，还校验长度前缀本身能被解回 65535。

// wantWire 是 65535 字节负载加上 1 字节 sign + 2 字节长度前缀后的总字节数。
const wantWire = Max2Byte + 3

func newBufEncoder() *encoder {
	// w 传 nil：只往内部 buf 写，不落盘
	return NewEncoder(nil)
}

// checkLenPrefix 校验 e.buf[1..2] 是大端 65535，且负载紧随其后。
func checkLenPrefix(t *testing.T, e *encoder, payload []byte) {
	t.Helper()
	if got := int(e.buf[1])<<8 | int(e.buf[2]); got != Max2Byte {
		t.Fatalf("2-byte length prefix = %d, want %d", got, Max2Byte)
	}
	if !bytes.Equal(e.buf[3:3+len(payload)], payload) {
		t.Fatal("payload bytes not written right after the length prefix")
	}
}

func TestEncodeStringAtMax2Byte(t *testing.T) {
	e := newBufEncoder()
	s := strings.Repeat("x", Max2Byte)

	n, err := e.encodeString(s)
	if err != nil {
		t.Fatalf("encodeString(%d bytes) must not error, got: %v", Max2Byte, err)
	}
	if n != wantWire || e.offset != wantWire {
		t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, wantWire)
	}
	checkLenPrefix(t, e, []byte(s))
}

func TestEncodeBytesAtMax2Byte(t *testing.T) {
	e := newBufEncoder()
	bs := bytes.Repeat([]byte("y"), Max2Byte)

	n, err := e.encodeBytes(bs)
	if err != nil {
		t.Fatalf("encodeBytes(%d bytes) must not error, got: %v", Max2Byte, err)
	}
	if n != wantWire || e.offset != wantWire {
		t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, wantWire)
	}
	checkLenPrefix(t, e, bs)
}

func TestEncodeArrayPayloadAtMax2Byte(t *testing.T) {
	e := newBufEncoder()
	payload := bytes.Repeat([]byte("z"), Max2Byte)

	n, err := e.encodeArray(payload)
	if err != nil {
		t.Fatalf("encodeArray(payload %d bytes) must not error, got: %v", Max2Byte, err)
	}
	if n != wantWire || e.offset != wantWire {
		t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, wantWire)
	}
	checkLenPrefix(t, e, payload)
}

func TestEncodeObjectPayloadAtMax2Byte(t *testing.T) {
	e := newBufEncoder()
	payload := bytes.Repeat([]byte("w"), Max2Byte)

	n, err := e.encodeObject(payload)
	if err != nil {
		t.Fatalf("encodeObject(payload %d bytes) must not error, got: %v", Max2Byte, err)
	}
	if n != wantWire || e.offset != wantWire {
		t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, wantWire)
	}
	checkLenPrefix(t, e, payload)
}

func TestEncodeTagPayloadAtMax2Byte(t *testing.T) {
	e := newBufEncoder()
	// tag + payload 合计正好 65535
	tag := bytes.Repeat([]byte("t"), 35)
	payload := bytes.Repeat([]byte("p"), Max2Byte-len(tag))

	n, err := e.encodeTag(payload, tag)
	if err != nil {
		t.Fatalf("encodeTag(tag %d + payload %d) must not error, got: %v",
			len(tag), len(payload), err)
	}
	if n != wantWire || e.offset != wantWire {
		t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, wantWire)
	}
	// 布局：sign(1) + len(2) + tag + payload
	if got := int(e.buf[1])<<8 | int(e.buf[2]); got != Max2Byte {
		t.Fatalf("2-byte length prefix = %d, want %d", got, Max2Byte)
	}
	if !bytes.Equal(e.buf[3:3+len(tag)], tag) {
		t.Fatal("tag bytes not written right after the length prefix")
	}
	if !bytes.Equal(e.buf[3+len(tag):], payload) {
		t.Fatal("payload bytes not written after the tag")
	}
}

// TestEncodeAtMax2BytePlusOneStillRejected 超过 65535 必须照旧报错——
// 防止把守卫一起放宽。
func TestEncodeAtMax2BytePlusOneStillRejected(t *testing.T) {
	if _, err := newBufEncoder().encodeString(strings.Repeat("x", Max2Byte+1)); err == nil {
		t.Fatal("string of Max2Byte+1 bytes must be rejected")
	}
	if _, err := newBufEncoder().encodeBytes(bytes.Repeat([]byte("y"), Max2Byte+1)); err == nil {
		t.Fatal("bytes of Max2Byte+1 must be rejected")
	}
	if _, err := newBufEncoder().encodeArray(bytes.Repeat([]byte("z"), Max2Byte+1)); err == nil {
		t.Fatal("array payload of Max2Byte+1 must be rejected")
	}
	if _, err := newBufEncoder().encodeObject(bytes.Repeat([]byte("w"), Max2Byte+1)); err == nil {
		t.Fatal("object payload of Max2Byte+1 must be rejected")
	}
	if _, err := newBufEncoder().encodeTag(
		bytes.Repeat([]byte("p"), Max2Byte), bytes.Repeat([]byte("t"), 2)); err == nil {
		t.Fatal("tag+payload of Max2Byte+2 must be rejected")
	}
}
