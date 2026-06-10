namespace MetaMessage.Core;

public static class SimpleValue
{
    public const int NULL = 0;
    public const int NULL_BOOL = 1;
    public const int NULL_INT = 2;
    public const int NULL_FLOAT = 3;
    public const int NULL_STRING = 4;
    public const int NULL_BYTES = 5;
    public const int FALSE = 6;
    public const int TRUE = 7;
    public const int CODE = 8;
    public const int MESSAGE = 9;
    public const int DATA = 10;
    public const int SUCCESS = 11;
    public const int ERROR = 12;
    public const int UNKNOWN = 13;
    public const int PAGE = 14;
    public const int LIMIT = 15;
    public const int OFFSET = 16;
    public const int TOTAL = 17;
    public const int ID = 18;
    public const int NAME = 19;
    public const int DESCRIPTION = 20;
    public const int TYPE = 21;
    public const int VERSION = 22;
    public const int STATUS = 23;
    public const int URL = 24;
    public const int CREATE_TIME = 25;
    public const int UPDATE_TIME = 26;
    public const int DELETE_TIME = 27;
    public const int ACCOUNT = 28;
    public const int TOKEN = 29;
    public const int EXPIRE_TIME = 30;
    public const int KEY = 31;

    private static readonly string[] _names = new string[]
    {
        "null", "null_bool", "null_int", "null_float", "null_string", "null_bytes",
        "false", "true",
        "code", "message", "data", "success", "error", "unknown",
        "page", "limit", "offset", "total",
        "id", "name", "description", "type", "version", "status",
        "url", "create_time", "update_time", "delete_time",
        "account", "token", "expire_time",
        "key"
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