import { MMBuffer } from './buffer';
import {
  Prefix,
  SimpleValue,
  getPrefix,
  getSuffix,
  stringLen,
  bytesLen,
  containerLen,
  intLen,
  floatLen,
  isArray,
  TagLen1Byte,
  TagLen2Byte,
  Max1Byte,
  Max2Byte,
  FloatPositiveNegativeMask,
} from './constants';
import { ValueType, typeToString } from '../ast/value-type';
import { Tag } from '../ast/tag';
import { MMValue, MMObject, MMArray, Node } from '../ast/ast';

export class MMDecoder {
  private buffer: MMBuffer;
  private data: Uint8Array;
  private offset: number;

  constructor(data: Uint8Array) {
    this.data = data;
    this.offset = 0;
    this.buffer = new MMBuffer(data.length);
    this.buffer.writeBytes(data);
    this.buffer.seek(0);
  }

  decode(): Node {
    return this.decodeNode(null, '');
  }

  decodeNode(tag: Tag | null, path: string): Node {
    if (this.offset >= this.data.length) {
      throw new Error(
        `Unexpected end of data ${this.offset} ${this.data.length}`,
      );
    }
    console.log('this.data', this.data, this.data.length);
    const b = this.readByte();
    const prefix = getPrefix(b) ?? 0;
    const suffix = getSuffix(b) ?? 0;

    switch (prefix) {
      case Prefix.Tag:
        return this.decodeTag(b, path);
      case Prefix.Simple:
        return this.decodeSimple(b, tag, path);
      case Prefix.PositiveInt:
        return this.decodePositiveInt(b, tag, path);
      case Prefix.NegativeInt:
        return this.decodeNegativeInt(b, tag, path);
      case Prefix.Float:
        return this.decodeFloat(b, tag, path);
      case Prefix.String:
        return this.decodeString(b, tag, path);
      case Prefix.Bytes:
        return this.decodeBytes(b, tag, path);
      case Prefix.Container:
        return this.decodeContainer(b, tag, path);
      default:
        throw new Error(`Unknown prefix: ${prefix}`);
    }
  }

  private readByte(): number {
    if (this.offset >= this.data.length) {
      throw new Error('Unexpected end of data');
    }
    const b = this.data[this.offset] ?? 0;
    this.offset++;
    return b;
  }

  private readBytes(n: number): Uint8Array {
    if (this.offset + n > this.data.length) {
      throw new Error('Unexpected end of data');
    }
    const result = this.data.slice(this.offset, this.offset + n);
    this.offset += n;
    return result;
  }

  private decodeTag(prefix: number, path: string): Node {
    const [l1, l2] = this.tagLen(prefix);
    let length = l2;

    switch (l1) {
      case 0:
        break;
      case 1:
        length = this.readByte();
        break;
      case 2:
        const lenBytes = this.readBytes(2);
        length = (lenBytes[0] ?? 0) | ((lenBytes[1] ?? 0) << 8);
        break;
      default:
        throw new Error('Invalid data');
    }

    const tag = new Tag();
    const b = this.readByte();
    let l = b;

    if (l < 254) {
    } else if (l < 257) {
      l = this.readByte();
    } else {
      const lenBytes = this.readBytes(2);
      l = (lenBytes[0] ?? 0) | ((lenBytes[1] ?? 0) << 8);
    }

    while (l > 0) {
      const n = this.decodeTagBytes(tag);
      if (n === 0) {
        throw new Error('Tag error');
      }
      if (n > l) {
        throw new Error('Tag overflow');
      }
      l -= n;
    }

    console.log('tagtag', tag);

    if (tag.isNull) {
      switch (tag.type) {
        case ValueType.DateTime:
        case ValueType.Date:
        case ValueType.Time:
          return new MMValue(new Date(0), tag);
        case ValueType.Int8:
        case ValueType.Int16:
        case ValueType.Int32:
        case ValueType.Int64:
        case ValueType.Uint:
        case ValueType.Uint8:
        case ValueType.Uint16:
        case ValueType.Uint32:
        case ValueType.Uint64:
        case ValueType.Int:
          return new MMValue(0, tag);
        case ValueType.Float32:
        case ValueType.Float64:
          return new MMValue(0.0, tag);
        case ValueType.Email:
        case ValueType.UUID:
        case ValueType.Decimal:
        case ValueType.URL:
        case ValueType.String:
          return new MMValue('', tag);
        case ValueType.BigInt:
          return new MMValue(BigInt(0), tag);
        case ValueType.Bool:
          return new MMValue(false, tag);
        case ValueType.Bytes:
          return new MMValue(new Uint8Array(0), tag);
        case ValueType.IP:
          return new MMValue(tag.version === 4 ? '0.0.0.0' : '::', tag);
        default:
          return this.decodeNode(tag, path);
      }
    } else {
      return this.decodeNode(tag, path);
      // switch (tag.type) {
      //   case ValueType.DateTime:
      //   case ValueType.Date:
      //   case ValueType.Time:
      //     return new MMValue(new Date(0), tag);
      //   case ValueType.Int8:
      //     console.log('ValueType.Int8ValueType.Int8ValueType.Int8ValueType.Int8')
      //     const int8 = this.readByte();
      //     return new MMValue(int8, tag);
      //   case ValueType.Int16:
      //   case ValueType.Int32:
      //   case ValueType.Int64:
      //   case ValueType.Uint:
      //   case ValueType.Uint8:
      //   case ValueType.Uint16:
      //   case ValueType.Uint32:
      //   case ValueType.Uint64:
      //   case ValueType.Int:
      //     return new MMValue(0, tag);
      //   case ValueType.Float32:
      //   case ValueType.Float64:
      //     return new MMValue(0.0, tag);
      //   case ValueType.Email:
      //   case ValueType.UUID:
      //   case ValueType.Decimal:
      //   case ValueType.URL:
      //   case ValueType.String:
      //     return new MMValue('', tag);
      //   case ValueType.BigInt:
      //     return new MMValue(BigInt(0), tag);
      //   case ValueType.Bool:
      //     return new MMValue(false, tag);
      //   case ValueType.Bytes:
      //     return new MMValue(new Uint8Array(0), tag);
      //   case ValueType.IP:
      //     return new MMValue(tag.version === 4 ? '0.0.0.0' : '::', tag);
      //   default:
      //     return this.decodeNode(tag, path);
      // }
    }
  }

  private tagLen(prefix: number): [number, number] {
    const suffix = prefix & 0x1f;
    if (suffix < TagLen1Byte) {
      return [0, suffix];
    } else if (suffix === TagLen1Byte) {
      return [1, 0];
    } else if (suffix === TagLen2Byte) {
      return [2, 0];
    }
    return [0, 0];
  }

  private decodeTagBytes(tag: Tag): number {
    const b = this.readByte();

    const prefix = b & 0xf8;
    const l = b & 0x07;

    switch (prefix) {
      case 0 << 3:
        tag.isNull = (l & 0x01) === 1;
        if (tag.isNull) {
          tag.nullable = true;
        }
        return 1;

      case 1 << 3:
        tag.example = (l & 0x01) === 1;
        return 1;

      case 2 << 3:
        let lenDesc = l;
        if (lenDesc <= 5) {
          const bs = this.readBytes(lenDesc);
          tag.desc = new TextDecoder('utf-8').decode(bs);
          return 1 + lenDesc;
        } else if (lenDesc <= 0xff) {
          const l2 = this.readByte();
          lenDesc = l2;
          const bs = this.readBytes(lenDesc);
          tag.desc = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + lenDesc;
        } else if (lenDesc <= 0xffff) {
          const l2 = this.readBytes(2);
          if (!l2 || l2.length < 2) {
            throw new Error('读取长度失败，数据不足');
          }
          lenDesc = (l2[0]! << 8) | l2[1]!;
          const bs = this.readBytes(lenDesc);
          tag.desc = new TextDecoder('utf-8').decode(bs);
          return 1 + 2 + lenDesc;
        }
        return 1;

      case 3 << 3:
        const typeB = this.readByte();
        tag.type = typeB;
        return 1 + 1;

      case 4 << 3:
        tag.raw = (l & 0x01) === 1;
        return 1;

      case 5 << 3:
        tag.nullable = (l & 0x01) === 1;
        return 1;

      case 6 << 3:
        tag.allowEmpty = (l & 0x01) === 1;
        return 1;

      case 7 << 3:
        tag.unique = (l & 0x01) === 1;
        return 1;

      case 8 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.default = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.default = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 9 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.min = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.min = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 10 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.max = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.max = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 11 << 3:
        if (l < 8) {
          for (let i = 0; i < l; i++) {
            const byteVal = this.readByte();
            tag.size = (tag.size << 8n) | BigInt(byteVal);
          }
          return 1 + l;
        }
        return 1;

      case 12 << 3:
        tag.type = ValueType.Enum;
        let lenEnum = l;
        if (lenEnum <= 5) {
          const bs = this.readBytes(lenEnum);
          tag.enum = new TextDecoder('utf-8').decode(bs);
          return 1 + lenEnum;
        } else if (lenEnum <= 0xff) {
          const l2 = this.readByte();
          lenEnum = l2;
          const bs = this.readBytes(lenEnum);
          tag.enum = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + lenEnum;
        } else if (lenEnum <= 0xffff) {
          const l2 = this.readBytes(2);
          if (!l2 || l2.length < 2) {
            throw new Error('读取长度失败，数据不足');
          }
          lenEnum = (l2[0]! << 8) | l2[1]!;
          const bs = this.readBytes(lenEnum);
          tag.enum = new TextDecoder('utf-8').decode(bs);
          return 1 + 2 + lenEnum;
        }
        return 1;

      case 13 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.pattern = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.pattern = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 14 << 3:
        const locBs = this.readBytes(l);
        const locNum =
          parseInt(new TextDecoder('utf-8').decode(locBs), 10) || 0;
        tag.location = locNum;
        return 1 + l;

      case 15 << 3:
        if (l < 8) {
          for (let i = 0; i < l; i++) {
            const byteVal = this.readByte();
            tag.version = (tag.version << 8) | byteVal;
          }
          return 1 + l;
        }
        return 1;

      case 16 << 3:
        if (l < 7) {
          tag.mime = this.mimeToString(l);
          return 1;
        } else {
          const l2 = this.readByte();
          tag.mime = this.mimeToString(l2);
          return 1 + 1;
        }

      case 17 << 3:
        let lenChildDesc = l;
        if (lenChildDesc <= 5) {
          const bs = this.readBytes(lenChildDesc);
          tag.childDesc = new TextDecoder('utf-8').decode(bs);
          return 1 + lenChildDesc;
        } else if (lenChildDesc <= 0xff) {
          const l2 = this.readByte();
          lenChildDesc = l2;
          const bs = this.readBytes(lenChildDesc);
          tag.childDesc = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + lenChildDesc;
        } else if (lenChildDesc <= 0xffff) {
          const l2 = this.readBytes(2);
          if (!l2 || l2.length < 2) {
            throw new Error('读取长度失败，数据不足');
          }
          lenChildDesc = (l2[0]! << 8) | l2[1]!;
          const bs = this.readBytes(lenChildDesc);
          tag.childDesc = new TextDecoder('utf-8').decode(bs);
          return 1 + 2 + lenChildDesc;
        }
        return 1;

      case 18 << 3:
        const childTypeB = this.readByte();
        tag.childType = childTypeB;
        return 1 + 1;

      case 19 << 3:
        tag.childRaw = (l & 0x01) === 1;
        return 1;

      case 20 << 3:
        tag.childNullable = (l & 0x01) === 1;
        return 1;

      case 21 << 3:
        tag.childAllowEmpty = (l & 0x01) === 1;
        return 1;

      case 22 << 3:
        tag.childUnique = (l & 0x01) === 1;
        return 1;

      case 23 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.childDefault = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.childDefault = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 24 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.childMin = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.childMin = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 25 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.childMax = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.childMax = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 26 << 3:
        if (l < 8) {
          for (let i = 0; i < l; i++) {
            const byteVal = this.readByte();
            tag.childSize = (tag.childSize << 8) | byteVal;
          }
          return 1 + l;
        }
        return 1;

      case 27 << 3:
        tag.childType = ValueType.Enum;
        let lenChildEnum = l;
        if (lenChildEnum <= 5) {
          const bs = this.readBytes(lenChildEnum);
          tag.childEnum = new TextDecoder('utf-8').decode(bs);
          return 1 + lenChildEnum;
        } else if (lenChildEnum <= 0xff) {
          const l2 = this.readByte();
          lenChildEnum = l2;
          const bs = this.readBytes(lenChildEnum);
          tag.childEnum = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + lenChildEnum;
        } else if (lenChildEnum <= 0xffff) {
          const l2 = this.readBytes(2);
          if (!l2 || l2.length < 2) {
            throw new Error('读取长度失败，数据不足');
          }
          lenChildEnum = (l2[0]! << 8) | l2[1]!;
          const bs = this.readBytes(lenChildEnum);
          tag.childEnum = new TextDecoder('utf-8').decode(bs);
          return 1 + 2 + lenChildEnum;
        }
        return 1;

      case 28 << 3:
        if (l < 7) {
          const bs = this.readBytes(l);
          tag.childPattern = new TextDecoder('utf-8').decode(bs);
          return 1 + l;
        } else {
          const l2 = this.readByte();
          const bs = this.readBytes(l2);
          tag.childPattern = new TextDecoder('utf-8').decode(bs);
          return 1 + 1 + l2;
        }

      case 29 << 3:
        const childLocBs = this.readBytes(l);
        const childLocNum =
          parseInt(new TextDecoder('utf-8').decode(childLocBs), 10) || 0;
        tag.childLocation = childLocNum;
        return 1 + l;

      case 30 << 3:
        if (l < 8) {
          for (let i = 0; i < l; i++) {
            const byteVal = this.readByte();
            tag.childVersion = (tag.childVersion << 8) | byteVal;
          }
          return 1 + l;
        }
        return 1;

      case 31 << 3:
        if (l < 7) {
          tag.childMime = this.mimeToString(l);
          return 1;
        } else {
          const l2 = this.readByte();
          tag.childMime = this.mimeToString(l2);
          return 1 + 1;
        }

      default:
        throw new Error('invalid data');
    }
  }

  private mimeToString(mime: number): string {
    const mimeTypes: Record<number, string> = {
      1: 'application/json',
      2: 'application/octet-stream',
      3: 'image/png',
      4: 'image/jpeg',
      5: 'image/gif',
      6: 'text/plain',
    };
    return mimeTypes[mime] || '';
  }

  private decodeBigInt(data: Uint8Array, n: number): string {
    const bits = this.bytesToBits(data);
    if (bits.length === 0) {
      return '';
    }

    const neg = bits[0] === 1;
    bits.shift();

    let numStr = '';
    let remaining = n;

    while (remaining > 0) {
      if (remaining >= 3 && bits.length >= 10) {
        const num = this.fromBits(bits.slice(0, 10));
        numStr += num.toString().padStart(3, '0');
        bits.splice(0, 10);
        remaining -= 3;
      } else if (remaining >= 2 && bits.length >= 7) {
        const num = this.fromBits(bits.slice(0, 7));
        numStr += num.toString().padStart(2, '0');
        bits.splice(0, 7);
        remaining -= 2;
      } else if (remaining >= 1 && bits.length >= 4) {
        const num = this.fromBits(bits.slice(0, 4));
        numStr += num.toString();
        bits.splice(0, 4);
        remaining -= 1;
      } else {
        break;
      }
    }

    if (neg) {
      numStr = '-' + numStr;
    }

    return numStr;
  }

  private bytesToBits(data: Uint8Array): number[] {
    const bits: number[] = [];
    for (const byte of data) {
      for (let i = 7; i >= 0; i--) {
        bits.push((byte >> i) & 1);
      }
    }
    return bits;
  }

  private fromBits(bits: number[]): number {
    let result = 0;
    for (const bit of bits) {
      result = (result << 1) | bit;
    }
    return result;
  }

  private decodeSimple(prefix: number, tag: Tag | null, path: string): MMValue {
    if (!tag) {
      tag = new Tag();
    }

    const value = prefix & 0x1f;
    let data: any;
    let text: string;

    switch (value) {
      case SimpleValue.False:
        tag.type = ValueType.Bool;
        data = false;
        text = 'false';
        break;
      case SimpleValue.True:
        tag.type = ValueType.Bool;
        data = true;
        text = 'true';
        break;
      case SimpleValue.NullBool:
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Bool;
        }
        data = false;
        text = 'false';
        break;
      case SimpleValue.NullFloat:
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Float64;
        }
        data = 0.0;
        text = '0.0';
        break;
      case SimpleValue.NullInt:
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Int;
        }
        data = 0;
        text = '0';
        break;
      case SimpleValue.NullString:
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.String;
        }
        data = '';
        text = '';
        break;
      case SimpleValue.NullBytes:
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Bytes;
        }
        data = new Uint8Array(0);
        text = '';
        break;
      default:
        throw new Error(`Unsupported simple value: ${typeToString(value)}`);
    }

    return new MMValue(data, tag);
  }

  private decodePositiveInt(
    prefix: number,
    tag: Tag | null,
    path: string,
  ): MMValue {
    const [l1, l2] = intLen(prefix & 0x1f);

    let v = BigInt(l2);
    if (l1 > 0) {
      const bytes = this.readBytes(l1);
      for (let i = 0; i < l1; i++) {
        v |= BigInt(bytes[i] ?? 0) << BigInt(i * 8);
      }
    }

    return this.convertIntValue(v, false, tag);
  }

  private decodeNegativeInt(
    prefix: number,
    tag: Tag | null,
    path: string,
  ): MMValue {
    const [l1, l2] = intLen(prefix & 0x1f);

    let v = BigInt(l2);
    if (l1 > 0) {
      const bytes = this.readBytes(l1);
      for (let i = 0; i < l1; i++) {
        v |= BigInt(bytes[i] ?? 0) << BigInt(i * 8);
      }
    }

    return this.convertIntValue(v, true, tag);
  }

  private convertIntValue(
    v: bigint,
    isNegative: boolean,
    tag: Tag | null,
  ): MMValue {
    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Int;
    }

    let data: any;
    let text: string;

    if (isNegative) {
      text = '-' + v.toString();
    } else {
      text = v.toString();
    }

    switch (tag.type) {
      case ValueType.Int:
        data = isNegative ? -v : v;
        break;

      case ValueType.Int8:
        data = isNegative ? -Number(v) : Number(v);
        break;

      case ValueType.Int16:
        data = isNegative ? -Number(v) : Number(v);
        break;

      case ValueType.Int32:
        data = isNegative ? -Number(v) : Number(v);
        break;

      case ValueType.Int64:
        data = isNegative ? -v : v;
        break;

      case ValueType.Uint:
        data = v;
        break;

      case ValueType.Uint8:
        data = Number(v);
        break;

      case ValueType.Uint16:
        data = Number(v);
        break;

      case ValueType.Uint32:
        data = Number(v);
        break;

      case ValueType.Uint64:
        data = v;
        break;

      case ValueType.DateTime:
        const ts = isNegative ? -Number(v) : Number(v);
        data = new Date(ts * 1000);
        text = data.toISOString();
        break;
      case ValueType.Date:
        const days = isNegative ? -Number(v) : Number(v);
        data = new Date(Date.UTC(1970, 0, 1 + days));
        text = data.toISOString().split('T')[0];
        break;
      case ValueType.Time:
        const seconds = Number(v);
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const secs = seconds % 60;
        data = new Date(Date.UTC(1970, 0, 1, hours, minutes, secs));
        text = data.toTimeString().split(' ')[0];
        break;
      case ValueType.Enum:
        data = Number(v);
        if (tag.enum) {
          const enums = tag.enum.split('|');
          text = enums[Number(v)] || text;
        }
        break;
      default:
        throw new Error(`Unsupported int type: ${typeToString(tag.type)}`);
    }

    return new MMValue(data, tag);
  }

  private decodeFloat(prefix: number, tag: Tag | null, path: string): MMValue {
    const [l1, l2] = floatLen(prefix);
    const p = Prefix.Float;

    let v = 0.0;
    let length = 1;

    if (prefix >= p && prefix <= p + 7) {
      v = (prefix & 0x0f) / 10;
    } else {
      const exp = this.readByte();

      let mantissa: bigint;

      switch (l1) {
        case 0:
          mantissa = BigInt(l2);
          break;
        case 1:
          const bytes1 = this.readBytes(1);
          mantissa = BigInt(bytes1[0] ?? 0);
          break;
        case 2:
          const bytes2 = this.readBytes(2);
          mantissa = (BigInt(bytes2[0] ?? 0) << 8n) | BigInt(bytes2[1] ?? 0);
          break;
        case 3:
          const bytes3 = this.readBytes(3);
          mantissa =
            (BigInt(bytes3[0] ?? 0) << 16n) |
            (BigInt(bytes3[1] ?? 0) << 8n) |
            BigInt(bytes3[2] ?? 0);
          break;
        case 4:
          const bytes4 = this.readBytes(4);
          mantissa =
            (BigInt(bytes4[0] ?? 0) << 24n) |
            (BigInt(bytes4[1] ?? 0) << 16n) |
            (BigInt(bytes4[2] ?? 0) << 8n) |
            BigInt(bytes4[3] ?? 0);
          break;
        case 5:
          const bytes5 = this.readBytes(5);
          mantissa =
            (BigInt(bytes5[0] ?? 0) << 32n) |
            (BigInt(bytes5[1] ?? 0) << 24n) |
            (BigInt(bytes5[2] ?? 0) << 16n) |
            (BigInt(bytes5[3] ?? 0) << 8n) |
            BigInt(bytes5[4] ?? 0);
          break;
        case 6:
          const bytes6 = this.readBytes(6);
          mantissa =
            (BigInt(bytes6[0] ?? 0) << 40n) |
            (BigInt(bytes6[1] ?? 0) << 32n) |
            (BigInt(bytes6[2] ?? 0) << 24n) |
            (BigInt(bytes6[3] ?? 0) << 16n) |
            (BigInt(bytes6[4] ?? 0) << 8n) |
            BigInt(bytes6[5] ?? 0);
          break;
        case 7:
          const bytes7 = this.readBytes(7);
          mantissa =
            (BigInt(bytes7[0] ?? 0) << 48n) |
            (BigInt(bytes7[1] ?? 0) << 40n) |
            (BigInt(bytes7[2] ?? 0) << 32n) |
            (BigInt(bytes7[3] ?? 0) << 24n) |
            (BigInt(bytes7[4] ?? 0) << 16n) |
            (BigInt(bytes7[5] ?? 0) << 8n) |
            BigInt(bytes7[6] ?? 0);
          break;
        case 8:
          const bytes8 = this.readBytes(8);
          mantissa =
            (BigInt(bytes8[0] ?? 0) << 56n) |
            (BigInt(bytes8[1] ?? 0) << 48n) |
            (BigInt(bytes8[2] ?? 0) << 40n) |
            (BigInt(bytes8[3] ?? 0) << 32n) |
            (BigInt(bytes8[4] ?? 0) << 24n) |
            (BigInt(bytes8[5] ?? 0) << 16n) |
            (BigInt(bytes8[6] ?? 0) << 8n) |
            BigInt(bytes8[7] ?? 0);
          break;
        default:
          throw new Error(`unsupported numerical length: ${l2}`);
      }

      const decimalStr = this.mantissaToDecimal(mantissa, (exp << 24) >> 24);
      v = parseFloat(decimalStr);
      length = l1 + 2;
    }

    if (prefix & FloatPositiveNegativeMask) {
      v = -v;
    }

    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Float64;
    }

    let data: any;
    let text: string;

    switch (tag.type) {
      case ValueType.Float32:
        data = v;
        text = v.toString();
        break;
      case ValueType.Float64:
        data = v;
        text = v.toString();
        break;
      case ValueType.Decimal:
        text = v.toString();
        data = text;
        break;
      default:
        throw new Error(`unsupported value types: ${tag.type}`);
    }

    return new MMValue(data, tag);
  }

  private mantissaToDecimal(mantissa: bigint, exp: number): string {
    const numStr = mantissa.toString();
    const decimalPos = numStr.length + exp;

    let result: string;
    if (decimalPos <= 0) {
      result = '0.' + '0'.repeat(-decimalPos) + numStr;
    } else if (decimalPos > 0 && decimalPos < numStr.length) {
      result =
        numStr.substring(0, decimalPos) + '.' + numStr.substring(decimalPos);
    } else {
      const trailingZeros = decimalPos - numStr.length;
      result = numStr + '0'.repeat(trailingZeros);
    }

    return result;
  }

  private decodeString(prefix: number, tag: Tag | null, path: string): MMValue {
    const [lenBytes, len] = stringLen(prefix & 0x1f);
    let length = len;

    if (lenBytes === 1) {
      length = this.readByte();
    } else if (lenBytes === 2) {
      const bytes = this.readBytes(2);
      length = (bytes[0] ?? 0) | ((bytes[1] ?? 0) << 8);
    }

    const bs = this.readBytes(length);
    const text = new TextDecoder('utf-8').decode(bs);

    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.String;
    }

    let data: any;

    switch (tag.type) {
      case ValueType.Email:
      case ValueType.String:
        data = text;
        break;
      case ValueType.URL:
        data = text;
        break;
      case ValueType.IP:
        data = text;
        break;
      default:
        throw new Error(`Unsupported string type: ${typeToString(tag.type)}`);
    }

    return new MMValue(data, tag);
  }

  private decodeBytes(prefix: number, tag: Tag | null, path: string): MMValue {
    const [lenBytes, len] = bytesLen(prefix & 0x1f);
    let length = len;

    if (lenBytes === 1) {
      length = this.readByte();
    } else if (lenBytes === 2) {
      const bytes = this.readBytes(2);
      length = (bytes[0] ?? 0) | ((bytes[1] ?? 0) << 8);
    }

    const bs = this.readBytes(length);

    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Bytes;
    }

    let data: any;
    let text = '';

    switch (tag.type) {
      case ValueType.BigInt:
        const lenB = bs[0] || 0;
        text = this.decodeBigInt(bs.slice(1), lenB);
        data = BigInt(text);
        break;
      case ValueType.Bytes:
        data = bs;
        text = Buffer.from(bs).toString('base64');
        break;
      case ValueType.UUID:
        data = bs;
        text = this.bytesToUUIDString(bs);
        break;
      case ValueType.IP:
        data = bs;
        text = this.bytesToIPString(bs);
        break;
      default:
        throw new Error(`Unsupported bytes type: ${tag.type}`);
    }

    return new MMValue(data, tag);
  }

  private bytesToUUIDString(bytes: Uint8Array): string {
    const hex = Array.from(bytes)
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('');
    return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
  }

  private bytesToIPString(bytes: Uint8Array): string {
    if (bytes.length === 4) {
      return bytes.join('.');
    } else if (bytes.length === 16) {
      const parts: string[] = [];
      for (let i = 0; i < 8; i++) {
        const part = bytes.slice(i * 2, i * 2 + 2);
        parts.push(
          part
            .reduce((acc, b, idx) => acc + (b << (idx === 0 ? 8 : 0)), 0)
            .toString(16),
        );
      }
      return parts.join(':').replace(/(:0)+:/g, '::');
    }
    return bytes.join('.');
  }

  private decodeContainer(prefix: number, tag: Tag | null, path: string): Node {
    if (isArray(prefix)) {
      return this.decodeArray(prefix, tag, path);
    }
    return this.decodeObject(prefix, tag, path);
  }

  private decodeArray(prefix: number, tag: Tag | null, path: string): MMArray {
    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = tag.size > 0 ? ValueType.Array : ValueType.Slice;
    }

    const [l1, l2] = containerLen(prefix & 0x1f);
    let length = l2;

    if (l1 === 1) {
      length = this.readByte();
    } else if (l1 === 2) {
      const bytes = this.readBytes(2);
      length = (bytes[0] ?? 0) | ((bytes[1] ?? 0) << 8);
    }

    const arr = new MMArray();
    arr.setTag(tag);
    arr.setPath(path);

    for (let i = 0; i < length; i++) {
      const itemTag = new Tag();
      itemTag.inherit(tag);
      const itemPath = `${path}[${i}]`;
      const item = this.decodeNode(itemTag, itemPath);
      arr.addElement(item);
    }

    return arr;
  }

  private decodeObject(
    prefix: number,
    tag: Tag | null,
    path: string,
  ): MMObject {
    if (!tag) {
      tag = new Tag();
    }
    if (tag.type === ValueType.Unknown) {
      tag.type = ValueType.Object;
    }

    const [l1, l2] = containerLen(prefix & 0x1f);
    let length = l2;

    if (l1 === 1) {
      length = this.readByte();
    } else if (l1 === 2) {
      const bytes = this.readBytes(2);
      length = (bytes[0] ?? 0) | ((bytes[1] ?? 0) << 8);
    }

    const lArray = this.readByte();
    const keysNode = this.decodeArray(lArray, tag, path);

    const obj = new MMObject();
    obj.setTag(tag);
    obj.setPath(path);

    const keys = keysNode.getElements().map((k) => (k as MMValue).getValue());

    for (let i = 0; i < keys.length; i++) {
      const key = keys[i] as string;
      const fieldTag = new Tag();
      fieldTag.inherit(tag);
      const fieldPath = path ? `${path}.${key}` : key;
      const value = this.decodeNode(fieldTag, fieldPath);
      obj.setProperty(key, value);
    }

    return obj;
  }
}
