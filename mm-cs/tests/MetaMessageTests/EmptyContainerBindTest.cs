using System.Collections.Generic;
using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;
using Xunit;
using static MetaMessage.Core.MetaMessage;

namespace MetaMessageTests;

// 空容器 example 哨兵与 is_null 绑定行为测试（与 Go zzrepro_test.go 对齐）
public class ZzTarget
{
    public List<int> Vec { get; set; } = new();
    public int[] Arr { get; set; } = System.Array.Empty<int>();
    public string? Ptr { get; set; }
}

public class EmptyContainerBindTest
{
    private static NodeScalar Key(string name) => new NodeScalar(name, name, Tag.Empty());

    private static Tag ExampleIntTag()
    {
        var t = Tag.NewTag();
        t.Example = true;
        t.Type = ValueType.I;
        return t;
    }

    private static Tag PlainArrTag()
    {
        var t = Tag.NewTag();
        t.Type = ValueType.Arr;
        return t;
    }

    [Fact]
    public void EmptyVecWithOnlyExampleBindsToEmptyList()
    {
        // 数组仅含 example 哨兵项：应绑定为长度 0 的非 null List
        var exampleItem = new NodeScalar(0, "0", ExampleIntTag());
        var arr = new NodeArray(new List<INode> { exampleItem }, PlainArrTag());
        var root = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("Vec"), arr)
        }, Tag.NewTag());

        var decoded = Decode<ZzTarget>(EncodeTree(root));

        Assert.NotNull(decoded.Vec);
        Assert.Empty(decoded.Vec);
    }

    [Fact]
    public void EmptyArrayWithOnlyExampleBindsToEmptyArray()
    {
        // 数组仅含 example 哨兵项：应绑定为长度 0 的非 null 数组
        var exampleItem = new NodeScalar(0, "0", ExampleIntTag());
        var arr = new NodeArray(new List<INode> { exampleItem }, PlainArrTag());
        var root = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("Arr"), arr)
        }, Tag.NewTag());

        var decoded = Decode<ZzTarget>(EncodeTree(root));

        Assert.NotNull(decoded.Arr);
        Assert.Empty(decoded.Arr);
    }

    [Fact]
    public void ExampleItemsSkippedAmongRealItems()
    {
        // example 哨兵与真实项混存：只绑定真实项
        var realTag = Tag.NewTag();
        realTag.Type = ValueType.I;
        var items = new List<INode>
        {
            new NodeScalar(1, "1", realTag),
            new NodeScalar(0, "0", ExampleIntTag()),
            new NodeScalar(2, "2", realTag)
        };
        var arr = new NodeArray(items, PlainArrTag());
        var root = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("Vec"), arr)
        }, Tag.NewTag());

        var decoded = Decode<ZzTarget>(EncodeTree(root));

        Assert.Equal(new List<int> { 1, 2 }, decoded.Vec);
    }

    [Fact]
    public void IsNullScalarBindsToNull()
    {
        // 编码带 is_null：目标属性绑定为 null
        var nullTag = Tag.NewTag();
        nullTag.IsNull = true;
        nullTag.Nullable = true;
        nullTag.Type = ValueType.Str;
        var nullScalar = new NodeScalar(null, "null", nullTag);
        var root = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("Ptr"), nullScalar)
        }, Tag.NewTag());

        var before = new ZzTarget { Ptr = "preset" };
        Decode(EncodeTree(root), before);

        Assert.Null(before.Ptr);
    }

    [Fact]
    public void DecodeToValueSkipsExampleAndMapsNull()
    {
        // DecodeToValue 动态路径：example 哨兵跳过、is_null → null
        var exampleItem = new NodeScalar(0, "0", ExampleIntTag());
        var arr = new NodeArray(new List<INode> { exampleItem }, PlainArrTag());

        var mapTag = Tag.NewTag();
        mapTag.Type = ValueType.Map;
        var map = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("k"), exampleItem)
        }, mapTag);

        var nullTag = Tag.NewTag();
        nullTag.IsNull = true;
        nullTag.Nullable = true;
        nullTag.Type = ValueType.Str;
        var nullScalar = new NodeScalar(null, "null", nullTag);

        var root = new NodeObject(new List<KeyValuePair<NodeScalar, INode>>
        {
            new KeyValuePair<NodeScalar, INode>(Key("Vec"), arr),
            new KeyValuePair<NodeScalar, INode>(Key("M"), map),
            new KeyValuePair<NodeScalar, INode>(Key("Ptr"), nullScalar)
        }, Tag.NewTag());

        var value = DecodeToValue(EncodeTree(root)) as Dictionary<string, object?>;

        Assert.NotNull(value);
        var list = Assert.IsType<List<object?>>(value!["Vec"]);
        Assert.Empty(list);
        var dict = Assert.IsType<Dictionary<string, object?>>(value["M"]);
        Assert.Empty(dict);
        Assert.Null(value["Ptr"]);
    }
}
