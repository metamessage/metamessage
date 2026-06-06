#!/usr/bin/env python3
import subprocess, sys

go_result = subprocess.run(
    ['go', 'run', 'harness/go/harness.go', '--encode', 'fixtures/03_tags/child_tags.jsonc'],
    capture_output=True, text=True)
go_bytes = bytes.fromhex(go_result.stdout.strip())

ts_result = subprocess.run(
    ['node', 'harness/typescript/harness.cjs', '--encode', 'fixtures/03_tags/child_tags.jsonc'],
    capture_output=True, text=True)
ts_bytes = bytes.fromhex(ts_result.stdout.strip())

print(f'Go length: {len(go_bytes)}')
print(f'TS length: {len(ts_bytes)}')

# Byte 0 is same (0xcf), bytes 1-2 differ (container length), then content should match
# Find where content diverges after byte 2
match_len = 0
for i in range(3, min(len(go_bytes), len(ts_bytes))):
    if go_bytes[i] == ts_bytes[i]:
        match_len += 1
    else:
        print(f'Content matches for {match_len} bytes from byte 3')
        print(f'First content divergence at byte {i}:')
        ctx_start = max(3, i - 10)
        ctx_end = min(len(go_bytes), i + 30)
        for j in range(ctx_start, ctx_end):
            marker = ' <--' if j == i else ''
            ts_val = f'0x{ts_bytes[j]:02x}' if j < len(ts_bytes) else 'OUT_OF_RANGE'
            print(f'  [{j:4d}] Go: 0x{go_bytes[j]:02x}  TS: {ts_val}{marker}')
        break