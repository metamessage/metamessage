"""
Tests for metamessage JSONC parser/generator.
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from metamessage import Encoder, Decoder, parse_jsonc, to_jsonc


def test_parse_basic_json():
    """Test parsing basic JSON."""
    jsonc = '''{
        "name": "Alice",
        "age": 30,
        "active": true,
        "score": 95.5,
        "tags": ["admin", "user"]
    }'''
    node = parse_jsonc(jsonc)
    assert node.get_type().name == 'Object'
    
    # Round-trip through binary
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    
    assert result['name'] == 'Alice'
    assert result['age'] == 30
    assert result['active'] == True
    assert abs(result['score'] - 95.5) < 1e-10
    assert result['tags'] == ['admin', 'user']


def test_parse_null_with_tag():
    """Test parsing null with is_null tag (no bare null values in MM)."""
    jsonc = '''{
        "data": null // mm: is_null
    }'''
    node = parse_jsonc(jsonc)
    assert node.fields[0].value.tag.is_null == True
    assert node.tag.is_null == False  # object should not inherit is_null
    
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    assert result['data'] is None


def test_parse_empty_structures():
    """Test parsing empty objects and arrays."""
    assert parse_jsonc('{}').fields == []
    assert parse_jsonc('[]').items == []


def test_parse_with_mm_tag():
    """Test parsing with mm: comment tags."""
    jsonc = '''{
        // mm: size=10
        "items": [
            {
                "id": 1,
                "label": "test"
            }
        ]
    }'''
    node = parse_jsonc(jsonc)
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    assert result == {'items': [{'id': 1, 'label': 'test'}]}
    assert node.fields[0].value.get_tag().size == 10


def test_parse_nested():
    """Test parsing nested structures."""
    jsonc = '''{
        "meta": {
            "inner": {
                "value": 42
            }
        }
    }'''
    node = parse_jsonc(jsonc)
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    assert result == {'meta': {'inner': {'value': 42}}}


def test_roundtrip_keeps_structure():
    """Test round-trip preserves structure."""
    jsonc = '''{
        "a": 1,
        "b": "hello",
        "c": true,
        "d": null // mm: is_null
    }'''
    node = parse_jsonc(jsonc)
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    assert result['a'] == 1
    assert result['b'] == 'hello'
    assert result['c'] == True
    assert result['d'] is None


def test_to_jsonc():
    """Test JSONC generation."""
    from metamessage.ir.tag import Tag, ValueType
    from metamessage.ir.types import Obj, Arr, Val, Field
    
    obj = Obj(fields=[
        Field(key='name', value=Val('Alice', 'Alice', Tag(type=ValueType.String))),
    ])
    jsonc_out = to_jsonc(obj)
    # Should contain the key
    assert '"name"' in jsonc_out
    assert 'Alice' in jsonc_out
    # Should not have trailing comma issues
    assert ',\n}' not in jsonc_out


def test_parse_complex_jsonc():
    """Test parsing more complex JSONC with various features."""
    jsonc = '''{
        "float_val": -0.5,
        "negative_int": -42,
        "nested_array": [[1, 2], [3, 4]],
        "empty_object": {},
        "empty_array": []
    }'''
    node = parse_jsonc(jsonc)
    enc = Encoder()
    b = enc.encode(node)
    result = Decoder(b).decode()
    assert abs(result['float_val'] - (-0.5)) < 1e-10
    assert result['negative_int'] == -42
    assert result['nested_array'] == [[1, 2], [3, 4]]
    assert result['empty_object'] == {}
    assert result['empty_array'] == []


def test_parse_with_inferred_types():
    """Test that type tags for inferred types (str, int, bool, f64, obj) are omitted in output."""
    from metamessage.ir.types import Obj, Field
    jsonc = '''{
        "id": 123,
        "name": "test",
        "active": true,
        "score": 1.5,
        "tags": ["a", "b"]
    }'''
    node = parse_jsonc(jsonc)
    jsonc_out = to_jsonc(node)
    # Inferred types should not have explicit type= tags
    assert 'type=str' not in jsonc_out
    assert 'type=int' not in jsonc_out
    assert 'type=bool' not in jsonc_out
    assert 'type=f64' not in jsonc_out
    assert 'type=obj' not in jsonc_out


def test_null_output_has_default():
    """Test that null values output as default values with is_null comment."""
    from metamessage.ir.tag import Tag, ValueType
    from metamessage.ir.types import Val, Obj, Field
    
    obj = Obj(fields=[
        Field(key='s', value=Val(None, 'null', Tag(type=ValueType.String, is_null=True))),
        Field(key='n', value=Val(None, 'null', Tag(type=ValueType.Int, is_null=True))),
        Field(key='b', value=Val(None, 'null', Tag(type=ValueType.Bool, is_null=True))),
    ])
    jsonc_out = to_jsonc(obj)
    print(f'Null output:\n{jsonc_out}')
    assert 'is_null' in jsonc_out
    # Should show defaults instead of null literal
    assert '""' in jsonc_out  # string default
    assert '": 0' in jsonc_out or '":0' in jsonc_out  # int default
    assert '"b": false' in jsonc_out or '"b":false' in jsonc_out  # bool default


if __name__ == '__main__':
    test_parse_basic_json()
    test_parse_null_with_tag()
    test_parse_empty_structures()
    test_parse_with_mm_tag()
    test_parse_nested()
    test_roundtrip_keeps_structure()
    test_to_jsonc()
    test_parse_complex_jsonc()
    test_parse_with_inferred_types()
    test_null_output_has_default()
    print("All JSONC tests passed!")
