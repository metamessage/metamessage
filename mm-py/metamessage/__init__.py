"""
metamessage - A binary message encoding library with schema support
"""
import builtins as _builtins
from dataclasses import dataclass as _dataclass, field as _field, MISSING as _MISSING, Field as _Field
import typing as _typing
import types as _types

from .ir.tag import Tag, TagKey, ValueType, mm_tag, def_tag, NewTag, MergeTag
from .ir.ast import NodeObject, NodeArray, NodeScalar, Field, NodeType, Node, NodeNull
from .core.encoder import Encoder
from .core.decoder import Decoder
from .core.value_to_node import value_to_node, node_to_value
from .core.mm import mm
from .jsonc import parse_jsonc, to_jsonc


def _infer_zero_field(annotation):
    """Infer a dataclasses.field() with a zero-value default from a type annotation.
    
    Handles plain types, generic types (list[X], dict[X,Y]), and Optional[T].
    Falls back to None for unknown types.
    """
    # Unwrap Optional[T] → T, mark as optional
    is_optional = False
    origin = _typing.get_origin(annotation)
    if origin in (_typing.Union, getattr(_types, 'UnionType', ())):
        args = _typing.get_args(annotation)
        if type(None) in args:
            is_optional = True
            args = [a for a in args if a is not type(None)]
            annotation = args[0] if args else None

    if is_optional:
        return _field(default=None)

    # Resolve origin of generic types: list[int] → list, dict[str,int] → dict
    origin = _typing.get_origin(annotation)
    raw = origin if origin is not None else annotation

    if raw is str or raw is type(None):
        return _field(default="")
    if raw is int:
        return _field(default=0)
    if raw is float:
        return _field(default=0.0)
    if raw is bool:
        return _field(default=False)
    if raw is bytes:
        return _field(default=b"")
    if raw is list:
        return _field(default_factory=list)
    if raw is dict:
        return _field(default_factory=dict)
    if raw is tuple:
        return _field(default_factory=tuple)
    if raw is set:
        return _field(default_factory=set)

    # Unknown type — let dataclass handle it naturally (required field)
    return None


def _parse_default_val(val: str, raw_type) -> any:
    """Parse a default_val string into the proper Python type based on annotation."""
    if raw_type is str or raw_type is type(None):
        return val
    if raw_type is int:
        try:
            return int(val)
        except (ValueError, TypeError):
            return None
    if raw_type is float:
        try:
            return float(val)
        except (ValueError, TypeError):
            return None
    if raw_type is bool:
        return val.lower() in ('true', '1', 'yes')
    if raw_type is bytes:
        return val.encode('utf-8')
    # Complex types (list, dict, etc.) — fall back to zero value
    return None


def _set_zero_value_defaults(cls):
    """Set zero-value defaults for all fields that don't have explicit defaults.

    This is called BEFORE dataclass() so that every field has a default,
    avoiding Python 3.11+ field ordering errors (required fields before optional).
    Fields that already have an explicit default (set by __set_name__ or directly)
    are left untouched.

    If a field's tag (from _MM_FIELD_REGISTRY) has a default_val, that value is
    used in preference to the type's zero value.
    """
    from .core.mm import _MM_FIELD_REGISTRY
    field_tags = _MM_FIELD_REGISTRY.get(cls, {})

    for fname, annotation in cls.__annotations__.items():
        existing = getattr(cls, fname, None)
        # Skip if field already has a non-MISSING default
        if isinstance(existing, _Field):
            if existing.default is not _MISSING or existing.default_factory is not _MISSING:
                continue
            # field(default=MISSING) — replace with zero value
        elif existing is not None and not isinstance(existing, _Field):
            # Already has a plain default value (e.g. `x: int = 5`)
            continue

        # Resolve raw type for generic annotations
        origin = _typing.get_origin(annotation)
        raw = origin if origin is not None else annotation

        # Unwrap Optional[T] (Union[T, None]) → T for default_val parsing
        if raw in (_typing.Union, getattr(_types, 'UnionType', ())):
            union_args = _typing.get_args(annotation)
            non_none = tuple(a for a in union_args if a is not type(None))
            if non_none:
                raw = non_none[0]

        # If tag has default_val, parse it as the Python default
        tag = field_tags.get(fname)
        if tag and tag.default_val:
            parsed = _parse_default_val(tag.default_val, raw)
            if parsed is not None:
                setattr(cls, fname, _field(default=parsed))
                continue

        # Fall back to zero value for the type
        zero = _infer_zero_field(annotation)
        if zero is not None:
            setattr(cls, fname, zero)


# Auto-apply @mm + @dataclass to any class with type annotations
_build_class = getattr(_builtins, '__build_class__', None)
if _build_class is not None:
    def _mm_build_class(func, name, *bases, **kwargs):
        cls = _build_class(func, name, *bases, **kwargs)
        if (isinstance(cls, type)
                and hasattr(cls, '__annotations__')
                and cls.__annotations__):
            # Skip classes from site-packages/third-party libraries to avoid
            # interfering with their dataclass field ordering (e.g. anyio).
            mod_file = func.__globals__.get('__file__', '') or ''
            if 'site-packages' in mod_file.replace('\\', '/'):
                return cls
            # Resolve PEP 563 string annotations (from __future__ import annotations)
            # so that downstream code (type inference, _is_optional, @dataclass)
            # sees real type objects instead of strings.
            if any(isinstance(v, str) for v in cls.__annotations__.values()):
                try:
                    hints = _typing.get_type_hints(cls)
                    cls.__annotations__ = hints
                except Exception:
                    pass
            if not hasattr(cls, '__dataclass_fields__'):
                try:
                    # Set zero-value defaults for all fields before @dataclass
                    # so that field ordering is never an issue.
                    _set_zero_value_defaults(cls)
                    _dataclass(cls)
                    from .core.mm import _MM_CLASS_REGISTRY, _MM_FIELD_REGISTRY
                    _MM_CLASS_REGISTRY.setdefault(cls, NewTag())
                    _MM_FIELD_REGISTRY.setdefault(cls, {})
                except Exception:
                    # If @dataclass fails (e.g. Python 3.13+ field ordering),
                    # set __dataclass_fields__ to prevent the outer decorator
                    # from re-processing and crashing with the same error.
                    if not hasattr(cls, '__dataclass_fields__'):
                        try:
                            cls.__dataclass_fields__ = {}
                        except Exception:
                            pass
        return cls
    _builtins.__build_class__ = _mm_build_class


def encode_from_value(value, tag=None):
    """Convert a Python value directly to MetaMessage binary format.

    Python equivalent of mm-ts encodeFromValue / Go ValueToNode+Encode.
    """
    node = value_to_node(value, tag)
    encoder = Encoder()
    return encoder.encode(node)


def encode_from_jsonc(jsonc):
    """Convert a JSONC string to MetaMessage binary format.

    Python equivalent of mm-ts encodeFromJsonc.
    """
    node = parse_jsonc(jsonc)
    encoder = Encoder()
    return encoder.encode(node)


def decode_to_value(data, target_type=None):
    """Decode MetaMessage binary format to a Python value.

    Python equivalent of mm-ts decodeToValue / Go Decode+Bind.
    """
    decoder = Decoder(data)
    node = decoder.decode_node()

    if target_type is None or target_type is dict or target_type is list:
        from .core.decoder import _node_to_python
        return _node_to_python(node)

    return node_to_value(node, target_type)


def decode_to_jsonc(data):
    """Decode MetaMessage binary format to a JSONC string.

    Python equivalent of mm-ts decodeToJsonc.
    """
    decoder = Decoder(data)
    node = decoder.decode_node()
    return to_jsonc(node)


def value_to_jsonc(value, tag=None):
    """Convert a Python value to a JSONC string.

    Python equivalent of mm-ts valueToJsonc.
    """
    node = value_to_node(value, tag)
    return to_jsonc(node)


def jsonc_to_value(jsonc, target_type=None):
    """Convert a JSONC string to a Python value.

    Python equivalent of mm-ts jsoncToValue.
    """
    node = parse_jsonc(jsonc)
    if target_type is None:
        return node.data if isinstance(node, NodeScalar) else _node_data(node)
    return node_to_value(node, target_type)


def _node_data(node):
    if isinstance(node, NodeNull):
        return None
    elif isinstance(node, NodeObject):
        return {f.key: _node_data(f.value) for f in node.fields}
    elif isinstance(node, NodeArray):
        return [_node_data(item) for item in node.items]
    elif isinstance(node, NodeScalar):
        return node.data
    return None


__all__ = [
    "Tag", "TagKey", "ValueType", "mm_tag", "def_tag", "NewTag", "MergeTag",
    "NodeObject", "NodeArray", "NodeScalar", "NodeNull", "Field", "NodeType", "Node",
    "Encoder", "Decoder",
    "parse_jsonc", "to_jsonc",
    "value_to_node", "node_to_value", "mm",
    "encode_from_value", "encode_from_jsonc",
    "decode_to_value", "decode_to_jsonc",
    "value_to_jsonc", "jsonc_to_value",
]