using MetaMessage.Core;
using MetaMessage.Jsonc;
using JsoncParser = MetaMessage.Jsonc.Jsonc;

Console.OutputEncoding = System.Text.Encoding.UTF8;

string separator = new string('-', 50);

// 示例 1: JSONC 解析与输出
Console.WriteLine(separator);
Console.WriteLine("示例 1: JSONC 解析与输出");
Console.WriteLine(separator);

string jsonc = @"{
    // mm: type=str; desc=姓名
    ""name"": ""Alice"",
    // mm: type=i; desc=年龄; min=0; max=150
    ""age"": 25,
    // mm: type=bool; desc=是否激活
    ""active"": true,
    // mm: type=array; child_type=i; desc=分数
    ""scores"": [95, 87, 92]
}";

Console.WriteLine("输入 JSONC:");
Console.WriteLine(jsonc);

var node = JsoncParser.ParseFromString(jsonc);
Console.WriteLine();
Console.WriteLine("格式化输出:");
Console.WriteLine(JsoncParser.ToString(node));
Console.WriteLine();
Console.WriteLine("紧凑输出:");
Console.WriteLine(JsoncParser.ToMinString(node));

// 示例 2: JSONC 绑定到对象
Console.WriteLine();
Console.WriteLine(separator);
Console.WriteLine("示例 2: JSONC 绑定到对象");
Console.WriteLine(separator);

var user = JsoncParser.BindFromString<User>(jsonc);
Console.WriteLine($"name: {user.name}");
Console.WriteLine($"age: {user.age}");
Console.WriteLine($"active: {user.active}");
Console.WriteLine($"scores: {string.Join(", ", user.scores)}");

// 示例 3: 对象转 JSONC
Console.WriteLine();
Console.WriteLine(separator);
Console.WriteLine("示例 3: 对象转 JSONC");
Console.WriteLine(separator);

var userObj = new User { name = "Bob", age = 28, active = true, scores = new List<int> { 100, 90, 80 } };
string outputJsonc = JsoncParser.ValueToNodeString(userObj);
Console.WriteLine(outputJsonc);

// 示例 4: JSONC 值提取
Console.WriteLine();
Console.WriteLine(separator);
Console.WriteLine("示例 4: JSONC 值提取");
Console.WriteLine(separator);

object? extracted = JsoncParser.ExtractValue(node);
Console.WriteLine($"提取值: {System.Text.Json.JsonSerializer.Serialize(extracted)}");

// 示例 5: JSONC ↔ Wire 格式往返
Console.WriteLine();
Console.WriteLine(separator);
Console.WriteLine("示例 5: JSONC ↔ Wire 格式往返");
Console.WriteLine(separator);

string simpleJsonc = @"{
    // mm: type=str; desc=姓名
    ""name"": ""Alice"",
    // mm: type=i; desc=年龄; min=0; max=150
    ""age"": 30
}";

Console.WriteLine("输入:");
Console.WriteLine(simpleJsonc);

var wireBytes = MetaMessage.Core.MetaMessage.EncodeFromJsonc(simpleJsonc);
Console.WriteLine($"Wire 编码: {BitConverter.ToString(wireBytes).Replace("-", "").ToLower()}");

string decodedJsonc = MetaMessage.Core.MetaMessage.DecodeToJsonc(wireBytes);
Console.WriteLine("解码 JSONC:");
Console.WriteLine(decodedJsonc);

// 示例 6: 从字节数组解析 JSONC
Console.WriteLine();
Console.WriteLine(separator);
Console.WriteLine("示例 6: 从字节数组解析 JSONC");
Console.WriteLine(separator);

byte[] bytes = System.Text.Encoding.UTF8.GetBytes(jsonc);
var nodeFromBytes = JsoncParser.ParseFromBytes(bytes);
Console.WriteLine(JsoncParser.ToString(nodeFromBytes));

Console.WriteLine();
Console.WriteLine("所有示例运行完成!");

public class User
{
    public string name { get; set; } = "";
    public int age { get; set; }
    public bool active { get; set; }
    public List<int> scores { get; set; } = new();
}