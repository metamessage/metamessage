"""Debug script to check Python tag decoding for constraints_tag"""
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'mm-py'))
from metamessage import encode_from_jsonc, decode_to_jsonc
from metamessage.core.decoder import Decoder
from metamessage.ir.tag import Tag, ValueType, TagKey

# Test: JSONC round-trip
print("=== JSONC round-trip ===")
jsonc_content = '''{
	// mm: min=1; max=100
	"score": 85,

	// mm: type=arr; size=3
	"items": [
		1,
		2,
		3,
	],

	// mm: min=1; max=64
	"username": "john_doe",
}'''

encoded = encode_from_jsonc(jsonc_content)
print(f"Encoded hex: {encoded.hex()}")

# Decode and inspect the items tag
decoder = Decoder(encoded)
node = decoder.decode_node()
for field in node.fields:
    if field.key == 'items':
        arr = field.value
        print(f"\nItems tag:")
        print(f"  type: {arr.tag.type}")
        print(f"  type value: {arr.tag.type.value}")
        print(f"  size: {arr.tag.size}")
        print(f"  is_inherit: {arr.tag.is_inherit}")
        print(f"  str(): '{arr.tag}'")
        print(f"  type==Arr: {arr.tag.type == ValueType.Arr}")
        print(f"  type in suppression list: {arr.tag.type in (ValueType.Str, ValueType.I, ValueType.F64, ValueType.Bool, ValueType.Obj, ValueType.Vec)}")
        break

# Also show how the __str__ handles it
t = Tag(type=ValueType.Arr, size=3)
print(f"\nTest tag with type=Arr, size=3:")
print(f"  str(): '{t}'")
print(f"  is_inherit: {t.is_inherit}")

t2 = Tag(type=ValueType.Arr, size=3, is_inherit=True)
print(f"\nTest tag with type=Arr, size=3, is_inherit=True:")
print(f"  str(): '{t2}'")

print(f"\nFull decoded JSONC:")
print(decode_to_jsonc(encoded))