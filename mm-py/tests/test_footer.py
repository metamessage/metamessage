"""
v0.5.0 尾定界专项测试。
golden hex / 尾部槽三档 / 篡改防御，期望字节与 Go 参考实现完全一致（编码器直接尾序写出，无出口反转）。
尾序布局：[inner wire][reverse(meta)][meta_len 尾槽][total 档位大端][tag sign]。
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from metamessage.core.encoder import Encoder
from metamessage.core.decoder import Decoder
from metamessage.core.footer import (
    footer_slot_bytes, parse_footer_slot, reverse_bytes,
    FOOTER_SLOT_DIRECT_MAX,
)
from metamessage.ir.tag import Tag, ValueType
from metamessage.ir.ast import NodeScalar


def _hx(bs):
    return ' '.join(f'{b:02X}' for b in bs)


def _raises(fn):
    try:
        fn()
    except Exception:
        return True
    raise AssertionError('expected error, got none')


def test_slot_tiers_canonical_widths():
    cases = [(1, 1), (253, 1), (254, 2), (255, 2), (256, 3), (65535, 3)]
    for meta_len, width in cases:
        slot = footer_slot_bytes(meta_len)
        assert len(slot) == width, f"metaLen={meta_len} want width {width}, got {len(slot)}"
        data = bytes([0x00, 0x11, 0x22, 0x33]) + slot
        got_len, got_w = parse_footer_slot(data, len(data))
        assert got_len == meta_len, f"parse metaLen {got_len} != {meta_len}"
        assert got_w == width, f"parse slotW {got_w} != {width}"


def test_slot_too_long_rejected():
    _raises(lambda: footer_slot_bytes(65536))


def test_golden_single_field_desc():
    e = Encoder()
    t = Tag(type=ValueType.Str)
    t.desc = 'ab'
    got = e.encode(NodeScalar('hi', 'hi', t))
    # tail-order: inner body [68 69 82], reverse(meta) [62 61 1A], slot [03], sign E7
    assert _hx(got) == '68 69 82 62 61 1A 03 E7'


def test_golden_multi_field_order():
    e = Encoder()
    t = Tag(type=ValueType.Str)
    t.example = True
    t.nullable = True
    t.desc = 'ab'
    got = e.encode(NodeScalar('x', 'x', t))
    assert _hx(got) == '78 81 62 61 1A 29 09 05 E8'
    node = Decoder(got).decode_node()
    assert node.tag.example is True
    assert node.tag.nullable is True
    assert node.tag.desc == 'ab'


def test_golden_ref_meta_at_tail():
    e = Encoder()
    t = Tag(type=ValueType.Str)
    t.ref_to = 2
    got = e.encode(NodeScalar('r', 'r', t))
    assert _hx(got) == '72 81 02 02 F9 03 E6'
    assert Decoder(got).decode_node().tag.ref_to == 2


def test_reverse_bytes_roundtrip():
    bs = bytes([1, 2, 3, 4, 5])
    assert reverse_bytes(reverse_bytes(bs)) == bs
    assert reverse_bytes(bs) == bytes([5, 4, 3, 2, 1])


def test_unicode_and_inherit_roundtrip_stable():
    # footer 编解码往返稳定：encode→decode→re-encode 逐字节还原。
    srcs = [
        '{\n  // mm: desc=数据; pattern=^\\d$\n  "v": "1"\n}',
        '{\n  // mm: child_type=i8\n  "a": [1, 2]\n}',
        '{\n  // mm: enums=a|b|c\n  "v": "b"\n}',
    ]
    from metamessage.jsonc.jsonc import parse_jsonc
    from metamessage.jsonc import to_jsonc
    for src in srcs:
        e = Encoder()
        wire = e.encode(parse_jsonc(src))
        node = Decoder(wire).decode_node()
        wire2 = Encoder().encode(parse_jsonc(to_jsonc(node)))
        assert _hx(wire2) == _hx(wire), f"roundtrip drift for {src!r}"


def test_decode_defense_tampered_footer():
    e = Encoder()
    t = Tag(type=ValueType.Str)
    t.desc = 'ab'
    golden = e.encode(NodeScalar('hi', 'hi', t))  # 68 69 82 62 61 1A 03 E7

    # 截断尾部（去掉末尾 tag sign）
    _raises(lambda: Decoder(golden[:-1]).decode())

    # meta_len 尾槽（倒数第 2 字节）声称巨大 meta 长度 -> inner 下溢
    a = bytearray(golden)
    a[-2] = 0x7F
    _raises(lambda: Decoder(bytes(a)).decode())

    # meta_len 尾槽夸大 meta 长度
    b = bytearray(golden)
    b[-2] = 0x05
    _raises(lambda: Decoder(bytes(b)).decode())

    # tag sign 篡改为 bogus 档位 -> 定界不再对齐
    c = bytearray(golden)
    c[-1] = 0xEB  # total 11 > 可用长度
    _raises(lambda: Decoder(bytes(c)).decode())


if __name__ == '__main__':
    for name, fn in sorted(globals().items()):
        if name.startswith('test_') and callable(fn):
            fn()
            print(f"  {name} OK")
