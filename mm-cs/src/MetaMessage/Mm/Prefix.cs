namespace MetaMessage.Mm;

public static class Prefix
{
    public const int SIMPLE = 0b000 << 5;
    public const int POSITIVE_INT = 0b001 << 5;
    public const int NEGATIVE_INT = 0b010 << 5;
    public const int FLOAT = 0b011 << 5;
    public const int STRING = 0b100 << 5;
    public const int BYTES = 0b101 << 5;
    public const int CONTAINER = 0b110 << 5;
    public const int TAG = 0b111 << 5;

    public const int PREFIX_MASK = 0b11100000;
    public const int SUFFIX_MASK = 0b00011111;

    public static int Of(int b)
    {
        return b & PREFIX_MASK;
    }

    public static int Suffix(int b)
    {
        return b & SUFFIX_MASK;
    }
}