# MetaMessage

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

**例**

```jsonc
{
    // mm: type=datetime; desc=作成日時
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## データ変換

JSONC、YAML、TOMLなどのテキスト形式への出力をサポートします。

**JSONC**

- 配列やオブジェクトの末尾にカンマを許可します

推奨コメントスタイル：

- 通常のコメントが許可されます
- フィールドの上にコメントを書くべきです
- mmタグは最終行に配置する必要があります
- mmタグと通常コメントの間に空行を入れて可読性を高めます

## 注意

- まだ多くのバグがあり、テストは不完全です。プロダクション利用は推奨しません
- 配列やスライスは複合型を許可しません。mapのキーは文字列で、値は複合型を許可しません
- 空の配列/スライスには自動的に例が挿入されます
- 整数や文字列は明示的な型指定が不要です
- 構造体やスライスは明示的な型指定が不要です
- 配列サイズが0より大きい場合、型指定は不要です
- 浮動小数点はNaN/Inf/-0をサポートしません
- エンコードは65535バイト（64KB）までサポートします。将来的に拡張される可能性があります
- 浮動小数点リテラルには小数点が必要です
- 整数リテラルに小数点を含めてはいけません

## データ型

datetime: デフォルトUTC 1970-01-01 00:00:00

## タグ

- is_null: null値で空のプレースホルダを使用します
- example: 配列やmapが空の場合に使用されるサンプルデータ
- min: 配列の最小容量、文字列やバイト配列の最小長、数値の最小値
- max: 配列の最大容量、文字列やバイト配列の最大長、数値の最大値
- size: 配列、文字列、バイト配列の固定長
- location: タイムゾーンオフセット。デフォルト0、範囲-12～14

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
go get github.com/metamessage/metamessage/pkg
```

#### サンプルコード

```go
package main

import (
    "fmt"
    "github.com/metamessage/metamessage/pkg"
)

func main() {
    type Person struct {
        Name string
        Age  int
    }

    p := Person{Name: "Alice", Age: 30}
    data, err := pkg.EncodeFromStruct(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    var decoded Person
    err = pkg.Decode(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := pkg.EncodeFromJSONC(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := pkg.DecodeToJSONC(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API概要

- `NewEncoder(w io.Writer) Encoder`: エンコーダーを作成
- `EncodeFromStruct(in any) ([]byte, error)`: 構造体からエンコード
- `EncodeFromJSONC(in string) ([]byte, error)`: JSONC文字列からエンコード
- `NewDecoder(r io.Reader) Decoder`: デコーダーを作成
- `Decode(in []byte, out any) error`: 構造体へデコード
- `DecodeToJSONC(in []byte) (string, error)`: JSONC文字列へデコード

### 例

`examples/`ディレクトリのサンプルコードを参照してください。
