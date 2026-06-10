# MetaMessage

MetaMessage (mm) is a structured data exchange protocol. It is self-describing, self-constraining, and self-exemplifying, enabling lossless data exchange. It is designed as a next-generation universal protocol that natively supports AI, humans, and machines.

- Human and AI friendly
- Export/import to JSONC (currently; YAML/TOML support planned)
- Suitable for configuration files and data exchange
- Works for traditional APIs and AI interaction scenarios
- Supports conversion between language structs/classes and MetaMessage
- Supports code generation for multiple languages
- Data carries type, constraint, description, and example without separate documentation
- All metadata can be updated with the data itself, without extra coordination
- Structures and values stay consistent across languages
- No structural loss; parsers adapt automatically and do not crash
- Can serialize to compact binary for faster decoding and smaller size

**Problems solved**

- Unknown types, such as not knowing whether a field is uint8
- Incomplete structure, such as null without inner type information
- No validation rules, so data legality cannot be checked
- No examples or descriptions, forcing reliance on separate docs
- Format changes require protocol adjustment and documentation resync

MetaMessage is naturally suited for AI understanding and interaction, solving ambiguity and imprecision in data. It replaces traditional API docs, verbal format agreements, and manual version sync by making data self-explanatory and independently evolvable.

[github.com](https://github.com/metamessage/metamessage)

[NuGet](https://www.nuget.org/packages/MetaMessage)

## 1. 安装

### 1.1 项目引用

```bash
dotnet add package MetaMessage
```

### 1.2 版本要求

- .NET 8.0 或更高版本

## 2. 核心 API

### 2.1 二进制序列化 (Encode / Decode)

```csharp
using MetaMessage.Core;

public class User
{
    public string name { get; set; } = "";
    public int age { get; set; }
    public bool active { get; set; }
    public List<int> scores { get; set; } = new();
}

var user = new User
{
    name = "Alice",
    age = 25,
    active = true,
    scores = new List<int> { 95, 87, 92 }
};

// 编码为二进制
byte[] data = MetaMessage.Encode(user);

// 解码到新对象
User decoded = MetaMessage.Decode<User>(data);

// 解码到已有对象
var target = new User();
MetaMessage.Decode(data, target);

// 解码为树结构
INode tree = MetaMessage.Decode(data);
```

### 2.2 使用 Tag 标注自定义序列化

```csharp
public class Product
{
    [MM("type=str; desc=产品名称; min=1; max=100")]
    public string Title { get; set; } = "";

    [MM("type=f64; desc=价格; min=0; max=999999.99")]
    public double Price { get; set; }

    [MM("type=i; desc=库存; default=0")]
    public int Stock { get; set; }

    [MM("type=bool; desc=是否上架")]
    public bool IsActive { get; set; }

    [MM("-")] // 排除此属性
    public string InternalNotes { get; set; } = "";
}

var product = new Product
{
    Title = "无线鼠标",
    Price = 99.99,
    Stock = 1000,
    IsActive = true
};

byte[] data = MetaMessage.Encode(product);
var restored = MetaMessage.Decode<Product>(data);
```

### 2.3 从值编码 (FromValue)

```csharp
// 直接编码匿名对象或值
byte[] data = MetaMessage.FromValue(new { name = "test", count = 42 }, "type=obj; desc=示例");

// 使用 Tag 字符串
byte[] tagged = MetaMessage.FromValue(42, "type=i; desc=数量; min=0; max=100");
```

### 2.4 解码为 JSONC 字符串

```csharp
byte[] data = MetaMessage.Encode(user);

// 将二进制数据解码为 JSONC 格式字符串
string jsonc = MetaMessage.DecodeToJsonc(data);
Console.WriteLine(jsonc);
```

### 2.5 对象转 JSONC

```csharp
// 将对象转换为 JSONC 字符串
string jsonc = MetaMessage.ValueToJsonc(user);
Console.WriteLine(jsonc);
```

### 2.6 校验 (Validate)

```csharp
using MetaMessage.Ir;

var tag = Tag.Parse("type=i; min=0; max=100");

// 校验值是否符合 Tag 约束
ValidationResult result = MetaMessage.Validate(42, tag);
Console.WriteLine(result.IsValid); // True

result = MetaMessage.Validate(-1, tag);
Console.WriteLine(result.IsValid); // False（超出最小值）
```

### 2.7 快捷 API（推荐）

`Mm` 类提供简洁的静态方法，适合快速编码/解码值，无需导入 `Core` 命名空间：

```csharp
using MetaMessage;

// 从值编码为二进制（可选 Tag 字符串）
byte[] wire = Mm.EncodeFromValue(new { name = "hello", count = 42 });
byte[] tagged = Mm.EncodeFromValue(42, "type=i; desc=数量");

// 从 JSONC 编码为二进制
byte[] data = Mm.EncodeFromJsonc(@"{""a"": 1, ""b"": 2}");

// 从二进制解码为纯值（Dictionary / List / 标量）
object? result = Mm.DecodeToValue(wire);
// result 为 Dictionary<string, object?> 或 List<object?> 或原始值

// 二进制解码为 JSONC 字符串
string jsonc = Mm.DecodeToJsonc(wire);

// 值转 JSONC（可选 Tag 字符串）
string j = Mm.ValueToJsonc(new { foo = "bar", num = 123 });
string j2 = Mm.ValueToJsonc("hello", "type=str; desc=名称");

// JSONC 字符串转纯值
object? val = Mm.JsoncToValue(j);
```

## 3. JSONC API

JSONC 是一种支持注释的 JSON 格式，使用 `// mm:` 语法附加 Tag 元数据。

### 3.1 解析与输出

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

// 格式化输出（带缩进）
Console.WriteLine(JsoncParser.ToString(node));

// 紧凑输出（无缩进）
Console.WriteLine(JsoncParser.ToMinString(node));

// 从字节数组解析
byte[] bytes = System.Text.Encoding.UTF8.GetBytes(jsonc);
var nodeFromBytes = JsoncParser.ParseFromBytes(bytes);

// 解析 JSONC 并转换为 INode
INode tree = JsoncParser.ParseFromJSONC(jsonc);

// 将 INode 转换为 JSONC 字符串
string output = JsoncParser.ToJSONC(tree);
```

### 3.2 JSONC 绑定到对象

```csharp
// 绑定 JSONC 到强类型对象（自动创建实例）
var user = JsoncParser.BindFromString<User>(jsonc);
Console.WriteLine($"name: {user.name}, age: {user.age}");

// 从字节数组绑定
var userFromBytes = JsoncParser.BindFromBytes<User>(bytes);

// 绑定到已有对象
var target = new User();
JsoncParser.BindFromString(jsonc, target);
JsoncParser.BindFromBytes(bytes, target);
```

### 3.3 对象转 JSONC

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

### 3.4 JSONC 值提取

```csharp
// 提取为纯 JSON 对象（字典、列表、原始值）
object? extracted = JsoncParser.ExtractValue(node);

if (extracted is Dictionary<string, object?> dict)
{
    Console.WriteLine(dict["name"]);
}
```

### 3.5 通过 MetaMessage 类快捷调用

```csharp
// 从 JSONC 字符串编码为二进制
byte[] data = MetaMessage.FromJSONC(jsonc);

// 从 JSONC 绑定到对象
var user = new User();
MetaMessage.BindFromJSONC(jsonc, user);

// 从 JSONC 解析为树
INode tree = MetaMessage.ParseFromJSONC(jsonc);
```

## 4. Tag 系统

Tag 是 MetaMessage 的核心元数据系统，用于描述字段的类型、约束和行为。

### 4.1 Tag 属性

```csharp
using MetaMessage.Ir;

var tag = Tag.Parse("type=str; desc=用户名; default=guest; min=1; max=50; " +
    "pattern=^[a-zA-Z]+$; nullable; enums=admin|user|guest");

// 基本属性
tag.Type;      // ValueType.Str
tag.Name;      // 字段名
tag.Desc;      // "用户名"

// 约束
tag.DefaultVal; // "guest"
tag.Min;        // "1"
tag.Max;        // "50"
tag.Pattern;    // "^[a-zA-Z]+$"

// 行为
tag.Nullable;   // true
tag.IsNull;     // false
tag.Example;    // false

// 枚举
tag.Enums;      // "admin|user|guest"

// 已废弃
tag.Deprecated; // false

// 子类型（用于数组/Map）
tag.ChildType;  // ValueType.Unknown
tag.ChildDesc;
tag.ChildMin;
tag.ChildMax;
```

### 4.2 Tag 字符串语法

```
type=<type>; name=<name>; desc=<description>; default=<value>
min=<min>; max=<max>; pattern=<regex>
nullable; example; deprecated; allow_empty; unique
enums=<val1>|<val2>|<val3>
is_null
child_type=<type>; child_desc=<desc>
child_min=<min>; child_max=<max>
child_default=<value>; child_pattern=<regex>
child_enums=<val1>|<val2>|<val3>
```

### 4.3 支持的 ValueType

| 值         | 描述        | C# 类型            |
| ---------- | ----------- | ------------------ |
| `unknown`  | 未知        | -                  |
| `str`      | 字符串      | `string`           |
| `bool`     | 布尔        | `bool`             |
| `i`        | 整数 (32位) | `int`              |
| `i8`       | 有符号 8位  | `sbyte`            |
| `i16`      | 有符号 16位 | `short`            |
| `i32`      | 有符号 32位 | `int`              |
| `i64`      | 有符号 64位 | `long`             |
| `u`        | 无符号整数  | `uint`             |
| `u8`       | 无符号 8位  | `byte`             |
| `u16`      | 无符号 16位 | `ushort`           |
| `u32`      | 无符号 32位 | `uint`             |
| `u64`      | 无符号 64位 | `ulong`            |
| `f32`      | 浮点 (32位) | `float`            |
| `f64`      | 浮点 (64位) | `double`           |
| `bytes`    | 字节数组    | `byte[]`           |
| `decimal`  | 高精度小数  | `decimal`/`string` |
| `bigint`   | 大整数      | `string`           |
| `datetime` | 日期时间    | `DateTime`         |
| `date`     | 日期        | `DateTime`         |
| `time`     | 时间        | `DateTime`         |
| `uuid`     | UUID        | `string`/`Guid`    |
| `url`      | URL         | `string`           |
| `email`    | 电子邮件    | `string`           |
| `ip`       | IP 地址     | `string`           |
| `enums`    | 枚举        | `enum`/`string`    |
| `obj`      | 对象        | `class`/`struct`   |
| `map`      | 字典        | `Dictionary<K,V>`  |
| `vec`      | 数组/列表   | `List<T>`/`T[]`    |
| `arr`      | 定长数组    | `T[]`              |
| `doc`      | 文档根      | -                  |
| `media`    | 媒体        | `byte[]`           |

## 5. 树节点类型

MetaMessage 使用 `INode` 接口表示结构化的数据树：

| 类型         | 描述      | 特性                                      |
| ------------ | --------- | ----------------------------------------- |
| `NodeScalar` | 标量值    | `Data`(原始值), `Text`(字符串表示), `Tag` |
| `MmArray`    | 数组/列表 | `Children`(子节点列表), `Tag`             |
| `MmMap`      | 对象/字典 | `Entries`(键值对列表), `Tag`              |
| `MmDoc`      | 文档根    | `Fields`(顶级字段), `Tag`                 |

```csharp
INode tree = MetaMessage.Decode(data);

switch (tree)
{
    case NodeScalar scalar:
        Console.WriteLine($"{scalar.Data} ({scalar.Text})");
        break;
    case MmArray array:
        foreach (var child in array.Children) { }
        break;
    case MmMap map:
        foreach (var entry in map.Entries)
        {
            Console.WriteLine($"{entry.Key.Text}: {entry.Value}");
        }
        break;
}
```

## 6. 低层 API

### 6.1 WireEncoder / WireDecoder

```csharp
// 手动编码
var encoder = new WireEncoder();
encoder.EncodeInt64(42);
encoder.EncodeString("hello");
encoder.EncodeBool(true);
encoder.EncodeFloatString("3.14");
byte[] data = encoder.ToByteArray();

// 带标签编码
var tag = Tag.Parse("type=str; desc=示例");
var payload = new WireEncoder();
payload.EncodeString("test");
encoder.EncodeTaggedPayload(payload.ToByteArray(), tag.ToBytes());

// 手动解码
var decoder = new WireDecoder(data);
INode tree = decoder.Decode();
```

### 6.2 从反射编码 (ReflectMmEncoder)

```csharp
// 将任意对象转换为 INode
INode tree = ReflectMmEncoder.ValueToNode(obj, "");
```

### 6.3 绑定 (ReflectMmBinder)

```csharp
// 将二进制数据绑定到已有对象
ReflectMmBinder.Bind(data, target);
```

## 7. 验证器

内置验证器支持对值进行 Tag 约束校验：

```csharp
using MetaMessage.Core;
using MetaMessage.Ir;

var tag = Tag.Parse("type=i; min=0; max=100; pattern=^\\d+$");

ValidationResult result = Validator.Validate(50, tag);
Console.WriteLine(result.IsValid); // True
Console.WriteLine(result.Error);   // null

result = Validator.Validate(200, tag);
Console.WriteLine(result.IsValid); // False
Console.WriteLine(result.Error);   // "超出最大值: 200"

// 支持所有 ValueType 的验证
Validator.Validate("hello", Tag.Parse("type=str; min=1; max=10"));
Validator.Validate(3.14, Tag.Parse("type=f64; min=0; max=10"));
```

## 8. 类型推断

使用 `[MM]` 属性标记类成员时，可以省略 `type=`；省略时将从属性类型自动推断：

```csharp
public class Config
{
    [MM("desc=服务器地址")]
    public string Host { get; set; } = "localhost";

    [MM("desc=端口号; min=1; max=65535")]
    public int Port { get; set; } = 8080;

    [MM("desc=启用SSL")]
    public bool UseSsl { get; set; } = true;
}

// 等价于手动指定 type=str, type=i, type=bool
```

## 9. 类型映射参考

| C# 类型           | 推断的 ValueType | 编码方式            |
| ----------------- | ---------------- | ------------------- |
| `string`          | `Str`            | UTF-8 字符串        |
| `bool`            | `Bool`           | Simple (TRUE/FALSE) |
| `int`             | `I`              | 正整数/负整数       |
| `long`            | `I64`            | 正整数/负整数       |
| `short`           | `I16`            | 正整数/负整数       |
| `sbyte`           | `I8`             | 正整数/负整数       |
| `byte`            | `U8`             | 正整数              |
| `uint`            | `U`              | 正整数              |
| `ushort`          | `U16`            | 正整数              |
| `ulong`           | `U64`            | 正整数 / BigInt     |
| `float`           | `F32`            | 浮点数字符串        |
| `double`          | `F64`            | 浮点数字符串        |
| `decimal`         | `Decimal`        | 浮点数字符串        |
| `byte[]`          | `Bytes`          | 字节数组            |
| `DateTime`        | `Datetime`       | Unix 时间戳         |
| `List<T>`         | `Vec`            | 容器                |
| `T[]`             | `Arr`            | 容器                |
| `Dictionary<K,V>` | `Map`            | 容器                |
| `class`/`struct`  | `Obj`            | Map 容器            |

## 10. 测试

```bash
cd mm-cs
dotnet test meta-message.sln
```

测试框架：xUnit + .NET Test SDK

## 11. 相关资源

- [Go 实现](https://github.com/example/meta-message)
