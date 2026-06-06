#!/usr/bin/env python3
import subprocess, sys

go_result = subprocess.run(['go', 'run', 'harness/go/harness.go', '--encode', 'fixtures/03_tags/child_tags.jsonc'], capture_output=True, text=True)
go_bytes = bytes.fromhex(go_result.stdout.strip())

ts_result = subprocess.run(['node', 'harness/typescript/harness.cjs', '--encode', 'fixtures/03_tags/child_tags.jsonc'], capture_output=True, text=True)
ts_bytes = bytes.fromhex(ts_result.stdout.strip())

print(f'Go length: {len(go_bytes)}')
print(f'TS length: {len(ts_bytes)}')

for i in range(min(len(go_bytes), len(ts_bytes))):
    if go_bytes[i] != ts_bytes[i]:
        print(f'First diff at byte {i}:')
        print(f'  Go: ...{go_bytes[max(0,i-5):i+20].hex()}...')
        print(f'  TS: ...{ts_bytes[max(0,i-5):i+20].hex()}...')
        context_start = max(0, i-10)
        context_end = min(len(go_bytes), i+30)
        for j in range(context_start, context_end):
            marker = ' <--' if j == i else ''
            ts_val = f'0x{ts_bytes[j]:02x}' if j < len(ts_bytes) else 'OUT_OF_RANGE'
            print(f'  [{j:4d}] Go: 0x{go_bytes[j]:02x}  TS: {ts_val}{marker}')
        break

if len(go_bytes) != len(ts_bytes):
    print(f'\nLength mismatch: Go={len(go_bytes)} TS={len(ts_bytes)}')
    if len(go_bytes) > len(ts_bytes):
        extra = go_bytes[len(ts_bytes):min(len(ts_bytes)+100, len(go_bytes))]
        print(f'Extra Go bytes: {extra.hex()}')