using ValueType = MetaMessage.Ir.ValueType;
using MetaMessage.Ir;

namespace MetaMessage.Core;

public static class TagFieldParser
{
    public static Tag Parse(byte[] data)
    {
        Tag tag = Tag.Empty();
        int offset = 0;

        while (offset < data.Length)
        {
            if (offset >= data.Length)
                break;

            byte header = data[offset++];

            var key = (TagKey)(header & 0xF8);
            int lenInfo = header & 0x07;

            switch (key)
            {
                case TagKey.KExample:
                    tag.Example = true;
                    break;

                case TagKey.KIsNull:
                    tag.IsNull = true;
                    tag.Nullable = true;
                    break;

                case TagKey.KNullable:
                    tag.Nullable = true;
                    break;

                case TagKey.KDeprecated:
                    tag.Deprecated = true;
                    break;

                case TagKey.KAllowEmpty:
                    tag.AllowEmpty = true;
                    break;

                case TagKey.KUnique:
                    tag.Unique = true;
                    break;

                case TagKey.KChildNullable:
                    tag.ChildNullable = true;
                    break;

                case TagKey.KChildAllowEmpty:
                    tag.ChildAllowEmpty = true;
                    break;

                case TagKey.KChildUnique:
                    tag.ChildUnique = true;
                    break;

                case TagKey.KDesc:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag desc data overflow");
                        tag.Desc = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KType:
                    {
                        if (offset >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.Type = (ValueType)data[offset++];
                        break;
                    }

                case TagKey.KDefault:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag default value data overflow");
                        tag.DefaultVal = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KMin:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag min data overflow");
                        tag.Min = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KMax:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag max data overflow");
                        tag.Max = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KSize:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.Size = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KEnum:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag enum data overflow");
                        tag.Enums = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        tag.Type = ValueType.Enums;
                        offset += strLen;
                        break;
                    }

                case TagKey.KPattern:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag pattern data overflow");
                        tag.Pattern = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KLocation:
                    {
                        int strLen = lenInfo;
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag location data overflow");
                        tag.Location = int.Parse(System.Text.Encoding.UTF8.GetString(data, offset, strLen));
                        offset += strLen;
                        break;
                    }

                case TagKey.KVersion:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.Version = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KMime:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.Mime = DecodeU64(data, ref offset, byteCount).ToString();
                        tag.Type = ValueType.Media;
                        break;
                    }

                case TagKey.KChildDesc:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child desc data overflow");
                        tag.ChildDesc = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildType:
                    {
                        if (offset >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.ChildType = (ValueType)data[offset++];
                        break;
                    }

                case TagKey.KChildDefaultVal:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child default value data overflow");
                        tag.ChildDefaultVal = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildMin:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child min data overflow");
                        tag.ChildMin = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildMax:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child max data overflow");
                        tag.ChildMax = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildSize:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.ChildSize = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KChildEnums:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child enum data overflow");
                        tag.ChildEnums = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildPattern:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child pattern data overflow");
                        tag.ChildPattern = System.Text.Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildLocation:
                    {
                        int strLen = lenInfo;
                        if (offset + strLen > data.Length)
                            throw new MmDecodeException("Tag child location data overflow");
                        tag.ChildLocation = int.Parse(System.Text.Encoding.UTF8.GetString(data, offset, strLen));
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildVersion:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.ChildVersion = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KChildMime:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.ChildMime = DecodeU64(data, ref offset, byteCount).ToString();
                        tag.ChildType = ValueType.Media;
                        break;
                    }

                case TagKey.KMore:
                    {
                        int byteCount = lenInfo + 1;
                        if (offset + byteCount > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        tag.More = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                default:
                    if (lenInfo <= 5)
                    {
                        if (offset + lenInfo > data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        offset += lenInfo;
                    }
                    else if (lenInfo == 6)
                    {
                        if (offset >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        int skipLen = data[offset++];
                        if (offset + skipLen > data.Length)
                            throw new MmDecodeException("Tag data overflow");
                        offset += skipLen;
                    }
                    else if (lenInfo == 7)
                    {
                        if (offset + 1 >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        int skipLen = (data[offset] << 8) | data[offset + 1];
                        offset += 2;
                        if (offset + skipLen > data.Length)
                            throw new MmDecodeException("Tag data overflow");
                        offset += skipLen;
                    }
                    break;
            }
        }

        return tag;
    }

    private static int DecodeStringLength(byte[] data, ref int offset, int lenInfo)
    {
        if (lenInfo <= 5)
        {
            return lenInfo;
        }

        if (lenInfo == 6)
        {
            if (offset >= data.Length)
                throw new MmDecodeException("Unexpected end of tag data");
            int len = data[offset];
            offset++;
            return len;
        }

        if (offset + 1 >= data.Length)
            throw new MmDecodeException("Unexpected end of tag data");
        int lenHi = data[offset];
        offset++;
        int lenLo = data[offset];
        offset++;
        return (lenHi << 8) | lenLo;
    }

    private static int DecodeStringLengthSimple(byte[] data, ref int offset, int lenInfo)
    {
        if (lenInfo < 7)
        {
            return lenInfo;
        }

        if (offset >= data.Length)
            throw new MmDecodeException("Unexpected end of tag data");
        int len = data[offset];
        offset++;
        return len;
    }

    private static ulong DecodeU64(byte[] data, ref int offset, int byteCount)
    {
        ulong result = 0;
        for (int i = 0; i < byteCount; i++)
        {
            if (offset >= data.Length)
                throw new MmDecodeException("Unexpected end of tag data");
            result = (result << 8) | data[offset];
            offset++;
        }
        return result;
    }
}