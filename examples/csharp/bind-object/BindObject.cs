using MetaMessage;
using System.Collections.Generic;

namespace Examples.BindObject;

public class User
{
    public string Name { get; set; }
    public int Age { get; set; }
    public bool Active { get; set; }
    public List<int> Scores { get; set; }
}

public class BindObject
{
    public static void Main(string[] args)
    {
        // JSONC 字符串
        string jsonc = @"{
    // mm: type=str; desc=姓名
    ""name"": ""Alice"",
    // mm: type=i; desc=年龄
    ""age"": 25,
    // mm: type=bool; desc=是否激活
    ""active"": true,
    // mm: type=array; child_type=i; desc=分数
    ""scores"": [95, 87, 92]
}";

        Console.WriteLine("Input JSONC:");
        Console.WriteLine(jsonc);

        // 从 JSONC 绑定到 User 对象
        User user = Jsonc.BindFromString<User>(jsonc);

        Console.WriteLine("\nBound to object:");
        Console.WriteLine($"Name: {user.Name}");
        Console.WriteLine($"Age: {user.Age}");
        Console.WriteLine($"Active: {user.Active}");
        Console.WriteLine($"Scores: {string.Join(", ", user.Scores)}");
    }
}
