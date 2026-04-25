use std::collections::HashMap;

#[derive(Debug, Clone, PartialEq)]
pub enum TokenType {
    EOF,
    LBrace,
    RBrace,
    LBracket,
    RBracket,
    Colon,
    Comma,
    String,
    Number,
    True,
    False,
    Null,
    LeadingComment,
    TrailingComment,
}

#[derive(Debug, Clone)]
pub struct Token {
    pub token_type: TokenType,
    pub literal: String,
    pub line: usize,
    pub column: usize,
}

pub struct Scanner {
    input: Vec<char>,
    position: usize,
    line: usize,
    column: usize,
    last_token: Option<Token>,
    current_token: Option<Token>,
}

impl Scanner {
    pub fn new(input: &str) -> Self {
        Self {
            input: input.chars().collect(),
            position: 0,
            line: 1,
            column: 1,
            last_token: None,
            current_token: None,
        }
    }

    pub fn next_token(&mut self) -> Token {
        self.skip_whitespace();

        if self.position >= self.input.len() {
            let token = Token {
                token_type: TokenType::EOF,
                literal: String::new(),
                line: self.line,
                column: self.column,
            };
            self.last_token = self.current_token.take();
            self.current_token = Some(token.clone());
            return token;
        }

        let ch = self.input[self.position];

        if ch == '/' {
            return self.scan_comment();
        }

        let token = match ch {
            '{' => self.create_token(TokenType::LBrace),
            '}' => self.create_token(TokenType::RBrace),
            '[' => self.create_token(TokenType::LBracket),
            ']' => self.create_token(TokenType::RBracket),
            ':' => self.create_token(TokenType::Colon),
            ',' => self.create_token(TokenType::Comma),
            '"' => self.scan_string(),
            _ if ch.is_ascii_digit() || ch == '-' => self.scan_number(),
            _ if ch.is_alphabetic() => self.scan_identifier(),
            _ => panic!("unexpected character: {} at line {}, column {}", ch, self.line, self.column),
        };

        self.last_token = self.current_token.take();
        self.current_token = Some(token.clone());
        token
    }

    fn create_token(&mut self, token_type: TokenType) -> Token {
        let token = Token {
            token_type,
            literal: String::new(),
            line: self.line,
            column: self.column,
        };
        self.advance(1);
        token
    }

    fn scan_comment(&mut self) -> Token {
        if self.position + 1 >= self.input.len() {
            return self.create_token(TokenType::LeadingComment);
        }

        let next = self.input[self.position + 1];
        let is_leading = self.last_token.is_none()
            || matches!(
                self.last_token.as_ref().map(|t| &t.token_type),
                Some(TokenType::Comma) | Some(TokenType::Colon) | Some(TokenType::LBrace) | Some(TokenType::LBracket)
            );

        let start_pos = self.position;
        let start_line = self.line;
        let start_column = self.column;

        if next == '/' {
            self.advance(2);
            while self.position < self.input.len() && self.input[self.position] != '\n' {
                self.advance(1);
            }
            let literal: String = self.input[start_pos..self.position].iter().collect();
            let token_type = if is_leading {
                TokenType::LeadingComment
            } else {
                TokenType::TrailingComment
            };
            let token = Token {
                token_type,
                literal,
                line: start_line,
                column: start_column,
            };
            self.last_token = self.current_token.take();
            self.current_token = Some(token.clone());
            return token;
        } else if next == '*' {
            self.advance(2);
            while self.position + 1 < self.input.len() {
                if self.input[self.position] == '*' && self.input[self.position + 1] == '/' {
                    self.advance(2);
                    break;
                }
                if self.input[self.position] == '\n' {
                    self.line += 1;
                    self.column = 0;
                }
                self.advance(1);
            }
            let literal: String = self.input[start_pos..self.position].iter().collect();
            let token_type = if is_leading {
                TokenType::LeadingComment
            } else {
                TokenType::TrailingComment
            };
            let token = Token {
                token_type,
                literal,
                line: start_line,
                column: start_column,
            };
            self.last_token = self.current_token.take();
            self.current_token = Some(token.clone());
            return token;
        }

        self.create_token(TokenType::LeadingComment)
    }

    fn scan_string(&mut self) -> Token {
        let start_line = self.line;
        let start_column = self.column;
        self.advance(1);

        let mut sb = String::new();
        while self.position < self.input.len() && self.input[self.position] != '"' {
            if self.input[self.position] == '\\' && self.position + 1 < self.input.len() {
                self.advance(1);
                let escaped = self.input[self.position];
                match escaped {
                    'n' => sb.push('\n'),
                    'r' => sb.push('\r'),
                    't' => sb.push('\t'),
                    'b' => sb.push('\x08'),
                    'f' => sb.push('\x0c'),
                    '"' => sb.push('"'),
                    '\\' => sb.push('\\'),
                    'u' => {
                        if self.position + 4 < self.input.len() {
                            self.advance(1);
                            let hex: String = self.input[self.position..self.position + 4].iter().collect();
                            if let Ok(unicode) = u32::from_str_radix(&hex, 16) {
                                sb.push(char::from_u32(unicode).unwrap_or('\u{FFFD}'));
                                self.advance(3);
                            }
                        }
                    }
                    _ => sb.push(escaped),
                }
            } else {
                if self.input[self.position] == '\n' {
                    self.line += 1;
                    self.column = 0;
                }
                sb.push(self.input[self.position]);
            }
            self.advance(1);
        }

        if self.position < self.input.len() {
            self.advance(1);
        }

        Token {
            token_type: TokenType::String,
            literal: sb,
            line: start_line,
            column: start_column,
        }
    }

    fn scan_number(&mut self) -> Token {
        let start_line = self.line;
        let start_column = self.column;
        let mut sb = String::new();

        if self.input[self.position] == '-' {
            sb.push('-');
            self.advance(1);
        }

        while self.position < self.input.len() {
            let ch = self.input[self.position];
            if ch.is_ascii_digit() || ch == '.' || ch == 'e' || ch == 'E' || ch == '+' || ch == '_' {
                if ch == '_' {
                    self.advance(1);
                    continue;
                }
                sb.push(ch);
                self.advance(1);
            } else {
                break;
            }
        }

        if self.position < self.input.len() {
            let ch = self.input[self.position];
            if ch == 'f' || ch == 'F' {
                sb.push(ch);
                self.advance(1);
            }
        }

        Token {
            token_type: TokenType::Number,
            literal: sb,
            line: start_line,
            column: start_column,
        }
    }

    fn scan_identifier(&mut self) -> Token {
        let start_line = self.line;
        let start_column = self.column;
        let mut sb = String::new();

        while self.position < self.input.len() && (self.input[self.position].is_alphanumeric() || self.input[self.position] == '_') {
            sb.push(self.input[self.position]);
            self.advance(1);
        }

        let identifier = sb.to_lowercase();
        let (token_type, literal) = match identifier.as_str() {
            "true" => (TokenType::True, "true".to_string()),
            "false" => (TokenType::False, "false".to_string()),
            "null" => (TokenType::Null, "null".to_string()),
            _ => (TokenType::String, sb.clone()),
        };

        Token {
            token_type,
            literal,
            line: start_line,
            column: start_column,
        }
    }

    fn skip_whitespace(&mut self) {
        while self.position < self.input.len() {
            let ch = self.input[self.position];
            if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
                if ch == '\n' {
                    self.line += 1;
                    self.column = 0;
                }
                self.advance(1);
            } else {
                break;
            }
        }
    }

    fn advance(&mut self, count: usize) {
        for _ in 0..count {
            if self.position < self.input.len() {
                if self.input[self.position] == '\n' {
                    self.line += 1;
                    self.column = 0;
                } else {
                    self.column += 1;
                }
                self.position += 1;
            }
        }
    }
}