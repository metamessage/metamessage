# Rust MetaMessage (mm-rs) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 Rust，測試要完整全面。

## 當前進度
- ✅ Wire 格式對齊完成 (4 tests)
- ✅ JSONC 基本對齊完成 (4 tests)
- ✅ JSONC MM 對齊完成 (3 tests)

## 項目結構
```
mm-rs/
├── src/
│   ├── jsonc/
│   │   ├── mod.rs
│   │   ├── scanner.rs
│   │   ├── parser.rs
│   │   ├── printer.rs
│   │   ├── ast.rs
│   │   ├── tag.rs
│   │   ├── to_string.rs
│   │   ├── value_type.rs
│   │   └── from_str.rs
│   └── mm/
│       ├── mod.rs
│       ├── constants.rs
│       ├── encoder.rs
│       ├── decoder.rs
│       └── types.rs
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
- ✅ 4 Wire 測試
- ✅ 4 JSONC 基本測試
- ✅ 3 JSONC MM 測試

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`
