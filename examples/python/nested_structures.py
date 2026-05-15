"""
Nested Structures Example

Shows encoding and decoding of nested dicts, lists, and mixed structures.
"""
from metamessage import encode_from_value, decode_to_value


def main():
    # =============================================
    # Deeply nested dict
    # =============================================
    
    data = {
        "user": {
            "profile": {
                "name": "Bob",
                "address": {
                    "city": "Beijing",
                    "zip": "100000"
                }
            },
            "scores": [85, 92, 78],
            "metadata": {
                "created": "2024-01-15",
                "tags": ["premium", "active"]
            }
        },
        "status": "ok"
    }
    
    binary = encode_from_value(data)
    result = decode_to_value(binary)
    
    print("=== Nested Structures ===")
    print(f"Encoded size: {len(binary)} bytes")
    print(f"Original == Decoded: {data == result}")
    
    # Access nested values
    print(f"\nUser name: {result['user']['profile']['name']}")
    print(f"City: {result['user']['profile']['address']['city']}")
    print(f"Scores: {result['user']['scores']}")
    print(f"Tags: {result['user']['metadata']['tags']}")
    
    # =============================================
    # Array of objects
    # =============================================
    
    users = [
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
        {"id": 3, "name": "Charlie"},
    ]
    
    binary = encode_from_value(users)
    result = decode_to_value(binary)
    
    print("\n=== Array of Objects ===")
    print(f"Count: {len(result)}")
    for user in result:
        print(f"  User {user['id']}: {user['name']}")


if __name__ == "__main__":
    main()
