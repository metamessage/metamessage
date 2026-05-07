package io.github.metamessage.mm;

import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

/**
 * Wire / schema value kinds. Ordinals must match Go {@code ast.ValueType} (internal/jsonc/ast/value.go).
 */
public enum ValueType {
    UNKNOWN,
    DOC,
    SLICE,
    ARRAY,
    STRUCT,
    MAP,
    STRING,
    BYTES,
    BOOL,
    INT,
    INT8,
    INT16,
    INT32,
    INT64,
    UINT,
    UINT8,
    UINT16,
    UINT32,
    UINT64,
    FLOAT32,
    FLOAT64,
    BIGINT,
    DATETIME,
    DATE,
    TIME,
    UUID,
    DECIMAL,
    IP,
    URL,
    EMAIL,
    ENUM,
    IMAGE,
    VIDEO;

    public byte code() {
        return (byte) ordinal();
    }

    public static ValueType fromCode(int b) {
        ValueType[] v = values();
        if (b < 0 || b >= v.length) {
            return UNKNOWN;
        }
        return v[b];
    }

    /** Same strings as Go {@code ValueType.String()} (e.g. {@code str}, {@code i}, {@code datetime}). */
    public String wireName() {
        return switch (this) {
            case UNKNOWN -> "unknown";
            case DOC -> "doc";
            case SLICE -> "slice";
            case ARRAY -> "arr";
            case STRUCT -> "obj";
            case MAP -> "map";
            case STRING -> "str";
            case BYTES -> "bytes";
            case BOOL -> "bool";
            case INT -> "i";
            case INT8 -> "i8";
            case INT16 -> "i16";
            case INT32 -> "i32";
            case INT64 -> "i64";
            case UINT -> "u";
            case UINT8 -> "u8";
            case UINT16 -> "u16";
            case UINT32 -> "u32";
            case UINT64 -> "u64";
            case FLOAT32 -> "f32";
            case FLOAT64 -> "f64";
            case BIGINT -> "bi";
            case DATETIME -> "datetime";
            case DATE -> "date";
            case TIME -> "time";
            case UUID -> "uuid";
            case DECIMAL -> "decimal";
            case IP -> "ip";
            case URL -> "url";
            case EMAIL -> "email";
            case ENUM -> "enum";
            case IMAGE -> "image";
            case VIDEO -> "video";
        };
    }

    private static final Map<String, ValueType> WIRE = new HashMap<>();

    static {
        for (ValueType v : values()) {
            WIRE.put(v.wireName(), v);
        }
    }

    public static ValueType parseWireName(String s) {
        if (s == null || s.isEmpty()) {
            return UNKNOWN;
        }
        return WIRE.getOrDefault(s.trim().toLowerCase(Locale.ROOT), UNKNOWN);
    }
}
