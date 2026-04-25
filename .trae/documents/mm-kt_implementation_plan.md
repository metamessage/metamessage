# Kotlin MetaMessage (mm-kt) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 Kotlin，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成 (5 tests)
- ✅ JSONC 基本對齊完成 (23 tests)
- ✅ JSONC MM 對齊完成 (10 tests)

## 項目結構
```
mm-kt/
├── src/main/kotlin/io/metamessage/
│   ├── mm/
│   │   ├── Constants.kt
│   │   ├── Types.kt
│   │   ├── Prefix.kt
│   │   ├── SimpleValue.kt
│   │   ├── MMBuffer.kt
│   │   ├── MMEncoder.kt
│   │   ├── MMDecoder.kt
│   │   └── MetaMessage.kt
│   └── jsonc/
│       ├── JsoncScanner.kt
│       ├── JsoncToken.kt
│       ├── JsoncParser.kt
│       ├── JsoncNode.kt
│       ├── JsoncPrinter.kt
│       ├── JsoncTag.kt
│       └── ValueType.kt
└── src/test/kotlin/io/metamessage/
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
- `floatLenExtraBytes` - 將 when 改為 if-else

## 測試覆蓋
- ✅ 5 Wire 測試
- ✅ 23 JSONC 基本測試
- ✅ 10 JSONC MM 測試

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
