# MetaMessage

MetaMessage 的 C# 实现，提供高性能的二进制序列化、JSONC 解析和对象绑定功能。

## 1. 安装

### 1.1 项目引用

将 MetaMessage 添加到您的项目中：

```bash
# 使用 .NET CLI
dotnet add reference ../../mm-cs/src/MetaMessage/MetaMessage.csproj
```

### 1.2 版本要求

- .NET 8.0 或更高版本

## 2. API 使用

### 2.1 JSONC 解析与输出

```csharp
using MetaMessage.Jsonc;
using JsoncParser = MetaMessage.Jsonc.Jsonc;

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

// 解析 JSONC 到节点树
var node = JsoncParser.ParseFromString(jsonc);

// 格式化输出
Console.WriteLine(JsoncParser.ToString(node));

// 紧凑输出
Console.WriteLine(JsoncParser.ToMinString(node));

// 从字节数组解析
byte[] bytes = Encoding.UTF8.GetBytes(jsonc);
var nodeFromBytes = JsoncParser.ParseFromBytes(bytes);
```

### 2.2 JSONC 绑定到对象

```csharp
public class User
{
    public string name { get; set; } = "";
    public int age { get; set; }
    public bool active { get; set; }
    public List<int> scores { get; set; } = new();
}

// 绑定 JSONC 到强类型对象
var user = JsoncParser.BindFromString<User>(jsonc);
Console.WriteLine($"name: {user.name}, age: {user.age}");

// 绑定到已有对象
var target = new User();
JsoncParser.BindFromString(jsonc, target);
```

### 2.3 对象转 JSONC

```csharp
var user = new User
{
    name = "Bob",
    age = 28,
    active = true,
    scores = new List<int> { 100, 90, 80 }
};

string output = JsoncParser.ValueToNodeString(user);
Console.WriteLine(output);
```

### 2.4 JSONC 值提取

```csharp
object? extracted = JsoncParser.ExtractValue(node);
// 返回纯 JSON 对象（字典、列表、原始值）
```

## 3. 完整示例

查看 `examples/csharp/` 目录下的完整运行示例：

```bash
cd examples/csharp
dotnet run
```

示例包含以下功能：

- JSONC 解析与格式化输出
- JSONC 绑定到强类型对象
- 对象转 JSONC
- JSONC 值提取
- 从字节数组解析 JSONC

## 4. 测试方法

```bash
cd mm-cs
dotnet test
```

测试框架：xUnit + .NET Test SDK

## 5. 相关资源

- [Go 实现](https://github.com/example/meta-message)
