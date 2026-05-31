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

通过标签`type=`定義數據類型，格式`type=类型标识`，例如`type=i8`表示int8類型

- doc: 文檔類型，未來為附加文檔預留，暫未啟用
- obj: object，對象/結構體，複合結構，對應多語言struct/object
- map: map，键值映射，鍵：僅字符串，值：不允許復合類型
- vec: 動態數組/切片，不允許復合類型
- arr: array，定長數組，不允許復合類型
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
- bigint: bigint
- f32: float32，不支持 NaN / Inf / -0; 字面量必須帶小數點, 比如 0.0
- f64: float64，不支持 NaN / Inf / -0; 字面量必須帶小數點, 比如 0.0
- decimal: 十進制小數，需傳入小數, 比如 0.0
- datetime: 默認utc 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: 唯一標識
- str: string，字符串
- bytes: 字節數組
- ip: IP，支持 IPv4/IPv6
- url: 网址，符合标准URL格式
- email: 邮箱，符合标准邮箱格式
- enums: 枚舉，枚舉值是字符串，以`|`分割
- media: 媒體，底層是bytes

### jsonc 示例

标量

```jsonc

// doc TODO

// obj，對象/結構體，複合結構，對應多語言struct/object，不需要標注類型
// mm: desc=對象/結構體
{
    "name": "Ed",
    "age": 30
}

// map，键值映射，鍵：僅字符串，值：不允許復合類型
// mm: desc=map映射; type=map
{
    "name": "Ed",
    "age": "30"
}

// vec，動態數組，不需要標注類型
// mm: desc=動態數組
[1, 2, 3]

// arr，定長數組，需要標注類型
// mm: desc=定長數組; type=arr
[1, 2, 3]

// 只允許小寫false、true
// mm: desc=布爾值
false

// 數字；數字不能帶點；默認是int，不需要標注類型
// mm: desc=int
0

// mm: desc=int8; type=i8
0

// mm: desc=int16; type=i16
0

// mm: desc=int32; type=i32
0

// mm: desc=int64; type=i64
0

// mm: desc=uint; type=u
0

// mm: desc=uint8; type=u8
0

// mm: desc=uint16; type=u16
0

// mm: desc=uint32; type=u32
0

// mm: desc=uint64; type=u64
0

// 表現為數字，但需要標注類型；
// mm: desc=bigint; type=bigint
0

// mm: desc=float32; type=f32
0.0

// 數字帶點；float必須帶點；默認是float64，不需要標注類型；
// mm: desc=float64
0.0

// 表現為小數，但需要標注類型；
// mm: desc=decimal; type=decimal
0.0

// 雙引號包裹，默認是字符串，不需要標注類型；
// mm: desc=字符串
""

// 表現為字符串，但需要標注類型；
// mm: desc=字節數組; type=bytes
""

// 表現為字符串，但需要標注類型；
// mm: desc=datetime; type=datetime
""

// 表現為字符串，但需要標注類型；
// mm: desc=date; type=date
""

// 表現為字符串，但需要標注類型；
// mm: desc=time; type=time
""

// 表現為字符串，但需要標注類型；
// mm: desc=uuid; type=uuid
""

// 表現為字符串，但需要標注類型；
// mm: desc=ip; type=ip
""

// 表現為字符串，但需要標注類型；
// mm: desc=url; type=url
""

// 表現為字符串，但需要標注類型；
// mm: desc=email; type=email
""

// 表現為字符串，但需要標注類型；
// mm: desc=enums; type=enums
""

// 表現為字符串，但需要標注類型；
// mm: desc=media; type=media
""
```

## 標籤

標籤是编程语言结构体的注解、标签或属性，或是文本格式的注释

- is_null: 值為null，並使用零值佔位。在jsonc等文本格式中，使用零值佔位
- example: 示例數據，用於數組、切片、map 類型為空時，自動生成的一個空值示例
- deprecated: 廢棄，不建議使用
- name: 僅用於各語言的數據對象，用於命名。在jsonc等文本格式中，此標籤無效，不建議使用。
- desc: 摘要，適用所有類型。最大長度 65535 比特
- type: 數據類型。在文本格式中，字符串、整數（int）、小數（float64）、切片、對象（或類似結構）等沒有歧義時可以不用標注類型，比如當數組 size > 0 時不需要標記類型。在編程語言中，若數組、map 等可以判斷出來的類型，那麼也可以不用標注類型
- nullable: 是否可為null，適用所有類型
- allow_empty: 除布爾類型外，其他類型默認不允許為空，當設置allow_empty後，可以為空，並允許通過一些規則。
- unique: 僅適用切片或數組，表示元素不可重複
- default_val: 默認值，尚未啟用
- min: 在數組中表示最小容量，在字符串、字節數組中表示最小長度，在數字類型（整數、小數、bigint）種表示最小值
- max: 在數組中表示最大容量，在字符串、字節數組中表示最大長度，在數字類型（整數、小數、bigint）種表示最大值
- size: 在數組中表示容量，在字符串、字節數組中表示固定長度
- enums: 當有這個標籤時，默認此數值是enum類型。這裡的enum類型表現為字符串形式，不接受其他形式。
- pattern: 正則，適用於字符串。
- location: 時區偏移量，默認值 0，僅適用於時間類型，取值範圍：-12 ～ 14
- version: 在uuid中限定版本；在ip中可以限制ipv4或ipv6
- mime: 文檔類型
- child_desc:
- child_type:
- child_nullable:
- child_allow_empty:
- child_unique:
- child_default_val:
- child_min:
- child_max:
- child_size:
- child_enums:
- child_pattern:
- child_location:
- child_version:
- child_mime:

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

其他語言相同邏輯

```go
// 錯誤寫法
// ID 已經有int64類型，不需要再標注類型i64
// 在mm中，不建議同時使用json標籤。一般mm會自動處理，如需指定name，請在mm中`name=`指定
// Age 建議直接使用原生類型uin8，可以省略type=u8
type User struct {
	ID       int64  `mm:"type=i64;desc=用戶ID" json:"id"`
	Name     string `mm:"type=str;desc=用戶名稱;min=1;max=50" json:"name"`
	Email    string `mm:"type=email;desc=電子郵箱" json:"email"`
	Age      int    `mm:"type=u8;desc=年齡;min=0;max=150" json:"age"`
	IsActive bool   `mm:"type=bool;desc=是否激活" json:"is_active"`
}

// 正確寫法
// Email由於沒有原生類型，所以需要使用type=email標籤指定
type User struct {
	ID       int64  `mm:"desc=用戶ID"`
	Name     string `mm:"desc=用戶名稱;min=1;max=50"`
	Email    string `mm:"type=email;desc=電子郵箱"`
	Age      uin8   `mm:"desc=年齡;min=0;max=150"`
	IsActive bool   `mm:"desc=是否激活"`
}

user := User{}

// 這裡可以指定最外層標籤
tag := "desc=用戶"
_, _ = EncodeFromValue(user, tag) {
```

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

[npmjs.com](https://www.npmjs.com/package/metamessage)

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

[pypi.org](https://pypi.org/project/metamessage/)

[mm-py](./mm-py)

#### JavaScript

[npmjs.com](https://www.npmjs.com/package/metamessage)

```javascript
const { encode, decode } = require("metamessage");

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

#### C#

[NuGet](https://www.nuget.org/packages/MetaMessage)

[mm-cs](./mm-cs)

#### Rust

[crates.io](https://crates.io/crates/metamessage)

[mm-rs](./mm-rs)

#### Swift

[mm-swift](./mm-swift)

#### PHP

[packagist.org](https://packagist.org/packages/metamessage/metamessage)

[mm-php](./mm-php)

### 示例

查看 `examples/` 目錄中的示例代碼。

## 測試

### 跨語言一致性測試

運行所有可用語言的 harness 對全部 fixtures 進行解析，並比較輸出是否一致：

```bash
./tests/run_cross_lang.sh
```

### 單個語言測試

運行某個語言的 harness 對單個 fixture 進行解析測試：

```bash
# Go（無需構建）
go run ./tests/harness/go/harness.go tests/fixtures/01_primitive/boolean.jsonc

# Python（無需構建）
python3 ./tests/harness/python/harness.py tests/fixtures/01_primitive/boolean.jsonc

# TypeScript（需先構建）
cd mm-ts && npm run build --silent && cd - && \
node ./tests/harness/typescript/harness.cjs tests/fixtures/01_primitive/boolean.jsonc

# Rust
cd tests/harness/rust && cargo run -- ../../fixtures/01_primitive/boolean.jsonc

# C（需先構建）
mkdir -p tests/harness/c/build && cd tests/harness/c/build && \
cmake .. -DCMAKE_BUILD_TYPE=Release >/dev/null && make -j4 >/dev/null && \
./mm_harness_c ../../../fixtures/01_primitive/boolean.jsonc

# C++（需先構建）
g++ -std=c++17 -I mm-cpp/src -o tests/harness/cpp/build/mm_harness_cpp \
  tests/harness/cpp/harness.cpp mm-cpp/src/jsonc/scanner.cpp && \
./tests/harness/cpp/build/mm_harness_cpp tests/fixtures/01_primitive/boolean.jsonc

# C#（需先構建）
dotnet build tests/harness/csharp/harness.csproj --nologo -v q && \
dotnet run --project tests/harness/csharp/harness.csproj --no-build -- \
  tests/fixtures/01_primitive/boolean.jsonc

# Kotlin（需先構建 mm-kt）
cd mm-kt && mvn compile -q -DskipTests && cd - && \
java -cp "$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-compiler/1.9.22/kotlin-compiler-1.9.22.jar:\
$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar:\
$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar:\
$HOME/.m2/repository/org/jetbrains/intellij/deps/trove4j/1.0.20221201/trove4j-1.0.20221201.jar:\
$HOME/.m2/repository/org/jetbrains/annotations/13.0/annotations-13.0.jar" \
  org.jetbrains.kotlin.cli.jvm.K2JVMCompiler \
  tests/harness/kotlin/harness.kt \
  -cp "mm-kt/target/classes:\
$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar:\
$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar" \
  -d /tmp/harness.jar -no-stdlib -no-reflect && \
java -cp "/tmp/harness.jar:mm-kt/target/classes:\
$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar" \
  HarnessKt tests/fixtures/01_primitive/boolean.jsonc

# Swift
cd tests/harness/swift && swift run ../../fixtures/01_primitive/boolean.jsonc
```
