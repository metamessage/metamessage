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
TOKEN_EOF = "EOF"


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

    def add_token(self, type_: str, literal: Any = None):
        self.tokens.append(Token(type_, literal, self.line))

    def scan_tokens(self):
        while not self.is_at_end():
            c = self.advance()
            
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
                    self._comment()
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

    def _comment(self):
        start = self.pos - 2
        while self.peek() and self.peek() != '\n':
            self.advance()
        value = self.source[start:self.pos]
        self.add_token(TOKEN_COMMENT, value)

    def _block_comment(self):
        self.advance()  # skip *
        start = self.pos - 4
        while not self.is_at_end():
            if self.peek() == '*' and self.peek_next() == '/':
                self.advance()
                self.advance()
                break
            self.advance()
        value = self.source[start:self.pos]
        self.add_token(TOKEN_COMMENT, value)


# ===== Parser =====

class Parser:
    def __init__(self, tokens: List[Token]):
        self.tokens = tokens
        self.pos = 0
        self.pending_comments: List[str] = []

    def peek(self):
        if self.pos >= len(self.tokens):
            return Token(TOKEN_EOF)
        return self.tokens[self.pos]

    def next(self):
        t = self.peek()
        self.pos += 1
        return t

    def _collect_comment(self, comment: str):
        self.pending_comments.append(comment)

    def _consume_tag(self) -> Optional[Tag]:
        if not self.pending_comments:
            return None

        tags = []
        for c in self.pending_comments:
            t = self._parse_mm_tag(c)
            if t is not None:
                tags.append(t)

        self.pending_comments = []

        if not tags:
            return None

        merged = tags[0]
        for t in tags[1:]:
            merged = MergeTag(merged, t)
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
        # Collect any leading comments before value
        while self.peek().type == TOKEN_COMMENT:
            self._collect_comment(self.peek().literal)
            self.next()

        tok = self.peek()

        if tok.type == TOKEN_LBRACE:
            return self._parse_object(path)
        elif tok.type == TOKEN_LBRACKET:
            return self._parse_array(path)
        else:
            return self._parse_value(path)

    def _parse_object(self, path: str) -> Obj:
        self.next()  # consume {

        tag = self._consume_tag()

        fields = []
        while self.peek().type != TOKEN_RBRACE and self.peek().type != TOKEN_EOF:
            if self.peek().type == TOKEN_COMMENT:
                self._collect_comment(self.peek().literal)
                self.next()
                continue
            
            key_token = self.next()
            if key_token.type != TOKEN_STRING:
                break
            
            if self.peek().type == TOKEN_COLON:
                self.next()

            value_path = f"{path}.{key_token.literal}" if path else key_token.literal
            value = self.parse(value_path)
            
            # After parsing value, check for inline comment and merge into value tag
            if self.peek().type == TOKEN_COMMENT:
                inline_tag = self._parse_mm_tag(self.peek().literal)
                if inline_tag is not None and hasattr(value, 'tag'):
                    value.tag = MergeTag(value.tag, inline_tag)
                # Don't collect inline comments into pending - they're applied to the value
                self.next()
            
            fields.append(Field(key=key_token.literal, value=value))

            if self.peek().type == TOKEN_COMMA:
                self.next()

        tag2 = self._consume_tag()
        if tag is None:
            tag = tag2
        elif tag2 is not None:
            tag = MergeTag(tag, tag2)

        if self.peek().type == TOKEN_RBRACE:
            self.next()

        return Obj(fields=fields, tag=tag or NewTag(), path=path)

    def _parse_array(self, path: str) -> Arr:
        self.next()  # consume [

        tag = self._consume_tag()

        items = []
        index = 0
        while self.peek().type != TOKEN_RBRACKET and self.peek().type != TOKEN_EOF:
            if self.peek().type == TOKEN_COMMENT:
                self._collect_comment(self.peek().literal)
                self.next()
                continue

            item_path = f"{path}[{index}]" if path else str(index)
            item = self.parse(item_path)
            
            # After parsing item, check for inline comment and merge into item tag
            if self.peek().type == TOKEN_COMMENT:
                inline_tag = self._parse_mm_tag(self.peek().literal)
                if inline_tag is not None and hasattr(item, 'tag'):
                    item.tag = MergeTag(item.tag, inline_tag)
                self.next()
            
            items.append(item)
            index += 1

            if self.peek().type == TOKEN_COMMA:
                self.next()

        tag2 = self._consume_tag()
        if tag is None:
            tag = tag2
        elif tag2 is not None:
            tag = MergeTag(tag, tag2)

        if self.peek().type == TOKEN_RBRACKET:
            self.next()

        return Arr(items=items, tag=tag or NewTag(), path=path)

    def _parse_value(self, path: str) -> Val:
        tag = self._consume_tag()
        tok = self.next()
        
        if tag is None:
            tag = NewTag()

        if tok.type == TOKEN_STRING:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.String
            return Val(data=tok.literal, text=str(tok.literal), tag=tag, path=path)
        elif tok.type == TOKEN_NUMBER:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Float64 if isinstance(tok.literal, float) else ValueType.Int
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
            # In MM, null is represented via is_null tag, not bare null value
            tag.is_null = True
            tag.nullable = True
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.String
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
    ValueType.Object,
    ValueType.Slice,
    ValueType.Array,
    ValueType.String,
    ValueType.Int,
    ValueType.Float64,
    ValueType.Bool,
}


def _get_tag_str(tag) -> str:
    """Get tag string, omitting inferred types."""
    if tag is None:
        return ""
    
    # Build tag string manually to filter inferred types
    parts = []
    
    if tag.type != ValueType.Unknown and tag.type not in _INFERRED_TYPES:
        parts.append(f"type={tag.type.to_str()}")
    
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
        parts.append(f"child_type={tag.child_type.to_str()}")
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


def _to_jsonc(node: Node, indent: int = 0) -> str:
    INDENT = "    "

    if isinstance(node, Obj):
        tag_str = _get_tag_str(node.get_tag())

        lines = ["{"]
        if tag_str:
            lines[-1] += f" // mm: {tag_str}"
        lines[-1] += "\n"

        for i, f in enumerate(node.fields):
            val_tag_str = ""
            if f.value is not None:
                t = f.value.get_tag()
                if t is not None:
                    val_tag_str = _get_tag_str(t)

            if val_tag_str:
                lines.append(f'{INDENT * (indent + 1)}// mm: {val_tag_str}')
                if not isinstance(f.value, (Obj, Arr)):
                    lines[-1] += "\n"
                    lines.append(f'{INDENT * (indent + 1)}{json.dumps(f.key)}: ')
                else:
                    lines.append("\n")
                    lines.append(f'{INDENT * (indent + 1)}{json.dumps(f.key)}: ')
            else:
                lines.append(f'{INDENT * (indent + 1)}{json.dumps(f.key)}: ')
            lines.append(_to_jsonc(f.value, indent + 1))
            if i < len(node.fields) - 1:
                lines[-1] += ","
            lines[-1] += "\n"

        lines.append(f"{INDENT * indent}}}")
        return "".join(lines)
    
    elif isinstance(node, Arr):
        if not node.items:
            return "[]"

        tag_str = _get_tag_str(node.get_tag())

        lines = ["["]
        if tag_str:
            lines[-1] += f" // mm: {tag_str}"
        lines[-1] += "\n"

        for i, item in enumerate(node.items):
            item_tag_str = _get_tag_str(item.get_tag())
            if item_tag_str:
                lines.append(f"{INDENT * (indent + 1)}// mm: {item_tag_str}\n")
            lines.append(f"{INDENT * (indent + 1)}")
            lines.append(_to_jsonc(item, indent + 1))
            if i < len(node.items) - 1:
                lines[-1] += ","
            lines[-1] += "\n"

        lines.append(f"{INDENT * indent}]")
        return "".join(lines)
    
    elif isinstance(node, Val):
        if node.tag is None:
            return json.dumps(node.text) if node.text else "null"

        val_type = node.tag.type

        # For null values (is_null=True), use default value for type
        if node.tag.is_null:
            return _default_value_str(val_type)

        return _value_to_json_str(val_type, node.text)

    return "null"


def _default_value_str(val_type: ValueType) -> str:
    """Return the default JSON representation for a type (used for null values)."""
    if val_type in (ValueType.String, ValueType.Bytes, ValueType.DateTime,
                     ValueType.Date, ValueType.Time, ValueType.UUID,
                     ValueType.IP, ValueType.URL, ValueType.Email,
                     ValueType.Enum, ValueType.Decimal):
        return json.dumps("")
    elif val_type == ValueType.BigInt:
        return "0"
    elif val_type == ValueType.Bool:
        return "false"
    elif val_type in (ValueType.Int, ValueType.Uint, ValueType.Int8, ValueType.Int16,
                     ValueType.Int32, ValueType.Int64, ValueType.Uint8, ValueType.Uint16,
                     ValueType.Uint32, ValueType.Uint64):
        return "0"
    elif val_type in (ValueType.Float32, ValueType.Float64):
        return "0.0"
    elif val_type == ValueType.Unknown:
        return "null"
    return "null"


def _value_to_json_str(val_type: ValueType, text: str) -> str:
    """Convert a value to its JSON string representation based on type."""
    if val_type in (ValueType.String, ValueType.Bytes, ValueType.DateTime,
                     ValueType.Date, ValueType.Time, ValueType.UUID,
                     ValueType.IP, ValueType.URL, ValueType.Email,
                     ValueType.Enum, ValueType.Decimal):
        return json.dumps(text)
    elif val_type == ValueType.BigInt:
        return text if text else "0"
    elif val_type == ValueType.Bool:
        return text if text else "false"
    elif val_type in (ValueType.Int, ValueType.Uint, ValueType.Int8, ValueType.Int16,
                     ValueType.Int32, ValueType.Int64, ValueType.Uint8, ValueType.Uint16,
                     ValueType.Uint32, ValueType.Uint64):
        return text if text else "0"
    elif val_type in (ValueType.Float32, ValueType.Float64):
        return text if text else "0.0"
    elif val_type == ValueType.Unknown:
        return json.dumps(text) if text else "null"
    return json.dumps(text) if text else "null"


def to_jsonc(node: Node) -> str:
    """Convert a Node tree to JSONC string."""
    return _to_jsonc(node, 0)