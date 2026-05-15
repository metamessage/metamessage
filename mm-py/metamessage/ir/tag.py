from typing import Any, Optional
from enum import IntEnum
from dataclasses import dataclass
import re


class ValueType(IntEnum):
    Unknown = 0
    Doc = 1
    Slice = 2
    Array = 3
    Object = 4
    Map = 5
    String = 6
    Bytes = 7
    Bool = 8
    Int = 9
    Int8 = 10
    Int16 = 11
    Int32 = 12
    Int64 = 13
    Uint = 14
    Uint8 = 15
    Uint16 = 16
    Uint32 = 17
    Uint64 = 18
    Float32 = 19
    Float64 = 20
    BigInt = 21
    DateTime = 22
    Date = 23
    Time = 24
    UUID = 25
    Decimal = 26
    IP = 27
    URL = 28
    Email = 29
    Enum = 30
    Image = 31
    Video = 32

    def to_str(self) -> str:
        mapping = {
            ValueType.Unknown: "unknown",
            ValueType.Doc: "doc",
            ValueType.Array: "arr",
            ValueType.Slice: "slice",
            ValueType.Object: "obj",
            ValueType.Map: "map",
            ValueType.String: "str",
            ValueType.Bytes: "bytes",
            ValueType.Bool: "bool",
            ValueType.Int: "i",
            ValueType.Int8: "i8",
            ValueType.Int16: "i16",
            ValueType.Int32: "i32",
            ValueType.Int64: "i64",
            ValueType.Uint: "u",
            ValueType.Uint8: "u8",
            ValueType.Uint16: "u16",
            ValueType.Uint32: "u32",
            ValueType.Uint64: "u64",
            ValueType.Float32: "f32",
            ValueType.Float64: "f64",
            ValueType.BigInt: "bi",
            ValueType.DateTime: "datetime",
            ValueType.Date: "date",
            ValueType.Time: "time",
            ValueType.UUID: "uuid",
            ValueType.Decimal: "decimal",
            ValueType.IP: "ip",
            ValueType.URL: "url",
            ValueType.Email: "email",
            ValueType.Enum: "enum",
            ValueType.Image: "image",
            ValueType.Video: "video",
        }
        return mapping.get(self, str(self.value))


_str_to_value_type = {
    "unknown": ValueType.Unknown, "doc": ValueType.Doc, "arr": ValueType.Array,
    "slice": ValueType.Slice, "obj": ValueType.Object, "map": ValueType.Map,
    "str": ValueType.String, "bytes": ValueType.Bytes, "bool": ValueType.Bool,
    "i": ValueType.Int, "i8": ValueType.Int8, "i16": ValueType.Int16, "i32": ValueType.Int32, "i64": ValueType.Int64,
    "u": ValueType.Uint, "u8": ValueType.Uint8, "u16": ValueType.Uint16, "u32": ValueType.Uint32, "u64": ValueType.Uint64,
    "f32": ValueType.Float32, "f64": ValueType.Float64,
    "bi": ValueType.BigInt, "datetime": ValueType.DateTime, "date": ValueType.Date, "time": ValueType.Time,
    "uuid": ValueType.UUID, "decimal": ValueType.Decimal, "ip": ValueType.IP, "url": ValueType.URL,
    "email": ValueType.Email, "enum": ValueType.Enum,
    "image": ValueType.Image, "video": ValueType.Video,
}


def parse_value_type(s: str) -> ValueType:
    s = s.lower().strip()
    if s in _str_to_value_type:
        return _str_to_value_type[s]
    return ValueType.Unknown


class TagKey(IntEnum):
    IsNull = 0 << 3
    Example = 1 << 3
    Desc = 2 << 3
    Type = 3 << 3
    Raw = 4 << 3
    Nullable = 5 << 3
    AllowEmpty = 6 << 3
    Unique = 7 << 3
    Default = 8 << 3
    Min = 9 << 3
    Max = 10 << 3
    Size = 11 << 3
    Enum = 12 << 3
    Pattern = 13 << 3
    Location = 14 << 3
    Version = 15 << 3
    Mime = 16 << 3
    ChildDesc = 17 << 3
    ChildType = 18 << 3
    ChildRaw = 19 << 3
    ChildNullable = 20 << 3
    ChildAllowEmpty = 21 << 3
    ChildUnique = 22 << 3
    ChildDefault = 23 << 3
    ChildMin = 24 << 3
    ChildMax = 25 << 3
    ChildSize = 26 << 3
    ChildEnum = 27 << 3
    ChildPattern = 28 << 3
    ChildLocation = 29 << 3
    ChildVersion = 30 << 3
    ChildMime = 31 << 3


_MAX1BYTE = 0xFF
_MAX2BYTE = 0xFFFF
_MAX3BYTE = 0xFFFFFF
_MAX4BYTE = 0xFFFFFFFF
_MAX5BYTE = 0xFFFFFFFFFF
_MAX6BYTE = 0xFFFFFFFFFFFF
_MAX7BYTE = 0xFFFFFFFFFFFFFF
_MAX8BYTE = 0xFFFFFFFFFFFFFFFF

_IntLenMask = 0b11111
_IntLen1Byte = _IntLenMask - 7
_IntLen2Byte = _IntLenMask - 6
_IntLen3Byte = _IntLenMask - 5
_IntLen4Byte = _IntLenMask - 4
_IntLen5Byte = _IntLenMask - 3
_IntLen6Byte = _IntLenMask - 2
_IntLen7Byte = _IntLenMask - 1
_IntLen8Byte = _IntLenMask


def _encode_uint64(buf: bytearray, key: int, uv: int):
    if uv <= _MAX1BYTE:
        buf.extend([key, uv])
    elif uv <= _MAX2BYTE:
        buf.extend([key | 1, uv >> 8, uv & 0xFF])
    elif uv <= _MAX3BYTE:
        buf.extend([key | 2, uv >> 16, uv >> 8, uv & 0xFF])
    elif uv <= _MAX4BYTE:
        buf.extend([key | 3, uv >> 24, uv >> 16, uv >> 8, uv & 0xFF])
    elif uv <= _MAX5BYTE:
        buf.extend([key | 4, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv & 0xFF])
    elif uv <= _MAX6BYTE:
        buf.extend([key | 5, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv & 0xFF])
    elif uv <= _MAX7BYTE:
        buf.extend([key | 6, uv >> 48, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv & 0xFF])
    elif uv <= _MAX8BYTE:
        buf.extend([key | 7, uv >> 56, uv >> 48, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv & 0xFF])


def NewTag() -> 'Tag':
    return Tag()


@dataclass
class Tag:
    name: str = ""
    is_null: bool = False
    example: bool = False
    desc: str = ""
    type: ValueType = ValueType.Unknown
    raw: bool = False
    nullable: bool = False
    allow_empty: bool = False
    unique: bool = False
    default: str = ""
    min: str = ""
    max: str = ""
    size: int = 0
    enum: str = ""
    pattern: str = ""
    location: Any = None
    version: int = 0
    mime: str = ""
    child_desc: str = ""
    child_type: ValueType = ValueType.Unknown
    child_raw: bool = False
    child_nullable: bool = False
    child_allow_empty: bool = False
    child_unique: bool = False
    child_default: str = ""
    child_min: str = ""
    child_max: str = ""
    child_size: int = 0
    child_enum: str = ""
    child_pattern: str = ""
    child_location: Any = None
    child_version: int = 0
    child_mime: str = ""
    is_inherit: bool = False

    def inherit(self, tag: 'Tag'):
        if tag is None:
            return
        self.is_inherit = True
        if tag.child_desc:
            self.desc = tag.child_desc
        if tag.child_type != ValueType.Unknown:
            self.type = tag.child_type
        if tag.child_raw:
            self.raw = tag.child_raw
        if tag.child_nullable:
            self.nullable = tag.child_nullable
        if tag.child_allow_empty:
            self.allow_empty = tag.child_allow_empty
        if tag.child_unique:
            self.unique = tag.child_unique
        if tag.child_default:
            self.default = tag.child_default
        if tag.child_min:
            self.min = tag.child_min
        if tag.child_max:
            self.max = tag.child_max
        if tag.child_size:
            self.size = tag.child_size
        if tag.child_enum:
            self.enum = tag.child_enum
        if tag.child_pattern:
            self.pattern = tag.child_pattern
        if tag.child_location is not None:
            self.location = tag.child_location
        if tag.child_version:
            self.version = tag.child_version
        if tag.child_mime:
            self.mime = tag.child_mime

    def bytes(self) -> bytes:
        buf = bytearray()

        if self.example:
            buf.append(TagKey.Example | 1)

        if self.is_null:
            buf.append(TagKey.IsNull | 1)

        if self.nullable and not self.is_inherit:
            if not self.is_null:
                buf.append(TagKey.Nullable | 1)

        if self.desc and not self.is_inherit:
            desc_bytes = self.desc.encode('utf-8')
            l = len(desc_bytes)
            if l <= 5:
                buf.append(TagKey.Desc | l)
                buf.extend(desc_bytes)
            elif l <= 256:
                buf.append(TagKey.Desc | 6)
                buf.append(l)
                buf.extend(desc_bytes)
            elif l <= 65536:
                buf.append(TagKey.Desc | 7)
                buf.extend([l >> 8, l & 0xFF])
                buf.extend(desc_bytes)

        if self.type != ValueType.Unknown and not self.is_inherit:
            if self.type not in (ValueType.String, ValueType.Bytes, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Object, ValueType.Slice):
                if not (self.type == ValueType.Array and self.size > 0) and not (self.type == ValueType.Enum and self.enum):
                    buf.append(TagKey.Type)
                    buf.append(self.type)

        if self.raw and not self.is_inherit:
            buf.append(TagKey.Raw | 1)

        if self.allow_empty and not self.is_inherit:
            buf.append(TagKey.AllowEmpty | 1)

        if self.unique and not self.is_inherit:
            buf.append(TagKey.Unique | 1)

        if self.default and not self.is_inherit:
            l = len(self.default)
            if l < 7:
                buf.append(TagKey.Default | l)
            else:
                buf.append(TagKey.Default | 7)
                buf.append(l)
            buf.extend(self.default.encode('utf-8'))

        if self.min and not self.is_inherit:
            l = len(self.min)
            if l < 7:
                buf.append(TagKey.Min | l)
            else:
                buf.append(TagKey.Min | 7)
                buf.append(l)
            buf.extend(self.min.encode('utf-8'))

        if self.max and not self.is_inherit:
            l = len(self.max)
            if l < 7:
                buf.append(TagKey.Max | l)
            else:
                buf.append(TagKey.Max | 7)
                buf.append(l)
            buf.extend(self.max.encode('utf-8'))

        if self.size and not self.is_inherit:
            _encode_uint64(buf, TagKey.Size, self.size)

        if self.enum and not self.is_inherit:
            enum_bytes = self.enum.encode('utf-8')
            l = len(enum_bytes)
            if l <= 5:
                buf.append(TagKey.Enum | l)
                buf.extend(enum_bytes)
            elif l <= 256:
                buf.append(TagKey.Enum | 6)
                buf.append(l)
                buf.extend(enum_bytes)
            elif l <= 65536:
                buf.append(TagKey.Enum | 7)
                buf.extend([l >> 8, l & 0xFF])
                buf.extend(enum_bytes)

        if self.pattern and not self.is_inherit:
            l = len(self.pattern)
            if l < 7:
                buf.append(TagKey.Pattern | l)
            else:
                buf.append(TagKey.Pattern | 7)
                buf.append(l)
            buf.extend(self.pattern.encode('utf-8'))

        location_offset_hour = 0
        if self.location is not None:
            try:
                location_offset_hour = int(self.location)
            except (ValueError, TypeError):
                pass
        if location_offset_hour != 0 and not self.is_inherit:
            loc_str = str(location_offset_hour)
            buf.append(TagKey.Location | len(loc_str))
            buf.extend(loc_str.encode('utf-8'))

        if self.version and not self.is_inherit:
            _encode_uint64(buf, TagKey.Version, self.version)

        if self.mime and not self.is_inherit:
            l = len(self.mime)
            if l < 7:
                buf.append(TagKey.Mime | l)
            else:
                buf.append(TagKey.Mime | 7)
                buf.append(l)

        if self.child_desc:
            child_desc_bytes = self.child_desc.encode('utf-8')
            l = len(child_desc_bytes)
            if l <= 5:
                buf.append(TagKey.ChildDesc | l)
                buf.extend(child_desc_bytes)
            elif l <= 256:
                buf.append(TagKey.ChildDesc | 6)
                buf.append(l)
                buf.extend(child_desc_bytes)
            elif l <= 65536:
                buf.append(TagKey.ChildDesc | 7)
                buf.extend([l >> 8, l & 0xFF])
                buf.extend(child_desc_bytes)

        if self.child_type != ValueType.Unknown:
            if self.child_type not in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Object, ValueType.Slice):
                if not (self.child_type == ValueType.Array and self.child_size > 0) and not (self.child_type == ValueType.Enum and self.child_enum):
                    buf.append(TagKey.ChildType)
                    buf.append(self.child_type)

        if self.child_raw:
            buf.append(TagKey.ChildRaw | 1)

        if self.child_nullable:
            buf.append(TagKey.ChildNullable | 1)

        if self.child_allow_empty:
            buf.append(TagKey.ChildAllowEmpty | 1)

        if self.child_unique:
            buf.append(TagKey.ChildUnique | 1)

        if self.child_default:
            l = len(self.child_default)
            if l < 7:
                buf.append(TagKey.ChildDefault | l)
            else:
                buf.append(TagKey.ChildDefault | 7)
                buf.append(l)
            buf.extend(self.child_default.encode('utf-8'))

        if self.child_min:
            l = len(self.child_min)
            if l < 7:
                buf.append(TagKey.ChildMin | l)
            else:
                buf.append(TagKey.ChildMin | 7)
                buf.append(l)
            buf.extend(self.child_min.encode('utf-8'))

        if self.child_max:
            l = len(self.child_max)
            if l < 7:
                buf.append(TagKey.ChildMax | l)
            else:
                buf.append(TagKey.ChildMax | 7)
                buf.append(l)
            buf.extend(self.child_max.encode('utf-8'))

        if self.child_size:
            _encode_uint64(buf, TagKey.ChildSize, self.child_size)

        if self.child_enum:
            child_enum_bytes = self.child_enum.encode('utf-8')
            l = len(child_enum_bytes)
            if l <= 5:
                buf.append(TagKey.ChildEnum | l)
                buf.extend(child_enum_bytes)
            elif l <= 256:
                buf.append(TagKey.ChildEnum | 6)
                buf.append(l)
                buf.extend(child_enum_bytes)
            elif l <= 65536:
                buf.append(TagKey.ChildEnum | 7)
                buf.extend([l >> 8, l & 0xFF])
                buf.extend(child_enum_bytes)

        if self.child_pattern:
            l = len(self.child_pattern)
            if l < 7:
                buf.append(TagKey.ChildPattern | l)
            else:
                buf.append(TagKey.ChildPattern | 7)
                buf.append(l)
            buf.extend(self.child_pattern.encode('utf-8'))

        child_location_offset_hour = 0
        if self.child_location is not None:
            try:
                child_location_offset_hour = int(self.child_location)
            except (ValueError, TypeError):
                pass
        if child_location_offset_hour != 0:
            loc_str = str(child_location_offset_hour)
            buf.append(TagKey.ChildLocation | len(loc_str))
            buf.extend(loc_str.encode('utf-8'))

        if self.child_version:
            _encode_uint64(buf, TagKey.ChildVersion, self.child_version)

        if self.child_mime:
            l = len(self.child_mime)
            if l < 7:
                buf.append(TagKey.ChildMime | l)
            else:
                buf.append(TagKey.ChildMime | 7)
                buf.append(l)

        return bytes(buf)

    def __str__(self) -> str:
        parts = []
        
        if self.type != ValueType.Unknown and not self.is_inherit:
            if self.type in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Object, ValueType.Slice):
                pass
            else:
                if not (self.type == ValueType.Array and self.size > 0) and not (self.type == ValueType.Enum and self.enum):
                    parts.append(f"type={self.type.to_str()}")

        if self.example:
            parts.append("example")
        
        if self.is_null:
            parts.append("is_null")
        
        if self.nullable and not self.is_inherit:
            if not self.is_null:
                parts.append("nullable")
        
        if self.desc and not self.is_inherit:
            parts.append(f'desc="{self.desc}"')
        
        if self.raw and not self.is_inherit:
            parts.append("raw")
        
        if self.allow_empty and not self.is_inherit:
            parts.append("allow_empty")
        
        if self.unique and not self.is_inherit:
            parts.append("unique")
        
        if self.default and not self.is_inherit:
            parts.append(f"default={self.default}")
        
        if self.min and not self.is_inherit:
            parts.append(f"min={self.min}")
        
        if self.max and not self.is_inherit:
            parts.append(f"max={self.max}")
        
        if self.size and not self.is_inherit:
            parts.append(f"size={self.size}")
        
        if self.enum and not self.is_inherit:
            parts.append(f"enum={self.enum}")
        
        if self.pattern and not self.is_inherit:
            parts.append(f"pattern={self.pattern}")
        
        location_offset_hour = 0
        if self.location is not None:
            try:
                location_offset_hour = int(self.location)
            except (ValueError, TypeError):
                pass
        if location_offset_hour != 0 and not self.is_inherit:
            parts.append(f"location={location_offset_hour}")

        if self.version and not self.is_inherit:
            parts.append(f"version={self.version}")
        
        if self.mime and not self.is_inherit:
            parts.append(f"mime={self.mime}")
        
        if self.child_desc:
            parts.append(f'child_desc="{self.child_desc}"')

        if self.child_type != ValueType.Unknown:
            if self.child_type not in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Object, ValueType.Slice):
                if not (self.child_type == ValueType.Array and self.child_size > 0) and not (self.child_type == ValueType.Enum and self.child_enum):
                    parts.append(f"child_type={self.child_type.to_str()}")

        if self.child_raw:
            parts.append("child_raw")
        
        if self.child_nullable:
            parts.append("child_nullable")
        
        if self.child_allow_empty:
            parts.append("child_allow_empty")
        
        if self.child_unique:
            parts.append("child_unique")
        
        if self.child_default:
            parts.append(f"child_default={self.child_default}")
        
        if self.child_min:
            parts.append(f"child_min={self.child_min}")
        
        if self.child_max:
            parts.append(f"child_max={self.child_max}")
        
        if self.child_size:
            parts.append(f"child_size={self.child_size}")
        
        if self.child_enum:
            parts.append(f"child_enum={self.child_enum}")
        
        if self.child_pattern:
            parts.append(f"child_pattern={self.child_pattern}")
        
        child_location_offset_hour = 0
        if self.child_location is not None:
            try:
                child_location_offset_hour = int(self.child_location)
            except (ValueError, TypeError):
                pass
        if child_location_offset_hour != 0:
            parts.append(f"child_location={child_location_offset_hour}")

        if self.child_version:
            parts.append(f"child_version={self.child_version}")
        
        if self.child_mime:
            parts.append(f"child_mime={self.child_mime}")
        
        return "; ".join(parts)


TAG_KEY_MAP = {
    "is_null": "is_null",
    "example": "example",
    "desc": "desc",
    "type": "type",
    "raw": "raw",
    "nullable": "nullable",
    "allow_empty": "allow_empty",
    "unique": "unique",
    "default": "default",
    "min": "min",
    "max": "max",
    "size": "size",
    "enum": "enum",
    "pattern": "pattern",
    "location": "location",
    "version": "version",
    "mime": "mime",
    "child_desc": "child_desc",
    "child_type": "child_type",
    "child_raw": "child_raw",
    "child_nullable": "child_nullable",
    "child_allow_empty": "child_allow_empty",
    "child_unique": "child_unique",
    "child_default": "child_default",
    "child_min": "child_min",
    "child_max": "child_max",
    "child_size": "child_size",
    "child_enum": "child_enum",
    "child_pattern": "child_pattern",
    "child_location": "child_location",
    "child_version": "child_version",
    "child_mime": "child_mime",
}


def _split_tag(tag_str: str) -> list:
    parts = []
    current = ""
    in_quotes = False
    for c in tag_str:
        if c == '"':
            in_quotes = not in_quotes
            current += c
        elif c == ';' and not in_quotes:
            parts.append(current.strip())
            current = ""
        else:
            current += c
    if current.strip():
        parts.append(current.strip())
    return parts


def mm_tag(tag_str: str) -> Tag:
    if not tag_str:
        return Tag()
    
    tag_str = tag_str.strip()
    tag_str = tag_str.lstrip("//")
    tag_str = tag_str.strip()
    tag_str = tag_str.lstrip("mm:")
    tag_str = tag_str.strip()
    if not tag_str:
        return Tag()
    
    tag = Tag()
    
    parts = _split_tag(tag_str)
    for p in parts:
        if not p:
            continue
        
        k, v = p, ""
        if "=" in p:
            kv = p.split("=", 1)
            k = kv[0].strip()
            v = kv[1].strip() if len(kv) > 1 else ""
        
        k = k.lower()
        # strip quotes from value
        if len(v) >= 2 and v.startswith('"') and v.endswith('"'):
            v = v[1:-1]

        if k == "is_null":
            tag.is_null = True
            tag.nullable = True
        elif k == "example":
            tag.example = True
        elif k == "desc":
            tag.desc = v
        elif k == "type":
            tag.type = parse_value_type(v)
        elif k == "raw":
            tag.raw = True
        elif k == "nullable":
            tag.nullable = True
        elif k == "allow_empty":
            tag.allow_empty = True
        elif k == "unique":
            tag.unique = True
        elif k == "default":
            tag.default = v
        elif k == "min":
            tag.min = v
        elif k == "max":
            tag.max = v
        elif k == "size":
            try:
                tag.size = int(v)
            except ValueError:
                pass
        elif k == "enum":
            tag.type = ValueType.Enum
            tag.enum = v
        elif k == "pattern":
            tag.pattern = v
        elif k == "location":
            try:
                tag.location = int(v)
            except ValueError:
                tag.location = v
        elif k == "version":
            try:
                tag.version = int(v)
            except ValueError:
                pass
        elif k == "mime":
            tag.mime = v
        elif k == "child_desc":
            tag.child_desc = v
        elif k == "child_type":
            tag.child_type = parse_value_type(v)
        elif k == "child_raw":
            tag.child_raw = True
        elif k == "child_nullable":
            tag.child_nullable = True
        elif k == "child_allow_empty":
            tag.child_allow_empty = True
        elif k == "child_unique":
            tag.child_unique = True
        elif k == "child_default":
            tag.child_default = v
        elif k == "child_min":
            tag.child_min = v
        elif k == "child_max":
            tag.child_max = v
        elif k == "child_size":
            try:
                tag.child_size = int(v)
            except ValueError:
                pass
        elif k == "child_enum":
            tag.child_enum = v
            tag.child_type = ValueType.Enum
        elif k == "child_pattern":
            tag.child_pattern = v
        elif k == "child_location":
            try:
                tag.child_location = int(v)
            except ValueError:
                tag.child_location = v
        elif k == "child_version":
            try:
                tag.child_version = int(v)
            except ValueError:
                pass
        elif k == "child_mime":
            tag.child_mime = v
    
    return tag


def def_tag(**kwargs) -> Tag:
    return Tag(**kwargs)


def MergeTag(dst: Tag, src: Tag) -> Tag:
    if src is None:
        return dst
    if dst is None:
        return src

    if src.is_null:
        dst.is_null = src.is_null
    if src.example:
        dst.example = src.example
    if src.desc:
        dst.desc = src.desc
    if src.type != ValueType.Unknown:
        dst.type = src.type
    if src.raw:
        dst.raw = True
    if src.nullable:
        dst.nullable = True
    if src.allow_empty:
        dst.allow_empty = True
    if src.unique:
        dst.unique = True
    if src.default:
        dst.default = src.default
    if src.min:
        dst.min = src.min
    if src.max:
        dst.max = src.max
    if src.size:
        dst.size = src.size
    if src.enum:
        dst.enum = src.enum
    if src.pattern:
        dst.pattern = src.pattern
    if src.location is not None:
        dst.location = src.location
    if src.version:
        dst.version = src.version
    if src.mime:
        dst.mime = src.mime
    if src.child_desc:
        dst.child_desc = src.child_desc
    if src.child_type != ValueType.Unknown:
        dst.child_type = src.child_type
    if src.child_raw:
        dst.child_raw = True
    if src.child_nullable:
        dst.child_nullable = True
    if src.child_allow_empty:
        dst.child_allow_empty = True
    if src.child_unique:
        dst.child_unique = True
    if src.child_default:
        dst.child_default = src.child_default
    if src.child_min:
        dst.child_min = src.child_min
    if src.child_max:
        dst.child_max = src.child_max
    if src.child_size:
        dst.child_size = src.child_size
    if src.child_enum:
        dst.child_enum = src.child_enum
    if src.child_pattern:
        dst.child_pattern = src.child_pattern
    if src.child_location is not None:
        dst.child_location = src.child_location
    if src.child_version:
        dst.child_version = src.child_version
    if src.child_mime:
        dst.child_mime = src.child_mime

    return dst
