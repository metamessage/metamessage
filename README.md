# MetaMessage

- [README 中文](README.md)
- [README English](README.en.md)
- [README 日本語](README.ja.md)
- [README 한국어](README.ko.md)
- [README Español](README.es.md)
- [README Français](README.fr.md)
- [README Deutsch](README.de.md)
- [README Русский](README.ru.md)
- [README Tiếng Việt](README.vi.md)
- [README Bahasa Indonesia](README.id.md)
- [README ไทย](README.th.md)

MetaMessage（mm）是一個結構化數據交換協議。自描述、自約束、自示例，實現了無損數據交換，是一種原生適配 AI + 人類 + 機器的下一代通用協議。

- 人類、AI閱讀友好
- 可以導出和導入為 JSONC、YAML、TOML 等文本格式，方便閱讀
- 同時適合配置文件、數據交互
- 適合傳統API、AI交互全場景
- 支持各語言結構體、類等和MetaMessage互轉
- 支持各語言結構體、類等以及數據代碼自動生成
- 數據自帶類型、約束、描述、示例，無需額外文檔或約定
- 類型、約束、描述、示例等所有信息可隨時更新，隨數據同步下發，無需主動通知
- 各語言之間數據結構及值完全一致
- 數據結構永不丟失，程序可自動適應，解析永遠不會崩潰
- 可以序列化為緊湊的二進制，解析更快，體積更小

**解決了以下問題**

- 類型不明確，如無法確定類型為uint8，無法進行數據結構重建
- 結構不完整，如null無法判斷內部結構，類型丟失
- 沒有約束規則，數據合法性無法驗證
- 沒有示例，沒有描述，必須依靠外部文檔，數據分離
- 格式變動敏感，編解碼需重新適配，文檔需重新同步

MetaMessage 天生適合 AI 的理解和交互場景，解決了理解歧義、數據不精確等痛點問題。它會取代傳統的接口文檔、口頭格式約定、手動版本同步等這些傳統協作方式，讓數據本身具備完整的自解釋能力，可獨立迭代演化且不影響各端的正常生成解析，適用於接口交互、配置管理、AI 數據對接等一切涉及數據流轉的場景。

注意:

- 編碼支持上限：65535 字節（64KB）。目前僅完成基礎類型支持，待後續文檔類型完全支持後，該限制可能會擴展
- 目前處於開發與測試階段，尚未達到生產穩定性，不建議在生產環境中使用
- API 與行為仍可能調整，使用時請留意版本更新

[github.com](https://github.com/metamessage/metamessage)

[pkg.go.dev](https://pkg.go.dev/github.com/metamessage/metamessage)

## 可相互轉換的文本格式

### JSONC

- 允許數組或對象的末尾字段以逗號結尾
- 允許普通註釋
- 註釋應寫在字段上方
- mm tag 必須在最後一行
- mm tag 需與普通註釋之間留一個空行以增強可讀性

**例**

```jsonc
{
  // mm: type=datetime; desc=創建時間
  "create_time": "2026-01-01 00:00:00",
}
```

### YAML

### TOML

## 數據類型

通过标签`type=`定義數據類型，格式`type=类型标识`，例如`type=i`表示整形

- doc: 文檔類型，未來為附加文檔預留，暫未啟用
- vec: 動態數組/切片，不允許復合類型
- arr: array，定長數組，不允許復合類型
- obj: object，對象/結構體，複合結構，對應多語言struct/object
- map: map，键值映射，鍵：僅字符串，值：不允許復合類型
- str: string，字符串
- bytes: 字節數組
- bool: 布爾值，取值：true/false，全小寫
- i: int，字面量不能包含小數點
- i8: int8
- i16: int16
- i32: int32
- i64: int64
- u: uint
- u8: uint8
- u16: uint16
- u32: uint32
- u64: uint64
- f32: float32，不支持 NaN / Inf / -0; 字面量必須帶小數點, 比如 0.0
- f64: float64，不支持 NaN / Inf / -0; 字面量必須帶小數點, 比如 0.0
- bigint: bigint
- datetime: 默認utc 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: 唯一標識
- decimal: 十進制小數，需傳入字符串
- ip: IP，支持 IPv4/IPv6
- url: 网址，符合标准URL格式
- email: 邮箱，符合标准邮箱格式
- enum: 枚舉，枚舉值是字符串，以`|`分割
- image: 圖片，底層是bytes
- video: 視頻，底層是bytes

## 標籤

標籤是编程语言结构体的注解、标签或属性，或是文本格式的注释

- is\_null: 值為null，並使用空值佔位
- desc: 摘要，適用所有類型。最大長度 65535 比特
- type: 數據類型。在文本格式中，字符串、整數（int）、小數（float64）、切片、對象（或類似結構）等沒有歧義時可以不用標注類型，比如當數組 size > 0 時不需要標記類型。在編程語言中，若數組、map 等可以判斷出來的類型，那麼也可以不用標注類型
- raw: 在一些編程語言中，數據類型通常使用包裝類型，如java。默認使用包裝類型，若不希望使用，可以設置為raw。待定，後續可能刪除此標籤
- nullable: 是否可為null，適用所有類型
- allow\_empty: 除布爾類型外，其他類型默認不允許為空，當設置allow\_empty後，可以為空，並允許通過一些規則。
- unique: 僅適用切片或數組，表示元素不可重複
- default: 默認值，尚未啟用
- example: 示例數據，用於數組、切片、map 類型為空時，自動生成的一個空值示例
- min: 在數組中表示最小容量，在字符串、字節數組中表示最小長度，在數字類型（整數、小數、bigint）種表示最小值
- max: 在數組中表示最大容量，在字符串、字節數組中表示最大長度，在數字類型（整數、小數、bigint）種表示最大值
- size: 在數組中表示容量，在字符串、字節數組中表示固定長度
- enum: 當有這個標籤時，默認此數值是enum類型。這裡的enum類型表現為字符串形式，不接受其他形式。
- pattern: 正則，適用於字符串。
- location: 時區偏移量，默認值 0，僅適用於時間類型，取值範圍：-12 ～ 14
- version: 在uuid中限定版本；在ip中可以限制ipv4或ipv6
- mime: 文檔類型，尚未啟用

## 使用方法

### CLI 工具

項目提供了一個命令行工具 `mm`，用於編碼、解碼和代碼生成。

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### 構建工具

```bash
make
```

#### 使用示例

1. **編碼 JSONC 到 MetaMessage**
   ```bash
   ./mm -encode -in input.jsonc -out output.mm
   ```
   或從 stdin 讀取：
   ```bash
   cat input.jsonc | ./mm -encode > output.mm
   ```
2. **解碼 MetaMessage 到 JSONC**
   ```bash
   ./mm -decode -in input.mm -out output.jsonc
   ```
   或從 stdin 讀取：
   ```bash
   cat input.mm | ./mm -decode > output.jsonc
   ```
3. **從 JSONC 生成結構體及數據代碼**

   支持 go, java, ts, kt, py, js, cs, rs, swift, php
   ```bash
   ./mm -generate -lang go -in input.jsonc -out output.go
   ```
   ```bash
   ./mm -generate -lang java -in input.jsonc -out output.java
   ```
   ```bash
   ./mm -generate -lang ts -in input.jsonc -out output.ts
   ```
   ```bash
   ./mm -generate -lang kt -in input.jsonc -out output.kt
   ```
   ```bash
   ./mm -generate -lang py -in input.jsonc -out output.py
   ```
   ```bash
   ./mm -generate -lang js -in input.jsonc -out output.js
   ```
   ```bash
   ./mm -generate -lang cs -in input.jsonc -out output.cs
   ```
   ```bash
   ./mm -generate -lang rs -in input.jsonc -out output.rs
   ```
   ```bash
   ./mm -generate -lang swift -in input.jsonc -out output.swift
   ```
   ```bash
   ./mm -generate -lang php -in input.jsonc -out output.php
   ```

#### 選項說明

- -encode, -e: 編碼模式
- -decode, -d: 解碼模式
- -generate, -g: 生成代碼模式
- -in, -i: 輸入文件路徑（空則從 stdin 讀取）
- -out, -o: 輸出文件路徑（空則輸出到 stdout）
- -force, -f: 強制覆蓋輸出文件
- -lang, -l: 生成目標語言（僅用於 generate 模式，支持 go, java, ts, kt, py, js, cs, rs, swift, php）

### 庫使用

項目提供 Go 庫用於程序調用。

#### 安裝

```bash
go get github.com/metamessage/metamessage
```

#### 示例代碼

```go
package main

import (
    "fmt"
    mm "github.com/metamessage/metamessage"
)

func main() {
    // 從結構體編碼
    type Person struct {
        Name  string
        Age   int
    }

    p := Person{Name: "Alice", Age: 30}
    data, err := mm.EncodeFromValue(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    // 解碼到結構體
    var decoded Person
    err = mm.DecodeToValue(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    // 從 JSONC 編碼
    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := mm.EncodeFromJsonc(jsoncStr)
    if err != nil {
        panic(err)
    }

    // 解碼到 JSONC
    jsoncOut, err := mm.DecodeToJsonc(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API 概覽

- `NewEncoder(w io.Writer) Encoder`: 創建編碼器
- `EncodeFromValue(in any) ([]byte, error)`: 從結構體編碼
- `EncodeFromJsonc(in string) ([]byte, error)`: 從 JSONC 字符串編碼
- `NewDecoder(r io.Reader) Decoder`: 創建解碼器
- `DecodeToValue(in []byte, out any) error`: 解碼到結構體
- `DecodeToJsonc(in []byte) (string, error)`: 解碼到 JSONC 字符串

### 其他語言示例

#### Java

[![jitpack.io](https://jitpack.io/v/metamessage/metamessage.svg)](https://jitpack.io/#metamessage/metamessage)

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

#### Kotlin

[![jitpack.io](https://jitpack.io/v/metamessage/metamessage.svg)](https://jitpack.io/#metamessage/metamessage)

```kotlin
import io.github.metamessage.MetaMessage
import io.github.metamessage.MM

@MM(desc="person")
class Person(var name: String = "Ed", var age: Uint8 = 30.toUint8())

fun main() {
    val person = Person()

    val wire = MetaMessage.encodeFromValue(person)

    val person = MetaMessage.decodeToValue(wire, Person::class.java)

    val jsonc = MetaMessage.valueToJsonc(person)

    val person = MetaMessage.jsoncToValue(jsoncOutput, Person::class.java)

    val wire = MetaMessage.encodeFromJsonc(jsonc)

    val jsonc = MetaMessage.decodeToJsonc(wire)
}
```

#### TypeScript

```typescript
import { encodeFromValue, decodeToValue, mm, ValueType } from "metamessage";

@mm({ desc: "" })
class Person {
  @mm({ desc: "" })
  name: string = "";
  @mm({ desc: "" })
  age: number = 0;
}
const person = { name: "Ed", age: 30 };
const wire = encodeFromValue(person);
const decoded = decodeToValue(wire, Person);
```

#### Python

```python
from metamessage import encode, decode

person = {"name": "Ed", "age": 30}
wire = encode(person)
decoded = decode(wire)
```

#### JavaScript

```javascript
const { encode, decode } = require("metamessage");

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

#### C\#

```csharp
using MetaMessage;

var person = new Person { Name = "Ed", Age = 30 };
byte[] wire = MetaMessage.Encode(person);
var decoded = MetaMessage.Decode<Person>(wire);
```

#### Rust

```rust
use metamessage::{encode, decode, Node};

let person = Node::Object(/* ... */);
let wire = encode(&person);
let decoded = decode(&wire).unwrap();
```

#### Swift

```swift
import MetaMessage

let person = Person(name: "Ed", age: 30)
let wire = MetaMessage.encodeFromValue(person)
let decoded = try MetaMessage.decodeToValue(wire)
```

#### PHP

```php
<?php
use io\metamessage\mm\MetaMessage;

$person = new Person();
$wire = MetaMessage::encode($person);
$decoded = MetaMessage::decode($wire, Person::class);
```

### 示例

查看 `examples/` 目錄中的示例代碼。
