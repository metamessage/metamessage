using MetaMessage.Ir;

namespace MetaMessage.Jsonc;

public static class Jsonc
{
    private static readonly JsoncPrinter DefaultPrinter = new(prettyPrint: true);

    public static IMmTree ParseFromString(string input)
    {
        var parser = new JsoncParser(input);
        return parser.Parse();
    }

    public static IMmTree ParseFromBytes(byte[] input)
    {
        var jsoncString = System.Text.Encoding.UTF8.GetString(input);
        return ParseFromString(jsoncString);
    }

    public static IMmTree ParseFromJSONC(string input)
    {
        return ParseFromString(input);
    }

    public static string ToJSONC(IMmTree tree)
    {
        if (tree == null)
            return "";
        return DefaultPrinter.Print(tree);
    }

    public static string ToString(IMmTree node)
    {
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static string ToMinString(IMmTree node)
    {
        var printer = new JsoncPrinter(prettyPrint: false);
        return printer.Print(node);
    }

    public static object? ExtractValue(IMmTree node)
    {
        switch (node)
        {
            case MmScalar scalar:
                return scalar.Data;
            case MmArray array:
                {
                    var list = new List<object?>();
                    foreach (var element in array.Children)
                    {
                        list.Add(ExtractValue(element));
                    }
                    return list;
                }
            case MmMap map:
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