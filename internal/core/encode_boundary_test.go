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
// v0.3.0 把 **Bytes / Container / Tag** 的长度槽扩到 4 字节
// （上限 Max4Byte = 0xFFFFFFFF ≈ 4.29GB）。档位分配见 constants.go：
//   - Bytes / Tag：内联 0..27，28/29/30/31 = 1/2/3/4 字节大端长度
//   - Container：内联 0..11，12/13/14/15 = 1/2/3/4 字节大端长度
//
// **String 不参与本次扩展**，维持 v0.2.x 语义：内联 0..29，30/31 = 1/2 字节，
// 上限 Max2Byte(65535)。注意 String 的 1 字节档边界是 `< Max1Byte`
// （即 255 落 2 字节档），这是 v0.2.x 的原样行为，不要"顺手"改成 `<=`，
// 否则会与尚未改动的其他语言移植字节不一致。
//
// 已知覆盖边界：上限 Max4Byte(0xFFFFFFFF) 本身需要 4GB 内存，无法在 CI 里
// 真实分配，因此「length == Max4Byte 不得落空」这一条靠最后一档统一使用 `<=`
// 的约定 + code review 保证，下面只验证到 4 字节档的入口（Max3Byte+1）。

// wantTier 返回 length 应落到的档位标记与长度字段字节数。
// kind 取：
//   - "string"   String，v0.2.x 语义，最多 2 字节长度
//   - "scalar"   Bytes / Tag，v0.3.0 四档
//   - "container" Container，v0.3.0 四档
func wantTier(kind string, length int) (marker byte, n int) {
	if kind == "string" {
		// String 维持 v0.2.x：内联 0..29，30 = 1 字节，31 = 2 字节。
		// 1 字节档边界是 `< Max1Byte`（255 落 2 字节档），与 v0.2.x 完全一致。
		switch {
		case length < StringLen1Byte:
			return byte(length), 0
		case length < Max1Byte:
			return byte(StringLen1Byte), 1
		default:
			return byte(StringLen2Byte), 2
		}
	}

	inlineCap, m1 := 27, BytesLen1Byte
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
	switch kind {
	case "string":
		return byte(StringLenMask)
	case "container":
		return byte(ContainerLenMask)
	default:
		return byte(BytesLenMask)
	}
}

// checkTierSuffix 校验 e.buf 里的尾定界布局：[payload][大端长度字段][sign]。
// v0.5.0 编码器一趟正向写出：body 先落盘，定界信息（大端长度 + sign）追加在尾部，
// 因此 sign 位于节点最后一个字节，长度字段紧贴其前。
func checkTierSuffix(t *testing.T, e *encoder, kind string, length int, payload []byte) {
	t.Helper()
	marker, n := wantTier(kind, length)
	total := int(e.offset)
	if got := e.buf[total-1] & tierMask(kind); got != marker {
		t.Fatalf("sign tier = %d, want %d (kind %s, length %d)", got, marker, kind, length)
	}
	if n > 0 {
		var gotLen int
		for i := 0; i < n; i++ {
			gotLen = gotLen<<8 | int(e.buf[total-1-n+i])
		}
		if gotLen != length {
			t.Fatalf("length field = %d, want %d", gotLen, length)
		}
	}
	if !bytes.Equal(e.buf[:len(payload)], payload) {
		t.Fatal("payload bytes not written at the head of the node (before tail delimiters)")
	}
}

// Bytes / Container / Tag 的档位临界点：每个档位取「本档最大值」与
// 「下一档最小值」两个点。
var tierBoundaryLengths = []int{
	27, 28,                  // 内联 ↔ 1 字节
	Max1Byte, Max1Byte + 1,  // 1 字节 ↔ 2 字节
	Max2Byte, Max2Byte + 1,  // 2 字节 ↔ 3 字节
	Max3Byte, Max3Byte + 1,  // 3 字节 ↔ 4 字节
}

// String 只有内联 / 1 字节 / 2 字节三档，2 字节是最后一档。
// 254 → 1 字节，255 → 2 字节（v0.2.x 的 `<` 边界，特意钉住）。
var stringTierBoundaryLengths = []int{
	29, 30,                    // 内联 ↔ 1 字节
	Max1Byte - 1, Max1Byte,    // 1 字节 ↔ 2 字节
	Max2Byte,                  // 最后一档上限，必须仍能编码
}

func newBufEncoder() *encoder {
	// w 传 nil：只往内部 buf 写，不落盘
	return NewEncoder(nil)
}

// ---- String：维持 v0.2.x，最多 2 字节长度 ----

func TestEncodeStringTierBoundaries(t *testing.T) {
	for _, l := range stringTierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			s := strings.Repeat("x", l)

			n, err := e.encodeString(s)
			if err != nil {
				t.Fatalf("encodeString(%d) must not error, got: %v", l, err)
			}
			_, nb := wantTier("string", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}
			checkTierSuffix(t, e, "string", l, []byte(s))
		})
	}
}

// TestEncodeStringAboveMax2ByteStillRejected String 不参与 4GB 扩展，
// 超过 Max2Byte(65535) 必须照旧报错——防止把 String 的守卫一起放宽。
func TestEncodeStringAboveMax2ByteStillRejected(t *testing.T) {
	if _, err := newBufEncoder().encodeString(strings.Repeat("x", Max2Byte+1)); err == nil {
		t.Fatal("string of Max2Byte+1 bytes must still be rejected")
	}
}

// ---- Bytes / Container / Tag：v0.3.0 四档 ----

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
			checkTierSuffix(t, e, "scalar", l, bs)
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
			checkTierSuffix(t, e, "container", l, payload)
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
			checkTierSuffix(t, e, "container", l, payload)
		})
	}
}

func TestEncodeTagPayloadTierBoundaries(t *testing.T) {
	meta := []byte("tt")
	slot, err := footerSlotBytes(len(meta))
	if err != nil {
		t.Fatalf("footerSlotBytes(%d): %v", len(meta), err)
	}
	for _, l := range tierBoundaryLengths {
		t.Run(fmt.Sprintf("len_%d", l), func(t *testing.T) {
			e := newBufEncoder()
			// 尾定界 tag 节点 total = len(inner)+len(meta)+len(slot)，钉住 total == l
			inner := bytes.Repeat([]byte("p"), l-len(meta)-len(slot))

			n, err := e.encodeTag(inner, meta)
			if err != nil {
				t.Fatalf("encodeTag(meta %d + inner %d) must not error, got: %v",
					len(meta), len(inner), err)
			}
			marker, nb := wantTier("scalar", l)
			if want := uint32(1 + nb + l); n != want || e.offset != want {
				t.Fatalf("wrote %d bytes (offset %d), want %d", n, e.offset, want)
			}

			// 尾定界布局：[inner][reverse(meta)][slot][大端长度字段][sign]
			total := int(e.offset)
			if got := e.buf[total-1] & tierMask("scalar"); got != marker {
				t.Fatalf("sign tier = %d, want %d", got, marker)
			}
			off := total - 1
			if nb > 0 {
				var gotLen int
				for i := 0; i < nb; i++ {
					gotLen = gotLen<<8 | int(e.buf[off-nb+i])
				}
				if gotLen != l {
					t.Fatalf("length field = %d, want %d", gotLen, l)
				}
				off -= nb
			}
			if !bytes.Equal(e.buf[off-len(slot):off], slot) {
				t.Fatal("footer slot not written before the length field")
			}
			off -= len(slot)
			if !bytes.Equal(e.buf[off-len(meta):off], reverseBytes(meta)) {
				t.Fatal("reversed meta not written before the footer slot")
			}
			off -= len(meta)
			if !bytes.Equal(e.buf[:off], inner) {
				t.Fatal("inner bytes not written at the head of the node")
			}
		})
	}
}

// TestEncodeAtTierCeilingStillWorks 守住「长度正好等于某档上限时不得落空」
// 这条历史 bug。String 的上限仍是 Max2Byte，Bytes 的 4 字节档入口用
// Max3Byte+1 验证。
func TestEncodeAtTierCeilingStillWorks(t *testing.T) {
	e := newBufEncoder()
	s := strings.Repeat("x", Max2Byte)
	n, err := e.encodeString(s)
	if err != nil {
		t.Fatalf("encodeString(%d) must not error, got: %v", Max2Byte, err)
	}
	if n == 0 {
		t.Fatalf("encodeString(%d) silently wrote 0 bytes (D6 regression)", Max2Byte)
	}
	checkTierSuffix(t, e, "string", Max2Byte, []byte(s))

	be := newBufEncoder()
	bs := bytes.Repeat([]byte("y"), Max3Byte+1)
	bn, err := be.encodeBytes(bs)
	if err != nil {
		t.Fatalf("encodeBytes(%d) must not error, got: %v", Max3Byte+1, err)
	}
	if bn == 0 {
		t.Fatalf("encodeBytes(%d) silently wrote 0 bytes (D6 regression)", Max3Byte+1)
	}
	checkTierSuffix(t, be, "scalar", Max3Byte+1, bs)
}
