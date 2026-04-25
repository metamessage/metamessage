import { MMBuffer } from './buffer';
import { Prefix, SimpleValue, getPrefix, getSuffix, stringLen, bytesLen, containerLen, intLen, floatLen, isArray } from './constants';

export interface DecodedValue {
  type: string;
  value: any;
}

export class MMDecoder {
  private buffer: MMBuffer;

  constructor(data: Uint8Array | ArrayBuffer | number[]) {
    let uint8Array: Uint8Array;
    if (data instanceof Uint8Array) {
      uint8Array = data;
    } else if (data instanceof ArrayBuffer) {
      uint8Array = new Uint8Array(data);
    } else {
      uint8Array = new Uint8Array(data);
    }
    
    const buffer = new MMBuffer(uint8Array.length);
    buffer.writeBytes(uint8Array);
    buffer.seek(0);
    this.buffer = buffer;
  }

  decode(): DecodedValue {
    const b = this.buffer.readUint8();
    const prefix = getPrefix(b);
    const suffix = getSuffix(b);

    switch (prefix) {
      case Prefix.Simple:
        return this.decodeSimple(suffix);
      case Prefix.PositiveInt:
        return this.decodePositiveInt(suffix);
      case Prefix.NegativeInt:
        return this.decodeNegativeInt(suffix);
      case Prefix.Float:
        return this.decodeFloat(suffix);
      case Prefix.String:
        return this.decodeString(suffix);
      case Prefix.Bytes:
        return this.decodeBytes(suffix);
      case Prefix.Container:
        return this.decodeContainer(b, suffix);
      case Prefix.Tag:
        return this.decodeTag(suffix);
      default:
        throw new Error(`Unknown prefix: ${prefix}`);
    }
  }

  private decodeSimple(value: number): DecodedValue {
    switch (value) {
      case SimpleValue.NullBool:
      case SimpleValue.NullInt:
      case SimpleValue.NullFloat:
      case SimpleValue.NullString:
      case SimpleValue.NullBytes:
        return { type: 'null', value: null };
      case SimpleValue.False:
        return { type: 'bool', value: false };
      case SimpleValue.True:
        return { type: 'bool', value: true };
      default:
        return { type: 'simple', value: value };
    }
  }

  private decodePositiveInt(suffix: number): DecodedValue {
    const [len, value] = intLen(suffix);
    
    if (len === 0) {
      return { type: 'int', value: BigInt(value) };
    }

    let result = 0n;
    for (let i = 0; i < len; i++) {
      const byte = this.buffer.readUint8();
      result |= BigInt(byte) << (BigInt(i) * 8n);
    }
    
    return { type: 'int', value: result };
  }

  private decodeNegativeInt(suffix: number): DecodedValue {
    const [len, value] = intLen(suffix);
    
    if (len === 0) {
      return { type: 'int', value: -BigInt(value) };
    }

    let result = 0n;
    for (let i = 0; i < len; i++) {
      const byte = this.buffer.readUint8();
      result |= BigInt(byte) << (BigInt(i) * 8n);
    }
    
    return { type: 'int', value: -result };
  }

  private decodeFloat(suffix: number): DecodedValue {
    const [len] = floatLen(suffix);
    
    if (len === 0) {
      return { type: 'float', value: 0.0 };
    } else if (len === 4) {
      const value = this.buffer.readFloat32();
      return { type: 'float', value };
    } else if (len === 8) {
      const value = this.buffer.readFloat64();
      return { type: 'float', value };
    } else {
      throw new Error(`Invalid float length: ${len}`);
    }
  }

  private decodeString(suffix: number): DecodedValue {
    const [lenBytes, len] = stringLen(suffix);
    let length = len;
    
    if (lenBytes === 1) {
      length = this.buffer.readUint8();
    } else if (lenBytes === 2) {
      length = this.buffer.readUint16LE();
    }
    
    let result = '';
    let remaining = length;
    
    while (remaining > 0) {
      const byte = this.buffer.readUint8();
      remaining--;
      
      if (byte < 0x80) {
        result += String.fromCharCode(byte);
      } else if (byte < 0xe0) {
        const byte2 = this.buffer.readUint8();
        remaining--;
        const code = ((byte & 0x1f) << 6) | (byte2 & 0x3f);
        result += String.fromCharCode(code);
      } else if (byte < 0xf0) {
        const byte2 = this.buffer.readUint8();
        const byte3 = this.buffer.readUint8();
        remaining -= 2;
        const code = ((byte & 0x0f) << 12) | ((byte2 & 0x3f) << 6) | (byte3 & 0x3f);
        result += String.fromCharCode(code);
      } else {
        const byte2 = this.buffer.readUint8();
        const byte3 = this.buffer.readUint8();
        const byte4 = this.buffer.readUint8();
        remaining -= 3;
        const code = ((byte & 0x07) << 18) | ((byte2 & 0x3f) << 12) | ((byte3 & 0x3f) << 6) | (byte4 & 0x3f);
        result += String.fromCharCode(code);
      }
    }
    
    return { type: 'string', value: result };
  }

  private decodeBytes(suffix: number): DecodedValue {
    const [lenBytes, len] = bytesLen(suffix);
    let length = len;
    
    if (lenBytes === 1) {
      length = this.buffer.readUint8();
    } else if (lenBytes === 2) {
      length = this.buffer.readUint16LE();
    }
    
    const value = this.buffer.readBytes(length);
    return { type: 'bytes', value };
  }

  private decodeContainer(b: number, suffix: number): DecodedValue {
    const [lenBytes, len] = containerLen(suffix);
    let length = len;
    
    if (lenBytes === 1) {
      length = this.buffer.readUint8();
    } else if (lenBytes === 2) {
      length = this.buffer.readUint16LE();
    }
    
    if (isArray(b)) {
      const value: any[] = [];
      for (let i = 0; i < length; i++) {
        const decoded = this.decode();
        value.push(decoded.value);
      }
      return { type: 'array', value };
    } else {
      const value: Record<string, any> = {};
      for (let i = 0; i < length; i++) {
        const keyDecoded = this.decode();
        const valueDecoded = this.decode();
        value[keyDecoded.value] = valueDecoded.value;
      }
      return { type: 'object', value };
    }
  }

  private decodeTag(suffix: number): DecodedValue {
    const [lenBytes, len] = stringLen(suffix);
    let length = len;
    
    if (lenBytes === 1) {
      length = this.buffer.readUint8();
    } else if (lenBytes === 2) {
      length = this.buffer.readUint16LE();
    }
    
    let tag = '';
    let remaining = length;
    
    while (remaining > 0) {
      const byte = this.buffer.readUint8();
      tag += String.fromCharCode(byte);
      remaining--;
    }
    
    const value = this.decode();
    return { type: 'tag', value: { tag, value: value.value } };
  }
}