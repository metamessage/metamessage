# AST Tag Validation - Product Requirement Document

## Overview
- **Summary**: 实现各语言在 AST node 建立时对 tag 的验证逻辑，包括 jsonc->ast 和 struct->ast 两个路径，确保与 Golang 的验证逻辑一致。
- **Purpose**: 确保各语言的 AST node 建立过程中能够正确验证 tag 规则，提高数据的一致性和可靠性。
- **Target Users**: 开发人员和库用户，需要在各语言中使用 MetaMessage 的 AST 功能。

## Goals
- 为 Java、TypeScript、Kotlin、Python、JavaScript、C#、Rust、Swift、PHP 等语言实现 AST node 建立时的 tag 验证逻辑
- 确保验证逻辑与 Golang 的实现一致
- 支持 jsonc->ast 和 struct->ast 两个路径的验证
- 为 JavaScript 实现 mmValue 以支持更多数据类型的验证

## Non-Goals (Out of Scope)
- 不修改 Golang 的原始验证逻辑
- 不改变现有 API 的使用方式
- 不涉及运行时的验证，仅关注 AST node 建立时的验证

## Background & Context
- Golang 已经实现了完整的 AST node 建立时的 tag 验证逻辑
- 各语言已经实现了基本的验证功能，但需要在 AST node 建立时集成验证逻辑
- JavaScript 需要通过 mmValue 实现更多数据类型的支持

## Functional Requirements
- **FR-1**: 实现 jsonc->ast 路径的 tag 验证
- **FR-2**: 实现 struct->ast 路径的 tag 验证
- **FR-3**: 支持所有数据类型的验证，包括数组、结构体、字符串、字节、布尔值、整数、浮点数、日期时间、UUID、IP、URL、邮箱、枚举、图像等
- **FR-4**: 为 JavaScript 实现 mmValue 以支持更多数据类型
- **FR-5**: 确保验证错误信息与 Golang 一致

## Non-Functional Requirements
- **NFR-1**: 验证逻辑与 Golang 实现一致
- **NFR-2**: 性能影响最小化
- **NFR-3**: 代码可读性和可维护性
- **NFR-4**: 跨语言一致性

## Constraints
- **Technical**: 各语言的语法和特性差异
- **Dependencies**: 依赖各语言的现有 MetaMessage 实现

## Assumptions
- 各语言已经实现了基本的 MetaMessage 功能
- 各语言已经有了 tag 结构的定义
- Golang 的验证逻辑是标准参考

## Acceptance Criteria

### AC-1: Java AST Tag Validation
- **Given**: Java 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-2: TypeScript AST Tag Validation
- **Given**: TypeScript 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-3: Kotlin AST Tag Validation
- **Given**: Kotlin 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-4: Python AST Tag Validation
- **Given**: Python 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-5: JavaScript AST Tag Validation
- **Given**: JavaScript 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，使用 mmValue 支持更多数据类型
- **Verification**: `programmatic`

### AC-6: C# AST Tag Validation
- **Given**: C# 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-7: Rust AST Tag Validation
- **Given**: Rust 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-8: Swift AST Tag Validation
- **Given**: Swift 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

### AC-9: PHP AST Tag Validation
- **Given**: PHP 代码使用 MetaMessage 解析 JSONC 或结构体
- **When**: 建立 AST node 时
- **Then**: 正确验证 tag 规则，与 Golang 验证结果一致
- **Verification**: `programmatic`

## Open Questions
- [ ] 各语言的 AST 建立流程是否与 Golang 一致？
- [ ] JavaScript 的 mmValue 实现需要支持哪些具体的数据类型？
- [ ] 各语言的错误处理机制是否与 Golang 一致？