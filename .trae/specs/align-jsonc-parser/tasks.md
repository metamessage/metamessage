# Tasks

- [x] T1: 深入分析各語言 JSONC 實現完整性
  - 逐一檢查 9 種語言的 Token、Scanner、Parser 的完整實現
  - 記錄每種語言與 Go 的差異清單
  - 依賴: 無

- [x] T2: Python JSONC 對齊 Go 實現
  - Token 增加 Column 字段
  - Scanner: _identifier 移除 lowercase、_number 返回字串、_line_comment 行為確認
  - Parser: 增加 CamelToSnake、tag.Name path、完整 type validation、depth 限制、tag2 merge 啟用
  - Printer: 確認 tag 輸出格式與順序

- [x] T3: PHP JSONC 對齊 Go 實現
  - Token 增加 Column 字段
  - Scanner 行為確認與調整
  - Parser 增加完整 type validation、CamelToSnake、depth 限制等

- [x] T4: TypeScript JSONC 對齊 Go 實現
  - parseObject 增加 CamelToSnake
  - 移除 hasChildFields 條件繼承邏輯（與 Go 保持一致）
  - 移除額外 revalidateValue 邏輯
  - Scanner 確認 scanString 行為

- [x] T5: Rust JSONC 對齊 Go 實現
  - Scanner: scan_identifier 移除 lowercase
  - Parser: 移除 inline comment 處理、增加完整 type validation、CamelToSnake
  - 確認 null 處理行為

- [x] T6: C JSONC 對齊 Go 實現
  - 全面檢查 Token/Scanner/Parser 差異
  - 增加 CamelToSnake、type validation、depth 限制等

- [x] T7: C++ JSONC 對齊 Go 實現
  - 全面檢查 Token/Scanner/Parser 差異
  - 增加 CamelToSnake、type validation、depth 限制等

- [x] T8: C# JSONC 對齊 Go 實現
  - 全面檢查 Token/Scanner/Parser 差異
  - 增加 CamelToSnake、type validation、depth 限制等

- [x] T9: Kotlin JSONC 對齊 Go 實現
  - 全面檢查 Token/Scanner/Parser 差異
  - 增加 CamelToSnake、type validation、depth 限制等

- [x] T10: Swift JSONC 對齊 Go 實現
  - 全面檢查 Token/Scanner/Parser 差異
  - 增加 CamelToSnake、type validation、depth 限制等

- [x] T11: 跨語言測試驗證
  - 運行 encode_decode_test.sh 確認所有語言 decode round-trip 全部通過
  - decode 還原一致性：34/34 通過 ✅
  - encode 字節一致性：C# WireEncoder 有獨立問題（非 JSONC 範疇）

# Task Dependencies

- T1 無依賴（探索性任務）
- T2–T10 可並行執行（各語言獨立修改）
- T11 依賴 T2–T10 完成