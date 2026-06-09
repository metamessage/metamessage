# Checklist

## 分析阶段
- [ ] 读取了 C# MetaMessage.cs 编码逻辑
- [ ] 读取了 Go encode.go 编码逻辑
- [ ] 读取了 C# Tag.cs 的 ToBytes()
- [ ] 读取了 Go tag.go 的 Bytes()
- [ ] 读取了 C# WireEncoder.cs 的 EncodeTaggedPayload
- [ ] 读取了 3 个 encode_diff 文件确认字节差异

## 修复阶段
- [ ] Tag.ToBytes() 序列化与 Go 完全一致
- [ ] IR 节点编码标签传播与 Go 完全一致
- [ ] WireEncoder 标签编码与 Go 完全一致

## 验证阶段
- [ ] 构建通过（`dotnet build` 成功）
- [ ] 所有 34 个编码测试通过（`encode_decode_test_cs.sh` 无编码失败）
- [ ] 所有 34 个解码测试仍然通过