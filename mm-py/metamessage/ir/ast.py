from typing import Any, List, Optional, TYPE_CHECKING
from dataclasses import dataclass, field
from enum import IntEnum

if TYPE_CHECKING:
    from .tag import Tag


class NodeType(IntEnum):
    Unknown = 0
    Object = 1
    Array = 2
    Value = 3
    Doc = 4

    def __str__(self) -> str:
        mapping = {
            NodeType.Unknown: "unknown",
            NodeType.Object: "object",
            NodeType.Array: "array",
            NodeType.Value: "value",
            NodeType.Doc: "doc",
        }
        return mapping.get(self, "unknown")


def parse_node_type(s: str) -> NodeType:
    s = s.lower().strip()
    mapping = {
        "unknown": NodeType.Unknown,
        "object": NodeType.Object,
        "array": NodeType.Array,
        "value": NodeType.Value,
        "doc": NodeType.Doc,
    }
    return mapping.get(s, NodeType.Unknown)


class Node:
    def get_tag(self) -> Optional['Tag']:
        return None

    def get_type(self) -> NodeType:
        return NodeType.Unknown

    def get_path(self) -> str:
        return ""

    def set_path(self, path: str):
        pass


@dataclass
class Field:
    key: str = ""
    value: 'Node' = None


@dataclass
class Obj(Node):
    fields: List[Field] = field(default_factory=list)
    tag: Optional['Tag'] = None
    path: str = ""

    def get_tag(self):
        return self.tag

    def get_type(self):
        return NodeType.Object

    def get_path(self):
        return self.path

    def set_path(self, path: str):
        self.path = path


@dataclass
class Arr(Node):
    items: List['Node'] = field(default_factory=list)
    tag: Optional['Tag'] = None
    path: str = ""

    def get_tag(self):
        return self.tag

    def get_type(self):
        return NodeType.Array

    def get_path(self):
        return self.path

    def set_path(self, path: str):
        self.path = path


@dataclass
class Val(Node):
    data: Any = None
    text: str = ""
    tag: Optional['Tag'] = None
    path: str = ""

    def get_tag(self):
        return self.tag

    def get_type(self):
        return NodeType.Value

    def get_path(self):
        return self.path

    def set_path(self, path: str):
        self.path = path


@dataclass
class Doc(Node):
    fields: List[Field] = field(default_factory=list)
    tag: Optional['Tag'] = None
    path: str = ""

    def get_tag(self):
        return self.tag

    def get_type(self):
        return NodeType.Doc

    def get_path(self):
        return self.path

    def set_path(self, path: str):
        self.path = path


Empty = ""
Null = "null"
TrueStr = "true"
FalseStr = "false"

SimpleCodeStr = "code"
SimpleMessageStr = "message"
SimpleDataStr = "data"
SimpleSuccessStr = "success"
SimpleErrorStr = "error"
SimpleUnknownStr = "unknown"

SimplePageStr = "page"
SimpleLimitStr = "limit"
SimpleOffsetStr = "offset"
SimpleTotalStr = "total"
SimpleIdStr = "id"
SimpleNameStr = "name"
SimpleDescriptionStr = "description"
SimpleTypeStr = "type"
SimpleVersionStr = "version"
SimpleStatusStr = "status"
SimpleUrlStr = "url"
SimpleCreateTimeStr = "create_time"
SimpleUpdateTimeStr = "update_time"
SimpleDeleteTimeStr = "delete_time"
SimpleAccountStr = "account"
SimpleTokenStr = "token"
SimpleExpireTimeStr = "expire_time"
SimpleKeyStr = "key"
SimpleValStr = "value"

import sys
BitSize = 32 << (1 if sys.maxsize > 2**32 else 0)