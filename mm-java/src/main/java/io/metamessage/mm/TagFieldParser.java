package io.metamessage.mm;

import java.nio.charset.StandardCharsets;

/** Parses one tag-field chunk (Go {@code decodeTagBytes}). */
final class TagFieldParser {

    private TagFieldParser() {}

    static int parseOne(Cursor c, MmTag tag) throws MmDecodeException {
        int start = c.pos;
        int b = c.read();
        int prefix = b & 0xF8;
        int low = b & 0x07;
        if (prefix == MmTag.TagKey.K_IS_NULL) {
            tag.isNull = (low & 1) == 1;
            if (tag.isNull) {
                tag.nullable = true;
            }
        } else if (prefix == MmTag.TagKey.K_EXAMPLE) {
            tag.example = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_DESC) {
            readSizedUtf8(c, tag, low, 0);
        } else if (prefix == MmTag.TagKey.K_TYPE) {
            tag.type = ValueType.fromCode(c.read());
        } else if (prefix == MmTag.TagKey.K_RAW) {
            tag.raw = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_NULLABLE) {
            tag.nullable = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_ALLOW_EMPTY) {
            tag.allowEmpty = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_UNIQUE) {
            tag.unique = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_DEFAULT) {
            tag.defaultValue = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_MIN) {
            tag.min = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_MAX) {
            tag.max = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_SIZE) {
            tag.size = readUintN(c, low);
        } else if (prefix == MmTag.TagKey.K_ENUM) {
            tag.type = ValueType.ENUM;
            tag.enumValues = readSizedUtf8Only(c, low);
        } else if (prefix == MmTag.TagKey.K_PATTERN) {
            tag.pattern = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_LOCATION) {
            tag.locationHours = Integer.parseInt(readAscii(c, low));
        } else if (prefix == MmTag.TagKey.K_VERSION) {
            tag.version = readUintN(c, low);
        } else if (prefix == MmTag.TagKey.K_MIME) {
            readMime(c, tag, low, true);
        } else if (prefix == MmTag.TagKey.K_CHILD_DESC) {
            readSizedUtf8(c, tag, low, 1);
        } else if (prefix == MmTag.TagKey.K_CHILD_TYPE) {
            tag.childType = ValueType.fromCode(c.read());
        } else if (prefix == MmTag.TagKey.K_CHILD_RAW) {
            tag.childRaw = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_CHILD_NULLABLE) {
            tag.childNullable = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_CHILD_ALLOW_EMPTY) {
            tag.childAllowEmpty = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_CHILD_UNIQUE) {
            tag.childUnique = (low & 1) == 1;
        } else if (prefix == MmTag.TagKey.K_CHILD_DEFAULT) {
            tag.childDefault = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_MIN) {
            tag.childMin = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_MAX) {
            tag.childMax = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_SIZE) {
            tag.childSize = readUintN(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_ENUM) {
            tag.childType = ValueType.ENUM;
            tag.childEnum = readSizedUtf8Only(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_PATTERN) {
            tag.childPattern = readShortUtf8(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_LOCATION) {
            tag.childLocationHours = Integer.parseInt(readAscii(c, low));
        } else if (prefix == MmTag.TagKey.K_CHILD_VERSION) {
            tag.childVersion = readUintN(c, low);
        } else if (prefix == MmTag.TagKey.K_CHILD_MIME) {
            readMime(c, tag, low, false);
        } else {
            throw new MmDecodeException("invalid tag field prefix: 0x" + Integer.toHexString(prefix));
        }
        return c.pos - start;
    }

    private static void readSizedUtf8(Cursor c, MmTag tag, int low, int mode) throws MmDecodeException {
        String s = readSizedUtf8Only(c, low);
        if (mode == 0) {
            tag.desc = s;
        } else {
            tag.childDesc = s;
        }
    }

    private static String readSizedUtf8Only(Cursor c, int low) throws MmDecodeException {
        if (low <= 5) {
            return new String(c.readBytes(low), StandardCharsets.UTF_8);
        }
        if (low == 6) {
            int l = c.read();
            return new String(c.readBytes(l), StandardCharsets.UTF_8);
        }
        int hi = c.read();
        int lo = c.read();
        int l = hi << 8 | lo;
        return new String(c.readBytes(l), StandardCharsets.UTF_8);
    }

    private static String readShortUtf8(Cursor c, int low) throws MmDecodeException {
        if (low < 7) {
            return new String(c.readBytes(low), StandardCharsets.UTF_8);
        }
        int l = c.read();
        return new String(c.readBytes(l), StandardCharsets.UTF_8);
    }

    private static String readAscii(Cursor c, int low) throws MmDecodeException {
        return new String(c.readBytes(low), StandardCharsets.US_ASCII);
    }

    private static int readUintN(Cursor c, int low) throws MmDecodeException {
        if (low >= 8) {
            throw new MmDecodeException("uint field length");
        }
        int v = 0;
        for (int i = 0; i < low; i++) {
            v = (v << 8) | c.read();
        }
        return v;
    }

    private static void readMime(Cursor c, MmTag tag, int low, boolean self) throws MmDecodeException {
        if (low < 7) {
            if (self) {
                tag.mime = MimeWire.toString(low);
            } else {
                tag.childMime = MimeWire.toString(low);
            }
        } else {
            int l2 = c.read();
            if (self) {
                tag.mime = MimeWire.toString(l2);
            } else {
                tag.childMime = MimeWire.toString(l2);
            }
        }
    }

    static final class Cursor {
        final byte[] d;
        int pos;

        Cursor(byte[] d, int pos) {
            this.d = d;
            this.pos = pos;
        }

        int read() throws MmDecodeException {
            if (pos >= d.length) {
                throw new MmDecodeException("eof");
            }
            return d[pos++] & 0xFF;
        }

        byte[] readBytes(int n) throws MmDecodeException {
            if (pos + n > d.length) {
                throw new MmDecodeException("eof");
            }
            byte[] r = new byte[n];
            System.arraycopy(d, pos, r, 0, n);
            pos += n;
            return r;
        }
    }
}
