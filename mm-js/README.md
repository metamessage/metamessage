# JavaScript MetaMessage 库使用说明

## 1. 安装

### npm 依赖
使用 npm 安装：

```bash
npm install metamessage
```

### 版本要求
- Node.js 14 或更高版本

## 2. 基本使用

### 2.1 导入模块

```javascript
const { encode, decode } = require('metamessage');
```

### 2.2 编码示例

```javascript
const person = { name: "Ed", age: 30 };
const wire = encode(person);
console.log('Encoded:', bytesToHex(wire));
```

### 2.3 解码示例

```javascript
const decoded = decode(wire);
console.log('Decoded:', decoded.value);
```

### 2.4 JSONC 解析示例

```javascript
const { parseFromString, toString } = require('metamessage/jsonc');

const jsonc = `{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25
}`;

// 解析 JSONC
const node = parseFromString(jsonc);

// 转换为字符串
const jsoncString = toString(node);
```

## 3. 测试方法

### 3.1 运行现有测试

```bash
# 在 mm-js 目录下运行
npm test
```

### 3.2 测试框架
- Jest

### 3.3 测试覆盖范围
- 编码测试
- 解码测试
- JSONC 解析测试

## 4. 常见问题

### 4.1 依赖问题
- **问题**: npm 依赖安装失败
  **解决**: 检查网络连接，或使用 npm 镜像

### 4.2 运行时问题
- **问题**: 编码/解码失败
  **解决**: 检查数据类型是否正确

## 5. 示例代码

查看 `examples/javascript/` 目录下的示例代码：
- `basic/` - 基本使用示例

## 6. 相关资源

- [Node.js 文档](https://nodejs.org/docs/)
- [Jest 文档](https://jestjs.io/docs/getting-started)
