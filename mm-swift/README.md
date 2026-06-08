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

## 1. 安装

### Swift Package Manager

在你的 `Package.swift` 文件中添加依赖：

```swift
// swift-tools-version:5.5
import PackageDescription

let package = Package(
    name: "YourProject",
    dependencies: [
        .package(url: "https://github.com/metamessage/metamessage.git", from: "1.0.0")
    ],
    targets: [
        .target(
            name: "YourProject",
            dependencies: ["MetaMessage"]
        ),
        .testTarget(
            name: "YourProjectTests",
            dependencies: ["YourProject"]
        )
    ]
)
```

### 版本要求

- Swift 5.5 或更高版本
- iOS 13.0+ / macOS 10.15+ / tvOS 13.0+ / watchOS 6.0+

## 2. 基本使用

### 2.1 导入模块

```swift
import MetaMessage
```

### 2.2 编码示例

```swift
struct Person {
    var name: String = "Ed"
    var age: Int = 30
}

let person = Person()
let wire = MetaMessage.encodeFromValue(person)
print("Encoded: \(wire)")
```

### 2.3 解码示例

```swift
if let decoded = try? MetaMessage.decodeToValue(wire) {
    print("Decoded: \(decoded)")
}
```

### 2.4 JSONC 解析示例

```swift
import MetaMessage.JSONC

let jsonc = """
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}
"""

// 解析 JSONC
if let node = try? JSONC.parse(jsonc) {
    print("Parsed: \(node)")
}
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-swift 目录下运行
swift test
```

### 3.2 测试框架

- XCTest

### 3.3 测试覆盖范围

- 编码测试
- 解码测试
- JSONC 解析测试

## 4. 常见问题

### 4.1 依赖问题

- **问题**: Swift Package Manager 依赖下载失败
  **解决**: 检查网络连接，或使用镜像源

### 4.2 编译问题

- **问题**: 编译错误
  **解决**: 确保 Swift 版本符合要求，并且依赖配置正确

### 4.3 运行时问题

- **问题**: 编码/解码失败
  **解决**: 检查数据结构是否正确

## 5. 示例代码

查看 `examples/swift/` 目录下的示例代码：

- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式
- `wire-to-jsonc/` - Wire 格式转 JSONC
- `bind-object/` - 对象绑定示例

## 6. 相关资源

- [Swift 文档](https://docs.swift.org/swift-book/)
- [Swift Package Manager 文档](https://swift.org/package-manager/)
- [XCTest 文档](https://developer.apple.com/documentation/xctest)
