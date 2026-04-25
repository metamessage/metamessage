# C# 版本 MetaMessage 实现计划

## 项目结构

创建 `mm-cs` 目录，包含以下结构：

```
mm-cs/
├── src/
│   └── MetaMessage/
│       ├── Jsonc/
│       │   ├── JcField.cs
│       │   ├── JcNode.cs
│       │   ├── Jsonc.cs
│       │   ├── JsoncBinder.cs
│       │   ├── JsoncParser.cs
│       │   ├── JsoncPrinter.cs
│       │   ├── JsoncScanner.cs
│       │   └── JsoncToken.cs
│       ├── Mm/
│           ├── BigIntWireCodec.cs
│           ├── CamelToSnake.cs
│           ├── FloatCodec.cs
│           ├── GrowableByteBuf.cs
│           ├── MM.cs (属性类)
│           ├── MetaMessage.cs (公共 API)
│           ├── MimeWire.cs
│           ├── MmDecodeException.cs
│           ├── MmTag.cs
│           ├── MmTree.cs
│           ├── Prefix.cs
│           ├── ReflectMmBinder.cs
│           ├── ReflectMmEncoder.cs
│           ├── SimpleValue.cs
│           ├── TagFieldParser.cs
│           ├── TimeUtil.cs
│           ├── TypeInference.cs
│           ├── ValueType.cs
│           ├── WireConstants.cs
│           ├── WireDecoder.cs
│           └── WireEncoder.cs
├── tests/
│   └── MetaMessageTests/
│       └── MetaMessageTest.cs
├── MetaMessage.csproj
└── README.md
```

## 实现步骤

### 1. 创建项目结构和基础文件
- 创建 `mm-cs` 目录
- 创建 `MetaMessage.csproj` 文件，配置项目依赖
- 创建 `src` 和 `tests` 目录结构

### 2. 实现核心常量和枚举
- `WireConstants.cs` - 定义编解码相关的常量
- `Prefix.cs` - 定义二进制格式的前缀
- `SimpleValue.cs` - 定义简单值类型
- `ValueType.cs` - 定义数据类型枚举

### 3. 实现工具类
- `GrowableByteBuf.cs` - 可增长的字节缓冲区
- `CamelToSnake.cs` - 驼峰命名转蛇形命名
- `TimeUtil.cs` - 时间处理工具
- `FloatCodec.cs` - 浮点数编解码
- `BigIntWireCodec.cs` - 大整数编解码
- `MimeWire.cs` - MIME 类型编解码

### 4. 实现核心数据结构
- `MmTag.cs` - 标签字段
- `MmTree.cs` - 数据树结构
- `TypeInference.cs` - 类型推断

### 5. 实现编解码逻辑
- `WireEncoder.cs` - 二进制编码
- `WireDecoder.cs` - 二进制解码
- `TagFieldParser.cs` - 标签字段解析

### 6. 实现反射和绑定逻辑
- `ReflectMmEncoder.cs` - 基于反射的编码
- `ReflectMmBinder.cs` - 基于反射的解码和绑定
- `MM.cs` - 属性类，用于字段元数据

### 7. 实现公共 API
- `MetaMessage.cs` - 公共接口，提供编码和解码方法

### 8. 实现 JSONC 相关功能
- `Jsonc.cs` - JSONC 核心功能
- `JsoncParser.cs` - JSONC 解析器
- `JsoncScanner.cs` - JSONC 扫描器
- `JsoncToken.cs` - JSONC 令牌
- `JsoncPrinter.cs` - JSONC 打印器
- `JsoncBinder.cs` - JSONC 绑定
- `JcNode.cs` - JSONC 节点
- `JcField.cs` - JSONC 字段

### 9. 实现测试用例
- `MetaMessageTest.cs` - 完整的测试用例，覆盖所有数据类型和标签场景

## 技术要点

1. **C# 特性**：使用 C# 特性（Attributes）代替 Java 的注解，用于字段元数据定义
2. **反射**：使用 C# 的反射 API 实现对象的自动编码和解码
3. **类型系统**：处理 C# 的类型系统，包括可空类型、泛型等
4. **二进制格式**：实现与 Go 版本兼容的二进制 wire 格式
5. **性能优化**：考虑 C# 特有的性能优化方式，如 Span<T>、内存池等

## 依赖项

- .NET 6.0 或更高版本
- 测试框架：xUnit 或 MSTest

## 验证计划

1. 运行单元测试，确保所有测试通过
2. 与 Go 版本进行互操作性测试
3. 性能测试，确保编码和解码性能符合预期

## 风险评估

1. **类型系统差异**：C# 和 Go 的类型系统存在差异，需要特别处理
2. **反射性能**：C# 的反射性能可能与 Go 有所不同，需要优化
3. **二进制格式兼容性**：确保与 Go 版本的二进制格式完全兼容

## 预期成果

- 完整的 C# 版本 MetaMessage 实现
- 与 Go 版本兼容的二进制编解码
- 完整的测试覆盖
- 良好的文档和示例