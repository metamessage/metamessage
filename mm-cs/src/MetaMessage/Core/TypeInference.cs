namespace MetaMessage.Core;

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
            return ValueType.I8;
        }
        else if (type == typeof(short) || type == typeof(short?))
        {
            return ValueType.I16;
        }
        else if (type == typeof(int) || type == typeof(int?))
        {
            return ValueType.I32;
        }
        else if (type == typeof(long) || type == typeof(long?))
        {
            return ValueType.I64;
        }
        else if (type == typeof(byte) || type == typeof(byte?))
        {
            return ValueType.U;
        }
        else if (type == typeof(ushort) || type == typeof(ushort?))
        {
            return ValueType.U16;
        }
        else if (type == typeof(uint) || type == typeof(uint?))
        {
            return ValueType.U32;
        }
        else if (type == typeof(ulong) || type == typeof(ulong?))
        {
            return ValueType.U64;
        }
        else if (type == typeof(float) || type == typeof(float?))
        {
            return ValueType.F32;
        }
        else if (type == typeof(double) || type == typeof(double?))
        {
            return ValueType.F64;
        }
        else if (type == typeof(decimal) || type == typeof(decimal?))
        {
            return ValueType.DECIMAL;
        }
        else if (type == typeof(string))
        {
            return ValueType.STR;
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
            return ValueType.ARR;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            return ValueType.VEC;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Dictionary<,>))
        {
            return ValueType.MAP;
        }
        else
        {
            return ValueType.OBJ;
        }
    }

    public static ValueType ValueTypeForComponent(Type componentType)
    {
        return ValueTypeForType(componentType);
    }
}