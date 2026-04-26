import { MMBuffer } from './buffer.js';
import { MMPrefix, MMSimpleValue, MMConstants, intLen, getPrefix } from './constants.js';
import { mm, isMM, getMMValue } from './types.js';

export class MMEncoder {
  constructor() {
    this.buffer = new MMBuffer();
  }

  encode(value) {
    if (isMM(value)) {
      const opts = value.options;
      const val = value.value;

      if (opts?.isNull) {
        return this.encodeNil();
      }

      switch (opts?.type) {
        case 'bool':
          return this.encodeBool(val);
        case 'int':
        case 'i':
        case 'i8':
        case 'i16':
        case 'i32':
        case 'i64':
        case 'int8':
        case 'int16':
        case 'int32':
        case 'int64':
          return this.encodeInt(val);
        case 'uint':
        case 'u':
        case 'u8':
        case 'u16':
        case 'u32':
        case 'u64':
        case 'uint8':
        case 'uint16':
        case 'uint32':
        case 'uint64':
          return this.encodeUInt(val);
        case 'f32':
        case 'float32':
          return this.encodeFloat32(val);
        case 'f64':
        case 'float':
        case 'float64':
        case 'double':
          return this.encodeFloat64(val);
        case 'str':
        case 'string':
          return this.encodeString(val);
        case 'bytes':
          return this.encodeBytes(val);
        case 'array':
        case 'arr':
          return this.encodeArray(val);
        case 'struct':
        case 'object':
        case 'map':
          return this.encodeObject(val);
        case 'bi':
        case 'bigint':
          return this.encodeBigInt(val);
        default:
          return this.encodeAuto(val);
      }
    }

    return this.encodeAuto(value);
  }

  encodeAuto(value) {
    if (value === null || value === undefined) {
      return this.encodeNil();
    }
    if (typeof value === 'boolean') {
      return this.encodeBool(value);
    }
    if (typeof value === 'number') {
      return this.encodeFloat64(value);
    }
    if (typeof value === 'string') {
      return this.encodeString(value);
    }
    if (value instanceof Uint8Array || Array.isArray(value) && typeof value[0] === 'number') {
      return this.encodeBytes(value);
    }
    if (Array.isArray(value)) {
      return this.encodeArray(value);
    }
    if (typeof value === 'bigint') {
      return this.encodeBigInt(value);
    }
    if (typeof value === 'object') {
      return this.encodeObject(value);
    }
    throw new Error(`Unsupported type: ${typeof value}`);
  }

  encodeNil() {
    this.buffer.write(MMSimpleValue.NullInt);
  }

  encodeBool(value) {
    if (value) {
      this.buffer.write(MMSimpleValue.TrueValue);
    } else {
      this.buffer.write(MMSimpleValue.FalseValue);
    }
  }

  encodeInt(value) {
    const v = BigInt(value);
    if (v >= 0n) {
      return this.encodeUInt(v);
    } else {
      let uv;
      if (v === -BigInt(Number.MIN_SAFE_INTEGER) - 1n) {
        uv = 9223372036854775808n;
      } else {
        uv = BigInt(-value) - 1n;
      }
      return this.encodeUIntWithSign(MMPrefix.NegativeInt, uv);
    }
  }

  encodeUInt(value) {
    return this.encodeUIntWithSign(MMPrefix.PositiveInt, BigInt(value));
  }

  encodeUIntWithSign(sign, value) {
    const v = value;
    const { extraBytes } = this.calcIntExtraBytes(v);

    switch (extraBytes) {
      case 0:
        this.buffer.write(sign | Number(v));
        break;
      case 1:
        this.buffer.write(sign | MMConstants.IntLen1Byte);
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 2:
        this.buffer.write(sign | MMConstants.IntLen2Byte);
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 3:
        this.buffer.write(sign | MMConstants.IntLen3Byte);
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 4:
        this.buffer.write(sign | MMConstants.IntLen4Byte);
        this.buffer.write(Number((v >> 24n) & 0xFFn));
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 5:
        this.buffer.write(sign | MMConstants.IntLen5Byte);
        this.buffer.write(Number((v >> 32n) & 0xFFn));
        this.buffer.write(Number((v >> 24n) & 0xFFn));
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 6:
        this.buffer.write(sign | MMConstants.IntLen6Byte);
        this.buffer.write(Number((v >> 40n) & 0xFFn));
        this.buffer.write(Number((v >> 32n) & 0xFFn));
        this.buffer.write(Number((v >> 24n) & 0xFFn));
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 7:
        this.buffer.write(sign | MMConstants.IntLen7Byte);
        this.buffer.write(Number((v >> 48n) & 0xFFn));
        this.buffer.write(Number((v >> 40n) & 0xFFn));
        this.buffer.write(Number((v >> 32n) & 0xFFn));
        this.buffer.write(Number((v >> 24n) & 0xFFn));
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
      case 8:
        this.buffer.write(sign | MMConstants.IntLen8Byte);
        this.buffer.write(Number((v >> 56n) & 0xFFn));
        this.buffer.write(Number((v >> 48n) & 0xFFn));
        this.buffer.write(Number((v >> 40n) & 0xFFn));
        this.buffer.write(Number((v >> 32n) & 0xFFn));
        this.buffer.write(Number((v >> 24n) & 0xFFn));
        this.buffer.write(Number((v >> 16n) & 0xFFn));
        this.buffer.write(Number((v >> 8n) & 0xFFn));
        this.buffer.write(Number(v & 0xFFn));
        break;
    }
  }

  calcIntExtraBytes(value) {
    const v = value >= 0n ? value : -value - 1n;

    if (v < 0xFFn) {
      return { extraBytes: 0 };
    } else if (v < 0xFFFFn) {
      return { extraBytes: 1 };
    } else if (v < 0xFFFFFFn) {
      return { extraBytes: 2 };
    } else if (v < 0xFFFFFFFFn) {
      return { extraBytes: 3 };
    } else if (v < 0xFFFFFFFFFFn) {
      return { extraBytes: 4 };
    } else if (v < 0xFFFFFFFFFFFFn) {
      return { extraBytes: 5 };
    } else if (v < 0xFFFFFFFFFFFFFFn) {
      return { extraBytes: 6 };
    } else if (v < 0xFFFFFFFFFFFFFFFFn) {
      return { extraBytes: 7 };
    }
    return { extraBytes: 8 };
  }

  encodeFloat32(value) {
    this.buffer.write(MMPrefix.PrefixFloat);
    this.buffer.writeFloat32(Number(value));
  }

  encodeFloat64(value) {
    this.buffer.write(MMPrefix.PrefixFloat);
    this.buffer.writeFloat64(Number(value));
  }

  encodeString(value) {
    const bytes = new TextEncoder().encode(value);
    const len = bytes.length;

    if (len < 254) {
      this.buffer.write(MMPrefix.PrefixString | len);
    } else if (len < 65536) {
      this.buffer.write(MMPrefix.PrefixString | MMConstants.StringLen1Byte);
      this.buffer.writeUint8(len & 0xFF);
    } else {
      this.buffer.write(MMPrefix.PrefixString | MMConstants.StringLen2Byte);
      this.buffer.writeUint8((len >> 8) & 0xFF);
      this.buffer.writeUint8(len & 0xFF);
    }
    this.buffer.writeBytes(bytes);
  }

  encodeBytes(value) {
    let bytes;
    if (value instanceof Uint8Array) {
      bytes = value;
    } else if (Array.isArray(value)) {
      bytes = new Uint8Array(value);
    } else {
      bytes = new Uint8Array(0);
    }
    const len = bytes.length;

    if (len < 254) {
      this.buffer.write(MMPrefix.PrefixBytes | len);
    } else if (len < 65536) {
      this.buffer.write(MMPrefix.PrefixBytes | MMConstants.BytesLen1Byte);
      this.buffer.writeUint8(len & 0xFF);
    } else {
      this.buffer.write(MMPrefix.PrefixBytes | MMConstants.BytesLen2Byte);
      this.buffer.writeUint8((len >> 8) & 0xFF);
      this.buffer.writeUint8(len & 0xFF);
    }
    this.buffer.writeBytes(bytes);
  }

  encodeArray(array) {
    const valBuf = new MMBuffer();

    for (const item of array) {
      const encoder = new MMEncoder();
      encoder.encode(item);
      valBuf.writeBytes(encoder.buffer.data);
    }

    const payload = valBuf.data;
    const len = payload.length;

    if (len < 254) {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerArray | len);
    } else if (len < 65536) {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerArray | MMConstants.ContainerLen1Byte);
      this.buffer.writeUint8(len & 0xFF);
    } else {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerArray | MMConstants.ContainerLen2Byte);
      this.buffer.writeUint8((len >> 8) & 0xFF);
      this.buffer.writeUint8(len & 0xFF);
    }
    this.buffer.writeBytes(payload);
  }

  encodeObject(obj) {
    const valBuf = new MMBuffer();

    for (const [key, value] of Object.entries(obj)) {
      const keyEncoder = new MMEncoder();
      keyEncoder.encodeString(key);
      valBuf.writeBytes(keyEncoder.buffer.data);

      const valueEncoder = new MMEncoder();
      valueEncoder.encode(value);
      valBuf.writeBytes(valueEncoder.buffer.data);
    }

    const payload = valBuf.data;
    const len = payload.length;

    if (len < 254) {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerObject | len);
    } else if (len < 65536) {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerObject | MMConstants.ContainerLen1Byte);
      this.buffer.writeUint8(len & 0xFF);
    } else {
      this.buffer.write(MMPrefix.Container | MMConstants.ContainerObject | MMConstants.ContainerLen2Byte);
      this.buffer.writeUint8((len >> 8) & 0xFF);
      this.buffer.writeUint8(len & 0xFF);
    }
    this.buffer.writeBytes(payload);
  }

  encodeBigInt(value) {
    return this.encodeInt(value);
  }

  get result() {
    return this.buffer.data;
  }
}