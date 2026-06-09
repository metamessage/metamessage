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

### Composer 依赖

将以下依赖添加到你的 `composer.json` 文件中：

```bash
composer require metamessage/metamessage
```

````

```json
{
  "require": {
    "metamessage/metamessage": "1.0.0"
  }
}
```

然后运行：

```bash
composer install
```

### 版本要求

- PHP 7.4 或更高版本

## 2. 基本使用

### 2.1 导入命名空间

```php
use io\metamessage\core\MetaMessage;
use function io\metamessage\encodeFromValue;
use function io\metamessage\decodeToValue;
use function io\metamessage\valueToJsonc;
use function io\metamessage\jsoncToValue;
```

### 2.2 类定义

```php
class Person {
    public $name = "Ed";
    public $age = 30;
}
```

### 2.3 快捷 API（推荐）

提供 namespace 级别的函数，类似 Python 版本的顶层 API，适合快速编码/解码：

```php
// PHP 值 → 二进制
$wire = encodeFromValue(["name" => "hello", "count" => 42]);

// 二进制 → PHP 值
$result = decodeToValue($wire);
// $result = ["name" => "hello", "count" => 42]

// PHP 值 → JSONC
$jsonc = valueToJsonc(["foo" => "bar", "num" => 123]);

// JSONC → PHP 值
$back = jsoncToValue($jsonc);

// JSONC → 二进制
$wire = encodeFromJsonc('{"a": 1, "b": 2}');

// 二进制 → JSONC
$jsonc = decodeToJsonc($wire);
```

所有函数均支持可选 `$tag` 参数（mm tag 字符串）：

```php
$wire = encodeFromValue(42, 'name=count; type=i; desc=数量');
$jsonc = valueToJsonc("hello", 'type=str; desc=名称');
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-php 目录下运行
composer test
```

### 3.2 测试框架

- PHPUnit

### 3.3 测试覆盖范围

- 编码测试
- 解码测试
- JSONC 解析测试
- 绑定测试

## 4. 常见问题

### 4.1 依赖问题

- **问题**: Composer 依赖安装失败
  **解决**: 检查网络连接，或使用 Composer 镜像

### 4.2 运行时问题

- **问题**: 编码/解码失败
  **解决**: 检查类定义是否正确，属性是否可访问

## 5. 示例代码

查看 `https://github.com/metamessage/metamessage/examples/php/` 目录下的示例代码：

- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式
- `wire-to-jsonc/` - Wire 格式转 JSONC
- `bind-object/` - 对象绑定示例

## 6. 相关资源

- [metamessage/mm-php](https://github.com/metamessage/mm-php)
- [metamessage/metamessage](https://github.com/metamessage/metamessage)
- [PHP 文档](https://www.php.net/docs.php)
- [Composer 文档](https://getcomposer.org/doc/)
- [PHPUnit 文档](https://phpunit.readthedocs.io/)
````
