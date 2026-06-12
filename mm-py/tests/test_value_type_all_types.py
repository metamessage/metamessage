"""
Test that all ValueType members work correctly with @mm decorator (with validation).
"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))
import datetime as dt
from uuid import UUID
from decimal import Decimal
from enum import Enum

from metamessage.core.mm import mm, get_mm_tag_for_class, get_mm_tag_for_field, _MM_CLASS_REGISTRY, _MM_FIELD_REGISTRY
from metamessage.ir.value_type import ValueType


def _clear():
    _MM_CLASS_REGISTRY.clear()
    _MM_FIELD_REGISTRY.clear()


def test_class_level_types():
    """Obj and Map as class decorators."""
    _clear()

    @mm(type=ValueType.Obj, desc="test")
    class A:
        pass
    assert get_mm_tag_for_class(A).type == ValueType.Obj

    @mm(type=ValueType.Map, desc="test")
    class B:
        pass
    assert get_mm_tag_for_class(B).type == ValueType.Map

    print("  class-level: Obj, Map OK")


def test_field_int_types():
    """All int-compatible ValueTypes with int annotation."""
    for vt in [ValueType.I, ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64,
               ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64,
               ValueType.Bigint]:
        _clear()

        @mm(desc="Test")
        class T:
            x: int = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: all int types OK")


def test_field_float_types():
    """F32 and F64 with float annotation."""
    for vt in [ValueType.F32, ValueType.F64]:
        _clear()

        @mm(desc="Test")
        class T:
            x: float = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: float types OK")


def test_field_list_types():
    """Vec and Arr with list annotation."""
    for vt in [ValueType.Vec, ValueType.Arr]:
        _clear()

        @mm(desc="Test")
        class T:
            x: list[int] = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: list types OK")


def test_field_string_types():
    """Str, Ip, Url, Email with str annotation."""
    for vt in [ValueType.Str, ValueType.Ip, ValueType.Url, ValueType.Email]:
        _clear()

        @mm(desc="Test")
        class T:
            x: str = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: string types OK")


def test_field_bool():
    """Bool with bool annotation."""
    _clear()

    @mm(desc="Test")
    class T:
        x: bool = mm(type=ValueType.Bool, desc="x")

    assert get_mm_tag_for_field(T, "x").type == ValueType.Bool
    print("  field: Bool OK")


def test_field_bytes_types():
    """Bytes and Media with bytes annotation."""
    for vt in [ValueType.Bytes, ValueType.Media]:
        _clear()

        @mm(desc="Test")
        class T:
            x: bytes = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: bytes/media OK")


def test_field_datetime_types():
    """Datetime, Date, Time with their types."""
    for vt, py_type in [
        (ValueType.Datetime, dt.datetime),
        (ValueType.Date, dt.date),
        (ValueType.Time, dt.time),
    ]:
        _clear()

        @mm(desc="Test")
        class T:
            x: py_type = mm(type=vt, desc="x")

        assert get_mm_tag_for_field(T, "x").type == vt

    print("  field: datetime/date/time OK")


def test_field_uuid():
    """Uuid with UUID annotation."""
    _clear()

    @mm(desc="Test")
    class T:
        x: UUID = mm(type=ValueType.Uuid, desc="x")

    assert get_mm_tag_for_field(T, "x").type == ValueType.Uuid
    print("  field: UUID OK")


def test_field_decimal():
    """Decimal with Decimal annotation."""
    _clear()

    @mm(desc="Test")
    class T:
        x: Decimal = mm(type=ValueType.Decimal, desc="x")

    assert get_mm_tag_for_field(T, "x").type == ValueType.Decimal
    print("  field: Decimal OK")


def test_field_enums():
    """Enums with Enum subclass annotation."""
    _clear()

    class Color(Enum):
        RED = 1
        GREEN = 2

    _clear()

    @mm(desc="Test")
    class T:
        x: Color = mm(type=ValueType.Enums, desc="x")

    assert get_mm_tag_for_field(T, "x").type == ValueType.Enums
    print("  field: Enums OK")


def test_value_type_parse_roundtrip():
    """ValueType -> str -> parse -> same type."""
    from metamessage.ir.value_type import parse_value_type
    for vt in ValueType:
        s = str(vt)
        parsed = parse_value_type(s)
        assert parsed == vt, f"{vt}: str={s!r}, parsed={parsed}"
    print("  str/parse roundtrip OK")


if __name__ == '__main__':
    test_class_level_types()
    test_field_int_types()
    test_field_float_types()
    test_field_list_types()
    test_field_string_types()
    test_field_bool()
    test_field_bytes_types()
    test_field_datetime_types()
    test_field_uuid()
    test_field_decimal()
    test_field_enums()
    test_value_type_parse_roundtrip()
    print()
    print("All value_type tests passed!")