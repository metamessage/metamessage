#!/usr/bin/env bash
# Compare a single language implementation against Go reference.
# Usage: ./compare_with_go.sh <lang_key>
#   lang_key: go, py, php, ts, rs, c, cpp, cs, kt, sw
#
# Tests:
#   Test 1 (encode): Encode fixture with both Go and <lang>, compare hex bytes
#   Test 2 (decode): Encode with Go, decode with <lang>, compare normalized JSONC

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

TARGET_LANG="$1"
if [ -z "$TARGET_LANG" ]; then
    echo "Usage: $0 <lang_key>"
    echo "  lang_key: go, py, php, ts, rs, c, cpp, cs, kt, sw"
    exit 1
fi

# Validate language key
VALID_LANGS="go py php ts rs c cpp cs kt sw"
found=0
for l in $VALID_LANGS; do
    if [ "$l" = "$TARGET_LANG" ]; then
        found=1
        break
    fi
done
if [ "$found" -eq 0 ]; then
    echo -e "${RED}Invalid language key: $TARGET_LANG${NC}"
    echo "Valid keys: $VALID_LANGS"
    exit 1
fi

# ---------------------------------------------------------------------------
# Normalize JSONC
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
# Language harnesses
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

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

ENC_PASS=0
ENC_FAIL=0
DEC_PASS=0
DEC_FAIL=0

# Map lang key to display name
case "$TARGET_LANG" in
    go)  LANG_NAME="Go";;
    py)  LANG_NAME="Python";;
    php) LANG_NAME="PHP";;
    ts)  LANG_NAME="TypeScript";;
    rs)  LANG_NAME="Rust";;
    c)   LANG_NAME="C";;
    cpp) LANG_NAME="C++";;
    cs)  LANG_NAME="C#";;
    kt)  LANG_NAME="Kotlin";;
    sw)  LANG_NAME="Swift";;
esac

echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  ${TARGET_LANG} ($LANG_NAME) vs Go comparison${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

# Build harnesses
echo -e "${CYAN}--- Building harnesses ---${NC}"

go_build
printf "  %-12s " "go"
echo -e "${GREEN}OK${NC}"

printf "  %-12s " "$TARGET_LANG"
if ${TARGET_LANG}_build; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}BUILD FAILED${NC}"
    exit 1
fi
echo ""

# Collect fixtures
FIXTURES=()
while IFS= read -r -d '' f; do
    FIXTURES+=("$f")
done < <(find "$FIXTURES_DIR" -name "*.jsonc" -type f -print0 | sort -z)

mkdir -p "$RESULTS_DIR"

# ===================================================================
# Test 1: Encode - compare wire-format bytes (Go vs target)
# ===================================================================
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Test 1: EncodeFromJsonc - bytes對比 (Go vs ${TARGET_LANG})${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    go_hex=$(go_encode "$fixture" 2>/dev/null) || true
    target_hex=$(${TARGET_LANG}_encode "$fixture" 2>/dev/null) || true

    if [ -z "$go_hex" ]; then
        echo -e " ${RED}GO-ERR${NC}"
        ENC_FAIL=$((ENC_FAIL + 1))
        continue
    fi

    if [ -z "$target_hex" ]; then
        echo -e " ${RED}${TARGET_LANG}-ERR${NC}"
        ENC_FAIL=$((ENC_FAIL + 1))
        continue
    fi

    if [ "$go_hex" = "$target_hex" ]; then
        echo -e " ${GREEN}MATCH${NC}"
        ENC_PASS=$((ENC_PASS + 1))
    else
        echo -e " ${RED}DIFF${NC}"
        ENC_FAIL=$((ENC_FAIL + 1))

        diff_file="$RESULTS_DIR/${rel//\//_}.${TARGET_LANG}_vs_go.encode_diff"
        {
            echo "=== $rel Encode bytes mismatch: Go vs $TARGET_LANG ==="
            echo ""
            echo "--- Go (hex, ${#go_hex} chars) ---"
            echo "$go_hex"
            echo ""
            echo "--- $TARGET_LANG (hex, ${#target_hex} chars) ---"
            echo "$target_hex"
            echo ""
            echo "--- char-by-char diff ---"
            # Show side-by-side with markers for differences
            len_go=${#go_hex}
            len_target=${#target_hex}
            max_len=$(( len_go > len_target ? len_go : len_target ))
            for ((i=0; i<max_len; i+=2)); do
                g_char="${go_hex:$i:2}"
                t_char="${target_hex:$i:2}"
                if [ "$g_char" != "$t_char" ]; then
                    byte_pos=$((i/2))
                    if [ -n "$g_char" ] && [ -n "$t_char" ]; then
                        printf "  byte[%3d] Go=0x%s  %s=0x%s\n" "$byte_pos" "$g_char" "$TARGET_LANG" "$t_char"
                    elif [ -z "$g_char" ]; then
                        printf "  byte[%3d] Go=(end) %s=0x%s\n" "$byte_pos" "$TARGET_LANG" "$t_char"
                    else
                        printf "  byte[%3d] Go=0x%s  %s=(end)\n" "$byte_pos" "$g_char" "$TARGET_LANG" "$t_char"
                    fi
                fi
            done
        } > "$diff_file"
        echo -e "  ${YELLOW}Details: $diff_file${NC}"
    fi
done

# ===================================================================
# Test 2: Decode round-trip - encode with Go, decode with target
# ===================================================================
echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Test 2: DecodeToJsonc round-trip (Go encode → ${TARGET_LANG} decode)${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    # Encode with Go
    ref_hex=$(go_encode "$fixture" 2>/dev/null) || true
    if [ -z "$ref_hex" ]; then
        echo -e " ${RED}GO-ENC-ERR${NC}"
        DEC_FAIL=$((DEC_FAIL + 1))
        continue
    fi

    # Decode with target language
    target_output=$(${TARGET_LANG}_decode "$ref_hex" 2>/dev/null) || true
    if [ -z "$target_output" ]; then
        echo -e " ${RED}${TARGET_LANG}-DEC-ERR${NC}"
        DEC_FAIL=$((DEC_FAIL + 1))
        continue
    fi

    # Also decode with Go for reference JSONC
    go_output=$(go_decode "$ref_hex" 2>/dev/null) || true

    ref_norm=$(echo "$go_output" | normalize) || true
    target_norm=$(echo "$target_output" | normalize) || true

    if [ "$ref_norm" = "$target_norm" ]; then
        echo -e " ${GREEN}MATCH${NC}"
        DEC_PASS=$((DEC_PASS + 1))
    else
        echo -e " ${RED}DIFF${NC}"
        DEC_FAIL=$((DEC_FAIL + 1))

        diff_file="$RESULTS_DIR/${rel//\//_}.${TARGET_LANG}_vs_go.decode_diff"
        {
            echo "=== $rel Decode round-trip mismatch (Go encode → decode) ==="
            echo ""
            echo "--- Go decode (normalized) ---"
            echo "$ref_norm"
            echo ""
            echo "--- $TARGET_LANG decode (normalized) ---"
            echo "$target_norm"
            echo ""
            echo "--- unified diff ---"
            diff -u <(echo "$ref_norm") <(echo "$target_norm") 2>/dev/null || true
            echo ""
            echo "--- Go decode (raw) ---"
            echo "$go_output"
            echo ""
            echo "--- $TARGET_LANG decode (raw) ---"
            echo "$target_output"
        } > "$diff_file"
        echo -e "  ${YELLOW}Details: $diff_file${NC}"
    fi
done

# ===================================================================
# Summary
# ===================================================================
echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Summary: ${TARGET_LANG} vs Go${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""
echo -e "Fixtures:  ${#FIXTURES[@]}"
echo ""
echo -e "${CYAN}--- Encode (bytes一致性) ---${NC}"
echo -e "  ${GREEN}PASS (bytes match Go): $ENC_PASS${NC}"
echo -e "  ${RED}FAIL (bytes differ):     $ENC_FAIL${NC}"
echo ""
echo -e "${CYAN}--- Decode (還原一致性) ---${NC}"
echo -e "  ${GREEN}PASS (JSONC matches Go): $DEC_PASS${NC}"
echo -e "  ${RED}FAIL (JSONC differs):     $DEC_FAIL${NC}"
echo ""

TOTAL_FAIL=$((ENC_FAIL + DEC_FAIL))

if [ "$TOTAL_FAIL" -gt 0 ]; then
    echo -e "${YELLOW}Diff files in: $RESULTS_DIR/${NC}"
    find "$RESULTS_DIR" -type f \( -name "*.${TARGET_LANG}_vs_go.*" \) -exec echo "  {}" \;
fi

exit $TOTAL_FAIL