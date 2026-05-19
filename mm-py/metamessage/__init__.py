"""
metamessage - A binary message encoding library with schema support
"""
from .ir.tag import Tag, TagKey, ValueType, mm_tag, def_tag, NewTag, MergeTag
from .ir.types import Obj, Arr, Val, Field, NodeType, Node
from .core.encoder import Encoder
from .core.decoder import Decoder
from .core.value_to_node import value_to_node, node_to_value, mm
from .jsonc import parse_jsonc, to_jsonc


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
        return node.data if isinstance(node, Val) else _node_data(node)
    return node_to_value(node, target_type)


def _node_data(node):
    if isinstance(node, Obj):
        return {f.key: _node_data(f.value) for f in node.fields}
    elif isinstance(node, Arr):
        return [_node_data(item) for item in node.items]
    elif isinstance(node, Val):
        return node.data
    return None


__all__ = [
    "Tag", "TagKey", "ValueType", "mm_tag", "def_tag", "NewTag", "MergeTag",
    "Obj", "Arr", "Val", "Field", "NodeType", "Node",
    "Encoder", "Decoder",
    "parse_jsonc", "to_jsonc",
    "value_to_node", "node_to_value", "mm",
    "encode_from_value", "encode_from_jsonc",
    "decode_to_value", "decode_to_jsonc",
    "value_to_jsonc", "jsonc_to_value",
]