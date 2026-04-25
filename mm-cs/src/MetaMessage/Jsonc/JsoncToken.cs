namespace MetaMessage.Jsonc;

public enum JsoncTokenType
{
    EOF,

    LBrace,
    RBrace,
    LBracket,
    RBracket,
    Colon,
    Comma,

    String,
    Number,
    True,
    False,
    Null,

    LeadingComment,
    TrailingComment
}

public class JsoncToken
{
    public JsoncTokenType Type { get; set; }
    public string Literal { get; set; } = string.Empty;
    public int Line { get; set; }
    public int Column { get; set; }

    public override string ToString()
    {
        return Type switch
        {
            JsoncTokenType.EOF => "EOF",
            JsoncTokenType.LBrace => "{",
            JsoncTokenType.RBrace => "}",
            JsoncTokenType.LBracket => "[",
            JsoncTokenType.RBracket => "]",
            JsoncTokenType.Colon => ":",
            JsoncTokenType.Comma => ",",
            JsoncTokenType.String => $"String(\"{Literal}\")",
            JsoncTokenType.Number => $"Number({Literal})",
            JsoncTokenType.True => "True",
            JsoncTokenType.False => "False",
            JsoncTokenType.Null => "Null",
            JsoncTokenType.LeadingComment => $"LeadingComment({Literal})",
            JsoncTokenType.TrailingComment => $"TrailingComment({Literal})",
            _ => "Unknown"
        };
    }
}