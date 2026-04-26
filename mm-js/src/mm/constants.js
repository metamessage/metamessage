export const MMPrefix = {
  Simple: 0b000 << 5,
  PositiveInt: 0b001 << 5,
  NegativeInt: 0b010 << 5,
  PrefixFloat: 0b011 << 5,
  PrefixString: 0b100 << 5,
  PrefixBytes: 0b101 << 5,
  Container: 0b110 << 5,
  PrefixTag: 0b111 << 5
};

export const MMSimpleValue = {
  NullBool: 0,
  NullInt: 1,
  NullFloat: 2,
  NullString: 3,
  NullBytes: 4,
  FalseValue: 5,
  TrueValue: 6,
  Code: 7,
  Message: 8,
  Data: 9,
  Success: 10,
  Error: 11,
  Unknown: 12,
  Page: 13,
  Limit: 14,
  Offset: 15,
  Total: 16,
  Id: 17,
  Name: 18,
  Description: 19,
  TypeValue: 20,
  Version: 21,
  Status: 22,
  Url: 23,
  CreateTime: 24,
  UpdateTime: 25,
  DeleteTime: 26,
  Account: 27,
  Token: 28,
  ExpireTime: 29,
  Key: 30,
  Val: 31
};

export const MMConstants = {
  Max1Byte: 0xFFn,
  Max2Byte: 0xFFFFn,
  Max3Byte: 0xFFFFFFn,
  Max4Byte: 0xFFFFFFFFn,
  Max5Byte: 0xFFFFFFFFFFn,
  Max6Byte: 0xFFFFFFFFFFFFn,
  Max7Byte: 0xFFFFFFFFFFFFFFn,
  Max8Byte: 0xFFFFFFFFFFFFFFFFn,

  IntLenMask: 0b11111,
  IntLen1Byte: 0b11000,
  IntLen2Byte: 0b11001,
  IntLen3Byte: 0b11010,
  IntLen4Byte: 0b11011,
  IntLen5Byte: 0b11100,
  IntLen6Byte: 0b11101,
  IntLen7Byte: 0b11110,
  IntLen8Byte: 0b11111,

  FloatPositiveNegativeMask: 0b10000,
  FloatLenMask: 0b01111,
  FloatLen1Byte: 0b00111,
  FloatLen2Byte: 0b01000,
  FloatLen3Byte: 0b01001,
  FloatLen4Byte: 0b01010,
  FloatLen5Byte: 0b01011,
  FloatLen6Byte: 0b01100,
  FloatLen7Byte: 0b01101,
  FloatLen8Byte: 0b01110,

  StringLenMask: 0b11111,
  StringLen1Byte: 0b11110,
  StringLen2Byte: 0b11111,

  BytesLenMask: 0b11111,
  BytesLen1Byte: 0b11110,
  BytesLen2Byte: 0b11111,

  ContainerMask: 0b10000,
  ContainerObject: 0b00000,
  ContainerArray: 0b10000,
  ContainerLenMask: 0b01111,
  ContainerLen1Byte: 0b01110,
  ContainerLen2Byte: 0b01111,

  TagLenMask: 0b11111,
  TagLen1Byte: 0b11110,
  TagLen2Byte: 0b11111,

  PrefixMask: 0b11100000,
  SuffixMask: 0b00011111
};

export function getPrefix(b) {
  return b & MMConstants.PrefixMask;
}

export function getSuffix(b) {
  return b & MMConstants.SuffixMask;
}

export function intLen(b) {
  const l = b & MMConstants.IntLenMask;
  if (l < MMConstants.IntLen1Byte) {
    return { extraBytes: 0, len: l };
  } else {
    return { extraBytes: l - MMConstants.IntLen1Byte + 1, len: 0 };
  }
}

export function floatLen(b) {
  const l = b & MMConstants.FloatLenMask;
  if (l < MMConstants.FloatLen1Byte) {
    return { extraBytes: 0, len: l };
  } else {
    return { extraBytes: l - MMConstants.FloatLen1Byte + 1, len: 0 };
  }
}

export function stringLen(b) {
  const l = b & MMConstants.StringLenMask;
  if (l < MMConstants.StringLen1Byte) {
    return { extraBytes: 0, len: l };
  } else if (l === MMConstants.StringLen1Byte) {
    return { extraBytes: 1, len: l };
  } else {
    return { extraBytes: 2, len: l };
  }
}

export function bytesLen(b) {
  const l = b & MMConstants.BytesLenMask;
  if (l < MMConstants.BytesLen1Byte) {
    return { extraBytes: 0, len: l };
  } else if (l === MMConstants.BytesLen1Byte) {
    return { extraBytes: 1, len: l };
  } else {
    return { extraBytes: 2, len: l };
  }
}

export function containerLen(b) {
  const l = b & MMConstants.ContainerLenMask;
  if (l < MMConstants.ContainerLen1Byte) {
    return { extraBytes: 0, len: l };
  } else if (l === MMConstants.ContainerLen1Byte) {
    return { extraBytes: 1, len: l };
  } else {
    return { extraBytes: 2, len: l };
  }
}

export function tagLen(b) {
  const l = b & MMConstants.TagLenMask;
  if (l < MMConstants.TagLen1Byte) {
    return { extraBytes: 0, len: l };
  } else if (l === MMConstants.TagLen1Byte) {
    return { extraBytes: 1, len: l };
  } else {
    return { extraBytes: 2, len: l };
  }
}

export function isArray(b) {
  return (b & MMConstants.ContainerMask) === MMConstants.ContainerArray;
}

export const ValueType = {
  Unknown: 'unknown',
  Doc: 'doc',
  Slice: 'slice',
  Array: 'array',
  Struct: 'struct',
  Map: 'map',
  String: 'str',
  Bytes: 'bytes',
  Bool: 'bool',
  Int: 'int',
  Int8: 'i8',
  Int16: 'i16',
  Int32: 'i32',
  Int64: 'i64',
  Uint: 'uint',
  Uint8: 'u8',
  Uint16: 'u16',
  Uint32: 'u32',
  Uint64: 'u64',
  Float32: 'f32',
  Float64: 'f64',
  BigInt: 'bi',
  DateTime: 'datetime',
  Date: 'date',
  Time: 'time',
  Uuid: 'uuid',
  Decimal: 'decimal',
  Ip: 'ip',
  Url: 'url',
  Email: 'email',
  Enum: 'enum',
  Image: 'image',
  Video: 'video'
};