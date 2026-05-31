"""
JSONC Parser and Generator for metamessage.
Parses JSON with mm: comments and generates JSONC output.
"""

import json
from typing import Any, List, Optional
from dataclasses import dataclass

from ..ir.tag import Tag, ValueType, mm_tag, MergeTag, NewTag
from ..ir.ast import Obj, Arr, Val, Field, Node


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
TOKEN_EOF = "EOF"


class Lexer:
    def __init__(self, source: str):
        self.source = source
        self.pos = 0
        self.line = 1
        self.col = 1
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
            elif c == ':':
                self.add_token(TOKEN_COLON, ":")
            elif c in ' \t\n\r':
                pass
            elif c == '"':
                self._string()
            elif c == '/':
                if self.peek() == '/':
                    self._line_comment()
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
        self.advance()  # skip second /
        start = self.pos
        while self.peek() and self.peek() != '\n':
            self.advance()
        value = self.source[start:self.pos].strip()
        self.add_token(TOKEN_COMMENT, value)


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

            if tok.type == TOKEN_COMMENT:
                if self.pending:
                    last = self.pending[-1]
                    if tok.line - last.line > 1:
                        self.pending = []
                self.pending.append(tok)
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

            if tok.type == TOKEN_COMMENT:
                self.pending.append(tok)
                self.next()
                continue

            key_token = self.next()
            if key_token.type != TOKEN_STRING:
                break

            if self.peek().type == TOKEN_COLON:
                self.next()

            value_path = f"{path}.{key_token.literal}" if path else key_token.literal
            value = self._parse_value_or_container(value_path)

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

            if tok.type == TOKEN_COMMENT:
                self.pending.append(tok)
                self.next()
                continue

            item_path = f"{path}[{index}]" if path else str(index)
            item = self._parse_value_or_container(item_path)

            if tag is not None and item is not None:
                item_tag = item.get_tag()
                if item_tag is not None:
                    item_tag.inherit(tag)

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
            raise Exception("null is not supported")
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
    ValueType.Media,
}


def _get_tag_str(tag) -> str:
    """Get tag string, omitting inferred types."""
    if tag is None:
        return ""

    inh = tag.is_inherit

    # Build tag string manually to filter inferred types
    parts = []

    type_key = "child_type" if inh else "type"
    if tag.type != ValueType.Unknown and tag.type not in _INFERRED_TYPES:
        if not (tag.type == ValueType.Arr and tag.size > 0 or
                tag.type == ValueType.Enums and tag.enums):
            parts.append(f"{type_key}={str(tag.type)}")

    if tag.example:
        parts.append("example")
    if tag.is_null:
        parts.append("is_null")

    nullable_key = "child_nullable" if inh else "nullable"
    if tag.nullable and not tag.is_null:
        parts.append(nullable_key)

    desc_key = "child_desc" if inh else "desc"
    if tag.desc:
        parts.append(f'{desc_key}="{tag.desc}"')

    if tag.deprecated and not inh:
        parts.append("deprecated")

    allow_empty_key = "child_allow_empty" if inh else "allow_empty"
    if tag.allow_empty:
        parts.append(allow_empty_key)

    unique_key = "child_unique" if inh else "unique"
    if tag.unique:
        parts.append(unique_key)

    default_key = "child_default_val" if inh else "default_val"
    if tag.default_val:
        parts.append(f"{default_key}={tag.default_val}")

    min_key = "child_min" if inh else "min"
    if tag.min:
        parts.append(f"{min_key}={tag.min}")

    max_key = "child_max" if inh else "max"
    if tag.max:
        parts.append(f"{max_key}={tag.max}")

    size_key = "child_size" if inh else "size"
    if tag.size:
        parts.append(f"{size_key}={tag.size}")

    enums_key = "child_enums" if inh else "enums"
    if tag.enums:
        parts.append(f"{enums_key}={tag.enums}")

    pattern_key = "child_pattern" if inh else "pattern"
    if tag.pattern:
        parts.append(f"{pattern_key}={tag.pattern}")

    location_offset_hour = 0
    if tag.location is not None:
        try:
            location_offset_hour = int(tag.location)
        except (ValueError, TypeError):
            pass
    if location_offset_hour != 0 and not inh:
        parts.append(f"location={location_offset_hour}")

    version_key = "child_version" if inh else "version"
    if tag.version != 0:
        parts.append(f"{version_key}={tag.version}")

    mime_key = "child_mime" if inh else "mime"
    if tag.mime:
        parts.append(f"{mime_key}={tag.mime}")

    if tag.child_desc and not inh:
        parts.append(f'child_desc="{tag.child_desc}"')
    if tag.child_type != ValueType.Unknown and tag.child_type not in _INFERRED_TYPES and not inh:
        if not (tag.child_type == ValueType.Arr and tag.child_size > 0 or
                tag.child_type == ValueType.Enums and tag.child_enums):
            parts.append(f"child_type={str(tag.child_type)}")
    if tag.child_nullable and not inh:
        parts.append("child_nullable")
    if tag.child_allow_empty and not inh:
        parts.append("child_allow_empty")
    if tag.child_unique and not inh:
        parts.append("child_unique")
    if tag.child_default_val and not inh:
        parts.append(f"child_default_val={tag.child_default_val}")
    if tag.child_min and not inh:
        parts.append(f"child_min={tag.child_min}")
    if tag.child_max and not inh:
        parts.append(f"child_max={tag.child_max}")
    if tag.child_size and not inh:
        parts.append(f"child_size={tag.child_size}")
    if tag.child_enums and not inh:
        parts.append(f"child_enums={tag.child_enums}")
    if tag.child_pattern and not inh:
        parts.append(f"child_pattern={tag.child_pattern}")
    child_location_offset_hour = 0
    if tag.child_location is not None:
        try:
            child_location_offset_hour = int(tag.child_location)
        except (ValueError, TypeError):
            pass
    if child_location_offset_hour != 0 and not inh:
        parts.append(f"child_location={child_location_offset_hour}")
    if tag.child_version != 0 and not inh:
        parts.append(f"child_version={tag.child_version}")
    if tag.child_mime and not inh:
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
                        ValueType.Enums):
            b.append('""')
        elif val_type in (ValueType.Bool,):
            b.append("false")
        else:
            b.append("0")
        return

    if val_type in (ValueType.Str, ValueType.Bytes, ValueType.Datetime,
                    ValueType.Date, ValueType.Time, ValueType.Uuid,
                    ValueType.Ip, ValueType.Url, ValueType.Email,
                    ValueType.Enums, ValueType.Media):
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


def _tag_has_child(tag) -> bool:
    if tag is None:
        return False
    return (tag.child_desc != "" or
            tag.child_type != ValueType.Unknown or
            tag.child_nullable or
            tag.child_allow_empty or
            tag.child_unique or
            tag.child_default_val != "" or
            tag.child_min != "" or
            tag.child_max != "" or
            tag.child_size != 0 or
            tag.child_enums != "" or
            tag.child_pattern != "" or
            tag.child_version != 0 or
            tag.child_mime != "")


def write_array_jsonc(b: list, a: Arr, indent: int):
    b.append("[\n")
    for item in a.items:
        item_tag = item.get_tag()
        comment_tag = item_tag
        if a.tag is not None and _tag_has_child(a.tag):
            comment_tag = NewTag()
            if item_tag is not None:
                comment_tag = MergeTag(comment_tag, item_tag)
            comment_tag.inherit(a.tag)
        write_leading_comments(b, comment_tag, indent + 1)
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