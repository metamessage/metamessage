namespace MetaMessage.Core;

public static class TagFieldParser
{
    public static MmTag Parse(byte[] data)
    {
        MmTag tag = MmTag.Empty();
        int offset = 0;

        while (offset < data.Length)
        {
            if (offset >= data.Length)
                break;

            byte fieldType = data[offset++];

            switch (fieldType)
            {
                case 1: // Name
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int nameLen = data[offset++];
                    if (offset + nameLen > data.Length)
                        throw new MmDecodeException("Tag name data overflow");
                    tag.Name = System.Text.Encoding.UTF8.GetString(data, offset, nameLen);
                    offset += nameLen;
                    break;

                case 2: // Desc
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int descLen = data[offset++];
                    if (offset + descLen > data.Length)
                        throw new MmDecodeException("Tag desc data overflow");
                    tag.Desc = System.Text.Encoding.UTF8.GetString(data, offset, descLen);
                    offset += descLen;
                    break;

                case 3: // Type
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.Type = (ValueType)data[offset++];
                    break;

                case 4: // ChildType
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.ChildType = (ValueType)data[offset++];
                    break;

                case 5: // Nullable
                    tag.Nullable = true;
                    break;

                case 6: // Raw
                    tag.Raw = true;
                    break;

                case 7: // AllowEmpty
                    tag.AllowEmpty = true;
                    break;

                case 8: // Unique
                    tag.Unique = true;
                    break;

                case 9: // DefaultValue
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int defaultValueLen = data[offset++];
                    if (offset + defaultValueLen > data.Length)
                        throw new MmDecodeException("Tag default value data overflow");
                    tag.DefaultValue = System.Text.Encoding.UTF8.GetString(data, offset, defaultValueLen);
                    offset += defaultValueLen;
                    break;

                case 10: // Enum
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int enumLen = data[offset++];
                    if (offset + enumLen > data.Length)
                        throw new MmDecodeException("Tag enum data overflow");
                    tag.Enum = System.Text.Encoding.UTF8.GetString(data, offset, enumLen);
                    offset += enumLen;
                    break;

                case 11: // LocationHours
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    var locationLen = data[offset++];
                    if (offset + locationLen > data.Length)
                        throw new MmDecodeException("Tag location data overflow");
                    var locationStr = System.Text.Encoding.UTF8.GetString(data, offset, locationLen);
                    tag.LocationHours = int.Parse(locationStr);
                    offset += locationLen;
                    break;

                case 12: // Version
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.Version = data[offset++];
                    break;

                case 13: // Mime
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int mimeLen = data[offset++];
                    if (offset + mimeLen > data.Length)
                        throw new MmDecodeException("Tag mime data overflow");
                    tag.Mime = System.Text.Encoding.UTF8.GetString(data, offset, mimeLen);
                    offset += mimeLen;
                    break;

                case 14: // ChildDesc
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childDescLen = data[offset++];
                    if (offset + childDescLen > data.Length)
                        throw new MmDecodeException("Tag child desc data overflow");
                    tag.ChildDesc = System.Text.Encoding.UTF8.GetString(data, offset, childDescLen);
                    offset += childDescLen;
                    break;

                case 15: // ChildNullable
                    tag.ChildNullable = true;
                    break;

                case 28: // ChildAllowEmpty
                    tag.ChildAllowEmpty = true;
                    break;

                case 29: // ChildUnique
                    tag.ChildUnique = true;
                    break;

                case 16: // ChildEnum
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childEnumLen = data[offset++];
                    if (offset + childEnumLen > data.Length)
                        throw new MmDecodeException("Tag child enum data overflow");
                    tag.ChildEnum = System.Text.Encoding.UTF8.GetString(data, offset, childEnumLen);
                    offset += childEnumLen;
                    break;

                case 17: // Size
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.Size = data[offset++];
                    break;

                case 18: // Pattern
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int patternLen = data[offset++];
                    if (offset + patternLen > data.Length)
                        throw new MmDecodeException("Tag pattern data overflow");
                    tag.Pattern = System.Text.Encoding.UTF8.GetString(data, offset, patternLen);
                    offset += patternLen;
                    break;

                case 19: // ChildPattern
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childPatternLen = data[offset++];
                    if (offset + childPatternLen > data.Length)
                        throw new MmDecodeException("Tag child pattern data overflow");
                    tag.ChildPattern = System.Text.Encoding.UTF8.GetString(data, offset, childPatternLen);
                    offset += childPatternLen;
                    break;

                case 20: // ChildSize
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.ChildSize = data[offset++];
                    break;

                case 21: // Min
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int minLen = data[offset++];
                    if (offset + minLen > data.Length)
                        throw new MmDecodeException("Tag min data overflow");
                    tag.Min = System.Text.Encoding.UTF8.GetString(data, offset, minLen);
                    offset += minLen;
                    break;

                case 22: // Max
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int maxLen = data[offset++];
                    if (offset + maxLen > data.Length)
                        throw new MmDecodeException("Tag max data overflow");
                    tag.Max = System.Text.Encoding.UTF8.GetString(data, offset, maxLen);
                    offset += maxLen;
                    break;

                case 23: // ChildMin
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childMinLen = data[offset++];
                    if (offset + childMinLen > data.Length)
                        throw new MmDecodeException("Tag child min data overflow");
                    tag.ChildMin = System.Text.Encoding.UTF8.GetString(data, offset, childMinLen);
                    offset += childMinLen;
                    break;

                case 24: // ChildMax
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childMaxLen = data[offset++];
                    if (offset + childMaxLen > data.Length)
                        throw new MmDecodeException("Tag child max data overflow");
                    tag.ChildMax = System.Text.Encoding.UTF8.GetString(data, offset, childMaxLen);
                    offset += childMaxLen;
                    break;

                case 25: // ChildLocation
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.ChildLocation = data[offset++];
                    break;

                case 26: // ChildVersion
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.ChildVersion = data[offset++];
                    break;

                case 27: // ChildMime
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childMimeLen = data[offset++];
                    if (offset + childMimeLen > data.Length)
                        throw new MmDecodeException("Tag child mime data overflow");
                    tag.ChildMime = System.Text.Encoding.UTF8.GetString(data, offset, childMimeLen);
                    offset += childMimeLen;
                    break;

                default:
                    // 忽略未知的标签字段
                    break;
            }
        }

        return tag;
    }
}