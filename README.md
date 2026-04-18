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

- 可以導出和導入為 JSONC（截止目前。未來支持 YAML / TOML 等更多格式），方便閱讀

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

**例**  

```jsonc
{
    // mm: type=datetime; desc=創建時間
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## 數據轉換

支持輸出為 JSONC、YAML、TOML 等文本格式

**JSONC**  

- 允許數組或對象的末尾字段以逗號結尾

推薦的註釋風格：

- 允許普通註釋

- 註釋應寫在字段上方

- mm tag 必須在最後一行

- mm tag 需與普通註釋之間留一個空行以增強可讀性

## 注意事項

- 目前還有很多bug待修復，測試也不完整，不建議生產使用

- 數組、切片不允許復合類型；map中鍵值不允許非字符串，值不允許復合類型

- 空數組、切片，會自動插入一條示例數據

- 整數、字符串不需要標記類型

- 結構體、切片不需要標記類型

- 當數組 size > 0 時不需要標記類型

- 小數不支持NaN/Inf/-0

- 編碼支持上限：65535 字節（64KB）。待後續文檔類型完全支持後，可突破此限制

- 浮點數字面量必須帶小數點

- 整數字面量不能帶小數點

## 數據類型

datetime: 默認utc 1970-01-01 00:00:00

## 標籤

- is_null: 值為null，並使用空值佔位

- example: 示例數據，用於數組、map 類型為空時，自動生成的一個空值示例

- min: 在數組中表示最小容量，在字符串、字節數組中表示最小長度，在數字類型（整數、小數、bigint）種表示最小值

- max: 在數組中表示最大容量，在字符串、字節數組中表示最大長度，在數字類型（整數、小數、bigint）種表示最大值

- size: 在數組中表示容量，在字符串、字節數組中表示固定長度

- location: 時區偏移量，默認值 0，取值範圍：-12 ～ 14

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
go get github.com/metamessage/metamessage/pkg
```

#### 示例代碼

```go
package main

import (
    "fmt"
    "github.com/metamessage/metamessage/pkg"
)

func main() {
    // 從結構體編碼
    type Person struct {
        Name  string 
        Age   int   
    }
    
    p := Person{Name: "Alice", Age: 30}
    data, err := pkg.EncodeFromStruct(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %xn", data)
    
    // 解碼到結構體
    var decoded Person
    err = pkg.Decode(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+vn", decoded)
    
    // 從 JSONC 編碼
    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := pkg.EncodeFromJSONC(jsoncStr)
    if err != nil {
        panic(err)
    }
    
    // 解碼到 JSONC
    jsoncOut, err := pkg.DecodeToJSONC(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API 概覽

- `NewEncoder(w io.Writer) Encoder`: 創建編碼器

- `EncodeFromStruct(in any) ([]byte, error)`: 從結構體編碼

- `EncodeFromJSONC(in string) ([]byte, error)`: 從 JSONC 字符串編碼

- `NewDecoder(r io.Reader) Decoder`: 創建解碼器

- `Decode(in []byte, out any) error`: 解碼到結構體

- `DecodeToJSONC(in []byte) (string, error)`: 解碼到 JSONC 字符串

### 示例

查看 `examples/` 目錄中的示例代碼。
