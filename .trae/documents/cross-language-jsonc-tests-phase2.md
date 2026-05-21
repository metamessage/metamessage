# Cross-Language JSONC Tests — Phase 2: Fix & Execute

## Summary

Phase 1 has completed: 25 JSONC fixtures collected, 10 language harnesses created, `run_cross_lang.sh` cross-language comparison script written. However, the script has bugs preventing full execution (output truncated after 2 fixtures). This plan fixes the script, runs the full test suite, and fixes any discovered IR consistency or JSONC reversibility issues.

## Current State

- **Fixtures**: 25 `.jsonc` files in 6 categories under `tests/fixtures/`
- **Harnesses**: 10 language harnesses under `tests/harness/` — Go, Python, PHP, TypeScript, Rust, C, C++, C#, Kotlin, Swift
- **Buildable**: 8/10 (Kotlin blocked by version 1.7 < 1.9, Swift blocked by missing toolchain)
- **Previous run**: Script produced output for only 2/25 fixtures, then truncated
- **Root cause**: (1) C# `dotnet build` stdout not redirected, (2) `set -euo pipefail` causes early exit on normalize pipeline failure

## Proposed Changes

### 1. Fix `tests/run_cross_lang.sh`

#### 1a. Fix C# build output leakage (line 88)

**Problem**: `dotnet build ... 2>/dev/null` only redirects stderr; Chinese locale build messages go to stdout, polluting the output stream.

**Fix**: Change to redirect both stdout and stderr:

```bash
cs_build() {
    dotnet build "$SCRIPT_DIR/harness/csharp/harness.csproj" --nologo -v q >/dev/null 2>&1
}
```

#### 1b. Fix early exit from `set -euo pipefail`

**Problem**: `set -euo pipefail` at line 5 causes the script to exit immediately when:
- `normalize()` function returns non-zero (e.g., JSON parse error on malformed output)
- Pipeline `echo "${OUTPUTS[0]}" | normalize` fails
- Any harness returns non-zero (currently guarded with `|| true`, but pipefail can still trigger)

**Fix**: Replace `set -euo pipefail` with more targeted error handling:
- Remove `set -e` from global scope
- Add explicit error checks where needed
- Keep `set -o pipefail` for pipeline integrity
- Wrap normalize calls in conditional checks

```bash
set -o pipefail
# Remove: set -euo pipefail
```

#### 1c. Add reversibility testing

**Problem**: Current script only tests that all languages produce consistent IR output (normalized JSON). The user's request also asks to test JSONC reversibility — i.e., parse JSONC → re-print → parse again → ensure JSON is unchanged.

**Fix**: Add a second phase after IR consistency comparison. For each fixture:
1. Parse JSONC with Go (reference) → get output1
2. Feed output1 back as input → parse → get output2
3. Normalize and compare output1 == output2
4. Report any differences

Implementation approach: Add a `test_reversibility()` function that takes a fixture path, runs Go harness against it, feeds output back, and compares.

#### 1d. Improve normalize function robustness

**Problem**: `normalize()` uses `sys.exit(1)` on JSON parse error, which causes the pipeline to fail but error is swallowed.

**Fix**: Make normalize return the original input on failure (for debugging) and output to stderr for diagnostics:

```python
except Exception as e:
    print(f"// NORMALIZE_ERROR: {e}", file=sys.stderr)
    print(s, end='')  # return original on failure for debugging
```

#### 1e. Add detailed diff output for mismatches

**Problem**: When outputs differ, the diff file saves raw outputs but doesn't show a clean diff.

**Fix**: Use `diff` command to generate a side-by-side comparison of normalized outputs in diff files.

### 2. Run Full Test Suite

```bash
cd /Users/lizongying/IdeaProjects/meta-message
bash tests/run_cross_lang.sh 2>&1 | tee tests/results/summary.txt
```

Expected: 25 fixtures × 8 languages = all should either MATCH or have documented differences.

### 3. Fix Discovered Issues

Based on the test results, fix any of:
- IR inconsistency between language implementations
- JSONC reversibility failures (re-parse produces different JSON)
- Parser bugs in specific language implementations

Known potential issues to watch for:
- **`null_with_tag.jsonc`**: Go parser may reject bare `null`; some languages may handle differently
- **`trailing_comma.jsonc`**: Not all parsers may accept trailing commas
- **`deep_nested.jsonc`**: Some implementations may have recursion depth limits
- **Tag handling differences**: `mm:` tag comment parsing may differ between languages
- **ValueType enum inconsistencies**: Already fixed for Python, TypeScript, C# — but remaining languages may have similar issues

### 4. Update Result Documentation

After fixes, update `tests/results/summary.txt` with final pass/fail counts and document any known limitations.

## Assumptions & Decisions

- Go implementation is the reference/ground truth
- "Reversibility" = parse JSONC → re-print → parse re-printed → normalized JSON equals original normalized JSON
- "IR consistency" = all 8 available languages produce the same normalized JSON output for the same fixture
- Kotlin and Swift are skipped (blocked by environment) — this is acceptable
- Output normalization strips comments, whitespace, trailing commas; sorts keys via `json.dumps(sort_keys=True)`

## Verification

1. `run_cross_lang.sh` runs to completion and processes all 25 fixtures
2. No Chinese locale messages leak into output
3. All 8 languages produce output for each fixture (or documented as expected failure)
4. Cross-language IR comparison shows MATCH for most fixtures
5. Reversibility test passes for most fixtures
6. Any failing fixtures have documented diffs in `tests/results/`
7. All discovered bugs are fixed or documented as known limitations