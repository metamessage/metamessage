#!/usr/bin/env python3
"""Compare round-trip outputs to understand encoding differences"""
import subprocess
import os

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = 'fixtures/03_tags/child_tags.jsonc'

# Get Go hex
go_hex = subprocess.run(
    f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

# Get Python hex
py_hex = subprocess.run(
    f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

print(f"Go hex length: {len(go_hex)} chars = {len(go_hex)//2} bytes")
print(f"Py hex length: {len(py_hex)} chars = {len(py_hex)//2} bytes")

# Decode Go hex with Go decoder -> JSONC
go_decoded = subprocess.run(
    f'printf "{go_hex}" | cd {TESTS_DIR} && go run harness/go/harness.go --decode 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

# Decode Go hex with Python decoder -> JSONC
py_decoded = subprocess.run(
    f'printf "{go_hex}" | cd {TESTS_DIR} && python3 harness/python/harness.py --decode 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

# Decode Python hex with Go decoder -> JSONC  
go_decoded_py = subprocess.run(
    f'printf "{py_hex}" | cd {TESTS_DIR} && go run harness/go/harness.go --decode 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

print(f"\n=== Go encode -> Go decode ===  ({len(go_decoded)} chars)")
print(go_decoded[:500])
print("...")

print(f"\n=== Go encode -> Python decode ===  ({len(py_decoded)} chars)")
print(py_decoded[:500])
print("...")

print(f"\n=== Python encode -> Go decode ===  ({len(go_decoded_py)} chars)")
print(go_decoded_py[:500])
print("...")

# Also check: what does Python produce when encoding its OWN decoded output?
# If Python's encoder is correct, it should produce the SAME hex from the decoded output.
py_roundtrip_hex = subprocess.run(
    f'printf "{go_hex}" | cd {TESTS_DIR} && python3 harness/python/harness.py --decode 2>/dev/null | cd {TESTS_DIR} && python3 harness/python/harness.py --encode /dev/stdin 2>/dev/null',
    shell=True, capture_output=True, text=True
).stdout.strip()

print(f"\n=== Python round-trip (Go hex -> Py decode -> Py encode) ===")
if py_roundtrip_hex == py_hex:
    print("MATCHES Python's own encoding")
else:
    print(f"DIFFERS from Python's own encoding")
    print(f"  Roundtrip: {py_roundtrip_hex[:100]}...")
    print(f"  Original:  {py_hex[:100]}...")