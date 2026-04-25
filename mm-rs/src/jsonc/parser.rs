use crate::jsonc::scanner::{Scanner, Token, TokenType};
use crate::jsonc::ast::{Node, Object, Array, Value, Field, ValueData};
use crate::jsonc::tag::Tag;

pub struct Parser {
    scanner: Scanner,
    current_token: Option<Token>,
}

impl Parser {
    pub fn new(input: &str) -> Self {
        Self {
            scanner: Scanner::new(input),
            current_token: None,
        }
    }

    pub fn parse(&mut self) -> Result<Node, String> {
        self.current_token = Some(self.scanner.next_token());
        self.parse_value()
    }

    fn next_token(&mut self) {
        self.current_token = Some(self.scanner.next_token());
    }

    fn peek(&self) -> &Token {
        self.current_token.as_ref().expect("no current token")
    }

    fn parse_value(&mut self) -> Result<Node, String> {
        let token = self.peek();

        match &token.token_type {
            TokenType::LBrace => {
                self.next_token();
                self.parse_object()
            }
            TokenType::LBracket => {
                self.next_token();
                self.parse_array()
            }
            TokenType::String | TokenType::Number | TokenType::True | TokenType::False | TokenType::Null => {
                self.parse_primitive()
            }
            TokenType::EOF => Err("unexpected EOF".to_string()),
            _ => Err(format!("unexpected token: {:?}", token.token_type)),
        }
    }

    fn parse_primitive(&mut self) -> Result<Node, String> {
        let token = self.peek().clone();
        self.next_token();

        let mut tag = Tag::new();
        tag.value_type = crate::jsonc::value_type::ValueType::String;

        let data = match token.token_type {
            TokenType::String => ValueData::String(token.literal.clone()),
            TokenType::Number => {
                if let Ok(n) = token.literal.parse::<i64>() {
                    ValueData::Int(n)
                } else if let Ok(f) = token.literal.parse::<f64>() {
                    ValueData::Float(f)
                } else {
                    ValueData::String(token.literal.clone())
                }
            }
            TokenType::True => ValueData::Bool(true),
            TokenType::False => ValueData::Bool(false),
            TokenType::Null => ValueData::Null,
            _ => return Err(format!("unexpected token type: {:?}", token.token_type)),
        };

        Ok(Node::Value(Value {
            data,
            text: token.literal,
            tag: Some(tag),
        }))
    }

    fn parse_object(&mut self) -> Result<Node, String> {
        let mut fields = Vec::new();
        let tag = Tag::new();

        loop {
            let token = self.peek();

            match &token.token_type {
                TokenType::RBrace => {
                    self.next_token();
                    break;
                }
                TokenType::EOF => break,
                TokenType::Comma => {
                    self.next_token();
                }
                TokenType::LeadingComment => {
                    self.next_token();
                }
                TokenType::String => {
                    let key = token.literal.clone();
                    self.next_token();

                    if !matches!(self.peek().token_type, TokenType::Colon) {
                        return Err("expected colon after key".to_string());
                    }
                    self.next_token();

                    let value = self.parse_value()?;
                    fields.push(Field { key, value });
                }
                _ => {
                    return Err(format!("unexpected token in object: {:?}", token.token_type));
                }
            }
        }

        Ok(Node::Object(Object {
            fields,
            tag: Some(tag),
        }))
    }

    fn parse_array(&mut self) -> Result<Node, String> {
        let mut items = Vec::new();
        let tag = Tag::new();

        loop {
            let token = self.peek();

            match &token.token_type {
                TokenType::RBracket => {
                    self.next_token();
                    break;
                }
                TokenType::EOF => break,
                TokenType::Comma => {
                    self.next_token();
                }
                TokenType::LeadingComment => {
                    self.next_token();
                }
                TokenType::String | TokenType::Number | TokenType::True | TokenType::False | TokenType::Null => {
                    let value = self.parse_value()?;
                    items.push(value);
                }
                TokenType::LBrace => {
                    let value = self.parse_value()?;
                    items.push(value);
                }
                TokenType::LBracket => {
                    let value = self.parse_value()?;
                    items.push(value);
                }
                _ => {
                    return Err(format!("unexpected token in array: {:?}", token.token_type));
                }
            }
        }

        Ok(Node::Array(Array {
            items,
            tag: Some(tag),
        }))
    }
}