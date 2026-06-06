import re

with open('tests/results/03_tags_child_tags.jsonc.encode_diff') as f:
    content = f.read()

sections = re.findall(r'--- (\w+) \(hex\) ---\n([a-f0-9]+)', content)
hex_data = {s[0]: s[1] for s in sections}

lengths = {k: len(bytes.fromhex(v)) for k, v in hex_data.items()}
print("=== Language byte lengths ===")
for lang in ['go', 'py', 'php', 'ts', 'rs', 'cs', 'kt', 'sw']:
    print(f"  {lang}: {lengths.get(lang, 'N/A')}")

go_bytes = bytes.fromhex(hex_data.get('go', ''))

for lang in ['py', 'php', 'ts', 'rs', 'cs', 'kt', 'sw']:
    if lang not in hex_data:
        continue
    lb = bytes.fromhex(hex_data[lang])
    print(f"\n=== Go vs {lang} ===")
    print(f"  Go len={len(go_bytes)}, {lang} len={len(lb)}")
    max_len = max(len(go_bytes), len(lb))
    i = 0
    diff_count = 0
    while i < max_len:
        b1 = go_bytes[i] if i < len(go_bytes) else -1
        b2 = lb[i] if i < len(lb) else -1
        if b1 != b2:
            diff_count += 1
            chunk_start = i
            i += 1
            while i < max_len:
                b1 = go_bytes[i] if i < len(go_bytes) else -1
                b2 = lb[i] if i < len(lb) else -1
                if b1 == b2:
                    break
                i += 1
            chunk_end = i
            if diff_count <= 10:
                print(f"  Diff chunk #{diff_count}: offset {chunk_start}-{chunk_end-1} ({chunk_end-chunk_start} bytes)")
                if chunk_end <= len(go_bytes) and chunk_end <= len(lb):
                    print(f"    go: {go_bytes[chunk_start:chunk_end].hex()}")
                    print(f"    {lang}: {lb[chunk_start:chunk_end].hex()}")
                    try:
                        print(f"    go ascii: {go_bytes[chunk_start:chunk_end].decode('ascii', errors='replace')}")
                        print(f"    {lang} ascii: {lb[chunk_start:chunk_end].decode('ascii', errors='replace')}")
                    except:
                        pass
        else:
            i += 1
    print(f"  Total diff chunks: {diff_count}")