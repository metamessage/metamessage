from typing import Any, Optional, List
import re
import json
from dataclasses import dataclass, field

from .tag import Tag, ValueType
from .types import Obj, Arr, Val, Field


@dataclass
class Node:
    path: str = ""
    tag: Optional[Tag] = None


@dataclass
class Token:
    type: str
    literal: Any = None
    line: int = 0


TokenString = "STRING"
TokenNumber = "NUMBER"
TokenLBrace = "LBRACE"
TokenRBrace = "RBRACE"
TokenLBracket = "LBRACKET"
TokenRBracket = "RBRACKET"
TokenComma = "COMMA"
TokenColon = "COLON"
TokenTrue = "TRUE"
TokenFalse = "FALSE"
TokenNull = "NULL"
TokenComment = "COMMENT"
TokenEOF = "EOF"


class Lexer:
    def __init__(self, source: str):
        self.source = source
        self.pos = 0
        self.line = 1
        self.tokens: List[Token] = []

    def is_at_end(self):
        return self.pos >= len(self.source)

    def advance(self):
        if self.is_at_end():
            return None
        c = self.source[self.pos]
        self.pos += 1
        if c == '\n':
            self.line += 1
        return c

    def peek(self):
        if self.is_at_end():
            return None
        return self.source[self.pos]

    def peek_next(self):
        if self.pos + 1 >= len(self.source):
            return None
        return self.source[self.pos + 1]

    def add_token(self, type: str, literal: Any = None):
        self.tokens.append(Token(type, literal, self.line))

    def scan_tokens(self):
        while not self.is_at_end():
            start = self.pos
            c = self.advance()
            
            if c == '{':
                self.add_token(TokenLBrace, "{")
            elif c == '}':
                self.add_token(TokenRBrace, "}")
            elif c == '[':
                self.add_token(TokenLBracket, "[")
            elif c == ']':
                self.add_token(TokenRBracket, "]")
            elif c == ',':
                self.add_token(TokenComma, ",")
            elif c == ':':
                self.add_token(TokenColon, ":")
            elif c in ' \t\n\r':
                pass
            elif c == '"':
                self.string()
            elif c == '/':
                if self.peek() == '/':
                    self.comment()
                elif self.peek() == '*':
                    self.block_comment()
            elif c.isdigit() or (c == '-' and self.peek().isdigit()):
                start_pos = self.pos - 1
                self.number(start_pos)
            elif c.isalpha():
                start_pos = self.pos - 1
                self.identifier(start_pos)
            else:
                pass
        
        self.add_token(TokenEOF)
        return self.tokens

    def string(self):
        start = self.pos - 1
        while self.peek() != '"' and not self.is_at_end():
            if self.peek() == '\\':
                self.advance()
            self.advance()
        
        self.advance()
        value = self.source[start:self.pos]
        self.add_token(TokenString, json.loads(value))

    def number(self, start_pos: int):
        while self.peek() and (self.peek().isdigit() or self.peek() == '.' or self.peek() == 'e' or self.peek() == 'E' or self.peek() == '+' or self.peek() == '-'):
            self.advance()
        
        value = self.source[start_pos:self.pos]
        if '.' in value or 'e' in value or 'E' in value:
            self.add_token(TokenNumber, float(value))
        else:
            self.add_token(TokenNumber, int(value))

    def identifier(self, start_pos: int):
        while self.peek() and (self.peek().isalnum() or self.peek() == '_'):
            self.advance()
        
        value = self.source[start_pos:self.pos].lower()
        if value == "true":
            self.add_token(TokenTrue, True)
        elif value == "false":
            self.add_token(TokenFalse, False)
        elif value == "null":
            self.add_token(TokenNull, None)
        else:
            self.add_token(TokenString, value)

    def comment(self):
        while self.peek() and self.peek() != '\n':
            self.advance()
        
        start = self.pos - 2
        value = self.source[start:self.pos]
        self.add_token(TokenComment, value)

    def block_comment(self):
        self.advance()
        while not self.is_at_end():
            if self.peek() == '*' and self.peek_next() == '/':
                self.advance()
                self.advance()
                break
            self.advance()
        
        start = self.pos - 4
        value = self.source[start:self.pos]
        self.add_token(TokenComment, value)


class Parser:
    def __init__(self, tokens: List[Token]):
        self.tokens = tokens
        self.pos = 0
        self.pending_comments: List[Token] = []

    def peek(self):
        if self.pos >= len(self.tokens):
            return Token(TokenEOF)
        return self.tokens[self.pos]

    def next(self):
        t = self.peek()
        self.pos += 1
        return t

    def parse(self, path: str = "") -> Node:
        tok = self.peek()
        
        if tok.type == TokenLBrace:
            return self.parse_object(path)
        elif tok.type == TokenLBracket:
            return self.parse_array(path)
        else:
            return self.parse_value(path)

    def parse_object(self, path: str) -> Obj:
        self.next()
        
        tag = self.consume_comments()
        
        fields = []
        while self.peek().type != TokenRBrace and self.peek().type != TokenEOF:
            if self.peek().type == TokenComment:
                self.next()
                continue
            
            key_token = self.next()
            if key_token.type != TokenString:
                break
            
            self.next()
            
            value_path = f"{path}.{key_token.literal}" if path else key_token.literal
            value = self.parse(value_path)
            
            fields.append(Field(key=key_token.literal, value=value))
            
            if self.peek().type == TokenComma:
                self.next()
        
        if self.peek().type == TokenRBrace:
            self.next()
        
        return Obj(fields=fields, tag=tag, path=path)

    def parse_array(self, path: str) -> Arr:
        self.next()
        
        tag = self.consume_comments()
        
        items = []
        index = 0
        while self.peek().type != TokenRBracket and self.peek().type != TokenEOF:
            if self.peek().type == TokenComment:
                self.next()
                continue
            
            item_path = f"{path}.{index}" if path else str(index)
            item = self.parse(item_path)
            items.append(item)
            index += 1
            
            if self.peek().type == TokenComma:
                self.next()
        
        if self.peek().type == TokenRBracket:
            self.next()
        
        return Arr(items=items, tag=tag, path=path)

    def parse_value(self, path: str) -> Val:
        tok = self.next()
        
        tag = self.consume_comments()
        
        if tag is None:
            tag = Tag()
        
        if tok.type == TokenString:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.String
            return Val(data=tok.literal, text=str(tok.literal), tag=tag, path=path)
        elif tok.type == TokenNumber:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Int
            return Val(data=tok.literal, text=str(tok.literal), tag=tag, path=path)
        elif tok.type == TokenTrue:
            tag.type = ValueType.Bool
            return Val(data=True, text="true", tag=tag, path=path)
        elif tok.type == TokenFalse:
            tag.type = ValueType.Bool
            return Val(data=False, text="false", tag=tag, path=path)
        elif tok.type == TokenNull:
            tag.is_null = True
            return Val(data=None, text="null", tag=tag, path=path)
        else:
            tag.type = ValueType.Unknown
            return Val(data=None, text="", tag=tag, path=path)

    def consume_comments(self) -> Optional[Tag]:
        if len(self.pending_comments) == 0:
            return None
        
        comments = []
        for tc in self.pending_comments:
            if tc.type == TokenComment:
                comments.append(tc.literal)
        
        self.pending_comments = []
        
        if not comments:
            return None
        
        merged = Tag()
        for c in comments:
            t = mm_tag_from_comment(c)
            if t:
                merged = merge_tag(merged, t)
        
        return merged


def mm_tag_from_comment(comment: str) -> Optional[Tag]:
    comment = comment.strip()
    if comment.startswith("//"):
        comment = comment[2:].strip()
    elif comment.startswith("/*"):
        comment = comment[2:]
        if comment.endswith("*/"):
            comment = comment[:-2]
        comment = comment.strip()
    if not comment.startswith("mm:"):
        return None
    tag_str = comment[3:].strip()
    if not tag_str:
        return None
    return mm_tag(tag_str)


def merge_tag(a: Tag, b: Tag) -> Tag:
    if a.type != ValueType.Unknown:
        b.type = a.type
    if a.desc:
        b.desc = a.desc
    if a.nullable:
        b.nullable = True
    if a.is_null:
        b.is_null = True
    if a.default:
        b.default = a.default
    if a.min:
        b.min = a.min
    if a.max:
        b.max = a.max
    if a.size:
        b.size = a.size
    if a.enum:
        b.enum = a.enum
    if a.pattern:
        b.pattern = a.pattern
    if a.version:
        b.version = a.version
    if a.mime:
        b.mime = a.mime
    if a.child_type != ValueType.Unknown:
        b.child_type = a.child_type
    if a.child_desc:
        b.child_desc = a.child_desc
    if a.child_nullable:
        b.child_nullable = True
    if a.child_default:
        b.child_default = a.child_default
    if a.child_min:
        b.child_min = a.child_min
    if a.child_max:
        b.child_max = a.child_max
    if a.child_size:
        b.child_size = a.child_size
    if a.child_enum:
        b.child_enum = a.child_enum
    
    return b


def parse_jsonc(source: str, tag_str: str = "") -> Node:
    lexer = Lexer(source)
    tokens = lexer.scan_tokens()
    parser = Parser(tokens)
    return parser.parse("")


def struct_to_mm(data: Any, path: str = "") -> Node:
    if isinstance(data, dict):
        fields = []
        for k, v in data.items():
            field_path = f"{path}.{k}" if path else k
            fields.append(Field(key=k, value=struct_to_mm(v, field_path)))
        return Obj(fields=fields, tag=Tag(), path=path)
    elif isinstance(data, list):
        items = [struct_to_mm(item, f"{path}.{i}") for i, item in enumerate(data)]
        return Arr(items=items, tag=Tag(type=ValueType.Array), path=path)
    else:
        tag = Tag(type=infer_type(data))
        return Val(data=data, text=str(data), tag=tag, path=path)


def infer_type(data: Any) -> ValueType:
    if data is None:
        return ValueType.Unknown
    elif isinstance(data, bool):
        return ValueType.Bool
    elif isinstance(data, int):
        return ValueType.Int
    elif isinstance(data, float):
        return ValueType.Float64
    elif isinstance(data, str):
        return ValueType.String
    elif isinstance(data, bytes):
        return ValueType.Bytes
    elif isinstance(data, list):
        return ValueType.Array
    elif isinstance(data, dict):
        return ValueType.Struct
    else:
        return ValueType.Unknown


def mm_tag(tag_str: str) -> Optional[Tag]:
    if not tag_str:
        return None
    
    tag_str = tag_str.strip()
    if not tag_str:
        return None
    
    tag = Tag()
    
    type_map = {
        "unknown": ValueType.Unknown, "str": ValueType.String, "i": ValueType.Int,
        "i8": ValueType.Int8, "i16": ValueType.Int16, "i32": ValueType.Int32, "i64": ValueType.Int64,
        "u": ValueType.Uint, "u8": ValueType.Uint8, "u16": ValueType.Uint16, "u32": ValueType.Uint32, "u64": ValueType.Uint64,
        "f32": ValueType.Float32, "f64": ValueType.Float64, "bool": ValueType.Bool, "bytes": ValueType.Bytes,
        "bi": ValueType.BigInt, "datetime": ValueType.DateTime, "date": ValueType.Date, "time": ValueType.Time,
        "uuid": ValueType.UUID, "decimal": ValueType.Decimal, "ip": ValueType.IP, "url": ValueType.URL,
        "email": ValueType.Email, "enum": ValueType.Enum, "arr": ValueType.Array, "struct": ValueType.Struct,
    }
    
    parts = [p.strip() for p in tag_str.split(";")]
    for p in parts:
        if not p:
            continue
        
        k, v = p, ""
        if "=" in p:
            kv = p.split("=", 1)
            k = kv[0].strip()
            v = kv[1].strip() if len(kv) > 1 else ""
        
        k = k.lower()
        
        if k == "is_null":
            tag.is_null = True
            tag.nullable = True
        elif k == "example":
            tag.example = True
        elif k == "desc":
            tag.desc = v
        elif k == "type":
            tag.type = type_map.get(v, ValueType.Unknown)
        elif k == "raw":
            tag.raw = True
        elif k == "nullable":
            tag.nullable = True
        elif k == "allow_empty":
            tag.allow_empty = True
        elif k == "unique":
            tag.unique = True
        elif k == "default":
            tag.default = v
        elif k == "min":
            tag.min = v
        elif k == "max":
            tag.max = v
        elif k == "size":
            tag.size = int(v) if v.isdigit() else 0
        elif k == "enum":
            tag.type = ValueType.Enum
            tag.enum = v
        elif k == "pattern":
            tag.pattern = v
        elif k == "location":
            tag.location = v
        elif k == "version":
            tag.version = int(v) if v.isdigit() else 0
        elif k == "mime":
            tag.mime = v
        elif k == "child_desc":
            tag.child_desc = v
        elif k == "child_type":
            tag.child_type = type_map.get(v, ValueType.Unknown)
        elif k == "child_raw":
            tag.child_raw = True
        elif k == "child_nullable":
            tag.child_nullable = True
        elif k == "child_allow_empty":
            tag.child_allow_empty = True
        elif k == "child_unique":
            tag.child_unique = True
        elif k == "child_default":
            tag.child_default = v
        elif k == "child_min":
            tag.child_min = v
        elif k == "child_max":
            tag.child_max = v
        elif k == "child_size":
            tag.child_size = int(v) if v.isdigit() else 0
        elif k == "child_enum":
            tag.child_enum = v
            tag.child_type = ValueType.Enum
        elif k == "child_pattern":
            tag.child_pattern = v
        elif k == "child_location":
            tag.child_location = v
        elif k == "child_version":
            tag.child_version = int(v) if v.isdigit() else 0
        elif k == "child_mime":
            tag.child_mime = v
    
    return tag


def to_jsonc(node: Node, indent: int = 0) -> str:
    import json
    
    indent_str = "    "
    
    if isinstance(node, Obj):
        lines = ["{\n"]
        for f in node.fields:
            if f.value and hasattr(f.value, 'tag') and f.value.tag:
                tag_str = str(f.value.tag)
                if tag_str:
                    lines.append(f"{indent_str * (indent + 1)}// mm: {tag_str}\n")
            
            lines.append(f'{indent_str * (indent + 1)}"{f.key}": ')
            lines.append(to_jsonc(f.value, indent + 1))
            lines.append(",\n")
        
        lines.append(f"{indent_str * indent}}}")
        return "".join(lines)
    
    elif isinstance(node, Arr):
        if not node.items:
            return "[]"
        
        lines = ["[\n"]
        for item in node.items:
            lines.append(f"{indent_str * (indent + 1)}")
            lines.append(to_jsonc(item, indent + 1))
            lines.append(",\n")
        
        lines.append(f"{indent_str * indent}]")
        return "".join(lines)
    
    elif isinstance(node, Val):
        if node.tag is None:
            return str(node.text) if node.text else ""
        
        val_type = node.tag.type
        
        if val_type in (ValueType.String, ValueType.Bytes, ValueType.DateTime, ValueType.Date, 
                     ValueType.Time, ValueType.UUID, ValueType.IP, ValueType.URL, 
                     ValueType.Email, ValueType.Enum, ValueType.Decimal):
            return json.dumps(node.text)
        
        return node.text if node.text else ""
    
    return ""


def write_leading_comment(node: Node, indent: int = 0) -> str:
    indent_str = "    "
    
    if hasattr(node, 'tag') and node.tag:
        tag_str = str(node.tag)
        if tag_str:
            return f"{indent_str * indent}// mm: {tag_str}\n"
    
    return ""