from metamessage.jsonc import to_jsonc
from metamessage.tag import Tag, ValueType
from metamessage.types import Val

# Test: 類型感知打印 - String 應該加引號
v_str = Val(data='hello', text='hello', tag=Tag(type=ValueType.String))
print(f'String: {to_jsonc(v_str)}')

# Test: Int 不應加引號
v_int = Val(data=42, text='42', tag=Tag(type=ValueType.Int))
print(f'Int: {to_jsonc(v_int)}')

# Test: UUID 應該加引號
v_uuid = Val(data='abc123', text='abc123', tag=Tag(type=ValueType.UUID))
print(f'UUID: {to_jsonc(v_uuid)}')

# Test: DateTime 應該加引號
v_dt = Val(data='2024-01-01', text='2024-01-01', tag=Tag(type=ValueType.DateTime))
print(f'DateTime: {to_jsonc(v_dt)}')

# Test: Email 應該加引號
v_email = Val(data='test@test.com', text='test@test.com', tag=Tag(type=ValueType.Email))
print(f'Email: {to_jsonc(v_email)}')

# Test: URL 應該加引號
v_url = Val(data='https://example.com', text='https://example.com', tag=Tag(type=ValueType.URL))
print(f'URL: {to_jsonc(v_url)}')

# Test: Enum 應該加引號
v_enum = Val(data='a', text='a', tag=Tag(type=ValueType.Enum))
print(f'Enum: {to_jsonc(v_enum)}')

# Test: Bytes 應該加引號
v_bytes = Val(data=b'hello', text='hello', tag=Tag(type=ValueType.Bytes))
print(f'Bytes: {to_jsonc(v_bytes)}')

# Test: Date 應該加引號
v_date = Val(data='2024-01-01', text='2024-01-01', tag=Tag(type=ValueType.Date))
print(f'Date: {to_jsonc(v_date)}')

# Test: Time 應該加引號
v_time = Val(data='12:30:00', text='12:30:00', tag=Tag(type=ValueType.Time))
print(f'Time: {to_jsonc(v_time)}')

# Test: IP 應該加引號
v_ip = Val(data='127.0.0.1', text='127.0.0.1', tag=Tag(type=ValueType.IP))
print(f'IP: {to_jsonc(v_ip)}')

# Test: Decimal 應該加引號
v_decimal = Val(data='3.14', text='3.14', tag=Tag(type=ValueType.Decimal))
print(f'Decimal: {to_jsonc(v_decimal)}')