namespace MetaMessage.Jsonc;

public static class Jsonc
{
    public static IJsoncNode ParseFromString(string input)
    {
        var parser = new JsoncParser(input);
        return parser.Parse();
    }

    public static IJsoncNode ParseFromBytes(byte[] input)
    {
        var jsoncString = System.Text.Encoding.UTF8.GetString(input);
        return ParseFromString(jsoncString);
    }

    public static T BindFromString<T>(string input) where T : new()
    {
        var node = ParseFromString(input);
        var binder = new JsoncBinder();
        return binder.Bind<T>(node);
    }

    public static T BindFromBytes<T>(byte[] input) where T : new()
    {
        var node = ParseFromBytes(input);
        var binder = new JsoncBinder();
        return binder.Bind<T>(node);
    }

    public static void BindFromString(string input, object target)
    {
        var node = ParseFromString(input);
        var binder = new JsoncBinder();
        binder.Bind(node, target);
    }

    public static void BindFromBytes(byte[] input, object target)
    {
        var node = ParseFromBytes(input);
        var binder = new JsoncBinder();
        binder.Bind(node, target);
    }

    public static string StructToJSONCString(object value)
    {
        var binder = new JsoncBinder();
        var node = binder.StructToNode(value);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static string ToString(IJsoncNode node)
    {
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static string ToMinString(IJsoncNode node)
    {
        var printer = new JsoncPrinter(prettyPrint: false);
        return printer.Print(node);
    }
}