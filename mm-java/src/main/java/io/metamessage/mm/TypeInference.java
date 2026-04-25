package io.metamessage.mm;

import java.lang.reflect.Field;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;
import java.util.UUID;

final class TypeInference {
    private TypeInference() {}

    static MmTag forField(Field f) {
        MM ann = f.getAnnotation(MM.class);
        if (ann != null) {
            MmTag t = MmTag.fromAnnotation(ann);
            if (List.class.isAssignableFrom(f.getType()) && t.type == ValueType.UNKNOWN) {
                t.type = ValueType.SLICE;
            }
            return t;
        }
        MmTag t = MmTag.empty();
        Class<?> ft = f.getType();
        applyClass(t, ft);
        if (t.type == ValueType.SLICE && f.getGenericType() instanceof ParameterizedType pt
                && pt.getActualTypeArguments().length > 0
                && pt.getActualTypeArguments()[0] instanceof Class<?> ac) {
            t.childType = valueTypeForComponent(ac);
        }
        if (t.type == ValueType.SLICE && t.childType == ValueType.UNKNOWN) {
            t.childType = ValueType.STRING;
        }
        return t;
    }

    static ValueType valueTypeForComponent(Class<?> c) {
        MmTag t = MmTag.empty();
        applyClass(t, c);
        return t.type;
    }

    private static void applyClass(MmTag t, Class<?> ft) {
        if (ft.isPrimitive()) {
            t.type = primitiveType(ft);
        } else if (String.class.equals(ft)) {
            t.type = ValueType.STRING;
        } else if (Boolean.class.equals(ft) || boolean.class.equals(ft)) {
            t.type = ValueType.BOOL;
        } else if (Byte.class.equals(ft) || byte.class.equals(ft)) {
            t.type = ValueType.INT8;
        } else if (Short.class.equals(ft) || short.class.equals(ft)) {
            t.type = ValueType.INT16;
        } else if (Integer.class.equals(ft) || int.class.equals(ft)) {
            t.type = ValueType.INT;
        } else if (Long.class.equals(ft) || long.class.equals(ft)) {
            t.type = ValueType.INT64;
        } else if (Float.class.equals(ft) || float.class.equals(ft)) {
            t.type = ValueType.FLOAT32;
        } else if (Double.class.equals(ft) || double.class.equals(ft)) {
            t.type = ValueType.FLOAT64;
        } else if (byte[].class.equals(ft)) {
            t.type = ValueType.BYTES;
        } else if (BigInteger.class.equals(ft)) {
            t.type = ValueType.BIGINT;
        } else if (UUID.class.equals(ft)) {
            t.type = ValueType.UUID;
        } else if (LocalDateTime.class.equals(ft) || Instant.class.equals(ft)) {
            t.type = ValueType.DATETIME;
        } else if (LocalDate.class.equals(ft)) {
            t.type = ValueType.DATE;
        } else if (LocalTime.class.equals(ft)) {
            t.type = ValueType.TIME;
        } else if (List.class.isAssignableFrom(ft)) {
            t.type = ValueType.SLICE;
        } else if (java.util.Map.class.isAssignableFrom(ft)) {
            t.type = ValueType.MAP;
        } else {
            t.type = ValueType.STRUCT;
        }
    }

    private static ValueType primitiveType(Class<?> ft) {
        if (ft == int.class) {
            return ValueType.INT;
        }
        if (ft == long.class) {
            return ValueType.INT64;
        }
        if (ft == boolean.class) {
            return ValueType.BOOL;
        }
        if (ft == byte.class) {
            return ValueType.INT8;
        }
        if (ft == short.class) {
            return ValueType.INT16;
        }
        if (ft == float.class) {
            return ValueType.FLOAT32;
        }
        if (ft == double.class) {
            return ValueType.FLOAT64;
        }
        return ValueType.INT;
    }
}
