namespace MetaMessage.Mm;

public class MmDecodeException : Exception
{
    public MmDecodeException(string message) : base(message) { }
}

public class WireDecoder
{
    private byte[] _data;
    private int _offset;

    public WireDecoder(byte[] data)
    {
        _data = data;
        _offset = 0;
    }

    public IMmTree Decode()
    {
        if (_data == null || _data.Length == 0)
        {
            throw new MmDecodeException("Empty data");
        }

        return DecodeNext(null);
    }

    private IMmTree DecodeNext(MmTag? inherited)
    {
        if (_offset >= _data.Length)
        {
            throw new MmDecodeException("Unexpected end of data");
        }

        int first = _data[_offset++];
        int p = Prefix.Of(first);

        switch (p)
        {
            case Prefix.SIMPLE:
                return DecodeSimple(first, inherited);
            case Prefix.POSITIVE_INT:
                return DecodePositiveInt(first, inherited);
            case Prefix.NEGATIVE_INT:
                return DecodeNegativeInt(first, inherited);
            case Prefix.FLOAT:
                return DecodeFloat(first, inherited);
            case Prefix.STRING:
                return DecodeString(first, inherited);
            case Prefix.BYTES:
                return DecodeBytes(first, inherited);
            case Prefix.CONTAINER:
                return DecodeContainer(first, inherited);
            case Prefix.TAG:
                return DecodeTagged(first, inherited);
            default:
                throw new MmDecodeException($"Unknown prefix: {p}");
        }
    }

    private IMmTree DecodeSimple(int first, MmTag? inherited)
    {
        int val = first & Prefix.SUFFIX_MASK;
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();

        switch (val)
        {
            case SimpleValue.NULL_BOOL:
                return NullBool(tag);
            case SimpleValue.NULL_INT:
                return NullInt(tag);
            case SimpleValue.NULL_FLOAT:
                return NullFloat(tag);
            case SimpleValue.NULL_STRING:
                return NullString(tag);
            case SimpleValue.NULL_BYTES:
                return NullBytes(tag);
            case SimpleValue.FALSE:
                tag.Type = ValueType.BOOL;
                return new MmScalar(false, "false", tag);
            case SimpleValue.TRUE:
                tag.Type = ValueType.BOOL;
                return new MmScalar(true, "true", tag);
            case SimpleValue.NAN:
                tag.Type = ValueType.FLOAT64;
                return new MmScalar(double.NaN, "NaN", tag);
            case SimpleValue.INFINITY:
                tag.Type = ValueType.FLOAT64;
                return new MmScalar(double.PositiveInfinity, "Infinity", tag);
            case SimpleValue.NEG_INFINITY:
                tag.Type = ValueType.FLOAT64;
                return new MmScalar(double.NegativeInfinity, "-Infinity", tag);
            default:
                throw new MmDecodeException($"Unknown simple value: {val}");
        }
    }

    private IMmTree NullBool(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.BOOL;
        }
        if (tag.Type != ValueType.BOOL)
        {
            throw new MmDecodeException("null_bool type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(null, "null", tag);
    }

    private IMmTree NullInt(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.INT;
        }
        if (tag.Type != ValueType.INT && tag.Type != ValueType.INT8 && tag.Type != ValueType.INT16 && tag.Type != ValueType.INT32 && tag.Type != ValueType.INT64 && tag.Type != ValueType.UINT && tag.Type != ValueType.UINT16 && tag.Type != ValueType.UINT32 && tag.Type != ValueType.UINT64)
        {
            throw new MmDecodeException("null_int type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(null, "null", tag);
    }

    private IMmTree NullFloat(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.FLOAT64;
        }
        if (tag.Type != ValueType.FLOAT32 && tag.Type != ValueType.FLOAT64 && tag.Type != ValueType.DECIMAL)
        {
            throw new MmDecodeException("null_float type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(null, "null", tag);
    }

    private IMmTree NullString(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.STRING;
        }
        if (tag.Type != ValueType.STRING && tag.Type != ValueType.EMAIL && tag.Type != ValueType.URL)
        {
            throw new MmDecodeException("null_string type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(null, "null", tag);
    }

    private IMmTree NullBytes(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.BYTES;
        }
        if (tag.Type != ValueType.BYTES)
        {
            throw new MmDecodeException("null_bytes type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(null, "null", tag);
    }

    private IMmTree DecodePositiveInt(int first, MmTag? inherited)
    {
        long v = ReadUintBody(first);
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        tag.Type = ValueType.INT;
        return new MmScalar(v, v.ToString(), tag);
    }

    private IMmTree DecodeNegativeInt(int first, MmTag? inherited)
    {
        long v = ReadUintBody(first);
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        tag.Type = ValueType.INT;
        return new MmScalar(-v, (-v).ToString(), tag);
    }

    private IMmTree DecodeFloat(int first, MmTag? inherited)
    {
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.FLOAT64;
        }

        double val;
        int l = first & WireConstants.FLOAT_LEN_MASK;
        if (l < WireConstants.FLOAT_LEN_1)
        {
            int mantissa = l;
            val = mantissa / 10.0;
            if ((first & WireConstants.FLOAT_NEG_MASK) != 0)
            {
                val = -val;
            }
        }
        else
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            int exp = _data[_offset++];
            int l1 = FloatLenExtraBytes(first);
            long mantissa = 0;
            if (l1 == 0)
            {
                mantissa = 0;
            }
            else
            {
                for (int i = 0; i < l1; i++)
                {
                    if (_offset >= _data.Length)
                    {
                        throw new MmDecodeException("Unexpected end of data");
                    }
                    mantissa = (mantissa << 8) | _data[_offset++];
                }
            }
            string dec = FloatCodec.MantissaToDecimal(mantissa, exp);
            val = double.Parse(dec);
            if ((first & WireConstants.FLOAT_NEG_MASK) != 0)
            {
                val = -val;
            }
        }

        return new MmScalar(val, val.ToString(), tag);
    }

    private IMmTree DecodeString(int first, MmTag? inherited)
    {
        var (l1, l2) = StringLen(first);
        if (l1 == 1)
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = _data[_offset++];
        }
        else if (l1 == 2)
        {
            if (_offset + 1 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }

        if (_offset + l2 > _data.Length)
        {
            throw new MmDecodeException("String data overflow");
        }

        string s = System.Text.Encoding.UTF8.GetString(_data, _offset, l2);
        _offset += l2;

        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        tag.Type = ValueType.STRING;
        return new MmScalar(s, s, tag);
    }

    private IMmTree DecodeBytes(int first, MmTag? inherited)
    {
        var (l1, l2) = BytesLen(first);
        if (l1 == 1)
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = _data[_offset++];
        }
        else if (l1 == 2)
        {
            if (_offset + 1 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }

        if (_offset + l2 > _data.Length)
        {
            throw new MmDecodeException("Bytes data overflow");
        }

        byte[] bytes = new byte[l2];
        Array.Copy(_data, _offset, bytes, 0, l2);
        _offset += l2;

        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        tag.Type = ValueType.BYTES;
        return new MmScalar(bytes, Convert.ToBase64String(bytes), tag);
    }

    private IMmTree DecodeContainer(int first, MmTag? inherited)
    {
        int containerType = first & WireConstants.CONTAINER_MASK;
        var (l1, l2) = ContainerLen(first);

        if (l1 == 1)
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = _data[_offset++];
        }
        else if (l1 == 2)
        {
            if (_offset + 1 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }

        int end = _offset + l2;
        if (end > _data.Length)
        {
            throw new MmDecodeException("Container data overflow");
        }

        MmTag tag = inherited?.Copy() ?? MmTag.Empty();

        if (containerType == WireConstants.CONTAINER_ARRAY)
        {
            tag.Type = ValueType.SLICE;
            var children = new List<IMmTree>();
            while (_offset < end)
            {
                MmTag itemTag = MmTag.Empty();
                itemTag.InheritFromArrayParent(tag);
                children.Add(DecodeNext(itemTag));
            }
            return new MmArray(children, tag);
        }
        else // CONTAINER_MAP
        {
            tag.Type = ValueType.MAP;
            var entries = new List<KeyValuePair<MmScalar, IMmTree>>();
            while (_offset < end)
            {
                MmTag keyTag = MmTag.Empty();
                keyTag.Type = ValueType.STRING;
                var key = (MmScalar)DecodeNext(keyTag);
                var value = DecodeNext(null);
                entries.Add(new KeyValuePair<MmScalar, IMmTree>(key, value));
            }
            _offset = end;
            return new MmMap(entries, tag);
        }
    }

    private IMmTree DecodeTagged(int first, MmTag? inherited)
    {
        var (l1, l2) = TagLen(first);
        if (l1 == 1)
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = _data[_offset++];
        }
        else if (l1 == 2)
        {
            if (_offset + 1 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            l2 = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }

        int end = _offset + l2;
        if (end > _data.Length)
        {
            throw new MmDecodeException("Tag data overflow");
        }

        byte[] tagBytes = new byte[l2];
        Array.Copy(_data, _offset, tagBytes, 0, l2);
        _offset += l2;

        MmTag tag = TagFieldParser.Parse(tagBytes);
        if (inherited != null)
        {
        }

        return DecodeNext(tag);
    }

    private long ReadUintBody(int first)
    {
        int len = first & WireConstants.INT_LEN_MASK;
        long v = 0;

        if (len <= WireConstants.INT_LEN_8 - 1)
        {
            v = len;
        }
        else if (len == WireConstants.INT_LEN_1)
        {
            if (_offset >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = _data[_offset++];
        }
        else if (len == WireConstants.INT_LEN_2)
        {
            if (_offset + 1 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }
        else if (len == WireConstants.INT_LEN_3)
        {
            if (_offset + 2 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = (_data[_offset] << 16) | (_data[_offset + 1] << 8) | _data[_offset + 2];
            _offset += 3;
        }
        else if (len == WireConstants.INT_LEN_4)
        {
            if (_offset + 3 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = (_data[_offset] << 24) | (_data[_offset + 1] << 16) | (_data[_offset + 2] << 8) | _data[_offset + 3];
            _offset += 4;
        }
        else if (len == WireConstants.INT_LEN_5)
        {
            if (_offset + 4 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((long)_data[_offset] << 32) | ((long)_data[_offset + 1] << 24) | ((long)_data[_offset + 2] << 16) | ((long)_data[_offset + 3] << 8) | _data[_offset + 4];
            _offset += 5;
        }
        else if (len == WireConstants.INT_LEN_6)
        {
            if (_offset + 5 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((long)_data[_offset] << 40) | ((long)_data[_offset + 1] << 32) | ((long)_data[_offset + 2] << 24) | ((long)_data[_offset + 3] << 16) | ((long)_data[_offset + 4] << 8) | _data[_offset + 5];
            _offset += 6;
        }
        else if (len == WireConstants.INT_LEN_7)
        {
            if (_offset + 6 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((long)_data[_offset] << 48) | ((long)_data[_offset + 1] << 40) | ((long)_data[_offset + 2] << 32) | ((long)_data[_offset + 3] << 24) | ((long)_data[_offset + 4] << 16) | ((long)_data[_offset + 5] << 8) | _data[_offset + 6];
            _offset += 7;
        }
        else if (len == WireConstants.INT_LEN_8)
        {
            if (_offset + 7 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((long)_data[_offset] << 56) | ((long)_data[_offset + 1] << 48) | ((long)_data[_offset + 2] << 40) | ((long)_data[_offset + 3] << 32) | ((long)_data[_offset + 4] << 24) | ((long)_data[_offset + 5] << 16) | ((long)_data[_offset + 6] << 8) | _data[_offset + 7];
            _offset += 8;
        }

        return v;
    }

    private (int, int) StringLen(int first)
    {
        int len = first & WireConstants.STRING_LEN_MASK;
        if (len < WireConstants.STRING_LEN_1)
        {
            return (0, len);
        }
        else if (len == WireConstants.STRING_LEN_1)
        {
            return (1, 0);
        }
        else // len == WireConstants.STRING_LEN_2
        {
            return (2, 0);
        }
    }

    private (int, int) BytesLen(int first)
    {
        int len = first & WireConstants.BYTES_LEN_MASK;
        if (len < WireConstants.BYTES_LEN_1)
        {
            return (0, len);
        }
        else if (len == WireConstants.BYTES_LEN_1)
        {
            return (1, 0);
        }
        else // len == WireConstants.BYTES_LEN_2
        {
            return (2, 0);
        }
    }

    private (int, int) ContainerLen(int first)
    {
        int len = first & WireConstants.CONTAINER_LEN_MASK;
        if (len < WireConstants.CONTAINER_LEN_1)
        {
            return (0, len);
        }
        else if (len == WireConstants.CONTAINER_LEN_1)
        {
            return (1, 0);
        }
        else // len == WireConstants.CONTAINER_LEN_2
        {
            return (2, 0);
        }
    }

    private (int, int) TagLen(int first)
    {
        int len = first & WireConstants.TAG_LEN_MASK;
        if (len < WireConstants.TAG_LEN_1)
        {
            return (0, len);
        }
        else if (len == WireConstants.TAG_LEN_1)
        {
            return (1, 0);
        }
        else // len == WireConstants.TAG_LEN_2
        {
            return (2, 0);
        }
    }

    private int FloatLenExtraBytes(int first)
    {
        int l = first & WireConstants.FLOAT_LEN_MASK;
        return l < WireConstants.FLOAT_LEN_1 ? 0 : l - WireConstants.FLOAT_LEN_1 + 1;
    }
}