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

注意:

- 現在、エンコードは最大65535バイト（64KB）までサポートしています。ドキュメント型の完全サポート後にこの制限が拡張される可能性があります
- 本プロジェクトは開発およびテスト段階にあり、まだ本番環境での使用は推奨されません
- APIや挙動は変更される可能性があるため、バージョン更新にご注意ください

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

`type=` タグを使ってデータ型を定義します。形式は `type=typeIdentifier` です。例：`type=i` は整数を意味します。

- doc: エンコードは最大65535バイト（64KB）までサポートします。ドキュメント型の完全なサポート後、この制限を超える可能性があります
- vec: 配列/スライス（動的配列）、複合型を許可しません
- arr: 配列（固定長）、複合型を許可しません
- obj: オブジェクト/構造体、複合構造で多言語の struct/object に対応
- map: map のキーは文字列であり、値は複合型であってはいけません
- str: string
- bytes: バイト配列
- bool: boolean
- i: int; 整数リテラルに小数点を含めてはいけません
- i8: int8
- i16: int16
- i32: int32
- i64: int64
- u: uint
- u8: uint8
- u16: uint16
- u32: uint32
- u64: uint64
- f32: float32; 浮動小数点は NaN/Inf/-0 をサポートしません。浮動小数点リテラルには小数点が必要です。例：0.0
- f64: float64; 浮動小数点は NaN/Inf/-0 をサポートしません。浮動小数点リテラルには小数点が必要です。例：0.0
- bigint: bigint
- datetime: デフォルトUTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: 一意の識別子
- decimal: 10進小数。文字列として渡す必要があります
- ip: IP、IPv4/IPv6 をサポート
- url: URL、標準URL形式
- email: メールアドレス、標準形式
- enum: 列挙型。値は `|` で区切られた文字列です
- image: 画像。内部は bytes
- video: 動画。内部は bytes

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

他の言語でも同じロジックが適用されます

```go
// 誤った書き方
// ID はすでにint64型なので、type=i64 の指定は不要
// mmではjsonタグを併用しないこと。mmが自動的に命名を処理します。名前を指定する必要がある場合は mm で name= を使用
// Age はネイティブのuint8型を使用すれば type=u8 を省略可能
type User struct {
	ID       int64  `mm:"type=i64;desc=User ID" json:"id"`
	Name     string `mm:"type=str;desc=User Name;min=1;max=50" json:"name"`
	Email    string `mm:"type=email;desc=Email" json:"email"`
	Age      int    `mm:"type=u8;desc=Age;min=0;max=150" json:"age"`
	IsActive bool   `mm:"type=bool;desc=Is Active" json:"is_active"`
}

// 正しい書き方
// Email にはネイティブ型がないため、type=email の指定が必要
type User struct {
	ID       int64  `mm:"desc=User ID"`
	Name     string `mm:"desc=User Name;min=1;max=50"`
	Email    string `mm:"type=email;desc=Email"`
	Age      uint8  `mm:"desc=Age;min=0;max=150"`
	IsActive bool   `mm:"desc=Is Active"`
}

user := User{}

// ルートレベルのタグはここで指定可能
tag := "desc=User"
_, _ = EncodeFromValue(user, tag)
```

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
    data2, err := mm.EncodeFromJsonc(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := mm.DecodeToJsonc(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API概要

- `NewEncoder(w io.Writer) Encoder`: エンコーダーを作成
- `EncodeFromValue(in any) ([]byte, error)`: 構造体からエンコード
- `EncodeFromJsonc(in string) ([]byte, error)`: JSONC文字列からエンコード
- `NewDecoder(r io.Reader) Decoder`: デコーダーを作成
- `DecodeToValue(in []byte, out any) error`: 構造体へデコード
- `DecodeToJsonc(in []byte) (string, error)`: JSONC文字列へデコード

### 他の言語の例

#### Java

[jitpack.io](https://jitpack.io/#metamessage/metamessage/)

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

[jitpack.io](https://jitpack.io/#metamessage/metamessage/)

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
import { encodeFromValue, decodeToValue, mm, ValueType } from 'metamessage';

@mm({ desc: '' })
class Person {
    @mm({ desc: '' })
    name: string = ''
    @mm({ desc: '' })
    age: number = 0
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

### 例

`examples/`ディレクトリのサンプルコードを参照してください。
