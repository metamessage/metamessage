namespace MetaMessage.Jsonc;

public interface IJsoncNode
{
    JsoncTokenType TokenType { get; }
    JsoncTag? Tag { get; set; }
    JsoncComment? LeadingComment { get; set; }
    JsoncComment? TrailingComment { get; set; }
}

public class JsoncComment
{
    public string Text { get; set; } = string.Empty;
    public int Line { get; set; }
    public int Column { get; set; }
    public bool IsBlock { get; set; }
}

public abstract class JsoncNode : IJsoncNode
{
    public JsoncTokenType TokenType { get; set; }
    public JsoncTag? Tag { get; set; }
    public JsoncComment? LeadingComment { get; set; }
    public JsoncComment? TrailingComment { get; set; }
}

public class JsoncValue : JsoncNode
{
    public object? Value { get; set; }

    public JsoncValue()
    {
        TokenType = JsoncTokenType.String;
    }

    public JsoncValue(object? value, JsoncTokenType tokenType)
    {
        Value = value;
        TokenType = tokenType;
    }

    public string? GetString()
    {
        return Value as string;
    }

    public bool? GetBool()
    {
        if (Value is bool b)
        {
            return b;
        }
        return null;
    }

    public double? GetNumber()
    {
        if (Value is double d)
        {
            return d;
        }
        if (Value is int i)
        {
            return i;
        }
        if (Value is long l)
        {
            return l;
        }
        if (Value is float f)
        {
            return f;
        }
        if (Value is string s && double.TryParse(s, out double result))
        {
            return result;
        }
        return null;
    }

    public bool IsNull()
    {
        return Value == null || Value is string s && s == "null";
    }
}

public class JsoncObject : JsoncNode
{
    public Dictionary<string, IJsoncNode> Fields { get; set; } = new();

    public JsoncObject()
    {
        TokenType = JsoncTokenType.LBrace;
    }

    public void Add(string key, IJsoncNode node)
    {
        Fields[key] = node;
    }

    public IJsoncNode? Get(string key)
    {
        return Fields.TryGetValue(key, out var node) ? node : null;
    }

    public JsoncValue? GetValue(string key)
    {
        return Get(key) as JsoncValue;
    }

    public JsoncObject? GetObject(string key)
    {
        return Get(key) as JsoncObject;
    }

    public JsoncArray? GetArray(string key)
    {
        return Get(key) as JsoncArray;
    }
}

public class JsoncArray : JsoncNode
{
    public List<IJsoncNode> Elements { get; set; } = new();

    public JsoncArray()
    {
        TokenType = JsoncTokenType.LBracket;
    }

    public void Add(IJsoncNode node)
    {
        Elements.Add(node);
    }

    public IJsoncNode? Get(int index)
    {
        if (index >= 0 && index < Elements.Count)
        {
            return Elements[index];
        }
        return null;
    }

    public JsoncValue? GetValue(int index)
    {
        return Get(index) as JsoncValue;
    }

    public JsoncObject? GetObject(int index)
    {
        return Get(index) as JsoncObject;
    }

    public JsoncArray? GetArray(int index)
    {
        return Get(index) as JsoncArray;
    }
}