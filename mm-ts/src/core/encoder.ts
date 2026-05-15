import { MMBuffer } from './buffer';
import {
  Prefix,
  SimpleValue,
  StringLen1Byte,
  StringLen2Byte,
  BytesLen1Byte,
  BytesLen2Byte,
  ContainerLen1Byte,
  ContainerLen2Byte,
  ContainerObject,
  ContainerArray,
  IntLen1Byte,
  IntLen2Byte,
  IntLen3Byte,
  IntLen4Byte,
  IntLen5Byte,
  IntLen6Byte,
  IntLen7Byte,
  IntLen8Byte,
  Max1Byte,
  Max2Byte,
  Max3Byte,
  Max4Byte,
  Max5Byte,
  Max6Byte,
  Max7Byte,
  Max8Byte,
  TagLen1Byte,
  TagLen2Byte,
  FloatPositiveNegativeMask,
  FloatLen1Byte,
  FloatLen2Byte,
  FloatLen3Byte,
  FloatLen4Byte,
  FloatLen5Byte,
  FloatLen6Byte,
  FloatLen7Byte,
  FloatLen8Byte,
} from './constants';
import { ValueType } from '../ir/value-type';
import { Node, MMValue, MMObject, MMArray } from '../ir/ast';
import { Tag } from '../ir/tag';
import { ValueToNode } from './value-to-node';
import { toJSONC } from '../jsonc/printer';

export class MMEncoder {
  private buffer: MMBuffer;

  constructor() {
    this.buffer = new MMBuffer();
  }

  get result(): Uint8Array {
    return this.buffer.result;
  }

  reset(): void {
    this.buffer.reset();
  }

  encodeValue(value: any, tag?: Tag): Uint8Array {
    const node = ValueToNode(value, tag);
    // console.log('node1', (node as MMObject).getProperties().id?.getTag());
    // console.log('node2', (node as MMObject).getProperties().name?.getTag());
    // console.log('node3', (node as MMObject).getProperties().age?.getTag());
    const jsonc = toJSONC(node);
    // console.log('jsonc', jsonc);
    return this.encodeNode(node);
  }

  encodeNode(node: Node): Uint8Array {
    let n = 0;
    switch (node.getType()) {
      case 'object':
        n = this.encodeNodeObject(node as MMObject);
        break;
      case 'array':
        n = this.encodeNodeArray(node as MMArray);
        break;
      case 'value':
        n = this.encodeNodeValue(node as MMValue);
        break;
      default:
        throw new Error(`unsupported node type: ${node.getType()}`);
    }

    return this.buffer.getBytes(this.buffer.offset - n);
  }

  private encodeNodeObject(obj: MMObject): number {
    const tag = obj.getTag();
    const properties = obj.getProperties();

    const keysBuffer = new MMBuffer();
    const valuesBuffer = new MMBuffer();
    for (const [key, valueNode] of Object.entries(properties)) {
      keysBuffer.writeBytes(this.encodeStringToBytes(key));
      valuesBuffer.writeBytes(this.encodeNode(valueNode));
    }

    const keysEncoded = this.encodeArrayToBytes(keysBuffer.result);
    const valuesEncoded = valuesBuffer.result;

    const combined = new Uint8Array(keysEncoded.length + valuesEncoded.length);
    combined.set(keysEncoded, 0);
    combined.set(valuesEncoded, keysEncoded.length);

    const n = this.encodeObjectToBytes(combined);
    const n1 = this.encodeComment(
      this.buffer.getBytes(this.buffer.offset - n),
      tag,
    );
    if (n1 === 0) {
      return n;
    }

    return n1;
  }

  private encodeNodeArray(arr: MMArray): number {
    const tag = arr.getTag();
    const elements = arr.getElements();

    const valuesBuffer = new MMBuffer();
    let n = 0;
    for (const item of elements) {
      const temp = this.encodeNode(item);
      valuesBuffer.writeBytes(temp);
    }

    const arrayEncoded = this.encodeArrayToBytes(valuesBuffer.result);
    const n1 = this.encodeComment(arrayEncoded, tag);

    if (n1 === 0) {
      return n;
    }

    return n1;
  }

  private encodeNodeValue(val: MMValue): number {
    const tag = val.getTag();

    let n = 0;
    switch (tag.type) {
      case ValueType.DateTime:
        if (!tag.isNull) {
          n = this.encodeDateTime(val.getValue() as Date);
        }
        break;
      case ValueType.Date:
        if (!tag.isNull) {
          n = this.encodeDate(val.getValue() as Date);
        }
        break;
      case ValueType.Time:
        if (!tag.isNull) {
          n = this.encodeTime(val.getValue() as Date);
        }
        break;
      case ValueType.Int:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullInt);
        } else {
          n = this.encodeInt64(val.getValue() as bigint);
        }
        break;

      case ValueType.Int8:
      case ValueType.Int16:
      case ValueType.Int32:
      case ValueType.Int64:
        if (!tag.isNull) {
          n = this.encodeInt64(BigInt(val.getValue()));
        }
        break;

      case ValueType.Uint:
      case ValueType.Uint8:
      case ValueType.Uint16:
      case ValueType.Uint32:
      case ValueType.Uint64:
        if (!tag.isNull) {
          n = this.encodeUint64(BigInt(val.getValue()));
        }
        break;

      case ValueType.Float32:
        if (!tag.isNull) {
          n = this.encodeFloat(String(val.getValue()));
        }
        break;

      case ValueType.Float64:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullFloat);
        } else {
          n = this.encodeFloat(String(val.getValue()));
        }
        break;

      case ValueType.String:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullString);
        } else {
          n = this.encodeString(String(val.getValue()));
        }
        break;
      case ValueType.Email:
        if (!tag.isNull) {
          n = this.encodeString(String(val.getValue()));
        }
        break;
      case ValueType.UUID:
        if (!tag.isNull) {
          const data = val.getValue();
          if (data instanceof Uint8Array) {
            n = this.encodeBytes(data);
          } else {
            const uuidStr = String(data);
            const bytes = new Uint8Array(16);
            const parts = uuidStr.replace(/-/g, '');
            for (let i = 0; i < 16; i++) {
              bytes[i] = parseInt(parts.substr(i * 2, 2), 16);
            }
            n = this.encodeBytes(bytes);
          }
        }
        break;
      case ValueType.Decimal:
        if (!tag.isNull) {
          n = this.encodeFloat(String(val.getValue()));
        }
        break;
      case ValueType.URL:
        if (!tag.isNull) {
          n = this.encodeString(String(val.getValue()));
        }
        break;
      case ValueType.IP:
        if (!tag.isNull) {
          const data = val.getValue();
          if (data instanceof Uint8Array) {
            n = this.encodeBytes(data);
          } else {
            n = this.encodeString(String(data));
          }
        }
        break;
      case ValueType.Bytes:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullBytes);
        } else {
          const data = val.getValue();
          n = this.encodeBytes(
            data instanceof Uint8Array ? data : new Uint8Array(data),
          );
        }
        break;
      case ValueType.BigInt:
        n = this.encodeBigInt(String(val.getValue()));
        break;
      case ValueType.Bool:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullBool);
        } else {
          n = this.encodeBool(Boolean(val.getValue()));
        }
        break;
      case ValueType.Enum:
        if (!tag.isNull) {
          n = this.encodeInt64(BigInt(val.getValue() as number));
        }
        break;
      default:
        throw new Error(`type error: unsupported type: ${tag.type}`);
    }

    const n1 = this.encodeComment(
      this.buffer.getBytes(this.buffer.offset - n),
      tag,
    );
    if (n1 === 0) {
      return n;
    }

    return n1;
  }

  private encodeComment(payload: Uint8Array, tag: Tag): number {
    if (!(tag instanceof Tag)) {
      tag = new Tag(tag);
    }

    const tagBytes = tag.toBytes();
    if (tagBytes.length === 0) {
      return 0;
    }

    const n = this.encodeT(tagBytes);
    if (n === 0) {
      return n;
    }

    return this.encodeTag(
      payload,
      this.buffer.getBytes(this.buffer.offset - n),
    );
  }

  private encodeTag(bs: Uint8Array, tag: Uint8Array): number {
    const length = bs.length + tag.length;
    if (length > Max2Byte) {
      throw new Error(
        `tag+payload too long, max length: ${Max2Byte}, actual: ${length}`,
      );
    }

    const sign = Prefix.Tag;
    const savedOffset = this.buffer.offset;

    if (length < TagLen1Byte) {
      const prefix = sign | length;
      this.buffer.writeUint8(prefix);
      this.buffer.writeBytes(tag);
      this.buffer.writeBytes(bs);
    } else if (length < Max1Byte) {
      const prefix = sign | TagLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
      this.buffer.writeBytes(tag);
      this.buffer.writeBytes(bs);
    } else {
      const prefix = sign | TagLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
      this.buffer.writeBytes(tag);
      this.buffer.writeBytes(bs);
    }

    return this.buffer.offset - savedOffset;
  }

  private encodeT(bs: Uint8Array): number {
    const length = bs.length;
    if (length === 0) {
      return 0;
    }

    const savedOffset = this.buffer.offset;

    if (length < 254) {
      this.buffer.writeUint8(length);
      this.buffer.writeBytes(bs);
    } else if (length < 257) {
      this.buffer.writeUint8(254);
      this.buffer.writeUint8(length);
      this.buffer.writeBytes(bs);
    } else {
      this.buffer.writeUint8(255);
      this.buffer.writeUint16LE(length);
      this.buffer.writeBytes(bs);
    }

    return this.buffer.offset - savedOffset;
  }

  private encodeSimple(value: SimpleValue): number {
    const prefix = Prefix.Simple;
    this.buffer.writeUint8(prefix | value);
    return 1;
  }

  encodeBool(value: boolean): number {
    const suffix = value ? SimpleValue.True : SimpleValue.False;
    return this.encodeSimple(suffix);
  }

  encodeNil(): void {
    this.encodeSimple(SimpleValue.NullInt);
  }

  encodeInt(value: bigint): void {
    if (value >= 0) {
      this.encodeUInt(value);
    } else {
      this.encodeNegativeInt(-value);
    }
  }

  encodeInt64(value: bigint): number {
    let sign: number;
    let uv: bigint;

    if (value >= 0) {
      sign = Prefix.PositiveInt;
      uv = value;
    } else {
      sign = Prefix.NegativeInt;
      if (value === -9223372036854775808n) {
        uv = 9223372036854775808n;
      } else {
        uv = -value;
      }
    }

    return this.encodeIntWithSign(sign, uv);
  }

  encodeUint64(value: bigint): number {
    return this.encodeIntWithSign(Prefix.PositiveInt, value);
  }

  private encodeIntWithSign(sign: number, value: bigint): number {
    if (value <= 23n) {
      this.buffer.writeUint8(sign | Number(value));
      return 1;
    } else if (value <= Max1Byte) {
      this.buffer.writeUint8(sign | IntLen1Byte);
      this.buffer.writeUint8(Number(value));
      return 2;
    } else if (value <= Max2Byte) {
      this.buffer.writeUint8(sign | IntLen2Byte);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      return 3;
    } else if (value <= Max3Byte) {
      this.buffer.writeUint8(sign | IntLen3Byte);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      return 4;
    } else if (value <= Max4Byte) {
      this.buffer.writeUint8(sign | IntLen4Byte);
      this.buffer.writeUint32LE(Number(value));
      return 5;
    } else if (value <= Max5Byte) {
      this.buffer.writeUint8(sign | IntLen5Byte);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      return 6;
    } else if (value <= Max6Byte) {
      this.buffer.writeUint8(sign | IntLen6Byte);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      return 7;
    } else if (value <= Max7Byte) {
      this.buffer.writeUint8(sign | IntLen7Byte);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 48n) & 0xffn));
      return 8;
    } else if (value <= Max8Byte) {
      this.buffer.writeUint8(sign | IntLen8Byte);
      this.buffer.writeUint64LE(value);
      return 9;
    } else {
      throw new Error('Int too large');
      return 0;
    }
  }

  private encodeUInt(value: bigint): void {
    this.encodeIntWithSign(Prefix.PositiveInt, value);
  }

  private encodeNegativeInt(value: bigint): void {
    this.encodeIntWithSign(Prefix.NegativeInt, value);
  }

  encodeFloat(value: string): number {
    const result = this.parseStringToUint64(value);
    if (!result.success) {
      throw new Error(result.error || 'Failed to parse float');
    }

    const { isNegative, exponent, mantissa } = result;

    let sign = Prefix.Float;
    if (isNegative) {
      sign |= FloatPositiveNegativeMask;
    }

    if (exponent === -1 && mantissa <= 7) {
      sign |= Number(mantissa);
      this.buffer.writeUint8(sign);
      return 1;
    }

    if (mantissa <= Max1Byte) {
      sign |= FloatLen1Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number(mantissa));
      return 1 + 1 + 1;
    } else if (mantissa <= Max2Byte) {
      sign |= FloatLen2Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 2;
    } else if (mantissa <= Max3Byte) {
      sign |= FloatLen3Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 3;
    } else if (mantissa <= Max4Byte) {
      sign |= FloatLen4Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 4;
    } else if (mantissa <= BigInt(Max5Byte)) {
      sign |= FloatLen5Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 5;
    } else if (mantissa <= BigInt(Max6Byte)) {
      sign |= FloatLen6Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 6;
    } else if (mantissa <= BigInt(Max7Byte)) {
      sign |= FloatLen7Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 7;
    } else if (mantissa <= Max8Byte) {
      sign |= FloatLen8Byte;
      this.buffer.writeUint8(sign);
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 56n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      return 1 + 1 + 8;
    } else {
      throw new Error(`Unsupported mantissa length: ${mantissa}`);
    }
  }

  private parseStringToUint64(s: string): {
    success: boolean;
    isNegative: boolean;
    exponent: number;
    mantissa: bigint;
    error?: string;
  } {
    if (s === '') {
      return {
        success: false,
        isNegative: false,
        exponent: 0,
        mantissa: 0n,
        error: 'empty numeric string',
      };
    }

    let isNegative = s.startsWith('-');
    if (isNegative) {
      s = s.substring(1);
      if (s === '') {
        return {
          success: false,
          isNegative: false,
          exponent: 0,
          mantissa: 0n,
          error: 'invalid numeric string: only minus sign',
        };
      }
    }

    let expPart = '';
    const eIdx = s.search(/[eE]/);
    if (eIdx !== -1) {
      expPart = s.substring(eIdx + 1);
      s = s.substring(0, eIdx);
      if (expPart === '') {
        return {
          success: false,
          isNegative: false,
          exponent: 0,
          mantissa: 0n,
          error: 'missing exponent part in scientific notation',
        };
      }
    }

    const parts = s.split('.');
    let intPart = parts[0];
    let fracPart = '';
    if (parts.length > 1) {
      fracPart = parts[1]!;
    }

    if (intPart === '') {
      intPart = '0';
    }

    let baseExp = -fracPart.length;

    if (expPart !== '') {
      const exp = parseInt(expPart, 10);
      if (isNaN(exp)) {
        return {
          success: false,
          isNegative: false,
          exponent: 0,
          mantissa: 0n,
          error: 'invalid exponent',
        };
      }
      baseExp += exp;
    }

    if (baseExp < -128 || baseExp > 127) {
      return {
        success: false,
        isNegative: false,
        exponent: 0,
        mantissa: 0n,
        error: `final exponent out of range (-128 ~ 127): ${baseExp}`,
      };
    }

    let mantissaStr = (intPart + fracPart).replace(/^0+/, '');
    if (mantissaStr === '') {
      mantissaStr = '0';
    }

    const mantissa = BigInt(mantissaStr);

    return { success: true, isNegative, exponent: baseExp, mantissa };
  }

  encodeDouble(value: number): void {
    const prefix = Prefix.Float | 31;
    this.buffer.writeUint8(prefix);
    this.buffer.writeFloat64(value);
  }

  encodeString(value: string): number {
    let n = 0;
    const length = value.length;
    if (length < StringLen1Byte) {
      const prefix = Prefix.String | length;
      this.buffer.writeUint8(prefix);
      n += 1;
    } else if (length < 256) {
      const prefix = Prefix.String | StringLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
      n += 2;
    } else {
      const prefix = Prefix.String | StringLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
      n += 3;
    }

    for (let i = 0; i < length; i++) {
      const code = value.charCodeAt(i);
      if (code <= 0x7f) {
        this.buffer.writeUint8(code);
      } else if (code <= 0x7ff) {
        this.buffer.writeUint8(0xc0 | (code >> 6));
        this.buffer.writeUint8(0x80 | (code & 0x3f));
      } else if (code <= 0xffff) {
        this.buffer.writeUint8(0xe0 | (code >> 12));
        this.buffer.writeUint8(0x80 | ((code >> 6) & 0x3f));
        this.buffer.writeUint8(0x80 | (code & 0x3f));
      } else {
        this.buffer.writeUint8(0xf0 | (code >> 18));
        this.buffer.writeUint8(0x80 | ((code >> 12) & 0x3f));
        this.buffer.writeUint8(0x80 | ((code >> 6) & 0x3f));
        this.buffer.writeUint8(0x80 | (code & 0x3f));
      }
    }
    n += length;
    return n;
  }

  private encodeStringToBytes(value: string): Uint8Array {
    const startOffset = this.buffer.offset;
    this.encodeString(value);
    return this.buffer.getBytes(startOffset);
  }

  encodeBytes(value: Uint8Array): number {
    let n = 0;
    const length = value.length;
    if (length < BytesLen1Byte) {
      const prefix = Prefix.Bytes | length;
      this.buffer.writeUint8(prefix);
      n += 1;
    } else if (length < 256) {
      const prefix = Prefix.Bytes | BytesLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
      n += 2;
    } else {
      const prefix = Prefix.Bytes | BytesLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
      n += 3;
    }

    this.buffer.writeBytes(value);
    n += value.length;
    return n;
  }

  encodeBigInt(value: string): number {
    const startOffset = this.buffer.offset;
    this.buffer.writeUint8(value.length);
    this.encodeBigIntInternal(value);
    const bytes = this.buffer.getBytes(startOffset);
    return this.encodeBytes(bytes);
  }

  private encodeBigIntInternal(s: string): number {
    if (s === '') {
      return 0;
    }

    let neg = false;
    if (s.startsWith('-')) {
      neg = true;
      s = s.substring(1);
    }

    const bits: number[] = [];
    bits.push(neg ? 1 : 0);

    const n = s.length;
    let i = 0;
    for (; i < n; i += 3) {
      const rem = n - i;
      if (rem >= 3) {
        const num = this.atoi(s.substring(i, i + 3));
        bits.push(...this.toBits(num, 10));
      } else if (rem === 2) {
        const num = this.atoi(s.substring(i, i + 2));
        bits.push(...this.toBits(num, 7));
      } else {
        const num = this.atoi(s.substring(i, i + 1));
        bits.push(...this.toBits(num, 4));
      }
    }

    return this.writeBits(bits);
  }

  private atoi(s: string): number {
    let v = 0;
    for (const c of s) {
      v = v * 10 + (c.charCodeAt(0) - '0'.charCodeAt(0));
    }
    return v;
  }

  private toBits(v: number, n: number): number[] {
    const b: number[] = new Array(n);
    for (let i = 0; i < n; i++) {
      b[n - 1 - i] = (v >> i) & 1;
    }
    return b;
  }

  private writeBits(bits: number[]): number {
    let bt = 0;
    let bl = 0;
    let n = 0;

    for (const b of bits) {
      bt = (bt << 1) | (b & 1);
      bl++;
      if (bl === 8) {
        this.buffer.writeUint8(bt);
        n++;
        bt = 0;
        bl = 0;
      }
    }

    if (bl > 0) {
      bt <<= 8 - bl;
      this.buffer.writeUint8(bt);
      n++;
    }

    return n;
  }

  encodeDateTime(date: Date): number {
    const v = date.getTime() / 1000;
    return this.encodeInt64(BigInt(Math.floor(v)));
  }

  encodeDate(date: Date): number {
    const utcDate = new Date(date.getTime() + date.getTimezoneOffset() * 60000);
    const defaultTime = new Date('1970-01-01T00:00:00Z');
    const diff = utcDate.getTime() - defaultTime.getTime();
    const days = Math.floor(diff / (24 * 60 * 60 * 1000));
    return this.encodeInt64(BigInt(days));
  }

  encodeTime(date: Date): number {
    const utcDate = new Date(date.getTime() + date.getTimezoneOffset() * 60000);
    const seconds =
      utcDate.getHours() * 3600 +
      utcDate.getMinutes() * 60 +
      utcDate.getSeconds();
    return this.encodeUint64(BigInt(seconds));
  }

  private encodeArrayToBytes(value: Uint8Array): Uint8Array {
    const startOffset = this.buffer.offset;

    let prefix = Prefix.Container | ContainerArray;
    const length = value.length;
    if (length < ContainerLen1Byte) {
      prefix |= length;
      this.buffer.writeUint8(prefix);
    } else if (length < 256) {
      prefix |= ContainerLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
    } else {
      prefix |= ContainerLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
    }

    this.buffer.writeBytes(value);
    return this.buffer.getBytes(startOffset);
  }

  private encodeObjectToBytes(value: Uint8Array): number {
    const length = value.length;
    let n = 0;

    let prefix = Prefix.Container | ContainerObject;

    if (length < ContainerLen1Byte) {
      prefix |= length;
      this.buffer.writeUint8(prefix);
      n = 1 + length;
    } else if (length < 256) {
      prefix |= ContainerLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
      n = 2 + length;
    } else {
      prefix |= ContainerLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
      n = 3 + length;
    }

    this.buffer.writeBytes(value);
    return n;
  }
}
