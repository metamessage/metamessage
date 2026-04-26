package io.metamessage.mm;

import java.nio.charset.StandardCharsets;

/**
 * Low-level MM wire encoder; mirrors Go {@code internal/mm} encoder helpers.
 */
public final class WireEncoder {
    private final GrowableByteBuf buf = new GrowableByteBuf();

    public byte[] toByteArray() {
        return buf.copyRange(0, buf.length());
    }

    public byte[] finishTakeLast(int writtenFromEnd) {
        int end = buf.length();
        return buf.copyRange(end - writtenFromEnd, end);
    }

    public void reset() {
        buf.reset();
    }

    public int size() {
        return buf.length();
    }

    public int encodeSimple(int simpleValue) {
        int start = buf.length();
        buf.write((byte) (Prefix.SIMPLE | simpleValue));
        return buf.length() - start;
    }

    public int encodeBool(boolean v) {
        return encodeSimple(v ? SimpleValue.TRUE : SimpleValue.FALSE);
    }

    public int encodeInt64(long v) {
        if (v >= 0) {
            return encodeUintWithPrefix(Prefix.POSITIVE_INT, v);
        }
        long uv = v == Long.MIN_VALUE ? 1L << 63 : -v;
        return encodeUintWithPrefix(Prefix.NEGATIVE_INT, uv);
    }

    public int encodeUint64(long uv) {
        if (uv < 0) {
            throw new IllegalArgumentException("expected unsigned");
        }
        return encodeUintWithPrefix(Prefix.POSITIVE_INT, uv);
    }

    private int encodeUintWithPrefix(int prefix, long uv) {
        int start = buf.length();
        if (uv < WireConstants.INT_LEN_1) {
            buf.write((byte) (prefix | (int) uv));
        } else if (uv <= WireConstants.MAX_1) {
            buf.write((byte) (prefix | WireConstants.INT_LEN_1), (byte) uv);
        } else if (uv <= WireConstants.MAX_2) {
            buf.write((byte) (prefix | WireConstants.INT_LEN_2), (byte) (uv >> 8), (byte) uv);
        } else if (uv <= WireConstants.MAX_3) {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_3),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        } else if (uv <= WireConstants.MAX_4) {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_4),
                    (byte) (uv >> 24),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        } else if (uv <= WireConstants.MAX_5) {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_5),
                    (byte) (uv >> 32),
                    (byte) (uv >> 24),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        } else if (uv <= WireConstants.MAX_6) {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_6),
                    (byte) (uv >> 40),
                    (byte) (uv >> 32),
                    (byte) (uv >> 24),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        } else if (uv <= WireConstants.MAX_7) {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_7),
                    (byte) (uv >> 48),
                    (byte) (uv >> 40),
                    (byte) (uv >> 32),
                    (byte) (uv >> 24),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        } else {
            buf.write(
                    (byte) (prefix | WireConstants.INT_LEN_8),
                    (byte) (uv >> 56),
                    (byte) (uv >> 48),
                    (byte) (uv >> 40),
                    (byte) (uv >> 32),
                    (byte) (uv >> 24),
                    (byte) (uv >> 16),
                    (byte) (uv >> 8),
                    (byte) uv);
        }
        return buf.length() - start;
    }

    public int encodeFloatString(String s) {
        FloatParts p = FloatCodec.parseDecimalString(s);
        int start = buf.length();
        int sign = Prefix.FLOAT;
        if (p.negative()) {
            sign |= WireConstants.FLOAT_NEG_MASK;
        }
        if (p.exponent() == -1 && p.mantissa() <= 7) {
            buf.write((byte) (sign | (int) p.mantissa()));
        } else {
            long mantissa = p.mantissa();
            if (mantissa <= WireConstants.MAX_1) {
                buf.write((byte) (sign | WireConstants.FLOAT_LEN_1), p.exponent(), (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_2) {
                buf.write((byte) (sign | WireConstants.FLOAT_LEN_2), p.exponent(), (byte) (mantissa >> 8), (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_3) {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_3),
                        p.exponent(),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_4) {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_4),
                        p.exponent(),
                        (byte) (mantissa >> 24),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_5) {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_5),
                        p.exponent(),
                        (byte) (mantissa >> 32),
                        (byte) (mantissa >> 24),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_6) {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_6),
                        p.exponent(),
                        (byte) (mantissa >> 40),
                        (byte) (mantissa >> 32),
                        (byte) (mantissa >> 24),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            } else if (mantissa <= WireConstants.MAX_7) {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_7),
                        p.exponent(),
                        (byte) (mantissa >> 48),
                        (byte) (mantissa >> 40),
                        (byte) (mantissa >> 32),
                        (byte) (mantissa >> 24),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            } else {
                buf.write(
                        (byte) (sign | WireConstants.FLOAT_LEN_8),
                        p.exponent(),
                        (byte) (mantissa >> 56),
                        (byte) (mantissa >> 48),
                        (byte) (mantissa >> 40),
                        (byte) (mantissa >> 32),
                        (byte) (mantissa >> 24),
                        (byte) (mantissa >> 16),
                        (byte) (mantissa >> 8),
                        (byte) mantissa);
            }
        }
        return buf.length() - start;
    }

    public int encodeString(String s) {
        byte[] utf = s.getBytes(StandardCharsets.UTF_8);
        int length = utf.length;
        int start = buf.length();
        int sign = Prefix.STRING;
        if (length < WireConstants.STRING_LEN_1) {
            buf.write((byte) (sign | length));
            buf.writeAll(utf);
        } else if (length < WireConstants.MAX_1) {
            buf.write((byte) (sign | WireConstants.STRING_LEN_1), (byte) length);
            buf.writeAll(utf);
        } else if (length < WireConstants.MAX_2) {
            buf.write((byte) (sign | WireConstants.STRING_LEN_2), (byte) (length >> 8), (byte) length);
            buf.writeAll(utf);
        } else {
            throw new IllegalArgumentException("string too long");
        }
        return buf.length() - start;
    }

    public int encodeBytes(byte[] data) {
        int length = data.length;
        int start = buf.length();
        int sign = Prefix.BYTES;
        if (length < WireConstants.BYTES_LEN_1) {
            buf.write((byte) (sign | length));
            buf.writeAll(data);
        } else if (length < WireConstants.MAX_1) {
            buf.write((byte) (sign | WireConstants.BYTES_LEN_1), (byte) length);
            buf.writeAll(data);
        } else if (length < WireConstants.MAX_2) {
            buf.write((byte) (sign | WireConstants.BYTES_LEN_2), (byte) (length >> 8), (byte) length);
            buf.writeAll(data);
        } else {
            throw new IllegalArgumentException("bytes too long");
        }
        return buf.length() - start;
    }

    public int encodeArrayPayload(byte[] payload) {
        return encodeContainer(payload, Prefix.CONTAINER | WireConstants.CONTAINER_ARRAY);
    }

    public int encodeObjectPayload(byte[] payload) {
        return encodeContainer(payload, Prefix.CONTAINER | WireConstants.CONTAINER_MAP);
    }

    private int encodeContainer(byte[] payload, int baseSign) {
        int length = payload.length;
        int start = buf.length();
        if (length < WireConstants.CONTAINER_LEN_1) {
            buf.write((byte) (baseSign | length));
            buf.writeAll(payload);
        } else if (length < WireConstants.MAX_1) {
            buf.write((byte) (baseSign | WireConstants.CONTAINER_LEN_1), (byte) length);
            buf.writeAll(payload);
        } else if (length < WireConstants.MAX_2) {
            buf.write((byte) (baseSign | WireConstants.CONTAINER_LEN_2), (byte) (length >> 8), (byte) length);
            buf.writeAll(payload);
        } else {
            throw new IllegalArgumentException("container payload too long");
        }
        return buf.length() - start;
    }

    /** Prefixes inner tag bytes with length (Go {@code encodeT}). */
    public int encodeTagInner(byte[] tagBytes) {
        if (tagBytes.length == 0) {
            return 0;
        }
        if (tagBytes.length > WireConstants.MAX_2) {
            throw new IllegalArgumentException("tag too long");
        }
        int start = buf.length();
        int length = tagBytes.length;
        if (length < 254) {
            buf.write((byte) length);
            buf.writeAll(tagBytes);
        } else if (length < 257) {
            buf.write((byte) 254, (byte) length);
            buf.writeAll(tagBytes);
        } else {
            buf.write((byte) 255, (byte) (length >> 8), (byte) length);
            buf.writeAll(tagBytes);
        }
        return buf.length() - start;
    }

    /** @param rawTagFields same as Go {@code (*Tag).Bytes()} (no {@code encodeT} wrapper). */
    public int encodeTaggedPayload(byte[] payload, byte[] rawTagFields) {
        if (rawTagFields.length == 0) {
            buf.writeAll(payload);
            return payload.length;
        }
        WireEncoder tEnc = new WireEncoder();
        int tagEncodedLength = tEnc.encodeTagInner(rawTagFields);
        byte[] tagEncoded = tEnc.toByteArray();
        int length = tagEncodedLength + payload.length;
        if (length > WireConstants.MAX_2) {
            throw new IllegalArgumentException("tag+payload too long");
        }
        int start = buf.length();
        int sign = Prefix.TAG;
        if (length < WireConstants.TAG_LEN_1) {
            buf.write((byte) (sign | length));
            buf.writeAll(tagEncoded);
            buf.writeAll(payload);
        } else if (length < WireConstants.MAX_1) {
            buf.write((byte) (sign | WireConstants.TAG_LEN_1), (byte) length);
            buf.writeAll(tagEncoded);
            buf.writeAll(payload);
        } else {
            buf.write((byte) (sign | WireConstants.TAG_LEN_2), (byte) (length >> 8), (byte) length);
            buf.writeAll(tagEncoded);
            buf.writeAll(payload);
        }
        return buf.length() - start;
    }

    public int encodeBigIntDecimal(String s) {
        byte[] bits = BigIntWireCodec.encodeSignedDecimal(s);
        byte[] inner = new byte[1 + bits.length];
        inner[0] = (byte) s.length();
        System.arraycopy(bits, 0, inner, 1, bits.length);
        return encodeBytes(inner);
    }

    public int sliceFrom(int start) {
        return buf.length() - start;
    }

    public byte[] copyLast(int n) {
        int end = buf.length();
        return buf.copyRange(end - n, end);
    }
}
