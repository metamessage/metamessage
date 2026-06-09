// MetaMessage C# test harness - parse JSONC file and re-print to JSONC.
using MetaMessage.Jsonc;
using MetaMessage.Core;
using MetaMessage.Ir;
using MM = MetaMessage.Core.MetaMessage;

if (args.Length < 1)
{
    Console.Error.WriteLine("usage: harness [--encode|--decode|--debug] <file.jsonc>");
    Environment.Exit(1);
}

if (args[0] == "--debug")
{
    if (args.Length < 2)
    {
        Console.Error.WriteLine("usage: harness --debug <file.jsonc>");
        Environment.Exit(1);
    }
    string inputText = "";
    try
    {
        inputText = File.ReadAllText(args[1]);
    }
    catch (Exception ex)
    {
        Console.Error.WriteLine($"read error: {ex.Message}");
        Environment.Exit(1);
    }
    try
    {
        var tree = Jsonc.ParseFromString(inputText);
        DumpNode(tree, 0);
    }
    catch (Exception ex)
    {
        Console.Error.WriteLine($"debug error: {ex.Message}");
        Environment.Exit(1);
    }
    return;
}

if (args[0] == "--encode")
{
    if (args.Length < 2)
    {
        Console.Error.WriteLine("usage: harness --encode <file.jsonc>");
        Environment.Exit(1);
    }
    string inputText = "";
    try
    {
        inputText = File.ReadAllText(args[1]);
    }
    catch (Exception ex)
    {
        Console.Error.WriteLine($"read error: {ex.Message}");
        Environment.Exit(1);
    }
    try
    {
        byte[] wire = MM.EncodeFromJsonc(inputText);
        Console.Write(Convert.ToHexString(wire).ToLower());
    }
    catch (Exception ex)
    {
        Console.Error.WriteLine($"encode error: {ex.Message}");
        Environment.Exit(1);
    }
    return;
}

if (args[0] == "--decode")
{
    string hexStr = Console.In.ReadToEnd().Trim();
    try
    {
        byte[] wire = Convert.FromHexString(hexStr);
        string decoded = MM.DecodeToJsonc(wire);
        Console.Write(decoded);
    }
    catch (Exception ex)
    {
        Console.Error.WriteLine($"decode error: {ex.Message}");
        Environment.Exit(1);
    }
    return;
}

// Existing behavior
string fileContent = "";
try
{
    fileContent = File.ReadAllText(args[0]);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"read error: {ex.Message}");
    Environment.Exit(1);
}

INode node = null!;
try
{
    node = Jsonc.ParseFromString(fileContent);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"parse error: {ex.Message}");
    Environment.Exit(1);
}

string jsoncOutput = Jsonc.ToString(node);
Console.Write(jsoncOutput);

static void DumpNode(INode node, int indent)
{
    var prefix = new string(' ', indent * 2);
    switch (node)
    {
        case NodeScalar scalar:
            var tagBytes = scalar.Tag?.ToBytes() ?? Array.Empty<byte>();
            Console.Error.WriteLine($"{prefix}Scalar: text='{scalar.Text}' data='{scalar.Data}' tag.Type={scalar.Tag?.Type} tag.Desc='{scalar.Tag?.Desc}' tag.ToBytes().Length={tagBytes.Length} tag.ToBytes()={Convert.ToHexString(tagBytes).ToLower()}");
            break;
        case NodeArray arr:
            Console.Error.WriteLine($"{prefix}Array: count={arr.Children.Count} tag.Type={arr.Tag?.Type} tag.ToBytes().Length={(arr.Tag?.ToBytes() ?? Array.Empty<byte>()).Length}");
            for (int i = 0; i < arr.Children.Count; i++)
            {
                Console.Error.WriteLine($"{prefix}  [{i}]:");
                DumpNode(arr.Children[i], indent + 2);
            }
            break;
        case NodeObject obj:
            Console.Error.WriteLine($"{prefix}Object: count={obj.Entries.Count} tag.Type={obj.Tag?.Type} tag.ToBytes().Length={(obj.Tag?.ToBytes() ?? Array.Empty<byte>()).Length}");
            foreach (var kvp in obj.Entries)
            {
                Console.Error.WriteLine($"{prefix}  Key='{kvp.Key.Text}':");
                DumpNode(kvp.Value, indent + 2);
            }
            break;
    }
}