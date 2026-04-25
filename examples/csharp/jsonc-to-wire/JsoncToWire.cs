using MetaMessage;

namespace Examples.JsoncToWire;

public class JsoncToWire
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

        Console.WriteLine("Input JSONC:");
        Console.WriteLine(jsonc);

        // 解析 JSONC
        var node = Jsonc.ParseFromString(jsonc);
        Console.WriteLine("\nParsed:");
        Console.WriteLine(Jsonc.ToString(node));

        // 编码到 Wire 格式
        byte[] wire = MetaMessage.MetaMessage.Encode(node);
        Console.WriteLine("\nEncoded Wire:");
        Console.WriteLine(BitConverter.ToString(wire).Replace("-", "").ToLower());
    }
}
