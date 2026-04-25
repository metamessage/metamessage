import { encode, decode, mm } from '../../src/mm/index.js';

describe('MM Encoder/Decoder', () => {
  describe('Boolean', () => {
    test('should encode and decode true', () => {
      const encoded = encode(true);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('bool');
      expect(decoded.value).toBe(true);
    });

    test('should encode and decode false', () => {
      const encoded = encode(false);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('bool');
      expect(decoded.value).toBe(false);
    });
  });

  describe('Null', () => {
    test('should encode and decode null', () => {
      const encoded = encode(null);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('null');
      expect(decoded.value).toBeNull();
    });
  });

  describe('Integer', () => {
    test('should encode and decode 0', () => {
      const encoded = encode(0);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe('0');
    });

    test('should encode and decode positive number', () => {
      const encoded = encode(42);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe('42');
    });

    test('should encode and decode negative number', () => {
      const encoded = encode(-42);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe('-42');
    });

    test('should encode and decode bigint', () => {
      const bigValue = 9007199254740991n;
      const encoded = encode(bigValue);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe(bigValue.toString());
    });
  });

  describe('Float', () => {
    test('should encode and decode 0.0', () => {
      const encoded = encode(mm.float(0.0));
      const decoded = decode(encoded);
      expect(decoded.type).toBe('float');
      expect(decoded.value).toBe(0.0);
    });

    test('should encode and decode positive float', () => {
      const encoded = encode(mm.float(3.14));
      const decoded = decode(encoded);
      expect(decoded.type).toBe('float');
      expect(decoded.value).toBeCloseTo(3.14);
    });

    test('should encode and decode negative float', () => {
      const encoded = encode(mm.float(-3.14));
      const decoded = decode(encoded);
      expect(decoded.type).toBe('float');
      expect(decoded.value).toBeCloseTo(-3.14);
    });
  });

  describe('String', () => {
    test('should encode and decode empty string', () => {
      const encoded = encode('');
      const decoded = decode(encoded);
      expect(decoded.type).toBe('string');
      expect(decoded.value).toBe('');
    });

    test('should encode and decode short string', () => {
      const encoded = encode('hello');
      const decoded = decode(encoded);
      expect(decoded.type).toBe('string');
      expect(decoded.value).toBe('hello');
    });

    test('should encode and decode long string', () => {
      const longString = 'a'.repeat(300);
      const encoded = encode(longString);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('string');
      expect(decoded.value).toBe(longString);
    });
  });

  describe('Bytes', () => {
    test('should encode and decode empty bytes', () => {
      const encoded = encode(new Uint8Array());
      const decoded = decode(encoded);
      expect(decoded.type).toBe('bytes');
      expect(decoded.value.length).toBe(0);
    });

    test('should encode and decode bytes', () => {
      const bytes = new Uint8Array([1, 2, 3, 4, 5]);
      const encoded = encode(bytes);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('bytes');
      expect(Array.from(decoded.value)).toEqual([1, 2, 3, 4, 5]);
    });
  });

  describe('Array', () => {
    test('should encode and decode empty array', () => {
      const encoded = encode([]);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('array');
      expect(decoded.value.length).toBe(0);
    });

    test('should encode and decode array of numbers', () => {
      const encoded = encode([1, 2, 3, 4, 5]);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('array');
      expect(decoded.value.length).toBe(5);
      expect(decoded.value[0].toString()).toBe('1');
      expect(decoded.value[1].toString()).toBe('2');
      expect(decoded.value[2].toString()).toBe('3');
      expect(decoded.value[3].toString()).toBe('4');
      expect(decoded.value[4].toString()).toBe('5');
    });

    test('should encode and decode array of mixed types', () => {
      const encoded = encode([1, 'hello', true, null, mm.float(3.14)]);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('array');
      expect(decoded.value.length).toBe(5);
      expect(decoded.value[0].toString()).toBe('1');
      expect(decoded.value[1]).toBe('hello');
      expect(decoded.value[2]).toBe(true);
      expect(decoded.value[3]).toBeNull();
      expect(decoded.value[4]).toBeCloseTo(3.14);
    });
  });

  describe('Object', () => {
    test('should encode and decode empty object', () => {
      const encoded = encode({});
      const decoded = decode(encoded);
      expect(decoded.type).toBe('object');
      expect(Object.keys(decoded.value).length).toBe(0);
    });

    test('should encode and decode object with simple values', () => {
      const obj = { name: 'test', age: 25, active: true };
      const encoded = encode(obj);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('object');
      expect(decoded.value.name).toBe('test');
      expect(decoded.value.age.toString()).toBe('25');
      expect(decoded.value.active).toBe(true);
    });

    test('should encode and decode nested object', () => {
      const obj = { person: { name: 'test', age: 25 } };
      const encoded = encode(obj);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('object');
      expect(decoded.value.person.name).toBe('test');
      expect(decoded.value.person.age.toString()).toBe('25');
    });
  });

  describe('MM Value', () => {
    test('should encode and decode mm() wrapped value', () => {
      const wrapped = mm(42, { type: 'int', desc: 'test' });
      const encoded = encode(wrapped);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe('42');
    });

    test('should encode and decode mm.int()', () => {
      const wrapped = mm.int(123);
      const encoded = encode(wrapped);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('int');
      expect(decoded.value.toString()).toBe('123');
    });

    test('should encode and decode mm.str()', () => {
      const wrapped = mm.str('hello');
      const encoded = encode(wrapped);
      const decoded = decode(encoded);
      expect(decoded.type).toBe('string');
      expect(decoded.value).toBe('hello');
    });
  });
});