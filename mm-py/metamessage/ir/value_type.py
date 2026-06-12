"""
ValueType enum for MetaMessage.
"""
from enum import IntEnum


class ValueType(IntEnum):
    Unknown = 0
    Doc = 1
    Vec = 2
    Arr = 3
    Obj = 4
    Map = 5
    Str = 6
    Bytes = 7
    Bool = 8
    I = 9
    I8 = 10
    I16 = 11
    I32 = 12
    I64 = 13
    U = 14
    U8 = 15
    U16 = 16
    U32 = 17
    U64 = 18
    F32 = 19
    F64 = 20
    Bigint = 21
    Datetime = 22
    Date = 23
    Time = 24
    Uuid = 25
    Decimal = 26
    Ip = 27
    Url = 28
    Email = 29
    Enums = 30
    Media = 31

    def __str__(self) -> str:
        mapping = {
            ValueType.Unknown: "unknown",
            ValueType.Doc: "doc",
            ValueType.Arr: "arr",
            ValueType.Vec: "vec",
            ValueType.Obj: "obj",
            ValueType.Map: "map",
            ValueType.Str: "str",
            ValueType.Bytes: "bytes",
            ValueType.Bool: "bool",
            ValueType.I: "i",
            ValueType.I8: "i8",
            ValueType.I16: "i16",
            ValueType.I32: "i32",
            ValueType.I64: "i64",
            ValueType.U: "u",
            ValueType.U8: "u8",
            ValueType.U16: "u16",
            ValueType.U32: "u32",
            ValueType.U64: "u64",
            ValueType.F32: "f32",
            ValueType.F64: "f64",
            ValueType.Bigint: "bigint",
            ValueType.Datetime: "datetime",
            ValueType.Date: "date",
            ValueType.Time: "time",
            ValueType.Uuid: "uuid",
            ValueType.Decimal: "decimal",
            ValueType.Ip: "ip",
            ValueType.Url: "url",
            ValueType.Email: "email",
            ValueType.Enums: "enums",
            ValueType.Media: "media",
        }
        return mapping.get(self, "ValueType(%d)" % self.value)


_str_to_value_type = {
    "unknown": ValueType.Unknown, "doc": ValueType.Doc, "vec": ValueType.Vec,
    "arr": ValueType.Arr, "obj": ValueType.Obj, "map": ValueType.Map,
    "str": ValueType.Str, "bytes": ValueType.Bytes, "bool": ValueType.Bool,
    "i": ValueType.I, "i8": ValueType.I8, "i16": ValueType.I16, "i32": ValueType.I32, "i64": ValueType.I64,
    "u": ValueType.U, "u8": ValueType.U8, "u16": ValueType.U16, "u32": ValueType.U32, "u64": ValueType.U64,
    "f32": ValueType.F32, "f64": ValueType.F64,
    "bigint": ValueType.Bigint, "datetime": ValueType.Datetime, "date": ValueType.Date, "time": ValueType.Time,
    "uuid": ValueType.Uuid, "decimal": ValueType.Decimal, "ip": ValueType.Ip, "url": ValueType.Url,
    "email": ValueType.Email, "enums": ValueType.Enums,
    "media": ValueType.Media,
}


def parse_value_type(s: str) -> ValueType:
    s = s.lower().strip()
    if s in _str_to_value_type:
        return _str_to_value_type[s]
    raise ValueError("Invalid ValueType string: %s" % s)