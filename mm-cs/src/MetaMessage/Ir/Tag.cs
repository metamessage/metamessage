using System.Text;
using MetaMessage.Core;

namespace MetaMessage.Ir;

public static class TagConstants
{
    public const string TIsNull = "is_null";
    public const string TExample = "example";
    public const string TDeprecated = "deprecated";

    public const string TName = "name";
    public const string TDesc = "desc";
    public const string TType = "type";
    public const string TNullable = "nullable";
    public const string TAllowEmpty = "allow_empty";
    public const string TUnique = "unique";
    public const string TDefaultVal = "default_val";
    public const string TMin = "min";
    public const string TMax = "max";
    public const string TSize = "size";
    public const string TEnum = "enums";
    public const string TPattern = "pattern";
    public const string TLocation = "location";
    public const string TVersion = "version";
    public const string TMime = "mime";

    public const string TChildDesc = "child_desc";
    public const string TChildType = "child_type";
    public const string TChildNullable = "child_nullable";
    public const string TChildAllowEmpty = "child_allow_empty";
    public const string TChildUnique = "child_unique";
    public const string TChildDefaultVal = "child_default_val";
    public const string TChildMin = "child_min";
    public const string TChildMax = "child_max";
    public const string TChildSize = "child_size";
    public const string TChildEnums = "child_enums";
    public const string TChildPattern = "child_pattern";
    public const string TChildLocation = "child_location";
    public const string TChildVersion = "child_version";
    public const string TChildMime = "child_mime";
}

public enum TagKey : byte
{
    KIsNull = 0 << 3,
    KExample = 1 << 3,
    KDeprecated = 2 << 3,

    KDesc = 3 << 3,
    KType = 4 << 3,
    KNullable = 5 << 3,
    KAllowEmpty = 6 << 3,
    KUnique = 7 << 3,
    KDefault = 8 << 3,
    KMin = 9 << 3,
    KMax = 10 << 3,
    KSize = 11 << 3,
    KEnum = 12 << 3,
    KPattern = 13 << 3,
    KLocation = 14 << 3,
    KVersion = 15 << 3,
    KMime = 16 << 3,

    KChildDesc = 17 << 3,
    KChildType = 18 << 3,
    KChildNullable = 19 << 3,
    KChildAllowEmpty = 20 << 3,
    KChildUnique = 21 << 3,
    KChildDefaultVal = 22 << 3,
    KChildMin = 23 << 3,
    KChildMax = 24 << 3,
    KChildSize = 25 << 3,
    KChildEnums = 26 << 3,
    KChildPattern = 27 << 3,
    KChildLocation = 28 << 3,
    KChildVersion = 29 << 3,
    KChildMime = 30 << 3,

    KMore = 31 << 3
}

public class Tag
{
    public string Name { get; set; } = string.Empty;

    public bool IsNull { get; set; }
    public bool Example { get; set; }

    public string Desc { get; set; } = string.Empty;
    public ValueType Type { get; set; } = ValueType.Unknown;
    public bool Deprecated { get; set; }
    public bool Nullable { get; set; }
    public bool AllowEmpty { get; set; }
    public bool Unique { get; set; }
    public string DefaultVal { get; set; } = string.Empty;
    public string Min { get; set; } = string.Empty;
    public string Max { get; set; } = string.Empty;
    public int Size { get; set; }
    public string Enums { get; set; } = string.Empty;
    public string Pattern { get; set; } = string.Empty;
    public int Location { get; set; }
    public int Version { get; set; }
    public string Mime { get; set; } = string.Empty;
    public int More { get; set; }

    public string ChildDesc { get; set; } = string.Empty;
    public ValueType ChildType { get; set; } = ValueType.Unknown;
    public bool ChildNullable { get; set; }
    public bool ChildAllowEmpty { get; set; }
    public bool ChildUnique { get; set; }
    public string ChildDefaultVal { get; set; } = string.Empty;
    public string ChildMin { get; set; } = string.Empty;
    public string ChildMax { get; set; } = string.Empty;
    public int ChildSize { get; set; }
    public string ChildEnums { get; set; } = string.Empty;
    public string ChildPattern { get; set; } = string.Empty;
    public int ChildLocation { get; set; }
    public int ChildVersion { get; set; }
    public string ChildMime { get; set; } = string.Empty;

    public bool IsInherit { get; set; }

    public const int DefaultVersionValue = 0;

    public static Tag NewTag()
    {
        return new Tag
        {
            Version = DefaultVersionValue,
            ChildVersion = DefaultVersionValue,
            Location = 0,
            ChildLocation = 0
        };
    }

    public static Tag Empty()
    {
        return new Tag();
    }

    public void InheritFromArrayParent(Tag parent)
    {
        Inherit(parent);
    }

    public void Inherit(Tag parent)
    {
        IsInherit = true;

        if (!string.IsNullOrEmpty(parent.ChildDesc))
        {
            Desc = parent.ChildDesc;
        }

        if (parent.ChildType != ValueType.Unknown)
        {
            Type = parent.ChildType;
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

        if (!string.IsNullOrEmpty(parent.ChildDefaultVal))
        {
            DefaultVal = parent.ChildDefaultVal;
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

        if (!string.IsNullOrEmpty(parent.ChildEnums))
        {
            Enums = parent.ChildEnums;
            Type = ValueType.Enums;
        }

        if (!string.IsNullOrEmpty(parent.ChildPattern))
        {
            Pattern = parent.ChildPattern;
        }

        if (parent.ChildLocation != 0)
        {
            Location = parent.ChildLocation;
        }

        if (parent.ChildVersion != DefaultVersionValue)
        {
            Version = parent.ChildVersion;
        }

        if (!string.IsNullOrEmpty(parent.ChildMime))
        {
            Mime = parent.ChildMime;
            Type = ValueType.Media;
        }
    }

    public Tag Copy()
    {
        return new Tag
        {
            Name = Name,
            IsNull = IsNull,
            Example = Example,
            Desc = Desc,
            Type = Type,
            Deprecated = Deprecated,
            Nullable = Nullable,
            AllowEmpty = AllowEmpty,
            Unique = Unique,
            DefaultVal = DefaultVal,
            Min = Min,
            Max = Max,
            Size = Size,
            Enums = Enums,
            Pattern = Pattern,
            Location = Location,
            Version = Version,
            Mime = Mime,
            More = More,
            ChildDesc = ChildDesc,
            ChildType = ChildType,
            ChildNullable = ChildNullable,
            ChildAllowEmpty = ChildAllowEmpty,
            ChildUnique = ChildUnique,
            ChildDefaultVal = ChildDefaultVal,
            ChildMin = ChildMin,
            ChildMax = ChildMax,
            ChildSize = ChildSize,
            ChildEnums = ChildEnums,
            ChildPattern = ChildPattern,
            ChildLocation = ChildLocation,
            ChildVersion = ChildVersion,
            ChildMime = ChildMime,
            IsInherit = IsInherit
        };
    }

    private static string QuoteString(string s)
    {
        var sb = new StringBuilder();
        sb.Append('"');
        foreach (char c in s)
        {
            switch (c)
            {
                case '"':
                    sb.Append("\\\"");
                    break;
                case '\\':
                    sb.Append("\\\\");
                    break;
                case '\n':
                    sb.Append("\\n");
                    break;
                case '\r':
                    sb.Append("\\r");
                    break;
                case '\t':
                    sb.Append("\\t");
                    break;
                default:
                    sb.Append(c);
                    break;
            }
        }
        sb.Append('"');
        return sb.ToString();
    }

    private static string UnquoteString(string s)
    {
        if (s.Length < 2 || s[0] != '"' || s[s.Length - 1] != '"')
        {
            return s;
        }
        var inner = s.Substring(1, s.Length - 2);
        var sb = new StringBuilder();
        for (int i = 0; i < inner.Length; i++)
        {
            if (inner[i] == '\\' && i + 1 < inner.Length)
            {
                switch (inner[i + 1])
                {
                    case '"': sb.Append('"'); i++; break;
                    case '\\': sb.Append('\\'); i++; break;
                    case 'n': sb.Append('\n'); i++; break;
                    case 'r': sb.Append('\r'); i++; break;
                    case 't': sb.Append('\t'); i++; break;
                    default: sb.Append(inner[i]); break;
                }
            }
            else
            {
                sb.Append(inner[i]);
            }
        }
        return sb.ToString();
    }

    public override string ToString()
    {
        var b = new StringBuilder();
        bool first = true;

        void Add(string s)
        {
            if (!first)
            {
                b.Append("; ");
            }
            b.Append(s);
            first = false;
        }

        if (Type != ValueType.Unknown && !IsInherit)
        {
            if (Type == ValueType.Str ||
                Type == ValueType.I ||
                Type == ValueType.F64 ||
                Type == ValueType.Bool ||
                Type == ValueType.Obj ||
                Type == ValueType.Map ||
                Type == ValueType.Vec)
            {
            }
            else
            {
                if ((Type == ValueType.Enums && !string.IsNullOrEmpty(Enums)) ||
                    (Type == ValueType.Media && !string.IsNullOrEmpty(Mime)))
                {
                }
                else
                {
                    Add($"{TagConstants.TType}={Type.ToTypeString()}");
                }
            }
        }

        if (Example)
        {
            Add(TagConstants.TExample);
        }

        if (IsNull)
        {
            Add(TagConstants.TIsNull);
        }

        if (Nullable && !IsInherit)
        {
            if (!IsNull)
            {
                Add(TagConstants.TNullable);
            }
        }

        if (!string.IsNullOrEmpty(Desc) && !IsInherit)
        {
            Add($"{TagConstants.TDesc}={QuoteString(Desc)}");
        }

        if (Deprecated && !IsInherit)
        {
            Add("deprecated");
        }

        if (AllowEmpty && !IsInherit)
        {
            Add(TagConstants.TAllowEmpty);
        }

        if (Unique && !IsInherit)
        {
            Add(TagConstants.TUnique);
        }

        if (!string.IsNullOrEmpty(DefaultVal) && !IsInherit)
        {
            Add($"{TagConstants.TDefaultVal}={DefaultVal}");
        }

        if (!string.IsNullOrEmpty(Min) && !IsInherit)
        {
            Add($"{TagConstants.TMin}={Min}");
        }

        if (!string.IsNullOrEmpty(Max) && !IsInherit)
        {
            Add($"{TagConstants.TMax}={Max}");
        }

        if (Size != 0 && !IsInherit)
        {
            Add($"{TagConstants.TSize}={Size}");
        }

        if (!string.IsNullOrEmpty(Enums) && !IsInherit)
        {
            Add($"{TagConstants.TEnum}={Enums}");
        }

        if (!string.IsNullOrEmpty(Pattern) && !IsInherit)
        {
            Add($"{TagConstants.TPattern}={Pattern}");
        }

        if (Location != 0 && !IsInherit)
        {
            Add($"{TagConstants.TLocation}={Location}");
        }

        if (Version != DefaultVersionValue && !IsInherit)
        {
            Add($"{TagConstants.TVersion}={Version}");
        }

        if (!string.IsNullOrEmpty(Mime) && !IsInherit)
        {
            Add($"{TagConstants.TMime}={Mime}");
        }

        if (!string.IsNullOrEmpty(ChildDesc))
        {
            Add($"{TagConstants.TChildDesc}={QuoteString(ChildDesc)}");
        }

        if (ChildType != ValueType.Unknown)
        {
            if (ChildType == ValueType.Str ||
                ChildType == ValueType.I ||
                ChildType == ValueType.F64 ||
                ChildType == ValueType.Bool ||
                ChildType == ValueType.Obj ||
                ChildType == ValueType.Vec)
            {
            }
            else
            {
                if ((ChildType == ValueType.Enums && !string.IsNullOrEmpty(ChildEnums)) ||
                    (ChildType == ValueType.Media && !string.IsNullOrEmpty(ChildMime)))
                {
                }
                else
                {
                    Add($"{TagConstants.TChildType}={ChildType.ToTypeString()}");
                }
            }
        }

        if (ChildNullable)
        {
            Add(TagConstants.TChildNullable);
        }

        if (ChildAllowEmpty)
        {
            Add(TagConstants.TChildAllowEmpty);
        }

        if (ChildUnique)
        {
            Add(TagConstants.TChildUnique);
        }

        if (!string.IsNullOrEmpty(ChildDefaultVal))
        {
            Add($"{TagConstants.TChildDefaultVal}={ChildDefaultVal}");
        }

        if (!string.IsNullOrEmpty(ChildMin))
        {
            Add($"{TagConstants.TChildMin}={ChildMin}");
        }

        if (!string.IsNullOrEmpty(ChildMax))
        {
            Add($"{TagConstants.TChildMax}={ChildMax}");
        }

        if (ChildSize != 0)
        {
            Add($"{TagConstants.TChildSize}={ChildSize}");
        }

        if (!string.IsNullOrEmpty(ChildEnums))
        {
            Add($"{TagConstants.TChildEnums}={ChildEnums}");
        }

        if (!string.IsNullOrEmpty(ChildPattern))
        {
            Add($"{TagConstants.TChildPattern}={ChildPattern}");
        }

        if (ChildLocation != 0)
        {
            Add($"{TagConstants.TChildLocation}={ChildLocation}");
        }

        if (ChildVersion != DefaultVersionValue)
        {
            Add($"{TagConstants.TChildVersion}={ChildVersion}");
        }

        if (!string.IsNullOrEmpty(ChildMime))
        {
            Add($"{TagConstants.TChildMime}={ChildMime}");
        }

        var result = b.ToString();
        if (string.IsNullOrEmpty(result))
        {
            return "";
        }
        return result;
    }

    public static Tag Parse(string comment)
    {
        var r = NewTag();
        if (string.IsNullOrWhiteSpace(comment))
        {
            return r;
        }

        comment = comment.Trim();
        if (comment.StartsWith("//"))
        {
            comment = comment.Substring(2).Trim();
        }

        if (!comment.StartsWith("mm:"))
        {
            return r;
        }

        comment = comment.Substring(3).Trim();
        if (string.IsNullOrEmpty(comment))
        {
            return r;
        }

        var parts = comment.Split(';');
        foreach (var p in parts)
        {
            var pair = p.Trim();
            if (string.IsNullOrEmpty(pair))
            {
                continue;
            }

            string k, v;
            var eqIdx = pair.IndexOf('=');
            if (eqIdx >= 0)
            {
                k = pair.Substring(0, eqIdx).Trim();
                v = pair.Substring(eqIdx + 1).Trim();
            }
            else
            {
                k = pair.Trim();
                v = string.Empty;
            }

            if (v.Length >= 2 && v[0] == '"' && v[v.Length - 1] == '"')
            {
                v = UnquoteString(v);
            }

            var lower = k.ToLowerInvariant();
            switch (lower)
            {
                case TagConstants.TName:
                    r.Name = v;
                    break;

                case TagConstants.TIsNull:
                    r.IsNull = true;
                    r.Nullable = true;
                    break;

                case TagConstants.TExample:
                    r.Example = true;
                    break;

                case TagConstants.TDesc:
                    r.Desc = v;
                    break;

                case TagConstants.TType:
                    r.Type = ValueTypeExtensions.ParseValueType(v);
                    break;

                case TagConstants.TDeprecated:
                    r.Deprecated = true;
                    break;

                case TagConstants.TNullable:
                    r.Nullable = true;
                    break;

                case TagConstants.TAllowEmpty:
                    r.AllowEmpty = true;
                    break;

                case TagConstants.TUnique:
                    r.Unique = true;
                    break;

                case TagConstants.TDefaultVal:
                    r.DefaultVal = v;
                    break;

                case TagConstants.TPattern:
                    r.Pattern = v;
                    break;

                case TagConstants.TMin:
                    r.Min = v;
                    break;

                case TagConstants.TMax:
                    r.Max = v;
                    break;

                case TagConstants.TSize:
                    if (ulong.TryParse(v, out var sizeVal))
                    {
                        r.Size = (int)sizeVal;
                    }
                    break;

                case TagConstants.TEnum:
                    r.Type = ValueType.Enums;
                    r.Enums = v;
                    break;

                case TagConstants.TLocation:
                    if (int.TryParse(v, out var locVal))
                    {
                        if (locVal >= -12 && locVal <= 14)
                        {
                            r.Location = locVal;
                        }
                    }
                    break;

                case TagConstants.TVersion:
                    if (int.TryParse(v, out var verVal))
                    {
                        if (verVal >= 1 && verVal <= 10)
                        {
                            r.Version = verVal;
                        }
                    }
                    break;

                case TagConstants.TMime:
                    r.Mime = v;
                    r.Type = ValueType.Media;
                    break;

                case TagConstants.TChildDesc:
                    r.ChildDesc = v;
                    break;

                case TagConstants.TChildType:
                    r.ChildType = ValueTypeExtensions.ParseValueType(v);
                    break;

                case TagConstants.TChildNullable:
                    r.ChildNullable = true;
                    break;

                case TagConstants.TChildAllowEmpty:
                    r.ChildAllowEmpty = true;
                    break;

                case TagConstants.TChildUnique:
                    r.ChildUnique = true;
                    break;

                case TagConstants.TChildDefaultVal:
                    r.ChildDefaultVal = v;
                    break;

                case TagConstants.TChildPattern:
                    r.ChildPattern = v;
                    break;

                case TagConstants.TChildMin:
                    r.ChildMin = v;
                    break;

                case TagConstants.TChildMax:
                    r.ChildMax = v;
                    break;

                case TagConstants.TChildSize:
                    if (ulong.TryParse(v, out var childSizeVal))
                    {
                        r.ChildSize = (int)childSizeVal;
                    }
                    break;

                case TagConstants.TChildEnums:
                    r.ChildEnums = v;
                    r.ChildType = ValueType.Enums;
                    break;

                case TagConstants.TChildLocation:
                    if (int.TryParse(v, out var childLocVal))
                    {
                        if (childLocVal >= -12 && childLocVal <= 14)
                        {
                            r.ChildLocation = childLocVal;
                        }
                    }
                    break;

                case TagConstants.TChildVersion:
                    if (int.TryParse(v, out var childVerVal))
                    {
                        if (childVerVal >= 1 && childVerVal <= 10)
                        {
                            r.ChildVersion = childVerVal;
                        }
                    }
                    break;

                case TagConstants.TChildMime:
                    r.ChildMime = v;
                    r.ChildType = ValueType.Media;
                    break;
            }
        }

        return r;
    }

    private const ulong Max1Byte = 0xFF;
    private const ulong Max2Byte = 0xFFFF;
    private const ulong Max3Byte = 0xFFFFFF;
    private const ulong Max4Byte = 0xFFFFFFFF;
    private const ulong Max5Byte = 0xFFFFFFFFFF;
    private const ulong Max6Byte = 0xFFFFFFFFFFFF;
    private const ulong Max7Byte = 0xFFFFFFFFFFFFFF;
    private const ulong Max8Byte = 0xFFFFFFFFFFFFFFFF;

    private static void EncodeU64(List<byte> buf, TagKey sign, ulong uv)
    {
        if (uv <= Max1Byte)
        {
            buf.Add((byte)((byte)sign | 0));
            buf.Add((byte)uv);
        }
        else if (uv <= Max2Byte)
        {
            buf.Add((byte)((byte)sign | 1));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max3Byte)
        {
            buf.Add((byte)((byte)sign | 2));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max4Byte)
        {
            buf.Add((byte)((byte)sign | 3));
            buf.Add((byte)(uv >> 24));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max5Byte)
        {
            buf.Add((byte)((byte)sign | 4));
            buf.Add((byte)(uv >> 32));
            buf.Add((byte)(uv >> 24));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max6Byte)
        {
            buf.Add((byte)((byte)sign | 5));
            buf.Add((byte)(uv >> 40));
            buf.Add((byte)(uv >> 32));
            buf.Add((byte)(uv >> 24));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max7Byte)
        {
            buf.Add((byte)((byte)sign | 6));
            buf.Add((byte)(uv >> 48));
            buf.Add((byte)(uv >> 40));
            buf.Add((byte)(uv >> 32));
            buf.Add((byte)(uv >> 24));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
        else if (uv <= Max8Byte)
        {
            buf.Add((byte)((byte)sign | 7));
            buf.Add((byte)(uv >> 56));
            buf.Add((byte)(uv >> 48));
            buf.Add((byte)(uv >> 40));
            buf.Add((byte)(uv >> 32));
            buf.Add((byte)(uv >> 24));
            buf.Add((byte)(uv >> 16));
            buf.Add((byte)(uv >> 8));
            buf.Add((byte)uv);
        }
    }

    public byte[] ToBytes()
    {
        var bs = new List<byte>();

        if (Example)
        {
            bs.Add((byte)((byte)TagKey.KExample | 1));
        }

        if (IsNull)
        {
            bs.Add((byte)((byte)TagKey.KIsNull | 1));
        }

        if (Nullable && !IsInherit)
        {
            if (!IsNull)
            {
                bs.Add((byte)((byte)TagKey.KNullable | 1));
            }
        }

        if (!string.IsNullOrEmpty(Desc) && !IsInherit)
        {
            var descBytes = Encoding.UTF8.GetBytes(Desc);
            var l = descBytes.Length;
            if (l <= 5)
            {
                bs.Add((byte)((byte)TagKey.KDesc | (byte)l));
                bs.AddRange(descBytes);
            }
            else if (l <= 1 << 8)
            {
                bs.Add((byte)((byte)TagKey.KDesc | 6));
                bs.Add((byte)l);
                bs.AddRange(descBytes);
            }
            else if (l <= 1 << 16)
            {
                bs.Add((byte)((byte)TagKey.KDesc | 7));
                bs.Add((byte)(l >> 8));
                bs.Add((byte)l);
                bs.AddRange(descBytes);
            }
        }

        if (Type != ValueType.Unknown && !IsInherit)
        {
            if (Type == ValueType.Str ||
                Type == ValueType.Bytes ||
                Type == ValueType.I ||
                Type == ValueType.F64 ||
                Type == ValueType.Bool ||
                Type == ValueType.Obj ||
                Type == ValueType.Map ||
                Type == ValueType.Vec ||
                (Type == ValueType.Arr && Size > 0) || Type == ValueType.Arr)
            {
            }
            else
            {
                if ((Type == ValueType.Enums && !string.IsNullOrEmpty(Enums)) ||
                    (Type == ValueType.Media && !string.IsNullOrEmpty(Mime)))
                {
                }
                else
                {
                    bs.Add((byte)TagKey.KType);
                    bs.Add((byte)Type);
                }
            }
        }

        if (Deprecated && !IsInherit)
        {
            bs.Add((byte)((byte)TagKey.KDeprecated | 1));
        }

        if (AllowEmpty && !IsInherit)
        {
            bs.Add((byte)((byte)TagKey.KAllowEmpty | 1));
        }

        if (Unique && !IsInherit)
        {
            bs.Add((byte)((byte)TagKey.KUnique | 1));
        }

        if (!string.IsNullOrEmpty(DefaultVal) && !IsInherit)
        {
            var defaultBytes = Encoding.UTF8.GetBytes(DefaultVal);
            var l = defaultBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KDefault | (byte)l));
                bs.AddRange(defaultBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KDefault | 7));
                bs.Add((byte)l);
                bs.AddRange(defaultBytes);
            }
        }

        if (!string.IsNullOrEmpty(Min) && !IsInherit)
        {
            var minBytes = Encoding.UTF8.GetBytes(Min);
            var l = minBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KMin | (byte)l));
                bs.AddRange(minBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KMin | 7));
                bs.Add((byte)l);
                bs.AddRange(minBytes);
            }
        }

        if (!string.IsNullOrEmpty(Max) && !IsInherit)
        {
            var maxBytes = Encoding.UTF8.GetBytes(Max);
            var l = maxBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KMax | (byte)l));
                bs.AddRange(maxBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KMax | 7));
                bs.Add((byte)l);
                bs.AddRange(maxBytes);
            }
        }

        if (Size != 0 && !IsInherit)
        {
            EncodeU64(bs, TagKey.KSize, (ulong)Size);
        }

        if (!string.IsNullOrEmpty(Enums) && !IsInherit)
        {
            var enumBytes = Encoding.UTF8.GetBytes(Enums);
            var l = enumBytes.Length;
            if (l <= 5)
            {
                bs.Add((byte)((byte)TagKey.KEnum | (byte)l));
                bs.AddRange(enumBytes);
            }
            else if (l <= 1 << 8)
            {
                bs.Add((byte)((byte)TagKey.KEnum | 6));
                bs.Add((byte)l);
                bs.AddRange(enumBytes);
            }
            else if (l <= 1 << 16)
            {
                bs.Add((byte)((byte)TagKey.KEnum | 7));
                bs.Add((byte)(l >> 8));
                bs.Add((byte)l);
                bs.AddRange(enumBytes);
            }
        }

        if (!string.IsNullOrEmpty(Pattern) && !IsInherit)
        {
            var patternBytes = Encoding.UTF8.GetBytes(Pattern);
            var l = patternBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KPattern | (byte)l));
                bs.AddRange(patternBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KPattern | 7));
                bs.Add((byte)l);
                bs.AddRange(patternBytes);
            }
        }

        if (Location != 0 && !IsInherit)
        {
            var v = Location.ToString();
            var vBytes = Encoding.UTF8.GetBytes(v);
            bs.Add((byte)((byte)TagKey.KLocation | (byte)vBytes.Length));
            bs.AddRange(vBytes);
        }

        if (Version != DefaultVersionValue && !IsInherit)
        {
            EncodeU64(bs, TagKey.KVersion, (ulong)Version);
        }

        if (!string.IsNullOrEmpty(Mime) && !IsInherit)
        {
            int mimeIndex = MimeWire.ParseMIME(Mime);
            EncodeU64(bs, TagKey.KMime, (ulong)mimeIndex);
        }

        if (!string.IsNullOrEmpty(ChildDesc))
        {
            var childDescBytes = Encoding.UTF8.GetBytes(ChildDesc);
            var l = childDescBytes.Length;
            if (l <= 5)
            {
                bs.Add((byte)((byte)TagKey.KChildDesc | (byte)l));
                bs.AddRange(childDescBytes);
            }
            else if (l <= 1 << 8)
            {
                bs.Add((byte)((byte)TagKey.KChildDesc | 6));
                bs.Add((byte)l);
                bs.AddRange(childDescBytes);
            }
            else if (l <= 1 << 16)
            {
                bs.Add((byte)((byte)TagKey.KChildDesc | 7));
                bs.Add((byte)(l >> 8));
                bs.Add((byte)l);
                bs.AddRange(childDescBytes);
            }
        }

        if (ChildType != ValueType.Unknown)
        {
            if (ChildType == ValueType.Str ||
                ChildType == ValueType.I ||
                ChildType == ValueType.F64 ||
                ChildType == ValueType.Bool ||
                ChildType == ValueType.Obj ||
                ChildType == ValueType.Vec)
            {
            }
            else
            {
                if ((ChildType == ValueType.Enums && !string.IsNullOrEmpty(ChildEnums)) ||
                    (ChildType == ValueType.Media && !string.IsNullOrEmpty(ChildMime)))
                {
                }
                else
                {
                    bs.Add((byte)TagKey.KChildType);
                    bs.Add((byte)ChildType);
                }
            }
        }

        if (ChildNullable)
        {
            bs.Add((byte)((byte)TagKey.KChildNullable | 1));
        }

        if (ChildAllowEmpty)
        {
            bs.Add((byte)((byte)TagKey.KChildAllowEmpty | 1));
        }

        if (ChildUnique)
        {
            bs.Add((byte)((byte)TagKey.KChildUnique | 1));
        }

        if (!string.IsNullOrEmpty(ChildDefaultVal))
        {
            var childDefaultBytes = Encoding.UTF8.GetBytes(ChildDefaultVal);
            var l = childDefaultBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KChildDefaultVal | (byte)l));
                bs.AddRange(childDefaultBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KChildDefaultVal | 7));
                bs.Add((byte)l);
                bs.AddRange(childDefaultBytes);
            }
        }

        if (!string.IsNullOrEmpty(ChildMin))
        {
            var childMinBytes = Encoding.UTF8.GetBytes(ChildMin);
            var l = childMinBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KChildMin | (byte)l));
                bs.AddRange(childMinBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KChildMin | 7));
                bs.Add((byte)l);
                bs.AddRange(childMinBytes);
            }
        }

        if (!string.IsNullOrEmpty(ChildMax))
        {
            var childMaxBytes = Encoding.UTF8.GetBytes(ChildMax);
            var l = childMaxBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KChildMax | (byte)l));
                bs.AddRange(childMaxBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KChildMax | 7));
                bs.Add((byte)l);
                bs.AddRange(childMaxBytes);
            }
        }

        if (ChildSize != 0)
        {
            EncodeU64(bs, TagKey.KChildSize, (ulong)ChildSize);
        }

        if (!string.IsNullOrEmpty(ChildEnums))
        {
            var childEnumBytes = Encoding.UTF8.GetBytes(ChildEnums);
            var l = childEnumBytes.Length;
            if (l <= 5)
            {
                bs.Add((byte)((byte)TagKey.KChildEnums | (byte)l));
                bs.AddRange(childEnumBytes);
            }
            else if (l <= 1 << 8)
            {
                bs.Add((byte)((byte)TagKey.KChildEnums | 6));
                bs.Add((byte)l);
                bs.AddRange(childEnumBytes);
            }
            else if (l <= 1 << 16)
            {
                bs.Add((byte)((byte)TagKey.KChildEnums | 7));
                bs.Add((byte)(l >> 8));
                bs.Add((byte)l);
                bs.AddRange(childEnumBytes);
            }
        }

        if (!string.IsNullOrEmpty(ChildPattern))
        {
            var childPatternBytes = Encoding.UTF8.GetBytes(ChildPattern);
            var l = childPatternBytes.Length;
            if (l < 7)
            {
                bs.Add((byte)((byte)TagKey.KChildPattern | (byte)l));
                bs.AddRange(childPatternBytes);
            }
            else
            {
                bs.Add((byte)((byte)TagKey.KChildPattern | 7));
                bs.Add((byte)l);
                bs.AddRange(childPatternBytes);
            }
        }

        if (ChildLocation != 0)
        {
            var v = ChildLocation.ToString();
            var vBytes = Encoding.UTF8.GetBytes(v);
            bs.Add((byte)((byte)TagKey.KChildLocation | (byte)vBytes.Length));
            bs.AddRange(vBytes);
        }

        if (ChildVersion != DefaultVersionValue)
        {
            EncodeU64(bs, TagKey.KChildVersion, (ulong)ChildVersion);
        }

        if (!string.IsNullOrEmpty(ChildMime))
        {
            EncodeU64(bs, TagKey.KChildMime, (ulong)MimeWire.ParseMIME(ChildMime));
        }

        if (More != 0)
        {
            EncodeU64(bs, TagKey.KMore, (ulong)More);
        }

        return bs.ToArray();
    }

    public static Tag FromBytes(byte[] data)
    {
        var tag = NewTag();
        int offset = 0;

        while (offset < data.Length)
        {
            byte header = data[offset];
            offset++;

            var key = (TagKey)(header & 0xF8);
            int lenInfo = header & 0x07;

            switch (key)
            {
                case TagKey.KExample:
                    tag.Example = true;
                    break;

                case TagKey.KIsNull:
                    tag.IsNull = true;
                    tag.Nullable = true;
                    break;

                case TagKey.KNullable:
                    tag.Nullable = true;
                    break;

                case TagKey.KDeprecated:
                    tag.Deprecated = true;
                    break;

                case TagKey.KAllowEmpty:
                    tag.AllowEmpty = true;
                    break;

                case TagKey.KUnique:
                    tag.Unique = true;
                    break;

                case TagKey.KChildNullable:
                    tag.ChildNullable = true;
                    break;

                case TagKey.KChildAllowEmpty:
                    tag.ChildAllowEmpty = true;
                    break;

                case TagKey.KChildUnique:
                    tag.ChildUnique = true;
                    break;

                case TagKey.KDesc:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        tag.Desc = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KType:
                    {
                        tag.Type = (ValueType)data[offset];
                        offset++;
                        break;
                    }

                case TagKey.KDefault:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.DefaultVal = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KMin:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.Min = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KMax:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.Max = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KSize:
                    {
                        int byteCount = lenInfo + 1;
                        tag.Size = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KEnum:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        tag.Enums = Encoding.UTF8.GetString(data, offset, strLen);
                        tag.Type = ValueType.Enums;
                        offset += strLen;
                        break;
                    }

                case TagKey.KPattern:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.Pattern = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KLocation:
                    {
                        int strLen = lenInfo;
                        tag.Location = int.Parse(Encoding.UTF8.GetString(data, offset, strLen));
                        offset += strLen;
                        break;
                    }

                case TagKey.KVersion:
                    {
                        int byteCount = lenInfo + 1;
                        tag.Version = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KMime:
                    {
                        int byteCount = lenInfo + 1;
                        int mimeIndex = (int)DecodeU64(data, ref offset, byteCount);
                        tag.Mime = MimeWire.MimeToString(mimeIndex);
                        tag.Type = ValueType.Media;
                        break;
                    }

                case TagKey.KChildDesc:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        tag.ChildDesc = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildType:
                    {
                        tag.ChildType = (ValueType)data[offset];
                        offset++;
                        break;
                    }

                case TagKey.KChildDefaultVal:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.ChildDefaultVal = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildMin:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.ChildMin = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildMax:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.ChildMax = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildSize:
                    {
                        int byteCount = lenInfo + 1;
                        tag.ChildSize = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KChildEnums:
                    {
                        int strLen = DecodeStringLength(data, ref offset, lenInfo);
                        tag.ChildEnums = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildPattern:
                    {
                        int strLen = DecodeStringLengthSimple(data, ref offset, lenInfo);
                        tag.ChildPattern = Encoding.UTF8.GetString(data, offset, strLen);
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildLocation:
                    {
                        int strLen = lenInfo;
                        tag.ChildLocation = int.Parse(Encoding.UTF8.GetString(data, offset, strLen));
                        offset += strLen;
                        break;
                    }

                case TagKey.KChildVersion:
                    {
                        int byteCount = lenInfo + 1;
                        tag.ChildVersion = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }

                case TagKey.KChildMime:
                    {
                        int byteCount = lenInfo + 1;
                        int childMimeIndex = (int)DecodeU64(data, ref offset, byteCount);
                        tag.ChildMime = MimeWire.MimeToString(childMimeIndex);
                        tag.ChildType = ValueType.Media;
                        break;
                    }

                case TagKey.KMore:
                    {
                        int byteCount = lenInfo + 1;
                        tag.More = (int)DecodeU64(data, ref offset, byteCount);
                        break;
                    }
            }
        }

        return tag;
    }

    private static int DecodeStringLength(byte[] data, ref int offset, int lenInfo)
    {
        if (lenInfo <= 5)
        {
            return lenInfo;
        }

        if (lenInfo == 6)
        {
            int len = data[offset];
            offset++;
            return len;
        }

        int lenHi = data[offset];
        offset++;
        int lenLo = data[offset];
        offset++;
        return (lenHi << 8) | lenLo;
    }

    private static int DecodeStringLengthSimple(byte[] data, ref int offset, int lenInfo)
    {
        if (lenInfo < 7)
        {
            return lenInfo;
        }

        int len = data[offset];
        offset++;
        return len;
    }

    private static ulong DecodeU64(byte[] data, ref int offset, int byteCount)
    {
        ulong result = 0;
        for (int i = 0; i < byteCount; i++)
        {
            result = (result << 8) | data[offset];
            offset++;
        }
        return result;
    }

    public static Tag MergeTag(Tag? dst, Tag? src)
    {
        if (src == null)
        {
            return dst ?? Tag.NewTag();
        }

        if (dst == null)
        {
            return src;
        }

        if (src.IsNull)
        {
            dst.IsNull = src.IsNull;
        }

        if (src.Example)
        {
            dst.Example = src.Example;
        }

        if (!string.IsNullOrEmpty(src.Desc))
        {
            dst.Desc = src.Desc;
        }

        if (src.Type != ValueType.Unknown)
        {
            dst.Type = src.Type;
        }

        if (src.Deprecated)
        {
            dst.Deprecated = true;
        }

        if (src.Nullable)
        {
            dst.Nullable = true;
        }

        if (src.AllowEmpty)
        {
            dst.AllowEmpty = true;
        }

        if (src.Unique)
        {
            dst.Unique = true;
        }

        if (!string.IsNullOrEmpty(src.DefaultVal))
        {
            dst.DefaultVal = src.DefaultVal;
        }

        if (!string.IsNullOrEmpty(src.Min))
        {
            dst.Min = src.Min;
        }

        if (!string.IsNullOrEmpty(src.Max))
        {
            dst.Max = src.Max;
        }

        if (src.Size != 0)
        {
            dst.Size = src.Size;
        }

        if (!string.IsNullOrEmpty(src.Enums))
        {
            dst.Enums = src.Enums;
        }

        if (!string.IsNullOrEmpty(src.Pattern))
        {
            dst.Pattern = src.Pattern;
        }

        if (src.Location != 0)
        {
            dst.Location = src.Location;
        }

        if (src.Version != DefaultVersionValue)
        {
            dst.Version = src.Version;
        }

        if (!string.IsNullOrEmpty(src.Mime))
        {
            dst.Mime = src.Mime;
        }

        if (!string.IsNullOrEmpty(src.ChildDesc))
        {
            dst.ChildDesc = src.ChildDesc;
        }

        if (src.ChildType != ValueType.Unknown)
        {
            dst.ChildType = src.ChildType;
        }

        if (src.ChildNullable)
        {
            dst.ChildNullable = true;
        }

        if (src.ChildAllowEmpty)
        {
            dst.ChildAllowEmpty = true;
        }

        if (src.ChildUnique)
        {
            dst.ChildUnique = true;
        }

        if (!string.IsNullOrEmpty(src.ChildDefaultVal))
        {
            dst.ChildDefaultVal = src.ChildDefaultVal;
        }

        if (!string.IsNullOrEmpty(src.ChildMin))
        {
            dst.ChildMin = src.ChildMin;
        }

        if (!string.IsNullOrEmpty(src.ChildMax))
        {
            dst.ChildMax = src.ChildMax;
        }

        if (src.ChildSize != 0)
        {
            dst.ChildSize = src.ChildSize;
        }

        if (!string.IsNullOrEmpty(src.ChildEnums))
        {
            dst.ChildEnums = src.ChildEnums;
        }

        if (!string.IsNullOrEmpty(src.ChildPattern))
        {
            dst.ChildPattern = src.ChildPattern;
        }

        if (src.ChildLocation != 0)
        {
            dst.ChildLocation = src.ChildLocation;
        }

        if (src.ChildVersion != DefaultVersionValue)
        {
            dst.ChildVersion = src.ChildVersion;
        }

        if (!string.IsNullOrEmpty(src.ChildMime))
        {
            dst.ChildMime = src.ChildMime;
        }

        return dst;
    }
}