namespace MetaMessage.Mm;

public class WireEncoder
{
    private GrowableByteBuf _buf;

    public WireEncoder()
    {
        _buf = new GrowableByteBuf();
    }

    public void Reset()
    {
        _buf.Reset();
    }

    public byte[] ToByteArray()
    {
        return _buf.ToArray();
    }

    public int EncodeSimple(int value)
    {
        int start = _buf.Length;
        _buf.Write(Prefix.SIMPLE | value);
        return _buf.Length - start;
    }

    public int EncodeBool(bool value)
    {
        return EncodeSimple(value ? SimpleValue.TRUE : SimpleValue.FALSE);
    }

    public int EncodeInt64(long value)
    {
        int start = _buf.Length;
        if (value >= 0)
        {
            if (value <= WireConstants.INT_LEN_8 - 1)
            {
                _buf.Write(Prefix.POSITIVE_INT | (int)value);
            }
            else if (value <= WireConstants.MAX_1)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_1, (byte)value);
            }
            else if (value <= WireConstants.MAX_2)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_2, (byte)(value >> 8), (byte)value);
            }
            else if (value <= WireConstants.MAX_3)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_3, (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
            else if (value <= WireConstants.MAX_4)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_4, (byte)(value >> 24), (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
            else if (value <= WireConstants.MAX_5)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_5, (byte)(value >> 32), (byte)(value >> 24), (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
            else if (value <= WireConstants.MAX_6)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_6, (byte)(value >> 40), (byte)(value >> 32), (byte)(value >> 24), (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
            else if (value <= WireConstants.MAX_7)
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_7, (byte)(value >> 48), (byte)(value >> 40), (byte)(value >> 32), (byte)(value >> 24), (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
            else
            {
                _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_8, (byte)(value >> 56), (byte)(value >> 48), (byte)(value >> 40), (byte)(value >> 32), (byte)(value >> 24), (byte)(value >> 16), (byte)(value >> 8), (byte)value);
            }
        }
        else
        {
            long absValue = -value;
            if (absValue <= WireConstants.INT_LEN_8 - 1)
            {
                _buf.Write(Prefix.NEGATIVE_INT | (int)absValue);
            }
            else if (absValue <= WireConstants.MAX_1)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_1, (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_2)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_2, (byte)(absValue >> 8), (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_3)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_3, (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_4)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_4, (byte)(absValue >> 24), (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_5)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_5, (byte)(absValue >> 32), (byte)(absValue >> 24), (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_6)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_6, (byte)(absValue >> 40), (byte)(absValue >> 32), (byte)(absValue >> 24), (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
            else if (absValue <= WireConstants.MAX_7)
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_7, (byte)(absValue >> 48), (byte)(absValue >> 40), (byte)(absValue >> 32), (byte)(absValue >> 24), (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
            else
            {
                _buf.Write(Prefix.NEGATIVE_INT | WireConstants.INT_LEN_8, (byte)(absValue >> 56), (byte)(absValue >> 48), (byte)(absValue >> 40), (byte)(absValue >> 32), (byte)(absValue >> 24), (byte)(absValue >> 16), (byte)(absValue >> 8), (byte)absValue);
            }
        }
        return _buf.Length - start;
    }

    public int EncodeFloatString(string s)
    {
        var (negative, exponent, mantissa) = FloatCodec.ParseDecimalString(s);
        int start = _buf.Length;
        int sign = Prefix.FLOAT;
        if (negative)
        {
            sign |= WireConstants.FLOAT_NEG_MASK;
        }
        if (exponent == -1 && mantissa <= 7)
        {
            _buf.Write(sign | (int)mantissa);
        }
        else
        {
            if (mantissa <= WireConstants.MAX_1)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_1, (byte)exponent, (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_2)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_2, (byte)exponent, (byte)(mantissa >> 8), (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_3)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_3, (byte)exponent, (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_4)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_4, (byte)exponent, (byte)(mantissa >> 24), (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_5)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_5, (byte)exponent, (byte)(mantissa >> 32), (byte)(mantissa >> 24), (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_6)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_6, (byte)exponent, (byte)(mantissa >> 40), (byte)(mantissa >> 32), (byte)(mantissa >> 24), (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
            else if (mantissa <= WireConstants.MAX_7)
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_7, (byte)exponent, (byte)(mantissa >> 48), (byte)(mantissa >> 40), (byte)(mantissa >> 32), (byte)(mantissa >> 24), (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
            else
            {
                _buf.Write(sign | WireConstants.FLOAT_LEN_8, (byte)exponent, (byte)(mantissa >> 56), (byte)(mantissa >> 48), (byte)(mantissa >> 40), (byte)(mantissa >> 32), (byte)(mantissa >> 24), (byte)(mantissa >> 16), (byte)(mantissa >> 8), (byte)mantissa);
            }
        }
        return _buf.Length - start;
    }

    public int EncodeString(string s)
    {
        byte[] utf = System.Text.Encoding.UTF8.GetBytes(s);
        int length = utf.Length;
        int start = _buf.Length;
        int sign = Prefix.STRING;
        if (length < WireConstants.STRING_LEN_1)
        {
            _buf.WriteWithBytes(sign | length, utf);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteWithBytes(sign | WireConstants.STRING_LEN_1, (byte)length, utf);
        }
        else if (length < WireConstants.MAX_2)
        {
            _buf.WriteWithBytes(sign | WireConstants.STRING_LEN_2, (byte)(length >> 8), (byte)length, utf);
        }
        else
        {
            throw new Exception("String too long");
        }
        return _buf.Length - start;
    }

    public int EncodeBytes(byte[] bytes)
    {
        int length = bytes.Length;
        int start = _buf.Length;
        int sign = Prefix.BYTES;
        if (length < WireConstants.BYTES_LEN_1)
        {
            _buf.WriteWithBytes(sign | length, bytes);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteWithBytes(sign | WireConstants.BYTES_LEN_1, (byte)length, bytes);
        }
        else if (length < WireConstants.MAX_2)
        {
            _buf.WriteWithBytes(sign | WireConstants.BYTES_LEN_2, (byte)(length >> 8), (byte)length, bytes);
        }
        else
        {
            throw new Exception("Bytes too long");
        }
        return _buf.Length - start;
    }

    public int EncodeArrayPayload(byte[] payload)
    {
        int length = payload.Length;
        int start = _buf.Length;
        int sign = Prefix.CONTAINER | WireConstants.CONTAINER_ARRAY;
        if (length < WireConstants.CONTAINER_LEN_1)
        {
            _buf.WriteWithBytes(sign | length, payload);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteWithBytes(sign | WireConstants.CONTAINER_LEN_1, (byte)length, payload);
        }
        else if (length < WireConstants.MAX_2)
        {
            _buf.WriteWithBytes(sign | WireConstants.CONTAINER_LEN_2, (byte)(length >> 8), (byte)length, payload);
        }
        else
        {
            throw new Exception("Array payload too long");
        }
        return _buf.Length - start;
    }

    public int EncodeObjectPayload(byte[] payload)
    {
        int length = payload.Length;
        int start = _buf.Length;
        int sign = Prefix.CONTAINER | WireConstants.CONTAINER_MAP;
        if (length < WireConstants.CONTAINER_LEN_1)
        {
            _buf.WriteWithBytes(sign | length, payload);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteWithBytes(sign | WireConstants.CONTAINER_LEN_1, (byte)length, payload);
        }
        else if (length < WireConstants.MAX_2)
        {
            _buf.WriteWithBytes(sign | WireConstants.CONTAINER_LEN_2, (byte)(length >> 8), (byte)length, payload);
        }
        else
        {
            throw new Exception("Map payload too long");
        }
        return _buf.Length - start;
    }

    public int EncodeTagInner(byte[] tagBytes)
    {
        if (tagBytes.Length == 0)
        {
            return 0;
        }
        if (tagBytes.Length > WireConstants.MAX_2)
        {
            throw new Exception("Tag too long");
        }
        int start = _buf.Length;
        int length = tagBytes.Length;
        if (length < 254)
        {
            _buf.WriteWithBytes(length, tagBytes);
        }
        else if (length < 257)
        {
            _buf.WriteWithBytes(254, (byte)length, tagBytes);
        }
        else
        {
            _buf.WriteWithBytes(255, (byte)(length >> 8), (byte)length, tagBytes);
        }
        return _buf.Length - start;
    }

    public int EncodeTaggedPayload(byte[] payload, byte[] rawTagFields)
    {
        if (rawTagFields.Length == 0)
        {
            _buf.WriteAll(payload);
            return payload.Length;
        }
        var tEnc = new WireEncoder();
        tEnc.EncodeTagInner(rawTagFields);
        byte[] tagEncoded = tEnc.ToByteArray();
        int length = tagEncoded.Length + payload.Length;
        if (length > WireConstants.MAX_2)
        {
            throw new Exception("Tag+payload too long");
        }
        int start = _buf.Length;
        int sign = Prefix.TAG;
        if (length < WireConstants.TAG_LEN_1)
        {
            _buf.WriteWithMultipleBytes(sign | length, tagEncoded, payload);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteWithMultipleBytes(sign | WireConstants.TAG_LEN_1, (byte)length, tagEncoded, payload);
        }
        else
        {
            _buf.WriteWithMultipleBytes(sign | WireConstants.TAG_LEN_2, (byte)(length >> 8), (byte)length, tagEncoded, payload);
        }
        return _buf.Length - start;
    }

    public int EncodeBigIntDecimal(string s)
    {
        byte[] bits = BigIntWireCodec.EncodeSignedDecimal(s);
        byte[] inner = new byte[bits.Length + 1];
        inner[0] = (byte)s.Length;
        Array.Copy(bits, 0, inner, 1, bits.Length);
        return EncodeBytes(inner);
    }
}