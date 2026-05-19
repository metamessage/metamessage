"""
Decoder for metamessage binary format.
Based on Go implementation in internal/core/decode*.go
"""

import io
import struct
from datetime import datetime, timezone, date, time as dt_time
from typing import Any, Optional, Tuple

from ..ir.tag import Tag, TagKey, ValueType
from ..ir.types import Obj, Arr, Val, Field, Node, NodeType

Simple       = 0b000 << 5
PositiveInt  = 0b001 << 5
NegativeInt  = 0b010 << 5
PrefixFloat  = 0b011 << 5
PrefixString = 0b100 << 5
PrefixBytes  = 0b101 << 5
Container    = 0b110 << 5
PrefixTag    = 0b111 << 5

SimpleNullBool   = 0
SimpleNullInt    = 1
SimpleNullFloat  = 2
SimpleNullString = 3
SimpleNullBytes  = 4
SimpleFalse      = 5
SimpleTrue       = 6
SimpleCode       = 7
SimpleMessage    = 8
SimpleData       = 9
SimpleSuccess    = 10
SimpleError      = 11
SimpleUnknown    = 12
SimplePage       = 13
SimpleLimit      = 14
SimpleOffset     = 15
SimpleTotal      = 16
SimpleId         = 17
SimpleName       = 18
SimpleDescription = 19
SimpleType       = 20
SimpleVersion    = 21
SimpleStatus     = 22
SimpleUrl        = 23
SimpleCreateTime = 24
SimpleUpdateTime = 25
SimpleDeleteTime = 26
SimpleAccount    = 27
SimpleToken      = 28
SimpleExpireTime = 29
SimpleKey        = 30
SimpleVal        = 31

Max1Byte = 0xFF
Max2Byte = 0xFFFF

IntLenMask  = 0b11111
IntLen1Byte = IntLenMask - 7
IntLen2Byte = IntLenMask - 6
IntLen3Byte = IntLenMask - 5
IntLen4Byte = IntLenMask - 4
IntLen5Byte = IntLenMask - 3
IntLen6Byte = IntLenMask - 2
IntLen7Byte = IntLenMask - 1
IntLen8Byte = IntLenMask

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

StringLenMask  = 0b11111
StringLen1Byte = StringLenMask - 1
StringLen2Byte = StringLenMask

BytesLenMask  = 0b11111
BytesLen1Byte = BytesLenMask - 1
BytesLen2Byte = BytesLenMask

ContainerMask   = 0b10000
ContainerObject = 0b00000
ContainerArray  = 0b10000
ContainerLenMask  = 0b01111
ContainerLen1Byte = ContainerLenMask - 1
ContainerLen2Byte = ContainerLenMask

TagLenMask  = 0b11111
TagLen1Byte = TagLenMask - 1
TagLen2Byte = TagLenMask

PrefixMask = 0b11100000
SuffixMask = 0b00011111

DEFAULT_TIME = datetime(1970, 1, 1, tzinfo=timezone.utc)


def get_prefix(b: int) -> int:
    return b & PrefixMask


def _int_len(b: int):
    l = b & IntLenMask
    if l < IntLen1Byte:
        return 0, l
    else:
        return l - IntLen1Byte + 1, 0


def _float_len(b: int):
    l = b & FloatLenMask
    if l < FloatLen1Byte:
        return 0, l
    else:
        return l - FloatLen1Byte + 1, 0


def _string_len(b: int):
    l = b & StringLenMask
    if l < StringLen1Byte:
        return 0, l
    elif l == StringLen1Byte:
        return 1, 0
    else:
        return 2, 0


def _bytes_len(b: int):
    l = b & BytesLenMask
    if l < BytesLen1Byte:
        return 0, l
    elif l == BytesLen1Byte:
        return 1, 0
    else:
        return 2, 0


def _container_len(b: int):
    l = b & ContainerLenMask
    if l < ContainerLen1Byte:
        return 0, l
    elif l == ContainerLen1Byte:
        return 1, 0
    else:
        return 2, 0


def _tag_len(b: int):
    l = b & TagLenMask
    if l < TagLen1Byte:
        return 0, l
    elif l == TagLen1Byte:
        return 1, 0
    else:
        return 2, 0


def _is_array(b: int) -> bool:
    return (b & ContainerMask) == ContainerArray


def _mantissa_to_decimal(mantissa: int, exp: int) -> str:
    num_str = str(mantissa)
    decimal_pos = len(num_str) + exp

    if decimal_pos <= 0:
        result = "0." + "0" * (-decimal_pos) + num_str
    elif 0 < decimal_pos < len(num_str):
        result = num_str[:decimal_pos] + "." + num_str[decimal_pos:]
    else:
        trailing_zeros = decimal_pos - len(num_str)
        result = num_str + "0" * trailing_zeros

    return result


def _decode_bytes_bigint(data: bytes, n: int) -> str:
    bits = _bytes_to_bits(data)
    if not bits:
        return ""

    neg = bits[0] == 1
    bits = bits[1:]

    parts = []
    while n > 0:
        if n >= 3 and len(bits) >= 10:
            num = _from_bits(bits[:10])
            parts.append(f"{num:03d}")
            bits = bits[10:]
            n -= 3
        elif n >= 2 and len(bits) >= 7:
            num = _from_bits(bits[:7])
            parts.append(f"{num:02d}")
            bits = bits[7:]
            n -= 2
        elif n >= 1 and len(bits) >= 4:
            num = _from_bits(bits[:4])
            parts.append(str(num))
            bits = bits[4:]
            n -= 1
        else:
            break

    res = "".join(parts)
    if neg:
        res = "-" + res
    return res


def _bytes_to_bits(data: bytes):
    bits = []
    for bt in data:
        for i in range(7, -1, -1):
            bits.append((bt >> i) & 1)
    return bits


def _from_bits(bits) -> int:
    v = 0
    for bit in bits:
        v = (v << 1) | bit
    return v


class Decoder:
    def __init__(self, data: bytes):
        self.data = data
        self.offset = 0

    def _read_byte(self) -> int:
        if self.offset >= len(self.data):
            raise ValueError("buffer underflow")
        b = self.data[self.offset]
        self.offset += 1
        return b

    def _read_bytes(self, n: int) -> bytes:
        if self.offset + n > len(self.data):
            raise ValueError("buffer underflow")
        bs = self.data[self.offset:self.offset + n]
        self.offset += n
        return bs

    def decode(self) -> Any:
        node = self._decode_top()
        return _node_to_python(node)

    def decode_node(self) -> Node:
        return self._decode_top()

    def _decode_top(self) -> Node:
        node, _ = self._decode(None, "")
        return node

    def _decode(self, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        b = self._read_byte()
        prefix = get_prefix(b)

        if prefix == PrefixTag:
            return self._decode_tag(b, path)
        elif prefix == Simple:
            return self._decode_simple(b, tag, path)
        elif prefix == PositiveInt:
            return self._decode_positive_int(b, tag, path)
        elif prefix == NegativeInt:
            return self._decode_negative_int(b, tag, path)
        elif prefix == PrefixFloat:
            return self._decode_float(b, tag, path)
        elif prefix == PrefixString:
            return self._decode_string(b, tag, path)
        elif prefix == PrefixBytes:
            return self._decode_bytes(b, tag, path)
        elif prefix == Container:
            return self._decode_container(b, tag, path)
        else:
            raise ValueError("invalid prefix")

    def _decode_tag(self, prefix: int, path: str) -> Tuple[Node, int]:
        l1, l2 = _tag_len(prefix)

        if l1 == 1:
            l2 = self._read_byte()
        elif l1 == 2:
            bs = self._read_bytes(2)
            l2 = (bs[0] << 8) | bs[1]

        tag = Tag()

        b = self._read_byte()
        l = b
        if l >= 255:
            bs = self._read_bytes(2)
            l = (bs[0] << 8) | bs[1]
        elif l >= 254:
            l = self._read_byte()

        while l > 0:
            n = self._decode_tag_bytes(tag)
            if n == 0:
                raise ValueError("tag error")
            if n > l:
                raise ValueError("tag overflow")
            l -= n

        if tag.is_null:
            node = self._decode_null_value(tag, path, prefix)
            if node is not None:
                length = l1 + 1 + l2
                return node, length

        node, _, = self._decode(tag, path)
        length = l1 + 1 + l2
        return node, length

    def _decode_null_value(self, tag: Tag, path: str, prefix: int) -> Optional[Node]:
        v = ValueType

        if tag.type == v.DateTime:
            text = DEFAULT_TIME.strftime('%Y-%m-%d %H:%M:%S')
            if tag.location is not None:
                text = DEFAULT_TIME.strftime('%Y-%m-%d %H:%M:%S')
            return Val(data=DEFAULT_TIME, text=text, tag=tag, path=path), 0
        elif tag.type == v.Date:
            text = DEFAULT_TIME.strftime('%Y-%m-%d')
            return Val(data=DEFAULT_TIME, text=text, tag=tag, path=path), 0
        elif tag.type == v.Time:
            text = DEFAULT_TIME.strftime('%H:%M:%S')
            return Val(data=DEFAULT_TIME, text=text, tag=tag, path=path), 0
        elif tag.type == v.Int8:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Int16:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Int32:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Int64:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Uint:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Uint8:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Uint16:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Uint32:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Uint64:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.Float32:
            return Val(data=0.0, text='0.0', tag=tag, path=path), 0
        elif tag.type in (v.Email, v.UUID, v.Decimal):
            return Val(data='', text='', tag=tag, path=path), 0
        elif tag.type == v.BigInt:
            return Val(data=0, text='0', tag=tag, path=path), 0
        elif tag.type == v.URL:
            return Val(data='', text='', tag=tag, path=path), 0
        elif tag.type == v.IP:
            if tag.version == 4:
                text = '0.0.0.0'
            elif tag.version == 6:
                text = '::'
            else:
                text = ''
            return Val(data='', text=text, tag=tag, path=path), 0

        return None

    def _decode_tag_bytes(self, tag: Tag) -> int:
        b = self._read_byte()
        p = b & 0xF8
        l = b & 0x07

        if p == TagKey.IsNull:
            tag.is_null = (l & 0x01) == 1
            if tag.is_null:
                tag.nullable = True
            return 1
        elif p == TagKey.Example:
            tag.example = (l & 0x01) == 1
            return 1
        elif p == TagKey.Desc:
            n, s = self._read_length_str(l, True, True)
            tag.desc = s
            return n
        elif p == TagKey.Type:
            tb = self._read_byte()
            tag.type = ValueType(tb)
            return 2
        elif p == TagKey.Raw:
            tag.raw = (l & 0x01) == 1
            return 1
        elif p == TagKey.Nullable:
            tag.nullable = (l & 0x01) == 1
            return 1
        elif p == TagKey.Default:
            n, s = self._read_length_str_small(l)
            tag.default = s
            return n
        elif p == TagKey.Min:
            n, s = self._read_length_str_small(l)
            tag.min = s
            return n
        elif p == TagKey.Max:
            n, s = self._read_length_str_small(l)
            tag.max = s
            return n
        elif p == TagKey.Size:
            tag.size = self._read_varint(l)
            return 2 + l
        elif p == TagKey.Enum:
            tag.type = ValueType.Enum
            n, s = self._read_length_str(l, True, True)
            tag.enum = s
            return n
        elif p == TagKey.Pattern:
            n, s = self._read_length_str_small(l)
            tag.pattern = s
            return n
        elif p == TagKey.Location:
            bs = self._read_bytes(l)
            tag.location = int(bs.decode('utf-8'))
            return 1 + l
        elif p == TagKey.Version:
            tag.version = self._read_varint(l)
            return 2 + l
        elif p == TagKey.Mime:
            if l < 7:
                tag.mime = str(l)
            else:
                tb = self._read_byte()
                tag.mime = str(tb)
                return 2
            return 1
        elif p == TagKey.ChildDesc:
            n, s = self._read_length_str(l, True, True)
            tag.child_desc = s
            return n
        elif p == TagKey.ChildType:
            tb = self._read_byte()
            tag.child_type = ValueType(tb)
            return 2
        elif p == TagKey.ChildRaw:
            tag.child_raw = (l & 0x01) == 1
            return 1
        elif p == TagKey.ChildNullable:
            tag.child_nullable = (l & 0x01) == 1
            return 1
        elif p == TagKey.ChildDefault:
            n, s = self._read_length_str_small(l)
            tag.child_default = s
            return n
        elif p == TagKey.ChildMin:
            n, s = self._read_length_str_small(l)
            tag.child_min = s
            return n
        elif p == TagKey.ChildMax:
            n, s = self._read_length_str_small(l)
            tag.child_max = s
            return n
        elif p == TagKey.ChildSize:
            tag.child_size = self._read_varint(l)
            return 2 + l
        elif p == TagKey.ChildEnum:
            tag.child_type = ValueType.Enum
            n, s = self._read_length_str(l, True, True)
            tag.child_enum = s
            return n
        elif p == TagKey.ChildPattern:
            n, s = self._read_length_str_small(l)
            tag.child_pattern = s
            return n
        elif p == TagKey.ChildLocation:
            bs = self._read_bytes(l)
            tag.child_location = int(bs.decode('utf-8'))
            return 1 + l
        elif p == TagKey.ChildVersion:
            tag.child_version = self._read_varint(l)
            return 2 + l
        elif p == TagKey.ChildMime:
            if l < 7:
                tag.child_mime = str(l)
            else:
                tb = self._read_byte()
                tag.child_mime = str(tb)
                return 2
            return 1
        else:
            raise ValueError("invalid data")

    def _read_length_str(self, l: int, large: bool = False, huge: bool = False):
        if l <= 5:
            bs = self._read_bytes(l)
            return 1 + l, bs.decode('utf-8')
        elif not large:
            return 1, ""
        elif l == 6:
            l2 = self._read_byte()
            bs = self._read_bytes(l2)
            return 1 + 1 + l2, bs.decode('utf-8')
        elif l == 7 and huge:
            bs2 = self._read_bytes(2)
            l2 = (bs2[0] << 8) | bs2[1]
            bs = self._read_bytes(l2)
            return 1 + 2 + l2, bs.decode('utf-8')
        return 1, ""

    def _read_length_str_small(self, l: int):
        if l < 7:
            bs = self._read_bytes(l)
            return 1 + l, bs.decode('utf-8')
        else:
            l2 = self._read_byte()
            bs = self._read_bytes(l2)
            return 1 + 1 + l2, bs.decode('utf-8')

    def _read_varint(self, l: int) -> int:
        v = 0
        for _ in range(l + 1):
            b = self._read_byte()
            v = (v << 8) | b
        return v

    def _decode_simple(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        if tag is None:
            tag = Tag()

        suffix = b & SuffixMask

        if suffix == SimpleFalse:
            tag.type = ValueType.Bool
            return Val(data=False, text='false', tag=tag, path=path), 1
        elif suffix == SimpleTrue:
            tag.type = ValueType.Bool
            return Val(data=True, text='true', tag=tag, path=path), 1
        elif suffix == SimpleNullBool:
            tag.type = ValueType.Bool
            return Val(data=False, text='false', tag=tag, path=path), 1
        elif suffix == SimpleNullInt:
            tag.type = ValueType.Int
            return Val(data=0, text='0', tag=tag, path=path), 1
        elif suffix == SimpleNullFloat:
            tag.type = ValueType.Float64
            return Val(data=0.0, text='0.0', tag=tag, path=path), 1
        elif suffix == SimpleNullString:
            tag.type = ValueType.String
            return Val(data='', text='', tag=tag, path=path), 1
        elif suffix == SimpleNullBytes:
            tag.type = ValueType.Bytes
            return Val(data=b'', text='', tag=tag, path=path), 1
        elif suffix == SimpleCode:
            tag.type = ValueType.String
            return Val(data=None, text='code', tag=tag, path=path), 1
        elif suffix == SimpleMessage:
            tag.type = ValueType.String
            return Val(data=None, text='message', tag=tag, path=path), 1
        elif suffix == SimpleData:
            tag.type = ValueType.String
            return Val(data=None, text='data', tag=tag, path=path), 1
        elif suffix == SimpleSuccess:
            tag.type = ValueType.String
            return Val(data=None, text='success', tag=tag, path=path), 1
        elif suffix == SimpleError:
            tag.type = ValueType.String
            return Val(data=None, text='error', tag=tag, path=path), 1
        elif suffix == SimpleUnknown:
            tag.type = ValueType.String
            return Val(data=None, text='unknown', tag=tag, path=path), 1
        elif suffix == SimplePage:
            tag.type = ValueType.String
            return Val(data=None, text='page', tag=tag, path=path), 1
        elif suffix == SimpleLimit:
            tag.type = ValueType.String
            return Val(data=None, text='limit', tag=tag, path=path), 1
        elif suffix == SimpleOffset:
            tag.type = ValueType.String
            return Val(data=None, text='offset', tag=tag, path=path), 1
        elif suffix == SimpleTotal:
            tag.type = ValueType.String
            return Val(data=None, text='total', tag=tag, path=path), 1
        elif suffix == SimpleId:
            tag.type = ValueType.String
            return Val(data=None, text='id', tag=tag, path=path), 1
        elif suffix == SimpleName:
            tag.type = ValueType.String
            return Val(data=None, text='name', tag=tag, path=path), 1
        elif suffix == SimpleDescription:
            tag.type = ValueType.String
            return Val(data=None, text='description', tag=tag, path=path), 1
        elif suffix == SimpleType:
            tag.type = ValueType.String
            return Val(data=None, text='type', tag=tag, path=path), 1
        elif suffix == SimpleVersion:
            tag.type = ValueType.String
            return Val(data=None, text='version', tag=tag, path=path), 1
        elif suffix == SimpleStatus:
            tag.type = ValueType.String
            return Val(data=None, text='status', tag=tag, path=path), 1
        elif suffix == SimpleUrl:
            tag.type = ValueType.String
            return Val(data=None, text='url', tag=tag, path=path), 1
        elif suffix == SimpleCreateTime:
            tag.type = ValueType.String
            return Val(data=None, text='create_time', tag=tag, path=path), 1
        elif suffix == SimpleUpdateTime:
            tag.type = ValueType.String
            return Val(data=None, text='update_time', tag=tag, path=path), 1
        elif suffix == SimpleDeleteTime:
            tag.type = ValueType.String
            return Val(data=None, text='delete_time', tag=tag, path=path), 1
        elif suffix == SimpleAccount:
            tag.type = ValueType.String
            return Val(data=None, text='account', tag=tag, path=path), 1
        elif suffix == SimpleToken:
            tag.type = ValueType.String
            return Val(data=None, text='token', tag=tag, path=path), 1
        elif suffix == SimpleExpireTime:
            tag.type = ValueType.String
            return Val(data=None, text='expire_time', tag=tag, path=path), 1
        elif suffix == SimpleKey:
            tag.type = ValueType.String
            return Val(data=None, text='key', tag=tag, path=path), 1
        elif suffix == SimpleVal:
            tag.type = ValueType.String
            return Val(data=None, text='val', tag=tag, path=path), 1
        else:
            raise ValueError(f"unsupported simple value: {b}")

    def _decode_positive_int(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        l1, l2 = _int_len(b)

        if l1 == 0:
            v = l2
        else:
            bs = self._read_bytes(l1)
            v = 0
            for i in range(l1):
                v = (v << 8) | bs[i]

        if tag is None:
            tag = Tag()

        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Int

        data = None
        text = str(v)

        if tag.type in (ValueType.Int, ValueType.Int8, ValueType.Int16, ValueType.Int32, ValueType.Int64):
            data = v
        elif tag.type in (ValueType.Uint, ValueType.Uint8, ValueType.Uint16, ValueType.Uint32, ValueType.Uint64):
            data = v
        elif tag.type == ValueType.DateTime:
            if tag.is_null:
                data = None
                text = ""
            else:
                d = datetime.fromtimestamp(v, tz=timezone.utc)
                data = d
                text = d.strftime('%Y-%m-%d %H:%M:%S')
        elif tag.type == ValueType.Date:
            if tag.is_null:
                data = None
                text = ""
            else:
                d = DEFAULT_TIME + date.resolution * v
                data = d
                text = d.strftime('%Y-%m-%d')
        elif tag.type == ValueType.Time:
            if tag.is_null:
                data = None
                text = ""
            else:
                hour = v // 3600
                minute = (v % 3600) // 60
                second = v % 60
                d = datetime(1970, 1, 1, hour, minute, second, tzinfo=timezone.utc)
                data = d
                text = d.strftime('%H:%M:%S')
        elif tag.type == ValueType.Enum:
            if tag.is_null:
                data = -1
                text = ""
            else:
                if tag.enum:
                    enums = [e.strip() for e in tag.enum.split('|')]
                    if v >= len(enums):
                        raise ValueError("enum index out of range")
                    data = v
                    text = enums[v]
                else:
                    raise ValueError("only enum are supported")
        else:
            raise ValueError(f"unsupported value type: {tag.type}")

        node = Val(data=data, text=text, tag=tag, path=path)
        length = l1 + 1
        return node, length

    def _decode_negative_int(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        l1, l2 = _int_len(b)

        if l1 == 0:
            v = l2
        else:
            bs = self._read_bytes(l1)
            v = 0
            for i in range(l1):
                v = (v << 8) | bs[i]

        if tag is None:
            tag = Tag()

        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Int

        data = None
        text = "-" + str(v)

        if tag.type in (ValueType.Int, ValueType.Int8, ValueType.Int16, ValueType.Int32, ValueType.Int64):
            data = -v
        elif tag.type == ValueType.DateTime:
            if tag.is_null:
                data = None
                text = ""
            else:
                d = datetime.fromtimestamp(-v, tz=timezone.utc)
                data = d
                text = d.strftime('%Y-%m-%d %H:%M:%S')
        elif tag.type == ValueType.Date:
            if tag.is_null:
                data = None
                text = ""
            else:
                d = DEFAULT_TIME - date.resolution * v
                data = d
                text = d.strftime('%Y-%m-%d')
        else:
            raise ValueError(f"unsupported negative value type: {tag.type}")

        node = Val(data=data, text=text, tag=tag, path=path)
        length = l1 + 1
        return node, length

    def _decode_float(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        l1, l2 = _float_len(b)

        v = 0.0
        length = 0

        if PrefixFloat <= b <= PrefixFloat + 7:
            v = (b & 0xF) / 10.0
            length = 1
        elif (PrefixFloat | FloatPositiveNegativeMask) <= b <= (PrefixFloat | FloatPositiveNegativeMask) + 7:
            v = -(b & 0xF) / 10.0
            length = 1
        else:
            exp_byte = self._read_byte()
            exp = exp_byte if exp_byte < 128 else exp_byte - 256

            if l1 == 0:
                mantissa = l2
            else:
                bs = self._read_bytes(l1)
                mantissa = 0
                for i in range(l1):
                    mantissa = (mantissa << 8) | bs[i]

            v = float(_mantissa_to_decimal(mantissa, exp))
            length = l1 + 2
            if b & FloatPositiveNegativeMask:
                v = -v

        if tag is None:
            tag = Tag()
        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Float64

        if tag.type == ValueType.Float32:
            data = float(v)
        elif tag.type == ValueType.Float64:
            data = v
        elif tag.type == ValueType.Decimal:
            data = _mantissa_to_decimal(mantissa, exp) if 'mantissa' in dir() else str(v)
            if b & FloatPositiveNegativeMask:
                data = '-' + data
        else:
            raise ValueError(f"unsupported float value type: {tag.type}")

        text = str(v)
        node = Val(data=data, text=text, tag=tag, path=path)
        return node, length

    def _decode_string(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        l1, l2 = _string_len(b)

        if l1 == 1:
            l2 = self._read_byte()
        elif l1 == 2:
            bs = self._read_bytes(2)
            l2 = (bs[0] << 8) | bs[1]

        if l2 > 0:
            by = self._read_bytes(l2)
            text = by.decode('utf-8')
        else:
            text = ""

        if tag is None:
            tag = Tag()

        if tag.type == ValueType.Unknown:
            tag.type = ValueType.String

        if tag.type == ValueType.Email:
            data = text
        elif tag.type == ValueType.URL:
            data = text
        elif tag.type == ValueType.IP:
            data = text
            try:
                import ipaddress
                data = ipaddress.ip_address(text)
                text = str(data)
            except ValueError:
                pass
        elif tag.type == ValueType.String:
            data = text
        else:
            raise ValueError(f"unsupported string type: {tag.type}")

        node = Val(data=data, text=text, tag=tag, path=path)
        length = l1 + 1 + l2
        return node, length

    def _decode_bytes(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        l1, l2 = _bytes_len(b)

        if l1 == 1:
            l2 = self._read_byte()
        elif l1 == 2:
            bs = self._read_bytes(2)
            l2 = (bs[0] << 8) | bs[1]

        if l2 > 0:
            bs = self._read_bytes(l2)
        else:
            bs = b""

        if tag is None:
            tag = Tag()

        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Bytes

        if tag.type == ValueType.BigInt:
            text = _decode_bytes_bigint(bs[1:], bs[0])
            try:
                data = int(text)
            except ValueError:
                data = text
        elif tag.type == ValueType.Bytes:
            data = bs
            import base64
            text = base64.b64encode(bs).decode('utf-8')
        elif tag.type == ValueType.UUID:
            data = bytes(bs[:16])
            text = '-'.join([bs[:4].hex(), bs[4:6].hex(), bs[6:8].hex(), bs[8:10].hex(), bs[10:16].hex()])
        elif tag.type == ValueType.IP:
            data = bytes(bs)
            try:
                import ipaddress
                ip = ipaddress.ip_address(data)
                text = str(ip)
            except ValueError:
                text = ""
        else:
            raise ValueError(f"unsupported bytes type: {tag.type}")

        node = Val(data=data, text=text, tag=tag, path=path)
        length = l1 + 1 + l2
        return node, length

    def _decode_container(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        if _is_array(b):
            return self._decode_array(b, tag, path)
        return self._decode_object(b, tag, path)

    def _decode_array(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        if tag is None:
            tag = Tag()
            tag.type = ValueType.Slice
        if tag.type == ValueType.Unknown:
            if tag.size > 0:
                tag.type = ValueType.Array
            else:
                tag.type = ValueType.Slice

        l1, l2 = _container_len(b)

        if l1 == 1:
            l2 = self._read_byte()
        elif l1 == 2:
            bs = self._read_bytes(2)
            l2 = (bs[0] << 8) | bs[1]

        arr = Arr(tag=tag, path=path)
        index = 0
        while index < l2:
            child_tag = Tag()
            child_tag.inherit(tag)

            child_path = f"{path}[{index}]"
            node, ln = self._decode(child_tag, child_path)
            if ln <= 0:
                raise ValueError(f"{child_path}: decode error")

            arr.items.append(node)
            index += ln

        length = l1 + 1 + l2
        return arr, length

    def _decode_object(self, b: int, tag: Optional[Tag], path: str) -> Tuple[Node, int]:
        if tag is None:
            tag = Tag()
            tag.type = ValueType.Object
        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Object

        l1, l2 = _container_len(b)

        if l1 == 1:
            l2 = self._read_byte()
        elif l1 == 2:
            bs = self._read_bytes(2)
            l2 = (bs[0] << 8) | bs[1]

        obj = Obj(tag=tag, path=path)

        l_array = self._read_byte()
        keys_arr, l_keys = self._decode_array(l_array, tag, path)

        index = l_keys
        i = 0
        while index < l2:
            child_tag = Tag()
            child_tag.inherit(tag)

            key = keys_arr.items[i].text
            child_path = f"{path}.{key}"
            node, ln = self._decode(child_tag, child_path)
            if ln <= 0:
                raise ValueError(f"{child_path}: decode error")

            obj.fields.append(Field(key=key, value=node))
            index += ln
            i += 1

        length = l1 + 1 + l2
        return obj, length


def _node_to_python(node: Node) -> Any:
    if isinstance(node, Val):
        return node.data
    elif isinstance(node, Arr):
        return [_node_to_python(item) for item in node.items]
    elif isinstance(node, Obj):
        return {field.key: _node_to_python(field.value) for field in node.fields}
    else:
        return None