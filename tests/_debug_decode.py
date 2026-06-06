#!/usr/bin/env python3
"""Debug tool to decode and compare wire-format bytes"""
import subprocess
import os
import sys

TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_DIR = os.path.dirname(TESTS_DIR)
FIXTURE = sys.argv[1] if len(sys.argv) > 1 else 'fixtures/03_tags/child_tags.jsonc'

def run_encode(lang, cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    hex_out = result.stdout.strip()
    if not hex_out:
        print(f"{lang}: ERROR - {result.stderr}")
        return None
    try:
        return bytes.fromhex(hex_out)
    except ValueError as e:
        print(f"{lang}: HEX ERROR - {e}")
        return None

# Run all languages
results = {}
for lang, cmd in [
    ('Go', f'cd {TESTS_DIR} && go run harness/go/harness.go --encode {FIXTURE} 2>/dev/null'),
    ('Py', f'cd {TESTS_DIR} && python3 harness/python/harness.py --encode {FIXTURE} 2>/dev/null'),
    ('PHP', f'cd {TESTS_DIR} && php harness/php/harness.php --encode {FIXTURE} 2>/dev/null'),
    ('TS', f'cd {TESTS_DIR} && node harness/ts/harness.mjs --encode {FIXTURE} 2>/dev/null'),
    ('RS', f'cd {TESTS_DIR} && cargo run --manifest-path ../mm-rs/Cargo.toml -- encode {FIXTURE} 2>/dev/null'),
    ('CS', f'cd {TESTS_DIR} && dotnet run --project ../mm-cs/src/MetaMessage -- encode {FIXTURE} 2>/dev/null'),
    ('KT', f'cd {PROJECT_DIR}/mm-kt && mvn -q exec:java -Dexec.mainClass="com.metamessage.Main" -Dexec.args="encode ../tests/{FIXTURE}" 2>/dev/null'),
]:
    data = run_encode(lang, cmd)
    if data:
        results[lang] = data

print(f"\n{'='*60}")
print(f"File: {FIXTURE}")
print(f"{'='*60}")
for lang, data in results.items():
    print(f"{lang:4s}: {len(data):5d} bytes")

if 'Go' in results:
    go_data = results['Go']
    for lang, data in results.items():
        if lang == 'Go':
            continue
        match_len = sum(1 for i in range(min(len(go_data), len(data))) if go_data[i] == data[i])
        match_pct = 100.0 * match_len / max(len(go_data), len(data))
        print(f"\nGo vs {lang}: {match_len}/{max(len(go_data), len(data))} bytes match ({match_pct:.1f}%)")
        
        # Find first difference
        for i in range(min(len(go_data), len(data))):
            if go_data[i] != data[i]:
                ctx_start = max(0, i-20)
                ctx_end = min(len(go_data), i+40)
                print(f"  First diff at offset {i}:")
                print(f"    Go: {go_data[ctx_start:ctx_end].hex()}")
                print(f"    {lang}: {data[ctx_start:ctx_end].hex()}")
                break