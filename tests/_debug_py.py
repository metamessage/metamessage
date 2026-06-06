#!/usr/bin/env python3
"""Debug Python's encoding of child_tags.jsonc."""
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'mm-py'))
from metamessage import parse_jsonc, to_jsonc, encode_from_jsonc

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = os.path.join(SCRIPT_DIR, "fixtures", "03_tags", "child_tags.jsonc")

with open(FIXTURE, 'r') as f:
    source = f.read()

node = parse_jsonc(source)

# Check fields
if hasattr(node, 'fields'):
    print(f"Number of fields: {len(node.fields)}")
    for f in node.fields:
        val = f.value
        tag = val.get_tag()
        tp = tag.type if tag else 'N/A'
        inh = tag.is_inherit if tag else 'N/A'
        cd = tag.child_desc if tag else 'N/A'
        ct = tag.child_type if tag else 'N/A'
        ci = tag.is_inherit if tag else 'N/A'
        print(f"  {f.key}: type={tp}, child_desc={cd!r}, child_type={ct}, is_inherit={ci}")
        if hasattr(val, 'items'):
            for i, item in enumerate(val.items):
                item_tag = item.get_tag()
                print(f"    [{i}]: type={item_tag.type}, is_inherit={item_tag.is_inherit}, data={item.data if hasattr(item,'data') else 'N/A'}")
elif hasattr(node, 'items'):
    print(f"Node is array with {len(node.items)} items")
elif hasattr(node, 'data'):
    print(f"Node is value: {node.data}")

# Try encoding
wire = encode_from_jsonc(source)
print(f"\nEncoded wire length: {len(wire)}")
print(f"Encoded hex: {wire.hex()}")

# Also try Go-style encoding by parsing first
from metamessage.core.encoder import Encoder
encoder = Encoder()
wire2 = encoder.encode(node)
print(f"\nEncoder.encode(node) length: {len(wire2)}")
print(f"Encoded hex: {wire2.hex()}")
print(f"Match: {wire == wire2}")