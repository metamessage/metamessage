package io.github.metamessage.mm;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Objects;

/**
 * Runtime tag metadata; mirrors Go {@code ast.Tag}. Serialized with the same layout as {@code (*Tag).Bytes()}.
 */
public final class MmTag {
    public static final int DEFAULT_VERSION = 0;

    public String name = "";
    public boolean isNull;
    public boolean example;
    public String desc = "";
    public ValueType type = ValueType.UNKNOWN;
    public boolean raw;
    public boolean nullable;
    public boolean allowEmpty;
    public boolean unique;
    public String defaultValue = "";
    public String min = "";
    public String max = "";
    public int size;
    public String enumValues = "";
    public String pattern = "";
    /** Hours east of UTC for wire tag; 0 = omit (UTC). */
    public int locationHours;
    public int version = DEFAULT_VERSION;
    public String mime = "";

    public String childDesc = "";
    public ValueType childType = ValueType.UNKNOWN;
    public boolean childRaw;
    public boolean childNullable;
    public boolean childAllowEmpty;
    public boolean childUnique;
    public String childDefault = "";
    public String childMin = "";
    public String childMax = "";
    public int childSize;
    public String childEnum = "";
    public String childPattern = "";
    public int childLocationHours;
    public int childVersion = DEFAULT_VERSION;
    public String childMime = "";

    public boolean isInherit;

    public static MmTag empty() {
        return new MmTag();
    }

    public static MmTag fromAnnotation(MM mm) {
        MmTag t = new MmTag();
        if (mm == null) {
            return t;
        }
        t.name = mm.name();
        t.isNull = mm.isNull();
        t.example = mm.example();
        t.desc = mm.desc();
        t.type = mm.type();
        t.raw = mm.raw();
        t.nullable = mm.nullable();
        t.allowEmpty = mm.allowEmpty();
        t.unique = mm.unique();
        t.defaultValue = mm.defaultValue();
        t.min = mm.min();
        t.max = mm.max();
        t.size = mm.size();
        t.enumValues = mm.enumValues();
        if (!t.enumValues.isEmpty()) {
            t.type = ValueType.ENUM;
        }
        t.pattern = mm.pattern();
        t.locationHours = mm.location();
        t.version = mm.version();
        t.mime = mm.mime();

        t.childDesc = mm.childDesc();
        t.childType = mm.childType();
        t.childRaw = mm.childRaw();
        t.childNullable = mm.childNullable();
        t.childAllowEmpty = mm.childAllowEmpty();
        t.childUnique = mm.childUnique();
        t.childDefault = mm.childDefault();
        t.childMin = mm.childMin();
        t.childMax = mm.childMax();
        t.childSize = mm.childSize();
        t.childEnum = mm.childEnum();
        if (!t.childEnum.isEmpty()) {
            t.childType = ValueType.ENUM;
        }
        t.childPattern = mm.childPattern();
        t.childLocationHours = mm.childLocation();
        t.childVersion = mm.childVersion();
        t.childMime = mm.childMime();
        return t;
    }

    /** Apply array/slice child constraints to a value tag (Go {@code Tag.Inherit}). */
    public void inheritFromArrayParent(MmTag parent) {
        if (parent == null) {
            return;
        }
        desc = parent.childDesc;
        type = parent.childType;
        raw = parent.childRaw;
        nullable = parent.childNullable;
        allowEmpty = parent.childAllowEmpty;
        unique = parent.childUnique;
        defaultValue = parent.childDefault;
        min = parent.childMin;
        max = parent.childMax;
        size = parent.childSize;
        enumValues = parent.childEnum;
        pattern = parent.childPattern;
        locationHours = parent.childLocationHours;
        version = parent.childVersion;
        mime = parent.childMime;
        isInherit = true;
    }

    /** Merge like Go {@code ast.MergeTag(dst, src)} (for JSONC comment merging). */
    public static MmTag merge(MmTag dst, MmTag src) {
        if (src == null) {
            return dst;
        }
        if (dst == null) {
            return src.copy();
        }
        MmTag d = dst;
        if (src.isNull) {
            d.isNull = true;
        }
        if (src.example) {
            d.example = true;
        }
        if (src.desc != null && !src.desc.isEmpty()) {
            d.desc = src.desc;
        }
        if (src.type != ValueType.UNKNOWN) {
            d.type = src.type;
        }
        d.raw |= src.raw;
        d.nullable |= src.nullable;
        d.allowEmpty |= src.allowEmpty;
        d.unique |= src.unique;
        if (src.defaultValue != null && !src.defaultValue.isEmpty()) {
            d.defaultValue = src.defaultValue;
        }
        if (src.min != null && !src.min.isEmpty()) {
            d.min = src.min;
        }
        if (src.max != null && !src.max.isEmpty()) {
            d.max = src.max;
        }
        if (src.size != 0) {
            d.size = src.size;
        }
        if (src.enumValues != null && !src.enumValues.isEmpty()) {
            d.enumValues = src.enumValues;
        }
        if (src.pattern != null && !src.pattern.isEmpty()) {
            d.pattern = src.pattern;
        }
        if (src.locationHours != 0) {
            d.locationHours = src.locationHours;
        }
        if (src.version != DEFAULT_VERSION) {
            d.version = src.version;
        }
        if (src.mime != null && !src.mime.isEmpty()) {
            d.mime = src.mime;
        }
        if (src.childDesc != null && !src.childDesc.isEmpty()) {
            d.childDesc = src.childDesc;
        }
        if (src.childType != ValueType.UNKNOWN) {
            d.childType = src.childType;
        }
        d.childRaw |= src.childRaw;
        d.childNullable |= src.childNullable;
        d.childAllowEmpty |= src.childAllowEmpty;
        d.childUnique |= src.childUnique;
        if (src.childDefault != null && !src.childDefault.isEmpty()) {
            d.childDefault = src.childDefault;
        }
        if (src.childMin != null && !src.childMin.isEmpty()) {
            d.childMin = src.childMin;
        }
        if (src.childMax != null && !src.childMax.isEmpty()) {
            d.childMax = src.childMax;
        }
        if (src.childSize != 0) {
            d.childSize = src.childSize;
        }
        if (src.childEnum != null && !src.childEnum.isEmpty()) {
            d.childEnum = src.childEnum;
        }
        if (src.childPattern != null && !src.childPattern.isEmpty()) {
            d.childPattern = src.childPattern;
        }
        if (src.childLocationHours != 0) {
            d.childLocationHours = src.childLocationHours;
        }
        if (src.childVersion != DEFAULT_VERSION) {
            d.childVersion = src.childVersion;
        }
        if (src.childMime != null && !src.childMime.isEmpty()) {
            d.childMime = src.childMime;
        }
        return d;
    }

    public MmTag copy() {
        MmTag o = new MmTag();
        o.name = name;
        o.isNull = isNull;
        o.example = example;
        o.desc = desc;
        o.type = type;
        o.raw = raw;
        o.nullable = nullable;
        o.allowEmpty = allowEmpty;
        o.unique = unique;
        o.defaultValue = defaultValue;
        o.min = min;
        o.max = max;
        o.size = size;
        o.enumValues = enumValues;
        o.pattern = pattern;
        o.locationHours = locationHours;
        o.version = version;
        o.mime = mime;
        o.childDesc = childDesc;
        o.childType = childType;
        o.childRaw = childRaw;
        o.childNullable = childNullable;
        o.childAllowEmpty = childAllowEmpty;
        o.childUnique = childUnique;
        o.childDefault = childDefault;
        o.childMin = childMin;
        o.childMax = childMax;
        o.childSize = childSize;
        o.childEnum = childEnum;
        o.childPattern = childPattern;
        o.childLocationHours = childLocationHours;
        o.childVersion = childVersion;
        o.childMime = childMime;
        o.isInherit = isInherit;
        return o;
    }

    public byte[] toBytes() {
        TagByteWriter w = new TagByteWriter();
        if (example) {
            w.writeByte((byte) (TagKey.K_EXAMPLE | 1));
        }
        if (isNull) {
            w.writeByte((byte) (TagKey.K_IS_NULL | 1));
        }
        if (nullable && !isInherit) {
            if (!isNull) {
                w.writeByte((byte) (TagKey.K_NULLABLE | 1));
            }
        }
        if (!desc.isEmpty() && !isInherit) {
            writeSizedString(w, TagKey.K_DESC, desc);
        }
        if (type != ValueType.UNKNOWN && !isInherit) {
            if (shouldEmitExplicitType()) {
                w.writeByte((byte) TagKey.K_TYPE);
                w.writeByte(type.code());
            }
        }
        if (raw && !isInherit) {
            w.writeByte((byte) (TagKey.K_RAW | 1));
        }
        if (allowEmpty && !isInherit) {
            w.writeByte((byte) (TagKey.K_ALLOW_EMPTY | 1));
        }
        if (unique && !isInherit) {
            w.writeByte((byte) (TagKey.K_UNIQUE | 1));
        }
        if (!defaultValue.isEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_DEFAULT, defaultValue);
        }
        if (!min.isEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_MIN, min);
        }
        if (!max.isEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_MAX, max);
        }
        if (size != 0 && !isInherit) {
            encodeUint64(w, TagKey.K_SIZE, size);
        }
        if (!enumValues.isEmpty() && !isInherit) {
            writeSizedString(w, TagKey.K_ENUM, enumValues);
        }
        if (!pattern.isEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_PATTERN, pattern);
        }
        if (locationHours != 0 && !isInherit) {
            String v = Integer.toString(locationHours);
            w.writeByte((byte) (TagKey.K_LOCATION | v.length()));
            w.writeAscii(v);
        }
        if (version != DEFAULT_VERSION && !isInherit) {
            encodeUint64(w, TagKey.K_VERSION, version);
        }
        if (!mime.isEmpty() && !isInherit) {
            int m = MimeWire.parse(mime);
            if (m < 7) {
                w.writeByte((byte) (TagKey.K_MIME | m));
            } else {
                w.writeByte((byte) (TagKey.K_MIME | 7));
                w.writeByte((byte) m);
            }
        }
        if (!childDesc.isEmpty()) {
            writeSizedString(w, TagKey.K_CHILD_DESC, childDesc);
        }
        if (childType != ValueType.UNKNOWN) {
            if (shouldEmitChildType()) {
                w.writeByte((byte) TagKey.K_CHILD_TYPE);
                w.writeByte(childType.code());
            }
        }
        if (childRaw) {
            w.writeByte((byte) (TagKey.K_CHILD_RAW | 1));
        }
        if (childNullable) {
            w.writeByte((byte) (TagKey.K_CHILD_NULLABLE | 1));
        }
        if (childAllowEmpty) {
            w.writeByte((byte) (TagKey.K_CHILD_ALLOW_EMPTY | 1));
        }
        if (childUnique) {
            w.writeByte((byte) (TagKey.K_CHILD_UNIQUE | 1));
        }
        if (!childDefault.isEmpty()) {
            writeShortString(w, TagKey.K_CHILD_DEFAULT, childDefault);
        }
        if (!childMin.isEmpty()) {
            writeShortString(w, TagKey.K_CHILD_MIN, childMin);
        }
        if (!childMax.isEmpty()) {
            writeShortString(w, TagKey.K_CHILD_MAX, childMax);
        }
        if (childSize != 0) {
            encodeUint64(w, TagKey.K_CHILD_SIZE, childSize);
        }
        if (!childEnum.isEmpty()) {
            writeSizedString(w, TagKey.K_CHILD_ENUM, childEnum);
        }
        if (!childPattern.isEmpty()) {
            writeShortString(w, TagKey.K_CHILD_PATTERN, childPattern);
        }
        if (childLocationHours != 0) {
            String v = Integer.toString(childLocationHours);
            w.writeByte((byte) (TagKey.K_CHILD_LOCATION | v.length()));
            w.writeAscii(v);
        }
        if (childVersion != DEFAULT_VERSION) {
            encodeUint64(w, TagKey.K_CHILD_VERSION, childVersion);
        }
        if (!childMime.isEmpty()) {
            int m = MimeWire.parse(childMime);
            if (m < 7) {
                w.writeByte((byte) (TagKey.K_CHILD_MIME | m));
            } else {
                w.writeByte((byte) (TagKey.K_CHILD_MIME | 7));
                w.writeByte((byte) m);
            }
        }
        return w.toByteArray();
    }

    private boolean shouldEmitExplicitType() {
        return switch (type) {
            case STRING, BYTES, INT, FLOAT64, BOOL, STRUCT, SLICE -> false;
            default -> {
                if (type == ValueType.ARRAY && size > 0) {
                    yield false;
                }
                if (type == ValueType.ENUM && !enumValues.isEmpty()) {
                    yield false;
                }
                yield true;
            }
        };
    }

    private boolean shouldEmitChildType() {
        return switch (childType) {
            case STRING, INT, FLOAT64, BOOL, STRUCT, SLICE -> false;
            default -> {
                if (childType == ValueType.ARRAY && childSize > 0) {
                    yield false;
                }
                if (childType == ValueType.ENUM && !childEnum.isEmpty()) {
                    yield false;
                }
                yield true;
            }
        };
    }

    private static void writeSizedString(TagByteWriter w, int key, String s) {
        byte[] b = s.getBytes(StandardCharsets.UTF_8);
        int l = b.length;
        if (l <= 5) {
            w.writeByte((byte) (key | l));
            w.writeBytes(b);
        } else if (l <= 0xFF) {
            w.writeByte((byte) (key | 6));
            w.writeByte((byte) l);
            w.writeBytes(b);
        } else if (l <= 0xFFFF) {
            w.writeByte((byte) (key | 7));
            w.writeByte((byte) (l >> 8));
            w.writeByte((byte) l);
            w.writeBytes(b);
        }
    }

    private static void writeShortString(TagByteWriter w, int key, String s) {
        byte[] b = s.getBytes(StandardCharsets.UTF_8);
        int l = b.length;
        if (l < 7) {
            w.writeByte((byte) (key | l));
            w.writeBytes(b);
        } else {
            w.writeByte((byte) (key | 7));
            w.writeByte((byte) l);
            w.writeBytes(b);
        }
    }

    private static void encodeUint64(TagByteWriter buf, int sign, long uv) {
        if (uv < 0) {
            throw new IllegalArgumentException("unsigned expected");
        }
        if (uv <= WireConstants.MAX_1) {
            buf.writeByte((byte) sign);
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_2) {
            buf.writeByte((byte) (sign | 1));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_3) {
            buf.writeByte((byte) (sign | 2));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_4) {
            buf.writeByte((byte) (sign | 3));
            buf.writeByte((byte) (uv >> 24));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_5) {
            buf.writeByte((byte) (sign | 4));
            buf.writeByte((byte) (uv >> 32));
            buf.writeByte((byte) (uv >> 24));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_6) {
            buf.writeByte((byte) (sign | 5));
            buf.writeByte((byte) (uv >> 40));
            buf.writeByte((byte) (uv >> 32));
            buf.writeByte((byte) (uv >> 24));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else if (uv <= WireConstants.MAX_7) {
            buf.writeByte((byte) (sign | 6));
            buf.writeByte((byte) (uv >> 48));
            buf.writeByte((byte) (uv >> 40));
            buf.writeByte((byte) (uv >> 32));
            buf.writeByte((byte) (uv >> 24));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        } else {
            buf.writeByte((byte) (sign | 7));
            buf.writeByte((byte) (uv >> 56));
            buf.writeByte((byte) (uv >> 48));
            buf.writeByte((byte) (uv >> 40));
            buf.writeByte((byte) (uv >> 32));
            buf.writeByte((byte) (uv >> 24));
            buf.writeByte((byte) (uv >> 16));
            buf.writeByte((byte) (uv >> 8));
            buf.writeByte((byte) uv);
        }
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof MmTag mmTag)) {
            return false;
        }
        return isNull == mmTag.isNull
                && example == mmTag.example
                && raw == mmTag.raw
                && nullable == mmTag.nullable
                && allowEmpty == mmTag.allowEmpty
                && unique == mmTag.unique
                && size == mmTag.size
                && locationHours == mmTag.locationHours
                && version == mmTag.version
                && childRaw == mmTag.childRaw
                && childNullable == mmTag.childNullable
                && childAllowEmpty == mmTag.childAllowEmpty
                && childUnique == mmTag.childUnique
                && childSize == mmTag.childSize
                && childLocationHours == mmTag.childLocationHours
                && childVersion == mmTag.childVersion
                && isInherit == mmTag.isInherit
                && Objects.equals(name, mmTag.name)
                && desc.equals(mmTag.desc)
                && type == mmTag.type
                && defaultValue.equals(mmTag.defaultValue)
                && min.equals(mmTag.min)
                && max.equals(mmTag.max)
                && enumValues.equals(mmTag.enumValues)
                && pattern.equals(mmTag.pattern)
                && mime.equals(mmTag.mime)
                && childDesc.equals(mmTag.childDesc)
                && childType == mmTag.childType
                && childDefault.equals(mmTag.childDefault)
                && childMin.equals(mmTag.childMin)
                && childMax.equals(mmTag.childMax)
                && childEnum.equals(mmTag.childEnum)
                && childPattern.equals(mmTag.childPattern)
                && childMime.equals(mmTag.childMime);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
                name,
                isNull,
                example,
                desc,
                type,
                raw,
                nullable,
                allowEmpty,
                unique,
                defaultValue,
                min,
                max,
                size,
                enumValues,
                pattern,
                locationHours,
                version,
                mime,
                childDesc,
                childType,
                childRaw,
                childNullable,
                childAllowEmpty,
                childUnique,
                childDefault,
                childMin,
                childMax,
                childSize,
                childEnum,
                childPattern,
                childLocationHours,
                childVersion,
                childMime,
                isInherit);
    }

    private static final class TagByteWriter {
        private byte[] buf = new byte[64];
        private int len;

        void writeByte(byte b) {
            ensure(1);
            buf[len++] = b;
        }

        void writeBytes(byte[] b) {
            ensure(b.length);
            System.arraycopy(b, 0, buf, len, b.length);
            len += b.length;
        }

        void writeAscii(String s) {
            byte[] b = s.getBytes(StandardCharsets.US_ASCII);
            writeBytes(b);
        }

        byte[] toByteArray() {
            return Arrays.copyOf(buf, len);
        }

        private void ensure(int n) {
            if (len + n > buf.length) {
                buf = Arrays.copyOf(buf, Math.max(buf.length * 2, len + n));
            }
        }
    }

    public static final class TagKey {
        public static final int K_IS_NULL = 0 << 3;
        public static final int K_EXAMPLE = 1 << 3;
        public static final int K_DESC = 2 << 3;
        public static final int K_TYPE = 3 << 3;
        public static final int K_RAW = 4 << 3;
        public static final int K_NULLABLE = 5 << 3;
        public static final int K_ALLOW_EMPTY = 6 << 3;
        public static final int K_UNIQUE = 7 << 3;
        public static final int K_DEFAULT = 8 << 3;
        public static final int K_MIN = 9 << 3;
        public static final int K_MAX = 10 << 3;
        public static final int K_SIZE = 11 << 3;
        public static final int K_ENUM = 12 << 3;
        public static final int K_PATTERN = 13 << 3;
        public static final int K_LOCATION = 14 << 3;
        public static final int K_VERSION = 15 << 3;
        public static final int K_MIME = 16 << 3;
        public static final int K_CHILD_DESC = 17 << 3;
        public static final int K_CHILD_TYPE = 18 << 3;
        public static final int K_CHILD_RAW = 19 << 3;
        public static final int K_CHILD_NULLABLE = 20 << 3;
        public static final int K_CHILD_ALLOW_EMPTY = 21 << 3;
        public static final int K_CHILD_UNIQUE = 22 << 3;
        public static final int K_CHILD_DEFAULT = 23 << 3;
        public static final int K_CHILD_MIN = 24 << 3;
        public static final int K_CHILD_MAX = 25 << 3;
        public static final int K_CHILD_SIZE = 26 << 3;
        public static final int K_CHILD_ENUM = 27 << 3;
        public static final int K_CHILD_PATTERN = 28 << 3;
        public static final int K_CHILD_LOCATION = 29 << 3;
        public static final int K_CHILD_VERSION = 30 << 3;
        public static final int K_CHILD_MIME = 31 << 3;

        private TagKey() {}
    }
}
