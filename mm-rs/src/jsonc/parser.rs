use crate::ir::ast::{Array, Field, Node, Object, Value, ValueData};
use crate::ir::tag::Tag;
use crate::ir::value_type::ValueType;
use crate::jsonc::scanner::{Scanner, Token, TokenType};

const MAX_DEPTH: usize = 32;

pub struct Parser {
    toks: Vec<Token>,
    pos: usize,
    pending: Vec<Token>,
    depth: usize,
}

impl Parser {
    pub fn new(input: &str) -> Self {
        let mut scanner = Scanner::new(input);
        let mut toks = Vec::new();
        loop {
            let tok = scanner.next_token();
            let is_eof = tok.token_type == TokenType::EOF;
            toks.push(tok);
            if is_eof {
                break;
            }
        }
        Self {
            toks,
            pos: 0,
            pending: Vec::new(),
            depth: 0,
        }
    }

    fn peek(&self) -> &Token {
        if self.pos >= self.toks.len() {
            static EOF: Token = Token {
                token_type: TokenType::EOF,
                literal: String::new(),
                line: 0,
                column: 0,
            };
            return &EOF;
        }
        &self.toks[self.pos]
    }

    fn next(&mut self) -> Token {
        let tok = self.peek().clone();
        self.pos += 1;
        tok
    }

    fn consume_comments_for(&mut self, anchor_line: usize) -> Option<Tag> {
        if self.pending.is_empty() {
            return None;
        }

        let last = &self.pending[self.pending.len() - 1];
        if anchor_line - last.line > 1 {
            self.pending.clear();
            return None;
        }

        let mut result: Option<Tag> = None;
        for comment in &self.pending {
            if let Some(parsed) = Tag::parse(&comment.literal) {
                result = Some(Tag::merge(result, parsed));
            }
        }

        self.pending.clear();
        result
    }

    pub fn parse(&mut self) -> Result<Node, String> {
        loop {
            let tok = self.peek().clone();
            if tok.token_type == TokenType::EOF {
                return Err("empty input".to_string());
            }

            if tok.token_type == TokenType::LeadingComment {
                if !self.pending.is_empty() {
                    let last = &self.pending[self.pending.len() - 1];
                    if tok.line - last.line > 1 {
                        self.pending.clear();
                    }
                }
                self.pending.push(tok);
                self.next();
                continue;
            }

            if tok.token_type == TokenType::TrailingComment {
                self.next();
                continue;
            }

            return self
                .parse_node("")
                .and_then(|opt| opt.ok_or_else(|| "no value parsed".to_string()));
        }
    }

    fn parse_node(&mut self, path: &str) -> Result<Option<Node>, String> {
        let tok = self.next();

        match tok.token_type {
            TokenType::EOF => Ok(None),
            TokenType::LBrace => self.parse_object(tok.line, path).map(Some),
            TokenType::LBracket => self.parse_array(tok.line, path).map(Some),
            TokenType::String => {
                let mut tag = self.consume_comments_for(tok.line).unwrap_or_default();
                if tag.value_type == ValueType::Unknown {
                    tag.value_type = ValueType::String;
                }
                let text = tok.literal;
                let value = Node::Value(Value {
                    data: ValueData::String(text.clone()),
                    text,
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::Number => {
                let mut tag = self.consume_comments_for(tok.line).unwrap_or_default();
                let text = tok.literal;

                if tag.value_type == ValueType::Unknown {
                    if text.contains('.') {
                        tag.value_type = ValueType::Float64;
                    } else {
                        tag.value_type = ValueType::Int;
                    }
                }

                let data: ValueData;
                if text.contains('.') {
                    data = ValueData::Float(text.parse::<f64>().unwrap_or(0.0));
                } else if text.starts_with('-') {
                    if let Ok(ival) = text.parse::<i64>() {
                        data = ValueData::Int(ival);
                    } else {
                        data = ValueData::Int(i64::MIN);
                    }
                } else if let Ok(uval) = text.parse::<u64>() {
                    if uval > i64::MAX as u64 {
                        data = ValueData::Uint(uval);
                    } else {
                        data = ValueData::Int(uval as i64);
                    }
                } else {
                    data = ValueData::Int(0);
                }

                let value = Node::Value(Value {
                    data,
                    text,
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::True => {
                let mut tag = self.consume_comments_for(tok.line).unwrap_or_default();
                if tag.value_type == ValueType::Unknown {
                    tag.value_type = ValueType::Bool;
                }
                let value = Node::Value(Value {
                    data: ValueData::Bool(true),
                    text: "true".to_string(),
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::False => {
                let mut tag = self.consume_comments_for(tok.line).unwrap_or_default();
                if tag.value_type == ValueType::Unknown {
                    tag.value_type = ValueType::Bool;
                }
                let value = Node::Value(Value {
                    data: ValueData::Bool(false),
                    text: "false".to_string(),
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::Null => {
                let mut tag = self.consume_comments_for(tok.line).unwrap_or_default();
                tag.nullable = true;
                tag.is_null = true;
                let value = Node::Value(Value {
                    data: ValueData::Null,
                    text: "null".to_string(),
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::TrailingComment => Ok(None),
            _ => Err(format!("unexpected token: {:?}", tok.token_type)),
        }
    }

    fn parse_object(&mut self, open_line: usize, path: &str) -> Result<Node, String> {
        self.depth += 1;
        if self.depth > MAX_DEPTH {
            return Err("max depth exceeded".to_string());
        }

        let mut tag = self.consume_comments_for(open_line).unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            tag.value_type = ValueType::Struct;
        }

        let obj_path = path.to_string();
        let mut fields: Vec<Field> = Vec::new();

        loop {
            let tok = self.peek().clone();
            if tok.token_type == TokenType::EOF {
                break;
            }
            if tok.token_type == TokenType::RBrace {
                self.next();
                break;
            }

            if tok.token_type == TokenType::LeadingComment {
                if !self.pending.is_empty() {
                    let last = &self.pending[self.pending.len() - 1];
                    if tok.line - last.line > 1 {
                        self.pending.clear();
                    }
                }
                self.pending.push(tok);
                self.next();
                continue;
            }

            if tok.token_type == TokenType::TrailingComment {
                if let Some(parsed) = Tag::parse(&tok.literal) {
                    if let Some(last) = fields.last_mut() {
                        if let Some(ref mut existing) = last.value.get_tag_mut() {
                            let merged = Tag::merge(Some(existing.clone()), parsed);
                            **existing = merged;
                        }
                    }
                }
                self.next();
                continue;
            }

            let key_tok = self.next();
            if key_tok.token_type != TokenType::String {
                return Err("expected string key".to_string());
            }
            let key = key_tok.literal;

            let colon = self.next();
            if colon.token_type != TokenType::Colon {
                return Err("expected colon".to_string());
            }

            let child_path = format!("{}.{}", obj_path, key);
            if let Some(mut val) = self.parse_node(&child_path)? {
                if let Some(mut ct) = val.get_tag().cloned() {
                    ct.inherit(&tag);
                    if let Some(t) = val.get_tag_mut() {
                        *t = ct;
                    }
                }
                fields.push(Field { key, value: val });
            }

            if self.peek().token_type == TokenType::Comma {
                self.next();
            }
        }

        self.depth -= 1;
        Ok(Node::Object(Object {
            fields,
            tag: Some(tag),
            path: obj_path,
        }))
    }

    fn parse_array(&mut self, open_line: usize, path: &str) -> Result<Node, String> {
        self.depth += 1;
        if self.depth > MAX_DEPTH {
            return Err("max depth exceeded".to_string());
        }

        let mut tag = self.consume_comments_for(open_line).unwrap_or_default();
        if tag.value_type == ValueType::Unknown {
            if let Some(size) = tag.size {
                if size > 0 {
                    tag.value_type = ValueType::Array;
                } else {
                    tag.value_type = ValueType::Slice;
                }
            } else {
                tag.value_type = ValueType::Slice;
            }
        }

        let arr_path = path.to_string();
        let mut items: Vec<Node> = Vec::new();
        let mut index: usize = 0;

        loop {
            let tok = self.peek().clone();
            if tok.token_type == TokenType::EOF {
                break;
            }
            if tok.token_type == TokenType::RBracket {
                self.next();
                break;
            }

            if tok.token_type == TokenType::LeadingComment {
                if !self.pending.is_empty() {
                    let last = &self.pending[self.pending.len() - 1];
                    if tok.line - last.line > 1 {
                        self.pending.clear();
                    }
                }
                self.pending.push(tok);
                self.next();
                continue;
            }

            if tok.token_type == TokenType::TrailingComment {
                if let Some(parsed) = Tag::parse(&tok.literal) {
                    if let Some(last) = items.last_mut() {
                        if let Some(ref mut existing) = last.get_tag_mut() {
                            let merged = Tag::merge(Some(existing.clone()), parsed);
                            **existing = merged;
                        }
                    }
                }
                self.next();
                continue;
            }

            let item_path = format!("{}[{}]", arr_path, index);
            if let Some(mut item) = self.parse_node(&item_path)? {
                if let Some(mut ct) = item.get_tag().cloned() {
                    ct.inherit(&tag);
                    if let Some(t) = item.get_tag_mut() {
                        *t = ct;
                    }
                }
                items.push(item);
                index += 1;
            }

            if self.peek().token_type == TokenType::Comma {
                self.next();
            }
        }

        self.depth -= 1;
        Ok(Node::Array(Array {
            items,
            tag: Some(tag),
            path: arr_path,
        }))
    }
}
