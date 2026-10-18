# MetaMessage v0.4.0 Footer 線格式規格

版本：v0.4.0（硬切，不兼容 v0.3.x）  
實作基準：Go 參考實現 `internal/core`  
文件狀態：規格已凍結，供後續 9 個語言移植對齊。

---

## 1. 總體布局

### 1.1 帶 meta（tag）的節點

```
[PrefixTag | 檔位][total_len 大端][body 正向編碼][reverse(meta 正序字節)][尾部長度槽]
```

| 區段 | 寬度 | 說明 |
|------|------|------|
| `PrefixTag \| 檔位` | 1 字節 | 高位 3 位 = `0b111`；低位 5 位 = 長度檔位標記（同 v0.3.x）|
| `total_len` | 0/1/2/3/4 字節 | body + reverse(meta) + 尾槽 的總字節數，大端整數 |
| `body` | 可變 | 節點內部結構，正向編碼（嵌套子節點各自獨立封裝）|
| `reverse(meta)` | meta_len 字節 | tag 序列化後字節流的**整體字節逆序**（見 §1.2）|
| `尾部長度槽` | 1/2/3 字節 | 自描述 meta_len（見 §1.3）|

`total_len = len(body) + len(meta) + len(尾部槽)`

**長度檔位對照（沿用 TagLen 常量）**

| 頭字節低位 5 位 | 含義 |
|----------------|------|
| 0..27 | `total_len` 內聯在頭字節，寬度 0 |
| 28 | `total_len` 用 1 字節大端表示 |
| 29 | `total_len` 用 2 字節大端表示 |
| 30 | `total_len` 用 3 字節大端表示 |
| 31 | `total_len` 用 4 字節大端表示 |

### 1.2 meta 逆序編碼的精確定義

meta 正序字節流為 `Tag.Bytes()` 的輸出（字段邏輯順序不變：Example → IsNull → Nullable → Desc → Type → … → KMore）。

存儲區 `S = reverse(concat(m1..mk))`，即將整條 meta 字節流做**字節級逆序**後緊跟在 body 之後寫入。

解碼器等價操作：從 `S` 的末尾字節向前逐字節讀，得到的字節流恰好等於 `Tag.Bytes()` 的正序輸出。

Go 實作采用「切片翻轉 + 複用正向 `decodeTagBytes` 循環」路線；其他語言可二選一（翻轉切片或反向光標讀取），但產物字節必須與此規格逐字節一致。

### 1.3 尾部長度槽（Footer Slot）

放在節點最末端，解碼器從最後一字節開始讀：

| 末字節值 `v` | 含義 | 槽寬 |
|-------------|------|------|
| `v < 254` | meta_len = `v` | 1 |
| `v == 254` | 倒數第 2 字節（u8 大端）= meta_len | 2 |
| `v == 255` | 倒數第 3..2 字節（u16 大端）= meta_len；上限 Max2Byte = 65535 | 3 |

槽采用 canonical 最小寬度編碼，不允许多餘寬度。meta_len > 65535 時編碼器報 "tag too long"。

原 v0.3.x 頭部 `tag_len` 編碼（`encodeT`）廢除，由尾部槽完全取代。

### 1.4 無 meta 的節點

無 tag（meta 為空）的節點保持裸編碼，不添加 PrefixTag 封裝，無 footer：

```
[裸 body]
```

`Tag.Bytes()` 字段生成順序與 v0.3.x 相同；`encodeFooter` 對 meta 長度為零時返回 0（不封裝）。

---

## 2. 節點邊界穩定解碼保證

### 2.1 雙向定界流程

```
節點起點 start
  → 讀頭字節得 prefix，檔位得 l1（長度字段寬度）
  → 讀 total_len（若 l1>0）
  → node_end = start + 1 + l1 + total_len
  → 從 node_end 反向讀尾部槽得 meta_len, slot_w
  → footer_start = node_end - slot_w - meta_len
  → body_start  = start + 1 + l1
  → body_len    = footer_start - body_start
  → 正向遞歸解碼 body，斷言實際消耗字節 == body_len
  → 解碼 reverse(meta) 得 tag 字段
  → d.offset 設為 node_end，返回
```

「斷言實際消耗字節 == body_len」是 footer 模式的核心不變量。若正向解析消耗不等於 body_len，報 **"footer boundary mismatch"**。

### 2.2 越界防護

| 條件 | 處理 |
|------|------|
| `footer_start < body_start` | 報 "footer overflow" |
| `total_len > Max4Byte`（4 字節檔位超限）| 報 "tag+payload too long" |
| 尾部槽檔位越界（字節不足）| 報 EOF / "truncated" |
| meta_len = 0（空 meta）| 報 "empty meta"（不合法封裝）|
| body_len = 0 且非 is_null 合成 | 報 "empty body in tagged node" |

---

## 3. 邊界場景規則

### 3.1 數組（Array）

容器 body = 元素順序正向拼接（不變）。數組自身的 tag 進 footer。子元素帶 tag 時各自持有獨立 PrefixTag+footer 封裝，遞歸解碼天然處理。

`decodeArr` 循環上界由 `l2`（容器 payload 長度）控制，末尾加精確消耗校驗：`index != l2` → "invalid arr"。

### 3.2 對象（Object）

body = `[鍵數組][值1][值2]…`，鍵數組屬於結構信息，**保持在 body 頭部不遷移**（顯式決策）。

對象自身的 tag（含 `child_*` 繼承源）進 footer。`decodeObj` 循環上界 = `body_len`，鍵數組+各值的消耗總長必須等於 body_len；同時增加 key 數量邊界檢查（防止過長 total_len 觸發越界讀取）。

### 3.3 空節點

- 空數組 / 空對象：body_len = 0 合法，footer 仍正常解析 tag。
- `NodeNull` 帶 tag：body 為 `SimpleNull` 單字節 + footer。
- `NodeNull` 不帶 tag：裸 `SimpleNull` 單字節，無封裝。

### 3.4 is_null 合成表

is_null = true 時，根據 `tag.Type` 直接合成零值節點，不讀 body。

合成規則與 v0.3.x 解碼器相同；v0.4.0 新增：
- `ValueTypeEnum`：合成 `Data: -1`
- `ValueTypeMedia`：合成 `Data: []byte{}`

部分類型 is_null 寫 `SimpleNullX` 單字節 body（body_len=1），部分寫空 body（body_len=0）。兩種都支持。

### 3.5 引用信息（ref_id / ref_to）

KMore meta 塊（位圖 + 值）隨 meta 整體進 footer，無特殊規則。

順序性保證：任一節點的 meta 在自身 body 之前完成解析（因為 footer 先於 body 被切出），被引用節點（流中在前）的 ref_id 註冊必然早於後流引用節點的 ref_to 解析。與 v0.3.x header 模式的左到右性質一致。

### 3.6 nullable 指針

`EncodeFromValue(&struct, "")` 對指針類型自動注入 `nullable` tag（行為同 v0.3.x）。footer 模式下 nullable 字節進入 meta 逆序區。

---

## 4. 長度上限

| 字段 | 上限 | 超出處理 |
|------|------|----------|
| meta_len | Max2Byte = 65535 | 報 "tag too long" |
| total_len（Tag 節點） | Max4Byte ≈ 4.29 GB | 報 "tag+payload too long" |
| body_len | 由 total_len 間接限制 | — |

---

## 5. 十六進制示例

### 5.1 單字段標量（Desc="ab"，value="hi"）

meta 正序：`[0x1A 'a' 'b']`（KDesc 前綴 0x1A = 字節長 2 + 1 個 key 字節）  
meta 逆序：`['b' 'a' 0x1A]`  
尾部槽：`[0x03]`（meta_len=3）  
body：`[0x82 'h' 'i']`（PrefixString inline len=2）  
total_len = 3 + 3 + 1 = 7 → 內聯 → 頭字節 = `0xE0 | 7 = 0xE7`

完整字節：`E7 82 68 69 62 61 1A 03`

### 5.2 多字段標量（Example+Nullable+Desc="ab"，value="x"）

meta 正序：`[0x09][0x29][0x1A 'a' 'b']`（KExample=0x09，KNullable=0x29）  
meta 逆序：`['b' 'a' 0x1A 0x29 0x09]`（5 字節）  
尾部槽：`[0x05]`  
body：`[0x81 'x']`（2 字節）  
total_len = 2 + 5 + 1 = 8 → 內聯 → 頭字節 = `0xE8`

完整字節：`E8 81 78 62 61 1A 29 09 05`

### 5.3 引用 ref_to=2（KMore 位圖）

meta 正序：`[0xF9 0x02 0x02]`（KMore=0xF9，位圖=0x02，值=0x02）  
meta 逆序：`[0x02 0x02 0xF9]`  
尾部槽：`[0x03]`  
body：`[0x81 'r']`（2 字節）  
total_len = 2 + 3 + 1 = 6 → 內聯 → 頭字節 = `0xE6`

完整字節：`E6 81 72 02 02 F9 03`

### 5.4 無 tag 對象 `{"a": 1}`（裸編碼，無 footer）

`C4 D2 81 61 21`

- `0xC4`：Container|Object|inline_len=4  
- body：`D2 81 61 21`  
  - `0xD2`：Container|Array|inline_len=2（鍵數組）  
  - `0x81 61`：PrefixString inline_len=1，key="a"  
  - `0x21`：PositiveInt 1（value=1）

---

## 6. 遷移注意事項（v0.3.x → v0.4.0）

**硬切，無兼容模式，無 wire 版本字節。**

- 舊 v0.3.x 格式數據在新 v0.4.0 解碼器上會產生錯誤或錯位，不做靜默兼容。
- 10 個語言實現逐字節對齊 v0.3.x header 模式，進入「凍結待遷移」狀態。
- Go 新實作為唯一基準，各語言按 §5 hex 示例 + 本規格文檔對齊。
- 文本側（JSONC `mm:` 註釋語法、tag key 集合、`Tag.ToString`/`ParseMMTag`）零變更。
- 源碼級用戶 API（`metamessage.go`、`mm.EncodeFromValue`/`DecodeToValue` 等）簽名不變，僅字節級不兼容。

### 各語言移植配方（三件事）

1. **meta 正序序列化器保持不變**（各語言 `Tag.Bytes()` 對應物）。
2. **封裝層改為「頭總長 + 逆序 meta + 尾部槽」**。
3. **解碼器實現尾部反向讀取**（推薦「切片反轉後複用既有正向 tag 解析器」的最小改動路線）+ body 精確消耗斷言。

建議遷移順序：mm-ts → mm-py → mm-php → mm-cs → mm-kt → mm-rs → mm-swift → mm-c → mm-cpp。

每語言驗收標準：`tests/encode_decode_test_<lang>.sh` 全綠 + `compare_with_go.sh` 逐字節對齊 Go 基準 hex。
