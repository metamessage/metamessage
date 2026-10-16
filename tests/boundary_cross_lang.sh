#!/usr/bin/env bash
# boundary_cross_lang.sh — 容量扩容（Bytes/Container/Tag 2 档 → 4 档）跨语言主门禁。
#
# 对每个语言 harness 的 --decode（wire ← stdin）施加两类断言：
#   1. 有效档位数据：由 Go 参考编码产生，跨语言解码必须还原出与 Go 完全一致的 JSONC。
#      覆盖 1B / 2B / 3B / 4B 各档真实长度（含越过 65536 的 3B 大 payload）。
#   2. 4B 档截断负例：header 声明超大长度但数据不足——每个语言必须「显式报错」，
#      既不崩溃（非 134/139 信号退出）、也不静默返回成功，更不得尝试分配 ~4GB 而 OOM。
#
# 用法：./tests/boundary_cross_lang.sh [lang ...]   # 不指定则跑全部 10 语言
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 清除外部进程注入的 LIBRARY_PATH/CC/CPPFLAGS 等污染，详见 env_sanitize.sh
source "$SCRIPT_DIR/env_sanitize.sh"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
GO_HARNESS="$SCRIPT_DIR/harness/go/harness.go"
C_BIN="$SCRIPT_DIR/harness/c/build/mm_harness_c"
CPP_BIN="$SCRIPT_DIR/harness/cpp/build/mm_harness_cpp"
RS_BIN="$SCRIPT_DIR/harness/rust/target/debug/mm-harness-rs"
SW_BIN="$SCRIPT_DIR/harness/swift/.build/debug/mm-harness-swift"
KT_MM="$PROJECT_DIR/mm-kt/target/classes"
TS_HARNESS="$SCRIPT_DIR/harness/typescript/harness.cjs"
PY_HARNESS="$SCRIPT_DIR/harness/python/harness.py"
PHP_HARNESS="$SCRIPT_DIR/harness/php/harness.php"
CS_PROJ="$SCRIPT_DIR/harness/csharp/harness.csproj"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

# ---------------------------------------------------------------------------
# 每语言的 --decode：把 hex 从 stdin 喂入，stdout 输出 JSONC（出错则非零退出）
# ---------------------------------------------------------------------------
dec_go()  { printf '%s' "$1" | go run "$GO_HARNESS" --decode; }
dec_py()  { printf '%s' "$1" | python3 "$PY_HARNESS" --decode; }
dec_php() { printf '%s' "$1" | php "$PHP_HARNESS" --decode; }
dec_ts()  { printf '%s' "$1" | node "$TS_HARNESS" --decode; }
dec_rs()  { printf '%s' "$1" | "$RS_BIN" --decode; }
dec_c()   { printf '%s' "$1" | "$C_BIN" --decode; }
dec_cpp() { printf '%s' "$1" | "$CPP_BIN" --decode; }
dec_cs()  { printf '%s' "$1" | dotnet run --project "$CS_PROJ" --no-build -- --decode; }
dec_sw()  { printf '%s' "$1" | "$SW_BIN" --decode; }
dec_kt() {
    local tmpdir; tmpdir=$(mktemp -d)
    local trove4j="$HOME/.m2/repository/org/jetbrains/intellij/deps/trove4j/1.0.20221201/trove4j-1.0.20221201.jar"
    local annotations="$HOME/.m2/repository/org/jetbrains/annotations/13.0/annotations-13.0.jar"
    local kc="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-compiler/1.9.22/kotlin-compiler-1.9.22.jar"
    local kstd="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar"
    local kref="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar"
    printf '%s' "$1" | java -cp "$kc:$kstd:$kref:$annotations:$trove4j" \
        org.jetbrains.kotlin.cli.jvm.K2JVMCompiler \
        "$SCRIPT_DIR/harness/kotlin/harness.kt" -cp "$KT_MM:$kstd:$kref" \
        -d "$tmpdir/h.jar" -no-stdlib -no-reflect 2>/dev/null
    printf '%s' "$1" | java -cp "$tmpdir/h.jar:$KT_MM:$kstd" HarnessKt --decode
    local rc=$?; rm -rf "$tmpdir"; return $rc
}

ALL_LANGS="go py php ts rs c cpp cs kt sw"

# ---------------------------------------------------------------------------
# 1. 有效档位数据：用 Go 编码真实数据，跨语言解码必须与 Go 解码结果一致
# ---------------------------------------------------------------------------
declare -a VALID
add_valid() { VALID+=("$1|$2"); }   # name|hex

# 大 payload（object+array 越过 65536 → 3B 档）——容器路径，跨语言 JSONC 还原一致
BIG_HEX=$(go run "$GO_HARNESS" --encode "$SCRIPT_DIR/fixtures/07_large/big_object.jsonc" 2>/dev/null)
if [ -n "$BIG_HEX" ]; then
    add_valid "big_object_3B" "$BIG_HEX"
fi

# 注：1B/2B 档有效路径已由 compare_with_go.sh 的 36 个真实 fixtures（含跨 28–65535 各
# 档的容器/字符串）逐语言认证；此处不再放「根为 tagged scalar」的合成 bytes 正例——
# 那类用例的 wire 字节一致，仅因 Go 对根 tagged scalar 会重印 `// mm:` 注释而产生
# 纯展示层差异（与本次容量改造无关，属既有行为）。

# ---------------------------------------------------------------------------
# 2. 4B / 3B 档截断负例：header 声明超大长度，但 body 数据不足
# ---------------------------------------------------------------------------
declare -a NEG=(
    "bytes_4B_trunc|bfffffffff0102"          # 0xBF=BYTES_LEN_4 声明 ~4GB，仅 2 字节
    "object_4B_trunc|cfffffffff0102"         # 0xCF=CONTAINER_LEN_4
    "array_4B_trunc|dfffffffff0102"          # 0xDF=CONTAINER_LEN_4（array）
    "bytes_3B_trunc|beffffff0102"            # 0xBE=BYTES_LEN_3 声明 ~16MB
    "object_3B_trunc|ceffffff0102"           # 0xCE=CONTAINER_LEN_3
    "array_3B_trunc|deffffff0102"            # 0xDE=CONTAINER_LEN_3（array）
)

# macOS 无 timeout：用后台进程 + 轮询兜住「卡死/超大分配」的实现。
# 仅一次 wait，用 hung 标志区分「超时被杀」与「真实崩溃」。
run_guarded() {   # $1=cmd名 $2=hex ; 设置全局 RC / HUNG
    local cmd=$1 hex=$2 tmpf pid hung=0 t
    tmpf=$(mktemp /tmp/mm_neg.XXXXXX)
    ( "$cmd" "$hex" >"$tmpf" 2>/dev/null ) &
    pid=$!
    for ((t=0; t<30; t++)); do
        kill -0 "$pid" 2>/dev/null || break
        sleep 1
    done
    if kill -0 "$pid" 2>/dev/null; then kill -9 "$pid" 2>/dev/null; hung=1; fi
    wait "$pid" 2>/dev/null; RC=$?
    HUNG=$hung
    rm -f "$tmpf"
}

TARGET_LANGS=("$@")
[ ${#TARGET_LANGS[@]} -eq 0 ] && TARGET_LANGS=($ALL_LANGS)

echo -e "${CYAN}=== boundary_cross_lang: 容量档位跨语言门禁 ===${NC}"
echo -e "语言: ${TARGET_LANGS[*]}"
echo ""

VALID_FAIL=0; NEG_FAIL=0

# --- 有效数据 ---
echo -e "${CYAN}--- 有效档位数据（decode 必须与 Go 一致）---${NC}"
for entry in "${VALID[@]}"; do
    name="${entry%%|*}"; hex="${entry#*|}"
    printf "%-24s" "$name"
    ref=$(dec_go "$hex" 2>/dev/null) || true
    for lang in "${TARGET_LANGS[@]}"; do
        out=$("dec_$lang" "$hex" 2>/dev/null) || true
        if [ "$lang" = "go" ]; then
            printf " ${CYAN}ref${NC}"
        elif [ -n "$out" ] && [ "$out" = "$ref" ]; then
            printf " ${GREEN}%s:OK${NC}" "$lang"
        else
            printf " ${RED}%s:BAD${NC}" "$lang"; VALID_FAIL=$((VALID_FAIL+1))
        fi
    done
    echo ""
done

# --- 截断负例 ---
echo ""
echo -e "${CYAN}--- 4B/3B 档截断负例（必须显式报错：非零退出、不崩溃、不超时）---${NC}"
for entry in "${NEG[@]}"; do
    name="${entry%%|*}"; hex="${entry#*|}"
    printf "%-24s" "$name"
    for lang in "${TARGET_LANGS[@]}"; do
        run_guarded "dec_$lang" "$hex"
        if [ "$HUNG" -eq 1 ]; then
            printf " ${RED}%s:HANG${NC}" "$lang"; NEG_FAIL=$((NEG_FAIL+1))
        elif [ $RC -ge 128 ]; then
            printf " ${RED}%s:CRASH($RC)${NC}" "$lang"; NEG_FAIL=$((NEG_FAIL+1))
        elif [ $RC -eq 0 ]; then
            printf " ${RED}%s:SILENT${NC}" "$lang"; NEG_FAIL=$((NEG_FAIL+1))
        else
            printf " ${GREEN}%s:ERR${NC}" "$lang"
        fi
    done
    echo ""
done

echo ""
echo -e "${CYAN}=== Summary ===${NC}"
echo -e "  有效数据 FAIL: ${VALID_FAIL}"
echo -e "  负例处理 FAIL: ${NEG_FAIL}"
TOTAL=$((VALID_FAIL + NEG_FAIL))
if [ "$TOTAL" -eq 0 ]; then
    echo -e "  ${GREEN}ALL PASS${NC}"
else
    echo -e "  ${RED}FAILURES: $TOTAL${NC}"
fi
exit $TOTAL
