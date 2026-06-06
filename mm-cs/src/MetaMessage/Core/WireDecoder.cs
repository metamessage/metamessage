using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

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

    private IMmTree DecodeNext(Tag? inherited)
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

    private IMmTree DecodeSimple(int first, Tag? inherited)
    {
        int val = first & Prefix.SUFFIX_MASK;
        Tag tag = inherited?.Copy() ?? Tag.Empty();

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
                tag.Type = ValueType.Bool;
                return new MmScalar(false, "false", tag);
            case SimpleValue.TRUE:
                tag.Type = ValueType.Bool;
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
                    tag.Type = ValueType.Str;
                    return new MmScalar(name, name, tag);
                }
            default:
                throw new MmDecodeException($"Unknown simple value: {val}");
        }
    }

    static private MmScalar NullBool(Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Bool;
        }
        if (tag.Type != ValueType.Bool)
        {
            throw new MmDecodeException("null_bool type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(false, "false", tag);
    }

    static private MmScalar NullInt(Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
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

    static private MmScalar NullFloat(Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.F64;
        }
        if (tag.Type != ValueType.F32 && tag.Type != ValueType.F64 && tag.Type != ValueType.Decimal)
        {
            throw new MmDecodeException("null_float type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(0.0, "0.0", tag);
    }

    static private MmScalar NullString(Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Str;
        }
        tag.IsNull = true;
        return new MmScalar("", "", tag);
    }

    static private MmScalar NullBytes(Tag tag)
    {
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Bytes;
        }
        if (tag.Type != ValueType.Bytes)
        {
            throw new MmDecodeException("null_bytes type mismatch");
        }
        tag.IsNull = true;
        return new MmScalar(Array.Empty<byte>(), "", tag);
    }

    private IMmTree DecodePositiveInt(int first, Tag? inherited)
    {
        ulong uv = ReadUintBody(first);
        Tag tag = inherited?.Copy() ?? Tag.Empty();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.I;
        }

        (object data, string text) = ConvertPositiveIntValue(uv, tag);
        return new MmScalar(data, text, tag);
    }

    private IMmTree DecodeNegativeInt(int first, Tag? inherited)
    {
        ulong uv = ReadUintBody(first);
        Tag tag = inherited?.Copy() ?? Tag.Empty();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.I;
        }

        (object data, string text) = ConvertNegativeIntValue(uv, tag);
        return new MmScalar(data, text, tag);
    }

    private static (object data, string text) ConvertPositiveIntValue(ulong uv, Tag tag)
    {
        string text = uv.ToString();
        return tag.Type switch
        {
            ValueType.I8 => ((sbyte)uv, text),
            ValueType.I16 => ((short)uv, text),
            ValueType.I32 => ((int)uv, text),
            ValueType.I64 => ((long)uv, text),
            ValueType.U => ((uint)uv, text),
            ValueType.U8 => ((byte)uv, text),
            ValueType.U16 => ((ushort)uv, text),
            ValueType.U32 => ((uint)uv, text),
            ValueType.U64 => (uv, text),
            ValueType.Datetime => ((long)uv, DateTimeFromInt((long)uv, tag)),
            ValueType.Date => ((long)uv, DateFromInt((long)uv, tag)),
            ValueType.Time => ((long)uv, TimeFromInt((long)uv, tag)),
            ValueType.Enums => ((long)uv, EnumFromInt((int)uv, tag)),
            _ => ((long)uv, text)
        };
    }

    private static string DateTimeFromInt(long v, Tag tag)
    {
        long locationOffset = (tag?.Location ?? 0) * 3600L;
        long local = v + locationOffset;
        var dt = DateTime.UnixEpoch.AddSeconds(local);
        return dt.ToString("yyyy-MM-dd HH:mm:ss");
    }

    private static string DateFromInt(long v, Tag tag)
    {
        long locationOffset = (tag?.Location ?? 0) * 3600L;
        long local = v + locationOffset;
        var dt = DateTime.UnixEpoch.AddDays(local);
        return dt.ToString("yyyy-MM-dd");
    }

    private static string TimeFromInt(long v, Tag tag)
    {
        long locationOffset = (tag?.Location ?? 0) * 3600L;
        long local = v + locationOffset;
        var dt = DateTime.UnixEpoch.AddSeconds(local);
        return dt.ToString("HH:mm:ss");
    }

    private static string EnumFromInt(int v, Tag tag)
    {
        var enumStr = !string.IsNullOrEmpty(tag.Enums) ? tag.Enums : tag.ChildEnums;
        if (!string.IsNullOrEmpty(enumStr))
        {
            var parts = enumStr.Split('|');
            if (v >= 0 && v < parts.Length)
            {
                return parts[v].Trim();
            }
        }
        return v.ToString();
    }

    private static (object data, string text) ConvertNegativeIntValue(ulong uv, Tag tag)
    {
        string text = "-" + uv.ToString();
        return tag.Type switch
        {
            ValueType.I => (-(long)uv, text),
            ValueType.I8 => ((sbyte)(-(long)uv), text),
            ValueType.I16 => ((short)(-(long)uv), text),
            ValueType.I32 => ((int)(-(long)uv), text),
            ValueType.I64 => (-(long)uv, text),
            _ => (-(long)uv, text)
        };
    }

    private IMmTree DecodeFloat(int first, Tag? inherited)
    {
        Tag tag = inherited?.Copy() ?? Tag.Empty();
        if (tag.Type == ValueType.Unknown)
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

    private IMmTree DecodeString(int first, Tag? inherited)
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

        Tag tag = inherited?.Copy() ?? Tag.Empty();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Str;
        }

        // Type-specific string handling
        return tag.Type switch
        {
            ValueType.Email or ValueType.Url or ValueType.Ip or ValueType.Enums =>
                new MmScalar(s, s, tag),
            _ => new MmScalar(s, s, tag)
        };
    }

    private IMmTree DecodeBytes(int first, Tag? inherited)
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

        Tag tag = inherited?.Copy() ?? Tag.Empty();
        if (tag.Type == ValueType.Unknown)
        {
            tag.Type = ValueType.Bytes;
        }

        // Type-specific bytes handling
        return tag.Type switch
        {
            ValueType.Uuid => BytesToUuidResult(bytes, tag),
            ValueType.Media =>
                new MmScalar(bytes, Convert.ToBase64String(bytes), tag),
            ValueType.Bigint =>
                new MmScalar(bytes, BigIntWireCodec.DecodeSignedDecimal(bytes), tag),
            _ => new MmScalar(bytes, Convert.ToBase64String(bytes), tag)
        };
    }

    private static MmScalar BytesToUuidResult(byte[] bytes, Tag tag)
    {
        if (bytes.Length != 16)
        {
            throw new MmDecodeException($"UUID bytes must be 16 bytes, got {bytes.Length}");
        }
        var sb = new System.Text.StringBuilder();
        for (int i = 0; i < 16; i++)
        {
            sb.Append(bytes[i].ToString("x2"));
            if (i == 3 || i == 5 || i == 7 || i == 9)
                sb.Append('-');
        }
        return new MmScalar(bytes, sb.ToString(), tag);
    }

    private IMmTree DecodeContainer(int first, Tag? inherited)
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

        Tag tag = inherited?.Copy() ?? Tag.Empty();

        if (containerType == WireConstants.CONTAINER_ARRAY)
        {
            tag.Type = ValueType.Vec;
            var children = new List<IMmTree>();
            while (_offset < end)
            {
                Tag itemTag = Tag.Empty();
                itemTag.InheritFromArrayParent(tag);
                children.Add(DecodeNext(itemTag));
            }
            return new MmArray(children, tag);
        }
        else // CONTAINER_MAP
        {
            tag.Type = ValueType.Map;
            var entries = new List<KeyValuePair<MmScalar, IMmTree>>();

            var firstElem = DecodeNext(tag);
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
                Tag valueTag = Tag.Empty();
                valueTag.Inherit(tag);
                var value = DecodeNext(valueTag);
                entries.Add(new KeyValuePair<MmScalar, IMmTree>(key, value));
                keyIdx++;
            }

            _offset = end;
            return new MmMap(entries, tag);
        }
    }

    private IMmTree DecodeTagged(int first, Tag? inherited)
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

        Tag tag = TagFieldParser.Parse(tagBytes);
        if (inherited != null)
        {
            tag.Inherit(inherited);
        }

        if (tag.IsNull)
        {
            if (tag.Type == ValueType.Unknown)
            {
                var payloadTree = DecodeNext(tag);
                _offset = end;
                return payloadTree;
            }
            _offset = end;
            return CreateNullValue(tag);
        }

        return DecodeNext(tag);
    }

    private static IMmTree CreateNullValue(Tag tag)
    {
        tag.Nullable = true;
        switch (tag.Type)
        {
            case ValueType.Bool:
                return new MmScalar(false, "false", tag);
            case ValueType.I:
            case ValueType.I8:
            case ValueType.I16:
            case ValueType.I32:
            case ValueType.I64:
                return new MmScalar(0L, "0", tag);
            case ValueType.U:
            case ValueType.U8:
            case ValueType.U16:
            case ValueType.U32:
            case ValueType.U64:
                return new MmScalar(0UL, "0", tag);
            case ValueType.F32:
                return new MmScalar(0.0f, "0.0", tag);
            case ValueType.F64:
            case ValueType.Decimal:
                return new MmScalar(0.0, "0.0", tag);
            case ValueType.Str:
            case ValueType.Email:
                return new MmScalar("", "", tag);
            case ValueType.Bytes:
                return new MmScalar(Array.Empty<byte>(), "", tag);
            case ValueType.Datetime:
                tag.Type = ValueType.Datetime;
                return new MmScalar(DateTime.UnixEpoch, "0001-01-01 00:00:00", tag);
            default:
                return new MmScalar("", "", tag);
        }
    }

    private ulong ReadUintBody(int first)
    {
        int len = first & WireConstants.INT_LEN_MASK;
        ulong v = 0;

        if (len < WireConstants.INT_LEN_1)
        {
            v = (ulong)len;
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
            v = (ulong)((_data[_offset] << 8) | _data[_offset + 1]);
            _offset += 2;
        }
        else if (len == WireConstants.INT_LEN_3)
        {
            if (_offset + 2 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = (ulong)((_data[_offset] << 16) | (_data[_offset + 1] << 8) | _data[_offset + 2]);
            _offset += 3;
        }
        else if (len == WireConstants.INT_LEN_4)
        {
            if (_offset + 3 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = (ulong)((_data[_offset] << 24) | (_data[_offset + 1] << 16) | (_data[_offset + 2] << 8) | _data[_offset + 3]);
            _offset += 4;
        }
        else if (len == WireConstants.INT_LEN_5)
        {
            if (_offset + 4 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((ulong)_data[_offset] << 32) | ((ulong)_data[_offset + 1] << 24) | ((ulong)_data[_offset + 2] << 16) | ((ulong)_data[_offset + 3] << 8) | _data[_offset + 4];
            _offset += 5;
        }
        else if (len == WireConstants.INT_LEN_6)
        {
            if (_offset + 5 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((ulong)_data[_offset] << 40) | ((ulong)_data[_offset + 1] << 32) | ((ulong)_data[_offset + 2] << 24) | ((ulong)_data[_offset + 3] << 16) | ((ulong)_data[_offset + 4] << 8) | _data[_offset + 5];
            _offset += 6;
        }
        else if (len == WireConstants.INT_LEN_7)
        {
            if (_offset + 6 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((ulong)_data[_offset] << 48) | ((ulong)_data[_offset + 1] << 40) | ((ulong)_data[_offset + 2] << 32) | ((ulong)_data[_offset + 3] << 24) | ((ulong)_data[_offset + 4] << 16) | ((ulong)_data[_offset + 5] << 8) | _data[_offset + 6];
            _offset += 7;
        }
        else if (len == WireConstants.INT_LEN_8)
        {
            if (_offset + 7 >= _data.Length)
            {
                throw new MmDecodeException("Unexpected end of data");
            }
            v = ((ulong)_data[_offset] << 56) | ((ulong)_data[_offset + 1] << 48) | ((ulong)_data[_offset + 2] << 40) | ((ulong)_data[_offset + 3] << 32) | ((ulong)_data[_offset + 4] << 24) | ((ulong)_data[_offset + 5] << 16) | ((ulong)_data[_offset + 6] << 8) | _data[_offset + 7];
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