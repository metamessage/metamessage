#!/usr/bin/env python3
import subprocess
import os

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = 'fixtures/03_tags/child_tags.jsonc'

def get_hex(lang, cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    raw = result.stdout.strip()
    return bytes.fromhex(raw) if raw else None

go = get_hex('Go', f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null')
py = get_hex('Py', f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null')

# Key array is same for both: de f3 + 243 bytes
# Root container: cf + 2-byte len
# After root container header (3 bytes), key array starts: de f3 (2 bytes) + 243 bytes
# So value area starts at offset 3 + 2 + 243 = 248

key_array_len = 0xf3
value_start = 5 + key_array_len  # 3 for root header + 2 for key array header + 243

go_values = go[value_start:]
py_values = py[value_start:]

print(f"Go values ({len(go_values)} bytes):")
print(f"  {go_values.hex()}")
print()
print(f"Py values ({len(py_values)} bytes):")
print(f"  {py_values.hex()}")

# Let's find the first difference in values
print()
min_len = min(len(go_values), len(py_values))
for i in range(min_len):
    if go_values[i] != py_values[i]:
        ctx_start = max(0, i - 4)
        ctx_end = min(len(go_values), i + 40)
        print(f"First difference at offset {i} (absolute {value_start + i}):")
        print(f"  Go: {go_values[ctx_start:ctx_end].hex()}")
        print(f"  Py: {py_values[ctx_start:ctx_end].hex()}")
        print(f"  Go byte: 0x{go_values[i]:02x} vs Py byte: 0x{py_values[i]:02x}")
        break

# The values should correspond to fields in order:
# names, scores, count, old_field, i_arr, i8_arr, ... enum_arr
# Let's look for field boundaries
print()
print("=== Field structure analysis ===")

# names (child_desc="item name") - should have tag for each item
# Let's look at field 2 (count=42, simple int)
# count is the 3rd field

# Let me look at i8_arr which has child_type=i8
i8_hex = "69385f617272"
go_i8_off = go.hex().index(i8_hex)
py_i8_off = py.hex().index(i8_hex)
print(f"\n--- i8_arr field ---")
print(f"Go i8_arr at absolute offset: {go_i8_off}")
print(f"Py i8_arr at absolute offset: {py_i8_off}")

# Look at the values after i8_arr key
print()
# Find the value by looking for the encoded value bytes
# i8_arr has values [1, 2, 3]
# With inherited tag child_type=i8, each value should be encoded as i8 (single byte int)
# Without tag, values would be encoded as regular ints

# Let me print more detail
# Count field is at hex offset ~58 (just counting)