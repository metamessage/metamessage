"""
Tests for the @mm tag decorator (with validation).
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from metamessage.core.mm import (
    mm, get_mm_tag_for_class, get_mm_tag_for_field,
    _MM_CLASS_REGISTRY, _MM_FIELD_REGISTRY,
)
from metamessage.ir.tag import Tag, ValueType, NewTag


def setup_method():
    """Clear registries before each test."""
    _MM_CLASS_REGISTRY.clear()
    _MM_FIELD_REGISTRY.clear()


# ============================================================================
# Class decorator tests (Obj / Map only)
# ============================================================================

def test_class_decorator_obj():
    """@mm with keyword args as class decorator (Obj, default)."""
    setup_method()

    @mm(desc="User information", nullable=True, deprecated=True)
    class User:
        pass

    tag = get_mm_tag_for_class(User)
    assert tag is not None
    assert tag.desc == "User information"
    assert tag.type == ValueType.Obj
    assert tag.nullable is True
    assert tag.deprecated is True
    print("  class decorator Obj OK")


def test_class_decorator_map():
    """@mm with Map type as class decorator."""
    setup_method()

    @mm(type=ValueType.Map, desc="Config map")
    class Config:
        pass

    tag = get_mm_tag_for_class(Config)
    assert tag is not None
    assert tag.desc == "Config map"
    assert tag.type == ValueType.Map
    print("  class decorator Map OK")


def test_class_decorator_empty():
    """@mm() with no args defaults to Obj."""
    setup_method()

    @mm()
    class Empty:
        pass

    tag = get_mm_tag_for_class(Empty)
    assert tag is not None
    assert tag.type == ValueType.Obj
    print("  class decorator empty OK")


def test_class_decorator_string():
    """@mm with tag string as class decorator."""
    setup_method()

    @mm("desc=User info; type=obj")
    class User:
        pass

    tag = get_mm_tag_for_class(User)
    assert tag is not None
    assert tag.desc == "User info"
    print("  class decorator string OK")


# ============================================================================
# Field decorator tests
# ============================================================================

def test_field_decorator():
    """@mm as field-level decorator with type annotations."""
    setup_method()

    @mm(desc="User")
    class User:
        id: int = mm(desc="User ID", type=ValueType.I64)
        name: str = mm(desc="User name")

    id_tag = get_mm_tag_for_field(User, "id")
    assert id_tag is not None
    assert id_tag.desc == "User ID"
    assert id_tag.type == ValueType.I64

    name_tag = get_mm_tag_for_field(User, "name")
    assert name_tag is not None
    assert name_tag.desc == "User name"
    assert name_tag.type == ValueType.Str  # auto-inferred from str

    print("  field decorator OK")


def test_field_auto_infer_type():
    """Field type auto-inferred from annotation."""
    setup_method()

    @mm(desc="Test")
    class Test:
        count: int = mm(desc="Count")
        price: float = mm(desc="Price")
        flag: bool = mm(desc="Flag")
        data: bytes = mm(desc="Data")
        items: list[int] = mm(desc="Items")

    assert get_mm_tag_for_field(Test, "count").type == ValueType.I64
    assert get_mm_tag_for_field(Test, "price").type == ValueType.F64
    assert get_mm_tag_for_field(Test, "flag").type == ValueType.Bool
    assert get_mm_tag_for_field(Test, "data").type == ValueType.Bytes
    assert get_mm_tag_for_field(Test, "items").type == ValueType.Vec
    print("  field auto-infer type OK")


def test_field_explicit_type_compatible():
    """Explicit type that is compatible with annotation."""
    setup_method()

    @mm(desc="Test")
    class Test:
        age: int = mm(type=ValueType.U8, desc="Age")
        score: float = mm(type=ValueType.F32, desc="Score")
        arr: list[int] = mm(type=ValueType.Arr, desc="Fixed array")

    assert get_mm_tag_for_field(Test, "age").type == ValueType.U8
    assert get_mm_tag_for_field(Test, "score").type == ValueType.F32
    assert get_mm_tag_for_field(Test, "arr").type == ValueType.Arr
    print("  field explicit compatible type OK")


def test_field_optional_is_null():
    """is_null allowed with Optional type annotation."""
    setup_method()

    @mm(desc="Test")
    class Test:
        name: str | None = mm(is_null=True, desc="Optional name")
        age: int | None = mm(is_null=True, desc="Optional age")

    assert get_mm_tag_for_field(Test, "name").is_null is True
    assert get_mm_tag_for_field(Test, "age").is_null is True
    print("  field optional is_null OK")


def test_field_tags_allowed():
    """Allowed tags for specific types work correctly."""
    setup_method()

    @mm(desc="Test")
    class Test:
        # Str: pattern, min, max, size, allow_empty are allowed
        email: str = mm(desc="Email", pattern=r"^.+@.+$", min="5", max="100")
        # Bool: only base tags allowed
        active: bool = mm(desc="Active", nullable=True)
        # Vec: child_* tags allowed
        tags: list[str] = mm(desc="Tags", child_type=ValueType.Str, child_desc="Tag value")

    email_tag = get_mm_tag_for_field(Test, "email")
    assert email_tag.pattern == r"^.+@.+$"
    assert email_tag.min == "5"
    assert email_tag.max == "100"

    active_tag = get_mm_tag_for_field(Test, "active")
    assert active_tag.nullable is True

    tags_tag = get_mm_tag_for_field(Test, "tags")
    assert tags_tag.child_desc == "Tag value"
    assert tags_tag.child_type == ValueType.Str
    print("  field tags allowed OK")


# ============================================================================
# get_tag tests
# ============================================================================

def test_get_tag():
    """mm.get_tag() returns the Tag object."""
    t = mm(desc="hello").get_tag()
    assert t.desc == "hello"

    t2 = mm("desc=world").get_tag()
    assert t2.desc == "world"

    print("  get_tag OK")


# ============================================================================
# Not found tests
# ============================================================================

def test_get_mm_tag_for_field_not_found():
    """get_mm_tag_for_field returns None for unknown fields."""
    setup_method()

    @mm(desc="Test")
    class A:
        pass

    assert get_mm_tag_for_field(A, "nonexistent") is None
    print("  field not found OK")


def test_get_mm_tag_for_class_not_found():
    """get_mm_tag_for_class returns None for unknown classes."""
    class A:
        pass

    assert get_mm_tag_for_class(A) is None
    print("  class not found OK")


# ============================================================================
# Validation error tests
# ============================================================================

def test_error_class_with_field_type():
    """Class decorator rejects field-level types."""
    setup_method()
    try:
        @mm(type=ValueType.I64, desc="test")
        class Bad:
            pass
        assert False, "should have raised ValueError"
    except ValueError as e:
        assert "只允许 Obj 或 Map" in str(e)
    print("  error: class with field type OK")


def test_error_is_null_without_optional():
    """is_null=True rejected when type is not Optional."""
    setup_method()
    try:
        @mm(desc="Test")
        class Bad:
            name: str = mm(is_null=True, desc="name")
        assert False, "should have raised ValueError"
    except (ValueError, RuntimeError) as e:
        msg = str(e.__cause__) if isinstance(e, RuntimeError) and e.__cause__ else str(e)
        assert "不是 Optional" in msg, f"unexpected: {msg}"
    print("  error: is_null without Optional OK")


def test_error_wrong_tag_for_type():
    """Tag not allowed for the ValueType is rejected."""
    setup_method()
    try:
        @mm(desc="Test")
        class Bad:
            flag: bool = mm(min="1", desc="flag")
        assert False, "should have raised ValueError"
    except (ValueError, RuntimeError) as e:
        msg = str(e.__cause__) if isinstance(e, RuntimeError) and e.__cause__ else str(e)
        assert "不允许使用标签" in msg and "min" in msg, f"unexpected: {msg}"
    print("  error: wrong tag for type OK")


def test_error_type_mismatch():
    """Type annotation doesn't match explicit type."""
    setup_method()
    try:
        @mm(desc="Test")
        class Bad:
            name: str = mm(type=ValueType.I64, desc="name")
        assert False, "should have raised ValueError"
    except (ValueError, RuntimeError) as e:
        msg = str(e.__cause__) if isinstance(e, RuntimeError) and e.__cause__ else str(e)
        assert "类型不匹配" in msg, f"unexpected: {msg}"
    print("  error: type mismatch OK")


def test_error_class_level_on_field():
    """ValueType.Obj not allowed on a primitive field."""
    setup_method()
    try:
        @mm(desc="Test")
        class Bad:
            name: str = mm(type=ValueType.Obj, desc="name")
        assert False, "should have raised ValueError"
    except (ValueError, RuntimeError) as e:
        msg = str(e.__cause__) if isinstance(e, RuntimeError) and e.__cause__ else str(e)
        assert "只能用于类装饰器" in msg, f"unexpected: {msg}"
    print("  error: class-level type on field OK")


if __name__ == '__main__':
    test_class_decorator_obj()
    test_class_decorator_map()
    test_class_decorator_empty()
    test_class_decorator_string()
    test_field_decorator()
    test_field_auto_infer_type()
    test_field_explicit_type_compatible()
    test_field_optional_is_null()
    test_field_tags_allowed()
    test_get_tag()
    test_get_mm_tag_for_field_not_found()
    test_get_mm_tag_for_class_not_found()
    test_error_class_with_field_type()
    test_error_is_null_without_optional()
    test_error_wrong_tag_for_type()
    test_error_type_mismatch()
    test_error_class_level_on_field()
    print()
    print("All mm tests passed!")