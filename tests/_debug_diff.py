#!/usr/bin/env python3
"""Compare hex outputs between Go (reference) and other languages for child_tags.jsonc."""
import subprocess
import sys
import os

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
FIXTURE = os.path.join(SCRIPT_DIR, "fixtures", "03_tags", "child_tags.jsonc")

# Get hex from each language
langs = {
    "go": ["go", "run", os.path.join(SCRIPT_DIR, "harness/go/harness.go"), "--encode", FIXTURE],
    "py": ["python3", os.path.join(SCRIPT_DIR, "harness/python/harness.py"), "--encode", FIXTURE],
    "php": ["php", os.path.join(SCRIPT_DIR, "harness/php/harness.php"), "--encode", FIXTURE],
    "ts": ["node", os.path.join(SCRIPT_DIR, "harness/typescript/harness.cjs"), "--encode", FIXTURE],
    "rs": [os.path.join(SCRIPT_DIR, "harness/rust/target/debug/mm-harness-rs"), "--encode", FIXTURE],
    "cs": ["dotnet", "run", "--project", os.path.join(SCRIPT_DIR, "harness/csharp/harness.csproj"), "--no-build", "--", "--encode", FIXTURE],
    "kt": None,
    "sw": None,
}

def get_hex(lang, cmd):
    if cmd is None:
        return None, "SKIP"
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
        if result.returncode != 0:
            return None, f"ERR: {result.stderr[:200]}"
        hex_str = result.stdout.strip()
        return bytes.fromhex(hex_str), None
    except Exception as e:
        return None, str(e)

go_bytes, err = get_hex("go", langs["go"])
if err:
    print(f"Go error: {err}")
    sys.exit(1)

print(f"Go output ({len(go_bytes)} bytes):")
print(f"  {go_bytes.hex()}")
print()

for lang_name in ["py", "php", "ts", "rs", "cs", "kt", "sw"]:
    cmd = langs[lang_name]
    if cmd is None:
        print(f"{lang_name}: SKIP")
        continue
    b, err = get_hex(lang_name, cmd)
    if err:
        print(f"{lang_name}: {err}")
        continue
    print(f"{lang_name} output ({len(b)} bytes):")
    if b == go_bytes:
        print(f"  MATCH! (length: {len(b)})")
    else:
        print(f"  DIFF! (Go: {len(go_bytes)} bytes, {lang_name}: {len(b)} bytes)")
        # Find first diff
        min_len = min(len(go_bytes), len(b))
        for i in range(min_len):
            if go_bytes[i] != b[i]:
                ctx_start = max(0, i - 8)
                ctx_end = min(len(go_bytes), i + 32)
                print(f"  First diff at offset {i}:")
                print(f"    Go: {go_bytes[ctx_start:ctx_end].hex()}")
                print(f"    {lang_name}: {b[ctx_start:ctx_end].hex()}")
                print(f"    Go byte: 0x{go_bytes[i]:02x} vs {lang_name} byte: 0x{b[i]:02x}")
                break
        if len(go_bytes) != len(b):
            print(f"  Length diff: Go={len(go_bytes)}, {lang_name}={len(b)}")
            # Show the extra bytes from the longer one
            if len(go_bytes) > len(b):
                extra = go_bytes[len(b):]
                print(f"  Go extra bytes ({len(extra)}): {extra.hex()}")
            else:
                extra = b[len(go_bytes):]
                print(f"  {lang_name} extra bytes ({len(extra)}): {extra.hex()}")
    print()