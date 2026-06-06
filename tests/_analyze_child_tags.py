#!/usr/bin/env python3
"""Clean byte-level comparison of Go vs Python encoding for child_tags.jsonc"""
import subprocess
import os

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = 'fixtures/03_tags/child_tags.jsonc'

def get_hex(lang, cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    raw = result.stdout.strip()
    return bytes.fromhex(raw) if raw else None

go = get_hex('Go', f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null')
py = get_hex('Py', f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null')
php = get_hex('PHP', f'cd {TESTS_DIR} && php harness/php/harness.php --encode {FIXTURE} 2>/dev/null')

print(f"Go  : {len(go)} bytes")
print(f"Py  : {len(py)} bytes")
print(f"PHP : {len(php)} bytes")
print()

# Find all differences
print("=== Byte-level differences ===")
for i in range(min(len(go), len(py))):
    if go[i] != py[i]:
        ctx_start = max(0, i - 2)
        ctx_end = min(len(go), i + 20)
        print(f"Offset {i}:")
        print(f"  Go: {' '.join(f'{b:02x}' for b in go[ctx_start:ctx_end])}")
        print(f"  Py: {' '.join(f'{b:02x}' for b in py[ctx_start:ctx_end])}")
        print()

# Compare PHP vs Go to see if PHP matches
print("=== PHP vs Go differences ===")
php_diffs = 0
for i in range(min(len(go), len(php))):
    if go[i] != php[i]:
        php_diffs += 1
        if php_diffs <= 5:
            ctx_start = max(0, i - 2)
            ctx_end = min(len(go), i + 20)
            print(f"Offset {i}:")
            print(f"  Go:  {' '.join(f'{b:02x}' for b in go[ctx_start:ctx_end])}")
            print(f"  PHP: {' '.join(f'{b:02x}' for b in php[ctx_start:ctx_end])}")
            print()
print(f"Total PHP vs Go diffs: {php_diffs}")

# Count total differences
total_diffs = sum(1 for i in range(min(len(go), len(py))) if go[i] != py[i])
print(f"\nTotal Go vs Py diffs: {total_diffs}")