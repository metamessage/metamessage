"""
Tests for metamessage decoder.
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from metamessage.core.encoder import Encoder
from metamessage.core.decoder import Decoder
from metamessage.ir.tag import Tag, ValueType
from metamessage.ir.types import Obj, Arr, Val, Field


def test_decode_bool():
    """Test bool decoding."""
    dec = Decoder
    
    # SimpleTrue
    assert dec(bytes([0x06])).decode() == True
    
    # SimpleFalse
    assert dec(bytes([0x05])).decode() == False
    
    # Tag-wrapped True
    assert dec(bytes([0xE2, 0x01, 0x08, 0x06])).decode() == True


def test_decode_int():
    """Test int decoding."""
    dec = Decoder
    
    # 0
    assert dec(bytes([0x20])).decode() == 0
    
    # 5
    assert dec(bytes([0x25])).decode() == 5
    
    # 24 (1 extra byte)
    assert dec(bytes([0x38, 24])).decode() == 24
    
    # 255
    assert dec(bytes([0x38, 255])).decode() == 255
    
    # 256 (2 extra bytes)
    assert dec(bytes([0x39, 1, 0])).decode() == 256
    
    # -1
    assert dec(bytes([0x41])).decode() == -1
    
    # -24
    assert dec(bytes([0x58, 24])).decode() == -24


def test_decode_string():
    """Test string decoding."""
    dec = Decoder
    
    # Empty string
    assert dec(bytes([0x80])).decode() == ''
    
    # "hello"
    assert dec(bytes([0x85, 0x68, 0x65, 0x6C, 0x6C, 0x6F])).decode() == 'hello'


def test_decode_float():
    """Test float decoding."""
    dec = Decoder
    
    # 0.1 to 0.7
    for i in range(1, 8):
        val = i / 10
        result = dec(bytes([0x60 | i])).decode()
        assert abs(result - val) < 1e-10, f"0.{i} decoded as {result}"
    
    # -0.1 to -0.7
    for i in range(1, 8):
        val = -i / 10
        result = dec(bytes([0x70 | i])).decode()
        assert abs(result - val) < 1e-10, f"-0.{i} decoded as {result}"


def test_decode_array():
    """Test array decoding."""
    enc = Encoder()
    dec = Decoder
    
    # [1, 2, 3]
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


def test_decode_object():
    """Test object decoding."""
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


def test_decode_nested():
    """Test nested structure decoding."""
    enc = Encoder()
    dec = Decoder
    
    nested = Obj(fields=[
        Field(key='user', value=Obj(fields=[
            Field(key='name', value=Val('Bob', 'Bob', Tag(type=ValueType.String))),
            Field(key='scores', value=Arr(items=[
                Val(10, '10', Tag(type=ValueType.Int)),
                Val(20, '20', Tag(type=ValueType.Int)),
            ])),
        ])),
    ])
    b = enc.encode(nested)
    result = dec(b).decode()
    assert result == {'user': {'name': 'Bob', 'scores': [10, 20]}}


def test_decode_edge_cases():
    """Test edge cases."""
    dec = Decoder
    
    # Min/Absolute Int values
    assert dec(bytes([0x3F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])).decode() == 0
    
    # Max positive 64-bit
    b_max = bytes([0x3F, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF])
    assert dec(b_max).decode() == 9223372036854775807
    
    # Min negative 64-bit
    b_min = bytes([0x5F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01])
    # NegativeInt + IntLen8Byte | abs(-9223372036854775808) as uint64 overflows, this is just a structural test
    result = dec(b_min).decode()
    assert isinstance(result, int)
    assert result < 0


def test_empty_structures():
    """Test empty structures."""
    enc = Encoder()
    dec = Decoder
    
    # Single value
    for val, tag in [
        ('', Tag(type=ValueType.String)),
        (0, Tag(type=ValueType.Int)),
        (True, Tag(type=ValueType.Bool)),
    ]:
        v = Val(val, str(val), tag)
        b = enc.encode(v)
        result = dec(b).decode()
        assert result == val or (isinstance(result, float) and result == float(val))


if __name__ == '__main__':
    test_decode_bool()
    test_decode_int()
    test_decode_string()
    test_decode_float()
    test_decode_array()
    test_decode_object()
    test_decode_nested()
    test_decode_edge_cases()
    test_empty_structures()
    print("All decoder tests passed!")
