# TypeScript MetaMessage (mm-ts) 实现计划

## 概述

基于 Go 版本的功能完整移植到 TypeScript，排除 CLI 功能，测试要完整全面，包括所有数据类型和 tag。

由于 TypeScript 没有 Go 的 struct tag 系统，实现将采用 `mm(value, {tag})` 模式来处理类型元数据。

## 项目结构

```
mm-ts/
├── package.json
├── tsconfig.json
├── src/
│   ├── mm/
│   │   ├── constants.ts      - Wire 格式常量
│   │   ├── buffer.ts         - 字节缓冲区
│   │   ├── encoder.ts        - 编码器
│   │   ├── decoder.ts        - 解码器
│   │   ├── types.ts          - 类型定义和 mm() API
│   │   ├── index.ts          - 公共 API
│   │   └── utils.ts          - 工具函数
│   └── jsonc/
│       ├── scanner.ts        - 词法分析器
│       ├── ast.ts            - AST 节点定义
│       ├── parser.ts         - 语法分析器
│       ├── tag.ts            - Tag 定义和解析
│       ├── binder.ts         - 对象绑定器
│       ├── printer.ts        - 输出格式化器
│       └── index.ts          - JSONC 模块导出
└── test/
    ├── mm/
    │   ├── encoder.test.ts
    │   └── decoder.test.ts
    └── jsonc/
        ├── scanner.test.ts
        ├── parser.test.ts
        └── tag.test.ts
```

## 关键技术考虑

### 1. BigInt 处理

TypeScript 使用 JavaScript BigInt 处理大整数：

* `int64` → `BigInt`

* `uint64` → `BigInt`

* 编码时使用 BigInt 直接转换

* 解码时返回 BigInt 值

### 2. Tag 系统设计

由于 TypeScript 没有原生 tag 系统，采用函数式 API：

```typescript
// 设计示例
mm(123)                           // 普通数值
mm(123, { type: 'int' })          // 指定类型
mm('hello', { type: 'str' })      // 字符串
mm([1, 2, 3], { type: 'array' }) // 数组
mm({ a: 1 }, { type: 'struct' })  // 对象
```

### 3. 类型定义

* 使用 TypeScript 接口和类型别名

* 提供完整的类型定义

* 支持泛型

### 4. 模块格式

* 使用 ES Modules

* 兼容 Node.js 和浏览器环境

* 生成类型声明文件 (.d.ts)

## 待完成任务

### 1. 项目初始化

* 创建 package.json

* 配置 tsconfig.json

* 配置 ESLint 和 Prettier

### 2. MM 模块实现

#### 2.1 constants.ts - Wire 格式常量

* 前缀枚举（Simple, PositiveInt, NegativeInt, Float, String, Bytes, Container, Tag）

* 简单值枚举（null, true, false 等）

* 长度计算常量

* 前缀/后缀掩码

#### 2.2 buffer.ts - 字节缓冲区

* `MMBuffer` 类

* 读写方法：`writeUint8`, `writeInt8`, `writeUint16LE`, `writeInt16LE` 等

* `writeFloat32`, `writeFloat64`

* `readUint8`, `readInt8`, `readUint16LE`, `readInt16LE` 等

* `readFloat32`, `readFloat64`

* `readBytes`, `writeBytes`

#### 2.3 types.ts - 类型定义和 mm() API

```typescript
// 类型接口
export interface MMOptions {
  type?: string;
  desc?: string;
  nullable?: boolean;
  default?: any;
  min?: string;
  max?: string;
  size?: number;
  enum?: string;
  pattern?: string;
  location?: number;
  version?: number;
  mime?: string;
  [key: string]: any;
}

export interface MMValue<T = any> {
  value: T;
  options: MMOptions;
}

// 类型包装器
export function mm<T>(value: T, options: MMOptions = {}): MMValue<T> {
  return { value, options };
}

// 类型工厂函数
mm.int = (v: number | BigInt) => mm(v, { type: 'int' });
mm.str = (v: string) => mm(v, { type: 'str' });
mm.bool = (v: boolean) => mm(v, { type: 'bool' });
mm.float = (v: number) => mm(v, { type: 'float' });
mm.bytes = (v: Uint8Array | number[]) => mm(v, { type: 'bytes' });
mm.array = <T>(v: T[], options?: MMOptions) => mm(v, { ...options, type: 'array' });
mm.struct = <T>(v: T, options?: MMOptions) => mm(v, { ...options, type: 'struct' });
mm.bigint = (v: BigInt) => mm(v, { type: 'bigint' });
```

#### 2.4 encoder.ts - 编码器

* `MMEncoder` 类

* `encodeBool`, `encodeNil`

* `encodeInt` (使用 BigInt)

* `encodeUInt` (使用 BigInt)

* `encodeFloat`, `encodeDouble`

* `encodeString`

* `encodeBytes`

* `encodeArray`

* `encodeObject` / `encodeStruct`

* 处理 int64/uint64 使用 BigInt

#### 2.5 decoder.ts - 解码器

* `MMDecoder` 类

* `decode()` - 返回 `DecodedValue`

* `decodeBool`, `decodeNil`

* `decodeInt` - 返回 BigInt

* `decodeUInt` - 返回 BigInt

* `decodeFloat`, `decodeDouble`

* `decodeString`

* `decodeBytes`

* `decodeArray`

* `decodeObject`

#### 2.6 index.ts - 公共 API

```typescript
export function encode(value: any): Uint8Array;
export function decode(data: Uint8Array | ArrayBuffer | number[]): DecodedValue;
export function toJSONC(data: Uint8Array | ArrayBuffer | number[]): string;
export function fromJSONC(jsoncString: string): Uint8Array;
export { mm, MMValue, MMOptions } from './types';
export { MMEncoder, MMDecoder, MMBuffer } from './encoder';
export { DecodedValue } from './decoder';
export * as constants from './constants';
```

### 3. JSONC 模块实现

#### 3.1 scanner.ts - 词法分析器

* Token 类型（LCURLY, RCURLY, LBRACKET, RBRACKET, COLON, COMMA, STRING, NUMBER, TRUE, FALSE, NULL, LINECOMMENT, BLOCKCOMMENT）

* `JSONCScanner` 类

* `nextToken()` 方法

* 追踪行号和列位置

* 正确处理 `//` 行注释和 `/* */` 块注释

#### 3.2 ast.ts - AST 节点定义

```typescript
export interface JSONCNode {
  getType(): string;
  getTag(): JSONCTag | null;
  getPath(): string;
}

export class JSONCValue implements JSONCNode { ... }
export class JSONCObject implements JSONCNode { ... }
export class JSONCArray implements JSONCNode { ... }
export class JSONCDoc implements JSONCNode { ... }
```

#### 3.3 tag.ts - Tag 定义和解析

```typescript
export class JSONCTag {
  constructor()
  inherit(tag: JSONCTag): void
  toString(): string
}

export function parseMMTag(tagStr: string): JSONCTag;
```

#### 3.4 parser.ts - 语法分析器

* `JSONCParser` 类

* `parse()` 方法

* 递归下降解析

* 处理注释并关联到节点

* 创建完整的 AST

#### 3.5 printer.ts - 输出格式化器

* `JSONCPrinter` 类

* `print(node: JSONCNode): string` - 格式化输出

* `printCompact(node: JSONCNode): string` - 紧凑输出

* 处理所有节点类型

#### 3.6 binder.ts - 对象绑定器

* `JSONCBinder` 类

* `bind<T>(node: JSONCNode, Type: new () => T): T` 方法

* 将 AST 绑定到 TypeScript 对象

### 4. 测试模块

#### 4.1 MM 编码器/解码器测试

* `test/mm/encoder.test.ts`

* `test/mm/decoder.test.ts`

测试覆盖：

* Bool (true, false)

* Nil / Null

* Int (BigInt): 0, 23, 24, 123456, -7890

* Int8, Int16, Int32, Int64 (BigInt) 边界值

* UInt (BigInt): 0, 987654, 255, 65535, 4294967295, 18446744073709551615

* UInt8, UInt16, UInt32, UInt64 (BigInt) 边界值

* Float: 0.0, 3.14, -3.14, 极值

* Double: 0.0, 3.14159265359, 极值

* String: "", "hello", "hello world", 长字符串 (300+ chars)

* Bytes: empty, \[0x01, 0x02, 0x03], 长 bytes (300+)

* Array: \[Bool], \[Int], \[String], \[Float], \[Double], 空数组, 大数组 (1000 items)

* Object/Struct: 简单对象, 嵌套对象

#### 4.2 JSONC 测试

* `test/jsonc/scanner.test.ts` - Token 识别, 注释, 位置追踪

* `test/jsonc/parser.test.ts` - 有效 JSONC, 带注释的 JSONC, 嵌套结构

* `test/jsonc/tag.test.ts` - Tag 解析, 复杂 tag, 继承

## 参考实现

* Go 版本: `/Users/lizongying/IdeaProjects/meta-message/internal/mm/`

* JavaScript 版本: `/Users/lizongying/IdeaProjects/meta-message/mm-js/`

* Rust 版本: `/Users/lizongying/IdeaProjects/meta-message/mm-rs/`

* Swift 版本: `/Users/lizongying/IdeaProjects/meta-message/mm-swift/`

## 实现顺序

1. 项目初始化（package.json, tsconfig.json, 目录结构）
2. MM/constants.ts - 常量定义
3. MM/buffer.ts - 缓冲区
4. MM/types.ts - 类型定义和 mm() API
5. MM/encoder.ts - 编码器
6. MM/decoder.ts - 解码器
7. MM/index.ts - 公共 API
8. JSONC/scanner.ts - 词法分析器
9. JSONC/ast.ts - AST 节点
10. JSONC/tag.ts - Tag 解析
11. JSONC/parser.ts - 语法分析器
12. JSONC/printer.ts - 格式化输出
13. JSONC/binder.ts - 对象绑定
14. 测试文件实现
15. 运行测试验证

## 关键注意事项

* TypeScript BigInt 没有无符号上限，但 wire 格式需要处理 uint64

* 编码时 int64/uint64 都使用 BigInt

* 解码时返回 BigInt，调用者需自行转换

* 使用 ES Modules (import/export)

* 错误处理使用 Error 类型

* 缓冲区使用 Uint8Array

* 测试使用 Jest 框架

* 生成完整的 TypeScript 类型定义

* 确保类型安全性和类型推断

## 技术栈

* TypeScript 5.0+

* Jest 29.7+

* ESLint + Prettier

* ts-jest (TypeScript 测试支持)

* Node.js 18.0+

## 预期输出

* 完整的 TypeScript 实现

* 类型声明文件 (.d.ts)

* 完整的测试覆盖

* 与 Go 版本功能一致的 API

* 支持所有数据类型和 tag 系统

