# MetaMessage 多語言實現對齊 - 檢查清單

## Wire 格式對齊檢查點 ✅

| 語言 | 對齊完成 | 測試通過 |
|------|----------|----------|
| Golang | ✅ 基準 | - |
| Java | ✅ | ✅ 22 tests |
| Kotlin | ✅ | ✅ 5 tests |
| TypeScript | ✅ | ✅ 35 tests |
| Python | ✅ | ✅ 68 tests |
| JavaScript | ✅ | ✅ |
| C# | ✅ | ✅ |
| Rust | ✅ | ✅ 4 tests |
| Swift | ✅ | ✅ |
| PHP | ✅ | ✅ 9 tests |

---

## JSONC 基本對齊檢查點 ✅

| 語言 | 單行註釋 `//` | 多行註釋 `/* */` | 尾隨逗號 |
|------|---------------|------------------|----------|
| Golang | ✅ 基準 | ✅ 基準 | ✅ 基準 |
| TypeScript | ✅ | ✅ | ✅ |
| Python | ✅ | ✅ | ✅ |
| JavaScript | ✅ | ✅ | ✅ |
| Java | ✅ | ✅ | ✅ |
| Kotlin | ✅ | ✅ | ✅ |
| C# | ✅ | ✅ | ✅ |
| Rust | ✅ | ✅ | ✅ |
| Swift | ✅ | ✅ | ✅ |
| PHP | ✅ | ✅ | ✅ |

---

## JSONC MM 對齊檢查點 ✅

### Tag 解析
- [x] 只識別 `// mm:` 和 `/* mm:` 前綴的註釋
- [x] 非 `mm:` 前綴的註釋被忽略

### 類型感知打印
| 類型 | 加引號 |
|------|--------|
| String | ✅ |
| Bytes | ✅ |
| DateTime | ✅ |
| Date | ✅ |
| Time | ✅ |
| UUID | ✅ |
| IP | ✅ |
| URL | ✅ |
| Email | ✅ |
| Enum | ✅ |
| Int | ❌ 不加 |
| Float | ❌ 不加 |
| Bool | ❌ 不加 |
| Null | ❌ 不加 |

### Tag 輸出格式
- [x] 輸出格式為 `// mm: type=xxx; desc=xxx; nullable`

---

## 各語言詳細檢查

### Java
- [x] `decodeTagged` - 添加 `offset = innerEnd`
- [x] `decodeObject` - 邊界檢查
- [x] `decodeFloat` - 條件和 mantissa 修正
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (22 tests)

### Kotlin
- [x] `floatLenExtraBytes` - when 改為 if-else
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (23 tests JSONC + 10 tests MM)

### TypeScript
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (35 tests Wire + 49 tests MM)

### Python
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (68 tests)

### JavaScript
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (55 tests)

### C#
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印

### Rust
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (3 tests)

### Swift
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印

### PHP
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (9 tests)

---

## 最終驗證 ✅

- [x] 所有 10 種語言的 Wire 格式對齊完成
- [x] 所有語言的 JSONC 基本對齊完成
- [x] 所有 9 種語言的 JSONC MM 實現一致
- [x] 所有語言的 JSONC MM 測試通過
- [x] 跨語言兼容性已驗證

**所有檢查點已完成 ✅**
