# Tag Min/Max Range Validation - Product Requirement Document

## Overview
- **Summary**: 修复各语言验证器中 tag.min 和 tag.max 的范围验证，确保它们也在对应类型的范围内，与 Golang 实现保持一致。
- **Purpose**: 确保 tag.min 和 tag.max 值也符合对应类型的范围限制，避免无效的验证规则。
- **Target Users**: 开发人员和库用户，需要在各语言中使用 MetaMessage 的验证功能。

## Goals
- 为每种整数类型实现 tag.min 和 tag.max 的范围验证
- 确保验证逻辑与 Golang 实现一致
- 修复所有语言的类似问题

## Non-Goals (Out of Scope)
- 不修改其他验证逻辑
- 不改变现有 API 的使用方式

## Background & Context
- 在 Golang 中，对于每种整数类型，它使用对应的位大小来解析和验证 tag.min 和 tag.max
- 例如，对于 uint16 类型，它使用 strconv.ParseUint(t.Min, 10, 16) 来解析 min 值
- 这样如果 min 值超出 uint16 的范围就会返回错误
- 但在其他语言的实现中，只是简单地使用 parseInt 或类似函数解析，没有验证范围

## Functional Requirements
- **FR-1**: 对于每种整数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **FR-2**: 对于每种浮点数类型，验证 tag.min 和 tag.max 是否在该类型的范围内
- **FR-3**: 确保错误消息与 Golang 一致

## Non-Functional Requirements
- **NFR-1**: 验证逻辑与 Golang 实现一致
- **NFR-2**: 性能影响最小化
- **NFR-3**: 代码可读性和可维护性

## Constraints
- **Technical**: 各语言的语法和特性差异
- **Dependencies**: 依赖各语言的现有 MetaMessage 实现

## Assumptions
- 各语言已经实现了基本的 MetaMessage 功能
- Golang 的验证逻辑是标准参考

## Acceptance Criteria

### AC-1: TypeScript Tag Range Validation
- **Given**: TypeScript 代码使用 MetaMessage 验证 uint16 类型，tag.min 为 70000
- **When**: 调用验证方法时
- **Then**: 返回验证失败，因为 70000 超出了 uint16 的范围（0 到 65535）
- **Verification**: `programmatic`

### AC-2: All Languages Tag Range Validation
- **Given**: 所有语言代码使用 MetaMessage 验证各类型，tag.min 或 tag.max 超出对应类型的范围
- **When**: 调用验证方法时
- **Then**: 返回验证失败，与 Golang 验证结果一致
- **Verification**: `programmatic`

## Open Questions
- [ ] 各语言如何处理超出范围的浮点数？
- [ ] 如何统一错误消息格式？
