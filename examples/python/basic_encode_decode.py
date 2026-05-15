"""
Basic Encode/Decode Example

Shows how to encode and decode simple Python values using MetaMessage.
"""
from metamessage import encode_from_value, decode_to_value


def main():
    # =============================================
    # Basic scalar values
    # =============================================
    
    # String
    binary = encode_from_value("hello world")
    result = decode_to_value(binary)
    print(f"String: encode_from_value('hello world') -> decode_to_value -> {result!r}")
    
    # Integer
    binary = encode_from_value(42)
    result = decode_to_value(binary)
    print(f"Int: encode_from_value(42) -> {result}")
    
    # Float
    binary = encode_from_value(3.14)
    result = decode_to_value(binary)
    print(f"Float: encode_from_value(3.14) -> {result}")
    
    # Boolean
    binary = encode_from_value(True)
    result = decode_to_value(binary)
    print(f"Bool: encode_from_value(True) -> {result}")
    
    # Bytes
    binary = encode_from_value(b"raw bytes")
    result = decode_to_value(binary)
    print(f"Bytes: encode_from_value(b'raw bytes') -> {result}")
    
    # None (must provide a tag with type)
    from metamessage import Tag, ValueType, value_to_node
    tag = Tag(type=ValueType.String, nullable=True)
    node = value_to_node(None, tag)
    from metamessage import Encoder, Decoder
    binary = Encoder().encode(node)
    decoded = Decoder(binary).decode()
    print(f"None: encode/decode with nullable tag -> {decoded!r}")


if __name__ == "__main__":
    main()
