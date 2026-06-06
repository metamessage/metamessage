#!/usr/bin/env bash
# Convenience wrapper: compare Kotlin against Go reference.
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$SCRIPT_DIR/compare_with_go.sh" kt