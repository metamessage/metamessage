# MetaMessage

MetaMessage 是一种自描述、可约束、可示例化的结构化数据交换协议，支持无损数据交换。它适用于配置文件、API 交互和 AI 数据交换场景。

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

- Node.js 18+
- TypeScript 5+

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
@mm({ desc: '用户' })
class User {
  @mm({ type: ValueType.Int64, desc: '用户ID', nullable: false })
  id: bigint = 0n;

  @mm({ desc: '昵称' })
  name: string = '';

  @mm({ type: ValueType.Uint8, desc: '年龄' })
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

### 4.1 依赖安装失败

- 检查网络连接
- 使用国内 npm 镜像

### 4.2 TypeScript 编译错误

- 确认 `tsconfig.json` 配置正确
- 确保 `experimentalDecorators` 开启（如果使用装饰器）

### 4.3 运行时失败

- 检查输入数据类型是否正确
- 确认导入的 API 名称和参数匹配

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
