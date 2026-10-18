namespace MetaMessage.Core;

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

    // v0.5.0 尾定界：一趟正向直接按 [body BE][sign] 写出，档位/内联语义与 v0.3.x 一致，
    // 仅物理落盘顺序由头部形态改为尾部形态（sign 落在尾字节，长度天然已知，零回写）。
    public int EncodeInt64(long value)
    {
        int start = _buf.Length;
        if (value >= 0)
        {
            WriteIntTail((ulong)value, Prefix.POSITIVE_INT);
        }
        else
        {
            ulong uv = value == long.MinValue ? 9223372036854775808UL : (ulong)(-value);
            WriteIntTail(uv, Prefix.NEGATIVE_INT);
        }
        return _buf.Length - start;
    }

    private void WriteIntTail(ulong uv, int prefix)
    {
        if (uv <= (ulong)(WireConstants.INT_LEN_1 - 1))
        {
            _buf.Write(prefix | (int)uv);
            return;
        }
        int width;
        if (uv <= (ulong)WireConstants.MAX_1) width = 1;
        else if (uv <= (ulong)WireConstants.MAX_2) width = 2;
        else if (uv <= (ulong)WireConstants.MAX_3) width = 3;
        else if (uv <= (ulong)WireConstants.MAX_4) width = 4;
        else if (uv <= (ulong)WireConstants.MAX_5) width = 5;
        else if (uv <= (ulong)WireConstants.MAX_6) width = 6;
        else if (uv <= (ulong)WireConstants.MAX_7) width = 7;
        else width = 8;
        // body 大端先写，sign(档位) 落尾
        for (int i = width - 1; i >= 0; i--)
            _buf.Write((byte)(uv >> (8 * i)));
        _buf.Write(prefix | (WireConstants.INT_LEN_1 - 1 + width));
    }

    public int EncodeUInt64(ulong value)
    {
        int start = _buf.Length;
        for (int i = 7; i >= 0; i--)
            _buf.Write((byte)(value >> (8 * i)));
        _buf.Write(Prefix.POSITIVE_INT | WireConstants.INT_LEN_8);
        return _buf.Length - start;
    }

    public int EncodeInt8(sbyte value) => EncodeInt64(value);

    public int EncodeInt16(short value) => EncodeInt64(value);

    public int EncodeInt32(int value) => EncodeInt64(value);

    public int EncodeUint8(byte value) => EncodeInt64(value);

    public int EncodeUint16(ushort value) => EncodeInt64(value);

    public int EncodeUint32(uint value) => EncodeInt64(value);

    public int EncodeUint64(ulong value) => EncodeUInt64(value);

    public int EncodeFloat(float value) => EncodeFloatString(value.ToString("G"));

    public int EncodeDouble(double value) => EncodeFloatString(value.ToString("G"));

    public int EncodeDatetime(System.DateTime value) => EncodeInt64(new System.DateTimeOffset(value).ToUnixTimeSeconds());

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
            int width;
            if (mantissa <= WireConstants.MAX_1) width = 1;
            else if (mantissa <= WireConstants.MAX_2) width = 2;
            else if (mantissa <= WireConstants.MAX_3) width = 3;
            else if (mantissa <= WireConstants.MAX_4) width = 4;
            else if (mantissa <= WireConstants.MAX_5) width = 5;
            else if (mantissa <= WireConstants.MAX_6) width = 6;
            else if (mantissa <= WireConstants.MAX_7) width = 7;
            else width = 8;
            // 尾定界：body [exp][mantissa BE w字节]，sign(档位) 落尾
            _buf.Write((byte)exponent);
            for (int i = width - 1; i >= 0; i--)
                _buf.Write((byte)(mantissa >> (8 * i)));
            _buf.Write(sign | (WireConstants.FLOAT_LEN_1 - 1 + width));
        }
        return _buf.Length - start;
    }

    public int EncodeString(string s)
    {
        byte[] utf = System.Text.Encoding.UTF8.GetBytes(s);
        int length = utf.Length;
        int start = _buf.Length;
        int sign = Prefix.STRING;
        // 尾定界：[body][lenBE 档位][sign]
        if (length < WireConstants.STRING_LEN_1)
        {
            _buf.WriteAll(utf);
            _buf.Write(sign | length);
        }
        else if (length < WireConstants.MAX_1)
        {
            _buf.WriteAll(utf);
            _buf.Write((byte)length);
            _buf.Write(sign | WireConstants.STRING_LEN_1);
        }
        // The last bucket must use <=: MAX_2 (65535) itself is representable by the
        // 2-byte length field. Using < makes length == MAX_2 fall through every
        // branch and get reported as "too long" even though it fits on the wire.
        else if (length <= WireConstants.MAX_2)
        {
            _buf.WriteAll(utf);
            _buf.Write((byte)(length >> 8));
            _buf.Write((byte)length);
            _buf.Write(sign | WireConstants.STRING_LEN_2);
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
        // 尾定界：[body][lenBE 档位][sign]
        if (length < WireConstants.BYTES_LEN_1)
        {
            _buf.WriteAll(bytes);
            _buf.Write(sign | length);
        }
        else if (length <= WireConstants.MAX_4)
        {
            int width;
            // v0.3.0：1 字节档边界对齐 Go 参考（length == 255 落 1 字节档）
            if (length <= WireConstants.MAX_1) width = 1;
            else if (length <= WireConstants.MAX_2) width = 2;
            else if (length <= WireConstants.MAX_3) width = 3;
            else width = 4;
            WriteBodyTierSign(bytes, length, width, sign | (WireConstants.BYTES_LEN_1 - 1 + width));
        }
        // 最后一档必须 <=：length == MAX_4 本身可表达，否则静默丢数据
        else
        {
            throw new Exception("Bytes too long");
        }
        return _buf.Length - start;
    }

    // 尾定界通用写出：body 先，再大端 length(width 字节)，最后 sign(含档位)。
    private void WriteBodyTierSign(byte[] body, long length, int width, int signWithTier)
    {
        _buf.WriteAll(body);
        for (int i = width - 1; i >= 0; i--)
            _buf.Write((byte)(length >> (8 * i)));
        _buf.Write(signWithTier);
    }

    public int EncodeArrayPayload(byte[] payload)
    {
        return EncodeContainerTail(payload, Prefix.CONTAINER | WireConstants.CONTAINER_ARRAY, "Array payload too long");
    }

    public int EncodeObjectPayload(byte[] payload)
    {
        return EncodeContainerTail(payload, Prefix.CONTAINER | WireConstants.CONTAINER_MAP, "Map payload too long");
    }

    // 尾定界容器：[payload][lenBE 档位][sign]（sign bit4=Array/Object）。
    private int EncodeContainerTail(byte[] payload, int baseSign, string errName)
    {
        int length = payload.Length;
        int start = _buf.Length;
        if (length < WireConstants.CONTAINER_LEN_1)
        {
            _buf.WriteAll(payload);
            _buf.Write(baseSign | length);
        }
        else if (length <= WireConstants.MAX_4)
        {
            int width;
            if (length <= WireConstants.MAX_1) width = 1;
            else if (length <= WireConstants.MAX_2) width = 2;
            else if (length <= WireConstants.MAX_3) width = 3;
            else width = 4;
            WriteBodyTierSign(payload, length, width, baseSign | (WireConstants.CONTAINER_LEN_1 - 1 + width));
        }
        else
        {
            throw new Exception(errName);
        }
        return _buf.Length - start;
    }

    // v0.5.0 尾定界 tag 包装段（一趟正向，零回写），与 Go encode_tag.go 逐字节对齐：
    //
    //   [inner 已编码 payload wire][reverse(meta 正序)][meta_len 尾部槽][total_len 档位大端][tag sign]
    //
    // meta 为 Tag.ToBytes() 正序字节；total = len(payload)+len(meta)+len(尾槽)，与解码侧
    // tailframe 反推式 inner = total - ml - len(尾槽) 一致（尾槽与前缀槽同一套三档）。
    public int EncodeTaggedPayload(byte[] payload, byte[] meta)
    {
        if (meta.Length == 0)
        {
            _buf.WriteAll(payload);
            return payload.Length;
        }
        byte[] slot = Footer.FooterSlotBytes(meta.Length);
        int length = payload.Length + meta.Length + slot.Length;
        if ((long)length > (long)WireConstants.MAX_4)
        {
            throw new Exception("Tag+payload too long");
        }
        int start = _buf.Length;
        int sign = Prefix.TAG;
        byte[] tier;
        if (length < WireConstants.TAG_LEN_1)
        {
            sign |= length;
            tier = Array.Empty<byte>();
        }
        else if (length <= WireConstants.MAX_1)
        {
            sign |= WireConstants.TAG_LEN_1;
            tier = new byte[] { (byte)length };
        }
        else if (length <= WireConstants.MAX_2)
        {
            sign |= WireConstants.TAG_LEN_2;
            tier = new byte[] { (byte)(length >> 8), (byte)length };
        }
        else if (length <= WireConstants.MAX_3)
        {
            sign |= WireConstants.TAG_LEN_3;
            tier = new byte[] { (byte)(length >> 16), (byte)(length >> 8), (byte)length };
        }
        // 最后一档必须 <=：total == MAX_4 本身可表达，否则静默丢档位字节。
        else
        {
            sign |= WireConstants.TAG_LEN_4;
            tier = new byte[] { (byte)(length >> 24), (byte)(length >> 16), (byte)(length >> 8), (byte)length };
        }
        _buf.WriteAll(payload);
        _buf.WriteAll(Footer.ReverseBytes(meta));
        _buf.WriteAll(slot);
        _buf.WriteAll(tier);
        _buf.Write(sign);
        return _buf.Length - start;
    }

    public int EncodeBigIntDecimal(string s)
    {
        byte[] bits = BigIntWireCodec.EncodeSignedDecimal(s);
        return EncodeBytes(bits);
    }

    public static bool TryEncodeSimpleByName(WireEncoder enc, string name)
    {
        int? val = SimpleValue.NameToValue(name.ToLower());
        if (val.HasValue)
        {
            enc.EncodeSimple(val.Value);
            return true;
        }
        return false;
    }
}