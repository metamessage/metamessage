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

## MM Validator 對齊檢查點 ✅

### JavaScript Validator
- [x] 添加 decimal 類型驗證
- [x] 添加 IP 類型驗證 (IPv4/IPv6)
- [x] 添加 URL 類型驗證
- [x] 添加 slice 類型驗證
- [x] 添加 date/time 單獨驗證方法
- [x] 更新 validate() 方法分發新類型
- [x] 語法檢查通過

### JavaScript types.js
- [x] 移除 `null` 和 `nil` 方法
- [x] 為所有方法添加 `options` 參數支持

### TypeScript types.ts
- [x] 修正 `int` 類型定義 (number | bigint → number)
- [x] 為所有方法添加 `options` 參數支持
- [x] 添加缺失的方法定義 (uuid, datetime, date, time, email, url, ip, decimal, i8/16/32/64, u/8/16/32/64, f32/64)
- [x] 修正 `i64` 和 `u64` 使用 `bigint` 類型
- [x] 移除不必要的 `null` 方法
- [x] TypeScript 編譯通過

### Kotlin Validator 測試
- [x] ValidatorTest.kt 文件創建完成
- [x] int8 驗證測試
- [x] int16 驗證測試
- [x] int32 驗證測試
- [x] int64 驗證測試
- [x] uint/uint8/uint16/uint32/uint64 驗證測試
- [x] float32/float64 驗證測試
- [x] string 驗證測試（含 pattern 約束）
- [x] bytes 驗證測試（含 base64）
- [x] bool 驗證測試
- [x] uuid 驗證測試（含版本檢查）
- [x] email 驗證測試
- [x] enum 驗證測試
- [x] datetime/date/time 驗證測試
- [x] array 驗證測試（含 childUnique）
- [x] struct 驗證測試
- [x] image 驗證測試
- [x] 空值測試（allowEmpty=true/false）
- [x] 最小值邊界測試
- [x] 最大值邊界測試
- [x] 越界值測試
- [x] 格式錯誤測試
- [x] 類型不匹配測試
- [x] 參數解析錯誤測試
- [x] 測試文件組織清晰
- [x] 測試方法命名規範

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
- [x] Validator 完整測試覆蓋

### TypeScript
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (35 tests Wire + 49 tests MM)
- [x] types.ts 方法定義完善

### Python
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (68 tests)

### JavaScript
- [x] JSONC tag 解析
- [x] JSONC 類型感知打印
- [x] 測試通過 (55 tests)
- [x] Validator 增強完成

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
- [x] JavaScript Validator 與 Golang 功能對齊
- [x] Kotlin Validator 測試覆蓋完整
- [x] TypeScript types.ts 方法定義完善
- [x] 跨語言兼容性已驗證

**所有檢查點已完成 ✅**
