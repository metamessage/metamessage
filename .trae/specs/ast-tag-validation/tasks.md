# AST Tag Validation - The Implementation Plan

## [x] 任务 1：Java AST Tag 验证实现
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 实现 Java 在 AST node 建立时的 tag 验证逻辑
  - 支持 jsonc->ast 和 struct->ast 两个路径
  - 确保验证逻辑与 Golang 一致
- **Acceptance Criteria Addressed**: [AC-1]
- **Test Requirements**:
  - `programmatic` TR-1.1: 验证 Java 代码能够正确验证 JSONC 解析时的 tag 规则
  - `programmatic` TR-1.2: 验证 Java 代码能够正确验证结构体解析时的 tag 规则
  - `programmatic` TR-1.3: 验证 Java 验证结果与 Golang 一致
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 2：TypeScript AST Tag 验证实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 实现 TypeScript 在 AST node 建立时的 tag 验证逻辑
  - 参考 Golang 的 validate.go 实现
  - 确保验证逻辑与 Golang 一致
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3, AC-4, AC-5
- **Test Requirements**:
  - `programmatic` TR-2.1: 验证 TypeScript 代码能够正确验证 JSONC 解析时的 tag 规则
  - `programmatic` TR-2.2: 验证 TypeScript 代码能够正确验证结构体解析时的 tag 规则
  - `programmatic` TR-2.3: 验证 TypeScript 验证结果与 Golang 一致
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 7：Rust AST Tag 验证实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 实现 Rust 在 AST node 建立时的 tag 验证逻辑
  - 参考 Golang 的 validate.go 实现
  - 确保验证逻辑与 Golang 一致
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3, AC-4, AC-5
- **Test Requirements**:
  - `programmatic` TR-7.1: 验证 Rust 代码能够正确验证 JSONC 解析时的 tag 规则
  - `programmatic` TR-7.2: 验证 Rust 代码能够正确验证结构体解析时的 tag 规则
  - `programmatic` TR-7.3: 验证 Rust 验证结果与 Golang 一致
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 8：Swift AST Tag 验证实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 实现 Swift 在 AST node 建立时的 tag 验证逻辑
  - 参考 Golang 的 validate.go 实现
  - 确保验证逻辑与 Golang 一致
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3, AC-4, AC-5
- **Test Requirements**:
  - `programmatic` TR-8.1: 验证 Swift 代码能够正确验证 JSONC 解析时的 tag 规则
  - `programmatic` TR-8.2: 验证 Swift 代码能够正确验证结构体解析时的 tag 规则
  - `programmatic` TR-8.3: 验证 Swift 验证结果与 Golang 一致
- **Notes**: 参考 Golang 的 validate.go 实现

## [x] 任务 9：PHP AST Tag 验证实现
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 实现 PHP 在 AST node 建立时的 tag 验证逻辑
  - 参考 Golang 的 validate.go 实现
  - 确保验证逻辑与 Golang 一致
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3, AC-4, AC-5
- **Test Requirements**:
  - `programmatic` TR-9.1: 验证 PHP 代码能够正确验证 JSONC 解析时的 tag 规则
  - `programmatic` TR-9.2: 验证 PHP 代码能够正确验证结构体解析时的 tag 规则
  - `programmatic` TR-9.3: 验证 PHP 验证结果与 Golang 一致
- **Notes**: 参考 Golang 的 validate.go 实现