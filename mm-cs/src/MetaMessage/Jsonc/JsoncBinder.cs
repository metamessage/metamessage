using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

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
            return new JsoncValue { Value = intVal.ToString(), TokenType = JsoncTokenType.Number };
        }

        if (value is long longVal)
        {
            return new JsoncValue { Value = longVal.ToString(), TokenType = JsoncTokenType.Number };
        }

        if (value is double doubleVal)
        {
            return new JsoncValue { Value = doubleVal.ToString("G"), TokenType = JsoncTokenType.Number };
        }

        if (value is float floatVal)
        {
            return new JsoncValue { Value = floatVal.ToString("G"), TokenType = JsoncTokenType.Number };
        }

        if (value is System.DateTime dt)
        {
            return new JsoncValue { Value = dt.ToString("yyyy-MM-dd HH:mm:ss"), TokenType = JsoncTokenType.String };
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
                continue;

            var propValue = prop.GetValue(value);
            var propNode = StructToNode(propValue!);
            obj.Add(prop.Name, propNode);
        }

        return obj;
    }

    /// <summary>
    /// Convert JsoncNode tree to IMmTree for use with the MetaMessage IR layer.
    /// </summary>
    public IMmTree ToMmTree(IJsoncNode node)
    {
        if (node is JsoncObject obj)
        {
            var fields = new List<KeyValuePair<MmScalar, IMmTree>>();
            foreach (var kvp in obj.Fields)
            {
                var keyNode = new MmScalar(kvp.Key, kvp.Key, Tag.NewTag());
                var valNode = ToMmTree(kvp.Value);
                fields.Add(new KeyValuePair<MmScalar, IMmTree>(keyNode, valNode));
            }
            return new MmMap(fields, obj.Tag ?? Tag.NewTag()) { Path = obj.Path };
        }

        if (node is JsoncArray arr)
        {
            var children = new List<IMmTree>();
            foreach (var item in arr.Elements)
            {
                children.Add(ToMmTree(item));
            }
            return new MmArray(children, arr.Tag ?? Tag.NewTag()) { Path = arr.Path };
        }

        if (node is JsoncDoc doc)
        {
            var fields = new List<KeyValuePair<MmScalar, IMmTree>>();
            foreach (var kvp in doc.Fields)
            {
                var keyNode = new MmScalar(kvp.Key, kvp.Key, Tag.NewTag());
                var valNode = ToMmTree(kvp.Value);
                fields.Add(new KeyValuePair<MmScalar, IMmTree>(keyNode, valNode));
            }
            return new MmDoc(fields, doc.Tag ?? Tag.NewTag()) { Path = doc.Path };
        }

        if (node is JsoncValue val)
        {
            var text = val.Value?.ToString() ?? "";
            var tag = val.Tag ?? Tag.NewTag();
            return new MmScalar(val.Value, text, tag) { Path = val.Path };
        }

        throw new Exception($"Unknown node type: {node.GetType()}");
    }

    /// <summary>
    /// Convert IMmTree to JsoncNode tree for serialization.
    /// </summary>
    public IJsoncNode FromMmTree(IMmTree tree)
    {
        if (tree is MmMap mmMap)
        {
            var obj = new JsoncObject { Tag = mmMap.Tag, Path = mmMap.Path };
            foreach (var entry in mmMap.Entries)
            {
                var key = entry.Key.Text;
                var valNode = FromMmTree(entry.Value);
                obj.Add(key, valNode);
            }
            return obj;
        }

        if (tree is MmDoc mmDoc)
        {
            var doc = new JsoncDoc { Tag = mmDoc.Tag, Path = mmDoc.Path };
            foreach (var entry in mmDoc.Fields)
            {
                var key = entry.Key.Text;
                var valNode = FromMmTree(entry.Value);
                doc.Add(key, valNode);
            }
            return doc;
        }

        if (tree is MmArray mmArray)
        {
            var arr = new JsoncArray { Tag = mmArray.Tag, Path = mmArray.Path };
            foreach (var child in mmArray.Children)
            {
                arr.Add(FromMmTree(child));
            }
            return arr;
        }

        if (tree is MmScalar mmScalar)
        {
            JsoncTokenType tokenType;
            object? value;
            if (mmScalar.Tag != null)
            {
                switch (mmScalar.Tag.Type)
                {
                    case ValueType.Str:
                    case ValueType.Bytes:
                    case ValueType.Media:
                    case ValueType.Image:
                    case ValueType.Video:
                    case ValueType.Uuid:
                    case ValueType.Bigint:
                    case ValueType.Datetime:
                    case ValueType.Date:
                    case ValueType.Time:
                    case ValueType.Enums:
                    case ValueType.Ip:
                    case ValueType.Url:
                    case ValueType.Email:
                        tokenType = JsoncTokenType.String;
                        value = mmScalar.Text;
                        break;
                    case ValueType.Bool:
                        tokenType = mmScalar.Text == "true" ? JsoncTokenType.True : JsoncTokenType.False;
                        value = mmScalar.Data;
                        break;
                    default:
                        tokenType = JsoncTokenType.Number;
                        value = mmScalar.Data;
                        break;
                }
            }
            else
            {
                tokenType = JsoncTokenType.String;
                value = mmScalar.Text;
            }

            return new JsoncValue
            {
                Value = value,
                TokenType = tokenType,
                Tag = mmScalar.Tag,
                Path = mmScalar.Path
            };
        }

        throw new Exception($"Unknown IMmTree type: {tree.GetType()}");
    }

    private void BindObject(JsoncObject obj, object target)
    {
        var type = target.GetType();
        var properties = type.GetProperties();

        foreach (var prop in properties)
        {
            if (!prop.CanWrite)
                continue;

            var key = GetJsonKey(prop.Name);
            if (!obj.Fields.TryGetValue(key, out var node))
                continue;

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
                    var itemValue = ConvertJsoncValue(array.Elements[i], elementType);
                    arr.SetValue(itemValue, i);
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
                    var itemValue = ConvertJsoncValue(array.Elements[i], elementType);
                    list.Add(itemValue);
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
    }

    private object? ConvertJsoncValue(IJsoncNode node, Type targetType)
    {
        if (node is JsoncValue value)
        {
            if (targetType == typeof(string))
                return value.Value?.ToString() ?? "";
            if (targetType == typeof(bool))
                return value.Value is bool b ? b : (bool.TryParse(value.Value?.ToString(), out bool br) ? br : false);
            if (targetType == typeof(int))
                return value.Value is int i ? i : (int.TryParse(value.Value?.ToString(), out int ir) ? ir : 0);
            if (targetType == typeof(long))
                return value.Value is long l ? l : (long.TryParse(value.Value?.ToString(), out long lr) ? lr : 0L);
            if (targetType == typeof(double))
                return value.Value is double d ? d : (double.TryParse(value.Value?.ToString(), out double dr) ? dr : 0.0);
            if (targetType == typeof(float))
                return value.Value is float f ? f : (float.TryParse(value.Value?.ToString(), out float fr) ? fr : 0f);
            if (targetType == typeof(byte[]))
                return value.Value is string s ? Convert.FromBase64String(s) : Array.Empty<byte>();
            return value.Value;
        }
        if (node is JsoncObject childObj)
        {
            var childTarget = Activator.CreateInstance(targetType)!;
            BindObject(childObj, childTarget);
            return childTarget;
        }
        if (node is JsoncArray childArray && targetType.IsGenericType && targetType.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = targetType.GetGenericArguments()[0];
            var listType = typeof(System.Collections.Generic.List<>).MakeGenericType(elementType);
            var list = (System.Collections.IList)Activator.CreateInstance(listType)!;
            foreach (var element in childArray.Elements)
            {
                list.Add(ConvertJsoncValue(element, elementType));
            }
            return list;
        }
        return null;
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
                var itemValue = ConvertJsoncValue(array.Elements[i], elementType);
                arr.SetValue(itemValue, i);
            }
            type.GetProperty("Length")!.SetValue(target, arr);
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = type.GetGenericArguments()[0];
            var list = (System.Collections.IList)target;
            for (int i = 0; i < array.Elements.Count; i++)
            {
                var itemValue = ConvertJsoncValue(array.Elements[i], elementType);
                list.Add(itemValue);
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
                prop.SetValue(target, b);
            else if (bool.TryParse(rawValue.ToString(), out bool result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(int))
        {
            if (rawValue is int i)
                prop.SetValue(target, i);
            else if (rawValue is double d)
                prop.SetValue(target, (int)d);
            else if (int.TryParse(rawValue.ToString(), out int result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(long))
        {
            if (rawValue is long l)
                prop.SetValue(target, l);
            else if (rawValue is double d)
                prop.SetValue(target, (long)d);
            else if (long.TryParse(rawValue.ToString(), out long result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(double))
        {
            if (rawValue is double d)
                prop.SetValue(target, d);
            else if (double.TryParse(rawValue.ToString(), out double result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(float))
        {
            if (rawValue is float f)
                prop.SetValue(target, f);
            else if (float.TryParse(rawValue.ToString(), out float result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(byte[]))
        {
            if (rawValue is string str)
                prop.SetValue(target, Convert.FromBase64String(str));
            return;
        }

        if (propType == typeof(System.DateTime))
        {
            if (rawValue is string str && System.DateTime.TryParse(str, out var dt))
                prop.SetValue(target, dt);
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