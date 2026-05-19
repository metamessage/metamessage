namespace MetaMessage.Core;

public static class BigIntWireCodec
{
    public static byte[] EncodeSignedDecimal(string s)
    {
        if (string.IsNullOrEmpty(s))
        {
            return new byte[] { 0 };
        }

        bool negative = s[0] == '-';
        string absValue = negative ? s.Substring(1) : s;
        absValue = absValue.TrimStart('0');
        if (string.IsNullOrEmpty(absValue))
        {
            absValue = "0";
        }

        int len = absValue.Length;
        var bw = new BitWriter();

        bw.WriteBit(negative ? 1 : 0);

        for (int i = 0; i < len;)
        {
            int rem = len - i;
            if (rem >= 3)
            {
                int groupVal = int.Parse(absValue.Substring(i, 3));
                bw.WriteBits(groupVal, 10);
                i += 3;
            }
            else if (rem == 2)
            {
                int groupVal = int.Parse(absValue.Substring(i, 2));
                bw.WriteBits(groupVal, 7);
                i += 2;
            }
            else
            {
                int groupVal = int.Parse(absValue.Substring(i, 1));
                bw.WriteBits(groupVal, 4);
                i += 1;
            }
        }

        byte[] encoded = bw.ToArray();
        var result = new byte[encoded.Length + 1];
        result[0] = (byte)len;
        Array.Copy(encoded, 0, result, 1, encoded.Length);
        return result;
    }

    public static string DecodeSignedDecimal(byte[] data)
    {
        if (data == null || data.Length <= 1)
        {
            return "0";
        }

        int n = data[0];
        var payload = new byte[data.Length - 1];
        Array.Copy(data, 1, payload, 0, payload.Length);

        if (payload.Length == 0)
        {
            return "0";
        }

        var br = new BitReader(payload);
        int totalBits = payload.Length * 8;

        int sign = br.ReadBit();

        var sb = new System.Text.StringBuilder();
        int pos = 0;
        while (pos < n)
        {
            int rem = n - pos;
            if (rem >= 3)
            {
                sb.Append(br.ReadBits(10).ToString("D3"));
                pos += 3;
            }
            else if (rem >= 2)
            {
                sb.Append(br.ReadBits(7).ToString("D2"));
                pos += 2;
            }
            else
            {
                sb.Append(br.ReadBits(4).ToString());
                pos += 1;
            }
        }

        if (sb.Length == 0)
        {
            return "0";
        }

        string result = sb.ToString().TrimStart('0');
        if (string.IsNullOrEmpty(result))
        {
            result = "0";
        }

        return sign == 1 ? "-" + result : result;
    }

    private class BitWriter
    {
        private readonly List<byte> _bytes = new();
        private int _bitPos;

        public void WriteBit(int bit)
        {
            if (_bitPos == 0)
            {
                _bytes.Add(0);
            }

            int byteIndex = _bytes.Count - 1;
            if (bit != 0)
            {
                _bytes[byteIndex] |= (byte)(1 << (7 - _bitPos));
            }

            _bitPos++;
            if (_bitPos == 8)
            {
                _bitPos = 0;
            }
        }

        public void WriteBits(int value, int numBits)
        {
            for (int i = numBits - 1; i >= 0; i--)
            {
                WriteBit((value >> i) & 1);
            }
        }

        public byte[] ToArray()
        {
            return _bytes.ToArray();
        }
    }

    private class BitReader
    {
        private readonly byte[] _data;
        private int _byteIndex;
        private int _bitPos;

        public BitReader(byte[] data)
        {
            _data = data;
            _byteIndex = 0;
            _bitPos = 0;
        }

        public int ReadBit()
        {
            int bit = (_data[_byteIndex] >> (7 - _bitPos)) & 1;
            _bitPos++;
            if (_bitPos == 8)
            {
                _bitPos = 0;
                _byteIndex++;
            }
            return bit;
        }

        public int ReadBits(int numBits)
        {
            int value = 0;
            for (int i = 0; i < numBits; i++)
            {
                value = (value << 1) | ReadBit();
            }
            return value;
        }
    }
}