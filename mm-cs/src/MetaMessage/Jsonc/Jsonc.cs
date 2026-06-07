using MetaMessage.Ir;

namespace MetaMessage.Jsonc;

public static class Jsonc
{
    private static readonly JsoncPrinter DefaultPrinter = new(prettyPrint: true);

    public static INode ParseFromString(string input)
    {
        var parser = new JsoncParser(input);
        return parser.Parse();
    }

    public static INode ParseFromBytes(byte[] input)
    {
        var jsoncString = System.Text.Encoding.UTF8.GetString(input);
        return ParseFromString(jsoncString);
    }

    public static INode ParseFromJSONC(string input)
    {
        return ParseFromString(input);
    }

    public static string ToJSONC(INode tree)
    {
        if (tree == null)
            return "";
        return DefaultPrinter.Print(tree);
    }

    public static string ToString(INode node)
    {
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static string ToMinString(INode node)
    {
        var printer = new JsoncPrinter(prettyPrint: false);
        return printer.Print(node);
    }

    public static object? ExtractValue(INode node)
    {
        switch (node)
        {
            case NodeScalar scalar:
                return scalar.Data;
            case NodeArray array:
                {
                    var list = new List<object?>();
                    foreach (var element in array.Children)
                    {
                        list.Add(ExtractValue(element));
                    }
                    return list;
                }
            case NodeObject map:
                {
                    var dict = new Dictionary<string, object?>();
                    foreach (var kvp in map.Entries)
                    {
                        dict[kvp.Key.Text] = ExtractValue(kvp.Value);
                    }
                    return dict;
                }
            default:
                return null;
        }
    }
}