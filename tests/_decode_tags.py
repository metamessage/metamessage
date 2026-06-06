#!/usr/bin/env python3
"""Detailed wire format decoder to understand tag structure"""
import subprocess
import os
import sys

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = 'fixtures/03_tags/child_tags.jsonc'

def get_hex(lang):
    cmds = {
        'Go': f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null',
        'Py': f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null',
        'PHP': f'cd {TESTS_DIR} && php harness/php/harness.php --encode {FIXTURE} 2>/dev/null',
        'TS': f'cd {TESTS_DIR} && node harness/ts/harness.mjs --encode {FIXTURE} 2>/dev/null',
    }
    result = subprocess.run(cmds[lang], shell=True, capture_output=True, text=True)
    return bytes.fromhex(result.stdout.strip())

# Get all hex
go = get_hex('Go')
py = get_hex('Py')
php = get_hex('PHP')
ts = get_hex('TS')

print(f"Go : {len(go)} bytes")
print(f"Py : {len(py)} bytes")
print(f"PHP: {len(php)} bytes")
print(f"TS : {len(ts)} bytes")

# Check if PHP or TS matches Go
for lang, data in [('PHP', php), ('TS', ts)]:
    match = sum(1 for i in range(min(len(go), len(data))) if go[i] == data[i])
    print(f"\nGo vs {lang}: {match}/{len(go)} match")

# Show the first value after key array (the "names" array)
# Go key array: starts at offset 3 (de), length at offset 4 (f3=243), content at offset 5
# So values start at offset 3 + 2 + 243 = 248
go_val_start = 3 + 2 + go[4]  # offset 3: de, offset 4: length
print(f"\nGo values start at offset {go_val_start}")
print(f"Go first value (names): {go[go_val_start:go_val_start+80].hex()}")

# Py key array: starts at offset 3 (de), length at offset 4 (c4=196)
py_val_start = 3 + 2 + py[4]
print(f"\nPy values start at offset {py_val_start}")
print(f"Py first value (names): {py[py_val_start:py_val_start+60].hex()}")

# PHP key array: starts at offset 3 (de), length at offset 4
php_val_start = 3 + 2 + php[4]
print(f"\nPHP values start at offset {php_val_start}")
print(f"PHP first value (names): {php[php_val_start:php_val_start+80].hex()}")

# TS key array: starts at offset 3 (de), length at offset 4
ts_val_start = 3 + 2 + ts[4]
print(f"\nTS values start at offset {ts_val_start}")
print(f"TS first value (names): {ts[ts_val_start:ts_val_start+80].hex()}")

# Manually parse Go's "names" value to understand its structure
# Tag wrapper: fe 32 (TagLen1Byte, length=50)
# Tag content: 0b 8e 09 69 74 65 6d 20 6e 61 6d 65 (child_desc "item name")
# Array: de 24 (len=36)
print("\n\n=== Manual decode of Go 'names' value ===")
d = go[go_val_start:]
print(f"Tag sign: 0x{d[0]:02x}")
print(f"  PrefixTag: {(d[0] & 0xE0) == 0xE0}")
print(f"  Total length (tag+payload): {d[0] & 0x1F}")
tag_sign = d[0]
total_len = d[0] & 0x1F
if total_len >= 30:  # TagLen1Byte
    extra = 1 if total_len == 30 else 2
    total_len = int.from_bytes(d[1:1+extra], 'big')
    tag_start = 1 + extra
else:
    tag_start = 1

print(f"  Actual total length: {total_len}")
print(f"  Tag starts at offset {tag_start}")

# Parse tag content: [length][content]
tag_len_byte = d[tag_start]
print(f"  Tag length byte: 0x{tag_len_byte:02x}")
tag_content_len = tag_len_byte
print(f"  Tag content length: {tag_content_len}")
tag_content = d[tag_start+1:tag_start+1+tag_content_len]
print(f"  Tag content: {tag_content.hex()}")
# Try to decode tag key-value pairs
pos = 0
while pos < len(tag_content):
    key_byte = tag_content[pos]
    key_type = key_byte & 0xF8
    key_len = key_byte & 0x07
    key_names = {
        0x18: 'KDesc', 0x88: 'KChildDesc', 0x20: 'KType',
        0x90: 'KChildType', 0x50: 'KNullable', 0xD0: 'KChildNullable',
        0x60: 'KAllowEmpty', 0xE0: 'KChildAllowEmpty'
    }
    kn = key_names.get(key_type, f'0x{key_type:02x}')
    if key_len <= 5:
        val = tag_content[pos+1:pos+1+key_len]
        print(f"    [{kn}] len={key_len} val={val.hex()} = '{val.decode('utf-8', errors='replace')}'")
        pos += 1 + key_len
    elif key_len == 6:
        val_len = tag_content[pos+1]
        val = tag_content[pos+2:pos+2+val_len]
        print(f"    [{kn}] ext_len len={val_len} val={val.hex()} = '{val.decode('utf-8', errors='replace')}'")
        pos += 2 + val_len
    else:
        print(f"    [{kn}] unknown format, key=0x{key_byte:02x}")
        pos += 1

# Now parse array content
array_start = tag_start + 1 + tag_content_len  # 1 for tag_len_byte
array_sign = d[array_start]
print(f"\n  Array sign: 0x{array_sign:02x}")
array_len_type = array_sign & 0x0F
if array_len_type < 14:
    array_len = array_len_type
    arr_content_start = array_start + 1
elif array_len_type == 14:  # ContainerLen1Byte
    array_len = d[array_start + 1]
    arr_content_start = array_start + 2
else:
    array_len = (d[array_start + 1] << 8) | d[array_start + 2]
    arr_content_start = array_start + 3
print(f"  Array length: {array_len}")
print(f"  Array content start offset in this block: {arr_content_start - array_start}")

# Parse each array element
arr_content = d[arr_content_start:arr_content_start + array_len]
print(f"\n  Array content ({len(arr_content)} bytes): {arr_content.hex()}")
arr_pos = 0
elem_idx = 0
while arr_pos < len(arr_content):
    elem = arr_content[arr_pos]
    if (elem & 0xE0) == 0xE0:  # PrefixTag
        elem_total = elem & 0x1F  # assume < 30
        tag_len_byte = arr_content[arr_pos + 1]
        tag_c_len = tag_len_byte
        tag_c = arr_content[arr_pos+2:arr_pos+2+tag_c_len]
        pay_start = arr_pos + 2 + tag_c_len
        elem_total_bytes = 1 + 1 + tag_c_len + (arr_content[pay_start] & 0x1F) + 1  # approx
        print(f"\n    [{elem_idx}] Tagged element (total={elem_total}):")
        print(f"      Tag length byte: {tag_len_byte}")
        print(f"      Tag content ({tag_c_len} bytes): {tag_c.hex()}")
        # Parse tag content
        p = 0
        while p < len(tag_c):
            kb = tag_c[p]
            kt = kb & 0xF8
            kl = kb & 0x07
            kn = {0x18: 'KDesc', 0x88: 'KChildDesc', 0x20: 'KType'}.get(kt, f'0x{kt:02x}')
            if kl <= 5:
                v = tag_c[p+1:p+1+kl]
                print(f"        [{kn}] val={v.hex()}='{v.decode('utf-8', errors='replace')}'")
                p += 1 + kl
            elif kl == 6:
                vl = tag_c[p+1]
                v = tag_c[p+2:p+2+vl]
                print(f"        [{kn}] ext_len={vl} val={v.hex()}='{v.decode('utf-8', errors='replace')}'")
                p += 2 + vl
            else:
                p += 1
        # Print payload
        payload = arr_content[pay_start:]
        print(f"      Payload: {payload.hex()}")
        elem_idx += 1
        # Calculate actual element size
        # tag wrapper: 1 (sign) + tag_len_byte(1) + tag_content... 
        # but we need to know where this element ends
        arr_pos = arr_content_start + 1  # oversimplified, just break
        break
    else:
        print(f"\n    [{elem_idx}] Untagged element, first byte: 0x{elem:02x}")
        # Try to guess element size
        arr_pos += 1
        elem_idx += 1
        break