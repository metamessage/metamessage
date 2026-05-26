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

    public IJsoncNode Parse()
    {
        while (true)
        {
            var tok = _scanner.NextToken();
            if (tok.Type == JsoncTokenType.EOF)
                return new JsoncValue { TokenType = JsoncTokenType.Null };

            if (tok.Type == JsoncTokenType.LeadingComment)
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

            if (tok.Type == JsoncTokenType.TrailingComment)
            {
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

    private static void MergeNodeTag(IJsoncNode n, Tag parsed)
    {
        if (n == null || parsed == null)
            return;
        var existing = n.Tag;
        n.Tag = Tag.MergeTag(existing, parsed);
    }

    private IJsoncNode ParseValue(string path, JsoncToken firstTok, bool example)
    {
        var tok = firstTok;
        while (true)
        {
            switch (tok.Type)
            {
                case JsoncTokenType.EOF:
                    return new JsoncValue { TokenType = JsoncTokenType.Null };

                case JsoncTokenType.LeadingComment:
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
                        tok = Advance();
                        continue;
                    }

                case JsoncTokenType.TrailingComment:
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

    private IJsoncNode ParseString(JsoncToken tok, string path, bool example)
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

        var val = new JsoncValue
        {
            Value = text,
            TokenType = JsoncTokenType.String,
            Tag = tag,
            Path = path
        };

        // Handle trailing comments
        var next = Peek();
        if (next.Type == JsoncTokenType.TrailingComment)
        {
            var trailing = Advance();
            var parsed = ParseCommentsToTag(trailing.Literal);
            if (parsed != null)
            {
                MergeNodeTag(val, parsed);
            }
        }

        return val;
    }

    private IJsoncNode ParseNumber(JsoncToken tok, string path, bool example)
    {
        var tag = ConsumeCommentsFor(tok.Line);
        var text = tok.Literal;

        if (tag == null)
            tag = Tag.NewTag();

        if (text.Contains("."))
        {
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.F64;
        }
        else
        {
            if (tag.Type == ValueType.Unknown)
                tag.Type = ValueType.I;
        }

        var val = new JsoncValue
        {
            Value = double.Parse(text, System.Globalization.CultureInfo.InvariantCulture),
            TokenType = JsoncTokenType.Number,
            Tag = tag,
            Path = path
        };

        var next = Peek();
        if (next.Type == JsoncTokenType.TrailingComment)
        {
            var trailing = Advance();
            var parsed = ParseCommentsToTag(trailing.Literal);
            if (parsed != null)
            {
                MergeNodeTag(val, parsed);
            }
        }

        return val;
    }

    private IJsoncNode ParseBool(JsoncToken tok, string path, bool example)
    {
        var tag = ConsumeCommentsFor(tok.Line);

        if (tag == null)
            tag = Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
            tag.Type = ValueType.Bool;

        bool boolVal = tok.Type == JsoncTokenType.True;

        var val = new JsoncValue
        {
            Value = boolVal,
            TokenType = tok.Type,
            Tag = tag,
            Path = path
        };

        var next = Peek();
        if (next.Type == JsoncTokenType.TrailingComment)
        {
            var trailing = Advance();
            var parsed = ParseCommentsToTag(trailing.Literal);
            if (parsed != null)
            {
                MergeNodeTag(val, parsed);
            }
        }

        return val;
    }

    private IJsoncNode ParseObject(int openLine, string path)
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

        var obj = new JsoncObject
        {
            Tag = tag,
            Path = path
        };

        IJsoncNode? lastValue = null;

        while (true)
        {
            var tok = Advance();
            if (tok.Type == JsoncTokenType.EOF)
                break;
            if (tok.Type == JsoncTokenType.RBrace)
                break;

            if (tok.Type == JsoncTokenType.LeadingComment)
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

            if (tok.Type == JsoncTokenType.TrailingComment)
            {
                if (lastValue != null)
                {
                    var parsed = ParseCommentsToTag(tok.Literal);
                    if (parsed != null)
                    {
                        MergeNodeTag(lastValue, parsed);
                    }
                }
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

            obj.Add(keyStr, val);
            lastValue = val;

            var nextTok = Peek();
            if (nextTok.Type == JsoncTokenType.Comma)
            {
                Advance();
            }
        }

        if (!tag!.Example)
        {
            switch (tag!.Type)
            {
                case ValueType.Map:
                    Validator.Validate(obj, tag);
                    break;
                case ValueType.Obj:
                    Validator.Validate(obj, tag);
                    break;
            }
        }

        _depth--;
        return obj;
    }

    private IJsoncNode ParseArray(int openLine, string path)
    {
        _depth++;
        if (_depth > MaxDepth)
            throw new Exception($"max depth: {MaxDepth}");

        Tag tag = ConsumeCommentsFor(openLine) ?? Tag.NewTag();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = tag.Size > 0 ? ValueType.Arr : ValueType.Vec;
        }

        if (!string.IsNullOrEmpty(tag.Name))
        {
            path = $"{path}.{tag.Name}";
        }

        var arr = new JsoncArray
        {
            Tag = tag,
            Path = path
        };

        IJsoncNode? lastItem = null;
        int i = 0;

        while (true)
        {
            var tok = Advance();
            if (tok.Type == JsoncTokenType.EOF)
                break;
            if (tok.Type == JsoncTokenType.RBracket)
                break;

            if (tok.Type == JsoncTokenType.LeadingComment)
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

            if (tok.Type == JsoncTokenType.TrailingComment)
            {
                if (lastItem != null)
                {
                    var parsed = ParseCommentsToTag(tok.Literal);
                    if (parsed != null)
                    {
                        MergeNodeTag(lastItem, parsed);
                    }
                }
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

            arr.Add(item);
            lastItem = item;
            i++;

            var nextTok = Peek();
            if (nextTok.Type == JsoncTokenType.Comma)
            {
                Advance();
            }
        }

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