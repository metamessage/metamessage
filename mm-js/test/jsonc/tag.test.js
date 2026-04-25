import { parseJSONC } from '../../src/jsonc/scanner.js';

describe('JSONCTag', () => {
  describe('parseMMTag', () => {
    test('parses empty tag', () => {
      const result = parseJSONC('{}');
      expect(result).not.toBeNull();
    });

    test('parses is_null tag', () => {
      const result = parseJSONC(`{
        // mm:is_null
        "value": null
      }`);
      expect(result).not.toBeNull();
    });

    test('parses type tag', () => {
      const result = parseJSONC(`{
        // mm:type=str
        "name": "test"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses desc tag', () => {
      const result = parseJSONC(`{
        // mm:desc="test description"
        "name": "test"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses nullable tag', () => {
      const result = parseJSONC(`{
        // mm:nullable
        "value": null
      }`);
      expect(result).not.toBeNull();
    });

    test('parses default tag', () => {
      const result = parseJSONC(`{
        // mm:default=value
        "value": "test"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses min max tags', () => {
      const result = parseJSONC(`{
        // mm:min=1;max=100
        "value": 50
      }`);
      expect(result).not.toBeNull();
    });

    test('parses size tag', () => {
      const result = parseJSONC(`{
        // mm:size=10
        "value": "test"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses enum tag', () => {
      const result = parseJSONC(`{
        // mm:enum=a|b|c
        "value": "a"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses pattern tag', () => {
      const result = parseJSONC(`{
        // mm:pattern=^[a-z]+$
        "value": "abc"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses location tag', () => {
      const result = parseJSONC(`{
        // mm:location=8
        "timezone": "+8"
      }`);
      expect(result).not.toBeNull();
    });

    test('parses version tag', () => {
      const result = parseJSONC(`{
        // mm:version=4
        "config": {}
      }`);
      expect(result).not.toBeNull();
    });

    test('parses child_type tag', () => {
      const result = parseJSONC(`{
        // mm:child_type=str
        "items": []
      }`);
      expect(result).not.toBeNull();
    });

    test('parses complex tag', () => {
      const result = parseJSONC(`{
        // mm:type=array;size=5;child_type=i;child_nullable;desc="array of ints"
        "numbers": [1, 2, 3]
      }`);
      expect(result).not.toBeNull();
    });

    test('parses multiple tags on same field', () => {
      const result = parseJSONC(`{
        // mm:type=str
        // mm:desc="name field"
        // mm:nullable
        "name": "test"
      }`);
      expect(result).not.toBeNull();
    });
  });
});