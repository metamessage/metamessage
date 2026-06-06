#!/usr/bin/env python3
"""Debug why Python can't decode its own hex"""
import subprocess
import os

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = 'fixtures/03_tags/child_tags.jsonc'

# Get hex from different languages
py_hex = subprocess.run(
    f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

go_hex = subprocess.run(
    f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

# Try to decode Python hex with Go's decoder
with open('/tmp/py_hex.txt', 'w') as f:
    f.write(py_hex)

go_decode_py = subprocess.run(
    f'cd {TESTS_DIR} && go run harness/go/harness.go --decode < /tmp/py_hex.txt 2>&1',
    shell=True, capture_output=True, text=True
)
print(f"Go decode Python hex: stdout={repr(go_decode_py.stdout[:500])}")
print(f"Go decode Python hex: stderr={repr(go_decode_py.stderr[:500])}")
print()

# Compare Python vs Go hex side by side
import binascii
go_bytes = bytes.fromhex(go_hex)
py_bytes = bytes.fromhex(py_hex)

# Find which fields are different by parsing at a high level
# Let me manually compare region by region
# Top-level container (3 bytes) then tag wrapper or key array

# The Go encoding structure for the top level (no mm: tag):
# cf 03 f6 = container prefix (len=0x3f6=1014)
# then key array: starts with dc or similar (array prefix), then keys
# then values follow

# Python structure:
# cf 02 c4 = container prefix (len=0x2c4=708)
# then key array...

# Skip container header
go_pos = 3
py_pos = 3

# Next should be the tag wrapper for the top-level object
# OR the key array (if no tag)
# Since there's no mm: tag on the top level, there should be no tag wrapper.

# Go bytes at pos 3: de f3 ... (PrefixTag!)
# Python bytes at pos 3: de f3 ... (PrefixTag!)

# Wait, BOTH have a tag wrapper! That means the top-level has a tag.
# de = PrefixTag | TagLen1Byte = 0xDE
# f3 = 243 bytes total (tag + payload)

# The tag itself comes first within these 243 bytes.
# Let's examine the tag structure
print("=== Go: Tag content at start of payload ===")
# Tag starts with: [tag_len_byte][tag_content]
go_tag_len_byte = go_bytes[go_pos + 2]  # skip de f3
go_tag_content_len = go_tag_len_byte
go_tag_content = go_bytes[go_pos+3:go_pos+3+go_tag_content_len]
print(f"Tag length byte: {go_tag_len_byte}")
print(f"Tag content ({go_tag_content_len} bytes): {go_tag_content.hex()}")
print()

print("=== Py: Tag content at start of payload ===")
py_tag_len_byte = py_bytes[py_pos + 2]
py_tag_content_len = py_tag_len_byte
py_tag_content = py_bytes[py_pos+3:py_pos+3+py_tag_content_len]
print(f"Tag length byte: {py_tag_len_byte}")
print(f"Tag content ({py_tag_content_len} bytes): {py_tag_content.hex()}")
print()

# After tag content, the payload starts
# Payload should be: key_array + values
go_payload_start = go_pos + 3 + go_tag_content_len
py_payload_start = py_pos + 3 + py_tag_content_len

go_payload_len = int.from_bytes(go_bytes[1:3], 'big') - (go_payload_start - 3)
py_payload_len = int.from_bytes(py_bytes[1:3], 'big') - (py_payload_start - 3)

print(f"Go payload: start={go_payload_start}, remaining={go_payload_len}")
print(f"Py payload: start={py_payload_start}, remaining={py_payload_len}")

# Parse key array
# Key array is a container: [prefix][length][contents]
print()
print("=== Key array ===")
go_key_prefix = go_bytes[go_payload_start]
py_key_prefix = py_bytes[py_payload_start]
print(f"Go key prefix: {go_key_prefix:02x} {'Array' if go_key_prefix & 0x10 else 'Object'}")
print(f"Py key prefix: {py_key_prefix:02x} {'Array' if py_key_prefix & 0x10 else 'Object'}")

# Parse key array length
go_key_len_type = go_key_prefix & 0x0F
py_key_len_type = py_key_prefix & 0x0F

if go_key_len_type < 14:
    go_key_len = go_key_len_type
    go_key_start = go_payload_start + 1
else:
    go_key_len = go_bytes[go_payload_start+1]
    go_key_start = go_payload_start + 2

if py_key_len_type < 14:
    py_key_len = py_key_len_type
    py_key_start = py_payload_start + 1
else:
    py_key_len = py_bytes[py_payload_start+1]
    py_key_start = py_payload_start + 2

print(f"Go key array: length={go_key_len}, start={go_key_start}")
print(f"Py key array: length={py_key_len}, start={py_key_start}")

# After key array, values begin
go_val_start = go_key_start + go_key_len
py_val_start = py_key_start + py_key_len
print(f"Go values start: {go_val_start}")
print(f"Py values start: {py_val_start}")
print(f"Go remaining after keys: {len(go_bytes) - go_val_start}")
print(f"Py remaining after keys: {len(py_bytes) - py_val_start}")

# Now try to decode with Python's decoder
# The issue is Python can't decode its own bytes
# Let me try using Python's decoder on the Python bytes
print()
print("=== Attempting Python decode of Python bytes ===")

# Python's decoder decode function
sys.path.insert(0, os.path.join(TESTS_DIR, '..'))
from mm_py.metamessage.core.decoder import decode
from mm_py.metamessage.ir.tag import NewTag

try:
    result, length = decode(py_bytes)
    print(f"Python decode own bytes: SUCCESS, length={length}")
    print(f"Result type: {type(result)}")
except Exception as e:
    print(f"Python decode own bytes: FAILED - {e}")
    import traceback
    traceback.print_exc()

# Try to decode with small steps
print()
print("=== Step-by-step decode of Python first 10 bytes ===")
for i in range(10):
    print(f"  byte {i}: {py_bytes[i]:02x} ({py_bytes[i]:08b}) prefix={py_bytes[i] & 0xE0:02x}")

# Let me also try decoding byte by byte with the Python decoder
# to find where it fails