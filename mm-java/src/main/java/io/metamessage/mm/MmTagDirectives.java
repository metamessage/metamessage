package io.github.metamessage.mm;

/**
 * Renders an {@link MmTag} as a semicolon-separated directive list (Go {@code (*Tag).String()}).
 */
public final class MmTagDirectives {

    private MmTagDirectives() {}

    public static String format(MmTag t) {
        if (t == null) {
            return "";
        }
        StringBuilder b = new StringBuilder();
        boolean first = true;
        if (t.type != ValueType.UNKNOWN && !t.isInherit) {
            if (shouldEmitType(t)) {
                add(b, first, "type=" + t.type.wireName());
                first = false;
            }
        }
        if (t.example) {
            add(b, first, "example");
            first = false;
        }
        if (t.isNull) {
            add(b, first, "is_null");
            first = false;
        }
        if (t.nullable && !t.isInherit && !t.isNull) {
            add(b, first, "nullable");
            first = false;
        }
        if (t.desc != null && !t.desc.isEmpty() && !t.isInherit) {
            add(b, first, "desc=" + quoteIfNeeded(t.desc));
            first = false;
        }
        if (t.raw && !t.isInherit) {
            add(b, first, "raw");
            first = false;
        }
        if (t.allowEmpty && !t.isInherit) {
            add(b, first, "allow_empty");
            first = false;
        }
        if (t.unique && !t.isInherit) {
            add(b, first, "unique");
            first = false;
        }
        if (t.defaultValue != null && !t.defaultValue.isEmpty() && !t.isInherit) {
            add(b, first, "default=" + t.defaultValue);
            first = false;
        }
        if (t.min != null && !t.min.isEmpty() && !t.isInherit) {
            add(b, first, "min=" + t.min);
            first = false;
        }
        if (t.max != null && !t.max.isEmpty() && !t.isInherit) {
            add(b, first, "max=" + t.max);
            first = false;
        }
        if (t.size != 0 && !t.isInherit) {
            add(b, first, "size=" + t.size);
            first = false;
        }
        if (t.enumValues != null && !t.enumValues.isEmpty() && !t.isInherit) {
            add(b, first, "enum=" + t.enumValues);
            first = false;
        }
        if (t.pattern != null && !t.pattern.isEmpty() && !t.isInherit) {
            add(b, first, "pattern=" + t.pattern);
            first = false;
        }
        if (t.locationHours != 0 && !t.isInherit) {
            add(b, first, "location=" + t.locationHours);
            first = false;
        }
        if (t.version != MmTag.DEFAULT_VERSION && !t.isInherit) {
            add(b, first, "version=" + t.version);
            first = false;
        }
        if (t.mime != null && !t.mime.isEmpty() && !t.isInherit) {
            add(b, first, "mime=" + t.mime);
            first = false;
        }
        if (t.childDesc != null && !t.childDesc.isEmpty()) {
            add(b, first, "child_desc=" + quoteIfNeeded(t.childDesc));
            first = false;
        }
        if (t.childType != ValueType.UNKNOWN) {
            if (shouldEmitChildType(t)) {
                add(b, first, "child_type=" + t.childType.wireName());
                first = false;
            }
        }
        if (t.childRaw) {
            add(b, first, "child_raw");
            first = false;
        }
        if (t.childNullable) {
            add(b, first, "child_nullable");
            first = false;
        }
        if (t.childAllowEmpty) {
            add(b, first, "child_allow_empty");
            first = false;
        }
        if (t.childUnique) {
            add(b, first, "child_unique");
            first = false;
        }
        if (t.childDefault != null && !t.childDefault.isEmpty()) {
            add(b, first, "child_default=" + t.childDefault);
            first = false;
        }
        if (t.childMin != null && !t.childMin.isEmpty()) {
            add(b, first, "child_min=" + t.childMin);
            first = false;
        }
        if (t.childMax != null && !t.childMax.isEmpty()) {
            add(b, first, "child_max=" + t.childMax);
            first = false;
        }
        if (t.childSize != 0) {
            add(b, first, "child_size=" + t.childSize);
            first = false;
        }
        if (t.childEnum != null && !t.childEnum.isEmpty()) {
            add(b, first, "child_enum=" + t.childEnum);
            first = false;
        }
        if (t.childPattern != null && !t.childPattern.isEmpty()) {
            add(b, first, "child_pattern=" + t.childPattern);
            first = false;
        }
        if (t.childLocationHours != 0) {
            add(b, first, "child_location=" + t.childLocationHours);
            first = false;
        }
        if (t.childVersion != MmTag.DEFAULT_VERSION) {
            add(b, first, "child_version=" + t.childVersion);
            first = false;
        }
        if (t.childMime != null && !t.childMime.isEmpty()) {
            add(b, first, "child_mime=" + t.childMime);
        }
        return b.toString();
    }

    private static boolean shouldEmitType(MmTag t) {
        return switch (t.type) {
            case STRING, BYTES, INT, FLOAT64, BOOL, STRUCT, SLICE -> false;
            case ARRAY -> t.size == 0;
            case ENUM -> t.enumValues == null || t.enumValues.isEmpty();
            default -> true;
        };
    }

    private static boolean shouldEmitChildType(MmTag t) {
        return switch (t.childType) {
            case STRING, INT, FLOAT64, BOOL, STRUCT, SLICE -> false;
            case ARRAY -> t.childSize == 0;
            case ENUM -> t.childEnum == null || t.childEnum.isEmpty();
            default -> true;
        };
    }

    private static void add(StringBuilder b, boolean first, String s) {
        if (!first) {
            b.append("; ");
        }
        b.append(s);
    }

    private static String quoteIfNeeded(String s) {
        if (s.contains(";") || s.contains(" ") || s.contains("=")) {
            return '"' + s.replace("\"", "\\\"") + '"';
        }
        return s;
    }
}
