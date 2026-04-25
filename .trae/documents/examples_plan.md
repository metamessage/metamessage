# 各语言使用示例添加计划

## 概述
为项目中各语言（Java, Kotlin, TypeScript, Python, JavaScript, C#, Rust, Swift, PHP）添加使用示例，统一放在 `examples/` 目录下。

## 目标
基于 README 中的 Go 示例，为每个语言创建以下示例：

1. **Basic**: 基础使用（结构体编解码）
2. **JSONC to Wire**: JSONC 转 Wire 格式
3. **Wire to JSONC**: Wire 转 JSONC
4. **Bind Object**: 对象绑定

## 项目结构
```
examples/
├── go/ (现有)
│   ├── basic/
│   ├── int/
│   ├── one/
│   └── two/
├── java/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── kotlin/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── typescript/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── python/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── javascript/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── csharp/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── rust/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── swift/
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
└── php/
    ├── basic/
    ├── jsonc-to-wire/
    ├── wire-to-jsonc/
    └── bind-object/
```

## 基础示例内容（基于 README）

### Basic 示例
```jsonc
{
    // mm: type=str; desc=姓名
    "Name": "Ed",
    // mm: type=i; desc=年纪
    "Age": 30
}
```

### JSONC to Wire 示例
```jsonc
{
    // mm: type=datetime; desc=创建时间
    "create_time": "2026-01-01 00:00:00",
    // mm: type=str; desc=用户名称
    "user_name": "Alice",
    // mm: type=bool; desc=是否激活
    "is_active": true,
    // mm: type=array; child_type=i
    "scores": [95, 87, 92]
}
```

## 各语言示例内容

### Java
- 创建 Maven 项目结构
- 示例代码：
  - Basic.java: 结构体编解码
  - JsoncToWire.java: JSONC 转 Wire
  - WireToJsonc.java: Wire 转 JSONC
  - BindObject.java: 对象绑定

### Kotlin
- 创建 Maven 项目结构
- 示例代码：
  - Basic.kt: 结构体编解码
  - JsoncToWire.kt: JSONC 转 Wire
  - WireToJsonc.kt: Wire 转 JSONC
  - BindObject.kt: 对象绑定

### TypeScript
- 创建 Node.js 项目结构
- 示例代码：
  - basic.ts: 结构体编解码
  - jsonc-to-wire.ts: JSONC 转 Wire
  - wire-to-jsonc.ts: Wire 转 JSONC
  - bind-object.ts: 对象绑定

### Python
- 示例代码：
  - basic.py: 结构体编解码
  - jsonc_to_wire.py: JSONC 转 Wire
  - wire_to_jsonc.py: Wire 转 JSONC
  - bind_object.py: 对象绑定

### JavaScript
- 示例代码：
  - basic.js: 结构体编解码
  - jsonc-to-wire.js: JSONC 转 Wire
  - wire-to-jsonc.js: Wire 转 JSONC
  - bind-object.js: 对象绑定

### C#
- 创建 .NET 项目结构
- 示例代码：
  - Basic.cs: 结构体编解码
  - JsoncToWire.cs: JSONC 转 Wire
  - WireToJsonc.cs: Wire 转 JSONC
  - BindObject.cs: 对象绑定

### Rust
- 创建 Cargo 项目结构
- 示例代码：
  - basic.rs: 结构体编解码
  - jsonc_to_wire.rs: JSONC 转 Wire
  - wire_to_jsonc.rs: Wire 转 JSONC
  - bind_object.rs: 对象绑定

### Swift
- 创建 Swift 项目结构
- 示例代码：
  - Basic.swift: 结构体编解码
  - JsoncToWire.swift: JSONC 转 Wire
  - WireToJsonc.swift: Wire 转 JSONC
  - BindObject.swift: 对象绑定

### PHP
- 示例代码：
  - basic.php: 结构体编解码
  - jsonc-to-wire.php: JSONC 转 Wire
  - wire-to-jsonc.php: Wire 转 JSONC
  - bind-object.php: 对象绑定

## 实施步骤

1. **创建目录结构**
   - 在 `examples/` 下创建各语言的目录和子目录

2. **为每个语言添加示例**
   - 添加各语言的基础示例代码
   - 使用 README 中的示例 JSONC 内容
   - 确保各语言示例逻辑一致

3. **添加必要的项目文件**
   - Java: pom.xml
   - Kotlin: pom.xml
   - TypeScript: package.json
   - C#: .csproj
   - Rust: Cargo.toml
   - Swift: Package.swift

4. **更新 README**
   - 在根目录 README 中添加各语言示例说明

## 考虑事项
- 示例代码需保持简洁明了
- 确保各语言示例逻辑一致
- 添加注释说明关键步骤
- 示例应展示核心功能（JSONC 解析、编码、解码、绑定）
- 使用统一的 Person 结构体和 JSONC 数据
