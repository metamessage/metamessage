using MetaMessage;

namespace Examples.Basic;

public class Person
{
    public string Name { get; set; } = "Ed";
    public int Age { get; set; } = 30;
}

public class Basic
{
    public static void Main(string[] args)
    {
        // 创建 Person 对象
        Person person = new Person();
        Console.WriteLine($"Original: Name={person.Name}, Age={person.Age}");

        // 编码到 Wire 格式
        byte[] wire = MetaMessage.MetaMessage.Encode(person);
        Console.WriteLine($"Encoded: {BitConverter.ToString(wire).Replace("-", "").ToLower()}");

        // 从 Wire 解码
        Person decoded = MetaMessage.MetaMessage.Decode<Person>(wire);
        Console.WriteLine($"Decoded: Name={decoded.Name}, Age={decoded.Age}");
    }
}
