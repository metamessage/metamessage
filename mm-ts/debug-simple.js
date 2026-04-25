// 直接在脚本中实现一个简单的对象编码调试器

// 模拟 MMBuffer 类
class MMBuffer {
  constructor(initialCapacity = 64) {
    this.buffer = new Uint8Array(initialCapacity);
    this.position = 0;
    this.capacity = initialCapacity;
  }

  writeUint8(value) {
    if (this.position >= this.capacity) {
      this.grow();
    }
    this.buffer[this.position++] = value;
  }

  writeUint16LE(value) {
    if (this.position + 1 >= this.capacity) {
      this.grow();
    }
    this.buffer[this.position++] = value & 0xff;
    this.buffer[this.position++] = (value >> 8) & 0xff;
  }

  writeUint32LE(value) {
    if (this.position + 3 >= this.capacity) {
      this.grow();
    }
    this.buffer[this.position++] = value & 0xff;
    this.buffer[this.position++] = (value >> 8) & 0xff;
    this.buffer[this.position++] = (value >> 16) & 0xff;
    this.buffer[this.position++] = (value >> 24) & 0xff;
  }

  writeFloat32(value) {
    if (this.position + 3 >= this.capacity) {
      this.grow();
    }
    const view = new DataView(this.buffer.buffer, this.position);
    view.setFloat32(0, value, true);
    this.position += 4;
  }

  writeFloat64(value) {
    if (this.position + 7 >= this.capacity) {
      this.grow();
    }
    const view = new DataView(this.buffer.buffer, this.position);
    view.setFloat64(0, value, true);
    this.position += 8;
  }

  writeBytes(value) {
    if (this.position + value.length >= this.capacity) {
      this.grow(this.position + value.length);
    }
    for (let i = 0; i < value.length; i++) {
      this.buffer[this.position++] = value[i];
    }
  }

  writeUint64LE(value) {
    if (this.position + 7 >= this.capacity) {
      this.grow();
    }
    for (let i = 0; i < 8; i++) {
      this.buffer[this.position++] = Number((value >> (BigInt(i) * 8n)) & 0xffn);
    }
  }

  grow(minCapacity) {
    const newCapacity = minCapacity || this.capacity * 2;
    const newBuffer = new Uint8Array(newCapacity);
    newBuffer.set(this.buffer.subarray(0, this.position));
    this.buffer = newBuffer;
    this.capacity = newCapacity;
  }

  get result() {
    return this.buffer.subarray(0, this.position);
  }

  reset() {
    this.position = 0;
  }
}

// 常量定义
const Prefix = {
  Simple: 0b00000000,
  PositiveInt: 0b00100000,
  NegativeInt: 0b01000000,
  Float: 0b01100000,
  String: 0b10000000,
  Bytes: 0b10100000,
  Container: 0b11000000,
  Tag: 0b11100000
};

const SimpleValue = {
  False: 0b00000,
  True: 0b00001,
  NullInt: 0b00010
};

const StringLen1Byte = 30;
const StringLen2Byte = 31;

const BytesLen1Byte = 30;
const BytesLen2Byte = 31;

const ContainerMap = 0b00000;
const ContainerArray = 0b10000;
const ContainerLen1Byte = 14;
const ContainerLen2Byte = 15;

const IntLen1Byte = 24;
const IntLen2Byte = 25;
const IntLen3Byte = 26;
const IntLen4Byte = 27;
const IntLen5Byte = 28;
const IntLen6Byte = 29;
const IntLen7Byte = 30;
const IntLen8Byte = 31;

const Max1Byte = 0xff;
const Max2Byte = 0xffff;
const Max3Byte = 0xffffff;
const Max4Byte = 0xffffffff;
const Max5Byte = 0xffffffffffn;
const Max6Byte = 0xffffffffffffn;
const Max7Byte = 0xffffffffffffffn;
const Max8Byte = 0xffffffffffffffffn;

// 简单的编码器类
class MMEncoder {
  constructor() {
    this.buffer = new MMBuffer();
  }

  get result() {
    return this.buffer.result;
  }

  reset() {
    this.buffer.reset();
  }

  encode(value) {
    if (value === null || value === undefined) {
      this.encodeNil();
    } else if (typeof value === 'boolean') {
      this.encodeBool(value);
    } else if (typeof value === 'number') {
      if (Number.isInteger(value)) {
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
      this.encodeObject(value);
    } else {
      throw new Error(`Unsupported type: ${typeof value}`);
    }
  }

  encodeBool(value) {
    const prefix = Prefix.Simple;
    const suffix = value ? SimpleValue.True : SimpleValue.False;
    this.buffer.writeUint8(prefix | suffix);
  }

  encodeNil() {
    const prefix = Prefix.Simple;
    const suffix = SimpleValue.NullInt;
    this.buffer.writeUint8(prefix | suffix);
  }

  encodeInt(value) {
    if (value >= 0) {
      this.encodeUInt(value);
    } else {
      this.encodeNegativeInt(-value);
    }
  }

  encodeUInt(value) {
    if (value < 0) {
      throw new Error('UInt cannot be negative');
    }

    if (value <= 31n) {
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

  encodeNegativeInt(value) {
    if (value <= 31n) {
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

  encodeFloat(value) {
    const prefix = Prefix.Float | 27; // FloatLen4Byte
    this.buffer.writeUint8(prefix);
    this.buffer.writeFloat32(value);
  }

  encodeDouble(value) {
    const prefix = Prefix.Float | 31; // FloatLen8Byte
    this.buffer.writeUint8(prefix);
    this.buffer.writeFloat64(value);
  }

  encodeString(value) {
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

  encodeBytes(value) {
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

  encodeArray(value) {
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

  encodeObject(value) {
    const entries = Object.entries(value);
    const length = entries.length;
    let prefix = Prefix.Container | ContainerMap;
    
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
    
    console.log('Encoding object entries:', entries);
    for (const [key, val] of entries) {
      console.log('Encoding key:', key);
      this.encodeString(key);
      console.log('Encoding value:', val);
      this.encode(val);
    }
  }
}

// 简单的解码器类
class MMDecoder {
  constructor(data) {
    if (Array.isArray(data)) {
      data = new Uint8Array(data);
    } else if (data instanceof ArrayBuffer) {
      data = new Uint8Array(data);
    }
    this.buffer = new MMBuffer(data.length);
    this.buffer.writeBytes(data);
    this.buffer.reset();
  }

  decode() {
    const b = this.buffer.readUint8();
    const prefix = b & 0b11100000;
    const suffix = b & 0b00011111;

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

  decodeSimple(suffix) {
    switch (suffix) {
      case SimpleValue.False:
        return { type: 'bool', value: false };
      case SimpleValue.True:
        return { type: 'bool', value: true };
      case SimpleValue.NullInt:
        return { type: 'null', value: null };
      default:
        throw new Error(`Unknown simple value: ${suffix}`);
    }
  }

  decodePositiveInt(suffix) {
    if (suffix < 24) {
      return { type: 'int', value: BigInt(suffix) };
    }

    let length = suffix - 24 + 1;
    let result = 0n;
    for (let i = 0; i < length; i++) {
      const byte = this.buffer.readUint8();
      result |= BigInt(byte) << (BigInt(i) * 8n);
    }
    
    return { type: 'int', value: result };
  }

  decodeNegativeInt(suffix) {
    if (suffix < 24) {
      return { type: 'int', value: -BigInt(suffix) };
    }

    let length = suffix - 24 + 1;
    let result = 0n;
    for (let i = 0; i < length; i++) {
      const byte = this.buffer.readUint8();
      result |= BigInt(byte) << (BigInt(i) * 8n);
    }
    
    return { type: 'int', value: -result };
  }

  decodeFloat(suffix) {
    const length = suffix - 24 + 1;
    if (length === 4) {
      const value = this.buffer.readFloat32();
      return { type: 'float', value };
    } else if (length === 8) {
      const value = this.buffer.readFloat64();
      return { type: 'float', value };
    } else {
      throw new Error(`Unknown float length: ${length}`);
    }
  }

  decodeString(suffix) {
    let length = suffix;
    if (suffix === 30) {
      length = this.buffer.readUint8();
    } else if (suffix === 31) {
      length = this.buffer.readUint16LE();
    }

    let result = '';
    let remaining = length;
    while (remaining > 0) {
      const byte = this.buffer.readUint8();
      remaining--;
      if (byte < 0x80) {
        result += String.fromCharCode(byte);
      } else if (byte < 0xc0) {
        throw new Error('Invalid UTF-8');
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

  decodeBytes(suffix) {
    let length = suffix;
    if (suffix === 30) {
      length = this.buffer.readUint8();
    } else if (suffix === 31) {
      length = this.buffer.readUint16LE();
    }

    const result = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
      result[i] = this.buffer.readUint8();
    }
    
    return { type: 'bytes', value: result };
  }

  decodeContainer(b, suffix) {
    let length = suffix;
    if (suffix === 14) {
      length = this.buffer.readUint8();
    } else if (suffix === 15) {
      length = this.buffer.readUint16LE();
    }

    const isArray = (b & 0b10000) === 0b10000;
    if (isArray) {
      const value = [];
      for (let i = 0; i < length; i++) {
        const decoded = this.decode();
        value.push(decoded);
      }
      return { type: 'array', value };
    } else {
      const value = {};
      for (let i = 0; i < length; i++) {
        const keyDecoded = this.decode();
        const valueDecoded = this.decode();
        console.log('Decoded key:', keyDecoded);
        console.log('Decoded value:', valueDecoded);
        value[keyDecoded.value] = valueDecoded.value;
      }
      return { type: 'object', value };
    }
  }

  decodeTag(suffix) {
    let length = suffix;
    if (suffix === 30) {
      length = this.buffer.readUint8();
    } else if (suffix === 31) {
      length = this.buffer.readUint16LE();
    }

    let tag = '';
    let remaining = length;
    while (remaining > 0) {
      const byte = this.buffer.readUint8();
      remaining--;
      if (byte < 0x80) {
        tag += String.fromCharCode(byte);
      } else if (byte < 0xc0) {
        throw new Error('Invalid UTF-8');
      } else if (byte < 0xe0) {
        const byte2 = this.buffer.readUint8();
        remaining--;
        const code = ((byte & 0x1f) << 6) | (byte2 & 0x3f);
        tag += String.fromCharCode(code);
      } else if (byte < 0xf0) {
        const byte2 = this.buffer.readUint8();
        const byte3 = this.buffer.readUint8();
        remaining -= 2;
        const code = ((byte & 0x0f) << 12) | ((byte2 & 0x3f) << 6) | (byte3 & 0x3f);
        tag += String.fromCharCode(code);
      } else {
        const byte2 = this.buffer.readUint8();
        const byte3 = this.buffer.readUint8();
        const byte4 = this.buffer.readUint8();
        remaining -= 3;
        const code = ((byte & 0x07) << 18) | ((byte2 & 0x3f) << 12) | ((byte3 & 0x3f) << 6) | (byte4 & 0x3f);
        tag += String.fromCharCode(code);
      }
    }

    const value = this.decode();
    return { type: 'tag', value: { tag, value: value.value } };
  }
}

// 添加 MMBuffer 的 read 方法
MMBuffer.prototype.readUint8 = function() {
  if (this.position >= this.capacity) {
    throw new Error('Buffer underflow');
  }
  return this.buffer[this.position++];
};

MMBuffer.prototype.readUint16LE = function() {
  if (this.position + 1 >= this.capacity) {
    throw new Error('Buffer underflow');
  }
  const value = this.buffer[this.position] | (this.buffer[this.position + 1] << 8);
  this.position += 2;
  return value;
};

MMBuffer.prototype.readUint32LE = function() {
  if (this.position + 3 >= this.capacity) {
    throw new Error('Buffer underflow');
  }
  const value = this.buffer[this.position] | (this.buffer[this.position + 1] << 8) | (this.buffer[this.position + 2] << 16) | (this.buffer[this.position + 3] << 24);
  this.position += 4;
  return value;
};

MMBuffer.prototype.readFloat32 = function() {
  if (this.position + 3 >= this.capacity) {
    throw new Error('Buffer underflow');
  }
  const view = new DataView(this.buffer.buffer, this.position);
  const value = view.getFloat32(0, true);
  this.position += 4;
  return value;
};

MMBuffer.prototype.readFloat64 = function() {
  if (this.position + 7 >= this.capacity) {
    throw new Error('Buffer underflow');
  }
  const view = new DataView(this.buffer.buffer, this.position);
  const value = view.getFloat64(0, true);
  this.position += 8;
  return value;
};

// 测试对象编码
const obj = { name: 'test', age: 25, active: true };
console.log('Original object:', obj);

const encoder = new MMEncoder();
encoder.encode(obj);
const encoded = encoder.result;
console.log('Encoded bytes:', Array.from(encoded));

const decoder = new MMDecoder(encoded);
const decoded = decoder.decode();
console.log('Decoded object:', decoded);
console.log('Decoded age:', decoded.value.age);
console.log('Decoded age type:', typeof decoded.value.age);
console.log('Decoded age toString():', decoded.value.age.toString());

// 测试单个字段编码
console.log('\nDebugging individual fields:');

// 测试 name 字段
const nameEncoder = new MMEncoder();
nameEncoder.encode('name');
const nameEncoded = nameEncoder.result;
console.log('Encoded "name":', Array.from(nameEncoded));
const nameDecoder = new MMDecoder(nameEncoded);
const nameDecoded = nameDecoder.decode();
console.log('Decoded "name":', nameDecoded);

// 测试 test 字段
const testEncoder = new MMEncoder();
testEncoder.encode('test');
const testEncoded = testEncoder.result;
console.log('Encoded "test":', Array.from(testEncoded));
const testDecoder = new MMDecoder(testEncoded);
const testDecoded = testDecoder.decode();
console.log('Decoded "test":', testDecoded);

// 测试 age 字段
const ageEncoder = new MMEncoder();
ageEncoder.encode(25);
const ageEncoded = ageEncoder.result;
console.log('Encoded 25:', Array.from(ageEncoded));
const ageDecoder = new MMDecoder(ageEncoded);
const ageDecoded = ageDecoder.decode();
console.log('Decoded 25:', ageDecoded);

// 测试 active 字段
const activeEncoder = new MMEncoder();
activeEncoder.encode(true);
const activeEncoded = activeEncoder.result;
console.log('Encoded true:', Array.from(activeEncoded));
const activeDecoder = new MMDecoder(activeEncoded);
const activeDecoded = activeDecoder.decode();
console.log('Decoded true:', activeDecoded);