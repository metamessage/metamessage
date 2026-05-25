using MetaMessage.Core;
using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

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
        IJsoncNode? result = null;

        while (true)
        {
            var token = NextToken();

            if (token.Type == JsoncTokenType.EOF)
                break;

            if (token.Type == JsoncTokenType.LeadingComment)
            {
                if (_pendingComments.Count > 0)
                {
                    var last = _pendingComments.Last();
                    if (token.Line - last.Line > 1)
                    {
                        _pendingComments.Clear();
                    }
                }
                _pendingComments.Add(token);
                continue;
            }

            if (token.Type == JsoncTokenType.TrailingComment)
            {
                continue;
            }

            result = ParseValue();

            while (true)
            {
                var next = PeekToken();
                if (next.Type == JsoncTokenType.EOF)
                    break;
                if (next.Type == JsoncTokenType.TrailingComment)
                {
                    NextToken();
                    var parsed = ParseCommentsToTag(next.Literal);
                    if (parsed != null && result != null)
                    {
                        var existing = result.Tag;
                        if (existing != null)
                        {
                            result.Tag = Tag.MergeTag(existing, parsed);
                        }
                        else
                        {
                            result.Tag = parsed;
                        }
                    }
                }
                else
                {
                    break;
                }
            }

            break;
        }

        return result!;
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
                IsBlock = trailingToken.IsBlock
            };
            var parsed = ParseCommentsToTag(trailingToken.Literal);
            if (parsed != null)
            {
                if (valueNode.Tag != null)
                {
                    valueNode.Tag = Tag.MergeTag(valueNode.Tag, parsed);
                }
                else
                {
                    valueNode.Tag = parsed;
                }
            }
        }

        if (valueNode.Tag != null)
        {
            var mmTag = valueNode.Tag;
            var validationValue = ConvertForValidation(valueNode.Value, mmTag.Type);
            if (validationValue != null || valueNode.Value == null)
            {
                var result = Validator.Validate(validationValue!, mmTag);
                if (!result.IsValid)
                {
                    throw new Exception(string.Join(", ", result.Errors) ?? "Value validation failed");
                }
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
            var mmTag = obj.Tag;
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
                var commentToken = PeekToken();
                NextToken();
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
                            IsBlock = trailingToken.IsBlock
                        };
                        var parsed = ParseCommentsToTag(trailingToken.Literal);
                        if (parsed != null)
                        {
                            if (lastValue.Tag != null)
                            {
                                lastValue.Tag = Tag.MergeTag(lastValue.Tag, parsed);
                            }
                            else
                            {
                                lastValue.Tag = parsed;
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
                IsBlock = trailingToken.IsBlock
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
            var mmTag = array.Tag;
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
                var commentToken = PeekToken();
                NextToken();
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
                            IsBlock = trailingToken.IsBlock
                        };
                        var parsed = ParseCommentsToTag(trailingToken.Literal);
                        if (parsed != null)
                        {
                            if (lastValue.Tag != null)
                            {
                                lastValue.Tag = Tag.MergeTag(lastValue.Tag, parsed);
                            }
                            else
                            {
                                lastValue.Tag = parsed;
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
                IsBlock = trailingToken.IsBlock
            };
        }

        return array;
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

    private Tag? ConsumeCommentsFor(int anchorLine, JsoncNode? node = null)
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
                IsBlock = _pendingComments[0].IsBlock
            };
        }

        Tag? mergedTag = null;
        foreach (var commentToken in _pendingComments)
        {
            var parsed = ParseCommentsToTag(commentToken.Literal);
            if (parsed != null)
            {
                if (mergedTag == null)
                {
                    mergedTag = parsed;
                }
                else
                {
                    mergedTag = Tag.MergeTag(mergedTag, parsed);
                }
            }
        }

        _pendingComments.Clear();
        return mergedTag;
    }

    private static object? ConvertForValidation(object? value, ValueType type)
    {
        if (value == null) return null;

        switch (type)
        {
            case ValueType.I:
            case ValueType.I8:
            case ValueType.I16:
            case ValueType.I32:
                if (value is string s) return int.Parse(s);
                return Convert.ToInt32(value);
            case ValueType.I64:
                if (value is string s64) return long.Parse(s64);
                return Convert.ToInt64(value);
            case ValueType.U:
            case ValueType.U8:
            case ValueType.U16:
            case ValueType.U32:
                if (value is string us) return uint.Parse(us);
                return Convert.ToUInt32(value);
            case ValueType.U64:
                if (value is string u64s) return ulong.Parse(u64s);
                return Convert.ToUInt64(value);
            case ValueType.F32:
            case ValueType.F64:
                if (value is string fs) return double.Parse(fs);
                return Convert.ToDouble(value);
            case ValueType.Bool:
                return value;
            default:
                return value;
        }
    }
}