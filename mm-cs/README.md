# C# MetaMessage 库使用说明

## 1. 安装

### NuGet 依赖
使用 NuGet 安装：

```bash
# 使用 .NET CLI
dotnet add package MetaMessage

# 或使用 Package Manager Console
Install-Package MetaMessage
```

### 版本要求
- .NET 6.0 或更高版本

## 2. 基本使用

### 2.1 类定义

```csharp
using MetaMessage;

public class Person
{
    public string Name { get; set; } = "Ed";
    public int Age { get; set; } = 30;
}
```

### 2.2 编码示例

```csharp
var person = new Person();
byte[] wire = MetaMessage.MetaMessage.Encode(person);
Console.WriteLine($"Encoded: {BitConverter.ToString(wire).Replace("-", "").ToLower()}");
```

### 2.3 解码示例

```csharp
var decoded = MetaMessage.MetaMessage.Decode<Person>(wire);
Console.WriteLine($"Decoded: Name={decoded.Name}, Age={decoded.Age}");
```

### 2.4 JSONC 解析示例

```csharp
using MetaMessage.Jsonc;

string jsonc = @"{
    // mm: type=str; desc=姓名
    ""name"": ""Alice"",
    // mm: type=i; desc=年龄
    ""age"": 25
}";

// 解析 JSONC
var node = Jsonc.ParseFromString(jsonc);

// 绑定到对象
var person = Jsonc.BindFromString<Person>(jsonc);
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-cs 目录下运行
dotnet test
```

### 3.2 测试框架
- xUnit
- .NET Test SDK

### 3.3 测试覆盖范围
- 编码测试
- 解码测试
- JSONC 解析测试
- 绑定测试

## 4. 常见问题

### 4.1 依赖问题
- **问题**: NuGet 包安装失败
  **解决**: 检查网络连接，或使用 NuGet 镜像

### 4.2 编译问题
- **问题**: 找不到命名空间
  **解决**: 确保包引用正确，并且项目已重新生成

### 4.3 运行时问题
- **问题**: 编码/解码失败
  **解决**: 检查类定义是否正确，属性是否可访问

## 5. 示例代码

查看 `examples/csharp/` 目录下的示例代码：
- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式
- `wire-to-jsonc/` - Wire 格式转 JSONC
- `bind-object/` - 对象绑定示例

## 6. 相关资源

- [.NET 文档](https://docs.microsoft.com/en-us/dotnet/)
- [xUnit 文档](https://xunit.net/docs/)
- [NuGet 文档](https://docs.microsoft.com/en-us/nuget/)
