import { parseJSONC, printJSONC, printJSONCCompact, parseMMTag, JSONCValueType } from '../../src/jsonc/index';

describe('JSONC Parser', () => {
  test('should parse empty object', () => {
    const input = '{}';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse object with properties', () => {
    const input = '{ "name": "test", "age": 25 }';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse array', () => {
    const input = '[1, 2, 3]';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('array');
  });

  test('should parse nested structure', () => {
    const input = '{ "person": { "name": "test", "age": 25, "hobbies": ["reading", "gaming"] } }';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse with line comments', () => {
    const input = '{ // This is a comment\n  "name": "test" // Another comment\n}';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse with block comments', () => {
    const input = '{ /* This is a block comment */ "name": "test" }';
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse mm: tags from line comments', () => {
    const input = `{ // mm: type=str;desc=name field
  "name": "test"
}`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse mm: tags from block comments', () => {
    const input = `{ /* mm: type=str;desc=name field */ "name": "test" }`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should ignore non-mm: comments', () => {
    const input = `{ // regular comment
  "name": "test"
}`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse mm: type for UUID', () => {
    const input = `{ /* mm: type=uuid;desc=user id */ "id": "550e8400-e29b-41d4-a716-446655440000" }`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse mm: type for DateTime', () => {
    const input = `{ /* mm: type=datetime;desc=creation time */ "created_at": "2024-01-01T00:00:00Z" }`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });

  test('should parse mm: type for Email', () => {
    const input = `{ /* mm: type=email;desc=user email */ "email": "test@example.com" }`;
    const doc = parseJSONC(input);
    const root = doc.getRoot();
    expect(root.getType()).toBe('object');
  });
});

describe('JSONC Printer', () => {
  test('should print object', () => {
    const input = '{ "name": "test", "age": 25 }';
    const doc = parseJSONC(input);
    const printed = printJSONC(doc);
    expect(printed).toContain('name');
    expect(printed).toContain('test');
    expect(printed).toContain('age');
    expect(printed).toContain('25');
  });

  test('should print compact', () => {
    const input = '{ "name": "test", "age": 25 }';
    const doc = parseJSONC(input);
    const printed = printJSONCCompact(doc);
    expect(printed).toBe('{"name":"test","age":25}');
  });

  test('should print desc tag as comment', () => {
    const input = `{ /* mm: type=str;desc=name field */ "name": "test" }`;
    const doc = parseJSONC(input);
    const printed = printJSONC(doc);
    expect(printed).toContain('// mm:');
    expect(printed).toContain('desc=name field');
  });

  test('should quote UUID type values', () => {
    const input = `{ /* mm: type=uuid;desc=id */ "id": "550e8400-e29b-41d4-a716-446655440000" }`;
    const doc = parseJSONC(input);
    const printed = printJSONC(doc);
    expect(printed).toContain('"550e8400-e29b-41d4-a716-446655440000"');
  });

  test('should quote Email type values', () => {
    const input = `{ /* mm: type=email;desc=email */ "email": "test@example.com" }`;
    const doc = parseJSONC(input);
    const printed = printJSONC(doc);
    expect(printed).toContain('"test@example.com"');
  });

  test('should not quote int values', () => {
    const input = `{ "count": 42 }`;
    const doc = parseJSONC(input);
    const printed = printJSONCCompact(doc);
    expect(printed).toBe('{"count":42}');
  });

  test('should not quote float values', () => {
    const input = `{ "price": 3.14 }`;
    const doc = parseJSONC(input);
    const printed = printJSONCCompact(doc);
    expect(printed).toBe('{"price":3.14}');
  });
});

describe('JSONC Tag Parser', () => {
  test('should parse simple tag', () => {
    const tagStr = 'type=str';
    const tag = parseMMTag(tagStr);
    expect(tag.type).toBe(JSONCValueType.String);
  });

  test('should parse multiple tags', () => {
    const tagStr = 'type=str;desc=test;nullable';
    const tag = parseMMTag(tagStr);
    expect(tag.type).toBe(JSONCValueType.String);
    expect(tag.desc).toBe('test');
    expect(tag.nullable).toBe(true);
  });

  test('should parse is_null tag', () => {
    const tagStr = 'is_null';
    const tag = parseMMTag(tagStr);
    expect(tag.isNull).toBe(true);
    expect(tag.nullable).toBe(true);
  });

  test('should parse empty tag', () => {
    const tagStr = '';
    const tag = parseMMTag(tagStr);
    expect(tag.type).toBe(JSONCValueType.Unknown);
  });

  test('should parse all type abbreviations', () => {
    const types: [string, JSONCValueType][] = [
      ['i', JSONCValueType.Int],
      ['i8', JSONCValueType.Int8],
      ['i16', JSONCValueType.Int16],
      ['i32', JSONCValueType.Int32],
      ['i64', JSONCValueType.Int64],
      ['u', JSONCValueType.Uint],
      ['u8', JSONCValueType.Uint8],
      ['u16', JSONCValueType.Uint16],
      ['u32', JSONCValueType.Uint32],
      ['u64', JSONCValueType.Uint64],
      ['f32', JSONCValueType.Float32],
      ['f64', JSONCValueType.Float64],
      ['bool', JSONCValueType.Bool],
      ['bytes', JSONCValueType.Bytes],
      ['bi', JSONCValueType.BigInt],
      ['datetime', JSONCValueType.DateTime],
      ['date', JSONCValueType.Date],
      ['time', JSONCValueType.Time],
      ['uuid', JSONCValueType.UUID],
      ['decimal', JSONCValueType.Decimal],
      ['ip', JSONCValueType.IP],
      ['url', JSONCValueType.URL],
      ['email', JSONCValueType.Email],
      ['enum', JSONCValueType.Enum],
      ['arr', JSONCValueType.Array],
      ['struct', JSONCValueType.Struct],
    ];

    for (const [abbr, expectedType] of types) {
      const tag = parseMMTag(`type=${abbr}`);
      expect(tag.type).toBe(expectedType);
    }
  });

  test('tag toString should format correctly', () => {
    const tag = parseMMTag('type=str;desc=test;nullable');
    const str = tag.toString();
    expect(str).toContain('type=str');
    expect(str).toContain('desc=test');
    expect(str).toContain('nullable');
  });
});