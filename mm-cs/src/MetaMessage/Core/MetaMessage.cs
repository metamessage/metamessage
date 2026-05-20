using MetaMessage.Jsonc;
using JsoncParser = MetaMessage.Jsonc.Jsonc;

namespace MetaMessage.Core;

public static class MetaMessage
{
    public static byte[] Encode(object obj)
    {
        return ReflectMmEncoder.Encode(obj);
    }

    public static void Decode(byte[] data, object target)
    {
        ReflectMmBinder.Bind(data, target);
    }

    public static T Decode<T>(byte[] data) where T : new()
    {
        var obj = new T();
        ReflectMmBinder.Bind(data, obj);
        return obj;
    }

    public static IMmTree DecodeToTree(byte[] data)
    {
        var decoder = new WireDecoder(data);
        return decoder.Decode();
    }

    public static byte[] EncodeTree(IMmTree tree)
    {
        var encoder = new WireEncoder();
        EncodeTreeValue(encoder, tree, MmTag.Empty());
        return encoder.ToByteArray();
    }

    public static string DecodeToJsonc(byte[] data)
    {
        var tree = DecodeToTree(data);
        var node = TreeToJsoncNode(tree);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static byte[] EncodeFromJsonc(string jsonc)
    {
        var node = JsoncParser.ParseFromString(jsonc);
        return EncodeFromJsoncNode(node);
    }

    public static string ValueToJsonc(object value)
    {
        var binder = new JsoncBinder();
        var node = binder.StructToNode(value);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static object? JsoncToValue(string jsonc)
    {
        var node = JsoncParser.ParseFromString(jsonc);
        return ExtractValueFromJsoncNode(node);
    }

    public static ValidationResult Validate(object value, MmTag tag)
    {
        return Validator.Validate(value, tag);
    }

    private static void EncodeTreeValue(WireEncoder encoder, IMmTree tree, MmTag inherited)
    {
        switch (tree)
        {
            case MmScalar scalar:
                EncodeScalarTree(encoder, scalar, inherited);
                break;
            case MmArray array:
                EncodeArrayTree(encoder, array, inherited);
                break;
            case MmMap map:
                EncodeMapTree(encoder, map, inherited);
                break;
        }
    }

    private static void EncodeScalarTree(WireEncoder encoder, MmScalar scalar, MmTag inherited)
    {
        var tag = scalar.Tag?.Copy() ?? inherited.Copy();
        if (scalar.Data == null || tag.IsNull)
        {
            tag.IsNull = true;
            var nullPayload = new WireEncoder();
            EncodeNullScalarPayload(nullPayload, tag);
            encoder.EncodeTaggedPayload(nullPayload.ToByteArray(), tag.ToBytes());
            return;
        }

        var payload = new WireEncoder();
        EncodeScalarPayload(payload, scalar.Data, tag);
        encoder.EncodeTaggedPayload(payload.ToByteArray(), tag.ToBytes());
    }

    private static void EncodeNullScalarPayload(WireEncoder encoder, MmTag tag)
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

    private static void EncodeScalarPayload(WireEncoder encoder, object value, MmTag tag)
    {
        switch (tag.Type)
        {
            case ValueType.BOOL:
                encoder.EncodeBool((bool)value);
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
            case ValueType.FLOAT64:
            case ValueType.DECIMAL:
                encoder.EncodeFloatString(value.ToString()!);
                break;
            case ValueType.STRING:
            case ValueType.EMAIL:
            case ValueType.URL:
                encoder.EncodeString(value as string ?? value.ToString()!);
                break;
            case ValueType.BYTES:
                encoder.EncodeBytes(value as byte[] ?? Array.Empty<byte>());
                break;
            case ValueType.BIGINT:
                encoder.EncodeBigIntDecimal(value.ToString()!);
                break;
            case ValueType.UUID:
                encoder.EncodeBytes(UuidToBytes(value.ToString()!));
                break;
            case ValueType.DATETIME:
                encoder.EncodeInt64(TimeUtil.EpochSeconds((DateTime)value));
                break;
            case ValueType.DATE:
                encoder.EncodeInt64(TimeUtil.DaysSinceEpochUtc((DateTime)value));
                break;
            case ValueType.TIME:
                encoder.EncodeInt64(TimeUtil.SecondsOfDay((DateTime)value));
                break;
            case ValueType.ENUM:
                encoder.EncodeInt64(Convert.ToInt64(value));
                break;
            default:
                if (value is string s)
                {
                    encoder.EncodeString(s);
                }
                else if (value is byte[] b)
                {
                    encoder.EncodeBytes(b);
                }
                else if (value is bool bb)
                {
                    encoder.EncodeBool(bb);
                }
                else if (value is double || value is float)
                {
                    encoder.EncodeFloatString(value.ToString()!);
                }
                else
                {
                    encoder.EncodeInt64(Convert.ToInt64(value));
                }
                break;
        }
    }

    private static void EncodeArrayTree(WireEncoder encoder, MmArray array, MmTag inherited)
    {
        var tag = array.Tag?.Copy() ?? inherited.Copy();
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();

        foreach (var child in array.Children)
        {
            elementEncoder.Reset();
            var itemTag = MmTag.Empty();
            itemTag.InheritFromArrayParent(tag);
            EncodeTreeValue(elementEncoder, child, itemTag);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        elementEncoder.Reset();
        elementEncoder.EncodeArrayPayload(body.ToArray());
        encoder.EncodeTaggedPayload(elementEncoder.ToByteArray(), tag.ToBytes());
    }

    private static void EncodeMapTree(WireEncoder encoder, MmMap map, MmTag inherited)
    {
        var tag = map.Tag?.Copy() ?? inherited.Copy();
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var entry in map.Entries)
        {
            tmp.Reset();
            if (!TryEncodeSimpleByName(tmp, entry.Key.Text))
            {
                tmp.EncodeString(entry.Key.Text);
            }
            keysPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            var valueTag = MmTag.Empty();
            EncodeTreeValue(tmp, entry.Value, valueTag);
            valsPacked.WriteAll(tmp.ToByteArray());
        }

        tmp.Reset();
        tmp.EncodeArrayPayload(keysPacked.ToArray());
        var mapBody = new GrowableByteBuf();
        mapBody.WriteAll(tmp.ToByteArray());
        mapBody.WriteAll(valsPacked.ToArray());

        tmp.Reset();
        tmp.EncodeObjectPayload(mapBody.ToArray());
        encoder.EncodeTaggedPayload(tmp.ToByteArray(), tag.ToBytes());
    }

    private static bool TryEncodeSimpleByName(WireEncoder encoder, string name)
    {
        return WireEncoder.TryEncodeSimpleByName(encoder, name);
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

    private static JsoncTag? _treeTag; // re-entry guard

    private static byte[] EncodeFromJsoncNode(IJsoncNode node)
    {
        var encoder = new WireEncoder();
        EncodeJsoncNodeValue(encoder, node);
        return encoder.ToByteArray();
    }

    private static void EncodeJsoncNodeValue(WireEncoder encoder, IJsoncNode node)
    {
        var payload = new WireEncoder();
        switch (node)
        {
            case JsoncValue jsoncValue:
                EncodeJsoncScalarPayload(payload, jsoncValue);
                break;
            case JsoncArray jsoncArray:
                EncodeJsoncArrayPayload(payload, jsoncArray);
                break;
            case JsoncObject jsoncObject:
                EncodeJsoncObjectPayload(payload, jsoncObject);
                break;
        }

        var mmTag = MmTag.Empty();
        if (node?.Tag != null)
        {
            mmTag = ConvertJsoncTagToMmTag(node.Tag);
        }
        encoder.EncodeTaggedPayload(payload.ToByteArray(), mmTag.ToBytes());
    }

    private static MmTag ConvertJsoncTagToMmTag(JsoncTag tag)
    {
        var mmTag = new MmTag();
        mmTag.Type = (Core.ValueType)(int)tag.Type;
        mmTag.Desc = tag.Desc ?? string.Empty;
        mmTag.Nullable = tag.Nullable;
        mmTag.IsNull = tag.IsNull;
        mmTag.Min = tag.MinValue ?? string.Empty;
        mmTag.Max = tag.MaxValue ?? string.Empty;
        mmTag.DefaultValue = tag.DefaultValue ?? string.Empty;
        return mmTag;
    }

    private static JsoncTag ConvertMmTagToJsoncTag(MmTag mmTag)
    {
        var tag = new JsoncTag();
        tag.Type = (Jsonc.ValueType)(int)mmTag.Type;
        tag.Desc = mmTag.Desc;
        tag.Nullable = mmTag.Nullable;
        tag.IsNull = mmTag.IsNull;
        if (!string.IsNullOrEmpty(mmTag.Min) && mmTag.Min != "0")
            tag.MinValue = mmTag.Min;
        if (!string.IsNullOrEmpty(mmTag.Max) && mmTag.Max != "0")
            tag.MaxValue = mmTag.Max;
        if (!string.IsNullOrEmpty(mmTag.DefaultValue))
            tag.DefaultValue = mmTag.DefaultValue;
        return tag;
    }

    private static void EncodeJsoncScalarPayload(WireEncoder encoder, JsoncValue value)
    {
        switch (value.TokenType)
        {
            case JsoncTokenType.Null:
                encoder.EncodeSimple(SimpleValue.NULL_STRING);
                break;
            case JsoncTokenType.True:
                encoder.EncodeSimple(SimpleValue.TRUE);
                break;
            case JsoncTokenType.False:
                encoder.EncodeSimple(SimpleValue.FALSE);
                break;
            case JsoncTokenType.Number:
                if (value.Value is double d)
                {
                    if (double.IsNaN(d) || double.IsInfinity(d))
                    {
                        encoder.EncodeFloatString(d.ToString());
                    }
                    else if (d == Math.Floor(d) && d >= long.MinValue && d <= long.MaxValue)
                    {
                        encoder.EncodeInt64((long)d);
                    }
                    else
                    {
                        encoder.EncodeFloatString(d.ToString());
                    }
                }
                else if (value.Value is int i)
                {
                    encoder.EncodeInt64(i);
                }
                else if (value.Value is long l)
                {
                    encoder.EncodeInt64(l);
                }
                else if (value.Value is float f)
                {
                    encoder.EncodeFloatString(f.ToString());
                }
                else
                {
                    encoder.EncodeFloatString(value.Value?.ToString() ?? "0");
                }
                break;
            case JsoncTokenType.String:
            default:
                encoder.EncodeString(value.Value?.ToString() ?? "");
                break;
        }
    }

    private static void EncodeJsoncArrayPayload(WireEncoder encoder, JsoncArray array)
    {
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();

        foreach (var element in array.Elements)
        {
            elementEncoder.Reset();
            EncodeJsoncNodeValue(elementEncoder, element);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        elementEncoder.Reset();
        elementEncoder.EncodeArrayPayload(body.ToArray());
        encoder.EncodeArrayPayload(elementEncoder.ToByteArray());
    }

    private static void EncodeJsoncObjectPayload(WireEncoder encoder, JsoncObject obj)
    {
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var kvp in obj.Fields)
        {
            tmp.Reset();
            if (!TryEncodeSimpleByName(tmp, kvp.Key))
            {
                tmp.EncodeString(kvp.Key);
            }
            keysPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            EncodeJsoncNodeValue(tmp, kvp.Value);
            valsPacked.WriteAll(tmp.ToByteArray());
        }

        tmp.Reset();
        tmp.EncodeArrayPayload(keysPacked.ToArray());
        var mapBody = new GrowableByteBuf();
        mapBody.WriteAll(tmp.ToByteArray());
        mapBody.WriteAll(valsPacked.ToArray());

        encoder.EncodeObjectPayload(mapBody.ToArray());
    }

    private static IJsoncNode TreeToJsoncNode(IMmTree tree)
    {
        switch (tree)
        {
            case MmScalar scalar:
                return ScalarToJsoncValue(scalar);
            case MmArray array:
                return ArrayToJsoncArray(array);
            case MmMap map:
                return MapToJsoncObject(map);
            default:
                return new JsoncValue { Value = null, TokenType = JsoncTokenType.Null };
        }
    }

    private static IJsoncNode ScalarToJsoncValue(MmScalar scalar)
    {
        JsoncValue result;
        if (scalar.Tag.IsNull)
        {
            switch (scalar.Tag.Type)
            {
                case ValueType.BOOL:
                    result = new JsoncValue { Value = false, TokenType = JsoncTokenType.False };
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
                case ValueType.BIGINT:
                case ValueType.FLOAT32:
                case ValueType.FLOAT64:
                case ValueType.DECIMAL:
                    result = new JsoncValue { Value = 0L, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.STRING:
                case ValueType.EMAIL:
                case ValueType.URL:
                case ValueType.ENUM:
                case ValueType.DATETIME:
                case ValueType.DATE:
                case ValueType.TIME:
                case ValueType.UUID:
                    result = new JsoncValue { Value = "", TokenType = JsoncTokenType.String };
                    break;
                case ValueType.BYTES:
                    result = new JsoncValue { Value = "", TokenType = JsoncTokenType.String };
                    break;
                default:
                    result = new JsoncValue { Value = "", TokenType = JsoncTokenType.String };
                    break;
            }
        }
        else if (scalar.Data == null)
        {
            result = new JsoncValue { Value = null, TokenType = JsoncTokenType.Null };
        }
        else
        {
            switch (scalar.Tag.Type)
            {
                case ValueType.BOOL:
                    result = new JsoncValue
                    {
                        Value = scalar.Data,
                        TokenType = (bool)scalar.Data ? JsoncTokenType.True : JsoncTokenType.False
                    };
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
                    result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.FLOAT32:
                case ValueType.FLOAT64:
                case ValueType.DECIMAL:
                    result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.BIGINT:
                case ValueType.DATETIME:
                case ValueType.DATE:
                case ValueType.TIME:
                case ValueType.UUID:
                case ValueType.EMAIL:
                case ValueType.URL:
                case ValueType.STRING:
                    result = new JsoncValue { Value = scalar.Data?.ToString(), TokenType = JsoncTokenType.String };
                    break;
                case ValueType.BYTES:
                    result = new JsoncValue
                    {
                        Value = scalar.Data is byte[] b ? Convert.ToBase64String(b) : scalar.Data?.ToString(),
                        TokenType = JsoncTokenType.String
                    };
                    break;
                case ValueType.ENUM:
                    result = new JsoncValue { Value = scalar.Data?.ToString(), TokenType = JsoncTokenType.String };
                    break;
                default:
                    if (scalar.Data is string s)
                    {
                        result = new JsoncValue { Value = s, TokenType = JsoncTokenType.String };
                    }
                    else if (scalar.Data is bool bb)
                    {
                        result = new JsoncValue
                        {
                            Value = bb,
                            TokenType = bb ? JsoncTokenType.True : JsoncTokenType.False
                        };
                    }
                    else if (scalar.Data is double || scalar.Data is float)
                    {
                        result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    }
                    else if (scalar.Data is int || scalar.Data is long)
                    {
                        result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    }
                    else
                    {
                        result = new JsoncValue { Value = scalar.Data?.ToString(), TokenType = JsoncTokenType.String };
                    }
                    break;
            }
        }

        if (scalar.Tag != null)
        {
            result.Tag = ConvertMmTagToJsoncTag(scalar.Tag);
        }
        return result;
    }

    private static IJsoncNode ArrayToJsoncArray(MmArray array)
    {
        var jsoncArray = new JsoncArray();
        foreach (var child in array.Children)
        {
            jsoncArray.Add(TreeToJsoncNode(child));
        }
        if (array.Tag != null)
        {
            jsoncArray.Tag = ConvertMmTagToJsoncTag(array.Tag);
        }
        return jsoncArray;
    }

    private static IJsoncNode MapToJsoncObject(MmMap map)
    {
        var jsoncObject = new JsoncObject();
        foreach (var entry in map.Entries)
        {
            var key = entry.Key.Text;
            jsoncObject.Add(key, TreeToJsoncNode(entry.Value));
        }
        return jsoncObject;
    }

    private static object? ExtractValueFromJsoncNode(IJsoncNode node)
    {
        switch (node)
        {
            case JsoncValue value:
                return ExtractScalarValue(value);
            case JsoncArray array:
                {
                    var list = new List<object?>();
                    foreach (var element in array.Elements)
                    {
                        list.Add(ExtractValueFromJsoncNode(element));
                    }
                    return list;
                }
            case JsoncObject obj:
                {
                    var dict = new Dictionary<string, object?>();
                    foreach (var kvp in obj.Fields)
                    {
                        dict[kvp.Key] = ExtractValueFromJsoncNode(kvp.Value);
                    }
                    return dict;
                }
            default:
                return null;
        }
    }

    private static object? ExtractScalarValue(JsoncValue value)
    {
        switch (value.TokenType)
        {
            case JsoncTokenType.Null:
                return null;
            case JsoncTokenType.True:
                return true;
            case JsoncTokenType.False:
                return false;
            case JsoncTokenType.Number:
                return value.Value;
            case JsoncTokenType.String:
            default:
                return value.Value?.ToString();
        }
    }
}