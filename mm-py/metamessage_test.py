import unittest
from datetime import datetime, date, time as dt_time
from metamessage import Tag, ValueType, Obj, Arr, Val, Field, Encoder, Decoder, parse_jsonc, to_jsonc


class TestParseJSONC(unittest.TestCase):
    def test_parse_simple_object(self):
        source = '{"name": "John", "age": 30}'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Obj)
        self.assertEqual(len(result.fields), 2)
        
        name_field = next(f for f in result.fields if f.key == "name")
        self.assertEqual(name_field.value.data, "John")
        
        age_field = next(f for f in result.fields if f.key == "age")
        self.assertEqual(age_field.value.data, 30)

    def test_parse_simple_array(self):
        source = '["a", "b", "c"]'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Arr)
        self.assertEqual(len(result.items), 3)

    def test_parse_nested_object(self):
        source = '{"user": {"name": "John", "age": 30}}'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Obj)
        
        user_field = next(f for f in result.fields if f.key == "user")
        self.assertIsInstance(user_field.value, Obj)

    def test_parse_array_of_objects(self):
        source = '[{"id": 1}, {"id": 2}]'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Arr)
        self.assertEqual(len(result.items), 2)

    def test_parse_with_comments(self):
        source = '''{
            "name": "John", // mm: type=str;desc=name
            "age": 30
        }'''
        result = parse_jsonc(source)
        self.assertIsInstance(result, Obj)

    def test_parse_number_types(self):
        source = '{"int": 42, "float": 3.14}'
        result = parse_jsonc(source)
        
        int_field = next(f for f in result.fields if f.key == "int")
        self.assertEqual(int_field.value.data, 42)
        
        float_field = next(f for f in result.fields if f.key == "float")
        self.assertEqual(float_field.value.data, 3.14)

    def test_parse_bool_null(self):
        source = '{"active": true, "deleted": false, "data": null}'
        result = parse_jsonc(source)
        
        active_field = next(f for f in result.fields if f.key == "active")
        self.assertEqual(active_field.value.data, True)
        
        deleted_field = next(f for f in result.fields if f.key == "deleted")
        self.assertEqual(deleted_field.value.data, False)
        
        data_field = next(f for f in result.fields if f.key == "data")
        self.assertIsNone(data_field.value.data)

    def test_parse_empty_object(self):
        source = '{}'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Obj)
        self.assertEqual(len(result.fields), 0)

    def test_parse_empty_array(self):
        source = '[]'
        result = parse_jsonc(source)
        self.assertIsInstance(result, Arr)
        self.assertEqual(len(result.items), 0)

    def test_parse_string_with_spaces(self):
        source = '{"name": "Hello World"}'
        result = parse_jsonc(source)
        
        name_field = next(f for f in result.fields if f.key == "name")
        self.assertEqual(name_field.value.data, "Hello World")


class TestEncoderValueTypes(unittest.TestCase):
    def setUp(self):
        self.encoder = Encoder()

    def test_encode_string(self):
        t = Tag(type=ValueType.Str)
        v = Val(data="hello", text="hello", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_string_with_tag(self):
        t = Tag(type=ValueType.Str, desc="test desc")
        v = Val(data="hello", text="hello", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_int(self):
        t = Tag(type=ValueType.I)
        v = Val(data=42, text="42", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_int8(self):
        t = Tag(type=ValueType.I8)
        v = Val(data=10, text="10", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_int16(self):
        t = Tag(type=ValueType.I16)
        v = Val(data=100, text="100", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_int32(self):
        t = Tag(type=ValueType.I32)
        v = Val(data=1000, text="1000", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_int64(self):
        t = Tag(type=ValueType.I64)
        v = Val(data=100000, text="100000", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uint(self):
        t = Tag(type=ValueType.U)
        v = Val(data=42, text="42", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uint8(self):
        t = Tag(type=ValueType.U8)
        v = Val(data=10, text="10", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uint16(self):
        t = Tag(type=ValueType.U16)
        v = Val(data=100, text="100", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uint32(self):
        t = Tag(type=ValueType.U32)
        v = Val(data=1000, text="1000", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uint64(self):
        t = Tag(type=ValueType.U64)
        v = Val(data=100000, text="100000", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_float32(self):
        t = Tag(type=ValueType.F32)
        v = Val(data=3.14, text="3.14", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_float64(self):
        t = Tag(type=ValueType.F64)
        v = Val(data=3.14, text="3.14", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_bool_true(self):
        t = Tag(type=ValueType.Bool)
        v = Val(data=True, text="true", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_bool_false(self):
        t = Tag(type=ValueType.Bool)
        v = Val(data=False, text="false", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_bytes(self):
        t = Tag(type=ValueType.Bytes)
        v = Val(data=b"hello", text="hello", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_bigint(self):
        t = Tag(type=ValueType.Bigint)
        v = Val(data="12345678901234567890", text="12345678901234567890", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_datetime(self):
        t = Tag(type=ValueType.Datetime)
        v = Val(data=datetime.now(), text="2024-01-01T00:00:00Z", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_date(self):
        t = Tag(type=ValueType.Date)
        v = Val(data=date.today(), text="2024-01-01", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_time(self):
        t = Tag(type=ValueType.Time)
        v = Val(data=dt_time(12, 30, 0), text="12:30:00", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_uuid(self):
        t = Tag(type=ValueType.Uuid)
        v = Val(data="0123456789abcdef", text="0123456789abcdef0123456789abcdef", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_decimal(self):
        t = Tag(type=ValueType.Decimal)
        v = Val(data="3.14159", text="3.14159", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_url(self):
        t = Tag(type=ValueType.Url)
        v = Val(data="https://example.com", text="https://example.com", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_email(self):
        t = Tag(type=ValueType.Email)
        v = Val(data="test@example.com", text="test@example.com", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_ip_v4(self):
        t = Tag(type=ValueType.Ip, version=4)
        v = Val(data="127.0.0.1", text="127.0.0.1", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_ip_v6(self):
        t = Tag(type=ValueType.Ip, version=6)
        v = Val(data="::1", text="::1", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_enum(self):
        t = Tag(type=ValueType.Enum, enum="a|b|c")
        v = Val(data=0, text="a", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)


class TestEncoderComplexTypes(unittest.TestCase):
    def setUp(self):
        self.encoder = Encoder()

    def test_encode_array_empty(self):
        arr = Arr(items=[], tag=Tag(type=ValueType.Arr))
        result = self.encoder.encode(arr)
        self.assertIsInstance(result, bytes)

    def test_encode_array_with_strings(self):
        arr = Arr(
            items=[
                Val(data="a", text="a", tag=Tag(type=ValueType.Str)),
                Val(data="b", text="b", tag=Tag(type=ValueType.Str)),
            ],
            tag=Tag(type=ValueType.Arr, child_type=ValueType.Str)
        )
        result = self.encoder.encode(arr)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_array_with_ints(self):
        arr = Arr(
            items=[
                Val(data=1, text="1", tag=Tag(type=ValueType.I)),
                Val(data=2, text="2", tag=Tag(type=ValueType.I)),
            ],
            tag=Tag(type=ValueType.Arr, child_type=ValueType.I)
        )
        result = self.encoder.encode(arr)
        self.assertIsInstance(result, bytes)

    def test_encode_object_empty(self):
        obj = Obj(fields=[], tag=Tag())
        result = self.encoder.encode(obj)
        self.assertIsInstance(result, bytes)

    def test_encode_object_with_fields(self):
        obj = Obj(
            fields=[
                Field(
                    key="name",
                    value=Val(data="John", text="John", tag=Tag(type=ValueType.Str))
                ),
                Field(
                    key="age",
                    value=Val(data=30, text="30", tag=Tag(type=ValueType.I))
                ),
            ],
            tag=Tag(name="person")
        )
        result = self.encoder.encode(obj)
        self.assertIsInstance(result, bytes)
        self.assertGreater(len(result), 0)

    def test_encode_object_nested(self):
        inner = Obj(
            fields=[
                Field(
                    key="city",
                    value=Val(data="Beijing", text="Beijing", tag=Tag(type=ValueType.Str))
                ),
            ],
            tag=Tag(name="address")
        )
        outer = Obj(
            fields=[
                Field(
                    key="name",
                    value=Val(data="John", text="John", tag=Tag(type=ValueType.Str))
                ),
                Field(
                    key="address",
                    value=inner
                ),
            ],
            tag=Tag(name="person")
        )
        result = self.encoder.encode(outer)
        self.assertIsInstance(result, bytes)


class TestEncoderEdgeCases(unittest.TestCase):
    def setUp(self):
        self.encoder = Encoder()

    def test_encode_negative_int(self):
        t = Tag(type=ValueType.I)
        v = Val(data=-42, text="-42", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_large_int(self):
        t = Tag(type=ValueType.I64)
        v = Val(data=9223372036854775807, text="9223372036854775807", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_negative_float(self):
        t = Tag(type=ValueType.F64)
        v = Val(data=-3.14, text="-3.14", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_zero_values(self):
        t = Tag(type=ValueType.I)
        v = Val(data=0, text="0", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)

    def test_encode_empty_string(self):
        t = Tag(type=ValueType.Str)
        v = Val(data="", text="", tag=t)
        result = self.encoder.encode(v)
        self.assertIsInstance(result, bytes)


class TestTag(unittest.TestCase):
    def test_tag_basic(self):
        t = Tag(type=ValueType.Str, desc="test")
        self.assertEqual(t.type, ValueType.Str)
        self.assertEqual(t.desc, "test")

    def test_tag_bytes(self):
        t = Tag(type=ValueType.Str)
        b = t.bytes()
        self.assertIsInstance(b, bytes)

    def test_tag_inherit(self):
        parent = Tag(child_type=ValueType.I)
        child = Tag()
        child.inherit(parent)
        self.assertEqual(child.type, ValueType.I)


class TestToJSONC(unittest.TestCase):
    def test_to_jsonc_value_string(self):
        v = Val(data="hello", text="hello", tag=Tag(type=ValueType.Str))
        result = to_jsonc(v)
        self.assertEqual(result, '"hello"')

    def test_to_jsonc_value_int(self):
        v = Val(data=42, text="42", tag=Tag(type=ValueType.I))
        result = to_jsonc(v)
        self.assertEqual(result, "42")

    def test_to_jsonc_value_float(self):
        v = Val(data=3.14, text="3.14", tag=Tag(type=ValueType.F64))
        result = to_jsonc(v)
        self.assertEqual(result, "3.14")

    def test_to_jsonc_value_bool(self):
        v = Val(data=True, text="true", tag=Tag(type=ValueType.Bool))
        result = to_jsonc(v)
        self.assertEqual(result, "true")

    def test_to_jsonc_object(self):
        obj = Obj(
            fields=[
                Field(key="name", value=Val(data="John", text="John", tag=Tag(type=ValueType.Str))),
                Field(key="age", value=Val(data=30, text="30", tag=Tag(type=ValueType.I))),
            ],
            tag=Tag(name="person")
        )
        result = to_jsonc(obj)
        self.assertIn('"name"', result)
        self.assertIn('"age"', result)
        self.assertIn('{', result)
        self.assertIn('}', result)

    def test_to_jsonc_array(self):
        arr = Arr(
            items=[
                Val(data="a", text="a", tag=Tag(type=ValueType.Str)),
                Val(data="b", text="b", tag=Tag(type=ValueType.Str)),
            ],
            tag=Tag(type=ValueType.Arr)
        )
        result = to_jsonc(arr)
        self.assertIn('[', result)
        self.assertIn(']', result)

    def test_to_jsonc_object_with_tag(self):
        obj = Obj(
            fields=[
                Field(key="name", value=Val(data="John", text="John", tag=Tag(type=ValueType.Str, desc="user name"))),
            ],
            tag=Tag(name="person")
        )
        result = to_jsonc(obj)
        self.assertIn('// mm:', result)
        self.assertIn('desc=', result)


class TestEncoderDecoder(unittest.TestCase):
    def setUp(self):
        self.encoder = Encoder()

    def test_encode_decode_string(self):
        t = Tag(type=ValueType.Str)
        v = Val(data="hello", text="hello", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, "hello")

    def test_encode_decode_int(self):
        t = Tag(type=ValueType.I)
        v = Val(data=42, text="42", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, 42)

    def test_encode_decode_int_negative(self):
        t = Tag(type=ValueType.I)
        v = Val(data=-42, text="-42", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, -42)

    def test_encode_decode_float(self):
        t = Tag(type=ValueType.F64)
        v = Val(data=3.14, text="3.14", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertAlmostEqual(decoded, 3.14, places=2)

    def test_encode_decode_bool(self):
        # Test true
        t = Tag(type=ValueType.Bool)
        v = Val(data=True, text="true", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, True)
        
        # Test false
        v = Val(data=False, text="false", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, False)

    def test_encode_decode_bytes(self):
        t = Tag(type=ValueType.Bytes)
        v = Val(data=b"hello", text="hello", tag=t)
        encoded = self.encoder.encode(v)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertEqual(decoded, b"hello")

    def test_encode_decode_array(self):
        arr = Arr(
            items=[
                Val(data="a", text="a", tag=Tag(type=ValueType.Str)),
                Val(data="b", text="b", tag=Tag(type=ValueType.Str)),
            ],
            tag=Tag(type=ValueType.Arr, child_type=ValueType.Str)
        )
        encoded = self.encoder.encode(arr)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertIsInstance(decoded, list)
        self.assertEqual(len(decoded), 2)

    def test_encode_decode_object(self):
        obj = Obj(
            fields=[
                Field(
                    key="name",
                    value=Val(data="John", text="John", tag=Tag(type=ValueType.Str))
                ),
                Field(
                    key="age",
                    value=Val(data=30, text="30", tag=Tag(type=ValueType.I))
                ),
            ],
            tag=Tag(name="person")
        )
        encoded = self.encoder.encode(obj)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertIsInstance(decoded, dict)
        self.assertIn("name", decoded)
        self.assertIn("age", decoded)

    def test_encode_decode_nested_object(self):
        inner = Obj(
            fields=[
                Field(
                    key="city",
                    value=Val(data="Beijing", text="Beijing", tag=Tag(type=ValueType.Str))
                ),
            ],
            tag=Tag(name="address")
        )
        outer = Obj(
            fields=[
                Field(
                    key="name",
                    value=Val(data="John", text="John", tag=Tag(type=ValueType.Str))
                ),
                Field(
                    key="address",
                    value=inner
                ),
            ],
            tag=Tag(name="person")
        )
        encoded = self.encoder.encode(outer)
        decoder = Decoder(encoded)
        decoded = decoder.decode()
        self.assertIsInstance(decoded, dict)
        self.assertIn("name", decoded)
        self.assertIn("address", decoded)
        self.assertIsInstance(decoded["address"], dict)


if __name__ == "__main__":
    unittest.main()