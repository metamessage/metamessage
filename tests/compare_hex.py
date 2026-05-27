import sys

cpp_hex = "ce6bde158a637265617465645f617485656d61696c83756964f3111e0d6372656174696f6e2074696d65201620f402201d9074657374406578616d706c652e636f6dfe29022019be2435353065383430302d653239622d343164342d613731362d343436363535343430303030"
go_hex = "ce59de158a637265617465645f617485656d61696c83756964f7111e0d6372656174696f6e2074696d6520163b65920080f402201d9074657374406578616d706c652e636f6df4022019b0550e8400e29b41d4a716446655440000"

cpp = bytes.fromhex(cpp_hex)
go = bytes.fromhex(go_hex)

print(f"Length: C++={len(cpp)}, Go={len(go)}")
print()

# Find match prefix
match_len = 0
for i in range(min(len(cpp), len(go))):
    if cpp[i] == go[i]:
        match_len += 1
    else:
        break

print(f"Matching prefix: {match_len} bytes")
print(f"C++ first {match_len}: {cpp[:match_len].hex()}")
print(f"Go  first {match_len}: {go[:match_len].hex()}")
print()

# Show differences
print("=== Differences ===")
i = match_len
while i < max(len(cpp), len(go)):
    cb = cpp[i] if i < len(cpp) else 0
    gb = go[i] if i < len(go) else 0
    if cb != gb:
        # Show a block of context
        start = max(0, i - 2)
        end = min(max(len(cpp), len(go)), i + 8)
        cpp_block = cpp[start:end]
        go_block = go[start:end]

        cpp_str = " ".join(f"{b:02x}" for b in cpp_block)
        go_str = " ".join(f"{b:02x}" for b in go_block)

        marker_start = i - start
        marker = " " * (marker_start * 3) + "^^"

        print(f"Pos {i}:")
        print(f"  C++: {cpp_str}")
        print(f"  Go:  {go_str}")
        print(f"       {marker}")
        print()

        # Skip ahead to next matching position
        # Find next matching byte
        next_match = i + 1
        while next_match < min(len(cpp), len(go)):
            if cpp[next_match] == go[next_match]:
                break
            next_match += 1
        if next_match < min(len(cpp), len(go)):
            # Check if there's a long match stretch ahead
            j = next_match
            match_run = 0
            while j < min(len(cpp), len(go)) and cpp[j] == go[j]:
                match_run += 1
                j += 1
            if match_run >= 2:
                i = next_match
                continue
        i += 1
    else:
        i += 1