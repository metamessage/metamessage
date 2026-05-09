# MM-Kotlin 深度对齐检查清单

## AST 结构对齐检查点 ✅

- [x] Node 接口与 Golang 一致
- [x] Value 类字段与 Golang 一致
- [x] Object 类字段与 Golang 一致
- [x] Array 类字段与 Golang 一致
- [x] Field 类字段与 Golang 一致
- [x] Tag 类字段与 Golang 一致
- [x] ValueType 枚举与 Golang 一致
- [x] NodeType 枚举与 Golang 一致
- [x] Doc 类与 Golang 一致

## Tag 验证逻辑对齐检查点 ✅

- [x] validateString 方法与 Golang 一致
- [x] validateInt8/16/32/64 方法与 Golang 一致
- [x] validateUint/8/16/32/64 方法与 Golang 一致
- [x] validateFloat32/64 方法与 Golang 一致
- [x] validateBool 方法与 Golang 一致
- [x] validateBytes 方法与 Golang 一致
- [x] validateUUID 方法与 Golang 一致
- [x] validateEmail 方法与 Golang 一致
- [x] validateEnum 方法与 Golang 一致
- [x] validateDateTime/Date/Time 方法与 Golang 一致
- [x] validateDecimal 方法与 Golang 一致
- [x] validateIP 方法与 Golang 一致
- [x] validateURL 方法与 Golang 一致
- [x] validateArray/Slice 方法与 Golang 一致
- [x] validateStruct 方法与 Golang 一致
- [x] validateMap 方法与 Golang 一致
- [x] validateImage 方法与 Golang 一致

## JSONC 解析器对齐检查点 ✅

- [x] 字符串解析与 Golang 一致
- [x] 数字解析与 Golang 一致
- [x] 布尔值解析与 Golang 一致
- [x] 对象解析与 Golang 一致
- [x] 数组解析与 Golang 一致
- [x] 注释解析与 Golang 一致
- [x] Tag 解析与 Golang 一致
- [x] 深度限制与 Golang 一致
- [x] 错误处理与 Golang 一致

## JSONC 打印机对齐检查点 ✅

- [x] 字符串打印与 Golang 一致
- [x] 数字打印与 Golang 一致
- [x] 布尔值打印与 Golang 一致
- [x] 对象打印与 Golang 一致
- [x] 数组打印与 Golang 一致
- [x] 注释打印与 Golang 一致
- [x] Tag 输出格式与 Golang 一致
- [x] 类型感知打印与 Golang 一致

## MM 编解码对齐检查点 ✅

- [x] 整数编码与 Golang 一致
- [x] 浮点数编码与 Golang 一致
- [x] 布尔值编码与 Golang 一致
- [x] 字符串编码与 Golang 一致
- [x] 字节数组编码与 Golang 一致
- [x] 对象编码与 Golang 一致
- [x] 数组编码与 Golang 一致
- [x] Tag 编码与 Golang 一致
- [x] 整数解码与 Golang 一致
- [x] 浮点数解码与 Golang 一致
- [x] 布尔值解码与 Golang 一致
- [x] 字符串解码与 Golang 一致
- [x] 字节数组解码与 Golang 一致
- [x] 对象解码与 Golang 一致
- [x] 数组解码与 Golang 一致
- [x] Tag 解码与 Golang 一致

## 测试验证检查点 ✅

- [x] AST 结构测试
- [x] Tag 验证测试
- [x] JSONC 解析测试
- [x] JSONC 打印测试
- [x] MM 编码测试
- [x] MM 解码测试
- [x] 跨语言兼容性测试