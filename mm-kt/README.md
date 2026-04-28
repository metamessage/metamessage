# Kotlin MetaMessage 库使用说明

## 1. 安装

### Maven 依赖
将以下依赖添加到你的 `pom.xml` 文件中：

```xml
<dependency>
    <groupId>io.metamessage</groupId>
    <artifactId>mm-kt</artifactId>
    <version>1.0.0</version>
</dependency>
```

### 版本要求
- JDK 11 或更高版本
- Kotlin 1.5 或更高版本
- Maven 3.6 或更高版本

## 2. 基本使用

### 2.1 类定义
使用 `@MM` 注解标记需要编解码的类：

```kotlin
import io.metamessage.mm.MM

@MM
class Person(var name: String = "Ed", var age: Int = 30)
```

### 2.2 编码示例

```kotlin
import io.metamessage.mm.MetaMessage

val person = Person()
val wire = MetaMessage.encode(person)
println("Encoded: ${bytesToHex(wire)}")
```

### 2.3 解码示例

```kotlin
val decoded = MetaMessage.decode(wire, Person::class.java)
println("Decoded: Name=${decoded.name}, Age=${decoded.age}")
```

### 2.4 JSONC 解析示例

```kotlin
import io.metamessage.jsonc.Jsonc

val jsonc = """
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}
"""

// 解析 JSONC
val node = Jsonc.parseFromString(jsonc)

// 绑定到对象
val person = Jsonc.bindFromString(jsonc, Person::class.java)
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-kt 目录下运行
mvn test
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
