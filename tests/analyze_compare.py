#!/usr/bin/env python3
import subprocess, os

PROJECT_DIR = "/Users/lizongying/IdeaProjects/meta-message"
FIXTURE = os.path.join(PROJECT_DIR, "tests/fixtures/03_tags/child_tags.jsonc")

def get_hex(lang):
    if lang == "go":
        cmd = f"cd {PROJECT_DIR} && go run tests/harness/go/harness.go --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "rs":
        cmd = f"cd {PROJECT_DIR} && tests/harness/rust/target/debug/mm-harness-rs --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "py":
        cmd = f"cd {PROJECT_DIR} && python3 tests/harness/python/harness.py --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    elif lang == "ts":
        cmd = f"cd {PROJECT_DIR} && node tests/harness/typescript/harness.cjs --encode {FIXTURE} 2>/dev/null | xxd -p -c 999999"
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
    hex_str = result.stdout.strip()
    hex_str = "".join(c for c in hex_str if c in "0123456789abcdefABCDEF")
    return hex_str

go = get_hex("go")
rs = get_hex("rs")
py = get_hex("py")
ts = get_hex("ts")

go_bytes = bytes.fromhex(go)
rs_bytes = bytes.fromhex(rs)
py_bytes = bytes.fromhex(py)
ts_bytes = bytes.fromhex(ts)

print(f"Go length: {len(go_bytes)} bytes")
print(f"RS length: {len(rs_bytes)} bytes")
print(f"PY length: {len(py_bytes)} bytes")
print(f"TS length: {len(ts_bytes)} bytes")

print()
print("First 20 bytes comparison (Go, RS, PY, TS):")
for i in range(20):
    g = go_bytes[i] if i < len(go_bytes) else -1
    r = rs_bytes[i] if i < len(rs_bytes) else -1
    p = py_bytes[i] if i < len(py_bytes) else -1
    t = ts_bytes[i] if i < len(ts_bytes) else -1
    marker = " <--" if g != r or g != p or g != t else ""
    print(f"  [{i:4d}] Go={g:3d}(0x{g:02x}) RS={r:3d}(0x{r:02x}) PY={p:3d}(0x{p:02x}) TS={t:3d}(0x{t:02x}){marker}")

print()
print("Header analysis:")
for lang_name, lang_bytes in [("Go", go_bytes), ("RS", rs_bytes), ("PY", py_bytes), ("TS", ts_bytes)]:
    obj_len = (lang_bytes[1] << 8) | lang_bytes[2]
    inner_arr_byte = lang_bytes[3]
    inner_arr_len_sig = inner_arr_byte & 0x0F
    if inner_arr_len_sig < 14:
        inner_len = inner_arr_len_sig
        inner_arr_header = 1
    elif inner_arr_len_sig == 14:
        inner_len = lang_bytes[4]
        inner_arr_header = 2
    else:
        inner_len = (lang_bytes[4] << 8) | lang_bytes[5]
        inner_arr_header = 3
    print(f"  {lang_name}: obj_len=0x{obj_len:04x}({obj_len}) inner_arr_header={inner_arr_header} inner_arr_len={inner_len}")
    
    # key array starts right after inner header
    key_start = 3 + inner_arr_header
    key_bytes = lang_bytes[key_start:key_start+inner_len]
    print(f"    Keys ({inner_len} bytes) at byte {key_start}: ...{key_bytes[-10:].hex()}")
    
    # values start after keys
    val_start = key_start + inner_len
    print(f"    Values start at byte {val_start}")

# Now find the first 5 diffs between Go and RS in the value section
go_obj_len = (go_bytes[1] << 8) | go_bytes[2]
rs_obj_len = (rs_bytes[1] << 8) | rs_bytes[2]

# Keys are the same for both (confirmed earlier)
go_inner_arr_len_sig = go_bytes[3] & 0x0F
if go_inner_arr_len_sig < 14:
    go_inner_len = go_inner_arr_len_sig
    go_inner_header = 1
elif go_inner_arr_len_sig == 14:
    go_inner_len = go_bytes[4]
    go_inner_header = 2
else:
    go_inner_len = (go_bytes[4] << 8) | go_bytes[5]
    go_inner_header = 3

go_key_start = 3 + go_inner_header
go_val_start = go_key_start + go_inner_len
go_val_payload = go_bytes[go_val_start:]

rs_inner_arr_len_sig = rs_bytes[3] & 0x0F
if rs_inner_arr_len_sig < 14:
    rs_inner_len = rs_inner_arr_len_sig
    rs_inner_header = 1
elif rs_inner_arr_len_sig == 14:
    rs_inner_len = rs_bytes[4]
    rs_inner_header = 2
else:
    rs_inner_len = (rs_bytes[4] << 8) | rs_bytes[5]
    rs_inner_header = 3

rs_key_start = 3 + rs_inner_header
rs_val_start = rs_key_start + rs_inner_len
rs_val_payload = rs_bytes[rs_val_start:]

py_inner_arr_len_sig = py_bytes[3] & 0x0F
if py_inner_arr_len_sig < 14:
    py_inner_len = py_inner_arr_len_sig
    py_inner_header = 1
elif py_inner_arr_len_sig == 14:
    py_inner_len = py_bytes[4]
    py_inner_header = 2
else:
    py_inner_len = (py_bytes[4] << 8) | py_bytes[5]
    py_inner_header = 3

py_key_start = 3 + py_inner_header
py_val_start = py_key_start + py_inner_len
py_val_payload = py_bytes[py_val_start:]

print()
print(f"Go value payload: {len(go_val_payload)} bytes")
print(f"RS value payload: {len(rs_val_payload)} bytes")
print(f"PY value payload: {len(py_val_payload)} bytes")
print()

# Find first 30 diffs in value payload
min_val_len = min(len(go_val_payload), len(rs_val_payload))
print(f"First 30 byte diffs in value section (Go vs RS):")
diff_count = 0
for i in range(min_val_len):
    if go_val_payload[i] != rs_val_payload[i]:
        # Show context
        ctx_start = max(0, i-3)
        ctx_end = min(min_val_len, i+4)
        go_ctx = " ".join(f"{b:02x}" for b in go_val_payload[ctx_start:ctx_end])
        rs_ctx = " ".join(f"{b:02x}" for b in rs_val_payload[ctx_start:ctx_end])
        print(f"  ValByte {i:4d}: Go=0x{go_val_payload[i]:02x} RS=0x{rs_val_payload[i]:02x}")
        print(f"           Go: {go_ctx}")
        print(f"           RS: {rs_ctx}")
        diff_count += 1
        if diff_count >= 30:
            break