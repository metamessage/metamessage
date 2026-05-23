export enum ValueType {
  Unknown,
  Doc,
  Arr,
  Vec,
  Obj,
  Map,
  Str,
  Bytes,
  Bool,
  I,
  I8,
  I16,
  I32,
  I64,
  U,
  U8,
  U16,
  U32,
  U64,
  F32,
  F64,
  Bigint,
  Datetime,
  Date,
  Time,
  Uuid,
  Decimal,
  Ip,
  Url,
  Email,
  Enums,
  Image,
  Video,
}

// 你要的字符串名字，按顺序对应
export const ValueTypeStr = [
  'unknown',
  'doc',
  'arr',
  'vec',
  'obj',
  'map',
  'str',
  'bytes',
  'bool',
  'i',
  'i8',
  'i16',
  'i32',
  'i64',
  'u',
  'u8',
  'u16',
  'u32',
  'u64',
  'f32',
  'f64',
  'bigint',
  'datetime',
  'date',
  'time',
  'uuid',
  'decimal',
  'ip',
  'url',
  'email',
  'enums',
  'image',
  'video',
];

export function typeToString(value: ValueType): string {
  return ValueTypeStr[value as number] ?? 'unknown';
}

export function stringToType(value: string): ValueType {
  switch (value) {
    case 'doc':
      return ValueType.Doc;
    case 'arr':
      return ValueType.Arr;
    case 'vec':
      return ValueType.Vec;
    case 'obj':
      return ValueType.Obj;
    case 'map':
      return ValueType.Map;
    case 'str':
      return ValueType.Str;
    case 'bytes':
      return ValueType.Bytes;
    case 'bool':
      return ValueType.Bool;
    case 'i':
      return ValueType.I;
    case 'i8':
      return ValueType.I8;
    case 'i16':
      return ValueType.I16;
    case 'i32':
      return ValueType.I32;
    case 'i64':
      return ValueType.I64;
    case 'u':
      return ValueType.U;
    case 'u8':
      return ValueType.U8;
    case 'u16':
      return ValueType.U16;
    case 'u32':
      return ValueType.U32;
    case 'u64':
      return ValueType.U64;
    case 'f32':
      return ValueType.F32;
    case 'f64':
      return ValueType.F64;
    case 'bigint':
      return ValueType.Bigint;
    case 'datetime':
      return ValueType.Datetime;
    case 'date':
      return ValueType.Date;
    case 'time':
      return ValueType.Time;
    case 'uuid':
      return ValueType.Uuid;
    case 'decimal':
      return ValueType.Decimal;
    case 'ip':
      return ValueType.Ip;
    case 'url':
      return ValueType.Url;
    case 'email':
      return ValueType.Email;
    case 'enums':
      return ValueType.Enums;
    case 'image':
      return ValueType.Image;
    case 'video':
      return ValueType.Video;
    default:
      return ValueType.Unknown;
  }
}
