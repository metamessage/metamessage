"""Simple MetaMessage wire format decoder for debugging"""

import struct
import sys

PREFIX_SIMPLE = 0b000 << 5
PREFIX_POSITIVE_INT = 0b001 << 5
PREFIX_NEGATIVE_INT = 0b010 << 5
PREFIX_FLOAT = 0b011 << 5
PREFIX_STRING = 0b100 << 5
PREFIX_BYTES = 0b101 << 5
PREFIX_CONTAINER = 0b110 << 5
PREFIX_TAG = 0b111 << 5

PREFIX_MASK = 0b11100000
SUFFIX_MASK = 0b00011111

# C# ValueType enum (Go uses same values)
VT = {
    0: "Unknown", 1: "Doc", 2: "Vec", 3: "Arr", 4: "Obj", 5: "Map",
    6: "Str", 7: "Bytes", 8: "Bool", 9: "I", 10: "I8", 11: "I16", 12: "I32",
    13: "I64", 14: "U", 15: "U8", 16: "U16", 17: "U32", 18: "U64",
    19: "F32", 20: "F64", 21: "Bigint", 22: "Datetime", 23: "Date", 24: "Time",
    25: "Uuid", 26: "Decimal", 27: "Ip", 28: "Url", 29: "Email",
    30: "Enums", 31: "Media"
}

TAG_KEY = {
    0: "KIsNull", 8: "KExample", 16: "KDeprecated",
    24: "KDesc", 32: "KType", 40: "KNullable", 48: "KAllowEmpty",
    56: "KUnique", 64: "KDefault", 72: "KMin", 80: "KMax",
    88: "KSize", 96: "KEnum", 104: "KPattern", 112: "KLocation",
    120: "KVersion", 128: "KMime",
    136: "KChildDesc", 144: "KChildType", 152: "KChildNullable",
    160: "KChildAllowEmpty", 168: "KChildUnique", 176: "KChildDefaultVal",
    184: "KChildMin", 192: "KChildMax", 200: "KChildSize",
    208: "KChildEnums", 216: "KChildPattern", 224: "KChildLocation",
    232: "KChildVersion", 240: "KChildMime",
    248: "KMore"
}

def prefix_name(p):
    if p == PREFIX_SIMPLE: return "SIMPLE"
    if p == PREFIX_POSITIVE_INT: return "POS_INT"
    if p == PREFIX_NEGATIVE_INT: return "NEG_INT"
    if p == PREFIX_FLOAT: return "FLOAT"
    if p == PREFIX_STRING: return "STRING"
    if p == PREFIX_BYTES: return "BYTES"
    if p == PREFIX_CONTAINER: return "CONTAINER"
    if p == PREFIX_TAG: return "TAG"
    return f"UNKNOWN({p>>5})"

def decode_suffix(first_byte, suffix):
    """Get the raw suffix value"""
    return first_byte & SUFFIX_MASK

def decode_value(data, offset, indent=0):
    """Decode a single wire value and return (value_info, new_offset)"""
    if offset >= len(data):
        return ("END", offset)
    
    pfx = data[offset] & PREFIX_MASK
    suffix = data[offset] & SUFFIX_MASK
    tag_key = data[offset] & 0b11111000  # For tag key decoding
    
    if pfx == PREFIX_SIMPLE:
        val = suffix
        offset += 1
        return (f"SIMPLE({val})", offset)
    
    elif pfx == PREFIX_POSITIVE_INT:
        offset += 1
        if suffix <= 23:  # Inline small int
            return (f"INT({suffix})", offset)
        byte_count = 31 - suffix
        val = 0
        for _ in range(byte_count):
            val = (val << 8) | data[offset]
            offset += 1
        return (f"INT({val})", offset)
    
    elif pfx == PREFIX_NEGATIVE_INT:
        offset += 1
        if suffix <= 23:
            return (f"NEG_INT({-suffix})", offset)
        byte_count = 31 - suffix
        val = 0
        for _ in range(byte_count):
            val = (val << 8) | data[offset]
            offset += 1
        return (f"NEG_INT({-val})", offset)
    
    elif pfx == PREFIX_FLOAT:
        # Simplified: just skip
        offset += 1
        neg = (suffix & 0b10000) != 0
        len_suffix = suffix & 0b01111
        if len_suffix <= 7:  # Inline
            return (f"FLOAT({'-' if neg else '+'}{len_suffix})", offset)
        byte_count = 15 - len_suffix
        val_bytes = data[offset:offset+byte_count]
        offset += byte_count
        return (f"FLOAT({val_bytes.hex()})", offset)
    
    elif pfx == PREFIX_STRING:
        offset += 1
        if suffix <= 29:  # Short string
            length = suffix
        elif suffix == 30:  # 1-byte length
            length = data[offset]
            offset += 1
        else:  # suffix == 31, 2-byte length
            length = struct.unpack_from('>H', data, offset)[0]
            offset += 2
        s = data[offset:offset+length].decode('utf-8', errors='replace')
        offset += length
        return (f"STR({length}): {repr(s)}", offset)
    
    elif pfx == PREFIX_BYTES:
        offset += 1
        if suffix <= 29:
            length = suffix
        elif suffix == 30:
            length = data[offset]
            offset += 1
        else:
            length = struct.unpack_from('>H', data, offset)[0]
            offset += 2
        b = data[offset:offset+length]
        offset += length
        return (f"BYTES({length}): {b.hex()}", offset)
    
    elif pfx == PREFIX_CONTAINER:
        offset += 1
        is_arr = (suffix & 0b10000) != 0
        len_suffix = suffix & 0b01111
        if len_suffix <= 13:
            length = len_suffix
        elif len_suffix == 14:
            length = data[offset]
            offset += 1
        else:  # 15, 2-byte length
            length = struct.unpack_from('>H', data, offset)[0]
            offset += 2
        cont_type = "ARR" if is_arr else "MAP"
        return (f"{cont_type}(len={length})", offset, length)
    
    elif pfx == PREFIX_TAG:
        offset += 1
        if suffix <= 29:
            total_len = suffix
        elif suffix == 30:
            total_len = data[offset]
            offset += 1
        else:
            total_len = struct.unpack_from('>H', data, offset)[0]
            offset += 2
        return (f"TAG(total_len={total_len})", offset, total_len)
    
    return (f"UNKNOWN({data[offset]:02x})", offset+1)


def decode_tag_bytes(data, offset, indent=0):
    """Decode tag inner bytes and return (tag_info, new_offset)"""
    tag_fields = {}
    while offset < len(data):
        pfx = data[offset]
        tag_key_val = pfx & 0b11111000  # Upper 5 bits
        tag_suffix = pfx & 0b00000111  # Lower 3 bits (length or special)
        
        if tag_key_val == 0:  # KIsNull (0 << 3)
            tag_fields["IsNull"] = True
            offset += 1
        elif tag_key_val == 8:  # KExample
            tag_fields["Example"] = True
            offset += 1
        elif tag_key_val == 16:  # KDeprecated
            tag_fields["Deprecated"] = True
            offset += 1
        elif tag_key_val == 24:  # KDesc
            # Parse string (suffix is length if <= 5)
            if tag_suffix <= 5:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["Desc"] = s
                offset += 1 + tag_suffix
            elif tag_suffix == 6:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["Desc"] = s
                offset += 2 + l
            else:
                l = struct.unpack_from('>H', data, offset+1)[0]
                s = data[offset+3:offset+3+l].decode('utf-8', errors='replace')
                tag_fields["Desc"] = s
                offset += 3 + l
        elif tag_key_val == 32:  # KType
            v = data[offset+1]
            tag_fields["Type"] = VT.get(v, f"Unknown({v})")
            offset += 2
        elif tag_key_val == 40:  # KNullable
            tag_fields["Nullable"] = True
            offset += 1
        elif tag_key_val == 48:  # KAllowEmpty
            tag_fields["AllowEmpty"] = True
            offset += 1
        elif tag_key_val == 56:  # KUnique
            tag_fields["Unique"] = True
            offset += 1
        elif tag_key_val == 64:  # KDefault
            if tag_suffix < 7:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["Default"] = s
                offset += 1 + tag_suffix
            else:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["Default"] = s
                offset += 2 + l
        elif tag_key_val == 72:  # KMin
            if tag_suffix < 7:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["Min"] = s
                offset += 1 + tag_suffix
            else:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["Min"] = s
                offset += 2 + l
        elif tag_key_val == 80:  # KMax
            if tag_suffix < 7:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["Max"] = s
                offset += 1 + tag_suffix
            else:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                offset += 2 + l
        elif tag_key_val == 88:  # KSize
            # Variable length uint
            kw = tag_suffix
            v = 0
            for _ in range(kw+1):
                offset += 1
                v = (v << 8) | data[offset]
            offset += 1
            tag_fields["Size"] = v
        elif tag_key_val == 96:  # KEnum
            if tag_suffix <= 5:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["Enum"] = s
                offset += 1 + tag_suffix
            elif tag_suffix == 6:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["Enum"] = s
                offset += 2 + l
            else:
                l = struct.unpack_from('>H', data, offset+1)[0]
                s = data[offset+3:offset+3+l].decode('utf-8', errors='replace')
                tag_fields["Enum"] = s
                offset += 3 + l
        elif tag_key_val == 144:  # KChildType
            v = data[offset+1]
            tag_fields["ChildType"] = VT.get(v, f"Unknown({v})")
            offset += 2
        elif tag_key_val == 136:  # KChildDesc
            if tag_suffix <= 5:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["ChildDesc"] = s
                offset += 1 + tag_suffix
            elif tag_suffix == 6:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["ChildDesc"] = s
                offset += 2 + l
        elif tag_key_val == 208:  # KChildEnums
            if tag_suffix <= 5:
                s = data[offset+1:offset+1+tag_suffix].decode('utf-8', errors='replace')
                tag_fields["ChildEnums"] = s
                offset += 1 + tag_suffix
            elif tag_suffix == 6:
                l = data[offset+1]
                s = data[offset+2:offset+2+l].decode('utf-8', errors='replace')
                tag_fields["ChildEnums"] = s
                offset += 2 + l
        else:
            # Unknown tag key - skip
            offset += 1
            # Try to skip gracefully
            if tag_key_val > 0:
                pass
            break
    
    return tag_fields, offset


def decode_toplevel(data, indent=0):
    """Try to decode and print the top-level structure"""
    offset = 0
    level = 0
    
    def p(msg, *args):
        prefix = "  " * (level + indent)
        print(f"{prefix}{msg}", *args)
    
    while offset < len(data):
        pfx = data[offset] & PREFIX_MASK
        suffix = data[offset] & SUFFIX_MASK
        pname = prefix_name(pfx)
        
        if pfx == PREFIX_TAG:
            # TAG wrapper
            if suffix <= 29:
                total_len = suffix
                tag_inner_len = data[offset+1]
                tag_end = offset + 1 + 1 + tag_inner_len
                # Tag inner bytes
                tag_fields, _ = decode_tag_bytes(data, offset+2, level)
                inner_offset = tag_end  # after tag inner
                remaining = total_len - (1 + tag_inner_len)  # -1 for the first byte of tag inner (length byte)
                # Actually total_len includes: tag_inner_len_byte + tag_inner_bytes + payload
                # More precisely: total_len = 1 + tag_inner_len + payload_len
                payload_len = total_len - 1 - tag_inner_len
                payload_offset = tag_end
                
                p(f"TAG(total={total_len}, tag_inner_len={tag_inner_len}, tag={tag_fields}, payload_len={payload_len}) @ {offset}")
                
                # Decode payload recursively
                if payload_len > 0:
                    level += 1
                    decode_toplevel(data[payload_offset:payload_offset+payload_len], level)
                    level -= 1
                
                offset = offset + 1 + total_len
                
            elif suffix == 30:
                total_len = data[offset+1]
                tag_inner_len = data[offset+2]
                tag_end = offset + 2 + 1 + tag_inner_len
                tag_fields, _ = decode_tag_bytes(data, offset+3, level)
                payload_len = total_len - 2 - tag_inner_len  # -2 for the length byte and inner byte
                payload_offset = tag_end
                
                p(f"TAG_LEN1(total={total_len}, tag={tag_fields}) @ {offset}")
                if payload_len > 0:
                    level += 1
                    decode_toplevel(data[payload_offset:payload_offset+payload_len], level)
                    level -= 1
                offset = offset + 2 + total_len
            else:  # suffix == 31
                total_len = struct.unpack_from('>H', data, offset+1)[0]
                tag_inner_len = data[offset+3]
                tag_end = offset + 3 + 1 + tag_inner_len
                tag_fields, _ = decode_tag_bytes(data, offset+4, level)
                payload_len = total_len - 3 - tag_inner_len
                payload_offset = tag_end
                
                p(f"TAG_LEN2(total={total_len}, tag={tag_fields}) @ {offset}")
                if payload_len > 0:
                    level += 1
                    decode_toplevel(data[payload_offset:payload_offset+payload_len], level)
                    level -= 1
                offset = offset + 3 + total_len
        
        elif pfx == PREFIX_CONTAINER:
            is_arr = (suffix & 0b10000) != 0
            len_suffix = suffix & 0b01111
            if len_suffix <= 13:
                length = len_suffix
                offset += 1
            elif len_suffix == 14:
                length = data[offset+1]
                offset += 2
            else:
                length = struct.unpack_from('>H', data, offset+1)[0]
                offset += 3
            cont_type = "ARR" if is_arr else "MAP"
            p(f"{cont_type}(payload_len={length}) @ orig_offset")
            if length > 0:
                level += 1
                decode_toplevel(data[offset:offset+length], level)
                level -= 1
            offset += length
        
        elif pfx == PREFIX_SIMPLE:
            val = suffix
            p(f"SIMPLE({val}) @ {offset}")
            offset += 1
        
        elif pfx in (PREFIX_POSITIVE_INT, PREFIX_NEGATIVE_INT):
            # Simplified integer decoding
            start = offset
            offset += 1
            if suffix <= 23:
                val = suffix
            else:
                byte_count = 31 - suffix
                val = 0
                for _ in range(byte_count):
                    val = (val << 8) | data[offset]
                    offset += 1
            if pfx == PREFIX_NEGATIVE_INT and suffix <= 23:
                val = -val
            p(f"INT({val}) @ {start}")
        
        elif pfx == PREFIX_STRING:
            start = offset
            offset += 1
            if suffix <= 29:
                length = suffix
            elif suffix == 30:
                length = data[offset]
                offset += 1
            else:
                length = struct.unpack_from('>H', data, offset)[0]
                offset += 2
            s = data[offset:offset+length].decode('utf-8', errors='replace')
            offset += length
            p(f"STR({length}): {repr(s)} @ {start}")
        
        elif pfx == PREFIX_BYTES:
            start = offset
            offset += 1
            if suffix <= 29:
                length = suffix
            elif suffix == 30:
                length = data[offset]
                offset += 1
            else:
                length = struct.unpack_from('>H', data, offset)[0]
                offset += 2
            b = data[offset:offset+length].hex()
            offset += length
            p(f"BYTES({length}): {b} @ {start}")
        
        elif pfx == PREFIX_FLOAT:
            start = offset
            offset += 1
            neg = (suffix & 0b10000) != 0
            len_sfx = suffix & 0b01111
            if len_sfx <= 7:
                p(f"FLOAT(neg={neg},val={len_sfx}) @ {start}")
            else:
                byte_count = 15 - len_sfx
                val_bytes = data[offset:offset+byte_count]
                offset += byte_count
                p(f"FLOAT(neg={neg},{val_bytes.hex()}) @ {start}")
        
        else:
            p(f"UNKNOWN({data[offset]:02x}) @ {offset}")
            offset += 1


# Test with Go output
go_hex = input().strip()
data = bytes.fromhex(go_hex)
print(f"Total length: {len(data)}")
print("=" * 60)
decode_toplevel(data)