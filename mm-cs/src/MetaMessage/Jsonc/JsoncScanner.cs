namespace MetaMessage.Jsonc;

public class JsoncScanner
{
    private readonly string _input;
    private int _position;
    private int _line;
    private int _column;
    private JsoncToken? _lastToken;
    private JsoncToken? _currentToken;

    public JsoncScanner(string input)
    {
        _input = input;
        _position = 0;
        _line = 1;
        _column = 1;
    }

    public JsoncToken CurrentToken => _currentToken ?? throw new InvalidOperationException("No current token");

    public JsoncToken LastToken => _lastToken ?? throw new InvalidOperationException("No last token");

    public int Position => _position;

    public JsoncToken NextToken()
    {
        SkipWhitespace();
        if (_position >= _input.Length)
        {
            _lastToken = _currentToken;
            _currentToken = new JsoncToken
            {
                Type = JsoncTokenType.EOF,
                Line = _line,
                Column = _column
            };
            return _currentToken;
        }

        char ch = _input[_position];

        if (ch == '/')
        {
            return ScanComment();
        }

        var token = ch switch
        {
            '{' => CreateToken(JsoncTokenType.LBrace),
            '}' => CreateToken(JsoncTokenType.RBrace),
            '[' => CreateToken(JsoncTokenType.LBracket),
            ']' => CreateToken(JsoncTokenType.RBracket),
            ':' => CreateToken(JsoncTokenType.Colon),
            ',' => CreateToken(JsoncTokenType.Comma),
            '"' => ScanString(),
            _ when IsDigit(ch) || ch == '-' => ScanNumber(),
            _ when IsLetter(ch) => ScanIdentifier(),
            _ => throw new Exception($"Unexpected character: {ch} at line {_line}, column {_column}")
        };

        _lastToken = _currentToken;
        _currentToken = token;
        return token;
    }

    private JsoncToken CreateToken(JsoncTokenType type, string literal = "")
    {
        var token = new JsoncToken
        {
            Type = type,
            Literal = literal,
            Line = _line,
            Column = _column
        };
        Advance(1);
        return token;
    }

    private JsoncToken ScanComment()
    {
        if (_input.Length <= _position + 1)
        {
            return CreateToken(JsoncTokenType.LeadingComment, "/");
        }

        char next = _input[_position + 1];
        string comment;
        bool isLeading = _lastToken == null || _lastToken.Type == JsoncTokenType.Comma ||
                          _lastToken.Type == JsoncTokenType.Colon ||
                          _lastToken.Type == JsoncTokenType.LBrace ||
                          _lastToken.Type == JsoncTokenType.LBracket;

        if (next == '/')
        {
            int startPos = _position;
            int startLine = _line;
            int startColumn = _column;
            Advance(2);
            while (_position < _input.Length && _input[_position] != '\n')
            {
                Advance(1);
            }
            comment = _input.Substring(startPos, _position - startPos);
            var token = new JsoncToken
            {
                Type = isLeading ? JsoncTokenType.LeadingComment : JsoncTokenType.TrailingComment,
                Literal = comment,
                Line = startLine,
                Column = startColumn
            };
            _lastToken = _currentToken;
            _currentToken = token;
            return token;
        }
        else if (next == '*')
        {
            int startPos = _position;
            int startLine = _line;
            int startColumn = _column;
            Advance(2);
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
            comment = _input.Substring(startPos, _position - startPos);
            var token = new JsoncToken
            {
                Type = isLeading ? JsoncTokenType.LeadingComment : JsoncTokenType.TrailingComment,
                Literal = comment,
                Line = startLine,
                Column = startColumn
            };
            _lastToken = _currentToken;
            _currentToken = token;
            return token;
        }

        return CreateToken(JsoncTokenType.LeadingComment, "/");
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
                Advance(1);
                char escaped = _input[_position];
                switch (escaped)
                {
                    case 'n':
                        sb.Append('\n');
                        break;
                    case 'r':
                        sb.Append('\r');
                        break;
                    case 't':
                        sb.Append('\t');
                        break;
                    case 'b':
                        sb.Append('\b');
                        break;
                    case 'f':
                        sb.Append('\f');
                        break;
                    case '"':
                        sb.Append('"');
                        break;
                    case '\\':
                        sb.Append('\\');
                        break;
                    case 'u':
                        if (_position + 4 < _input.Length)
                        {
                            Advance(1);
                            string hex = _input.Substring(_position, 4);
                            if (int.TryParse(hex, System.Globalization.NumberStyles.HexNumber, null, out int unicode))
                            {
                                sb.Append((char)unicode);
                                Advance(3);
                            }
                            else
                            {
                                sb.Append(escaped);
                            }
                        }
                        else
                        {
                            sb.Append(escaped);
                        }
                        break;
                    default:
                        sb.Append(escaped);
                        break;
                }
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

        bool isFloat = false;
        if (_position < _input.Length && (_input[_position] == 'f' || _input[_position] == 'F'))
        {
            sb.Append(_input[_position]);
            Advance(1);
            isFloat = true;
        }

        if (!isFloat && sb.ToString().Contains('.'))
        {
            isFloat = true;
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