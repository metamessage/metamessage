# MetaMessage

**⚠️ 重要提示**

这是一个**只读镜像**，本仓库仅作为 Composer 包的发布源，如有問題，请在主仓库的 Issues 中提交

**请勿直接向此仓库提交任何代码！**

## 🔄 同步规则

- 本仓库的所有代码，均通过自动化脚本从主仓库 `metamessage/metamessage` 的 `/mm-php` 目录**单向同步**而来。
- 任何直接提交到本仓库的 Pull Request 或推送，都会被**强制拒绝**。

## 🧑‍💻 正确的贡献方式

如果你需要修改 PHP SDK 的代码，请遵循以下流程：

1. 克隆主仓库：

   ```bash
   git clone https://github.com/metamessage/metamessage.git
   cd metamessage/mm-php
   ```

2. 在 `mm-php` 目录下完成开发和测试。

3. 将你的改动提交到主仓库

4. 本镜像仓库会通过 GitHub Actions 等方式自动同步主仓库的变更。

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
use io\metamessage\mm\MetaMessage;
```

### 2.2 类定义

```php
class Person {
    public $name = "Ed";
    public $age = 30;
}
```

### 2.3 编码示例

```php
$person = new Person();
$wire = MetaMessage::encode($person);
echo "Encoded: " . bin2hex(implode(array_map('chr', $wire))) . "\n";
```

### 2.4 解码示例

```php
$decoded = MetaMessage::decode($wire, Person::class);
echo "Decoded: Name={$decoded->name}, Age={$decoded->age}\n";
```

### 2.5 JSONC 解析示例

```php
use io\metamessage\jsonc\Jsonc;

$jsonc = '{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}';

// 解析 JSONC
$node = Jsonc::parseFromString($jsonc);

// 绑定到对象
$person = Jsonc::bindFromString($jsonc, Person::class);
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
