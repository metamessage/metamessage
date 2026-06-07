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
from metamessage.ir.ast import NodeObject, Arr, NodeScalar, Field


def test_encode_bool():
    """Test bool encoding matches Go encode_bool_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # True
    b = enc.encode(NodeScalar(True, 'true', Tag(type=ValueType.Bool)))
    assert len(b) == 1
    assert b[0] == 0x06  # SimpleTrue
    assert dec(b).decode() == True
    
    # False
    b = enc.encode(NodeScalar(False, 'false', Tag(type=ValueType.Bool)))
    assert len(b) == 1
    assert b[0] == 0x05  # SimpleFalse
    assert dec(b).decode() == False


def test_encode_i():
    """Test int encoding matches Go encode_int_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Small positive int (fits in 5 bits)
    b = enc.encode(NodeScalar(0, '0', Tag(type=ValueType.I)))
    assert b[0] == 0x20  # PositiveInt | 0
    
    b = enc.encode(NodeScalar(5, '5', Tag(type=ValueType.I)))
    assert b[0] == 0x25  # PositiveInt | 5
    
    b = enc.encode(NodeScalar(23, '23', Tag(type=ValueType.I)))
    assert b[0] == 0x37  # PositiveInt | 23
    
    # Medium int (1 extra byte)
    b = enc.encode(NodeScalar(24, '24', Tag(type=ValueType.I)))
    assert b[0] == 0x38  # PositiveInt | IntLen1Byte (24)
    assert b[1] == 24
    
    b = enc.encode(NodeScalar(255, '255', Tag(type=ValueType.I)))
    assert b[0] == 0x38
    assert b[1] == 255
    
    # Large int (2 extra bytes)
    b = enc.encode(NodeScalar(256, '256', Tag(type=ValueType.I)))
    assert b[0] == 0x39  # PositiveInt | IntLen2Byte (25)
    assert len(b) == 3
    
    # Negative int
    b = enc.encode(NodeScalar(-1, '-1', Tag(type=ValueType.I)))
    assert b[0] & 0xE0 == 0x40  # NegativeInt
    assert dec(b).decode() == -1
    
    b = enc.encode(NodeScalar(-24, '-24', Tag(type=ValueType.I)))
    assert b[0] & 0xE0 == 0x40
    assert dec(b).decode() == -24
    
    # Edge cases
    b = enc.encode(NodeScalar(9223372036854775807, '9223372036854775807', Tag(type=ValueType.I)))
    assert dec(b).decode() == 9223372036854775807
    
    b = enc.encode(NodeScalar(-9223372036854775808, '-9223372036854775808', Tag(type=ValueType.I)))
    assert dec(b).decode() == -9223372036854775808


def test_encode_string():
    """Test string encoding matches Go encode_string_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Short string
    b = enc.encode(NodeScalar('hello', 'hello', Tag(type=ValueType.Str)))
    # String prefix (0x80) | length 5 = 0x85
    assert b[0] == 0x85
    assert b[1:] == b'hello'
    assert dec(b).decode() == 'hello'
    
    # Empty string
    b = enc.encode(NodeScalar('', '', Tag(type=ValueType.Str)))
    assert b[0] == 0x80  # String | 0
    assert dec(b).decode() == ''
    
    # Longer string requiring extra length byte
    s = 'a' * 30
    b = enc.encode(NodeScalar(s, s, Tag(type=ValueType.Str)))
    assert b[0] == 0x9E  # String | StringLen1Byte (0x1E = 30)
    assert b[1] == 30
    
    s2 = 'b' * 300
    b = enc.encode(NodeScalar(s2, s2, Tag(type=ValueType.Str)))
    assert b[0] == 0x9F  # String | StringLen2Byte (0x1F = 31)
    assert (b[1] << 8) | b[2] == 300


def test_encode_float():
    """Test float encoding matches Go encode_float_test.go"""
    enc = Encoder()
    dec = Decoder
    
    # Simple short floats: 0.1 to 0.7
    for i in range(1, 8):
        val = i / 10
        b = enc.encode(NodeScalar(val, str(val), Tag(type=ValueType.F64)))
        # PrefixFloat (0x60) | i (without negative flag)
        assert b[0] == 0x60 | i, f"0.{i} should encode as 0x{0x60|i:02x}, got 0x{b[0]:02x}"
        result = dec(b).decode()
        assert abs(result - val) < 1e-10, f"0.{i} roundtrip: expected {val}, got {result}"
    
    # -0.1 to -0.7
    for i in range(1, 8):
        val = -i / 10
        b = enc.encode(NodeScalar(val, str(val), Tag(type=ValueType.F64)))
        # PrefixFloat (0x60) | FloatPositiveNegativeMask (0x10) | i
        assert b[0] == 0x60 | 0x10 | i, f"-0.{i} should encode as 0x{0x60|0x10|i:02x}"
        result = dec(b).decode()
        assert abs(result - val) < 1e-10, f"-0.{i} roundtrip: expected {val}, got {result}"
    
    # Larger floats
    val = 3.14
    b = enc.encode(NodeScalar(val, '3.14', Tag(type=ValueType.F64)))
    result = dec(b).decode()
    assert abs(result - 3.14) < 1e-10


def test_encode_bytes():
    """Test bytes encoding."""
    enc = Encoder()
    dec = Decoder
    
    b = enc.encode(NodeScalar(b'\x00\x01\x02', '', Tag(type=ValueType.Bytes)))
    assert b[0] == 0xA3  # PrefixBytes | 3
    assert dec(b).decode() == b'\x00\x01\x02'
    
    # Empty bytes
    b = enc.encode(NodeScalar(b'', '', Tag(type=ValueType.Bytes)))
    assert b[0] == 0xA0  # PrefixBytes | 0
    assert dec(b).decode() == b''


def test_encode_object():
    """Test object encoding matches Go encode_object_test.go"""
    enc = Encoder()
    dec = Decoder
    
    obj = NodeObject(fields=[
        Field(key='name', value=NodeScalar('Alice', 'Alice', Tag(type=ValueType.Str))),
        Field(key='age', value=NodeScalar(30, '30', Tag(type=ValueType.I))),
    ])
    b = enc.encode(obj)
    result = dec(b).decode()
    assert result == {'name': 'Alice', 'age': 30}
    
    # Empty object
    obj = NodeObject(fields=[])
    b = enc.encode(obj)
    result = dec(b).decode()
    assert result == {}
    
    # Nested object
    obj = NodeObject(fields=[
        Field(key='meta', value=NodeObject(fields=[
            Field(key='count', value=NodeScalar(5, '5', Tag(type=ValueType.I))),
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
        NodeScalar(1, '1', Tag(type=ValueType.I)),
        NodeScalar(2, '2', Tag(type=ValueType.I)),
        NodeScalar(3, '3', Tag(type=ValueType.I)),
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
        Arr(items=[NodeScalar(1, '1', Tag(type=ValueType.I))]),
        Arr(items=[NodeScalar(2, '2', Tag(type=ValueType.I))]),
    ])
    b = enc.encode(arr)
    result = dec(b).decode()
    assert result == [[1], [2]]


def test_roundtrip_complex():
    """Complex roundtrip test."""
    enc = Encoder()
    dec = Decoder
    
    obj = NodeObject(fields=[
        Field(key='id', value=NodeScalar(42, '42', Tag(type=ValueType.I))),
        Field(key='name', value=NodeScalar('Bob', 'Bob', Tag(type=ValueType.Str))),
        Field(key='active', value=NodeScalar(True, 'true', Tag(type=ValueType.Bool))),
        Field(key='score', value=NodeScalar(98.5, '98.5', Tag(type=ValueType.F64))),
        Field(key='tags', value=Arr(items=[
            NodeScalar('admin', 'admin', Tag(type=ValueType.Str)),
            NodeScalar('user', 'user', Tag(type=ValueType.Str)),
        ])),
        Field(key='meta', value=NodeObject(fields=[
            Field(key='created', value=NodeScalar('2024-01-15', '2024-01-15', Tag(type=ValueType.Str))),
            Field(key='count', value=NodeScalar(100, '100', Tag(type=ValueType.I))),
        ])),
        Field(key='data', value=NodeScalar(None, 'null', Tag(type=ValueType.Str, is_null=True))),
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
    assert result['data'] == ''


def test_null_values():
    """Test null encoding/decoding for different types."""
    enc = Encoder()
    dec = Decoder
    
    expected = {
        'bool': False,
        'int': 0,
        'float': 0.0,
        'string': '',
    }
    for typ, tag in [
        ('bool', Tag(type=ValueType.Bool, is_null=True)),
        ('int', Tag(type=ValueType.I, is_null=True)),
        ('float', Tag(type=ValueType.F64, is_null=True)),
        ('string', Tag(type=ValueType.Str, is_null=True)),
    ]:
        v = NodeScalar(None, 'null', tag)
        b = enc.encode(v)
        result = dec(b).decode()
        assert result == expected[typ], f"Null {typ} should decode as {expected[typ]}"


if __name__ == '__main__':
    test_encode_bool()
    test_encode_i()
    test_encode_string()
    test_encode_float()
    test_encode_bytes()
    test_encode_object()
    test_encode_array()
    test_roundtrip_complex()
    test_null_values()
    print("All encoder tests passed!")
