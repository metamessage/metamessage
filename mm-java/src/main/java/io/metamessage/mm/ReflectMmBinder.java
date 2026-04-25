package io.metamessage.mm;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.math.BigInteger;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.UUID;

final class ReflectMmBinder {

    static <T> T bind(MmTree tree, Class<T> clazz) throws ReflectiveOperationException {
        if (!(tree instanceof MmTree.MmObject obj)) {
            throw new IllegalArgumentException("root must be object");
        }
        T inst = clazz.getDeclaredConstructor().newInstance();
        Map<String, Field> byKey = fieldsByJsonKey(clazz);
        for (Map.Entry<String, MmTree> e : obj.fields()) {
            Field f = byKey.get(e.getKey());
            if (f == null) {
                continue;
            }
            f.setAccessible(true);
            f.set(inst, materialize(f, e.getValue()));
        }
        return inst;
    }

    private static Map<String, Field> fieldsByJsonKey(Class<?> c) {
        Map<String, Field> m = new java.util.LinkedHashMap<>();
        for (Field f : c.getDeclaredFields()) {
            if (Modifier.isStatic(f.getModifiers())) {
                continue;
            }
            MmTag t = TypeInference.forField(f);
            MM mm = f.getAnnotation(MM.class);
            if (mm != null && "-".equals(mm.name())) {
                continue;
            }
            String key = fieldKey(f, t, mm);
            m.put(key, f);
        }
        return m;
    }

    private static String fieldKey(Field f, MmTag ft, MM mm) {
        if (mm != null && mm.name() != null && !mm.name().isEmpty() && !"-".equals(mm.name())) {
            return mm.name();
        }
        if (ft.name != null && !ft.name.isEmpty()) {
            return ft.name;
        }
        return CamelToSnake.convert(f.getName());
    }

    private static Object materialize(Field f, MmTree node) throws ReflectiveOperationException {
        Class<?> ft = f.getType();
        if (node instanceof MmTree.MmScalar sc) {
            return scalarToField(ft, sc);
        }
        if (node instanceof MmTree.MmArray arr) {
            return listFrom(arr, f);
        }
        if (node instanceof MmTree.MmObject o) {
            if (java.util.Map.class.isAssignableFrom(ft)) {
                return mapFrom(o, f);
            }
            return bind(o, ft);
        }
        throw new IllegalStateException("node " + node);
    }

    @SuppressWarnings("unchecked")
    private static Object listFrom(MmTree.MmArray arr, Field f) throws ReflectiveOperationException {
        Class<?> raw = List.class;
        java.lang.reflect.Type gt = f.getGenericType();
        Class<?> elemClass = String.class;
        if (gt instanceof java.lang.reflect.ParameterizedType pt && pt.getActualTypeArguments().length > 0) {
            java.lang.reflect.Type a0 = pt.getActualTypeArguments()[0];
            if (a0 instanceof Class<?> ac) {
                elemClass = ac;
            }
        }
        List<Object> out = new ArrayList<>();
        for (MmTree ch : arr.items()) {
            if (ch instanceof MmTree.MmScalar sc) {
                out.add(scalarToField(elemClass, sc));
            } else if (ch instanceof MmTree.MmObject mo) {
                Object sub = elemClass.getDeclaredConstructor().newInstance();
                Map<String, Field> byKey = fieldsByJsonKey(elemClass);
                for (Map.Entry<String, MmTree> e : mo.fields()) {
                    Field sf = byKey.get(e.getKey());
                    if (sf != null) {
                        sf.setAccessible(true);
                        sf.set(sub, materialize(sf, e.getValue()));
                    }
                }
                out.add(sub);
            }
        }
        return out;
    }
    
    @SuppressWarnings("unchecked")
    private static Object mapFrom(MmTree.MmObject obj, Field f) throws ReflectiveOperationException {
        Class<?> raw = java.util.Map.class;
        java.lang.reflect.Type gt = f.getGenericType();
        Class<?> keyClass = String.class;
        Class<?> valueClass = Object.class;
        if (gt instanceof java.lang.reflect.ParameterizedType pt && pt.getActualTypeArguments().length > 2) {
            java.lang.reflect.Type a0 = pt.getActualTypeArguments()[0];
            java.lang.reflect.Type a1 = pt.getActualTypeArguments()[1];
            if (a0 instanceof Class<?> ac0) {
                keyClass = ac0;
            }
            if (a1 instanceof Class<?> ac1) {
                valueClass = ac1;
            }
        }
        java.util.Map<Object, Object> out = new java.util.HashMap<>();
        for (Map.Entry<String, MmTree> e : obj.fields()) {
            String key = e.getKey();
            MmTree valueNode = e.getValue();
            Object value;
            if (valueNode instanceof MmTree.MmScalar sc) {
                value = scalarToField(valueClass, sc);
            } else if (valueNode instanceof MmTree.MmObject mo) {
                value = bind(mo, valueClass);
            } else if (valueNode instanceof MmTree.MmArray ma) {
                // 处理值为数组的情况
                value = listFrom(ma, f);
            } else {
                value = null;
            }
            out.put(key, value);
        }
        return out;
    }

    private static Object scalarToField(Class<?> ft, MmTree.MmScalar sc) {
        Object d = sc.data();
        
        // 使用类型转换器将 MM 类型转换为 Java 类型
        Object javaValue = TypeConverterRegistry.fromMmType(d, sc.tag(), ft);
        
        if (ft == int.class || ft == Integer.class) {
            return javaValue instanceof Number n ? n.intValue() : Integer.parseInt(sc.text());
        }
        if (ft == long.class || ft == Long.class) {
            return javaValue instanceof Number n ? n.longValue() : Long.parseLong(sc.text());
        }
        if (ft == boolean.class || ft == Boolean.class) {
            return Boolean.TRUE.equals(javaValue);
        }
        if (ft == String.class) {
            return sc.text();
        }
        if (ft == double.class || ft == Double.class) {
            return javaValue instanceof Number n ? n.doubleValue() : Double.parseDouble(sc.text());
        }
        if (ft == float.class || ft == Float.class) {
            return javaValue instanceof Number n ? n.floatValue() : Float.parseFloat(sc.text());
        }
        if (ft == LocalDateTime.class) {
            if (javaValue instanceof LocalDateTime ldt) {
                return ldt;
            }
            return LocalDateTime.parse(sc.text().replace(' ', 'T'));
        }
        if (ft == UUID.class && javaValue instanceof UUID u) {
            return u;
        }
        if (ft == BigInteger.class && javaValue instanceof BigInteger bi) {
            return bi;
        }
        return javaValue;
    }
}
