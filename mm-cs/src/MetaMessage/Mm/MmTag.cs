namespace MetaMessage.Mm;

public class MmTag
{
    public string Name { get; set; } = string.Empty;
    public string Desc { get; set; } = string.Empty;
    public ValueType Type { get; set; } = ValueType.UNKNOWN;
    public ValueType ChildType { get; set; } = ValueType.UNKNOWN;
    public bool Nullable { get; set; } = false;
    public bool Raw { get; set; } = false;
    public bool AllowEmpty { get; set; } = false;
    public bool Unique { get; set; } = false;
    public string DefaultValue { get; set; } = string.Empty;
    public List<string> EnumValues { get; set; } = new List<string>();
    public int LocationHours { get; set; } = 0;
    public int Version { get; set; } = 0;
    public string Mime { get; set; } = string.Empty;
    public string ChildDesc { get; set; } = string.Empty;
    public bool ChildNullable { get; set; } = false;
    public List<string> ChildEnum { get; set; } = new List<string>();
    public bool IsNull { get; set; } = false;

    public static MmTag Empty()
    {
        return new MmTag();
    }

    public MmTag Copy()
    {
        return new MmTag
        {
            Name = Name,
            Desc = Desc,
            Type = Type,
            ChildType = ChildType,
            Nullable = Nullable,
            Raw = Raw,
            AllowEmpty = AllowEmpty,
            Unique = Unique,
            DefaultValue = DefaultValue,
            EnumValues = new List<string>(EnumValues),
            LocationHours = LocationHours,
            Version = Version,
            Mime = Mime,
            ChildDesc = ChildDesc,
            ChildNullable = ChildNullable,
            ChildEnum = new List<string>(ChildEnum),
            IsNull = IsNull
        };
    }

    public void InheritFromArrayParent(MmTag parent)
    {
        if (Type == ValueType.UNKNOWN)
        {
            Type = parent.ChildType;
        }
        if (ChildType == ValueType.UNKNOWN)
        {
            ChildType = parent.ChildType;
        }
        if (string.IsNullOrEmpty(Desc))
        {
            Desc = parent.ChildDesc;
        }
        Nullable |= parent.ChildNullable;
    }

    public byte[] ToBytes()
    {
        var result = new List<byte>();

        if (!string.IsNullOrEmpty(Name))
        {
            result.Add(1);
            var nameBytes = System.Text.Encoding.UTF8.GetBytes(Name);
            result.Add((byte)nameBytes.Length);
            result.AddRange(nameBytes);
        }

        if (!string.IsNullOrEmpty(Desc))
        {
            result.Add(2);
            var descBytes = System.Text.Encoding.UTF8.GetBytes(Desc);
            result.Add((byte)descBytes.Length);
            result.AddRange(descBytes);
        }

        if (Type != ValueType.UNKNOWN)
        {
            result.Add(3);
            result.Add((byte)Type);
        }

        if (ChildType != ValueType.UNKNOWN)
        {
            result.Add(4);
            result.Add((byte)ChildType);
        }

        if (Nullable)
        {
            result.Add(5);
        }

        if (Raw)
        {
            result.Add(6);
        }

        if (AllowEmpty)
        {
            result.Add(7);
        }

        if (Unique)
        {
            result.Add(8);
        }

        if (!string.IsNullOrEmpty(DefaultValue))
        {
            result.Add(9);
            var defaultValueBytes = System.Text.Encoding.UTF8.GetBytes(DefaultValue);
            result.Add((byte)defaultValueBytes.Length);
            result.AddRange(defaultValueBytes);
        }

        if (EnumValues.Count > 0)
        {
            result.Add(10);
            result.Add((byte)EnumValues.Count);
            foreach (var enumValue in EnumValues)
            {
                var enumValueBytes = System.Text.Encoding.UTF8.GetBytes(enumValue);
                result.Add((byte)enumValueBytes.Length);
                result.AddRange(enumValueBytes);
            }
        }

        if (LocationHours != 0)
        {
            result.Add(11);
            result.AddRange(BitConverter.GetBytes(LocationHours));
        }

        if (Version != 0)
        {
            result.Add(12);
            result.Add((byte)Version);
        }

        if (!string.IsNullOrEmpty(Mime))
        {
            result.Add(13);
            var mimeBytes = System.Text.Encoding.UTF8.GetBytes(Mime);
            result.Add((byte)mimeBytes.Length);
            result.AddRange(mimeBytes);
        }

        if (!string.IsNullOrEmpty(ChildDesc))
        {
            result.Add(14);
            var childDescBytes = System.Text.Encoding.UTF8.GetBytes(ChildDesc);
            result.Add((byte)childDescBytes.Length);
            result.AddRange(childDescBytes);
        }

        if (ChildNullable)
        {
            result.Add(15);
        }

        if (ChildEnum.Count > 0)
        {
            result.Add(16);
            result.Add((byte)ChildEnum.Count);
            foreach (var childEnumValue in ChildEnum)
            {
                var childEnumValueBytes = System.Text.Encoding.UTF8.GetBytes(childEnumValue);
                result.Add((byte)childEnumValueBytes.Length);
                result.AddRange(childEnumValueBytes);
            }
        }

        return result.ToArray();
    }
}