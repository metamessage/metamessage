# TypeScript MetaMessage 库使用说明

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
import { encode, decode, mm, bind } from 'metamessage';
```

or

```typescript
const { encode, decode, mm, bind } = require('metamessage');
```

### 2.2 编码示例

```typescript
// 使用 mm() 函数创建带类型信息的对象
const person = {
  name: mm("Ed", { desc: "" }),
  email: mm("Ed@gmail.com", { desc: "", type: "email" }),
  score: mm(90, { desc: "", type: "uint8" }),
  age: mm(30, { desc: "" })
};

const wire = encode(person);
console.log('Encoded:', bytesToHex(wire));
```

### 2.3 解码示例

```typescript
const decoded = decode(wire);
console.log('Decoded:', decoded.value);

// 绑定到模板
const bound = bind(wire, person);
console.log('Bound:', bound);
```

### 2.4 JSONC 解析示例

```typescript
import { toJSONC, fromJSONC, mm } from 'metamessage';

const person = {
  name: mm("Ed", { desc: "" }),
  age: mm(30, { desc: "" })
};

// 编码为 Wire 格式，再转换为 JSONC
const wire = encode(person);
const jsoncString = toJSONC(wire);
console.log('JSONC:', jsoncString);

// 从 JSONC 转换回 Wire 格式，再解码
const wireFromJsonc = fromJSONC(jsoncString);
const decoded = decode(wireFromJsonc);
console.log('Decoded from JSONC:', decoded.value);
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
