#!/usr/bin/env bash
# env_sanitize.sh — neutralise externally-injected build environment
#
# Cross-language harness scripts (encode_decode_test.sh, compare_with_go.sh,
# run_cross_lang.sh, boundary_cross_lang.sh) invoke gcc/clang/cargo/mvn/
# dotnet/swiftc/node/php in the same shell session. When the caller has
# exported compiler-selection or search-path variables (commonly via
# ~/.zshenv, direnv, CI wrappers, or `source activate` scripts), the child
# build steps inherit them and break in confusing ways:
#
#   * CC/CXX pinned to a Homebrew clang while the harness expects Apple
#     clang → CMake configures a mismatched toolchain.
#   * LIBRARY_PATH / LD_LIBRARY_PATH pointing to unrelated trees → linker
#     picks up the wrong libstdc++ or LLVM dylibs.
#   * CPPFLAGS / CFLAGS / CXXFLAGS / LDFLAGS with -isystem or -isysroot
#     entries → headers from a stale SDK shadow the current one.
#   * PKG_CONFIG_PATH drift → cmake find_package resolves to a different
#     version than the port's own CMakeLists expects.
#   * GOFLAGS with -mod=vendor against a module that isn't vendored.
#
# Sourced scripts should not modify these; they should reset them to a
# clean baseline and let each port's build system pick defaults.
#
# This file is meant to be sourced, not executed:
#   source "$(dirname "$0")/env_sanitize.sh"
#
# Anything the caller legitimately needs (e.g. JAVA_HOME, RUSTUP_HOME,
# CARGO_HOME, DOTNET_ROOT, NVM_DIR) is intentionally left alone.

# NOTE: this file must not `set -o errexit/nounset/pipefail` because it is
# sourced — the options would leak into the harness scripts and change their
# error-handling semantics. Each harness script sets its own shell options.

# --- C/C++ toolchain selection -----------------------------------------------
unset CC CXX
unset CFLAGS CXXFLAGS CPPFLAGS LDFLAGS
unset LIBRARY_PATH LD_LIBRARY_PATH
unset DYLD_LIBRARY_PATH DYLD_FALLBACK_LIBRARY_PATH
unset DYLD_FRAMEWORK_PATH DYLD_FALLBACK_FRAMEWORK_PATH
unset INCLUDE PATH_TO_INCLUDE
unset CPATH

# --- pkg-config --------------------------------------------------------------
unset PKG_CONFIG_PATH PKG_CONFIG_LIBDIR PKG_CONFIG_SYSROOT_DIR

# --- Go ----------------------------------------------------------------------
unset GOFLAGS GOEXPERIMENT
# GOOS / GOARCH: only override if the caller forced a cross-compile value;
# native builds should let `go env` pick defaults.
unset GOOS GOARCH

# --- Rust / Cargo ------------------------------------------------------------
unset RUSTFLAGS RUSTC_WRAPPER CARGO_BUILD_TARGET
# Leave RUSTUP_HOME / CARGO_HOME / PATH untouched so rustup-managed
# toolchains still resolve.

# --- Java / Maven / Gradle ---------------------------------------------------
unset MAVEN_OPTS MAVEN_ARGS GRADLE_OPTS
# Leave JAVA_HOME alone: some ports legitimately pin a specific JDK.

# --- Node / TypeScript -------------------------------------------------------
unset NODE_OPTIONS npm_config_prefix
# Leave NVM_DIR / PATH alone.

# --- Python ------------------------------------------------------------------
unset PYTHONPATH PYTHONHOME PYTHONDONTWRITEBYTECODE
# Leave VIRTUAL_ENV alone if it was activated explicitly by the caller.

# --- PHP ---------------------------------------------------------------------
unset PHP_INI_SCAN_DIR

# --- .NET --------------------------------------------------------------------
# Leave DOTNET_ROOT / DOTNET_CLI_HOME alone; the SDK needs them to boot.
unset DOTNET_SYSTEM_GLOBALIZATION_INVARIANT

# --- macOS SDK root ----------------------------------------------------------
# If SDKROOT was set to a specific SDK by an outer shell (e.g. an older
# Xcode), unset it and let each toolchain resolve via `xcrun --show-sdk-path`.
unset SDKROOT
