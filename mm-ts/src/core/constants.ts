export enum Prefix {
  Simple = 0b000 << 5,
  PositiveInt = 0b001 << 5,
  NegativeInt = 0b010 << 5,
  Float = 0b011 << 5,
  String = 0b100 << 5,
  Bytes = 0b101 << 5,
  Container = 0b110 << 5,
  Tag = 0b111 << 5,
}

export function prefixToString(prefix: Prefix): string {
  switch (prefix) {
    case Prefix.Simple:
      return 'Simple';
    case Prefix.PositiveInt:
      return 'PositiveInt';
    case Prefix.NegativeInt:
      return 'NegativeInt';
    case Prefix.Float:
      return 'Float';
    case Prefix.String:
      return 'String';
    case Prefix.Bytes:
      return 'Bytes';
    case Prefix.Container:
      return 'Container';
    case Prefix.Tag:
      return 'Tag';
    default:
      return `Prefix(0x${(prefix as number).toString(16).padStart(2, '0')})`;
  }
}

export enum SimpleValue {
  NullBool = 0,
  NullInt,
  NullFloat,
  NullString,
  NullBytes,
  False,
  True,
  Code,
  Message,
  Data,
  Success,
  Error,
  Unknown,
  Page,
  Limit,
  Offset,
  Total,
  Id,
  Name,
  Description,
  Type,
  Version,
  Status,
  Url,
  CreateTime,
  UpdateTime,
  DeleteTime,
  Account,
  Token,
  ExpireTime,
  Key,
  Val,
}

export function simpleValueToString(value: SimpleValue): string {
  switch (value) {
    case SimpleValue.NullBool:
      return 'null_bool';
    case SimpleValue.NullInt:
      return 'null_int';
    case SimpleValue.NullFloat:
      return 'null_float';
    case SimpleValue.NullString:
      return 'null_string';
    case SimpleValue.NullBytes:
      return 'null_bytes';
    case SimpleValue.False:
      return 'false';
    case SimpleValue.True:
      return 'true';
    case SimpleValue.Code:
      return 'code';
    case SimpleValue.Message:
      return 'message';
    case SimpleValue.Data:
      return 'data';
    case SimpleValue.Success:
      return 'success';
    case SimpleValue.Error:
      return 'error';
    case SimpleValue.Unknown:
      return 'unknown';
    case SimpleValue.Page:
      return 'page';
    case SimpleValue.Limit:
      return 'limit';
    case SimpleValue.Offset:
      return 'offset';
    case SimpleValue.Total:
      return 'total';
    case SimpleValue.Id:
      return 'id';
    case SimpleValue.Name:
      return 'name';
    case SimpleValue.Description:
      return 'description';
    case SimpleValue.Type:
      return 'type';
    case SimpleValue.Version:
      return 'version';
    case SimpleValue.Status:
      return 'status';
    case SimpleValue.Url:
      return 'url';
    case SimpleValue.CreateTime:
      return 'create_time';
    case SimpleValue.UpdateTime:
      return 'update_time';
    case SimpleValue.DeleteTime:
      return 'delete_time';
    case SimpleValue.Account:
      return 'account';
    case SimpleValue.Token:
      return 'token';
    case SimpleValue.ExpireTime:
      return 'expire_time';
    case SimpleValue.Key:
      return 'key';
    case SimpleValue.Val:
      return 'value';
    default:
      return `SimpleValue(${value})`;
  }
}

export function isValidSimpleValue(value: SimpleValue): boolean {
  return value < 32;
}

export const Max1Byte = 0xff;
export const Max2Byte = 0xffff;
export const Max3Byte = 0xffffff;
export const Max4Byte = 0xffffffff;
export const Max5Byte = 0xffffffffff;
export const Max6Byte = 0xffffffffffff;
export const Max7Byte = 0xffffffffffffff;
export const Max8Byte = 0xffffffffffffffffn;

export const IntLenMask = 0b11111;
export const IntLen1Byte = 24;
export const IntLen2Byte = 25;
export const IntLen3Byte = 26;
export const IntLen4Byte = 27;
export const IntLen5Byte = 28;
export const IntLen6Byte = 29;
export const IntLen7Byte = 30;
export const IntLen8Byte = 31;

export const FloatPositiveNegativeMask = 0b10000;
export const FloatLenMask = 0b01111;
export const FloatLen1Byte = FloatLenMask - 7;
export const FloatLen2Byte = FloatLenMask - 6;
export const FloatLen3Byte = FloatLenMask - 5;
export const FloatLen4Byte = FloatLenMask - 4;
export const FloatLen5Byte = FloatLenMask - 3;
export const FloatLen6Byte = FloatLenMask - 2;
export const FloatLen7Byte = FloatLenMask - 1;
export const FloatLen8Byte = FloatLenMask;

export const StringLenMask = 0b11111;
export const StringLen1Byte = StringLenMask - 1;
export const StringLen2Byte = StringLenMask;

export const BytesLenMask = 0b11111;
export const BytesLen1Byte = BytesLenMask - 1;
export const BytesLen2Byte = BytesLenMask;

export const ContainerMask = 0b10000;
export const ContainerObject = 0b00000;
export const ContainerArray = 0b10000;
export const ContainerLenMask = 0b01111;
export const ContainerLen1Byte = ContainerLenMask - 1;
export const ContainerLen2Byte = ContainerLenMask;

export const TagLenMask = 0b11111;
export const TagLen1Byte = TagLenMask - 1;
export const TagLen2Byte = TagLenMask;

export const TagPayload1Byte = BytesLenMask - 1;
export const TagPayload2Byte = BytesLenMask;

export const PrefixMask = 0b11100000;
export const SuffixMask = 0b00011111;

export function getPrefix(b: number): Prefix {
  return b & PrefixMask as Prefix;
}

export function getSuffix(b: number): number {
  return b & SuffixMask;
}

export function tagLen(b: number): [number, number] {
  const l = b & TagLenMask;
  switch (l) {
    case TagLen1Byte:
      return [1, 0];
    case TagLen2Byte:
      return [2, 0];
    default:
      return [0, l];
  }
}

export function containerLen(b: number): [number, number] {
  const l = b & ContainerLenMask;
  switch (l) {
    case ContainerLen1Byte:
      return [1, 0];
    case ContainerLen2Byte:
      return [2, 0];
    default:
      return [0, l];
  }
}

export function isArray(b: number): boolean {
  return (b & ContainerMask) === ContainerArray;
}

export function stringLen(b: number): [number, number] {
  const l = b & StringLenMask;
  if (l < StringLen1Byte) {
    return [0, l];
  } else if (l === StringLen1Byte) {
    return [1, l];
  } else {
    return [2, l];
  }
}

export function bytesLen(b: number): [number, number] {
  const l = b & BytesLenMask;
  if (l < BytesLen1Byte) {
    return [0, l];
  } else if (l === BytesLen1Byte) {
    return [1, l];
  } else {
    return [2, l];
  }
}

export function intLen(b: number): [number, number] {
  const l = b & IntLenMask;
  if (l < 24) {
    return [0, l];
  } else {
    return [l - 24 + 1, 0];
  }
}

export function floatLen(b: number): [number, number] {
  const l = b & FloatLenMask;
  if (l < FloatLen1Byte) {
    return [0, l];
  } else {
    return [l - FloatLen1Byte + 1, 0];
  }
}