"""
JSONC Parser and Generator for metamessage.
Parses JSON with mm: comments and generates JSONC output.
"""

import json
from typing import Any, List, Optional
from dataclasses import dataclass

from ..ir.tag import Tag, ValueType, mm_tag, MergeTag, NewTag
from ..ir.types import Obj, Arr, Val, Field, Node


# ===== Tokenizer =====

@dataclass
class Token:
    type: str
    literal: Any = None
    line: int = 0


TOKEN_STRING = "STRING"
TOKEN_NUMBER = "NUMBER"
TOKEN_LBRACE = "LBRACE"
TOKEN_RBRACE = "RBRACE"
TOKEN_LBRACKET = "LBRACKET"
TOKEN_RBRACKET = "RBRACKET"
TOKEN_COMMA = "COMMA"
TOKEN_COLON = "COLON"
TOKEN_TRUE = "TRUE"
TOKEN_FALSE = "FALSE"
TOKEN_NULL = "NULL"
TOKEN_COMMENT = "COMMENT"
TOKEN_LEADING_COMMENT = "LEADING_COMMENT"
TOKEN_TRAILING_COMMENT = "TRAILING_COMMENT"
TOKEN_EOF = "EOF"


class Lexer:
    def __init__(self, source: str):
        self.source = source
        self.pos = 0
        self.line = 1
        self.col = 1
        self.new_line = True
        self.tokens: List[Token] = []

    def is_at_end(self):
        return self.pos >= len(self.source)

    def advance(self):
        if self.is_at_end():
            return None
        c = self.source[self.pos]
        self.pos += 1
        if c == '\n':
            self.new_line = True
            self.line += 1
            self.col = 1
        else:
            self.col += 1
        return c

    def peek(self):
        if self.is_at_end():
            return None
        return self.source[self.pos]

    def peek_next(self):
        if self.pos + 1 >= len(self.source):
            return None
        return self.source[self.pos + 1]

    def add_token(self, type_: str, literal: Any = None):
        self.tokens.append(Token(type_, literal, self.line))

    def scan_tokens(self):
        while not self.is_at_end():
            c = self.advance()
            start_line = self.line
            
            if c == '{':
                self.add_token(TOKEN_LBRACE, "{")
            elif c == '}':
                self.add_token(TOKEN_RBRACE, "}")
            elif c == '[':
                self.add_token(TOKEN_LBRACKET, "[")
            elif c == ']':
                self.add_token(TOKEN_RBRACKET, "]")
            elif c == ',':
                self.add_token(TOKEN_COMMA, ",")
                self.new_line = False
            elif c == ':':
                self.add_token(TOKEN_COLON, ":")
                self.new_line = False
            elif c in ' \t\n\r':
                pass
            elif c == '"':
                self._string()
            elif c == '/':
                if self.peek() == '/':
                    self._line_comment()
                elif self.peek() == '*':
                    self._block_comment()
            elif c in '0123456789' or c == '-':
                start_pos = self.pos - 1
                self._number(start_pos)
            elif c.isalpha():
                start_pos = self.pos - 1
                self._identifier(start_pos)
            else:
                pass

        self.add_token(TOKEN_EOF)
        return self.tokens

    def _string(self):
        s = self._read_string()
        self.add_token(TOKEN_STRING, s)

    def _read_string(self) -> str:
        result = ""
        while not self.is_at_end():
            c = self.advance()
            if c == '"':
                return result
            elif c == '\\':
                if self.is_at_end():
                    break
                n = self.advance()
                if n == '"':
                    result += '"'
                elif n == '\\':
                    result += '\\'
                elif n == '/':
                    result += '/'
                elif n == 'b':
                    result += '\b'
                elif n == 'f':
                    result += '\f'
                elif n == 'n':
                    result += '\n'
                elif n == 'r':
                    result += '\r'
                elif n == 't':
                    result += '\t'
                elif n == 'u':
                    hex_str = ""
                    for _ in range(4):
                        hex_str += self.advance()
                    result += chr(int(hex_str, 16))
                else:
                    result += n
            else:
                result += c
        return result

    def _number(self, start_pos: int):
        while self.peek() and self.peek() in '0123456789.eE+-':
            self.advance()
        
        value = self.source[start_pos:self.pos]
        if '.' in value or 'e' in value or 'E' in value:
            self.add_token(TOKEN_NUMBER, float(value))
        else:
            self.add_token(TOKEN_NUMBER, int(value))

    def _identifier(self, start_pos: int):
        while self.peek() and (self.peek().isalnum() or self.peek() == '_'):
            self.advance()
        
        value = self.source[start_pos:self.pos].lower()
        if value == "true":
            self.add_token(TOKEN_TRUE, True)
        elif value == "false":
            self.add_token(TOKEN_FALSE, False)
        elif value == "null":
            self.add_token(TOKEN_NULL, None)
        else:
            self.add_token(TOKEN_STRING, value)

    def _line_comment(self):
        comment_type = TOKEN_LEADING_COMMENT if self.new_line else TOKEN_TRAILING_COMMENT
        self.advance()  # skip second /
        start = self.pos
        while self.peek() and self.peek() != '\n':
            self.advance()
        value = self.source[start:self.pos].strip()
        self.add_token(comment_type, value)

    def _block_comment(self):
        self.advance()  # skip *
        comment_type = TOKEN_LEADING_COMMENT if self.new_line else TOKEN_TRAILING_COMMENT
        start = self.pos
        while not self.is_at_end():
            if self.peek() == '*' and self.peek_next() == '/':
                self.advance()
                self.advance()
                break
            self.advance()
        value = self.source[start:self.pos - 2].strip()
        self.add_token(comment_type, value)


# ===== Parser =====

class Parser:
    def __init__(self, tokens: List[Token]):
        self.tokens = tokens
        self.pos = 0
        self.pending: List[Token] = []
        self.depth = 0

    def peek(self):
        if self.pos >= len(self.tokens):
            return Token(TOKEN_EOF)
        return self.tokens[self.pos]

    def next(self):
        t = self.peek()
        self.pos += 1
        return t

    def _consume_comments_for(self, anchor_line: int) -> Optional[Tag]:
        if not self.pending:
            return None

        last = self.pending[-1]
        if anchor_line - last.line > 1:
            self.pending = []
            return None

        merged = None
        for ct in self.pending:
            parsed = self._parse_mm_tag(ct.literal)
            if parsed is not None:
                merged = MergeTag(merged, parsed)

        self.pending = []
        return merged

    @staticmethod
    def _parse_mm_tag(comment: str) -> Optional[Tag]:
        comment = comment.strip()
        if comment.startswith("//"):
            comment = comment[2:].strip()
        elif comment.startswith("/*"):
            comment = comment[2:]
            if comment.endswith("*/"):
                comment = comment[:-2].strip()

        if not comment.startswith("mm:"):
            return None

        tag_str = comment[3:].strip()
        if not tag_str:
            return None

        return mm_tag(tag_str)

    def parse(self, path: str = "") -> Node:
        while True:
            tok = self.peek()
            if tok.type == TOKEN_EOF:
                return None

            if tok.type in (TOKEN_LEADING_COMMENT, TOKEN_COMMENT):
                if self.pending:
                    last = self.pending[-1]
                    if tok.line - last.line > 1:
                        self.pending = []
                self.pending.append(tok)
                self.next()
                continue

            if tok.type == TOKEN_TRAILING_COMMENT:
                self.next()
                continue

            return self._parse_value_or_container(path)

    def _parse_value_or_container(self, path: str) -> Node:
        tok = self.peek()

        if tok.type == TOKEN_LBRACE:
            return self._parse_object(tok.line, path)
        elif tok.type == TOKEN_LBRACKET:
            return self._parse_array(tok.line, path)
        else:
            return self._parse_value(path, tok.line)

    def _parse_object(self, anchor_line: int, path: str) -> Obj:
        self.next()  # consume {
        self.depth += 1

        tag = self._consume_comments_for(anchor_line)

        fields = []
        while self.peek().type != TOKEN_RBRACE and self.peek().type != TOKEN_EOF:
            tok = self.peek()

            if tok.type in (TOKEN_LEADING_COMMENT, TOKEN_COMMENT):
                self.pending.append(tok)
                self.next()
                continue

            if tok.type == TOKEN_TRAILING_COMMENT:
                self.next()
                continue

            key_token = self.next()
            if key_token.type != TOKEN_STRING:
                break

            if self.peek().type == TOKEN_COLON:
                self.next()

            value_path = f"{path}.{key_token.literal}" if path else key_token.literal
            value = self._parse_value_or_container(value_path)

            if self.peek().type == TOKEN_TRAILING_COMMENT:
                inline_tag = self._parse_mm_tag(self.peek().literal)
                if inline_tag is not None and hasattr(value, 'tag'):
                    value.tag = MergeTag(value.tag, inline_tag)
                self.next()

            fields.append(Field(key=key_token.literal, value=value))

            if self.peek().type == TOKEN_COMMA:
                self.next()

        tag2 = self._consume_comments_for(anchor_line)
        if tag is None:
            tag = tag2
        elif tag2 is not None:
            tag = MergeTag(tag, tag2)

        if self.peek().type == TOKEN_RBRACE:
            self.next()

        self.depth -= 1
        return Obj(fields=fields, tag=tag or NewTag(), path=path)

    def _parse_array(self, anchor_line: int, path: str) -> Arr:
        self.next()  # consume [
        self.depth += 1

        tag = self._consume_comments_for(anchor_line)

        items = []
        index = 0
        while self.peek().type != TOKEN_RBRACKET and self.peek().type != TOKEN_EOF:
            tok = self.peek()

            if tok.type in (TOKEN_LEADING_COMMENT, TOKEN_COMMENT):
                self.pending.append(tok)
                self.next()
                continue

            if tok.type == TOKEN_TRAILING_COMMENT:
                self.next()
                continue

            item_path = f"{path}[{index}]" if path else str(index)
            item = self._parse_value_or_container(item_path)

            if self.peek().type == TOKEN_TRAILING_COMMENT:
                inline_tag = self._parse_mm_tag(self.peek().literal)
                if inline_tag is not None and hasattr(item, 'tag'):
                    item.tag = MergeTag(item.tag, inline_tag)
                self.next()

            items.append(item)
            index += 1

            if self.peek().type == TOKEN_COMMA:
                self.next()

        tag2 = self._consume_comments_for(anchor_line)
        if tag is None:
            tag = tag2
        elif tag2 is not None:
            tag = MergeTag(tag, tag2)

        if self.peek().type == TOKEN_RBRACKET:
            self.next()

        self.depth -= 1
        return Arr(items=items, tag=tag or NewTag(), path=path)

    def _parse_value(self, path: str, anchor_line: int) -> Val:
        tag = self._consume_comments_for(anchor_line)
        tok = self.next()

        if tag is None:
            tag = NewTag()

        if tok.type == TOKEN_STRING:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Str
            return Val(data=tok.literal, text=str(tok.literal), tag=tag, path=path)
        elif tok.type == TOKEN_NUMBER:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.F64 if isinstance(tok.literal, float) else ValueType.I
            return Val(data=tok.literal, text=str(tok.literal), tag=tag, path=path)
        elif tok.type == TOKEN_TRUE:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Bool
            return Val(data=True, text="true", tag=tag, path=path)
        elif tok.type == TOKEN_FALSE:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Bool
            return Val(data=False, text="false", tag=tag, path=path)
        elif tok.type == TOKEN_NULL:
            tag.is_null = True
            tag.nullable = True
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Str
            return Val(data=None, text="null", tag=tag, path=path)
        else:
            return Val(data=None, text="", tag=tag, path=path)


# ===== Public API =====

def parse_jsonc(source: str) -> Node:
    """Parse JSONC source string into a Node tree."""
    lexer = Lexer(source)
    tokens = lexer.scan_tokens()
    parser = Parser(tokens)
    return parser.parse("")


# Types that can be inferred and don't need explicit tags
_INFERRED_TYPES = {
    ValueType.Obj,
    ValueType.Vec,
    ValueType.Arr,
    ValueType.Str,
    ValueType.I,
    ValueType.F64,
    ValueType.Bool,
}


def _get_tag_str(tag) -> str:
    """Get tag string, omitting inferred types."""
    if tag is None:
        return ""
    
    # Build tag string manually to filter inferred types
    parts = []
    
    if tag.type != ValueType.Unknown and tag.type not in _INFERRED_TYPES:
        parts.append(f"type={str(tag.type)}")
    
    if tag.example:
        parts.append("example")
    if tag.is_null:
        parts.append("is_null")
    if tag.nullable and not tag.is_null and not tag.is_inherit:
        parts.append("nullable")
    if tag.desc:
        parts.append(f'desc="{tag.desc}"')
    if tag.raw:
        parts.append("raw")
    if tag.allow_empty:
        parts.append("allow_empty")
    if tag.unique:
        parts.append("unique")
    if tag.default:
        parts.append(f"default={tag.default}")
    if tag.min:
        parts.append(f"min={tag.min}")
    if tag.max:
        parts.append(f"max={tag.max}")
    if tag.size:
        parts.append(f"size={tag.size}")
    if tag.enum:
        parts.append(f"enum={tag.enum}")
    if tag.pattern:
        parts.append(f"pattern={tag.pattern}")
    if tag.version:
        parts.append(f"version={tag.version}")
    if tag.mime:
        parts.append(f"mime={tag.mime}")
    
    if tag.child_desc:
        parts.append(f'child_desc="{tag.child_desc}"')
    if tag.child_type != ValueType.Unknown and tag.child_type not in _INFERRED_TYPES:
        parts.append(f"child_type={str(tag.child_type)}")
    if tag.child_raw:
        parts.append("child_raw")
    if tag.child_nullable:
        parts.append("child_nullable")
    if tag.child_allow_empty:
        parts.append("child_allow_empty")
    if tag.child_unique:
        parts.append("child_unique")
    if tag.child_default:
        parts.append(f"child_default={tag.child_default}")
    if tag.child_min:
        parts.append(f"child_min={tag.child_min}")
    if tag.child_max:
        parts.append(f"child_max={tag.child_max}")
    if tag.child_size:
        parts.append(f"child_size={tag.child_size}")
    if tag.child_enum:
        parts.append(f"child_enum={tag.child_enum}")
    if tag.child_pattern:
        parts.append(f"child_pattern={tag.child_pattern}")
    if tag.child_version:
        parts.append(f"child_version={tag.child_version}")
    if tag.child_mime:
        parts.append(f"child_mime={tag.child_mime}")
    
    return "; ".join(parts)


INDENT = "\t"


def write_indent(b: list, indent: int):
    b.append(INDENT * indent)


def write_value_jsonc(b: list, v) -> None:
    if v is None:
        return
    if v.tag is None:
        return

    val_type = v.tag.type

    if v.tag.is_null:
        if val_type in (ValueType.Str, ValueType.Bytes, ValueType.Datetime,
                        ValueType.Date, ValueType.Time, ValueType.Uuid,
                        ValueType.Ip, ValueType.Url, ValueType.Email,
                        ValueType.Enum, ValueType.Decimal):
            b.append('""')
        elif val_type in (ValueType.I, ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64,
                          ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64,
                          ValueType.Bigint):
            b.append("0")
        elif val_type == ValueType.Bool:
            b.append("false")
        elif val_type in (ValueType.F32, ValueType.F64):
            b.append("0.0")
        else:
            b.append("null")
        return

    if val_type in (ValueType.Str, ValueType.Bytes, ValueType.Datetime,
                    ValueType.Date, ValueType.Time, ValueType.Uuid,
                    ValueType.Ip, ValueType.Url, ValueType.Email,
                    ValueType.Enum):
        b.append(json.dumps(v.text))
    elif val_type in (ValueType.I, ValueType.I8, ValueType.I16, ValueType.I32, ValueType.I64,
                      ValueType.U, ValueType.U8, ValueType.U16, ValueType.U32, ValueType.U64,
                      ValueType.Bigint, ValueType.Decimal, ValueType.Bool):
        b.append(v.text)
    elif val_type in (ValueType.F32, ValueType.F64):
        b.append(v.text)
    else:
        b.append(v.text)


def write_leading_comments(b: list, tag, indent: int):
    tag_str = _get_tag_str(tag)
    if tag_str:
        b.append("\n")
        write_indent(b, indent)
        b.append(f"// mm: {tag_str}\n")


def write_node_jsonc(b: list, n: Node, indent: int):
    if isinstance(n, Val):
        write_value_jsonc(b, n)
    elif isinstance(n, Obj):
        write_object_jsonc(b, n, indent)
    elif isinstance(n, Arr):
        write_array_jsonc(b, n, indent)


def write_object_jsonc(b: list, o: Obj, indent: int):
    b.append("{\n")
    for f in o.fields:
        write_leading_comments(b, f.value.get_tag(), indent + 1)
        write_indent(b, indent + 1)
        b.append(json.dumps(f.key))
        b.append(": ")
        write_node_jsonc(b, f.value, indent + 1)
        b.append(",\n")
    write_indent(b, indent)
    b.append("}")


def write_array_jsonc(b: list, a: Arr, indent: int):
    b.append("[\n")
    for item in a.items:
        write_leading_comments(b, item.get_tag(), indent + 1)
        write_indent(b, indent + 1)
        write_node_jsonc(b, item, indent + 1)
        b.append(",\n")
    write_indent(b, indent)
    b.append("]")


def to_jsonc(node: Node) -> str:
    if node is None:
        return ""
    b: List[str] = []
    write_leading_comments(b, node.get_tag(), 0)
    write_node_jsonc(b, node, 0)
    return "".join(b)