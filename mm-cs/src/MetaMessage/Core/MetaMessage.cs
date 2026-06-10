using ValueType = MetaMessage.Ir.ValueType;
using MetaMessage.Jsonc;
using MetaMessage.Ir;
using JsoncParser = MetaMessage.Jsonc.Jsonc;

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

    public static INode Decode(byte[] data)
    {
        var decoder = new WireDecoder(data);
        return decoder.Decode();
    }

    public static INode DecodeToTree(byte[] data)
    {
        return Decode(data);
    }

    public static INode ParseFromJSONC(string input)
    {
        return JsoncParser.ParseFromString(input);
    }

    public static byte[] EncodeTree(INode tree)
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

    public static byte[] EncodeFromJsonc(string jsonc)
    {
        var tree = JsoncParser.ParseFromString(jsonc);
        return EncodeFromJsoncTree(tree);
    }

    public static object? JsoncToValue(string jsonc)
    {
        var tree = JsoncParser.ParseFromString(jsonc);
        return ExtractValueFromTree(tree);
    }

    public static string DecodeToJsonc(byte[] data)
    {
        var tree = Decode(data);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(tree);
    }

    public static string ValueToJsonc(object value)
    {
        var binder = new Binder();
        var tree = binder.StructToNode(value);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(tree);
    }

    public static object? DecodeToValue(byte[] data)
    {
        var tree = Decode(data);
        return ExtractValueFromTree(tree);
    }

    public static string ValueToJsonc(object value, string tagString)
    {
        var node = ReflectMmEncoder.ValueToNode(value, tagString);
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(node);
    }

    public static ValidationResult Validate(object value, Tag tag)
    {
        return Validator.Validate(value, tag);
    }

    private static string TreeToJsoncString(INode tree)
    {
        var printer = new JsoncPrinter(prettyPrint: true);
        return printer.Print(tree);
    }

    private static void EncodeTreeValue(WireEncoder encoder, INode tree, Tag inherited)
    {
        switch (tree)
        {
            case NodeScalar scalar:
                EncodeScalarTree(encoder, scalar, inherited);
                break;
            case NodeArray array:
                EncodeArrayTree(encoder, array, inherited);
                break;
            case NodeObject map:
                EncodeMapTree(encoder, map, inherited);
                break;
        }
    }

    private static void EncodeScalarTree(WireEncoder encoder, NodeScalar scalar, Tag inherited)
    {
        var tag = scalar.Tag?.Copy() ?? inherited.Copy();
        if (scalar.Data == null || tag.IsNull)
        {
            tag.IsNull = true;
            var nullPayload = new WireEncoder();
            EncodeNullScalarPayload(nullPayload, tag);
            encoder.EncodeTaggedPayload(nullPayload.ToByteArray(), EncodeTagBytesWithPrefix(tag.ToBytes()));
            return;
        }

        var payload = new WireEncoder();
        EncodeScalarPayload(payload, scalar.Data, tag);
        encoder.EncodeTaggedPayload(payload.ToByteArray(), EncodeTagBytesWithPrefix(tag.ToBytes()));
    }

    private static void EncodeNullScalarPayload(WireEncoder encoder, Tag tag)
    {
        switch (tag.Type)
        {
            case ValueType.Unknown:
                encoder.EncodeSimple(SimpleValue.NULL);
                return;
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
            case ValueType.Bigint:
                encoder.EncodeSimple(SimpleValue.NULL_BYTES);
                return;
            case ValueType.F32:
            case ValueType.F64:
                encoder.EncodeSimple(SimpleValue.NULL_FLOAT);
                return;
            case ValueType.Str:
            case ValueType.Email:
            case ValueType.Url:
            case ValueType.Enums:
            case ValueType.Uuid:
                encoder.EncodeSimple(SimpleValue.NULL_STRING);
                return;
            case ValueType.Datetime:
            case ValueType.Date:
            case ValueType.Time:
                encoder.EncodeSimple(SimpleValue.NULL_INT);
                return;
            case ValueType.Decimal:
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

    private static void EncodeScalarPayload(WireEncoder encoder, object data, Tag tag)
    {
        switch (tag.Type)
        {
            case ValueType.Bool:
                encoder.EncodeBool((bool)data);
                return;
            case ValueType.I:
                encoder.EncodeInt64(Convert.ToInt64(data));
                return;
            case ValueType.I8:
                encoder.EncodeInt8(Convert.ToSByte(data));
                return;
            case ValueType.I16:
                encoder.EncodeInt16(Convert.ToInt16(data));
                return;
            case ValueType.I32:
                encoder.EncodeInt32(Convert.ToInt32(data));
                return;
            case ValueType.I64:
                encoder.EncodeInt64(Convert.ToInt64(data));
                return;
            case ValueType.U:
                encoder.EncodeUint32(Convert.ToUInt32(data));
                return;
            case ValueType.U8:
                encoder.EncodeUint8(Convert.ToByte(data));
                return;
            case ValueType.U16:
                encoder.EncodeUint16(Convert.ToUInt16(data));
                return;
            case ValueType.U32:
                encoder.EncodeUint32(Convert.ToUInt32(data));
                return;
            case ValueType.U64:
                encoder.EncodeUint64(Convert.ToUInt64(data));
                return;
            case ValueType.F32:
                encoder.EncodeFloat(Convert.ToSingle(data));
                return;
            case ValueType.F64:
                encoder.EncodeDouble(Convert.ToDouble(data));
                return;
            case ValueType.Str:
            case ValueType.Email:
            case ValueType.Url:
            case ValueType.Enums:
            case ValueType.Uuid:
                encoder.EncodeString(data.ToString() ?? "");
                return;
            case ValueType.Datetime:
                if (data is DateTime dt)
                {
                    encoder.EncodeDatetime(dt);
                }
                else if (data is long l)
                {
                    encoder.EncodeDatetime(DateTimeOffset.FromUnixTimeSeconds(l).DateTime);
                }
                else
                {
                    encoder.EncodeString(data.ToString() ?? "");
                }
                return;
            case ValueType.Date:
            case ValueType.Time:
                encoder.EncodeString(data.ToString() ?? "");
                return;
            case ValueType.Bytes:
                if (data is byte[] bytes)
                {
                    encoder.EncodeBytes(bytes);
                }
                else if (data is string s)
                {
                    encoder.EncodeBytes(Convert.FromBase64String(s));
                }
                else
                {
                    encoder.EncodeBytes(Array.Empty<byte>());
                }
                return;
            case ValueType.Decimal:
                encoder.EncodeString(data.ToString() ?? "0");
                return;
            case ValueType.Bigint:
                encoder.EncodeBigIntDecimal(data.ToString() ?? "0");
                return;
            case ValueType.Arr:
            case ValueType.Vec:
            case ValueType.Obj:
            case ValueType.Map:
                // Should not reach here for scalars
                encoder.EncodeString(data.ToString() ?? "");
                return;
            default:
                encoder.EncodeString(data.ToString() ?? "");
                return;
        }
    }

    private static void EncodeArrayTree(WireEncoder encoder, NodeArray array, Tag inherited)
    {
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();
        var arrayTag = array.Tag ?? inherited;

        foreach (var element in array.Children)
        {
            elementEncoder.Reset();
            var inheritedTag = Tag.Empty();
            inheritedTag.InheritFromArrayParent(arrayTag);
            if (element.Tag != null)
            {
                element.Tag = Tag.MergeTag(inheritedTag, element.Tag);
            }
            else
            {
                element.Tag = inheritedTag;
            }
            EncodeTreeValue(elementEncoder, element, inheritedTag);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        encoder.EncodeArrayPayload(body.ToArray());
    }

    private static void EncodeMapTree(WireEncoder encoder, NodeObject map, Tag inherited)
    {
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var kvp in map.Entries)
        {
            tmp.Reset();
            tmp.EncodeString(kvp.Key.Text);
            keysPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            EncodeTreeValue(tmp, kvp.Value, inherited);
            valsPacked.WriteAll(tmp.ToByteArray());
        }

        tmp.Reset();
        tmp.EncodeArrayPayload(keysPacked.ToArray());
        var mapBody = new GrowableByteBuf();
        mapBody.WriteAll(tmp.ToByteArray());
        mapBody.WriteAll(valsPacked.ToArray());

        encoder.EncodeObjectPayload(mapBody.ToArray());
    }

    private static byte[] EncodeFromJsoncTree(INode tree)
    {
        var encoder = new WireEncoder();
        EncodeJsoncTreeValue(encoder, tree);
        return encoder.ToByteArray();
    }

    private static byte[] EncodeTagBytesWithPrefix(byte[] rawTagBytes)
    {
        if (rawTagBytes.Length == 0)
            return Array.Empty<byte>();

        var enc = new WireEncoder();
        enc.EncodeTagInner(rawTagBytes);
        return enc.ToByteArray();
    }

    private static void EncodeJsoncTreeValue(WireEncoder encoder, INode tree)
    {
        var mmTag = tree?.Tag ?? Tag.Empty();
        var payload = new WireEncoder();
        switch (tree)
        {
            case NodeScalar scalar:
                if (mmTag.IsNull)
                {
                    EncodeNullScalarPayload(payload, mmTag);
                }
                else
                {
                    EncodeJsoncScalarPayload(payload, scalar);
                }
                break;
            case NodeArray array:
                EncodeJsoncArrayPayload(payload, array);
                break;
            case NodeObject map:
                EncodeJsoncMapPayload(payload, map);
                break;
            case NodeNull:
                payload.EncodeSimple(SimpleValue.NULL);
                break;
        }

        encoder.EncodeTaggedPayload(payload.ToByteArray(), EncodeTagBytesWithPrefix(mmTag.ToBytes()));
    }

    private static void EncodeJsoncScalarPayload(WireEncoder encoder, NodeScalar scalar)
    {
        var valTag = scalar.Tag;
        if (valTag != null && valTag.Type == ValueType.Bigint)
        {
            encoder.EncodeBigIntDecimal(scalar.Text);
            return;
        }

        if (scalar.Data == null)
        {
            encoder.EncodeSimple(SimpleValue.NULL_STRING);
            return;
        }

        // Check if the tag indicates a special type that needs type-specific encoding.
        if (valTag != null && valTag.Type != ValueType.Unknown && valTag.Type != ValueType.Str)
        {
            switch (valTag.Type)
            {
                case ValueType.Datetime:
                    if (DateTime.TryParse(scalar.Text, out var dt))
                    {
                        var loc = valTag?.Location ?? 0;
                        long unixTime;
                        if (loc == 0)
                        {
                            // UTC
                            dt = DateTime.SpecifyKind(dt, DateTimeKind.Utc);
                            unixTime = new DateTimeOffset(dt).ToUnixTimeSeconds();
                        }
                        else
                        {
                            // Apply timezone offset
                            var offset = TimeSpan.FromHours(loc);
                            var dto = new DateTimeOffset(dt, offset);
                            unixTime = dto.ToUnixTimeSeconds();
                        }
                        encoder.EncodeInt64(unixTime);
                    }
                    else
                    {
                        encoder.EncodeString(scalar.Text);
                    }
                    return;

                case ValueType.Date:
                    if (DateTime.TryParse(scalar.Text, out var dateVal))
                    {
                        var loc = valTag?.Location ?? 0;
                        DateTime utcDate;
                        if (loc == 0)
                        {
                            utcDate = DateTime.SpecifyKind(dateVal, DateTimeKind.Utc);
                        }
                        else
                        {
                            var offset = TimeSpan.FromHours(loc);
                            var dto = new DateTimeOffset(dateVal, offset);
                            utcDate = dto.UtcDateTime;
                        }
                        var epoch = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);
                        var days = (int)(utcDate - epoch).TotalDays;
                        encoder.EncodeInt64(days);
                    }
                    else
                    {
                        encoder.EncodeString(scalar.Text);
                    }
                    return;

                case ValueType.Time:
                    if (TimeSpan.TryParse(scalar.Text, out var ts))
                    {
                        encoder.EncodeInt64((long)ts.TotalSeconds);
                    }
                    else
                    {
                        encoder.EncodeString(scalar.Text);
                    }
                    return;

                case ValueType.Uuid:
                    {
                        var uuidStr = scalar.Text?.Replace("-", "") ?? "";
                        if (uuidStr.Length == 32 && System.Text.RegularExpressions.Regex.IsMatch(uuidStr, "^[0-9a-fA-F]{32}$"))
                        {
                            var uuidBytes = new byte[16];
                            for (int i = 0; i < 16; i++)
                            {
                                uuidBytes[i] = Convert.ToByte(uuidStr.Substring(i * 2, 2), 16);
                            }
                            encoder.EncodeBytes(uuidBytes);
                        }
                        else
                        {
                            encoder.EncodeString(scalar.Text ?? "");
                        }
                        return;
                    }

                case ValueType.Ip:
                    {
                        var ipStr = scalar.Text ?? "";
                        if (valTag.Version != 0 && System.Net.IPAddress.TryParse(ipStr, out var ipAddr))
                        {
                            var ipBytes = ipAddr.GetAddressBytes();
                            encoder.EncodeBytes(ipBytes);
                        }
                        else
                        {
                            encoder.EncodeString(ipStr);
                        }
                        return;
                    }

                case ValueType.Enums:
                    {
                        var enumStr = scalar.Text ?? "";
                        if (!string.IsNullOrEmpty(valTag.Enums))
                        {
                            var parts = valTag.Enums.Split('|');
                            int idx = Array.IndexOf(parts, enumStr);
                            if (idx >= 0)
                            {
                                encoder.EncodeInt64(idx);
                            }
                            else
                            {
                                encoder.EncodeString(enumStr);
                            }
                        }
                        else
                        {
                            encoder.EncodeString(enumStr);
                        }
                        return;
                    }

                case ValueType.Media:
                case ValueType.Bytes:
                    {
                        var b64 = scalar.Text ?? "";
                        try
                        {
                            var rawBytes = Convert.FromBase64String(b64);
                            encoder.EncodeBytes(rawBytes);
                        }
                        catch
                        {
                            encoder.EncodeString(b64);
                        }
                        return;
                    }

                case ValueType.Bigint:
                    encoder.EncodeBigIntDecimal(scalar.Text);
                    return;

                case ValueType.Decimal:
                    encoder.EncodeFloatString(scalar.Text);
                    return;

                case ValueType.Email:
                case ValueType.Url:
                    encoder.EncodeString(scalar.Text ?? "");
                    return;
            }
        }

        switch (scalar.Data)
        {
            case double d:
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
                break;
            case bool b:
                encoder.EncodeSimple(b ? SimpleValue.TRUE : SimpleValue.FALSE);
                break;
            case string strVal:
                {
                    var numTag = scalar.Tag;
                    if (numTag != null && TypeInference.IsIntegerType(numTag.Type))
                    {
                        if (long.TryParse(strVal, out long lVal))
                            encoder.EncodeInt64(lVal);
                        else
                            encoder.EncodeString(strVal);
                    }
                    else if (numTag != null && TypeInference.IsFloatType(numTag.Type))
                    {
                        encoder.EncodeFloatString(strVal);
                    }
                    else
                    {
                        encoder.EncodeString(strVal);
                    }
                }
                break;
            default:
                {
                    var strVal = scalar.Text;
                    var numTag = scalar.Tag;
                    if (numTag != null && TypeInference.IsIntegerType(numTag.Type))
                    {
                        if (long.TryParse(strVal, out long lVal))
                            encoder.EncodeInt64(lVal);
                        else
                            encoder.EncodeString(strVal);
                    }
                    else if (numTag != null && TypeInference.IsFloatType(numTag.Type))
                    {
                        encoder.EncodeFloatString(strVal);
                    }
                    else
                    {
                        encoder.EncodeString(strVal ?? "");
                    }
                }
                break;
        }
    }

    private static void EncodeJsoncArrayPayload(WireEncoder encoder, NodeArray array)
    {
        var body = new GrowableByteBuf();
        var elementEncoder = new WireEncoder();
        var arrayTag = array.Tag;

        foreach (var element in array.Children)
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
            EncodeJsoncTreeValue(elementEncoder, element);
            body.WriteAll(elementEncoder.ToByteArray());
        }

        encoder.EncodeArrayPayload(body.ToArray());
    }

    private static void EncodeJsoncMapPayload(WireEncoder encoder, NodeObject map)
    {
        var keysPacked = new GrowableByteBuf();
        var valsPacked = new GrowableByteBuf();
        var tmp = new WireEncoder();

        foreach (var kvp in map.Entries)
        {
            tmp.Reset();
            tmp.EncodeString(kvp.Key.Text);
            keysPacked.WriteAll(tmp.ToByteArray());

            tmp.Reset();
            EncodeJsoncTreeValue(tmp, kvp.Value);
            valsPacked.WriteAll(tmp.ToByteArray());
        }

        tmp.Reset();
        tmp.EncodeArrayPayload(keysPacked.ToArray());
        var mapBody = new GrowableByteBuf();
        mapBody.WriteAll(tmp.ToByteArray());
        mapBody.WriteAll(valsPacked.ToArray());

        encoder.EncodeObjectPayload(mapBody.ToArray());
    }

    private static object? ExtractValueFromTree(INode tree)
    {
        switch (tree)
        {
            case NodeScalar scalar:
                return scalar.Data;
            case NodeArray array:
                {
                    var list = new List<object?>();
                    foreach (var element in array.Children)
                    {
                        list.Add(ExtractValueFromTree(element));
                    }
                    return list;
                }
            case NodeObject map:
                {
                    var dict = new Dictionary<string, object?>();
                    foreach (var kvp in map.Entries)
                    {
                        dict[kvp.Key.Text] = ExtractValueFromTree(kvp.Value);
                    }
                    return dict;
                }
            default:
                return null;
        }
    }
}