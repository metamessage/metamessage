"""
ValueToNode / NodeToValue API Example

Shows low-level node operations: creating nodes manually,
converting between Python values and nodes, and working
with the intermediate node representation.
"""
from metamessage import (
    value_to_node, node_to_value, encode_from_value, decode_to_value,
    Tag, ValueType, Obj, Arr, Val, Field, mm,
    Encoder, Decoder
)


def main():
    # =============================================
    # value_to_node: Python value → Node tree
    # =============================================
    
    print("=== value_to_node (Python → Node) ===")
    
    # String value
    node = value_to_node("hello")
    print(f"String → Val: data={node.data!r}, text={node.text!r}, type={node.tag.type}")
    
    # Integer value
    node = value_to_node(42)
    print(f"Int → Val: data={node.data!r}, text={node.text!r}, type={node.tag.type}")
    
    # Dict value
    node = value_to_node({"name": "Alice", "age": 30})
    print(f"Dict → Obj: fields={[f.key for f in node.fields]}")
    for f in node.fields:
        print(f"  {f.key}: {f.value.data!r} ({f.value.tag.type.name})")
    
    # List value
    node = value_to_node([1, 2, 3])
    print(f"List → Arr: items={[item.data for item in node.items]}")
    
    # =============================================
    # node_to_value: Node tree → Python value
    # =============================================
    
    print("\n=== node_to_value (Node → Python) ===")
    
    # Create a node manually
    obj = Obj(
        fields=[
            Field(key="x", value=Val(data=10, text="10", tag=Tag(type=ValueType.Int))),
            Field(key="y", value=Val(data=20, text="20", tag=Tag(type=ValueType.Int))),
        ],
        tag=Tag(name="point")
    )
    
    result = node_to_value(obj, dict)
    print(f"Obj Node → dict: {result}")
    
    from typing import Any
    result = node_to_value(obj, Any)
    print(f"Obj Node → Any: {result}")
    
    # =============================================
    # Full pipeline: value → node → binary → node → value
    # =============================================
    
    print("\n=== Full Pipeline ===")
    
    original = {"message": "Hello MetaMessage!", "count": 42}
    
    # Step 1: Value → Node
    node = value_to_node(original)
    print(f"1. Value → Node: {type(node).__name__} with {len(node.fields)} fields")
    
    # Step 2: Node → Binary (via Encoder)
    encoder = Encoder()
    binary = encoder.encode(node)
    print(f"2. Node → Binary: {len(binary)} bytes")
    
    # Step 3: Binary → Value (via Decoder)
    decoder = Decoder(binary)
    decoded = decoder.decode()
    print(f"3. Binary → Value: {decoded}")
    
    # Or use the high-level functions
    binary2 = encode_from_value(original)
    decoded2 = decode_to_value(binary2)
    print(f"   High-level: {decoded2}")
    print(f"   Match: {original == decoded2}")


if __name__ == "__main__":
    main()
