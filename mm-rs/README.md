# Rust MetaMessage 库使用说明

## 1. 安装

### Cargo 依赖
将以下依赖添加到你的 `Cargo.toml` 文件中：

```toml
[dependencies]
metamessage = "1.0.0"
```

### 版本要求
- Rust 1.56 或更高版本

## 2. 基本使用

### 2.1 导入模块

```rust
use metamessage::{encode, decode, Node};
```

### 2.2 编码示例

```rust
let person = Node::Object(/* 构造对象 */);
let wire = encode(&person);
println!("Encoded: {:?}", wire);
```

### 2.3 解码示例

```rust
let decoded = decode(&wire).unwrap();
println!("Decoded: {:?}", decoded);
```

### 2.4 JSONC 解析示例

```rust
use metamessage::{parse_jsonc, to_jsonc_string};

let jsonc = r#"
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}
"#;

// 解析 JSONC
let node = parse_jsonc(jsonc).unwrap();

// 转换为字符串
let jsonc_string = to_jsonc_string(&node);
println!("JSONC: {}", jsonc_string);
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-rs 目录下运行
cargo test
```

### 3.2 测试框架
- Rust 标准测试框架

### 3.3 测试覆盖范围
- 编码测试
- 解码测试
- JSONC 解析测试

## 4. 常见问题

### 4.1 依赖问题
- **问题**: Cargo 依赖下载失败
  **解决**: 检查网络连接，或使用 Cargo 镜像

### 4.2 编译问题
- **问题**: 编译错误
  **解决**: 确保 Rust 版本符合要求，并且依赖配置正确

### 4.3 运行时问题
- **问题**: 编码/解码失败
  **解决**: 检查数据结构是否正确

## 5. 示例代码

查看 `examples/rust/` 目录下的示例代码：
- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式
- `wire-to-jsonc/` - Wire 格式转 JSONC
- `bind-object/` - 对象绑定示例

## 6. 相关资源

- [Rust 文档](https://doc.rust-lang.org/)
- [Cargo 文档](https://doc.rust-lang.org/cargo/)
