# Swift MetaMessage (mm-swift) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 Swift，排除 CLI 功能，測試要完整全面，包括所有數據類型和 tag。

## 當前進度
- ✅ Package.swift - 項目結構已完成
- ✅ MM/Constants.swift - Wire 格式常量已完成
- 🔄 MM/MMEncoder.swift - 編碼器（部分完成，需完善 array 類型）
- ✅ MM/MMDecoder.swift - 解碼器已完成
- ❌ JSONC 模块 - 未開始

## 項目結構
```
mm-swift/
├── Package.swift
├── Sources/
│   └── MetaMessage/
│       └── MM/
│           ├── Constants.swift
│           ├── Prefix.swift
│           ├── SimpleValue.swift
│           ├── MMBuffer.swift
│           ├── MMEncoder.swift
│           ├── MMDecoder.swift
│           └── MetaMessage.swift
│       └── JSONC/
│           ├── ValueType.swift
│           ├── JSONCScanner.swift
│           ├── JSONCAST.swift
│           ├── JSONCParser.swift
│           ├── JSONCTag.swift
│           ├── JSONCBinder.swift
│           └── JSONCPrinter.swift
└── Tests/
    └── MetaMessageTests/
        ├── MMEncoderTests.swift
        ├── MMDecoderTests.swift
        ├── JSONCScannerTests.swift
        ├── JSONCParserTests.swift
        ├── JSONCTagTests.swift
        └── MetaMessageIntegrationTests.swift
```

## 待完成任務

### 1. 完成 MMEncoder.swift 編碼器
**文件**: `Sources/MetaMessage/MM/MMEncoder.swift`

需要實現以下類型的 Array 編碼：
- `encodeArray(_ array: [Bool])`
- `encodeArrayStrings(_ array: [String])`
- `encodeArrayInt(_ array: [Int])`
- `encodeArrayUInt(_ array: [UInt])`
- `encodeArrayFloat(_ array: [Float])`
- `encodeArrayDouble(_ array: [Double])`
- `encodeArrayData(_ array: [Data])`

每個 Array 編碼包含：
1. 計算元素數量（使用 PositiveInt tag）
2. 計算總長度
3. 編碼長度
4. 遞歸編碼每個元素

### 2. 完成 JSONC 模块

#### 2.1 ValueType.swift - 值類型枚舉
**文件**: `Sources/MetaMessage/JSONC/ValueType.swift`
```swift
public enum ValueType {
    case null
    case bool
    case int
    case uint
    case float
    case string
    case array
    case object
}
```

#### 2.2 JSONCScanner.swift - 詞法分析器
**文件**: `Sources/MetaMessage/JSONC/JSONCScanner.swift`

功能：
- 支持 `{`, `}`, `[`, `]`, `:`, `,` 符號
- 支持字符串（雙引號）
- 支持數字（整數、浮點數）
- 支持 `true`, `false`, `null`
- 支持行註釋 `//`
- 支持塊註釋 `/* */`
- 追蹤行號和列位置
- 跳過註釋和空白字符

#### 2.3 JSONCAST.swift - AST 節點定義
**文件**: `Sources/MetaMessage/JSONC/JSONCAST.swift`

節點類型：
- `JSONCNode` - 基礎協議
- `ValueNode` - 值節點（包含具體值）
- `ObjectNode` - 對象節點（鍵值對集合）
- `ArrayNode` - 數組節點
- `CommentNode` - 註釋節點
- `CommentType` - 註釋類型（行、塊）

#### 2.4 JSONCParser.swift - 語法分析器
**文件**: `Sources/MetaMessage/JSONC/JSONCParser.swift`

功能：
- 解析 JSONC 文法（支持註釋）
- 遞歸下降解析
- 創建 AST
- 錯誤報告

#### 2.5 JSONCTag.swift - Tag 定義和解析
**文件**: `Sources/MetaMessage/JSONC/JSONCTag.swift`

功能：
- TagMetadata 結構（type, desc, optional 等）
- 解析 `// mm:type=str;desc=...` 格式
- 從註釋中提取 tag 信息

#### 2.6 JSONCBinder.swift - 對象綁定器
**文件**: `Sources/MetaMessage/JSONC/JSONCBinder.swift`

功能：
- 將 AST 綁定到 Swift 對象
- 使用 Mirror 反射機制
- 類型轉換和驗證
- 支持嵌套對象和數組

#### 2.7 JSONCPrinter.swift - 輸出格式化器
**文件**: `Sources/MetaMessage/JSONC/JSONCPrinter.swift`

功能：
- 將 AST 轉換為 JSONC 字符串
- 格式化輸出（縮進）
- 處理註釋輸出

### 3. MetaMessage.swift - 公共 API
**文件**: `Sources/MetaMessage/MM/MetaMessage.swift`

提供統一的消息編碼/解碼接口

### 4. 測試模塊

#### 4.1 MMEncoderTests.swift - 編碼器測試
測試所有數據類型：
- Bool
- Int (Int8, Int16, Int32, Int64)
- UInt (UInt8, UInt16, UInt32, UInt64)
- Float
- Double
- String
- Data
- Array (所有類型)
- Object (嵌套)

#### 4.2 MMDecoderTests.swift - 解碼器測試
測試所有前綴和值的解碼

#### 4.3 JSONCScannerTests.swift - 詞法分析器測試
- 測試所有 Token 類型
- 測試註釋跳過
- 測試位置追蹤
- 測試錯誤處理

#### 4.4 JSONCParserTests.swift - 語法分析器測試
- 測試有效 JSONC 輸入
- 測試帶註釋的 JSONC
- 測試嵌套結構
- 測試錯誤恢復

#### 4.5 JSONCTagTests.swift - Tag 解析測試
- 測試標準 tag 格式
- 測試多個 tag
- 測試特殊字符

#### 4.6 MetaMessageIntegrationTests.swift - 集成測試
- 完整的消息編碼/解碼循環
- 邊界條件測試
- 性能測試

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`

## 實現順序
1. 完成 MMEncoder.swift 的 array 編碼
2. 實現 JSONCScanner.swift
3. 實現 JSONCAST.swift
4. 實現 JSONCParser.swift
5. 實現 JSONCTag.swift
6. 實現 JSONCBinder.swift
7. 實現 JSONCPrinter.swift
8. 實現 MetaMessage.swift
9. 實現所有測試
10. 運行測試驗證

## 關鍵注意事項
- 使用 Swift 反射 Mirror 進行類型檢查
- 錯誤處理使用 Result 或 throws
- 遵循 Swift API 設計規範
- 確保測試覆蓋所有數據類型和邊界條件