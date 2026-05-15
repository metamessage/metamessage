"""
Encoder for metamessage binary format.
Based on Go implementation in internal/core/encode*.go
"""

import io
import math
import struct
from datetime import datetime, date, time as dt_time, timezone
from typing import Any, Optional, Union

from ..ir.tag import Tag, ValueType
from ..ir.types import Obj, Arr, Val, Field, Node, NodeType

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

# Bytes constants
BytesLenMask  = 0b11111
BytesLen1Byte = BytesLenMask - 1  # 30
BytesLen2Byte = BytesLenMask      # 31

# Container constants
ContainerMask    = 0b10000
ContainerObject  = 0b00000
ContainerArray   = 0b10000
ContainerLenMask  = 0b01111
ContainerLen1Byte = ContainerLenMask - 1  # 14 (for length >= 12 and < 256)
ContainerLen2Byte = ContainerLenMask      # 15 (for length >= 256 and < 65536)

# Tag constants
TagLenMask  = 0b11111
TagLen1Byte = TagLenMask - 1  # 30 (for length >= 26 and < 256)
TagLen2Byte = TagLenMask      # 31 (for length >= 256 and < 65536)

DEFAULT_BUF_SIZE = 1024
MAX_CAP = 1024 * 1024 * 1024  # 1GB


def get_prefix(b: int) -> int:
    return b & 0b11100000


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
        l = len(s)
        if l == 0:
            return 0
        self._ensure_capacity(l)
        encoded = s.encode('utf-8')
        self.buf[self.offset:self.offset + l] = encoded
        self.offset += l
        return l

    def _write_bytes_with_prefix(self, bs: bytes, *prefix: int) -> int:
        n = 0
        for p in prefix:
            n += self._write_byte(p)
        n += self._write_bytes(bs)
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

    def _encode_int(self, sign: int, uv: int) -> int:
        if uv < IntLen1Byte:
            return self._write_byte(sign | uv)
        elif uv <= Max1Byte:
            return self._write_bytes_with_prefix(
                bytes([uv & 0xFF]),
                sign | IntLen1Byte
            )
        elif uv <= Max2Byte:
            return self._write_bytes_with_prefix(
                bytes([(uv >> 8) & 0xFF, uv & 0xFF]),
                sign | IntLen2Byte
            )
        elif uv <= Max3Byte:
            return self._write_bytes_with_prefix(
                bytes([(uv >> 16) & 0xFF, (uv >> 8) & 0xFF, uv & 0xFF]),
                sign | IntLen3Byte
            )
        elif uv <= Max4Byte:
            return self._write_bytes_with_prefix(
                struct.pack('>I', uv),
                sign | IntLen4Byte
            )
        elif uv <= Max5Byte:
            return self._write_bytes_with_prefix(
                struct.pack('>Q', uv)[3:],
                sign | IntLen5Byte
            )
        elif uv <= Max6Byte:
            return self._write_bytes_with_prefix(
                struct.pack('>Q', uv)[2:],
                sign | IntLen6Byte
            )
        elif uv <= Max7Byte:
            return self._write_bytes_with_prefix(
                struct.pack('>Q', uv)[1:],
                sign | IntLen7Byte
            )
        elif uv <= Max8Byte:
            return self._write_bytes_with_prefix(
                struct.pack('>Q', uv),
                sign | IntLen8Byte
            )
        else:
            raise ValueError("integer value too large")

    def _encode_uint64(self, uv: int) -> int:
        return self._encode_int(PositiveInt, uv)

    def _encode_int64(self, v: int) -> int:
        if v >= 0:
            return self._encode_int(PositiveInt, v)
        else:
            if v == -9223372036854775808:
                uv = 9223372036854775808
            else:
                uv = -v
            return self._encode_int(NegativeInt, uv)

    def _encode_big_int(self, s: str) -> int:
        # Write length byte first
        self._ensure_capacity(1)
        # We'll write the length byte after encoding the big int bytes
        # Encode big int (simplified - storing as string for now)
        encoded = s.encode('utf-8')
        if len(encoded) > Max2Byte:
            raise ValueError("big int string too long")
        # Write length prefix
        n = self._write_byte(len(encoded))
        n += self._write_bytes(encoded)
        # Wrap in bytes encoding
        return self._encode_bytes(self.buf[self.offset - n:self.offset])

    # ===== Floats =====

    @staticmethod
    def _parse_float(s: str):
        """Parse float string into (is_negative, exponent, mantissa)."""
        is_negative = s.startswith('-')
        s = s.lstrip('-')

        if 'e' in s or 'E' in s:
            import re
            m = re.match(r'^([\d.]+)[eE]([+-]?\d+)$', s)
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
            return self._write_bytes_with_prefix(
                bytes([exp_byte, mantissa & 0xFF]),
                sign | FloatLen1Byte
            )
        elif mantissa <= Max2Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte, (mantissa >> 8) & 0xFF, mantissa & 0xFF]),
                sign | FloatLen2Byte
            )
        elif mantissa <= Max3Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte, (mantissa >> 16) & 0xFF, (mantissa >> 8) & 0xFF, mantissa & 0xFF]),
                sign | FloatLen3Byte
            )
        elif mantissa <= Max4Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte]) + struct.pack('>I', mantissa),
                sign | FloatLen4Byte
            )
        elif mantissa <= Max5Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[3:],
                sign | FloatLen5Byte
            )
        elif mantissa <= Max6Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[2:],
                sign | FloatLen6Byte
            )
        elif mantissa <= Max7Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa)[1:],
                sign | FloatLen7Byte
            )
        elif mantissa <= Max8Byte:
            return self._write_bytes_with_prefix(
                bytes([exp_byte]) + struct.pack('>Q', mantissa),
                sign | FloatLen8Byte
            )
        else:
            raise ValueError(f"mantissa too large: {mantissa}")

    # ===== Strings =====

    def _encode_string(self, s: str) -> int:
        length = len(s)
        if length > Max2Byte:
            raise ValueError(f"string too long: {length}")

        sign = PrefixString
        encoded = s.encode('utf-8')
        if length < StringLen1Byte:  # length < 30
            n = self._write_byte(sign | length)
            n += self._write_bytes(encoded)
            return n
        elif length < Max1Byte:  # length < 256
            n = self._write_byte(sign | StringLen1Byte)
            n += self._write_byte(length)
            n += self._write_bytes(encoded)
            return n
        else:  # length < 65536
            n = self._write_byte(sign | StringLen2Byte)
            n += self._write_bytes(bytes([(length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_bytes(encoded)
            return n

    # ===== Bytes =====

    def _encode_bytes(self, bs: bytes) -> int:
        length = len(bs)
        if length > Max2Byte:
            raise ValueError(f"bytes too long: {length}")

        sign = PrefixBytes
        if length < BytesLen1Byte:  # length < 30
            n = self._write_byte(sign | length)
            n += self._write_bytes(bs)
            return n
        elif length < Max1Byte:  # length < 256
            n = self._write_byte(sign | BytesLen1Byte)
            n += self._write_byte(length)
            n += self._write_bytes(bs)
            return n
        else:  # length < 65536
            n = self._write_byte(sign | BytesLen2Byte)
            n += self._write_bytes(bytes([(length >> 8) & 0xFF, length & 0xFF]))
            n += self._write_bytes(bs)
            return n

    # ===== Containers =====

    def _encode_array(self, bs: bytes) -> int:
        """Encode raw array bytes with container prefix."""
        length = len(bs)
        if length > Max2Byte:
            raise ValueError(f"array too long: {length}")

        sign = Container | ContainerArray
        if length < ContainerLen1Byte:  # length < 14
            return self._write_bytes_with_prefix(bs, sign | length)
        elif length < Max1Byte:  # length < 256
            return self._write_bytes_with_prefix(bs, sign | ContainerLen1Byte, length)
        else:  # length < 65536
            return self._write_bytes_with_prefix(
                bs, sign | ContainerLen2Byte,
                (length >> 8) & 0xFF, length & 0xFF
            )

    def _encode_object(self, bs: bytes) -> int:
        """Encode raw object bytes with container prefix."""
        length = len(bs)
        if length > Max2Byte:
            raise ValueError(f"object too long: {length}")

        sign = Container | ContainerObject
        if length < ContainerLen1Byte:  # length < 14
            return self._write_bytes_with_prefix(bs, sign | length)
        elif length < Max1Byte:  # length < 256
            return self._write_bytes_with_prefix(bs, sign | ContainerLen1Byte, length)
        else:  # length < 65536
            return self._write_bytes_with_prefix(
                bs, sign | ContainerLen2Byte,
                (length >> 8) & 0xFF, length & 0xFF
            )

    # ===== Tags =====

    def _encode_t(self, bs: bytes) -> int:
        """Encode tag bytes with length prefix."""
        length = len(bs)
        if length == 0:
            return 0
        if length > Max2Byte:
            raise ValueError(f"tag too long: {length}")

        if length < 254:
            return self._write_bytes_with_prefix(bs, length)
        elif length < 257:
            return self._write_bytes_with_prefix(bs, 254, length)
        else:
            return self._write_bytes_with_prefix(
                bs, 255,
                (length >> 8) & 0xFF, length & 0xFF
            )

    def _encode_tag(self, payload: bytes, tag_bytes: bytes) -> int:
        """Encode payload with tag metadata.
        
        Format: [PrefixTag|total_len][tag_data_len][tag_data...][payload...]
        total_len = 1 (tag_data_len byte) + len(tag_data) + len(payload)
        """
        if len(tag_bytes) == 0:
            return 0

        total_length = 1 + len(tag_bytes) + len(payload)
        if total_length > Max2Byte:
            raise ValueError(f"tag+payload too long: {total_length}")

        sign = PrefixTag
        if total_length < TagLen1Byte:  # total_length < 30
            n = self._write_byte(sign | total_length)
        elif total_length < Max1Byte:  # total_length < 256
            n = self._write_bytes_with_prefix(b'', sign | TagLen1Byte, total_length)
        else:  # total_length < 65536
            n = self._write_bytes_with_prefix(
                b'', sign | TagLen2Byte,
                (total_length >> 8) & 0xFF, total_length & 0xFF
            )
        
        # Write tag data with length prefix
        n += self._write_byte(len(tag_bytes))
        n += self._write_bytes(tag_bytes)
        
        # Write payload
        n += self._write_bytes(payload)
        return n

    def _encode_comment(self, payload: bytes, tag: Optional[Tag]) -> int:
        """Wrap payload with tag metadata if tag has properties."""
        if tag is None:
            return 0
        tag_bytes = tag.bytes()
        if len(tag_bytes) == 0:
            return 0
        return self._encode_tag(payload, tag_bytes)

    # ===== Node encoding =====

    def _encode_node_object(self, obj: Obj) -> int:
        """Encode an Obj node.
        
        Binary layout: [container_prefix][key_array][values...]
        Where key_array is an array of encoded field keys.
        """
        tag = obj.get_tag()
        tag_bytes = tag.bytes() if tag else b''

        # Strategy: encode everything sequentially to buf, tracking offsets
        # Format: [key_array][value1][value2]...
        # But we need key_array before values. So encode values first,
        # track them, then encode key array, then combine.

        # Step 1: Encode all values and keys to temp storage
        start = self.offset
        values_data = bytearray()
        keys_data = bytearray()

        for field in obj.fields:
            # Encode value
            n = self._encode_any_node(field.value)
            values_data.extend(bytes(self.buf[self.offset - n:self.offset]))

            # Encode key
            nk = self._encode_string(field.key)
            keys_data.extend(bytes(self.buf[self.offset - nk:self.offset]))

        # Step 2: Encode key array
        nk = self._encode_array(bytes(keys_data))

        # Now buffer has: [values][keys][array_header]
        # We need: [array_header][keys][values]
        # So let's build the correct data and overwrite

        # Read the key array encoding
        key_array = bytes(self.buf[self.offset - nk:self.offset])

        # Build final combined data: key_array + values
        combined = key_array + bytes(values_data)

        # Rewind to before we started
        self.offset = start

        # Write final object
        n = self._encode_object(combined)

        # Wrap with tag
        if len(tag_bytes) > 0:
            payload = bytes(self.buf[self.offset - n:self.offset])
            n = self._encode_tag(payload, tag_bytes)

        return n

    def _encode_node_array(self, arr: Arr) -> int:
        """Encode an Arr node."""
        tag = arr.get_tag()
        tag_bytes = tag.bytes() if tag else b''

        start = self.offset

        for item in arr.items:
            self._encode_any_node(item)

        # Now buffer has all items sequentially
        data = bytes(self.buf[start:self.offset])

        # Rewind
        self.offset = start

        # Wrap in array container
        n = self._encode_array(data)

        # Wrap with tag
        if len(tag_bytes) > 0:
            payload = bytes(self.buf[self.offset - n:self.offset])
            n = self._encode_tag(payload, tag_bytes)

        return n

    def _encode_any_node(self, node: Node) -> int:
        if isinstance(node, Obj):
            return self._encode_node_object(node)
        elif isinstance(node, Arr):
            return self._encode_node_array(node)
        elif isinstance(node, Val):
            return self._encode_node_value(node)
        else:
            raise ValueError(f"unsupported node type: {type(node)}")

    def _encode_node_value(self, val: Val) -> int:
        """Encode a Val node."""
        tag = val.get_tag()
        if tag is None:
            tag = Tag()

        tag_bytes = tag.bytes()
        start = self.offset

        # Handle null first - if is_null, encode null simple value
        if tag.is_null:
            # Use String null as default for unknown types
            self._encode_simple(SimpleNullString)
        elif tag.type == ValueType.DateTime:
            self._encode_datetime(val.data)
        elif tag.type == ValueType.Date:
                self._encode_date(val.data)
        elif tag.type == ValueType.Time:
                self._encode_time(val.data)
        elif tag.type in (ValueType.Int, ValueType.Int8, ValueType.Int16, ValueType.Int32, ValueType.Int64):
                self._encode_int64(int(val.data))
        elif tag.type in (ValueType.Uint, ValueType.Uint8, ValueType.Uint16, ValueType.Uint32, ValueType.Uint64):
                self._encode_uint64(int(val.data))
        elif tag.type == ValueType.Float32:
                self._encode_float(val.text)
        elif tag.type == ValueType.Float64:
                self._encode_float(val.text)
        elif tag.type == ValueType.String:
                self._encode_string(val.text)
        elif tag.type == ValueType.Email:
                self._encode_string(val.text)
        elif tag.type == ValueType.UUID:
            data = val.data
            if isinstance(data, str):
                import binascii
                raw = data.replace('-', '')
                data_bytes = binascii.unhexlify(raw)
                self._encode_bytes(data_bytes)
            elif isinstance(data, (bytes, bytearray)):
                self._encode_bytes(bytes(data[:16]))
            else:
                self._encode_bytes(bytes(data)[:16] if hasattr(data, '__len__') else bytes(16))
        elif tag.type == ValueType.Decimal:
                self._encode_float(val.text)
        elif tag.type == ValueType.URL:
                self._encode_string(val.text)
        elif tag.type == ValueType.IP:
            ip = val.data
            if isinstance(ip, str):
                pass
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
            data = val.data
            if isinstance(data, str):
                data = data.encode('utf-8')
            elif not isinstance(data, (bytes, bytearray)):
                data = bytes(data)
            self._encode_bytes(bytes(data))
        elif tag.type == ValueType.BigInt:
                self._encode_big_int(val.text)
        elif tag.type == ValueType.Bool:
                self._encode_bool(bool(val.data))
        elif tag.type == ValueType.Enum:
                self._encode_int64(int(val.data))
        elif tag.type == ValueType.Unknown:
            # Unknown type - encode as simple value or string
            if val.data is None:
                self._encode_simple(SimpleNullString)
            elif isinstance(val.data, str):
                self._encode_string(val.text)
            elif isinstance(val.data, bool):
                self._encode_bool(val.data)
            elif isinstance(val.data, int):
                self._encode_int64(val.data)
            elif isinstance(val.data, float):
                self._encode_float(val.text)
            else:
                self._encode_string(str(val.data))
        else:
            raise ValueError(f"unsupported value type: {tag.type}")

        n = self.offset - start
        if n == 0:
            raise ValueError(f"failed to encode value type: {tag.type}")

        # Wrap with tag if tag has properties
        if len(tag_bytes) > 0:
            payload = bytes(self.buf[start:self.offset])
            self.offset = start
            n = self._encode_tag(payload, tag_bytes)

        return n

    # ===== Date/time encoding =====

    def _encode_datetime(self, t):
        if isinstance(t, (datetime, dt_time)):
            v = int(t.timestamp())
        else:
            v = int(t)
        self._encode_int64(v)

    def _encode_date(self, t):
        if isinstance(t, (datetime, date)):
            d = t if isinstance(t, date) else t.date()
            ref = date(1970, 1, 1)
            days = (d - ref).days
        else:
            days = int(t)
        sign = PositiveInt if days >= 0 else NegativeInt
        uv = abs(days)
        self._encode_int(sign, uv)

    def _encode_time(self, t):
        if isinstance(t, dt_time):
            v = t.hour * 3600 + t.minute * 60 + t.second
        elif isinstance(t, datetime):
            v = t.hour * 3600 + t.minute * 60 + t.second
        else:
            v = int(t)
        self._encode_uint64(v)

    # ===== Public API =====

    def get_encoded_bytes(self, written: int) -> bytes:
        return bytes(self.buf[self.offset - written:self.offset])

    def encode(self, node: Node) -> bytes:
        """Encode a Node tree to binary format."""
        self.offset = 0
        self._encode_any_node(node)
        result = bytes(self.buf[:self.offset])
        self.offset = 0
        return result

    def encode_stream(self, data: Any) -> int:
        """Encode Python data to binary and write to writer."""
        if self.w is None:
            raise ValueError("writer cannot be None")

        node = self._struct_to_mm(data)
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
            return Obj(fields=fields, tag=Tag(), path=path)
        elif isinstance(data, list):
            items = [self._struct_to_mm(item, f"{path}.{i}") for i, item in enumerate(data)]
            return Arr(items=items, tag=Tag(type=ValueType.Array), path=path)
        elif isinstance(data, bool):
            return Val(data=data, text="true" if data else "false", tag=Tag(type=ValueType.Bool), path=path)
        elif isinstance(data, int):
            return Val(data=data, text=str(data), tag=Tag(type=ValueType.Int), path=path)
        elif isinstance(data, float):
            return Val(data=data, text=str(data), tag=Tag(type=ValueType.Float64), path=path)
        elif isinstance(data, str):
            return Val(data=data, text=data, tag=Tag(type=ValueType.String), path=path)
        elif isinstance(data, bytes):
            return Val(data=data, text="", tag=Tag(type=ValueType.Bytes), path=path)
        elif data is None:
            return Val(data=None, text="null", tag=Tag(is_null=True), path=path)
        else:
            return Val(data=data, text=str(data), tag=Tag(type=ValueType.String), path=path)
