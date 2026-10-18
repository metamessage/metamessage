"""
Encoder for metamessage binary format.
Based on Go implementation in internal/core/encode*.go
"""

import io
import math
import re
import struct
from datetime import datetime, date, time as dt_time, timezone
from typing import Any, Optional, Union

from ..ir.tag import Tag, ValueType
from ..ir.ast import NodeObject, NodeArray, NodeScalar, Field, Node, NodeType, NodeNull
from .footer import footer_slot_bytes, reverse_bytes

# ===== Constants (matching Go internal/core/constants.go) =====

# Prefix types (3 bits)
Simple       = 0b000 << 5
PositiveInt  = 0b001 << 5
NegativeInt  = 0b010 << 5
PrefixFloat  = 0b011 << 5
PrefixString = 0b100 << 5
PrefixBytes  = 0b101 << 5
Container    = 0b110 << 5
PrefixTag    = 0b111 << 5

# Simple values
SimpleNull       = 0
SimpleNullBool   = 1
SimpleNullInt    = 2
SimpleNullFloat  = 3
SimpleNullString = 4
SimpleNullBytes  = 5
SimpleFalse      = 6
SimpleTrue       = 7

SimpleCode       = 8
SimpleMessage    = 9
SimpleData       = 10
SimpleSuccess    = 11
SimpleError      = 12
SimpleUnknown    = 13
SimplePage       = 14
SimpleLimit      = 15
SimpleOffset     = 16
SimpleTotal      = 17
SimpleId         = 18
SimpleName       = 19
SimpleDescription = 20
SimpleType       = 21
SimpleVersion    = 22
SimpleStatus     = 23
SimpleUrl        = 24
SimpleCreateTime = 25
SimpleUpdateTime = 26
SimpleDeleteTime = 27
SimpleAccount    = 28
SimpleToken      = 29
SimpleExpireTime = 30
SimpleKey        = 31

# Integer length constants
Max1Byte = 0xFF
Max2Byte = 0xFFFF
Max3Byte = 0xFFFFFF
Max4Byte = 0xFFFFFFFF
Max5Byte = 0xFFFFFFFFFF
Max6Byte = 0xFFFFFFFFFFFF
Max7Byte = 0xFFFFFFFFFFFFFF
Max8Byte = 0xFFFFFFFFFFFFFFFF

IntLenMask  = 0b11111
IntLen1Byte = IntLenMask - 7  # 24
IntLen2Byte = IntLenMask - 6  # 25
IntLen3Byte = IntLenMask - 5  # 26
IntLen4Byte = IntLenMask - 4  # 27
IntLen5Byte = IntLenMask - 3  # 28
IntLen6Byte = IntLenMask - 2  # 29
IntLen7Byte = IntLenMask - 1  # 30
IntLen8Byte = IntLenMask       # 31

# Float constants
FloatPositiveNegativeMask = 0b10000
FloatLenMask  = 0b01111
FloatLen1Byte = FloatLenMask - 7
FloatLen2Byte = FloatLenMask - 6
FloatLen3Byte = FloatLenMask - 5
FloatLen4Byte = FloatLenMask - 4
FloatLen5Byte = FloatLenMask - 3
FloatLen6Byte = FloatLenMask - 2
FloatLen7Byte = FloatLenMask - 1
FloatLen8Byte = FloatLenMask

# String constants
StringLenMask  = 0b11111
StringLen1Byte = StringLenMask - 1  # 30 (for length >= 24 and < 256)
StringLen2Byte = StringLenMask      # 31 (for length >= 256 and < 65536)

# Bytes constants（v0.3.0：4 档；内联 0..27，档位标记 28/29/30/31）
BytesLenMask  = 0b11111
BytesLen1Byte = BytesLenMask - 3  # 28
BytesLen2Byte = BytesLenMask - 2  # 29
BytesLen3Byte = BytesLenMask - 1  # 30
BytesLen4Byte = BytesLenMask      # 31

# Container constants（v0.3.0：4 档；内联 0..11，档位标记 12/13/14/15）
ContainerMask    = 0b10000
ContainerObject  = 0b00000
ContainerArray   = 0b10000
ContainerLenMask  = 0b01111
ContainerLen1Byte = ContainerLenMask - 3  # 12
ContainerLen2Byte = ContainerLenMask - 2  # 13
ContainerLen3Byte = ContainerLenMask - 1  # 14
ContainerLen4Byte = ContainerLenMask      # 15

# Tag constants（v0.3.0：4 档；tag+payload 合计槽；内联 0..27）
TagLenMask  = 0b11111
TagLen1Byte = TagLenMask - 3  # 28
TagLen2Byte = TagLenMask - 2  # 29
TagLen3Byte = TagLenMask - 1  # 30
TagLen4Byte = TagLenMask      # 31

DEFAULT_BUF_SIZE = 1024
MAX_CAP = 0xFFFFFFFF  # v0.3.0：~4GB，对齐 Max4Byte

_RE_SCI_FLOAT = re.compile(r'^([\d.]+)[eE]([+-]?\d+)$')


def get_prefix(b: int) -> int:
    return b & 0b11100000


def _to_bits(v: int, n: int):
    b = [0] * n
    for i in range(n):
        b[n - 1 - i] = (v >> i) & 1
    return b


def _write_bits(write_byte_fn, bits) -> int:
    bt = 0
    bl = 0
    n = 0
    for b in bits:
        bt = (bt << 1) | b
        bl += 1
        if bl == 8:
            write_byte_fn(bt)
            n += 1
            bt = 0
            bl = 0
    if bl > 0:
        bt <<= (8 - bl)
        write_byte_fn(bt)
        n += 1
    return n


def _encode_big_int_to_buf(buf: bytearray, s: str):
    """Encode big int string into bits and write to buf. Shared by module function and method."""
    if not s:
        return
    neg = s.startswith('-')
    if neg:
        s = s[1:]

    bits = [1 if neg else 0]
    n = len(s)
    i = 0
    while i < n:
        rem = n - i
        if rem >= 3:
            num = int(s[i:i + 3])
            bits.extend(_to_bits(num, 10))
            i += 3
        elif rem == 2:
            num = int(s[i:i + 2])
            bits.extend(_to_bits(num, 7))
            i += 2
        else:
            num = int(s[i:i + 1])
            bits.extend(_to_bits(num, 4))
            i += 1

    _write_bits(lambda b: buf.append(b), bits)


def _encode_big_int(write_byte_fn, s: str) -> int:
    buf = bytearray()
    _encode_big_int_to_buf(buf, s)
    for b in buf:
        write_byte_fn(b)
    return len(buf)


class Encoder:
    """
    Binary encoder for metamessage format.
    """

    def __init__(self, w: Optional[io.IOBase] = None):
        self.w = w
        self.buf = bytearray(DEFAULT_BUF_SIZE)
        self.offset = 0
        self.max_cap = MAX_CAP

    def reset(self, w: Optional[io.IOBase] = None):
        if w is not None:
            self.w = w
        self.offset = 0

    def _ensure_capacity(self, needed: int):
        required = self.offset + needed
        if required > self.max_cap:
            raise ValueError("maximum size exceeded")
        if required > len(self.buf):
            new_cap = len(self.buf) * 2
            if new_cap > self.max_cap or new_cap < required:
                new_cap = required
            new_buf = bytearray(new_cap)
            new_buf[:self.offset] = self.buf[:self.offset]
            self.buf = new_buf

    def _write_byte(self, b: int) -> int:
        self._ensure_capacity(1)
        self.buf[self.offset] = b & 0xFF
        self.offset += 1
        return 1

    def _write_bytes(self, bs: bytes) -> int:
        l = len(bs)
        if l == 0:
            return 0
        self._ensure_capacity(l)
        self.buf[self.offset:self.offset + l] = bs
        self.offset += l
        return l

    def _write_string(self, s: str) -> int:
        encoded = s.encode('utf-8')
        l = len(encoded)
        if l == 0:
            return 0
        self._ensure_capacity(l)
        self.buf[self.offset:self.offset + l] = encoded
        self.offset += l
        return l

    def _write_bytes_with_prefix(self, bs: bytes, *prefix: int) -> int:
        n = 0
        for p in prefix:
            n += self._write_byte(p)
        n += self._write_bytes(bs)
        return n

    # v0.5.0 尾定界：先写 body，再追加大端长度/尾 sign（一趟正向，无回写）。
    def _write_bytes_with_suffix(self, bs: bytes, *suffix: int) -> int:
        n = self._write_bytes(bs)
        for s in suffix:
            n += self._write_byte(s)
        return n

    def _write_string_with_prefix(self, s: str, *prefix: int) -> int:
        n = 0
        for p in prefix:
            n += self._write_byte(p)
        n += self._write_string(s)
        return n

    # ===== Simple values =====

    def _encode_simple(self, val: int) -> int:
        return self._write_byte(Simple | val)

    def _encode_bool(self, v: bool) -> int:
        val = SimpleTrue if v else SimpleFalse
        return self._encode_simple(val)

    # ===== Integers =====

    def _encode_i(self, sign: int, uv: int) -> int:
        if uv < IntLen1Byte:
            return self._write_byte(sign | uv)
        elif uv <= Max1Byte:
            return self._write_bytes_with_suffix(
                bytes([uv & 0xFF]),
                sign | IntLen1Byte
            )
        elif uv <= Max2Byte:
            return self._write_bytes_with_suffix(
                bytes([(uv >> 8) & 0xFF, uv & 0xFF]),
                sign | IntLen2Byte
            )
        elif uv <= Max3Byte:
            return self._write_bytes_with_suffix(
                bytes([(uv >> 16) & 0xFF, (uv >> 8) & 0xFF, uv & 0xFF]),
                sign | IntLen3Byte
            )
        elif uv <= Max4Byte:
            return self._write_bytes_with_suffix(
                struct.pack('>I', uv),
                sign | IntLen4Byte
            )
        elif uv <= Max5Byte:
            return self._write_bytes_with_suffix(
                struct.pack('>Q', uv)[3:],
                sign | IntLen5Byte
            )
        elif uv <= Max6Byte:
            return self._write_bytes_with_suffix(
                struct.pack('>Q', uv)[2:],
                sign | IntLen6Byte
            )
        elif uv <= Max7Byte:
            return self._write_bytes_with_suffix(
                struct.pack('>Q', uv)[1:],
                sign | IntLen7Byte
            )
        elif uv <= Max8Byte:
            return self._write_bytes_with_suffix(
                struct.pack('>Q', uv),
                sign | IntLen8Byte
            )
        else:
            raise ValueError("integer value too large")

    def _encode_u64(self, uv: int) -> int:
        return self._encode_i(PositiveInt, uv)

    def _encode_i64(self, v: int) -> int:
        if v >= 0:
            return self._encode_i(PositiveInt, v)
        else:
            if v == -9223372036854775808:
                uv = 9223372036854775808
            else:
                uv = -v
            return self._encode_i(NegativeInt, uv)

    def _encode_big_int(self, s: str) -> int:
        if not s:
            return 0
        neg = s.startswith('-')
        digits = s[1:] if neg else s
        n = len(digits)

        payload = bytearray()
        _encode_big_int_to_buf(payload, s)

        full_payload = bytearray(1 + len(payload))
        full_payload[0] = n
        full_payload[1:] = payload

        return self._encode_bytes(bytes(full_payload))

    # ===== Floats =====

    @staticmethod
    def _parse_float(s: str):
        """Parse float string into (is_negative, exponent, mantissa)."""
        is_negative = s.startswith('-')
        s = s.lstrip('-')

        if 'e' in s or 'E' in s:
            m = _RE_SCI_FLOAT.match(s)
            if not m:
                raise ValueError(f"invalid float scientific notation: {s}")
            num_part = m.group(1)
            exp_part = int(m.group(2))
            # Split num part
            if '.' in num_part:
                int_part, frac_part = num_part.split('.')
            else:
                int_part = num_part
                frac_part = ''
            base_exp = -len(frac_part)
            exp = exp_part + base_exp
            mantissa_str = (int_part + frac_part).lstrip('0') or '0'
            mantissa = int(mantissa_str)
            if exp < -128 or exp > 127:
                raise ValueError(f"exponent out of range: {exp}")
            return is_negative, exp, mantissa

        if '.' not in s:
            return is_negative, 0, int(s)

        parts = s.split('.')
        int_part = parts[0]
        frac_part = parts[1] if len(parts) > 1 else ''
        exp = -len(frac_part)
        mantissa_str = (int_part + frac_part).lstrip('0') or '0'
        mantissa = int(mantissa_str)
        return is_negative, exp, mantissa

    def _encode_float(self, s: str) -> int:
        is_negative, exponent, mantissa = self._parse_float(s)

        sign = PrefixFloat
        if is_negative:
            sign |= FloatPositiveNegativeMask

        # Special short encoding for common decimals: 0.1-0.7 (exponent=-1, mantissa<=7)
        if exponent == -1 and mantissa <= 7:
            return self._write_byte(sign | mantissa)

        # Encode exponent and mantissa
        exp_byte = exponent & 0xFF

        if mantissa <= Max1Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte, mantissa & 0xFF]),
                sign | FloatLen1Byte
            )
        elif mantissa <= Max2Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte, (mantissa >> 8) & 0xFF, mantissa & 0xFF]),
                sign | FloatLen2Byte
            )
        elif mantissa <= Max3Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte, (mantissa >> 16) & 0xFF, (mantissa >> 8) & 0xFF, mantissa & 0xFF]),
                sign | FloatLen3Byte
            )
        elif mantissa <= Max4Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte]) + struct.pack('>I', mantissa),
                sign | FloatLen4Byte
            )
        elif mantissa <= Max5Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[3:],
                sign | FloatLen5Byte
            )
        elif mantissa <= Max6Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[2:],
                sign | FloatLen6Byte
            )
        elif mantissa <= Max7Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[1:],
                sign | FloatLen7Byte
            )
        elif mantissa <= Max8Byte:
            return self._write_bytes_with_suffix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa),
                sign | FloatLen8Byte
            )
        else:
            raise ValueError(f"mantissa too large: {mantissa}")

    # ===== Strings =====

    # v0.5.0 尾定界：先写 UTF-8 body，再写长度大端，最后写 sign。
    def _encode_string(self, s: str) -> int:
        encoded = s.encode('utf-8')
        length = len(encoded)
        if length > Max2Byte:
            raise ValueError(f"string too long: {length}")

        sign = PrefixString
        if length < StringLen1Byte:  # length < 30
            n = self._write_bytes(encoded)
            n += self._write_byte(sign | length)
            return n
        elif length < Max1Byte:  # length < 255
            n = self._write_bytes(encoded)
            n += self._write_byte(length)
            n += self._write_byte(sign | StringLen1Byte)
            return n
        else:  # length <= 65535
            n = self._write_bytes(encoded)
            n += self._write_bytes(bytes([(length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_byte(sign | StringLen2Byte)
            return n

    # ===== Bytes =====

    # v0.5.0 尾定界：先写 body，再写长度大端，最后写 sign。
    def _encode_bytes(self, bs: bytes) -> int:
        length = len(bs)
        if length > Max4Byte:
            raise ValueError(f"bytes too long: {length}")

        sign = PrefixBytes
        if length < BytesLen1Byte:  # inline 0..27
            n = self._write_bytes(bs)
            n += self._write_byte(sign | length)
            return n
        elif length <= Max1Byte:  # 1 字节长度槽
            n = self._write_bytes(bs)
            n += self._write_byte(length)
            n += self._write_byte(sign | BytesLen1Byte)
            return n
        elif length <= Max2Byte:  # 2 字节长度槽
            n = self._write_bytes(bs)
            n += self._write_bytes(bytes([(length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_byte(sign | BytesLen2Byte)
            return n
        elif length <= Max3Byte:  # 3 字节长度槽
            n = self._write_bytes(bs)
            n += self._write_bytes(bytes([(length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_byte(sign | BytesLen3Byte)
            return n
        else:  # 4 字节长度槽
            n = self._write_bytes(bs)
            n += self._write_bytes(bytes([(length >> 24) & 0xFF, (length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_byte(sign | BytesLen4Byte)
            return n

    # ===== Containers =====

    # v0.5.0 尾定界：先写 payload（子节点 wire 正序），再写长度大端，最后写 sign。
    def _encode_array(self, bs: bytes) -> int:
        """Encode raw array bytes with container suffix (tail-delimited)."""
        length = len(bs)
        if length > Max4Byte:
            raise ValueError(f"array too long: {length}")

        sign = Container | ContainerArray
        if length < ContainerLen1Byte:  # inline 0..11
            return self._write_bytes_with_suffix(bs, sign | length)
        elif length <= Max1Byte:  # 1 字节
            return self._write_bytes_with_suffix(bs, length, sign | ContainerLen1Byte)
        elif length <= Max2Byte:  # 2 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen2Byte
            )
        elif length <= Max3Byte:  # 3 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen3Byte
            )
        else:  # 4 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 24) & 0xFF, (length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen4Byte
            )

    def _encode_object(self, bs: bytes) -> int:
        """Encode raw object bytes with container suffix (tail-delimited)."""
        length = len(bs)
        if length > Max4Byte:
            raise ValueError(f"object too long: {length}")

        sign = Container | ContainerObject
        if length < ContainerLen1Byte:  # inline 0..11
            return self._write_bytes_with_suffix(bs, sign | length)
        elif length <= Max1Byte:  # 1 字节
            return self._write_bytes_with_suffix(bs, length, sign | ContainerLen1Byte)
        elif length <= Max2Byte:  # 2 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen2Byte
            )
        elif length <= Max3Byte:  # 3 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen3Byte
            )
        else:  # 4 字节
            return self._write_bytes_with_suffix(
                bs,
                (length >> 24) & 0xFF, (length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF,
                sign | ContainerLen4Byte
            )

    # ===== Tags =====

    # encode_tag 尾定界写出 tag 包装段（一趟正向，零回写）：
    #   [inner 已编码 wire][reverse(meta 正序)][meta_len 尾部槽][total_len 档位大端][tag sign]
    # bs 为已完成尾定界编码的 inner 字节；meta 为 Tag.bytes() 正序字节。
    # total = len(bs)+len(meta)+len(尾槽)，与解码侧（footer.py node_tail_start）反推式
    # inner = total - ml - len(尾槽) 一致（尾槽与前缀槽为同一套 canonical 三档）。
    def _encode_tag(self, bs: bytes, meta: bytes) -> int:
        if len(meta) == 0:
            return 0

        slot = footer_slot_bytes(len(meta))
        length = len(bs) + len(meta) + len(slot)
        if length > Max4Byte:
            raise ValueError(f"tag+payload too long: {length}")

        sign = PrefixTag
        tier = b""
        if length < TagLen1Byte:
            sign |= length
        elif length <= Max1Byte:
            sign |= TagLen1Byte
            tier = bytes([length & 0xFF])
        elif length <= Max2Byte:
            sign |= TagLen2Byte
            tier = bytes([(length >> 8) & 0xFF, length & 0xFF])
        elif length <= Max3Byte:
            sign |= TagLen3Byte
            tier = bytes([(length >> 16) & 0xFF, (length >> 8) & 0xFF, length & 0xFF])
        else:
            # 最后一档必须覆盖到 Max4Byte
            sign |= TagLen4Byte
            tier = bytes([
                (length >> 24) & 0xFF, (length >> 16) & 0xFF,
                (length >> 8) & 0xFF, length & 0xFF,
            ])

        saved = self.offset
        self._write_bytes(bs)
        self._write_bytes(reverse_bytes(meta))
        self._write_bytes(slot)
        self._write_bytes(tier)
        self._write_byte(sign)
        return self.offset - saved

    # v0.5.0 尾定界：payload（已是 wire 形态）先出，meta 正序字节交给 encode_tag
    # 逆序后置 + 尾部槽，一趟写完，无前缀长度回写。
    def _encode_comment(self, payload: bytes, tag: Optional[Tag]) -> int:
        if tag is None:
            return 0
        meta = bytes(tag.bytes())
        if len(meta) == 0:
            return 0
        return self._encode_tag(payload, meta)

    # ===== Node encoding =====

    def _encode_node_object(self, obj: NodeObject) -> int:
        tag = obj.get_tag()
        buf_key = bytearray()
        buf = bytearray()

        for field in obj.fields:
            n = self._encode_any_node(field.value)
            encoded_sub = self.buf[self.offset - n:self.offset]
            buf.extend(encoded_sub)

            nk = self._encode_string(field.key)
            encoded_key = self.buf[self.offset - nk:self.offset]
            buf_key.extend(encoded_key)

        nk = self._encode_array(bytes(buf_key))
        encoded_key_array = self.buf[self.offset - nk:self.offset]

        buf_all = bytearray(len(encoded_key_array) + len(buf))
        buf_all[:len(encoded_key_array)] = encoded_key_array
        buf_all[len(encoded_key_array):] = buf

        n = self._encode_object(bytes(buf_all))

        n1 = self._encode_comment(bytes(self.buf[self.offset - n:self.offset]), tag)
        if n1 == 0:
            return n
        n = n1
        return n

    def _encode_node_array(self, arr: NodeArray) -> int:   
        tag = arr.get_tag()
        buf = bytearray()

        for item in arr.items:
            n = self._encode_any_node(item)
            encoded_sub = self.buf[self.offset - n:self.offset]
            buf.extend(encoded_sub)

        n = self._encode_array(bytes(buf))

        n1 = self._encode_comment(bytes(self.buf[self.offset - n:self.offset]), tag)
        if n1 == 0:
            return n
        n = n1
        return n

    def _encode_any_node(self, node: Node) -> int:
        if isinstance(node, NodeObject):
            return self._encode_node_object(node)
        elif isinstance(node, NodeArray):
            return self._encode_node_array(node)
        elif isinstance(node, NodeScalar):
            return self._encode_node_value(node)
        elif isinstance(node, NodeNull):
            return self._encode_node_null(node)
        else:
            raise ValueError(f"unsupported node type: {type(node)}")

    def _encode_node_null(self, val: NodeNull) -> int:
        n = self._encode_simple(SimpleNull)
        tag = val.get_tag()
        n1 = self._encode_comment(bytes(self.buf[self.offset - n:self.offset]), tag)
        if n1 == 0:
            return n
        return n1

    def _encode_node_value(self, val: NodeScalar) -> int:
        tag = val.get_tag()
        if tag is None:
            tag = Tag()

        start = self.offset

        if tag.type == ValueType.Datetime:
            if not tag.is_null:
                self._encode_datetime(val.data, tag.location)
        elif tag.type == ValueType.Date:
            if not tag.is_null:
                self._encode_date(val.data)
        elif tag.type == ValueType.Time:
            if not tag.is_null:
                self._encode_time(val.data)
        elif tag.type == ValueType.I:
            if tag.is_null:
                self._encode_simple(SimpleNullInt)
            else:
                try:
                    self._encode_i64(int(val.data))
                except ValueError:
                    self._encode_big_int(val.text)
        elif tag.type in (ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64):
            if not tag.is_null:
                self._encode_i64(int(val.data))
        elif tag.type in (ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64):
            if not tag.is_null:
                self._encode_u64(int(val.data))
        elif tag.type == ValueType.F32:
            if not tag.is_null:
                self._encode_float(val.text)
        elif tag.type == ValueType.F64:
            if tag.is_null:
                self._encode_simple(SimpleNullFloat)
            else:
                self._encode_float(val.text)
        elif tag.type == ValueType.Str:
            if tag.is_null:
                self._encode_simple(SimpleNullString)
            else:
                self._encode_string(val.text)
        elif tag.type == ValueType.Email:
            if not tag.is_null:
                self._encode_string(val.text)
        elif tag.type == ValueType.Uuid:
            if not tag.is_null:
                data = val.data
                if isinstance(data, (bytes, bytearray)):
                    self._encode_bytes(bytes(data[:16]))
                elif isinstance(data, str):
                    import binascii
                    raw = data.replace('-', '')
                    data_bytes = binascii.unhexlify(raw)
                    self._encode_bytes(data_bytes)
                else:
                    self._encode_bytes(bytes(data)[:16])
        elif tag.type == ValueType.Decimal:
            if not tag.is_null:
                self._encode_float(val.text)
        elif tag.type == ValueType.Url:
            if not tag.is_null:
                self._encode_string(val.text)
        elif tag.type == ValueType.Ip:
            if not tag.is_null:
                ip = val.data
                ver = tag.version
                if ver == 0:
                    self._encode_string(val.text)
                elif ver == 4:
                    if isinstance(ip, (bytes, bytearray)):
                        self._encode_bytes(bytes(ip[:4]))
                    else:
                        self._encode_string(val.text)
                elif ver == 6:
                    if isinstance(ip, (bytes, bytearray)) and len(ip) >= 16:
                        self._encode_bytes(bytes(ip[:16]))
                    else:
                        self._encode_string(val.text)
                else:
                    raise ValueError(f"unsupported IP version: {ver}")
        elif tag.type == ValueType.Bytes:
            if tag.is_null:
                self._encode_simple(SimpleNullBytes)
            else:
                data = val.data
                if isinstance(data, str):
                    import base64
                    data = base64.b64decode(data)
                elif not isinstance(data, (bytes, bytearray)):
                    data = bytes(data)
                self._encode_bytes(bytes(data))
        elif tag.type == ValueType.Media:
            if not tag.is_null:
                data = val.data
                if isinstance(data, str):
                    import base64
                    data = base64.b64decode(data)
                self._encode_bytes(bytes(data))
        elif tag.type == ValueType.Bigint:
            if not tag.is_null:
                self._encode_big_int(val.text)
        elif tag.type == ValueType.Bool:
            if tag.is_null:
                self._encode_simple(SimpleNullBool)
            else:
                self._encode_bool(bool(val.data))
        elif tag.type == ValueType.Enums:
            if not tag.is_null:
                data = val.data
                if isinstance(data, str):
                    if tag.enums:
                        enum_list = [e.strip() for e in tag.enums.split('|')]
                        try:
                            idx = enum_list.index(data)
                            data = idx
                        except ValueError:
                            raise ValueError(f"enum value '{data}' not found in enums: {tag.enums}")
                    else:
                        raise ValueError(f"enum value '{data}' is string but no enums defined in tag")
                self._encode_i64(int(data))
        else:
            raise ValueError(f"unsupported value type: {tag.type}")

        n = self.offset - start
        if n == 0 and not tag.is_null:
            raise ValueError(f"failed to encode value type: {tag.type}")

        n1 = self._encode_comment(bytes(self.buf[start:self.offset]), tag)
        if n1 == 0:
            return n
        n = n1

        return n

    # ===== Date/time encoding =====

    def _encode_datetime(self, t, location=0):
        if location is None:
            location = 0
        if isinstance(t, (datetime, dt_time)):
            if t.tzinfo is None:
                if location != 0:
                    from datetime import timedelta as _td2
                    t = t.replace(tzinfo=timezone(_td2(hours=location)))
                else:
                    t = t.replace(tzinfo=timezone.utc)
            v = int(t.timestamp())
        elif isinstance(t, str):
            for fmt in ['%Y-%m-%d %H:%M:%S', '%Y-%m-%dT%H:%M:%S', '%Y-%m-%d']:
                try:
                    dt = datetime.strptime(t, fmt)
                    if location != 0:
                        from datetime import timedelta as _td3
                        dt = dt.replace(tzinfo=timezone(_td3(hours=location)))
                    else:
                        dt = dt.replace(tzinfo=timezone.utc)
                    v = int(dt.timestamp())
                    break
                except ValueError:
                    continue
            else:
                raise ValueError(f"cannot parse datetime string: {t}")
        else:
            v = int(t)
        self._encode_i64(v)

    def _encode_date(self, t):
        if isinstance(t, (datetime, date)):
            d = t if isinstance(t, date) else t.date()
            ref = date(1970, 1, 1)
            days = (d - ref).days
        elif isinstance(t, str):
            d = datetime.strptime(t, '%Y-%m-%d').date()
            ref = date(1970, 1, 1)
            days = (d - ref).days
        else:
            days = int(t)
        sign = PositiveInt if days >= 0 else NegativeInt
        uv = abs(days)
        self._encode_i(sign, uv)

    def _encode_time(self, t):
        if isinstance(t, dt_time):
            v = t.hour * 3600 + t.minute * 60 + t.second
        elif isinstance(t, datetime):
            utc = t.utctimetuple()
            v = utc.tm_hour * 3600 + utc.tm_min * 60 + utc.tm_sec
        elif isinstance(t, str):
            parts = t.split(':')
            v = int(parts[0]) * 3600 + int(parts[1]) * 60 + int(parts[2])
        else:
            v = int(t)
        self._encode_u64(v)

    # ===== Public API =====

    def get_encoded_bytes(self, written: int) -> bytes:
        return bytes(self.buf[self.offset - written:self.offset])

    def encode(self, node: Node) -> bytes:
        self.offset = 0
        n = self._encode_any_node(node)
        result = bytes(self.buf[self.offset - n:self.offset])
        self.offset = 0
        return result

    def encode_stream(self, data: Any) -> int:
        """Encode Python data to binary and write to writer."""
        if self.w is None:
            raise ValueError("writer cannot be None")

        from .value_to_node import value_to_node
        node = value_to_node(data)
        encoded = self.encode(node)
        n = self.w.write(encoded)
        if n != len(encoded):
            raise ValueError(f"short write: wrote {n} bytes, expected {len(encoded)}")
        return n

    def _struct_to_mm(self, data: Any, path: str = "") -> Node:
        """Convert Python data to Node tree."""
        if isinstance(data, dict):
            fields = []
            for k, v in data.items():
                field_path = f"{path}.{k}" if path else k
                fields.append(Field(key=k, value=self._struct_to_mm(v, field_path)))
            return NodeObject(fields=fields, tag=Tag(), path=path)
        elif isinstance(data, list):
            items = [self._struct_to_mm(item, f"{path}.{i}") for i, item in enumerate(data)]
            return NodeArray(items=items, tag=Tag(type=ValueType.Arr), path=path)
        elif isinstance(data, bool):
            return NodeScalar(data=data, text="true" if data else "false", tag=Tag(type=ValueType.Bool), path=path)
        elif isinstance(data, int):
            return NodeScalar(data=data, text=str(data), tag=Tag(type=ValueType.I), path=path)
        elif isinstance(data, float):
            return NodeScalar(data=data, text=str(data), tag=Tag(type=ValueType.F64), path=path)
        elif isinstance(data, str):
            return NodeScalar(data=data, text=data, tag=Tag(type=ValueType.Str), path=path)
        elif isinstance(data, bytes):
            return NodeScalar(data=data, text="", tag=Tag(type=ValueType.Bytes), path=path)
        elif data is None:
            return NodeScalar(data=None, text="null", tag=Tag(is_null=True), path=path)
        else:
            return NodeScalar(data=data, text=str(data), tag=Tag(type=ValueType.Str), path=path)
