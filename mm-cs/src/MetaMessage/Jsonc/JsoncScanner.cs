namespace MetaMessage.Jsonc;

public class JsoncScanner
{
    private readonly string _input;
    private int _position;
    private int _line;
    private int _column;
    private bool _newLine;

    public JsoncScanner(string input)
    {
        _input = input;
        _position = 0;
        _line = 1;
        _column = 1;
    }

    public JsoncToken NextToken()
    {
        SkipWhitespace();
        if (_position >= _input.Length)
        {
            return new JsoncToken
            {
                Type = JsoncTokenType.EOF,
                Line = _line,
                Column = _column
            };
        }

        char ch = _input[_position];

        if (ch == '/')
        {
            return ScanComment();
        }

        switch (ch)
        {
            case '{':
                Advance(1);
                return NewToken(JsoncTokenType.LBrace);
            case '}':
                Advance(1);
                return NewToken(JsoncTokenType.RBrace);
            case '[':
                Advance(1);
                return NewToken(JsoncTokenType.LBracket);
            case ']':
                Advance(1);
                return NewToken(JsoncTokenType.RBracket);
            case ':':
                Advance(1);
                _newLine = false;
                return NewToken(JsoncTokenType.Colon);
            case ',':
                Advance(1);
                _newLine = false;
                return NewToken(JsoncTokenType.Comma);
            case '"':
                return ScanString();
            default:
                if (IsDigit(ch) || ch == '-')
                {
                    return ScanNumber();
                }
                if (IsLetter(ch))
                {
                    return ScanIdentifier();
                }
                throw new Exception($"Unexpected character: {ch} at line {_line}, column {_column}");
        }
    }

    private JsoncToken NewToken(JsoncTokenType type, string literal = "")
    {
        return new JsoncToken
        {
            Type = type,
            Literal = literal,
            Line = _line,
            Column = _column
        };
    }

    private JsoncToken ScanComment()
    {
        if (_input.Length <= _position + 1)
        {
            Advance(1);
            return NewToken(JsoncTokenType.LeadingComment, "/");
        }

        char next = _input[_position + 1];
        int startLine = _line;
        int startColumn = _column;

        if (next == '/')
        {
            Advance(2);
            int contentStart = _position;
            while (_position < _input.Length && _input[_position] != '\n')
            {
                Advance(1);
            }
            string content = _input.Substring(contentStart, _position - contentStart).Trim();
            return new JsoncToken
            {
                Type = _newLine ? JsoncTokenType.LeadingComment : JsoncTokenType.TrailingComment,
                Literal = content,
                Line = startLine,
                Column = startColumn
            };
        }

        if (next == '*')
        {
            Advance(2);
            int contentStart = _position;
            while (_position < _input.Length - 1)
            {
                if (_input[_position] == '*' && _input[_position + 1] == '/')
                {
                    Advance(2);
                    break;
                }
                if (_input[_position] == '\n')
                {
                    _line++;
                    _column = 1;
                }
                Advance(1);
            }
            string content = _input.Substring(contentStart, _position - contentStart - 2).Trim();
            return new JsoncToken
            {
                Type = _newLine ? JsoncTokenType.LeadingComment : JsoncTokenType.TrailingComment,
                Literal = content,
                IsBlock = true,
                Line = startLine,
                Column = startColumn
            };
        }

        Advance(1);
        return NewToken(JsoncTokenType.LeadingComment, "/");
    }

    private JsoncToken ScanString()
    {
        int startLine = _line;
        int startColumn = _column;
        Advance(1);
        var sb = new System.Text.StringBuilder();
        while (_position < _input.Length && _input[_position] != '"')
        {
            if (_input[_position] == '\\' && _position + 1 < _input.Length)
            {
                sb.Append('\\');
                Advance(1);
                sb.Append(_input[_position]);
            }
            else
            {
                if (_input[_position] == '\n')
                {
                    _line++;
                    _column = 1;
                }
                sb.Append(_input[_position]);
            }
            Advance(1);
        }

        if (_position < _input.Length)
        {
            Advance(1);
        }

        return new JsoncToken
        {
            Type = JsoncTokenType.String,
            Literal = sb.ToString(),
            Line = startLine,
            Column = startColumn
        };
    }

    private JsoncToken ScanNumber()
    {
        int startLine = _line;
        int startColumn = _column;
        var sb = new System.Text.StringBuilder();

        if (_input[_position] == '-')
        {
            sb.Append('-');
            Advance(1);
        }

        while (_position < _input.Length && (IsDigit(_input[_position]) || _input[_position] == '.' ||
               _input[_position] == 'e' || _input[_position] == 'E' ||
               _input[_position] == '+' || _input[_position] == '_'))
        {
            if (_input[_position] == '_')
            {
                Advance(1);
                continue;
            }
            sb.Append(_input[_position]);
            Advance(1);
        }

        return new JsoncToken
        {
            Type = JsoncTokenType.Number,
            Literal = sb.ToString(),
            Line = startLine,
            Column = startColumn
        };
    }

    private JsoncToken ScanIdentifier()
    {
        int startLine = _line;
        int startColumn = _column;
        var sb = new System.Text.StringBuilder();

        while (_position < _input.Length && (IsLetter(_input[_position]) || IsDigit(_input[_position])))
        {
            sb.Append(_input[_position]);
            Advance(1);
        }

        string identifier = sb.ToString();
        return identifier.ToLower() switch
        {
            "true" => new JsoncToken { Type = JsoncTokenType.True, Literal = "true", Line = startLine, Column = startColumn },
            "false" => new JsoncToken { Type = JsoncTokenType.False, Literal = "false", Line = startLine, Column = startColumn },
            "null" => new JsoncToken { Type = JsoncTokenType.Null, Literal = "null", Line = startLine, Column = startColumn },
            _ => new JsoncToken { Type = JsoncTokenType.String, Literal = identifier, Line = startLine, Column = startColumn }
        };
    }

    private void SkipWhitespace()
    {
        while (_position < _input.Length && IsWhitespace(_input[_position]))
        {
            if (_input[_position] == '\n')
            {
                _newLine = true;
                _line++;
                _column = 1;
            }
            Advance(1);
        }
    }

    private void Advance(int count)
    {
        _position += count;
        _column += count;
    }

    private static bool IsWhitespace(char ch)
    {
        return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r';
    }

    private static bool IsDigit(char ch)
    {
        return ch >= '0' && ch <= '9';
    }

    private static bool IsLetter(char ch)
    {
        return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z');
    }
}