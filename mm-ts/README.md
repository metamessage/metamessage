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

[npmjs.com](https://www.npmjs.com/package/metamessage)

[github.com](https://github.com/metamessage/metamessage)

## 核心优势

- 自描述数据：数据本身携带类型、约束、描述和示例
- JSONC 支持：可直接解析带注释的 JSONC 格式
- 多语言兼容：支持对象/类与 MetaMessage 之间的转换
- 无结构损失：解析器自动适配，不会因缺失字段崩溃

## 1. 安装

使用 npm 安装最新版本：

```bash
npm install metamessage@latest
```

版本要求：

- Node.js
- TypeScript

## 2. 快速开始

### 2.1 导入模块

```ts
import { encodeFromValue, decodeToValue, mm, ValueType } from 'metamessage';
```

如果你使用 CommonJS：

```js
const {
  encodeFromValue,
  decodeToValue,
  mm,
  ValueType,
} = require('metamessage');
```

主要函数：

- `encodeFromValue(value, tag)`：将对象转换为 MetaMessage Wire 格式
- `encodeFromJsonc(jsonc: string): Uint8Array`：将 JSONC 格式转换为 MetaMessage Wire 格式
- `decodeToValue(wire, type)`：将 MetaMessage Wire 格式转换为对象
- `decodeToJsonc(wire: Uint8Array): string`：将 MetaMessage Wire 格式转换为 JSONC 格式
- `valueToJsonc(value: any, tag?: Tag): string`：将对象转换为 JSONC 格式
- `jsoncToValue<T>(jsonc: string, type: Constructor<T> | T): T`：将 JSONC 格式转换为对象
- `mm`：用于定义字段的装饰器
- `ValueType`：定义字段类型的枚举值。

### 2.2 对象编码

```ts
const person = {
  name: mm('Ed', { desc: '姓名' }),
  email: mm('Ed@gmail.com', { desc: '邮箱', type: 'email' }),
  score: mm(90, { desc: '成绩', type: 'uint8' }),
  age: mm(30, { desc: '年龄' }),
};

const wire = encodeFromValue(person);
console.log('wire', wire);
```

### 2.3 类实例编码

```ts
import { encodeFromValue, decodeToValue, mm, ValueType } from 'metamessage';

@mm({ desc: '用户' })
class User {
  @mm({ type: ValueType.I64, desc: '用户ID', nullable: false })
  id: bigint = 0n;

  @mm({ desc: '昵称' })
  name: string = '';

  @mm({ type: ValueType.U8, desc: '年龄' })
  age: number = 0;
}

const u = new User();
u.id = 666n;
u.name = 'abc';
u.age = 20;

const wire = encodeFromValue(u);
console.log('wire', wire);
```

> 如果使用装饰器，请确保 `tsconfig.json` 中启用 `experimentalDecorators`。

### 2.4 解码示例

```ts
const decoded = decodeToValue(wire, User);
console.log('decoded', decoded);
```

### 2.5 JSONC 示例

```ts
import { encodeFromJsonc, decodeToJsonc } from 'metamessage';

const jsonc = `
// mm: desc="用户"
{
  // mm: type=i64; desc="用户ID"
  "id": 666,

  // mm: desc="昵称"
  "name": "abc",

  // mm: type=u8
  "age": 20
}
`;

const wire = encodeFromJsonc(jsonc);
const jsoncString = decodeToJsonc(wire);
console.log('JSONC result:\n', jsoncString);
```

## 3. 运行测试

在 `mm-ts` 目录下运行：

```bash
npm test
```

额外命令：

```bash
npm run test:type
npm run build
npm run lint
npm run format
```

## 4. 常见问题

- 确认 `tsconfig.json` 配置正确
- 确保 `experimentalDecorators` 开启（如果使用装饰器）

## 5. 示例目录

示例代码位于 `examples/typescript/`：

- `basic/`：基础示例
- `jsonc-to-wire/`：JSONC 转 MetaMessage Wire
- `wire-to-jsonc/`：Wire 转 JSONC

## 6. 参考资源

- [MetaMessage GitHub](https://github.com/metamessage/metamessage)
- [npm package](https://www.npmjs.com/package/metamessage)
- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)
- [Jest 官方文档](https://jestjs.io/docs/getting-started)
