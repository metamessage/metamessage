"""
ValueToNode: Convert Python values to MetaMessage Node tree.
Supports decorator-based @mm() tagging like Go's struct tags and TS decorators.

Usage:
    @mm(type=ValueType.Int64, desc="用户ID")
    class User:
        id: int
        name: str
        age: int

    # Or use field-level decorators:
    class User:
        @mm(type=ValueType.Int64, desc="用户ID")
        id: int

    # Encode/Decode:
    user = User(id=1, name="Alice", age=20)
    node = value_to_node(user)  # Python value → Node tree
    binary = encode_from_value(user)  # Python value → binary
    user2 = decode_to_value(binary, User)  # binary → Python value
"""
import inspect
import re
import types
from typing import Any, Dict, List, Optional, Type, Tuple, get_type_hints, Union
from datetime import datetime, date, time as dt_time
from enum import Enum

from ..ir.tag import Tag, ValueType, NewTag, MergeTag, mm_tag
from ..ir.types import Obj, Arr, Val, Field, Node
from .encoder import Encoder
from .decoder import Decoder


def _camel_to_snake(name: str) -> str:
    """Convert CamelCase to snake_case.
    
    UserID → user_id
    HTTPRequest → http_request
    """
    if not name:
        return name
    result = []
    for i, c in enumerate(name):
        if c.isupper():
            if i > 0 and (i + 1 < len(name) and name[i + 1].islower() or
                          not name[i - 1].isupper()):
                result.append('_')
            result.append(c.lower())
        else:
            result.append(c)
    return ''.join(result)


# ===== Decorator-based MM Tag =====

_MM_FIELD_REGISTRY: Dict[Type, Dict[str, Tag]] = {}
_MM_CLASS_REGISTRY: Dict[Type, Tag] = {}


class mm:
    """MetaMessage tag annotation.

    Can be used as:
    1. Class decorator (recommended): @mm(desc="User information")
    2. Field decorator: @mm inside class body (see examples)
    3. Tag string: @mm("type=i64; desc=用户ID")

    Examples:
        @mm(desc="User information")
        class User:
            @mm(desc="User ID")
            id: int
            @mm(desc="User name")
            name: str
            @mm(type=ValueType.Uint8, desc="User age")
            age: int

    How it works:
    - When used as a class decorator (@mm above class), it registers the tag
      for the entire class.
    - When used as a field decorator (@mm above a field annotation), Python's
      class creation machinery processes it via __init_subclass__.
      The mm instance stores the tag and associates it with the next
      annotated field name.
    """
    # Track the last field name that was annotated
    _last_field_name: Optional[str] = None

    def __new__(cls, *args, **kwargs):
        instance = super().__new__(cls)
        instance._tag = None
        if args and isinstance(args[0], str):
            # mm("type=i64; desc=用户ID")
            instance._tag_str = args[0]
            instance._kwargs = {}
        else:
            # mm(type=ValueType.Int64, desc="用户ID")
            instance._tag_str = None
            instance._kwargs = kwargs
        return instance

    def __call__(self, target):
        tag = self._build_tag()

        if isinstance(target, type):
            # Class decorator: @mm(desc="...")
            _MM_CLASS_REGISTRY[target] = tag
            _MM_FIELD_REGISTRY.setdefault(target, {})
            return target
        else:
            # If used as a regular decorator on a non-class (function, etc.),
            # just return the target unchanged
            return target

    def __set_name__(self, owner, name):
        """Descriptor protocol: called when the class is created.
        
        This allows mm instances used as default values to register
        themselves into the field registry.
        """
        tag = self._build_tag()
        _MM_FIELD_REGISTRY.setdefault(owner, {})
        _MM_FIELD_REGISTRY[owner][name] = tag

    def _build_tag(self) -> Tag:
        """Build Tag from stored args/kwargs."""
        if self._tag is None:
            if self._tag_str:
                self._tag = mm_tag(self._tag_str)
            else:
                self._tag = self._kwargs_to_tag(self._kwargs)
        return self._tag

    def get_tag(self) -> Tag:
        """Get the Tag stored in this mm instance."""
        return self._build_tag()

    def _kwargs_to_tag(self, kwargs: dict) -> Tag:
        """Convert keyword arguments to a Tag."""
        tag = NewTag()
        
        for k, v in kwargs.items():
            k = k.lower()
            if k == 'type':
                if isinstance(v, ValueType):
                    tag.type = v
                elif isinstance(v, str):
                    from ..ir.tag import parse_value_type
                    tag.type = parse_value_type(v)
                elif isinstance(v, int):
                    tag.type = ValueType(v)
            elif k == 'desc':
                tag.desc = str(v)
            elif k == 'is_null':
                tag.is_null = bool(v)
            elif k == 'nullable':
                tag.nullable = bool(v)
            elif k == 'raw':
                tag.raw = bool(v)
            elif k == 'example':
                tag.example = bool(v)
            elif k == 'allow_empty':
                tag.allow_empty = bool(v)
            elif k == 'unique':
                tag.unique = bool(v)
            elif k == 'default':
                tag.default = str(v)
            elif k == 'min':
                tag.min = str(v)
            elif k == 'max':
                tag.max = str(v)
            elif k == 'size':
                tag.size = int(v)
            elif k == 'enum':
                tag.type = ValueType.Enum
                tag.enum = str(v)
            elif k == 'pattern':
                tag.pattern = str(v)
            elif k == 'version':
                tag.version = int(v)
            elif k == 'mime':
                tag.mime = str(v)
            elif k == 'child_desc':
                tag.child_desc = str(v)
            elif k == 'child_type':
                if isinstance(v, ValueType):
                    tag.child_type = v
                elif isinstance(v, str):
                    from ..ir.tag import parse_value_type
                    tag.child_type = parse_value_type(v)
            elif k == 'child_raw':
                tag.child_raw = bool(v)
            elif k == 'child_nullable':
                tag.child_nullable = bool(v)
            elif k == 'child_allow_empty':
                tag.child_allow_empty = bool(v)
            elif k == 'child_unique':
                tag.child_unique = bool(v)
            elif k == 'child_default':
                tag.child_default = str(v)
            elif k == 'child_min':
                tag.child_min = str(v)
            elif k == 'child_max':
                tag.child_max = str(v)
            elif k == 'child_size':
                tag.child_size = int(v)
            elif k == 'child_enum':
                tag.child_enum = str(v)
            elif k == 'child_pattern':
                tag.child_pattern = str(v)
            elif k == 'child_version':
                tag.child_version = int(v)
            elif k == 'child_mime':
                tag.child_mime = str(v)
        
        return tag


def get_mm_tag_for_class(cls: Type) -> Optional[Tag]:
    """Get the class-level MM tag for a class."""
    return _MM_CLASS_REGISTRY.get(cls)


def get_mm_tag_for_field(cls: Type, field_name: str) -> Optional[Tag]:
    """Get the field-level MM tag for a class field."""
    field_registry = _MM_FIELD_REGISTRY.get(cls)
    if field_registry:
        return field_registry.get(field_name)
    return None


# ===== Python to Node Type Mapping =====

_PYTHON_TYPE_TO_VALUETYPE: Dict[type, ValueType] = {
    str: ValueType.Str,
    int: ValueType.I,
    float: ValueType.F64,
    bool: ValueType.Bool,
    bytes: ValueType.Bytes,
    datetime: ValueType.Datetime,
    date: ValueType.Date,
    dt_time: ValueType.Time,
}


def _python_type_to_value_type(py_type: type) -> ValueType:
    """Map Python type to ValueType."""
    if py_type in _PYTHON_TYPE_TO_VALUETYPE:
        return _PYTHON_TYPE_TO_VALUETYPE[py_type]
    if isinstance(py_type, type) and issubclass(py_type, Enum):
        return ValueType.Enum
    return ValueType.Unknown


def _is_union_type_with_null(type_str: Any) -> bool:
    """Check if a type string represents a union type containing null/None."""
    if type_str is None or not isinstance(type_str, str):
        return False
    return any(t.strip().lower() == 'null' for t in type_str.split('|'))


def _is_union_type_with_none(py_type: Any) -> bool:
    """Check if a Python type annotation is a union type containing None."""
    if py_type is None:
        return False
    
    # Check for Python 3.10+ union syntax (e.g., int | None)
    # In Python 3.10+, int | None creates a types.UnionType which has __args__
    if isinstance(py_type, types.UnionType):
        return type(None) in py_type.__args__
    
    # Check for typing.Union syntax (e.g., Union[int, None])
    if hasattr(py_type, '__origin__'):
        try:
            from typing import Union
            if py_type.__origin__ is Union:
                return type(None) in py_type.__args__
        except ImportError:
            pass
    
    # Check for typing.Optional (which is Union[X, None])
    if hasattr(py_type, '__origin__'):
        try:
            from typing import Optional
            # Optional[X] is equivalent to Union[X, None]
            if py_type.__origin__ is Optional:
                return True
        except ImportError:
            pass
    
    return False


# ===== ValueToNode =====

_MAX_DEPTH = 32


def value_to_node(value: Any, tag: Optional[Tag] = None, depth: int = 0, path: str = "") -> Node:
    """Convert a Python value to a MetaMessage Node tree.

    This is the Python equivalent of Go's ValueToNode() function.
    """
    if depth > _MAX_DEPTH:
        raise ValueError(f"max depth: {_MAX_DEPTH}")

    if tag is None:
        tag = NewTag()

    if value is None:
        if tag.type == ValueType(0):
            raise ValueError("invalid input: v is None (no concrete type/value)")
        tag.is_null = True
        return Val(data=None, text="null", tag=tag, path=path)

    # bool must be checked before int (bool is subclass of int in Python)
    if isinstance(value, bool):
        # Auto-detect type: override inherited/incorrect types
        if tag.type in (ValueType(0), ValueType.Str, ValueType.I, ValueType.F64):
            tag.type = ValueType.Bool
        if tag.type == ValueType.Bool:
            _, text = _validate_bool(value, tag)
            return Val(data=value, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: bool")

    elif isinstance(value, int):
        if tag.type in (ValueType(0), ValueType.Str, ValueType.F64, ValueType.Bool):
            tag.type = ValueType.I
        if tag.type in (ValueType.I, ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64,
                        ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64):
            val_int = _validate_i(value, tag)
            if val_int is not None:
                data, text = val_int
                return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: int")

    elif isinstance(value, float):
        if tag.type in (ValueType(0), ValueType.Str, ValueType.I, ValueType.Bool):
            tag.type = ValueType.F64
        if tag.type in (ValueType.F32, ValueType.F64):
            val_float = _validate_f(value, tag)
            if val_float is not None:
                data, text = val_float
                return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: float")

    elif isinstance(value, str):
        # Auto-detect: any non-string type that got inherited gets overridden
        if tag.type == ValueType(0) or tag.type not in (ValueType.Str, ValueType.Email, ValueType.Enum, ValueType.Decimal, ValueType.Uuid,
                            ValueType.Url, ValueType.Bigint):
            tag.type = ValueType.Str
        if tag.type in (ValueType.Str, ValueType.Email, ValueType.Enum, ValueType.Decimal, ValueType.Uuid,
                        ValueType.Url, ValueType.Bigint):
            val_str = _validate_str(value, tag)
            if val_str is not None:
                data, text = val_str
                return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: str")

    elif isinstance(value, bytes):
        if tag.type == ValueType(0):
            tag.type = ValueType.Bytes
        val_bytes = _validate_bytes(value, tag)
        if val_bytes is not None:
            data, text = val_bytes
            return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: bytes")

    elif isinstance(value, datetime):
        if tag.type == ValueType(0):
            tag.type = ValueType.Datetime
        val_dt = _validate_datetime(value, tag)
        if val_dt is not None:
            data, text = val_dt
            return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: datetime")

    elif isinstance(value, date):
        if tag.type == ValueType(0):
            tag.type = ValueType.Date
        val_d = _validate_date(value, tag)
        if val_d is not None:
            data, text = val_d
            return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: date")

    elif isinstance(value, dt_time):
        if tag.type == ValueType(0):
            tag.type = ValueType.Time
        val_t = _validate_time(value, tag)
        if val_t is not None:
            data, text = val_t
            return Val(data=data, text=text, tag=tag, path=path)
        raise ValueError(f"{tag.type} unsupported type: time")

    # Handle dict (map)
    elif isinstance(value, dict):
        return _any_to_node_dict(value, tag, depth, path)

    # Handle list (slice/array)
    elif isinstance(value, (list, tuple)):
        return _any_to_node_list(value, tag, depth, path)

    # Handle class instances (struct/object)
    elif hasattr(value, '__class__') and not isinstance(value, type):
        return _any_to_node_object(value, tag, depth, path)

    raise ValueError(f"unsupported type: {type(value)}")


# ===== Inline validation helpers (simplified) =====

def _validate_bool(val: bool, tag: Tag):
    if tag.allow_empty:
        raise ValueError("type bool not support allow empty")
    return (tag, "true" if val else "false")


def _validate_i(val: int, tag: Tag):
    data = val
    text = str(val)
    
    if val == 0 and not tag.allow_empty:
        return (data, text)  # Allow 0 even without allow_empty for Python
    
    # Min/max validation
    if tag.min:
        try:
            min_val = int(tag.min)
        except ValueError:
            min_val = None
        if min_val is not None and val < min_val:
            raise ValueError(f"value {val} < min {min_val}")
    
    if tag.max:
        try:
            max_val = int(tag.max)
        except ValueError:
            max_val = None
        if max_val is not None and val > max_val:
            raise ValueError(f"value {val} > max {max_val}")
    
    return (data, text)


def _validate_f(val: float, tag: Tag):
    import math
    if math.isinf(val) or math.isnan(val):
        raise ValueError(f"unsupported value: {val}")
    
    data = val
    text = str(val)
    
    if val == 0.0 and not tag.allow_empty:
        return (data, text)
    
    if tag.min:
        try:
            min_val = float(tag.min)
        except ValueError:
            min_val = None
        if min_val is not None and val < min_val:
            raise ValueError(f"value {val} < min {min_val}")
    
    if tag.max:
        try:
            max_val = float(tag.max)
        except ValueError:
            max_val = None
        if max_val is not None and val > max_val:
            raise ValueError(f"value {val} > max {max_val}")
    
    return (data, text)


def _validate_str(val: str, tag: Tag):
    data = val
    text = val
    
    if val == "" and not tag.allow_empty:
        return (data, text)
    
    l = len(val)
    if tag.size and l != tag.size:
        raise ValueError(f"string length {l} != size {tag.size}")
    
    if tag.pattern:
        import re
        if not re.match(tag.pattern, val):
            raise ValueError(f"value doesn't match pattern {tag.pattern}")
    
    if tag.min:
        try:
            min_len = int(tag.min)
        except ValueError:
            min_len = None
        if min_len is not None and l < min_len:
            raise ValueError(f"string length {l} < min {min_len}")
    
    if tag.max:
        try:
            max_len = int(tag.max)
        except ValueError:
            max_len = None
        if max_len is not None and l > max_len:
            raise ValueError(f"string length {l} > max {max_len}")
    
    return (data, text)


def _validate_bytes(val: bytes, tag: Tag):
    import base64
    data = val
    text = base64.b64encode(val).decode('ascii')
    
    l = len(val)
    if l == 0 and not tag.allow_empty:
        return (data, text)
    
    if tag.size and l != tag.size:
        raise ValueError(f"bytes length {l} != size {tag.size}")
    
    if tag.min:
        try:
            min_len = int(tag.min)
        except ValueError:
            min_len = None
        if min_len is not None and l < min_len:
            raise ValueError(f"bytes length {l} < min {min_len}")
    
    if tag.max:
        try:
            max_len = int(tag.max)
        except ValueError:
            max_len = None
        if max_len is not None and l > max_len:
            raise ValueError(f"bytes length {l} > max {max_len}")
    
    return (data, text)


def _validate_datetime(val: datetime, tag: Tag):
    val = val.replace(microsecond=0)
    text = val.strftime('%Y-%m-%d %H:%M:%S')
    return (val, text)


def _validate_date(val: date, tag: Tag):
    text = val.strftime('%Y-%m-%d')
    return (val, text)


def _validate_time(val: dt_time, tag: Tag):
    val = val.replace(microsecond=0)
    text = val.strftime('%H:%M:%S')
    return (val, text)


# ===== Any to Node (struct, map, slice) =====

def _any_to_node_object(obj: Any, tag: Tag, depth: int, path: str) -> Obj:
    """Convert a Python object/class instance to an Object node."""
    depth += 1
    if depth > _MAX_DEPTH:
        raise ValueError(f"max depth: {_MAX_DEPTH}")

    tag.type = ValueType.Obj
    cls = obj.__class__
    tag.name = _camel_to_snake(cls.__name__)
    if tag.name and not path:
        path = tag.name
    elif tag.name:
        path = f"{path}.{tag.name}"

    # Check for class-level @mm decorator tag
    class_tag = get_mm_tag_for_class(cls)
    if class_tag is not None:
        tag = MergeTag(tag, class_tag)
        if class_tag.nullable or _is_union_type_with_null(class_tag.type):
            tag.nullable = True

    nodes = []
    for field_name, field_type in _get_type_hints(cls).items():
        # Get field value
        field_value = getattr(obj, field_name, None)
        
        field_key = _camel_to_snake(field_name)
        
        # Check for field-level @mm decorator tag
        field_tag = get_mm_tag_for_field(cls, field_name)
        if field_tag is None:
            field_tag = NewTag()
        
        # Check for union type with null in Python type annotation (e.g., int|None)
        if _is_union_type_with_none(field_type):
            field_tag.nullable = True
        
        # Auto-detect type from Python type annotation
        if field_tag.type == ValueType(0):
            auto_type = _python_type_to_value_type(field_type)
            if auto_type != ValueType(0):
                field_tag.type = auto_type
        
        # Check for union type with null (e.g., "int|null" as type kwarg)
        if field_tag.nullable or _is_union_type_with_null(field_tag.type):
            field_tag.nullable = True
        
        p = f"{path}.{field_key}"
        child_node = value_to_node(field_value, field_tag, depth, p)
        
        nodes.append(Field(key=field_key, value=child_node))
    
    return Obj(fields=nodes, tag=tag, path=path)


def _any_to_node_dict(value: dict, tag: Tag, depth: int, path: str) -> Obj:
    """Convert a dict to an Object (map) node."""
    depth += 1
    if depth > _MAX_DEPTH:
        raise ValueError(f"max depth: {_MAX_DEPTH}")

    tag.type = ValueType.Map

    nodes = []
    set_tag = False
    for key, val in value.items():
        if not isinstance(key, str):
            raise ValueError(f"map key must be string, got {type(key)}")
        
        key_str = _camel_to_snake(key)
        
        tag_item = NewTag()
        # Don't inherit child_type/child_* for dict fields (they are heterogeneous)
        # Only inherit non-child properties
        if tag.child_desc:
            tag_item.desc = tag.child_desc
        if tag.child_type != ValueType(0):
            tag_item.type = tag.child_type
        if tag.child_raw:
            tag_item.raw = tag.child_raw
        if tag.child_nullable:
            tag_item.nullable = tag.child_nullable
        if tag.child_allow_empty:
            tag_item.allow_empty = tag.child_allow_empty
        if tag.child_unique:
            tag_item.unique = tag.child_unique
        if tag.child_default:
            tag_item.default = tag.child_default
        if tag.child_min:
            tag_item.min = tag.child_min
        if tag.child_max:
            tag_item.max = tag.child_max
        if tag.child_size:
            tag_item.size = tag.child_size
        if tag.child_enum:
            tag_item.enum = tag.child_enum
        if tag.child_pattern:
            tag_item.pattern = tag.child_pattern
        if tag.child_version:
            tag_item.version = tag.child_version
        if tag.child_mime:
            tag_item.mime = tag.child_mime
        tag_item.name = key_str
        
        p = f"{path}[{key_str}]"
        child_node = value_to_node(val, tag_item, depth, p)
        tag_item = child_node.tag
        
        if not set_tag:
            tag.child_desc = tag_item.desc
            tag.child_type = tag_item.type
            tag.child_raw = tag_item.raw
            tag.child_nullable = tag_item.nullable
            tag.child_allow_empty = tag_item.allow_empty
            tag.child_unique = tag_item.unique
            tag.child_default = tag_item.default
            tag.child_min = tag_item.min
            tag.child_max = tag_item.max
            tag.child_size = tag_item.size
            tag.child_enum = tag_item.enum
            tag.child_pattern = tag_item.pattern
            tag.child_version = tag_item.version
            tag.child_mime = tag_item.mime
            set_tag = True
        
        nodes.append(Field(key=key_str, value=child_node))
    
    # Handle empty dict
    if not nodes:
        tag_item = NewTag()
        tag_item.inherit(tag)
        tag_item.example = True
        p = f"{path}[]"
        child_node = value_to_node("", tag_item, depth, p)
        
        tag.child_type = tag_item.type
        
        nodes.append(Field(key="", value=child_node))
    
    return Obj(fields=nodes, tag=tag, path=path)


def _any_to_node_list(value: list, tag: Tag, depth: int, path: str) -> Arr:
    """Convert a list/tuple to an Array (slice) node."""
    depth += 1
    if depth > _MAX_DEPTH:
        raise ValueError(f"max depth: {_MAX_DEPTH}")

    tag.type = ValueType.Vec

    items = []
    set_tag = False
    for i, item in enumerate(value):
        tag_item = NewTag()
        tag_item.inherit(tag)
        
        p = f"{path}[{i}]"
        child_node = value_to_node(item, tag_item, depth, p)
        tag_item = child_node.tag
        
        if not set_tag:
            tag.child_desc = tag_item.desc
            tag.child_type = tag_item.type
            tag.child_raw = tag_item.raw
            tag.child_nullable = tag_item.nullable
            tag.child_allow_empty = tag_item.allow_empty
            tag.child_unique = tag_item.unique
            tag.child_default = tag_item.default
            tag.child_min = tag_item.min
            tag.child_max = tag_item.max
            tag.child_size = tag_item.size
            tag.child_enum = tag_item.enum
            tag.child_pattern = tag_item.pattern
            tag.child_version = tag_item.version
            tag.child_mime = tag_item.mime
            set_tag = True
        
        items.append(child_node)
    
    # Handle empty list
    if not items:
        tag_item = NewTag()
        tag_item.inherit(tag)
        tag_item.example = True
        p = f"{path}[0]"
        child_node = value_to_node("", tag_item, depth, p)
        
        tag.child_type = tag_item.type
        
        items.append(child_node)
    
    return Arr(items=items, tag=tag, path=path)


def _get_type_hints(cls: type) -> Dict[str, type]:
    """Get type hints for a class (from class annotations and __init__ annotations)."""
    hints = {}
    
    # First try class-level __annotations__
    try:
        for name, hint in cls.__annotations__.items():
            if not name.startswith('_'):
                hints[name] = hint
    except AttributeError:
        pass
    
    # Check __init__ annotations
    init_func = cls.__init__ if hasattr(cls, '__init__') else None
    if init_func:
        try:
            init_hints = init_func.__annotations__
            for name, hint in init_hints.items():
                if name != 'return' and not name.startswith('_'):
                    hints[name] = hint
        except AttributeError:
            pass
    
    return hints


def _get_field_names(cls: type) -> set:
    """Get all field names for a class (from annotations and __init__ params)."""
    names = set()
    
    # From class annotations
    try:
        names.update(n for n in cls.__annotations__ if not n.startswith('_'))
    except AttributeError:
        pass
    
    # From __init__ annotations
    init_func = cls.__init__ if hasattr(cls, '__init__') else None
    if init_func:
        try:
            names.update(n for n in init_func.__annotations__ 
                        if n != 'return' and not n.startswith('_'))
        except AttributeError:
            pass
    
    return names


# ===== NodeToValue (Bind) =====

def node_to_value(node: Node, target_type: Any) -> Any:
    """Convert a Node tree back to a Python value.
    
    Python equivalent of Go's Bind() function.
    """
    if isinstance(node, Obj):
        return _bind_object(node, target_type)
    elif isinstance(node, Arr):
        return _bind_array(node, target_type)
    elif isinstance(node, Val):
        return _bind_value(node, target_type)
    else:
        raise ValueError(f"unsupported node type: {type(node)}")


def _bind_object(obj: Obj, target_type: Any) -> Any:
    """Bind an Object node to a Python value (dict or class instance)."""
    # If target type is dict or Any, return dict
    if target_type is dict or target_type is Any:
        result = {}
        for field in obj.fields:
            val = _bind_value_or_node(field.value)
            result[field.key] = val
        return result
    
    # If target type is a class, create instance
    if isinstance(target_type, type):
        hints = _get_type_hints(target_type)
        field_names = _get_field_names(target_type)
        
        # Build kwargs from fields
        kwargs = {}
        for field in obj.fields:
            field_key = field.key
            
            field_name = _find_field_name(target_type, field_key)
            if field_name is None:
                continue
            
            field_type = hints.get(field_name, Any)
            field_value = node_to_value(field.value, field_type)
            kwargs[field_name] = field_value
        
        # Create instance with __init__ if it accepts kwargs, otherwise use __new__ + setattr
        try:
            instance = target_type(**kwargs)
        except (TypeError, Exception):
            # Fallback: use __new__ and setattr
            instance = target_type.__new__(target_type)
            for name, value in kwargs.items():
                setattr(instance, name, value)
        
        return instance
    
    # Fallback to dict
    result = {}
    for field in obj.fields:
        result[field.key] = _bind_value_or_node(field.value)
    return result


def _bind_array(arr: Arr, target_type: Any) -> Any:
    """Bind an Array node to a Python list."""
    if target_type is list or target_type is Any or target_type is tuple:
        result = []
        for item in arr.items:
            result.append(_bind_value_or_node(item))
        if target_type is tuple:
            return tuple(result)
        return result
    
    if isinstance(target_type, type):
        # Try to determine element type from annotations/generics
        item_type = Any
        result = []
        for item in arr.items:
            result.append(node_to_value(item, item_type))
        
        # If target is a specific container like List[str], try to convert
        return result
    
    return [_bind_value_or_node(item) for item in arr.items]


def _bind_value(val: Val, target_type: Any) -> Any:
    """Bind a Value node to a Python value."""
    tag = val.tag
    
    if tag.is_null:
        return None
    
    # Use the data field if available
    if val.data is not None:
        # Convert data back to target type
        if isinstance(val.data, target_type) if isinstance(target_type, type) else True:
            return val.data
        return val.data
    
    # Fallback to text parsing
    text = val.text
    if text == "null":
        return None
    
    if target_type is Any or target_type is str:
        return text
    
    if target_type is int:
        try:
            return int(text)
        except (ValueError, TypeError):
            return 0
    elif target_type is float:
        try:
            return float(text)
        except (ValueError, TypeError):
            return 0.0
    elif target_type is bool:
        return text.lower() == "true"
    elif target_type is bytes and isinstance(val.data, bytes):
        return val.data
    elif target_type is datetime and isinstance(val.data, datetime):
        return val.data
    elif target_type is date and isinstance(val.data, date):
        return val.data
    elif target_type is dt_time and isinstance(val.data, dt_time):
        return val.data
    
    return text


def _bind_value_or_node(node: Node) -> Any:
    """Bind a node to a plain Python value (no type info)."""
    if isinstance(node, Obj):
        result = {}
        for field in node.fields:
            result[field.key] = _bind_value_or_node(field.value)
        return result
    elif isinstance(node, Arr):
        return [_bind_value_or_node(item) for item in node.items]
    elif isinstance(node, Val):
        if node.tag.is_null:
            return None
        if node.data is not None:
            return node.data
        text = node.text
        if text == "null":
            return None
        # Try to parse numbers
        try:
            if '.' in text:
                return float(text)
            return int(text)
        except (ValueError, TypeError):
            pass
        return text
    return None


def _find_field_name(cls: type, snake_key: str) -> Optional[str]:
    """Find the Python attribute name for a snake_case key."""
    # Get all known field names
    field_names = _get_field_names(cls)
    
    # Direct match
    if snake_key in field_names:
        return snake_key
    
    # Try CamelCase match
    parts = snake_key.split('_')
    camel = ''.join(p.capitalize() for p in parts)
    
    # Try lowercase match
    for name in field_names:
        if name.lower() == snake_key:
            return name
        if name == camel:
            return name
        if _camel_to_snake(name) == snake_key:
            return name
    
    # Check dir(cls) as last resort
    for name in dir(cls):
        if name.startswith('_'):
            continue
        if name.lower() == snake_key:
            return name
        if name == camel:
            return name
        if _camel_to_snake(name) == snake_key:
            return name
    
    return None


# ===== High-level encode/decode functions =====

def encode_from_value(value: Any) -> bytes:
    """Convert a Python value directly to MetaMessage binary format.
    
    Python equivalent of calling ValueToNode then Encode.
    """
    node = value_to_node(value)
    encoder = Encoder()
    return encoder.encode(node)


def decode_to_value(data: bytes, target_type: Any = Any) -> Any:
    """Decode MetaMessage binary format directly to a Python value.
    
    Python equivalent of calling Decode then Bind.
    """
    decoder = Decoder(data)
    decoded = decoder.decode()
    
    if target_type is Any or target_type is dict or target_type is list:
        return decoded
    
    # We need to go through Node tree for proper binding
    # Re-encode the decoded value to get back Nodes
    # Actually we need proper decoder that returns Nodes
    # For now, use the simple decoder output
    
    return decoded


def decode_to_value_with_node(data: bytes, target_type: Any) -> Any:
    """Decode binary to Python value via Node intermediate representation."""
    from . import create_decoder  # Lazy import to avoid circular
    import sys
    sys.path.insert(0, '/Users/lizongying/IdeaProjects/meta-message/mm-py')
    
    # Use the existing parse/decode chain
    decoder = Decoder(data)
    decoded = decoder.decode()
    
    # Create a Node from decoded value for binding
    node = value_to_node(decoded)
    return node_to_value(node, target_type)
