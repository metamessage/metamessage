export enum ValueType {
  Unknown,
  Doc,
  Array,
  Slice,
  Object,
  Map,
  String,
  Bytes,
  Bool,
  Int,
  Int8,
  Int16,
  Int32,
  Int64,
  Uint,
  Uint8,
  Uint16,
  Uint32,
  Uint64,
  Float32,
  Float64,
  BigInt,
  DateTime,
  Date,
  Time,
  UUID,
  Decimal,
  IP,
  URL,
  Email,
  Enum,
  Image,
  Video,
}

// 你要的字符串名字，按顺序对应
export const ValueTypeString = [
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
  'enum',
  'image',
  'video',
];

export function typeToString(value: ValueType): string {
  return ValueTypeString[value as number] ?? 'unknown';
}

export function stringToType(value: string): ValueType {
  switch (value) {
    case 'doc':
      return ValueType.Doc;
    case 'arr':
      return ValueType.Array;
    case 'vec':
      return ValueType.Slice;
    case 'obj':
      return ValueType.Object;
    case 'map':
      return ValueType.Map;
    case 'str':
      return ValueType.String;
    case 'bytes':
      return ValueType.Bytes;
    case 'bool':
      return ValueType.Bool;
    case 'i':
      return ValueType.Int;
    case 'i8':
      return ValueType.Int8;
    case 'i16':
      return ValueType.Int16;
    case 'i32':
      return ValueType.Int32;
    case 'i64':
      return ValueType.Int64;
    case 'u':
      return ValueType.Uint;
    case 'u8':
      return ValueType.Uint8;
    case 'u16':
      return ValueType.Uint16;
    case 'u32':
      return ValueType.Uint32;
    case 'u64':
      return ValueType.Uint64;
    case 'f32':
      return ValueType.Float32;
    case 'f64':
      return ValueType.Float64;
    case 'bigint':
      return ValueType.BigInt;
    case 'datetime':
      return ValueType.DateTime;
    case 'date':
      return ValueType.Date;
    case 'time':
      return ValueType.Time;
    case 'uuid':
      return ValueType.UUID;
    case 'decimal':
      return ValueType.Decimal;
    case 'ip':
      return ValueType.IP;
    case 'url':
      return ValueType.URL;
    case 'email':
      return ValueType.Email;
    case 'enum':
      return ValueType.Enum;
    case 'image':
      return ValueType.Image;
    case 'video':
      return ValueType.Video;
    default:
      return ValueType.Unknown;
  }
}
