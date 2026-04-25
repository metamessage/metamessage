import { JSONCScanner, TokenType, Token } from './scanner';
import { JSONCValue, JSONCObject, JSONCArray, JSONCDoc, JSONCTag, parseMMTag } from './ast';

export interface JSONCNode {
  getType(): string;
  getTag(): any;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: any): void;
}

export class JSONCParser {
  private scanner: JSONCScanner;
  private currentToken: Token;

  constructor(input: string) {
    this.scanner = new JSONCScanner(input);
    this.currentToken = this.scanner.nextToken();
  }

  parse(): JSONCDoc {
    const root = this.parseValue();
    this.expect(TokenType.EOF);
    return new JSONCDoc(root);
  }

  private parseValue(): JSONCNode {
    this.skipComments();
    
    switch (this.currentToken.type) {
      case TokenType.LCURLY:
        return this.parseObject();
      case TokenType.LBRACKET:
        return this.parseArray();
      case TokenType.STRING:
        return this.parseString();
      case TokenType.NUMBER:
        return this.parseNumber();
      case TokenType.TRUE:
        return this.parseBoolean(true);
      case TokenType.FALSE:
        return this.parseBoolean(false);
      case TokenType.NULL:
        return this.parseNull();
      default:
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
    }
  }

  private parseObject(): JSONCObject {
    this.expect(TokenType.LCURLY);
    const obj = new JSONCObject();
    
    this.skipComments();
    
    while (this.currentToken.type !== (TokenType.RCURLY as any)) {
      this.skipComments();
      const keyToken = this.currentToken;
      this.expect(TokenType.STRING);
      const key = keyToken.value;
      
      this.skipComments();
      this.expect(TokenType.COLON);
      
      this.skipComments();
      const value = this.parseValue();
      
      // Check for comments after value that might contain tags
      this.skipComments();
      
      obj.setProperty(key, value);
      
      if (this.currentToken.type === (TokenType.COMMA as any)) {
        this.expect(TokenType.COMMA);
        this.skipComments();
      } else if (this.currentToken.type !== (TokenType.RCURLY as any)) {
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
      }
    }
    
    this.expect(TokenType.RCURLY);
    return obj;
  }

  private parseArray(): JSONCArray {
    this.expect(TokenType.LBRACKET);
    const array = new JSONCArray();
    
    this.skipComments();
    
    while (this.currentToken.type !== (TokenType.RBRACKET as any)) {
      this.skipComments();
      const value = this.parseValue();
      array.addElement(value);
      
      this.skipComments();
      
      if (this.currentToken.type === (TokenType.COMMA as any)) {
        this.expect(TokenType.COMMA);
        this.skipComments();
      } else if (this.currentToken.type !== (TokenType.RBRACKET as any)) {
        throw new Error(`Unexpected token: ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
      }
    }
    
    this.expect(TokenType.RBRACKET);
    return array;
  }

  private parseString(): JSONCValue {
    const token = this.currentToken;
    this.expect(TokenType.STRING);
    return new JSONCValue(token.value);
  }

  private parseNumber(): JSONCValue {
    const token = this.currentToken;
    this.expect(TokenType.NUMBER);
    const value = parseFloat(token.value);
    return new JSONCValue(value);
  }

  private parseBoolean(value: boolean): JSONCValue {
    this.expect(value ? TokenType.TRUE : TokenType.FALSE);
    return new JSONCValue(value);
  }

  private parseNull(): JSONCValue {
    this.expect(TokenType.NULL);
    return new JSONCValue(null);
  }

  private skipComments(): void {
    while (this.currentToken.type === TokenType.LINECOMMENT || this.currentToken.type === TokenType.BLOCKCOMMENT) {
      this.consumeToken();
    }
  }

  private expect(type: TokenType): void {
    if (this.currentToken.type !== type) {
      throw new Error(`Expected ${type}, got ${this.currentToken.type} at ${this.currentToken.line}:${this.currentToken.column}`);
    }
    this.consumeToken();
  }

  private consumeToken(): void {
    this.currentToken = this.scanner.nextToken();
  }
}

export function parseJSONC(input: string): JSONCDoc {
  const parser = new JSONCParser(input);
  return parser.parse();
}