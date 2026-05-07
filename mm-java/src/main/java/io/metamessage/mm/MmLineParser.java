package io.github.metamessage.mm;

import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

/**
 * Parses {@code mm:} comment bodies (Go {@code ast.ParseMMTag}).
 */
public final class MmLineParser {

    private MmLineParser() {}

    public static MmTag parse(String input) {
        MmTag r = MmTag.empty();
        String tag = input.trim();
        if (tag.startsWith("//")) {
            tag = tag.substring(2).trim();
        }
        if (tag.startsWith("mm:")) {
            tag = tag.substring(3).trim();
        }
        if (tag.isEmpty()) {
            return r;
        }
        for (String part : splitTag(tag)) {
            String p = part.trim();
            if (p.isEmpty()) {
                continue;
            }
            String k;
            String v;
            int eq = p.indexOf('=');
            if (eq >= 0) {
                k = p.substring(0, eq).trim();
                v = p.substring(eq + 1).trim();
            } else {
                k = p.trim();
                v = "";
            }
            String lower = k.toLowerCase(Locale.ROOT);
            try {
                apply(r, lower, v);
            } catch (RuntimeException ignored) {
                // ignore unknown keys like Go default branch
            }
        }
        return r;
    }

    private static List<String> splitTag(String tag) {
        List<String> parts = new ArrayList<>();
        for (String s : tag.split(";")) {
            parts.add(s.trim());
        }
        return parts;
    }

    @SuppressWarnings("PMD.CyclomaticComplexity")
    private static void apply(MmTag r, String lower, String v) {
        switch (lower) {
            case "is_null" -> {
                r.isNull = true;
                r.nullable = true;
            }
            case "example" -> r.example = true;
            case "desc" -> r.desc = unquote(v);
            case "type" -> {
                ValueType t = ValueType.parseWireName(v);
                if (t == ValueType.UNKNOWN) {
                    throw new IllegalArgumentException("type");
                }
                r.type = t;
            }
            case "raw" -> r.raw = true;
            case "nullable" -> r.nullable = true;
            case "allow_empty" -> r.allowEmpty = true;
            case "unique" -> r.unique = true;
            case "default" -> r.defaultValue = v;
            case "pattern" -> r.pattern = v;
            case "min" -> r.min = v;
            case "max" -> r.max = v;
            case "size" -> r.size = (int) Long.parseLong(v);
            case "enum" -> {
                r.type = ValueType.ENUM;
                r.enumValues = v;
            }
            case "location" -> {
                int d = Integer.parseInt(v);
                r.locationHours = d;
            }
            case "version" -> {
                int d = Integer.parseInt(v);
                if (d < 1 || d > 10) {
                    throw new IllegalArgumentException("version");
                }
                r.version = d;
            }
            case "mime" -> r.mime = v;
            case "child_desc" -> r.childDesc = unquote(v);
            case "child_type" -> {
                ValueType t = ValueType.parseWireName(v);
                if (t == ValueType.UNKNOWN) {
                    throw new IllegalArgumentException("child_type");
                }
                r.childType = t;
            }
            case "child_raw" -> r.childRaw = true;
            case "child_nullable" -> r.childNullable = true;
            case "child_allow_empty" -> r.childAllowEmpty = true;
            case "child_unique" -> r.childUnique = true;
            case "child_default" -> r.childDefault = v;
            case "child_min" -> r.childMin = v;
            case "child_max" -> r.childMax = v;
            case "child_size" -> r.childSize = (int) Long.parseLong(v);
            case "child_enum" -> {
                r.childEnum = v;
                r.childType = ValueType.ENUM;
            }
            case "child_pattern" -> r.childPattern = v;
            case "child_location" -> r.childLocationHours = Integer.parseInt(v);
            case "child_version" -> {
                int d = Integer.parseInt(v);
                if (d < 1 || d > 10) {
                    throw new IllegalArgumentException("child_version");
                }
                r.childVersion = d;
            }
            case "child_mime" -> r.childMime = v;
            default -> {
            }
        }
    }

    private static String unquote(String v) {
        if (v.length() >= 2 && v.charAt(0) == '"' && v.charAt(v.length() - 1) == '"') {
            return v.substring(1, v.length() - 1).replace("\\\"", "\"");
        }
        return v;
    }
}
