# MetaMessage 项目工作规范

## 项目概述
MetaMessage 是一个结构化数据交换协议，支持多种语言（C、C++、Go、JavaScript、TypeScript、Java、Kotlin、C#、Swift、PHP、Rust、Python等）。

## 目录结构
```
.
├── mm-c/           # C语言实现
├── mm-cpp/         # C++实现（使用STL）
├── mm-cs/          # C#实现
├── mm-java/        # Java实现
├── mm-js/          # JavaScript实现
├── mm-kt/          # Kotlin实现
├── mm-php/         # PHP实现
├── mm-py/          # Python实现
├── mm-rs/          # Rust实现
├── mm-swift/       # Swift实现
├── mm-ts/          # TypeScript实现
├── internal/       # Go内部包
│   └── mm/         # Go mm包
│       └── x64dbg_test.go  # x64dbg通信测试
├── src/            # x64dbg MCP插件源码
│   ├── MCPx64dbg.cpp  # C++ HTTP插件
│   └── x64dbg.py      # Python MCP服务器
├── examples/       # 各语言示例
└── internal/jsonc/ # JSONC解析器
```

## 重要规则

### 1. 不要删除任何文件
- **绝对不能**删除 go.mod、go.sum
- **绝对不能**删除 src/ 目录下的任何文件
- 复制文件时使用 `cp -r` 但不要覆盖已有模块文件

### 2. x64dbg相关工作
- x64dbg MCP文件在 src/ 目录
- CMakeLists.txt 用于编译x64dbg插件（需要x64dbg SDK）
- 不要创建额外的Go模块目录（会导致go mod冲突）

### 3. 测试要求
- Go测试: `go test ./internal/mm/... -v`
- C++测试: 需要CMake + MSVC，在Windows环境编译

### 4. Git工作流
- 提交前检查 `git status`
- 推送到: https://github.com/ddkwork/metamessage
- 使用 GitHub API 创建 PR

### 5. mm标签格式
```go
type MyStruct struct {
    Field1 string `mm:"type=string; allow_empty"`
    Field2 int    `mm:"type=int"`
    Field3 bool   `mm:"type=bool"`
}
```

## 常用命令
```bash
# Go测试
go test ./internal/mm/... -v

# 构建C++（需要Windows）
cmake -S . -B build
cmake --build build

# 查看状态
git status
git diff --stat HEAD
```
