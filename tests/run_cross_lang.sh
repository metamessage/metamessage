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

# Compare raw JSONC output directly. Go is the reference format.

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
    java -cp "$tmpdir/harness.jar:$mm_classes:$kotlin_stdlib" HarnessKt "$1" 2>/dev/null
    local rc=$?
    rm -rf "$tmpdir"
    return $rc
}

# --- Swift ---
sw_build() {
    [ -f "$SCRIPT_DIR/harness/swift/.build/debug/mm-harness-swift" ] && return 0
    cd "$SCRIPT_DIR/harness/swift" && swift build --quiet 2>/dev/null
}
sw_run() { "$SCRIPT_DIR/harness/swift/.build/debug/mm-harness-swift" "$1" 2>/dev/null; }

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

printf "%-45s" ""
for _lang in "${AVAILABLE[@]}"; do
    printf " ${CYAN}%-6s${NC}" "$_lang"
done
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

    # Compare raw outputs: all must match Go (first language)
    if [ "$fixture_ok" -eq 1 ] && [ ${#OUTPUTS[@]} -gt 0 ]; then
        ref_out="${OUTPUTS[0]}"
        all_match=1

        for ((i=1; i<${#OUTPUTS[@]}; i++)); do
            if [ "${OUTPUTS[$i]}" != "$ref_out" ]; then
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
            diff_file="$RESULTS_DIR/${rel//\//_}.diff"
            {
                echo "=== $rel ==="
                echo ""
                for ((i=0; i<${#AVAILABLE[@]}; i++)); do
                    lang="${AVAILABLE[$i]}"
                    echo "--- $lang ---"
                    echo "${OUTPUTS[$i]}"
                    echo ""
                    if [ "$i" -gt 0 ] && [ "${OUTPUTS[$i]}" != "$ref_out" ]; then
                        echo "--- diff: go vs $lang ---"
                        diff -u <(echo "$ref_out") <(echo "${OUTPUTS[$i]}") 2>/dev/null || true
                        echo ""
                    fi
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
# Reversibility phase 1: per-language round-trip
#   parse(print(parse(input))) == print(parse(input)) for each language
# ---------------------------------------------------------------------------
echo ""
echo -e "${CYAN}=== Reversibility: per-language round-trip ===${NC}"
echo ""

REV_PASS=0
REV_FAIL=0

printf "%-45s" ""
for _lang in "${AVAILABLE[@]}"; do
    printf " ${CYAN}%-6s${NC}" "$_lang"
done
echo ""

for fixture in "${FIXTURES[@]}"; do
    rel="${fixture#$FIXTURES_DIR/}"
    printf "%-45s" "$rel"

    rev_match=1
    for lang in "${AVAILABLE[@]}"; do
        output1=$("${lang}_run" "$fixture" 2>/dev/null) || true
        if [ -z "$output1" ]; then
            printf " ${RED}FAIL${NC}  "
            REV_FAIL=$((REV_FAIL + 1))
            rev_match=0
            continue
        fi

        tmpfile=$(mktemp /tmp/mm_rev_test.XXXXXX)
        echo "$output1" > "$tmpfile"
        output2=$("${lang}_run" "$tmpfile" 2>/dev/null) || true
        rm -f "$tmpfile"

        if [ -z "$output2" ]; then
            printf " ${RED}FAIL${NC}  "
            REV_FAIL=$((REV_FAIL + 1))
            rev_match=0
            continue
        fi

        if [ "$output1" = "$output2" ]; then
            printf " ${GREEN}OK${NC}    "
        else
            printf " ${RED}DIFF${NC}  "
            REV_FAIL=$((REV_FAIL + 1))
            rev_match=0
            rev_file="$RESULTS_DIR/${lang}_${rel//\//_}.rev_diff"
            {
                echo "=== $rel ($lang) per-language reversibility failure ==="
                echo ""
                echo "--- round 1 ---"
                echo "$output1"
                echo ""
                echo "--- round 2 ---"
                echo "$output2"
                echo ""
                echo "--- diff ---"
                diff -u <(echo "$output1") <(echo "$output2") 2>/dev/null || true
            } > "$rev_file"
        fi
    done

    if [ "$rev_match" -eq 1 ]; then
        printf " ${GREEN}MATCH${NC}"
    else
        printf " ${RED}DIFF${NC}"
    fi
    echo ""
done

# ---------------------------------------------------------------------------
# Reversibility phase 2: cross-language round-trip
#   parse_lang(print_go(input)) == print_go(input) for each non-Go language
# ---------------------------------------------------------------------------
# echo ""
# echo -e "${CYAN}=== Reversibility: cross-language round-trip (Go output → each lang) ===${NC}"
# echo ""

# XREV_PASS=0
# XREV_FAIL=0

# printf "%-45s" ""
# for _lang in "${AVAILABLE[@]}"; do
#     printf " ${CYAN}%-6s${NC}" "$_lang"
# done
# echo ""

# for fixture in "${FIXTURES[@]}"; do
#     rel="${fixture#$FIXTURES_DIR/}"
#     printf "%-45s" "$rel"

#     go_output=$(go_run "$fixture" 2>/dev/null) || true
#     if [ -z "$go_output" ]; then
#         for _lang in "${AVAILABLE[@]}"; do
#             printf " ${RED}FAIL${NC}  "
#         done
#         echo ""
#         continue
#     fi

#     for lang in "${AVAILABLE[@]}"; do
#         if [ "$lang" = "go" ]; then
#             printf " ${GREEN}ref${NC}   "
#             continue
#         fi

#         tmpfile=$(mktemp /tmp/mm_xrev_test.XXXXXX)
#         echo "$go_output" > "$tmpfile"
#         output2=$("${lang}_run" "$tmpfile" 2>/dev/null) || true
#         rm -f "$tmpfile"

#         if [ -z "$output2" ]; then
#             printf " ${RED}FAIL${NC}  "
#             XREV_FAIL=$((XREV_FAIL + 1))
#             continue
#         fi

#         if [ "$go_output" = "$output2" ]; then
#             printf " ${GREEN}OK${NC}    "
#         else
#             printf " ${RED}DIFF${NC}  "
#             XREV_FAIL=$((XREV_FAIL + 1))
#             xrev_file="$RESULTS_DIR/${lang}_${rel//\//_}.xrev_diff"
#             {
#                 echo "=== $rel ($lang) cross-language reversibility failure ==="
#                 echo ""
#                 echo "--- Go output ---"
#                 echo "$go_output"
#                 echo ""
#                 echo "--- $lang re-parse ---"
#                 echo "$output2"
#                 echo ""
#                 echo "--- diff ---"
#                 diff -u <(echo "$go_output") <(echo "$output2") 2>/dev/null || true
#             } > "$xrev_file"
#         fi
#     done
#     echo ""
# done

# ---------------------------------------------------------------------------
# Reversibility Summary (per-language)
# ---------------------------------------------------------------------------
REV_TOTAL=$(( ${#FIXTURES[@]} * ${#AVAILABLE[@]} ))
REV_PASS=$(( REV_TOTAL - REV_FAIL ))

# ---------------------------------------------------------------------------
# Cross-language Reversibility Summary
# ---------------------------------------------------------------------------
XREV_TOTAL=$(( ${#FIXTURES[@]} * (${#AVAILABLE[@]} - 1) ))
XREV_PASS=$(( XREV_TOTAL - XREV_FAIL ))

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
echo -e "${CYAN}=== Reversibility: per-language round-trip ===${NC}"
echo -e "Fixtures:  ${#FIXTURES[@]}"
echo -e "Languages: ${#AVAILABLE[@]} (${AVAILABLE[*]})"
echo -e "  ${GREEN}PASS: $REV_PASS${NC}"
echo -e "  ${RED}FAIL: $REV_FAIL${NC}"
echo ""
# echo -e "${CYAN}=== Reversibility: cross-language round-trip ===${NC}"
# echo -e "Fixtures:        ${#FIXTURES[@]}"
# echo -e "Non-Go languages: $(( ${#AVAILABLE[@]} - 1 ))"
# echo -e "  ${GREEN}PASS: $XREV_PASS${NC}"
# echo -e "  ${RED}FAIL: $XREV_FAIL${NC}"
# echo ""

TOTAL_FAIL=$((FAIL + REV_FAIL + XREV_FAIL))

if [ "$TOTAL_FAIL" -gt 0 ]; then
    echo -e "${YELLOW}Differences in: $RESULTS_DIR/${NC}"
    find "$RESULTS_DIR" -type f -exec echo "  {}" \;
fi

exit $TOTAL_FAIL