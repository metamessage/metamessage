namespace MetaMessage.Ir;

public enum ValueType
{
    Unknown,

    Doc,
    Vec,
    Arr,
    Obj,
    Map,

    Str,
    Bytes,
    Bool,

    I,
    I8,
    I16,
    I32,
    I64,
    U,
    U8,
    U16,
    U32,
    U64,

    F32,
    F64,

    Bigint,
    Datetime,
    Date,
    Time,

    Uuid,
    Decimal,
    Ip,
    Url,
    Email,

    Enums,

    Media
}

public static class ValueTypeConstants
{
    public const string UnknownStr = "unknown";
    public const string DocStr = "doc";
    public const string ArrStr = "arr";
    public const string VecStr = "vec";
    public const string ObjStr = "obj";
    public const string MapStr = "map";
    public const string StrStr = "str";
    public const string BytesStr = "bytes";
    public const string BoolStr = "bool";
    public const string IStr = "i";
    public const string I8Str = "i8";
    public const string I16Str = "i16";
    public const string I32Str = "i32";
    public const string I64Str = "i64";
    public const string UStr = "u";
    public const string U8Str = "u8";
    public const string U16Str = "u16";
    public const string U32Str = "u32";
    public const string U64Str = "u64";
    public const string F32Str = "f32";
    public const string F64Str = "f64";
    public const string BigintStr = "bigint";
    public const string DatetimeStr = "datetime";
    public const string DateStr = "date";
    public const string TimeStr = "time";
    public const string UuidStr = "uuid";
    public const string DecimalStr = "decimal";
    public const string IpStr = "ip";
    public const string UrlStr = "url";
    public const string EmailStr = "email";
    public const string EnumStr = "enums";
    public const string MediaStr = "media";
}

public static class ValueTypeExtensions
{
    private static readonly Dictionary<string, ValueType> StringToValueType = new()
    {
        { ValueTypeConstants.UnknownStr, ValueType.Unknown },
        { ValueTypeConstants.DocStr, ValueType.Doc },
        { ValueTypeConstants.ArrStr, ValueType.Arr },
        { ValueTypeConstants.VecStr, ValueType.Vec },
        { ValueTypeConstants.ObjStr, ValueType.Obj },
        { ValueTypeConstants.MapStr, ValueType.Map },
        { ValueTypeConstants.StrStr, ValueType.Str },
        { ValueTypeConstants.BytesStr, ValueType.Bytes },
        { ValueTypeConstants.BoolStr, ValueType.Bool },
        { ValueTypeConstants.IStr, ValueType.I },
        { ValueTypeConstants.I8Str, ValueType.I8 },
        { ValueTypeConstants.I16Str, ValueType.I16 },
        { ValueTypeConstants.I32Str, ValueType.I32 },
        { ValueTypeConstants.I64Str, ValueType.I64 },
        { ValueTypeConstants.UStr, ValueType.U },
        { ValueTypeConstants.U8Str, ValueType.U8 },
        { ValueTypeConstants.U16Str, ValueType.U16 },
        { ValueTypeConstants.U32Str, ValueType.U32 },
        { ValueTypeConstants.U64Str, ValueType.U64 },
        { ValueTypeConstants.F32Str, ValueType.F32 },
        { ValueTypeConstants.F64Str, ValueType.F64 },
        { ValueTypeConstants.BigintStr, ValueType.Bigint },
        { ValueTypeConstants.DatetimeStr, ValueType.Datetime },
        { ValueTypeConstants.DateStr, ValueType.Date },
        { ValueTypeConstants.TimeStr, ValueType.Time },
        { ValueTypeConstants.UuidStr, ValueType.Uuid },
        { ValueTypeConstants.DecimalStr, ValueType.Decimal },
        { ValueTypeConstants.IpStr, ValueType.Ip },
        { ValueTypeConstants.UrlStr, ValueType.Url },
        { ValueTypeConstants.EmailStr, ValueType.Email },
        { ValueTypeConstants.EnumStr, ValueType.Enums },
        { ValueTypeConstants.MediaStr, ValueType.Media }
    };

    public static string ToTypeString(this ValueType vt)
    {
        return vt switch
        {
            ValueType.Unknown => ValueTypeConstants.UnknownStr,
            ValueType.Doc => ValueTypeConstants.DocStr,
            ValueType.Arr => ValueTypeConstants.ArrStr,
            ValueType.Vec => ValueTypeConstants.VecStr,
            ValueType.Obj => ValueTypeConstants.ObjStr,
            ValueType.Map => ValueTypeConstants.MapStr,
            ValueType.Str => ValueTypeConstants.StrStr,
            ValueType.Bytes => ValueTypeConstants.BytesStr,
            ValueType.Bool => ValueTypeConstants.BoolStr,
            ValueType.I => ValueTypeConstants.IStr,
            ValueType.I8 => ValueTypeConstants.I8Str,
            ValueType.I16 => ValueTypeConstants.I16Str,
            ValueType.I32 => ValueTypeConstants.I32Str,
            ValueType.I64 => ValueTypeConstants.I64Str,
            ValueType.U => ValueTypeConstants.UStr,
            ValueType.U8 => ValueTypeConstants.U8Str,
            ValueType.U16 => ValueTypeConstants.U16Str,
            ValueType.U32 => ValueTypeConstants.U32Str,
            ValueType.U64 => ValueTypeConstants.U64Str,
            ValueType.F32 => ValueTypeConstants.F32Str,
            ValueType.F64 => ValueTypeConstants.F64Str,
            ValueType.Bigint => ValueTypeConstants.BigintStr,
            ValueType.Datetime => ValueTypeConstants.DatetimeStr,
            ValueType.Date => ValueTypeConstants.DateStr,
            ValueType.Time => ValueTypeConstants.TimeStr,
            ValueType.Uuid => ValueTypeConstants.UuidStr,
            ValueType.Decimal => ValueTypeConstants.DecimalStr,
            ValueType.Ip => ValueTypeConstants.IpStr,
            ValueType.Url => ValueTypeConstants.UrlStr,
            ValueType.Email => ValueTypeConstants.EmailStr,
            ValueType.Enums => ValueTypeConstants.EnumStr,
            ValueType.Media => ValueTypeConstants.MediaStr,
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
            ValueType.Str => true,
            ValueType.Bytes => true,
            ValueType.Datetime => true,
            ValueType.Date => true,
            ValueType.Time => true,
            ValueType.Uuid => true,
            ValueType.Ip => true,
            ValueType.Url => true,
            ValueType.Email => true,
            ValueType.Enums => true,
            _ => false
        };
    }
}