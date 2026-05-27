#!/usr/bin/env bash
# Cross-language EncodeFromJsonc / DecodeToJsonc test runner
# Tests that all language implementations produce identical wire-format bytes
# and that the round-trip (encode -> decode) restores the original JSONC.
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

ENC_PASS=0
ENC_FAIL=0
DEC_PASS=0
DEC_FAIL=0

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
    print(s.strip(), end='')
" 2>/dev/null
}

# ---------------------------------------------------------------------------
# Language harnesses: each function builds (if needed), then runs.
# Exit 0 if harness is ready; non-zero to skip this language.
# ---------------------------------------------------------------------------

# --- Go ---
go_build() { return 0; }
go_run() { go run "$SCRIPT_DIR/harness/go/harness.go" "$@"; }
go_encode() { go_run --encode "$1"; }
go_decode() { printf '%s' "$1" | go_run --decode; }

# --- Python ---
py_build() { return 0; }
py_run() { python3 "$SCRIPT_DIR/harness/python/harness.py" "$@"; }
py_encode() { py_run --encode "$1"; }
py_decode() { printf '%s' "$1" | py_run --decode; }

# --- PHP ---
php_build() {
    [ -f "$PROJECT_DIR/mm-php/vendor/autoload.php" ] || return 1
    return 0
}
php_run() { php "$SCRIPT_DIR/harness/php/harness.php" "$@"; }
php_encode() { php_run --encode "$1"; }
php_decode() { printf '%s' "$1" | php_run --decode; }

# --- TypeScript ---
ts_build() {
    cd "$PROJECT_DIR/mm-ts" && npm run build --silent 2>/dev/null
}
ts_run() { node "$SCRIPT_DIR/harness/typescript/harness.cjs" "$@"; }
ts_encode() { ts_run --encode "$1"; }
ts_decode() { printf '%s' "$1" | ts_run --decode; }

# --- Rust ---
rs_build() {
    cd "$SCRIPT_DIR/harness/rust" && cargo build --quiet 2>/dev/null
}
rs_run() { "$SCRIPT_DIR/harness/rust/target/debug/mm-harness-rs" "$@"; }
rs_encode() { rs_run --encode "$1"; }
rs_decode() { printf '%s' "$1" | rs_run --decode; }

# --- C ---
c_build() {
    mkdir -p "$SCRIPT_DIR/harness/c/build"
    cd "$SCRIPT_DIR/harness/c/build" && cmake .. -DCMAKE_BUILD_TYPE=Release >/dev/null 2>&1 && make -j$(sysctl -n hw.ncpu 2>/dev/null || echo 4) >/dev/null 2>&1
}
c_run() { "$SCRIPT_DIR/harness/c/build/mm_harness_c" "$@"; }
c_encode() { c_run --encode "$1"; }
c_decode() { printf '%s' "$1" | c_run --decode; }

# --- C++ ---
cpp_build() {
    mkdir -p "$SCRIPT_DIR/harness/cpp/build"
    g++ -std=c++17 -I "$PROJECT_DIR/mm-cpp/src" -o "$SCRIPT_DIR/harness/cpp/build/mm_harness_cpp" "$SCRIPT_DIR/harness/cpp/harness.cpp" "$PROJECT_DIR/mm-cpp/src/jsonc/scanner.cpp" 2>/dev/null
}
cpp_run() { "$SCRIPT_DIR/harness/cpp/build/mm_harness_cpp" "$@"; }
cpp_encode() { cpp_run --encode "$1"; }
cpp_decode() { printf '%s' "$1" | cpp_run --decode; }

# --- C# ---
cs_build() {
    dotnet build "$SCRIPT_DIR/harness/csharp/harness.csproj" --nologo -v q >/dev/null 2>&1
}
cs_run() { dotnet run --project "$SCRIPT_DIR/harness/csharp/harness.csproj" --no-build -- "$@"; }
cs_encode() { cs_run --encode "$1"; }
cs_decode() { printf '%s' "$1" | cs_run --decode; }

# --- Kotlin ---
kt_build() {
    cd "$PROJECT_DIR/mm-kt" && mvn compile -q -DskipTests 2>/dev/null
}
kt_run() {
    local tmpdir=$(mktemp -d)
    local trove4j="$HOME/.m2/repository/org/jetbrains/intellij/deps/trove4j/1.0.20221201/trove4j-1.0.20221201.jar"
    local annotations="$HOME/.m2/repository/org/jetbrains/annotations/13.0/annotations-13.0.jar"
    local kotlin_compiler="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-compiler/1.9.22/kotlin-compiler-1.9.22.jar"
    local kotlin_stdlib="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar"
    local kotlin_reflect="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar"
    if [ ! -f "$kotlin_compiler" ] || [ ! -f "$kotlin_stdlib" ] || [ ! -f "$trove4j" ]; then
        rm -rf "$tmpdir"
        return 1
    fi
    local compiler_cp="$kotlin_compiler:$kotlin_stdlib:$kotlin_reflect:$annotations:$trove4j"
    local mm_classes="$PROJECT_DIR/mm-kt/target/classes"
    java -cp "$compiler_cp" org.jetbrains.kotlin.cli.jvm.K2JVMCompiler \
        "$SCRIPT_DIR/harness/kotlin/harness.kt" \
        -cp "$mm_classes:$kotlin_stdlib:$kotlin_reflect" \
        -d "$tmpdir/harness.jar" \
        -no-stdlib -no-reflect 2>/dev/null && \
    java -cp "$tmpdir/harness.jar:$mm_classes:$kotlin_stdlib" HarnessKt "$@"
    local rc=$?
    rm -rf "$tmpdir"
    return $rc
}
kt_encode() { kt_run --encode "$1"; }
kt_decode() { printf '%s' "$1" | kt_run --decode; }

# --- Swift ---
sw_build() {
    cd "$SCRIPT_DIR/harness/swift" && swift build --quiet 2>/dev/null
}
sw_run() { swift run --package-path "$SCRIPT_DIR/harness/swift" --skip-build mm-harness-swift "$@" 2>/dev/null; }
sw_encode() { sw_run --encode "$1"; }
sw_decode() { printf '%s' "$1" | sw_run --decode; }

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
    if ${lang}_build; then
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

echo -e "Available languages (${#AVAILABLE[@]}): ${AVAILABLE[*]}"
echo ""

# ---------------------------------------------------------------------------
# Collect fixtures
# ---------------------------------------------------------------------------
FIXTURES=()
while IFS= read -r -d '' f; do
    FIXTURES+=("$f")
done < <(find "$FIXTURES_DIR" -name "*.jsonc" -type f -print0 | sort -z)

# ===================================================================
# Test 1: EncodeFromJsonc - compare wire-format bytes across languages
# ===================================================================
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Test 1: EncodeFromJsonc - bytes一致性對比${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

printf "%-45s" ""
for _lang in "${AVAILABLE[@]}"; do
    printf " ${CYAN}%-6s${NC}" "$_lang"
done
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    HEX_OUTPUTS=()
    fixture_ok=1

    for lang in "${AVAILABLE[@]}"; do
        hex=$(${lang}_encode "$fixture" 2>/dev/null) || true

        if [ -z "$hex" ]; then
            printf " ${RED}FAIL${NC}  "
            fixture_ok=0
        else
            HEX_OUTPUTS+=("$hex")
            printf " ${GREEN}OK${NC}    "
        fi
    done

    # Compare hex outputs across all successful languages
    if [ "$fixture_ok" -eq 1 ] && [ ${#HEX_OUTPUTS[@]} -gt 0 ]; then
        ref_hex="${HEX_OUTPUTS[0]}"
        all_match=1

        for ((i=1; i<${#HEX_OUTPUTS[@]}; i++)); do
            if [ "${HEX_OUTPUTS[$i]}" != "$ref_hex" ]; then
                all_match=0
                break
            fi
        done

        if [ "$all_match" -eq 1 ]; then
            printf " ${GREEN}MATCH${NC}"
            ENC_PASS=$((ENC_PASS + 1))
        else
            printf " ${RED}DIFF${NC}"
            ENC_FAIL=$((ENC_FAIL + 1))
            diff_file="$RESULTS_DIR/${rel//\//_}.encode_diff"
            {
                echo "=== $rel EncodeFromJsonc bytes mismatch ==="
                echo ""
                for ((i=0; i<${#AVAILABLE[@]}; i++)); do
                    lang="${AVAILABLE[$i]}"
                    echo "--- $lang (hex) ---"
                    echo "${HEX_OUTPUTS[$i]}"
                    echo ""
                done
            } > "$diff_file"
        fi
    else
        printf " ${RED}FAIL${NC}"
        ENC_FAIL=$((ENC_FAIL + 1))
    fi

    echo ""
done

# ===================================================================
# Test 2: DecodeToJsonc round-trip - encode with Go, decode with all
# ===================================================================
echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Test 2: DecodeToJsonc round-trip - 還原一致性對比${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

printf "%-45s" ""
for _lang in "${AVAILABLE[@]}"; do
    printf " ${CYAN}%-6s${NC}" "$_lang"
done
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    # Step 1: Encode with Go (reference) to get hex bytes
    ref_hex=$(go_encode "$fixture" 2>/dev/null) || true
    if [ -z "$ref_hex" ]; then
        printf " ${RED}NO-ENC${NC}\n"
        DEC_FAIL=$((DEC_FAIL + 1))
        continue
    fi

    # Step 2: Decode hex with each language and collect JSONC output
    JSONC_OUTPUTS=()
    fixture_ok=1

    for lang in "${AVAILABLE[@]}"; do
        output=$(${lang}_decode "$ref_hex" 2>/dev/null) || true

        if [ -z "$output" ]; then
            printf " ${RED}FAIL${NC}  "
            fixture_ok=0
        else
            JSONC_OUTPUTS+=("$output")
            printf " ${GREEN}OK${NC}    "
        fi
    done

    # Compare normalized JSONC outputs
    if [ "$fixture_ok" -eq 1 ] && [ ${#JSONC_OUTPUTS[@]} -gt 0 ]; then
        ref_norm=$(echo "${JSONC_OUTPUTS[0]}" | normalize) || true
        all_match=1

        for ((i=1; i<${#JSONC_OUTPUTS[@]}; i++)); do
            norm=$(echo "${JSONC_OUTPUTS[$i]}" | normalize) || true
            if [ "$norm" != "$ref_norm" ]; then
                all_match=0
                break
            fi
        done

        if [ "$all_match" -eq 1 ]; then
            printf " ${GREEN}MATCH${NC}"
            DEC_PASS=$((DEC_PASS + 1))
        else
            printf " ${RED}DIFF${NC}"
            DEC_FAIL=$((DEC_FAIL + 1))
            diff_file="$RESULTS_DIR/${rel//\//_}.decode_diff"
            {
                echo "=== $rel DecodeToJsonc round-trip mismatch ==="
                echo "--- reference hex (from Go encode) ---"
                echo "$ref_hex"
                echo ""
                for ((i=0; i<${#AVAILABLE[@]}; i++)); do
                    lang="${AVAILABLE[$i]}"
                    norm_out=$(echo "${JSONC_OUTPUTS[$i]}" | normalize) || true
                    echo "--- $lang (normalized) ---"
                    echo "$norm_out"
                    echo ""
                    if [ "$i" -gt 0 ] && [ "$norm_out" != "$ref_norm" ]; then
                        echo "--- diff: ${AVAILABLE[0]} vs $lang ---"
                        diff -u <(echo "$ref_norm") <(echo "$norm_out") 2>/dev/null || true
                        echo ""
                    fi
                done
            } > "$diff_file"
        fi
    else
        printf " ${RED}FAIL${NC}"
        DEC_FAIL=$((DEC_FAIL + 1))
    fi

    echo ""
done

# ===================================================================
# Summary
# ===================================================================
echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Summary${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""
echo -e "Fixtures:  ${#FIXTURES[@]}"
echo -e "Languages: ${#AVAILABLE[@]} (${AVAILABLE[*]})"
echo ""
echo -e "${CYAN}--- EncodeFromJsonc (bytes一致性) ---${NC}"
echo -e "  ${GREEN}PASS (all languages produce identical bytes): $ENC_PASS${NC}"
echo -e "  ${RED}FAIL (bytes mismatch or error):               $ENC_FAIL${NC}"
echo ""
echo -e "${CYAN}--- DecodeToJsonc (還原一致性) ---${NC}"
echo -e "  ${GREEN}PASS (all languages round-trip match):       $DEC_PASS${NC}"
echo -e "  ${RED}FAIL (round-trip mismatch or error):           $DEC_FAIL${NC}"
echo ""

TOTAL_FAIL=$((ENC_FAIL + DEC_FAIL))

if [ "$TOTAL_FAIL" -gt 0 ]; then
    echo -e "${YELLOW}Differences in: $RESULTS_DIR/${NC}"
    find "$RESULTS_DIR" -type f \( -name "*.encode_diff" -o -name "*.decode_diff" \) -exec echo "  {}" \;
fi

exit $TOTAL_FAIL