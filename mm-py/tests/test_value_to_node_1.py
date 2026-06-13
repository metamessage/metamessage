"""
Tests for value_to_node, node_to_value, and @mm decorator.
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from typing import Any
from metamessage import Encoder, Decoder, decode_to_jsonc
from metamessage.core.value_to_node import (
    value_to_node, node_to_value, encode_from_value, decode_to_value,
    _camel_to_snake
)
from metamessage.core.mm import mm, get_mm_tag_for_class, get_mm_tag_for_field
from metamessage.ir.tag import Tag, ValueType, NewTag

def test_value_to_node_basic():
    """Test basic value to node conversion."""
    # String
    # node = value_to_node("hello")
    # assert node.text == "hello"
    # assert node.tag.type == ValueType.Str

    # Email
    node = value_to_node("abc@example.com", tag=Tag(type=ValueType.Email))
    assert node.text == "abc@example.com"
    assert node.tag.type == ValueType.Email

    # Int
    node = value_to_node(42)
    assert node.data == 42
    assert node.text == "42"
    assert node.tag.type == ValueType.I

    # Float
    node = value_to_node(3.14)
    assert node.data == 3.14
    assert node.tag.type == ValueType.F64

    # Bool
    node = value_to_node(True)
    assert node.data == True
    assert node.text == "true"
    assert node.tag.type == ValueType.Bool

    # None (needs a tag with type)
    node = value_to_node(None)
    print(node.get_type().name)
    assert node.get_type().name == 'Null'

    print("  basic value_to_node OK")


def test_value_to_node_dict():
    """Test dict to node conversion."""
    node = value_to_node({"name": "Alice", "age": 30})
    assert node.get_type().name == 'Object'
    assert len(node.fields) == 2

    fields = {f.key: f.value for f in node.fields}
    assert fields["name"].text == "Alice"
    assert fields["name"].tag.type == ValueType.Str
    assert fields["age"].data == 30

    node = value_to_node({"name": "Alice", "age": 30}, tag=Tag(type=ValueType.Map))
    assert node.get_type().name == 'Object'
    assert node.get_tag().type == ValueType.Map

    print("  dict to node OK")


def test_value_to_node_list():
    """Test list to node conversion."""
    node = value_to_node([1, 2, 3])
    assert node.get_type().name == 'Array'
    assert len(node.items) == 3
    assert node.items[0].data == 1
    assert node.items[1].data == 2
    assert node.items[2].data == 3

    node = value_to_node([1, 2, 3], tag=Tag(type=ValueType.Arr))
    assert node.get_type().name == 'Array'
    assert node.get_tag().type == ValueType.Arr

    print("  list to node OK")


def test_value_to_node_nested():
    """Test nested structures."""
    node = value_to_node({
        "user": {
            "name": "Bob",
            "scores": [85, 90, 78]
        }
    })
    user_field = node.fields[0]
    assert user_field.key == "user"
    assert len(user_field.value.fields) == 2

    scores_field = [f for f in user_field.value.fields if f.key == "scores"][0]
    assert len(scores_field.value.items) == 3
    assert scores_field.value.items[0].data == 85

    print("  nested to node OK")


def test_encode_from_value():
    """Test encoding a Python value directly to binary and decoding back."""
    data = {"name": "Alice", "age": 30, "active": True}
    binary = encode_from_value(data)

    result = decode_to_value(binary)
    assert result["name"] == "Alice"
    assert result["age"] == 30
    assert result["active"] == True
    print("  encode/decode dict OK")


def test_node_to_value():
    """Test node back to Python value."""
    node = value_to_node({"x": 10, "y": 20})
    result = node_to_value(node, dict)
    assert result == {"x": 10, "y": 20}

    result2 = node_to_value(node, Any)
    assert result2 == {"x": 10, "y": 20}

    print("  node_to_value OK")


def test_auto_mm_class():
    """Test auto-apply of @mm + @dataclass for plain class with annotations."""

    class User:
        id: int
        name: str
        age: int

    user = User(id=1, name="Alice", age=20)
    node = value_to_node(user)
    assert len(node.fields) == 3

    # Auto-applied @mm() has default (empty) tag
    cls_tag = get_mm_tag_for_class(User)
    assert cls_tag is not None
    assert cls_tag.desc == ""

    fields = {f.key: f.value for f in node.fields}
    assert fields["id"].data == 1
    assert fields["name"].text == "Alice"
    assert fields["age"].data == 20

    print("  auto-apply @mm class OK")


def test_mm_decorator_with_tag_string():
    """Test @mm with tag string."""
    @mm("desc=用户信息")
    class User:
        id: int

    cls_tag = get_mm_tag_for_class(User)
    assert cls_tag is not None
    assert cls_tag.desc == "用户信息"

    print("  @mm with string OK")


def test_complex_object_generic_list():
    """Test complex object with generic list[T] annotation — auto-infer child_type."""
    @mm(desc="User")
    class User:
        name: str
        age: int
        tags: list[str]  # Should auto-infer child_type=ValueType.Str

    user = User(name="Charlie", age=25, tags=["admin", "user"])
    node = value_to_node(user)

    fields = {f.key: f.value for f in node.fields}
    assert fields["name"].text == "Charlie"
    assert fields["age"].data == 25

    tags_node = fields["tags"]
    assert len(tags_node.items) == 2
    assert tags_node.items[0].text == "admin"
    assert tags_node.items[1].text == "user"
    assert tags_node.tag.child_type == ValueType.Str

    print("  complex object with list[str] OK")


def test_complex_object_explicit_child_type():
    """Test complex object with bare list and explicit child_type — works."""
    @mm(desc="User")
    class User:
        name: str
        age: int
        tags: list = mm(child_type=ValueType.I)

    user = User(name="Charlie", age=25, tags=[1, 2, 3])
    node = value_to_node(user)

    fields = {f.key: f.value for f in node.fields}
    assert fields["name"].text == "Charlie"
    assert fields["age"].data == 25

    tags_node = fields["tags"]
    assert len(tags_node.items) == 3
    assert tags_node.items[0].data == 1
    assert tags_node.items[1].data == 2
    assert tags_node.items[2].data == 3

    print("  complex object with bare list + explicit child_type OK")


def test_bare_list_without_child_type_raises():
    """Test bare list without child_type raises ValueError."""

    @mm(desc="Bad")
    class Bad:
        tags: list  # Missing child_type!

    try:
        value_to_node(Bad(tags=[]))
        assert False, "expected ValueError"
    except ValueError as e:
        assert "must specify child_type" in str(e)

    print("  bare list without child_type raises ValueError OK")


def test_node_to_value_class():
    """Test binding node back to class instance."""
    class Point:
        x: int
        y: int

    node = value_to_node({"x": 10, "y": 20})
    point = node_to_value(node, Point)
    assert point.x == 10
    assert point.y == 20

    print("  node_to_value class OK")


def test_roundtrip_class():
    """Full round-trip: class -> node -> binary -> node -> class."""
    class Person:
        name: str
        age: int

    person = Person(name="Alice", age=30)
    binary = encode_from_value(person)
    result = decode_to_value(binary)
    assert result["name"] == "Alice"
    assert result["age"] == 30

    print("  class roundtrip OK")


def test_mm_field_decorator():
    """Test @mm with field-level annotations (default value syntax)."""
    @mm(desc="User")
    class User:
        id: int = mm(desc="User ID")
        name: str = mm(desc="User name")
        age: int = mm(desc="User age", type=ValueType.U8)

    id_tag = get_mm_tag_for_field(User, "id")
    assert id_tag is not None
    assert id_tag.desc == "User ID"

    name_tag = get_mm_tag_for_field(User, "name")
    assert name_tag is not None
    assert name_tag.desc == "User name"

    age_tag = get_mm_tag_for_field(User, "age")
    assert age_tag is not None
    assert age_tag.desc == "User age"
    assert age_tag.type == ValueType.U8

    user = User(id=42, name="Bob", age=25)
    binary = encode_from_value(user)
    result = decode_to_value(binary)
    assert result["id"] == 42
    assert result["name"] == "Bob"
    assert result["age"] == 25

    print("  mm field decorator OK")


def test_mm_field_with_constraints():
    """Test field-level mm with validation constraints."""
    class Product:
        id: int|None = mm(desc="Product ID", min=1, max=99999)
        name: str = mm(desc="Product name", min=1, max=100)
        price: float = mm(desc="Price", min=0.0, max=999999.99)

    product = Product(id=1001, name="Laptop", price=999.99)
    binary = encode_from_value(product)
    result = decode_to_value(binary)
    assert result["id"] == 1001
    assert result["name"] == "Laptop"
    assert result["price"] == 999.99


    result = decode_to_value(binary, Product)
    print(result.id)

    print(decode_to_jsonc(binary))

    print("  mm field with constraints OK")


def test_mm_field_mixed_with_plain():
    """Test mixing field-level mm with plain fields."""
    class Order:
        id: int = mm(desc="Order ID")
        item: str
        quantity: int = mm(desc="Quantity", min=1)

    order = Order(id=1, item="Widget", quantity=5)
    binary = encode_from_value(order)
    result = decode_to_value(binary)
    assert result["id"] == 1
    assert result["item"] == "Widget"
    assert result["quantity"] == 5

    print("  mm field mixed with plain OK")


if __name__ == '__main__':
    test_value_to_node_basic()
    test_value_to_node_dict()
    test_value_to_node_list()
    test_value_to_node_nested()
    test_encode_from_value()
    test_node_to_value()
    test_auto_mm_class()
    test_mm_decorator_with_tag_string()
    test_complex_object_generic_list()
    test_complex_object_explicit_child_type()
    test_bare_list_without_child_type_raises()
    test_node_to_value_class()
    test_roundtrip_class()
    test_mm_field_decorator()
    test_mm_field_with_constraints()
    test_mm_field_mixed_with_plain()
    # print()
    # print("All value_to_node tests passed!")
