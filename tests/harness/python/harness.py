#!/usr/bin/env python3
"""MetaMessage Python test harness - parse JSONC file and re-print to JSONC."""
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', '..', 'mm-py'))
from metamessage import parse_jsonc, to_jsonc

if len(sys.argv) < 2:
    sys.stderr.write("usage: harness.py <file.jsonc>\n")
    sys.exit(1)

with open(sys.argv[1], 'r') as f:
    data = f.read()

try:
    node = parse_jsonc(data)
    output = to_jsonc(node)
    sys.stdout.write(output)
except Exception as e:
    sys.stderr.write(f"parse error: {e}\n")
    sys.exit(1)