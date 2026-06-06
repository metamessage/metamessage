#!/usr/bin/env python3
"""Compare Go vs Python encoding by decoding both and comparing ASTs."""
import sys
import os
import subprocess

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'mm-py'))
from metamessage import parse_jsonc, to_jsonc
from metamessage.core.decoder import Decoder
from metamessage.core.encoder import Encoder

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = os.path.join(SCRIPT_DIR, "fixtures", "03_tags", "child_tags.jsonc")

with open(FIXTURE, 'r') as f:
    source = f.read()

# Get Go's hex
result = subprocess.run(
    ["go", "run", os.path.join(SCRIPT_DIR, "harness/go/harness.go"), "--encode", FIXTURE],
    capture_output=True, text=True, timeout=30
)
go_hex = result.stdout.strip()
go_bytes = bytes.fromhex(go_hex)
print(f"Go encoded: {len(go_bytes)} bytes")

# Get Python's hex
result = subprocess.run(
    ["python3", os.path.join(SCRIPT_DIR, "harness/python/harness.py"), "--encode", FIXTURE],
    capture_output=True, text=True, timeout=30
)
py_hex = result.stdout.strip()
py_bytes = bytes.fromhex(py_hex)
print(f"Python encoded: {len(py_bytes)} bytes")

# Decode Go's output with Python's decoder
go_decoder = Decoder(go_bytes)
go_node = go_decoder.decode_node()
print(f"\n=== Go-encoded decoded by Python ===")
print(f"Node type: {type(go_node).__name__}")
if hasattr(go_node, 'fields'):
    print(f"Number of fields: {len(go_node.fields)}")
    for f in go_node.fields:
        tag = f.value.get_tag()
        tp = tag.type if tag else 'N/A'
        ct = tag.child_type if tag else 'N/A'
        cd = tag.child_desc if tag else 'N/A'
        print(f"  {f.key}: type={tp}, child_type={ct}, child_desc={cd!r}")
        if hasattr(f.value, 'items'):
            for i, item in enumerate(f.value.items):
                item_tag = item.get_tag()
                dt = item.data if hasattr(item, 'data') else '?'
                print(f"    [{i}]: type={item_tag.type}, is_inherit={item_tag.is_inherit}, data={dt}")

# Decode Python's output with Python's decoder
py_decoder = Decoder(py_bytes)
py_node = py_decoder.decode_node()
print(f"\n=== Python-encoded decoded by Python ===")
print(f"Node type: {type(py_node).__name__}")
if hasattr(py_node, 'fields'):
    print(f"Number of fields: {len(py_node.fields)}")
    for f in py_node.fields:
        tag = f.value.get_tag()
        tp = tag.type if tag else 'N/A'
        ct = tag.child_type if tag else 'N/A'
        cd = tag.child_desc if tag else 'N/A'
        print(f"  {f.key}: type={tp}, child_type={ct}, child_desc={cd!r}")
        if hasattr(f.value, 'items'):
            for i, item in enumerate(f.value.items):
                item_tag = item.get_tag()
                dt = item.data if hasattr(item, 'data') else '?'
                print(f"    [{i}]: type={item_tag.type}, is_inherit={item_tag.is_inherit}, data={dt}")

# Compare lengths
print(f"\n=== Comparison ===")
print(f"Go length: {len(go_bytes)}, Python length: {len(py_bytes)}, Diff: {abs(len(go_bytes) - len(py_bytes))}")

# Find first byte difference
min_len = min(len(go_bytes), len(py_bytes))
for i in range(min_len):
    if go_bytes[i] != py_bytes[i]:
        ctx_start = max(0, i - 4)
        ctx_end = min(len(go_bytes), i + 40)
        print(f"First byte diff at offset {i}:")
        print(f"  Go: {go_bytes[ctx_start:ctx_end].hex()}")
        print(f"  Py: {py_bytes[ctx_start:ctx_end].hex()}")
        break