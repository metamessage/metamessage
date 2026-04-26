namespace MetaMessage.Mm;

public static class ReflectMmEncoder
{
    public static byte[] Encode(object root)
    {
        var encoder = new WireEncoder();
        EncodeValue(encoder, root, RootTagForClass(root.GetType()));
        return encoder.ToByteArray();
    }

    private static MmTag RootTagForClass(Type className)
    {
        var attributes = className.GetCustomAttributes(typeof(MM), false);
        MmTag tag;
        if (attributes.Length > 0)
        {
            var mmAttribute = (MM)attributes[0];
            tag = MmTagExtensions.FromAttribute(mmAttribute);
        }
        else
        {
            tag = MmTag.Empty();
        }

        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.STRUCT;
        }
        if (string.IsNullOrEmpty(tag.Name))
        {
            tag.Name = CamelToSnake.Convert(className.Name);
        }
        return tag;
    }

    private static void EncodeValue(WireEncoder encoder, object? value, MmTag tag)
    {
        var workTag = tag.Copy();
        if (value == null)
        {
            if (!workTag.Nullable && !workTag.IsNull)
            {
                throw new Exception("Null for non-nullable");
            }
            workTag.IsNull = true;
        }

        if (workTag.Type == ValueType.STRUCT)
        {
            if (value == null)
            {
                throw new Exception("Null struct");
            }
            EncodeObject(encoder, value, workTag);
            return;
        }

        if (workTag.Type == ValueType.SLICE || workTag.Type == ValueType.ARRAY)
        {
            EncodeList(encoder, value, workTag);
            return;
        }

        if (workTag.Type == ValueType.UNKNOWN)
        {
            workTag.Type = InferTypeFromValue(value);
        }

        var payload = new WireEncoder();
        EncodeScalarPayload(payload, value, workTag);
        encoder.EncodeTaggedPayload(payload.ToByteArray(), workTag.ToBytes());
    }

    private static ValueType InferTypeFromValue(object? value)
    {
        if (value == null)
        {
            return ValueType.STRING;
        }

        Type type = value.GetType();
        return TypeInference.ValueTypeForType(type);
    }

    private static void EncodeObject(WireEncoder encoder, object obj, MmTag objTag)
    {
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        var properties = obj.GetType().GetProperties();
        foreach (var property in properties)
        {
            if (!property.CanRead)
                continue;

            var attributes = property.GetCustomAttributes(typeof(MM), false);
            MM? mmAttribute = null;
            if (attributes.Length > 0)
            {
                mmAttribute = (MM)attributes[0];
                if (mmAttribute.Name == "-")
                    continue;
            }

            var fieldType = InferFieldType(property, mmAttribute);
            var fieldValue = property.GetValue(obj);
            var fieldKey = GetFieldKey(property, fieldType, mmAttribute);

            tmp.Reset();
            EncodeValue(tmp, fieldValue, fieldType);
            valsPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            tmp.EncodeString(fieldKey);
            keysPacked.WriteAll(tmp.ToByteArray());
        }

        tmp.Reset();
        tmp.EncodeArrayPayload(keysPacked.ToArray());
        var mapBody = new GrowableByteBuf();
        mapBody.WriteAll(tmp.ToByteArray());
        mapBody.WriteAll(valsPacked.ToArray());

        tmp.Reset();
        tmp.EncodeObjectPayload(mapBody.ToArray());
        encoder.EncodeTaggedPayload(tmp.ToByteArray(), objTag.ToBytes());
    }

    private static MmTag InferFieldType(System.Reflection.PropertyInfo property, MM? attribute)
    {
        var fieldType = MmTag.Empty();
        fieldType.Type = TypeInference.ValueTypeForType(property.PropertyType);

        if (attribute != null)
        {
            var attrTag = MmTagExtensions.FromAttribute(attribute);
            if (attrTag.Type != ValueType.UNKNOWN)
            {
                fieldType.Type = attrTag.Type;
            }
            if (attrTag.ChildType != ValueType.UNKNOWN)
            {
                fieldType.ChildType = attrTag.ChildType;
            }
            if (!string.IsNullOrEmpty(attrTag.Name))
            {
                fieldType.Name = attrTag.Name;
            }
            MergeAnnotations(fieldType, attrTag);
        }

        return fieldType;
    }

    private static void MergeAnnotations(MmTag dst, MmTag src)
    {
        if (!string.IsNullOrEmpty(src.Desc))
            dst.Desc = src.Desc;
        dst.Nullable |= src.Nullable;
        dst.Raw |= src.Raw;
        dst.AllowEmpty |= src.AllowEmpty;
        dst.Unique |= src.Unique;
        if (!string.IsNullOrEmpty(src.DefaultValue))
            dst.DefaultValue = src.DefaultValue;
        if (src.EnumValues.Count > 0)
        {
            dst.EnumValues = src.EnumValues;
            dst.Type = ValueType.ENUM;
        }
        dst.LocationHours = src.LocationHours;
        dst.Version = src.Version;
        if (!string.IsNullOrEmpty(src.Mime))
            dst.Mime = src.Mime;
        dst.ChildDesc = src.ChildDesc;
        if (src.ChildType != ValueType.UNKNOWN)
            dst.ChildType = src.ChildType;
        dst.ChildNullable |= src.ChildNullable;
        if (src.ChildEnum.Count > 0)
        {
            dst.ChildEnum = src.ChildEnum;
            dst.ChildType = ValueType.ENUM;
        }
    }

    private static string GetFieldKey(System.Reflection.PropertyInfo property, MmTag fieldType, MM? attribute)
    {
        if (attribute != null && !string.IsNullOrEmpty(attribute.Name) && attribute.Name != "-")
        {
            return attribute.Name;
        }
        if (!string.IsNullOrEmpty(fieldType.Name) && fieldType.Name != "-")
        {
            return fieldType.Name;
        }
        return CamelToSnake.Convert(property.Name);
    }

    private static void EncodeScalarPayload(WireEncoder encoder, object? value, MmTag tag)
    {
        if (tag.IsNull)
        {
            switch (tag.Type)
            {
                case ValueType.BOOL:
                    encoder.EncodeSimple(SimpleValue.NULL_BOOL);
                    return;
                case ValueType.INT:
                case ValueType.INT8:
                case ValueType.INT16:
                case ValueType.INT32:
                case ValueType.INT64:
                case ValueType.UINT:
                case ValueType.UINT16:
                case ValueType.UINT32:
                case ValueType.UINT64:
                    encoder.EncodeSimple(SimpleValue.NULL_INT);
                    return;
                case ValueType.FLOAT32:
                case ValueType.FLOAT64:
                case ValueType.DECIMAL:
                    encoder.EncodeSimple(SimpleValue.NULL_FLOAT);
                    return;
                case ValueType.STRING:
                case ValueType.EMAIL:
                case ValueType.URL:
                    encoder.EncodeSimple(SimpleValue.NULL_STRING);
                    return;
                case ValueType.BYTES:
                    encoder.EncodeSimple(SimpleValue.NULL_BYTES);
                    return;
                default:
                    encoder.EncodeSimple(SimpleValue.NULL_STRING);
                    return;
            }
        }

        switch (tag.Type)
        {
            case ValueType.BOOL:
                encoder.EncodeBool((bool)value!);
                break;
            case ValueType.INT:
            case ValueType.INT8:
            case ValueType.INT16:
            case ValueType.INT32:
            case ValueType.INT64:
            case ValueType.UINT:
            case ValueType.UINT16:
            case ValueType.UINT32:
            case ValueType.UINT64:
                encoder.EncodeInt64(Convert.ToInt64(value));
                break;
            case ValueType.FLOAT32:
                encoder.EncodeFloatString(value!.ToString()!);
                break;
            case ValueType.FLOAT64:
            case ValueType.DECIMAL:
                encoder.EncodeFloatString(value!.ToString()!);
                break;
            case ValueType.STRING:
            case ValueType.EMAIL:
            case ValueType.URL:
                encoder.EncodeString(value as string ?? "");
                break;
            case ValueType.BYTES:
                encoder.EncodeBytes(value as byte[] ?? Array.Empty<byte>());
                break;
            case ValueType.BIGINT:
                encoder.EncodeBigIntDecimal(value as string ?? "0");
                break;
            case ValueType.UUID:
                encoder.EncodeBytes(UuidToBytes(value as string ?? ""));
                break;
            case ValueType.DATETIME:
                var dateTime = (DateTime)value!;
                encoder.EncodeInt64(TimeUtil.EpochSeconds(dateTime));
                break;
            case ValueType.DATE:
                var date = (DateTime)value!;
                encoder.EncodeInt64(TimeUtil.DaysSinceEpochUtc(date));
                break;
            case ValueType.TIME:
                var time = (DateTime)value!;
                encoder.EncodeInt64(TimeUtil.SecondsOfDay(time));
                break;
            case ValueType.ENUM:
                encoder.EncodeInt64(Convert.ToInt64(value));
                break;
            default:
                throw new Exception($"Unsupported scalar type: {tag.Type}");
        }
    }

    private static void EncodeList(WireEncoder encoder, object? list, MmTag tag)
    {
        if (list == null)
        {
            var nt = tag.Copy();
            nt.IsNull = true;
            encoder.EncodeTaggedPayload(Array.Empty<byte>(), nt.ToBytes());
            return;
        }

        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();

        if (list is System.Collections.IEnumerable enumerable)
        {
            foreach (var item in enumerable)
            {
                elementEncoder.Reset();
                var itemTag = MmTag.Empty();
                itemTag.InheritFromArrayParent(tag);
                if (itemTag.Type == ValueType.UNKNOWN && item != null)
                {
                    itemTag.Type = InferTypeFromValue(item);
                    if (itemTag.Type == ValueType.UNKNOWN)
                    {
                        itemTag.Type = ValueType.STRUCT;
                    }
                }
                EncodeValue(elementEncoder, item, itemTag);
                body.WriteAll(elementEncoder.ToByteArray());
            }
        }

        elementEncoder.Reset();
        elementEncoder.EncodeArrayPayload(body.ToArray());
        encoder.EncodeTaggedPayload(elementEncoder.ToByteArray(), tag.ToBytes());
    }

    private static byte[] UuidToBytes(string uuid)
    {
        uuid = uuid.Replace("-", "");
        byte[] bytes = new byte[16];
        for (int i = 0; i < 32; i += 2)
        {
            bytes[i / 2] = Convert.ToByte(uuid.Substring(i, 2), 16);
        }
        return bytes;
    }
}