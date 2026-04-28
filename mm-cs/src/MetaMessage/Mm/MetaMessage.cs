namespace MetaMessage.Mm;

public static class MetaMessage
{
    public static byte[] Encode(object obj)
    {
        return ReflectMmEncoder.Encode(obj);
    }

    public static void Decode(byte[] data, object target)
    {
        ReflectMmBinder.Bind(data, target);
    }

    public static T Decode<T>(byte[] data) where T : new()
    {
        var obj = new T();
        ReflectMmBinder.Bind(data, obj);
        return obj;
    }

    public static IMmTree DecodeToTree(byte[] data)
    {
        var decoder = new WireDecoder(data);
        return decoder.Decode();
    }

    public static byte[] EncodeTree(IMmTree tree)
    {
        // 实现从 MmTree 编码为二进制数据
        // 这里需要实现一个编码器，将 MmTree 转换为二进制格式
        throw new NotImplementedException("EncodeTree not implemented yet");
    }

    public static ValidationResult Validate(object value, MmTag tag)
    {
        return Validator.Validate(value, tag);
    }
}