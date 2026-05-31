#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

tmpdir=$(mktemp -d)
trove4j="$HOME/.m2/repository/org/jetbrains/intellij/deps/trove4j/1.0.20221201/trove4j-1.0.20221201.jar"
annotations="$HOME/.m2/repository/org/jetbrains/annotations/13.0/annotations-13.0.jar"
kotlin_compiler="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-compiler/1.9.22/kotlin-compiler-1.9.22.jar"
kotlin_stdlib="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar"
kotlin_reflect="$HOME/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar"
compiler_cp="$kotlin_compiler:$kotlin_stdlib:$kotlin_reflect:$annotations:$trove4j"
mm_classes="$PROJECT_DIR/mm-kt/target/classes"

java -cp "$compiler_cp" org.jetbrains.kotlin.cli.jvm.K2JVMCompiler \
    "$SCRIPT_DIR/harness/kotlin/harness.kt" \
    -cp "$mm_classes:$kotlin_stdlib:$kotlin_reflect" \
    -d "$tmpdir/harness.jar" \
    -no-stdlib -no-reflect 2>&1

echo "=== child_tags.jsonc ==="
kt_hex=$(java -cp "$tmpdir/harness.jar:$mm_classes:$kotlin_stdlib" HarnessKt --encode "$SCRIPT_DIR/fixtures/03_tags/child_tags.jsonc" 2>&1)
go_hex=$(go run "$SCRIPT_DIR/harness/go/harness.go" --encode "$SCRIPT_DIR/fixtures/03_tags/child_tags.jsonc" 2>/dev/null)
echo "kt: $( [ -z "$kt_hex" ] && echo 'FAIL' || echo "OK (${#kt_hex} chars)" )"
echo "go: OK (${#go_hex} chars)"

if [ -n "$kt_hex" ] && [ "$kt_hex" = "$go_hex" ]; then
    echo "RESULT: MATCH!"
else
    echo "RESULT: DIFF"
fi

echo ""
echo "=== mime_tag.jsonc ==="
kt_hex=$(java -cp "$tmpdir/harness.jar:$mm_classes:$kotlin_stdlib" HarnessKt --encode "$SCRIPT_DIR/fixtures/03_tags/mime_tag.jsonc" 2>&1)
go_hex=$(go run "$SCRIPT_DIR/harness/go/harness.go" --encode "$SCRIPT_DIR/fixtures/03_tags/mime_tag.jsonc" 2>/dev/null)
echo "kt: $( [ -z "$kt_hex" ] && echo 'FAIL' || echo "OK (${#kt_hex} chars)" )"
echo "go: OK (${#go_hex} chars)"
if [ -n "$kt_hex" ] && [ "$kt_hex" = "$go_hex" ]; then
    echo "RESULT: MATCH!"
else
    echo "RESULT: DIFF"
fi

rm -rf "$tmpdir"