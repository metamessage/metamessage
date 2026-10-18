import { encodeNode, decode, fromJSONC, decodeToJSONC, mm } from '../../src/core/index.js';
import { ValueType } from '../../src/ir/value-type.js';
import { Tag } from '../../src/ir/tag.js';
import { NodeScalar } from '../../src/ir/ast.js';
import { footerSlotBytes, parseFooterSlot, reverseBytes } from '../../src/core/footer.js';

// v0.5.0 尾定界模式专项测试：golden hex / 尾部槽三档 / 篡改防御。
// 期望字节与 Go internal/core 参考实现完全一致（编码器直接尾序写出，无出口反转）。
// 尾序布局：[inner wire][reverse(meta)][meta_len 尾槽][total 档位大端][tag sign]。

function toHex(bs: Uint8Array): string {
  return Array.from(bs)
    .map((b) => b.toString(16).padStart(2, '0').toUpperCase())
    .join(' ');
}

describe('MM Footer (v0.5.0)', () => {
  test('slot tiers: canonical widths', () => {
    const cases: [number, number][] = [
      [1, 1],
      [253, 1],
      [254, 2],
      [255, 2],
      [256, 3],
      [65535, 3],
    ];
    for (const [metaLen, width] of cases) {
      const slot = footerSlotBytes(metaLen);
      expect(slot.length).toBe(width);
      const data = new Uint8Array([0x00, 0x11, 0x22, 0x33, ...slot]);
      const [gotLen, gotW] = parseFooterSlot(data, data.length);
      expect(gotLen).toBe(metaLen);
      expect(gotW).toBe(width);
    }
  });

  test('slot too long (> 65535) rejected', () => {
    expect(() => footerSlotBytes(65536)).toThrow(/tag too long/);
  });

  test('golden single field: str "hi", desc="ab"', () => {
    const tag = new Tag();
    tag.type = ValueType.Str;
    tag.desc = 'ab';
    const got = encodeNode(mm.str('hi', tag));
    // tail-order: inner body [68 69 82], reverse(meta) [62 61 1A], slot [03], sign E7
    expect(toHex(got)).toBe('68 69 82 62 61 1A 03 E7');
  });

  test('golden multi field order: Example+Nullable+Desc', () => {
    const tag = new Tag();
    tag.type = ValueType.Str;
    tag.example = true;
    tag.nullable = true;
    tag.desc = 'ab';
    const got = encodeNode(mm.str('x', tag));
    expect(toHex(got)).toBe('78 81 62 61 1A 29 09 05 E8');
    const back = decode(got) as NodeScalar;
    expect(back.getTag().example).toBe(true);
    expect(back.getTag().nullable).toBe(true);
    expect(back.getTag().desc).toBe('ab');
  });

  test('golden ref meta at tail: refTo=2', () => {
    const tag = new Tag();
    tag.type = ValueType.Str;
    tag.refTo = 2;
    const got = encodeNode(mm.str('r', tag));
    expect(toHex(got)).toBe('72 81 02 02 F9 03 E6');
    expect((decode(got) as NodeScalar).getTag().refTo).toBe(2);
  });

  test('reverseBytes roundtrip', () => {
    const bs = new Uint8Array([1, 2, 3, 4, 5]);
    expect(Array.from(reverseBytes(reverseBytes(bs)))).toEqual([1, 2, 3, 4, 5]);
    expect(Array.from(reverseBytes(bs))).toEqual([5, 4, 3, 2, 1]);
  });

  test('empty containers / is_null / unicode roundtrip via footer', () => {
    const srcs = [
      '{\n  // mm: desc=arr; allow_empty\n  "a": [],\n}',
      '{\n  // mm: desc=数据; pattern=^\\d$\n  "v": "1",\n}',
      '{\n  // mm: child_type=i8\n  "a": [1, 2],\n}',
      '{\n  // mm: enums=a|b|c\n  "v": "b",\n}',
    ];
    for (const src of srcs) {
      const wire = fromJSONC(src);
      // footer 编解码往返稳定：decode→JSONC→re-encode 应逐字节还原。
      const jsonc = decodeToJSONC(wire);
      const wire2 = fromJSONC(jsonc);
      expect(toHex(wire2)).toBe(toHex(wire));
    }
  });

  test('decode defense: tampered footer fails', () => {
    const tag = new Tag();
    tag.type = ValueType.Str;
    tag.desc = 'ab';
    const golden = encodeNode(mm.str('hi', tag)); // 68 69 82 62 61 1A 03 E7

    // truncated tail (drops the trailing tag sign) -> boundary mismatch
    expect(() => decode(golden.slice(0, golden.length - 1))).toThrow();

    // meta_len slot (2nd-to-last) claims empty meta
    const a = Uint8Array.from(golden);
    a[a.length - 2] = 0x00;
    expect(() => decode(a)).toThrow();

    // meta_len slot overstates meta length -> inner underflow
    const b = Uint8Array.from(golden);
    b[b.length - 2] = 0x05;
    expect(() => decode(b)).toThrow();

    // tag sign tampered to a bogus tier -> total/frames no longer line up
    const c = Uint8Array.from(golden);
    c[c.length - 1] = 0xeb; // total 11 > available
    expect(() => decode(c)).toThrow();
  });
});
