using Xunit;
using static MetaMessage.Mm.MetaMessage;

namespace MetaMessageTests;

public class MetaMessageTest
{
    [Fact]
    public void TestEncodeDecodeBasicTypes()
    {
        // 测试基本类型
        var testObj = new TestBasicTypes
        {
            BoolValue = true,
            IntValue = 42,
            LongValue = 1234567890,
            FloatValue = 3.14f,
            DoubleValue = 2.71828,
            StringValue = "Hello, World!"
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestBasicTypes>(encoded);

        Assert.Equal(testObj.BoolValue, decoded.BoolValue);
        Assert.Equal(testObj.IntValue, decoded.IntValue);
        Assert.Equal(testObj.LongValue, decoded.LongValue);
        Assert.Equal(testObj.FloatValue, decoded.FloatValue, 5);
        Assert.Equal(testObj.DoubleValue, decoded.DoubleValue, 5);
        Assert.Equal(testObj.StringValue, decoded.StringValue);
    }

    [Fact]
    public void TestEncodeDecodeNullableTypes()
    {
        // 测试可空类型
        var testObj = new TestNullableTypes
        {
            NullableBool = null,
            NullableInt = null,
            NullableString = null
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestNullableTypes>(encoded);

        Assert.Null(decoded.NullableBool);
        Assert.Null(decoded.NullableInt);
        Assert.Null(decoded.NullableString);
    }

    [Fact]
    public void TestEncodeDecodeArrayTypes()
    {
        // 测试数组类型
        var testObj = new TestArrayTypes
        {
            IntArray = new int[] { 1, 2, 3, 4, 5 },
            StringList = new List<string> { "a", "b", "c" }
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestArrayTypes>(encoded);

        Assert.Equal(testObj.IntArray.Length, decoded.IntArray.Length);
        for (int i = 0; i < testObj.IntArray.Length; i++)
        {
            Assert.Equal(testObj.IntArray[i], decoded.IntArray[i]);
        }

        Assert.Equal(testObj.StringList.Count, decoded.StringList.Count);
        for (int i = 0; i < testObj.StringList.Count; i++)
        {
            Assert.Equal(testObj.StringList[i], decoded.StringList[i]);
        }
    }

    [Fact]
    public void TestEncodeDecodeStructTypes()
    {
        // 测试结构体类型
        var testObj = new TestStructTypes
        {
            Name = "Test",
            Nested = new NestedStruct
            {
                Value = 42
            }
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestStructTypes>(encoded);

        Assert.Equal(testObj.Name, decoded.Name);
        Assert.NotNull(decoded.Nested);
        Assert.Equal(testObj.Nested.Value, decoded.Nested.Value);
    }

    [Fact]
    public void TestEncodeDecodeDateTimeTypes()
    {
        // 测试日期时间类型
        var testObj = new TestDateTimeTypes
        {
            DateTimeValue = DateTime.Now
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestDateTimeTypes>(encoded);

        // 允许一定的时间误差
        Assert.InRange(decoded.DateTimeValue, testObj.DateTimeValue.AddSeconds(-1), testObj.DateTimeValue.AddSeconds(1));
    }

    [Fact]
    public void TestEncodeDecodeEnumTypes()
    {
        // 测试枚举类型
        var testObj = new TestEnumTypes
        {
            EnumValue = TestEnum.Value2
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestEnumTypes>(encoded);

        Assert.Equal(testObj.EnumValue, decoded.EnumValue);
    }

    [Fact]
    public void TestEncodeDecodeBytesTypes()
    {
        // 测试字节数组类型
        var testObj = new TestBytesTypes
        {
            BytesValue = new byte[] { 1, 2, 3, 4, 5 }
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestBytesTypes>(encoded);

        Assert.Equal(testObj.BytesValue.Length, decoded.BytesValue.Length);
        for (int i = 0; i < testObj.BytesValue.Length; i++)
        {
            Assert.Equal(testObj.BytesValue[i], decoded.BytesValue[i]);
        }
    }

    [Fact]
    public void TestRoundtripAllTypes()
    {
        // 测试所有类型的往返编码解码
        var testObj = new TestAllTypes
        {
            BoolValue = true,
            IntValue = 42,
            LongValue = 1234567890,
            FloatValue = 3.14f,
            DoubleValue = 2.71828,
            StringValue = "Hello, World!",
            NullableBool = true,
            NullableInt = 100,
            NullableString = "Nullable",
            IntArray = new int[] { 1, 2, 3 },
            StringList = new List<string> { "a", "b" },
            DateTimeValue = DateTime.Now,
            EnumValue = TestEnum.Value1,
            BytesValue = new byte[] { 1, 2, 3 },
            Nested = new NestedStruct { Value = 999 }
        };

        var encoded = Encode(testObj);
        var decoded = Decode<TestAllTypes>(encoded);

        Assert.Equal(testObj.BoolValue, decoded.BoolValue);
        Assert.Equal(testObj.IntValue, decoded.IntValue);
        Assert.Equal(testObj.LongValue, decoded.LongValue);
        Assert.Equal(testObj.FloatValue, decoded.FloatValue, 5);
        Assert.Equal(testObj.DoubleValue, decoded.DoubleValue, 5);
        Assert.Equal(testObj.StringValue, decoded.StringValue);
        Assert.Equal(testObj.NullableBool, decoded.NullableBool);
        Assert.Equal(testObj.NullableInt, decoded.NullableInt);
        Assert.Equal(testObj.NullableString, decoded.NullableString);
        Assert.Equal(testObj.IntArray.Length, decoded.IntArray.Length);
        Assert.Equal(testObj.StringList.Count, decoded.StringList.Count);
        Assert.InRange(decoded.DateTimeValue, testObj.DateTimeValue.AddSeconds(-1), testObj.DateTimeValue.AddSeconds(1));
        Assert.Equal(testObj.EnumValue, decoded.EnumValue);
        Assert.Equal(testObj.BytesValue.Length, decoded.BytesValue.Length);
        Assert.NotNull(decoded.Nested);
        Assert.Equal(testObj.Nested.Value, decoded.Nested.Value);
    }

    // 测试类定义
    public class TestBasicTypes
    {
        public bool BoolValue { get; set; }
        public int IntValue { get; set; }
        public long LongValue { get; set; }
        public float FloatValue { get; set; }
        public double DoubleValue { get; set; }
        public string StringValue { get; set; } = string.Empty;
    }

    public class TestNullableTypes
    {
        public bool? NullableBool { get; set; }
        public int? NullableInt { get; set; }
        public string? NullableString { get; set; }
    }

    public class TestArrayTypes
    {
        public int[] IntArray { get; set; } = Array.Empty<int>();
        public List<string> StringList { get; set; } = new List<string>();
    }

    public class NestedStruct
    {
        public int Value { get; set; }
    }

    public class TestStructTypes
    {
        public string Name { get; set; } = string.Empty;
        public NestedStruct Nested { get; set; } = new NestedStruct();
    }

    public class TestDateTimeTypes
    {
        public DateTime DateTimeValue { get; set; }
    }

    public enum TestEnum
    {
        Value1,
        Value2,
        Value3
    }

    public class TestEnumTypes
    {
        public TestEnum EnumValue { get; set; }
    }

    public class TestBytesTypes
    {
        public byte[] BytesValue { get; set; } = Array.Empty<byte>();
    }

    public class TestAllTypes
    {
        public bool BoolValue { get; set; }
        public int IntValue { get; set; }
        public long LongValue { get; set; }
        public float FloatValue { get; set; }
        public double DoubleValue { get; set; }
        public string StringValue { get; set; } = string.Empty;
        public bool? NullableBool { get; set; }
        public int? NullableInt { get; set; }
        public string? NullableString { get; set; }
        public int[] IntArray { get; set; } = Array.Empty<int>();
        public List<string> StringList { get; set; } = new List<string>();
        public DateTime DateTimeValue { get; set; }
        public TestEnum EnumValue { get; set; }
        public byte[] BytesValue { get; set; } = Array.Empty<byte>();
        public NestedStruct Nested { get; set; } = new NestedStruct();
    }
}