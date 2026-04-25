namespace MetaMessage.Mm;

public static class WireConstants
{
    public const int MAX_1 = 0xFF;
    public const int MAX_2 = 0xFFFF;
    public const int MAX_3 = 0xFFFFFF;
    public const uint MAX_4 = 0xFFFFFFFF;
    public const long MAX_5 = 0xFFFFFFFFFF;
    public const long MAX_6 = 0xFFFFFFFFFFFF;
    public const long MAX_7 = 0xFFFFFFFFFFFFFF;
    public const ulong MAX_8 = 0xFFFFFFFFFFFFFFFF;

    public const int INT_LEN_MASK = 0b11111;
    public const int INT_LEN_1 = INT_LEN_MASK - 7;
    public const int INT_LEN_2 = INT_LEN_MASK - 6;
    public const int INT_LEN_3 = INT_LEN_MASK - 5;
    public const int INT_LEN_4 = INT_LEN_MASK - 4;
    public const int INT_LEN_5 = INT_LEN_MASK - 3;
    public const int INT_LEN_6 = INT_LEN_MASK - 2;
    public const int INT_LEN_7 = INT_LEN_MASK - 1;
    public const int INT_LEN_8 = INT_LEN_MASK;

    public const int FLOAT_NEG_MASK = 0b10000;
    public const int FLOAT_LEN_MASK = 0b01111;
    public const int FLOAT_LEN_1 = FLOAT_LEN_MASK - 7;
    public const int FLOAT_LEN_2 = FLOAT_LEN_MASK - 6;
    public const int FLOAT_LEN_3 = FLOAT_LEN_MASK - 5;
    public const int FLOAT_LEN_4 = FLOAT_LEN_MASK - 4;
    public const int FLOAT_LEN_5 = FLOAT_LEN_MASK - 3;
    public const int FLOAT_LEN_6 = FLOAT_LEN_MASK - 2;
    public const int FLOAT_LEN_7 = FLOAT_LEN_MASK - 1;
    public const int FLOAT_LEN_8 = FLOAT_LEN_MASK;

    public const int STRING_LEN_MASK = 0b11111;
    public const int STRING_LEN_1 = STRING_LEN_MASK - 1;
    public const int STRING_LEN_2 = STRING_LEN_MASK;

    public const int BYTES_LEN_MASK = 0b11111;
    public const int BYTES_LEN_1 = BYTES_LEN_MASK - 1;
    public const int BYTES_LEN_2 = BYTES_LEN_MASK;

    public const int CONTAINER_MASK = 0b10000;
    public const int CONTAINER_MAP = 0b00000;
    public const int CONTAINER_ARRAY = 0b10000;
    public const int CONTAINER_LEN_MASK = 0b01111;
    public const int CONTAINER_LEN_1 = CONTAINER_LEN_MASK - 1;
    public const int CONTAINER_LEN_2 = CONTAINER_LEN_MASK;

    public const int TAG_LEN_MASK = 0b11111;
    public const int TAG_LEN_1 = TAG_LEN_MASK - 1;
    public const int TAG_LEN_2 = TAG_LEN_MASK;
}