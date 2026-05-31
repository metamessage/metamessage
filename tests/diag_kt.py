#!/usr/bin/env python3
import subprocess, tempfile, os

SCRIPT_DIR = "/Users/lizongying/IdeaProjects/meta-message/tests"
PROJECT_DIR = "/Users/lizongying/IdeaProjects/meta-message"
HOME = os.path.expanduser("~")

trove4j = f"{HOME}/.m2/repository/org/jetbrains/intellij/deps/trove4j/1.0.20221201/trove4j-1.0.20221201.jar"
annotations = f"{HOME}/.m2/repository/org/jetbrains/annotations/13.0/annotations-13.0.jar"
kc = f"{HOME}/.m2/repository/org/jetbrains/kotlin/kotlin-compiler/1.9.22/kotlin-compiler-1.9.22.jar"
ks = f"{HOME}/.m2/repository/org/jetbrains/kotlin/kotlin-stdlib/1.9.22/kotlin-stdlib-1.9.22.jar"
kr = f"{HOME}/.m2/repository/org/jetbrains/kotlin/kotlin-reflect/1.9.22/kotlin-reflect-1.9.22.jar"
cp = f"{kc}:{ks}:{kr}:{annotations}:{trove4j}"
mc = f"{PROJECT_DIR}/mm-kt/target/classes"

tmpdir = tempfile.mkdtemp()

# Compile harness
r = subprocess.run([
    "java", "-cp", cp,
    "org.jetbrains.kotlin.cli.jvm.K2JVMCompiler",
    f"{SCRIPT_DIR}/harness/kotlin/harness.kt",
    "-cp", f"{mc}:{ks}:{kr}",
    "-d", f"{tmpdir}/harness.jar",
    "-no-stdlib", "-no-reflect"
], capture_output=True, text=True)
print("COMPILE stderr:", r.stderr[:200])

# Run encode
r = subprocess.run([
    "java", "-cp", f"{tmpdir}/harness.jar:{mc}:{ks}",
    "HarnessKt", "--encode",
    f"{SCRIPT_DIR}/fixtures/03_tags/child_tags.jsonc"
], capture_output=True, text=True)
print("STDOUT length:", len(r.stdout))
print("STDERR:", r.stderr[:500])
print("Return code:", r.returncode)

if r.stdout:
    print("STDOUT first 100:", r.stdout[:100])

# Cleanup
import shutil
shutil.rmtree(tmpdir)