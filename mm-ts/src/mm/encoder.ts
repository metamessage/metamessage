import { MMBuffer } from './buffer';
import { Prefix, SimpleValue, StringLen1Byte, StringLen2Byte, BytesLen1Byte, BytesLen2Byte, ContainerLen1Byte, ContainerLen2Byte, ContainerObject, ContainerArray, IntLen1Byte, IntLen2Byte, IntLen3Byte, IntLen4Byte, IntLen5Byte, IntLen6Byte, IntLen7Byte, IntLen8Byte, Max1Byte, Max2Byte, Max3Byte, Max4Byte, Max5Byte, Max6Byte, Max7Byte, Max8Byte } from './constants';
import { MMValue } from './types';

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

  encode(value: any): void {
    if (value === null || value === undefined) {
      this.encodeNil();
    } else if (typeof value === 'boolean') {
      this.encodeBool(value);
    } else if (typeof value === 'number') {
      if (value === 0 && value.toString().includes('.')) {
        this.encodeDouble(value);
      } else if (Number.isInteger(value)) {
        this.encodeInt(BigInt(value));
      } else {
        this.encodeDouble(value);
      }
    } else if (typeof value === 'bigint') {
      this.encodeInt(value);
    } else if (typeof value === 'string') {
      this.encodeString(value);
    } else if (value instanceof Uint8Array) {
      this.encodeBytes(value);
    } else if (Array.isArray(value)) {
      this.encodeArray(value);
    } else if (typeof value === 'object') {
      if ('value' in value && 'options' in value) {
        this.encodeMMValue(value);
      } else {
        this.encodeObject(value);
      }
    } else {
      throw new Error(`Unsupported type: ${typeof value}`);
    }
  }

  private encodeMMValue(value: MMValue<any>): void {
    const { value: val, options } = value;
    if (options.type === 'float') {
      this.encodeDouble(val);
    } else if (options.type === 'int') {
      this.encodeInt(BigInt(val));
    } else if (options.type === 'str') {
      this.encodeString(val);
    } else if (options.type === 'bool') {
      this.encodeBool(val);
    } else if (options.type === 'bytes') {
      this.encodeBytes(val);
    } else if (options.type === 'array') {
      this.encodeArray(val);
    } else if (options.type === 'struct') {
      this.encodeObject(val);
    } else {
      this.encode(val);
    }
  }

  encodeBool(value: boolean): void {
    const prefix = Prefix.Simple;
    const suffix = value ? SimpleValue.True : SimpleValue.False;
    this.buffer.writeUint8(prefix | suffix);
  }

  encodeNil(): void {
    const prefix = Prefix.Simple;
    const suffix = SimpleValue.NullInt;
    this.buffer.writeUint8(prefix | suffix);
  }

  encodeInt(value: bigint): void {
    if (value >= 0) {
      this.encodeUInt(value);
    } else {
      this.encodeNegativeInt(-value);
    }
  }

  private encodeUInt(value: bigint): void {
    if (value < 0) {
      throw new Error('UInt cannot be negative');
    }

    if (value <= 23n) {
      const prefix = Prefix.PositiveInt;
      const suffix = Number(value);
      this.buffer.writeUint8(prefix | suffix);
    } else if (value <= Max1Byte) {
      const prefix = Prefix.PositiveInt | IntLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value));
    } else if (value <= Max2Byte) {
      const prefix = Prefix.PositiveInt | IntLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number(value >> 8n));
    } else if (value <= Max3Byte) {
      const prefix = Prefix.PositiveInt | IntLen3Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
    } else if (value <= Max4Byte) {
      const prefix = Prefix.PositiveInt | IntLen4Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value));
    } else if (value <= Max5Byte) {
      const prefix = Prefix.PositiveInt | IntLen5Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
    } else if (value <= Max6Byte) {
      const prefix = Prefix.PositiveInt | IntLen6Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
    } else if (value <= Max7Byte) {
      const prefix = Prefix.PositiveInt | IntLen7Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 48n) & 0xffn));
    } else if (value <= Max8Byte) {
      const prefix = Prefix.PositiveInt | IntLen8Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint64LE(value);
    } else {
      throw new Error('UInt too large');
    }
  }

  private encodeNegativeInt(value: bigint): void {
    if (value <= 23n) {
      const prefix = Prefix.NegativeInt;
      const suffix = Number(value);
      this.buffer.writeUint8(prefix | suffix);
    } else if (value <= Max1Byte) {
      const prefix = Prefix.NegativeInt | IntLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value));
    } else if (value <= Max2Byte) {
      const prefix = Prefix.NegativeInt | IntLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number(value >> 8n));
    } else if (value <= Max3Byte) {
      const prefix = Prefix.NegativeInt | IntLen3Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(Number(value & 0xffn));
      this.buffer.writeUint8(Number((value >> 8n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 16n) & 0xffn));
    } else if (value <= Max4Byte) {
      const prefix = Prefix.NegativeInt | IntLen4Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value));
    } else if (value <= Max5Byte) {
      const prefix = Prefix.NegativeInt | IntLen5Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
    } else if (value <= Max6Byte) {
      const prefix = Prefix.NegativeInt | IntLen6Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
    } else if (value <= Max7Byte) {
      const prefix = Prefix.NegativeInt | IntLen7Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint32LE(Number(value & 0xffffffffn));
      this.buffer.writeUint8(Number((value >> 32n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 40n) & 0xffn));
      this.buffer.writeUint8(Number((value >> 48n) & 0xffn));
    } else if (value <= Max8Byte) {
      const prefix = Prefix.NegativeInt | IntLen8Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint64LE(value);
    } else {
      throw new Error('NegativeInt too large');
    }
  }

  encodeFloat(value: number): void {
    const prefix = Prefix.Float | 27; // FloatLen4Byte
    this.buffer.writeUint8(prefix);
    this.buffer.writeFloat32(value);
  }

  encodeDouble(value: number): void {
    const prefix = Prefix.Float | 31; // FloatLen8Byte
    this.buffer.writeUint8(prefix);
    this.buffer.writeFloat64(value);
  }

  encodeString(value: string): void {
    const length = value.length;
    if (length < StringLen1Byte) {
      const prefix = Prefix.String | length;
      this.buffer.writeUint8(prefix);
    } else if (length < 256) {
      const prefix = Prefix.String | StringLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
    } else {
      const prefix = Prefix.String | StringLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
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
  }

  encodeBytes(value: Uint8Array | number[]): void {
    const length = value.length;
    if (length < BytesLen1Byte) {
      const prefix = Prefix.Bytes | length;
      this.buffer.writeUint8(prefix);
    } else if (length < 256) {
      const prefix = Prefix.Bytes | BytesLen1Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint8(length);
    } else {
      const prefix = Prefix.Bytes | BytesLen2Byte;
      this.buffer.writeUint8(prefix);
      this.buffer.writeUint16LE(length);
    }
    this.buffer.writeBytes(value);
  }

  encodeArray(value: any[]): void {
    const length = value.length;
    let prefix = Prefix.Container | ContainerArray;
    
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
    
    for (const item of value) {
      this.encode(item);
    }
  }

  encodeObject(value: Record<string, any>): void {
    const entries = Object.entries(value);
    const length = entries.length;
    let prefix = Prefix.Container | ContainerObject;
    
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
    
    for (const [key, val] of entries) {
      this.encodeString(key);
      this.encode(val);
    }
  }
}