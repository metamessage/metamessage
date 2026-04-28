# Per-Type Validation - Implementation Plan

## [x] 任务 1：TypeScript Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 mm-ts/src/mm/validator.ts 中的统一 `validateInt` 方法拆分为独立的验证方法
  - 实现 `validateInt8`、`validateInt16`、`validateInt32`、`validateInt64`
  - 实现 `validateUint`、`validateUint8`、`validateUint16`、`validateUint32`、`validateUint64`
  - 实现 `validateFloat32`、`validateFloat64`、`validateBigInt`
  - 每个方法严格按照 Golang 的范围进行验证
- **Acceptance Criteria Addressed**: AC-1, AC-2
- **Test Requirements**:
  - `programmatic` TR-1.1: int8 范围验证（-128 到 127）
  - `programmatic` TR-1.2: int16 范围验证（-32768 到 32767）
  - `programmatic` TR-1.3: int32 范围验证（-2147483648 到 2147483647）
  - `programmatic` TR-1.4: int64 范围验证
  - `programmatic` TR-1.5: uint8 范围验证（0 到 255）
  - `programmatic` TR-1.6: uint16 范围验证（0 到 65535）
  - `programmatic` TR-1.7: uint32 范围验证
  - `programmatic` TR-1.8: uint64 范围验证
  - `programmatic` TR-1.9: float32 范围验证
  - `programmatic` TR-1.10: float64 范围验证
  - `programmatic` TR-1.11: big.Int 验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 2：Kotlin Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 mm-kt 中的验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数和 BigInt 的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-2.1: 各整数类型范围验证
  - `programmatic` TR-2.2: 浮点数类型范围验证
  - `programmatic` TR-2.3: big.Int 验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 3：Python Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 Python 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数和 BigInt 的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-3.1: 各整数类型范围验证
  - `programmatic` TR-3.2: 浮点数类型范围验证
  - `programmatic` TR-3.3: big.Int 验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 4：JavaScript Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 JavaScript 验证器实现改为每类型独立验证
  - 使用 mmValue 实现更多数据类型支持
  - 实现各整数类型的独立验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-4.1: 各整数类型范围验证
  - `programmatic` TR-4.2: mmValue 数据类型支持
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 5：Java Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 Java 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数和 BigInteger 的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-5.1: 各整数类型范围验证
  - `programmatic` TR-5.2: 浮点数类型范围验证
  - `programmatic` TR-5.3: BigInteger 验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 6：C# Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 C# 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数类型的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-6.1: 各整数类型范围验证
  - `programmatic` TR-6.2: 浮点数类型范围验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 7：Rust Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 Rust 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数类型的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-7.1: 各整数类型范围验证
  - `programmatic` TR-7.2: 浮点数类型范围验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 8：Swift Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 Swift 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数类型的验证方法
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-8.1: 各整数类型范围验证
  - `programmatic` TR-8.2: 浮点数类型范围验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 9：PHP Per-Type Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 将 PHP 验证器实现改为每类型独立验证
  - 实现各整数类型的独立验证方法
  - 实现浮点数类型的验证方法
  - 实现 BigInteger 的验证方法（使用 GMP 或 BCMath）
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-9.1: 各整数类型范围验证
  - `programmatic` TR-9.2: 浮点数类型范围验证
  - `programmatic` TR-9.3: BigInteger 验证
- **Notes**: 参考 Golang 的 validate.go 实现

## Task Dependencies
- 所有任务之间没有依赖关系，可以并行实现