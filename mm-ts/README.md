# MetaMessage

MetaMessage (mm) is a structured data exchange protocol. It is self-describing, self-constraining, and self-exemplifying, enabling lossless data exchange. It is designed as a next-generation universal protocol that natively supports AI, humans, and machines.

- Human and AI friendly
- Export/import to JSONC (currently; YAML/TOML support planned)
- Suitable for configuration files and data exchange
- Works for traditional APIs and AI interaction scenarios
- Supports conversion between language structs/classes and MetaMessage
- Supports code generation for multiple languages
- Data carries type, constraint, description, and example without separate documentation
- All metadata can be updated with the data itself, without extra coordination
- Structures and values stay consistent across languages
- No structural loss; parsers adapt automatically and do not crash
- Can serialize to compact binary for faster decoding and smaller size

**Problems solved**

- Unknown types, such as not knowing whether a field is uint8
- Incomplete structure, such as null without inner type information
- No validation rules, so data legality cannot be checked
- No examples or descriptions, forcing reliance on separate docs
- Format changes require protocol adjustment and documentation resync

MetaMessage is naturally suited for AI understanding and interaction, solving ambiguity and imprecision in data. It replaces traditional API docs, verbal format agreements, and manual version sync by making data self-explanatory and independently evolvable.

[meta-message](https://github.com/metamessage/metamessage)

## 1. 安装

### npm 依赖

使用 npm 安装：

```bash
npm i metamessage@latest
```

### 版本要求

- Node.js 18 或更高版本
- TypeScript 5.0 或更高版本

## 2. 基本使用

### 2.1 引用

```typescript
import { encodeFromValue, decodeToValue, mm, ValueType } from 'metamessage';
```

or

```typescript
const { encodeFromValue, decodeToValue, mm, ValueType } = require('metamessage');
```

### 2.2 编码示例

javascript

```javascript
// 使用 mm() 函数创建带类型信息的对象
const person = {
  name: mm("Ed", { desc: "" }),
  email: mm("Ed@gmail.com", { desc: "", type: "email" }),
  score: mm(90, { desc: "", type: "uint8" }),
  age: mm(30, { desc: "" })
};

const wire = encodeFromValue(person);
console.log('wire', wire);
```

typescript

```typescript
@mm({ desc: '用户' })
class User {
  @mm({ type: ValueType.Int64, desc: '用户ID', nullable: false })
  id: bigint = 0n;
  @mm({ desc: '昵称' })
  name: string = '';
  @mm({ type: ValueType.Uint8 })
  age: number = 0;
}

const u = new User();
u.id = 666n;
u.name = 'abc';
u.age = 20;

const wire = encodeFromValue(u);
console.log('wire', wire);
```

### 2.3 解码示例

typescript

```typescript
```

typescript

```typescript
const decoded = decodeToValue(wire, User);
console.log('Decoded:', decoded);
```

### 2.4 JSONC 解析示例

javascript

```javascript
import { encodeFromJSONC, decodeToJSONC } from 'metamessage';

const jsonc = `
        // mm: desc="用户"
        {
        
                // mm: type=i64; desc="用户ID"
                "id": 666,
        
                // mm: desc="昵称"
                "name": "abc",
        
                // mm: type=u8
                "age": 20,
        }
`
const wire = encodeFromJSONC(jsonc);
const jsoncString = decodeToJSONC(wire);
console.log('JSONC:', jsoncString);
```

typescript

```typescript
import { encodeFromJSONC, decodeToJSONC } from 'metamessage';

const jsonc = `
        // mm: desc="用户"
        {
        
                // mm: type=i64; desc="用户ID"
                "id": 666,
        
                // mm: desc="昵称"
                "name": "abc",
        
                // mm: type=u8
                "age": 20,
        }
`
const wire = encodeFromJSONC(jsonc);
const jsoncString = decodeToJSONC(wire);
console.log('JSONC:', jsoncString);
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-ts 目录下运行
npm test
```

### 3.2 测试框架

- Jest
- TypeScript Compiler

### 3.3 测试覆盖范围

- 编码测试
- 解码测试
- JSONC 解析测试

## 4. 常见问题

### 4.1 依赖问题

- **问题**: npm 依赖安装失败
  **解决**: 检查网络连接，或使用 npm 镜像

### 4.2 编译问题

- **问题**: TypeScript 编译错误
  **解决**: 确保 tsconfig.json 配置正确

### 4.3 运行时问题

- **问题**: 编码/解码失败
  **解决**: 检查数据类型是否正确

## 5. 示例代码

查看 `examples/typescript/` 目录下的示例代码：

- `basic/` - 基本使用示例
- `jsonc-to-wire/` - JSONC 转 Wire 格式

## 6. 相关资源

- [metamessage](https://www.npmjs.com/package/metamessage)
- [TypeScript 文档](https://www.typescriptlang.org/docs/)
- [Node.js 文档](https://nodejs.org/docs/)
- [Jest 文档](https://jestjs.io/docs/getting-started)
