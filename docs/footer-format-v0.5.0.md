# MetaMessage v0.5.0 全尾部定界（直接尾序写出）線格式規格

版本：v0.5.0（硬切，不兼容 v0.4.0 / v0.3.x）
實作基準：Go 參考實現 `internal/core`（各 `encode_*.go` 原語直接尾序寫出，解碼剝幀 `tailframe.go`，尾部槽 `footer.go`）
文件狀態：規格已凍結，取代 `footer-format-v0.4.0.md`（後者標記 superseded）。

---

## 0. 設計原則（意圖與邊界）

- **唯一目的**：讓編碼「不用回寫長度」——編碼器原語**一趟正向直接按尾定界順序寫出**：先寫 body，再寫大端長度，最後寫 sign。body 寫完時其長度天然已知，直接追加定界字節即可，**既無回寫、也無事後反轉整幀**；解碼**從後往前**定界。
- **編碼邏輯零改動**：v0.3.x 的**檔位/內聯語義原樣保留**（`constants.go` 的 `IntLen/FloatLen/StringLen/BytesLen/ContainerLen/TagLen` 低 5 位檔位規則不動），節點**寬度不變**。本版本**不新增**任何獨立 `payload_len` 五檔槽，也不丟棄任何內聯壓縮檔位。
- **實現方式 = 原語物理寫序尾化（非反轉層）**：
  - 各 `encode*` 函數僅調整**輸出順序**（body→lenBE→sign），檔位判斷/內聯壓縮邏輯一字不動。
  - 容器先遞歸寫出子節點（已是 wire 形態），正序拼接成 payload 後再追加 lenBE + sign。
  - tag 段先寫 inner（wire），再寫 `reverse(meta)`、meta_len 尾槽、total_len 大端、tag sign。
  - `Decode()` 入口調用 `Untailify(wire)`：從尾反向定界還原成头部定界形態後，餵給**未經改動的正向解码器**。

---

## 1. 節點布局（wire 形態）

記 head 形態（v0.3.x 邏輯形態，僅供對照）→ wire 形態（落盤字節，編碼器直接產出）。所有類型**保持節點總寬度不變**（僅物理寫序尾化）。

### 1.1 Simple
head `[sign]`（單字節）→ wire 原樣透傳 `[sign]`。無 body、無長度槽。

### 1.2 Int
- 內聯檔（幅值 fit 在 sign 低 5 位，`IntLen` 返回 `w==0`）：head `[sign]` 單字節 → wire `[sign]` 原樣。
- 帶檔（`w = IntLen 返回的 1..8`）：head `[sign][body w 字節大端]` → wire `[body w 字節][sign]`。
  長度由 sign 低 5 位檔位反推，**無獨立長度字節**。

### 1.3 Float
- 內聯檔（`exponent==-1 && mantissa<=7`，`FloatLen` 返回 `w==0`）：head `[sign]`（如 `-0.1 = 0x71`）→ wire 原樣單字節。
- 帶檔（`w = 1..8`）：head `[sign][exp int8][mantissa 大端 w 字節]` → wire `[exp][mantissa][sign]`。
  sign 低 4 位承載檔位（bit4=負號），body 實寬 = `w+1`（含 exp）。

### 1.4 String / Bytes
- 內聯檔（`w==0`，長度值在 sign 低 5 位）：head `[sign][body inline]` → wire `[body inline][sign]`（**非單字節**，body 仍存在）。
- 帶檔（String：`w=1/2`；Bytes：`w=1..4`）：head `[sign][lenBE w][body]` → wire `[body][lenBE][sign]`，`lenBE` 緊鄰 sign 左側。

### 1.5 Container（Array / Object）
- head `[sign][可選 lenBE][payload]`，payload = 子節點 head 形態正序拼接。
- wire `[子節點 wire 形態正序拼接][可選 lenBE][sign]`。
  子節點**遞歸換位**後仍**正序**排列（換位不改子節點先後順序）。
- sign bit4 = Array(1)/Object(0)；`ContainerLen` 低 4 位為檔位（內聯 0..11，檔位 12..15 = 1..4 字節）。

### 1.6 帶 meta 的節點（tag 包裝段）
meta 非空時值節點外套一層 `PrefixTag(0b111)` 包裝段：
```
head: [tag sign][可選 total_len 大端][meta 前綴槽][meta 正序][inner(值節點 head 形態)]
wire: [inner wire 形態][reverse(meta 正序)][meta_len 尾槽][可選 total_len 大端 BE][tag sign]
```
- `total_len`（head 語義 = `len(tag段簽名後) = 前綴槽寬 + meta + inner`）由 `TagLen` 低 5 位檔位決定是否佔用大端字節（內聯 `total<28` 無字節）。
- `reverse(meta)`：meta 正序字節流整體字節逆序，長度 = `meta_len`；反轉回正序後複用既有正向 `decodeTagBytes` 字段循環。
- `meta_len 尾槽`：見 §3.2，位於 tag sign 左側、reverse(meta) 右側。
- 空 inner 合法（如 `is_null` 合成場景）：`inner wire` 區寬 0。

---

## 2. 類型前綴與檔位（沿用 constants.go，僅物理位移）

高 3 位類型前綴**完全沿用 v0.3.x**，低 5 位**仍承載檔位/內聯長度**（不同於舊版草案「低 5 位歸零 + 五檔槽」方案）：

| 高 3 位 | 常量 | 低 5 位用途（不變） |
|--------|------|--------------------|
| 0b000 | Simple | 整個字節 = SimpleValue id（0..31）|
| 0b001 | PositiveInt | `IntLen`：0..23 內聯幅值；24..31 = 1..8 字節 body |
| 0b010 | NegativeInt | 同 PositiveInt（取負）|
| 0b011 | PrefixFloat | bit4=負號；低 4 位 `FloatLen`：0=內聯(exp=-1,mant<=7)，1..8=mantissa 字節數 |
| 0b100 | PrefixString | `StringLen`：<30 內聯長度；30=1 字節 lenBE；31=2 字節 lenBE |
| 0b101 | PrefixBytes | `BytesLen`：<28 內聯；28..31 = 1..4 字節 lenBE |
| 0b110 | Container | bit4=Array(1)/Object(0)；低 4 位 `ContainerLen`：<12 內聯；12..15=1..4 字節 lenBE |
| 0b111 | PrefixTag | `TagLen`：<28 內聯 total；28..31 = 1..4 字節 total_len BE |

> 換位層據 sign 低 5 位檔位（`tierOf`）反推 body/payload 寬度，從而在尾部完成定界。

---

## 3. 尾部長度槽

### 3.1 值節點：無獨立 payload_len 槽
值節點（§1.1–§1.5）**不設**任何獨立 payload_len 槽。節點寬度完全由 sign 低 5 位檔位 +（container/str/bytes 帶檔時的 lenBE 字節）確定，與头部定界一一對應。

### 3.2 meta_len 尾槽（僅 tag 包裝段，沿用 v0.4.0 `parseFooterSlot`）
從 tag sign 左側反向讀，末字節為標記或直接值：
| 末字節 `v` | 含義 | 槽寬 |
|-----------|------|------|
| `v <= 253` | meta_len = `v` | 1 |
| `v == 254` | 左鄰 1 字節（u8）= meta_len | 2 |
| `v == 255` | 左鄰 2 字節（u16 大端）= meta_len（≤ Max2Byte）| 3 |

常量：`FooterSlotDirectMax=253`、`FooterSlotMarker1=254`、`FooterSlotMarker2=255`。canonical 最小寬度。

**前綴槽與尾槽統一**（v0.5.0 根治）：头部 meta 前綴槽與尾槽使用**同一套 canonical 三档**（`ml<254→1`、`254..255→2`、`256..65535→3`）。
舊 v0.3.x 前綴槽實際只有兩档（`[254, ml&0xFF]` 對一切 `ml>255` mod 256 回繞 → 静默損壞），且 `ml<257→2` 邊界使 `ml=256` 與尾槽寬 3 不一致（+1 字節 quirk）。
v0.5.0 將兩者統一：反推式改用實際尾槽寬 `inner = total - ml - slotW`，正向解码器 meta 前綴的第三档（`b==255 → u16`，舊代碼 `l<257` 恒真導致的死代碼）與 meta 串属性 payload==7 的 2 字節档同步激活；`ml` 合法域 0..65535 現全區間可 round-trip。

---

## 4. 解碼算法（Untailify → 正向解码器）

`Decode()` 不在反向形態上直接建 IR，而是先**整體還原成头部定界形態**，復用零改動的正向解码器：

```
Untailify(wire):
    start = nodeTailStart(wire, len(wire))   // 據尾 sign 反推根節點起點
    if start != 0: 報 "footer boundary mismatch"
    return untailifyNode(wire, start, len(wire))   // 遞歸還原成 head 字節
```

`nodeTailStart(data, end)` 讀 `data[end-1]`（尾 sign）→ 按 §2 檔位反推節點起點：
- Simple / int 內聯 / float 內聯：`end-1`
- int 帶檔：`end-1-w`；float 帶檔：`end-2-w`（含 exp）
- str/bytes/container 內聯：`end-1-inline`；帶檔：`end-1-w-l`（l 由 sign 左側 lenBE 讀）
- tag：先讀 total_len（尾 sign 左側 tier），剝 meta_len 槽 → 反推 reverse(meta) 邊界 → inner 區在槽左側

`untailifyNode` 據此逐段還原成 head 字節；容器 payload 內子節點**從區域末端反向定界掃描**、還原後按正序拼接（寬度不变保證與头部等長）。還原得到的 head 形態交由既有 `decode/decodeTag/decodeString/decodeContainer/…` 正向解析。

---

## 5. 邊界與篡改防護

| 條件 | 處理 |
|------|------|
| `end` 越界 / wire 為空 | 報 truncated |
| 根節點起點 != 0 | 報 `footer boundary mismatch`（如砍掉尾 sign）|
| int/float/str/bytes/container 檔位字節越界 | 報 `truncated: tier bytes out of range` |
| tag `total_len` 反推 inner < 0 | 報 `tag total underflow` |
| 容器子節點反向定界邊界不吻合 | 報 `bad boundary` / body 與 tier 不符 |
| 前綴槽/尾槽寬度不一致或 total 跨檔 | 已根治：兩槽同一套三档，反推式 `inner = total - ml - slotW`；約束不符時報 `tag total underflow` / `inner want` 錯位錯誤 |

---

## 6. 引用 ref_id / ref_to 順序

換位後解码先還原成头部正向形態再走原正向解码器，`ref_to` 遍歷順序與 v0.3.x **完全一致**（前向），無需在反向幀上做回填 pass。原正向 ref 解析規則保持不變。

## 7. 流式解碼

`Decode(bytes)` 需持有整幀（根節點右邊界 = 幀長）。`Untailify` 需先看到尾 sign，故真·邊到邊增量在尾定界下不成立；`DecodeStream` 維持「先讀盡緩衝、再還原正向解析」語義。

---

## 8. 十六進制示例（Go 實現實測）

以 `internal/core/tailframe_test.go` 斷為準：

- float `-0.1`（內聯）：`71`（單字節，與 v0.3.x 相同）。
- int `5`（PositiveInt 內聯）：`25`（單字節）。
- int `300`（帶檔 w=2）：head `39 01 2C` → wire `01 2C 39`（body 原序 + sign 挪尾；sign = PositiveInt|IntLen2Byte）。
- string 52 字節（帶檔 w=1）：`[52 字節 body][34][9E]`（body 後緊 lenBE=52=0x34、再 sign=PrefixString|StringLen1Byte）。
- 容器 `{"a":1,"b":2}`：子節點換位後正序拼接，尾為 container sign。
- tag 包裝（root，desc=hello）：`[inner wire][reverse(meta)][meta_len 槽][可選 total BE][0xE0|tier]`，尾高 3 位必為 `0b111`。

---

## 9. 遷移注意事項（v0.3.x / v0.4.0 → v0.5.0）

**硬切，無兼容。**
- 舊 v0.4.0/v0.3.x 數據在新解码器上錯位/報錯，不靜默兼容。
- v0.3.x 各類型**編碼內容、檔位、內聯壓縮全部保留**；本版本僅把自描述信息物理尾化。
- **舊版草案的「五檔 payload_len 槽 / 丟棄內聯」方案作廢**，以本文為準。
- meta 正序序列化器（`Tag.Bytes()`）與文本側（JSONC `mm:` 註釋、tag key 集合、`Tag.ToString`/`ParseMMTag`）零變更。
- 源碼級 API（`metamessage.go`、`mm.EncodeFromValue`/`DecodeToValue`）簽名不變，僅字節級不兼容。
- golden hex：帶 body 的多字節節點字節序會變（sign 從頭挪到尾），全量重算；單字節內聯（Simple / int 內聯 / float 內聯）**hex 不變**。

### 各語言移植配方
1. meta 正序序列化器保持不變。
2. **編碼器原語改寫序**（非加反轉層）：把各 `encode*` 函數的輸出順序改為 body→lenBE→sign（int/float 為 body→sign）；容器先遞歸寫子節點再追加定界；tag 先寫 inner、再寫 reverse(meta)、meta_len 尾槽、total_len 大端、sign。檔位/內聯判斷邏輯一字不動。**不要在出口做整幀反轉**（那等於把回寫换个形式，反而多一遍拷贝）。
3. **解碼入口剝幀**（等價 `tailframe.go` 的 `Untailify`/`nodeTailStart`/`untailifyNode`）：從尾 sign 反向定界把 wire 還原成头部形態，再餵給各語言 v0.3.x 原正向解码器。容器子節點從區域末端反向掃描、還原後正序拼接。
4. tag 段按 §3.2 meta_len 槽；前綴槽與尾槽為同一套 canonical 三档（還原用 `metaPrefixBytes`，反推用實際槽寬 `slotW`）。
5. **不要**引入獨立 payload_len 五檔槽，**不要**丟棄內聯檔位。

遷移順序：Go 基準先行 → mm-ts → mm-py → mm-php → mm-cs → mm-kt → mm-rs → mm-swift → mm-c → mm-cpp。
每語言驗收：`tests/encode_decode_test_<lang>.sh` 全綠 + `compare_with_go.sh` 逐字節對齊 Go。

> **返工提示**：先前 mm-ts/py/php/cs 若按旧「五檔槽」草案或「出口整幀反轉」遷移，均与本規格不符，需按「編碼器直接尾序寫出 + 解碼入口 Untailify」思路返工重新對齊 0 diff。
