"""
Decoder for metamessage binary format.
Based on Go implementation in internal/core/decode*.go
"""

import io
import struct
from datetime import datetime, timezone, date, time as dt_time
from typing import Any, Optional

from ..ir.tag import Tag, TagKey, ValueType
from ..ir.types import Obj, Arr, Val, Field, Node, NodeType

# ===== Constants (matching Go internal/core/constants.go) =====

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


def get_prefix(b: int) -> int:
    return b & PrefixMask


def get_suffix(b: int) -> int:
    return b & SuffixMask


def tag_len(b: int):
    """Return (extra_bytes, prefix_length) for tag."""
    l = b & TagLenMask
    if l == TagLen1Byte:
        return 1, 0
    elif l == TagLen2Byte:
        return 2, 0
    else:
        return 0, l


def container_len(b: int):
    """Return (extra_bytes, prefix_length) for container."""
    l = b & ContainerLenMask
    if l == ContainerLen1Byte:
        return 1, 0
    elif l == ContainerLen2Byte:
        return 2, 0
    else:
        return 0, l


def is_array(b: int) -> bool:
    return (b & ContainerMask) == ContainerArray


def string_len(b: int):
    """Return (extra_bytes, prefix_length) for string."""
    l = b & StringLenMask
    if l < StringLen1Byte:
        return 0, l
    elif l == StringLen1Byte:
        return 1, 0
    else:  # StringLen2Byte
        return 2, 0


def bytes_len(b: int):
    """Return (extra_bytes, prefix_length) for bytes."""
    l = b & BytesLenMask
    if l < BytesLen1Byte:
        return 0, l
    elif l == BytesLen1Byte:
        return 1, 0
    else:  # BytesLen2Byte
        return 2, 0


def int_len(b: int):
    """Return (extra_bytes, prefix_length) for int."""
    l = b & IntLenMask
    if l < IntLen1Byte:
        return 0, l
    else:
        return l - IntLen1Byte + 1, 0


def float_len(b: int):
    """Return (extra_bytes, prefix_length) for float."""
    l = b & FloatLenMask
    if l < FloatLen1Byte:
        return 0, l
    else:
        return l - FloatLen1Byte + 1, 0


# ===== Tag Binary Decoder =====

def decode_tag(data: bytes) -> Tag:
    """Decode tag metadata from binary bytes."""
    tag = Tag()
    offset = 0

    while offset < len(data):
        key_byte = data[offset]
        offset += 1
        key = key_byte & 0xF8  # TagKey (upper 5 bits)
        length_code = key_byte & 0x07  # lower 3 bits

        if key == TagKey.IsNull:
            tag.is_null = True
            tag.nullable = True
        elif key == TagKey.Example:
            tag.example = True
        elif key == TagKey.Nullable:
            tag.nullable = True
        elif key == TagKey.Raw:
            tag.raw = True
        elif key == TagKey.AllowEmpty:
            tag.allow_empty = True
        elif key == TagKey.Unique:
            tag.unique = True
        elif key == TagKey.ChildRaw:
            tag.child_raw = True
        elif key == TagKey.ChildNullable:
            tag.child_nullable = True
        elif key == TagKey.ChildAllowEmpty:
            tag.child_allow_empty = True
        elif key == TagKey.ChildUnique:
            tag.child_unique = True
        elif key == TagKey.Type:
            tag.type = ValueType(data[offset]) if offset < len(data) else ValueType.Unknown
            offset += 1
        elif key == TagKey.ChildType:
            tag.child_type = ValueType(data[offset]) if offset < len(data) else ValueType.Unknown
            offset += 1
        elif key == TagKey.Size:
            extra, val = _decode_uint64_value(data, offset, length_code)
            tag.size = val
            offset += extra
        elif key == TagKey.ChildSize:
            extra, val = _decode_uint64_value(data, offset, length_code)
            tag.child_size = val
            offset += extra
        elif key == TagKey.Version:
            extra, val = _decode_uint64_value(data, offset, length_code)
            tag.version = val
            offset += extra
        elif key == TagKey.ChildVersion:
            extra, val = _decode_uint64_value(data, offset, length_code)
            tag.child_version = val
            offset += extra
        elif key == TagKey.Desc:
            length, consumed = _decode_length_offset(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.desc = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.Default:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.default = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.Min:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.min = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.Max:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.max = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.Enum:
            length, consumed = _decode_length_offset(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.enum = data[offset:offset + length].decode('utf-8')
                tag.type = ValueType.Enum
                offset += length
        elif key == TagKey.Pattern:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.pattern = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.Location:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                loc_str = data[offset:offset + length].decode('utf-8')
                try:
                    tag.location = int(loc_str)
                except ValueError:
                    tag.location = loc_str
                offset += length
        elif key == TagKey.Mime:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.mime = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildDesc:
            length, consumed = _decode_length_offset(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_desc = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildDefault:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_default = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildMin:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_min = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildMax:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_max = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildEnum:
            length, consumed = _decode_length_offset(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_enum = data[offset:offset + length].decode('utf-8')
                tag.child_type = ValueType.Enum
                offset += length
        elif key == TagKey.ChildPattern:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_pattern = data[offset:offset + length].decode('utf-8')
                offset += length
        elif key == TagKey.ChildLocation:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                loc_str = data[offset:offset + length].decode('utf-8')
                try:
                    tag.child_location = int(loc_str)
                except ValueError:
                    tag.child_location = loc_str
                offset += length
        elif key == TagKey.ChildMime:
            length, consumed = _decode_length_offset_small(length_code, data, offset)
            offset += consumed
            if length > 0:
                tag.child_mime = data[offset:offset + length].decode('utf-8')
                offset += length
        else:
            # Unknown key, skip based on length code
            # For flag-like keys (0 length), skip
            pass

    return tag


def _decode_length(length_code: int, data: bytes, offset: int) -> int:
    """Decode length from length_code and optional extra bytes."""
    if length_code <= 5:
        return length_code
    elif length_code == 6:
        # 1 byte length
        return data[offset]
    elif length_code == 7:
        # 2 byte length
        return (data[offset] << 8) | data[offset + 1]
    return 0


def _decode_length_offset(length_code: int, data: bytes, offset: int) -> tuple:
    """Decode length from length_code and optional extra bytes, returning (length, bytes_consumed)."""
    if length_code <= 5:
        return length_code, 0
    elif length_code == 6:
        return data[offset], 1
    elif length_code == 7:
        return (data[offset] << 8) | data[offset + 1], 2
    return 0, 0


def _decode_length_small(length_code: int, data: bytes, offset: int) -> int:
    """Decode length with small values (up to 6 bytes inline, 7=1 extra byte)."""
    if length_code < 7:
        return length_code
    else:
        return data[offset]


def _decode_length_offset_small(length_code: int, data: bytes, offset: int) -> tuple:
    """Decode length with small values, returning (length, bytes_consumed)."""
    if length_code < 7:
        return length_code, 0
    else:
        return data[offset], 1


def _decode_uint64_value(data: bytes, offset: int, extra_len: int):
    """Decode uint64 value from data at offset, given extra byte count."""
    if extra_len == 0:
        return 1, data[offset]
    elif extra_len == 1:
        return 2, data[offset] | (data[offset + 1] << 8)
    elif extra_len == 2:
        return 3, data[offset] | (data[offset + 1] << 8) | (data[offset + 2] << 16)
    elif extra_len == 3:
        return 4, struct.unpack_from('>I', data, offset)[0]
    elif extra_len == 4:
        return 5, struct.unpack_from('>Q', data, offset + 1)[0] | (data[offset] << 32)
    elif extra_len == 5:
        return 6, (data[offset] << 40) | (data[offset + 1] << 32) | struct.unpack_from('>I', data, offset + 2)[0]
    elif extra_len == 6:
        return 7, (data[offset] << 48) | (data[offset + 1] << 40) | (data[offset + 2] << 32) | struct.unpack_from('>I', data, offset + 3)[0]
    elif extra_len == 7:
        return 8, struct.unpack_from('>Q', data, offset)[0]
    else:
        return 1, 0


# ===== Main Decoder =====

class Decoder:
    def __init__(self, data: bytes):
        self.data = data
        self.offset = 0

    def decode(self) -> Any:
        """Decode binary data to Python objects. Returns decoded value."""
        if self.offset >= len(self.data):
            raise ValueError("buffer underflow")

        b = self._read_byte()
        prefix = get_prefix(b)

        if prefix == Simple:
            value = self._decode_simple(get_suffix(b))
            # Simple values are raw values, no tag wrapping
            return value
        elif prefix == PositiveInt:
            value = self._decode_int(get_suffix(b), False)
            return value
        elif prefix == NegativeInt:
            value = self._decode_int(get_suffix(b), True)
            return value
        elif prefix == PrefixFloat:
            return self._decode_float(get_suffix(b))
        elif prefix == PrefixString:
            return self._decode_string(get_suffix(b))
        elif prefix == PrefixBytes:
            return self._decode_bytes(get_suffix(b))
        elif prefix == Container:
            suffix = b & 0b00001111  # Container uses 4 bits for length
            return self._decode_container(b, suffix)
        elif prefix == PrefixTag:
            value = self._decode_tag_wrapper(get_suffix(b))
            return value
        else:
            raise ValueError(f"unknown prefix: {prefix}")

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
            # Common simple values like code, message, etc
            # Return as simple integer (can be mapped later)
            return suffix

    def _decode_int(self, suffix: int, is_negative: bool):
        extra, prefix_val = int_len(suffix)
        if extra == 0:
            value = prefix_val
        else:
            if extra == 1:
                value = self._read_byte()
            elif extra == 2:
                bs = self._read_bytes(2)
                value = (bs[0] << 8) | bs[1]
            elif extra == 3:
                bs = self._read_bytes(3)
                value = (bs[0] << 16) | (bs[1] << 8) | bs[2]
            elif extra == 4:
                bs = self._read_bytes(4)
                value = struct.unpack('>I', bs)[0]
            elif extra == 5:
                bs = self._read_bytes(5)
                value = (bs[0] << 32) | struct.unpack('>I', bs[1:])[0]
            elif extra == 6:
                bs = self._read_bytes(6)
                value = (bs[0] << 40) | (bs[1] << 32) | struct.unpack('>I', bs[2:])[0]
            elif extra == 7:
                bs = self._read_bytes(7)
                value = (bs[0] << 48) | (bs[1] << 40) | (bs[2] << 32) | struct.unpack('>I', bs[3:])[0]
            elif extra == 8:
                bs = self._read_bytes(8)
                value = struct.unpack('>Q', bs)[0]
            else:
                raise ValueError(f"invalid int length: {extra}")

        if is_negative:
            return -value
        return value

    def _decode_float(self, suffix: int):
        is_negative = (suffix & FloatPositiveNegativeMask) != 0
        l = suffix & FloatLenMask

        if l < FloatLen1Byte:
            mantissa = l
            exponent = -1
            result = mantissa * (10 ** exponent)
            return -result if is_negative else result

        extra = l - FloatLen1Byte + 1

        exponent_byte = self._read_byte()
        exponent = exponent_byte if exponent_byte < 128 else exponent_byte - 256

        if extra == 0:
            mantissa = 0
        elif extra == 1:
            mantissa = self._read_byte()
        elif extra == 2:
            bs = self._read_bytes(2)
            mantissa = (bs[0] << 8) | bs[1]
        elif extra == 3:
            bs = self._read_bytes(3)
            mantissa = (bs[0] << 16) | (bs[1] << 8) | bs[2]
        elif extra == 4:
            bs = self._read_bytes(4)
            mantissa = struct.unpack('>I', bs)[0]
        elif extra == 5:
            bs = self._read_bytes(5)
            mantissa = (bs[0] << 32) | struct.unpack('>I', bs[1:])[0]
        elif extra == 6:
            bs = self._read_bytes(6)
            mantissa = (bs[0] << 40) | (bs[1] << 32) | struct.unpack('>I', bs[2:])[0]
        elif extra == 7:
            bs = self._read_bytes(7)
            mantissa = (bs[0] << 48) | (bs[1] << 40) | (bs[2] << 32) | struct.unpack('>I', bs[3:])[0]
        elif extra == 8:
            bs = self._read_bytes(8)
            mantissa = struct.unpack('>Q', bs)[0]
        else:
            raise ValueError(f"invalid float length: {extra}")

        result = mantissa * (10 ** exponent)
        return -result if is_negative else result

    def _decode_string(self, suffix: int):
        extra, prefix_val = string_len(suffix)
        if extra == 0:
            length = prefix_val
        elif extra == 1:
            length = self._read_byte()
        else:  # extra == 2
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]

        if length > 0:
            bs = self._read_bytes(length)
            return bs.decode('utf-8')
        return ""

    def _decode_bytes(self, suffix: int):
        extra, prefix_val = bytes_len(suffix)
        if extra == 0:
            length = prefix_val
        elif extra == 1:
            length = self._read_byte()
        else:  # extra == 2
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]

        if length > 0:
            return self._read_bytes(length)
        return b""

    def _decode_container(self, b: int, suffix: int):
        is_arr = is_array(b)
        extra, prefix_val = container_len(suffix)
        if extra == 0:
            length = prefix_val
        elif extra == 1:
            length = self._read_byte()
        else:  # extra == 2
            bs = self._read_bytes(2)
            length = (bs[0] << 8) | bs[1]

        if is_arr:
            return self._decode_array_data(length)
        else:
            return self._decode_object_data(length)

    def _decode_array_data(self, length: int) -> list:
        """Decode array items from data with given byte length."""
        result = []
        start = self.offset
        end = start + length
        while self.offset < end:
            result.append(self.decode())
        return result

    def _decode_object_data(self, length: int) -> dict:
        """Decode object: first key array (as proper array), then values."""
        # First byte after header is the array prefix for keys
        l_array_byte = self._read_byte()
        l_array_offset = self.offset - 1
        temp_offset = self.offset
        keys = []
        
        # Manually decode the key array without relying on general decode
        # Read array content length
        arr_extra, arr_prefix = container_len(l_array_byte & ContainerLenMask)
        if arr_extra == 0:
            arr_len = arr_prefix
        elif arr_extra == 1:
            arr_len = self._read_byte()
        else:
            bs = self._read_bytes(2)
            arr_len = (bs[0] << 8) | bs[1]
        
        # Decode keys from array content
        arr_start = self.offset
        arr_end = arr_start + arr_len
        while self.offset < arr_end:
            keys.append(self.decode())
        
        # Now decode values for each key
        result = {}
        for key in keys:
            result[key] = self.decode()
        
        return result

    def _decode_tag_wrapper(self, suffix: int):
        """Decode a tag-wrapped value."""
        extra, prefix_val = tag_len(suffix)
        if extra == 0:
            tag_length = prefix_val
        elif extra == 1:
            tag_length = self._read_byte()
        else:  # extra == 2
            bs = self._read_bytes(2)
            tag_length = (bs[0] << 8) | bs[1]

        # Read tag bytes (first byte is length of tag data itself)
        tag_data_length_byte = self._read_byte()
        if tag_data_length_byte < 254:
            tag_data_len = tag_data_length_byte
            payload_remaining = tag_length - 1
        elif tag_data_length_byte == 254:
            tag_data_len = self._read_byte()
            payload_remaining = tag_length - 2
        else:  # 255
            bs = self._read_bytes(2)
            tag_data_len = (bs[0] << 8) | bs[1]
            payload_remaining = tag_length - 3

        if tag_data_len > 0:
            tag_bytes = self._read_bytes(tag_data_len)
        else:
            tag_bytes = b""

        tag = decode_tag(tag_bytes)

        # Now decode the actual value
        value = self.decode()

        return value
