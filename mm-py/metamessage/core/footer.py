# v0.5.0 尾定界辅助：
#   1. meta_len 尾部槽（footer_slot_bytes / parse_footer_slot）+ 字节逆序（reverse_bytes），
#      与 Go internal/core/footer.go 逐字节对齐。
#   2. 解码入口剥帧 untailify（node_tail_start / untailify_node / untailify_children），
#      与 Go internal/core/tailframe.go 逐字节对齐：从节点尾 sign 反向定界，把尾定界
#      wire 还原成「sign 在前」的头部形态，再交给原有正向解码器。
#
# 编码器各原语已改为直接按尾定界顺序写出（body→lenBE→sign），本文件不包含编码侧
# 换位逻辑，也不设独立 payload_len 槽（各类型 sign 低 5 位档位/内联语义承担定界）。
#
# meta_len 尾槽三档（canonical 最小宽度，末字节为标记或直接值）：
#   - metaLen <= 253：[metaLen]                     槽宽 1
#   - metaLen <= 255：[metaLen][254]                槽宽 2
#   - metaLen <= 65535：[hi][lo][255]（u16 大端）    槽宽 3

# ===== Prefix（3 位类型前缀，左移 5 位） =====
Simple = 0b000 << 5
PositiveInt = 0b001 << 5
NegativeInt = 0b010 << 5
PrefixFloat = 0b011 << 5
PrefixString = 0b100 << 5
PrefixBytes = 0b101 << 5
Container = 0b110 << 5
PrefixTag = 0b111 << 5

PrefixMask = 0b11100000

# ===== 档位掩码 =====
IntLenMask = 0b11111
IntLen1Byte = 24  # IntLenMask - 7

FloatLenMask = 0b01111
FloatLen1Byte = FloatLenMask - 7  # 8

StringLenMask = 0b11111
StringLen1Byte = StringLenMask - 1  # 30

BytesLenMask = 0b11111
BytesLen1Byte = BytesLenMask - 3  # 28

ContainerLenMask = 0b01111
ContainerLen1Byte = ContainerLenMask - 3  # 12

TagLenMask = 0b11111
TagLen1Byte = TagLenMask - 3  # 28

Max1Byte = 0xFF
Max2Byte = 0xFFFF

FOOTER_SLOT_DIRECT_MAX = 253  # 一档可直写的 meta 长度上限（尾字节 < 254）
FOOTER_SLOT_MARKER1 = 254     # 二档标记：倒数第 2 字节为 u8 长度
FOOTER_SLOT_MARKER2 = 255     # 三档标记：倒数第 3..2 字节为 u16 大端长度


def get_prefix(b: int) -> int:
    return b & PrefixMask


# ----- 各类型档位函数：返回 (字节宽 w, 内联值 inline)，w==0 表示内联档 -----

def _int_len(b: int):
    l = b & IntLenMask
    if l < IntLen1Byte:
        return 0, l
    return l - IntLen1Byte + 1, 0


def _float_len(b: int):
    l = b & FloatLenMask
    if l < FloatLen1Byte:
        return 0, l
    return l - FloatLen1Byte + 1, 0


def _string_len(b: int):
    l = b & StringLenMask
    if l < StringLen1Byte:
        return 0, l
    elif l == StringLen1Byte:
        return 1, 0
    return 2, 0


def _bytes_len(b: int):
    l = b & BytesLenMask
    if l < BytesLen1Byte:
        return 0, l
    return l - BytesLen1Byte + 1, 0


def _container_len(b: int):
    l = b & ContainerLenMask
    if l < ContainerLen1Byte:
        return 0, l
    return l - ContainerLen1Byte + 1, 0


def _tag_len(b: int):
    l = b & TagLenMask
    if l < TagLen1Byte:
        return 0, l
    return l - TagLen1Byte + 1, 0


# footer_slot_bytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
def footer_slot_bytes(meta_len: int) -> bytes:
    if meta_len <= FOOTER_SLOT_DIRECT_MAX:
        return bytes([meta_len])
    elif meta_len <= Max1Byte:
        return bytes([meta_len, FOOTER_SLOT_MARKER1])
    elif meta_len <= Max2Byte:
        return bytes([(meta_len >> 8) & 0xFF, meta_len & 0xFF, FOOTER_SLOT_MARKER2])
    raise ValueError(f"tag too long, max length: {Max2Byte}, actual: {meta_len}")


# parse_footer_slot 从节点尾端反向解析长度槽，返回 (meta_len, slot_width)。
# data 为完整 wire 缓冲，node_end 为节点（含封装）的结束偏移（不含）。
def parse_footer_slot(data: bytes, node_end: int):
    if node_end < 1:
        raise ValueError("footer slot out of range")
    marker = data[node_end - 1]
    if marker == FOOTER_SLOT_MARKER1:
        if node_end < 2:
            raise ValueError("footer slot truncated")
        return data[node_end - 2], 2
    elif marker == FOOTER_SLOT_MARKER2:
        if node_end < 3:
            raise ValueError("footer slot truncated")
        return ((data[node_end - 3] << 8) | data[node_end - 2]), 3
    return marker, 1


# reverse_bytes 返回 bs 的字节逆序拷贝。
def reverse_bytes(bs: bytes) -> bytes:
    return bytes(reversed(bs))


# meta_prefix_width 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，
# 宽度恒等于 footer_slot_bytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
def meta_prefix_width(ml: int) -> int:
    if ml < 254:
        return 1
    if ml <= 255:
        return 2
    return 3


# meta_prefix_bytes 还原正向 meta 前缀槽。
def meta_prefix_bytes(ml: int) -> bytes:
    if ml < 254:
        return bytes([ml])
    if ml <= Max1Byte:
        return bytes([254, ml])
    if ml <= Max2Byte:
        return bytes([255, (ml >> 8) & 0xFF, ml & 0xFF])
    raise ValueError(f"meta too long: {ml}")


# tier_of 依 sign 的 3 位类型前缀分派各档位函数，返回 (prefix, w, inline)。
def tier_of(sign: int):
    p = get_prefix(sign)
    if p == Simple:
        return p, 0, 0
    if p in (PositiveInt, NegativeInt):
        return (p,) + _int_len(sign)
    if p == PrefixFloat:
        return (p,) + _float_len(sign)
    if p == PrefixString:
        return (p,) + _string_len(sign)
    if p == PrefixBytes:
        return (p,) + _bytes_len(sign)
    if p == Container:
        return (p,) + _container_len(sign)
    if p == PrefixTag:
        return (p,) + _tag_len(sign)
    raise ValueError(f"invalid prefix {p}")


def _tier_bytes_be(length: int, w: int) -> bytes:
    if w == 0:
        return b""
    out = bytearray()
    for i in range(w - 1, -1, -1):
        out.append((length >> (8 * i)) & 0xFF)
    return bytes(out)


# read_tier_be 从 [end-w, end) 读大端档位长度。
def _read_tier_be(data: bytes, end: int, w: int) -> int:
    if end - w < 0:
        raise ValueError("truncated: tier bytes out of range")
    v = 0
    for i in range(end - w, end):
        v = v * 256 + data[i]
    return v


# node_tail_start 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
def node_tail_start(wire: bytes, end: int) -> int:
    if end < 1 or end > len(wire):
        raise ValueError(f"truncated: node end {end} out of range")
    sign = wire[end - 1]
    p, w, inline = tier_of(sign)
    if p == Simple:
        return end - 1
    if p in (PositiveInt, NegativeInt):
        return end - 1 if w == 0 else end - 1 - w
    if p == PrefixFloat:
        # float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
        return end - 1 if w == 0 else end - 2 - w
    if p in (PrefixString, PrefixBytes, Container):
        if w == 0:
            return end - 1 - inline
        return end - 1 - w - _read_tier_be(wire, end - 1, w)
    if p == PrefixTag:
        total = inline if w == 0 else _read_tier_be(wire, end - 1, w)
        slot_end = end - 1 - w
        ml, slot_w = parse_footer_slot(wire, slot_end)
        # 反推用实际尾槽宽 slot_w：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
        # 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
        inner = total - ml - slot_w
        if inner < 0:
            raise ValueError("tag total underflow")
        return slot_end - slot_w - ml - inner
    raise ValueError("invalid prefix")


# untailify_node：wire 形态 [start,end) → head 形态字节。
def untailify_node(wire: bytes, start: int, end: int) -> bytes:
    if start < 0 or end > len(wire) or start >= end:
        raise ValueError(f"untailify: bad range [{start},{end})")
    sign = wire[end - 1]
    p, w, inline = tier_of(sign)

    if p == Simple:
        return bytes([sign])
    if p in (PositiveInt, NegativeInt, PrefixFloat):
        if w == 0:
            return bytes([sign])
        body = wire[start:end - 1]
        return bytes([sign]) + body
    if p in (PrefixString, PrefixBytes, Container):
        l = inline if w == 0 else _read_tier_be(wire, end - 1, w)
        body = wire[start:end - 1 - w]
        if len(body) != l:
            raise ValueError(f"untailify: payload {len(body)}, tier says {l}")
        out = bytes([sign]) + _tier_bytes_be(l, w)
        if p == Container:
            return out + untailify_children(body)
        return out + body
    if p == PrefixTag:
        total = inline if w == 0 else _read_tier_be(wire, end - 1, w)
        slot_end = end - 1 - w
        ml, slot_w = parse_footer_slot(wire, slot_end)
        inner_len = total - ml - slot_w
        if inner_len < 0:
            raise ValueError("untailify tag: total underflow")
        meta_end = slot_end - slot_w
        meta = reverse_bytes(wire[meta_end - ml:meta_end])
        inner_wire = wire[start:meta_end - ml]
        if len(inner_wire) != inner_len:
            raise ValueError(f"untailify tag: inner {len(inner_wire)}, want {inner_len}")
        inner = b""
        if inner_len > 0:
            inner = untailify_node(inner_wire, 0, len(inner_wire))
        out = bytes([sign]) + _tier_bytes_be(total, w)
        out += meta_prefix_bytes(ml)
        out += meta
        return out + inner
    raise ValueError("untailify: unreachable")


# untailify_children：wire 容器 payload 内顺序 wire 子节点逐个还原。
# wire 子节点以尾 sign 结尾，从区域末端反向定界扫描，还原后按正序拼接。
def untailify_children(payload: bytes) -> bytes:
    nodes = []
    end = len(payload)
    while end > 0:
        start = node_tail_start(payload, end)
        if start < 0 or start >= end:
            raise ValueError(f"untailify children: bad boundary {start}..{end}")
        nodes.append(untailify_node(payload, start, end))
        end = start
    out = b""
    for i in range(len(nodes) - 1, -1, -1):
        out += nodes[i]
    return out


# untailify 解码入口：整帧 wire 形态 → head 形态。
def untailify(wire: bytes) -> bytes:
    if len(wire) == 0:
        raise ValueError("empty wire")
    start = node_tail_start(wire, len(wire))
    if start != 0:
        raise ValueError(f"footer boundary mismatch: node start {start}, want 0")
    return untailify_node(wire, start, len(wire))
