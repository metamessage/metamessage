# Cross-Language JSONC Tests — Phase 3: Fix Remaining Bugs

## Summary

Phase 2 已完成：脚本修复、Python/PHP 的 ValueType 命名修正。测试结果从 3 PASS / 22 FAIL 提升至 12 PASS / 13 FAIL。当前剩余 3 个核心 bug 导致 IR 一致性失败和 JSONC 可逆性失败。本 Phase 修复这些 bug，重新运行测试，确保所有 8 种语言实现产生一致的输出。

## Current State Analysis

### 测试框架状态
- 25 个 JSONC fixtures，6 大类别
- 8 种语言可用（Go, Python, PHP, TypeScript, Rust, C, C++, C#）
- Kotlin/Swift 跳过（环境限制）
- 测试脚本 `run_cross_lang.sh` 已修复并能运行至完成
- 结果保存在 `tests/results/` 目录

### Go 参考实现的行为（作为一致性基准）

Go 的 `writeObjectJSONC` 在 [internal/jsonc/jsonc.go](file:///Users/lizongying/IdeaProjects/meta-message/internal/jsonc/jsonc.go#L71-L89) 中的模式：
```
// mm: type=datetime; desc="creation time"    ← 标签在独立行
    "created_at": "2024-01-01 00:00:00",      ← key: value
```

关键函数 `writeLeadingComments`（[jsonc.go:91-98](file:///Users/lizongying/IdeaProjects/meta-message/internal/jsonc/jsonc.go#L91-L98)）：标签总是打印在值**之前**的独立行，而不是 inline。

Go 将值文本原样存储（不做 Date 对象转换），重新打印时输出与输入相同的文本。

### 已识别的 Bug

---

## Bug 1: C++ JSONC Printer — 标签注释打印在 key:value 中间（inline）

**文件**: [mm-cpp/src/jsonc/printer.hpp](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp)

**根因**: `printValue()` 第 28 行在打印值之前先调用 `printTag()`，而 `printObject()` 第 168 行先打印 `"key": `，再调用 `printValue()`。结果标签被插入到冒号和值之间：

```
"created_at":  // mm: type=datetime; desc="creation time""2024-01-01 00:00:00",
```

**期望输出**（参照 Go）:
```

    // mm: type=datetime; desc="creation time"
    "created_at": "2024-01-01 00:00:00",
```

**修复方案**: 参照 Go 的 `writeLeadingComments` 模式重构：

1. **新增 `printLeadingComment` 函数**（在 `printTag` 之后）：
```cpp
inline void printLeadingComment(std::ostringstream& os, const ir::Tag* tag, int indent) {
    if (tag == nullptr) return;
    std::string ts = tag->toString();
    if (!ts.empty()) {
        os << "\n";
        printIndent(os, indent);
        os << "// mm: " << ts << "\n";
    }
}
```

2. **修改 `printObject`**（[printer.hpp:139-194](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp#L139-L194)）：在每个 field 的 key 之前调用 `printLeadingComment(os, field.value->getTag(), indent + 1)`

3. **修改 `printArray`**（[printer.hpp:86-137](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp#L86-L137)）：在每个 item 之前调用 `printLeadingComment(os, arr->items[i]->getTag(), multiline ? indent + 1 : indent)`

4. **从 `printValue` 中移除 `printTag` 调用**（[printer.hpp:27-28](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp#L27-L28)）

5. **从 `printArray` 中移除 `printTag` 调用**（[printer.hpp:87](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp#L87)）

6. **修改 `toJSONC`**（[printer.hpp:196-213](file:///Users/lizongying/IdeaProjects/meta-message/mm-cpp/src/jsonc/printer.hpp#L196-L213)）：在顶层打印标签

7. **修改 `printObject`**：将 `printTag(os, obj->getTag())` 改为调用 `printLeadingComment`

### 影响范围
- 修复 `type_tags.jsonc`、`constraints_tag.jsonc`、`nullable_tag.jsonc`、`desc_tag.jsonc`、`null_with_tag.jsonc` 等所有带标签的 fixture

---

## Bug 2: TypeScript Datetime/Date/Time 格式化错误

**文件**: [mm-ts/src/jsonc/printer.ts](file:///Users/lizongying/IdeaProjects/meta-message/mm-ts/src/jsonc/printer.ts)

**根因**: 第 81-86 行 `val.toString()` 对 JS Date 对象调用 `Date.prototype.toString()`，产生原生格式如 `"Mon Jan 01 2024 00:00:00 GMT+0800 (China Standard Time)"`，而非 `"2024-01-01 00:00:00"`。

这是因为 parser.ts 第 167 行 `new Date(text)` 将 datetime 字符串转换为了 JS Date 对象，丢失了原始文本格式。

**修复方案**: 在 printer.ts 的 `valueToStringOnly()` 中，为 DateTime/Date/Time 类型添加日期格式化：

```typescript
case ValueType.DateTime:
case ValueType.Date:
case ValueType.Time:
    if (val instanceof Date) {
        const pad = (n: number) => String(n).padStart(2, '0');
        return `"${val.getFullYear()}-${pad(val.getMonth() + 1)}-${pad(val.getDate())} ${pad(val.getHours())}:${pad(val.getMinutes())}:${pad(val.getSeconds())}"`;
    }
    return `"${val}"`;
```

### 影响范围
- 修复 `type_tags.jsonc`（包含 datetime 字段）

---

## Bug 3: TypeScript `validateArr` 传入了空数组

**文件**: [mm-ts/src/jsonc/parser.ts](file:///Users/lizongying/IdeaProjects/meta-message/mm-ts/src/jsonc/parser.ts)

**根因**: 第 737 行 `tag.validateArr([])` 始终传入空数组 `[]`，无论实际解析了多少元素。`validateArr`（[tag.ts:1739-1780](file:///Users/lizongying/IdeaProjects/meta-message/mm-ts/src/ir/tag.ts#L1753-L1758)）在 `length === 0` 且没有 `allowEmpty` 时返回 `"type array not allow empty"` 错误。

**修复方案**: 传入实际解析的元素：
```typescript
// 第 737 行，将：
const result = tag.validateArr([]);
// 改为：
const result = tag.validateArr(arr.getElements());
```

### 影响范围
- 修复 `constraints_tag.jsonc`（包含 `"items": [1, 2, 3]` 带 `size=10` 标签）

---

## Verification Steps

修复完成后，执行以下验证：

```bash
# 1. 构建所有 harnesses
cd /Users/lizongying/IdeaProjects/meta-message

# 2. 测试 C++ 修复：验证 type_tags.jsonc 输出格式正确
g++ -std=c++17 -I mm-cpp/src -o /tmp/test_cpp mm-cpp/src/jsonc/scanner.cpp tests/harness/cpp/harness.cpp 2>/dev/null
/tmp/test_cpp tests/fixtures/03_tags/type_tags.jsonc

# 3. 测试 TypeScript 修复
cd mm-ts && npm run build --silent 2>/dev/null && cd ..
node tests/harness/typescript/harness.cjs tests/fixtures/03_tags/type_tags.jsonc

# 4. 运行全部测试
bash tests/run_cross_lang.sh

# 5. 检查结果文件
ls -la tests/results/
```

预期结果：
- C++ 和 TypeScript 输出与 Go 匹配（归一化后 JSON 一致）
- `constraints_tag.jsonc` 在所有语言中解析成功
- 可逆性测试大部分通过

## Assumptions & Decisions

- Go 实现仍为参考基准
- C++ 标签输出格式对齐 Go 的 `writeLeadingComments` 模式
- TypeScript datetime 使用 `YYYY-MM-DD HH:mm:ss` 格式（与 fixture 输入格式一致）
- TypeScript `validateArr` 传入 `Node[]`（TypeScript 允许，`validateArr` 只访问 `.length` 和索引值）
- 仅修复已识别的 bug，不做额外的重构或优化