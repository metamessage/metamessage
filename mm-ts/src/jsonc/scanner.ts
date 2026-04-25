export enum TokenType {
  LCURLY = '{',
  RCURLY = '}',
  LBRACKET = '[',
  RBRACKET = ']',
  COLON = ':',
  COMMA = ',',
  STRING = 'STRING',
  NUMBER = 'NUMBER',
  TRUE = 'true',
  FALSE = 'false',
  NULL = 'null',
  LINECOMMENT = 'LINECOMMENT',
  BLOCKCOMMENT = 'BLOCKCOMMENT',
  EOF = 'EOF',
}

export interface Token {
  type: TokenType;
  value: string;
  line: number;
  column: number;
}

export class JSONCScanner {
  private input: string;
  private position: number;
  private line: number;
  private column: number;
  private length: number;

  constructor(input: string) {
    this.input = input;
    this.position = 0;
    this.line = 1;
    this.column = 1;
    this.length = input.length;
  }

  nextToken(): Token {
    this.skipWhitespace();
    
    if (this.position >= this.length) {
      return this.createToken(TokenType.EOF, '');
    }

    const char = this.input[this.position];

    switch (char) {
      case '{':
        return this.consumeChar(TokenType.LCURLY);
      case '}':
        return this.consumeChar(TokenType.RCURLY);
      case '[':
        return this.consumeChar(TokenType.LBRACKET);
      case ']':
        return this.consumeChar(TokenType.RBRACKET);
      case ':':
        return this.consumeChar(TokenType.COLON);
      case ',':
        return this.consumeChar(TokenType.COMMA);
      case '"':
        return this.consumeString();
      case '0':
      case '1':
      case '2':
      case '3':
      case '4':
      case '5':
      case '6':
      case '7':
      case '8':
      case '9':
      case '-':
        return this.consumeNumber();
      case 't':
        return this.consumeKeyword('true', TokenType.TRUE);
      case 'f':
        return this.consumeKeyword('false', TokenType.FALSE);
      case 'n':
        return this.consumeKeyword('null', TokenType.NULL);
      case '/':
        return this.consumeComment();
      default:
        throw new Error(`Unexpected character: ${char} at ${this.line}:${this.column}`);
    }
  }

  private consumeChar(type: TokenType): Token {
    const value = this.input[this.position] as string;
    const token = this.createToken(type, value);
    this.position++;
    this.column++;
    return token;
  }

  private consumeString(): Token {
    const start = this.position;
    const startLine = this.line;
    const startColumn = this.column;
    this.position++;
    this.column++;

    while (this.position < this.length) {
      const char = this.input[this.position];
      
      if (char === '"') {
        this.position++;
        this.column++;
        const value = this.input.substring(start + 1, this.position - 1);
        return {
          type: TokenType.STRING,
          value,
          line: startLine,
          column: startColumn,
        };
      } else if (char === '\\') {
        this.position++;
        this.column++;
        if (this.position < this.length) {
          this.position++;
          this.column++;
        }
      } else if (char === '\n') {
        this.line++;
        this.column = 1;
        this.position++;
      } else {
        this.position++;
        this.column++;
      }
    }

    throw new Error(`Unclosed string at ${startLine}:${startColumn}`);
  }

  private consumeNumber(): Token {
    const start = this.position;
    const startLine = this.line;
    const startColumn = this.column;

    if (this.input[this.position] === '-') {
      this.position++;
      this.column++;
    }

    if (this.position < this.length && this.input[this.position] === '0') {
      this.position++;
      this.column++;
    } else if (this.position < this.length && /[1-9]/.test(this.input[this.position] as string)) {
      this.position++;
      this.column++;
      while (this.position < this.length && /[0-9]/.test(this.input[this.position] as string)) {
        this.position++;
        this.column++;
      }
    } else {
      throw new Error(`Invalid number at ${startLine}:${startColumn}`);
    }

    if (this.position < this.length && this.input[this.position] === '.') {
      this.position++;
      this.column++;
      if (this.position < this.length && /[0-9]/.test(this.input[this.position] as string)) {
        while (this.position < this.length && /[0-9]/.test(this.input[this.position] as string)) {
          this.position++;
          this.column++;
        }
      } else {
        throw new Error(`Invalid number at ${startLine}:${startColumn}`);
      }
    }

    if (this.position < this.length && (this.input[this.position] === 'e' || this.input[this.position] === 'E')) {
      this.position++;
      this.column++;
      if (this.position < this.length && (this.input[this.position] === '+' || this.input[this.position] === '-')) {
        this.position++;
        this.column++;
      }
      if (this.position < this.length && /[0-9]/.test(this.input[this.position] as string)) {
        while (this.position < this.length && /[0-9]/.test(this.input[this.position] as string)) {
          this.position++;
          this.column++;
        }
      } else {
        throw new Error(`Invalid number at ${startLine}:${startColumn}`);
      }
    }

    const value = this.input.substring(start, this.position);
    return {
      type: TokenType.NUMBER,
      value,
      line: startLine,
      column: startColumn,
    };
  }

  private consumeKeyword(expected: string, type: TokenType): Token {
    const start = this.position;
    const startLine = this.line;
    const startColumn = this.column;

    for (let i = 0; i < expected.length; i++) {
      if (this.position + i >= this.length || this.input[this.position + i] !== expected[i]) {
        throw new Error(`Unexpected token at ${startLine}:${startColumn}`);
      }
    }

    this.position += expected.length;
    this.column += expected.length;

    return {
      type,
      value: expected,
      line: startLine,
      column: startColumn,
    };
  }

  private consumeComment(): Token {
    const start = this.position;
    const startLine = this.line;
    const startColumn = this.column;
    this.position++;
    this.column++;

    if (this.position < this.length && this.input[this.position] === '/') {
      // Line comment
      this.position++;
      this.column++;
      while (this.position < this.length && this.input[this.position] !== '\n') {
        this.position++;
        this.column++;
      }
      const value = this.input.substring(start, this.position);
      return {
        type: TokenType.LINECOMMENT,
        value,
        line: startLine,
        column: startColumn,
      };
    } else if (this.position < this.length && this.input[this.position] === '*') {
      // Block comment
      this.position++;
      this.column++;
      let foundEnd = false;
      while (this.position + 1 < this.length) {
        if (this.input[this.position] === '*' && this.input[this.position + 1] === '/') {
          this.position += 2;
          this.column += 2;
          foundEnd = true;
          break;
        } else if (this.input[this.position] === '\n') {
          this.line++;
          this.column = 1;
          this.position++;
        } else {
          this.position++;
          this.column++;
        }
      }
      if (!foundEnd) {
        throw new Error(`Unclosed block comment at ${startLine}:${startColumn}`);
      }
      const value = this.input.substring(start, this.position);
      return {
        type: TokenType.BLOCKCOMMENT,
        value,
        line: startLine,
        column: startColumn,
      };
    } else {
      throw new Error(`Unexpected character after '/' at ${startLine}:${startColumn}`);
    }
  }

  private skipWhitespace(): void {
    while (this.position < this.length) {
      const char = this.input[this.position];
      if (char === ' ' || char === '\t' || char === '\r') {
        this.position++;
        this.column++;
      } else if (char === '\n') {
        this.line++;
        this.column = 1;
        this.position++;
      } else {
        break;
      }
    }
  }

  private createToken(type: TokenType, value: string): Token {
    return {
      type,
      value,
      line: this.line,
      column: this.column,
    };
  }
}

export function jsoncScanner(input: string): JSONCScanner {
  return new JSONCScanner(input);
}