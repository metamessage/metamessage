#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
FIXTURES_DIR="$SCRIPT_DIR/fixtures"
SWIFT_BIN="/tmp/mm-harness-swift"

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

PASS=0
FAIL=0
DIFF=0

for f in $(find "$FIXTURES_DIR" -name "*.jsonc" -type f | sort); do
    name=$(basename "$f")
    REF_HEX=$(go run "$SCRIPT_DIR/harness/go/harness.go" --encode "$f" 2>/dev/null)

    SW_OUT=$(echo "$REF_HEX" | "$SWIFT_BIN" --decode 2>&1) || true
    GO_OUT=$(printf '%s' "$REF_HEX" | go run "$SCRIPT_DIR/harness/go/harness.go" --decode 2>/dev/null) || true

    if [ -z "$SW_OUT" ]; then
        echo "FAIL: $name (empty output)"
        FAIL=$((FAIL + 1))
        continue
    fi

    SW_NORM=$(echo "$SW_OUT" | normalize)
    GO_NORM=$(echo "$GO_OUT" | normalize)

    if [ "$SW_NORM" = "$GO_NORM" ]; then
        echo "MATCH: $name"
        PASS=$((PASS + 1))
    else
        echo "DIFF:  $name"
        echo "  Swift: $SW_NORM"
        echo "  Go:    $GO_NORM"
        DIFF=$((DIFF + 1))
    fi
done

echo "----------------------------------------"
echo "Results: PASS=$PASS FAIL=$FAIL DIFF=$DIFF"