# TypeScript MM Tag 功能实现计划

## [x] Task 1: 增强 mm 函数支持 tag 和自动类型推断
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 更新 `MMOptions` 接口，添加 `tag` 字段
  - 确保 `mm` 函数正确处理 tag 选项
  - 实现自动类型推断，字符串和整数类型不需要显式标注 `type`
- **Acceptance Criteria Addressed**: [AC-1]
- **Test Requirements**:
  - `programmatic` TR-1.1: `mm("Ed", { tag: { desc: "姓名" } })` 返回包含 tag 的 MMValue，自动推断类型为 "str"
  - `programmatic` TR-1.2: `mm(30, { tag: { desc: "年龄" } })` 返回包含 tag 的 MMValue，自动推断类型为 "int"
  - `programmatic` TR-1.3: `mm` 函数的所有方法都支持 tag 选项

## [x] Task 2: 确保 encode 函数支持 MMValue 结构
- **Priority**: P0
- **Depends On**: Task 1
- **Description**:
  - 验证 `encode` 函数已正确处理 MMValue 结构
  - 确保 tag 信息被保留
- **Acceptance Criteria Addressed**: [AC-2]
- **Test Requirements**:
  - `programmatic` TR-2.1: `encode({ name: mm("Ed", { tag: { desc: "姓名" } }) })` 成功执行
  - `programmatic` TR-2.2: 编码结果包含正确的类型信息

## [x] Task 3: 实现 bind 函数支持 MMValue 结构
- **Priority**: P0
- **Depends On**: Task 1, Task 2
- **Description**:
  - 在 `mm/index.ts` 中添加 `bind` 函数
  - 实现绑定到 MMValue 结构的功能
  - 保留 tag 信息
- **Acceptance Criteria Addressed**: [AC-3]
- **Test Requirements**:
  - `programmatic` TR-3.1: `bind(data, { name: mm("", { tag: { desc: "姓名" } }) })` 成功执行
  - `programmatic` TR-3.2: 绑定结果包含正确的 tag 信息

## [x] Task 4: 更新示例代码
- **Priority**: P1
- **Depends On**: Task 1, Task 2, Task 3
- **Description**:
  - 更新 `examples/typescript/basic/basic.ts`
  - 使用 `mm()` 函数包装值
  - 添加 tag 使用示例
- **Acceptance Criteria Addressed**: [AC-4]
- **Test Requirements**:
  - `human-judgment` TR-4.1: 示例代码使用 `mm()` 函数
  - `human-judgment` TR-4.2: 示例代码包含 tag 使用

## [x] Task 5: 运行测试验证
- **Priority**: P0
- **Depends On**: Task 1, Task 2, Task 3, Task 4
- **Description**:
  - 运行 TypeScript 测试
  - 验证所有功能正常工作
- **Acceptance Criteria Addressed**: [AC-1, AC-2, AC-3, AC-4]
- **Test Requirements**:
  - `programmatic` TR-5.1: 所有测试通过 - ✅ 49/49 测试通过
  - `programmatic` TR-5.2: 示例代码可正常运行 - ✅ 示例运行成功
