// MetaMessage C# test harness - parse JSONC file and re-print to JSONC.
using MetaMessage.Jsonc;

if (args.Length < 1)
{
    Console.Error.WriteLine("usage: harness <file.jsonc>");
    Environment.Exit(1);
}

string input = "";
try
{
    input = File.ReadAllText(args[0]);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"read error: {ex.Message}");
    Environment.Exit(1);
}

IJsoncNode node = null!;
try
{
    node = Jsonc.ParseFromString(input);
}
catch (Exception ex)
{
    Console.Error.WriteLine($"parse error: {ex.Message}");
    Environment.Exit(1);
}

string output = Jsonc.ToString(node);
Console.Write(output);