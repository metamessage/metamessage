# Java MetaMessage (mm-java) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 Java，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成 (22 tests)
- ✅ JSONC 基本對齊完成 (15 tests)
- ✅ JSONC MM 對齊完成 (22 tests)

## 項目結構
```
mm-java/
├── src/main/java/io/metamessage/
│   ├── mm/
│   │   ├── Constants.java
│   │   ├── Types.java
│   │   ├── Prefix.java
│   │   ├── SimpleValue.java
│   │   ├── MMBuffer.java
│   │   ├── MMEncoder.java
│   │   ├── MMDecoder.java
│   │   └── MetaMessage.java
│   └── jsonc/
│       ├── JsoncScanner.java
│       ├── JsoncToken.java
│       ├── JsoncParser.java
│       ├── JsoncNode.java
│       ├── JsoncPrinter.java
│       ├── JsoncTag.java
│       └── ValueType.java
└── src/test/java/io/metamessage/
```

## 已完成功能

### Wire 格式
- ✅ Bool 編碼/解碼
- ✅ Int 編碼/解碼
- ✅ Float/Double 編碼/解碼
- ✅ String 編碼/解碼
- ✅ Bytes 編碼/解碼
- ✅ Array 編碼/解碼
- ✅ Object/Struct 編碼/解碼

### JSONC
- ✅ Scanner - 詞法分析，支持 `//` 和 `/* */` 註釋
- ✅ Parser - 語法分析，創建 AST
- ✅ Printer - 格式化輸出
- ✅ Tag 解析 - 支持 `// mm:` 前綴
- ✅ 類型感知打印 - String/UUID 等加引號，Int/Float 不加

## 已修復問題
- `decodeTagged` - 添加 `offset = innerEnd`
- `decodeObject` - 邊界檢查
- `decodeFloat` - 條件和 mantissa 修正

## 測試覆蓋
- ✅ 22 Wire 測試
- ✅ 15 JSONC 基本測試
- ✅ 22 JSONC MM 測試

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
