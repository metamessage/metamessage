# Tasks

- [ ] T1: 读取关键源文件，理解 C# 和 Go 编码差异
  - 读取 C# MetaMessage.cs 编码逻辑（EncodeTreeValue、EncodeArrayTree、EncodeMapTree、EncodeScalarTree）
  - 读取 Go encode.go 编码逻辑
  - 读取 C# Tag.cs 的 ToBytes() 方法
  - 读取 Go tag.go 的 Bytes() 方法
  - 读取 C# WireEncoder.cs 的 EncodeTaggedPayload 方法
  - 读取 encode_diff 文件（03_tags_desc_tag.jsonc.cs_vs_go.encode_diff、03_tags_type_tags.jsonc.cs_vs_go.encode_diff、06_complex_config.jsonc.cs_vs_go.encode_diff）了解具体字节差异
  - 依赖: 无

- [ ] T2: 修复 Tag.ToBytes() 序列化对齐 Go
  - 对比 C# Tag.ToBytes() 与 Go Tag.Bytes() 的字节生成逻辑
  - 修复 tag 字段编码（key、length、value）以匹配 Go 输出
  - 处理 desc、type、child tags 等字段的序列化差异
  - 依赖: T1

- [ ] T3: 修复 IR 节点编码标签传播逻辑
  - 对比 C# MetaMessage.cs 与 Go encode.go 的标签传播逻辑
  - 修复 EncodeTreeValue/EncodeArrayTree/EncodeMapTree 中的 tag inheritance 和 propagation
  - 依赖: T1

- [ ] T4: 修复 WireEncoder 标签编码逻辑
  - 对比 C# WireEncoder.cs 与 Go wire format 编码器的标签编码
  - 修复 EncodeTaggedPayload 的前缀标记和长度编码
  - 依赖: T1

- [ ] T5: 构建并验证测试
  - 运行 `dotnet build tests/harness/csharp/harness.csproj --nologo` 确认构建通过
  - 运行 `bash tests/encode_decode_test_cs.sh` 确认所有编码测试通过
  - 依赖: T2, T3, T4

# Task Dependencies

- T1 无依赖（探索性任务）
- T2、T3、T4 依赖 T1
- T5 依赖 T2、T3、T4