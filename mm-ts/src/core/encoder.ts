import { MMBuffer } from './buffer';
import {
  Prefix,
  SimpleValue,
  StringLen1Byte,
  StringLen2Byte,
  BytesLen1Byte,
  BytesLen2Byte,
  BytesLen3Byte,
  BytesLen4Byte,
  ContainerLen1Byte,
  ContainerLen2Byte,
  ContainerLen3Byte,
  ContainerLen4Byte,
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
  TagLen3Byte,
  TagLen4Byte,
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
import { Node, NodeScalar, NodeObject, NodeArray, NodeNull } from '../ir/ast';
import { Tag } from '../ir/tag';
import { ValueToNode } from './value-to-node';
import { toJSONC } from '../jsonc/printer';
import { footerSlotBytes, reverseBytes } from './footer';

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
    // console.log('node1', (node as NodeObject).getProperties().id?.getTag());
    // console.log('node2', (node as NodeObject).getProperties().name?.getTag());
    // console.log('node3', (node as NodeObject).getProperties().age?.getTag());
    const jsonc = toJSONC(node);
    // console.log('jsonc', jsonc);
    return this.encodeNode(node);
  }

  encodeNode(node: Node): Uint8Array {
    let n = 0;
    switch (node.getType()) {
      case 'object':
        n = this.encodeNodeObject(node as NodeObject);
        break;
      case 'array':
        n = this.encodeNodeArray(node as NodeArray);
        break;
      case 'value':
        n = this.encodeNodeValue(node as NodeScalar);
        break;
      case 'null':
        {
          const tag = node.getTag();
          const n = this.encodeSimple(SimpleValue.SimpleNull);
          const buf = this.buffer.getBytes(this.buffer.offset - n);
          const n1 = this.encodeComment(buf, tag);
          if (n1 === 0) {
            return buf;
          }
          return this.buffer.getBytes(this.buffer.offset - n1);
        }
      default:
        throw new Error(`unsupported node type: ${node.getType()}`);
    }

    return this.buffer.getBytes(this.buffer.offset - n);
  }

  private encodeNodeObject(obj: NodeObject): number {
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

  private encodeNodeArray(arr: NodeArray): number {
    const tag = arr.getTag();
    const elements = arr.getElements();

    const valuesBuffer = new MMBuffer();
    for (const item of elements) {
      const temp = this.encodeNode(item);
      valuesBuffer.writeBytes(temp);
    }

    const arrayEncoded = this.encodeArrayToBytes(valuesBuffer.result);
    const savedOffset = this.buffer.offset;
    this.buffer.writeBytes(arrayEncoded);
    const n = this.buffer.offset - savedOffset;

    const n1 = this.encodeComment(arrayEncoded, tag);

    if (n1 === 0) {
      return n;
    }

    return n1;
  }

  private encodeNodeValue(val: NodeScalar): number {
    const tag = val.getTag();

    let n = 0;
    switch (tag.type) {
      case ValueType.Datetime:
        if (!tag.isNull) {
          n = this.encodeDateTime(val.getValue() as Date, tag.location || 0);
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
      case ValueType.I:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullInt);
        } else {
          n = this.encodeInt64(BigInt(val.getValue()));
        }
        break;

      case ValueType.I8:
      case ValueType.I16:
      case ValueType.I32:
      case ValueType.I64:
        if (!tag.isNull) {
          n = this.encodeInt64(BigInt(val.getValue()));
        }
        break;

      case ValueType.U:
      case ValueType.U8:
      case ValueType.U16:
      case ValueType.U32:
      case ValueType.U64:
        if (!tag.isNull) {
          n = this.encodeU64(BigInt(val.getValue()));
        }
        break;

      case ValueType.F32:
        if (!tag.isNull) {
          n = this.encodeFloat(String(val.getValue()));
        }
        break;

      case ValueType.F64:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullFloat);
        } else {
          n = this.encodeFloat(String(val.getValue()));
        }
        break;

      case ValueType.Str:
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
      case ValueType.Uuid:
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
      case ValueType.Url:
        if (!tag.isNull) {
          n = this.encodeString(String(val.getValue()));
        }
        break;
      case ValueType.Ip:
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
      case ValueType.Media:
        if (!tag.isNull) {
          const data = val.getValue();
          let bytes: Uint8Array;
          if (typeof data === 'string') {
            const binaryStr = atob(data);
            bytes = new Uint8Array(binaryStr.length);
            for (let i = 0; i < binaryStr.length; i++) {
              bytes[i] = binaryStr.charCodeAt(i);
            }
          } else if (data instanceof Uint8Array) {
            bytes = data;
          } else {
            bytes = new Uint8Array(data as ArrayBuffer);
          }
          n = this.encodeBytes(bytes);
        }
        break;
      case ValueType.Bigint:
        n = this.encodeBigInt(String(val.getValue()));
        break;
      case ValueType.Bool:
        if (tag.isNull) {
          n = this.encodeSimple(SimpleValue.NullBool);
        } else {
          n = this.encodeBool(Boolean(val.getValue()));
        }
        break;
      case ValueType.Enums:
        if (!tag.isNull) {
          n = this.encodeInt64(BigInt(val.getValue() as number));
        }
        break;
      case ValueType.Unknown:
        const raw = val.getValue();
        if (typeof raw === 'bigint') {
          n = this.encodeInt64(raw);
        } else if (typeof raw === 'number') {
          n = this.encodeInt64(BigInt(raw));
        } else if (typeof raw === 'boolean') {
          n = this.encodeBool(raw);
        } else if (typeof raw === 'string') {
          n = this.encodeString(raw);
        } else {
          throw new Error(`type error: unsupported type: ${tag.type}`);
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

  // v0.5.0 尾定界：payload（已是 wire 形态）先出，meta 正序字节交给 encodeTag
  // 逆序后置 + 尾部槽，一趟写完，无前缀长度回写。
  private encodeComment(payload: Uint8Array, tag: Tag): number {
    if (!(tag instanceof Tag)) {
      tag = new Tag(tag);
    }

    const meta = tag.toBytes();
    if (meta.length === 0) {
      return 0;
    }

    return this.encodeTag(payload, meta);
  }

  // encodeTag 尾定界写出 tag 包装段（一趟正向，零回写）：
  //
  //   [inner 已编码 wire][reverse(meta 正序)][meta_len 尾部槽][total_len 档位大端][tag sign]
  //
  // bs 为已完成尾定界编码的 inner 字节；meta 为 Tag.toBytes() 正序字节。
  // total = len(bs)+len(meta)+len(尾槽)，与解码侧（footer.ts nodeTailStart）反推式
  // inner = total - ml - len(尾槽) 一致（尾槽与前缀槽为同一套 canonical 三档）。
  private encodeTag(bs: Uint8Array, meta: Uint8Array): number {
    if (meta.length === 0) {
      return 0;
    }

    const slot = footerSlotBytes(meta.length);
    const length = bs.length + meta.length + slot.length;
    if (length > Max4Byte) {
      throw new Error(
        `tag+payload too long, max length: ${Max4Byte}, actual: ${length}`,
      );
    }

    let sign: number = Prefix.Tag;
    let tier: number[] = [];
    if (length < TagLen1Byte) {
      sign |= length;
    } else if (length <= Max1Byte) {
      sign |= TagLen1Byte;
      tier = [length & 0xff];
    } else if (length <= Max2Byte) {
      sign |= TagLen2Byte;
      tier = [(length >> 8) & 0xff, length & 0xff];
    } else if (length <= Max3Byte) {
      sign |= TagLen3Byte;
      tier = [(length >>> 16) & 0xff, (length >> 8) & 0xff, length & 0xff];
    } else {
      // 最后一档必须覆盖到 Max4Byte
      sign |= TagLen4Byte;
      tier = [
        (length >>> 24) & 0xff,
        (length >>> 16) & 0xff,
        (length >>> 8) & 0xff,
        length & 0xff,
      ];
    }

    const savedOffset = this.buffer.offset;
    const rmq = reverseBytes(meta);
    this.buffer.writeBytes(bs);
    this.buffer.writeBytes(rmq);
    this.buffer.writeBytes(slot);
    for (const t of tier) {
      this.buffer.writeUint8(t);
    }
    this.buffer.writeUint8(sign);

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
    this.encodeSimple(SimpleValue.SimpleNull);
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

  encodeU64(value: bigint): number {
    return this.encodeIntWithSign(Prefix.PositiveInt, value);
  }

  // v0.5.0 尾定界：body（大端正序）先写，sign（含档位）追加在尾部。内联档为单字节。
  private encodeIntWithSign(sign: number, value: bigint): number {
    if (value <= 23n) {
      this.buffer.writeUint8(sign | Number(value));
      return 1;
    } else if (value <= Max1Byte) {
      this.buffer.writeUint8(Number(value));
      this.buffer.writeUint8(sign | IntLen1Byte);
      return 2;
    } else if (value <= Max2Byte) {
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen2Byte);
      return 3;
    } else if (value <= Max3Byte) {
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen3Byte);
      return 4;
    } else if (value <= Max4Byte) {
      this.buffer.writeUint8(Number((value >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen4Byte);
      return 5;
    } else if (value <= Max5Byte) {
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen5Byte);
      return 6;
    } else if (value <= Max6Byte) {
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen6Byte);
      return 7;
    } else if (value <= Max7Byte) {
      this.buffer.writeUint8(Number((value >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen7Byte);
      return 8;
    } else if (value <= Max8Byte) {
      this.buffer.writeUint8(Number((value >> 56n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(sign | IntLen8Byte);
      return 9;
    } else {
      throw new Error('Int too large');
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

    // v0.5.0 尾定界：[exp][mantissa 大端正序][sign]。
    if (mantissa <= Max1Byte) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number(mantissa));
      this.buffer.writeUint8(sign | FloatLen1Byte);
      return 1 + 1 + 1;
    } else if (mantissa <= Max2Byte) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen2Byte);
      return 1 + 1 + 2;
    } else if (mantissa <= Max3Byte) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen3Byte);
      return 1 + 1 + 3;
    } else if (mantissa <= Max4Byte) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen4Byte);
      return 1 + 1 + 4;
    } else if (mantissa <= BigInt(Max5Byte)) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen5Byte);
      return 1 + 1 + 5;
    } else if (mantissa <= BigInt(Max6Byte)) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen6Byte);
      return 1 + 1 + 6;
    } else if (mantissa <= BigInt(Max7Byte)) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen7Byte);
      return 1 + 1 + 7;
    } else if (mantissa <= Max8Byte) {
      this.buffer.writeUint8(exponent);
      this.buffer.writeUint8(Number((mantissa >> 56n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 48n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 24n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 16n) & 0xffn));
      this.buffer.writeUint8(Number((mantissa >> 8n) & 0xffn));
      this.buffer.writeUint8(Number(mantissa & 0xffn));
      this.buffer.writeUint8(sign | FloatLen8Byte);
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

  // v0.5.0 尾定界：先写 UTF-8 body，再写长度大端，最后写 sign。
  // 长度字段 = UTF-8 字节数（与 Go 参考实现 len(s) 逐字节对齐）。
  encodeString(value: string): number {
    const bytes = new TextEncoder().encode(value);
    const length = bytes.length;
    // String 维持 v0.2.x：上限 Max2Byte(65535)，最多 2 字节长度字段。
    // 1 字节档边界保持 `< 256`（255 落 2 字节档），与基线一致。
    if (length > Max2Byte) {
      throw new Error(
        `string too long, max length: ${Max2Byte}, actual: ${length}`,
      );
    }

    const savedOffset = this.buffer.offset;
    this.buffer.writeBytes(bytes);
    if (length < StringLen1Byte) {
      this.buffer.writeUint8(Prefix.String | length);
    } else if (length < 256) {
      this.buffer.writeUint8(length);
      this.buffer.writeUint8(Prefix.String | StringLen1Byte);
    } else {
      this.buffer.writeUint16BE(length);
      this.buffer.writeUint8(Prefix.String | StringLen2Byte);
    }

    return this.buffer.offset - savedOffset;
  }

  private encodeStringToBytes(value: string): Uint8Array {
    const startOffset = this.buffer.offset;
    this.encodeString(value);
    return this.buffer.getBytes(startOffset);
  }

  // v0.5.0 尾定界：先写 body，再写长度大端，最后写 sign。
  encodeBytes(value: Uint8Array): number {
    const length = value.length;
    if (length > Max4Byte) {
      throw new Error(
        `bytes too long, max length: ${Max4Byte}, actual: ${length}`,
      );
    }

    const savedOffset = this.buffer.offset;
    this.buffer.writeBytes(value);
    if (length < BytesLen1Byte) {
      this.buffer.writeUint8(Prefix.Bytes | length);
    } else if (length <= Max1Byte) {
      this.buffer.writeUint8(length);
      this.buffer.writeUint8(Prefix.Bytes | BytesLen1Byte);
    } else if (length <= Max2Byte) {
      this.buffer.writeUint8((length >> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(Prefix.Bytes | BytesLen2Byte);
    } else if (length <= Max3Byte) {
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(Prefix.Bytes | BytesLen3Byte);
    } else {
      // 最后一档必须覆盖到 Max4Byte
      this.buffer.writeUint8((length >>> 24) & 0xff);
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(Prefix.Bytes | BytesLen4Byte);
    }

    return this.buffer.offset - savedOffset;
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

  encodeDateTime(date: Date, locationOffset: number = 0): number {
    const v = date.getTime() / 1000 - locationOffset * 3600;
    return this.encodeInt64(BigInt(Math.floor(v)));
  }

  encodeDate(date: Date): number {
    const defaultTime = new Date('1970-01-01T00:00:00Z');
    const diff = date.getTime() - defaultTime.getTime();
    const days = Math.floor(diff / (24 * 60 * 60 * 1000));
    return this.encodeInt64(BigInt(days));
  }

  encodeTime(date: Date): number {
    const utcDate = new Date(date.getTime() + date.getTimezoneOffset() * 60000);
    const seconds =
      utcDate.getHours() * 3600 +
      utcDate.getMinutes() * 60 +
      utcDate.getSeconds();
    return this.encodeU64(BigInt(seconds));
  }

  // v0.5.0 尾定界：先写 payload（子节点 wire 正序），再写长度大端，最后写 sign。
  private encodeArrayToBytes(value: Uint8Array): Uint8Array {
    const startOffset = this.buffer.offset;

    let prefix = Prefix.Container | ContainerArray;
    const length = value.length;
    if (length > Max4Byte) {
      throw new Error(
        `array too long, max length: ${Max4Byte}, actual: ${length}`,
      );
    }

    this.buffer.writeBytes(value);
    if (length < ContainerLen1Byte) {
      this.buffer.writeUint8(prefix | length);
    } else if (length <= Max1Byte) {
      this.buffer.writeUint8(length);
      this.buffer.writeUint8(prefix | ContainerLen1Byte);
    } else if (length <= Max2Byte) {
      this.buffer.writeUint8((length >> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen2Byte);
    } else if (length <= Max3Byte) {
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen3Byte);
    } else {
      // 最后一档必须覆盖到 Max4Byte
      this.buffer.writeUint8((length >>> 24) & 0xff);
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen4Byte);
    }

    return this.buffer.getBytes(startOffset);
  }

  // v0.5.0 尾定界：先写 payload（子节点 wire 正序），再写长度大端，最后写 sign。
  private encodeObjectToBytes(value: Uint8Array): number {
    const length = value.length;

    if (length > Max4Byte) {
      throw new Error(
        `object too long, max length: ${Max4Byte}, actual: ${length}`,
      );
    }

    const prefix = Prefix.Container | ContainerObject;
    const savedOffset = this.buffer.offset;

    this.buffer.writeBytes(value);
    if (length < ContainerLen1Byte) {
      this.buffer.writeUint8(prefix | length);
    } else if (length <= Max1Byte) {
      this.buffer.writeUint8(length);
      this.buffer.writeUint8(prefix | ContainerLen1Byte);
    } else if (length <= Max2Byte) {
      this.buffer.writeUint8((length >> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen2Byte);
    } else if (length <= Max3Byte) {
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen3Byte);
    } else {
      // 最后一档必须覆盖到 Max4Byte
      this.buffer.writeUint8((length >>> 24) & 0xff);
      this.buffer.writeUint8((length >>> 16) & 0xff);
      this.buffer.writeUint8((length >>> 8) & 0xff);
      this.buffer.writeUint8(length & 0xff);
      this.buffer.writeUint8(prefix | ContainerLen4Byte);
    }

    return this.buffer.offset - savedOffset;
  }
}
