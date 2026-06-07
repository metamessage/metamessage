from metamessage.jsonc import to_jsonc
from metamessage.tag import Tag, ValueType
from metamessage.types import NodeScalar

# Test: 類型感知打印 - String 應該加引號
v_str = NodeScalar(data='hello', text='hello', tag=Tag(type=ValueType.Str))
print(f'String: {to_jsonc(v_str)}')

# Test: Int 不應加引號
v_int = NodeScalar(data=42, text='42', tag=Tag(type=ValueType.I))
print(f'Int: {to_jsonc(v_int)}')

# Test: UUID 應該加引號
v_uuid = NodeScalar(data='abc123', text='abc123', tag=Tag(type=ValueType.Uuid))
print(f'UUID: {to_jsonc(v_uuid)}')

# Test: DateTime 應該加引號
v_dt = NodeScalar(data='2024-01-01', text='2024-01-01', tag=Tag(type=ValueType.Datetime))
print(f'DateTime: {to_jsonc(v_dt)}')

# Test: Email 應該加引號
v_email = NodeScalar(data='test@test.com', text='test@test.com', tag=Tag(type=ValueType.Email))
print(f'Email: {to_jsonc(v_email)}')

# Test: URL 應該加引號
v_url = NodeScalar(data='https://example.com', text='https://example.com', tag=Tag(type=ValueType.Url))
print(f'URL: {to_jsonc(v_url)}')

# Test: Enum 應該加引號
v_enum = NodeScalar(data='a', text='a', tag=Tag(type=ValueType.Enums))
print(f'Enum: {to_jsonc(v_enum)}')

# Test: Bytes 應該加引號
v_bytes = NodeScalar(data=b'hello', text='hello', tag=Tag(type=ValueType.Bytes))
print(f'Bytes: {to_jsonc(v_bytes)}')

# Test: Date 應該加引號
v_date = NodeScalar(data='2024-01-01', text='2024-01-01', tag=Tag(type=ValueType.Date))
print(f'Date: {to_jsonc(v_date)}')

# Test: Time 應該加引號
v_time = NodeScalar(data='12:30:00', text='12:30:00', tag=Tag(type=ValueType.Time))
print(f'Time: {to_jsonc(v_time)}')

# Test: IP 應該加引號
v_ip = NodeScalar(data='127.0.0.1', text='127.0.0.1', tag=Tag(type=ValueType.Ip))
print(f'IP: {to_jsonc(v_ip)}')

# Test: Decimal 應該加引號
v_decimal = NodeScalar(data='3.14', text='3.14', tag=Tag(type=ValueType.Decimal))
print(f'Decimal: {to_jsonc(v_decimal)}')