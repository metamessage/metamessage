# Tag Min/Max Range Validation - Verification Checklist

## TypeScript 验证
- [x] uint16 类型，tag.min 为 70000 时验证失败
- [x] int8 类型，tag.min 为 -130 时验证失败
- [x] uint8 类型，tag.max 为 260 时验证失败
- [x] int16 类型的边界值验证
- [x] int32 类型的边界值验证
- [x] int64 类型的边界值验证
- [x] uint 类型的边界值验证
- [x] uint32 类型的边界值验证
- [x] uint64 类型的边界值验证
- [x] float32 类型的边界值验证
- [x] float64 类型的边界值验证

## Kotlin 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## Python 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## JavaScript 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## Java 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## C# 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## Rust 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## Swift 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## PHP 验证
- [x] 各整数类型的 tag.min/tag.max 范围验证
- [x] 各浮点数类型的 tag.min/tag.max 范围验证

## 跨语言一致性验证
- [x] 所有语言的 tag.min/tag.max 范围验证结果一致
- [x] 所有语言的错误消息一致
