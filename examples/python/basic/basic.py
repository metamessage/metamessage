#!/usr/bin/env python3
import sys
sys.path.insert(0, '../../../mm-py')
from mm import encode, decode

# 创建对象
person = {
    "name": "Ed",
    "age": 30
}

print('Original:', person)

# 编码到 Wire 格式
wire = encode(person)
print('Encoded:', bytes_to_hex(wire))

# 从 Wire 解码
decoded = decode(wire)
print('Decoded:', decoded.value)

def bytes_to_hex(bytes_data):
    return ''.join(f'{b:02x}' for b in bytes_data)
