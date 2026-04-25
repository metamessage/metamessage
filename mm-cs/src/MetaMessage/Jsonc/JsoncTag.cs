namespace MetaMessage.Jsonc;

public class JsoncTag
{
    public ValueType Type { get; set; } = ValueType.Unknown;
    public string? Desc { get; set; }
    public bool Nullable { get; set; }
    public string? DefaultValue { get; set; }
    public string? MinValue { get; set; }
    public string? MaxValue { get; set; }
    public string? Size { get; set; }
    public List<string>? EnumValues { get; set; }
    public string? Pattern { get; set; }
    public string? Location { get; set; }
    public string? Version { get; set; }
    public string? Mime { get; set; }
    public ValueType? ChildType { get; set; }
    public string? ChildDesc { get; set; }
    public string? KeyDesc { get; set; }
    public string? ValueDesc { get; set; }
    public string? EleDesc { get; set; }

    public static JsoncTag Parse(string comment)
    {
        var tag = new JsoncTag();
        if (string.IsNullOrWhiteSpace(comment))
        {
            return tag;
        }

        comment = comment.Trim();
        if (!comment.StartsWith("//") && !comment.StartsWith("/*"))
        {
            return tag;
        }

        if (comment.StartsWith("//"))
        {
            comment = comment[2..].Trim();
        }
        else if (comment.StartsWith("/*"))
        {
            comment = comment[2..].Trim();
            if (comment.EndsWith("*/"))
            {
                comment = comment[..^2].Trim();
            }
        }

        if (!comment.StartsWith("mm:"))
        {
            return tag;
        }

        comment = comment[3..].Trim();
        var pairs = comment.Split(';', StringSplitOptions.RemoveEmptyEntries);
        foreach (var pair in pairs)
        {
            var idx = pair.IndexOf('=');
            if (idx < 0)
            {
                continue;
            }

            var key = pair[..idx].Trim();
            var value = pair[(idx + 1)..].Trim();

            switch (key.ToLower())
            {
                case "type":
                    tag.Type = ValueTypeExtensions.ParseValueType(value);
                    break;
                case "desc":
                    tag.Desc = value;
                    break;
                case "nullable":
                    tag.Nullable = value.Equals("true", StringComparison.OrdinalIgnoreCase);
                    break;
                case "default":
                    tag.DefaultValue = value;
                    break;
                case "min":
                    tag.MinValue = value;
                    break;
                case "max":
                    tag.MaxValue = value;
                    break;
                case "size":
                    tag.Size = value;
                    break;
                case "enum":
                    tag.EnumValues = value.Split(',', StringSplitOptions.RemoveEmptyEntries)
                        .Select(x => x.Trim()).ToList();
                    break;
                case "pattern":
                    tag.Pattern = value;
                    break;
                case "location":
                    tag.Location = value;
                    break;
                case "version":
                    tag.Version = value;
                    break;
                case "mime":
                    tag.Mime = value;
                    break;
                case "child_type":
                    tag.ChildType = ValueTypeExtensions.ParseValueType(value);
                    break;
                case "child_desc":
                    tag.ChildDesc = value;
                    break;
                case "key_desc":
                    tag.KeyDesc = value;
                    break;
                case "value_desc":
                    tag.ValueDesc = value;
                    break;
                case "ele_desc":
                    tag.EleDesc = value;
                    break;
            }
        }

        return tag;
    }

    public override string ToString()
    {
        var parts = new List<string>();
        if (Type != ValueType.Unknown)
        {
            parts.Add($"type={Type.ToTypeString()}");
        }
        if (!string.IsNullOrEmpty(Desc))
        {
            parts.Add($"desc={Desc}");
        }
        if (Nullable)
        {
            parts.Add("nullable=true");
        }
        if (!string.IsNullOrEmpty(DefaultValue))
        {
            parts.Add($"default={DefaultValue}");
        }
        if (!string.IsNullOrEmpty(MinValue))
        {
            parts.Add($"min={MinValue}");
        }
        if (!string.IsNullOrEmpty(MaxValue))
        {
            parts.Add($"max={MaxValue}");
        }
        if (!string.IsNullOrEmpty(Size))
        {
            parts.Add($"size={Size}");
        }
        if (EnumValues != null && EnumValues.Count > 0)
        {
            parts.Add($"enum={string.Join(",", EnumValues)}");
        }
        if (!string.IsNullOrEmpty(Pattern))
        {
            parts.Add($"pattern={Pattern}");
        }
        if (!string.IsNullOrEmpty(Location))
        {
            parts.Add($"location={Location}");
        }
        if (!string.IsNullOrEmpty(Version))
        {
            parts.Add($"version={Version}");
        }
        if (!string.IsNullOrEmpty(Mime))
        {
            parts.Add($"mime={Mime}");
        }
        if (ChildType.HasValue)
        {
            parts.Add($"child_type={ChildType.Value.ToTypeString()}");
        }
        if (!string.IsNullOrEmpty(ChildDesc))
        {
            parts.Add($"child_desc={ChildDesc}");
        }
        if (!string.IsNullOrEmpty(KeyDesc))
        {
            parts.Add($"key_desc={KeyDesc}");
        }
        if (!string.IsNullOrEmpty(ValueDesc))
        {
            parts.Add($"value_desc={ValueDesc}");
        }
        if (!string.IsNullOrEmpty(EleDesc))
        {
            parts.Add($"ele_desc={EleDesc}");
        }

        return "// mm:" + string.Join(";", parts);
    }
}