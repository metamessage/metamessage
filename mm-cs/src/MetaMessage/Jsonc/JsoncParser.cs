using MetaMessage.Mm;

namespace MetaMessage.Jsonc;

public class JsoncParser
{
    private readonly JsoncScanner _scanner;
    private JsoncToken? _currentToken;
    private JsoncToken? _lastToken;

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
        _lastToken = _currentToken;
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
        var token = NextToken();
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

        ProcessTrailingComment(valueNode);

        // 验证值
        if (valueNode.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(valueNode.Tag);
            var result = MmValidator.Validate(valueNode.Value, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Value validation failed");
            }
        }

        return valueNode;
    }

    private JsoncObject ParseObject()
    {
        NextToken();
        var obj = new JsoncObject();

        ProcessLeadingComment(obj);

        // 验证结构体 tag
        if (obj.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(obj.Tag);
            var result = MmValidator.Validate(obj, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Struct validation failed");
            }
        }

        while (PeekToken().Type != JsoncTokenType.RBrace && PeekToken().Type != JsoncTokenType.EOF)
        {
            if (PeekToken().Type == JsoncTokenType.LeadingComment)
            {
                NextToken();
                continue;
            }

            if (PeekToken().Type != JsoncTokenType.String)
            {
                throw new Exception($"Expected string key at line {PeekToken().Line}, column {PeekToken().Column}");
            }

            var keyToken = NextToken();
            var key = keyToken.Literal;

            if (PeekToken().Type != JsoncTokenType.Colon)
            {
                throw new Exception($"Expected colon at line {PeekToken().Line}, column {PeekToken().Column}");
            }
            NextToken();

            var value = ParseValue();
            obj.Add(key, value);

            if (PeekToken().Type == JsoncTokenType.Comma)
            {
                NextToken();
                ProcessTrailingComment(value);
            }
        }

        if (PeekToken().Type != JsoncTokenType.RBrace)
        {
            throw new Exception($"Expected closing brace at line {PeekToken().Line}, column {PeekToken().Column}");
        }
        NextToken();

        ProcessTrailingComment(obj);
        return obj;
    }

    private JsoncArray ParseArray()
    {
        NextToken();
        var array = new JsoncArray();

        ProcessLeadingComment(array);

        // 验证数组 tag
        if (array.Tag != null)
        {
            var mmTag = ConvertJsoncTagToMmTag(array.Tag);
            var result = MmValidator.Validate(array, mmTag);
            if (!result.IsValid)
            {
                throw new Exception(string.Join(", ", result.Errors) ?? "Array validation failed");
            }
        }

        while (PeekToken().Type != JsoncTokenType.RBracket && PeekToken().Type != JsoncTokenType.EOF)
        {
            if (PeekToken().Type == JsoncTokenType.LeadingComment)
            {
                NextToken();
                continue;
            }

            var value = ParseValue();
            array.Add(value);

            if (PeekToken().Type == JsoncTokenType.Comma)
            {
                NextToken();
                ProcessTrailingComment(value);
            }
        }

        if (PeekToken().Type != JsoncTokenType.RBracket)
        {
            throw new Exception($"Expected closing bracket at line {PeekToken().Line}, column {PeekToken().Column}");
        }
        NextToken();

        ProcessTrailingComment(array);
        return array;
    }

    private void ProcessLeadingComment(JsoncNode node)
    {
        if (_lastToken != null && _lastToken.Type == JsoncTokenType.LeadingComment)
        {
            var comment = new JsoncComment
            {
                Text = _lastToken.Literal,
                Line = _lastToken.Line,
                Column = _lastToken.Column,
                IsBlock = _lastToken.Literal.StartsWith("/*")
            };
            node.LeadingComment = comment;

            if (_lastToken.Literal.Contains("mm:"))
            {
                node.Tag = JsoncTag.Parse(_lastToken.Literal);
            }
        }
    }

    private void ProcessTrailingComment(JsoncNode node)
    {
        if (PeekToken().Type == JsoncTokenType.TrailingComment)
        {
            var comment = new JsoncComment
            {
                Text = _currentToken!.Literal,
                Line = _currentToken.Line,
                Column = _currentToken.Column,
                IsBlock = _currentToken.Literal.StartsWith("/*")
            };
            node.TrailingComment = comment;
        }
    }

    private MmTag ConvertJsoncTagToMmTag(JsoncTag jsoncTag)
    {
        var mmTag = new MmTag();

        // 转换类型
        switch (jsoncTag.Type)
        {
            case ValueType.String:
                mmTag.Type = ValueType.STRING;
                break;
            case ValueType.Number:
                mmTag.Type = ValueType.INT;
                break;
            case ValueType.Boolean:
                mmTag.Type = ValueType.BOOL;
                break;
            case ValueType.Null:
                mmTag.Type = ValueType.NULL;
                break;
            case ValueType.Object:
                mmTag.Type = ValueType.STRUCT;
                break;
            case ValueType.Array:
                mmTag.Type = ValueType.ARRAY;
                break;
        }

        // 转换其他属性
        mmTag.Nullable = jsoncTag.Nullable;
        if (!string.IsNullOrEmpty(jsoncTag.MinValue) && double.TryParse(jsoncTag.MinValue, out double min))
        {
            mmTag.Min = min;
        }
        if (!string.IsNullOrEmpty(jsoncTag.MaxValue) && double.TryParse(jsoncTag.MaxValue, out double max))
        {
            mmTag.Max = max;
        }
        if (!string.IsNullOrEmpty(jsoncTag.Size) && int.TryParse(jsoncTag.Size, out int size))
        {
            mmTag.Size = size;
        }
        mmTag.Pattern = jsoncTag.Pattern;
        mmTag.EnumValues = jsoncTag.EnumValues;
        mmTag.DefaultValue = jsoncTag.DefaultValue;

        return mmTag;
    }
}