import { parseJSONC, printJSONC, printJSONCCompact, parseMMTag } from '../../src/jsonc/index';

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
});

describe('JSONC Tag Parser', () => {
  test('should parse simple tag', () => {
    const tagStr = 'type=str';
    const tag = parseMMTag(tagStr);
    expect(tag.get('type')).toBe('str');
  });

  test('should parse multiple tags', () => {
    const tagStr = 'type=str;desc=test;nullable=true';
    const tag = parseMMTag(tagStr);
    expect(tag.get('type')).toBe('str');
    expect(tag.get('desc')).toBe('test');
    expect(tag.get('nullable')).toBe('true');
  });

  test('should parse empty tag', () => {
    const tagStr = '';
    const tag = parseMMTag(tagStr);
    expect(tag.size()).toBe(0);
  });
});