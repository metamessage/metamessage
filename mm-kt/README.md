# MetaMessage

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
