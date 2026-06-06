#!/usr/bin/env python3
import subprocess
import sys
import os

PROJECT_DIR = "/Users/lizongying/IdeaProjects/meta-message"
FIXTURE = os.path.join(PROJECT_DIR, "tests/fixtures/03_tags/child_tags.jsonc")

def hex_to_bytes(h):
    return bytes.fromhex(h)

def get_encode_hex(lang):
    """Run encode for a specific language and get the hex output."""
    with open(FIXTURE, 'r') as f:
        content = f.read()

    script = os.path.join(PROJECT_DIR, "tests/encode_decode_test.sh")
    # Just call the encode directly
    if lang == "go":
        cmd = f"cd {PROJECT_DIR} && go run tests/harness/go/harness.go --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "rs":
        cmd = f"cd {PROJECT_DIR} && tests/harness/rust/target/debug/mm-harness-rs --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "py":
        cmd = f"cd {PROJECT_DIR} && python3 tests/harness/python/harness.py --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "ts":
        cmd = f"cd {PROJECT_DIR} && node tests/harness/typescript/harness.cjs --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "c":
        cmd = f"cd {PROJECT_DIR} && tests/harness/c/build/mm_harness_c --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "cs":
        cmd = f"cd {PROJECT_DIR} && tests/harness/cs/build/mm_harness_cs --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "kt":
        cmd = f"cd {PROJECT_DIR} && tests/harness/kt/build/mm_harness_kt --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "php":
        cmd = f"cd {PROJECT_DIR} && php tests/harness/php/harness.php --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    else:
        raise ValueError(f"Unknown language: {lang}")

    result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
    if result.returncode != 0:
        print(f"Error running {lang}: {result.stderr}", file=sys.stderr)
        return None

    hex_str = result.stdout.strip()
    # Remove any non-hex chars
    hex_str = ''.join(c for c in hex_str if c in '0123456789abcdefABCDEF')
    return hex_str

def analyze_diffs(go_hex, other_hex, label):
    """Compare bytes and identify differences."""
    go_bytes = hex_to_bytes(go_hex)
    other_bytes = hex_to_bytes(other_hex)

    print(f"\n=== {label} vs Go ===")
    print(f"Go length: {len(go_bytes)} bytes")
    print(f"{label} length: {len(other_bytes)} bytes")

    min_len = min(len(go_bytes), len(other_bytes))
    diffs = []
    for i in range(min_len):
        if go_bytes[i] != other_bytes[i]:
            diffs.append((i, go_bytes[i], other_bytes[i]))

    if len(go_bytes) != len(other_bytes):
        print(f"LENGTH MISMATCH: Go has {len(go_bytes)} bytes, {label} has {len(other_bytes)} bytes")
        if len(go_bytes) > len(other_bytes):
            print(f"Go extra bytes at end: {go_bytes[min_len:].hex()}")
        else:
            print(f"{label} extra bytes at end: {other_bytes[min_len:].hex()}")

    print(f"Total byte diffs: {len(diffs)}")

    # Print first 30 diffs with context
    for i, (pos, gb, ob) in enumerate(diffs[:40]):
        # Show context (10 bytes before)
        ctx_start = max(0, pos - 5)
        ctx_end = min(min_len, pos + 5)
        go_ctx = go_bytes[ctx_start:ctx_end]
        ob_ctx = other_bytes[ctx_start:ctx_end]
        marker = 5 if pos - ctx_start < 5 else pos - ctx_start
        go_str = ' '.join(f'{b:02x}' for b in go_ctx)
        ob_str = ' '.join(f'{b:02x}' for b in ob_ctx)
        print(f"  Byte {pos}: Go={gb:02x} vs {label}={ob:02x}")
        print(f"    Go ctx:  {go_str}  (offset from ctx_start: {pos - ctx_start})")
        print(f"    {label} ctx: {ob_str}")

    return diffs

def main():
    go_hex = get_encode_hex("go")
    if not go_hex:
        print("Failed to get Go hex")
        return

    languages = ["rs", "py", "ts", "c", "cs", "kt", "php"]
    for lang in languages:
        other_hex = get_encode_hex(lang)
        if other_hex:
            analyze_diffs(go_hex, other_hex, lang.upper())

if __name__ == "__main__":
    main()