namespace MetaMessage.Core;

public static class SimpleValue
{
    public const int NULL_BOOL = 0;
    public const int NULL_INT = 1;
    public const int NULL_FLOAT = 2;
    public const int NULL_STRING = 3;
    public const int NULL_BYTES = 4;
    public const int FALSE = 5;
    public const int TRUE = 6;
    public const int CODE = 7;
    public const int MESSAGE = 8;
    public const int DATA = 9;
    public const int SUCCESS = 10;
    public const int ERROR = 11;
    public const int UNKNOWN = 12;
    public const int PAGE = 13;
    public const int LIMIT = 14;
    public const int OFFSET = 15;
    public const int TOTAL = 16;
    public const int ID = 17;
    public const int NAME = 18;
    public const int DESCRIPTION = 19;
    public const int TYPE = 20;
    public const int VERSION = 21;
    public const int STATUS = 22;
    public const int URL = 23;
    public const int CREATE_TIME = 24;
    public const int UPDATE_TIME = 25;
    public const int DELETE_TIME = 26;
    public const int ACCOUNT = 27;
    public const int TOKEN = 28;
    public const int EXPIRE_TIME = 29;
    public const int KEY = 30;
    public const int VAL = 31;

    private static readonly string[] _names = new string[]
    {
        "null_bool", "null_int", "null_float", "null_string", "null_bytes",
        "false", "true",
        "code", "message", "data", "success", "error", "unknown",
        "page", "limit", "offset", "total",
        "id", "name", "description", "type", "version", "status",
        "url", "create_time", "update_time", "delete_time",
        "account", "token", "expire_time",
        "key", "val"
    };

    private static readonly Dictionary<string, int> _valueMap = new(StringComparer.OrdinalIgnoreCase);

    static SimpleValue()
    {
        for (int i = 0; i < _names.Length; i++)
        {
            _valueMap[_names[i]] = i;
        }
    }

    public static string NameOf(int value)
    {
        if (value >= 0 && value < _names.Length)
        {
            return _names[value];
        }
        return "unknown";
    }

    public static int? NameToValue(string name)
    {
        if (_valueMap.TryGetValue(name, out int value))
        {
            return value;
        }
        return null;
    }
}