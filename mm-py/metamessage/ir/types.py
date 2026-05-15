from typing import Any, Dict, List, Optional
from dataclasses import dataclass, field


@dataclass
class Obj:
    fields: List['Field'] = field(default_factory=list)
    tag: Optional['Tag'] = None
    path: str = ""


@dataclass
class Arr:
    items: List[Any] = field(default_factory=list)
    tag: Optional['Tag'] = None
    path: str = ""


@dataclass
class Val:
    data: Any = None
    text: str = ""
    tag: Optional['Tag'] = None
    path: str = ""


@dataclass
class Field:
    key: str = ""
    value: Any = None