using ValueType = MetaMessage.Ir.ValueType;
namespace MetaMessage.Core;

public static class TypeInference
{
    public static ValueType ValueTypeForType(Type type)
    {
        if (type == typeof(bool) || type == typeof(bool?))
        {
            return ValueType.Bool;
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
            return ValueType.Decimal;
        }
        else if (type == typeof(string))
        {
            return ValueType.Str;
        }
        else if (type == typeof(byte[]))
        {
            return ValueType.Bytes;
        }
        else if (type == typeof(DateTime) || type == typeof(DateTime?))
        {
            return ValueType.Datetime;
        }
        else if (type.IsEnum)
        {
            return ValueType.Enums;
        }
        else if (type.IsArray)
        {
            return ValueType.Arr;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            return ValueType.Vec;
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Dictionary<,>))
        {
            return ValueType.Map;
        }
        else
        {
            return ValueType.Obj;
        }
    }

    public static ValueType ValueTypeForComponent(Type componentType)
    {
        return ValueTypeForType(componentType);
    }

    public static bool IsIntegerType(ValueType vt)
    {
        return vt == ValueType.I ||
               vt == ValueType.I8 ||
               vt == ValueType.I16 ||
               vt == ValueType.I32 ||
               vt == ValueType.I64 ||
               vt == ValueType.U ||
               vt == ValueType.U8 ||
               vt == ValueType.U16 ||
               vt == ValueType.U32 ||
               vt == ValueType.U64 ||
               vt == ValueType.Enums;
    }

    public static bool IsFloatType(ValueType vt)
    {
        return vt == ValueType.F32 ||
               vt == ValueType.F64 ||
               vt == ValueType.Decimal;
    }
}