# TypeScript MetaMessage (mm-ts) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 TypeScript，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成 (35 tests)
- ✅ JSONC 基本對齊完成 (35 tests)
- ✅ JSONC MM 對齊完成 (49 tests)

## 項目結構
```
mm-ts/
├── src/
│   ├── jsonc/
│   │   ├── scanner.ts
│   │   ├── parser.ts
│   │   ├── printer.ts
│   │   ├── ast.ts
│   │   ├── tag.ts
│   │   ├── binder.ts
│   │   └── index.ts
│   └── mm/
│       ├── encoder.ts
│       ├── decoder.ts
│       └── index.ts
└── tests/
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

## 測試覆蓋
- ✅ 35 Wire 測試
- ✅ 35 JSONC 基本測試
- ✅ 49 JSONC MM 測試

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
