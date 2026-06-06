#!/usr/bin/env python3
"""Debug Python encoder by comparing with Go for specific types"""
import subprocess
import os
import sys

# Add mm-py to path
sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), '..'))

from mm_py.metamessage.core.encoder import Encoder
from mm_py.metamessage.ir.tag import Tag, ValueType
from mm_py.metamessage.ir.ast import Val, Obj, Arr, Field

def encode_node(node):
    e = Encoder()
    return e.encode(node).hex()

# Test 1: Encode a simple i8 value
print("=== Test 1: Simple i8 value ===")
tag = Tag(type=ValueType.I8, is_inherit=False)
val = Val(data=1, text="1", tag=tag)
hex_out = encode_node(val)
print(f"Python (type=i8, is_inherit=False): {hex_out}")

# Test 2: Encode with inherited type
tag2 = Tag(type=ValueType.I8, is_inherit=True)
val2 = Val(data=1, text="1", tag=tag2)
hex_out2 = encode_node(val2)
print(f"Python (type=i8, is_inherit=True):  {hex_out2}")

# Test 3: What happens with default int type
tag3 = Tag(type=ValueType.I, is_inherit=False)
val3 = Val(data=1, text="1", tag=tag3)
hex_out3 = encode_node(val3)
print(f"Python (type=i, is_inherit=False):  {hex_out3}")

# Test 4: type=unknown but inheriting to i8  
tag4 = Tag(type=ValueType.I, is_inherit=False)
tag4.inherit(Tag(child_type=ValueType.I8))
val4 = Val(data=1, text="1", tag=tag4)
hex_out4 = encode_node(val4)
print(f"Python (inherit I->I8):              {hex_out4}")
print(f"  tag.type={tag4.type}, is_inherit={tag4.is_inherit}")
print(f"  tag.bytes()={tag4.bytes().hex() if tag4.bytes() else '(empty)'}")

# Test 5: Same test with a big integer
print()
print("=== Test 5: Bigint encoding ===")
tag5 = Tag(type=ValueType.Bigint, is_inherit=True)
val5 = Val(data="12345678901234567890", text="12345678901234567890", tag=tag5)
hex_out5 = encode_node(val5)
print(f"Python bigint (is_inherit=True): {hex_out5}")

tag5b = Tag(type=ValueType.Bigint, is_inherit=False)
val5b = Val(data="12345678901234567890", text="12345678901234567890", tag=tag5b)
hex_out5b = encode_node(val5b)
print(f"Python bigint (is_inherit=False): {hex_out5b}")

# Test 6: What does Go produce for the same?
print()
print("=== Test 6: Go encoding for comparison ===")
go_result = subprocess.run(
    'cd /Users/lizongying/IdeaProjects/meta-message/tests && go run -e - <<\'EOF\'\npackage main\n\nimport (\n\t"fmt"\n\t"math/big"\n\t"encoding/hex"\n\t"github.com/metamessage/metamessage/internal/core"\n\t"github.com/metamessage/metamessage/internal/ir"\n)\n\nfunc main() {\n\t// i8 value, inherited\n\ttag := ir.NewTag()\n\ttag.Type = ir.ValueTypeI8\n\ttag.IsInherit = true\n\tval := &ir.Value{Data: int8(1), Text: "1", Tag: tag}\n\tbs, _ := core.NewEncoder(nil).Encode(val)\n\tfmt.Println(hex.EncodeToString(bs))\n}\nEOF',
    shell=True, capture_output=True, text=True
)
print(f"Go encode result: {go_result.stdout.strip()}")
print(f"Go encode stderr: {go_result.stderr[:200] if go_result.stderr else ''}")

# Let me also test the Python encoder directly with a known fixture
print()
print("=== Test 7: Python encoder tag behavior ===")
# Create a tag with child_type=i8
parent_tag = Tag(type=ValueType.Arr, child_type=ValueType.I8, child_desc="test items")
tag = Tag()
tag.inherit(parent_tag)
print(f"After inherit:")
print(f"  type={tag.type}")
print(f"  is_inherit={tag.is_inherit}")
print(f"  desc={tag.desc}")
print(f"  bytes()={tag.bytes().hex() if tag.bytes() else '(empty)'}")

# Now encode a value with this tag
val = Val(data=1, text="1", tag=tag)
e = Encoder()
hex_out = e.encode(val).hex()
print(f"  encoded hex: {hex_out}")

# Now try without inheritance
tag_no_inherit = Tag(type=ValueType.I8)
val_no_inherit = Val(data=1, text="1", tag=tag_no_inherit)
e2 = Encoder()
hex_out2 = e2.encode(val_no_inherit).hex()  
print(f"Without inherit:")
print(f"  bytes()={tag_no_inherit.bytes().hex()}")
print(f"  encoded hex: {hex_out2}")