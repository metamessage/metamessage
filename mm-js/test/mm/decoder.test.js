import { MMEncoder } from '../../src/mm/encoder.js';
import { MMDecoder } from '../../src/mm/decoder.js';

describe('MMDecoder', () => {
  describe('decode nil', () => {
    test('decodes nil', () => {
      const encoder = new MMEncoder();
      encoder.encodeNil();
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('null');
      expect(result.value).toBe(null);
    });
  });

  describe('decode bool', () => {
    test('decodes true', () => {
      const encoder = new MMEncoder();
      encoder.encodeBool(true);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('bool');
      expect(result.value).toBe(true);
    });

    test('decodes false', () => {
      const encoder = new MMEncoder();
      encoder.encodeBool(false);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('bool');
      expect(result.value).toBe(false);
    });
  });

  describe('decode int', () => {
    test('decodes 0', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(0);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('int');
      expect(result.value).toBe(0n);
    });

    test('decodes 23', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(23);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('int');
      expect(result.value).toBe(23n);
    });

    test('decodes large positive int', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(123456);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('int');
      expect(result.value).toBe(123456n);
    });

    test('decodes negative int', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(-7890);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('int');
      expect(result.value).toBe(-7890n);
    });

    test('decodes BigInt', () => {
      const encoder = new MMEncoder();
      encoder.encodeInt(12345678901234567890n);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('int');
      expect(result.value).toBe(12345678901234567890n);
    });
  });

  describe('decode float', () => {
    test('decodes 0.0', () => {
      const encoder = new MMEncoder();
      encoder.encodeFloat64(0.0);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('float');
      expect(result.value).toBe(0.0);
    });

    test('decodes 3.14', () => {
      const encoder = new MMEncoder();
      encoder.encodeFloat64(3.14);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('float');
      expect(result.value).toBeCloseTo(3.14, 5);
    });
  });

  describe('decode string', () => {
    test('decodes empty string', () => {
      const encoder = new MMEncoder();
      encoder.encodeString('');
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('string');
      expect(result.value).toBe('');
    });

    test('decodes hello', () => {
      const encoder = new MMEncoder();
      encoder.encodeString('hello');
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('string');
      expect(result.value).toBe('hello');
    });
  });

  describe('decode bytes', () => {
    test('decodes bytes array', () => {
      const encoder = new MMEncoder();
      encoder.encodeBytes(new Uint8Array([0x01, 0x02, 0x03]));
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('bytes');
      expect(Array.from(result.value)).toEqual([0x01, 0x02, 0x03]);
    });
  });

  describe('decode array', () => {
    test('decodes array of strings', () => {
      const encoder = new MMEncoder();
      encoder.encodeArray(['a', 'b', 'c']);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('array');
      expect(result.value.length).toBe(3);
    });

    test('decodes array of ints', () => {
      const encoder = new MMEncoder();
      encoder.encodeArray([10, 20, 30]);
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('array');
      expect(result.value.length).toBe(3);
    });
  });

  describe('decode object', () => {
    test('decodes simple object', () => {
      const encoder = new MMEncoder();
      encoder.encodeObject({ name: 'test' });
      const decoder = new MMDecoder(encoder.buffer.data);
      const result = decoder.decode();
      expect(result.type).toBe('object');
      expect(result.value.name).toBeDefined();
    });
  });
});