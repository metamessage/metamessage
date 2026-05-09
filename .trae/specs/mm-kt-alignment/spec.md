# MM-Kotlin 深度对齐 Golang 实现 - 规格文档

## Overview

### Summary
对 mm-kt 模块进行深度对齐，确保 Kotlin 实现与 Golang 基准实现完全一致。

### Purpose
Golang 是 MetaMessage 的基准实现，需要确保 Kotlin 版本在所有功能上与 Golang 保持一致，包括：
- AST 结构和验证逻辑
- JSONC 解析和打印
- MM 编解码
- Tag 处理和验证

### Target Users
MetaMessage 开发者和使用者，需要跨语言兼容性保证。

## Goals
- 确保 Kotlin 的 AST 结构与 Golang 完全一致
- 确保 JSONC 解析器行为与 Golang 完全一致
- 确保 MM 编解码逻辑与 Golang 完全一致
- 确保 Tag 验证逻辑与 Golang 完全一致
- 提供完整的测试覆盖

## Non-Goals (Out of Scope)
- 不修改其他语言的实现
- 不添加 Golang 中不存在的新功能
- 不进行性能优化（除非是为了对齐）

## Background & Context

当前状态：
- Kotlin 实现已完成基本功能
- 但在细节上可能与 Golang 存在差异
- 需要进行深度对比和对齐

## Functional Requirements

### FR-1: AST 结构对齐
Kotlin 的 AST 类（Node、Value、Object、Array、Field）应与 Golang 完全一致。

### FR-2: Tag 验证对齐
Tag 类的验证方法应与 Golang 的 validate.go 完全一致。

### FR-3: JSONC 解析对齐
JSONC 解析器行为应与 Golang 的 jsonc/parser/parser.go 完全一致。

### FR-4: JSONC 打印对齐
JSONC 打印机行为应与 Golang 实现完全一致。

### FR-5: MM 编解码对齐
MM 编解码逻辑应与 Golang 的 internal/mm/ 实现完全一致。

## Non-Functional Requirements

### NFR-1: 测试覆盖
所有对齐工作应有对应的测试验证。

### NFR-2: 代码质量
代码应符合 Kotlin 最佳实践，结构清晰。

## Constraints

### Technical
- Kotlin 版本兼容性
- 依赖库版本限制

### Dependencies
- 依赖 Golang 基准实现作为参考

## Assumptions
- Golang 实现是正确的基准
- 所有差异都需要对齐到 Golang

## Acceptance Criteria

### AC-1: AST 结构一致
- **Given**: Kotlin 和 Golang 有相同的输入
- **When**: 创建 AST 节点
- **Then**: 两者的结构和字段完全一致
- **Verification**: `programmatic`

### AC-2: Tag 验证行为一致
- **Given**: 相同的输入和 Tag 配置
- **When**: 执行验证
- **Then**: 验证结果（valid/error）完全一致
- **Verification**: `programmatic`

### AC-3: JSONC 解析结果一致
- **Given**: 相同的 JSONC 输入
- **When**: 解析 JSONC
- **Then**: 生成的 AST 完全一致
- **Verification**: `programmatic`

### AC-4: JSONC 打印结果一致
- **Given**: 相同的 AST 输入
- **When**: 打印为 JSONC
- **Then**: 输出字符串完全一致
- **Verification**: `programmatic`

### AC-5: MM 编解码一致
- **Given**: 相同的输入数据
- **When**: 编码为 Wire 格式，然后解码
- **Then**: 解码结果与原始输入一致
- **Verification**: `programmatic`

## Open Questions
- [ ] 是否需要对齐测试文件结构？
- [ ] 是否需要添加性能测试？