using ValueType = MetaMessage.Ir.ValueType;
using MetaMessage.Jsonc;
using MetaMessage.Ir;
using JsoncParser = MetaMessage.Jsonc.Jsonc;
using System.Text.Json;

namespace MetaMessage.Core;

public static class MetaMessage
{
    public static byte[] Encode(object obj)
    {
        var node = ReflectMmEncoder.ValueToNode(obj, "");
        return EncodeTree(node);
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

    public static IMmTree Decode(byte[] data)
    {
        var decoder = new WireDecoder(data);
        return decoder.Decode();
    }

    public static IMmTree DecodeToTree(byte[] data)
    {
        return Decode(data);
    }

    public static IMmTree ParseFromJSONC(string input)
    {
        return JsoncToTree(input);
    }

    public static byte[] EncodeTree(IMmTree tree)
    {
        var encoder = new WireEncoder();
        EncodeTreeValue(encoder, tree, Tag.Empty());
        return encoder.ToByteArray();
    }

    public static byte[] FromValue(object value, string tagString)
    {
        var node = ReflectMmEncoder.ValueToNode(value, tagString);
        return EncodeTree(node);
    }

    public static byte[] FromJSONC(string jsonc)
    {
        return EncodeFromJsonc(jsonc);
    }

    public static void BindFromJSONC(string input, object output)
    {
        var tree = ParseFromJSONC(input);
        var encoded = EncodeTree(tree);
        ReflectMmBinder.Bind(encoded, output);
    }

    public static string ValueToJSONC(object value, string name)
    {
        var node = ReflectMmEncoder.ValueToNode(value, name);
        return TreeToJsoncString(node);
    }

    public static string ValueToJsonc(object value)
    {
        return ValueToJSONC(value, "");
    }

    public static void PrintJSONC(IMmTree node)
    {
        Console.WriteLine(TreeToJsoncString(node));
    }

    public static string Dump(IMmTree node)
    {
        var options = new JsonSerializerOptions { WriteIndented = true };
        return JsonSerializer.Serialize(node, options);
    }

    public static string DecodeToJsonc(byte[] data)
    {
        var tree = DecodeToTree(data);
        return TreeToJsoncString(tree);
    }

    public static byte[] EncodeFromJsonc(string jsonc)
    {
        var node = JsoncParser.ParseFromString(jsonc);
        return EncodeFromJsoncNode(node);
    }

    public static object? JsoncToValue(string jsonc)
    {
        var node = JsoncParser.ParseFromString(jsonc);
        return ExtractValueFromJsoncNode(node);
    }

    public static ValidationResult Validate(object value, Tag tag)
    {
        return Validator.Validate(value, tag);
    }

    private static string TreeToJsoncString(IMmTree tree)
    {
        var node = TreeToJsoncNode(tree);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    private static IMmTree JsoncToTree(string input)
    {
        var jsoncNode = JsoncParser.ParseFromString(input);
        return JsoncNodeToTree(jsoncNode);
    }

    private static IMmTree JsoncNodeToTree(IJsoncNode node)
    {
        switch (node)
        {
            case JsoncValue value:
                return JsoncValueToScalar(value);
            case JsoncArray array:
                {
                    var children = new List<IMmTree>();
                    foreach (var element in array.Elements)
                    {
                        children.Add(JsoncNodeToTree(element));
                    }
                    return new MmArray(children, array.Tag?.Copy() ?? Tag.Empty());
                }
            case JsoncObject obj:
                {
                    var entries = new List<KeyValuePair<MmScalar, IMmTree>>();
                    foreach (var kvp in obj.Fields)
                    {
                        var keyScalar = new MmScalar(kvp.Key, kvp.Key, Tag.Empty());
                        entries.Add(new KeyValuePair<MmScalar, IMmTree>(keyScalar, JsoncNodeToTree(kvp.Value)));
                    }
                    return new MmMap(entries, obj.Tag?.Copy() ?? Tag.Empty());
                }
            default:
                throw new Exception($"unsupported jsonc node type: {node?.GetType()}");
        }
    }

    private static MmScalar JsoncValueToScalar(JsoncValue value)
    {
        object? data = null;
        string text = "null";

        switch (value.TokenType)
        {
            case JsoncTokenType.Null:
                data = null;
                text = "null";
                break;
            case JsoncTokenType.True:
                data = true;
                text = "true";
                break;
            case JsoncTokenType.False:
                data = false;
                text = "false";
                break;
            case JsoncTokenType.Number:
                data = value.Value;
                text = value.Value?.ToString() ?? "0";
                break;
            case JsoncTokenType.String:
            default:
                data = value.Value?.ToString() ?? "";
                text = value.Value?.ToString() ?? "";
                break;
        }

        return new MmScalar(data, text, value.Tag?.Copy() ?? Tag.Empty());
    }

    private static void EncodeTreeValue(WireEncoder encoder, IMmTree tree, Tag inherited)
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

    private static void EncodeScalarTree(WireEncoder encoder, MmScalar scalar, Tag inherited)
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

    private static void EncodeNullScalarPayload(WireEncoder encoder, Tag tag)
    {
        switch (tag.Type)
        {
            case ValueType.Bool:
                encoder.EncodeSimple(SimpleValue.NULL_BOOL);
                return;
            case ValueType.I:
            case ValueType.I8:
            case ValueType.I16:
            case ValueType.I32:
            case ValueType.I64:
            case ValueType.U:
            case ValueType.U16:
            case ValueType.U32:
            case ValueType.U64:
                encoder.EncodeSimple(SimpleValue.NULL_INT);
                return;
            case ValueType.F32:
            case ValueType.F64:
            case ValueType.Decimal:
                encoder.EncodeSimple(SimpleValue.NULL_FLOAT);
                return;
            case ValueType.Str:
            case ValueType.Email:
            case ValueType.Url:
                encoder.EncodeSimple(SimpleValue.NULL_STRING);
                return;
            case ValueType.Bytes:
                encoder.EncodeSimple(SimpleValue.NULL_BYTES);
                return;
            default:
                encoder.EncodeSimple(SimpleValue.NULL_STRING);
                return;
        }
    }

    private static void EncodeScalarPayload(WireEncoder encoder, object value, Tag tag)
    {
        switch (tag.Type)
        {
            case ValueType.Bool:
                encoder.EncodeBool((bool)value);
                break;
            case ValueType.I:
            case ValueType.I8:
            case ValueType.I16:
            case ValueType.I32:
            case ValueType.I64:
            case ValueType.U:
            case ValueType.U16:
            case ValueType.U32:
                encoder.EncodeInt64(Convert.ToInt64(value));
                break;
            case ValueType.U64:
                if (value is ulong ulVal && ulVal > long.MaxValue)
                {
                    encoder.EncodeUInt64(ulVal);
                }
                else
                {
                    encoder.EncodeInt64(Convert.ToInt64(value));
                }
                break;
            case ValueType.F32:
            case ValueType.F64:
            case ValueType.Decimal:
                encoder.EncodeFloatString(value.ToString()!);
                break;
            case ValueType.Str:
            case ValueType.Email:
            case ValueType.Url:
                encoder.EncodeString(value as string ?? value.ToString()!);
                break;
            case ValueType.Bytes:
                encoder.EncodeBytes(value as byte[] ?? Array.Empty<byte>());
                break;
            case ValueType.Media:
                if (value is byte[] mediaBytes)
                {
                    encoder.EncodeBytes(mediaBytes);
                }
                else if (value is string mediaStr)
                {
                    encoder.EncodeBytes(Convert.FromBase64String(mediaStr));
                }
                else
                {
                    encoder.EncodeBytes(Array.Empty<byte>());
                }
                break;
            case ValueType.Bigint:
                encoder.EncodeBigIntDecimal(value.ToString()!);
                break;
            case ValueType.Uuid:
                encoder.EncodeBytes(UuidToBytes(value.ToString()!));
                break;
            case ValueType.Datetime:
                encoder.EncodeInt64(TimeUtil.EpochSeconds((DateTime)value));
                break;
            case ValueType.Date:
                encoder.EncodeInt64(TimeUtil.DaysSinceEpochUtc((DateTime)value));
                break;
            case ValueType.Time:
                encoder.EncodeInt64(TimeUtil.SecondsOfDay((DateTime)value));
                break;
            case ValueType.Enums:
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

    private static void EncodeArrayTree(WireEncoder encoder, MmArray array, Tag inherited)
    {
        var tag = array.Tag?.Copy() ?? inherited.Copy();
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();

        foreach (var child in array.Children)
        {
            elementEncoder.Reset();
            var itemTag = Tag.Empty();
            itemTag.InheritFromArrayParent(tag);
            EncodeTreeValue(elementEncoder, child, itemTag);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        elementEncoder.Reset();
        elementEncoder.EncodeArrayPayload(body.ToArray());
        encoder.EncodeTaggedPayload(elementEncoder.ToByteArray(), tag.ToBytes());
    }

    private static void EncodeMapTree(WireEncoder encoder, MmMap map, Tag inherited)
    {
        var tag = map.Tag?.Copy() ?? inherited.Copy();
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var entry in map.Entries)
        {
            tmp.Reset();
            tmp.EncodeString(entry.Key.Text);
            keysPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            var valueTag = Tag.Empty();
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

    private static byte[] EncodeFromJsoncNode(IJsoncNode node)
    {
        var encoder = new WireEncoder();
        EncodeJsoncNodeValue(encoder, node);
        return encoder.ToByteArray();
    }

    private static void EncodeJsoncNodeValue(WireEncoder encoder, IJsoncNode node)
    {
        var mmTag = node?.Tag ?? Tag.Empty();
        var payload = new WireEncoder();
        switch (node)
        {
            case JsoncValue jsoncValue:
                if (mmTag.IsNull)
                {
                    EncodeNullScalarPayload(payload, mmTag);
                }
                else
                {
                    EncodeJsoncScalarPayload(payload, jsoncValue);
                }
                break;
            case JsoncArray jsoncArray:
                EncodeJsoncArrayPayload(payload, jsoncArray);
                break;
            case JsoncObject jsoncObject:
                EncodeJsoncObjectPayload(payload, jsoncObject);
                break;
        }

        encoder.EncodeTaggedPayload(payload.ToByteArray(), mmTag.ToBytes());
    }

    private static void EncodeJsoncScalarPayload(WireEncoder encoder, JsoncValue value)
    {
        var valTag = value.Tag;
        if (valTag != null && valTag.Type == ValueType.Bigint)
        {
            encoder.EncodeBigIntDecimal(value.Value?.ToString() ?? "0");
            return;
        }
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
                    var strVal = value.Value?.ToString() ?? "0";
                    var numTag = value.Tag;
                    if (numTag != null && TypeInference.IsIntegerType(numTag.Type))
                    {
                        if (long.TryParse(strVal, out long lVal))
                            encoder.EncodeInt64(lVal);
                        else
                            encoder.EncodeFloatString(strVal);
                    }
                    else if (numTag != null && TypeInference.IsFloatType(numTag.Type))
                    {
                        encoder.EncodeFloatString(strVal);
                    }
                    else
                    {
                        if (long.TryParse(strVal, out long lVal))
                            encoder.EncodeInt64(lVal);
                        else
                            encoder.EncodeFloatString(strVal);
                    }
                }
                break;
            case JsoncTokenType.String:
            default:
                var tag = value.Tag;
                if (tag != null && tag.Type == ValueType.Enums && !string.IsNullOrEmpty(tag.Enums))
                {
                    var enumValues = tag.Enums.Split('|');
                    var strVal = value.Value?.ToString() ?? "";
                    int index = Array.IndexOf(enumValues, strVal);
                    if (index >= 0)
                    {
                        encoder.EncodeInt64(index);
                    }
                    else
                    {
                        encoder.EncodeString(strVal);
                    }
                }
                else if (tag != null && tag.Type == ValueType.Datetime)
                {
                    var dt = System.DateTime.Parse(value.Value?.ToString() ?? "", null, System.Globalization.DateTimeStyles.AssumeUniversal);
                    if (tag.Location != 0)
                    {
                        dt = dt.AddHours(-tag.Location);
                    }
                    encoder.EncodeInt64(TimeUtil.EpochSeconds(dt));
                }
                else if (tag != null && tag.Type == ValueType.Date)
                {
                    var dt = System.DateTime.Parse(value.Value?.ToString() ?? "", null, System.Globalization.DateTimeStyles.AssumeUniversal);
                    encoder.EncodeInt64(TimeUtil.DaysSinceEpochUtc(dt));
                }
                else if (tag != null && tag.Type == ValueType.Time)
                {
                    var rawVal = value.Value?.ToString() ?? "";
                    // Parse "HH:mm:ss" manually to avoid timezone issues with DateTime.Parse
                    int hour = 0, minute = 0, sec = 0;
                    var parts = rawVal.Split(':');
                    if (parts.Length >= 2)
                    {
                        int.TryParse(parts[0], out hour);
                        int.TryParse(parts[1], out minute);
                        if (parts.Length >= 3)
                            int.TryParse(parts[2], out sec);
                    }
                    long secs = hour * 3600L + minute * 60L + sec;
                    encoder.EncodeInt64(secs);
                }
                else if (tag != null && tag.Type == ValueType.Uuid)
                {
                    encoder.EncodeBytes(UuidToBytes(value.Value?.ToString() ?? ""));
                }
                else if (tag != null && (tag.Type == ValueType.Bytes || tag.Type == ValueType.Image || tag.Type == ValueType.Video || tag.Type == ValueType.Media))
                {
                    var base64Str = value.Value?.ToString() ?? "";
                    if (string.IsNullOrEmpty(base64Str))
                    {
                        encoder.EncodeBytes(Array.Empty<byte>());
                    }
                    else
                    {
                        encoder.EncodeBytes(Convert.FromBase64String(base64Str));
                    }
                }
                else
                {
                    encoder.EncodeString(value.Value?.ToString() ?? "");
                }
                break;
        }
    }

    private static void EncodeJsoncArrayPayload(WireEncoder encoder, JsoncArray array)
    {
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();
        var arrayTag = array.Tag;

        foreach (var element in array.Elements)
        {
            elementEncoder.Reset();
            if (arrayTag != null)
            {
                var inherited = Tag.Empty();
                inherited.InheritFromArrayParent(arrayTag);
                if (element.Tag != null)
                {
                    element.Tag = Tag.MergeTag(inherited, element.Tag);
                }
                else
                {
                    element.Tag = inherited;
                }
            }
            EncodeJsoncNodeValue(elementEncoder, element);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        encoder.EncodeArrayPayload(body.ToArray());
    }

    private static void EncodeJsoncObjectPayload(WireEncoder encoder, JsoncObject obj)
    {
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var kvp in obj.Fields)
        {
            tmp.Reset();
            tmp.EncodeString(kvp.Key);
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
                case ValueType.Bool:
                    result = new JsoncValue { Value = false, TokenType = JsoncTokenType.False };
                    break;
                case ValueType.I:
                case ValueType.I8:
                case ValueType.I16:
                case ValueType.I32:
                case ValueType.I64:
                case ValueType.U:
                case ValueType.U16:
                case ValueType.U32:
                case ValueType.U64:
                case ValueType.Bigint:
                case ValueType.F32:
                case ValueType.F64:
                case ValueType.Decimal:
                    result = new JsoncValue { Value = 0L, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.Str:
                case ValueType.Email:
                case ValueType.Url:
                case ValueType.Enums:
                case ValueType.Datetime:
                case ValueType.Date:
                case ValueType.Time:
                case ValueType.Uuid:
                    result = new JsoncValue { Value = "", TokenType = JsoncTokenType.String };
                    break;
                case ValueType.Bytes:
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
                case ValueType.Bool:
                    result = new JsoncValue
                    {
                        Value = scalar.Data,
                        TokenType = (bool)scalar.Data ? JsoncTokenType.True : JsoncTokenType.False
                    };
                    break;
                case ValueType.I:
                case ValueType.I8:
                case ValueType.I16:
                case ValueType.I32:
                case ValueType.I64:
                case ValueType.U:
                case ValueType.U16:
                case ValueType.U32:
                case ValueType.U64:
                    result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.F32:
                case ValueType.F64:
                case ValueType.Decimal:
                    result = new JsoncValue { Value = scalar.Data, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.Bigint:
                    result = new JsoncValue { Value = scalar.Text, TokenType = JsoncTokenType.Number };
                    break;
                case ValueType.Datetime:
                case ValueType.Date:
                case ValueType.Time:
                case ValueType.Uuid:
                case ValueType.Email:
                case ValueType.Url:
                case ValueType.Str:
                    result = new JsoncValue { Value = scalar.Text, TokenType = JsoncTokenType.String };
                    break;
                case ValueType.Bytes:
                    result = new JsoncValue
                    {
                        Value = scalar.Data is byte[] b ? Convert.ToBase64String(b) : scalar.Data?.ToString(),
                        TokenType = JsoncTokenType.String
                    };
                    break;
                case ValueType.Media:
                    result = new JsoncValue
                    {
                        Value = scalar.Text,
                        TokenType = JsoncTokenType.String
                    };
                    break;
                case ValueType.Enums:
                    result = new JsoncValue { Value = scalar.Text, TokenType = JsoncTokenType.String };
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
            result.Tag = scalar.Tag;
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
            jsoncArray.Tag = array.Tag;
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