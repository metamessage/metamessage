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
│   └── 06_complex/    # Real-world scenarios: person, product, config
├── harness/           # Language-specific CLI harnesses (one per language)
├── runner.go          # Go reference test runner
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