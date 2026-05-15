# MetaMessage

MetaMessage - 高性能二进制消息编码库，支持模式标记（schema tags）、JSONC 注释语法，以及 Python 装饰器驱动的序列化。

## 功能特性

- **二进制编码/解码** — 将 Python 值高效编码为紧凑二进制格式
- **`@mm` 装饰器** — 类似 Go 结构体标签的 Python 装饰器语法
- **JSONC 支持** — 解析带 `// mm:` 注释的 JSON 格式
- **Node 中间表示** — 通过 `value_to_node` / `node_to_value` 操作节点树
- **类型自动推断** — 根据 Python 类型注解自动推断 MetaMessage 类型
- **嵌套结构** — 支持字典、列表、类实例的递归转换

## 安装

### pip 安装

```bash
pip install metamessage
```

### 版本要求

- Python 3.7 或更高版本

## 快速开始

### 基础编码/解码

```python
from metamessage import encode_from_value, decode_to_value

# 编码 Python 值 → 二进制
binary = encode_from_value({"name": "Alice", "age": 30})
print(f"Encoded: {len(binary)} bytes")

# 解码二进制 → Python 值
result = decode_to_value(binary)
print(result)  # {'name': 'Alice', 'age': 30}
```

### 字符串/数字/布尔值

```python
from metamessage import encode_from_value, decode_to_value

for val in ["hello", 42, 3.14, True, b"bytes"]:
    binary = encode_from_value(val)
    result = decode_to_value(binary)
    print(f"{val!r} → {result!r}")
```

## @mm 装饰器

`@mm` 装饰器为 Python 类添加 MetaMessage 模式信息，类似 Go 的结构体标签。

`mm` 支持两种用法：
- **类级装饰器**：`@mm(desc="...")` 放在 `class` 定义上方
- **字段级标记**：`field: type = mm(desc="...")` 作为字段默认值（推荐）

### 类级 + 字段级装饰器（推荐）

```python
from metamessage import mm, ValueType, encode_from_value, decode_to_value

@mm(desc="User information")
class User:
    id: int = mm(desc="User ID")
    name: str = mm(desc="User name")
    age: int = mm(desc="User age")

    def __init__(self, id: int = 0, name: str = "", age: int = 0):
        self.id = id
        self.name = name
        self.age = age

user = User(id=1, name="Alice", age=30)
binary = encode_from_value(user)
# binary -> {'id': 1, 'name': 'Alice', 'age': 30}
```

### 带类型和约束

```python
@mm(desc="Product")
class Product:
    id: int = mm(desc="Product ID", type=ValueType.Int64, min=1)
    name: str = mm(desc="Name", min=1, max=100)
    price: float = mm(desc="Price", min=0.0, max=999999.99)

    def __init__(self, id: int = 0, name: str = "", price: float = 0.0):
        self.id = id
        self.name = name
        self.price = price
```

### 仅类级装饰器

```python
@mm(desc="User information")
class User:
    id: int
    name: str
    age: int

    def __init__(self, id: int = 0, name: str = "", age: int = 0):
        self.id = id
        self.name = name
        self.age = age
```

### 字符串标签语法

```python
@mm("desc=Product; type=i64")
class Product:
    id: int
    name: str
    price: float

    def __init__(self, id: int = 0, name: str = "", price: float = 0.0):
        self.id = id
        self.name = name
        self.price = price
```

## Node 中间表示

MetaMessage 使用 Node 树作为中间表示（IR）：

| 节点类型 | Python 类 | 描述 |
|---------|----------|------|
| `Val` | 值节点 | 存储标量值（字符串、数字、布尔等） |
| `Obj` | 对象节点 | 键值对集合（对应 dict / 类实例） |
| `Arr` | 数组节点 | 有序元素列表 |

### value_to_node / node_to_value

```python
from metamessage import value_to_node, node_to_value

# Python 值 → Node 树
node = value_to_node({"x": 10, "y": 20})

# Node 树 → Python 值
result = node_to_value(node, dict)
print(result)  # {'x': 10, 'y': 20}
```

### 手动构建 Node

```python
from metamessage import Tag, ValueType, Obj, Arr, Val, Field, Encoder, Decoder

obj = Obj(
    fields=[
        Field(key="name", value=Val(data="John", text="John",
                                     tag=Tag(type=ValueType.String))),
        Field(key="age", value=Val(data=30, text="30",
                                    tag=Tag(type=ValueType.Int))),
    ],
    tag=Tag(name="person")
)

encoder = Encoder()
binary = encoder.encode(obj)
decoder = Decoder(binary)
result = decoder.decode()
print(result)  # {'name': 'John', 'age': 30}
```

## JSONC 支持

JSONC 是带 `// mm:` 注释的 JSON 格式，用于声明模式标记：

### 解析 JSONC

```python
from metamessage import parse_jsonc, to_jsonc

jsonc = """{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 30
}"""

node = parse_jsonc(jsonc)
for field in node.fields:
    print(f"{field.key}: type={field.value.tag.type}, desc={field.value.tag.desc}")

# 转回 JSONC 字符串
output = to_jsonc(node)
print(output)
```

### JSONC -> 二进制 -> 解码

```python
from metamessage import parse_jsonc, encode_from_value, decode_to_value

# 1. 解析 JSONC 获取 Node
node = parse_jsonc(jsonc_source)

# 2. 直接编码字典
binary = encode_from_value({"name": "Alice", "age": 30})

# 3. 解码
result = decode_to_value(binary)
```

## 完整示例

查看 `examples/python/` 目录下的示例代码：

| 文件 | 说明 |
|------|------|
| `basic_encode_decode.py` | 基本编码/解码 |
| `dict_and_list.py` | 字典和列表操作 |
| `nested_structures.py` | 嵌套结构 |
| `mm_decorator.py` | @mm 装饰器使用 |
| `jsonc_roundtrip.py` | JSONC 完整流程 |
| `value_to_node_api.py` | Node API 使用 |

运行示例：

```bash
# 在 mm-py 目录下
PYTHONPATH=. python3 examples/python/basic_encode_decode.py
PYTHONPATH=. python3 examples/python/dict_and_list.py
PYTHONPATH=. python3 examples/python/nested_structures.py
```

## Tag 参数说明

### 基本参数

| 参数 | 类型 | 说明 |
|------|------|------|
| `type` | `ValueType` | 数据类型 |
| `desc` | `str` | 描述 |
| `nullable` | `bool` | 是否可空 |
| `allow_empty` | `bool` | 是否允许空值 |
| `default` | `str` | 默认值 |
| `raw` | `bool` | 是否原始模式 |
| `example` | `bool` | 是否为示例 |
| `unique` | `bool` | 是否唯一 |

### 约束参数

| 参数 | 适用类型 | 说明 |
|------|---------|------|
| `min` | int/float/bytes/string | 最小值/长度 |
| `max` | int/float/bytes/string | 最大值/长度 |
| `size` | array/bytes/string | 固定大小 |
| `pattern` | string | 正则表达式 |
| `enum` | enum | 枚举值列表（`"a\|b\|c"`） |
| `version` | uuid/ip | 版本号 |
| `mime` | bytes | MIME 类型 |
| `location` | datetime | 时区偏移（-12 ~ +14） |

### 子元素参数

| 参数 | 说明 |
|------|------|
| `child_type` | 子元素类型 |
| `child_desc` | 子元素描述 |
| `child_raw` | 子元素原始模式 |
| `child_nullable` | 子元素可空 |
| `child_allow_empty` | 子元素允许空值 |
| `child_unique` | 子元素唯一 |
| `child_default` | 子元素默认值 |
| `child_min` | 子元素最小值 |
| `child_max` | 子元素最大值 |
| `child_size` | 子元素大小 |
| `child_enum` | 子元素枚举 |
| `child_pattern` | 子元素模式 |
| `child_version` | 子元素版本 |
| `child_mime` | 子元素 MIME |

## API 参考

### 高层 API

| 函数 | 签名 | 说明 |
|------|------|------|
| `encode_from_value` | `(value: Any) -> bytes` | Python 值 → 二进制 |
| `decode_to_value` | `(data: bytes, target_type=Any) -> Any` | 二进制 → Python 值 |
| `value_to_node` | `(value: Any, tag=None, depth=0, path="") -> Node` | Python 值 → Node 树 |
| `node_to_value` | `(node: Node, target_type) -> Any` | Node 树 → Python 值 |

### 核心类

| 类 | 说明 |
|------|------|
| `Encoder` | 编码器：`Encoder().encode(node) -> bytes` |
| `Decoder` | 解码器：`Decoder(data).decode() -> Any` |
| `Tag` | 模式标记：类型、约束、描述等 |
| `ValueType` | 类型枚举：`String`, `Int`, `Float64`, `Bool` 等 |

### JSONC API

| 函数 | 签名 | 说明 |
|------|------|------|
| `parse_jsonc` | `(source: str) -> Node` | 解析 JSONC 字符串为 Node |
| `to_jsonc` | `(node: Node) -> str` | 将 Node 转回 JSONC 字符串 |

### 装饰器

| 装饰器 | 说明 |
|--------|------|
| `@mm(...)` | 类级标记装饰器 |

## 测试

```bash
# 运行主测试套件
python3 metamessage_test.py

# 运行分类测试
python3 tests/test_encoder.py
python3 tests/test_decoder.py
python3 tests/test_jsonc.py
python3 tests/test_value_to_node.py
```

## 许可

MIT
