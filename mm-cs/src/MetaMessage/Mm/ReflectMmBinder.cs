namespace MetaMessage.Mm;

public static class ReflectMmBinder
{
    public static void Bind(byte[] data, object target)
    {
        var decoder = new WireDecoder(data);
        var tree = decoder.Decode();
        BindTree(tree, target);
    }

    private static void BindTree(IMmTree tree, object target)
    {
        if (tree is MmMap map)
        {
            BindMap(map, target);
        }
        else if (tree is MmScalar scalar)
        {
            // 处理标量类型
        }
    }

    private static void BindMap(MmMap map, object target)
    {
        var type = target.GetType();
        foreach (var entry in map.Entries)
        {
            var key = entry.Key.Data as string;
            var value = entry.Value;

            if (key == null)
                continue;

            var property = FindProperty(type, key);
            if (property == null)
                continue;

            if (!property.CanWrite)
                continue;

            var propertyValue = ConvertValue(value, property.PropertyType);
            property.SetValue(target, propertyValue);
        }
    }

    private static System.Reflection.PropertyInfo? FindProperty(Type type, string key)
    {
        // 首先尝试直接匹配属性名
        var property = type.GetProperty(key);
        if (property != null)
            return property;

        // 尝试驼峰命名匹配
        var camelKey = SnakeToCamel(key);
        property = type.GetProperty(camelKey);
        if (property != null)
            return property;

        return null;
    }

    private static string SnakeToCamel(string snakeCase)
    {
        var parts = snakeCase.Split('_');
        var result = new System.Text.StringBuilder();
        for (int i = 0; i < parts.Length; i++)
        {
            if (i == 0)
                result.Append(parts[i].ToLower());
            else
                result.Append(char.ToUpper(parts[i][0])).Append(parts[i].Substring(1).ToLower());
        }
        return result.ToString();
    }

    private static object? ConvertValue(IMmTree tree, Type targetType)
    {
        if (tree is MmScalar scalar)
        {
            return ConvertScalar(scalar, targetType);
        }
        else if (tree is MmArray array)
        {
            return ConvertArray(array, targetType);
        }
        else if (tree is MmMap map)
        {
            return ConvertMap(map, targetType);
        }
        return null;
    }

    private static object? ConvertScalar(MmScalar scalar, Type targetType)
    {
        if (scalar.Tag.IsNull)
        {
            return null;
        }

        var data = scalar.Data;
        if (data == null)
            return null;

        // 处理可空类型
        if (targetType.IsNullableType())
        {
            targetType = Nullable.GetUnderlyingType(targetType)!;
        }

        // 处理基本类型
        if (targetType == typeof(bool))
        {
            return Convert.ToBoolean(data);
        }
        else if (targetType == typeof(sbyte))
        {
            return Convert.ToSByte(data);
        }
        else if (targetType == typeof(short))
        {
            return Convert.ToInt16(data);
        }
        else if (targetType == typeof(int))
        {
            return Convert.ToInt32(data);
        }
        else if (targetType == typeof(long))
        {
            return Convert.ToInt64(data);
        }
        else if (targetType == typeof(byte))
        {
            return Convert.ToByte(data);
        }
        else if (targetType == typeof(ushort))
        {
            return Convert.ToUInt16(data);
        }
        else if (targetType == typeof(uint))
        {
            return Convert.ToUInt32(data);
        }
        else if (targetType == typeof(ulong))
        {
            return Convert.ToUInt64(data);
        }
        else if (targetType == typeof(float))
        {
            return Convert.ToSingle(data);
        }
        else if (targetType == typeof(double))
        {
            return Convert.ToDouble(data);
        }
        else if (targetType == typeof(decimal))
        {
            return Convert.ToDecimal(data);
        }
        else if (targetType == typeof(string))
        {
            return data.ToString();
        }
        else if (targetType == typeof(byte[]))
        {
            return data as byte[];
        }
        else if (targetType == typeof(DateTime))
        {
            if (data is long timestamp)
            {
                return TimeUtil.FromEpochSeconds(timestamp);
            }
            return DateTime.Parse(data.ToString()!);
        }
        else if (targetType.IsEnum)
        {
            return Enum.ToObject(targetType, data);
        }

        return data;
    }

    private static object? ConvertArray(MmArray array, Type targetType)
    {
        if (targetType.IsArray)
        {
            var elementType = targetType.GetElementType()!;
            var elements = new List<object?>();
            foreach (var item in array.Children)
            {
                elements.Add(ConvertValue(item, elementType));
            }
            var arrayInstance = Array.CreateInstance(elementType, elements.Count);
            for (int i = 0; i < elements.Count; i++)
            {
                arrayInstance.SetValue(elements[i], i);
            }
            return arrayInstance;
        }
        else if (targetType.IsGenericType && targetType.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = targetType.GetGenericArguments()[0];
            var listType = typeof(List<>).MakeGenericType(elementType);
            var list = Activator.CreateInstance(listType);
            var addMethod = listType.GetMethod("Add");

            foreach (var item in array.Children)
            {
                var convertedValue = ConvertValue(item, elementType);
                addMethod?.Invoke(list, new object?[] { convertedValue });
            }

            return list;
        }

        return null;
    }

    private static object? ConvertMap(MmMap map, Type targetType)
    {
        var instance = Activator.CreateInstance(targetType);
        BindMap(map, instance!);
        return instance;
    }

    private static bool IsNullableType(this Type type)
    {
        return type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Nullable<>);
    }
}