namespace MetaMessage.Jsonc;

public enum ValueType
{
    Unknown,

    Doc,
    Slice,
    Array,
    Struct,
    Map,

    String,
    Bytes,
    Bool,

    Int,
    Int8,
    Int16,
    Int32,
    Int64,
    Uint,
    Uint8,
    Uint16,
    Uint32,
    Uint64,

    Float32,
    Float64,

    BigInt,
    DateTime,
    Date,
    Time,

    UUID,
    Decimal,
    IP,
    URL,
    Email,

    Enum,

    Image,
    Video
}

public static class ValueTypeConstants
{
    public const string UnknownStr = "unknown";
    public const string DocStr = "doc";
    public const string ArrayStr = "arr";
    public const string SliceStr = "slice";
    public const string ObjectStr = "obj";
    public const string MapStr = "map";
    public const string StringStr = "str";
    public const string BytesStr = "bytes";
    public const string BoolStr = "bool";
    public const string IntStr = "i";
    public const string Int8Str = "i8";
    public const string Int16Str = "i16";
    public const string Int32Str = "i32";
    public const string Int64Str = "i64";
    public const string UintStr = "u";
    public const string Uint8Str = "u8";
    public const string Uint16Str = "u16";
    public const string Uint32Str = "u32";
    public const string Uint64Str = "u64";
    public const string Float32Str = "f32";
    public const string Float64Str = "f64";
    public const string BigIntStr = "bi";
    public const string DateTimeStr = "datetime";
    public const string DateStr = "date";
    public const string TimeStr = "time";
    public const string UUIDStr = "uuid";
    public const string DecimalStr = "decimal";
    public const string IPStr = "ip";
    public const string URLStr = "url";
    public const string EmailStr = "email";
    public const string EnumStr = "enum";
    public const string ImageStr = "image";
    public const string VideoStr = "video";
}

public static class ValueTypeExtensions
{
    private static readonly Dictionary<string, ValueType> StringToValueType = new()
    {
        { ValueTypeConstants.UnknownStr, ValueType.Unknown },
        { ValueTypeConstants.DocStr, ValueType.Doc },
        { ValueTypeConstants.ArrayStr, ValueType.Array },
        { ValueTypeConstants.SliceStr, ValueType.Slice },
        { ValueTypeConstants.ObjectStr, ValueType.Struct },
        { ValueTypeConstants.MapStr, ValueType.Map },
        { ValueTypeConstants.StringStr, ValueType.String },
        { ValueTypeConstants.BytesStr, ValueType.Bytes },
        { ValueTypeConstants.BoolStr, ValueType.Bool },
        { ValueTypeConstants.IntStr, ValueType.Int },
        { ValueTypeConstants.Int8Str, ValueType.Int8 },
        { ValueTypeConstants.Int16Str, ValueType.Int16 },
        { ValueTypeConstants.Int32Str, ValueType.Int32 },
        { ValueTypeConstants.Int64Str, ValueType.Int64 },
        { ValueTypeConstants.UintStr, ValueType.Uint },
        { ValueTypeConstants.Uint8Str, ValueType.Uint8 },
        { ValueTypeConstants.Uint16Str, ValueType.Uint16 },
        { ValueTypeConstants.Uint32Str, ValueType.Uint32 },
        { ValueTypeConstants.Uint64Str, ValueType.Uint64 },
        { ValueTypeConstants.Float32Str, ValueType.Float32 },
        { ValueTypeConstants.Float64Str, ValueType.Float64 },
        { ValueTypeConstants.BigIntStr, ValueType.BigInt },
        { ValueTypeConstants.DateTimeStr, ValueType.DateTime },
        { ValueTypeConstants.DateStr, ValueType.Date },
        { ValueTypeConstants.TimeStr, ValueType.Time },
        { ValueTypeConstants.UUIDStr, ValueType.UUID },
        { ValueTypeConstants.DecimalStr, ValueType.Decimal },
        { ValueTypeConstants.IPStr, ValueType.IP },
        { ValueTypeConstants.URLStr, ValueType.URL },
        { ValueTypeConstants.EmailStr, ValueType.Email },
        { ValueTypeConstants.EnumStr, ValueType.Enum },
        { ValueTypeConstants.ImageStr, ValueType.Image },
        { ValueTypeConstants.VideoStr, ValueType.Video }
    };

    public static string ToTypeString(this ValueType vt)
    {
        return vt switch
        {
            ValueType.Unknown => ValueTypeConstants.UnknownStr,
            ValueType.Doc => ValueTypeConstants.DocStr,
            ValueType.Array => ValueTypeConstants.ArrayStr,
            ValueType.Slice => ValueTypeConstants.SliceStr,
            ValueType.Struct => ValueTypeConstants.ObjectStr,
            ValueType.Map => ValueTypeConstants.MapStr,
            ValueType.String => ValueTypeConstants.StringStr,
            ValueType.Bytes => ValueTypeConstants.BytesStr,
            ValueType.Bool => ValueTypeConstants.BoolStr,
            ValueType.Int => ValueTypeConstants.IntStr,
            ValueType.Int8 => ValueTypeConstants.Int8Str,
            ValueType.Int16 => ValueTypeConstants.Int16Str,
            ValueType.Int32 => ValueTypeConstants.Int32Str,
            ValueType.Int64 => ValueTypeConstants.Int64Str,
            ValueType.Uint => ValueTypeConstants.UintStr,
            ValueType.Uint8 => ValueTypeConstants.Uint8Str,
            ValueType.Uint16 => ValueTypeConstants.Uint16Str,
            ValueType.Uint32 => ValueTypeConstants.Uint32Str,
            ValueType.Uint64 => ValueTypeConstants.Uint64Str,
            ValueType.Float32 => ValueTypeConstants.Float32Str,
            ValueType.Float64 => ValueTypeConstants.Float64Str,
            ValueType.BigInt => ValueTypeConstants.BigIntStr,
            ValueType.DateTime => ValueTypeConstants.DateTimeStr,
            ValueType.Date => ValueTypeConstants.DateStr,
            ValueType.Time => ValueTypeConstants.TimeStr,
            ValueType.UUID => ValueTypeConstants.UUIDStr,
            ValueType.Decimal => ValueTypeConstants.DecimalStr,
            ValueType.IP => ValueTypeConstants.IPStr,
            ValueType.URL => ValueTypeConstants.URLStr,
            ValueType.Email => ValueTypeConstants.EmailStr,
            ValueType.Enum => ValueTypeConstants.EnumStr,
            ValueType.Image => ValueTypeConstants.ImageStr,
            ValueType.Video => ValueTypeConstants.VideoStr,
            _ => ValueTypeConstants.UnknownStr
        };
    }

    public static ValueType ParseValueType(string s)
    {
        s = s.ToLower();
        if (StringToValueType.TryGetValue(s, out var vt))
        {
            return vt;
        }
        return ValueType.Unknown;
    }

    public static bool NeedsQuotes(this ValueType vt)
    {
        return vt switch
        {
            ValueType.String => true,
            ValueType.Bytes => true,
            ValueType.DateTime => true,
            ValueType.Date => true,
            ValueType.Time => true,
            ValueType.UUID => true,
            ValueType.IP => true,
            ValueType.URL => true,
            ValueType.Email => true,
            ValueType.Enum => true,
            _ => false
        };
    }
}