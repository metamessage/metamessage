# JSONC 解析器多語言對齊 Go 實現 Spec

## Why

目前各語言（Python、PHP、TypeScript、Rust、C、C++、C#、Kotlin、Swift）的 JSONC 解析器實作與 Go 參考實現存在行為差異，導致跨語言測試可能因解析器細節不一致而失敗。需要將各語言解析器全面對齊 Go 實現。

## What Changes

1. **Token 定義對齊** — 所有語言 Token 結構體增加 `Column` 字段（缺失者），確保與 Go 的 `{Type, Literal, Line, Column}` 一致
2. **Scanner 行為對齊** — scanString 轉義處理、scanComment TrimSpace、scanLiteral 不 lower case 關鍵字等
3. **Parser 邏輯對齊** — `parseCommentsToTag` 使用 prefix 檢查、`consumeCommentsFor` 間隔邏輯、type validation、CamelToSnake、tag inheritance、depth 限制、null 報錯
4. **JSONC Printer (ToJSONC) 對齊** — tag 序列化格式、indent 行為一致
5. **測試覆蓋** — 每種語言增加 parser_test 對齊 Go 測試案例

## Impact

- Affected specs: Scanner, Parser, Token, JSONC Printer in all languages
- Affected code: 9 語言的 jsonc/ 目錄下所有 scanner/parser/token/printer 文件

## Key Reference Behaviors (Go)

### Token (`internal/jsonc/token/token.go`)
```go
type Token struct {
    Type    Type   // 13 types: EOF, LBrace, RBrace, LBracket, RBracket, Colon, Comma, String, Number, True, False, Null, Comment
    Literal string
    Line    int
    Column  int
}
```

### Scanner (`internal/jsonc/scanner/scanner.go`)
- `scanString`: 遇到 `\` 時寫入轉義序列原樣 (`buf.WriteRune(ch)` + `buf.WriteRune(s.next())`)，不做實際轉義
- `scanComment`: `//` 註釋內容用 `strings.TrimSpace` 去除前後空格
- `scanLiteral`: **不 lower case**，直接比對 `ir.True`/`ir.False`/`ir.Null`，其餘返回 Number
- `skipWhitespace`: 使用 `unicode.IsSpace`

### Parser (`internal/jsonc/parser/parser.go`)
- `parseCommentsToTag`: 使用 `strings.CutPrefix(cs, "mm:")` 檢查前綴
- `consumeCommentsFor`: 檢查 `anchorLine - last.Line > 1`，使用 `ir.MergeTag` 合併多個 tag
- `parse`: 呼叫 `p.parse("", false, tag)` 傳遞 `example` flag
- `parseObject`: key 用 `utils.CamelToSnake` 轉換，支援 Map/Obj 類型，`tag.Inherit()`，`tag.ValidateObj()`
- `parseArray`: 支援 Arr/Vec 類型，`tag.ValidateArr/ValidateVec`
- Type validation: 完整覆蓋 Str/Bytes/Datetime/Date/Time/Uuid/Decimal/Ip/Url/Email/Enum/Media/I/I8/I16/I32/I64/U/U8/U16/U32/U64/F32/F64/Bigint/Bool
- Max depth: 32
- Null: 返回 `"null is not supported"`

### JSONC Printer (`internal/jsonc/jsonc.go`)
- `ToJSONC`: 格式化輸出，`// mm: tag` 註釋寫入，`\t` indent，trailing comma
- Tag 輸出順序：desc → type → nullable → is_null → deprecated → allow_empty → unique → default_val → min → max → size → enums → pattern → location → version → mime → child_* 屬性

## Per-Language Differences

### Python
- Token 缺少 `Column` 字段
- `_identifier` lowercases before matching（Go 不做 lowercase）
- `_number` 返回 float/int（Go 返回 string literal）
- `_parse_mm_tag` 先 strip 再檢查前綴
- `_parse_object` 無 `CamelToSnake`、無 `tag.Name` path
- `_parse_array` 的 tag2 merge 被註解掉
- `_parse_value` 無完整 type validation
- 無 max depth 檢查

### PHP
- Token 缺少 `Column` 字段
- Scanner 可能缺少完整轉義處理
- Parser 缺少完整 type validation

### TypeScript
- `parseObject` 無 `CamelToSnake`
- inherit 使用 `hasChildFields` 條件判斷（Go 無此條件）
- 額外的 `revalidateValue` 邏輯（Go 無）
- Scanner 的 consumeString 有基本轉義處理

### Rust
- `scan_identifier` lowercases before matching（Go 不做 lowercase）
- Parser 只有 Datetime/Uuid 等基本類型驗證
- 無 `CamelToSnake`
- Inline comment handling（檢查 next token 是否為同行的 comment）— Go 無此行為

### C/C++/C#/Kotlin/Swift
- 需逐一檢查 Token/Scanner/Parser 與 Go 的差異
- 共同缺失：CamelToSnake、完整 type validation、depth 限制、tag merge 邏輯

## ADDED Requirements

### Requirement: Token Column Field
The system SHALL include a `Column` field in Token structs for all languages where it's missing.

#### Scenario: Token creation
- **WHEN** scanner creates a token
- **THEN** the token SHALL include `Column` (starting column position), matching Go's `Token.Column`

### Requirement: Scanner scanString Behavior
The system SHALL implement scanString to write escape sequences as raw text without interpreting them, matching Go behavior.

#### Scenario: Escaped characters in string
- **WHEN** scanner encounters `\"`, `\\`, `\n` in a string
- **THEN** the raw escape sequence SHALL be preserved as-is in the literal (e.g. `\"` stays as `\"`, not `"`)

### Requirement: Scanner scanComment Trim
The system SHALL trim whitespace from comment content using `TrimSpace`-equivalent, matching Go behavior.

#### Scenario: Comment with surrounding spaces
- **WHEN** scanner scans `//  mm: desc="hello"  `
- **THEN** the comment literal SHALL be `mm: desc="hello"` (trimmed)

### Requirement: Scanner scanLiteral No Lowercasing
The system SHALL NOT lowercase identifier text before matching true/false/null keywords, matching Go behavior.

### Requirement: Parser parseCommentsToTag
The system SHALL use prefix-based checking (equivalent of `strings.CutPrefix(cs, "mm:")`) to detect mm: tags.

### Requirement: Parser consumeCommentsFor
The system SHALL:
- Check `anchorLine - lastCommentLine > 1` to clear stale pending comments
- Merge multiple pending comments into a single Tag using `MergeTag`-equivalent
- Clear pending after consumption

### Requirement: Parser Object Key CamelToSnake
The system SHALL convert object keys from CamelCase to snake_case during parsing.

#### Scenario: CamelCase key
- **WHEN** parser encounters key `"userName"`
- **THEN** the field key SHALL be stored as `user_name`

### Requirement: Parser Tag Inheritance
The system SHALL support tag inheritance:
- For Map type objects, children SHALL inherit parent tag via `Inherit`-equivalent
- For Array children, SHALL inherit parent tag via `Inherit`-equivalent

### Requirement: Parser Depth Limit
The system SHALL enforce a maximum parsing depth of 32.

#### Scenario: Deep nesting
- **WHEN** JSONC nesting depth exceeds 32
- **THEN** parser SHALL return an error

### Requirement: Parser Null Literal Error
The system SHALL return an error with message `"null is not supported"` when encountering a null literal.

### Requirement: Type Validation
The system SHALL validate parsed values against tag type constraints for all supported ValueTypes: Str, Bytes, Datetime, Date, Time, Uuid, Decimal, Ip, Url, Email, Enum, Media, I/I8/I16/I32/I64, U/U8/U16/U32/U64, F32/F64, Bigint, Bool.

### Requirement: JSONC Printer Tag Output
The system SHALL output tags in `// mm: key=val` format, matching Go's output order and format.

## MODIFIED Requirements

### Requirement: Existing Scanner (all languages)
Scanner implementations SHALL be modified to match Go's exact behavior for:
- String scanning (raw escape preservation)
- Comment scanning (trim)
- Literal scanning (no lowercase, true/false/null keywords)
- Token position tracking (line + column)

### Requirement: Existing Parser (all languages)
Parser implementations SHALL be modified to include:
- `parseCommentsToTag` with mm: prefix detection
- `consumeCommentsFor` with line gap logic
- Object key conversion (CamelToSnake)
- Tag inheritance (Inherit)
- Full type validation per tag type
- 32-level depth limit
- Null literal error
- `tag.Name` based path construction

## REMOVED Requirements

### Requirement: Inline Comment After Value (Rust)
**Reason**: Go parser does not handle comments on the same line as a value token in parse_node
**Migration**: Remove the post-value inline comment handling in Rust's parse_node

### Requirement: hasChildFields Conditional Inherit (TypeScript)
**Reason**: Go parser does not condition inherit on `hasChildFields`; inherit always applies for Map/Array
**Migration**: Remove the `hasChildFields` check in TypeScript's parseObject/parseArray, apply inherit unconditionally for Map/Array