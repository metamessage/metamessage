# MetaMessage 多语言示例

本目录包含 MetaMessage 各语言库的使用示例。

## 目录结构

```
examples/
├── go/                    # Go 示例 (已存在)
├── java/                  # Java 示例
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── kotlin/                # Kotlin 示例
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── typescript/            # TypeScript 示例
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── python/                # Python 示例
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── javascript/            # JavaScript 示例
│   ├── basic/
│   ├── jsonc-to-wire/
│   ├── wire-to-jsonc/
│   └── bind-object/
├── csharp/                # C# 示例
│   └── ...
├── rust/                  # Rust 示例
│   └── ...
├── swift/                 # Swift 示例
│   └── ...
└── php/                   # PHP 示例
    └── ...
```

## 示例类型

### 1. Basic - 基础使用
演示如何使用结构体/类的基础编解码。

### 2. JSONC to Wire - JSONC 到 Wire 格式
演示如何将 JSONC 格式的数据编码为 Wire 二进制格式。

### 3. Wire to JSONC - Wire 格式到 JSONC
演示如何从 Wire 二进制格式解码回 JSONC 格式。

### 4. Bind Object - 对象绑定
演示如何将 JSONC 数据直接绑定到对象/结构体。

## 基础示例数据

### Person 示例
```jsonc
{
    // mm: type=str; desc=姓名
    "name": "Ed",
    // mm: type=i; desc=年龄
    "age": 30
}
```

### 完整示例
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

## 参考

- Go 文档: `/go/...`
- Java 文档: `/java/...`
- Kotlin 文档: `/kotlin/...`
- TypeScript 文档: `/typescript/...`
- Python 文档: `/python/...`
- JavaScript 文档: `/javascript/...`
- C# 文档: `/csharp/...`
- Rust 文档: `/rust/...`
- Swift 文档: `/swift/...`
- PHP 文档: `/php/...`
