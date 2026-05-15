"""
Tests for value_to_node, node_to_value, and @mm decorator.
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from typing import Any
from metamessage import Encoder, Decoder
from metamessage.core.value_to_node import (
    value_to_node, node_to_value, encode_from_value, decode_to_value,
    mm, get_mm_tag_for_class, get_mm_tag_for_field,
    _camel_to_snake
)
from metamessage.ir.tag import Tag, ValueType, NewTag


def test_camel_to_snake():
    """Test CamelCase to snake_case conversion."""
    assert _camel_to_snake("User") == "user"
    assert _camel_to_snake("UserID") == "user_id"
    assert _camel_to_snake("HTTPRequest") == "http_request"
    assert _camel_to_snake("") == ""
    print("  camel_to_snake OK")


def test_value_to_node_basic():
    """Test basic value to node conversion."""
    # String
    node = value_to_node("hello")
    assert node.text == "hello"
    assert node.tag.type == ValueType.String

    # Int
    node = value_to_node(42)
    assert node.data == 42
    assert node.text == "42"
    assert node.tag.type == ValueType.Int

    # Float
    node = value_to_node(3.14)
    assert node.data == 3.14
    assert node.tag.type == ValueType.Float64

    # Bool
    node = value_to_node(True)
    assert node.data == True
    assert node.text == "true"
    assert node.tag.type == ValueType.Bool

    # None (needs a tag with type)
    tag = NewTag()
    tag.type = ValueType.String
    node = value_to_node(None, tag)
    assert node.tag.is_null == True
    assert node.data is None

    print("  basic value_to_node OK")


def test_value_to_node_dict():
    """Test dict to node conversion."""
    node = value_to_node({"name": "Alice", "age": 30})
    assert node.get_type().name == 'Object'
    assert len(node.fields) == 2

    fields = {f.key: f.value for f in node.fields}
    assert fields["name"].text == "Alice"
    assert fields["name"].tag.type == ValueType.String
    assert fields["age"].data == 30

    print("  dict to node OK")


def test_value_to_node_list():
    """Test list to node conversion."""
    node = value_to_node([1, 2, 3])
    assert node.get_type().name == 'Array'
    assert len(node.items) == 3
    assert node.items[0].data == 1
    assert node.items[1].data == 2
    assert node.items[2].data == 3

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


def test_mm_decorator_class():
    """Test @mm as class decorator."""
    @mm(desc="用户信息")
    class User:
        def __init__(self, id: int = 0, name: str = "", age: int = 0):
            self.id = id
            self.name = name
            self.age = age

    cls_tag = get_mm_tag_for_class(User)
    assert cls_tag is not None
    assert cls_tag.desc == "用户信息"

    user = User(id=1, name="Alice", age=20)
    node = value_to_node(user)
    assert len(node.fields) == 3

    fields = {f.key: f.value for f in node.fields}
    assert fields["id"].data == 1
    assert fields["name"].text == "Alice"
    assert fields["age"].data == 20

    print("  @mm class decorator OK")


def test_mm_decorator_with_tag_string():
    """Test @mm with tag string."""
    @mm("desc=用户信息")
    class User:
        def __init__(self, id: int = 0):
            self.id = id

    cls_tag = get_mm_tag_for_class(User)
    assert cls_tag is not None
    assert cls_tag.desc == "用户信息"

    print("  @mm with string OK")


def test_complex_object():
    """Test complex object with nested objects."""
    @mm(desc="User")
    class User:
        def __init__(self, name: str = "", age: int = 0, tags: list = None):
            self.name = name
            self.age = age
            self.tags = tags or []

    user = User(name="Charlie", age=25, tags=["admin", "user"])
    node = value_to_node(user)

    fields = {f.key: f.value for f in node.fields}
    assert fields["name"].text == "Charlie"
    assert fields["age"].data == 25

    tags_node = fields["tags"]
    assert len(tags_node.items) == 2
    assert tags_node.items[0].text == "admin"
    assert tags_node.items[1].text == "user"

    print("  complex object OK")


def test_node_to_value_class():
    """Test binding node back to class instance."""
    class Point:
        def __init__(self, x: int = 0, y: int = 0):
            self.x = x
            self.y = y

    node = value_to_node({"x": 10, "y": 20})
    point = node_to_value(node, Point)
    assert point.x == 10
    assert point.y == 20

    print("  node_to_value class OK")


def test_roundtrip_class():
    """Full round-trip: class -> node -> binary -> node -> class."""
    class Person:
        def __init__(self, name: str = "", age: int = 0):
            self.name = name
            self.age = age

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
        age: int = mm(desc="User age", type=ValueType.Uint8)

        def __init__(self, id: int = 0, name: str = "", age: int = 0):
            self.id = id
            self.name = name
            self.age = age

    id_tag = get_mm_tag_for_field(User, "id")
    assert id_tag is not None
    assert id_tag.desc == "User ID"

    name_tag = get_mm_tag_for_field(User, "name")
    assert name_tag is not None
    assert name_tag.desc == "User name"

    age_tag = get_mm_tag_for_field(User, "age")
    assert age_tag is not None
    assert age_tag.desc == "User age"
    assert age_tag.type == ValueType.Uint8

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

        def __init__(self, id: int = 0, name: str = "", price: float = 0.0):
            self.id = id
            self.name = name
            self.price = price

    product = Product(id=1001, name="Laptop", price=999.99)
    binary = encode_from_value(product)
    result = decode_to_value(binary)
    assert result["id"] == 1001
    assert result["name"] == "Laptop"
    assert result["price"] == 999.99

    print("  mm field with constraints OK")


def test_mm_field_mixed_with_plain():
    """Test mixing field-level mm with plain fields."""
    class Order:
        id: int = mm(desc="Order ID")
        item: str
        quantity: int = mm(desc="Quantity", min=1)

        def __init__(self, id: int = 0, item: str = "", quantity: int = 0):
            self.id = id
            self.item = item
            self.quantity = quantity

    order = Order(id=1, item="Widget", quantity=5)
    binary = encode_from_value(order)
    result = decode_to_value(binary)
    assert result["id"] == 1
    assert result["item"] == "Widget"
    assert result["quantity"] == 5

    print("  mm field mixed with plain OK")


if __name__ == '__main__':
    test_camel_to_snake()
    test_value_to_node_basic()
    test_value_to_node_dict()
    test_value_to_node_list()
    test_value_to_node_nested()
    test_encode_from_value()
    test_node_to_value()
    test_mm_decorator_class()
    test_mm_decorator_with_tag_string()
    test_complex_object()
    test_node_to_value_class()
    test_roundtrip_class()
    test_mm_field_decorator()
    test_mm_field_with_constraints()
    test_mm_field_mixed_with_plain()
    print()
    print("All value_to_node tests passed!")
