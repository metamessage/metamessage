# Checklist

## Token 定義
- [x] Python Token 增加 Column 字段
- [x] PHP Token 增加 Column 字段
- [x] Rust Token 確認 Column 字段存在
- [x] C/C++/C#/Kotlin/Swift Token Column 字段確認
- [x] 所有語言 Token Type 枚舉值與 Go 一致（13 種類型）

## Scanner 行為
- [x] 所有語言 scanString 不做轉義（raw escape preservation）
- [x] 所有語言 scanComment trim whitespace
- [x] 所有語言 scanLiteral 不 lowercase 關鍵字
- [x] 所有語言 scanLiteral 正確識別 true/false/null/Number

## Parser 邏輯
- [x] 所有語言 parseCommentsToTag 使用 prefix 檢查 "mm:"
- [x] 所有語言 consumeCommentsFor 檢查 anchorLine gap > 1
- [x] 所有語言 consumeCommentsFor 使用 MergeTag 合併
- [x] 所有語言 parseObject 使用 CamelToSnake 轉換 key
- [x] 所有語言 parseObject 支援 tag.Name path 構建
- [x] 所有語言 parseObject/parseArray 支援 tag.Inherit
- [x] 所有語言 parseObject 物件驗證（ValidateObj/ValidateMap）
- [x] 所有語言 parseArray 陣列驗證（ValidateArr/ValidateVec）
- [x] 所有語言完整 type validation（Str/Bytes/Datetime/Date/Time/Uuid/Decimal/Ip/Url/Email/Enum/Media/I/I8/I16/I32/I64/U/U8/U16/U32/U64/F32/F64/Bigint/Bool）
- [x] 所有語言 max depth = 32 限制
- [x] 所有語言 null literal 返回 "null is not supported" 錯誤

## JSONC Printer
- [x] 所有語言 Printer tag 輸出格式 "// mm: ..."
- [x] 所有語言 Printer tag 輸出順序與 Go 一致
- [x] 所有語言 Printer indent 使用 \t

## 特殊語言修復
- [x] Python 移除 _identifier 中的 lowercase
- [x] Python _number 返回 string literal
- [x] Python 啟用 _parse_array 中的 tag2 merge
- [x] TypeScript 移除 hasChildFields 條件繼承
- [x] TypeScript 移除 revalidateValue 邏輯
- [x] Rust 移除 inline comment 處理（scan_identifier lowercase）
- [x] Rust 移除 parse_node 中的 inline comment 邏輯

## 測試
- [ ] 每種語言增加 parser_test 對齊 Go 測試案例
- [x] encode_decode_test.sh decode round-trip 全部通過
- [ ] C# WireEncoder 編碼問題（非 JSONC 範疇，需單獨處理）