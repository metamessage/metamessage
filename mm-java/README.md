# Java MetaMessage 库使用说明

## 1. 安装

### Maven 依赖
将以下依赖添加到你的 `pom.xml` 文件中：

```xml
<dependency>
    <groupId>io.metamessage</groupId>
    <artifactId>mm-java</artifactId>
    <version>1.0.0</version>
</dependency>
```

### 版本要求
- JDK 11 或更高版本
- Maven 3.6 或更高版本

## 2. 基本使用

### 2.1 结构体定义
使用 `@MM` 注解标记需要编解码的类：

```java
import io.metamessage.mm.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}
```

### 2.2 编码示例

```java
import io.metamessage.mm.MetaMessage;

Person person = new Person();
byte[] wire = MetaMessage.encode(person);
System.out.println("Encoded: " + bytesToHex(wire));
```

### 2.3 解码示例

```java
Person decoded = MetaMessage.decode(wire, Person.class);
System.out.println("Decoded: Name=" + decoded.name + ", Age=" + decoded.age);
```

### 2.4 JSONC 解析示例

```java
import io.metamessage.jsonc.Jsonc;

String jsonc = """
{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}
""";

// 解析 JSONC
var node = Jsonc.parseFromString(jsonc);

// 绑定到对象
Person person = Jsonc.bindFromString(jsonc, Person.class);
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-java 目录下运行
mvn test
```

### 3.2 测试框架
- JUnit 5
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
  **解决**: 检查类定义是否正确，字段是否可访问

## 5. 示例代码

查看 `examples/java/` 目录下的示例代码：
- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式
- `wire-to-jsonc/` - Wire 格式转 JSONC
- `bind-object/` - 对象绑定示例

## 6. 相关资源

- [Java 文档](https://docs.oracle.com/en/java/)
- [Maven 文档](https://maven.apache.org/guides/index.html)
- [JUnit 5 文档](https://junit.org/junit5/docs/current/user-guide/)
