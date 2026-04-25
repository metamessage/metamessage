from typing import Any, Dict, List, Optional
import io
import struct
from dataclasses import dataclass, field

from .tag import Tag, ValueType
from .types import Obj, Arr, Val, Field


Simple = 0b000 << 5
PositiveInt = 0b001 << 5
NegativeInt = 0b010 << 5
PrefixFloat = 0b011 << 5
PrefixString = 0b100 << 5
PrefixBytes = 0b101 << 5
Container = 0b110 << 5
PrefixTag = 0b111 << 5

SimpleNullBool = 0
SimpleNullInt = 1
SimpleNullFloat = 2
SimpleNullString = 3
SimpleNullBytes = 4
SimpleFalse = 5
SimpleTrue = 6

Max1Byte = 0xFF
Max2Byte = 0xFFFF
Max3Byte = 0xFFFFFF
Max4Byte = 0xFFFFFFFF
Max5Byte = 0xFFFFFFFFFF
Max6Byte = 0xFFFFFFFFFFFF
Max7Byte = 0xFFFFFFFFFFFFFF
Max8Byte = 0xFFFFFFFFFFFFFFFF

IntLenMask = 0b11111
IntLen1Byte = IntLenMask - 7
IntLen2Byte = IntLenMask - 6
IntLen3Byte = IntLenMask - 5
IntLen4Byte = IntLenMask - 4
IntLen5Byte = IntLenMask - 3
IntLen6Byte = IntLenMask - 2
IntLen7Byte = IntLenMask - 1
IntLen8Byte = IntLenMask

FloatPositiveNegativeMask = 0b10000
FloatLenMask = 0b01111
FloatLen1Byte = FloatLenMask - 7
FloatLen2Byte = FloatLenMask - 6
FloatLen3Byte = FloatLenMask - 5
FloatLen4Byte = FloatLenMask - 4
FloatLen5Byte = FloatLenMask - 3
FloatLen6Byte = FloatLenMask - 2
FloatLen7Byte = FloatLenMask - 1
FloatLen8Byte = FloatLenMask

StringLenMask = 0b11111
StringLen1Byte = StringLenMask - 1
StringLen2Byte = StringLenMask

BytesLenMask = 0b11111
BytesLen1Byte = BytesLenMask - 1
BytesLen2Byte = BytesLenMask

ContainerMask = 0b10000
ContainerMap = 0b00000
ContainerArray = 0b10000
ContainerLenMask = 0b01111
ContainerLen1Byte = 24
ContainerLen2Byte = 25

TagLenMask = 0b11111
TagLen1Byte = TagLenMask - 1
TagLen2Byte = TagLenMask

MAX_2BYTE = 0xFFFF


class Encoder:
    def __init__(self, w: Optional[io.BytesIO] = None):
        self.buf = io.BytesIO()
        self.offset = 0
        if w is not None:
            self.reset(w)

    def reset(self, w):
        self.buf = io.BytesIO()
        self.offset = 0
        self.w = w

    def encode(self, node) -> bytes:
        if isinstance(node, Obj):
            self._encode_obj(node)
        elif isinstance(node, Arr):
            self._encode_arr(node)
        elif isinstance(node, Val):
            self._encode_val(node)
        else:
            raise ValueError(f"unsupported type {type(node)}")

        result = self.buf.getvalue()
        self.buf = io.BytesIO()
        self.offset = 0
        return result

    def encode_stream(self, data: Any, path: str = "") -> bytes:
        node = struct_to_mm(data, path)
        return self.encode(node)

    def _write_byte(self, b: int):
        self.buf.write(bytes([b]))
        self.offset += 1
        return 1

    def _write_bytes(self, bs: bytes):
        self.buf.write(bs)
        self.offset += len(bs)
        return len(bs)

    def _write_bytes_with_prefix(self, bs: bytes, *prefix: int):
        total = 0
        for p in prefix:
            total += self._write_byte(p)
        total += self._write_bytes(bs)
        return total

    def _encode_simple(self, val: int) -> int:
        return self._write_byte(Simple | val)

    def _encode_bool(self, v: bool) -> int:
        val = SimpleTrue if v else SimpleFalse
        return self._encode_simple(val)

    def _encode_int(self, sign: int, uv: int) -> int:
        if uv < IntLen1Byte:
            return self._write_byte(sign | uv)
        elif uv <= Max1Byte:
            n = self._write_byte(sign | IntLen1Byte)
            n += self._write_byte(uv & 0xFF)
            return n
        elif uv <= Max2Byte:
            n = self._write_byte(sign | IntLen2Byte)
            n += self._write_bytes(bytes([(uv >> 8) & 0xFF, uv & 0xFF]))
            return n
        elif uv <= Max3Byte:
            n = self._write_byte(sign | IntLen3Byte)
            n += self._write_bytes(bytes([(uv >> 16) & 0xFF, (uv >> 8) & 0xFF, uv & 0xFF]))
            return n
        elif uv <= Max4Byte:
            n = self._write_byte(sign | IntLen4Byte)
            n += self._write_bytes(struct.pack(">I", uv)[1:])
            return n
        elif uv <= Max5Byte:
            n = self._write_byte(sign | IntLen5Byte)
            n += self._write_bytes(struct.pack(">Q", uv)[3:])
            return n
        elif uv <= Max6Byte:
            n = self._write_byte(sign | IntLen6Byte)
            n += self._write_bytes(struct.pack(">Q", uv)[2:])
            return n
        elif uv <= Max7Byte:
            n = self._write_byte(sign | IntLen7Byte)
            n += self._write_bytes(struct.pack(">Q", uv)[1:])
            return n
        elif uv <= Max8Byte:
            n = self._write_byte(sign | IntLen8Byte)
            n += self._write_bytes(struct.pack(">Q", uv))
            return n
        else:
            raise ValueError("invalid byte size")
        return 0

    def _encode_uint64(self, uv: int) -> int:
        return self._encode_int(PositiveInt, uv)

    def _encode_int64(self, v: int) -> int:
        sign = PositiveInt if v >= 0 else NegativeInt
        uv = abs(v)
        if v < 0 and v == -9223372036854775808:
            uv = 9223372036854775808
        return self._encode_int(sign, uv)

    def _encode_float(self, s: str) -> int:
        is_negative, exponent, mantissa = parse_float(s)
        
        sign = PrefixFloat
        if is_negative:
            sign |= FloatPositiveNegativeMask

        if exponent == -1 and mantissa <= 7:
            return self._write_byte(sign | mantissa)
        elif mantissa <= Max1Byte:
            self._write_byte(sign | FloatLen1Byte)
            self._write_byte(exponent & 0xFF)
            self._write_byte(mantissa & 0xFF)
            return 3
        elif mantissa <= Max2Byte:
            self._write_byte(sign | FloatLen2Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(bytes([(mantissa >> 8) & 0xFF, mantissa & 0xFF]))
            return 4
        elif mantissa <= Max3Byte:
            self._write_byte(sign | FloatLen3Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(bytes([(mantissa >> 16) & 0xFF, (mantissa >> 8) & 0xFF, mantissa & 0xFF]))
            return 5
        elif mantissa <= Max4Byte:
            self._write_byte(sign | FloatLen4Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(struct.pack(">I", mantissa)[1:])
            return 5
        elif mantissa <= Max5Byte:
            self._write_byte(sign | FloatLen5Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(struct.pack(">Q", mantissa)[3:])
            return 6
        elif mantissa <= Max6Byte:
            self._write_byte(sign | FloatLen6Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(struct.pack(">Q", mantissa)[2:])
            return 7
        elif mantissa <= Max7Byte:
            self._write_byte(sign | FloatLen7Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(struct.pack(">Q", mantissa)[1:])
            return 8
        elif mantissa <= Max8Byte:
            self._write_byte(sign | FloatLen8Byte)
            self._write_byte(exponent & 0xFF)
            self._write_bytes(struct.pack(">Q", mantissa))
            return 9
        else:
            raise ValueError(f"unsupported mantissa: {mantissa}")

    def _encode_string(self, s: str) -> int:
        length = len(s)
        if length > MAX_2BYTE:
            raise ValueError("string too long")

        sign = PrefixString
        encoded = s.encode('utf-8')
        if length < StringLen1Byte:
            n = self._write_byte(sign | length)
            n += self._write_bytes(encoded)
            return n
        elif length < Max1Byte:
            n = self._write_byte(sign | StringLen1Byte)
            n += self._write_byte(length)
            n += self._write_bytes(encoded)
            return n
        elif length < Max2Byte:
            n = self._write_byte(sign | StringLen2Byte)
            n += self._write_bytes(bytes([length >> 8, length]))
            n += self._write_bytes(encoded)
            return n
        
        return 0

    def _encode_bytes(self, bs: bytes) -> int:
        length = len(bs)
        if length > MAX_2BYTE:
            raise ValueError("bytes too long")

        sign = PrefixBytes
        if length < BytesLen1Byte:
            n = self._write_byte(sign | length)
            n += self._write_bytes(bs)
            return n
        elif length < Max1Byte:
            n = self._write_byte(sign | BytesLen1Byte)
            n += self._write_byte(length)
            n += self._write_bytes(bs)
            return n
        elif length < Max2Byte:
            n = self._write_byte(sign | BytesLen2Byte)
            n += self._write_bytes(bytes([length >> 8, length]))
            n += self._write_bytes(bs)
            return n
        
        return 0

    def _encode_array(self, items: list) -> int:
        length = len(items)
        if length > MAX_2BYTE:
            raise ValueError("array too long")

        sign = Container | ContainerArray
        if length < 24:
            n = self._write_byte(sign | length)
        elif length < 256:
            n = self._write_byte(sign | 24)
            n += self._write_byte(length)
        elif length < 65536:
            n = self._write_byte(sign | 25)
            n += self._write_bytes(bytes([length >> 8, length]))
        else:
            return 0
        
        for item in items:
            if isinstance(item, Obj):
                n += self._encode_obj(item)
            elif isinstance(item, Arr):
                n += self._encode_arr(item)
            elif isinstance(item, Val):
                n += self._encode_val(item)
            else:
                raise ValueError(f"unsupported type {type(item)}")
        
        return n

    def _encode_map(self, fields: list) -> int:
        length = len(fields)
        if length > MAX_2BYTE:
            raise ValueError("map too long")

        sign = Container | ContainerMap
        if length < 24:
            n = self._write_byte(sign | length)
        elif length < 256:
            n = self._write_byte(sign | 24)
            n += self._write_byte(length)
        elif length < 65536:
            n = self._write_byte(sign | 25)
            n += self._write_bytes(bytes([length >> 8, length]))
        else:
            return 0
        
        for field in fields:
            # Encode key
            n += self._encode_string(field.key)
            # Encode value
            val = field.value
            if isinstance(val, Obj):
                n += self._encode_obj(val)
            elif isinstance(val, Arr):
                n += self._encode_arr(val)
            elif isinstance(val, Val):
                n += self._encode_val(val)
            else:
                raise ValueError(f"unsupported type {type(val)}")
        
        return n

    def _encode_t(self, bs: bytes) -> int:
        length = len(bs)
        if length == 0:
            return 0
        
        if length > MAX_2BYTE:
            raise ValueError("tag too long")

        if length < 254:
            return self._write_bytes_with_prefix(bs, length)
        
        if length < 257:
            return self._write_bytes_with_prefix(bs, 254, length)
        
        return self._write_bytes_with_prefix(bs, 255, length >> 8, length)

    def _encode_tag(self, payload: bytes, tag_bytes: bytes) -> int:
        if len(tag_bytes) == 0:
            return 0

        length = len(payload) + len(tag_bytes)
        if length > MAX_2BYTE:
            raise ValueError(f"tag+payload too long, max: {MAX_2BYTE}, actual: {length}")

        sign = PrefixTag
        if length < TagLen1Byte:
            n = self._write_bytes_with_prefix(tag_bytes, sign | length)
            n += self._write_bytes(payload)
            return n
        elif length < Max1Byte:
            n = self._write_bytes_with_prefix(tag_bytes, sign | TagLen1Byte, length)
            n += self._write_bytes(payload)
            return n
        elif length < Max2Byte:
            n = self._write_bytes_with_prefix(tag_bytes, sign | TagLen2Byte, length >> 8, length)
            n += self._write_bytes(payload)
            return n
        
        return 0

    def _encode_comment(self, payload: bytes, tag) -> int:
        if tag is None or not isinstance(tag, Tag):
            return 0
        tag_bytes = tag.bytes()
        if len(tag_bytes) == 0:
            return 0
        return self._encode_tag(payload, tag_bytes)

    def _encode_obj(self, obj: Obj) -> int:
        return self._encode_map(obj.fields)

    def _encode_arr(self, arr: Arr) -> int:
        return self._encode_array(arr.items)

    def _encode_val(self, val: Val) -> int:
        tag = val.tag if val.tag else Tag()
        val_type = tag.type

        n = 0

        if val_type == ValueType.DateTime:
            if not tag.is_null:
                n = self._encode_datetime(val.data)
        elif val_type == ValueType.Date:
            if not tag.is_null:
                n = self._encode_date(val.data)
        elif val_type == ValueType.Time:
            if not tag.is_null:
                n = self._encode_time(val.data)
        elif val_type in (ValueType.Int, ValueType.Int8, ValueType.Int16, ValueType.Int32, ValueType.Int64):
            if tag.is_null:
                n = self._encode_simple(SimpleNullInt)
            else:
                n = self._encode_int64(int(val.data))
        elif val_type in (ValueType.Uint, ValueType.Uint8, ValueType.Uint16, ValueType.Uint32, ValueType.Uint64):
            if tag.is_null:
                n = self._encode_simple(SimpleNullInt)
            else:
                n = self._encode_uint64(int(val.data))
        elif val_type == ValueType.Float32:
            if not tag.is_null:
                n = self._encode_float(val.text)
        elif val_type == ValueType.Float64:
            if tag.is_null:
                n = self._encode_simple(SimpleNullFloat)
            else:
                n = self._encode_float(val.text)
        elif val_type == ValueType.String:
            if tag.is_null:
                n = self._encode_simple(SimpleNullString)
            else:
                n = self._encode_string(val.text)
        elif val_type == ValueType.Email:
            if not tag.is_null:
                n = self._encode_string(val.text)
        elif val_type == ValueType.UUID:
            if not tag.is_null:
                data = val.data
                if isinstance(data, str):
                    data = data.encode('utf-8')
                arr = data[:16] if len(data) >= 16 else data
                n = self._encode_bytes(arr)
        elif val_type == ValueType.Decimal:
            if not tag.is_null:
                n = self._encode_float(val.text)
        elif val_type == ValueType.URL:
            if not tag.is_null:
                n = self._encode_string(val.text)
        elif val_type == ValueType.IP:
            if not tag.is_null:
                ip = val.data
                if isinstance(ip, str):
                    ip = ip.encode('utf-8')
                if tag.version == 0:
                    n = self._encode_string(val.text)
                elif tag.version == 4:
                    n = self._encode_bytes(ip)
                elif tag.version == 6:
                    if len(val.text) < 16:
                        n = self._encode_string(val.text)
                    else:
                        n = self._encode_bytes(ip)
                else:
                    raise ValueError(f"unsupported IP version: {tag.version}")
        elif val_type == ValueType.Bytes:
            if tag.is_null:
                n = self._encode_simple(SimpleNullBytes)
            else:
                n = self._encode_bytes(val.data)
        elif val_type == ValueType.BigInt:
            if not tag.is_null:
                n = self._encode_string(val.text)
        elif val_type == ValueType.Bool:
            if tag.is_null:
                n = self._encode_simple(SimpleNullBool)
            else:
                n = self._encode_bool(bool(val.data))
        elif val_type == ValueType.Enum:
            if not tag.is_null:
                n = self._encode_int64(int(val.data))
        else:
            raise ValueError(f"unsupported type: {val_type}")

        return n

    def _encode_datetime(self, t) -> int:
        v = int(t.timestamp())
        return self._encode_int64(v)

    def _encode_date(self, t) -> int:
        import datetime
        if hasattr(t, 'date'):
            dt = t
        else:
            dt = t
        days = (dt - datetime.date(1970, 1, 1)).days
        return self._encode_int64(days)

    def _encode_time(self, t) -> int:
        v = t.hour * 3600 + t.minute * 60 + t.second
        return self._encode_uint64(v)


def struct_to_mm(data: Any, path: str = "") -> Any:
    if isinstance(data, dict):
        fields = []
        for k, v in data.items():
            field_path = f"{path}.{k}" if path else k
            fields.append(Field(key=k, value=struct_to_mm(v, field_path)))
        return Obj(fields=fields, tag=Tag())
    elif isinstance(data, (list, tuple)):
        items = [struct_to_mm(item, f"{path}.{i}") for i, item in enumerate(data)]
        return Arr(items=items, tag=Tag(type=ValueType.Array))
    else:
        tag = Tag(type=infer_value_type(data))
        return Val(data=data, text=str(data), tag=tag)


def infer_value_type(data: Any) -> ValueType:
    if data is None:
        return ValueType.Unknown
    elif isinstance(data, bool):
        return ValueType.Bool
    elif isinstance(data, int):
        return ValueType.Int
    elif isinstance(data, float):
        return ValueType.Float64
    elif isinstance(data, str):
        return ValueType.String
    elif isinstance(data, bytes):
        return ValueType.Bytes
    elif isinstance(data, (list, tuple)):
        return ValueType.Array
    elif isinstance(data, dict):
        return ValueType.Struct
    else:
        return ValueType.Unknown


def parse_float(s: str):
    is_negative = s.startswith('-')
    s = s.lstrip('-')
    
    if '.' not in s:
        return is_negative, 0, int(s)
    
    parts = s.split('.')
    int_part = parts[0]
    frac_part = parts[1] if len(parts) > 1 else ''
    
    exponent = -len(frac_part)
    mantissa = int(int_part + frac_part) if int_part else int(frac_part)
    
    return is_negative, exponent, mantissa


class Decoder:
    def __init__(self, data: bytes):
        self.buf = io.BytesIO(data)
        self.offset = 0
    
    def decode(self):
        b = self._read_byte()
        prefix = b >> 5
        
        if prefix == 0b000:  # Simple
            suffix = b & 0b00011111
            return self._decode_simple(suffix)
        elif prefix == 0b001:  # PositiveInt
            suffix = b & 0b00011111
            return self._decode_int(suffix, False)
        elif prefix == 0b010:  # NegativeInt
            suffix = b & 0b00011111
            return self._decode_int(suffix, True)
        elif prefix == 0b011:  # Float
            suffix = b & 0b00011111
            return self._decode_float(suffix)
        elif prefix == 0b100:  # String
            suffix = b & 0b00011111
            return self._decode_string(suffix)
        elif prefix == 0b101:  # Bytes
            suffix = b & 0b00011111
            return self._decode_bytes(suffix)
        elif prefix == 0b110:  # Container
            # 容器类型的长度部分是 4 位，因为容器类型的标记占用了其中的 1 位
            suffix = b & 0b00001111
            return self._decode_container(b, suffix)
        elif prefix == 0b111:  # Tag
            suffix = b & 0b00011111
            return self._decode_tag(suffix)
        else:
            raise ValueError(f"unknown prefix: {prefix}")
    
    def _read_byte(self) -> int:
        b = self.buf.read(1)
        if not b:
            raise ValueError("buffer underflow")
        self.offset += 1
        return b[0]
    
    def _read_bytes(self, n: int) -> bytes:
        bs = self.buf.read(n)
        if len(bs) < n:
            raise ValueError("buffer underflow")
        self.offset += n
        return bs
    
    def _decode_simple(self, suffix: int):
        if suffix == SimpleNullBool:
            return None
        elif suffix == SimpleNullInt:
            return None
        elif suffix == SimpleNullFloat:
            return None
        elif suffix == SimpleNullString:
            return None
        elif suffix == SimpleNullBytes:
            return None
        elif suffix == SimpleFalse:
            return False
        elif suffix == SimpleTrue:
            return True
        else:
            return suffix
    
    def _decode_int(self, suffix: int, is_negative: bool):
        if suffix < IntLen1Byte:
            value = suffix
        elif suffix == IntLen1Byte:
            value = self._read_byte()
        elif suffix == IntLen2Byte:
            bs = self._read_bytes(2)
            value = (bs[0] << 8) | bs[1]
        elif suffix == IntLen3Byte:
            bs = self._read_bytes(3)
            value = (bs[0] << 16) | (bs[1] << 8) | bs[2]
        elif suffix == IntLen4Byte:
            bs = self._read_bytes(4)
            value = struct.unpack('>I', bs)[0]
        elif suffix == IntLen5Byte:
            bs = self._read_bytes(5)
            value = (bs[0] << 32) | struct.unpack('>I', bs[1:])[0]
        elif suffix == IntLen6Byte:
            bs = self._read_bytes(6)
            value = (bs[0] << 40) | (bs[1] << 32) | struct.unpack('>I', bs[2:])[0]
        elif suffix == IntLen7Byte:
            bs = self._read_bytes(7)
            value = (bs[0] << 48) | (bs[1] << 40) | (bs[2] << 32) | struct.unpack('>I', bs[3:])[0]
        elif suffix == IntLen8Byte:
            bs = self._read_bytes(8)
            value = struct.unpack('>Q', bs)[0]
        else:
            raise ValueError(f"invalid int suffix: {suffix}")
        
        if is_negative:
            return -value
        return value
    
    def _decode_float(self, suffix: int):
        is_negative = (suffix & FloatPositiveNegativeMask) != 0
        len_mask = suffix & FloatLenMask
        
        if len_mask < FloatLen1Byte:
            # 处理特殊情况：exponent == -1 and mantissa <= 7
            mantissa = len_mask
            exponent = -1
        elif len_mask == FloatLen1Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            mantissa = self._read_byte()
        elif len_mask == FloatLen2Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(2)
            mantissa = (bs[0] << 8) | bs[1]
        elif len_mask == FloatLen3Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(3)
            mantissa = (bs[0] << 16) | (bs[1] << 8) | bs[2]
        elif len_mask == FloatLen4Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(4)
            mantissa = struct.unpack('>I', bs)[0]
        elif len_mask == FloatLen5Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(5)
            mantissa = (bs[0] << 32) | struct.unpack('>I', bs[1:])[0]
        elif len_mask == FloatLen6Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(6)
            mantissa = (bs[0] << 40) | (bs[1] << 32) | struct.unpack('>I', bs[2:])[0]
        elif len_mask == FloatLen7Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(7)
            mantissa = (bs[0] << 48) | (bs[1] << 40) | (bs[2] << 32) | struct.unpack('>I', bs[3:])[0]
        elif len_mask == FloatLen8Byte:
            exponent = self._read_byte()
            # 将无符号字节转换为有符号整数
            if exponent > 127:
                exponent -= 256
            bs = self._read_bytes(8)
            mantissa = struct.unpack('>Q', bs)[0]
        else:
            raise ValueError(f"invalid float len mask: {len_mask}")
        
        value = mantissa * (10 ** exponent)
        if is_negative:
            value = -value
        return value
    
    def _decode_string(self, suffix: int):
        if suffix < StringLen1Byte:
            length = suffix
        elif suffix == StringLen1Byte:
            length = self._read_byte()
        elif suffix == StringLen2Byte:
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]
        else:
            raise ValueError(f"invalid string suffix: {suffix}")
        
        if length > 0:
            bs = self._read_bytes(length)
            return bs.decode('utf-8')
        return ""
    
    def _decode_bytes(self, suffix: int):
        if suffix < BytesLen1Byte:
            length = suffix
        elif suffix == BytesLen1Byte:
            length = self._read_byte()
        elif suffix == BytesLen2Byte:
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]
        else:
            raise ValueError(f"invalid bytes suffix: {suffix}")
        
        if length > 0:
            return self._read_bytes(length)
        return b""
    
    def _decode_container(self, b: int, suffix: int):
        is_array = (b & ContainerMask) == ContainerArray
        
        if suffix < 24:
            length = suffix
        elif suffix == 24:
            length = self._read_byte()
        elif suffix == 25:
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]
        else:
            raise ValueError(f"invalid container suffix: {suffix}")
        
        if is_array:
            result = []
            for _ in range(length):
                result.append(self.decode())
            return result
        else:
            result = {}
            for _ in range(length):
                key = self.decode()
                value = self.decode()
                result[key] = value
            return result
    
    def _decode_tag(self, suffix: int):
        if suffix < TagLen1Byte:
            length = suffix
        elif suffix == TagLen1Byte:
            length = self._read_byte()
        elif suffix == TagLen2Byte:
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]
        else:
            raise ValueError(f"invalid tag suffix: {suffix}")
        
        if length > 0:
            tag_bytes = self._read_bytes(length)
            tag_str = tag_bytes.decode('utf-8')
            value = self.decode()
            return {"tag": tag_str, "value": value}
        return self.decode()