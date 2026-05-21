using MetaMessage.Core;
using MmVT = MetaMessage.Core.ValueType;

namespace MetaMessage.Jsonc;

public class JsoncParser
{
    private readonly JsoncScanner _scanner;
    private JsoncToken? _currentToken;
    private readonly List<JsoncToken> _pendingComments = new();

    public JsoncParser(string input)
    {
        _scanner = new JsoncScanner(input);
    }

    public IJsoncNode Parse()
    {
        NextToken();
        return ParseValue();
    }

    private JsoncToken NextToken()
    {
        _currentToken = _scanner.NextToken();
        return _currentToken;
    }

    private JsoncToken PeekToken()
    {
        return _currentToken ?? throw new Exception("No current token");
    }

    private JsoncNode ParseValue()
    {
        var token = PeekToken();

        switch (token.Type)
        {
            case JsoncTokenType.LBrace:
                return ParseObject();
            case JsoncTokenType.LBracket:
                return ParseArray();
            case JsoncTokenType.String:
            case JsoncTokenType.Number:
            case JsoncTokenType.True:
            case JsoncTokenType.False:
            case JsoncTokenType.Null:
                return ParsePrimitive();
            default:
                throw new Exception($"Unexpected token: {token.Type} at line {token.Line}, column {token.Column}");
        }
    }

    private JsoncNode ParsePrimitive()
    {
        var token = PeekToken();
        var valueNode = new JsoncValue();

        switch (token.Type)
        {
            case JsoncTokenType.String:
                valueNode.TokenType = JsoncTokenType.String;
                valueNode.Value = token.Literal;
                break;
            case JsoncTokenType.Number:
                valueNode.TokenType = JsoncTokenType.Number;
                if (double.TryParse(token.Literal, out double numValue))
                {
                    valueNode.Value = numValue;
                }
                else
                {
                    valueNode.Value = token.Literal;
                }
                break;
            case JsoncTokenType.True:
                valueNode.TokenType = JsoncTokenType.True;
                valueNode.Value = true;
                break;
            case JsoncTokenType.False:
                valueNode.TokenType = JsoncTokenType.False;
                valueNode.Value = false;
                break;
            case JsoncTokenType.Null:
                valueNode.TokenType = JsoncTokenType.Null;
                valueNode.Value = null;
                break;
        }

        NextToken();

        var tag = ConsumeCommentsFor(token.Line, valueNode);
        if (tag != null)
        {
            valueNode.Tag = tag;
        }

        if (PeekToken().Type == JsoncTokenType.TrailingComment)
        {
            var trailingToken = NextToken();
            valueNode.TrailingComment = new JsoncComment
            {
                Text = trailingToken.Literal,
                Line = trailingToken.Line,
                Column = trailingToken.Column,
                IsBlock = trailingToken.Literal.StartsWith("/*")
            };
            if (trailingToken.Literal.Contains("mm:") && valueNode.Tag == null)
            {
                valueNode.Tag = JsoncTag.Parse(trailingToken.Literal);
            }
        }

        if (valueNode.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(valueNode.Tag);
            var result = Validator.Validate(valueNode.Value, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Value validation failed");
            }
        }

        return valueNode;
    }

    private JsoncObject ParseObject()
    {
        var openingBrace = PeekToken();
        int openingBraceLine = openingBrace.Line;
        NextToken();
        var obj = new JsoncObject();

        var tag = ConsumeCommentsFor(openingBraceLine, obj);
        if (tag != null)
        {
            obj.Tag = tag;
        }

        if (obj.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(obj.Tag);
            var result = Validator.Validate(obj, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Struct validation failed");
            }
        }

        IJsoncNode? lastValue = null;

        while (PeekToken().Type != JsoncTokenType.RBrace && PeekToken().Type != JsoncTokenType.EOF)
        {
            if (PeekToken().Type == JsoncTokenType.LeadingComment)
            {
                var commentToken = NextToken();
                if (_pendingComments.Count > 0)
                {
                    var lastPending = _pendingComments.Last();
                    if (commentToken.Line - lastPending.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(commentToken);
                continue;
            }

            if (PeekToken().Type != JsoncTokenType.String)
            {
                throw new Exception($"Expected string key at line {PeekToken().Line}, column {PeekToken().Column}");
            }

            var keyToken = PeekToken();
            var key = keyToken.Literal;
            NextToken();

            if (PeekToken().Type != JsoncTokenType.Colon)
            {
                throw new Exception($"Expected colon at line {PeekToken().Line}, column {PeekToken().Column}");
            }
            NextToken();

            var value = ParseValue();
            obj.Add(key, value);
            lastValue = value;

            if (PeekToken().Type == JsoncTokenType.Comma)
            {
                NextToken();
                if (PeekToken().Type == JsoncTokenType.TrailingComment)
                {
                    var trailingToken = NextToken();
                    if (lastValue != null)
                    {
                        lastValue.TrailingComment = new JsoncComment
                        {
                            Text = trailingToken.Literal,
                            Line = trailingToken.Line,
                            Column = trailingToken.Column,
                            IsBlock = trailingToken.Literal.StartsWith("/*")
                        };
                        if (trailingToken.Literal.Contains("mm:"))
                        {
                            var trailingTag = JsoncTag.Parse(trailingToken.Literal);
                            if (lastValue.Tag != null)
                            {
                                MergeTag(lastValue.Tag, trailingTag);
                            }
                            else
                            {
                                lastValue.Tag = trailingTag;
                            }
                        }
                    }
                }
            }
        }

        if (PeekToken().Type != JsoncTokenType.RBrace)
        {
            throw new Exception($"Expected closing brace at line {PeekToken().Line}, column {PeekToken().Column}");
        }
        NextToken();

        if (PeekToken().Type == JsoncTokenType.TrailingComment)
        {
            var trailingToken = NextToken();
            obj.TrailingComment = new JsoncComment
            {
                Text = trailingToken.Literal,
                Line = trailingToken.Line,
                Column = trailingToken.Column,
                IsBlock = trailingToken.Literal.StartsWith("/*")
            };
        }

        return obj;
    }

    private JsoncArray ParseArray()
    {
        var openingBracket = PeekToken();
        int openingBracketLine = openingBracket.Line;
        NextToken();
        var array = new JsoncArray();

        var tag = ConsumeCommentsFor(openingBracketLine, array);
        if (tag != null)
        {
            array.Tag = tag;
        }

        if (array.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(array.Tag);
            var result = Validator.Validate(array, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Array validation failed");
            }
        }

        IJsoncNode? lastValue = null;

        while (PeekToken().Type != JsoncTokenType.RBracket && PeekToken().Type != JsoncTokenType.EOF)
        {
            if (PeekToken().Type == JsoncTokenType.LeadingComment)
            {
                var commentToken = NextToken();
                if (_pendingComments.Count > 0)
                {
                    var lastPending = _pendingComments.Last();
                    if (commentToken.Line - lastPending.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(commentToken);
                continue;
            }

            var value = ParseValue();
            array.Add(value);
            lastValue = value;

            if (PeekToken().Type == JsoncTokenType.Comma)
            {
                NextToken();
                if (PeekToken().Type == JsoncTokenType.TrailingComment)
                {
                    var trailingToken = NextToken();
                    if (lastValue != null)
                    {
                        lastValue.TrailingComment = new JsoncComment
                        {
                            Text = trailingToken.Literal,
                            Line = trailingToken.Line,
                            Column = trailingToken.Column,
                            IsBlock = trailingToken.Literal.StartsWith("/*")
                        };
                        if (trailingToken.Literal.Contains("mm:"))
                        {
                            var trailingTag = JsoncTag.Parse(trailingToken.Literal);
                            if (lastValue.Tag != null)
                            {
                                MergeTag(lastValue.Tag, trailingTag);
                            }
                            else
                            {
                                lastValue.Tag = trailingTag;
                            }
                        }
                    }
                }
            }
        }

        if (PeekToken().Type != JsoncTokenType.RBracket)
        {
            throw new Exception($"Expected closing bracket at line {PeekToken().Line}, column {PeekToken().Column}");
        }
        NextToken();

        if (PeekToken().Type == JsoncTokenType.TrailingComment)
        {
            var trailingToken = NextToken();
            array.TrailingComment = new JsoncComment
            {
                Text = trailingToken.Literal,
                Line = trailingToken.Line,
                Column = trailingToken.Column,
                IsBlock = trailingToken.Literal.StartsWith("/*")
            };
        }

        return array;
    }

    private JsoncTag? ConsumeCommentsFor(int anchorLine, JsoncNode? node = null)
    {
        if (_pendingComments.Count == 0)
        {
            return null;
        }

        var lastPending = _pendingComments.Last();
        if (anchorLine - lastPending.Line > 1)
        {
            _pendingComments.Clear();
            return null;
        }

        if (node != null)
        {
            node.LeadingComment = new JsoncComment
            {
                Text = _pendingComments[0].Literal,
                Line = _pendingComments[0].Line,
                Column = _pendingComments[0].Column,
                IsBlock = _pendingComments[0].Literal.StartsWith("/*")
            };
        }

        JsoncTag? mergedTag = null;
        foreach (var commentToken in _pendingComments)
        {
            if (commentToken.Literal.Contains("mm:"))
            {
                var tag = JsoncTag.Parse(commentToken.Literal);
                if (mergedTag == null)
                {
                    mergedTag = tag;
                }
                else
                {
                    MergeTag(mergedTag, tag);
                }
            }
        }

        _pendingComments.Clear();
        return mergedTag;
    }

    private static void MergeTag(JsoncTag target, JsoncTag source)
    {
        if (source.Type != ValueType.Unknown) target.Type = source.Type;
        if (source.Desc != null) target.Desc = source.Desc;
        if (source.Nullable) target.Nullable = true;
        if (source.DefaultValue != null) target.DefaultValue = source.DefaultValue;
        if (source.MinValue != null) target.MinValue = source.MinValue;
        if (source.MaxValue != null) target.MaxValue = source.MaxValue;
        if (source.Size != null) target.Size = source.Size;
        if (source.EnumValues != null) target.EnumValues = source.EnumValues;
        if (source.Pattern != null) target.Pattern = source.Pattern;
        if (source.Location != null) target.Location = source.Location;
        if (source.Version != null) target.Version = source.Version;
        if (source.Mime != null) target.Mime = source.Mime;
        if (source.ChildType.HasValue) target.ChildType = source.ChildType;
        if (source.ChildDesc != null) target.ChildDesc = source.ChildDesc;
        if (source.KeyDesc != null) target.KeyDesc = source.KeyDesc;
        if (source.ValueDesc != null) target.ValueDesc = source.ValueDesc;
        if (source.EleDesc != null) target.EleDesc = source.EleDesc;
    }

    private MmTag ConvertJsoncTagToMmTag(JsoncTag jsoncTag)
    {
        var mmTag = new MmTag();

        switch (jsoncTag.Type)
        {
            case ValueType.String:
                mmTag.Type = MmVT.STRING;
                break;
            case ValueType.Bytes:
                mmTag.Type = MmVT.BYTES;
                break;
            case ValueType.Bool:
                mmTag.Type = MmVT.BOOL;
                break;
            case ValueType.Int:
                mmTag.Type = MmVT.INT;
                break;
            case ValueType.Int8:
                mmTag.Type = MmVT.INT8;
                break;
            case ValueType.Int16:
                mmTag.Type = MmVT.INT16;
                break;
            case ValueType.Int32:
                mmTag.Type = MmVT.INT32;
                break;
            case ValueType.Int64:
                mmTag.Type = MmVT.INT64;
                break;
            case ValueType.Uint:
                mmTag.Type = MmVT.UINT;
                break;
            case ValueType.Uint8:
                mmTag.Type = MmVT.UINT8;
                break;
            case ValueType.Uint16:
                mmTag.Type = MmVT.UINT16;
                break;
            case ValueType.Uint32:
                mmTag.Type = MmVT.UINT32;
                break;
            case ValueType.Uint64:
                mmTag.Type = MmVT.UINT64;
                break;
            case ValueType.Float32:
                mmTag.Type = MmVT.FLOAT32;
                break;
            case ValueType.Float64:
                mmTag.Type = MmVT.FLOAT64;
                break;
            case ValueType.Decimal:
                mmTag.Type = MmVT.DECIMAL;
                break;
            case ValueType.BigInt:
                mmTag.Type = MmVT.BIGINT;
                break;
            case ValueType.DateTime:
                mmTag.Type = MmVT.DATETIME;
                break;
            case ValueType.Date:
                mmTag.Type = MmVT.DATE;
                break;
            case ValueType.Time:
                mmTag.Type = MmVT.TIME;
                break;
            case ValueType.UUID:
                mmTag.Type = MmVT.UUID;
                break;
            case ValueType.IP:
                mmTag.Type = MmVT.IP;
                break;
            case ValueType.URL:
                mmTag.Type = MmVT.URL;
                break;
            case ValueType.Email:
                mmTag.Type = MmVT.EMAIL;
                break;
            case ValueType.Enum:
                mmTag.Type = MmVT.ENUM;
                break;
            case ValueType.Image:
                mmTag.Type = MmVT.IMAGE;
                break;
            case ValueType.Video:
                mmTag.Type = MmVT.VIDEO;
                break;
            case ValueType.Doc:
                mmTag.Type = MmVT.DOC;
                break;
            case ValueType.Struct:
                mmTag.Type = MmVT.OBJ;
                break;
            case ValueType.Array:
                mmTag.Type = MmVT.ARRAY;
                break;
            case ValueType.Slice:
                mmTag.Type = MmVT.VEC;
                break;
            case ValueType.Map:
                mmTag.Type = MmVT.MAP;
                break;
            default:
                mmTag.Type = MmVT.UNKNOWN;
                break;
        }

        mmTag.Nullable = jsoncTag.Nullable;
        mmTag.IsNull = jsoncTag.IsNull;
        mmTag.Min = jsoncTag.MinValue ?? string.Empty;
        mmTag.Max = jsoncTag.MaxValue ?? string.Empty;
        mmTag.Enum = jsoncTag.EnumValues != null ? string.Join("|", jsoncTag.EnumValues) : string.Empty;
        mmTag.DefaultValue = jsoncTag.DefaultValue ?? string.Empty;

        return mmTag;
    }
}