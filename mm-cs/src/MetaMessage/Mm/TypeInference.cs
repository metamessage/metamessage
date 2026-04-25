namespace MetaMessage.Mm;

public static class TypeInference
{
    public static ValueType ValueTypeForType(Type type)
    {
        if (type == typeof(bool) || type == typeof(bool?))
        {
            return ValueType.BOOL;
        }
        else if (type == typeof(sbyte) || type == typeof(sbyte?))
        {
            return ValueType.INT8;
        }
        else if (type == typeof(short) || type == typeof(short?))
        {
            return ValueType.INT16;
        }
        else if (type == typeof(int) || type == typeof(int?))
        {
            return ValueType.INT32;
        }
        else if (type == typeof(long) || type == typeof(long?))
        {
            return ValueType.INT64;
        }
        else if (type == typeof(byte) || type == typeof(byte?))
        {
            return ValueType.UINT;
        }
        else if (type == typeof(ushort) || type == typeof(ushort?))
        {
            return ValueType.UINT16;
        }
        else if (type == typeof(uint) || type == typeof(uint?))
        {
            return ValueType.UINT32;
        }
        else if (type == typeof(ulong) || type == typeof(ulong?))
        {
            return ValueType.UINT64;
        }
        else if (type == typeof(float) || type == typeof(float?))
        {
            return ValueType.FLOAT32;
        }
        else if (type == typeof(double) || type == typeof(double?))
        {
            return ValueType.FLOAT64;
        }
        else if (type == typeof(decimal) || type == typeof(decimal?))
        {
            return ValueType.DECIMAL;
        }
        else if (type == typeof(string))
        {
            return ValueType.STRING;
        }
        else if (type == typeof(byte[]))
        {
            return ValueType.BYTES;
        }
        else if (type == typeof(DateTime) || type == typeof(DateTime?))
        {
            return ValueType.DATETIME;
        }
        else if (type.IsEnum)
        {
            return ValueType.ENUM;
        }
        else if (type.IsArray)
        {
            return ValueType.ARRAY;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            return ValueType.SLICE;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Dictionary<,>))
        {
            return ValueType.MAP;
        }
        else
        {
            return ValueType.STRUCT;
        }
    }

    public static ValueType ValueTypeForComponent(Type componentType)
    {
        return ValueTypeForType(componentType);
    }
}