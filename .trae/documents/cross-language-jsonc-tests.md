# Cross-Language JSONC Test Framework

## Summary

Create a `tests/` directory with shared JSONC fixture files, and build a Go-based test runner that validates JSONC round-trip (parse → output) consistency across all 11 language implementations. When IR or JSONC output diverges between languages, identify and fix the root cause.

## Current State Analysis

### Language Implementations

| Dir                                    | Language   | Build System  | Has JSONC parser? | Has JSONC printer? | Test runner          |
| -------------------------------------- | ---------- | ------------- | ----------------- | ------------------ | -------------------- |
| `/` (root)                             | Go         | `go test`     | Yes               | Yes                | `go test`            |
| `mm-c/`                                | C          | CMake         | Yes               | Yes                | `ctest`              |
| `mm-cpp/`                              | C++        | CMake         | Yes               | Yes                | `ctest`              |
| `mm-cs/`                               | C#         | dotnet        | Yes               | Yes                | `dotnet test`        |
| `mm-kt/`                               | Kotlin     | Maven         | Yes               | Yes                | `mvn test`           |
| `mm-php/`                              | PHP        | Composer      | Yes               | Yes                | `vendor/bin/phpunit` |
| `mm-py/`                               | Python     | pip/pyproject | Yes               | Yes                | `python -m pytest`   |
| `mm-rs/`                               | Rust       | Cargo         | Yes               | Yes                | `cargo test`         |
| `mm-swift/`                            | Swift      | SwiftPM       | Yes               | Yes                | `swift test`         |
| `mm-ts/`                               | TypeScript | npm           | Yes               | Yes                | `npm test`           |
| (Also: JS via `npm`, but shares mm-ts) | <br />     | <br />        | <br />            | <br />             | <br />               |

### Common IR Structure

All implementations share the same IR model:

* **Node types**: `Object`, `Array`, `Value`, `Doc`

* **Object**: list of `Field{key, value}`

* **Array**: list of `Node` items

* **Value**: `data`, `text`, `Tag`

* **Tag**: type, desc, default, min, max, nullable, etc.

### Existing Test Coverage (JSONC-specific)

| Language   | JSONC parse tests              | JSONC print tests | Round-trip tests             | Has standalone CLI entry |
| ---------- | ------------------------------ | ----------------- | ---------------------------- | ------------------------ |
| Go         | Yes (`jsonc_test.go`)          | Yes               | Yes                          | Yes (`cmd/mm/mm.go`)     |
| C          | Build-only (`test_jsonc.c`)    | Build-only        | Via `test_comprehensive.c`   | No                       |
| C++        | Build-only (`test_jsonc.cpp`)  | Build-only        | Via `test_comprehensive.cpp` | No                       |
| C#         | Yes (scanner/printer tests)    | Yes               | Via `ComprehensiveTests.cs`  | No                       |
| Kotlin     | Yes (`JsoncTest.kt`)           | Yes               | Via `MetaMessageTest.kt`     | No                       |
| PHP        | Yes (`JsoncTest.php`)          | Yes               | Via `MetaMessageTest.php`    | No                       |
| Python     | Yes (`test_jsonc.py`)          | Yes               | Yes                          | No                       |
| Rust       | Yes (in `lib.rs`)              | Yes               | Yes                          | No                       |
| Swift      | Yes (`JSONCParserTests.swift`) | Yes               | Yes                          | No                       |
| TypeScript | Yes (`parser.test.ts`)         | Yes               | Yes                          | No                       |

### Key Observation

Most languages have `parseJsonc()` and `toJsonc()` APIs but no standalone CLI entry point that reads a JSONC file and outputs the round-trip result. For cross-language testing, each language needs a minimal CLI harness or the Go runner needs to call each language's test framework programmatically.

## Proposed Changes

### Phase 1: Create `tests/fixtures/` with JSONC input files

**File:** **`tests/fixtures/`** **directory**

Create categorized JSONC fixture files covering all edge cases found across existing tests. Each file is a valid JSONC input. Categories:

```
tests/
├── fixtures/
│   ├── 01_primitive/
│   │   ├── empty_object.jsonc
│   │   ├── empty_array.jsonc
│   │   ├── string.jsonc
│   │   ├── integer.jsonc
│   │   ├── float.jsonc
│   │   ├── boolean.jsonc
│   │   └── null_with_tag.jsonc
│   ├── 02_nested/
│   │   ├── nested_object.jsonc
│   │   ├── nested_array.jsonc
│   │   └── mixed_nested.jsonc
│   ├── 03_tags/
│   │   ├── type_tags.jsonc
│   │   ├── desc_tag.jsonc
│   │   ├── nullable_tag.jsonc
│   │   ├── constraints_tag.jsonc
│   │   └── complex_tags.jsonc
│   ├── 04_comments/
│   │   ├── line_comments.jsonc
│   │   ├── block_comments.jsonc
│   │   └── mixed_comments.jsonc
│   ├── 05_edge_cases/
│   │   ├── trailing_comma.jsonc
│   │   ├── negative_numbers.jsonc
│   │   ├── empty_string.jsonc
│   │   └── deep_nested.jsonc
│   └── 06_complex/
│       ├── person.jsonc
│       ├── product.jsonc
│       └── config.jsonc
```

Each fixture file contains a valid JSONC string that all parsers should handle identically.

### Phase 2: Create Go reference test runner

**File:** **`tests/runner.go`**

A Go program that:

1. Reads all `.jsonc` fixture files from `tests/fixtures/`
2. For each fixture:
   a. Parse JSONC → Go IR
   b. Print Go IR → JSONC string
   c. Parse the re-printed JSONC again → Go IR2
   d. Compare IR == IR2 (structural equality)
   e. Normalize and compare the original and re-printed JSONC
3. Output: table of results (pass/fail per fixture)

This establishes the Go reference behavior first.

### Phase 3: Create language-specific test harnesses

For each language, create a small CLI tool at `tests/harness/<lang>/` that:

* Reads a JSONC file path as argument

* Parses it to IR

* Prints IR back to JSONC (to stdout)

* Exits with code 0 on success, non-zero on error

The harness directory:

```
tests/harness/
├── mm-go/        → uses `mm.ParseJsonc()` + `mm.ToJsonc()`
├── mm-c/         → C program linked to libmmc
├── mm-cpp/       → C++ program linked to libmmcpp
├── mm-cs/        → dotnet console app
├── mm-kt/        → Kotlin main class
├── mm-php/       → PHP CLI script
├── mm-py/        → Python script
├── mm-rs/        → Rust binary
├── mm-swift/     → Swift executable
└── mm-ts/        → Node.js script
```

### Phase 4: Cross-language comparison script

**File:** **`tests/run_cross_lang.sh`** (or `Makefile` target)

A script that:

1. Builds all harnesses
2. For each fixture, runs each language harness, captures JSONC output
3. Normalizes all outputs (whitespace-insensitive comparison)
4. Reports differences between languages
5. Flags any fixture where languages produce inconsistent output

### Phase 5: Fix discovered issues

After running the cross-language tests, fix any inconsistencies found in language implementations. Common potential issues:

* Type inference differences (e.g., int vs i vs i64)

* Tag serialization format differences (`// mm:` prefix handling)

* Null handling (some allow bare `null`, others require `is_null` tag)

* Float precision differences

* Comment preservation differences

* Field ordering differences

## Assumptions & Decisions

1. **Go as reference**: Go is the primary implementation and will serve as the reference for IR comparisons
2. **JSONC text comparison**: Compare canonical JSON (not IR trees directly) since each language has its own IR node types. The comparison normalizes whitespace and field ordering.
3. **Incremental approach**: Phase 1+2 first (fixtures + Go runner), then add language harnesses incrementally.
4. **Each language harness is minimal**: \~30-50 lines, just enough to call `parse → print` and exit.
5. **Not all languages need simultaneous building on CI**: The runner gracefully skips languages that can't be built in the current environment.

## Verification Steps

1. Run `go run tests/runner.go` — all fixtures pass round-trip in Go
2. Build and run each language harness manually: `tests/harness/mm-py/run.sh tests/fixtures/01_primitive/string.jsonc`
3. Run `tests/run_cross_lang.sh` — compare outputs across all available languages
4. For any mismatched output, file the issue and fix in the corresponding language implementation

