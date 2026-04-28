from typing import Any, Optional
from enum import IntEnum
from dataclasses import dataclass


class ValueType(IntEnum):
    Unknown = 0
    Doc = 1
    Slice = 2
    Array = 3
    Struct = 4
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
        buf.extend([key | 1, uv >> 8, uv])
    elif uv <= _MAX3BYTE:
        buf.extend([key | 2, uv >> 16, uv >> 8, uv])
    elif uv <= _MAX4BYTE:
        buf.extend([key | 3, uv >> 24, uv >> 16, uv >> 8, uv])
    elif uv <= _MAX5BYTE:
        buf.extend([key | 4, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv])
    elif uv <= _MAX6BYTE:
        buf.extend([key | 5, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv])
    elif uv <= _MAX7BYTE:
        buf.extend([key | 6, uv >> 48, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv])
    elif uv <= _MAX8BYTE:
        buf.extend([key | 7, uv >> 56, uv >> 48, uv >> 40, uv >> 32, uv >> 24, uv >> 16, uv >> 8, uv])


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
        self.desc = tag.child_desc
        self.type = tag.child_type
        self.raw = tag.child_raw
        self.nullable = tag.child_nullable
        self.allow_empty = tag.child_allow_empty
        self.unique = tag.child_unique
        self.default = tag.child_default
        self.min = tag.child_min
        self.max = tag.child_max
        self.size = tag.child_size
        self.enum = tag.child_enum
        self.pattern = tag.child_pattern
        self.location = tag.child_location
        self.version = tag.child_version
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
            l = len(self.desc)
            if l <= 5:
                buf.append(TagKey.Desc | l)
                buf.extend(self.desc.encode())
            elif l <= 256:
                buf.append(TagKey.Desc | 6)
                buf.append(l)
                buf.extend(self.desc.encode())
            elif l <= 65536:
                buf.append(TagKey.Desc | 7)
                buf.extend([l >> 8, l])
                buf.extend(self.desc.encode())

        if self.type != ValueType.Unknown and not self.is_inherit:
            if self.type not in (ValueType.String, ValueType.Bytes, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Struct, ValueType.Slice):
                if not (self.type == ValueType.Array and self.size > 0 or self.type == ValueType.Enum and self.enum):
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
            buf.extend(self.default.encode())

        if self.min and not self.is_inherit:
            l = len(self.min)
            if l < 7:
                buf.append(TagKey.Min | l)
            else:
                buf.append(TagKey.Min | 7)
                buf.append(l)
            buf.extend(self.min.encode())

        if self.max and not self.is_inherit:
            l = len(self.max)
            if l < 7:
                buf.append(TagKey.Max | l)
            else:
                buf.append(TagKey.Max | 7)
                buf.append(l)
            buf.extend(self.max.encode())

        if self.size and not self.is_inherit:
            _encode_uint64(buf, TagKey.Size, self.size)

        if self.enum and not self.is_inherit:
            l = len(self.enum)
            if l <= 5:
                buf.append(TagKey.Enum | l)
                buf.extend(self.enum.encode())
            elif l <= 256:
                buf.append(TagKey.Enum | 6)
                buf.append(l)
                buf.extend(self.enum.encode())
            elif l <= 65536:
                buf.append(TagKey.Enum | 7)
                buf.extend([l >> 8, l])
                buf.extend(self.enum.encode())

        if self.pattern and not self.is_inherit:
            l = len(self.pattern)
            if l < 7:
                buf.append(TagKey.Pattern | l)
            else:
                buf.append(TagKey.Pattern | 7)
                buf.append(l)
            buf.extend(self.pattern.encode())

        if self.location is not None and not self.is_inherit:
            loc_str = str(self.location)
            buf.append(TagKey.Location | len(loc_str))
            buf.extend(loc_str.encode())

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
            l = len(self.child_desc)
            if l <= 5:
                buf.append(TagKey.ChildDesc | l)
                buf.extend(self.child_desc.encode())
            elif l <= 256:
                buf.append(TagKey.ChildDesc | 6)
                buf.append(l)
                buf.extend(self.child_desc.encode())
            elif l <= 65536:
                buf.append(TagKey.ChildDesc | 7)
                buf.extend([l >> 8, l])
                buf.extend(self.child_desc.encode())

        if self.child_type != ValueType.Unknown:
            if self.child_type not in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Struct, ValueType.Slice):
                if not (self.child_type == ValueType.Array and self.child_size > 0 or self.child_type == ValueType.Enum and self.child_enum):
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
            buf.extend(self.child_default.encode())

        if self.child_min:
            l = len(self.child_min)
            if l < 7:
                buf.append(TagKey.ChildMin | l)
            else:
                buf.append(TagKey.ChildMin | 7)
                buf.append(l)
            buf.extend(self.child_min.encode())

        if self.child_max:
            l = len(self.child_max)
            if l < 7:
                buf.append(TagKey.ChildMax | l)
            else:
                buf.append(TagKey.ChildMax | 7)
                buf.append(l)
            buf.extend(self.child_max.encode())

        if self.child_size:
            _encode_uint64(buf, TagKey.ChildSize, self.child_size)

        if self.child_enum:
            l = len(self.child_enum)
            if l <= 5:
                buf.append(TagKey.ChildEnum | l)
                buf.extend(self.child_enum.encode())
            elif l <= 256:
                buf.append(TagKey.ChildEnum | 6)
                buf.append(l)
                buf.extend(self.child_enum.encode())
            elif l <= 65536:
                buf.append(TagKey.ChildEnum | 7)
                buf.extend([l >> 8, l])
                buf.extend(self.child_enum.encode())

        if self.child_pattern:
            l = len(self.child_pattern)
            if l < 7:
                buf.append(TagKey.ChildPattern | l)
            else:
                buf.append(TagKey.ChildPattern | 7)
                buf.append(l)
            buf.extend(self.child_pattern.encode())

        if self.child_location is not None:
            loc_str = str(self.child_location)
            buf.append(TagKey.ChildLocation | len(loc_str))
            buf.extend(loc_str.encode())

        if self.child_version:
            _encode_uint64(buf, TagKey.ChildVersion, self.child_version)

        if self.child_mime:
            l = len(self.child_mime)
            if l < 7:
                buf.append(TagKey.ChildMime | l)
            else:
                buf.append(TagKey.ChildMime | 7)
                buf.append(l)
            buf.extend(self.child_mime.encode())

        return bytes(buf)

    def __str__(self) -> str:
        parts = []
        
        if self.type != ValueType.Unknown and not self.is_inherit:
            if self.type in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Struct, ValueType.Slice):
                pass
            else:
                if not (self.type == ValueType.Array and self.size > 0 or self.type == ValueType.Enum and self.enum):
                    parts.append(f"type={self.type.name}")
        
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
        
        if self.version and not self.is_inherit:
            parts.append(f"version={self.version}")
        
        if self.mime and not self.is_inherit:
            parts.append(f"mime={self.mime}")
        
        if self.child_desc:
            parts.append(f"child_desc=\"{self.child_desc}")
        
        if self.child_type != ValueType.Unknown:
            if self.child_type in (ValueType.String, ValueType.Int, ValueType.Float64, ValueType.Bool, ValueType.Struct, ValueType.Slice):
                pass
            else:
                if not (self.child_type == ValueType.Array and self.child_size > 0 or self.child_type == ValueType.Enum and self.child_enum):
                    parts.append(f"child_type={self.child_type.name}")
        
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
        
        if self.child_version:
            parts.append(f"child_version={self.child_version}")
        
        if self.child_mime:
            parts.append(f"child_mime={self.child_mime}")
        
        return "; ".join(parts)


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
    
    parts = [p.strip() for p in tag_str.split(";")]
    type_map = {
        "unknown": ValueType.Unknown, "str": ValueType.String, "i": ValueType.Int,
        "i8": ValueType.Int8, "i16": ValueType.Int16, "i32": ValueType.Int32, "i64": ValueType.Int64,
        "u": ValueType.Uint, "u8": ValueType.Uint8, "u16": ValueType.Uint16, "u32": ValueType.Uint32, "u64": ValueType.Uint64,
        "f32": ValueType.Float32, "f64": ValueType.Float64, "bool": ValueType.Bool, "bytes": ValueType.Bytes,
        "bi": ValueType.BigInt, "datetime": ValueType.DateTime, "date": ValueType.Date, "time": ValueType.Time,
        "uuid": ValueType.UUID, "decimal": ValueType.Decimal, "ip": ValueType.IP, "url": ValueType.URL,
        "email": ValueType.Email, "enum": ValueType.Enum, "arr": ValueType.Array, "obj": ValueType.Struct,
    }
    
    for p in parts:
        if not p:
            continue
        
        k, v = p, ""
        if "=" in p:
            kv = p.split("=", 1)
            k = kv[0].strip()
            v = kv[1].strip() if len(kv) > 1 else ""
        
        k = k.lower()
        
        if k == "is_null":
            tag.is_null = True
            tag.nullable = True
        elif k == "example":
            tag.example = True
        elif k == "desc":
            tag.desc = v
        elif k == "type":
            tag.type = type_map.get(v, ValueType.Unknown)
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
            tag.size = int(v)
        elif k == "enum":
            tag.type = ValueType.Enum
            tag.enum = v
        elif k == "pattern":
            tag.pattern = v
        elif k == "location":
            tag.location = v
        elif k == "version":
            tag.version = int(v)
        elif k == "mime":
            tag.mime = v
        elif k == "child_desc":
            tag.child_desc = v
        elif k == "child_type":
            tag.child_type = type_map.get(v, ValueType.Unknown)
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
            tag.child_size = int(v)
        elif k == "child_enum":
            tag.child_enum = v
            tag.child_type = ValueType.Enum
        elif k == "child_pattern":
            tag.child_pattern = v
        elif k == "child_location":
            tag.child_location = v
        elif k == "child_version":
            tag.child_version = int(v)
        elif k == "child_mime":
            tag.child_mime = v
    
    return tag


def def_tag(**kwargs):
    return Tag(**kwargs)