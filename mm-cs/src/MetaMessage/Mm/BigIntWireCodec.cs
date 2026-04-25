namespace MetaMessage.Mm;

public static class BigIntWireCodec
{
    public static byte[] EncodeSignedDecimal(string s)
    {
        if (string.IsNullOrEmpty(s))
        {
            return Array.Empty<byte>();
        }

        bool negative = s[0] == '-';
        string absValue = negative ? s.Substring(1) : s;
        absValue = absValue.TrimStart('0');
        if (string.IsNullOrEmpty(absValue))
        {
            absValue = "0";
        }

        int length = absValue.Length;
        byte[] result = new byte[length + 1];
        result[0] = (byte)(negative ? 1 : 0);

        for (int i = 0; i < length; i++)
        {
            result[i + 1] = (byte)(absValue[i] - '0');
        }

        return result;
    }

    public static string DecodeSignedDecimal(byte[] bytes)
    {
        if (bytes == null || bytes.Length == 0)
        {
            return "0";
        }

        bool negative = bytes[0] == 1;
        int length = bytes.Length - 1;
        if (length == 0)
        {
            return "0";
        }

        char[] chars = new char[length];
        for (int i = 0; i < length; i++)
        {
            chars[i] = (char)('0' + bytes[i + 1]);
        }

        string value = new string(chars).TrimStart('0');
        if (string.IsNullOrEmpty(value))
        {
            value = "0";
        }

        return negative ? "-" + value : value;
    }
}