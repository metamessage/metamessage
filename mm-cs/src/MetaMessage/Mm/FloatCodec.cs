namespace MetaMessage.Mm;

public static class FloatCodec
{
    public static (bool Negative, int Exponent, long Mantissa) ParseDecimalString(string s)
    {
        if (string.IsNullOrEmpty(s))
        {
            throw new ArgumentException("Empty numeric string");
        }

        bool neg = s[0] == '-';
        if (neg)
        {
            s = s.Substring(1);
            if (string.IsNullOrEmpty(s))
            {
                throw new ArgumentException("Invalid numeric string: only minus sign");
            }
        }

        string? expPart = null;
        int eIdx = s.IndexOf('e');
        if (eIdx == -1)
        {
            eIdx = s.IndexOf('E');
        }
        if (eIdx != -1)
        {
            expPart = s.Substring(eIdx + 1);
            s = s.Substring(0, eIdx);
            if (string.IsNullOrEmpty(expPart))
            {
                throw new ArgumentException("Missing exponent part in scientific notation");
            }
        }

        int dot = s.IndexOf('.');
        string intPart, fracPart;
        if (dot == -1)
        {
            intPart = string.IsNullOrEmpty(s) ? "0" : s;
            fracPart = "";
        }
        else
        {
            intPart = s.Substring(0, dot);
            fracPart = s.Substring(dot + 1);
        }

        string intPartFinal = string.IsNullOrEmpty(intPart) ? "0" : intPart;
        int baseExp = -fracPart.Length;
        if (expPart != null)
        {
            baseExp += int.Parse(expPart);
        }

        if (baseExp < -128 || baseExp > 127)
        {
            throw new ArgumentException($"Final exponent out of range: {baseExp}");
        }

        string mantissaStr = (intPartFinal + fracPart).TrimStart('0');
        if (string.IsNullOrEmpty(mantissaStr))
        {
            mantissaStr = "0";
        }

        long mantissa = long.Parse(mantissaStr);
        if (mantissaStr.Length > 19)
        {
            throw new ArgumentException($"Mantissa overflow (exceeds int64 max): {mantissaStr}");
        }

        return (neg, baseExp, mantissa);
    }

    public static string MantissaToDecimal(long mantissa, int exp)
    {
        string numStr = mantissa.ToString();
        int decimalPos = numStr.Length + exp;

        if (decimalPos <= 0)
        {
            return "0." + new string('0', -decimalPos) + numStr;
        }
        else if (decimalPos > 0 && decimalPos < numStr.Length)
        {
            return numStr.Substring(0, decimalPos) + "." + numStr.Substring(decimalPos);
        }
        else
        {
            return numStr + new string('0', Math.Max(0, decimalPos - numStr.Length));
        }
    }
}