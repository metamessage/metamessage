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
                return new NodeScalar(null, "null", Tag.NewTag());

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


            var val = ParseValue("", tok, false);
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

        _pendingComments.Clear();
        return outTag;
    }

    private static Tag? ParseCommentsToTag(string cs)
    {
        if (string.IsNullOrWhiteSpace(cs))
            return null;
        var trimmed = cs.TrimStart();
        if (trimmed.StartsWith("mm:"))
        {
            return Tag.Parse(cs);
        }
        return null;
    }

    private INode ParseValue(string path, JsoncToken firstTok, bool example)
    {
        var tok = firstTok;
        while (true)
        {
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
                    return ParseObject(tok.Line, path);

                case JsoncTokenType.LBracket:
                    return ParseArray(tok.Line, path);

                case JsoncTokenType.String:
                    return ParseString(tok, path, example);

                case JsoncTokenType.Number:
                    return ParseNumber(tok, path, example);

                case JsoncTokenType.True:
                case JsoncTokenType.False:
                    return ParseBool(tok, path, example);

                case JsoncTokenType.Null:
                    throw new Exception("null is not supported");

                default:
                    throw new Exception($"unexpected token {tok.Type} at line {tok.Line}");
            }
        }
    }

    private INode ParseString(JsoncToken tok, string path, bool example)
    {
        var tag = ConsumeCommentsFor(tok.Line);
        var text = tok.Literal;

        if (tag == null)
            tag = Tag.NewTag();

        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Str;

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

            default:
                // type-specific parsing handled by tag type
                break;
        }

        return new NodeScalar(text, text, tag) { Path = path };
    }

    private INode ParseNumber(JsoncToken tok, string path, bool example)
    {
        var tag = ConsumeCommentsFor(tok.Line);
        var text = tok.Literal;

        if (tag == null)
            tag = Tag.NewTag();

        object data;
        if (text.Contains("."))
        {
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.F64;
            data = double.Parse(text);
        }
        else
        {
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.I;
            if (long.TryParse(text, out long longVal))
            {
                data = (double)longVal;
            }
            else
            {
                data = text;
            }
        }

        return new NodeScalar(data, text, tag) { Path = path };
    }

    private INode ParseBool(JsoncToken tok, string path, bool example)
    {
        var tag = ConsumeCommentsFor(tok.Line);

        if (tag == null)
            tag = Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Bool;

        bool boolVal = tok.Type == JsoncTokenType.True;

        return new NodeScalar(boolVal, boolVal ? "true" : "false", tag) { Path = path };
    }

    private INode ParseObject(int openLine, string path)
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

        INode? lastValue = null;

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

            var fieldPath = tag!.Type == ValueType.Map
                ? $"{path}[{keyStr}]"
                : $"{path}.{keyStr}";

            var exampleMode = tag!.Example;
            var val = ParseValue(fieldPath, Advance(), exampleMode);
            if (val == null)
                continue;

            // for map
            var childTag = val.Tag;
            if (childTag != null && tag != null && childTag.Type == ValueType.Map)
            {
                childTag.Inherit(tag);
            }

            var keyScalar = new NodeScalar(keyStr, keyStr, Tag.Empty()) { Path = fieldPath };
            entries.Add(new KeyValuePair<NodeScalar, INode>(keyScalar, val));
            lastValue = val;

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

    private INode ParseArray(int openLine, string path)
    {
        _depth++;
        if (_depth > MaxDepth)
            throw new Exception($"max depth: {MaxDepth}");

        Tag tag = ConsumeCommentsFor(openLine) ?? Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
        {
            // Always Vec (not Arr). Go behavior: size is independent of type.
            tag.Type = ValueType.Vec;
        }

        if (!string.IsNullOrEmpty(tag.Name))
        {
            path = $"{path}.{tag.Name}";
        }

        var children = new List<INode>();

        INode? lastItem = null;
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

            var itemPath = $"{path}[{i}]";
            var exampleMode = tag!.Example;
            var item = ParseValue(itemPath, tok, exampleMode);
            if (item == null)
                continue;

            var childTag = item.Tag;
            if (childTag != null && tag != null)
            {
                childTag.Inherit(tag);
            }

            children.Add(item);
            lastItem = item;
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