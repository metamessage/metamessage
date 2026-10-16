# 更新日誌

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
