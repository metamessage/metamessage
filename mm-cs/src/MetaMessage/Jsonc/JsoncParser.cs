using System.Net;
using MetaMessage.Core;
using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

namespace MetaMessage.Jsonc;

public class JsoncParser
{
    private readonly JsoncScanner _scanner;
    private readonly List<JsoncToken> _pendingComments = new();
    private int _depth;
    private const int MaxDepth = 32;

    public JsoncParser(string input)
    {
        _scanner = new JsoncScanner(input);
    }

    public INode Parse()
    {
        while (true)
        {
            var tok = _scanner.NextToken();
            if (tok.Type == JsoncTokenType.EOF)
                throw new Exception("no value parsed");

            if (tok.Type == JsoncTokenType.Comment)
            {
                if (_pendingComments.Count > 0)
                {
                    var last = _pendingComments.Last();
                    if (tok.Line - last.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(tok);
                continue;
            }

            Tag? tag;
            if ((tag = ConsumeCommentsFor(tok.Line)) == null)
                tag = Tag.NewTag();

            var val = ParseValue("", tok, false, tag);
            return val;
        }
    }

    private JsoncToken Advance()
    {
        return _scanner.NextToken();
    }

    private JsoncToken Peek()
    {
        var tok = _scanner.NextToken();
        _scanner.Unread();
        return tok;
    }

    private Tag? ConsumeCommentsFor(int anchorLine)
    {
        if (_pendingComments.Count == 0)
            return null;

        var last = _pendingComments.Last();
        if (anchorLine - last.Line > 1)
        {
            _pendingComments.Clear();
            return null;
        }

        Tag? outTag = null;
        foreach (var ct in _pendingComments)
        {
            var parsed = ParseCommentsToTag(ct.Literal);
            if (parsed != null)
            {
                outTag = Tag.MergeTag(outTag, parsed);
            }
        }

        // Note: intentionally not clearing _pendingComments here, matching Go behavior
        return outTag;
    }

    private static Tag? ParseCommentsToTag(string cs)
    {
        if (string.IsNullOrEmpty(cs))
            return null;
        // Go uses strings.CutPrefix(cs, "mm:") - exact prefix match
        if (cs.StartsWith("mm:"))
        {
            return Tag.Parse(cs);
        }
        return null;
    }

    private INode ParseValue(string path, JsoncToken firstTok, bool example, Tag? tag)
    {
        var tok = firstTok;
        while (true)
        {
            // If no tag provided, try to get from pending comments
            if (tag == null)
            {
                tag = ConsumeCommentsFor(tok.Line);
            }

            if (tag == null)
                tag = Tag.NewTag();

            switch (tok.Type)
            {
                case JsoncTokenType.EOF:
                    return new NodeScalar(null, "null", Tag.NewTag());

                case JsoncTokenType.Comment:
                    {
                        if (_pendingComments.Count > 0)
                        {
                            var last = _pendingComments.Last();
                            if (tok.Line - last.Line > 1)
                            {
                                _pendingComments.Clear();
                            }
                        }
                        _pendingComments.Add(tok);
                    }
                    tok = Advance();
                    continue;

                case JsoncTokenType.LBrace:
                    return ParseObject(tok.Line, path, tag);

                case JsoncTokenType.LBracket:
                    return ParseArray(tok.Line, path, tag);

                case JsoncTokenType.String:
                    return ParseString(tok, path, example, tag);

                case JsoncTokenType.Number:
                    return ParseNumber(tok, path, example, tag);

                case JsoncTokenType.True:
                case JsoncTokenType.False:
                    return ParseBool(tok, path, example, tag);

                case JsoncTokenType.Null:
                    if (tag.Type != ValueType.Unknown)
                        throw new Exception($"null is not supported for type {tag.Type}");
                    return new NodeNull(tag);

                default:
                    throw new Exception($"unexpected token {tok.Type} at line {tok.Line}");
            }
        }
    }

    private INode ParseString(JsoncToken tok, string path, bool example, Tag tag)
    {
        var text = tok.Literal;

        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Str;

        // Check simple types that force Str type
        switch (text)
        {
            case "code":
            case "message":
            case "data":
            case "success":
            case "error":
            case "unknown":
            case "page":
            case "limit":
            case "offset":
            case "total":
            case "id":
            case "name":
            case "description":
            case "type":
            case "version":
            case "status":
            case "url":
            case "create_time":
            case "update_time":
            case "delete_time":
            case "account":
            case "token":
            case "expire_time":
            case "key":
            case "value":
                tag.Type = ValueType.Str;
                break;
        }

        object? data = text;
        var ex = example || tag.Example;

        if (tag.IsNull)
        {
            switch (tag.Type)
            {
                case ValueType.Str:
                case ValueType.Decimal:
                case ValueType.Email:
                case ValueType.Url:
                    if (text != "")
                        throw new Exception($"invalid string: \"{text}\", valid: \"\"");
                    data = "";
                    break;

                case ValueType.Bytes:
                case ValueType.Media:
                    if (text != "")
                        throw new Exception($"invalid bytes: \"{text}\", valid: \"\"");
                    data = Array.Empty<byte>();
                    break;

                case ValueType.Datetime:
                    {
                        var defaultDt = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);
                        var defaultStr = defaultDt.ToString("yyyy-MM-dd HH:mm:ss");
                        if (text != defaultStr)
                            throw new Exception($"invalid datetime: \"{text}\", valid: \"{defaultStr}\"");
                        data = defaultDt;
                        break;
                    }

                case ValueType.Date:
                    {
                        var defaultDt = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);
                        var defaultStr = defaultDt.ToString("yyyy-MM-dd");
                        if (text != defaultStr)
                            throw new Exception($"invalid date: \"{text}\", valid: \"{defaultStr}\"");
                        data = defaultDt;
                        break;
                    }

                case ValueType.Time:
                    {
                        var defaultDt = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);
                        var defaultStr = defaultDt.ToString("HH:mm:ss");
                        if (text != defaultStr)
                            throw new Exception($"invalid time: \"{text}\", valid: \"{defaultStr}\"");
                        data = defaultDt;
                        break;
                    }

                case ValueType.Uuid:
                    if (text != "")
                        throw new Exception($"invalid uuid: \"{text}\", valid: \"\"");
                    data = Guid.Empty;
                    break;

                case ValueType.Ip:
                    if (text != "")
                        throw new Exception($"invalid ip: \"{text}\", valid: \"\"");
                    data = IPAddress.Any;
                    break;

                case ValueType.Enums:
                    if (text != "")
                        throw new Exception($"invalid enums: \"{text}\", valid: \"\"");
                    data = -1;
                    break;

                default:
                    throw new Exception($"unsupported type {tag.Type} for null string");
            }

            return new NodeScalar(data, text, tag) { Path = path };
        }

        // Type dispatch for non-null strings
        switch (tag.Type)
        {
            case ValueType.Str:
                {
                    if (!ex)
                    {
                        var result = Validator.Validate(text, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = text;
                    break;
                }

            case ValueType.Bytes:
                {
                    byte[] bytes;
                    try
                    {
                        bytes = Convert.FromBase64String(text);
                    }
                    catch (Exception exBytes)
                    {
                        throw new Exception($"invalid base64 bytes \"{text}\": {exBytes.Message}");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(bytes, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = bytes;
                    break;
                }

            case ValueType.Datetime:
                {
                    DateTime dt;
                    if (!DateTime.TryParseExact(text, "yyyy-MM-dd HH:mm:ss",
                            System.Globalization.CultureInfo.InvariantCulture,
                            System.Globalization.DateTimeStyles.None, out dt))
                    {
                        throw new Exception($"invalid datetime \"{text}\"");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(dt, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = dt;
                    break;
                }

            case ValueType.Date:
                {
                    DateTime dt;
                    if (!DateTime.TryParseExact(text, "yyyy-MM-dd",
                            System.Globalization.CultureInfo.InvariantCulture,
                            System.Globalization.DateTimeStyles.None, out dt))
                    {
                        throw new Exception($"invalid date \"{text}\"");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(dt, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = dt;
                    break;
                }

            case ValueType.Time:
                {
                    DateTime dt;
                    if (!DateTime.TryParseExact(text, "HH:mm:ss",
                            System.Globalization.CultureInfo.InvariantCulture,
                            System.Globalization.DateTimeStyles.None, out dt))
                    {
                        throw new Exception($"invalid time \"{text}\"");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(dt, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = dt;
                    break;
                }

            case ValueType.Uuid:
                {
                    if (!ex)
                    {
                        var result = Validator.Validate(text, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = text;
                    break;
                }

            case ValueType.Ip:
                {
                    IPAddress ip;
                    if (!IPAddress.TryParse(text, out ip!))
                    {
                        throw new Exception($"invalid ip \"{text}\"");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(ip, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = ip;
                    break;
                }

            case ValueType.Url:
                {
                    Uri uri;
                    if (!Uri.TryCreate(text, UriKind.Absolute, out uri!))
                    {
                        throw new Exception($"invalid url \"{text}\"");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(uri, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = uri;
                    break;
                }

            case ValueType.Email:
                {
                    if (!ex)
                    {
                        var result = Validator.Validate(text, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = text;
                    break;
                }

            case ValueType.Enums:
                {
                    if (string.IsNullOrEmpty(tag.Enums))
                        throw new Exception("enum empty");

                    if (!ex)
                    {
                        var result = Validator.Validate(text, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = text;
                    break;
                }

            case ValueType.Media:
                {
                    byte[] bytes;
                    try
                    {
                        bytes = Convert.FromBase64String(text);
                    }
                    catch (Exception exMedia)
                    {
                        throw new Exception($"invalid base64 media \"{text}\": {exMedia.Message}");
                    }
                    if (!ex)
                    {
                        var result = Validator.Validate(bytes, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = bytes;
                    break;
                }

            case ValueType.Decimal:
                {
                    if (!ex)
                    {
                        var result = Validator.Validate(text, tag);
                        if (!result.IsValid)
                            throw new Exception(string.Join("; ", result.Errors));
                    }
                    data = text;
                    break;
                }

            default:
                throw new Exception($"unsupported type {tag.Type} for string literal");
        }

        return new NodeScalar(data, text, tag) { Path = path };
    }

    private INode ParseNumber(JsoncToken tok, string path, bool example, Tag tag)
    {
        var text = tok.Literal;
        object? data = text;
        var ex = example || tag.Example;

        if (text.Contains("."))
        {
            // Float literal
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.F64;

            if (tag.IsNull)
            {
                if (text != "0.0")
                    throw new Exception($"invalid float: {text}, valid: 0.0");
                switch (tag.Type)
                {
                    case ValueType.F32:
                        data = 0.0f;
                        break;
                    case ValueType.F64:
                        data = 0.0;
                        break;
                    case ValueType.Decimal:
                        data = "";
                        break;
                    default:
                        throw new Exception($"unsupported numeric type {tag.Type} for float literal");
                }
                return new NodeScalar(data, text, tag) { Path = path };
            }

            switch (tag.Type)
            {
                case ValueType.F32:
                    {
                        float f32;
                        if (!float.TryParse(text, System.Globalization.NumberStyles.Float,
                                System.Globalization.CultureInfo.InvariantCulture, out f32))
                        {
                            throw new Exception($"invalid float32 \"{text}\"");
                        }
                        if (!ex)
                        {
                            var result = Validator.Validate((double)f32, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = f32;
                        break;
                    }

                case ValueType.F64:
                    {
                        double f64;
                        if (!double.TryParse(text, System.Globalization.NumberStyles.Float,
                                System.Globalization.CultureInfo.InvariantCulture, out f64))
                        {
                            throw new Exception($"invalid float64 \"{text}\"");
                        }
                        if (!ex)
                        {
                            var result = Validator.Validate(f64, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = f64;
                        break;
                    }

                case ValueType.Decimal:
                    {
                        if (!ex)
                        {
                            var result = Validator.Validate(text, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = text;
                        break;
                    }

                default:
                    throw new Exception($"unsupported numeric type {tag.Type} for float literal");
            }
        }
        else if (text.StartsWith("-"))
        {
            // Negative integer literal
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.I;

            if (tag.IsNull)
            {
                if (text != "0")
                    throw new Exception($"invalid int: {text}, valid: 0");
                switch (tag.Type)
                {
                    case ValueType.I: data = 0; break;
                    case ValueType.I8: data = (sbyte)0; break;
                    case ValueType.I16: data = (short)0; break;
                    case ValueType.I32: data = 0; break;
                    case ValueType.I64: data = 0L; break;
                    case ValueType.Bigint: data = System.Numerics.BigInteger.Zero; break;
                    default:
                        throw new Exception($"unsupported numeric type {tag.Type} for negative literal");
                }
                return new NodeScalar(data, text, tag) { Path = path };
            }

            switch (tag.Type)
            {
                case ValueType.I:
                    {
                        long val;
                        if (!long.TryParse(text, out val))
                            throw new Exception($"invalid int \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate((int)val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = (int)val;
                        break;
                    }

                case ValueType.I8:
                    {
                        sbyte val;
                        if (!sbyte.TryParse(text, out val))
                            throw new Exception($"invalid int8 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I16:
                    {
                        short val;
                        if (!short.TryParse(text, out val))
                            throw new Exception($"invalid int16 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I32:
                    {
                        int val;
                        if (!int.TryParse(text, out val))
                            throw new Exception($"invalid int32 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I64:
                    {
                        long val;
                        if (!long.TryParse(text, out val))
                            throw new Exception($"invalid int64 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.Bigint:
                    {
                        System.Numerics.BigInteger val;
                        if (!System.Numerics.BigInteger.TryParse(text, out val))
                            throw new Exception($"invalid bigint \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                default:
                    throw new Exception($"unsupported numeric type {tag.Type} for negative literal");
            }
        }
        else
        {
            // Non-negative integer literal
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.I;

            if (tag.IsNull)
            {
                if (text != "0")
                    throw new Exception($"invalid int: {text}, valid: 0");
                switch (tag.Type)
                {
                    case ValueType.I: data = 0; break;
                    case ValueType.I8: data = (sbyte)0; break;
                    case ValueType.I16: data = (short)0; break;
                    case ValueType.I32: data = 0; break;
                    case ValueType.I64: data = 0L; break;
                    case ValueType.U: data = 0u; break;
                    case ValueType.U8: data = (byte)0; break;
                    case ValueType.U16: data = (ushort)0; break;
                    case ValueType.U32: data = 0u; break;
                    case ValueType.U64: data = 0UL; break;
                    case ValueType.Bigint: data = System.Numerics.BigInteger.Zero; break;
                    default:
                        throw new Exception($"unsupported numeric type {tag.Type}");
                }
                return new NodeScalar(data, text, tag) { Path = path };
            }

            switch (tag.Type)
            {
                case ValueType.I:
                    {
                        long val;
                        if (!long.TryParse(text, out val))
                            throw new Exception($"invalid int \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate((int)val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = (int)val;
                        break;
                    }

                case ValueType.I8:
                    {
                        sbyte val;
                        if (!sbyte.TryParse(text, out val))
                            throw new Exception($"invalid int8 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I16:
                    {
                        short val;
                        if (!short.TryParse(text, out val))
                            throw new Exception($"invalid int16 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I32:
                    {
                        int val;
                        if (!int.TryParse(text, out val))
                            throw new Exception($"invalid int32 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.I64:
                    {
                        long val;
                        if (!long.TryParse(text, out val))
                            throw new Exception($"invalid int64 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.U:
                    {
                        ulong val;
                        if (!ulong.TryParse(text, out val))
                            throw new Exception($"invalid uint \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate((uint)val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = (uint)val;
                        break;
                    }

                case ValueType.U8:
                    {
                        byte val;
                        if (!byte.TryParse(text, out val))
                            throw new Exception($"invalid uint8 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.U16:
                    {
                        ushort val;
                        if (!ushort.TryParse(text, out val))
                            throw new Exception($"invalid uint16 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.U32:
                    {
                        uint val;
                        if (!uint.TryParse(text, out val))
                            throw new Exception($"invalid uint32 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.U64:
                    {
                        ulong val;
                        if (!ulong.TryParse(text, out val))
                            throw new Exception($"invalid uint64 \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                case ValueType.Bigint:
                    {
                        System.Numerics.BigInteger val;
                        if (!System.Numerics.BigInteger.TryParse(text, out val))
                            throw new Exception($"invalid bigint \"{text}\"");
                        if (!ex)
                        {
                            var result = Validator.Validate(val, tag);
                            if (!result.IsValid)
                                throw new Exception(string.Join("; ", result.Errors));
                        }
                        data = val;
                        break;
                    }

                default:
                    throw new Exception($"unsupported numeric type {tag.Type}");
            }
        }

        return new NodeScalar(data, text, tag) { Path = path };
    }

    private INode ParseBool(JsoncToken tok, string path, bool example, Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Bool;

        var boolVal = tok.Type == JsoncTokenType.True;

        if (tag.IsNull)
        {
            if (boolVal)
                throw new Exception("bool must false when bool is null");
        }
        else
        {
            var ex = example || tag.Example;
            if (!ex)
            {
                var result = Validator.Validate(boolVal, tag);
                if (!result.IsValid)
                    throw new Exception(string.Join("; ", result.Errors));
            }
        }

        return new NodeScalar(boolVal, boolVal ? "true" : "false", tag) { Path = path };
    }

    private INode ParseObject(int openLine, string path, Tag? existingTag = null)
    {
        _depth++;
        if (_depth > MaxDepth)
            throw new Exception($"max depth: {MaxDepth}");

        Tag tag = ConsumeCommentsFor(openLine) ?? Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Obj;

        if (!string.IsNullOrEmpty(tag.Name))
        {
            path = string.IsNullOrEmpty(path) ? tag.Name : $"{path}.{tag.Name}";
        }

        var entries = new List<KeyValuePair<NodeScalar, INode>>();

        while (true)
        {
            var tok = Advance();
            if (tok.Type == JsoncTokenType.EOF)
                break;
            if (tok.Type == JsoncTokenType.RBrace)
                break;

            if (tok.Type == JsoncTokenType.Comment)
            {
                if (_pendingComments.Count > 0)
                {
                    var last = _pendingComments.Last();
                    if (tok.Line - last.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(tok);
                continue;
            }

            if (tok.Type != JsoncTokenType.String)
                throw new Exception($"expect string key at line {tok.Line}");

            var keyStr = CamelToSnake.Convert(tok.Literal);

            var colonTok = Advance();
            if (colonTok.Type != JsoncTokenType.Colon)
                throw new Exception("expect colon");

            // Get child tag from comments (matching Go's consumeCommentsFor before parse)
            var childTag = ConsumeCommentsFor(tok.Line);
            if (childTag == null)
                childTag = Tag.NewTag();

            // Inherit from parent if Map type (matching Go)
            if (tag.Type == ValueType.Map)
            {
                childTag.Inherit(tag);
                if (childTag.Example)
                {
                    tag.IsEmpty = true;
                }
            }

            var fieldPath = tag!.Type == ValueType.Map
                ? $"{path}[{keyStr}]"
                : $"{path}.{keyStr}";

            var exampleMode = tag.Example;
            var val = ParseValue(fieldPath, Advance(), exampleMode, childTag);
            if (val == null)
                continue;

            var keyScalar = new NodeScalar(keyStr, keyStr, Tag.Empty()) { Path = fieldPath };
            entries.Add(new KeyValuePair<NodeScalar, INode>(keyScalar, val));

            var nextTok = Peek();
            if (nextTok.Type == JsoncTokenType.Comma)
            {
                Advance();
            }
        }

        var map = new NodeObject(entries, tag!) { Path = path };

        if (!tag!.Example)
        {
            switch (tag!.Type)
            {
                case ValueType.Map:
                    Validator.Validate(map, tag);
                    break;
                case ValueType.Obj:
                    Validator.Validate(map, tag);
                    break;
            }
        }

        _depth--;
        return map;
    }

    private INode ParseArray(int openLine, string path, Tag? existingTag = null)
    {
        _depth++;
        if (_depth > MaxDepth)
            throw new Exception($"max depth: {MaxDepth}");

        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Vec;
        }

        if (!string.IsNullOrEmpty(tag.Name))
        {
            path = $"{path}.{tag.Name}";
        }

        var children = new List<INode>();

        int i = 0;

        while (true)
        {
            var tok = Advance();
            if (tok.Type == JsoncTokenType.EOF)
                break;
            if (tok.Type == JsoncTokenType.RBracket)
                break;

            if (tok.Type == JsoncTokenType.Comment)
            {
                if (_pendingComments.Count > 0)
                {
                    var last = _pendingComments.Last();
                    if (tok.Line - last.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(tok);
                continue;
            }

            // Get child tag from comments, only if not on same line as open bracket (matching Go)
            Tag? childTag = null;
            if (openLine != tok.Line)
            {
                childTag = ConsumeCommentsFor(tok.Line);
            }
            if (childTag == null)
                childTag = Tag.NewTag();

            // Inherit from parent (matching Go: always inherit for array items)
            childTag.Inherit(tag);
            if (childTag.Example)
            {
                tag.IsEmpty = true;
            }

            var itemPath = $"{path}[{i}]";
            var exampleMode = tag.Example;
            var item = ParseValue(itemPath, tok, exampleMode, childTag);
            if (item == null)
                continue;

            children.Add(item);
            i++;

            var nextTok = Peek();
            if (nextTok.Type == JsoncTokenType.Comma)
            {
                Advance();
            }
        }

        var arr = new NodeArray(children, tag!) { Path = path };

        if (!tag!.Example)
        {
            switch (tag!.Type)
            {
                case ValueType.Arr:
                    Validator.Validate(arr, tag);
                    break;
                case ValueType.Vec:
                    Validator.Validate(arr, tag);
                    break;
            }
        }

        _depth--;
        return arr;
    }
}