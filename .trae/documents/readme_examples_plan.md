# 为各语言 README 添加编程语言实现示例计划

## 概述
为项目的多语言 README 文件添加各编程语言的示例代码。目前只有 Go 示例，需要为其他语言（Java, Kotlin, TypeScript, Python, JavaScript, C#, Rust, Swift, PHP）添加对应的示例代码。

## 目标 README 文件
- README.md (中文)
- README.en.md (英文)
- README.ja.md (日文)
- README.ko.md (韩文)
- README.es.md (西班牙文)
- README.fr.md (法文)
- README.de.md (德文)
- README.ru.md (俄文)
- README.vi.md (越南文)
- README.id.md (印尼文)
- README.th.md (泰文)

## 示例内容

### Go (已有)
```go
package main

import (
    "fmt"
    mm "github.com/metamessage/metamessage"
)

func main() {
    type Person struct {
        Name string
        Age  int
    }

    p := Person{Name: "Alice", Age: 30}
    data, err := mm.EncodeFromValue(p)
    // ...
}
```

### Java
```java
import io.github.metamessage.mm.MetaMessage;
import io.github.metamessage.mm.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}

public class Example {
    public static void main(String[] args) throws Exception {
        Person person = new Person();
        byte[] wire = MetaMessage.encode(person);
        Person decoded = MetaMessage.decode(wire, Person.class);
    }
}
```

### Kotlin
```kotlin
import io.github.metamessage.mm.MetaMessage
import io.github.metamessage.mm.MM

@MM
class Person(var name: String = "Ed", var age: Int = 30)

fun main() {
    val person = Person()
    val wire = MetaMessage.encode(person)
    val decoded = MetaMessage.decode(wire, Person::class.java)
}
```

### TypeScript
```typescript
import { encode, decode } from 'metamessage';

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

### Python
```python
from metamessage import encode, decode

person = {"name": "Ed", "age": 30}
wire = encode(person)
decoded = decode(wire)
```

### JavaScript
```javascript
const { encode, decode } = require('metamessage');

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

### C#
```csharp
using MetaMessage;

var person = new Person { Name = "Ed", Age = 30 };
byte[] wire = MetaMessage.Encode(person);
var decoded = MetaMessage.Decode<Person>(wire);
```

### Rust
```rust
use metamessage::{encode, decode, Node};

let person = Node::Object(/* ... */);
let wire = encode(&person);
let decoded = decode(&wire).unwrap();
```

### Swift
```swift
import MetaMessage

let person = Person(name: "Ed", age: 30)
let wire = MetaMessage.encode(person)
let decoded = try MetaMessage.decode(wire)
```

### PHP
```php
<?php
use io\metamessage\mm\MetaMessage;

$person = new Person();
$wire = MetaMessage::encode($person);
$decoded = MetaMessage::decode($wire, Person::class);
```

## 实施步骤

### 1. 分析现有 README 结构
- 查看 README.md 的 "Library Usage" 部分
- 确定示例代码应插入的位置

### 2. 为每种语言创建示例代码片段
- 为 10 种语言（Go, Java, Kotlin, TypeScript, Python, JavaScript, C#, Rust, Swift, PHP）创建示例
- 确保示例代码简洁明了
- 使用统一的 Person 示例数据结构

### 3. 更新各语言 README 文件
按顺序更新：
1. README.md (中文)
2. README.en.md (英文)
3. 其他语言 README 文件

### 4. 验证
- 确保代码格式正确
- 确保语言标识清晰

## 文件位置
- 计划文件: `.trae/documents/readme_examples_plan.md`
- 示例代码目录: `examples/`

## 考虑事项
- 不同语言的代码风格差异
- 注释应使用对应语言或英文
- 确保示例代码可读性
