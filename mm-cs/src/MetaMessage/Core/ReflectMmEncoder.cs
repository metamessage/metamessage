using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;
using System.Reflection;
using System.Collections;

namespace MetaMessage.Core;

public static class ReflectMmEncoder
{
    private const int MaxDepth = 32;

    public static INode ValueToNode(object root, string tagStr)
    {
        Tag? tag = null;
        if (!string.IsNullOrEmpty(tagStr))
        {
            tag = Tag.Parse(tagStr);
        }

        tag ??= Tag.NewTag();

        return ValueToNodeRecursive(root, tag, 0, "", false);
    }

    private static INode ValueToNodeRecursive(object? v, Tag tag, int depth, string path, bool example)
    {
        if (depth > MaxDepth)
        {
            throw new Exception($"max depth: {MaxDepth}");
        }

        tag ??= Tag.NewTag();

        if (v == null)
        {
            if (tag.Type != ValueType.Unknown)
            {
                tag.IsNull = true;
                return new NodeScalar(null, "null", tag.Copy());
            }

            // Try to infer from property context (will be handled by AnyToNode)
            throw new Exception("invalid input: v is null with unknown type");
        }

        object? data = null;
        string text = "null";

        var type = v.GetType();

        if (type == typeof(byte[]))
        {
            var bytes = (byte[])v;
            if (tag.Type == ValueType.Unknown)
            {
                tag.Type = ValueType.Vec;
            }

            switch (tag.Type)
            {
                case ValueType.Bytes:
                    data = bytes;
                    text = Convert.ToBase64String(bytes);
                    break;
                case ValueType.Vec:
                    return AnyToNode(v, tag, depth, path, false);
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }

            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(bool))
        {
            var b = (bool)v;
            if (tag.Type == ValueType.Unknown)
            {
                tag.Type = ValueType.Bool;
            }
            switch (tag.Type)
            {
                case ValueType.Bool:
                    data = b;
                    text = b ? "true" : "false";
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(int))
        {
            var val = (int)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.I;
            switch (tag.Type)
            {
                case ValueType.I:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(long))
        {
            var val = (long)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.I64;
            switch (tag.Type)
            {
                case ValueType.I64:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(short))
        {
            var val = (short)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.I16;
            switch (tag.Type)
            {
                case ValueType.I16:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(sbyte))
        {
            var val = (sbyte)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.I8;
            switch (tag.Type)
            {
                case ValueType.I8:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(byte))
        {
            var val = (byte)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.U8;
            switch (tag.Type)
            {
                case ValueType.U8:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(uint))
        {
            var val = (uint)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.U;
            switch (tag.Type)
            {
                case ValueType.U:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(ushort))
        {
            var val = (ushort)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.U16;
            switch (tag.Type)
            {
                case ValueType.U16:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(ulong))
        {
            var val = (ulong)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.U64;
            switch (tag.Type)
            {
                case ValueType.U64:
                    data = val;
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(float))
        {
            var val = (float)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.F32;
            switch (tag.Type)
            {
                case ValueType.F32:
                    if (float.IsInfinity(val) || float.IsNaN(val))
                        throw new Exception($"unsupported float value: {val}");
                    data = val;
                    text = val.ToString("G");
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(double))
        {
            var val = (double)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.F64;
            switch (tag.Type)
            {
                case ValueType.F64:
                    if (double.IsInfinity(val) || double.IsNaN(val))
                        throw new Exception($"unsupported double value: {val}");
                    data = val;
                    text = val.ToString("G");
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(decimal))
        {
            var val = (decimal)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.Decimal;
            switch (tag.Type)
            {
                case ValueType.Decimal:
                    data = val.ToString();
                    text = val.ToString();
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(string))
        {
            var val = (string)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.Str;
            switch (tag.Type)
            {
                case ValueType.Str:
                    data = val;
                    text = val;
                    break;
                case ValueType.Email:
                    data = val;
                    text = val;
                    break;
                case ValueType.Url:
                    data = val;
                    text = val;
                    break;
                case ValueType.Enums:
                    data = val;
                    text = val;
                    break;
                case ValueType.Uuid:
                    data = val;
                    text = val;
                    break;
                case ValueType.Decimal:
                    data = val;
                    text = val;
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        if (type == typeof(DateTime))
        {
            var val = (DateTime)v;
            if (tag.Type == ValueType.Unknown) tag.Type = ValueType.Datetime;
            switch (tag.Type)
            {
                case ValueType.Datetime:
                    data = val;
                    text = val.ToString("yyyy-MM-dd HH:mm:ss");
                    break;
                case ValueType.Date:
                    data = val;
                    text = val.ToString("yyyy-MM-dd");
                    break;
                case ValueType.Time:
                    data = val;
                    text = val.ToString("HH:mm:ss");
                    break;
                default:
                    throw new Exception($"unsupported type: {tag.Type}");
            }
            return new NodeScalar(data, text, tag.Copy());
        }

        return AnyToNode(v, tag, depth, path, false);
    }

    private static INode AnyToNode(object obj, Tag tag, int depth, string path, bool example)
    {
        depth++;
        if (depth > MaxDepth)
        {
            throw new Exception($"max depth: {MaxDepth}");
        }

        var val = obj;
        var typ = obj.GetType();

        if (tag == null)
        {
            tag = Tag.NewTag();
        }

        var valType = val.GetType();

        if (valType.IsGenericType && valType.GetGenericTypeDefinition() == typeof(Nullable<>))
        {
            tag.Nullable = true;
            var underlyingType = Nullable.GetUnderlyingType(valType);
            var prop = valType.GetProperty("HasValue");
            var hasValue = (bool)(prop!.GetValue(val) ?? false);
            if (!hasValue)
            {
                tag.IsNull = true;
                var nullScalar = new NodeScalar(null, "null", tag.Copy());
                nullScalar.Tag.IsNull = true;
                return nullScalar;
            }
            var getValue = valType.GetMethod("GetValueOrDefault", Type.EmptyTypes);
            val = getValue!.Invoke(val, null)!;
            typ = underlyingType;
        }

        if (tag.ToString() == "")
        {
            var mmAttr = typ!.GetCustomAttribute(typeof(MM), false) as MM;
            if (mmAttr != null)
            {
                var tagNode = mmAttr.ToTag();
                tag = Tag.MergeTag(tag, tagNode);
            }
        }

        tag.Type = ValueType.Obj;
        tag.Name = CamelToSnake.Convert(typ!.Name ?? "");
        if (!string.IsNullOrEmpty(tag.Name))
        {
            if (string.IsNullOrEmpty(path))
            {
                path = tag.Name;
            }
            else
            {
                path = $"{path}.{tag.Name}";
            }
        }

        if (typ!.IsGenericType && typ.GetGenericTypeDefinition() == typeof(Dictionary<,>))
        {
            tag.Type = ValueType.Map;
            var dict = (IDictionary)val!;
            var entries = new List<KeyValuePair<NodeScalar, INode>>();

            foreach (var key in dict.Keys)
            {
                var keyStr = key.ToString() ?? "";
                keyStr = CamelToSnake.Convert(keyStr);

                var tagItem = Tag.NewTag();
                tagItem.Inherit(tag);
                tagItem.Name = keyStr;

                var p = $"{path}[{keyStr}]";
                var valNode = ValueToNodeRecursive(dict[key], tagItem, depth, p, false);

                entries.Add(new KeyValuePair<NodeScalar, INode>(
                    new NodeScalar(keyStr, keyStr, Tag.Empty()),
                    valNode));
            }

            return new NodeObject(entries, tag.Copy());
        }

        if (val is IList list)
        {
            tag.Type = ValueType.Vec;

            var items = new List<INode>();
            foreach (var item in list)
            {
                var tagItem = Tag.NewTag();
                tagItem.Inherit(tag);

                var p = $"{path}[{items.Count}]";
                var itemNode = ValueToNodeRecursive(item, tagItem, depth, p, false);
                items.Add(itemNode);
            }

            return new NodeArray(items, tag.Copy());
        }

        if (typ!.IsArray && typ != typeof(byte[]))
        {
            tag.Type = ValueType.Arr;
            var arr = (Array)val!;
            tag.Size = arr.Length;

            var items = new List<INode>();
            for (int i = 0; i < arr.Length; i++)
            {
                var tagItem = Tag.NewTag();
                tagItem.Inherit(tag);

                var p = $"{path}[{i}]";
                var itemNode = ValueToNodeRecursive(arr.GetValue(i), tagItem, depth, p, false);
                items.Add(itemNode);
            }

            return new NodeArray(items, tag.Copy());
        }

        if (typ!.IsClass || (typ.IsValueType && !typ.IsPrimitive && !typ.IsEnum))
        {
            tag.Type = ValueType.Obj;

            var fields = new List<KeyValuePair<NodeScalar, INode>>();
            var properties = typ.GetProperties(BindingFlags.Public | BindingFlags.Instance);

            foreach (var property in properties)
            {
                if (!property.CanRead)
                    continue;

                var mmAttr = property.GetCustomAttribute(typeof(MM), false) as MM;
                if (mmAttr != null && mmAttr.Name == "-")
                    continue;

                var fieldKey = CamelToSnake.Convert(property.Name);
                if (mmAttr != null && !string.IsNullOrEmpty(mmAttr.Name) && mmAttr.Name != "-")
                {
                    fieldKey = mmAttr.Name;
                }

                var tagItem = Tag.NewTag();
                tagItem.Name = fieldKey;

                if (mmAttr != null)
                {
                    var attrTag = mmAttr.ToTag();
                    if (attrTag.Type != ValueType.Unknown)
                        tagItem.Type = attrTag.Type;
                    if (attrTag.ChildType != ValueType.Unknown)
                        tagItem.ChildType = attrTag.ChildType;
                    if (!string.IsNullOrEmpty(attrTag.Name))
                        tagItem.Name = attrTag.Name;
                    tagItem = Tag.MergeTag(tagItem, attrTag);
                }

                var p = $"{path}.{fieldKey}";
                var propVal = property.GetValue(val);

                if (tagItem.Type == ValueType.Unknown)
                {
                    tagItem.Type = InferTypeFromPropertyType(property.PropertyType);
                }

                var fieldNode = ValueToNodeRecursive(propVal, tagItem, depth, p, false);

                fields.Add(new KeyValuePair<NodeScalar, INode>(
                    new NodeScalar(fieldKey, fieldKey, Tag.Empty()),
                    fieldNode));
            }

            return new NodeObject(fields, tag.Copy());
        }

        if (typ!.IsEnum)
        {
            tag.Type = ValueType.Enums;
            var intVal = Convert.ToInt64(val);
            return new NodeScalar(intVal, intVal.ToString(), tag.Copy());
        }

        throw new Exception($"unsupported type: {typ.FullName}");
    }

    private static ValueType InferTypeFromPropertyType(Type propertyType)
    {
        var type = Nullable.GetUnderlyingType(propertyType) ?? propertyType;
        if (type == typeof(string))
            return ValueType.Str;
        if (type == typeof(bool))
            return ValueType.Bool;
        if (type == typeof(byte))
            return ValueType.U8;
        if (type == typeof(sbyte))
            return ValueType.I8;
        if (type == typeof(short))
            return ValueType.I16;
        if (type == typeof(ushort))
            return ValueType.U16;
        if (type == typeof(int))
            return ValueType.I;
        if (type == typeof(uint))
            return ValueType.U;
        if (type == typeof(long))
            return ValueType.I64;
        if (type == typeof(ulong))
            return ValueType.U64;
        if (type == typeof(float))
            return ValueType.F32;
        if (type == typeof(double))
            return ValueType.F64;
        if (type == typeof(decimal))
            return ValueType.Decimal;
        if (type == typeof(DateTime))
            return ValueType.Datetime;
        if (type == typeof(byte[]))
            return ValueType.Bytes;
        if (type.IsEnum)
            return ValueType.Enums;
        if (type.IsClass || (type.IsValueType && !type.IsPrimitive && !type.IsEnum))
            return ValueType.Obj;
        if (typeof(IList).IsAssignableFrom(type))
            return ValueType.Vec;
        return ValueType.Obj;
    }
}