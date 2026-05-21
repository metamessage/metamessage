namespace MetaMessage.Core;

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
    public string Enum { get; set; } = string.Empty;
    public int LocationHours { get; set; } = 0;
    public int Version { get; set; } = 0;
    public string Mime { get; set; } = string.Empty;
    public string Min { get; set; } = string.Empty;
    public string Max { get; set; } = string.Empty;
    public bool Example { get; set; } = false;
    public int Size { get; set; } = 0;
    public string Pattern { get; set; } = string.Empty;
    public string ChildDesc { get; set; } = string.Empty;
    public bool ChildNullable { get; set; } = false;
    public string ChildEnum { get; set; } = string.Empty;
    public bool ChildRaw { get; set; } = false;
    public bool ChildAllowEmpty { get; set; } = false;
    public bool ChildUnique { get; set; } = false;
    public string ChildDefault { get; set; } = string.Empty;
    public string ChildMin { get; set; } = string.Empty;
    public string ChildMax { get; set; } = string.Empty;
    public int ChildSize { get; set; } = 0;
    public string ChildPattern { get; set; } = string.Empty;
    public int ChildLocation { get; set; } = 0;
    public int ChildVersion { get; set; } = 0;
    public string ChildMime { get; set; } = string.Empty;
    public bool IsInherit { get; set; } = false;
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
            Enum = Enum,
            LocationHours = LocationHours,
            Version = Version,
            Mime = Mime,
            Min = Min,
            Max = Max,
            Example = Example,
            Size = Size,
            Pattern = Pattern,
            ChildDesc = ChildDesc,
            ChildNullable = ChildNullable,
            ChildEnum = ChildEnum,
            ChildRaw = ChildRaw,
            ChildAllowEmpty = ChildAllowEmpty,
            ChildUnique = ChildUnique,
            ChildDefault = ChildDefault,
            ChildMin = ChildMin,
            ChildMax = ChildMax,
            ChildSize = ChildSize,
            ChildPattern = ChildPattern,
            ChildLocation = ChildLocation,
            ChildVersion = ChildVersion,
            ChildMime = ChildMime,
            IsInherit = IsInherit,
            IsNull = IsNull
        };
    }

    public void InheritFromArrayParent(MmTag parent)
    {
        IsInherit = true;

        if (!string.IsNullOrEmpty(parent.ChildDesc))
        {
            Desc = parent.ChildDesc;
        }

        if (parent.ChildType != ValueType.UNKNOWN)
        {
            Type = parent.ChildType;
        }

        if (parent.ChildRaw)
        {
            Raw = parent.ChildRaw;
        }

        if (parent.ChildNullable)
        {
            Nullable = parent.ChildNullable;
        }

        if (parent.ChildAllowEmpty)
        {
            AllowEmpty = parent.ChildAllowEmpty;
        }

        if (parent.ChildUnique)
        {
            Unique = parent.ChildUnique;
        }

        if (!string.IsNullOrEmpty(parent.ChildDefault))
        {
            DefaultValue = parent.ChildDefault;
        }

        if (!string.IsNullOrEmpty(parent.ChildMin))
        {
            Min = parent.ChildMin;
        }

        if (!string.IsNullOrEmpty(parent.ChildMax))
        {
            Max = parent.ChildMax;
        }

        if (parent.ChildSize != 0)
        {
            Size = parent.ChildSize;
        }

        if (!string.IsNullOrEmpty(parent.ChildEnum))
        {
            Enum = parent.ChildEnum;
        }

        if (!string.IsNullOrEmpty(parent.ChildPattern))
        {
            Pattern = parent.ChildPattern;
        }

        if (parent.ChildLocation != 0)
        {
            LocationHours = parent.ChildLocation;
        }

        if (parent.ChildVersion != 0)
        {
            Version = parent.ChildVersion;
        }

        if (!string.IsNullOrEmpty(parent.ChildMime))
        {
            Mime = parent.ChildMime;
        }
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

        if (!string.IsNullOrEmpty(Enum))
        {
            result.Add(10);
            var enumBytes = System.Text.Encoding.UTF8.GetBytes(Enum);
            result.Add((byte)enumBytes.Length);
            result.AddRange(enumBytes);
        }

        if (LocationHours != 0)
        {
            result.Add(11);
            var locationStr = System.Text.Encoding.UTF8.GetBytes(LocationHours.ToString());
            result.Add((byte)locationStr.Length);
            result.AddRange(locationStr);
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

        if (!string.IsNullOrEmpty(ChildEnum))
        {
            result.Add(16);
            var childEnumBytes = System.Text.Encoding.UTF8.GetBytes(ChildEnum);
            result.Add((byte)childEnumBytes.Length);
            result.AddRange(childEnumBytes);
        }

        if (Size != 0)
        {
            result.Add(17);
            result.Add((byte)Size);
        }

        if (!string.IsNullOrEmpty(Pattern))
        {
            result.Add(18);
            var patternBytes = System.Text.Encoding.UTF8.GetBytes(Pattern);
            result.Add((byte)patternBytes.Length);
            result.AddRange(patternBytes);
        }

        if (!string.IsNullOrEmpty(ChildPattern))
        {
            result.Add(19);
            var childPatternBytes = System.Text.Encoding.UTF8.GetBytes(ChildPattern);
            result.Add((byte)childPatternBytes.Length);
            result.AddRange(childPatternBytes);
        }

        if (ChildSize != 0)
        {
            result.Add(20);
            result.Add((byte)ChildSize);
        }

        if (!string.IsNullOrEmpty(Min))
        {
            result.Add(21);
            var minBytes = System.Text.Encoding.UTF8.GetBytes(Min);
            result.Add((byte)minBytes.Length);
            result.AddRange(minBytes);
        }

        if (!string.IsNullOrEmpty(Max))
        {
            result.Add(22);
            var maxBytes = System.Text.Encoding.UTF8.GetBytes(Max);
            result.Add((byte)maxBytes.Length);
            result.AddRange(maxBytes);
        }

        if (!string.IsNullOrEmpty(ChildMin))
        {
            result.Add(23);
            var childMinBytes = System.Text.Encoding.UTF8.GetBytes(ChildMin);
            result.Add((byte)childMinBytes.Length);
            result.AddRange(childMinBytes);
        }

        if (!string.IsNullOrEmpty(ChildMax))
        {
            result.Add(24);
            var childMaxBytes = System.Text.Encoding.UTF8.GetBytes(ChildMax);
            result.Add((byte)childMaxBytes.Length);
            result.AddRange(childMaxBytes);
        }

        if (ChildLocation != 0)
        {
            result.Add(25);
            result.Add((byte)ChildLocation);
        }

        if (ChildVersion != 0)
        {
            result.Add(26);
            result.Add((byte)ChildVersion);
        }

        if (!string.IsNullOrEmpty(ChildMime))
        {
            result.Add(27);
            var childMimeBytes = System.Text.Encoding.UTF8.GetBytes(ChildMime);
            result.Add((byte)childMimeBytes.Length);
            result.AddRange(childMimeBytes);
        }

        return result.ToArray();
    }
}