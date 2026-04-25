export class MMBuffer {
  constructor(capacity = 1024) {
    this.buffer = new Uint8Array(capacity);
    this.offset = 0;
    this._view = new DataView(this.buffer.buffer);
  }

  reset() {
    this.offset = 0;
  }

  get count() {
    return this.offset;
  }

  get data() {
    return this.buffer.slice(0, this.offset);
  }

  get bytes() {
    return this.data;
  }

  ensureCapacity(needed) {
    if (this.offset + needed > this.buffer.length) {
      const newSize = Math.max(this.buffer.length * 2, this.offset + needed);
      const newBuffer = new Uint8Array(newSize);
      newBuffer.set(this.buffer);
      this.buffer = newBuffer;
      this._view = new DataView(this.buffer.buffer);
    }
  }

  write(byte) {
    this.ensureCapacity(1);
    this.buffer[this.offset] = byte;
    this.offset += 1;
  }

  writeBytes(bytes) {
    if (bytes instanceof Uint8Array) {
      this.ensureCapacity(bytes.length);
      this.buffer.set(bytes, this.offset);
      this.offset += bytes.length;
    } else if (ArrayBuffer.isView(bytes)) {
      const uint8 = new Uint8Array(bytes.buffer, bytes.byteOffset, bytes.byteLength);
      this.ensureCapacity(uint8.length);
      this.buffer.set(uint8, this.offset);
      this.offset += uint8.length;
    } else if (Array.isArray(bytes)) {
      this.ensureCapacity(bytes.length);
      for (const b of bytes) {
        this.buffer[this.offset++] = b;
      }
    }
  }

  writeUint8(value) {
    this.write(value & 0xFF);
  }

  writeInt8(value) {
    this.write(value & 0xFF);
  }

  writeUint16LE(value) {
    this.ensureCapacity(2);
    this._view.setUint16(this.offset, value, true);
    this.offset += 2;
  }

  writeInt16LE(value) {
    this.ensureCapacity(2);
    this._view.setInt16(this.offset, value, true);
    this.offset += 2;
  }

  writeUint32LE(value) {
    this.ensureCapacity(4);
    this._view.setUint32(this.offset, value, true);
    this.offset += 4;
  }

  writeInt32LE(value) {
    this.ensureCapacity(4);
    this._view.setInt32(this.offset, value, true);
    this.offset += 4;
  }

  writeFloat32(value) {
    this.ensureCapacity(4);
    this._view.setFloat32(this.offset, value, true);
    this.offset += 4;
  }

  writeFloat64(value) {
    this.ensureCapacity(8);
    this._view.setFloat64(this.offset, value, true);
    this.offset += 8;
  }

  writeBigInt64LE(value) {
    this.ensureCapacity(8);
    this._view.setBigInt64(this.offset, value, true);
    this.offset += 8;
  }

  writeBigUint64LE(value) {
    this.ensureCapacity(8);
    this._view.setBigUint64(this.offset, value, true);
    this.offset += 8;
  }

  peek() {
    if (this.offset >= this.buffer.length) {
      return null;
    }
    return this.buffer[this.offset];
  }

  read() {
    if (this.offset >= this.buffer.length) {
      return null;
    }
    return this.buffer[this.offset++];
  }

  readBytes(count) {
    if (this.offset + count > this.buffer.length) {
      return null;
    }
    const bytes = this.buffer.slice(this.offset, this.offset + count);
    this.offset += count;
    return bytes;
  }

  remaining() {
    return this.buffer.length - this.offset;
  }

  readUint8() {
    if (this.offset >= this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    return this.buffer[this.offset++];
  }

  readInt8() {
    if (this.offset >= this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this.buffer[this.offset++];
    return value << 24 >> 24;
  }

  readUint16LE() {
    if (this.offset + 2 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getUint16(this.offset, true);
    this.offset += 2;
    return value;
  }

  readInt16LE() {
    if (this.offset + 2 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getInt16(this.offset, true);
    this.offset += 2;
    return value;
  }

  readUint32LE() {
    if (this.offset + 4 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getUint32(this.offset, true);
    this.offset += 4;
    return value;
  }

  readInt32LE() {
    if (this.offset + 4 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getInt32(this.offset, true);
    this.offset += 4;
    return value;
  }

  readFloat32() {
    if (this.offset + 4 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getFloat32(this.offset, true);
    this.offset += 4;
    return value;
  }

  readFloat64() {
    if (this.offset + 8 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getFloat64(this.offset, true);
    this.offset += 8;
    return value;
  }

  readBigInt64LE() {
    if (this.offset + 8 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getBigInt64(this.offset, true);
    this.offset += 8;
    return value;
  }

  readBigUint64LE() {
    if (this.offset + 8 > this.buffer.length) {
      throw new Error('Unexpected end of data');
    }
    const value = this._view.getBigUint64(this.offset, true);
    this.offset += 8;
    return value;
  }
}

export class MMError extends Error {
  constructor(message) {
    super(message);
    this.name = 'MMError';
  }
}

export const errors = {
  UnexpectedEndOfData: new MMError('Unexpected end of data'),
  InvalidData: new MMError('Invalid data'),
  InvalidPrefix: new MMError('Invalid prefix'),
  InvalidTag: new MMError('Invalid tag'),
  TypeMismatch: new MMError('Type mismatch'),
  Overflow: new MMError('Overflow')
};