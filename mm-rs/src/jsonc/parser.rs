use crate::core::utils::camel_to_snake;
use crate::ir::ast::{Field, Node, NodeArray, NodeObject, NodeScalar, ValueData};
use crate::ir::tag::Tag;
use crate::ir::value_type::ValueType;
use crate::jsonc::scanner::{Scanner, Token, TokenType};
use chrono::NaiveDateTime;

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

            if tok.token_type == TokenType::Comment {
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
                    tag.value_type = ValueType::Str;
                }
                let text = tok.literal;

                if tag.is_null && !text.is_empty() {
                    return Err(format!("invalid string: {:?}, valid: {:?}", text, ""));
                }

                let data = match tag.value_type {
                    ValueType::Str => {
                        if tag.is_null {
                            if !text.is_empty() {
                                return Err(format!("invalid string: {:?}, valid: {:?}", text, ""));
                            }
                            ValueData::String(text.clone())
                        } else {
                            ValueData::String(text.clone())
                        }
                    }
                    ValueType::Bytes => {
                        let bytes = base64::Engine::decode(
                            &base64::engine::general_purpose::STANDARD,
                            &text,
                        )
                        .map_err(|e| format!("invalid base64 bytes {:?}: {}", text, e))?;
                        ValueData::Bytes(bytes)
                    }
                    ValueType::Datetime => {
                        let naive = NaiveDateTime::parse_from_str(&text, "%Y-%m-%d %H:%M:%S")
                            .map_err(|e| format!("invalid datetime '{}': {}", text, e))?;
                        let utc_ts = naive.and_utc().timestamp();
                        let timestamp = if let Some(loc) = tag.location {
                            utc_ts - (loc as i64 * 3600)
                        } else {
                            utc_ts
                        };
                        ValueData::Int(timestamp)
                    }
                    ValueType::Date => {
                        let _ = chrono::NaiveDate::parse_from_str(&text, "%Y-%m-%d")
                            .map_err(|e| format!("invalid date '{}': {}", text, e))?;
                        ValueData::String(text.clone())
                    }
                    ValueType::Time => {
                        let _ = chrono::NaiveTime::parse_from_str(&text, "%H:%M:%S")
                            .map_err(|e| format!("invalid time '{}': {}", text, e))?;
                        ValueData::String(text.clone())
                    }
                    ValueType::Uuid => {
                        let hex = text.replace('-', "");
                        if hex.len() != 32 {
                            return Err(format!("invalid uuid '{}'", text));
                        }
                        let bytes: Vec<u8> = (0..hex.len())
                            .step_by(2)
                            .map(|i| u8::from_str_radix(&hex[i..i + 2], 16))
                            .collect::<Result<Vec<_>, _>>()
                            .map_err(|e| format!("invalid uuid '{}': {}", text, e))?;
                        ValueData::Bytes(bytes)
                    }
                    ValueType::Decimal => {
                        if text.is_empty() {
                            return Err("invalid decimal".to_string());
                        }
                        ValueData::String(text.clone())
                    }
                    ValueType::Ip => {
                        let _ = text
                            .parse::<std::net::IpAddr>()
                            .map_err(|e| format!("invalid ip {:?}: {}", text, e))?;
                        ValueData::String(text.clone())
                    }
                    ValueType::Url => {
                        if !text.starts_with("http://") && !text.starts_with("https://") {
                            return Err(format!("invalid url {:?}", text));
                        }
                        ValueData::String(text.clone())
                    }
                    ValueType::Email => {
                        if !text.contains('@') || !text.contains('.') {
                            return Err(format!("invalid email {:?}", text));
                        }
                        ValueData::String(text.clone())
                    }
                    ValueType::Enum => {
                        if let Some(ref enums) = tag.enums {
                            if !enums.split('|').any(|e| e == text) {
                                return Err(format!(
                                    "invalid enum value {:?}, valid: {:?}",
                                    text, enums
                                ));
                            }
                        }
                        ValueData::String(text.clone())
                    }
                    ValueType::Media => {
                        let _bytes = base64::Engine::decode(
                            &base64::engine::general_purpose::STANDARD,
                            &text,
                        )
                        .map_err(|e| format!("invalid base64 media {:?}: {}", text, e))?;
                        ValueData::String(text.clone())
                    }
                    _ => {
                        return Err(format!(
                            "unsupported type {:?} for string literal",
                            tag.value_type
                        ));
                    }
                };
                let value = Node::Value(NodeScalar {
                    data,
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
                        tag.value_type = ValueType::F64;
                    } else {
                        tag.value_type = ValueType::I;
                    }
                }

                let is_float = text.contains('.');
                let is_negative = text.starts_with('-');

                // Validate type compatibility
                if is_float {
                    match tag.value_type {
                        ValueType::F32 | ValueType::F64 | ValueType::Decimal => {}
                        _ => {
                            return Err(format!(
                                "unsupported numeric type {:?} for float literal",
                                tag.value_type
                            ));
                        }
                    }
                } else if is_negative {
                    match tag.value_type {
                        ValueType::I
                        | ValueType::I8
                        | ValueType::I16
                        | ValueType::I32
                        | ValueType::I64
                        | ValueType::Bigint => {}
                        _ => {
                            return Err(format!(
                                "unsupported numeric type {:?} for negative literal",
                                tag.value_type
                            ));
                        }
                    }
                } else {
                    match tag.value_type {
                        ValueType::I
                        | ValueType::I8
                        | ValueType::I16
                        | ValueType::I32
                        | ValueType::I64
                        | ValueType::U
                        | ValueType::U8
                        | ValueType::U16
                        | ValueType::U32
                        | ValueType::U64
                        | ValueType::Bigint => {}
                        _ => {
                            return Err(format!("unsupported numeric type {:?}", tag.value_type));
                        }
                    }
                }

                let data: ValueData;
                if is_float {
                    data = ValueData::Float(text.parse::<f64>().unwrap_or(0.0));
                } else if is_negative {
                    if let Ok(ival) = text.parse::<i64>() {
                        data = ValueData::Int(ival);
                    } else {
                        data = ValueData::String(text.clone());
                    }
                } else if let Ok(uval) = text.parse::<u64>() {
                    if uval > i64::MAX as u64 {
                        data = ValueData::Uint(uval);
                    } else {
                        data = ValueData::Int(uval as i64);
                    }
                } else {
                    data = ValueData::String(text.clone());
                }

                let value = Node::Value(NodeScalar {
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
                match tag.value_type {
                    ValueType::Bool => {
                        if tag.is_null {
                            return Err("bool must false when bool is null".to_string());
                        }
                    }
                    _ => {
                        return Err(format!(
                            "unsupported type {:?} for boolean literal",
                            tag.value_type
                        ));
                    }
                }
                let value = Node::Value(NodeScalar {
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
                match tag.value_type {
                    ValueType::Bool => {}
                    _ => {
                        return Err(format!(
                            "unsupported type {:?} for boolean literal",
                            tag.value_type
                        ));
                    }
                }
                let value = Node::Value(NodeScalar {
                    data: ValueData::Bool(false),
                    text: "false".to_string(),
                    tag: Some(tag),
                    path: path.to_string(),
                });
                Ok(Some(value))
            }
            TokenType::Null => {
                let tag = self.consume_comments_for(tok.line).unwrap_or_default();
                if tag.value_type != ValueType::Unknown {
                    return Err(format!(
                        "null is not supported for type {:?}",
                        tag.value_type
                    ));
                }
                Ok(Some(Node::Value(NodeScalar {
                    data: ValueData::Null,
                    text: "null".to_string(),
                    tag: Some(tag),
                    path: path.to_string(),
                })))
            }
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
            tag.value_type = ValueType::Obj;
        }

        let obj_path = if let Some(ref name) = tag.name {
            if path.is_empty() {
                name.clone()
            } else {
                format!("{}.{}", path, name)
            }
        } else {
            path.to_string()
        };
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

            if tok.token_type == TokenType::Comment {
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

            let key_tok = self.next();
            if key_tok.token_type != TokenType::String {
                return Err("expected string key".to_string());
            }
            let key = camel_to_snake(&key_tok.literal);

            let colon = self.next();
            if colon.token_type != TokenType::Colon {
                return Err("expected colon".to_string());
            }

            let child_path = format!("{}.{}", obj_path, key);
            if let Some(mut val) = self.parse_node(&child_path)? {
                if let Some(mut ct) = val.get_tag().cloned() {
                    if tag.value_type == ValueType::Map {
                        ct.inherit(&tag);
                    }
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
        Ok(Node::Object(NodeObject {
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
            tag.value_type = ValueType::Vec;
        }

        let arr_path = if let Some(ref name) = tag.name {
            if path.is_empty() {
                name.clone()
            } else {
                format!("{}.{}", path, name)
            }
        } else {
            path.to_string()
        };
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

            if tok.token_type == TokenType::Comment {
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

            let item_path = format!("{}[{}]", arr_path, index);
            if let Some(mut item) = self.parse_node(&item_path)? {
                let mut ct = item.get_tag().cloned().unwrap_or_default();
                ct.inherit(&tag);
                if let Some(t) = item.get_tag_mut() {
                    *t = ct;
                }
                items.push(item);
                index += 1;
            }

            if self.peek().token_type == TokenType::Comma {
                self.next();
            }
        }

        self.depth -= 1;
        Ok(Node::Array(NodeArray {
            items,
            tag: Some(tag),
            path: arr_path,
        }))
    }
}
