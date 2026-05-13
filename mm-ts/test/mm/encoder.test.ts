import { encodeNode, decodeToValue, mm } from '../../src/mm/index.js';
import { ValueType } from '../../src/ast/value-type.js';
import { Tag } from '../../src/ast/tag.js';

// npm test test/mm/encoder.test.ts -- -t "should encode and decode array of numbers"
describe('MM Encoder/Decoder', () => {
  describe('Boolean', () => {
    test('should encode and decode true', () => {
      const encoded = encodeNode(mm.bool(true));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Bool);
      expect(decoded.value).toBe(true);
    });

    test('should encode and decode false', () => {
      const encoded = encodeNode(mm.bool(false));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Bool);
      expect(decoded.value).toBe(false);
    });
  });

  describe('Null', () => {
    test('should encode and decode zero int without allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.Int;
      const encoded = encodeNode(mm.i(0n, tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(0n);
    });
  });

  describe('Integer', () => {
    test('should encode and decode non-zero int', () => {
      const encoded = encodeNode(mm.i(42n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(42n);
    });

    test('should encode and decode 0 with allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.Int;
      tag.allowEmpty = true;
      const encoded = encodeNode(mm.i(0n, tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(0n);
    });

    test('should encode and decode positive number', () => {
      const encoded = encodeNode(mm.i(42n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(42n);
    });

    test('should encode and decode negative number', () => {
      const encoded = encodeNode(mm.i(-42n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(-42n);
    });

    test('should encode and decode bigint', () => {
      const bigValue = 9007199254740991n;
      const encoded = encodeNode(mm.bigint(bigValue));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.BigInt);
      expect(decoded.value).toBe(bigValue);
    });

    test('should encode and decode int8', () => {
      const node = mm.i8(127);
      const encoded = encodeNode(node);
      console.log('encoded', encoded);
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int8);
      expect(decoded.value).toBe(127);
    });

    test('should encode and decode int16', () => {
      const encoded = encodeNode(mm.i16(32767));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int16);
      expect(decoded.value).toBe(32767);
    });

    test('should encode and decode int32', () => {
      const encoded = encodeNode(mm.i32(2147483647));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int32);
      expect(decoded.value).toBe(2147483647);
    });

    test('should encode and decode int64', () => {
      const encoded = encodeNode(mm.i64(9223372036854775807n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int64);
      expect(decoded.value).toBe(9223372036854775807n);
    });

    test('should encode and decode uint', () => {
      const encoded = encodeNode(mm.u(42n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Uint);
      expect(decoded.value).toBe(42n);
    });

    test('should encode and decode uint8', () => {
      const encoded = encodeNode(mm.u8(255));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Uint8);
      expect(decoded.value).toBe(255);
    });

    test('should encode and decode uint16', () => {
      const encoded = encodeNode(mm.u16(65535));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Uint16);
      expect(decoded.value).toBe(65535);
    });

    test('should encode and decode uint32', () => {
      const encoded = encodeNode(mm.u32(4294967295));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Uint32);
      expect(decoded.value).toBe(4294967295);
    });

    test('should encode and decode uint64', () => {
      const encoded = encodeNode(mm.u64(18446744073709551615n));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Uint64);
      expect(decoded.value).toBe(18446744073709551615n);
    });
  });

  describe('Float', () => {
    test('should encode and decode non-zero float', () => {
      const encoded = encodeNode(mm.f64(3.14));
      const decoded = decodeToValue(encoded);
      console.log('decoded', decoded);
      expect(decoded.type).toBe(ValueType.Float64);
      expect(decoded.value).toBeCloseTo(3.14);
    });

    test('should encode and decode 0.0 with allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.Float64;
      tag.allowEmpty = true;
      const encoded = encodeNode(mm.f64(0.0, tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Float64);
      expect(decoded.value).toBe(0.0);
    });

    test('should encode and decode positive float', () => {
      const encoded = encodeNode(mm.f64(3.14));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Float64);
      expect(decoded.value).toBeCloseTo(3.14);
    });

    test('should encode and decode negative float', () => {
      const encoded = encodeNode(mm.f64(-3.14));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Float64);
      expect(decoded.value).toBeCloseTo(-3.14);
    });

    test('should encode and decode float32', () => {
      const encoded = encodeNode(mm.f32(3.14));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Float32);
      expect(decoded.value).toBeCloseTo(3.14);
    });
  });

  describe('String', () => {
    test('should encode and decode non-empty string', () => {
      const encoded = encodeNode(mm.str('hello'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.String);
      expect(decoded.value).toBe('hello');
    });

    test('should encode and decode empty string with allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.String;
      tag.allowEmpty = true;
      const encoded = encodeNode(mm.str('', tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.String);
      expect(decoded.value).toBe('');
    });

    test('should encode and decode short string', () => {
      const encoded = encodeNode(mm.str('hello'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.String);
      expect(decoded.value).toBe('hello');
    });

    test('should encode and decode long string', () => {
      const longString = 'a'.repeat(300);
      const encoded = encodeNode(mm.str(longString));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.String);
      expect(decoded.value).toBe(longString);
    });

    test('should encode and decode email', () => {
      const encoded = encodeNode(mm.email('test@example.com'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Email);
      expect(decoded.value).toBe('test@example.com');
    });

    test('should encode and decode url', () => {
      const encoded = encodeNode(mm.url('https://example.com'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.URL);
      expect(decoded.value).toBe('https://example.com');
    });

    test('should encode and decode uuid', () => {
      const encoded = encodeNode(
        mm.uuid('550e8400-e29b-41d4-a716-446655440000'),
      );
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.UUID);
      expect(Array.from(decoded.value)).toEqual([
        85, 14, 132, 0, 226, 155, 65, 212, 167, 22, 68, 102, 85, 68, 0, 0,
      ]);
    });

    test('should encode and decode decimal', () => {
      const encoded = encodeNode(mm.decimal('123.45'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Decimal);
      expect(decoded.value).toBe('123.45');
    });

    test('should encode and decode ip', () => {
      const encoded = encodeNode(mm.ip('192.168.1.1'));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.IP);
      expect(decoded.value).toBe('192.168.1.1');
    });
  });

  describe('Bytes', () => {
    test('should encode and decode non-empty bytes', () => {
      const bytes = new Uint8Array([1, 2, 3, 4, 5]);
      const encoded = encodeNode(mm.bytes(bytes));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Bytes);
      expect(Array.from(decoded.value)).toEqual([1, 2, 3, 4, 5]);
    });

    test('should encode and decode empty bytes with allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.Bytes;
      tag.allowEmpty = true;
      const encoded = encodeNode(mm.bytes(new Uint8Array(), tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Bytes);
      expect(decoded.value.length).toBe(0);
    });

    test('should encode and decode bytes', () => {
      const bytes = new Uint8Array([1, 2, 3, 4, 5]);
      const encoded = encodeNode(mm.bytes(bytes));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Bytes);
      expect(Array.from(decoded.value)).toEqual([1, 2, 3, 4, 5]);
    });
  });

  describe('Array', () => {
    test('should encode and decode non-empty array', () => {
      const encoded = encodeNode(mm.arr([mm.i(1n), mm.i(2n), mm.i(3n)]));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Array);
      expect(decoded.value.length).toBe(3);
    });

    test('should encode and decode empty array with allowEmpty', () => {
      const tag = new Tag();
      tag.type = ValueType.Array;
      tag.allowEmpty = true;
      const encoded = encodeNode(mm.arr([], undefined, tag));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Array);
      expect(decoded.value.length).toBe(0);
    });

    test('should encode and decode array of numbers', () => {
      const encoded = encodeNode(
        mm.arr([mm.i(1n), mm.i(2n), mm.i(3n), mm.i(4n), mm.i(5n)]),
      );
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Array);
      expect(decoded.value.length).toBe(5);
      expect(decoded.value[0].toString()).toBe('1');
      expect(decoded.value[1].toString()).toBe('2');
      expect(decoded.value[2].toString()).toBe('3');
      expect(decoded.value[3].toString()).toBe('4');
      expect(decoded.value[4].toString()).toBe('5');
    });

    // test('should encode and decode array of mixed types without null', () => {
    //   const encoded = encodeNode(
    //     mm.arr([mm.i(1n), mm.str('hello'), mm.bool(true), mm.f64(3.14)]),
    //   );
    //   const decoded = decodeToValue(encoded);
    //   expect(decoded.type).toBe(ValueType.Array);
    //   expect(decoded.value.length).toBe(4);
    //   expect(decoded.value[0].toString()).toBe('1');
    //   expect(decoded.value[1]).toBe('hello');
    //   expect(decoded.value[2]).toBe(true);
    //   expect(decoded.value[3]).toBeCloseTo(3.14);
    // });
  });

  describe('Object', () => {
    test('should encode and decode non-empty object', () => {
      const obj = { name: mm.str('test'), age: mm.i(25n) };
      const encoded = encodeNode(mm.obj(obj));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Object);
      expect(decoded.value.name).toBe('test');
    });

    test('should encode and decode empty object', () => {
      const encoded = encodeNode(mm.obj({}));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Object);
      expect(Object.keys(decoded.value).length).toBe(0);
    });

    test('should encode and decode object with simple values', () => {
      const obj = {
        name: mm.str('test'),
        age: mm.i(25n),
        active: mm.bool(true),
      };
      const encoded = encodeNode(mm.obj(obj));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Object);
      expect(decoded.value.name).toBe('test');
      expect(decoded.value.age.toString()).toBe('25');
      expect(decoded.value.active).toBe(true);
    });

    test('should encode and decode nested object', () => {
      const obj = { person: mm.obj({ name: mm.str('test'), age: mm.i(25n) }) };
      const encoded = encodeNode(mm.obj(obj));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Object);
      expect(decoded.value.person.name).toBe('test');
      expect(decoded.value.person.age).toBe(25n);
    });
  });

  describe('DateTime', () => {
    test('should encode and decode datetime', () => {
      const date = new Date('2024-01-01T12:00:00Z');
      const encoded = encodeNode(mm.datetime(date));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.DateTime);
    });

    test('should encode and decode date', () => {
      const date = new Date('2024-01-01T00:00:00Z');
      const encoded = encodeNode(mm.date(date));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Date);
    });

    test('should encode and decode time', () => {
      const date = new Date('2024-01-01T12:30:45Z');
      const encoded = encodeNode(mm.time(date));
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Time);
    });
  });

  describe('MM Value with Tag', () => {
    test('should encode and decode mm() wrapped value with custom tag', () => {
      const tag = new Tag();
      tag.type = ValueType.Int;
      tag.desc = 'test description';
      tag.min = '1';
      tag.max = '100';
      const wrapped = mm.i(42n, tag);
      const encoded = encodeNode(wrapped);
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(42n);
    });

    test('should encode and decode mm.i() with custom tag', () => {
      const tag = new Tag();
      tag.desc = 'custom integer';
      const wrapped = mm.i(123n, tag);
      const encoded = encodeNode(wrapped);
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.Int);
      expect(decoded.value).toBe(123n);
    });

    test('should encode and decode mm.str() with custom tag', () => {
      const tag = new Tag();
      tag.pattern = '^[a-z]+$';
      const wrapped = mm.str('hello', tag);
      const encoded = encodeNode(wrapped);
      const decoded = decodeToValue(encoded);
      expect(decoded.type).toBe(ValueType.String);
      expect(decoded.value).toBe('hello');
    });
  });
});
