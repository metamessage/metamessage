#!/usr/bin/env bash
# Cross-language JSONC test runner
# Builds and runs harnesses for all MetaMessage language implementations
# against shared fixtures, then compares normalized outputs.
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
FIXTURES_DIR="$SCRIPT_DIR/fixtures"
RESULTS_DIR="$SCRIPT_DIR/results"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

PASS=0
FAIL=0

# ---------------------------------------------------------------------------
# Normalize JSONC for comparison: remove comments, whitespace, trailing commas,
# then sort keys.
# ---------------------------------------------------------------------------
normalize() {
    python3 -c "
import sys, json, re
s = sys.stdin.read()
s = re.sub(r'//[^\n]*', '', s)
s = re.sub(r'/\*[\s\S]*?\*/', '', s)
s = re.sub(r',(\s*[}\]])', r'\1', s)
try:
    obj = json.loads(s)
    print(json.dumps(obj, separators=(',', ':'), sort_keys=True))
except Exception as e:
    import sys as _sys
    _sys.stderr.write(f'NORMALIZE_ERROR: {e}\n')
    # return original (stripped) input for debugging
    print(s.strip(), end='')
" 2>/dev/null
}

# ---------------------------------------------------------------------------
# Language harnesses: each function builds (if needed), then runs.
# Exit 0 if harness is ready; non-zero to skip this language.
# ---------------------------------------------------------------------------

# --- Go ---
go_build() { return 0; }
go_run() { go run "$SCRIPT_DIR/harness/go/harness.go" "$1"; }

# --- Python ---
py_build() { return 0; }
py_run() { python3 "$SCRIPT_DIR/harness/python/harness.py" "$1"; }

# --- PHP ---
php_build() {
    [ -f "$PROJECT_DIR/mm-php/vendor/autoload.php" ] || return 1
    return 0
}
php_run() { php "$SCRIPT_DIR/harness/php/harness.php" "$1"; }

# --- TypeScript ---
ts_build() {
    cd "$PROJECT_DIR/mm-ts" && npm run build --silent 2>/dev/null
}
ts_run() { node "$SCRIPT_DIR/harness/typescript/harness.cjs" "$1"; }

# --- Rust ---
rs_build() {
    cd "$SCRIPT_DIR/harness/rust" && cargo build --quiet 2>/dev/null
}
rs_run() { "$SCRIPT_DIR/harness/rust/target/debug/mm-harness-rs" "$1"; }

# --- C ---
c_build() {
    mkdir -p "$SCRIPT_DIR/harness/c/build"
    cd "$SCRIPT_DIR/harness/c/build" && cmake .. -DCMAKE_BUILD_TYPE=Release >/dev/null 2>&1 && make -j$(sysctl -n hw.ncpu 2>/dev/null || echo 4) >/dev/null 2>&1
}
c_run() { "$SCRIPT_DIR/harness/c/build/mm_harness_c" "$1"; }

# --- C++ ---
cpp_build() {
    mkdir -p "$SCRIPT_DIR/harness/cpp/build"
    g++ -std=c++17 -I "$PROJECT_DIR/mm-cpp/src" -o "$SCRIPT_DIR/harness/cpp/build/mm_harness_cpp" "$SCRIPT_DIR/harness/cpp/harness.cpp" "$PROJECT_DIR/mm-cpp/src/jsonc/scanner.cpp" 2>/dev/null
}
cpp_run() { "$SCRIPT_DIR/harness/cpp/build/mm_harness_cpp" "$1"; }

# --- C# ---
cs_build() {
    dotnet build "$SCRIPT_DIR/harness/csharp/harness.csproj" --nologo -v q >/dev/null 2>&1
}
cs_run() { dotnet run --project "$SCRIPT_DIR/harness/csharp/harness.csproj" --no-build -- "$1"; }

# --- Kotlin ---
kt_build() {
    # Require Kotlin 1.9+ (mm-kt was compiled with 1.9)
    local kt_ver=$(kotlin -version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1)
    if [ -z "$kt_ver" ] || [ "$(printf '%s\n' "1.9" "$kt_ver" | sort -V | head -1)" != "1.9" ]; then
        return 1
    fi
    cd "$PROJECT_DIR/mm-kt" && mvn compile -q -DskipTests 2>/dev/null
}
kt_run() {
    local tmpdir=$(mktemp -d)
    kotlinc "$SCRIPT_DIR/harness/kotlin/harness.kt" -cp "$PROJECT_DIR/mm-kt/target/classes" -d "$tmpdir/harness.jar" 2>/dev/null && \
    kotlin -cp "$tmpdir/harness.jar:$PROJECT_DIR/mm-kt/target/classes" HarnessKt "$1" 2>/dev/null
    local rc=$?
    rm -rf "$tmpdir"
    return $rc
}

# --- Swift ---
sw_build() {
    cd "$SCRIPT_DIR/harness/swift" && swift build --quiet 2>/dev/null
}
sw_run() { swift run --package-path "$SCRIPT_DIR/harness/swift" --skip-build "$1" 2>/dev/null; }

# Ordered list of language keys
LANGS="go py php ts rs c cpp cs kt sw"
LANG_NAMES="Go Python PHP TypeScript Rust C C++ C# Kotlin Swift"

# ---------------------------------------------------------------------------
# Build phase
# ---------------------------------------------------------------------------
echo -e "${CYAN}=== Building harnesses ===${NC}"

AVAILABLE=()
for lang in $LANGS; do
    printf "  %-12s " "$lang"
    if "${lang}_build"; then
        echo -e "${GREEN}OK${NC}"
        AVAILABLE+=("$lang")
    else
        echo -e "${YELLOW}SKIP${NC}"
    fi
done
echo ""

if [ ${#AVAILABLE[@]} -eq 0 ]; then
    echo -e "${RED}No harnesses available. Aborting.${NC}"
    exit 1
fi

# ---------------------------------------------------------------------------
# Collect fixtures
# ---------------------------------------------------------------------------
FIXTURES=()
while IFS= read -r -d '' f; do
    FIXTURES+=("$f")
done < <(find "$FIXTURES_DIR" -name "*.jsonc" -type f -print0 | sort -z)

# ---------------------------------------------------------------------------
# Test phase: IR consistency across languages
# ---------------------------------------------------------------------------
mkdir -p "$RESULTS_DIR"

echo -e "${CYAN}=== Running cross-language tests (${#FIXTURES[@]} fixtures x ${#AVAILABLE[@]} languages) ===${NC}"
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    OUTPUTS=()
    fixture_ok=1

    for lang in "${AVAILABLE[@]}"; do
        output=$("${lang}_run" "$fixture" 2>/dev/null) || true

        if [ -z "$output" ]; then
            printf " ${RED}FAIL${NC}  "
            fixture_ok=0
        else
            OUTPUTS+=("$output")
            printf " ${GREEN}OK${NC}    "
        fi
    done

    # Compare normalized outputs across all successful languages
    if [ "$fixture_ok" -eq 1 ] && [ ${#OUTPUTS[@]} -gt 0 ]; then
        ref_norm=$(echo "${OUTPUTS[0]}" | normalize) || true
        all_match=1

        for ((i=1; i<${#OUTPUTS[@]}; i++)); do
            norm=$(echo "${OUTPUTS[$i]}" | normalize) || true
            if [ "$norm" != "$ref_norm" ]; then
                all_match=0
                break
            fi
        done

        if [ "$all_match" -eq 1 ]; then
            printf " ${GREEN}MATCH${NC}"
            PASS=$((PASS + 1))
        else
            printf " ${RED}DIFF${NC}"
            FAIL=$((FAIL + 1))
            # Save diff details with unified diff
            diff_file="$RESULTS_DIR/${rel//\//_}.diff"
            {
                echo "=== $rel ==="
                echo ""
                # Generate per-language normalized outputs and diff
                for ((i=0; i<${#AVAILABLE[@]}; i++)); do
                    lang="${AVAILABLE[$i]}"
                    norm_out=$(echo "${OUTPUTS[$i]}" | normalize) || true
                    echo "--- $lang (normalized) ---"
                    echo "$norm_out"
                    echo ""
                    if [ "$i" -gt 0 ] && [ "$norm_out" != "$ref_norm" ]; then
                        echo "--- diff: go vs $lang ---"
                        diff -u <(echo "$ref_norm") <(echo "$norm_out") 2>/dev/null || true
                        echo ""
                    fi
                done
                echo "=== raw outputs ==="
                for ((i=0; i<${#AVAILABLE[@]}; i++)); do
                    echo ""
                    echo "--- ${AVAILABLE[$i]} (raw) ---"
                    echo "${OUTPUTS[$i]}"
                done
            } > "$diff_file"
        fi
    else
        printf " ${RED}FAIL${NC}"
        FAIL=$((FAIL + 1))
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Reversibility phase: parse JSONC → re-print → re-parse → compare JSON
# ---------------------------------------------------------------------------
echo ""
echo -e "${CYAN}=== Reversibility tests ===${NC}"
echo ""

REV_PASS=0
REV_FAIL=0

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    # Step 1: Parse original JSONC with Go (reference) → output1
    output1=$(go_run "$fixture" 2>/dev/null) || true
    if [ -z "$output1" ]; then
        printf " ${RED}FAIL (parse)${NC}\n"
        REV_FAIL=$((REV_FAIL + 1))
        continue
    fi

    # Step 2: Write output1 to temp file and re-parse → output2
    tmpfile=$(mktemp /tmp/mm_rev_test.XXXXXX)
    echo "$output1" > "$tmpfile"
    output2=$(go_run "$tmpfile" 2>/dev/null) || true
    rm -f "$tmpfile"

    if [ -z "$output2" ]; then
        printf " ${RED}FAIL (re-parse)${NC}\n"
        REV_FAIL=$((REV_FAIL + 1))
        continue
    fi

    # Step 3: Normalize and compare
    norm1=$(echo "$output1" | normalize) || true
    norm2=$(echo "$output2" | normalize) || true

    if [ "$norm1" = "$norm2" ]; then
        printf " ${GREEN}OK${NC}\n"
        REV_PASS=$((REV_PASS + 1))
    else
        printf " ${RED}DIFF${NC}\n"
        REV_FAIL=$((REV_FAIL + 1))
        rev_file="$RESULTS_DIR/${rel//\//_}.rev_diff"
        {
            echo "=== $rel reversibility failure ==="
            echo ""
            echo "--- round 1 (normalized) ---"
            echo "$norm1"
            echo ""
            echo "--- round 2 (normalized) ---"
            echo "$norm2"
            echo ""
            echo "--- diff ---"
            diff -u <(echo "$norm1") <(echo "$norm2") 2>/dev/null || true
            echo ""
            echo "--- round 1 (raw) ---"
            echo "$output1"
            echo ""
            echo "--- round 2 (raw) ---"
            echo "$output2"
        } > "$rev_file"
    fi
done

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo -e "${CYAN}=== Cross-language Summary ===${NC}"
echo -e "Fixtures:  ${#FIXTURES[@]}"
echo -e "Languages: ${#AVAILABLE[@]} (${AVAILABLE[*]})"
echo -e "  ${GREEN}PASS (all languages match): $PASS${NC}"
echo -e "  ${RED}FAIL (mismatch or error):    $FAIL${NC}"
echo ""
echo -e "${CYAN}=== Reversibility Summary ===${NC}"
echo -e "Fixtures: ${#FIXTURES[@]}"
echo -e "  ${GREEN}PASS (round-trip stable): $REV_PASS${NC}"
echo -e "  ${RED}FAIL (round-trip changed): $REV_FAIL${NC}"
echo ""

TOTAL_FAIL=$((FAIL + REV_FAIL))

if [ "$TOTAL_FAIL" -gt 0 ]; then
    echo -e "${YELLOW}Differences in: $RESULTS_DIR/${NC}"
    find "$RESULTS_DIR" -type f -exec echo "  {}" \;
fi

exit $TOTAL_FAIL