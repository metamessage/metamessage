import { JSONCScanner, TokenType } from '../../src/jsonc/scanner.js';

describe('JSONCScanner', () => {
  describe('basic tokens', () => {
    test('scans empty input', () => {
      const scanner = new JSONCScanner('');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.EOF);
    });

    test('scans LCURLY', () => {
      const scanner = new JSONCScanner('{');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.LCURLY);
      expect(token.literal).toBe('{');
    });

    test('scans RCURLY', () => {
      const scanner = new JSONCScanner('}');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.RCURLY);
      expect(token.literal).toBe('}');
    });

    test('scans LBRACKET', () => {
      const scanner = new JSONCScanner('[');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.LBRACKET);
      expect(token.literal).toBe('[');
    });

    test('scans RBRACKET', () => {
      const scanner = new JSONCScanner(']');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.RBRACKET);
      expect(token.literal).toBe(']');
    });

    test('scans COLON', () => {
      const scanner = new JSONCScanner(':');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.COLON);
      expect(token.literal).toBe(':');
    });

    test('scans COMMA', () => {
      const scanner = new JSONCScanner(',');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.COMMA);
      expect(token.literal).toBe(',');
    });
  });

  describe('string tokens', () => {
    test('scans simple string', () => {
      const scanner = new JSONCScanner('"hello"');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.STRING);
      expect(token.literal).toBe('hello');
    });

    test('scans empty string', () => {
      const scanner = new JSONCScanner('""');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.STRING);
      expect(token.literal).toBe('');
    });

    test('scans string with escape', () => {
      const scanner = new JSONCScanner('"hello\\"world"');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.STRING);
      expect(token.literal).toBe('hello\\"world');
    });
  });

  describe('number tokens', () => {
    test('scans positive integer', () => {
      const scanner = new JSONCScanner('123');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.NUMBER);
      expect(token.literal).toBe('123');
    });

    test('scans negative integer', () => {
      const scanner = new JSONCScanner('-456');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.NUMBER);
      expect(token.literal).toBe('-456');
    });

    test('scans float', () => {
      const scanner = new JSONCScanner('3.14');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.NUMBER);
      expect(token.literal).toBe('3.14');
    });
  });

  describe('literal tokens', () => {
    test('scans true', () => {
      const scanner = new JSONCScanner('true');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.TRUE);
      expect(token.literal).toBe('true');
    });

    test('scans false', () => {
      const scanner = new JSONCScanner('false');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.FALSE);
      expect(token.literal).toBe('false');
    });

    test('scans null', () => {
      const scanner = new JSONCScanner('null');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.NULL);
      expect(token.literal).toBe('null');
    });
  });

  describe('comment tokens', () => {
    test('scans line comment', () => {
      const scanner = new JSONCScanner('// this is a comment');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.LEADING_COMMENT);
      expect(token.literal).toBe('this is a comment');
    });

    test('scans block comment', () => {
      const scanner = new JSONCScanner('/* this is a block comment */');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.LEADING_COMMENT);
      expect(token.literal).toBe('this is a block comment');
    });

    test('scans multi-line block comment', () => {
      const scanner = new JSONCScanner('/* this is\na block\ncomment */');
      const token = scanner.nextToken();
      expect(token.type).toBe(TokenType.LEADING_COMMENT);
      expect(token.literal).toBe('this is\na block\ncomment');
    });
  });

  describe('whitespace handling', () => {
    test('skips whitespace', () => {
      const scanner = new JSONCScanner('   \t\n\r{}');
      const t1 = scanner.nextToken();
      const t2 = scanner.nextToken();
      expect(t1.type).toBe(TokenType.LCURLY);
      expect(t2.type).toBe(TokenType.RCURLY);
    });
  });

  describe('line and column tracking', () => {
    test('tracks line numbers', () => {
      const scanner = new JSONCScanner('"a"\n"b"\n"c"');
      const t1 = scanner.nextToken();
      const t2 = scanner.nextToken();
      const t3 = scanner.nextToken();
      expect(t1.line).toBe(1);
      expect(t2.line).toBe(2);
      expect(t3.line).toBe(3);
    });
  });

  describe('complex input', () => {
    test('scans JSON object', () => {
      const scanner = new JSONCScanner('{"name": "test", "value": 123}');
      const tokens = [];
      while (true) {
        const token = scanner.nextToken();
        tokens.push(token);
        if (token.type === TokenType.EOF) break;
      }
      expect(tokens.length).toBeGreaterThan(0);
      expect(tokens[0].type).toBe(TokenType.LCURLY);
      expect(tokens[tokens.length - 1].type).toBe(TokenType.EOF);
    });

    test('scans JSON array', () => {
      const scanner = new JSONCScanner('[1, 2, 3]');
      const tokens = [];
      while (true) {
        const token = scanner.nextToken();
        tokens.push(token);
        if (token.type === TokenType.EOF) break;
      }
      expect(tokens[0].type).toBe(TokenType.LBRACKET);
    });
  });
});