namespace MetaMessage.Jsonc;

public class JsoncBinder
{
    public T Bind<T>(IJsoncNode node) where T : new()
    {
        var result = new T();
        Bind(node, result);
        return result;
    }

    public void Bind(IJsoncNode node, object target)
    {
        if (node is JsoncObject obj && target != null)
        {
            BindObject(obj, target);
        }
        else if (node is JsoncArray array && target != null)
        {
            BindArray(array, target);
        }
        else if (node is JsoncValue value && target != null)
        {
            BindValue(value, target);
        }
    }

    public IJsoncNode StructToNode(object value)
    {
        if (value == null)
        {
            return new JsoncValue { Value = null, TokenType = JsoncTokenType.Null };
        }

        var type = value.GetType();

        if (type.IsArray)
        {
            var array = new JsoncArray();
            var elementType = type.GetElementType()!;
            var arr = (Array)value;
            for (int i = 0; i < arr.Length; i++)
            {
                array.Add(StructToNode(arr.GetValue(i)!));
            }
            return array;
        }

        if (value is System.Collections.IList list)
        {
            var array = new JsoncArray();
            foreach (var item in list)
            {
                array.Add(StructToNode(item));
            }
            return array;
        }

        if (value is string str)
        {
            return new JsoncValue { Value = str, TokenType = JsoncTokenType.String };
        }

        if (value is bool b)
        {
            return new JsoncValue { Value = b, TokenType = b ? JsoncTokenType.True : JsoncTokenType.False };
        }

        if (value is int intVal)
        {
            return new JsoncValue { Value = intVal, TokenType = JsoncTokenType.Number };
        }

        if (value is long longVal)
        {
            return new JsoncValue { Value = longVal, TokenType = JsoncTokenType.Number };
        }

        if (value is double doubleVal)
        {
            return new JsoncValue { Value = doubleVal, TokenType = JsoncTokenType.Number };
        }

        if (value is float floatVal)
        {
            return new JsoncValue { Value = floatVal, TokenType = JsoncTokenType.Number };
        }

        if (value is System.DateTime dt)
        {
            return new JsoncValue { Value = dt.ToString("o"), TokenType = JsoncTokenType.String };
        }

        if (value is byte[] bytes)
        {
            return new JsoncValue { Value = Convert.ToBase64String(bytes), TokenType = JsoncTokenType.String };
        }

        var obj = new JsoncObject();
        var properties = type.GetProperties();
        foreach (var prop in properties)
        {
            if (!prop.CanRead)
            {
                continue;
            }

            var propValue = prop.GetValue(value);
            var propNode = StructToNode(propValue!);
            obj.Add(prop.Name, propNode);
        }

        return obj;
    }

    private void BindObject(JsoncObject obj, object target)
    {
        var type = target.GetType();
        var properties = type.GetProperties();

        foreach (var prop in properties)
        {
            if (!prop.CanWrite)
            {
                continue;
            }

            var key = GetJsonKey(prop.Name);
            if (!obj.Fields.TryGetValue(key, out var node))
            {
                continue;
            }

            BindNodeToProperty(node, prop, target);
        }
    }

    private void BindNodeToProperty(IJsoncNode node, System.Reflection.PropertyInfo prop, object target)
    {
        var propType = prop.PropertyType;

        if (propType.IsArray)
        {
            if (node is JsoncArray array)
            {
                var elementType = propType.GetElementType()!;
                var arr = Array.CreateInstance(elementType, array.Elements.Count);
                for (int i = 0; i < array.Elements.Count; i++)
                {
                    var element = Array.CreateInstance(elementType, 1);
                    BindNodeToProperty(array.Elements[i], elementType.GetProperty("Item")!, element);
                    arr.SetValue(element.GetValue(0), i);
                }
                prop.SetValue(target, arr);
            }
            return;
        }

        if (propType.IsGenericType && propType.GetGenericTypeDefinition() == typeof(List<>))
        {
            if (node is JsoncArray array)
            {
                var elementType = propType.GetGenericArguments()[0];
                var listType = typeof(System.Collections.Generic.List<>).MakeGenericType(elementType);
                var list = (System.Collections.IList)Activator.CreateInstance(listType)!;
                for (int i = 0; i < array.Elements.Count; i++)
                {
                    var element = Activator.CreateInstance(elementType)!;
                    BindNodeToProperty(array.Elements[i], elementType.GetProperty("Item")!, element);
                    list.Add(element);
                }
                prop.SetValue(target, list);
            }
            return;
        }

        if (node is JsoncObject childObj)
        {
            var childTarget = Activator.CreateInstance(propType)!;
            BindObject(childObj, childTarget);
            prop.SetValue(target, childTarget);
            return;
        }

        if (node is JsoncValue value)
        {
            BindValue(value, prop, target);
            return;
        }

        if (node is JsoncArray)
        {
            return;
        }
    }

    private void BindArray(JsoncArray array, object target)
    {
        var type = target.GetType();
        if (type.IsArray)
        {
            var elementType = type.GetElementType()!;
            var arr = Array.CreateInstance(elementType, array.Elements.Count);
            for (int i = 0; i < array.Elements.Count; i++)
            {
                var element = Activator.CreateInstance(elementType)!;
                BindNodeToProperty(array.Elements[i], elementType.GetProperty("Item")!, element);
                arr.SetValue(element, i);
            }
            type.GetProperty("Length")!.SetValue(target, arr);
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = type.GetGenericArguments()[0];
            var list = (System.Collections.IList)target;
            for (int i = 0; i < array.Elements.Count; i++)
            {
                var element = Activator.CreateInstance(elementType)!;
                BindNodeToProperty(array.Elements[i], elementType.GetProperty("Item")!, element);
                list.Add(element);
            }
        }
    }

    private void BindValue(JsoncValue value, object target)
    {
        var type = target.GetType();
        BindValue(value, type.GetProperty("Item")!, target);
    }

    private void BindValue(JsoncValue value, System.Reflection.PropertyInfo prop, object target)
    {
        var propType = prop.PropertyType;
        var rawValue = value.Value;

        if (rawValue == null || (rawValue is string s && s == "null"))
        {
            if (IsNullable(propType))
            {
                prop.SetValue(target, null);
            }
            return;
        }

        if (propType == typeof(string))
        {
            prop.SetValue(target, rawValue.ToString());
            return;
        }

        if (propType == typeof(bool))
        {
            if (rawValue is bool b)
            {
                prop.SetValue(target, b);
            }
            else if (bool.TryParse(rawValue.ToString(), out bool result))
            {
                prop.SetValue(target, result);
            }
            return;
        }

        if (propType == typeof(int))
        {
            if (rawValue is int i)
            {
                prop.SetValue(target, i);
            }
            else if (rawValue is double d)
            {
                prop.SetValue(target, (int)d);
            }
            else if (int.TryParse(rawValue.ToString(), out int result))
            {
                prop.SetValue(target, result);
            }
            return;
        }

        if (propType == typeof(long))
        {
            if (rawValue is long l)
            {
                prop.SetValue(target, l);
            }
            else if (rawValue is double d)
            {
                prop.SetValue(target, (long)d);
            }
            else if (long.TryParse(rawValue.ToString(), out long result))
            {
                prop.SetValue(target, result);
            }
            return;
        }

        if (propType == typeof(double))
        {
            if (rawValue is double d)
            {
                prop.SetValue(target, d);
            }
            else if (double.TryParse(rawValue.ToString(), out double result))
            {
                prop.SetValue(target, result);
            }
            return;
        }

        if (propType == typeof(float))
        {
            if (rawValue is float f)
            {
                prop.SetValue(target, f);
            }
            else if (float.TryParse(rawValue.ToString(), out float result))
            {
                prop.SetValue(target, result);
            }
            return;
        }

        if (propType == typeof(byte[]))
        {
            if (rawValue is string str)
            {
                prop.SetValue(target, Convert.FromBase64String(str));
            }
            return;
        }

        if (propType == typeof(System.DateTime))
        {
            if (rawValue is string str && System.DateTime.TryParse(str, out var dt))
            {
                prop.SetValue(target, dt);
            }
            return;
        }
    }

    private string GetJsonKey(string propertyName)
    {
        return propertyName;
    }

    private bool IsNullable(Type type)
    {
        return !type.IsValueType || Nullable.GetUnderlyingType(type) != null;
    }
}