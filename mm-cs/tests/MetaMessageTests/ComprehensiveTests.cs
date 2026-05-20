using Xunit;
using MetaMessage.Core;
using MetaMessage.Jsonc;
using static MetaMessage.Core.MetaMessage;
using JsoncParser = MetaMessage.Jsonc.Jsonc;
using MmVT = MetaMessage.Core.ValueType;

namespace MetaMessageTests;

public class ComprehensiveTests
{
    // ========== SimpleValue Tests ==========

    [Fact]
    public void TestSimpleValueNameOf_All32()
    {
        Assert.Equal("null_bool", SimpleValue.NameOf(SimpleValue.NULL_BOOL));
        Assert.Equal("null_int", SimpleValue.NameOf(SimpleValue.NULL_INT));
        Assert.Equal("null_float", SimpleValue.NameOf(SimpleValue.NULL_FLOAT));
        Assert.Equal("null_string", SimpleValue.NameOf(SimpleValue.NULL_STRING));
        Assert.Equal("null_bytes", SimpleValue.NameOf(SimpleValue.NULL_BYTES));
        Assert.Equal("false", SimpleValue.NameOf(SimpleValue.FALSE));
        Assert.Equal("true", SimpleValue.NameOf(SimpleValue.TRUE));
        Assert.Equal("code", SimpleValue.NameOf(SimpleValue.CODE));
        Assert.Equal("message", SimpleValue.NameOf(SimpleValue.MESSAGE));
        Assert.Equal("data", SimpleValue.NameOf(SimpleValue.DATA));
        Assert.Equal("success", SimpleValue.NameOf(SimpleValue.SUCCESS));
        Assert.Equal("error", SimpleValue.NameOf(SimpleValue.ERROR));
        Assert.Equal("unknown", SimpleValue.NameOf(SimpleValue.UNKNOWN));
        Assert.Equal("page", SimpleValue.NameOf(SimpleValue.PAGE));
        Assert.Equal("limit", SimpleValue.NameOf(SimpleValue.LIMIT));
        Assert.Equal("offset", SimpleValue.NameOf(SimpleValue.OFFSET));
        Assert.Equal("total", SimpleValue.NameOf(SimpleValue.TOTAL));
        Assert.Equal("id", SimpleValue.NameOf(SimpleValue.ID));
        Assert.Equal("name", SimpleValue.NameOf(SimpleValue.NAME));
        Assert.Equal("description", SimpleValue.NameOf(SimpleValue.DESCRIPTION));
        Assert.Equal("type", SimpleValue.NameOf(SimpleValue.TYPE));
        Assert.Equal("version", SimpleValue.NameOf(SimpleValue.VERSION));
        Assert.Equal("status", SimpleValue.NameOf(SimpleValue.STATUS));
        Assert.Equal("url", SimpleValue.NameOf(SimpleValue.URL));
        Assert.Equal("create_time", SimpleValue.NameOf(SimpleValue.CREATE_TIME));
        Assert.Equal("update_time", SimpleValue.NameOf(SimpleValue.UPDATE_TIME));
        Assert.Equal("delete_time", SimpleValue.NameOf(SimpleValue.DELETE_TIME));
        Assert.Equal("account", SimpleValue.NameOf(SimpleValue.ACCOUNT));
        Assert.Equal("token", SimpleValue.NameOf(SimpleValue.TOKEN));
        Assert.Equal("expire_time", SimpleValue.NameOf(SimpleValue.EXPIRE_TIME));
        Assert.Equal("key", SimpleValue.NameOf(SimpleValue.KEY));
        Assert.Equal("val", SimpleValue.NameOf(SimpleValue.VAL));
    }

    [Fact]
    public void TestSimpleValueNameToValue_Roundtrip()
    {
        Assert.Equal(SimpleValue.NULL_BOOL, SimpleValue.NameToValue("null_bool"));
        Assert.Equal(SimpleValue.NULL_INT, SimpleValue.NameToValue("null_int"));
        Assert.Equal(SimpleValue.NULL_FLOAT, SimpleValue.NameToValue("null_float"));
        Assert.Equal(SimpleValue.NULL_STRING, SimpleValue.NameToValue("null_string"));
        Assert.Equal(SimpleValue.NULL_BYTES, SimpleValue.NameToValue("null_bytes"));
        Assert.Equal(SimpleValue.FALSE, SimpleValue.NameToValue("false"));
        Assert.Equal(SimpleValue.TRUE, SimpleValue.NameToValue("true"));
        Assert.Equal(SimpleValue.CODE, SimpleValue.NameToValue("code"));
        Assert.Equal(SimpleValue.MESSAGE, SimpleValue.NameToValue("message"));
        Assert.Equal(SimpleValue.NAME, SimpleValue.NameToValue("name"));
        Assert.Equal(SimpleValue.KEY, SimpleValue.NameToValue("key"));
        Assert.Equal(SimpleValue.VAL, SimpleValue.NameToValue("val"));

        Assert.Null(SimpleValue.NameToValue("nonexistent"));
    }

    [Fact]
    public void TestSimpleValueNameToValue_CaseInsensitive()
    {
        Assert.Equal(SimpleValue.TRUE, SimpleValue.NameToValue("TRUE"));
        Assert.Equal(SimpleValue.TRUE, SimpleValue.NameToValue("True"));
        Assert.Equal(SimpleValue.CODE, SimpleValue.NameToValue("CODE"));
        Assert.Equal(SimpleValue.NAME, SimpleValue.NameToValue("NAME"));
    }

    // ========== Wire Format Tests ==========

    [Fact]
    public void TestWireEncodeDecode_SimpleBool()
    {
        var enc = new WireEncoder();
        enc.EncodeBool(true);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(true, result.Data);
        Assert.Equal(MmVT.BOOL, result.Tag.Type);
    }

    [Fact]
    public void TestWireEncodeDecode_SimpleNullBool()
    {
        var enc = new WireEncoder();
        enc.EncodeSimple(SimpleValue.NULL_BOOL);
        var inner = new WireEncoder();
        var tag = MmTag.Empty();
        tag.Type = MmVT.BOOL;
        tag.IsNull = true;
        inner.EncodeTaggedPayload(enc.ToByteArray(), tag.ToBytes());
        var dec = new WireDecoder(inner.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(false, result.Data);
        Assert.Equal("false", result.Text);
        Assert.True(result.Tag.IsNull);
        Assert.Equal(MmVT.BOOL, result.Tag.Type);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(1)]
    [InlineData(7)]
    [InlineData(42)]
    [InlineData(255)]
    [InlineData(256)]
    [InlineData(65535)]
    [InlineData(65536)]
    [InlineData(1234567890)]
    [InlineData(long.MaxValue / 2)]
    public void TestWireEncodeDecode_PositiveInt(long value)
    {
        var enc = new WireEncoder();
        enc.EncodeInt64(value);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(value, Convert.ToInt64(result.Data));
    }

    [Theory]
    [InlineData(-1)]
    [InlineData(-7)]
    [InlineData(-42)]
    [InlineData(-255)]
    [InlineData(-256)]
    [InlineData(-65535)]
    [InlineData(-65536)]
    [InlineData(-1234567890)]
    [InlineData(long.MinValue / 2)]
    public void TestWireEncodeDecode_NegativeInt(long value)
    {
        var enc = new WireEncoder();
        enc.EncodeInt64(value);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(value, Convert.ToInt64(result.Data));
    }

    [Theory]
    [InlineData("0.1")]
    [InlineData("0.5")]
    [InlineData("1.5")]
    [InlineData("3.14")]
    [InlineData("-3.14")]
    [InlineData("100.5")]
    [InlineData("0.001")]
    [InlineData("12345.6789")]
    public void TestWireEncodeDecode_Float(string value)
    {
        var enc = new WireEncoder();
        enc.EncodeFloatString(value);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        double expected = double.Parse(value);
        double actual = Convert.ToDouble(result.Data);
        Assert.Equal(expected, actual, 5);
    }

    [Theory]
    [InlineData("")]
    [InlineData("hello")]
    [InlineData("Hello, World!")]
    [InlineData("中文测试")]
    [InlineData("a very long string that exceeds the small string limit for testing purposes")]
    public void TestWireEncodeDecode_String(string value)
    {
        var enc = new WireEncoder();
        enc.EncodeString(value);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(value, result.Data as string);
    }

    [Fact]
    public void TestWireEncodeDecode_Bytes()
    {
        var data = new byte[] { 0, 1, 127, 128, 255 };
        var enc = new WireEncoder();
        enc.EncodeBytes(data);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(data, result.Data as byte[]);
    }

    [Fact]
    public void TestWireEncodeDecode_EmptyBytes()
    {
        var data = Array.Empty<byte>();
        var enc = new WireEncoder();
        enc.EncodeBytes(data);
        var dec = new WireDecoder(enc.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Empty(result.Data as byte[]);
    }

    // ========== BigInt Tests ==========

    [Theory]
    [InlineData("0")]
    [InlineData("1")]
    [InlineData("9")]
    [InlineData("10")]
    [InlineData("99")]
    [InlineData("100")]
    [InlineData("999")]
    [InlineData("1000")]
    [InlineData("1234567890")]
    [InlineData("999999999999999999")]
    [InlineData("-1")]
    [InlineData("-999")]
    [InlineData("-1234567890123456789")]
    public void TestBigInt_Roundtrip(string value)
    {
        byte[] encoded = BigIntWireCodec.EncodeSignedDecimal(value);
        string decoded = BigIntWireCodec.DecodeSignedDecimal(encoded);
        Assert.Equal(value, decoded);
    }

    [Fact]
    public void TestBigInt_Zero()
    {
        byte[] encoded = BigIntWireCodec.EncodeSignedDecimal("0");
        string decoded = BigIntWireCodec.DecodeSignedDecimal(encoded);
        Assert.Equal("0", decoded);
    }

    [Fact]
    public void TestBigInt_PositiveMax()
    {
        string value = "999999999999999999999999999999";
        byte[] encoded = BigIntWireCodec.EncodeSignedDecimal(value);
        string decoded = BigIntWireCodec.DecodeSignedDecimal(encoded);
        Assert.Equal(value, decoded);
    }

    // ========== Tagged Payload Tests ==========

    [Fact]
    public void TestWireEncodeDecode_TaggedInt()
    {
        var enc = new WireEncoder();
        var tag = MmTag.Empty();
        tag.Type = MmVT.INT;
        tag.Name = "count";
        enc.EncodeTaggedPayload(new WireEncoder().ToByteArray(), tag.ToBytes());

        int payloadEnc = new WireEncoder().EncodeInt64(42);
        var payloadBytes = new byte[payloadEnc];
        var tmp = new WireEncoder();
        tmp.EncodeInt64(42);
        Array.Copy(tmp.ToByteArray(), 0, payloadBytes, 0, payloadEnc);

        byte[] tagBytes = tag.ToBytes();
        var full = new WireEncoder();
        full.EncodeTaggedPayload(tmp.ToByteArray(), tagBytes);
        var dec = new WireDecoder(full.ToByteArray());
        var result = (MmScalar)dec.Decode();
        Assert.Equal(42L, Convert.ToInt64(result.Data));
    }

    // ========== MetaMessage Public API Tests ==========

    [Fact]
    public void TestEncodeTree_DecodeToTree_Scalar()
    {
        var scalar = new MmScalar("hello", "hello", new MmTag { Type = MmVT.STRING });
        byte[] encoded = EncodeTree(scalar);
        var decoded = DecodeToTree(encoded);
        var resultScalar = Assert.IsType<MmScalar>(decoded);
        Assert.Equal("hello", resultScalar.Data);
    }

    [Fact]
    public void TestEncodeTree_DecodeToTree_Array()
    {
        var children = new List<IMmTree>
        {
            new MmScalar(1L, "1", new MmTag { Type = MmVT.INT }),
            new MmScalar(2L, "2", new MmTag { Type = MmVT.INT }),
            new MmScalar(3L, "3", new MmTag { Type = MmVT.INT }),
        };
        var array = new MmArray(children, new MmTag { Type = MmVT.SLICE });
        byte[] encoded = EncodeTree(array);
        var decoded = DecodeToTree(encoded);
        var resultArray = Assert.IsType<MmArray>(decoded);
        Assert.Equal(3, resultArray.Children.Count);
    }

    [Fact]
    public void TestEncodeTree_DecodeToTree_Map()
    {
        var entries = new List<KeyValuePair<MmScalar, IMmTree>>
        {
            new(new MmScalar("name", "name", new MmTag { Type = MmVT.STRING }),
                new MmScalar("test", "test", new MmTag { Type = MmVT.STRING })),
        };
        var map = new MmMap(entries, new MmTag { Type = MmVT.MAP });
        byte[] encoded = EncodeTree(map);
        var decoded = DecodeToTree(encoded);
        var resultMap = Assert.IsType<MmMap>(decoded);
        Assert.Single(resultMap.Entries);
    }

    [Fact]
    public void TestEncodeTree_DecodeToTree_NullValues()
    {
        var scalar = new MmScalar(null, "null", new MmTag { Type = MmVT.INT, IsNull = true });
        byte[] encoded = EncodeTree(scalar);
        var decoded = DecodeToTree(encoded);
        var resultScalar = Assert.IsType<MmScalar>(decoded);
        Assert.Equal(0L, resultScalar.Data);
        Assert.Equal("0", resultScalar.Text);
        Assert.True(resultScalar.Tag.IsNull);
    }

    // ========== Jsonc Tests ==========

    [Fact]
    public void TestJsoncParse_SimpleTypes()
    {
        var node = JsoncParser.ParseFromString("true");
        Assert.IsType<JsoncValue>(node);
        var val = (JsoncValue)node;
        Assert.Equal(JsoncTokenType.True, val.TokenType);

        node = JsoncParser.ParseFromString("false");
        val = (JsoncValue)node;
        Assert.Equal(JsoncTokenType.False, val.TokenType);

        node = JsoncParser.ParseFromString("null");
        val = (JsoncValue)node;
        Assert.Equal(JsoncTokenType.Null, val.TokenType);

        node = JsoncParser.ParseFromString("42");
        val = (JsoncValue)node;
        Assert.Equal(JsoncTokenType.Number, val.TokenType);
        Assert.Equal(42.0, val.Value);

        node = JsoncParser.ParseFromString("\"hello\"");
        val = (JsoncValue)node;
        Assert.Equal(JsoncTokenType.String, val.TokenType);
        Assert.Equal("hello", val.Value);
    }

    [Fact]
    public void TestJsoncParse_Object()
    {
        var node = JsoncParser.ParseFromString("{\"name\": \"test\", \"value\": 42}");
        var obj = Assert.IsType<JsoncObject>(node);
        Assert.Equal(2, obj.Fields.Count);
        Assert.Equal("test", ((JsoncValue)obj.Fields["name"]).Value);
        Assert.Equal(42.0, ((JsoncValue)obj.Fields["value"]).Value);
    }

    [Fact]
    public void TestJsoncParse_Array()
    {
        var node = JsoncParser.ParseFromString("[1, 2, 3]");
        var arr = Assert.IsType<JsoncArray>(node);
        Assert.Equal(3, arr.Elements.Count);
        Assert.Equal(1.0, ((JsoncValue)arr.Elements[0]).Value);
        Assert.Equal(2.0, ((JsoncValue)arr.Elements[1]).Value);
        Assert.Equal(3.0, ((JsoncValue)arr.Elements[2]).Value);
    }

    [Fact]
    public void TestJsoncParse_Nested()
    {
        var node = JsoncParser.ParseFromString("{\"items\": [{\"id\": 1}, {\"id\": 2}]}");
        var obj = Assert.IsType<JsoncObject>(node);
        var items = Assert.IsType<JsoncArray>(obj.Fields["items"]);
        Assert.Equal(2, items.Elements.Count);
        var item0 = Assert.IsType<JsoncObject>(items.Elements[0]);
        Assert.Equal(1.0, ((JsoncValue)item0.Fields["id"]).Value);
    }

    [Fact]
    public void TestJsoncParse_WithComments()
    {
        var node = JsoncParser.ParseFromString(@"{
  // This is a leading comment
  ""name"": ""test"", // trailing comment
  /* block comment */
  ""value"": 42
}");
        var obj = Assert.IsType<JsoncObject>(node);
        Assert.Equal(2, obj.Fields.Count);
        Assert.Equal("test", ((JsoncValue)obj.Fields["name"]).Value);
        Assert.Equal(42.0, ((JsoncValue)obj.Fields["value"]).Value);
    }

    [Fact]
    public void TestJsoncPrint_Object()
    {
        var node = JsoncParser.ParseFromString("{\"name\": \"test\", \"value\": 42}");
        var printer = new JsoncPrinter(prettyPrint: true);
        var result = printer.Print(node);
        Assert.Contains("\"name\"", result);
        Assert.Contains("\"test\"", result);
        Assert.Contains("\"value\"", result);
        Assert.Contains("42", result);
    }

    [Fact]
    public void TestJsoncPrint_Array()
    {
        var node = JsoncParser.ParseFromString("[1, 2, 3]");
        var printer = new JsoncPrinter(prettyPrint: true);
        var result = printer.Print(node);
        Assert.Contains("1", result);
        Assert.Contains("2", result);
        Assert.Contains("3", result);
    }

    [Fact]
    public void TestJsoncPrint_Roundtrip()
    {
        var original = "{\n\t\"name\": \"test\",\n\t\"value\": 42\n}";
        var node = JsoncParser.ParseFromString(original);
        var printer = new JsoncPrinter(prettyPrint: true);
        var printed = printer.Print(node);
        var reparsed = JsoncParser.ParseFromString(printed);
        var obj = Assert.IsType<JsoncObject>(reparsed);
        Assert.Equal("test", ((JsoncValue)obj.Fields["name"]).Value);
        Assert.Equal(42.0, ((JsoncValue)obj.Fields["value"]).Value);
    }

    [Fact]
    public void TestJsoncExtractValue_Object()
    {
        var node = JsoncParser.ParseFromString("{\"name\": \"test\", \"value\": 42, \"flag\": true}");
        var result = JsoncParser.ExtractValue(node);
        var dict = Assert.IsType<Dictionary<string, object?>>(result);
        Assert.Equal("test", dict["name"]);
        Assert.Equal(42.0, dict["value"]);
        Assert.Equal(true, dict["flag"]);
    }

    [Fact]
    public void TestJsoncExtractValue_Array()
    {
        var node = JsoncParser.ParseFromString("[1, \"two\", true, null]");
        var result = JsoncParser.ExtractValue(node);
        var list = Assert.IsType<List<object?>>(result);
        Assert.Equal(4, list.Count);
        Assert.Equal(1.0, list[0]);
        Assert.Equal("two", list[1]);
        Assert.Equal(true, list[2]);
        Assert.Null(list[3]);
    }

    // ========== MetaMessage JSONC API Tests ==========

    [Fact]
    public void TestEncodeFromJsonc_DecodeToJsonc_Roundtrip()
    {
        var jsonc = "{\"name\": \"test\", \"count\": 42, \"active\": true}";
        byte[] encoded = EncodeFromJsonc(jsonc);
        Assert.NotNull(encoded);
        Assert.True(encoded.Length > 0);

        string decoded = DecodeToJsonc(encoded);
        Assert.NotNull(decoded);
        Assert.Contains("test", decoded);
        Assert.Contains("42", decoded);
        Assert.Contains("true", decoded);
    }

    [Fact]
    public void TestEncodeFromJsonc_Array()
    {
        var jsonc = "[1, 2, 3]";
        byte[] encoded = EncodeFromJsonc(jsonc);
        Assert.NotNull(encoded);
        Assert.True(encoded.Length > 0);
    }

    [Fact]
    public void TestValueToJsonc_SimpleObject()
    {
        var obj = new { Name = "test", Count = 42 };
        string result = ValueToJsonc(obj);
        Assert.Contains("test", result);
        Assert.Contains("42", result);
    }

    [Fact]
    public void TestValueToJsonc_AnonymousType()
    {
        var obj = new { id = 1, name = "hello", active = true };
        string result = ValueToJsonc(obj);
        Assert.Contains("1", result);
        Assert.Contains("hello", result);
        Assert.Contains("true", result);
    }

    [Fact]
    public void TestJsoncToValue_Simple()
    {
        var jsonc = "{\"name\": \"test\", \"count\": 42}";
        object? result = JsoncToValue(jsonc);
        var dict = Assert.IsType<Dictionary<string, object?>>(result);
        Assert.Equal("test", dict["name"]);
        Assert.Equal(42.0, dict["count"]);
    }

    [Fact]
    public void TestJsoncToValue_WithComments()
    {
        var jsonc = @"{
  // comment
  ""key"": ""value""
}";
        object? result = JsoncToValue(jsonc);
        var dict = Assert.IsType<Dictionary<string, object?>>(result);
        Assert.Equal("value", dict["key"]);
    }

    // ========== Edge Cases ==========

    [Fact]
    public void TestEncodeDecode_NegativeFloat()
    {
        var testObj = new BasicFloatObj { Value = -3.14f };
        var encoded = Encode(testObj);
        var decoded = Decode<BasicFloatObj>(encoded);
        Assert.Equal(testObj.Value, decoded.Value, 4);
    }

    [Fact]
    public void TestEncodeDecode_SmallFloat()
    {
        var testObj = new BasicFloatObj { Value = 0.1f };
        var encoded = Encode(testObj);
        var decoded = Decode<BasicFloatObj>(encoded);
        Assert.Equal(testObj.Value, decoded.Value, 4);
    }

    [Fact]
    public void TestEncodeDecode_VeryLongString()
    {
        var longStr = new string('x', 500);
        var testObj = new BasicStringObj { Value = longStr };
        var encoded = Encode(testObj);
        var decoded = Decode<BasicStringObj>(encoded);
        Assert.Equal(longStr, decoded.Value);
    }

    [Fact]
    public void TestEncodeDecode_SpecialCharacters()
    {
        var special = "Line1\nLine2\tTab\u0000Null";
        var testObj = new BasicStringObj { Value = special };
        var encoded = Encode(testObj);
        var decoded = Decode<BasicStringObj>(encoded);
        Assert.Equal(special, decoded.Value);
    }

    [Fact]
    public void TestEncodeDecode_DeeplyNested()
    {
        var obj = new DeepNested
        {
            Name = "level1",
            Child = new DeepNested
            {
                Name = "level2",
                Child = new DeepNested
                {
                    Name = "level3"
                }
            }
        };
        var encoded = Encode(obj);
        var decoded = Decode<DeepNested>(encoded);
        Assert.Equal("level1", decoded.Name);
        Assert.NotNull(decoded.Child);
        Assert.Equal("level2", decoded.Child.Name);
        Assert.NotNull(decoded.Child.Child);
        Assert.Equal("level3", decoded.Child.Child.Name);
    }

    [Fact]
    public void TestEncodeDecode_NestedArray()
    {
        var obj = new NestedArrayObj
        {
            Name = "matrix",
            Rows = new List<NestedArrayRow>
            {
                new() { Values = new int[] { 1, 2 } },
                new() { Values = new int[] { 3, 4 } },
            }
        };
        var encoded = Encode(obj);
        var decoded = Decode<NestedArrayObj>(encoded);
        Assert.Equal("matrix", decoded.Name);
        Assert.Equal(2, decoded.Rows.Count);
        Assert.Equal(2, decoded.Rows[0].Values.Length);
        Assert.Equal(1, decoded.Rows[0].Values[0]);
        Assert.Equal(2, decoded.Rows[0].Values[1]);
    }

    [Fact]
    public void TestEncodeDecode_AllIntegerTypes()
    {
        var obj = new AllIntObj
        {
            ByteVal = 255,
            SByteVal = -128,
            ShortVal = -32768,
            UShortVal = 65535,
            UIntVal = 4294967295,
            ULongVal = 18446744073709551615
        };
        var encoded = Encode(obj);
        var decoded = Decode<AllIntObj>(encoded);
        Assert.Equal(obj.ByteVal, decoded.ByteVal);
        Assert.Equal(obj.SByteVal, decoded.SByteVal);
        Assert.Equal(obj.ShortVal, decoded.ShortVal);
        Assert.Equal(obj.UShortVal, decoded.UShortVal);
        Assert.Equal(obj.UIntVal, decoded.UIntVal);
        Assert.Equal(obj.ULongVal, decoded.ULongVal);
    }

    [Fact]
    public void TestEncodeDecode_EmptyStruct()
    {
        var obj = new EmptyStruct();
        var encoded = Encode(obj);
        var decoded = Decode<EmptyStruct>(encoded);
        Assert.NotNull(decoded);
    }

    [Fact]
    public void TestValidate_RequiredField()
    {
        var tag = MmTag.Empty();
        tag.Type = MmVT.STRING;
        tag.Nullable = false;

        var validResult = Validator.Validate("hello", tag);
        Assert.True(validResult.IsValid);

        var nullResult = Validator.Validate(null, tag);
        Assert.False(nullResult.IsValid);
    }

    [Fact]
    public void TestDecoder_EmptyData_Throws()
    {
        Assert.Throws<MmDecodeException>(() =>
        {
            var dec = new WireDecoder(Array.Empty<byte>());
            dec.Decode();
        });
    }

    [Fact]
    public void TestEncodeDecode_MixedNullableState()
    {
        var obj = new MixedNullObj
        {
            AlwaysPresent = "hello",
            SometimesNull = null,
            SometimesPresent = "world"
        };
        var encoded = Encode(obj);
        var decoded = Decode<MixedNullObj>(encoded);
        Assert.Equal("hello", decoded.AlwaysPresent);
        Assert.Null(decoded.SometimesNull);
        Assert.Equal("world", decoded.SometimesPresent);
    }

    // ========== Test Model Classes ==========

    public class BasicFloatObj
    {
        public float Value { get; set; }
    }

    public class BasicStringObj
    {
        public string Value { get; set; } = string.Empty;
    }

    public class DeepNested
    {
        public string Name { get; set; } = string.Empty;
        public DeepNested? Child { get; set; }
    }

    public class NestedArrayRow
    {
        public int[] Values { get; set; } = Array.Empty<int>();
    }

    public class NestedArrayObj
    {
        public string Name { get; set; } = string.Empty;
        public List<NestedArrayRow> Rows { get; set; } = new();
    }

    public class AllIntObj
    {
        public byte ByteVal { get; set; }
        public sbyte SByteVal { get; set; }
        public short ShortVal { get; set; }
        public ushort UShortVal { get; set; }
        public uint UIntVal { get; set; }
        public ulong ULongVal { get; set; }
    }

    public class EmptyStruct { }

    public class MixedNullObj
    {
        public string AlwaysPresent { get; set; } = string.Empty;
        public string? SometimesNull { get; set; }
        public string? SometimesPresent { get; set; }
    }
}