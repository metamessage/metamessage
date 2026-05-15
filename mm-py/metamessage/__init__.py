"""
metamessage - A binary message encoding library with schema support
"""
from .ir.tag import Tag, TagKey, ValueType, mm_tag, def_tag, NewTag, MergeTag
from .ir.types import Obj, Arr, Val, Field, NodeType, Node
from .core.encoder import Encoder
from .core.decoder import Decoder
from .jsonc import parse_jsonc, to_jsonc
from .core.value_to_node import value_to_node, node_to_value, encode_from_value, decode_to_value, mm
__all__ = [
    "Tag", "TagKey", "ValueType", "mm_tag", "def_tag", "NewTag", "MergeTag",
    "Obj", "Arr", "Val", "Field", "NodeType", "Node",
    "Encoder", "Decoder",
    "parse_jsonc", "to_jsonc",
    "value_to_node", "node_to_value", "encode_from_value", "decode_to_value", "mm",
]