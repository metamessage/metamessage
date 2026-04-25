namespace MetaMessage.Mm;

[AttributeUsage(AttributeTargets.Class | AttributeTargets.Property, AllowMultiple = false)]
public class MM : Attribute
{
    public string Name { get; set; } = string.Empty;
    public string Desc { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string ChildType { get; set; } = string.Empty;
    public bool Nullable { get; set; } = false;
    public bool Raw { get; set; } = false;
    public bool AllowEmpty { get; set; } = false;
    public bool Unique { get; set; } = false;
    public string DefaultValue { get; set; } = string.Empty;
    public string[] EnumValues { get; set; } = Array.Empty<string>();
    public int LocationHours { get; set; } = 0;
    public int Version { get; set; } = 0;
    public string Mime { get; set; } = string.Empty;
    public string ChildDesc { get; set; } = string.Empty;
    public bool ChildNullable { get; set; } = false;
    public string[] ChildEnum { get; set; } = Array.Empty<string>();
}

public static class MmTagExtensions
{
    public static MmTag FromAttribute(MM attribute)
    {
        var tag = MmTag.Empty();
        
        if (!string.IsNullOrEmpty(attribute.Name))
            tag.Name = attribute.Name;
        if (!string.IsNullOrEmpty(attribute.Desc))
            tag.Desc = attribute.Desc;
        if (!string.IsNullOrEmpty(attribute.Type))
        {
            if (Enum.TryParse<ValueType>(attribute.Type, true, out var type))
                tag.Type = type;
        }
        if (!string.IsNullOrEmpty(attribute.ChildType))
        {
            if (Enum.TryParse<ValueType>(attribute.ChildType, true, out var childType))
                tag.ChildType = childType;
        }
        tag.Nullable = attribute.Nullable;
        tag.Raw = attribute.Raw;
        tag.AllowEmpty = attribute.AllowEmpty;
        tag.Unique = attribute.Unique;
        if (!string.IsNullOrEmpty(attribute.DefaultValue))
            tag.DefaultValue = attribute.DefaultValue;
        if (attribute.EnumValues != null && attribute.EnumValues.Length > 0)
            tag.EnumValues = new List<string>(attribute.EnumValues);
        tag.LocationHours = attribute.LocationHours;
        tag.Version = attribute.Version;
        if (!string.IsNullOrEmpty(attribute.Mime))
            tag.Mime = attribute.Mime;
        if (!string.IsNullOrEmpty(attribute.ChildDesc))
            tag.ChildDesc = attribute.ChildDesc;
        tag.ChildNullable = attribute.ChildNullable;
        if (attribute.ChildEnum != null && attribute.ChildEnum.Length > 0)
            tag.ChildEnum = new List<string>(attribute.ChildEnum);
        
        return tag;
    }
}