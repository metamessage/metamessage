#!/usr/bin/env python3
"""MetaMessage Python test harness - parse JSONC file and re-print to JSONC."""
import sys
import os
import binascii

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', '..', 'mm-py'))
from metamessage import parse_jsonc, to_jsonc, encode_from_jsonc, decode_to_jsonc

def main():
    if len(sys.argv) < 2:
        sys.stderr.write("usage: harness [--encode|--decode] <file.jsonc>\n")
        sys.exit(1)

    if sys.argv[1] == "--encode":
        if len(sys.argv) < 3:
            sys.stderr.write("usage: harness --encode <file.jsonc>\n")
            sys.exit(1)
        with open(sys.argv[2], 'r') as f:
            data = f.read()
        try:
            wire = encode_from_jsonc(data)
            sys.stdout.write(binascii.hexlify(wire).decode())
        except Exception as e:
            sys.stderr.write(f"encode error: {e}\n")
            sys.exit(1)
        return

    if sys.argv[1] == "--decode":
        hex_str = sys.stdin.read().strip()
        try:
            wire = binascii.unhexlify(hex_str)
            output = decode_to_jsonc(wire)
            sys.stdout.write(output)
        except Exception as e:
            sys.stderr.write(f"decode error: {e}\n")
            sys.exit(1)
        return

    # Existing behavior
    with open(sys.argv[1], 'r') as f:
        data = f.read()
    try:
        node = parse_jsonc(data)
        output = to_jsonc(node)
        sys.stdout.write(output)
    except Exception as e:
        sys.stderr.write(f"parse error: {e}\n")
        sys.exit(1)

if __name__ == '__main__':
    main()