# Cross-Language MetaMessage Tests

Shared JSONC fixtures and test runners for validating consistency across all MetaMessage language implementations.

## Structure

```
tests/
├── fixtures/          # Shared JSONC test input files
│   ├── 01_primitive/  # Basic types: string, int, float, bool, null, empty
│   ├── 02_nested/     # Nested objects, arrays, mixed
│   ├── 03_tags/       # mm: tag annotations
│   ├── 04_comments/   # Line/block/mixed comments
│   ├── 05_edge_cases/ # Trailing comma, negative numbers, empty string, deep nesting
│   ├── 06_complex/    # Real-world scenarios: person, product, config
│   └── 07_large/      # v0.3.0 容量扩容：>65536 字节大 payload（3/4 字节长度档）
├── harness/           # Language-specific CLI harnesses (one per language)
├── runner.go          # Go reference test runner
├── compare_with_go.sh # 单语言 vs Go 字节/还原一致性对比
├── boundary_cross_lang.sh # v0.3.0 容量档位跨语言门禁（大 payload + 截断负例）
└── run_cross_lang.sh  # Cross-language comparison script
```

## Quick Start

### Run Go reference tests

```bash
go run tests/runner.go
```

### Run all cross-language tests

```bash
./tests/run_cross_lang.sh
```

### Compare a single language against the Go reference

对全部 fixture 验证：编码字节与 Go 完全一致 + 还原 JSONC 与 Go 一致。

```bash
./tests/compare_with_go.sh <lang>   # lang: ts c cpp cs kt php py rs sw
```

### Capacity boundary gate (v0.3.0)

验证 4 字节长度档扩容：大 payload（`07_large/big_object.jsonc`，越过 65536→3B 档）在所有语言 decode 与 Go 一致；向 stdin 喂入截断的 3B/4B 档位 wire 时，所有语言必须显式报错（非零退出、不崩溃(非 128+ 信号)、不静默(exit 0)、不超时）。

```bash
./tests/boundary_cross_lang.sh          # 全部 10 语言
./tests/boundary_cross_lang.sh cs php   # 指定语言
```

## Fixture Format

Each `.jsonc` file contains a valid JSONC input that all parsers should handle. Comments may include `mm:` tag annotations.

## Known Limitations

| Issue | Fixture | Description |
|-------|---------|-------------|
| Null literals | `01_primitive/null_with_tag.jsonc` | Go parser rejects bare `null`; requires `is_null` tag implicitly. Some other languages handle null natively. |
| Empty arrays/slices | — | Go validator rejects empty slices by default; auto-inserts example values. |
| Zero as empty | — | Go validator treats `0` / `0.0` as "empty" for int/float types without `allow_empty`. |

## Fixture Categories

### 01_primitive
Basic JSON value types: objects, arrays, strings, integers, floats, booleans.

### 02_nested
Nested structures: object-in-object, array-in-object, mixed nesting.

### 03_tags
MetaMessage `mm:` tag annotations: type, desc, nullable, constraints (min/max/size).

### 04_comments
JSONC comment handling: line comments (`//`), block comments (`/* */`), mixed.

### 05_edge_cases
Borderline cases: trailing commas, negative numbers, empty strings, deep nesting (>3 levels).

### 06_complex
Real-world scenarios combining multiple features: person, product, config objects.

### 07_large
v0.3.0 容量扩容验证。`big_object.jsonc` 含一个 70000 元素内联整数数组，编码后超过 65536 字节，触发 Object/Array 的 3 字节长度档（首字节 0xCE）。用于跨语言验证大长度槽的编码字节一致与还原无损。