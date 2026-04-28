# Python MetaMessage 库使用说明

## 1. 安装

### pip 安装
使用 pip 安装：

```bash
pip install metamessage
```

### 版本要求
- Python 3.7 或更高版本

## 2. 基本使用

### 2.1 导入模块

```python
from metamessage import encode, decode
```

### 2.2 编码示例

```python
person = {"name": "Ed", "age": 30}
wire = encode(person)
print('Encoded:', bytes_to_hex(wire))
```

### 2.3 解码示例

```python
decoded = decode(wire)
print('Decoded:', decoded.value)
```

### 2.4 JSONC 解析示例

```python
from metamessage import parse_jsonc, to_jsonc_string

jsonc = '''
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}
'''

# 解析 JSONC
node = parse_jsonc(jsonc)

# 转换为字符串
jsonc_string = to_jsonc_string(node)
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-py 目录下运行
python -m pytest metamessage_test.py -v
```

### 3.2 测试框架
- pytest

### 3.3 测试覆盖范围
- 编码测试
- 解码测试
- JSONC 解析测试

## 4. 常见问题

### 4.1 依赖问题
- **问题**: pip 安装失败
  **解决**: 检查网络连接，或使用 pip 镜像

### 4.2 运行时问题
- **问题**: 编码/解码失败
  **解决**: 检查数据类型是否正确

## 5. 示例代码

查看 `examples/python/` 目录下的示例代码：
- `basic/` - 基本使用示例

## 6. 相关资源

- [Python 文档](https://docs.python.org/3/)
- [pytest 文档](https://docs.pytest.org/en/latest/)
