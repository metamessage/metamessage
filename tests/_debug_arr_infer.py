#!/usr/bin/env python3
"""Debug array type inference in Python decoder."""
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'mm-py'))
from metamessage import encode_from_jsonc, decode_to_jsonc
from metamessage.core.decoder import Decoder

# Test case 1: array with explicit size=3 tag (the problem case)
tagged_json = '// mm: size=3\n[1, 2, 3]'
print('Test 1: Tagged array with size=3')
print('Input:', repr(tagged_json))

wire = encode_from_jsonc(tagged_json)
print('Wire hex:', wire.hex())

d = Decoder(wire)
node = d.decode_node()
print('Node type:', type(node).__name__)
print('Tag type:', node.tag.type, '(value:', int(node.tag.type), ')')
print('Tag size:', node.tag.size)
print('Tag is_inherit:', node.tag.is_inherit)
print('Tag str:', str(node.tag))

output = decode_to_jsonc(wire)
print('Decoded output:', repr(output))
print()

# Test case 2: tagged array field in object
tagged_obj = '// mm: size=3\n{"items": [1, 2, 3]}'
print('Test 2: Tagged array field in object')
print('Input:', repr(tagged_obj))

wire2 = encode_from_jsonc(tagged_obj)
print('Wire hex:', wire2.hex())

d2 = Decoder(wire2)
node2 = d2.decode_node()
print('Node type:', type(node2).__name__)
if hasattr(node2, 'fields'):
    for f in node2.fields:
        print(f'  Field "{f.key}":')
        print(f'    Value type: {type(f.value).__name__}')
        print(f'    Tag type: {f.value.tag.type} (value: {int(f.value.tag.type)})')
        print(f'    Tag size: {f.value.tag.size}')
        print(f'    Tag is_inherit: {f.value.tag.is_inherit}')
        print(f'    Tag str: "{str(f.value.tag)}"')

output2 = decode_to_jsonc(wire2)
print('Decoded output:', repr(output2))
print()

# Test case 3: no tag array (should be vec)
plain_json = '[1, 2, 3]'
print('Test 3: Untagged array')
wire3 = encode_from_jsonc(plain_json)
d3 = Decoder(wire3)
node3 = d3.decode_node()
print('Tag type:', node3.tag.type, '(value:', int(node3.tag.type), ')')
print('Tag str:', str(node3.tag))
print('Decoded output:', repr(decode_to_jsonc(wire3)))