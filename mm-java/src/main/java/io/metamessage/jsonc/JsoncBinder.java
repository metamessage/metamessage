package io.github.metamessage.jsonc;

import io.github.metamessage.mm.MmTag;
import io.github.metamessage.mm.ValueType;
import java.lang.reflect.Array;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.math.BigInteger;
import java.net.URI;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

/**
 * Binds a {@link JcNode} tree into Java objects (Go {@code internal/jsonc/bind.go}).
 */
public final class JsoncBinder {
    private JsoncBinder() {}

    public static void bind(JcNode node, Object out) {
        if (out == null) {
            throw new JsoncException("out must be non-null");
        }
        bind0(node, out);
    }

    private static void bind0(JcNode node, Object out) {
        if (node instanceof JcNode.JcObject o) {
            ValueType vt = o.tag() != null ? o.tag().type : ValueType.STRUCT;
            if (vt == ValueType.MAP) {
                if (!(out instanceof Map)) {
                    throw new JsoncException("map node requires a java.util.Map target");
                }
                @SuppressWarnings("unchecked")
                Map<String, Object> m = (Map<String, Object>) out;
                bindMap(m, o);
            } else {
                bindStruct(o, out);
            }
        } else if (node instanceof JcNode.JcArray a) {
            if (out.getClass().isArray()) {
                bindAsJavaArray(a, out);
            } else if (out instanceof List) {
                @SuppressWarnings("unchecked")
                List<Object> list = (List<Object>) out;
                fillList(a, list, null);
            } else {
                throw new JsoncException("array binding target must be a Java array or java.util.List");
            }
        } else if (node instanceof JcNode.JcValue) {
            throw new JsoncException("root value binding is not supported");
        } else {
            throw new JsoncException("unsupported node: " + node);
        }
    }

    @SuppressWarnings("unchecked")
    private static void bindMap(Map<String, Object> m, JcNode.JcObject o) {
        if (!m.isEmpty()) {
            throw new JsoncException("target map must be empty");
        }
        for (JcField jf : o.fields()) {
            m.put(jf.key(), materializeNode(jf.value()));
        }
    }

    private static void bindStruct(JcNode.JcObject o, Object out) {
        Class<?> c = out.getClass();
        for (JcField jf : o.fields()) {
            Field f = findField(c, jf.key());
            if (f == null) {
                throw new JsoncException("struct has no field for key: " + jf.key());
            }
            f.setAccessible(true);
            if (Modifier.isStatic(f.getModifiers()) || Modifier.isFinal(f.getModifiers())) {
                throw new JsoncException("cannot set field: " + f.getName());
            }
            JcNode ch = jf.value();
            try {
                if (ch instanceof JcNode.JcValue v) {
                    setScalar(f, out, v);
                } else if (f.getType().isArray() && ch instanceof JcNode.JcArray) {
                    Object arr = newArray(f.getType().getComponentType(), ((JcNode.JcArray) ch).items().size());
                    bind0((JcNode.JcArray) ch, arr);
                    f.set(out, arr);
                } else if (List.class.isAssignableFrom(f.getType()) && ch instanceof JcNode.JcArray) {
                    @SuppressWarnings({"unchecked", "rawtypes"})
                    List list = new ArrayList();
                    Class<?> elem = listElementClass(f);
                    fillList((JcNode.JcArray) ch, list, elem);
                    f.set(out, list);
                } else if (Map.class.isAssignableFrom(f.getType()) && ch instanceof JcNode.JcObject) {
                    Map<String, Object> map = new HashMap<>();
                    bindMap(map, (JcNode.JcObject) ch);
                    f.set(out, map);
                } else {
                    Object nested;
                    if (f.getType().getPackage() == null) {
                        throw new JsoncException("unsupported field type: " + f.getType());
                    }
                    try {
                        nested = f.getType().getDeclaredConstructor().newInstance();
                    } catch (Exception e) {
                        throw new JsoncException("need default constructor: " + f.getType().getName(), e);
                    }
                    bind0(ch, nested);
                    f.set(out, nested);
                }
            } catch (IllegalAccessException e) {
                throw new JsoncException("cannot set: " + f.getName(), e);
            }
        }
    }

    private static Object newArray(Class<?> comp, int n) {
        if (comp.isPrimitive()) {
            return Array.newInstance(comp, n);
        }
        return Array.newInstance(comp, n);
    }

    private static Class<?> listElementClass(Field f) {
        if (!(f.getGenericType() instanceof ParameterizedType p)) {
            return Object.class;
        }
        Type e = p.getActualTypeArguments()[0];
        if (e instanceof Class<?> cl) {
            return cl;
        }
        return Object.class;
    }

    private static void fillList(JcNode.JcArray a, List<Object> out, Class<?> elem) {
        if (!out.isEmpty()) {
            out.clear();
        }
        for (JcNode it : a.items()) {
            if (it instanceof JcNode.JcValue) {
                out.add(materializeNode(it));
            } else if (elem == null || elem == Object.class) {
                out.add(materializeNode(it));
            } else {
                try {
                    Object nested = elem.getDeclaredConstructor().newInstance();
                    bind0(it, nested);
                    out.add(nested);
                } catch (Exception e) {
                    throw new JsoncException("list element: " + elem.getName(), e);
                }
            }
        }
    }

    private static void bindAsJavaArray(JcNode.JcArray a, Object out) {
        Class<?> comp = out.getClass().getComponentType();
        int alen = Array.getLength(out);
        List<JcNode> items = a.items();
        MmTag tag = a.tag();
        if (tag != null && tag.size > 0 && tag.size != alen) {
            throw new JsoncException("array length mismatch: target " + alen + " tag size " + tag.size);
        }
        if (items.size() != alen) {
            throw new JsoncException("item count " + items.size() + " expected " + alen);
        }
        for (int i = 0; i < alen; i++) {
            JcNode it = items.get(i);
            if (it instanceof JcNode.JcValue v) {
                setInArray(out, i, v, comp);
            } else {
                try {
                    Object el = comp.getDeclaredConstructor().newInstance();
                    Array.set(out, i, el);
                    bind0(it, el);
                } catch (Exception e) {
                    throw new JsoncException("array element " + i, e);
                }
            }
        }
    }

    private static void setInArray(Object arr, int i, JcNode.JcValue v, Class<?> comp) {
        MmTag tag = v.tag() != null ? v.tag() : MmTag.empty();
        if (tag.isNull) {
            if (comp.isPrimitive()) {
                throw new JsoncException("null in primitive array");
            }
            Array.set(arr, i, null);
            return;
        }
        Object d = v.data();
        if (comp == int.class) {
            Array.setInt(arr, i, ((Number) d).intValue());
        } else if (comp == long.class) {
            Array.setLong(arr, i, ((Number) d).longValue());
        } else if (comp == double.class) {
            Array.setDouble(arr, i, ((Number) d).doubleValue());
        } else if (comp == float.class) {
            Array.setFloat(arr, i, ((Number) d).floatValue());
        } else if (comp == boolean.class) {
            Array.setBoolean(arr, i, (Boolean) d);
        } else if (comp == byte.class) {
            Array.setByte(arr, i, ((Number) d).byteValue());
        } else if (comp == short.class) {
            Array.setShort(arr, i, ((Number) d).shortValue());
        } else {
            Array.set(arr, i, d);
        }
    }

    private static Field findField(Class<?> c, String jsonKey) {
        String n = toStructFieldName(jsonKey);
        for (Class<?> w = c; w != null && w != Object.class; w = w.getSuperclass()) {
            for (Field f : w.getDeclaredFields()) {
                if (f.getName().equals(n)) {
                    return f;
                }
            }
        }
        return null;
    }

    /**
     * Go struct field: upper only first char (e.g. {@code user_name} &rarr; {@code User_name}).
     */
    static String toStructFieldName(String key) {
        if (key == null || key.isEmpty()) {
            return "";
        }
        return Character.toUpperCase(key.charAt(0)) + key.substring(1);
    }

    private static Object materializeNode(JcNode n) {
        if (n instanceof JcNode.JcValue v) {
            MmTag t = v.tag();
            if (t != null && t.isNull) {
                return null;
            }
            return v.data();
        }
        if (n instanceof JcNode.JcObject o) {
            Map<String, Object> map = new HashMap<>();
            bindMap(map, o);
            return map;
        }
        if (n instanceof JcNode.JcArray a) {
            List<Object> list = new ArrayList<>();
            fillList(a, list, null);
            return list;
        }
        return null;
    }

    private static void setScalar(Field f, Object parent, JcNode.JcValue v) {
        MmTag tag = v.tag() != null ? v.tag() : MmTag.empty();
        if (tag.isNull) {
            try {
                if (f.getType().isPrimitive()) {
                    throw new JsoncException("null for primitive: " + f.getName());
                }
                f.set(parent, null);
            } catch (IllegalAccessException e) {
                throw new JsoncException(e.getMessage(), e);
            }
            return;
        }
        Object d = v.data();
        try {
            Class<?> t = f.getType();
            if (t == String.class) {
                f.set(parent, d != null ? d.toString() : v.text());
            } else if (t == int.class || t == Integer.class) {
                f.setInt(parent, ((Number) d).intValue());
            } else if (t == long.class || t == Long.class) {
                f.setLong(parent, ((Number) d).longValue());
            } else if (t == boolean.class || t == Boolean.class) {
                f.setBoolean(parent, (Boolean) d);
            } else if (t == double.class || t == Double.class) {
                f.setDouble(parent, ((Number) d).doubleValue());
            } else if (t == float.class || t == Float.class) {
                f.setFloat(parent, ((Number) d).floatValue());
            } else if (t == byte.class || t == Byte.class) {
                f.setByte(parent, ((Number) d).byteValue());
            } else if (t == short.class || t == Short.class) {
                f.setShort(parent, ((Number) d).shortValue());
            } else if (t == BigInteger.class) {
                f.set(parent, d);
            } else if (t == UUID.class) {
                f.set(parent, d instanceof UUID u ? u : UUID.fromString(v.text()));
            } else if (t == LocalDateTime.class) {
                f.set(parent, d);
            } else if (t == LocalDate.class) {
                f.set(parent, d);
            } else if (t == LocalTime.class) {
                f.set(parent, d);
            } else if (t == URI.class) {
                f.set(parent, d);
            } else if (t == byte[].class) {
                f.set(parent, d);
            } else {
                f.set(parent, d);
            }
        } catch (IllegalAccessException | ClassCastException e) {
            throw new JsoncException("cannot set field " + f.getName() + " for " + tag.type, e);
        }
    }
}
