"""
Dict and List Example

Shows how to encode/decode dicts and lists (maps and slices/arrays).
"""
from metamessage import encode_from_value, decode_to_value


def main():
    # =============================================
    # Dict (Map)
    # =============================================
    
    person = {
        "name": "Alice",
        "age": 30,
        "active": True,
        "scores": [85, 92, 78],
    }
    
    binary = encode_from_value(person)
    result = decode_to_value(binary)
    
    print("=== Dict (Map) ===")
    print(f"Original: {person}")
    print(f"Encoded size: {len(binary)} bytes")
    print(f"Decoded: {result}")
    
    # =============================================
    # List (Slice)
    # =============================================
    
    items = ["apple", "banana", "cherry"]
    binary = encode_from_value(items)
    result = decode_to_value(binary)
    
    print("\n=== List (Slice) ===")
    print(f"Original: {items}")
    print(f"Encoded size: {len(binary)} bytes")
    print(f"Decoded: {result}")
    
    # =============================================
    # Mixed types in list
    # =============================================
    
    mixed = [1, 2, 3, 4, 5]
    binary = encode_from_value(mixed)
    result = decode_to_value(binary)
    
    print("\n=== Mixed List ===")
    print(f"Original: {mixed}")
    print(f"Decoded: {result}")


if __name__ == "__main__":
    main()
