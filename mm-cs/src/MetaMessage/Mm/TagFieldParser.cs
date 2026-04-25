namespace MetaMessage.Mm;

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

                case 10: // EnumValues
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int enumCount = data[offset++];
                    tag.EnumValues = new List<string>();
                    for (int i = 0; i < enumCount; i++)
                    {
                        if (offset >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        int enumValueLen = data[offset++];
                        if (offset + enumValueLen > data.Length)
                            throw new MmDecodeException("Tag enum value data overflow");
                        tag.EnumValues.Add(System.Text.Encoding.UTF8.GetString(data, offset, enumValueLen));
                        offset += enumValueLen;
                    }
                    break;

                case 11: // LocationHours
                    if (offset + 3 >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    tag.LocationHours = (data[offset] << 24) | (data[offset + 1] << 16) | (data[offset + 2] << 8) | data[offset + 3];
                    offset += 4;
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

                case 16: // ChildEnum
                    if (offset >= data.Length)
                        throw new MmDecodeException("Unexpected end of tag data");
                    int childEnumCount = data[offset++];
                    tag.ChildEnum = new List<string>();
                    for (int i = 0; i < childEnumCount; i++)
                    {
                        if (offset >= data.Length)
                            throw new MmDecodeException("Unexpected end of tag data");
                        int childEnumValueLen = data[offset++];
                        if (offset + childEnumValueLen > data.Length)
                            throw new MmDecodeException("Tag child enum value data overflow");
                        tag.ChildEnum.Add(System.Text.Encoding.UTF8.GetString(data, offset, childEnumValueLen));
                        offset += childEnumValueLen;
                    }
                    break;

                default:
                    // 忽略未知的标签字段
                    break;
            }
        }

        return tag;
    }
}