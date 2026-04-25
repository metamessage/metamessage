import { MMEncoder } from '../../src/mm/encoder.js';
import { mm } from '../../src/mm/types.js';

describe('MMEncoder', () => {
  describe('encodeNil', () => {
    test('encodes nil as null byte', () => {
      const encoder = new MMEncoder();
      encoder.encodeNil();
      expect(encoder.buffer.data[0]).toBe(1);
    });
  });

  describe('encodeBool', () => {
    test('encodes true', () => {
      const encoder = new MMEncoder();
      encoder.encodeBool(true);
      expect(encoder.buffer.data[0]).toBe(6);
    });

    test('encodes false', () => {
      const encoder = new MMEncoder();
      encoder.encodeBool(false);
      expect(encoder.buffer.data[0]).toBe(5);
    });
  });

  describe('encodeInt', () => {
    test('encodes 0', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(0);
      expect(encoder.buffer.data[0]).toBe(0b001_00000);
    });

    test('encodes small positive int (23)', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(23);
      expect(encoder.buffer.data[0]).toBe(0b001_10111);
    });

    test('encodes small negative int (-1)', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(-1);
      const data = encoder.buffer.data;
      expect(data[0] & 0b111_00000).toBe(0b010_00000);
    });

    test('encodes large positive int', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(123456);
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });

    test('encodes large negative int', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(-7890);
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });

    test('encodes BigInt', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(12345678901234567890n);
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });
  });

  describe('encodeFloat', () => {
    test('encodes 0.0', () => {
      const encoder = new MMEncoder();
      encoder.encodeFloat64(0.0);
      expect(encoder.buffer.data[0]).toBe(0b011_00000);
      expect(encoder.buffer.data.length).toBe(9);
    });

    test('encodes 3.14', () => {
      const encoder = new MMEncoder();
      encoder.encodeFloat64(3.14);
      expect(encoder.buffer.data[0]).toBe(0b011_00000);
      expect(encoder.buffer.data.length).toBe(9);
    });

    test('encodes negative float', () => {
      const encoder = new MMEncoder();
      encoder.encodeFloat64(-3.14);
      expect(encoder.buffer.data[0]).toBe(0b011_00000);
      expect(encoder.buffer.data.length).toBe(9);
    });
  });

  describe('encodeString', () => {
    test('encodes empty string', () => {
      const encoder = new MMEncoder();
      encoder.encodeString('');
      expect(encoder.buffer.data[0]).toBe(0b100_00000);
    });

    test('encodes hello', () => {
      const encoder = new MMEncoder();
      encoder.encodeString('hello');
      const data = encoder.buffer.data;
      expect(data[0] & 0b111_00000).toBe(0b100_00000);
      expect(data.length).toBe(6);
    });

    test('encodes long string', () => {
      const encoder = new MMEncoder();
      const longString = 'a'.repeat(300);
      encoder.encodeString(longString);
      expect(encoder.buffer.data.length).toBeGreaterThan(300);
    });
  });

  describe('encodeBytes', () => {
    test('encodes empty bytes', () => {
      const encoder = new MMEncoder();
      encoder.encodeBytes(new Uint8Array(0));
      expect(encoder.buffer.data[0]).toBe(0b101_00000);
    });

    test('encodes bytes array', () => {
      const encoder = new MMEncoder();
      encoder.encodeBytes(new Uint8Array([0x01, 0x02, 0x03]));
      expect(encoder.buffer.data.length).toBeGreaterThan(3);
    });
  });

  describe('encodeArray', () => {
    test('encodes empty array', () => {
      const encoder = new MMEncoder();
      encoder.encodeArray([]);
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b110_00000);
    });

    test('encodes array of strings', () => {
      const encoder = new MMEncoder();
      encoder.encodeArray(['a', 'b', 'c']);
      expect(encoder.buffer.data.length).toBeGreaterThan(3);
    });

    test('encodes array of ints', () => {
      const encoder = new MMEncoder();
      encoder.encodeArray([10, 20, 30]);
      expect(encoder.buffer.data.length).toBeGreaterThan(3);
    });

    test('encodes large array', () => {
      const encoder = new MMEncoder();
      const largeArray = Array(1000).fill(1);
      encoder.encodeArray(largeArray);
      expect(encoder.buffer.data.length).toBeGreaterThan(1000);
    });
  });

  describe('encodeObject', () => {
    test('encodes empty object', () => {
      const encoder = new MMEncoder();
      encoder.encodeObject({});
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b110_00000);
    });

    test('encodes simple object', () => {
      const encoder = new MMEncoder();
      encoder.encodeObject({ name: 'test', value: 123 });
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });

    test('encodes nested object', () => {
      const encoder = new MMEncoder();
      encoder.encodeObject({ outer: { inner: 'value' } });
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });
  });

  describe('encode with mm wrapper', () => {
    test('mm.int wraps int type', () => {
      const encoder = new MMEncoder();
      encoder.encode(mm.int(42));
      expect(encoder.buffer.data.length).toBeGreaterThan(0);
    });

    test('mm.str wraps string type', () => {
      const encoder = new MMEncoder();
      encoder.encode(mm.str('hello'));
      expect(encoder.buffer.data.length).toBeGreaterThan(5);
    });

    test('mm.bool wraps bool type', () => {
      const encoder = new MMEncoder();
      encoder.encode(mm.bool(true));
      expect(encoder.buffer.data[0]).toBe(6);
    });

    test('mm.bigint wraps bigint type', () => {
      const encoder = new MMEncoder();
      encoder.encode(mm.bigint(12345678901234567890n));
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });

    test('mm.null encodes null', () => {
      const encoder = new MMEncoder();
      encoder.encode(mm.null());
      expect(encoder.buffer.data[0]).toBe(1);
    });
  });

  describe('encodeAuto', () => {
    test('auto-detects boolean', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto(true);
      expect(encoder.buffer.data[0]).toBe(6);
    });

    test('auto-detects number', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto(3.14);
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b011_00000);
    });

    test('auto-detects string', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto('hello');
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b100_00000);
    });

    test('auto-detects array', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto(['a', 'b', 'c']);
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b110_00000);
    });

    test('auto-detects object', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto({ a: 1 });
      expect(encoder.buffer.data[0] & 0b111_00000).toBe(0b110_00000);
    });

    test('auto-detects BigInt', () => {
      const encoder = new MMEncoder();
      encoder.encodeAuto(12345678901234567890n);
      expect(encoder.buffer.data.length).toBeGreaterThan(1);
    });
  });
});