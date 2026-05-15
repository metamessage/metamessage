from typing import Any, List, Optional
from dataclasses import dataclass, field
from enum import IntEnum


class NodeType(IntEnum):
    Unknown = 0
    Object = 1
    Array = 2
    Value = 3
    Doc = 4


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
    value: Any = None  # Node


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
    items: List[Any] = field(default_factory=list)  # List[Node]
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
