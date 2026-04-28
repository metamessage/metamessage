# TypeScript MM Tag 功能实现计划

## 概述

为 TypeScript 实现更精确的类型支持，使用 `mm()` 函数包装值并支持 tag 功能，使示例代码从 `const person = { name: "Ed", age: 30 }` 变为 `const person = { name: mm("Ed", {tag}), age: mm(30, {tag}) }`。

## 目标

* 增强 `mm` 函数支持 tag 功能

* 确保 `encode` 函数支持接受 `mm()` 包装的结构

* 确保 `bind` 功能支持绑定到 `mm()` 包装的结构

* 更新示例代码展示新用法

## 背景

当前 TypeScript 实现已经有基本的 `mm` 函数和 `MMValue` 接口，但缺少 tag 相关的完整支持。

## 功能需求

### FR-1: 增强 mm 函数支持 tag 和自动类型推断
- **FR-1.1**: `mm` 函数应接受 `tag` 选项
- **FR-1.2**: 提供设置和读取 tag 的方法
- **FR-1.3**: 自动推断类型，字符串和整数类型不需要显式标注 `type`

### FR-2: 确保 encode 函数支持 MMValue 结构
- **FR-2.1**: `encode` 函数应正确处理 `mm()` 包装的对象
- **FR-2.2**: 保留 tag 信息

### FR-3: 实现 bind 功能支持 MMValue 结构
- **FR-3.1**: 实现 `bind` 函数，支持绑定到 `mm()` 包装的结构
- **FR-3.2**: 保留 tag 信息

### FR-4: 更新示例代码
- **FR-4.1**: 更新 `examples/typescript/basic/basic.ts` 示例
- **FR-4.2**: 添加 tag 使用示例，展示自动类型推断

## 非功能需求

* **NFR-1**: 类型安全 - 保持 TypeScript 类型定义清晰

* **NFR-2**: 向后兼容 - 确保现有代码仍能正常工作

## 约束

* **技术**: TypeScript 4.5+, CommonJS 模块系统

* **依赖**: 无额外依赖

## 假设

* 现有 `mm` 函数和 `MMValue` 接口结构保持不变

* tag 信息应包含在 `MMOptions` 中

## 验收标准

### AC-1: mm 函数支持 tag 和自动类型推断
- **Given**: 调用 `mm("Ed", { tag: { desc: "姓名" } })`
- **When**: 检查返回值
- **Then**: 返回包含 tag 信息的 `MMValue` 对象，自动推断类型为 "str"
- **Verification**: `programmatic`

### AC-1.1: mm 函数自动推断整数类型
- **Given**: 调用 `mm(30, { tag: { desc: "年龄" } })`
- **When**: 检查返回值
- **Then**: 返回包含 tag 信息的 `MMValue` 对象，自动推断类型为 "int"
- **Verification**: `programmatic`

### AC-2: encode 函数支持 MMValue
- **Given**: 传入 `{ name: mm("Ed", { tag: { desc: "姓名" } }) }`
- **When**: 调用 `encode` 函数
- **Then**: 成功编码，保留 tag 信息，自动推断类型为 "str"
- **Verification**: `programmatic`

### AC-3: bind 函数支持 MMValue
- **Given**: 传入 Wire 数据和 `{ name: mm("", { tag: { desc: "姓名" } }) }`
- **When**: 调用 `bind` 函数
- **Then**: 成功绑定到 MMValue 结构，自动推断类型为 "str"
- **Verification**: `programmatic`

### AC-4: 示例代码更新

* **Given**: 查看 `examples/typescript/basic/basic.ts`

* **When**: 检查代码内容

* **Then**: 使用 `mm()` 函数包装值

* **Verification**: `human-judgment`

## 开放问题

* [ ] 是否需要为 tag 提供专门的类型定义？

* [ ] bind 函数的具体实现方式？

