package core

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// v0.3.0 长度档位边界测试。
//
// 背景（D6 规则）：每个编码器的最后一档必须写成 `<=`。历史教训是 string/bytes/
// array/object/tag 的最后一档曾写成 `length < Max2Byte`，而守卫是
// `length > Max2Byte`，于是 length 正好等于上限时守卫不报错、所有分支又全部
// 落空 ⇒ 函数静默返回 (0, nil)，一个字节都不写，数据无声丢失。
//
// v0.3.0 把长度槽扩到 4 字节（上限 Max4Byte = 0xFFFFFFFF ≈ 4.29GB）。档位分配
// 见 constants.go：
//   - String / Bytes / Tag：内联 0..27，28/29/30/31 = 1/2/3/4 字节大端长度
//   - Container：内联 0..11，12/13/14/15 = 1/2/3/4 字节大端长度
//
// 下面的表驱动用例对每个档位取「本档最大值」与「下一档最小值」两个临界点，
// 校验：suffix 档位标记正确、长度字段是大端且能解回原长度、负载紧跟其后。
//
// 已知覆盖边界：上限 Max4Byte(0xFFFFFFFF) 本身需要 4GB 内存，无法在 CI 里
// 真实分配，因此「length == Max4Byte 不得落空」这一条靠最后一档统一使用 `<=`
// 的约定 + code review 保证，下面只验证到 4 字节档的入口（Max3Byte+1）。

// wantTier 返回 length 应落到的档位标记与长度字段字节数。
// kind 取 "scalar"（String/Bytes/Tag，5-bit suffix）或 "container"（4-bit 长度域）。
func wantTier(kind string, length int) (marker byte, n int) {
	inlineCap, m1 := 27, StringLen1Byte
	if kind == "container" {
		inlineCap, m1 = 11, ContainerLen1Byte
	}
	switch {
	case length <= inlineCap:
		return byte(length), 0
	case length <= Max1Byte:
		return byte(m1), 1
	case length <= Max2Byte:
		return byte(m1) + 1, 2
	case length <= Max3Byte:
		return byte(m1) + 2, 3
	default:
		return byte(m1) + 3, 4
	}
}

func tierMask(kind string) byte {
	if kind == "container" {
		return byte(ContainerLenMask)
	}
	return byte(StringLenMask)
}

// checkTierPrefix 校验 e.buf 里的档位标记 + 大端长度字段 + 负载布局。
func checkTierPrefix(t *testing.T, e *encoder, kind string, length int, payload []byte) {
	t.Helper()
	marker, n := wantTier(kind, length)
	if got := e.buf[0] & tierMask(kind); got != marker {
		t.Fatalf("suffix = %d, want %d (length %d)", got, marker, length)
	}
	off := 1
	if n > 0 {
		var gotLen int
		for i := 0; i < n; i++ {
			gotLen = gotLen<<8 | int(e.buf[off+i])
		}
		if gotLen != length {
			t.Fatalf("length field = %d, want %d", gotLen, length)
		}
		off += n
	}
	if !bytes.Equal(e.buf[off:off+len(payload)], payload) {
		t.Fatal("payload bytes not written right after the header/length field")
	}
}

// 每个档位取「本档最大值」与「下一档最小值」两个临界点。
var tierBoundaryLengths = []int{
	27, 28,                  // 内联 ↔ 1 字节
	Max1Byte, Max1Byte + 1,  // 1 字节 ↔ 2 字节
	Max2Byte, Max2Byte + 1,  // 2 字节 ↔ 3 字节
	Max3Byte, Max3Byte + 1,  // 3 字节 ↔ 4 字节
}

func newBufEncoder() *encoder {
	// w 传 nil：只往内部 buf 写，不落盘
	return NewEncoder(nil)
}

func TestEncodeStringTierBoundaries(t *testing.T) {
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			s := strings.Repeat("x", l)

			n, err := e.encodeString(s)
			if err != nil {
				t.Fatalf("encodeString(%d) must not error, got: %v", l, err)
			}
			_, nb := wantTier("scalar", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}
			checkTierPrefix(t, e, "scalar", l, []byte(s))
		})
	}
}

func TestEncodeBytesTierBoundaries(t *testing.T) {
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			bs := bytes.Repeat([]byte("y"), l)

			n, err := e.encodeBytes(bs)
			if err != nil {
				t.Fatalf("encodeBytes(%d) must not error, got: %v", l, err)
			}
			_, nb := wantTier("scalar", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}
			checkTierPrefix(t, e, "scalar", l, bs)
		})
	}
}

func TestEncodeArrayPayloadTierBoundaries(t *testing.T) {
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			payload := bytes.Repeat([]byte("z"), l)

			n, err := e.encodeArray(payload)
			if err != nil {
				t.Fatalf("encodeArray(payload %d) must not error, got: %v", l, err)
			}
			_, nb := wantTier("container", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}
			checkTierPrefix(t, e, "container", l, payload)
		})
	}
}

func TestEncodeObjectPayloadTierBoundaries(t *testing.T) {
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			payload := bytes.Repeat([]byte("w"), l)

			n, err := e.encodeObject(payload)
			if err != nil {
				t.Fatalf("encodeObject(payload %d) must not error, got: %v", l, err)
			}
			_, nb := wantTier("container", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}
			checkTierPrefix(t, e, "container", l, payload)
		})
	}
}

func TestEncodeTagPayloadTierBoundaries(t *testing.T) {
	tag := []byte("tt")
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			// tag + payload 合计正好 l
			payload := bytes.Repeat([]byte("p"), l-len(tag))

			n, err := e.encodeTag(payload, tag)
			if err != nil {
				t.Fatalf("encodeTag(tag %d + payload %d) must not error, got: %v",
					len(tag), len(payload), err)
			}
			_, nb := wantTier("scalar", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}

			// 布局：sign(1) + 长度字段(n) + tag + payload
			marker, n2 := wantTier("scalar", l)
			if n2 != nb {
				t.Fatalf("tier byte width mismatch: %d vs %d", n2, nb)
			}
			if got := e.buf[0] & tierMask("scalar"); got != marker {
				t.Fatalf("suffix = %d, want %d", got, marker)
			}
			off := 1
			if nb > 0 {
				var gotLen int
				for i := 0; i < nb; i++ {
					gotLen = gotLen<<8 | int(e.buf[off+i])
				}
				if gotLen != l {
					t.Fatalf("length field = %d, want %d", gotLen, l)
				}
				off += nb
			}
			if !bytes.Equal(e.buf[off:off+len(tag)], tag) {
				t.Fatal("tag bytes not written right after the length field")
			}
			if !bytes.Equal(e.buf[off+len(tag):off+len(tag)+len(payload)], payload) {
				t.Fatal("payload bytes not written after the tag")
			}
		})
	}
}

// TestEncodeAtTierCeilingStillWorks 守住「长度正好等于某档上限时不得落空」
// 这条历史 bug。Max2Byte(65535) 与 Max3Byte(16777215) 可真实分配，逐个校验。
func TestEncodeAtTierCeilingStillWorks(t *testing.T) {
	for _, l := range []int{Max2Byte, Max3Byte} {
		e := newBufEncoder()
		s := strings.Repeat("x", l)
		n, err := e.encodeString(s)
		if err != nil {
			t.Fatalf("encodeString(%d) must not error, got: %v", l, err)
		}
		if n == 0 {
			t.Fatalf("encodeString(%d) silently wrote 0 bytes (D6 regression)", l)
		}
		checkTierPrefix(t, e, "scalar", l, []byte(s))
	}
}
