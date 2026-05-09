# MetaMessage 多語言實現對齊 - 任務計劃

## 當前狀態：全部完成 ✅

---

## 1. Wire 格式對齊

| 語言 | 狀態 | 測試 |
|------|------|------|
| Java | ✅ 完成 | 22 tests |
| Kotlin | ✅ 完成 | 5 tests |
| TypeScript | ✅ 完成 | 35 tests |
| Python | ✅ 完成 | 68 tests |
| JavaScript | ✅ 完成 | - |
| C# | ✅ 完成 | - |
| Rust | ✅ 完成 | 4 tests |
| Swift | ✅ 完成 | - |
| PHP | ✅ 完成 | 9 tests |

---

## 2. JSONC 基本對齊

| 語言 | 狀態 | 測試 |
|------|------|------|
| Golang | ✅ 基準 | - |
| TypeScript | ✅ 完成 | 35 tests |
| Python | ✅ 完成 | 68 tests |
| JavaScript | ✅ 完成 | - |
| Java | ✅ 完成 | 15 tests |
| Kotlin | ✅ 完成 | 23 tests |
| C# | ✅ 完成 | - |
| Rust | ✅ 完成 | 4 tests |
| Swift | ✅ 完成 | - |
| PHP | ✅ 完成 | 9 tests |

---

## 3. JSONC MM 對齊

| 語言 | Tag 解析 | 類型感知打印 | 測試 |
|------|----------|--------------|------|
| Golang | ✅ | ✅ 基準 | - |
| Kotlin | ✅ | ✅ | 10 tests |
| Java | ✅ | ✅ | 22 tests |
| TypeScript | ✅ | ✅ | 49 tests |
| Python | ✅ | ✅ | 68 tests |
| JavaScript | ✅ | ✅ | 55 tests |
| C# | ✅ | ✅ | - |
| Rust | ✅ | ✅ | 3 tests |
| Swift | ✅ | ✅ | - |
| PHP | ✅ | ✅ | 9 tests |

---

## 4. MM Validator 對齊 ✅

### JavaScript Validator 增強
- [x] 添加 decimal 類型驗證
- [x] 添加 IP 類型驗證 (IPv4/IPv6)
- [x] 添加 URL 類型驗證
- [x] 添加 slice 類型驗證
- [x] 添加 date/time 單獨驗證方法
- [x] 更新 types.js 添加 options 參數支持
- [x] 更新 types.ts 添加完整方法定義

### Kotlin Validator 測試
- [x] 創建 ValidatorTest.kt 測試文件
- [x] 整數類型測試 (int8/16/32/64, uint/uint8/16/32/64)
- [x] 浮點類型測試 (float32/float64)
- [x] 字符串和特殊類型測試 (string, uuid, email, enum, datetime)
- [x] 數組和結構體測試 (array, struct, image)
- [x] 邊界條件測試
- [x] 錯誤處理測試

---

## 已修復的問題

### Kotlin
- `floatLenExtraBytes` - 修正 when 邏輯為 if-else

### Java
- `decodeTagged` - 添加 `offset = innerEnd`
- `decodeObject` - 添加邊界檢查 `i < keys.items().size()`
- `decodeFloat` - 修正條件判斷和 mantissa 計算

### TypeScript/JavaScript/Python/C#/Rust/Swift/PHP
- 對齊 Golang 基準實現

### JavaScript types.js
- 移除不必要的 `null` 和 `nil` 方法
- 為所有方法添加 `options` 參數支持

### TypeScript types.ts
- 修正 `int` 類型定義 (number | bigint → number)
- 為所有方法添加 `options` 參數支持
- 添加缺失的方法定義 (uuid, datetime, date, time, email, url, ip, decimal, i8/16/32/64, u/8/16/32/64, f32/64)
- 修正 `i64` 和 `u64` 使用 `bigint` 類型
- 移除不必要的 `null` 方法

---

## 驗證標準

| 語言 | Wire 測試 | JSONC 基本測試 | JSONC MM 測試 | Validator 測試 |
|------|-----------|----------------|---------------|----------------|
| Golang | - | 基準 | 基準 | - |
| Java | 22 | 15 | 22 | - |
| Kotlin | 5 | 23 | 10 | ✅ 完整 |
| TypeScript | 35 | 35 | 49 | - |
| Python | 68 | 68 | 68 | - |
| JavaScript | - | - | 55 | ✅ 增強 |
| C# | - | - | - | - |
| Rust | 4 | 4 | 3 | - |
| Swift | - | - | - | - |
| PHP | 9 | 9 | 9 | - |

**所有任務已完成 ✅**
