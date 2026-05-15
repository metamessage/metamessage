"""
Tests for metamessage encoder.
Based on Go tests in internal/core/encode_*_test.go
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from metamessage.core.encoder import Encoder
from metamessage.core.decoder import Decoder
from metamessage.ir.tag import Tag, ValueType
from metamessage.ir.types import Obj, Arr, Val, Field


def test_encode_bool():
    """Test bool encoding matches Go encode_bool_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # True
    b = enc.encode(Val(True, 'true', Tag(type=ValueType.Bool)))
    assert len(b) == 1
    assert b[0] == 0x06  # SimpleTrue
    assert dec(b).decode() == True
    
    # False
    b = enc.encode(Val(False, 'false', Tag(type=ValueType.Bool)))
    assert len(b) == 1
    assert b[0] == 0x05  # SimpleFalse
    assert dec(b).decode() == False


def test_encode_int():
    """Test int encoding matches Go encode_int_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Small positive int (fits in 5 bits)
    b = enc.encode(Val(0, '0', Tag(type=ValueType.Int)))
    assert b[0] == 0x20  # PositiveInt | 0
    
    b = enc.encode(Val(5, '5', Tag(type=ValueType.Int)))
    assert b[0] == 0x25  # PositiveInt | 5
    
    b = enc.encode(Val(23, '23', Tag(type=ValueType.Int)))
    assert b[0] == 0x37  # PositiveInt | 23
    
    # Medium int (1 extra byte)
    b = enc.encode(Val(24, '24', Tag(type=ValueType.Int)))
    assert b[0] == 0x38  # PositiveInt | IntLen1Byte (24)
    assert b[1] == 24
    
    b = enc.encode(Val(255, '255', Tag(type=ValueType.Int)))
    assert b[0] == 0x38
    assert b[1] == 255
    
    # Large int (2 extra bytes)
    b = enc.encode(Val(256, '256', Tag(type=ValueType.Int)))
    assert b[0] == 0x39  # PositiveInt | IntLen2Byte (25)
    assert len(b) == 3
    
    # Negative int
    b = enc.encode(Val(-1, '-1', Tag(type=ValueType.Int)))
    assert b[0] & 0xE0 == 0x40  # NegativeInt
    assert dec(b).decode() == -1
    
    b = enc.encode(Val(-24, '-24', Tag(type=ValueType.Int)))
    assert b[0] & 0xE0 == 0x40
    assert dec(b).decode() == -24
    
    # Edge cases
    b = enc.encode(Val(9223372036854775807, '9223372036854775807', Tag(type=ValueType.Int)))
    assert dec(b).decode() == 9223372036854775807
    
    b = enc.encode(Val(-9223372036854775808, '-9223372036854775808', Tag(type=ValueType.Int)))
    assert dec(b).decode() == -9223372036854775808


def test_encode_string():
    """Test string encoding matches Go encode_string_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Short string
    b = enc.encode(Val('hello', 'hello', Tag(type=ValueType.String)))
    # String prefix (0x80) | length 5 = 0x85
    assert b[0] == 0x85
    assert b[1:] == b'hello'
    assert dec(b).decode() == 'hello'
    
    # Empty string
    b = enc.encode(Val('', '', Tag(type=ValueType.String)))
    assert b[0] == 0x80  # String | 0
    assert dec(b).decode() == ''
    
    # Longer string requiring extra length byte
    s = 'a' * 30
    b = enc.encode(Val(s, s, Tag(type=ValueType.String)))
    assert b[0] == 0x9E  # String | StringLen1Byte (0x1E = 30)
    assert b[1] == 30
    
    s2 = 'b' * 300
    b = enc.encode(Val(s2, s2, Tag(type=ValueType.String)))
    assert b[0] == 0x9F  # String | StringLen2Byte (0x1F = 31)
    assert (b[1] << 8) | b[2] == 300


def test_encode_float():
    """Test float encoding matches Go encode_float_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Simple short floats: 0.1 to 0.7
    for i in range(1, 8):
        val = i / 10
        b = enc.encode(Val(val, str(val), Tag(type=ValueType.Float64)))
        # PrefixFloat (0x60) | i (without negative flag)
        assert b[0] == 0x60 | i, f"0.{i} should encode as 0x{0x60|i:02x}, got 0x{b[0]:02x}"
        result = dec(b).decode()
        assert abs(result - val) < 1e-10, f"0.{i} roundtrip: expected {val}, got {result}"
    
    # -0.1 to -0.7
    for i in range(1, 8):
        val = -i / 10
        b = enc.encode(Val(val, str(val), Tag(type=ValueType.Float64)))
        # PrefixFloat (0x60) | FloatPositiveNegativeMask (0x10) | i
        assert b[0] == 0x60 | 0x10 | i, f"-0.{i} should encode as 0x{0x60|0x10|i:02x}"
        result = dec(b).decode()
        assert abs(result - val) < 1e-10, f"-0.{i} roundtrip: expected {val}, got {result}"
    
    # Larger floats
    val = 3.14
    b = enc.encode(Val(val, '3.14', Tag(type=ValueType.Float64)))
    result = dec(b).decode()
    assert abs(result - 3.14) < 1e-10


def test_encode_bytes():
    """Test bytes encoding."""
    enc = Encoder()
    dec = Decoder
    
    b = enc.encode(Val(b'\x00\x01\x02', '', Tag(type=ValueType.Bytes)))
    assert b[0] == 0xA3  # PrefixBytes | 3
    assert dec(b).decode() == b'\x00\x01\x02'
    
    # Empty bytes
    b = enc.encode(Val(b'', '', Tag(type=ValueType.Bytes)))
    assert b[0] == 0xA0  # PrefixBytes | 0
    assert dec(b).decode() == b''


def test_encode_object():
    """Test object encoding matches Go encode_object_test.go"""
    enc = Encoder()
    dec = Decoder
    
    obj = Obj(fields=[
        Field(key='name', value=Val('Alice', 'Alice', Tag(type=ValueType.String))),
        Field(key='age', value=Val(30, '30', Tag(type=ValueType.Int))),
    ])
    b = enc.encode(obj)
    result = dec(b).decode()
    assert result == {'name': 'Alice', 'age': 30}
    
    # Empty object
    obj = Obj(fields=[])
    b = enc.encode(obj)
    result = dec(b).decode()
    assert result == {}
    
    # Nested object
    obj = Obj(fields=[
        Field(key='meta', value=Obj(fields=[
            Field(key='count', value=Val(5, '5', Tag(type=ValueType.Int))),
        ])),
    ])
    b = enc.encode(obj)
    result = dec(b).decode()
    assert result == {'meta': {'count': 5}}


def test_encode_array():
    """Test array encoding."""
    enc = Encoder()
    dec = Decoder
    
    arr = Arr(items=[
        Val(1, '1', Tag(type=ValueType.Int)),
        Val(2, '2', Tag(type=ValueType.Int)),
        Val(3, '3', Tag(type=ValueType.Int)),
    ])
    b = enc.encode(arr)
    result = dec(b).decode()
    assert result == [1, 2, 3]
    
    # Empty array
    arr = Arr(items=[])
    b = enc.encode(arr)
    result = dec(b).decode()
    assert result == []
    
    # Nested array
    arr = Arr(items=[
        Arr(items=[Val(1, '1', Tag(type=ValueType.Int))]),
        Arr(items=[Val(2, '2', Tag(type=ValueType.Int))]),
    ])
    b = enc.encode(arr)
    result = dec(b).decode()
    assert result == [[1], [2]]


def test_roundtrip_complex():
    """Complex roundtrip test."""
    enc = Encoder()
    dec = Decoder
    
    obj = Obj(fields=[
        Field(key='id', value=Val(42, '42', Tag(type=ValueType.Int))),
        Field(key='name', value=Val('Bob', 'Bob', Tag(type=ValueType.String))),
        Field(key='active', value=Val(True, 'true', Tag(type=ValueType.Bool))),
        Field(key='score', value=Val(98.5, '98.5', Tag(type=ValueType.Float64))),
        Field(key='tags', value=Arr(items=[
            Val('admin', 'admin', Tag(type=ValueType.String)),
            Val('user', 'user', Tag(type=ValueType.String)),
        ])),
        Field(key='meta', value=Obj(fields=[
            Field(key='created', value=Val('2024-01-15', '2024-01-15', Tag(type=ValueType.String))),
            Field(key='count', value=Val(100, '100', Tag(type=ValueType.Int))),
        ])),
        Field(key='data', value=Val(None, 'null', Tag(type=ValueType.String, is_null=True))),
    ])
    b = enc.encode(obj)
    result = dec(b).decode()
    
    assert result['id'] == 42
    assert result['name'] == 'Bob'
    assert result['active'] == True
    assert abs(result['score'] - 98.5) < 1e-10
    assert result['tags'] == ['admin', 'user']
    assert result['meta']['created'] == '2024-01-15'
    assert result['meta']['count'] == 100
    assert result['data'] is None


def test_null_values():
    """Test null encoding/decoding for different types."""
    enc = Encoder()
    dec = Decoder
    
    for typ, tag in [
        ('bool', Tag(type=ValueType.Bool, is_null=True)),
        ('int', Tag(type=ValueType.Int, is_null=True)),
        ('float', Tag(type=ValueType.Float64, is_null=True)),
        ('string', Tag(type=ValueType.String, is_null=True)),
    ]:
        v = Val(None, 'null', tag)
        b = enc.encode(v)
        result = dec(b).decode()
        assert result is None, f"Null {typ} should decode as None"


if __name__ == '__main__':
    test_encode_bool()
    test_encode_int()
    test_encode_string()
    test_encode_float()
    test_encode_bytes()
    test_encode_object()
    test_encode_array()
    test_roundtrip_complex()
    test_null_values()
    print("All encoder tests passed!")
