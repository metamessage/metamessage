# C# JSONC 实现计划

## 概述

JSONC 是带有注释支持的 JSON 扩展格式，支持在 JSON 中添加行注释(`//`)和块注释(`/**/`)，并通过注释添加元数据标签(`// mm:type=str;desc=...`)。本实现将基于 Go 版本的功能完整移植到 C#。

## 项目结构

```
mm-cs/src/MetaMessage/Jsonc/
├── Jsonc.cs                 # 主入口类
├── JsoncToken.cs            # Token 类型定义
├── JsoncScanner.cs          # 词法分析器
├── JsoncParser.cs            # 语法分析器
├── JsoncNode.cs              # AST 节点基类
├── JsoncObject.cs            # 对象节点
├── JsoncArray.cs            # 数组节点
├── JsoncValue.cs            # 值节点
├── JsoncTag.cs              # 元数据标签
├── JsoncBinder.cs           # 对象绑定器
├── JsoncPrinter.cs          # 输出格式化器
└── ValueType.cs             # 值类型枚举
```

## 实现步骤

### 1. 创建目录和基础文件

- 创建 `Jsonc/` 目录
- 实现 `ValueType.cs` - 值类型枚举
- 实现 `JsoncToken.cs` - Token 类型定义
- 实现 `JsoncTag.cs` - 元数据标签类

### 2. 实现词法分析器

- 实现 `JsoncScanner.cs`
  - 支持符号: `{`, `}`, `[`, `]`, `:`, `,`
  - 支持字面量: 字符串、数字、true、false
  - 支持行注释: `// ...`
  - 支持块注释: `/* ... */`
  - 跟踪行号和列号

### 3. 实现语法分析器

- 实现 `JsoncParser.cs`
  - 解析对象 `{}`
  - 解析数组 `[]`
  - 解析值（字符串、数字、布尔、null）
  - 处理注释中的标签元数据
  - 支持类型推断和验证

### 4. 实现 AST 节点

- 实现 `JsoncNode.cs` - 节点接口
- 实现 `JsoncObject.cs` - 对象节点
- 实现 `JsoncArray.cs` - 数组节点
- 实现 `JsoncValue.cs` - 值节点

### 5. 实现对象绑定器

- 实现 `JsoncBinder.cs`
  - 将 AST 绑定到 C# 对象
  - 支持结构体、映射、数组、切片
  - 支持各种值类型转换

### 6. 实现输出格式化器

- 实现 `JsoncPrinter.cs`
  - 将 AST 转换回 JSONC 字符串
  - 支持缩进格式化
  - 保留注释和标签

### 7. 实现主入口类

- 实现 `Jsonc.cs`
  - `ParseFromString(in string)` -> `IJcNode`
  - `ParseFromBytes(in byte[])` -> `IJcNode`
  - `BindFromString(in string, out T)` -> `T`
  - `BindFromBytes(in byte[], out T)` -> `T`
  - `ValueToJSONCString(value)` -> `string`
  - `ToString(node)` -> `string`

### 8. 创建测试用例

- 创建 `JsoncTest.cs`
- 测试基本类型解析
- 测试注释解析
- 测试标签解析
- 测试对象绑定
- 测试结构体转换

## 核心功能详解

### JsoncToken 类型

```csharp
public enum JsoncTokenType
{
    EOF,
    LBrace, RBrace,
    LBracket, RBracket,
    Colon, Comma,
    String, Number,
    True, False, Null,
    LeadingComment, TrailingComment
}
```

### JsoncTag 元数据标签

支持以下标签属性:
- `type` - 值类型 (str, i, i8, i16, i32, i64, u, u8, u16, u32, u64, f32, f64, bool, bytes, uuid, email, url, ip, datetime, date, time, decimal, enum, etc.)
- `desc` - 描述
- `nullable` - 可空
- `default` - 默认值
- `min/max` - 最小/最大值
- `size` - 大小
- `enum` - 枚举值
- `pattern` - 正则表达式
- `location` - 时区
- `version` - 版本
- `mime` - MIME类型
- 子元素标签 (child_type, child_desc 等)

### 支持的数据类型

| JSONC Type | C# Type |
|------------|---------|
| str | string |
| i/i8/i16/i32/i64 | int, int8, int16, int32, int64 |
| u/u8/u16/u32/u64 | uint, uint8, uint16, uint32, uint64 |
| f32/f64 | float, double |
| bool | bool |
| bytes | byte[] |
| uuid | string |
| email | string |
| url | string (System.Uri) |
| ip | string (System.Net.IPAddress) |
| datetime/date/time | DateTime |
| decimal | string |
| enum | string |
| struct | class/struct |
| slice | List<T> |
| arr | T[] |
| map | Dictionary<string, T> |

### 注释格式

```jsonc
{
    // mm:type=str;desc=用户名
    "user_name": "example",
    // mm:type=i;desc=年龄;min=0;max=150
    "age": 25
}
```

## 依赖关系

- `MetaMessage.Mm` - 使用已有的 Mm 模块
- `System.Text.RegularExpressions` - 正则表达式
- `System.Net` - IP 地址
- `System.Net.IPAddress` - IP 地址解析

## 实现注意事项

1. **性能优化**: 使用字符串 Builder 进行字符串拼接
2. **错误处理**: 详细的错误信息和位置
3. **类型安全**: 充分利用 C# 的类型系统
4. **可扩展性**: 设计支持未来扩展
5. **一致性**: 保持与 Go 版本的功能一致