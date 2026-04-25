import { MMBuffer, MMError } from './buffer.js';
import { MMPrefix, MMSimpleValue, MMConstants, intLen, floatLen, stringLen, bytesLen, containerLen, tagLen, isArray, getPrefix } from './constants.js';

export class MMDecoder {
  constructor(data) {
    this.buffer = new MMBuffer(data.length);
    this.buffer.writeBytes(data);
    this.buffer.offset = 0;
  }

  decode() {
    const byte = this.buffer.peek();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const prefix = getPrefix(byte);

    switch (prefix) {
      case MMPrefix.Simple:
        return this.decodeSimple();
      case MMPrefix.PositiveInt:
        return this.decodePositiveInt();
      case MMPrefix.NegativeInt:
        return this.decodeNegativeInt();
      case MMPrefix.PrefixFloat:
        return this.decodeFloat();
      case MMPrefix.PrefixString:
        return this.decodeString();
      case MMPrefix.PrefixBytes:
        return this.decodeBytes();
      case MMPrefix.Container:
        return this.decodeContainer();
      case MMPrefix.PrefixTag:
        return this.decodeTag();
      default:
        throw MMError.InvalidPrefix;
    }
  }

  decodeSimple() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    switch (byte) {
      case MMSimpleValue.TrueValue:
        return { type: 'bool', value: true };
      case MMSimpleValue.FalseValue:
        return { type: 'bool', value: false };
      case MMSimpleValue.NullBool:
      case MMSimpleValue.NullInt:
      case MMSimpleValue.NullFloat:
      case MMSimpleValue.NullString:
      case MMSimpleValue.NullBytes:
        return { type: 'null', value: null };
      default:
        return { type: 'null', value: null };
    }
  }

  decodePositiveInt() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const { extraBytes, len } = intLen(byte);

    let value = 0n;
    for (let i = 0; i < extraBytes; i++) {
      const b = this.buffer.read();
      if (b === null) {
        throw MMError.UnexpectedEndOfData;
      }
      value = (value << BigInt(8)) | BigInt(b);
    }

    if (extraBytes === 0) {
      const suffix = byte & 0x1F;
      value = BigInt(suffix);
    }

    return { type: 'int', value: value };
  }

  decodeNegativeInt() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const { extraBytes, len } = intLen(byte);

    let value = 0n;
    for (let i = 0; i < extraBytes; i++) {
      const b = this.buffer.read();
      if (b === null) {
        throw MMError.UnexpectedEndOfData;
      }
      value = (value << BigInt(8)) | BigInt(b);
    }

    if (extraBytes === 0) {
      const suffix = byte & 0x1F;
      value = BigInt(suffix);
    }

    return { type: 'int', value: -value };
  }

  decodeFloat() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const value = this.buffer.readFloat64();
    return { type: 'float', value };
  }

  decodeString() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const { extraBytes, len } = stringLen(byte);

    let totalLen = len;
    if (extraBytes > 0) {
      const extra = this.buffer.readBytes(extraBytes);
      if (extra === null) {
        throw MMError.UnexpectedEndOfData;
      }
      for (let i = 0; i < extraBytes; i++) {
        totalLen |= extra[i] << (i * 8);
      }
    }

    const strBytes = this.buffer.readBytes(totalLen);
    if (strBytes === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const str = new TextDecoder().decode(strBytes);
    return { type: 'string', value: str };
  }

  decodeBytes() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const { extraBytes, len } = bytesLen(byte);

    let totalLen = len;
    if (extraBytes > 0) {
      const extra = this.buffer.readBytes(extraBytes);
      if (extra === null) {
        throw MMError.UnexpectedEndOfData;
      }
      for (let i = 0; i < extraBytes; i++) {
        totalLen |= extra[i] << (i * 8);
      }
    }

    const data = this.buffer.readBytes(totalLen);
    if (data === null) {
      throw MMError.UnexpectedEndOfData;
    }

    return { type: 'bytes', value: data };
  }

  decodeContainer() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const isArrayContainer = isArray(byte);

    const { extraBytes, len } = containerLen(byte);

    let totalLen = len;
    if (extraBytes > 0) {
      const extra = this.buffer.readBytes(extraBytes);
      if (extra === null) {
        throw MMError.UnexpectedEndOfData;
      }
      for (let i = 0; i < extraBytes; i++) {
        totalLen |= extra[i] << (i * 8);
      }
    }

    if (isArrayContainer) {
      const elements = [];
      for (let i = 0; i < totalLen; i++) {
        const element = this.decode();
        elements.push(element);
      }
      return { type: 'array', value: elements };
    } else {
      const dict = {};
      let i = 0;
      while (this.buffer.offset < this.buffer.buffer.length && i < Math.floor(totalLen / 2)) {
        const keyValue = this.decode();
        const value = this.decode();
        if (keyValue.type === 'string') {
          dict[keyValue.value] = value;
        }
        i++;
      }
      return { type: 'object', value: dict };
    }
  }

  decodeTag() {
    const byte = this.buffer.read();
    if (byte === null) {
      throw MMError.UnexpectedEndOfData;
    }

    const { extraBytes } = tagLen(byte);

    if (extraBytes > 0) {
      const extra = this.buffer.readBytes(extraBytes);
      if (extra === null) {
        throw MMError.UnexpectedEndOfData;
      }
    }

    return this.decode();
  }
}

export class DecodedValue {
  constructor(type, value) {
    this.type = type;
    this.value = value;
  }

  toString() {
    return `DecodedValue(${this.type}, ${JSON.stringify(this.value)})`;
  }
}