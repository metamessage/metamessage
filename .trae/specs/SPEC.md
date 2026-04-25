# MetaMessage 多語言實現對齊 - 規格文檔

## Why

Golang 是基準實現，需要確保 10 種語言（go, java, ts, kt, py, js, cs, rs, swift, php）的 Wire 格式和 JSONC 實現完全一致，確保跨語言兼容性。

## 當前狀態

### 已完成的工作

| 模組 | 狀態 | 說明 |
|------|------|------|
| Wire 格式對齊 | ✅ 完成 | 10 種語言全部對齊 |
| JSONC 基本對齊 | ✅ 完成 | 所有語言支持 `//` 和 `/* */` 註釋 |
| Kotlin JSONC 實現 | ✅ 完成 | 完整實現 + 23 tests |
| PHP JSONC 實現 | ✅ 完成 | 完整實現 + 9 tests |
| JSONC MM 對齊 | ✅ 完成 | 所有語言支持 `mm:` tag 和類型感知打印 |

---

## 已完成的對齊工作

### 1. Wire 格式對齊

| 語言 | 狀態 | 測試 |
|------|------|------|
| Java | ✅ | 22 tests |
| Kotlin | ✅ | 5 tests |
| TypeScript | ✅ | 35 tests |
| Python | ✅ | 68 tests |
| JavaScript | ✅ | - |
| C# | ✅ | - |
| Rust | ✅ | 4 tests |
| Swift | ✅ | - |
| PHP | ✅ | 9 tests |
| Golang | ✅ 基準 | - |

### 2. JSONC 基本對齊

| 語言 | 狀態 | 測試 |
|------|------|------|
| Golang | ✅ 基準 | - |
| TypeScript | ✅ | 35 tests |
| Python | ✅ | 68 tests |
| JavaScript | ✅ | - |
| Java | ✅ | 15 tests |
| Kotlin | ✅ | 23 tests |
| C# | ✅ | - |
| Rust | ✅ | 4 tests |
| Swift | ✅ | - |
| PHP | ✅ | 9 tests |

### 3. JSONC MM 對齊

MM 特定 JSONC 格式：註釋中的 `mm:` tag 和類型感知打印。

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

## MM JSONC 特性

### 1. Tag 解析
- 只識別 `// mm:` 和 `/* mm:` 前綴的註釋
- 非 `mm:` 前綴的註釋被忽略

### 2. 類型感知打印
| 類型 | 加引號 |
|------|--------|
| String, Bytes, DateTime, Date, Time, UUID, IP, URL, Email, Enum | ✅ |
| Int, Float, Bool, Null | ❌ |

### 3. Tag 輸出格式
```
// mm: type=xxx; desc=xxx; nullable
```

---

## 受影響的代碼路徑

| 語言 | 路徑 |
|------|------|
| Java | `mm-java/src/main/java/io/metamessage/` |
| Kotlin | `mm-kt/src/main/kotlin/io/metamessage/` |
| TypeScript | `mm-ts/src/jsonc/` |
| Python | `mm-py/` |
| JavaScript | `mm-js/src/jsonc/` |
| C# | `mm-cs/src/MetaMessage/` |
| Rust | `mm-rs/src/jsonc/` |
| Swift | `mm-swift/Sources/MetaMessage/JSONC/` |
| PHP | `mm-php/src/io/metamessage/` |
| Golang | `internal/mm/`, `internal/jsonc/` |

---

## ADDED Requirements

### Requirement: 多語言 Wire 格式一致性
所有 10 種語言的 Wire 格式編碼/解碼實現應該完全一致。

### Requirement: 多語言 JSONC 一致性
所有 9 種語言的 JSONC 解析和打印實現應該完全一致。

### Requirement: MM JSONC 類型感知打印
MM JSONC 應該根據 tag.type 決定是否為字符串加引號。

## MODIFIED Requirements
無

## REMOVED Requirements
無
