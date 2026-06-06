#!/usr/bin/env python3
"""Debug tool: encode child_tags.jsonc with Go and Python, then compare the ASTs."""

import subprocess
import sys
import os

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SCRIPT_DIR = os.path.join(PROJECT_DIR, "tests")
FIXTURE = os.path.join(PROJECT_DIR, "tests/fixtures/03_tags/child_tags.jsonc")

def get_go_encode():
    result = subprocess.run(
        ["go", "run", os.path.join(SCRIPT_DIR, "harness/go/harness.go"), "--encode", FIXTURE],
        capture_output=True, text=True, cwd=PROJECT_DIR
    )
    return result.stdout.strip()

def get_py_encode():
    result = subprocess.run(
        ["python3", os.path.join(SCRIPT_DIR, "harness/python/harness.py"), "--encode", FIXTURE],
        capture_output=True, text=True, cwd=PROJECT_DIR
    )
    return result.stdout.strip()

go_hex = get_go_encode()
py_hex = get_py_encode()

go_bytes = bytes.fromhex(go_hex)
py_bytes = bytes.fromhex(py_hex)

print(f"Go encoded length: {len(go_bytes)} bytes")
print(f"Py encoded length: {len(py_bytes)} bytes")
print(f"Go hex: {go_hex[:200]}...")
print(f"Py hex: {py_hex[:200]}...")
print()

# Find where they first differ
min_len = min(len(go_bytes), len(py_bytes))
for i in range(min_len):
    if go_bytes[i] != py_bytes[i]:
        print(f"First difference at byte offset {i}")
        context_start = max(0, i - 10)
        context_end = min(len(go_bytes), i + 20)
        print(f"  Go[{context_start}:{context_end}]: {go_bytes[context_start:context_end].hex()}")
        print(f"  Py[{context_start}:{context_end}]: {py_bytes[context_start:context_end].hex()}")
        print(f"  Go at {i}: 0x{go_bytes[i]:02x}")
        print(f"  Py at {i}: 0x{py_bytes[i]:02x}")
        break

print()

# Now try to decode both using Go decoder
# First, let's use a simpler approach: extract the raw hex of key array+values
# and compare structure

# Count total bytes for Go and Python
print(f"Total bytes: Go={len(go_bytes)}, Py={len(py_bytes)}")
print(f"Difference: {abs(len(go_bytes) - len(py_bytes))} bytes")
print()

# Now let's use Python's decoder to try decoding the Go bytes
# Add the mm-py path
sys.path.insert(0, os.path.join(PROJECT_DIR, "mm-py"))
sys.path.insert(0, os.path.join(PROJECT_DIR, "mm-py/metamessage"))

from metamessage.core.decoder import Decoder
from metamessage.ir.ast import Obj, Arr, Val, Field

from metamessage.ir.tag import Tag, ValueType

def print_node(node, indent=0):
    prefix = "  " * indent
    if isinstance(node, Obj):
        print(f"{prefix}Object with {len(node.fields)} fields, tag={node.get_tag()}")
        for f in node.fields:
            print(f"{prefix}  Field '{f.key}':")
            print_node(f.value, indent + 2)
    elif isinstance(node, Arr):
        print(f"{prefix}Array with {len(node.items)} items, tag={node.get_tag()}")
        for item in node.items:
            print_node(item, indent + 2)
    elif isinstance(node, Val):
        tag = node.get_tag()
        print(f"{prefix}Val: data={node.data}, text='{node.text}', tag={tag}")

print("=== Decoding Go bytes with Python decoder ===")
go_decoder = Decoder(go_bytes)
go_node = go_decoder.decode_node()
print_node(go_node)

print()
print("=== Decoding Py bytes with Python decoder ===")
py_decoder = Decoder(py_bytes)
try:
    py_node = py_decoder.decode_node()
    print_node(py_node)
except Exception as e:
    print(f"Error decoding Python bytes: {e}")
    print(f"Python hex first 200: {py_hex[:200]}")
    print(f"Python hex: {py_hex}")