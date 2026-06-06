#!/usr/bin/env python3
"""Analyze byte-level differences between Go and other languages for child_tags.jsonc encode."""

import sys
import re

go_hex = "cf03f6def3856e616d65738673636f72657385636f756e74896f6c645f6669656c6485695f6172728669385f617272876931365f617272876933325f617272876936345f61727285755f6172728675385f617272877531365f617272877533325f617272877536345f617272876633325f617272876636345f617272877374725f61727288626f6f6c5f6172728962797465735f6172728a626967696e745f6172728c6461746574696d655f61727288646174655f6172728874696d655f61727288757569645f6172728b646563696d616c5f6172728669705f6172728775726c5f61727289656d61696c5f61727288656e756d5f617272fe320b8e096974656d206e616d65de24f20b8e096974656d206e616d6585616c696365f00b8e096974656d206e616d6583626f62fe2b0799b931c3313030de21ea0799b931c33130303850ea0799b931c3313030385aea0799b931c3313030385fe603423432382ae90111866c6567616379d3212223f402900ade0fe402900a21e402900a22e402900a23f402900bde0fe402900b21e402900b22e402900b23f402900cde0fe402900c21e402900c22e402900c23f402900dde0fe402900d21e402900d22e402900d23f402900ede0fe402900e21e402900e22e402900e23f402900fde0fe402900f21e402900f22e402900f23f4029010de0fe402901021e402901022e402901023f4029011de0fe402901121e402901122e402901123f4029012de0fe402901221e402901222e402901223fa029013de15e602901368ff0fe602901368ff19e602901368ff23d968ff0f68ff1968ff23d6816181628163d3060506f9029007de14e9029007a568656c6c6fe9029007a5776f726c64fe23029015de1eee029015aa140f6e462a062b3535a0ee029015aa147b747282315fad80a0f7029016de12e80290163b65920080e80290163b666d8948f3029017de0ee6029017394d0be6029017394db1f4029018de0fe602901839afc8e70290183a01517ffe2f029019de2af4029019b0550e8400e29b41d4a716446655440000f4029019b06ba7b8109dad11d180b400c04fd430c8f702901ade12e802901a6afb04cb2fe802901a6afb0425d4fe2202901bde1def02901b8b3139322e3136382e312e31ec02901b8831302e302e302e31fe3202901cde2df702901c9368747470733a2f2f6578616d706c652e636f6df402901c9068747470733a2f2f746573742e6f7267fe2d02901dde28f402901d9075736572406578616d706c652e636f6df202901d8e61646d696e40746573742e6f7267fe4c10d60e7265647c677265656e7c626c7565de39f210d60e7265647c677265656e7c626c756520f210d60e7265647c677265656e7c626c756521f210d60e7265647c677265656e7c626c756522f210d60e7265647c677265656e7c626c756523f210d60e7265647c677265656e7c626c756524f210d60e7265647c677265656e7c626c756525f210d60e7265647c677265656e7c626c756526f210d60e7265647c677265656e7c626c756527"

# After the tag section, the rest is the encoded object bytes.
# Let me look at known byte patterns.
# "names" = 6e616d6573
# "scores" = 73636f726573
# "count" = 636f756e74

# Let's compare Go vs Rust at key positions
rust_hex = "cf0323def3856e616d65738673636f72657385636f756e74896f6c645f6669656c6485695f6172728669385f617272876931365f617272876933325f617272876936345f61727285755f6172728675385f617272877531365f617272877533325f617272877536345f617272876633325f617272876636345f617272877374725f61727288626f6f6c5f6172728962797465735f6172728a626967696e745f6172728c6461746574696d655f61727288646174655f6172728874696d655f61727288757569645f6172728b646563696d616c5f6172728669705f6172728775726c5f61727289656d61696c5f61727288656e756d5f617272f70b8e096974656d206e616d65da85616c69636583626f62ef0799b931c3313030d63850385a385fe603423432382ae90111866c6567616379d3212223e702900ad3212223e702900bd3212223e702900cd3212223e702900dd3212223e702900ed3212223e702900fd3212223e7029010d3212223e7029011d3212223e7029012d3212223ed029013d968ff0f68ff1968ff23d968ff0f68ff1968ff23d6816181628163d3060506f7029007de1288614756736247383d88643239796247513dfe23029015de1e3fab54a98ceb1f0ad2943938373635343332313039383736353433323130fe2d029016de2893323032342d30312d30312030303a30303a303093323032342d30362d31352031323a33303a3030fb029017de168a323032342d30312d30318a323032342d30362d3135f7029018de128831323a33303a30308832333a35393a3539fe51029019de4c9e2435353065383430302d653239622d343164342d613731362d3434363635353434303030309e2436626137623831302d396461642d313164312d383062342d303063303466643433306338ee02901ada6afb04cb2f6afb0425d4fa02901bde158b3139322e3136382e312e318831302e302e302e31fe2a02901cde259368747470733a2f2f6578616d706c652e636f6d9068747470733a2f2f746573742e6f7267fe2502901dde209075736572406578616d706c652e636f6d8e61646d696e40746573742e6f7267f510d60e7265647c677265656e7c626c7565d3202122"

def hex_to_bytes(h):
    return bytes.fromhex(h)

def print_byte_comparison(label_a, hex_a, label_b, hex_b, start=0, length=200):
    ba = hex_to_bytes(hex_a)
    bb = hex_to_bytes(hex_b)
    print(f"=== Byte comparison: {label_a} vs {label_b} ===")
    print(f"Total length: {label_a}={len(ba)}, {label_b}={len(bb)}")
    
    for i in range(start, min(start + length, min(len(ba), len(bb)))):
        if ba[i] != bb[i]:
            print(f"  Byte {i}: {label_a}=0x{ba[i]:02x} ({ba[i]}), {label_b}=0x{bb[i]:02x} ({bb[i]}) {'<-- DIFF' if i < start + length else ''}")
    
    if len(ba) != len(bb):
        print(f"\nLength mismatch: {len(ba)} vs {len(bb)} (delta: {len(ba) - len(bb)})")
        if len(ba) > len(bb):
            diff_start = min(len(ba), len(bb))
            print(f"  Extra bytes in {label_a} starting at {diff_start}:")
            remaining = ba[diff_start:]
            print(f"  {remaining.hex()}")
        else:
            diff_start = min(len(ba), len(bb))
            print(f"  Extra bytes in {label_b} starting at {diff_start}:")
            remaining = bb[diff_start:]
            print(f"  {remaining.hex()}")

print("=" * 70)
print("ANALYSIS OF CHILD_TAGS ENCODE DIFF")
print("=" * 70)

# Compare Go vs Rust
print_byte_comparison("Go", go_hex, "Rust", rust_hex, 0, 500)

# Let's also look at the Python version
py_hex = "cf02c8def3856e616d65738673636f72657385636f756e74896f6c645f6669656c6485695f6172728669385f617272876931365f617272876933325f617272876936345f61727285755f6172728675385f617272877531365f617272877533325f617272877536345f617272876633325f617272876636345f617272877374725f61727288626f6f6c5f6172728962797465735f6172728a626967696e745f6172728c6461746574696d655f61727288646174655f6172728874696d655f61727288757569645f6172728b646563696d616c5f6172728669705f6172728775726c5f61727289656d61696c5f61727288656e756d5f617272f70b8e096974656d206e616d65da85616c69636583626f62ef0799b931c3313030d63850385a385fe603423432382ae90111866c6567616379d3212223e702900ad3212223e702900bd3212223e702900cd3212223e702900dd3212223e702900ed3212223e702900fd3212223e7029010d3212223e7029011d3212223e7029012d3212223ed029013d968ff0f68ff1968ff23d968ff0f68ff1968ff23d6816181628163d3060506f7029007de12a8614756736247383da8643239796247513dfe2f029015de2a140f6e462a062b3535a0aa140f6e462a062b3535a0147b747282315fad80a0aa147b747282315fad80a0ee029016da3b659200803b666d8948ea029017d6394d0b394db1eb029018d739afc83a01517ffe27029019de22b0550e8400e29b41d4a716446655440000b06ba7b8109dad11d180b400c04fd430c8ee02901ada6afb04cb2f6afb0425d4fa02901bde158b3139322e3136382e312e318831302e302e302e31fe2a02901cde259368747470733a2f2f6578616d706c652e636f6d9068747470733a2f2f746573742e6f7267fe2502901dde209075736572406578616d706c652e636f6d8e61646d696e40746573742e6f7267f510d60e7265647c677265656e7c626c7565d3202122"

# Compare PHP (which matches Go) vs Rust
php_hex = "cf03f6def3856e616d65738673636f72657385636f756e74896f6c645f6669656c6485695f6172728669385f617272876931365f617272876933325f617272876936345f61727285755f6172728675385f617272877531365f617272877533325f617272877536345f617272876633325f617272876636345f617272877374725f61727288626f6f6c5f6172728962797465735f6172728a626967696e745f6172728c6461746574696d655f61727288646174655f6172728874696d655f61727288757569645f6172728b646563696d616c5f6172728669705f6172728775726c5f61727289656d61696c5f61727288656e756d5f617272fe320b8e096974656d206e616d65de24f20b8e096974656d206e616d6585616c696365f00b8e096974656d206e616d6583626f62fe2b0799b931c3313030de21ea0799b931c33130303850ea0799b931c3313030385aea0799b931c3313030385fe603423432382ae90111866c6567616379d3212223f402900ade0fe402900a21e402900a22e402900a23f402900bde0fe402900b21e402900b22e402900b23f402900cde0fe402900c21e402900c22e402900c23f402900dde0fe402900d21e402900d22e402900d23f402900ede0fe402900e21e402900e22e402900e23f402900fde0fe402900f21e402900f22e402900f23f4029010de0fe402901021e402901022e402901023f4029011de0fe402901121e402901122e402901123f4029012de0fe402901221e402901222e402901223fa029013de15e602901368ff0fe602901368ff19e602901368ff23d968ff0f68ff1968ff23d6816181628163d3060506f9029007de14e9029007a568656c6c6fe9029007a5776f726c64fe23029015de1eee029015aa140f6e462a062b3535a0ee029015aa147b747282315fad80a0f7029016de12e80290163b65920080e80290163b666d8948f3029017de0ee6029017394d0be6029017394db1f4029018de0fe602901839afc8e70290183a01517ffe2f029019de2af4029019b0550e8400e29b41d4a716446655440000f4029019b06ba7b8109dad11d180b400c04fd430c8f702901ade12e802901a6afb04cb2fe802901a6afb0425d4fe2202901bde1def02901b8b3139322e3136382e312e31ec02901b8831302e302e302e31fe3202901cde2df702901c9368747470733a2f2f6578616d706c652e636f6df402901c9068747470733a2f2f746573742e6f7267fe2d02901dde28f402901d9075736572406578616d706c652e636f6df202901d8e61646d696e40746573742e6f7267fe4c10d60e7265647c677265656e7c626c7565de39f210d60e7265647c677265656e7c626c756520f210d60e7265647c677265656e7c626c756521f210d60e7265647c677265656e7c626c756522f210d60e7265647c677265656e7c626c756523f210d60e7265647c677265656e7c626c756524f210d60e7265647c677265656e7c626c756525f210d60e7265647c677265656e7c626c756526f210d60e7265647c677265656e7c626c756527"

print("\n\n")
print("=" * 70)
print("PHP vs Go (should be identical)")
print("=" * 70)
print_byte_comparison("Go", go_hex, "PHP", php_hex, 0, 50)

# Let's understand the wire format structure
print("\n\n")
print("=" * 70)
print("WIRE FORMAT STRUCTURE ANALYSIS")
print("=" * 70)

ba = hex_to_bytes(go_hex)
print(f"\nGo bytes (first 100):")
for i in range(0, min(100, len(ba))):
    print(f"  [{i:3d}] 0x{ba[i]:02x} ({ba[i]:3d}) '{chr(ba[i]) if 32 <= ba[i] < 127 else '.'}'", end="")
    if i > 0 and (i+1) % 4 == 0:
        print()
    else:
        print(" ", end="")
print()

# Let's identify key sections
print("\n\nKey patterns in Go hex:")
patterns = {
    "names": b"names",
    "scores": b"scores",
    "count": b"count",
    "old_field": b"old_field",
    "i_arr": b"i_arr",
    "i8_arr": b"i8_arr",
    "i16_arr": b"i16_arr",
    "i32_arr": b"i32_arr",
    "i64_arr": b"i64_arr",
    "u_arr": b"u_arr",
    "u8_arr": b"u8_arr",
    "u16_arr": b"u16_arr",
    "u32_arr": b"u32_arr",
    "u64_arr": b"u64_arr",
    "f32_arr": b"f32_arr",
    "f64_arr": b"f64_arr",
    "str_arr": b"str_arr",
    "bool_arr": b"bool_arr",
    "bytes_arr": b"bytes_arr",
    "bigint_arr": b"bigint_arr",
    "datetime_arr": b"datetime_arr",
    "date_arr": b"date_arr",
    "time_arr": b"time_arr",
    "uuid_arr": b"uuid_arr",
    "decimal_arr": b"decimal_arr",
    "ip_arr": b"ip_arr",
    "url_arr": b"url_arr",
    "email_arr": b"email_arr",
    "enum_arr": b"enum_arr",
    "alice": b"alice",
    "bob": b"bob",
    "item name": b"item name",
    "legacy": b"legacy",
}

for name, pattern in sorted(patterns.items(), key=lambda x: ba.find(x[1])):
    pos = ba.find(pattern)
    if pos >= 0:
        context_before = ba[max(0,pos-3):pos]
        print(f"  '{name}' at byte {pos}, prefix: {context_before.hex()}")

# Now let's understand the tag format
# Tag byte structure: 
# TagKey (upper 5 bits) | length (lower 3 bits) or extended length

print("\n\nTag analysis:")
known_tag_keys = {
    0: "KIsNull",
    8: "KExample",
    16: "KDeprecated",
    24: "KDesc",
    32: "KType",
    40: "KNullable",
    48: "KAllowEmpty",
    56: "KUnique",
    64: "KDefault",
    72: "KMin",
    80: "KMax",
    88: "KSize",
    96: "KEnum",
    104: "KPattern",
    112: "KLocation",
    120: "KVersion",
    128: "KMime",
    136: "KChildDesc",
    144: "KChildType",
    152: "KChildNullable",
    160: "KChildAllowEmpty",
    168: "KChildUnique",
    176: "KChildDefaultVal",
    184: "KChildMin",
    192: "KChildMax",
    200: "KChildSize",
    208: "KChildEnums",
    216: "KChildPattern",
    224: "KChildLocation",
    232: "KChildVersion",
    240: "KChildMime",
    248: "KMore",
}

# Parse the tag section
# First, let's find the structure
# cf 03 = tag prefix + 2-byte length
# f6 de f3 ... = tag data + payload

# After cf 03, the next 2 bytes are the tag+payload length (big-endian)
# 0xf6de = 63198 (very large, this must be tag + payload length)

# Actually wait, let me re-read the encode_tag logic
# sign = PrefixTag (0xE0)
# if length < Max1Byte (0xFF): sign |= TagLen1Byte (0x1E)... 
# Actually TagLen1Byte = 0b11111 - 1 = 30 = 0x1E
# TagLen2Byte = 0b11111 = 31 = 0x1F
# So byte 0: 0xE0 | 0x1F = 0xFF for 2-byte length

# But byte 0 is 0xcf = 0b1100 1111 which is Container, not Tag!

# Wait, let me re-read encodeComment:
# encodeT encodes the tag bytes with length, THEN encodeTag wraps payload+tag
# So the outermost structure is:
# [TagPrefix + tagLen][tag_bytes][payload_data]
# But for the root object tag, it might be different...

# Actually, let me re-read. If the tag is empty (no bytes), encodeT returns 0 and encodeComment returns 0.
# Then the payload is just the object bytes.

# But if there are tags, encodeComment wraps.

# Hmm, let me think about this differently. 

# Let me look at the first few bytes as a Container:
# 0xcf = 0b1100 1111
# Prefix = 0b110 = Container
# ContainerArray = 0x10
# ContainerLen: 0x0F = 15 → not a standard length

# Wait, ContainerLenMask = 0b01111 = 15
# ContainerLen1Byte = 15 - 1 = 14 = 0x0E
# ContainerLen2Byte = 15 = 0x0F

# 0xcf & 0x0F = 0x0F = ContainerLen2Byte → 2 more bytes for length

# Wait no, that's & ContainerLenMask (0b01111). 
# 0xcf & 0x0F = 0x0F. No wait, 0xcf = 0b1100 1111, 0x0F = 0b0000 1111, so 0xcf & 0x0F = 0x0F.

# ContainerLen:
# For ContainerLen2Byte (15): 2 more bytes follow for the length

# OK so the first byte 0xcf means:
# Prefix = Container
# It's an array (ContainerArray bit set)
# Length uses 2-byte encoding

# Bytes 1-2: length (big-endian)
# 0x03 0xf6 = 0x03f6 = 1014

# Then bytes 3 onward: the content

# So the structure is: Container(array) header + 2-byte length(0x03f6=1014) + content

# That makes more sense! The root object is a key array + values.

# Content at byte 3:
# 0xde = TagPrefix | TagLen2Byte = 0xDE = 0b1101 1110
# Wait, TagPrefix = 0b111 << 5 = 0xE0
# 0xDE & 0xE0 = 0xC0, which is Container
# Hmm, that's not right either.

# Wait, let me reconsider. PrefixTag = 0b111 << 5 = 0xE0
# TagLen2Byte = 0b11111 = 31 = 0x1F
# So TagPrefix | TagLen2Byte = 0xFF
# Hmm but the byte is 0xcf which is 0b1100 1111

# Let me check Container again:
# Container = 0b110 << 5 = 0xC0
# ContainerArray = 0b10000 = 0x10
# ContainerObject = 0b00000 = 0x00
# ContainerLen2Byte = 0b01111 = 0x0F

# Container | ContainerArray | ContainerLen2Byte = 0xC0 | 0x10 | 0x0F = 0xCF = 0b1100 1111 ✓

# So byte 0 is the root object: it's a Container of type Array (key array), with 2-byte length
# Bytes 1-2: length = 0x03f6 = 1014 bytes of content

# Inside the content at byte 3:
# 0xde = 0b1101 1110
# Top 3 bits: 0b110 = Container
# ContainerArray: 0x10
# ContainerLen: 0x0E = ContainerLen1Byte (14)
# Wait, 0xDE = 0b1101 1110
# & 0xE0 = 0xC0 = Container ✓
# & 0x10 = 0x10 = ContainerArray ✓  
# & 0x0F = 0x0E = 14 = ContainerLen1Byte → 1 more byte for length

# Hmm, but this is inside the content. Let me re-read the encode flow:
# encodeNodeObject:
# 1. encode keys array → encodeArray(bufKey.Bytes())
# 2. encode combined (keyArray + values) → encodeObject(bufAll)
# 3. wrap with tag → encodeComment(payload, tag)

# encodeComment:
# - encodeT(tag.Bytes()) → if tag has content, encode with length prefix
# - encodeTag(payload, tagBytes) → wrap tag + payload with TagPrefix

# So the outermost wrapping is encodeTag which starts with TagPrefix byte.
# But byte 0 is Container, not TagPrefix.

# That means the root object has NO tag (empty tag). So encodeComment returns n1=0, and the outer n stays as the raw object bytes.

# So the bytes are:
# Byte 0: Container(Array) | 0x0F → 2-byte length
# Bytes 1-2: length (0x03f6)
# Content starts at byte 3

# 0xde at byte 3:
# 0xDE = Container | ContainerArray | ContainerLen1Byte 
# = 0xC0 | 0x10 | 0x0E = 0xDE ✓
# Length byte at position 4: 0xf3 = 243 bytes for this inner array

# Wait, but the outer is an array of keys, and then values follow...

# Actually, looking at the structure more carefully:
# encodeNodeObject does:
# 1. Loop over fields, for each:
#    a. Encode the value (object/array/value)
#    b. Encode the key (string)
# 2. After loop:
#    a. encodeArray(bufKey.Bytes()) → wraps all keys in an array
#    b. encodeObject(combined) → wraps [keyArray + allValues] in an object
#    c. encodeComment(payload, tag) → if tag exists, wrap with tag

# For the root object with no tag:
# - The keys array is encoded
# - Then keyArray + values are combined into an object
# - Since no tag, just the object bytes are returned

# So the bytes are:
# [Object header][keyArray content + values content]

# Object header at byte 0: 0xcf
# = Container | ContainerObject | ContainerLen2Byte
# Wait I said ContainerArray before. Let me re-check.
# ContainerObject = 0b00000 = 0x00
# ContainerArray = 0b10000 = 0x10
# 0xCF & 0x10 = 0x10 = ContainerArray!
# So the outer container is an ARRAY, not an object?

# Wait, that's because encodeNodeObject calls:
# encodeArray(bufKey.Bytes()) for the keys, then combines with values.
# The result is stored, then encodeObject wraps it.
# Then encodeComment wraps the whole thing.

# But I said no tag for root, so encodeComment returns 0.

# Let me re-read encodeNodeObject more carefully:
# n = e.encodeArray(bufKey.Bytes()) → n is bytes written for key array
# encodedKeyArray = e.buf[e.offset-n : e.offset]
# bufAll = encodedKeyArray + buf.Bytes() (all values)
# n = e.encodeObject(bufAll) → n is bytes written for object

# Then:
# n1 = e.encodeComment(..., tag)
# if n1 == 0: return n

# So the final n is the object bytes (not the comment).

# The object header (from encodeObject):
# sign = Container | ContainerObject
# If length < ContainerLen1Byte: sign |= length
# etc.

# ContainerObject = 0b00000 = 0x00
# 0xcf & 0x10 = 0x10 → ContainerArray!

# Hmm, that's wrong. Let me re-check.
# 0xcf = 0b1100 1111
# 0x10 = 0b0001 0000
# 0xcf & 0x10 = 0b0000 0000 = 0x00
# So it's NOT ContainerArray!

# Oh wait, 0xCF = 0b1100 1111
# & 0xEF = ... no.
# ContainerMask = 0b10000 = 16 = 0x10
# ContainerArray = 0b10000 = 0x10
# ContainerObject = 0b00000 = 0x00

# 0xCF = 0b1100 1111
# 0xCF & 0x10 = 0b0000 0000 = 0x00
# So it's ContainerObject. Good.

# OK so the structure is:
# Byte 0: Container | ContainerObject | ContainerLen2Byte (0x0F) = 0xCF
# Bytes 1-2: big-endian length of object payload
# Then the payload...

# Inside the payload at byte 3:
# 0xDE = Container | ContainerArray | ContainerLen1Byte = 0xC0 | 0x10 | 0x0E = 0xDE
# This is the keys array!

# OK, this is getting complex. Let me take a different approach.
# I'll write a Python script that dumps both Go and Rust wire format
# and highlights the differences.

print("\n\nDetailed byte diff (Go vs Rust):")
bg = hex_to_bytes(go_hex)
br = hex_to_bytes(rust_hex)
max_len = min(len(bg), len(br))
diffs = []
for i in range(max_len):
    if bg[i] != br[i]:
        diffs.append((i, bg[i], br[i]))

for idx, (pos, gb, rb) in enumerate(diffs):
    print(f"  Diff {idx}: byte {pos}: Go=0x{gb:02x} Rust=0x{rb:02x}")
    # Show context
    start = max(0, pos - 5)
    end = min(len(bg), pos + 10)
    print(f"    Context: Go:   {bg[start:end].hex()}")
    print(f"    Context: Rust: {br[start:end].hex()}")
    
    # Try to interpret these bytes
    if start <= 50 and end <= 150:
        pass  # Already in the tag section
    print()