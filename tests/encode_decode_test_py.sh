#!/usr/bin/env bash
# Convenience wrapper: compare Python against Go reference.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$SCRIPT_DIR/compare_with_go.sh" py