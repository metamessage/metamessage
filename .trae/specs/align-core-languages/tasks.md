# Tasks

- [ ] Task 1: C# 新增 FromJSONC 別名
  - 在 [MetaMessage.cs](file:///Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Core/MetaMessage.cs) 中新增 `FromJSONC` 和 `FromValue` 等方法別名（保留現有方法不刪除）
  - 確保 `public static byte[] FromJSONC(string jsonc)` 委派給現有實現

- [ ] Task 2: Swift 公開 Top-Level MetaMessage API
  - 在 [MetaMessage.swift](file:///Users/lizongying/IdeaProjects/meta-message/mm-swift/Sources/MetaMessage/Core/MetaMessage.swift) 中新增公開的靜態方法：
    - `public static func fromJSONC(_ input: String) throws -> Data`（對應 Go `FromJSONC`）
    - `public static func decode(_ data: Data) throws -> Node`（對應 Go `Decode`）
    - `public static func fromValue(_ value: Any, tag: String) throws -> Data`（對應 Go `FromValue`）
    - `public static func valueToJSONC(_ value: Any, name: String) throws -> String`（對應 Go `ValueToJSONC`）
    - `public static func parseJSONC(_ input: String) throws -> Node`（對應 Go `ParseFromJSONC`）
  - 確保 `fromJSONC` 和 `toJSONC`（現有全域函數）仍可正常運作

- [ ] Task 3: Swift Wire 常數可從 Core 存取
  - 確認 [Ir/Constants.swift](file:///Users/lizongying/IdeaProjects/meta-message/mm-swift/Sources/MetaMessage/Ir/Constants.swift) 中的 `MMPrefix`、`MMSimpleValue`、`MMConstants` 從 Core 模組可存取
  - 若需要，在 Core 中新增重新匯出（typealias 或 extension）

- [ ] Task 4: 驗證跨語言編碼一致性
  - 執行 `encode_decode_test_cs.sh` 確認 C# vs Go 無 DIFF
  - 執行 `encode_decode_test_sw.sh` 確認 Swift vs Go 無 DIFF
  - 執行 `encode_decode_test_rs.sh` 確認 Rust vs Go 無 DIFF

# Task Dependencies
- [Task 1], [Task 2], [Task 3] 可並行執行
- [Task 4] 依賴所有前置任務完成