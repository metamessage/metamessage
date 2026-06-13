"""
JSONC Parser and Generator for metamessage.
Parses JSON with mm: comments and generates JSONC output.
"""

import json
from typing import Any, List, Optional
from dataclasses import dataclass

from ..ir.tag import Tag, ValueType, mm_tag, MergeTag, NewTag
from ..ir.ast import NodeObject, NodeArray, NodeScalar, Field, Node, NodeNull
from ..ir.validator import MmValidator


# ===== Tokenizer =====

@dataclass
class Token:
    type: str
    literal: Any = None
    line: int = 0
    column: int = 0


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

    def add_token(self, type_: str, literal: Any = None, line: int = None, column: int = None):
        self.tokens.append(Token(type_, literal, line or self.line, column or self.col))

    def scan_tokens(self):
        while not self.is_at_end():
            start_line = self.line
            start_col = self.col
            c = self.advance()

            if c == '{':
                self.add_token(TOKEN_LBRACE, "{", start_line, start_col)
            elif c == '}':
                self.add_token(TOKEN_RBRACE, "}", start_line, start_col)
            elif c == '[':
                self.add_token(TOKEN_LBRACKET, "[", start_line, start_col)
            elif c == ']':
                self.add_token(TOKEN_RBRACKET, "]", start_line, start_col)
            elif c == ',':
                self.add_token(TOKEN_COMMA, ",", start_line, start_col)
            elif c == ':':
                self.add_token(TOKEN_COLON, ":", start_line, start_col)
            elif c in ' \t\n\r':
                pass
            elif c == '"':
                self._string(start_line, start_col)
            elif c == '/':
                if self.peek() == '/':
                    self._line_comment(start_line, start_col)
            elif c in '0123456789' or c == '-':
                start_pos = self.pos - 1
                self._number(start_pos, start_line, start_col)
            elif c.isalpha():
                start_pos = self.pos - 1
                self._identifier(start_pos, start_line, start_col)
            else:
                pass

        self.add_token(TOKEN_EOF)
        return self.tokens

    def _string(self, start_line: int, start_col: int):
        s = self._read_string()
        self.add_token(TOKEN_STRING, s, start_line, start_col)

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

    def _number(self, start_pos: int, start_line: int, start_col: int):
        while self.peek() and self.peek() in '0123456789.eE+-':
            self.advance()

        value = self.source[start_pos:self.pos]
        self.add_token(TOKEN_NUMBER, value, start_line, start_col)

    def _identifier(self, start_pos: int, start_line: int, start_col: int):
        while self.peek() and (self.peek().isalnum() or self.peek() == '_'):
            self.advance()

        value = self.source[start_pos:self.pos]
        if value == "true":
            self.add_token(TOKEN_TRUE, True, start_line, start_col)
        elif value == "false":
            self.add_token(TOKEN_FALSE, False, start_line, start_col)
        elif value == "null":
            self.add_token(TOKEN_NULL, None, start_line, start_col)
        else:
            self.add_token(TOKEN_STRING, value, start_line, start_col)

    def _line_comment(self, start_line: int, start_col: int):
        self.advance()  # skip second /
        start = self.pos
        while self.peek() and self.peek() != '\n':
            self.advance()
        value = self.source[start:self.pos].strip()
        self.add_token(TOKEN_COMMENT, value, start_line, start_col)


# ===== Parser =====

MAX_DEPTH = 32


def camel_to_snake(name: str) -> str:
    """Convert CamelCase to snake_case.
    
    UserID → user_id
    HTTPRequest → http_request
    """
    if not name:
        return name
    result = []
    for i, c in enumerate(name):
        if c.isupper():
            if i > 0 and (i + 1 < len(name) and name[i + 1].islower() or
                          not name[i - 1].isupper()):
                result.append('_')
            result.append(c.lower())
        else:
            result.append(c)
    return ''.join(result)


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
        val = None
        while True:
            tok = self.peek()
            if tok.type == TOKEN_EOF:
                if val is None:
                    raise Exception("no value parsed")
                return val

            if tok.type == TOKEN_COMMENT:
                if self.pending:
                    last = self.pending[-1]
                    if tok.line - last.line > 1:
                        self.pending = []
                self.pending.append(tok)
                self.next()
                continue

            val = self._parse_value_or_container(path)

    def _parse_value_or_container(self, path: str) -> Node:
        tok = self.peek()

        if tok.type == TOKEN_LBRACE:
            return self._parse_object(tok.line, path)
        elif tok.type == TOKEN_LBRACKET:
            return self._parse_array(tok.line, path)
        else:
            return self._parse_value(path, tok.line)

    def _parse_object(self, anchor_line: int, path: str) -> NodeObject:
        self.next()  # consume {
        self.depth += 1
        if self.depth > MAX_DEPTH:
            raise Exception(f"max depth: {MAX_DEPTH}")

        tag = self._consume_comments_for(anchor_line)
        if tag is None:
            tag = NewTag()

        if tag.name:
            path = f"{path}.{tag.name}" if path else tag.name

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

            key = camel_to_snake(key_token.literal)
            value_path = f"{path}.{key}" if path else key
            value = self._parse_value_or_container(value_path)

            fields.append(Field(key=key, value=value))

            if self.peek().type == TOKEN_COMMA:
                self.next()

        # Clear stale trailing comments before collecting tag2.
        # Use the closing } line as anchor instead of the opening { line.
        if self.pending:
            first = self.pending[0]
            closing_line = self.peek().line  # line of }
            if closing_line - first.line > 1:
                self.pending = []

        tag2 = self._consume_comments_for(self.peek().line)
        if tag is None:
            tag = tag2
        elif tag2 is not None:
            tag = MergeTag(tag, tag2)

        if self.peek().type == TOKEN_RBRACE:
            self.next()

        self.depth -= 1
        return NodeObject(fields=fields, tag=tag or NewTag(), path=path)

    def _parse_array(self, anchor_line: int, path: str) -> NodeArray:
        self.next()  # consume [
        self.depth += 1
        if self.depth > MAX_DEPTH:
            raise Exception(f"max depth: {MAX_DEPTH}")

        tag = self._consume_comments_for(anchor_line)

        if tag is None:
            tag = NewTag()

        if tag.type == ValueType.Unknown:
            tag.type = ValueType.Vec

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
        arr = NodeArray(items=items, tag=tag or NewTag(), path=path)
        r = MmValidator.validate_arr(arr, tag)

        if r.error:
            raise Exception(f"Validation error at {path}: {r.error}")
        else:
            return NodeArray(items=items, tag=tag or NewTag(), path=path)

    def _parse_value(self, path: str, anchor_line: int) -> NodeScalar:
        if self.depth > MAX_DEPTH:
            raise Exception(f"max depth: {MAX_DEPTH}")

        tag = self._consume_comments_for(anchor_line)
        tok = self.next()

        if tag is None:
            tag = NewTag()

        if tok.type == TOKEN_STRING:
            text = str(tok.literal)
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Str

            if tag.is_null and not text:
                data = ""
            elif not text:
                data = text
            else:
                data, text = self._validate_string_value(text, tag, path)
            return NodeScalar(data=data, text=text, tag=tag, path=path)
        elif tok.type == TOKEN_NUMBER:
            text = tok.literal
            if '.' in text:
                if tag.type == ValueType.Unknown:
                    tag.type = ValueType.F64
                f64_val = float(text)
                if tag.is_null and f64_val == 0.0:
                    data = 0.0
                else:
                    r = MmValidator.validate_f64(f64_val, tag)
                    if not r.valid:
                        raise Exception(f"Validation error at {path}: {r.error}")
                    data = r.data if r.data is not None else f64_val
                    text = r.text or text
            else:
                if tag.type == ValueType.Unknown:
                    tag.type = ValueType.I
                i_val = int(text)
                if tag.is_null and i_val == 0:
                    data = 0
                else:
                    r = MmValidator.validate_i(i_val, tag)
                    if not r.valid:
                        raise Exception(f"Validation error at {path}: {r.error}")
                    data = r.data if r.data is not None else i_val
                    text = r.text or text
            return NodeScalar(data=data, text=text, tag=tag, path=path)
        elif tok.type == TOKEN_TRUE:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Bool
            if tag.is_null:
                raise Exception(f"bool must false when bool is null")
            r = MmValidator.validate_bool(True, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return NodeScalar(data=True, text="true", tag=tag, path=path)
        elif tok.type == TOKEN_FALSE:
            if tag.type == ValueType.Unknown:
                tag.type = ValueType.Bool
            if not tag.is_null:
                r = MmValidator.validate_bool(False, tag)
                if not r.valid:
                    raise Exception(f"Validation error at {path}: {r.error}")
            return NodeScalar(data=False, text="false", tag=tag, path=path)
        elif tok.type == TOKEN_NULL:
            if tag.type != ValueType.Unknown:
                raise Exception(f"null is not supported for type {tag.type}")
            return NodeNull(tag=tag, path=path)
        else:
            return NodeScalar(data=None, text="", tag=tag, path=path)

    def _validate_string_value(self, text: str, tag: Tag, path: str):
        """Dispatch string token validation based on tag.type, matching Go parser behavior."""
        from datetime import datetime, timezone, timedelta

        if tag.type == ValueType.Str:
            r = MmValidator.validate_str(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Datetime:
            tz = timezone.utc
            if tag.location is not None:
                try:
                    tz = timezone(timedelta(hours=int(tag.location)))
                except (ValueError, TypeError):
                    pass
            try:
                d = datetime.strptime(text, "%Y-%m-%d %H:%M:%S")
            except ValueError:
                raise Exception(f"Validation error at {path}: invalid datetime {repr(text)}")
            d = d.replace(tzinfo=tz)
            r = MmValidator.validate_datetime(d, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else d, r.text or text

        elif tag.type == ValueType.Date:
            tz = timezone.utc
            if tag.location is not None:
                try:
                    tz = timezone(timedelta(hours=int(tag.location)))
                except (ValueError, TypeError):
                    pass
            try:
                d = datetime.strptime(text, "%Y-%m-%d").date()
            except ValueError:
                raise Exception(f"Validation error at {path}: invalid date {repr(text)}")
            r = MmValidator.validate_date(d, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else d, r.text or text

        elif tag.type == ValueType.Time:
            tz = timezone.utc
            if tag.location is not None:
                try:
                    tz = timezone(timedelta(hours=int(tag.location)))
                except (ValueError, TypeError):
                    pass
            try:
                t = datetime.strptime(text, "%H:%M:%S").time()
            except ValueError:
                raise Exception(f"Validation error at {path}: invalid time {repr(text)}")
            # Wrap time in datetime for validation
            d = datetime.combine(datetime.now().date(), t, tzinfo=tz)
            r = MmValidator.validate_time(d, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else t, r.text or text

        elif tag.type == ValueType.Uuid:
            r = MmValidator.validate_uuid(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Ip:
            r = MmValidator.validate_ip(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Url:
            r = MmValidator.validate_url(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Email:
            r = MmValidator.validate_email(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Enums:
            r = MmValidator.validate_enum(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text

        elif tag.type == ValueType.Bytes:
            import base64
            try:
                val = base64.b64decode(text)
            except Exception:
                raise Exception(f"Validation error at {path}: invalid base64 bytes {repr(text)}")
            r = MmValidator.validate_bytes(val, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else val, r.text or text

        elif tag.type == ValueType.Media:
            import base64
            try:
                val = base64.b64decode(text)
            except Exception:
                raise Exception(f"Validation error at {path}: invalid base64 media {repr(text)}")
            r = MmValidator.validate_image(val, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else val, r.text or text

        else:
            # Default to str validation for unknown types
            r = MmValidator.validate_str(text, tag)
            if not r.valid:
                raise Exception(f"Validation error at {path}: {r.error}")
            return r.data if r.data is not None else text, r.text or text


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
    ValueType.Str,
    ValueType.I,
    ValueType.F64,
    ValueType.Bool,
}


def _get_tag_str(tag) -> str:
    """Get tag string, omitting inferred types."""
    if tag is None:
        return ""

    inh = tag.is_inherit

    # Build tag string manually to filter inferred types
    parts = []

    if tag.type != ValueType.Unknown and tag.type not in _INFERRED_TYPES and not inh:
        if not (tag.type == ValueType.Enums and tag.enums) and not (tag.type == ValueType.Media and tag.mime):
            parts.append(f"type={str(tag.type)}")

    if tag.example:
        parts.append("example")
    if tag.is_null:
        parts.append("is_null")

    nullable_key = "child_nullable" if inh else "nullable"
    if tag.nullable and not tag.is_null and not inh:
        parts.append(nullable_key)

    desc_key = "child_desc" if inh else "desc"
    if tag.desc and not inh:
        parts.append(f'{desc_key}="{tag.desc}"')

    if tag.deprecated and not inh:
        parts.append("deprecated")

    allow_empty_key = "child_allow_empty" if inh else "allow_empty"
    if tag.allow_empty and not inh:
        parts.append(allow_empty_key)

    unique_key = "child_unique" if inh else "unique"
    if tag.unique and not inh:
        parts.append(unique_key)

    default_key = "child_default_val" if inh else "default_val"
    if tag.default_val and not inh:
        parts.append(f"{default_key}={tag.default_val}")

    min_key = "child_min" if inh else "min"
    if tag.min and not inh:
        parts.append(f"{min_key}={tag.min}")

    max_key = "child_max" if inh else "max"
    if tag.max and not inh:
        parts.append(f"{max_key}={tag.max}")

    size_key = "child_size" if inh else "size"
    if tag.size and not inh:
        parts.append(f"{size_key}={tag.size}")

    enums_key = "child_enums" if inh else "enums"
    if tag.enums and not inh:
        parts.append(f"{enums_key}={tag.enums}")

    pattern_key = "child_pattern" if inh else "pattern"
    if tag.pattern and not inh:
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
    if tag.version != 0 and not inh:
        parts.append(f"{version_key}={tag.version}")

    mime_key = "child_mime" if inh else "mime"
    if tag.mime and not inh:
        parts.append(f"{mime_key}={tag.mime}")

    if tag.child_desc and not inh:
        parts.append(f'child_desc="{tag.child_desc}"')
    if tag.child_type != ValueType.Unknown and tag.child_type not in _INFERRED_TYPES and not inh:
        if not (tag.child_type == ValueType.Enums and tag.child_enums):
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
    if isinstance(n, NodeNull):
        b.append("null")
    elif isinstance(n, NodeScalar):
        write_value_jsonc(b, n)
    elif isinstance(n, NodeObject):
        write_object_jsonc(b, n, indent)
    elif isinstance(n, NodeArray):
        write_array_jsonc(b, n, indent)


def write_object_jsonc(b: list, o: NodeObject, indent: int):
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


def write_array_jsonc(b: list, a: NodeArray, indent: int):
    b.append("[\n")
    for item in a.items:
        item_tag = item.get_tag()
        # Only write the item's own comments, not inherited parent child_type
        write_leading_comments(b, item_tag, indent + 1)
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