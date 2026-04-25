package io.metamessage.mm;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * MetaMessage field / type metadata, equivalent to Go struct tag {@code mm:"..."}.
 * Maps to the same binary tag payload as {@code ast.Tag.Bytes()} in the Go implementation.
 */
@Target({ElementType.FIELD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface MM {
    /**
     * When non-empty, overrides the wire key for this field (Go {@code name=} in mm tag).
     * Use {@code "-"} to skip encoding this field (reserved for future use).
     */
    String name() default "";

    ValueType type() default ValueType.UNKNOWN;

    boolean isNull() default false;

    boolean example() default false;

    String desc() default "";

    boolean raw() default false;

    boolean nullable() default false;

    boolean allowEmpty() default false;

    boolean unique() default false;

    /** Go {@code default=}. */
    String defaultValue() default "";

    String min() default "";

    String max() default "";

    int size() default 0;

    /** Pipe-separated enum labels (Go {@code enum=}). */
    String enumValues() default "";

    String pattern() default "";

    /** Timezone offset hours for datetime/date/time; 0 means UTC (omitted on wire when 0). */
    int location() default 0;

    /** UUID / IP wire variant; 0 means default. */
    int version() default 0;

    String mime() default "";

    String childDesc() default "";

    ValueType childType() default ValueType.UNKNOWN;

    boolean childRaw() default false;

    boolean childNullable() default false;

    boolean childAllowEmpty() default false;

    boolean childUnique() default false;

    String childDefault() default "";

    String childMin() default "";

    String childMax() default "";

    int childSize() default 0;

    String childEnum() default "";

    String childPattern() default "";

    int childLocation() default 0;

    int childVersion() default 0;

    String childMime() default "";
}
