"""
metamessage - A binary message encoding library with schema support
"""
import builtins as _builtins
from dataclasses import dataclass as _dataclass

from .ir.tag import Tag, TagKey, ValueType, mm_tag, def_tag, NewTag, MergeTag
from .ir.ast import NodeObject, NodeArray, NodeScalar, Field, NodeType, Node, NodeNull
from .core.encoder import Encoder
from .core.decoder import Decoder
from .core.value_to_node import value_to_node, node_to_value
from .core.mm import mm
from .jsonc import parse_jsonc, to_jsonc


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
                    import typing
                    hints = typing.get_type_hints(cls)
                    cls.__annotations__ = hints
                except Exception:
                    pass
            if not hasattr(cls, '__dataclass_fields__'):
                try:
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