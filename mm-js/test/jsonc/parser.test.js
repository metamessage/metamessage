import { parseJSONC } from '../../src/jsonc/scanner.js';

describe('JSONCParser', () => {
  describe('parse empty structures', () => {
    test('parses empty object', () => {
      const result = parseJSONC('{}');
      expect(result).not.toBeNull();
      expect(result.getType()).toBe('object');
    });

    test('parses empty array', () => {
      const result = parseJSONC('[]');
      expect(result).not.toBeNull();
      expect(result.getType()).toBe('array');
    });
  });

  describe('parse simple values', () => {
    test('parses string', () => {
      const result = parseJSONC('"hello"');
      expect(result).not.toBeNull();
      expect(result.data).toBe('hello');
    });

    test('parses number', () => {
      const result = parseJSONC('123');
      expect(result).not.toBeNull();
      expect(result.data).toBe(123n);
    });

    test('parses float', () => {
      const result = parseJSONC('3.14');
      expect(result).not.toBeNull();
      expect(result.data).toBeCloseTo(3.14, 5);
    });

    test('parses true', () => {
      const result = parseJSONC('true');
      expect(result).not.toBeNull();
      expect(result.data).toBe(true);
    });

    test('parses false', () => {
      const result = parseJSONC('false');
      expect(result).not.toBeNull();
      expect(result.data).toBe(false);
    });

    test('parses null', () => {
      const result = parseJSONC('null');
      expect(result).not.toBeNull();
      expect(result.data).toBeNull();
    });
  });

  describe('parse object', () => {
    test('parses simple object', () => {
      const result = parseJSONC('{"name": "test"}');
      expect(result.getType()).toBe('object');
      expect(result.fields.length).toBe(1);
      expect(result.fields[0].key).toBe('name');
    });

    test('parses object with multiple fields', () => {
      const result = parseJSONC('{"name": "test", "value": 123}');
      expect(result.getType()).toBe('object');
      expect(result.fields.length).toBe(2);
    });

    test('parses nested object', () => {
      const result = parseJSONC('{"outer": {"inner": "value"}}');
      expect(result.getType()).toBe('object');
      expect(result.fields[0].value.getType()).toBe('object');
    });
  });

  describe('parse array', () => {
    test('parses simple array', () => {
      const result = parseJSONC('[1, 2, 3]');
      expect(result.getType()).toBe('array');
      expect(result.items.length).toBe(3);
    });

    test('parses empty array', () => {
      const result = parseJSONC('[]');
      expect(result.getType()).toBe('array');
      expect(result.items.length).toBe(0);
    });

    test('parses nested array', () => {
      const result = parseJSONC('[[1, 2], [3, 4]]');
      expect(result.getType()).toBe('array');
      expect(result.items[0].getType()).toBe('array');
    });
  });

  describe('parse with comments', () => {
    test('parses object with line comment', () => {
      const result = parseJSONC(`{
        // mm:type=str
        "name": "test"
      }`);
      expect(result).not.toBeNull();
      expect(result.getType()).toBe('object');
    });

    test('parses object with block comment', () => {
      const result = parseJSONC(`{
        /* mm:type=str */
        "name": "test"
      }`);
      expect(result).not.toBeNull();
      expect(result.getType()).toBe('object');
    });
  });

  describe('parse complex structures', () => {
    test('parses complex JSON', () => {
      const result = parseJSONC(`{
        "id": 1,
        "name": "test",
        "active": true,
        "items": [1, 2, 3],
        "metadata": {
          "created": "2024-01-01",
          "tags": ["a", "b"]
        }
      }`);
      expect(result.getType()).toBe('object');
      expect(result.fields.length).toBe(5);
    });
  });
});