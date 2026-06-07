// MetaMessage C# test harness - parse JSONC file and re-print to JSONC.
using MetaMessage.Jsonc;
using MetaMessage.Core;
using MetaMessage.Ir;
using MM = MetaMessage.Core.MetaMessage;

if (args.Length < 1)
{
    Console.Error.WriteLine("usage: harness [--encode|--decode] <file.jsonc>");
    Environment.Exit(1);
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