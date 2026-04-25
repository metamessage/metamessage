# C# MetaMessage (mm-cs) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 C#，排除 CLI 功能，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成
- ✅ JSONC 基本對齊完成
- ✅ JSONC MM 對齊完成

## 項目結構
```
mm-cs/
├── src/
│   └── MetaMessage/
│       ├── Constants.cs
│       ├── Types.cs
│       ├── Prefix.cs
│       ├── SimpleValue.cs
│       ├── MMBuffer.cs
│       ├── MMEncoder.cs
│       ├── MMDecoder.cs
│       ├── MetaMessage.cs
│       └── Jsonc/
│           ├── JsoncScanner.cs
│           ├── JsoncToken.cs
│           ├── JsoncParser.cs
│           ├── JsoncNode.cs
│           ├── JsoncPrinter.cs
│           ├── JsoncTag.cs
│           └── ValueType.cs
└── tests/
```

## 已完成功能

### Wire 格式
- ✅ Bool 編碼/解碼
- ✅ Int (8/16/32/64) 編碼/解碼
- ✅ UInt (8/16/32/64) 編碼/解碼
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

## 待完成功能
- 單元測試覆蓋

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
