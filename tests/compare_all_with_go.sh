#!/usr/bin/env bash
# Compare all available languages against Go reference.
# This runs compare_with_go.sh for each language that can be built.
#
# Usage: ./compare_all_with_go.sh [lang1 lang2 ...]
#   If no languages specified, tests all available languages.
#
# Examples:
#   ./compare_all_with_go.sh              # Test all languages
#   ./compare_all_with_go.sh cpp py       # Test only C++ and Python
#   ./compare_all_with_go.sh cpp          # Test only C++

set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
COMPARE_SCRIPT="$SCRIPT_DIR/compare_with_go.sh"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# If specific languages provided, use those; otherwise test all
if [ $# -gt 0 ]; then
    TARGET_LANGS=("$@")
else
    TARGET_LANGS=(py php ts rs c cpp cs kt sw)
fi

OVERALL_ENC_PASS=0
OVERALL_ENC_FAIL=0
OVERALL_DEC_PASS=0
OVERALL_DEC_FAIL=0

echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Cross-Language vs Go Comparison${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

for lang in "${TARGET_LANGS[@]}"; do
    echo ""
    echo -e "${YELLOW}============================================================${NC}"
    echo -e "${YELLOW}  Testing: $lang${NC}"
    echo -e "${YELLOW}============================================================${NC}"
    echo ""

    # Run the comparison script, capture exit code
    "$COMPARE_SCRIPT" "$lang"
    rc=$?

    # Parse results from output (last few lines)
    # Look for the summary lines
done

echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Overall Summary${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""
echo -e "${YELLOW}All diff files saved to: $RESULTS_DIR/${NC}"
echo ""
echo -e "${YELLOW}To view detailed byte-level diffs:${NC}"
echo "  cat $RESULTS_DIR/<fixture>.<lang>_vs_go.encode_diff"
echo ""
echo -e "${YELLOW}To view JSONC decode diffs:${NC}"
echo "  cat $RESULTS_DIR/<fixture>.<lang>_vs_go.decode_diff"
echo ""
echo -e "${YELLOW}To test a single language:${NC}"
echo "  $COMPARE_SCRIPT <lang_key>"
exit 0