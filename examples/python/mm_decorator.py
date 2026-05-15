"""
@mm Decorator Example

Shows how to use the @mm decorator to annotate Python classes
with MetaMessage schema information (equivalent to Go's struct tags).

Usage:
    1. Class-level decorator:
       @mm(desc="User information")
       class User: ...
    
    2. Field-level decorator (via default value syntax):
       class User:
           id: int = mm(desc="User ID")
           name: str = mm(desc="User name")
           age: int = mm(desc="User age")
"""
from metamessage import mm, ValueType, encode_from_value, decode_to_value, node_to_value
from typing import Any


def main():
    # =============================================
    # Class-level @mm decorator + field-level mm
    # =============================================
    # This is like Go's: type User struct { ... } with struct tags
    
    @mm(desc="User information")
    class User:
        id: int = mm(desc="User ID")
        name: str = mm(desc="User name")
        age: int = mm(desc="User age")
        
        def __init__(self, id: int = 0, name: str = "", age: int = 0):
            self.id = id
            self.name = name
            self.age = age
    
    # Create and encode
    user = User(id=1, name="Alice", age=30)
    binary = encode_from_value(user)
    
    print("=== Class + Field Level @mm ===")
    print(f"Encoded size: {len(binary)} bytes")
    
    # Decode (returns dict by default)
    result = decode_to_value(binary)
    print(f"Decoded: {result}")
    
    # =============================================
    # Class with field-level type annotations
    # =============================================
    
    @mm("desc=Product")
    class Product:
        id: int = mm(desc="Product ID", type=ValueType.Int64)
        name: str = mm(desc="Product name")
        price: float = mm(desc="Price", type=ValueType.Float64)
        
        def __init__(self, id: int = 0, name: str = "", price: float = 0.0):
            self.id = id
            self.name = name
            self.price = price
    
    product = Product(id=1001, name="Laptop", price=999.99)
    binary = encode_from_value(product)
    result = decode_to_value(binary)
    
    print("\n=== Product Class ===")
    print(f"Original: {product.id}, {product.name}, {product.price}")
    print(f"Decoded: {result}")
    
    # =============================================
    # dataclass-compatible with constraints
    # =============================================
    
    @mm(desc="Point")
    class Point:
        x: int = mm(desc="X coordinate", min=-1000, max=1000)
        y: int = mm(desc="Y coordinate", min=-1000, max=1000)
        
        def __init__(self, x: int = 0, y: int = 0):
            self.x = x
            self.y = y
    
    p = Point(x=10, y=20)
    binary = encode_from_value(p)
    result = decode_to_value(binary)
    
    print("\n=== Point Class ===")
    print(f"Original: ({p.x}, {p.y})")
    print(f"Decoded: {result}")


if __name__ == "__main__":
    main()
