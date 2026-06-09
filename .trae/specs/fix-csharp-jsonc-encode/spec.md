# 修复 C# JSONC 编码测试失败 Spec

## Why

C# JSONC 编码测试 34 个中有 20 个因字节不匹配而失败（所有 34 个解码测试均通过）。C# 实现的 wire format 编码与 Go 参考实现不一致，问题可能涉及标签传播、Tag 序列化以及 IR 节点编码逻辑。

## What Changes

1. **Tag 序列化对齐** — 修复 C# `Tag.ToBytes()` 以生成与 Go `Tag.Bytes()` 完全相同的字节表示，包括 desc、type、child tags 等字段的编码对齐
2. **IR 节点编码标签传播对齐** — 修复 `MetaMessage.cs` 中的 `EncodeTreeValue`/`EncodeArrayTree`/`EncodeMapTree` 方法，确保 tag inheritance 和 propagation 逻辑与 Go `encode.go` 一致
3. **WireEncoder 标签负载编码对齐** — 修复 `WireEncoder.cs` 中的 `EncodeTaggedPayload` 方法，确保 tag 编码前缀、长度编码与 Go 一致
4. **测试验证** — 构建通过后运行 `encode_decode_test_cs.sh` 确认所有编码测试通过

## Impact

- Affected code: `mm-cs/src/MetaMessage/Ir/Tag.cs`、`mm-cs/src/MetaMessage/Core/MetaMessage.cs`、`mm-cs/src/MetaMessage/Core/WireEncoder.cs`
- Affected specs: C# JSONC 编码层

## ADDED Requirements

### Requirement: Tag.ToBytes() 与 Go Tag.Bytes() 字节完全一致
C# 的 `Tag.ToBytes()` 方法 SHALL 生成与 Go 的 `Tag.Bytes()` 完全相同的字节序列。

#### Scenario: 包含 desc 字段的 Tag
- **WHEN** Tag 包含 Description 字段
- **THEN** `ToBytes()` 的字节输出 SHALL 与 Go 的 `Bytes()` 输出完全一致

#### Scenario: 包含 type 字段的 Tag
- **WHEN** Tag 包含 Type 字段
- **THEN** `ToBytes()` 的字节输出 SHALL 与 Go 的 `Bytes()` 输出完全一致

#### Scenario: 包含 child tags 的 Tag
- **WHEN** Tag 包含子标签
- **THEN** `ToBytes()` 的字节输出 SHALL 与 Go 的 `Bytes()` 输出完全一致

### Requirement: IR 节点编码标签传播
C# 的 `EncodeTreeValue`/`EncodeArrayTree`/`EncodeMapTree` SHALL 在编码时正确传播 inherited tag，与 Go 的 `encodeNodeValue`/`encodeArray`/`encodeMap` 行为一致。

### Requirement: WireEncoder TaggedPayload 编码
C# 的 `EncodeTaggedPayload` SHALL 在编码 tag + payload 时使用与 Go 一致的前缀标记和长度编码方式。

## MODIFIED Requirements

### Requirement: C# MetaMessage.cs 编码逻辑
`MetaMessage.cs` 中的 `EncodeTreeValue`、`EncodeArrayTree`、`EncodeMapTree`、`EncodeScalarTree` 方法 SHALL 修改以匹配 Go `encode.go` 的标签处理逻辑。

### Requirement: C# WireEncoder.cs
`WireEncoder.cs` SHALL 修改以匹配 Go wire format 编码器的标签编码行为。

## REMOVED Requirements

无