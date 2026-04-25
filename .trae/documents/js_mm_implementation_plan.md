# JavaScript MetaMessage (mm-js) 實現計劃

## 概述
基於 Go 版本的功能完整移植到 JavaScript，排除 CLI 功能，測試要完整全面，包括所有數據類型和 tag。

由於 JavaScript 沒有 Go 的 struct tag 系統，實現將採用 `mm(value, {tag})` 模式來處理類型元數據。

## 項目結構
```
mm-js/
├── package.json
├── src/
│   ├── mm/
│   │   ├── constants.js        - Wire 格式常量
│   │   ├── buffer.js           - 字節緩衝區
│   │   ├── encoder.js          - 編碼器
│   │   ├── decoder.js          - 解碼器
│   │   ├── index.js            - 公共 API
│   │   └── types.js            - 類型定義（使用 BigInt）
│   └── jsonc/
│       ├── scanner.js          - 詞法分析器
│       ├── ast.js              - AST 節點定義
│       ├── parser.js           - 語法分析器
│       ├── tag.js              - Tag 定義和解析
│       ├── binder.js           - 對象綁定器
│       └── printer.js           - 輸出格式化器
└── test/
    ├── mm/
    │   ├── encoder.test.js
    │   └── decoder.test.js
    └── jsonc/
        ├── scanner.test.js
        ├── parser.test.js
        └── tag.test.js
```

## 關鍵技術考慮

### 1. BigInt 處理
JavaScript 使用 BigInt 處理大整數：
- `int64` → `BigInt`
- `uint64` → `BigInt`
- 編碼時使用 BigInt 直接轉換
- 解碼時返回 BigInt 值

### 2. Tag 系統設計
由於 JS 沒有原生 tag 系統，採用函數式 API：
```javascript
// 設計示例
mm(123)                           // 普通數值
mm(123, { type: 'int' })          // 指定類型
mm('hello', { type: 'str' })      // 字符串
mm([1, 2, 3], { type: 'array' }) // 數組
mm({ a: 1 }, { type: 'struct' })  // 對象
```

### 3. 模塊格式
- 使用 ES Modules (ESM)
- 兼容 Node.js 和瀏覽器環境

## 待完成任務

### 1. 項目初始化
- 創建 package.json
- 配置 ES Modules

### 2. MM 模塊實現

#### 2.1 constants.js - Wire 格式常量
- 前綴枚舉（Simple, PositiveInt, NegativeInt, Float, String, Bytes, Container, Tag）
- 簡單值枚舉（null, true, false 等）
- 長度計算常量
- 前綴/後綴掩碼

#### 2.2 buffer.js - 字節緩衝區
- `MMBuffer` 類
- `writeUInt8`, `writeInt8`, `writeUInt16LE`, `writeInt16LE` 等
- `writeFloat32`, `writeFloat64`
- `readUInt8`, `readInt8`, `readUInt16LE`, `readInt16LE` 等
- `readFloat32`, `readFloat64`
- `readBytes`, `writeBytes`

#### 2.3 types.js - 類型包裝
```javascript
// 類型包裝器
export function mm(value, options = {}) {
  return { value, options }
}

// 類型工廠函數
mm.int = (v) => mm(v, { type: 'int' })
mm.str = (v) => mm(v, { type: 'str' })
mm.bool = (v) => mm(v, { type: 'bool' })
mm.float = (v) => mm(v, { type: 'float' })
mm.bytes = (v) => mm(v, { type: 'bytes' })
mm.array = (v) => mm(v, { type: 'array' })
mm.struct = (v) => mm(v, { type: 'struct' })
mm.bigint = (v) => mm(v, { type: 'bigint' })
mm.uuid = (v) => mm(v, { type: 'uuid' })
mm.datetime = (v) => mm(v, { type: 'datetime' })
mm.date = (v) => mm(v, { type: 'date' })
mm.time = (v) => mm(v, { type: 'time' })
mm.email = (v) => mm(v, { type: 'email' })
mm.url = (v) => mm(v, { type: 'url' })
mm.ip = (v) => mm(v, { type: 'ip' })
mm.decimal = (v) => mm(v, { type: 'decimal' })
```

#### 2.4 encoder.js - 編碼器
- `MMEncoder` 類
- `encodeBool`, `encodeNil`
- `encodeInt` (使用 BigInt)
- `encodeUInt` (使用 BigInt)
- `encodeFloat`, `encodeDouble`
- `encodeString`
- `encodeBytes`
- `encodeArray`
- `encodeObject` / `encodeStruct`
- 處理 int64/uint64 使用 BigInt

#### 2.5 decoder.js - 解碼器
- `MMDecoder` 類
- `decode()` - 返回 `DecodedValue`
- `decodeBool`, `decodeNil`
- `decodeInt` - 返回 BigInt
- `decodeUInt` - 返回 BigInt
- `decodeFloat`, `decodeDouble`
- `decodeString`
- `decodeBytes`
- `decodeArray`
- `decodeObject`

#### 2.6 index.js - 公共 API
```javascript
export function encode(value, options = {}) { ... }
export function decode(buffer) { ... }
export function toJSONC(buffer) { ... }
export function fromJSONC(jsoncString) { ... }
export { mm } from './types.js'
```

### 3. JSONC 模塊實現

#### 3.1 scanner.js - 詞法分析器
- Token 類型（LCURLY, RCURLY, LBRACKET, RBRACKET, COLON, COMMA, STRING, NUMBER, TRUE, FALSE, NULL, LINECOMMENT, BLOCKCOMMENT）
- `JSONCScanner` 類
- `nextToken()` 方法
- 追蹤行號和列位置
- 正確處理 `//` 行註釋和 `/* */` 塊註釋

#### 3.2 ast.js - AST 節點定義
```javascript
export class JSONCNode {}
export class JSONCValue extends JSONCNode {}
export class JSONCObject extends JSONCNode {}
export class JSONCArray extends JSONCNode {}
export class JSONCDoc extends JSONCNode {}
```

#### 3.3 tag.js - Tag 定義和解析
```javascript
export class JSONCTag {
  constructor()
  inherit(tag)
  toString()
}

export function parseMMTag(tagStr) { ... }
```

#### 3.4 parser.js - 語法分析器
- `JSONCParser` 類
- `parse()` 方法
- 遞歸下降解析
- 處理註釋並關聯到節點
- 創建完整的 AST

#### 3.5 printer.js - 輸出格式化器
- `JSONCPrinter` 類
- `print(node)` - 格式化輸出
- `printCompact(node)` - 緊湊輸出
- 處理所有節點類型

#### 3.6 binder.js - 對象綁定器
- `JSONCBinder` 類
- `bind(node, Type)` 方法
- 將 AST 綁定到 JS 對象

### 4. 測試模塊

#### 4.1 MM 編碼器/解碼器測試
- `test/mm/encoder.test.js`
- `test/mm/decoder.test.js`

測試覆蓋：
- Bool (true, false)
- Nil / Null
- Int (BigInt): 0, 23, 24, 123456, -7890
- Int8, Int16, Int32, Int64 (BigInt) 邊界值
- UInt (BigInt): 0, 987654, 255, 65535, 4294967295, 18446744073709551615
- UInt8, UInt16, UInt32, UInt64 (BigInt) 邊界值
- Float: 0.0, 3.14, -3.14, 極值
- Double: 0.0, 3.14159265359, 極值
- String: "", "hello", "hello world", 長字符串 (300+ chars)
- Bytes: empty, [0x01, 0x02, 0x03], 長 bytes (300+)
- Array: [Bool], [Int], [String], [Float], [Double], 空數組, 大數組 (1000 items)
- Object/Struct: 簡單對象, 嵌套對象

#### 4.2 JSONC 測試
- `test/jsonc/scanner.test.js` - Token 識別, 註釋, 位置追蹤
- `test/jsonc/parser.test.js` - 有效 JSONC, 帶註釋的 JSONC, 嵌套結構
- `test/jsonc/tag.test.js` - Tag 解析, 複雜 tag, 繼承

## 參考實現
- Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`
- Go JSONC: `/Users/lizongying/IdeaProjects/meta-message/internal/jsonc/`

## 實現順序
1. 項目初始化（package.json, 目錄結構）
2. MM/constants.js - 常量定義
3. MM/buffer.js - 緩衝區
4. MM/types.js - 類型包裝
5. MM/encoder.js - 編碼器
6. MM/decoder.js - 解碼器
7. MM/index.js - 公共 API
8. JSONC/scanner.js - 詞法分析器
9. JSONC/ast.js - AST 節點
10. JSONC/tag.js - Tag 解析
11. JSONC/parser.js - 語法分析器
12. JSONC/printer.js - 格式化輸出
13. JSONC/binder.js - 對象綁定
14. 測試文件實現
15. 運行測試驗證

## 關鍵注意事項
- JavaScript BigInt 沒有無符號上限，但 wire 格式需要處理 uint64
- 編碼時 int64/uint64 都使用 BigInt
- 解碼時返回 BigInt，調用者需自行轉換
- 使用 ES Modules (import/export)
- 錯誤處理使用 Error 類型
- 緩衝區使用 Uint8Array
- 測試使用 Jest 框架