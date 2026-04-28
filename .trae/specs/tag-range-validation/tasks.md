# Tag Min/Max Range Validation - Implementation Plan

## [x] 任务 1：TypeScript Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 TypeScript 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 确保错误消息与 Golang 一致
- **Acceptance Criteria Addressed**: AC-1, AC-2
- **Test Requirements**:
  - `programmatic` TR-1.1: uint16 类型，tag.min 为 70000 时验证失败
  - `programmatic` TR-1.2: int8 类型，tag.min 为 -130 时验证失败
  - `programmatic` TR-1.3: uint8 类型，tag.max 为 260 时验证失败
  - `programmatic` TR-1.4: 其他类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现，使用对应类型的范围进行验证

## [x] 任务 2：Kotlin Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 Kotlin 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-2.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 3：Python Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 Python 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-3.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 4：JavaScript Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 JavaScript 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-4.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 5：Java Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 Java 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-5.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 6：C# Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 C# 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-6.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 7：Rust Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 Rust 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-7.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 8：Swift Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 Swift 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-8.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 9：PHP Tag Range Validation 实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 修复 PHP 验证器中 tag.min 和 tag.max 的范围验证
  - 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
  - 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-9.1: 各类型的边界值验证
- **Notes**: 参考 Golang 的 validate.go 实现

## Task Dependencies
- 所有任务之间没有依赖关系，可以并行实现
