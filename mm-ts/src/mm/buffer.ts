export class MMBuffer {
  private buffer: ArrayBuffer;
  private view: DataView;
  private position: number;
  private capacity: number;

  constructor(initialCapacity: number = 1024) {
    this.capacity = initialCapacity;
    this.buffer = new ArrayBuffer(this.capacity);
    this.view = new DataView(this.buffer);
    this.position = 0;
  }

  get length(): number {
    return this.position;
  }

  get result(): Uint8Array {
    return new Uint8Array(this.buffer, 0, this.position);
  }

  reset(): void {
    this.position = 0;
  }

  private ensureCapacity(required: number): void {
    if (this.position + required <= this.capacity) {
      return;
    }

    const newCapacity = Math.max(this.capacity * 2, this.position + required);
    const newBuffer = new ArrayBuffer(newCapacity);
    const newView = new DataView(newBuffer);
    
    for (let i = 0; i < this.position; i++) {
      newView.setUint8(i, this.view.getUint8(i));
    }

    this.buffer = newBuffer;
    this.view = newView;
    this.capacity = newCapacity;
  }

  writeUint8(value: number): void {
    this.ensureCapacity(1);
    this.view.setUint8(this.position++, value);
  }

  writeInt8(value: number): void {
    this.ensureCapacity(1);
    this.view.setInt8(this.position++, value);
  }

  writeUint16LE(value: number): void {
    this.ensureCapacity(2);
    this.view.setUint16(this.position, value, true);
    this.position += 2;
  }

  writeInt16LE(value: number): void {
    this.ensureCapacity(2);
    this.view.setInt16(this.position, value, true);
    this.position += 2;
  }

  writeUint32LE(value: number): void {
    this.ensureCapacity(4);
    this.view.setUint32(this.position, value, true);
    this.position += 4;
  }

  writeInt32LE(value: number): void {
    this.ensureCapacity(4);
    this.view.setInt32(this.position, value, true);
    this.position += 4;
  }

  writeUint64LE(value: bigint): void {
    this.ensureCapacity(8);
    const low = Number(value & 0xffffffffn);
    const high = Number(value >> 32n);
    this.view.setUint32(this.position, low, true);
    this.view.setUint32(this.position + 4, high, true);
    this.position += 8;
  }

  writeInt64LE(value: bigint): void {
    this.ensureCapacity(8);
    const low = Number(value & 0xffffffffn);
    const high = Number(value >> 32n);
    this.view.setInt32(this.position, low, true);
    this.view.setInt32(this.position + 4, high, true);
    this.position += 8;
  }

  writeFloat32(value: number): void {
    this.ensureCapacity(4);
    this.view.setFloat32(this.position, value, true);
    this.position += 4;
  }

  writeFloat64(value: number): void {
    this.ensureCapacity(8);
    this.view.setFloat64(this.position, value, true);
    this.position += 8;
  }

  writeBytes(bytes: Uint8Array | number[]): void {
    const length = bytes.length;
    this.ensureCapacity(length);
    for (let i = 0; i < length; i++) {
      this.view.setUint8(this.position++, bytes[i] as number);
    }
  }

  readUint8(): number {
    if (this.position >= this.capacity) {
      throw new Error('Buffer underflow');
    }
    return this.view.getUint8(this.position++);
  }

  readInt8(): number {
    if (this.position >= this.capacity) {
      throw new Error('Buffer underflow');
    }
    return this.view.getInt8(this.position++);
  }

  readUint16LE(): number {
    if (this.position + 2 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getUint16(this.position, true);
    this.position += 2;
    return value;
  }

  readInt16LE(): number {
    if (this.position + 2 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getInt16(this.position, true);
    this.position += 2;
    return value;
  }

  readUint32LE(): number {
    if (this.position + 4 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getUint32(this.position, true);
    this.position += 4;
    return value;
  }

  readInt32LE(): number {
    if (this.position + 4 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getInt32(this.position, true);
    this.position += 4;
    return value;
  }

  readUint64LE(): bigint {
    if (this.position + 8 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const low = this.view.getUint32(this.position, true);
    const high = this.view.getUint32(this.position + 4, true);
    this.position += 8;
    return (BigInt(high) << 32n) | BigInt(low);
  }

  readInt64LE(): bigint {
    if (this.position + 8 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const low = this.view.getInt32(this.position, true);
    const high = this.view.getInt32(this.position + 4, true);
    this.position += 8;
    return (BigInt(high) << 32n) | BigInt(low);
  }

  readFloat32(): number {
    if (this.position + 4 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getFloat32(this.position, true);
    this.position += 4;
    return value;
  }

  readFloat64(): number {
    if (this.position + 8 > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const value = this.view.getFloat64(this.position, true);
    this.position += 8;
    return value;
  }

  readBytes(length: number): Uint8Array {
    if (this.position + length > this.capacity) {
      throw new Error('Buffer underflow');
    }
    const result = new Uint8Array(this.buffer, this.position, length);
    this.position += length;
    return result;
  }

  seek(position: number): void {
    if (position < 0 || position > this.capacity) {
      throw new Error('Invalid position');
    }
    this.position = position;
  }

  getPosition(): number {
    return this.position;
  }

  slice(start: number, end?: number): MMBuffer {
    const newBuffer = new MMBuffer(end ? end - start : this.position - start);
    const slice = new Uint8Array(this.buffer, start, end ? end - start : this.position - start);
    newBuffer.writeBytes(slice);
    return newBuffer;
  }
}

export class MMError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'MMError';
  }
}