using MetaMessage.Ir;

namespace MetaMessage.Core;

public static class ReflectMmBinder
{
    public static void Bind(byte[] data, object target)
    {
        var decoder = new WireDecoder(data);
        var tree = decoder.Decode();
        BindTree(tree, target);
    }

    private static void BindTree(INode tree, object target)
    {
        if (tree is NodeObject map)
        {
            BindMap(map, target);
        }
        else if (tree is NodeScalar scalar)
        {
            // 处理标量类型
        }
        else if (tree is NodeNull)
        {
            // NodeNull - nothing to bind
        }
    }

    // example 哨兵项判定：空容器编码时注入的占位节点，绑定时必须跳过
    private static bool IsExampleNode(INode? node)
    {
        return node is not null && node.Tag.Example;
    }

    // is_null 标记判定：编码带 is_null 时目标绑定为 null
    private static bool IsNullNode(INode? node)
    {
        return node is not null && node.Tag.IsNull;
    }

    public static object? BindDynamic(INode tree)
    {
        if (tree is NodeNull)
            return null;

        if (IsNullNode(tree))
            return null;

        if (tree is NodeScalar scalar)
            return scalar.Data;

        if (tree is NodeArray array)
        {
            var list = new List<object?>();
            foreach (var child in array.Children)
            {
                if (IsExampleNode(child))
                    continue;
                list.Add(BindDynamic(child));
            }
            return list;
        }

        if (tree is NodeObject map)
        {
            var dict = new Dictionary<string, object?>();
            foreach (var entry in map.Entries)
            {
                if (IsExampleNode(entry.Value))
                    continue;
                dict[entry.Key.Text] = BindDynamic(entry.Value);
            }
            return dict;
        }

        return null;
    }

    private static void BindMap(NodeObject map, object target)
    {
        var type = target.GetType();
        foreach (var entry in map.Entries)
        {
            // 跳过 example 哨兵项（空容器占位），不绑定到目标对象
            if (IsExampleNode(entry.Value))
                continue;

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
            result.Append(char.ToUpper(parts[i][0])).Append(parts[i].Substring(1).ToLower());
        }
        return result.ToString();
    }

    private static object? ConvertValue(INode tree, Type targetType)
    {
        // Handle object type dynamically (like Go's interface{})
        if (targetType == typeof(object))
        {
            return BindDynamic(tree);
        }

        if (tree is NodeNull)
        {
            return null;
        }

        // is_null 编码：绑定为 null
        if (IsNullNode(tree))
        {
            return null;
        }

        if (tree is NodeScalar scalar)
        {
            return ConvertScalar(scalar, targetType);
        }
        else if (tree is NodeArray array)
        {
            return ConvertArray(array, targetType);
        }
        else if (tree is NodeObject map)
        {
            return ConvertMap(map, targetType);
        }
        return null;
    }

    private static object? ConvertScalar(NodeScalar scalar, Type targetType)
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
            var underlyingType = Enum.GetUnderlyingType(targetType);
            var convertedData = Convert.ChangeType(data, underlyingType);
            return Enum.ToObject(targetType, convertedData);
        }

        return data;
    }

    private static object? ConvertArray(NodeArray array, Type targetType)
    {
        if (targetType.IsArray)
        {
            var elementType = targetType.GetElementType()!;
            var elements = new List<object?>();
            foreach (var item in array.Children)
            {
                // 跳过 example 哨兵项：全为 example 时产出长度为 0 的非 null 数组
                if (IsExampleNode(item))
                    continue;
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
                // 跳过 example 哨兵项：全为 example 时产出空的非 null List
                if (IsExampleNode(item))
                    continue;
                var convertedValue = ConvertValue(item, elementType);
                addMethod?.Invoke(list, new object?[] { convertedValue });
            }

            return list;
        }

        return null;
    }

    private static object? ConvertMap(NodeObject map, Type targetType)
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