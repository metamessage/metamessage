using MetaMessage.Core;
using MetaMessage.Jsonc;

// Test 1: Direct bigint encoding
Console.WriteLine("=== Direct BigIntWireCodec ===");
byte[] raw = BigIntWireCodec.EncodeSignedDecimal("12345678901234567890");
Console.WriteLine($"  EncodeSignedDecimal: hex={BitConverter.ToString(raw).Replace("-","").ToLower()}, len={raw.Length}");

var enc = new WireEncoder();
enc.EncodeBigIntDecimal("12345678901234567890");
byte[] wireResult = enc.ToByteArray();
Console.WriteLine($"  EncodeBigIntDecimal: hex={BitConverter.ToString(wireResult).Replace("-","").ToLower()}, len={wireResult.Length}");

// Test 2: Full EncodeFromJsonc flow
Console.WriteLine("\n=== Full EncodeFromJsonc ===");
string jsonc = @"
{
  // mm: child_type=bigint
  ""bigint_arr"": [12345678901234567890, 98765432109876543210]
}";
byte[] fullResult = MetaMessage.Core.MetaMessage.EncodeFromJsonc(jsonc);
string fullHex = BitConverter.ToString(fullResult).Replace("-", "").ToLower();
Console.WriteLine($"  Full encode: hex={fullHex}, len={fullResult.Length}");

// Test 3: Encode with the individual element approach
Console.WriteLine("\n=== Element-by-element ===");
var node = MetaMessage.Jsonc.Jsonc.ParseFromString(jsonc);
if (node is JsoncDoc doc && doc.Fields.Count > 0 && doc.Fields[0].Value is JsoncArray arr)
{
    Console.WriteLine($"  Array has {arr.Elements.Count} elements");
    foreach (var elem in arr.Elements)
    {
        if (elem is JsoncValue val)
        {
            Console.WriteLine($"  Element tag type={val.Tag?.Type}, value={val.Value}, value type={val.Value?.GetType()}");
        }
    }
}

// Test 4: Compare with Go reference
string goHex = "ce28db8a626967696e745f617272fb029015de16aa140f6e462a062b3535a0aa147b747282315fad80a0";
Console.WriteLine($"\n=== Compare ===");
Console.WriteLine($"  Go:  {goHex}");
Console.WriteLine($"  C#:  {fullHex}");
Console.WriteLine($"  Match: {goHex == fullHex}");

if (goHex != fullHex)
{
    // Detailed byte-by-byte comparison
    int maxLen = Math.Max(goHex.Length, fullHex.Length);
    Console.WriteLine("\n  Byte-by-byte diff:");
    for (int i = 0; i < maxLen; i += 2)
    {
        string gb = i < goHex.Length ? goHex.Substring(i, Math.Min(2, goHex.Length - i)) : "";
        string cb = i < fullHex.Length ? fullHex.Substring(i, Math.Min(2, fullHex.Length - i)) : "";
        if (gb != cb)
        {
            Console.WriteLine($"    byte[{i/2}]: Go=0x{gb,-4} C#=0x{cb}");
        }
    }
}