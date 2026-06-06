#!/usr/bin/env python3
"""Decode and compare wire format bytes between languages."""

import sys
import re

def hex_to_bytes(h):
    return bytes.fromhex(h.strip())

# Extract hex from encode_diff file
with open('/Users/lizongying/IdeaProjects/meta-message/tests/results/03_tags_child_tags.jsonc.encode_diff') as f:
    content = f.read()

# Parse hex blocks for each language
langs = {}
current_lang = None
hex_lines = []

for line in content.split('\n'):
    m = re.match(r'--- (\w+) \(hex\) ---', line)
    if m:
        if current_lang and hex_lines:
            langs[current_lang] = ''.join(hex_lines).replace(' ', '').replace('\n', '')
        current_lang = m.group(1)
        hex_lines = []
    elif current_lang and line.strip() and not line.startswith('---') and not line.startswith('==')):
        hex_lines.append(line.strip())
    
if current_lang and hex_lines:
    langs[current_lang] = ''.join(hex_lines).replace(' ', '').replace('\n', '')

# Decode wire format
PREFIX_TAG = 0xE0
PREFIX_CONTAINER = 0xC0
PREFIX_STRING = 0x80
PREFIX_BYTES = 0xA0
PREFIX_POSITIVE_INT = 0x20
PREFIX_NEGATIVE_INT = 0x40
PREFIX_FLOAT = 0x60
PREFIX_SIMPLE = 0x00

CONTAINER_ARRAY = 0x10
CONTAINER_MAP = 0x00

TAG_LEN_1 = 30
TAG_LEN_2 = 31
CONTAINER_LEN_1 = 14
CONTAINER_LEN_2 = 15

def get_prefix(b):
    return b & 0xE0

def decode_container(data, offset=0, indent=0):
    """Decode a container."""
    prefix = '  ' * indent
    if offset >= len(data):
        return offset
    
    b = data[offset]
    p = get_prefix(b)
    
    if p != PREFIX_CONTAINER:
        print(f"{prefix}[{offset}] Not a container: 0x{b:02x}")
        return offset
    
    is_array = b & CONTAINER_ARRAY
    container_type = "Array" if is_array else "Object"
    
    len_code = b & 0x0F
    if len_code < CONTAINER_LEN_1:
        length = len_code
        extra = 0
    elif len_code == CONTAINER_LEN_1:
        length = data[offset + 1]
        extra = 1
    else:
        length = (data[offset + 1] << 8) | data[offset + 2]
        extra = 2
    
    print(f"{prefix}[{offset}] Container({container_type}) len={length}")
    
    payload_start = offset + 1 + extra
    payload_end = payload_start + length
    
    # Parse payload contents
    pos = payload_start
    while pos < payload_end:
        pos = decode_node(data, pos, indent + 1)
    
    return payload_end

def decode_tag(b):
    """Decode a tag key byte."""
    tag_keys = {
        0: "IsNull", 8: "Example", 16: "Deprecated",
        24: "Desc", 32: "Type", 40: "Nullable",
        48: "AllowEmpty", 56: "Unique",
        64: "Default", 72: "Min", 80: "Max", 88: "Size",
        96: "Enum", 104: "Pattern", 112: "Location", 120: "Version", 128: "Mime",
        136: "ChildDesc", 144: "ChildType", 152: "ChildNullable",
        160: "ChildAllowEmpty", 168: "ChildUnique",
        176: "ChildDefaultVal", 184: "ChildMin", 192: "ChildMax", 200: "ChildSize",
        208: "ChildEnums", 216: "ChildPattern", 224: "ChildLocation", 232: "ChildVersion", 240: "ChildMime",
        248: "More"
    }
    key = b & 0xF8  # Upper 5 bits
    len_val = b & 0x07  # Lower 3 bits for inline length
    name = tag_keys.get(key, f"Unknown(0x{key:02x})")
    return key, len_val, name

def decode_tag_content(data, offset, indent=0):
    """Decode tag content (from Tag.Bytes())."""
    prefix = '  ' * indent
    pos = offset
    
    while pos < len(data):
        b = data[pos]
        key, len_val, name = decode_tag(b)
        
        # Boolean tags (len_val = 1 means true)
        if name in ("IsNull", "Example", "Nullable", "Deprecated", "AllowEmpty", "Unique",
                    "ChildNullable", "ChildAllowEmpty", "ChildUnique"):
            if len_val == 1:
                print(f"{prefix}[{pos}] Tag: {name} = true")
            else:
                print(f"{prefix}[{pos}] Tag: {name} with len={len_val}")
            pos += 1
            continue
        
        # Type tags (1 byte value)
        if name in ("Type", "ChildType"):
            val = data[pos + 1]
            val_names = {0: "Unknown", 1: "Doc", 2: "Vec", 3: "Arr", 4: "Obj", 5: "Map",
                        6: "Str", 7: "Bytes", 8: "Bool", 9: "I", 10: "I8", 11: "I16",
                        12: "I32", 13: "I64", 14: "U", 15: "U8", 16: "U16", 17: "U32",
                        18: "U64", 19: "F32", 20: "F64", 21: "Bigint", 22: "Datetime",
                        23: "Date", 24: "Time", 25: "Uuid", 26: "Decimal", 27: "Ip",
                        28: "Url", 29: "Email", 30: "Enum", 31: "Media"}
            val_name = val_names.get(val, f"ValueType({val})")
            print(f"{prefix}[{pos}] Tag: {name} = {val} ({val_name})")
            pos += 2
            continue
        
        # Size/Version/Mime tags (encodeU64 format)
        if name in ("Size", "Version", "Mime", "ChildSize", "ChildVersion", "ChildMime", "More"):
            int_len = (b & 0x07) + 1
            if int_len <= 0 or int_len > 8:
                int_len = 1
            val = 0
            for i in range(int_len):
                val = (val << 8) | data[pos + 1 + i]
            print(f"{prefix}[{pos}] Tag: {name} = {val}")
            pos += 1 + int_len
            continue
        
        # String tags with inline length (len_val < 6) or extended length (len_val == 6 or 7)
        if len_val <= 5:
            # Inline length
            string_len = len_val
            text = data[pos+1:pos+1+string_len].decode('utf-8', errors='replace')
            print(f"{prefix}[{pos}] Tag: {name} = \"{text}\"")
            pos += 1 + string_len
        elif len_val == 6:
            # 1-byte length
            string_len = data[pos + 1]
            text = data[pos+2:pos+2+string_len].decode('utf-8', errors='replace')
            print(f"{prefix}[{pos}] Tag: {name} = \"{text}\"")
            pos += 2 + string_len
        elif len_val == 7:
            # 2-byte length
            string_len = (data[pos + 1] << 8) | data[pos + 2]
            text = data[pos+3:pos+3+string_len].decode('utf-8', errors='replace')
            print(f"{prefix}[{pos}] Tag: {name} = \"{text}\"")
            pos += 3 + string_len
        else:
            print(f"{prefix}[{pos}] Unknown tag byte: 0x{b:02x}")
            pos += 1
            break
    
    return pos

def decode_node(data, offset=0, indent=0):
    """Decode a node (value, array, or object)."""
    prefix = '  ' * indent
    if offset >= len(data):
        return offset
    
    b = data[offset]
    p = get_prefix(b)
    
    if p == PREFIX_TAG:
        # Tag wrapping
        len_code = b & 0x1F
        if len_code < TAG_LEN_1:
            total_len = len_code
            tag_data_start = offset + 1
            extra = 0
        elif len_code == TAG_LEN_1:
            total_len = data[offset + 1]
            extra = 1
            tag_data_start = offset + 2
        else:
            total_len = (data[offset + 1] << 8) | data[offset + 2]
            extra = 2
            tag_data_start = offset + 3
        
        # Read tag content
        # First byte of tag data is the tag length (from encodeT)
        tag_len = data[tag_data_start]
        tag_content_start = tag_data_start + 1
        tag_content_end = tag_content_start + tag_len
        
        print(f"{prefix}[{offset}] TagWrapper(total_len={total_len}, tag_len={tag_len})")
        
        # Decode tag content
        decode_tag_content(data, tag_content_start, indent + 1)
        
        # Payload starts after tag content
        payload_start = tag_content_end
        print(f"{prefix}[{offset}] Payload starts at {payload_start}")
        
        # The payload is a container or value
        if payload_start < len(data):
            payload_end = decode_node(data, payload_start, indent + 1)
            return payload_end
        
        return tag_content_end
    
    elif p == PREFIX_CONTAINER:
        return decode_container(data, offset, indent)
    
    elif p == PREFIX_STRING:
        # String value
        len_code = b & 0x1F
        if len_code < 30:
            string_len = len_code
            extra = 0
        elif len_code == 30:
            string_len = data[offset + 1]
            extra = 1
        else:
            string_len = (data[offset + 1] << 8) | data[offset + 2]
            extra = 2
        
        text = data[offset+1+extra:offset+1+extra+string_len].decode('utf-8', errors='replace') if string_len > 0 else ""
        print(f"{prefix}[{offset}] String({string_len}): \"{text}\"")
        return offset + 1 + extra + string_len
    
    elif p == PREFIX_BYTES:
        len_code = b & 0x1F
        if len_code < 30:
            bytes_len = len_code
            extra = 0
        elif len_code == 30:
            bytes_len = data[offset + 1]
            extra = 1
        else:
            bytes_len = (data[offset + 1] << 8) | data[offset + 2]
            extra = 2
        print(f"{prefix}[{offset}] Bytes({bytes_len})")
        return offset + 1 + extra + bytes_len
    
    elif p == PREFIX_POSITIVE_INT:
        # Positive integer
        int_val = b & 0x1F
        if int_val < 24:
            print(f"{prefix}[{offset}] PositiveInt(inline): {int_val}")
            return offset + 1
        int_len = int_val - 24 + 1
        val = 0
        for i in range(int_len):
            val = (val << 8) | data[offset + 1 + i]
        print(f"{prefix}[{offset}] PositiveInt({int_len} bytes): {val}")
        return offset + 1 + int_len
    
    elif p == PREFIX_NEGATIVE_INT:
        int_val = b & 0x1F
        if int_val < 24:
            print(f"{prefix}[{offset}] NegativeInt(inline): -{int_val}")
            return offset + 1
        int_len = int_val - 24 + 1
        val = 0
        for i in range(int_len):
            val = (val << 8) | data[offset + 1 + i]
        print(f"{prefix}[{offset}] NegativeInt({int_len} bytes): -{val}")
        return offset + 1 + int_len
    
    elif p == PREFIX_FLOAT:
        mantissa_len_code = b & 0x0F
        is_neg = b & 0x10
        neg_str = "negative " if is_neg else ""
        if mantissa_len_code <= 7:
            print(f"{prefix}[{offset}] Float(inline): {neg_str}exp=-1 mantissa={mantissa_len_code}")
            return offset + 1
        mantissa_len = mantissa_len_code - 8 + 1
        exponent = data[offset + 1]
        mantissa = 0
        for i in range(mantissa_len):
            mantissa = (mantissa << 8) | data[offset + 2 + i]
        print(f"{prefix}[{offset}] Float({mantissa_len} bytes): {neg_str}exp={exponent} mantissa={mantissa}")
        return offset + 2 + mantissa_len
    
    elif p == PREFIX_SIMPLE:
        simple_val = b & 0x1F
        simple_names = {0: "NullBool", 1: "NullInt", 2: "NullFloat", 3: "NullString", 4: "NullBytes",
                        5: "False", 6: "True", 7: "Code", 8: "Message", 9: "Data",
                        10: "Success", 11: "Error", 12: "Unknown(Pageable)"}
        name = simple_names.get(simple_val, f"Simple({simple_val})")
        print(f"{prefix}[{offset}] Simple: {name}")
        return offset + 1
    
    else:
        print(f"{prefix}[{offset}] Unknown prefix: 0x{b:02x}")
        return offset + 1

# Parse and compare Go vs Python
print("=" * 70)
print("DECODING GO WIRE FORMAT")
print("=" * 70)
go_bytes = hex_to_bytes(langs.get('go', ''))
decode_node(go_bytes)

print("\n\n")
print("=" * 70)
print("DECODING PY WIRE FORMAT")
print("=" * 70)
py_bytes = hex_to_bytes(langs.get('py', ''))
decode_node(py_bytes)

# Find where the "names" field starts (key string)
print("\n\n")
print("=" * 70)
print("NAMES FIELD COMPARISON")
print("=" * 70)

# Find "names" in the Go bytes
names_pos = go_bytes.find(b'names')
if names_pos >= 0:
    print(f"Go 'names' at byte {names_pos}")
    # Show context
    start = max(0, names_pos - 10)
    end = min(len(go_bytes), names_pos + 80)
    print(f"  Context: {go_bytes[start:end].hex()}")

names_pos_py = py_bytes.find(b'names')
if names_pos_py >= 0:
    print(f"Py 'names' at byte {names_pos_py}")
    start = max(0, names_pos_py - 10)
    end = min(len(py_bytes), names_pos_py + 80)
    print(f"  Context: {py_bytes[start:end].hex()}")