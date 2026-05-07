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

MetaMessage（mm）は構造化データ交換プロトコルです。自己記述、自己制約、自己例示を備え、損失のないデータ交換を実現します。AI、人間、機械にネイティブに対応する次世代の汎用プロトコルです。

- 人間とAIにやさしい
- JSONCへのエクスポート/インポートをサポート（現時点ではJSONC、将来的にはYAML/TOMLなどを予定）
- 設定ファイルやデータ交換に適している
- 従来のAPIとAIインタラクションの両方に対応
- 各言語の構造体/クラスとMetaMessageの相互変換をサポート
- 複数言語向けのコード生成をサポート
- データ自身が型、制約、説明、例を持ち、追加のドキュメントが不要
- 型、制約、説明、例はデータとともに更新され、別途同期する必要がない
- 言語間でデータ構造と値が完全に一致する
- 構造の損失がなく、パーサーは自動適応し、クラッシュしない
- コンパクトなバイナリにシリアライズでき、解析が高速でサイズが小さい

**解決する課題**

- uint8であるかどうかなど型が不明確
- nullの内部構造が判別できず構造が不完全になる
- 制約ルールがなくデータの正当性を検証できない
- 例や説明がなく、別ドキュメントに依存する必要がある
- フォーマット変更に敏感で、エンコード/デコードやドキュメントの再同期が必要

MetaMessageはAIの理解と対話に適しており、曖昧さや不正確さを解消します。従来のインターフェースドキュメント、口頭でのフォーマット合意、手動のバージョン同期を置き換え、データ自体が自己説明的に進化できるようにします。

注意: 現在開発・テスト中であり、本番環境での使用は推奨されません

[meta-message](https://github.com/metamessage/metamessage)

## テキスト形式

### JSONC

- 配列やオブジェクトの末尾にカンマを許可します
- 通常のコメントが許可されます
- フィールドの上にコメントを書くべきです
- mmタグは最終行に配置する必要があります
- mmタグと通常コメントの間に空行を入れて可読性を高めます

**例**

```jsonc
{
    // mm: type=datetime; desc=作成日時
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## データ型

- doc: エンコードは最大65535バイト（64KB）までサポートします。ドキュメント型の完全なサポート後、この制限を超える可能性があります
- slice: 配列、スライスは複合型を許可しません
- array: arr
- struct:
- map: mapのキーは文字列である必要があり、値は複合型であってはいけません
- string: str
- bytes:
- bool:
- int: i; 整数リテラルに小数点を含めてはいけません
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32; 浮動小数点はNaN/Inf/-0をサポートしません; 浮動小数点リテラルには小数点が必要です。例：0.0
- float64: f64
- bigint: bi
- datetime: デフォルトUTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid
- decimal
- ip
- url
- email
- enum
- image
- video

## タグ

タグはプログラミング言語構造体の注釈、ラベル、または属性、またはテキスト形式のコメントです

- is_null: null値で空のプレースホルダを使用します

- desc: 概要。すべてのタイプに適用。最大長 65535 ビット

- type: データタイプ。テキスト形式では、曖昧でない場合は文字列、整数（int）、小数（float64）、スライス、オブジェクト（または類似構造）などに明示的な型タグは不要です。例えば、配列 size > 0 の場合は型タグ不要。プログラミング言語では、配列、map などが判別できれば、型タグ不要です

- raw: 一部のプログラミング言語では、データ型は通常 Java などのラッパータイプを使用します。デフォルトではラッパータイプを使用しますが、使用したくない場合は raw に設定します。決定待機中、将来的に削除される可能性があります

- nullable: null が許可されるかどうか。すべてのタイプに適用

- allow_empty: ブール型を除き、他のタイプはデフォルトで空を許可しません。allow_empty が設定されると、いくつかのルールに従って空を許可します

- unique: スライスまたは配列にのみ適用。要素の重複がないことを示します

- default: デフォルト値。まだ有効ではありません

- example: 配列、スライス、map が空の場合に使用されるサンプルデータで、自動的に空値の例が生成されます

- min: 配列の最小容量、文字列やバイト配列の最小長、数値（整数、小数、bigint）の最小値

- max: 配列の最大容量、文字列やバイト配列の最大長、数値（整数、小数、bigint）の最大値

- size: 配列の容量、文字列やバイト配列の固定長

- enum: このタグが存在する場合、値はデフォルトでenum型です。ここのenum型は文字列形式で表現され、他の形式を受け付けません

- pattern: 正規表現。文字列に適用します

- location: タイムゾーンオフセット。デフォルト0、datetime型にのみ適用、範囲-12～14

- version: uuidでバージョンを制限します。ipではipv4またはipv6を制限できます

- mime: ドキュメントタイプ。まだ有効ではありません

## 使い方

### CLIツール

このプロジェクトはエンコード、デコード、コード生成のためのコマンドラインツール `mm` を提供します。

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### ビルド

```bash
make
```

#### 使用例

1. JSONCをMetaMessageにエンコード

```bash
./mm -encode -in input.jsonc -out output.mm
```

またはstdinから読み取る：

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. MetaMessageをJSONCにデコード

```bash
./mm -decode -in input.mm -out output.jsonc
```

またはstdinから読み取る：

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. JSONCから構造体およびコードを生成

go、java、ts、kt、py、js、cs、rs、swift、phpをサポート

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

#### オプション

- -encode, -e: エンコードモード
- -decode, -d: デコードモード
- -generate, -g: コード生成モード
- -in, -i: 入力ファイルパス（空の場合stdinから読み取る）
- -out, -o: 出力ファイルパス（空の場合stdoutへ出力）
- -force, -f: 出力ファイルを強制上書き
- -lang, -l: 生成対象言語（go, java, ts, kt, py, js, cs, rs, swift, php）

### ライブラリ利用

このプロジェクトはプログラムで使用するためのGoライブラリを提供します。

#### インストール

```bash
go get github.com/metamessage/metamessage
```

#### サンプルコード

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
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    var decoded Person
    err = mm.DecodeToValue(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := mm.EncodeFromJSONC(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := mm.DecodeToJSONC(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API概要

- `NewEncoder(w io.Writer) Encoder`: エンコーダーを作成
- `EncodeFromValue(in any) ([]byte, error)`: 構造体からエンコード
- `EncodeFromJSONC(in string) ([]byte, error)`: JSONC文字列からエンコード
- `NewDecoder(r io.Reader) Decoder`: デコーダーを作成
- `DecodeToValue(in []byte, out any) error`: 構造体へデコード
- `DecodeToJSONC(in []byte) (string, error)`: JSONC文字列へデコード

### 他の言語の例

#### Java

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

#### Kotlin

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

#### TypeScript

```typescript
import { encode, decode } from 'metamessage';

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
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
const { encode, decode } = require('metamessage');

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
let wire = MetaMessage.encode(person)
let decoded = try MetaMessage.decode(wire)
```

#### PHP

```php
<?php
use io\metamessage\mm\MetaMessage;

$person = new Person();
$wire = MetaMessage::encode($person);
$decoded = MetaMessage::decode($wire, Person::class);
```

### 例

`examples/`ディレクトリのサンプルコードを参照してください。
