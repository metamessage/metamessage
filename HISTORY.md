# 更新日誌

## v0.4.0

* **不兼容變更（硬切）**：線格式從 header 模式改為 footer 模式；每個帶 meta 的節點佈局改為 `[PrefixTag|檔位][total_len][body 正向][reverse(meta 字節)][尾部槽]`
* meta（tag）信息從節點頭部移至節點尾部；原 `encodeT`（頭部 tag_len 編碼）廢除，改由尾部三檔自描述長度槽取代
* 解碼器新增「body 精確消耗斷言」不變量：正向解析實際消耗必須等於 body_len，否則報 footer boundary mismatch
* 對象鍵數組保留在 body 頭部不遷移（顯式決策）；`decodeObj` 新增 key 數量邊界檢查，防止過長 total_len 觸發越界讀取
* 合成表新增 `ValueTypeEnum`（Data:-1）和 `ValueTypeMedia`（Data:空字節）的 is_null 處理
* 文本側（JSONC `mm:` 註釋語法、tag key 集合、`Tag.ToString`/`ParseMMTag`）零變更；源碼級用戶 API 簽名不變
* 舊 v0.3.x 格式數據在 v0.4.0 解碼器上會產生錯誤，不做靜默兼容；詳見 `docs/footer-format-v0.4.0.md`
* Go 參考實現已升級至 v0.4.0；9 個語言移植（c / cpp / cs / kt / php / py / rs / swift / ts）進入凍結待遷移狀態

## v0.3.0

* 容量擴容（破壞性變更）：Bytes / Object / Array / Tag 的長度槽從 2 檔擴展到 4 檔，單值長度上限從 65535 字節提升到 4294967295 字節（約 4GB）
* 編碼器內部緩衝 maxCap 全語言從 1GB 提升到 0xFFFFFFFF，並補充超限的顯式報錯
* String 維持 v0.2.x 的 2 檔語義，上限仍為 65535 字節，未受本次變更影響
* Go 參考實現與 9 個語言移植（c / cpp / cs / kt / php / py / rs / swift / ts）同步升級，跨語言 wire 格式保持一致
* 破壞性說明：舊數據中 Bytes 長度 28–29、Container payload 12–13 在新解碼器下會被解釋為擴展長度標記而錯位；≤27 的 Bytes、≤11 的 Container、全部 String 與全部 Tag 的編碼不變

## v0.1.5

* 修復jsonc解析錯誤問題

## v0.1.0

* 初始代碼
