# Per-Type Validation - Product Requirement Document

## Overview
- **Summary**: 修复各语言的验证器实现，确保每个类型都有独立的验证方法，严格按照 Golang 的实现方式验证各类型的范围。
- **Purpose**: 确保各语言的验证结果与 Golang 一致，特别是对于有明确范围限制的类型（如 int8、int16、uint8 等）。
- **Target Users**: 开发人员和库用户，需要在各语言中使用 MetaMessage 的验证功能。

## Why
在 Golang 中，每个类型都有独立的验证方法（如 `ValidateInt8`、`ValidateInt16`、`ValidateUint8` 等），这些方法不仅检查类型，还验证值的范围是否在类型允许的范围内。然而，其他语言的实现（如 TypeScript）使用了统一的 `validateInt` 方法，只检查了值是否为整数，但没有验证值是否在特定类型的范围内。

例如，对于 `int8` 类型：
- 有效范围是 -128 到 127
- 如果值是 200，在 TypeScript 中会被认为有效（因为它是整数）
- 但在 Golang 中会被拒绝（因为 200 超出了 int8 的范围）

## What Changes
- **TypeScript**: 将统一的 `validateInt` 方法拆分为 `validateInt8`、`validateInt16`、`validateInt32`、`validateInt64`、`validateUint`、`validateUint8`、`validateUint16`、`validateUint32`、`validateUint64`，并添加范围检查
- **Kotlin**: 同样实现各类型的独立验证方法
- **Python**: 同样实现各类型的独立验证方法
- **JavaScript**: 同样实现各类型的独立验证方法，并使用 mmValue 支持更多数据类型
- **Java**: 同样实现各类型的独立验证方法
- **C#**: 同样实现各类型的独立验证方法
- **Rust**: 同样实现各类型的独立验证方法
- **Swift**: 同样实现各类型的独立验证方法
- **PHP**: 同样实现各类型的独立验证方法

## Impact
- **Affected Specs**: `ast-tag-validation`
- **Affected Code**:
  - `mm-ts/src/mm/validator.ts`
  - `mm-kt/src/main/kotlin/io/metamessage/mm/MmValidator.kt`
  - `mm-py/metamessage/validator.py`
  - `mm-js/src/mm/validator.js`
  - `mm-java/src/main/java/io/metamessage/mm/MmValidator.java`
  - `mm-cs/src/MetaMessage/Mm/MmValidator.cs`
  - `mm-rs/src/mm/validator.rs`
  - `mm-swift/Sources/MetaMessage/MM/MmValidator.swift`
  - `mm-php/src/io/metamessage/mm/MmValidator.php`

## ADDED Requirements

### Requirement: Per-Type Integer Validation
各语言的验证器必须为每个整数类型实现独立的验证方法，严格验证值的范围。

#### Scenario: Int8 Validation
- **Given**: Tag 类型为 `int8`，值为 `200`
- **When**: 调用验证方法时
- **Then**: 返回验证失败，因为 200超出了 int8 的范围（-128 到 127）
- **Verification**: `programmatic`

#### Scenario: Uint8 Validation
- **Given**: Tag 类型为 `uint8`，值为 `300`
- **When**: 调用验证方法时
- **Then**: 返回验证失败，因为 300 超出了 uint8 的范围（0 到 255）
- **Verification**: `programmatic`

### Requirement: Per-Type Float Validation
各语言的验证器必须为每个浮点类型实现独立的验证方法。

#### Scenario: Float32 Validation
- **Given**: Tag 类型为 `float32`
- **When**: 验证浮点值时
- **Then**: 使用 float32 的范围和精度进行验证
- **Verification**: `programmatic`

### Requirement: Per-Type BigInt Validation
各语言的验证器必须实现 `big.Int` 类型的验证方法。

#### Scenario: BigInt Validation
- **Given**: Tag 类型为 `big.Int`
- **When**: 验证大整数时
- **Then**: 使用大整数验证方法，支持任意大小的整数
- **Verification**: `programmatic`

## MODIFIED Requirements

### Requirement: Existing Type Validation
所有语言的验证器必须严格按照 Golang 的实现方式，为每个类型实现独立的验证方法，不再使用统一的验证方法处理多个类型。

## REMOVED Requirements

### Requirement: Unified Int Validation
移除使用统一 `validateInt` 方法处理所有整数类型的实现。

## Acceptance Criteria

### AC-1: TypeScript Per-Type Validation
- **Given**: TypeScript 代码使用 MetaMessage 验证不同整数类型
- **When**: 验证 `int8` 值为 200 或 `uint8` 值为 300 时
- **Then**: 返回验证失败，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-2: All Languages Per-Type Validation
- **Given**: 所有语言代码使用 MetaMessage 验证不同整数类型
- **When**: 验证各类型的值超出其范围时
- **Then**: 返回验证失败，与 Golang 验证结果一致
- **Verification**: `programmatic`

## Open Questions
- [ ] 各语言的浮点数范围如何精确定义？
- [ ] BigInt 在某些语言中可能没有原生支持，如何处理？
- [ ] 测试用例如何设计以覆盖所有边界情况？
