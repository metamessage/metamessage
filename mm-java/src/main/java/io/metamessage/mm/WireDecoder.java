package io.github.metamessage.mm;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.UUID;

public final class WireDecoder {
    private final byte[] data;
    private int offset;

    public WireDecoder(byte[] data) {
        this.data = data;
        this.offset = 0;
    }

    public MmTree decode() throws MmDecodeException {
        Decoded d = decodeNode(null);
        if (offset != data.length) {
            throw new MmDecodeException("trailing bytes at " + offset + " len " + data.length);
        }
        return d.node;
    }

    private record Decoded(MmTree node, int consumed) {}

    private Decoded decodeNode(MmTag inherited) throws MmDecodeException {
        int start = offset;
        if (offset >= data.length) {
            throw new MmDecodeException("eof");
        }
        int b = data[offset++] & 0xFF;
        int prefix = b & Prefix.PREFIX_MASK;
        Decoded d =
                switch (prefix) {
                    case Prefix.TAG -> decodeTagged(b, inherited, start);
                    case Prefix.SIMPLE -> new Decoded(decodeSimple(b, inherited), offset - start);
                    case Prefix.POSITIVE_INT -> decodePositiveInt(b, inherited, start);
                    case Prefix.NEGATIVE_INT -> decodeNegativeInt(b, inherited, start);
                    case Prefix.FLOAT -> decodeFloat(b, inherited, start);
                    case Prefix.STRING -> decodeString(b, inherited, start);
                    case Prefix.BYTES -> decodeBytes(b, inherited, start);
                    case Prefix.CONTAINER -> decodeContainer(b, inherited, start);
                    default -> throw new MmDecodeException("invalid prefix");
                };
        return d;
    }

    private Decoded decodeTagged(int firstByte, MmTag inherited, int start) throws MmDecodeException {
        int[] tl = tagOuterLen(firstByte);
        int l1 = tl[0];
        int l2 = tl[1];
        if (l1 == 1) {
            l2 = data[offset++] & 0xFF;
        } else if (l1 == 2) {
            l2 = (data[offset++] & 0xFF) << 8 | (data[offset++] & 0xFF);
        }
        int innerStart = offset;
        int innerEnd = innerStart + l2;
        if (innerEnd > data.length) {
            throw new MmDecodeException("tag frame past eof");
        }
        int tb = data[offset++] & 0xFF;
        int innerFieldLen = readInnerLen(tb);
        int fieldsEnd = offset + innerFieldLen;
        if (fieldsEnd > innerEnd) {
            throw new MmDecodeException("tag fields overflow");
        }
        MmTag tag = MmTag.empty();
        while (offset < fieldsEnd) {
            int n = TagFieldParser.parseOne(new TagFieldParser.Cursor(data, offset), tag);
            if (n <= 0) {
                throw new MmDecodeException("tag error");
            }
            offset += n;
        }
        MmTree node;
        if (tag.isNull) {
            MmTree synthetic = nullScalarForTag(tag);
            if (synthetic != null) {
                node = synthetic;
            } else {
                Decoded inner = decodeNode(tag);
                node = inner.node;
            }
        } else {
            Decoded inner = decodeNode(tag);
            node = inner.node;
        }
        // Ensure offset reaches innerEnd
        offset = innerEnd;
        return new Decoded(node, offset - start);
    }

    private MmTree nullScalarForTag(MmTag tag) {
        return switch (tag.type) {
            case DATETIME -> {
                ZoneId z = zoneForHours(tag.locationHours);
                LocalDateTime dt = LocalDateTime.ofInstant(Instant.EPOCH, z != null ? z : ZoneOffset.UTC);
                yield new MmTree.MmScalar(dt, dt.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")), tag);
            }
            case DATE -> {
                LocalDate d = LocalDate.of(1970, 1, 1);
                yield new MmTree.MmScalar(d, d.toString(), tag);
            }
            case TIME -> new MmTree.MmScalar(LocalTime.MIDNIGHT, LocalTime.MIDNIGHT.toString(), tag);
            case INT8 -> new MmTree.MmScalar((byte) 0, "0", tag);
            case INT16 -> new MmTree.MmScalar((short) 0, "0", tag);
            case INT32 -> new MmTree.MmScalar(0, "0", tag);
            case INT64 -> new MmTree.MmScalar(0L, "0", tag);
            case UINT, UINT8, UINT16, UINT32 -> new MmTree.MmScalar(0, "0", tag);
            case UINT64 -> new MmTree.MmScalar(0L, "0", tag);
            case FLOAT32 -> new MmTree.MmScalar(0f, "0.0", tag);
            case FLOAT64 -> new MmTree.MmScalar(0d, "0.0", tag);
            case EMAIL, UUID, DECIMAL -> new MmTree.MmScalar("", "", tag);
            case BIGINT -> new MmTree.MmScalar(BigInteger.ZERO, "0", tag);
            case URL -> new MmTree.MmScalar("", "", tag);
            case IP -> new MmTree.MmScalar(null, ipNullText(tag.version), tag);
            default -> null;
        };
    }

    private static String ipNullText(int version) {
        return switch (version) {
            case 4 -> "0.0.0.0";
            case 6 -> "::";
            default -> "";
        };
    }

    private static ZoneId zoneForHours(int hours) {
        if (hours == 0) {
            return null;
        }
        return ZoneOffset.ofHours(hours);
    }

    private static int[] tagOuterLen(int firstByte) {
        int l = firstByte & WireConstants.TAG_LEN_MASK;
        if (l < WireConstants.TAG_LEN_1) {
            return new int[] {0, l};
        }
        if (l == WireConstants.TAG_LEN_1) {
            return new int[] {1, 0};
        }
        if (l == WireConstants.TAG_LEN_2) {
            return new int[] {2, 0};
        }
        return new int[] {0, l};
    }

    private int readInnerLen(int b) throws MmDecodeException {
        int l = b;
        if (l < 254) {
            return l;
        }
        if (l == 254) {
            int result = data[offset] & 0xFF;
            offset++;
            return result;
        }
        if (l == 255) {
            int result = (data[offset] & 0xFF) << 8 | (data[offset + 1] & 0xFF);
            offset += 2;
            return result;
        }
        return l;
    }

    private MmTree decodeSimple(int first, MmTag inherited) throws MmDecodeException {
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        int sv = first & Prefix.SUFFIX_MASK;
        return switch (sv) {
            case SimpleValue.FALSE -> {
                tag.type = ValueType.BOOL;
                yield new MmTree.MmScalar(false, "false", tag);
            }
            case SimpleValue.TRUE -> {
                tag.type = ValueType.BOOL;
                yield new MmTree.MmScalar(true, "true", tag);
            }
            case SimpleValue.NULL_BOOL -> nullBool(tag);
            case SimpleValue.NULL_INT -> nullInt(tag);
            case SimpleValue.NULL_FLOAT -> nullFloat(tag);
            case SimpleValue.NULL_STRING -> nullString(tag);
            case SimpleValue.NULL_BYTES -> nullBytes(tag);
            default -> throw new MmDecodeException("unsupported simple: " + sv);
        };
    }

    private MmTree nullBool(MmTag tag) throws MmDecodeException {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.BOOL;
        }
        if (tag.type != ValueType.BOOL) {
            throw new MmDecodeException("null_bool type mismatch");
        }
        return new MmTree.MmScalar(false, "false", tag);
    }

    private MmTree nullInt(MmTag tag) throws MmDecodeException {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.INT;
        }
        if (tag.type != ValueType.INT) {
            throw new MmDecodeException("null_int type mismatch");
        }
        return new MmTree.MmScalar(0, "0", tag);
    }

    private MmTree nullFloat(MmTag tag) throws MmDecodeException {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.FLOAT64;
        }
        if (tag.type != ValueType.FLOAT32 && tag.type != ValueType.FLOAT64) {
            throw new MmDecodeException("null_float type mismatch");
        }
        if (tag.type == ValueType.FLOAT32) {
            return new MmTree.MmScalar(0f, "0.0", tag);
        }
        return new MmTree.MmScalar(0d, "0.0", tag);
    }

    private MmTree nullString(MmTag tag) throws MmDecodeException {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.STRING;
        }
        if (tag.type != ValueType.STRING) {
            throw new MmDecodeException("null_string type mismatch");
        }
        return new MmTree.MmScalar("", "", tag);
    }

    private MmTree nullBytes(MmTag tag) throws MmDecodeException {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.BYTES;
        }
        if (tag.type != ValueType.BYTES) {
            throw new MmDecodeException("null_bytes type mismatch");
        }
        return new MmTree.MmScalar(new byte[0], "", tag);
    }

    private Decoded decodePositiveInt(int first, MmTag inherited, int start) throws MmDecodeException {
        long v = readUintBody(first);
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.INT;
        }
        return new Decoded(mapUintToTree(tag, v), offset - start);
    }

    private Decoded decodeNegativeInt(int first, MmTag inherited, int start) throws MmDecodeException {
        long v = readUintBody(first);
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.INT;
        }
        return new Decoded(mapNegativeInt(tag, v), offset - start);
    }

    private long readUintBody(int first) throws MmDecodeException {
        int l1 = intLenExtraBytes(first);
        int low = first & WireConstants.INT_LEN_MASK;
        long v;
        if (l1 == 0) {
            v = low;
        } else {
            v = 0;
            for (int i = 0; i < l1; i++) {
                v = (v << 8) | (data[offset++] & 0xFFL);
            }
        }
        return v;
    }

    private static int intLenExtraBytes(int first) {
        int l = first & WireConstants.INT_LEN_MASK;
        if (l < WireConstants.INT_LEN_1) {
            return 0;
        }
        return l - WireConstants.INT_LEN_1 + 1;
    }

    private MmTree mapUintToTree(MmTag tag, long v) throws MmDecodeException {
        return switch (tag.type) {
            case INT -> new MmTree.MmScalar((int) v, Long.toString(v), tag);
            case INT8 -> new MmTree.MmScalar((byte) v, Long.toString(v), tag);
            case INT16 -> new MmTree.MmScalar((short) v, Long.toString(v), tag);
            case INT32 -> new MmTree.MmScalar((int) v, Long.toString(v), tag);
            case INT64 -> new MmTree.MmScalar(v, Long.toString(v), tag);
            case UINT -> new MmTree.MmScalar((int) v, Long.toString(v), tag);
            case UINT8 -> new MmTree.MmScalar((short) v, Long.toString(v), tag);
            case UINT16 -> new MmTree.MmScalar((int) v, Long.toString(v), tag);
            case UINT32 -> new MmTree.MmScalar((int) v, Long.toString(v), tag);
            case UINT64 -> new MmTree.MmScalar(v, Long.toString(v), tag);
            case DATETIME -> decodeDateTime(tag, v);
            case DATE -> decodeDate(tag, v);
            case TIME -> decodeTime(tag, v);
            case ENUM -> decodeEnum(tag, v);
            default -> throw new MmDecodeException("unsupported int type: " + tag.type);
        };
    }

    private MmTree mapNegativeInt(MmTag tag, long v) throws MmDecodeException {
        return switch (tag.type) {
            case INT -> new MmTree.MmScalar((int) -v, "-" + v, tag);
            case INT8 -> new MmTree.MmScalar((byte) -v, "-" + v, tag);
            case INT16 -> new MmTree.MmScalar((short) -v, "-" + v, tag);
            case INT32 -> new MmTree.MmScalar((int) -v, "-" + v, tag);
            case INT64 -> new MmTree.MmScalar(-v, "-" + v, tag);
            default -> throw new MmDecodeException("unsupported neg int type: " + tag.type);
        };
    }

    private MmTree.MmScalar decodeDateTime(MmTag tag, long v) throws MmDecodeException {
        if (tag.isNull) {
            return new MmTree.MmScalar(null, "", tag);
        }
        Instant ins = Instant.ofEpochSecond(v);
        ZoneId z = zoneForHours(tag.locationHours);
        LocalDateTime ldt = LocalDateTime.ofInstant(ins, z != null ? z : ZoneOffset.UTC);
        return new MmTree.MmScalar(ldt, ldt.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")), tag);
    }

    private MmTree.MmScalar decodeDate(MmTag tag, long v) throws MmDecodeException {
        if (tag.isNull) {
            return new MmTree.MmScalar(null, "", tag);
        }
        if (v > Integer.MAX_VALUE) {
            throw new MmDecodeException("date overflow");
        }
        LocalDate d = TimeUtil.dateFromDays(v);
        return new MmTree.MmScalar(d, d.toString(), tag);
    }

    private MmTree.MmScalar decodeTime(MmTag tag, long v) throws MmDecodeException {
        if (tag.isNull) {
            return new MmTree.MmScalar(null, "", tag);
        }
        if (v > 86399) {
            throw new MmDecodeException("time out of range");
        }
        LocalTime t = TimeUtil.timeFromSeconds((int) v);
        return new MmTree.MmScalar(t, t.toString(), tag);
    }

    private MmTree.MmScalar decodeEnum(MmTag tag, long v) throws MmDecodeException {
        if (tag.isNull) {
            return new MmTree.MmScalar(-1, "", tag);
        }
        if (tag.enumValues == null || tag.enumValues.isEmpty()) {
            throw new MmDecodeException("enum without labels");
        }
        String[] parts = tag.enumValues.split("\\|");
        if (v >= parts.length) {
            throw new MmDecodeException("enum index out of range");
        }
        String label = parts[(int) v].trim();
        return new MmTree.MmScalar((int) v, label, tag);
    }

    private Decoded decodeFloat(int first, MmTag inherited, int start) throws MmDecodeException {
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.FLOAT64;
        }
        double val;
        int l = first & WireConstants.FLOAT_LEN_MASK;
        if (l < WireConstants.FLOAT_LEN_1) {
            val = (first & 0xF) / 10.0;
            if ((first & WireConstants.FLOAT_NEG_MASK) != 0) {
                val = -val;
            }
        } else {
            byte exp = (byte) data[offset++];
            int l1 = floatLenExtraBytes(first);
            long mantissa;
            if (l1 == 0) {
                mantissa = 0;
            } else {
                mantissa = 0;
                for (int i = 0; i < l1; i++) {
                    mantissa = (mantissa << 8) | (data[offset++] & 0xFF);
                }
            }
            String dec = FloatCodec.mantissaToDecimal(mantissa, exp);
            val = Double.parseDouble(dec);
            if ((first & WireConstants.FLOAT_NEG_MASK) != 0) {
                val = -val;
            }
        }
        MmTree node =
                switch (tag.type) {
                    case FLOAT32 -> new MmTree.MmScalar((float) val, Float.toString((float) val), tag);
                    case FLOAT64, DECIMAL -> new MmTree.MmScalar(val, Double.toString(val), tag);
                    default -> throw new MmDecodeException("bad float tag " + tag.type);
                };
        return new Decoded(node, offset - start);
    }

    private static int floatLenExtraBytes(int first) {
        int l = first & WireConstants.FLOAT_LEN_MASK;
        if (l < WireConstants.FLOAT_LEN_1) {
            return 0;
        }
        return l - WireConstants.FLOAT_LEN_1 + 1;
    }

    private Decoded decodeString(int first, MmTag inherited, int start) throws MmDecodeException {
        int[] sl = stringLen(first);
        int l2 = sl[1];
        if (sl[0] == 1) {
            l2 = data[offset++] & 0xFF;
        } else if (sl[0] == 2) {
            l2 = (data[offset++] & 0xFF) << 8 | (data[offset++] & 0xFF);
        }
        byte[] bs = l2 > 0 ? readBytes(l2) : new byte[0];
        String text = new String(bs, StandardCharsets.UTF_8);
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.STRING;
        }
        MmTree node =
                switch (tag.type) {
                    case STRING, EMAIL -> new MmTree.MmScalar(text, text, tag);
                    case URL -> new MmTree.MmScalar(text, text, tag);
                    case IP -> new MmTree.MmScalar(text, text, tag);
                    default -> throw new MmDecodeException("unsupported string type: " + tag.type);
                };
        return new Decoded(node, offset - start);
    }

    private static int[] stringLen(int first) {
        int l = first & WireConstants.STRING_LEN_MASK;
        if (l < WireConstants.STRING_LEN_1) {
            return new int[] {0, l};
        }
        if (l == WireConstants.STRING_LEN_1) {
            return new int[] {1, 0};
        }
        return new int[] {2, 0};
    }

    private Decoded decodeBytes(int first, MmTag inherited, int start) throws MmDecodeException {
        int[] bl = bytesLen(first);
        int l2 = bl[1];
        if (bl[0] == 1) {
            l2 = data[offset++] & 0xFF;
        } else if (bl[0] == 2) {
            l2 = (data[offset++] & 0xFF) << 8 | (data[offset++] & 0xFF);
        }
        byte[] bs = l2 > 0 ? readBytes(l2) : new byte[0];
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.BYTES;
        }
        MmTree node =
                switch (tag.type) {
                    case BYTES -> new MmTree.MmScalar(bs, "", tag);
                    case BIGINT -> bigintFromBytes(bs, tag);
                    case UUID -> {
                        if (bs.length != 16) {
                            throw new MmDecodeException("uuid length");
                        }
                        UUID u = uuidFromBytes(bs);
                        yield new MmTree.MmScalar(u, u.toString(), tag);
                    }
                    case IP -> new MmTree.MmScalar(bs, "", tag);
                    default -> throw new MmDecodeException("unsupported bytes type: " + tag.type);
                };
        return new Decoded(node, offset - start);
    }

    private static UUID uuidFromBytes(byte[] bs) {
        java.nio.ByteBuffer bb = java.nio.ByteBuffer.wrap(bs);
        long hi = bb.getLong();
        long lo = bb.getLong();
        return new UUID(hi, lo);
    }

    private static MmTree.MmScalar bigintFromBytes(byte[] bs, MmTag tag) throws MmDecodeException {
        if (bs.length == 0) {
            return new MmTree.MmScalar(BigInteger.ZERO, "0", tag);
        }
        int n = bs[0] & 0xFF;
        byte[] body = new byte[bs.length - 1];
        System.arraycopy(bs, 1, body, 0, body.length);
        List<Integer> bits = bigintBits(body);
        boolean neg = !bits.isEmpty() && bits.get(0) == 1;
        String digits = BigIntWireCodec.decodePositive(body, n);
        BigInteger bi = new BigInteger(neg ? "-" + digits : digits);
        return new MmTree.MmScalar(bi, bi.toString(), tag);
    }

    private static List<Integer> bigintBits(byte[] data) {
        List<Integer> bits = new ArrayList<>(data.length * 8);
        for (byte bt : data) {
            for (int i = 7; i >= 0; i--) {
                bits.add((bt >> i) & 1);
            }
        }
        return bits;
    }

    private static int[] bytesLen(int first) {
        int l = first & WireConstants.BYTES_LEN_MASK;
        if (l < WireConstants.BYTES_LEN_1) {
            return new int[] {0, l};
        }
        if (l == WireConstants.BYTES_LEN_1) {
            return new int[] {1, 0};
        }
        return new int[] {2, 0};
    }

    private Decoded decodeContainer(int first, MmTag inherited, int start) throws MmDecodeException {
        boolean isArray = (first & WireConstants.CONTAINER_MASK) == WireConstants.CONTAINER_ARRAY;
        if (isArray) {
            return decodeArray(first, inherited, start);
        }
        return decodeObject(first, inherited, start);
    }

    private Decoded decodeArray(int first, MmTag inherited, int start) throws MmDecodeException {
        int[] cl = containerLen(first);
        int l2 = cl[1];
        if (cl[0] == 1) {
            l2 = data[offset++] & 0xFF;
        } else if (cl[0] == 2) {
            l2 = (data[offset++] & 0xFF) << 8 | (data[offset++] & 0xFF);
        }
        int bodyStart = offset;
        int bodyEnd = bodyStart + l2;
        if (bodyEnd > data.length) {
            throw new MmDecodeException("array past eof");
        }
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = tag.size > 0 ? ValueType.ARRAY : ValueType.SLICE;
        }
        List<MmTree> items = new ArrayList<>();
        while (offset < bodyEnd) {
            MmTag elemTag = MmTag.empty();
            elemTag.inheritFromArrayParent(tag);
            Decoded el = decodeNode(elemTag);
            items.add(el.node);
        }
        if (offset != bodyEnd) {
            throw new MmDecodeException("array body misaligned");
        }
        return new Decoded(new MmTree.MmArray(tag, items), offset - start);
    }

    private Decoded decodeObject(int first, MmTag inherited, int start) throws MmDecodeException {
        int[] cl = containerLen(first);
        int l2 = cl[1];
        if (cl[0] == 1) {
            l2 = data[offset++] & 0xFF;
        } else if (cl[0] == 2) {
            l2 = (data[offset++] & 0xFF) << 8 | (data[offset++] & 0xFF);
        }
        int innerStart = offset;
        int innerEnd = innerStart + l2;
        if (innerEnd > data.length) {
            throw new MmDecodeException("object past eof");
        }
        MmTag tag = inherited != null ? inherited.copy() : MmTag.empty();
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.STRUCT;
        }
        int keyPrefixPos = offset;
        int keyPrefix = data[offset++] & 0xFF;
        Decoded keysDec = decodeArray(keyPrefix, null, keyPrefixPos);
        MmTree.MmArray keys = (MmTree.MmArray) keysDec.node;
        List<Map.Entry<String, MmTree>> fields = new ArrayList<>();
        int i = 0;
        while (offset < innerEnd && i < keys.items().size()) {
            MmTag elemTag = MmTag.empty();
            elemTag.inheritFromArrayParent(tag);
            Decoded val = decodeNode(elemTag);
            String key = ((MmTree.MmScalar) keys.items().get(i)).text();
            fields.add(Map.entry(key, val.node));
            i++;
        }
        // Ensure offset reaches innerEnd
        offset = innerEnd;
        return new Decoded(new MmTree.MmObject(tag, fields), offset - start);
    }

    private static int[] containerLen(int first) {
        int l = first & WireConstants.CONTAINER_LEN_MASK;
        if (l < WireConstants.CONTAINER_LEN_1) {
            return new int[] {0, l};
        }
        if (l == WireConstants.CONTAINER_LEN_1) {
            return new int[] {1, 0};
        }
        return new int[] {2, 0};
    }

    private byte[] readBytes(int n) throws MmDecodeException {
        if (offset + n > data.length) {
            throw new MmDecodeException("eof");
        }
        byte[] r = new byte[n];
        System.arraycopy(data, offset, r, 0, n);
        offset += n;
        return r;
    }
}
