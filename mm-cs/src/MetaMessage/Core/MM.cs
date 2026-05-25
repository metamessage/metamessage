using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;
namespace MetaMessage.Core;

[AttributeUsage(AttributeTargets.Class | AttributeTargets.Property, AllowMultiple = false)]
public class MM : Attribute
{
    public string Name { get; set; } = string.Empty;
    public string Desc { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string ChildType { get; set; } = string.Empty;
    public bool Nullable { get; set; } = false;
    public bool Deprecated { get; set; } = false;
    public bool AllowEmpty { get; set; } = false;
    public bool Unique { get; set; } = false;
    public string DefaultVal { get; set; } = string.Empty;
    public string Enums { get; set; } = string.Empty;
    public int Location { get; set; } = 0;
    public int Version { get; set; } = 0;
    public string Mime { get; set; } = string.Empty;
    public string ChildDesc { get; set; } = string.Empty;
    public bool ChildNullable { get; set; } = false;
    public string ChildEnums { get; set; } = string.Empty;

    public MM()
    {
    }

    public MM(
        string name = "",
        string desc = "",
        string type = "",
        string childType = "",
        bool nullable = false,
        bool deprecated = false,
        bool allowEmpty = false,
        bool unique = false,
        string defaultVal = "",
        string enums = "",
        int location = 0,
        int version = 0,
        string mime = "",
        string childDesc = "",
        bool childNullable = false,
        string childEnums = "")
    {
        Name = name;
        Desc = desc;
        Type = type;
        ChildType = childType;
        Nullable = nullable;
        Deprecated = deprecated;
        AllowEmpty = allowEmpty;
        Unique = unique;
        DefaultVal = defaultVal;
        Enums = enums;
        Location = location;
        Version = version;
        Mime = mime;
        ChildDesc = childDesc;
        ChildNullable = childNullable;
        ChildEnums = childEnums;
    }

    public Tag ToTag()
    {
        var tag = Tag.Empty();

        if (!string.IsNullOrEmpty(Name))
            tag.Name = Name;
        if (!string.IsNullOrEmpty(Desc))
            tag.Desc = Desc;
        if (!string.IsNullOrEmpty(Type))
        {
            if (Enum.TryParse<ValueType>(Type, true, out var type))
                tag.Type = type;
        }
        if (!string.IsNullOrEmpty(ChildType))
        {
            if (Enum.TryParse<ValueType>(ChildType, true, out var childType))
                tag.ChildType = childType;
        }
        tag.Nullable = Nullable;
        tag.Deprecated = Deprecated;
        tag.AllowEmpty = AllowEmpty;
        tag.Unique = Unique;
        if (!string.IsNullOrEmpty(DefaultVal))
            tag.DefaultVal = DefaultVal;
        if (!string.IsNullOrEmpty(Enums))
            tag.Enums = Enums;
        tag.Location = Location;
        tag.Version = Version;
        if (!string.IsNullOrEmpty(Mime))
            tag.Mime = Mime;
        if (!string.IsNullOrEmpty(ChildDesc))
            tag.ChildDesc = ChildDesc;
        tag.ChildNullable = ChildNullable;
        if (!string.IsNullOrEmpty(ChildEnums))
            tag.ChildEnums = ChildEnums;

        return tag;
    }

    public static MM FromTag(Tag tag)
    {
        return new MM(
            name: tag.Name,
            desc: tag.Desc,
            type: tag.Type != ValueType.Unknown ? tag.Type.ToString() : "",
            childType: tag.ChildType != ValueType.Unknown ? tag.ChildType.ToString() : "",
            nullable: tag.Nullable,
            deprecated: tag.Deprecated,
            allowEmpty: tag.AllowEmpty,
            unique: tag.Unique,
            defaultVal: tag.DefaultVal,
            enums: tag.Enums,
            location: tag.Location,
            version: tag.Version,
            mime: tag.Mime,
            childDesc: tag.ChildDesc,
            childNullable: tag.ChildNullable,
            childEnums: tag.ChildEnums);
    }
}

public static class TagExtensions
{
    public static Tag FromAttribute(MM attribute)
    {
        return attribute.ToTag();
    }
}