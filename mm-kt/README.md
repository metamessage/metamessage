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

[![jitpack.io](https://jitpack.io/v/metamessage/metamessage.svg)](https://jitpack.io/#metamessage/metamessage)

## 1. 安装

### 依赖

将以下依赖添加到你的 `pom.xml` 文件中：

```xml
  <repositories>
		<repository>
		    <id>jitpack.io</id>
		    <url>https://jitpack.io</url>
		</repository>
	</repositories>

  <dependency>
	    <groupId>com.github.metamessage</groupId>
	    <artifactId>metamessage</artifactId>
	    <version>v0.1.10</version>
	</dependency>
```

### 版本要求

## 2. 基本使用

### 2.1 类定义

使用 `@MM` 注解标记需要编解码的类：

```kotlin
import io.github.metamessage.MM

@MM(desc = "Person")
class Person(var name: String = "Ed", var age: Int = 30)
```

### 2.2 示例

```kotlin
import io.github.metamessage.MetaMessage

val person = Person()
val wire = MetaMessage.encodeFromValue(person)
println("Encoded: ${bytesToHex(wire)}")

val decoded = MetaMessage.decodeToValue(wire, Person::class.java)
println("Decoded: Name=${decoded.name}, Age=${decoded.age}")

val jsoncInput = """
{
    // mm: desc=姓名
    "name": "Alice",
    // mm: type=u8; desc=年龄
    "age": 25
}
"""

// 编码 JSONC
val wire = MetaMessage.encodeFromJsonc(jsoncInput)

// 从二进制解码回 JSONC
val jsoncOutput = MetaMessage.decodeToJsonc(wire)

// 绑定到 Kotlin 对象
val person = MetaMessage.jsoncToValue(jsoncOutput, Person::class.java)

// 从对象重新生成 JSONC
val jsoncFromValue = MetaMessage.valueToJsonc(person)
```

```java
import io.github.metamessage.MetaMessage;
import io.github.metamessage.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}

public class Example {
    public static void main(String[] args) throws Exception {
        Person person = new Person();
        byte[] wire = MetaMessage.encodeFromValue(person);
        Person decoded = MetaMessage.decodeToValue(wire, Person.class);
    }
}
```

`nullable` 根據是否指針類型自動設置，只有指針類型的字段才可以設置`is_null`

| 类型 | 允许出现位置 | 语义 | 对应类型 | 允许标签 |
| ---- | ---- | ---- | ---- | ---- |
| Obj | class | 自定义结构体模型 | 自定义 MM 模型类 | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty` |
| Map | class | 字典</br>键值集合 | `dict[K, V]` | `example`</br>`deprecated`</br>`name`</br>`type`</br>`desc`</br>`allow_empty` |
| Vec | field | 一维数组 | `list[T]` </br> `[T]` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`type`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`unique`</br>**全部 `child_*` 子元素标签** |
| Arr | field | 固定长度数组 | `list[T]` </br> `[T]` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`unique`</br>**全部 `child_*` 子元素标签** |
| Str | field | 文本 | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern` |
| Bytes | field | 二进制数据、文件、媒体流 | `bytes` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size` |
| Bool | field | 真假状态 | `bool` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`nullable` |
| I、I8、I16、I32、I64、U、U8、U16、U32、U64、Bigint | field | 数字ID、数值、计数 | `int` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size` |
| F32、F64 | field | 普通小数（非金额） | `float` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size` |
| Datetime | field | 完整时间戳 | `datetime` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`location` |
| Date | field | 年月日 | `date` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`location` |
| Time | field | 时分秒 | `time` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`location` |
| Uuid | field | 唯一ID | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern`</br>`version` |
| Decimal | field | Decimal | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern` |
| Ip | field | Ip | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern`</br>`version` |
| Url | field | Url | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern` |
| Email | field | Email | `str` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`pattern` |
| Enums | field | 状态、选项、固定枚举值 | 自定义 `Enum` 子类 | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`enums` |
| Media | field | 图片</br>音视频</br>文件二进制 | `bytes` | `example`</br>`deprecated`</br>`name`</br>`desc`</br>`allow_empty`</br>`min`</br>`max`</br>`size`</br>`mime` |


## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-kt 目录下运行
mvn test
mvn -f mm-kt/pom.xml test -Dtest=MetaMessageTest
```

### 3.2 测试框架

- JUnit 5
- Kotest (可选)
- Maven Surefire Plugin

### 3.3 测试覆盖范围

- 编码测试
- 解码测试
- JSONC 解析测试
- 绑定测试

## 4. 常见问题

### 4.1 依赖问题

- **问题**: Maven 依赖下载失败
  **解决**: 检查网络连接，或使用 Maven 镜像

### 4.2 编译问题

- **问题**: 找不到 @MM 注解
  **解决**: 确保依赖配置正确，并且 IDE 已刷新依赖

### 4.3 运行时问题

- **问题**: 编码/解码失败
  **解决**: 检查类定义是否正确，属性是否可访问

## 5. 示例代码

查看 `examples/kotlin/` 目录下的示例代码：

- `basic/` - 基本使用示例

## 6. 相关资源

- [Kotlin 文档](https://kotlinlang.org/docs/home.html)
- [Maven 文档](https://maven.apache.org/guides/index.html)
- [JUnit 5 文档](https://junit.org/junit5/docs/current/user-guide/)
