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

## 已修復的問題

### Kotlin
- `floatLenExtraBytes` - 修正 when 邏輯為 if-else

### Java
- `decodeTagged` - 添加 `offset = innerEnd`
- `decodeObject` - 添加邊界檢查 `i < keys.items().size()`
- `decodeFloat` - 修正條件判斷和 mantissa 計算

### TypeScript/JavaScript/Python/C#/Rust/Swift/PHP
- 對齊 Golang 基準實現

---

## 驗證標準

| 語言 | Wire 測試 | JSONC 基本測試 | JSONC MM 測試 |
|------|-----------|----------------|---------------|
| Golang | - | 基準 | 基準 |
| Java | 22 | 15 | 22 |
| Kotlin | 5 | 23 | 10 |
| TypeScript | 35 | 35 | 49 |
| Python | 68 | 68 | 68 |
| JavaScript | - | - | 55 |
| C# | - | - | - |
| Rust | 4 | 4 | 3 |
| Swift | - | - | - |
| PHP | 9 | 9 | 9 |

**所有任務已完成 ✅**
