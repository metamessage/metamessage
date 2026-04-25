using MetaMessage;

namespace Examples.WireToJsonc;

public class WireToJsonc
{
    public static void Main(string[] args)
    {
        // JSONC 字符串
        string jsonc = @"{
    // mm: type=datetime; desc=创建时间
    ""create_time"": ""2026-01-01 00:00:00"",
    // mm: type=str; desc=用户名称
    ""user_name"": ""Alice"",
    // mm: type=bool; desc=是否激活
    ""is_active"": true,
    // mm: type=array; child_type=i
    ""scores"": [95, 87, 92]
}";

        Console.WriteLine("Original JSONC:");
        Console.WriteLine(jsonc);

        // 解析并编码
        var node = Jsonc.ParseFromString(jsonc);
        byte[] wire = MetaMessage.MetaMessage.Encode(node);

        Console.WriteLine("\nEncoded Wire:");
        Console.WriteLine(BitConverter.ToString(wire).Replace("-", "").ToLower());

        // 从 Wire 解码到 JSONC
        var decodedNode = MetaMessage.MetaMessage.Decode<JsoncNode>(wire);
        string outputJsonc = Jsonc.ToString(decodedNode);

        Console.WriteLine("\nDecoded to JSONC:");
        Console.WriteLine(outputJsonc);
    }
}
