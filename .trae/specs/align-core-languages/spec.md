# 各語言 Core 層實現對齊 Go 規範

## Why

各語言 Core 層（編碼器/解碼器/API）存在與 Go 參考實作不一致之處，包括 API 命名不統一、編碼邏輯組織方式差異、類型處理方式不同等，可能導致跨語言編碼/解碼行為不一致。

## What Changes

### 1. Top-Level API 命名統一

**Go** (`mm.go`): `FromJSONC`, `Decode`, `FromValue`, `ValueToJSONC`, `ParseFromJSONC`

| 語言 | 現有 API 命名 | 需對齊為 |
|------|-------------|---------|
| **C#** (`MetaMessage.cs`) | `EncodeFromJsonc`, `DecodeToJsonc`, `FromValue`, `ValueToJsonc`, `ParseFromJSONC` | 統一命名風格，移除 `EncodeFromJsonc`/`DecodeToJsonc` 這類非標準命名（但保留向後相容別名） |
| **Swift** (`MetaMessage.swift`) | `fromJSONC(input)`, `toJSONC(data)`（全域函數） | 統一命名風格，與 Go 對齊 |
| **Rust** (`lib.rs`) | `encode_from_jsonc`, `decode_to_jsonc`, `parse_jsonc`, `to_jsonc_string` | 命名風格已對齊（Rust snake_case 慣例） |

### 2. C# MetaMessage.cs 編碼邏輯拆分

**Go** 將編碼邏輯拆分為獨立檔案（`encode_array.go`, `encode_bytes.go`, `encode_float.go` 等），**C# `MetaMessage.cs`** 將所有編碼樹邏輯集中在單一檔案（~711 行）。

建議將 `MetaMessage.cs` 中的樹狀編碼邏輯（`EncodeTreeValue`, `EncodeScalarTree`, `EncodeArrayTree`, `EncodeMapTree` 等）拆分為獨立檔案或移至 `WireEncoder.cs`，使 `MetaMessage.cs` 僅保留頂層 API。

### 3. C# 額外檔案清理

C# Core 層包含 Go 中不存在的額外檔案：
- `CamelToSnake.cs` — 應移至共用工具
- `TypeInference.cs` — 若僅在 Core 內部使用，可保留但不影響行為
- `TimeUtil.cs` — 應移至共用工具
- `TagFieldParser.cs` — 應移至 Ir 層
- `MimeWire.cs` — 應移至 Ir 層

**注意**：以上為重構建議，不影響編碼/解碼行為。不作強制要求。

### 4. Swift Core 層缺少公開 Top-Level API

**Go** 提供 `mm.go` 作為核心公開函數：
- `FromJSONC(s string) (bs []byte, err error)`
- `Decode(data []byte) (ir.Node, error)`
- `FromValue(v any, tag string) (bs []byte, err error)`
- `ValueToJSONC(value any, name string) (string, error)`
- `ParseFromJSONC(in string) (out ir.Node, err error)`

**Swift** 目前 `MetaMessage.swift` 的公開 API 僅為 `class MetaMessageEncoder`/`MetaMessageDecoder`（內部類別），以及全域函數 `fromJSONC()`/`toJSONC()`（定義在 MetaMessage.swift 檔案層級）。缺少公開的 `MetaMessage` API 封裝。

### 5. Swift Wire 常數位置

Swift 的 wire 格式常數（`MMPrefix`, `MMSimpleValue`, `MMConstants` 等）定義在 `Ir/Constants.swift`，而非 `Core/` 目錄下。Go 的 `constants.go` 在 `core` 套件中。

**受影響語言**：
- **Swift**: 將 Core 相關常數移至 Core 目錄（或至少在 Core 中重新匯出）

### 6. Tag Encoding/Decoding 邏輯位置

**Go**: `core/encode_tag.go`, `core/decode_tag.go` 在 core 層
**C#**: Tag encoding 在 `WireEncoder.cs` 的 `EncodeTagInner` 方法中；Tag decoding 在 `WireDecoder.cs` 中
**Swift**: Tag encoding 在 `Core/Encoder.swift` 的 `encodeTagToBytes` 中
**Rust**: Tag encoding 在 `core/encoder.rs` 中

各語言都已包含對應的 tag 編碼邏輯，但需確保行為一致。Tag 編碼行為一致性透過 `compare_with_go.sh` 測試驗證。

### 7. ValueToNode / Type Inference 一致性

**Go** (`value_to_node.go`): 將原生類型（bool/int/float/string/map/slice）轉換為 `ir.Node`
**C#**: `ReflectMmEncoder.ValueToNode` + `TypeInference.cs`
**Swift**: `ValueToNode.swift`
**Rust**: `value_to_node.rs`

確保類型推斷規則一致：
- `int` → `ValueType.I`
- `float64` → `ValueType.F64`
- `string` → `ValueType.Str`
- `bool` → `ValueType.Bool`
- `[]byte` → `ValueType.Bytes`
- `map` → `ValueType.Map`
- `slice` → `ValueType.Arr`
- `nil` → `ValueType.Unknown` (is_null)

### 8. Encoder 內部方法對齊

Go 的 `Encoder` 通過 `Encoder` 結構體提供編碼方法，內部方法委派給專用編碼函數：
- `encodeInt` → `encodeIntWriter`
- `encodeString` → `encodeStringWriter`
- `encodeBytes` → `encodeBytesWriter`
- `encodeFloat` → `encodeFloatWriter`
- `encodeObject` → `encodeObjectWriter`
- `encodeArray` → `encodeArrayWriter`

各語言已實現對應功能，但確保在邊界情況（如負數、浮點數 +/-0、NaN、Infinity、大整數）行為一致即可。

## Core 層行為對齊關注點

相比 IR 層（struct 欄位命名、enum 值順序等），Core 層的對齊重點在於：

1. **Wire 格式常量值一致**：Prefix, SimpleValue, IntLen, FloatLen 等位元運算結果確保一致
2. **編碼/解碼行為一致**：對相同 Node 輸入產生相同 wire bytes 輸出
3. **Top-Level API 命名一致**：便於跨語言開發者理解與使用
4. **公開 API 簽名一致**：main harness 使用的 API 各語言都能呼叫

## Impact

### ADDED Requirements

#### Requirement: Swift 公開 MetaMessage API

##### Scenario: Swift 提供 FromJSONC/Decode 公開函數
- **WHEN** 外部程式呼叫 `MM.FromJSONC(input)` 
- **THEN** 返回的 wire bytes 與 Go `core.FromJSONC()` 完全一致

##### Scenario: Swift 提供 FromValue/ValueToJSONC
- **WHEN** 外部程式呼叫 `MM.FromValue(value, tag)`
- **THEN** 返回的 wire bytes 與 Go `core.FromValue()` 完全一致

### MODIFIED Requirements

#### Requirement: C# Top-Level API 命名別名

##### Scenario: C# 保留舊 API 同時新增標準命名
- **WHEN** 使用 `MetaMessage.EncodeFromJsonc`
- **THEN** 仍可正常運作（向後相容）
- **AND** 新增 `MetaMessage.FromJSONC` 作為標準別名

#### Requirement: Swift Wire 常數移至 Core

##### Scenario: Swift Core 目錄包含 wire 常數
- **WHEN** `import MetaMessage`
- **THEN** 可透過 Core 模組存取 `MMPrefix`、`MMSimpleValue`、`MMConstants`

### REMOVED Requirements

無

## 檢查清單

- [ ] C# `MetaMessage.cs` 新增 `FromJSONC` 別名
- [ ] Swift 公開 `MetaMessage` API（`fromJSONC`, `toJSONC`, `fromValue`）
- [ ] Swift Wire 常數（`MMPrefix`, `MMSimpleValue`, `MMConstants`）可在 Core 中存取
- [ ] 所有語言的 Top-Level API 命名風格對齊 Go
- [ ] `compare_with_go.sh` cs/sw/rs 測試全部 PASS（encode + decode）