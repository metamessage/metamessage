# Swift MetaMessage (mm-swift) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 Swift，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成
- ✅ JSONC 基本對齊完成
- ✅ JSONC MM 對齊完成

## 項目結構
```
mm-swift/
├── Sources/
│   └── MetaMessage/
│       ├── MM/
│       │   ├── Constants.swift
│       │   ├── Prefix.swift
│       │   ├── SimpleValue.swift
│       │   ├── MMBuffer.swift
│       │   ├── MMEncoder.swift
│       │   ├── MMDecoder.swift
│       │   └── MetaMessage.swift
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

## 待完成功能
- 單元測試覆蓋

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
