using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;
namespace MetaMessage.Core;

public class Binder
{
    public T Bind<T>(INode node) where T : new()
    {
        var result = new T();
        Bind(node, result);
        return result;
    }

    public void Bind(INode node, object target)
    {
        if (node is NodeObject map && target != null)
        {
            BindObject(map, target);
        }
        else if (node is NodeArray array && target != null)
        {
            BindArray(array, target);
        }
        else if (node is NodeScalar scalar && target != null)
        {
            BindValue(scalar, target);
        }
    }

    public INode StructToNode(object value)
    {
        if (value == null)
        {
            return new NodeScalar(null, "null", Tag.NewTag());
        }

        var type = value.GetType();

        if (type.IsArray)
        {
            var arr = (Array)value;
            var children = new List<INode>();
            for (int i = 0; i < arr.Length; i++)
            {
                children.Add(StructToNode(arr.GetValue(i)!));
            }
            return new NodeArray(children, Tag.NewTag());
        }

        if (value is System.Collections.IList list)
        {
            var children = new List<INode>();
            foreach (var item in list)
            {
                children.Add(StructToNode(item));
            }
            return new NodeArray(children, Tag.NewTag());
        }

        if (value is string str)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.Str;
            return new NodeScalar(str, str, tag);
        }

        if (value is bool b)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.Bool;
            return new NodeScalar(b, b ? "true" : "false", tag);
        }

        if (value is int intVal)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.I;
            return new NodeScalar(intVal, intVal.ToString(), tag);
        }

        if (value is long longVal)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.I;
            return new NodeScalar(longVal, longVal.ToString(), tag);
        }

        if (value is double doubleVal)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.F64;
            return new NodeScalar(doubleVal, doubleVal.ToString("G"), tag);
        }

        if (value is float floatVal)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.F32;
            return new NodeScalar(floatVal, floatVal.ToString("G"), tag);
        }

        if (value is System.DateTime dt)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.Str;
            return new NodeScalar(dt.ToString("yyyy-MM-dd HH:mm:ss"), dt.ToString("yyyy-MM-dd HH:mm:ss"), tag);
        }

        if (value is byte[] bytes)
        {
            var tag = Tag.NewTag();
            tag.Type = ValueType.Bytes;
            return new NodeScalar(bytes, Convert.ToBase64String(bytes), tag);
        }

        var entries = new List<KeyValuePair<NodeScalar, INode>>();
        var properties = type.GetProperties();
        foreach (var prop in properties)
        {
            if (!prop.CanRead)
                continue;

            var propValue = prop.GetValue(value);
            var propNode = StructToNode(propValue!);
            var keyScalar = new NodeScalar(prop.Name, prop.Name, Tag.Empty());
            entries.Add(new KeyValuePair<NodeScalar, INode>(keyScalar, propNode));
        }

        return new NodeObject(entries, Tag.NewTag());
    }

    private void BindObject(NodeObject map, object target)
    {
        var type = target.GetType();
        var properties = type.GetProperties();

        foreach (var prop in properties)
        {
            if (!prop.CanWrite)
                continue;

            var key = GetJsonKey(prop.Name);
            var entry = map.Entries.FirstOrDefault(e => e.Key.Text == key);
            if (entry.Key == null)
                continue;

            BindNodeToProperty(entry.Value, prop, target);
        }
    }

    private void BindNodeToProperty(INode node, System.Reflection.PropertyInfo prop, object target)
    {
        var propType = prop.PropertyType;

        if (propType.IsArray)
        {
            if (node is NodeArray array)
            {
                var elementType = propType.GetElementType()!;
                var arr = Array.CreateInstance(elementType, array.Children.Count);
                for (int i = 0; i < array.Children.Count; i++)
                {
                    var itemValue = ConvertMmTreeValue(array.Children[i], elementType);
                    arr.SetValue(itemValue, i);
                }
                prop.SetValue(target, arr);
            }
            return;
        }

        if (propType.IsGenericType && propType.GetGenericTypeDefinition() == typeof(List<>))
        {
            if (node is NodeArray array)
            {
                var elementType = propType.GetGenericArguments()[0];
                var listType = typeof(System.Collections.Generic.List<>).MakeGenericType(elementType);
                var list = (System.Collections.IList)Activator.CreateInstance(listType)!;
                for (int i = 0; i < array.Children.Count; i++)
                {
                    var itemValue = ConvertMmTreeValue(array.Children[i], elementType);
                    list.Add(itemValue);
                }
                prop.SetValue(target, list);
            }
            return;
        }

        if (node is NodeObject childMap)
        {
            var childTarget = Activator.CreateInstance(propType)!;
            BindObject(childMap, childTarget);
            prop.SetValue(target, childTarget);
            return;
        }

        if (node is NodeScalar scalar)
        {
            BindValue(scalar, prop, target);
            return;
        }
    }

    private object? ConvertMmTreeValue(INode node, Type targetType)
    {
        if (node is NodeScalar scalar)
        {
            if (targetType == typeof(string))
                return scalar.Text;
            if (targetType == typeof(bool))
                return scalar.Data is bool b ? b : (bool.TryParse(scalar.Text, out bool br) ? br : false);
            if (targetType == typeof(int))
                return scalar.Data is int i ? i : (scalar.Data is double d ? (int)d : (int.TryParse(scalar.Text, out int ir) ? ir : 0));
            if (targetType == typeof(long))
                return scalar.Data is long l ? l : (scalar.Data is double d ? (long)d : (long.TryParse(scalar.Text, out long lr) ? lr : 0L));
            if (targetType == typeof(double))
                return scalar.Data is double d ? d : (double.TryParse(scalar.Text, out double dr) ? dr : 0.0);
            if (targetType == typeof(float))
                return scalar.Data is float f ? f : (scalar.Data is double d ? (float)d : (float.TryParse(scalar.Text, out float fr) ? fr : 0f));
            if (targetType == typeof(byte[]))
                return scalar.Data is byte[] bytes ? bytes : (scalar.Text != null ? Convert.FromBase64String(scalar.Text) : Array.Empty<byte>());
            return scalar.Data;
        }
        if (node is NodeObject childMap)
        {
            var childTarget = Activator.CreateInstance(targetType)!;
            BindObject(childMap, childTarget);
            return childTarget;
        }
        if (node is NodeArray childArray && targetType.IsGenericType && targetType.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = targetType.GetGenericArguments()[0];
            var listType = typeof(System.Collections.Generic.List<>).MakeGenericType(elementType);
            var list = (System.Collections.IList)Activator.CreateInstance(listType)!;
            foreach (var element in childArray.Children)
            {
                list.Add(ConvertMmTreeValue(element, elementType));
            }
            return list;
        }
        return null;
    }

    private void BindArray(NodeArray array, object target)
    {
        var type = target.GetType();
        if (type.IsArray)
        {
            var elementType = type.GetElementType()!;
            var arr = Array.CreateInstance(elementType, array.Children.Count);
            for (int i = 0; i < array.Children.Count; i++)
            {
                var itemValue = ConvertMmTreeValue(array.Children[i], elementType);
                arr.SetValue(itemValue, i);
            }
            type.GetProperty("Length")!.SetValue(target, arr);
        }
        else if (type.IsGenericType && type.GetGenericTypeDefinition() == typeof(List<>))
        {
            var elementType = type.GetGenericArguments()[0];
            var list = (System.Collections.IList)target;
            for (int i = 0; i < array.Children.Count; i++)
            {
                var itemValue = ConvertMmTreeValue(array.Children[i], elementType);
                list.Add(itemValue);
            }
        }
    }

    private void BindValue(Ir.NodeScalar scalar, object target)
    {
        var type = target.GetType();
        BindValue(scalar, type.GetProperty("Item")!, target);
    }

    private void BindValue(Ir.NodeScalar scalar, System.Reflection.PropertyInfo prop, object target)
    {
        var propType = prop.PropertyType;
        var rawValue = scalar.Data;

        if (rawValue == null || scalar.Text == "null")
        {
            if (IsNullable(propType))
            {
                prop.SetValue(target, null);
            }
            return;
        }

        if (propType == typeof(string))
        {
            prop.SetValue(target, scalar.Text);
            return;
        }

        if (propType == typeof(bool))
        {
            if (rawValue is bool b)
                prop.SetValue(target, b);
            else if (bool.TryParse(scalar.Text, out bool result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(int))
        {
            if (rawValue is int i)
                prop.SetValue(target, i);
            else if (rawValue is double d)
                prop.SetValue(target, (int)d);
            else if (int.TryParse(scalar.Text, out int result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(long))
        {
            if (rawValue is long l)
                prop.SetValue(target, l);
            else if (rawValue is double d)
                prop.SetValue(target, (long)d);
            else if (long.TryParse(scalar.Text, out long result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(double))
        {
            if (rawValue is double d)
                prop.SetValue(target, d);
            else if (double.TryParse(scalar.Text, out double result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(float))
        {
            if (rawValue is float f)
                prop.SetValue(target, f);
            else if (float.TryParse(scalar.Text, out float result))
                prop.SetValue(target, result);
            return;
        }

        if (propType == typeof(byte[]))
        {
            if (scalar.Data is byte[] bytes)
                prop.SetValue(target, bytes);
            else if (scalar.Text != null)
                prop.SetValue(target, Convert.FromBase64String(scalar.Text));
            return;
        }

        if (propType == typeof(System.DateTime))
        {
            if (scalar.Text != null && System.DateTime.TryParse(scalar.Text, out var dt))
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