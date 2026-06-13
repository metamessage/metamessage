"""
mm: MetaMessage tag decorator for Python classes and fields.

Usage:
    @mm(desc="User information")
    class User:
        id: int
        name: str

    class User:
        @mm(desc="User ID")
        id: int
"""
from __future__ import annotations

import enum
import types
import typing
from dataclasses import dataclass, field, fields as dc_fields
from datetime import datetime, date, time
from decimal import Decimal
from typing import Any, Dict, Optional, Type
from uuid import UUID

from ..ir.tag import Tag, ValueType, NewTag, mm_tag


# ---------------------------------------------------------------------------
# Tag fields allowed per ValueType (based on the MetaMessage type spec)
# ---------------------------------------------------------------------------
_BASE_FIELD_TAGS = {"is_null", "example", "deprecated", "name", "desc", "nullable", "allow_empty"}
_NUMERIC_TAGS = _BASE_FIELD_TAGS | {"min", "max", "size"}
_STRING_TAGS = _BASE_FIELD_TAGS | {"min", "max", "size", "pattern"}
_CONTAINER_TAGS = _BASE_FIELD_TAGS | {"type", "min", "max", "size", "unique"}
_VEC_ARR_TAGS = _CONTAINER_TAGS | {
    "child_desc", "child_type", "child_nullable", "child_allow_empty",
    "child_unique", "child_default_val", "child_min", "child_max",
    "child_size", "child_enums", "child_pattern", "child_location",
    "child_version", "child_mime",
}
_TIME_TAGS = _BASE_FIELD_TAGS | {"min", "max", "location"}
_UUID_TAGS = _STRING_TAGS | {"version"}
_DECIMAL_TAGS = _STRING_TAGS
_IP_TAGS = _STRING_TAGS | {"version"}

ALLOWED_TAGS: Dict[ValueType, set] = {
    # --- class-level ---
    ValueType.Obj: _BASE_FIELD_TAGS,
    ValueType.Map: _BASE_FIELD_TAGS | {"type"},

    # --- field-level ---
    ValueType.Vec: _VEC_ARR_TAGS,
    ValueType.Arr: _VEC_ARR_TAGS,
    ValueType.Str: _STRING_TAGS,
    ValueType.Bytes: _BASE_FIELD_TAGS | {"min", "max", "size"},
    ValueType.Bool: _BASE_FIELD_TAGS,
    ValueType.I: _NUMERIC_TAGS,
    ValueType.I8: _NUMERIC_TAGS,
    ValueType.I16: _NUMERIC_TAGS,
    ValueType.I32: _NUMERIC_TAGS,
    ValueType.I64: _NUMERIC_TAGS,
    ValueType.U: _NUMERIC_TAGS,
    ValueType.U8: _NUMERIC_TAGS,
    ValueType.U16: _NUMERIC_TAGS,
    ValueType.U32: _NUMERIC_TAGS,
    ValueType.U64: _NUMERIC_TAGS,
    ValueType.Bigint: _NUMERIC_TAGS,
    ValueType.F32: _NUMERIC_TAGS,
    ValueType.F64: _NUMERIC_TAGS,
    ValueType.Datetime: _TIME_TAGS,
    ValueType.Date: _TIME_TAGS,
    ValueType.Time: _TIME_TAGS,
    ValueType.Uuid: _UUID_TAGS,
    ValueType.Decimal: _DECIMAL_TAGS,
    ValueType.Ip: _IP_TAGS,
    ValueType.Url: _STRING_TAGS,
    ValueType.Email: _STRING_TAGS,
    ValueType.Enums: _BASE_FIELD_TAGS | {"enums"},
    ValueType.Media: _BASE_FIELD_TAGS | {"min", "max", "size", "mime"},
}

# ValueTypes that are valid at class level
_CLASS_LEVEL_TYPES = {ValueType.Obj, ValueType.Map}

# Python type → default ValueType
_PYTHON_TYPE_TO_VT: Dict[type, ValueType] = {
    str: ValueType.Str,
    int: ValueType.I64,
    float: ValueType.F64,
    bool: ValueType.Bool,
    bytes: ValueType.Bytes,
    list: ValueType.Vec,
    dict: ValueType.Map,
    datetime: ValueType.Datetime,
    date: ValueType.Date,
    time: ValueType.Time,
    UUID: ValueType.Uuid,
    Decimal: ValueType.Decimal,
}

# ---------------------------------------------------------------------------
# Registry
# ---------------------------------------------------------------------------
_MM_FIELD_REGISTRY: Dict[Type, Dict[str, Tag]] = {}
_MM_CLASS_REGISTRY: Dict[Type, Tag] = {}


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
def _is_optional(annotation: Any) -> bool:
    """Check if annotation is Optional[T] or T | None."""
    origin = typing.get_origin(annotation)
    if origin in (typing.Union, getattr(types, 'UnionType', ())):
        args = typing.get_args(annotation)
        return type(None) in args
    return False


def _unwrap_optional(annotation: Any) -> Any:
    """Remove Optional / None from a union, return the inner type."""
    if not _is_optional(annotation):
        return annotation
    args = [a for a in typing.get_args(annotation) if a is not type(None)]
    return args[0] if len(args) == 1 else annotation


def _get_origin_type(annotation: Any) -> Any:
    """Get the origin type from an annotation (e.g. list[int] → list)."""
    origin = typing.get_origin(annotation)
    return origin if origin is not None else annotation


def _infer_value_type(annotation: Any) -> Optional[ValueType]:
    """Infer ValueType from a Python type annotation. Returns None if unknown."""
    tp = _get_origin_type(_unwrap_optional(annotation))
    if tp is None or tp is type(None):
        return None
    if isinstance(tp, type) and issubclass(tp, enum.Enum):
        return ValueType.Enums
    return _PYTHON_TYPE_TO_VT.get(tp)


def _mm_fields_dict(inst: mm) -> Dict[str, Any]:
    """Return a dict of all mm dataclass field names → values."""
    result = {}
    for f in dc_fields(inst):
        if f.name.startswith('_'):
            continue
        result[f.name] = getattr(inst, f.name)
    return result


def _find_type_in_tag(tag: Tag) -> ValueType:
    """Determine the effective ValueType for a tag."""
    return tag.type if tag.type is not None else ValueType.Unknown


# ---------------------------------------------------------------------------
# mm decorator
# ---------------------------------------------------------------------------
@dataclass
class mm:
    """MetaMessage tag annotation.

    Can be used as:
    1. Class decorator: @mm(desc="User information")
    2. Field decorator: @mm inside class body
    3. Tag string: @mm("type=i64; desc=用户ID")

    Examples:
        @mm(desc="User information")
        class User:
            @mm(desc="User ID")
            id: int
            @mm(type=ValueType.U8, desc="User age")
            age: int
    """
    tag_str: str = ''
    name: str = ''
    desc: str = ''
    type: Any = None
    is_null: bool = False
    nullable: bool = False
    deprecated: bool = False
    example: bool = False
    allow_empty: bool = False
    unique: bool = False
    default_val: str = ''
    min: str = ''
    max: str = ''
    size: int = 0
    enums: str = ''
    pattern: str = ''
    location: int = 0
    version: int = 0
    mime: str = ''
    child_desc: str = ''
    child_type: Any = None
    child_nullable: bool = False
    child_allow_empty: bool = False
    child_unique: bool = False
    child_default_val: str = ''
    child_min: str = ''
    child_max: str = ''
    child_size: int = 0
    child_enums: str = ''
    child_pattern: str = ''
    child_location: int = 0
    child_version: int = 0
    child_mime: str = ''

    _tag: Optional[Tag] = field(default=None, init=False, repr=False)
    _validated: bool = field(default=False, init=False, repr=False)

    def __post_init__(self):
        self._tag = None

    # ---- public API ----

    def __call__(self, target):
        tag = self._build_tag()
        if isinstance(target, type):
            if tag.type is None or tag.type == ValueType.Unknown:
                tag.type = ValueType.Obj  # default class type
            self._validate_class_level(tag, target)
            _MM_CLASS_REGISTRY[target] = tag
            _MM_FIELD_REGISTRY.setdefault(target, {})
            return dataclass(target)  # apply @dataclass to the user's class
        return target

    def __set_name__(self, owner, name):
        tag = self._build_tag()
        annotation = owner.__annotations__.get(name)
        self._validate_field_level(tag, annotation, owner, name)
        _MM_FIELD_REGISTRY.setdefault(owner, {})
        _MM_FIELD_REGISTRY[owner][name] = tag

    def get_tag(self) -> Tag:
        return self._build_tag()

    # ---- tag building ----

    def _build_tag(self) -> Tag:
        if self._tag is None:
            if self.tag_str:
                self._tag = mm_tag(self.tag_str)
            else:
                self._tag = self._fields_to_tag()
        return self._tag

    def _fields_to_tag(self) -> Tag:
        tag = NewTag()
        # Fields that must be stored as strings on Tag
        _str_fields = {'name', 'desc', 'default_val', 'min', 'max',
                       'enums', 'pattern', 'mime',
                       'child_desc', 'child_default_val', 'child_min', 'child_max',
                       'child_enums', 'child_pattern', 'child_mime'}
        for k, v in {
            'name': self.name, 'desc': self.desc, 'is_null': self.is_null,
            'nullable': self.nullable, 'deprecated': self.deprecated,
            'example': self.example, 'allow_empty': self.allow_empty,
            'unique': self.unique, 'default_val': self.default_val,
            'min': self.min, 'max': self.max,
            'size': self.size, 'enums': self.enums,
            'pattern': self.pattern, 'version': self.version,
            'mime': self.mime,
            'child_desc': self.child_desc,
            'child_nullable': self.child_nullable,
            'child_allow_empty': self.child_allow_empty,
            'child_unique': self.child_unique,
            'child_default_val': self.child_default_val,
            'child_min': self.child_min, 'child_max': self.child_max,
            'child_size': self.child_size, 'child_enums': self.child_enums,
            'child_pattern': self.child_pattern,
            'child_version': self.child_version, 'child_mime': self.child_mime,
        }.items():
            if k in _str_fields:
                # Convert to str first so numeric falsy values like 0, 0.0
                # (e.g. min=0.0) become truthy "0.0" and are not skipped.
                v = str(v) if v is not None else v
                if v:
                    setattr(tag, k, str(v))
            elif v or (k == 'size' and v != 0):
                setattr(tag, k, v)

        if self.type is not None:
            if isinstance(self.type, ValueType):
                tag.type = self.type
            elif isinstance(self.type, str):
                from ..ir.tag import parse_value_type
                tag.type = parse_value_type(self.type)
            elif isinstance(self.type, int):
                tag.type = ValueType(self.type)
        if self.child_type is not None:
            if isinstance(self.child_type, ValueType):
                tag.child_type = self.child_type
            elif isinstance(self.child_type, str):
                from ..ir.tag import parse_value_type
                tag.child_type = parse_value_type(self.child_type)

        if self.enums:
            tag.type = ValueType.Enums
        if self.mime:
            tag.type = ValueType.Media
        if self.child_enums:
            tag.child_type = ValueType.Enums
        if self.child_mime:
            tag.child_type = ValueType.Media

        return tag

    # ---- validation ----

    def _user_set_fields(self) -> Dict[str, Any]:
        """Return dict of fields that the user explicitly set (non-default)."""
        defaults = {
            'tag_str': '', 'name': '', 'desc': '', 'type': None,
            'is_null': False, 'nullable': False, 'deprecated': False,
            'example': False, 'allow_empty': False, 'unique': False,
            'default_val': '', 'min': '', 'max': '', 'size': 0,
            'enums': '', 'pattern': '', 'location': 0, 'version': 0, 'mime': '',
            'child_desc': '', 'child_type': None, 'child_nullable': False,
            'child_allow_empty': False, 'child_unique': False,
            'child_default_val': '', 'child_min': '', 'child_max': '',
            'child_size': 0, 'child_enums': '', 'child_pattern': '',
            'child_location': 0, 'child_version': 0, 'child_mime': '',
        }
        result = {}
        for k, v in _mm_fields_dict(self).items():
            default = defaults.get(k)
            if v != default:
                result[k] = v
        return result

    def _validate_allowed_tags(self, vt: ValueType, context: str):
        """Validate that all user-set tags are allowed for this ValueType."""
        allowed = ALLOWED_TAGS.get(vt)
        if allowed is None:
            return  # unknown type, skip

        for field_name in self._user_set_fields():
            if field_name in ('tag_str', 'type', 'child_type'):
                continue
            if field_name not in allowed:
                raise ValueError(
                    f"{context}: ValueType.{vt.name} 不允许使用标签 '{field_name}'。"
                    f"\n  允许的标签: {', '.join(sorted(allowed))}"
                )

    def _validate_is_null(self, tag: Tag, annotation: Any, context: str):
        """Validate is_null: only allowed when the Python type is Optional."""
        if not tag.is_null:
            return
        if annotation is None:
            return  # no annotation, can't check
        if not _is_optional(annotation):
            raise ValueError(
                f"{context}: 设置了 is_null=True，但字段类型不是 Optional。"
                f"\n  当前类型: {annotation}"
                f"\n  请使用 Optional[{annotation}] 或 {annotation} | None"
            )

    def _validate_class_level(self, tag: Tag, cls: Type):
        vt = _find_type_in_tag(tag)
        if vt == ValueType.Unknown:
            vt = ValueType.Obj  # default class type

        if vt not in _CLASS_LEVEL_TYPES:
            raise ValueError(
                f"类装饰器 @mm 只允许 Obj 或 Map 类型，"
                f"当前为 ValueType.{vt.name}。"
                f"\n  类级别仅支持: ValueType.Obj, ValueType.Map"
                f"\n  请使用字段装饰器: field: type = mm(...)"
            )

        context = f"类 '{cls.__name__}'"
        self._validate_allowed_tags(vt, context)

    def _validate_field_level(self, tag: Tag, annotation: Any, owner: Type, name: str):
        context = f"字段 '{owner.__name__}.{name}'"

        vt = _find_type_in_tag(tag)
        inferred_vt = _infer_value_type(annotation) if annotation is not None else None

        # --- class-level types not allowed at field level (check first) ---
        if vt in _CLASS_LEVEL_TYPES and annotation is not None:
            tp = _get_origin_type(_unwrap_optional(annotation))
            if not isinstance(tp, type) or tp in _PYTHON_TYPE_TO_VT:
                raise ValueError(
                    f"{context}: ValueType.{vt.name} 只能用于类装饰器，不能用于字段。"
                    f"\n  字段类型: {annotation}"
                    f"\n  请使用类装饰器: @mm(...) class MyClass: ..."
                )

        # --- type validation ---
        if vt != ValueType.Unknown and inferred_vt is not None:
            if vt != inferred_vt:
                _int_types = {ValueType.I, ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64,
                              ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64,
                              ValueType.Bigint}
                _float_types = {ValueType.F32, ValueType.F64}
                _list_types = {ValueType.Vec, ValueType.Arr}
                _str_types = {ValueType.Str, ValueType.Ip, ValueType.Url,
                              ValueType.Email, ValueType.Uuid, ValueType.Decimal}
                _bytes_types = {ValueType.Bytes, ValueType.Media}

                compatible = (
                    (inferred_vt in _int_types and vt in _int_types) or
                    (inferred_vt in _float_types and vt in _float_types) or
                    (inferred_vt in _list_types and vt in _list_types) or
                    (inferred_vt in _str_types and vt in _str_types) or
                    (inferred_vt in _bytes_types and vt in _bytes_types)
                )
                if not compatible:
                    raise ValueError(
                        f"{context}: 类型不匹配。"
                        f"\n  字段 Python 类型: {annotation} → 应为 ValueType.{inferred_vt.name}"
                        f"\n  但 @mm 中指定了 type={vt.name}"
                        f"\n  建议: 移除 type 参数让其自动推断，或使用正确的 ValueType"
                    )

        if vt == ValueType.Unknown and inferred_vt is not None:
            tag.type = inferred_vt
            vt = inferred_vt

        # --- allowed tags ---
        if vt != ValueType.Unknown:
            self._validate_allowed_tags(vt, context)

        # --- is_null constraint ---
        self._validate_is_null(tag, annotation, context)


# ---------------------------------------------------------------------------
# Public helpers
# ---------------------------------------------------------------------------
def get_mm_tag_for_class(cls: Type) -> Optional[Tag]:
    """Get the class-level MM tag for a class."""
    return _MM_CLASS_REGISTRY.get(cls)


def get_mm_tag_for_field(cls: Type, field_name: str) -> Optional[Tag]:
    """Get the field-level MM tag for a class field."""
    field_registry = _MM_FIELD_REGISTRY.get(cls)
    if field_registry:
        return field_registry.get(field_name)
    return None