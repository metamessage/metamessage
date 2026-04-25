# Rust 版本 MetaMessage (mm-rs) 实现计划

## 概述

将 Go 版本的 MetaMessage 功能完整移植到 Rust（不包括 CLI 功能），包括二进制 wire 格式编解码、JSONC 解析/生成、反射式对象绑定，并实现完整全面的测试。

## 项目结构

```
mm-rs/
├── Cargo.toml
├── src/
│   ├── lib.rs
│   ├── main.rs (可选)
│   ├── mm/
│   │   ├── mod.rs
│   │   ├── constants.rs      # Wire 格式常量
│   │   ├── prefix.rs        # 前缀枚举
│   │   ├── simple_value.rs  # 简单值枚举
│   │   ├── encoder.rs      # 编码器
│   │   ├── decoder.rs       # 解码器
│   │   ├── encode_int.rs    # 整数编码
│   │   ├── encode_float.rs   # 浮点数编码
│   │   ├── encode_str.rs     # 字符串编码
│   │   ├── encode_bytes.rs  # 字节编码
│   │   ├── encode_bool.rs   # 布尔编码
│   │   ├── encode_array.rs  # 数组编码
│   │   ├── encode_map.rs    # 映射编码
│   │   ├── encode_tag.rs    # 标签编码
│   │   ├── decode_tag.rs    # 标签解码
│   │   └── utils.rs        # 工具函数
│   ├── jsonc/
│   │   ├── mod.rs
│   │   ├── scanner.rs       # 词法分析器
│   │   ├── parser.rs        # 语法分析器
│   │   ├── ast.rs           # AST 节点定义
│   │   ├── tag.rs           # Tag 定义和解析
│   │   ├── binder.rs        # 对象绑定器
│   │   ├── to_string.rs     # JSONC 输出
│   │   └── value_type.rs   # 值类型枚举
│   └── tests/
│       ├── mod.rs
│       ├── mm_test.rs       # MM 编码解码测试
│       ├── jsonc_test.rs     # JSONC 解析测试
│       ├── bind_test.rs      # 绑定测试
│       └── full_test.rs     # 完整集成测试
```

## 支持的数据类型

### 基本类型

| Wire Type | Rust Type | 说明 |
|-----------|-----------|------|
| Simple | bool | 布尔值 (true/false) |
| Simple | i8, i16, i32, i64 | 有符号整数 |
| Simple | u8, u16, u32, u64 | 无符号整数 |
| Simple | f32, f64 | 浮点数 |
| Simple | String | 字符串 |
| Simple | Vec<u8> | 字节数组 |
| Simple | () | null 值 |

### 特殊类型

| Wire Type | Rust Type | 说明 |
|-----------|-----------|------|
| bytes | Vec<u8> | 二进制数据 |
| bytes | [u8; 16] | UUID |
| bytes | net::IpAddr | IP 地址 |
| string | String | Email |
| string | String | URL |
| string | String | UUID (字符串格式) |
| string | String | Decimal |
| int | i64 | DateTime (Unix epoch) |
| int | i64 | Date (天数) |
| int | i64 | Time (秒) |
| int | i64 | BigInt |
| int | String | Enum |

### 容器类型

| Wire Type | Rust Type | 说明 |
|-----------|-----------|------|
| Container.Array | [T; N] | 固定长度数组 |
| Container.Slice | Vec<T> | 可变长度切片 |
| Container.Map | HashMap<String, T> | 映射 |
| Container.Struct | struct | 结构体 |

## 实现步骤

### 阶段 1: 项目基础

1. 创建 `Cargo.toml` 和项目结构
2. 实现 `constants.rs` - Wire 格式常量
3. 实现 `prefix.rs` - 前缀枚举
4. 实现 `simple_value.rs` - 简单值枚举

### 阶段 2: MM 编码器

5. 实现 `encoder.rs` - 主编码器结构
6. 实现 `encode_int.rs` - 整数编码（正整数、负整数、VarInt）
7. 实现 `encode_float.rs` - 浮点数编码（IEEE 754）
8. 实现 `encode_str.rs` - 字符串编码（UTF-8）
9. 实现 `encode_bytes.rs` - 字节数组编码
10. 实现 `encode_bool.rs` - 布尔编码
11. 实现 `encode_array.rs` - 数组编码
12. 实现 `encode_map.rs` - 映射编码
13. 实现 `encode_tag.rs` - 标签编码

### 阶段 3: MM 解码器

14. 实现 `decoder.rs` - 主解码器结构
15. 实现 `decode_tag.rs` - 标签解码
16. 实现各类型的解码方法

### 阶段 4: JSONC 解析器

17. 实现 `value_type.rs` - 值类型枚举
18. 实现 `scanner.rs` - 词法分析器
    - 支持 `{`, `}`, `[`, `]`, `:`, `,`
    - 支持字符串、数字、true、false、null
    - 支持行注释 `//` 和块注释 `/**/`
    - 跟踪行列位置
19. 实现 `parser.rs` - 语法分析器
    - 解析对象 `{}`
    - 解析数组 `[]`
    - 解析值
    - 处理注释和 Tag

### 阶段 5: JSONC AST 和 Tag

20. 实现 `ast.rs` - AST 节点
    - `Value` - 值节点
    - `Object` - 对象节点
    - `Array` - 数组节点
    - `Field` - 字段
21. 实现 `tag.rs` - Tag 定义
    - `mm:` 标签解析
    - 类型、描述、可空性、最小值、最大值等
    - 验证方法

### 阶段 6: JSONC 绑定和输出

22. 实现 `binder.rs` - 对象绑定器
    - Struct 到 JSONC
    - JSONC 到 Struct
    - 类型转换
23. 实现 `to_string.rs` - JSONC 输出
    - 格式化输出
    - 保留注释

### 阶段 7: 集成和公共 API

24. 实现 `lib.rs` - 公共 API
    - `encode()` - 编码为二进制
    - `decode()` - 解码二进制
    - `parse_jsonc()` - 解析 JSONC
    - `to_jsonc_string()` - 转换为 JSONC 字符串
    - `bind()` - 绑定到对象

### 阶段 8: 测试

25. 实现 `mm_test.rs`
    - 所有基本类型编码解码测试
    - 整数范围测试
    - 浮点数精度测试
    - 字符串编码测试
    - 字节数组测试
    - 数组和切片测试
    - 映射测试
    - 结构体测试

26. 实现 `jsonc_test.rs`
    - 基本 JSON 解析测试
    - 行注释解析测试
    - 块注释解析测试
    - Tag 解析测试
    - 复杂嵌套结构测试

27. 实现 `bind_test.rs`
    - Struct 到 JSONC 转换测试
    - JSONC 到 Struct 绑定测试
    - 所有数据类型绑定测试
    - 可选字段测试
    - 默认值测试

28. 实现 `full_test.rs`
    - 完整流程测试
    - 往返编解码测试
    - 错误处理测试
    - 边界条件测试

## Tag 格式规范

```
// mm:type=str;desc=用户名;nullable=false;min=1;max=100
// mm:type=i;desc=年龄;min=0;max=150
// mm:type=f64;desc=价格;min=0.0
// mm:type=datetime;desc=创建时间;location=Asia/Shanghai
```

### 支持的 Tag 属性

| 属性 | 说明 | 示例 |
|------|------|------|
| type | 值类型 | str, i, i8, i16, i32, i64, u, u8, u16, u32, u64, f32, f64, bool, bytes, uuid, email, url, ip, datetime, date, time, decimal, enum |
| desc | 描述 | 用户名 |
| nullable | 可空 | true, false |
| default | 默认值 | 0 |
| min | 最小值 | 0 |
| max | 最大值 | 100 |
| size | 大小 | 255 |
| enum | 枚举值 | pending\|processing\|done |
| pattern | 正则表达式 | `^[a-z]+$` |
| location | 时区 | Asia/Shanghai |
| version | 版本 | 1.0 |
| mime | MIME类型 | image/png |
| child_type | 子元素类型 | str |
| child_desc | 子元素描述 | 列表项 |
| key_desc | 键描述 | 键名 |
| value_desc | 值描述 | 键值 |
| ele_desc | 元素描述 | 元素 |

## 测试覆盖要求

### 必测场景

1. **基本类型**
   - bool: true, false
   - i8, i16, i32, i64: 0, 正数, 负数, 边界值
   - u8, u16, u32, u64: 0, 正数, 边界值
   - f32, f64: 0, 正数, 负数, 小数, 特殊值 (Inf, NaN)
   - String: 空字符串, ASCII, Unicode, 特殊字符

2. **容器类型**
   - Vec<T>: 空, 单元素, 多元素
   - [T; N]: 固定长度数组
   - HashMap<String, T>: 空, 单键值, 多键值
   - Struct: 嵌套结构体

3. **特殊类型**
   - DateTime: Unix epoch, 不同时区
   - Date: 1970-01-01, 2038-01-19
   - Time: 00:00:00, 23:59:59
   - UUID: 有效 UUID 格式
   - Email: 有效格式
   - URL: 有效格式
   - IP: IPv4, IPv6
   - BigInt: 大整数

4. **Tag 处理**
   - 所有 Tag 属性解析
   - 类型验证
   - 约束验证 (min/max/enum/pattern)

5. **JSONC 特性**
   - 行注释解析
   - 块注释解析
   - 注释中的 Tag 解析
   - 混合空白字符

6. **错误处理**
   - 无效输入
   - 截断数据
   - 类型不匹配
   - 超出范围的值

## Rust 特性使用

- **Serde**: 可选集成 serde 进行结构体序列化
- **thiserror**: 错误处理
- **chrono**: 日期时间处理 (或使用 std::time)
- **uuid**: UUID 处理
- **ipnetwork**: IP 地址处理
- **regex**: 正则表达式验证
- **itoa**: 快速整数转字符串
- **ryu**: 快速浮点数转字符串

## 实施注意事项

1. **所有权**: Rust 的所有权模型需要在编码时仔细处理
2. **生命周期**: 确保引用生命周期正确
3. **错误处理**: 使用 Result 类型进行错误传播
4. **性能**: 使用 Buffer pool 避免频繁内存分配
5. **安全性**: 处理未信任输入时注意边界条件
6. **一致性**: 保持与 Go 版本的行为一致

## 验收标准

1. 所有 Go 版本的功能都已移植
2. 所有数据类型都有对应的编码解码测试
3. 所有 Tag 属性都有测试覆盖
4. 测试覆盖率 >= 90%
5. 与 Go 版本进行往返测试确保一致性
6. 代码通过 clippy 检查