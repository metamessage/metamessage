namespace MetaMessage.Core;

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
            case SimpleValue.CODE:
            case SimpleValue.MESSAGE:
            case SimpleValue.DATA:
            case SimpleValue.SUCCESS:
            case SimpleValue.ERROR:
            case SimpleValue.UNKNOWN:
            case SimpleValue.PAGE:
            case SimpleValue.LIMIT:
            case SimpleValue.OFFSET:
            case SimpleValue.TOTAL:
            case SimpleValue.ID:
            case SimpleValue.NAME:
            case SimpleValue.DESCRIPTION:
            case SimpleValue.TYPE:
            case SimpleValue.VERSION:
            case SimpleValue.STATUS:
            case SimpleValue.URL:
            case SimpleValue.CREATE_TIME:
            case SimpleValue.UPDATE_TIME:
            case SimpleValue.DELETE_TIME:
            case SimpleValue.ACCOUNT:
            case SimpleValue.TOKEN:
            case SimpleValue.EXPIRE_TIME:
            case SimpleValue.KEY:
            case SimpleValue.VAL:
                {
                    string name = SimpleValue.NameOf(val);
                    tag.Type = ValueType.STR;
                    return new MmScalar(name, name, tag);
                }
            default:
                throw new MmDecodeException($"Unknown simple value: {val}");
        }
    }

    static private MmScalar NullBool(MmTag tag)
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
        return new MmScalar(false, "false", tag);
    }

    static private MmScalar NullInt(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.I;
        }
        if (tag.Type != ValueType.I && tag.Type != ValueType.I8 && tag.Type != ValueType.I16 && tag.Type != ValueType.I32 && tag.Type != ValueType.I64 && tag.Type != ValueType.U && tag.Type != ValueType.U16 && tag.Type != ValueType.U32 && tag.Type != ValueType.U64)
        {
            throw new MmDecodeException("null_int type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(0L, "0", tag);
    }

    static private MmScalar NullFloat(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.F64;
        }
        if (tag.Type != ValueType.F32 && tag.Type != ValueType.F64 && tag.Type != ValueType.DECIMAL)
        {
            throw new MmDecodeException("null_float type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(0.0, "0.0", tag);
    }

    static private MmScalar NullString(MmTag tag)
    {
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.STR;
        }
        tag.IsNull = true;
        return new MmScalar("", "", tag);
    }

    static private MmScalar NullBytes(MmTag tag)
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
        return new MmScalar(Array.Empty<byte>(), "", tag);
    }

    private IMmTree DecodePositiveInt(int first, MmTag? inherited)
    {
        long v = ReadUintBody(first);
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.I;
        }

        (object data, string text) = ConvertIntValue(v, tag);
        return new MmScalar(data, text, tag);
    }

    private IMmTree DecodeNegativeInt(int first, MmTag? inherited)
    {
        long v = ReadUintBody(first);
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.I;
        }

        (object data, string text) = ConvertIntValue(-v, tag);
        return new MmScalar(data, text, tag);
    }

    private static (object data, string text) ConvertIntValue(long v, MmTag tag)
    {
        return tag.Type switch
        {
            ValueType.DATETIME => DateTimeFromInt(v, tag),
            ValueType.DATE => DateFromInt(v, tag),
            ValueType.TIME => TimeFromInt(v, tag),
            ValueType.ENUM => EnumFromInt((int)v, tag),
            ValueType.I8 => ((sbyte)v, v.ToString()),
            ValueType.I16 => ((short)v, v.ToString()),
            ValueType.I32 => ((int)v, v.ToString()),
            ValueType.I64 => (v, v.ToString()),
            ValueType.U => ((uint)v, v.ToString()),
            ValueType.U8 => ((byte)v, v.ToString()),
            ValueType.U16 => ((ushort)v, v.ToString()),
            ValueType.U32 => ((uint)v, v.ToString()),
            ValueType.U64 => ((ulong)v, v.ToString()),
            _ => (v, v.ToString())
        };
    }

    private static (object data, string text) DateTimeFromInt(long v, MmTag tag)
    {
        if (v < 0)
        {
            return (TimeUtil.FromEpochSeconds(0), TimeUtil.FromEpochSeconds(0).ToString("yyyy-MM-dd HH:mm:ss"));
        }
        var dt = TimeUtil.FromEpochSeconds(v);
        return (dt, dt.ToString("yyyy-MM-dd HH:mm:ss"));
    }

    private static (object data, string text) DateFromInt(long v, MmTag tag)
    {
        var dt = TimeUtil.FromDaysSinceEpoch(v);
        return (dt, dt.ToString("yyyy-MM-dd"));
    }

    private static (object data, string text) TimeFromInt(long v, MmTag tag)
    {
        if (v < 0 || v > 86399)
        {
            throw new MmDecodeException("Time value out of range (0-86399)");
        }
        var dt = TimeUtil.FromSecondsOfDay(v);
        return (dt, dt.ToString("HH:mm:ss"));
    }

    private static (object data, string text) EnumFromInt(int v, MmTag tag)
    {
        if (!string.IsNullOrEmpty(tag.Enum))
        {
            var enumValues = tag.Enum.Split('|');
            if (v >= 0 && v < enumValues.Length)
            {
                return (v, enumValues[v].Trim());
            }
            throw new MmDecodeException($"Enum index {v} out of range for values: {tag.Enum}");
        }
        return (v, v.ToString());
    }

    private IMmTree DecodeFloat(int first, MmTag? inherited)
    {
        MmTag tag = inherited?.Copy() ?? MmTag.Empty();
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.F64;
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
            int exp = (sbyte)_data[_offset++];
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
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.STR;
        }

        // Type-specific string handling
        return tag.Type switch
        {
            ValueType.EMAIL or ValueType.URL or ValueType.IP or ValueType.ENUM =>
                new MmScalar(s, s, tag),
            _ => new MmScalar(s, s, tag)
        };
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
        if (tag.Type == ValueType.UNKNOWN)
        {
            tag.Type = ValueType.BYTES;
        }

        // Type-specific bytes handling
        return tag.Type switch
        {
            ValueType.UUID => BytesToUuidResult(bytes, tag),
            ValueType.IMAGE or ValueType.VIDEO =>
                new MmScalar(bytes, Convert.ToBase64String(bytes), tag),
            _ => new MmScalar(bytes, Convert.ToBase64String(bytes), tag)
        };
    }

    private static MmScalar BytesToUuidResult(byte[] bytes, MmTag tag)
    {
        if (bytes.Length != 16)
        {
            throw new MmDecodeException($"UUID bytes must be 16 bytes, got {bytes.Length}");
        }
        var guid = new Guid(bytes);
        return new MmScalar(bytes, guid.ToString(), tag);
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
            tag.Type = ValueType.VEC;
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

            var firstElem = DecodeNext(null);
            MmArray keyArray;
            if (firstElem is MmArray ka)
            {
                keyArray = ka;
            }
            else
            {
                throw new MmDecodeException("Expected key array in map container");
            }

            int keyIdx = 0;
            while (_offset < end && keyIdx < keyArray.Children.Count)
            {
                var key = (MmScalar)keyArray.Children[keyIdx];
                var value = DecodeNext(null);
                entries.Add(new KeyValuePair<MmScalar, IMmTree>(key, value));
                keyIdx++;
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

        int innerTagLen = _data[_offset++];
        int tagDataLen;
        if (innerTagLen < 254)
        {
            tagDataLen = innerTagLen;
        }
        else if (innerTagLen == 254)
        {
            if (_offset >= end)
            {
                throw new MmDecodeException("Unexpected end of tag data");
            }
            tagDataLen = _data[_offset++];
        }
        else
        {
            if (_offset + 1 >= end)
            {
                throw new MmDecodeException("Unexpected end of tag data");
            }
            tagDataLen = (_data[_offset] << 8) | _data[_offset + 1];
            _offset += 2;
        }

        byte[] tagBytes = new byte[tagDataLen];
        Array.Copy(_data, _offset, tagBytes, 0, tagDataLen);
        _offset += tagDataLen;

        MmTag tag = TagFieldParser.Parse(tagBytes);
        if (inherited != null)
        {
        }

        int limitEnd = end;
        return DecodeNext(tag);
    }

    private long ReadUintBody(int first)
    {
        int len = first & WireConstants.INT_LEN_MASK;
        long v = 0;

        if (len < WireConstants.INT_LEN_1)
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