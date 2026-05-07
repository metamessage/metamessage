package io.github.metamessage.mm;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.List;
import java.util.UUID;

final class ReflectMmEncoder {

    static byte[] encode(Object root) throws ReflectiveOperationException {
        WireEncoder enc = new WireEncoder();
        encodeValue(enc, root, rootTagForClass(root.getClass()));
        return enc.toByteArray();
    }

    private static MmTag rootTagForClass(Class<?> c) {
        MM ann = c.getAnnotation(MM.class);
        MmTag t = ann != null ? MmTag.fromAnnotation(ann) : MmTag.empty();
        if (t.type == ValueType.UNKNOWN) {
            t.type = ValueType.STRUCT;
        }
        if (t.name == null || t.name.isEmpty()) {
            t.name = CamelToSnake.convert(c.getSimpleName());
        }
        return t;
    }

    private static void encodeValue(WireEncoder enc, Object v, MmTag tag) throws ReflectiveOperationException {
        MmTag work = tag.copy();
        if (v == null) {
            if (!work.nullable && !work.isNull) {
                throw new IllegalArgumentException("null for non-nullable");
            }
            work.isNull = true;
        }
        if (work.type == ValueType.STRUCT) {
            if (v == null) {
                throw new IllegalArgumentException("null struct");
            }
            encodeStruct(enc, v, work);
            return;
        }
        if (work.type == ValueType.SLICE || work.type == ValueType.ARRAY) {
            encodeList(enc, (List<?>) v, work);
            return;
        }
        if (work.type == ValueType.MAP) {
            encodeObject(enc, (java.util.Map<?, ?>) v, work);
            return;
        }
        WireEncoder payload = new WireEncoder();
        encodeScalarPayload(payload, v, work);
        enc.encodeTaggedPayload(payload.toByteArray(), work.toBytes());
    }

    private static void encodeStruct(WireEncoder enc, Object o, MmTag objTag) throws ReflectiveOperationException {
        GrowableByteBuf keysPacked = new GrowableByteBuf();
        GrowableByteBuf valsPacked = new GrowableByteBuf();
        WireEncoder tmp = new WireEncoder();
        for (Field f : o.getClass().getDeclaredFields()) {
            if (Modifier.isStatic(f.getModifiers())) {
                continue;
            }
            f.setAccessible(true);
            MM mm = f.getAnnotation(MM.class);
            if (mm != null && "-".equals(mm.name())) {
                continue;
            }
            MmTag ft = TypeInference.forField(f);
            if (mm != null) {
                MmTag ann = MmTag.fromAnnotation(mm);
                if (ann.type != ValueType.UNKNOWN) {
                    ft.type = ann.type;
                }
                if (ann.childType != ValueType.UNKNOWN) {
                    ft.childType = ann.childType;
                }
                if (ann.name != null && !ann.name.isEmpty()) {
                    ft.name = ann.name;
                }
                mergeAnnotations(ft, ann);
            }
            String key = fieldKey(f, ft, mm);
            Object fv = f.get(o);
            tmp.reset();
            encodeValue(tmp, fv, ft);
            valsPacked.writeAll(tmp.toByteArray());
            tmp.reset();
            tmp.encodeString(key);
            keysPacked.writeAll(tmp.toByteArray());
        }
        tmp.reset();
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()));
        GrowableByteBuf mapBody = new GrowableByteBuf();
        mapBody.writeAll(tmp.toByteArray());
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()));
        tmp.reset();
        tmp.encodeObjectPayload(mapBody.copyRange(0, mapBody.length()));
        enc.encodeTaggedPayload(tmp.toByteArray(), objTag.toBytes());
    }

    private static void mergeAnnotations(MmTag dst, MmTag src) {
        if (src.desc != null && !src.desc.isEmpty()) {
            dst.desc = src.desc;
        }
        dst.nullable |= src.nullable;
        dst.raw |= src.raw;
        dst.allowEmpty |= src.allowEmpty;
        dst.unique |= src.unique;
        if (src.defaultValue != null && !src.defaultValue.isEmpty()) {
            dst.defaultValue = src.defaultValue;
        }
        if (src.enumValues != null && !src.enumValues.isEmpty()) {
            dst.enumValues = src.enumValues;
            dst.type = ValueType.ENUM;
        }
        dst.locationHours = src.locationHours;
        dst.version = src.version;
        if (src.mime != null && !src.mime.isEmpty()) {
            dst.mime = src.mime;
        }
        dst.childDesc = src.childDesc;
        if (src.childType != ValueType.UNKNOWN) {
            dst.childType = src.childType;
        }
        dst.childNullable |= src.childNullable;
        if (src.childEnum != null && !src.childEnum.isEmpty()) {
            dst.childEnum = src.childEnum;
            dst.childType = ValueType.ENUM;
        }
    }

    private static String fieldKey(Field f, MmTag ft, MM mm) {
        if (mm != null && mm.name() != null && !mm.name().isEmpty() && !"-".equals(mm.name())) {
            return mm.name();
        }
        if (ft.name != null && !ft.name.isEmpty() && !"-".equals(ft.name)) {
            return ft.name;
        }
        return CamelToSnake.convert(f.getName());
    }

    private static void encodeScalarPayload(WireEncoder w, Object v, MmTag tag) throws ReflectiveOperationException {
        if (tag.isNull && tag.type == ValueType.INT) {
            w.encodeSimple(SimpleValue.NULL_INT);
            return;
        }
        if (tag.isNull && tag.type == ValueType.STRING) {
            w.encodeSimple(SimpleValue.NULL_STRING);
            return;
        }
        if (tag.isNull && tag.type == ValueType.BYTES) {
            w.encodeSimple(SimpleValue.NULL_BYTES);
            return;
        }
        if (tag.isNull && (tag.type == ValueType.FLOAT32 || tag.type == ValueType.FLOAT64)) {
            w.encodeSimple(SimpleValue.NULL_FLOAT);
            return;
        }
        
        // 使用类型转换器将 Java 对象转换为 MM 类型
        Object mmValue = TypeConverterRegistry.toMmType(v, tag);
        
        switch (tag.type) {
            case BOOL -> w.encodeBool(Boolean.TRUE.equals(mmValue));
            case INT -> w.encodeInt64(((Number) mmValue).intValue());
            case INT8 -> w.encodeInt64(((Number) mmValue).byteValue());
            case INT16 -> w.encodeInt64(((Number) mmValue).shortValue());
            case INT32, UINT, UINT16 -> w.encodeInt64(((Number) mmValue).intValue());
            case INT64, UINT32, UINT64 -> w.encodeInt64(((Number) mmValue).longValue());
            case FLOAT32 -> w.encodeFloatString(Float.toString(((Number) mmValue).floatValue()));
            case FLOAT64, DECIMAL -> w.encodeFloatString(Double.toString(((Number) mmValue).doubleValue()));
            case STRING, EMAIL, URL -> w.encodeString(mmValue == null ? "" : mmValue.toString());
            case BYTES -> w.encodeBytes(mmValue == null ? new byte[0] : (byte[]) mmValue);
            case BIGINT -> w.encodeBigIntDecimal(mmValue == null ? "0" : ((BigInteger) mmValue).toString());
            case UUID -> w.encodeBytes(uuidBytes((UUID) mmValue));
            case DATETIME -> {
                long sec;
                if (mmValue instanceof LocalDateTime ldt) {
                    sec = ldt.toEpochSecond(ZoneOffset.UTC);
                } else if (mmValue instanceof Instant ins) {
                    sec = ins.getEpochSecond();
                } else {
                    throw new IllegalArgumentException("datetime");
                }
                w.encodeInt64(sec);
            }
            case DATE -> w.encodeInt64(TimeUtil.daysSinceEpochUtc((LocalDate) mmValue));
            case TIME -> w.encodeInt64(((LocalTime) mmValue).toSecondOfDay());
            case ENUM -> w.encodeInt64(((Number) mmValue).intValue());
            default -> throw new UnsupportedOperationException("scalar " + tag.type);
        }
    }

    private static void encodeList(WireEncoder enc, List<?> list, MmTag tag) throws ReflectiveOperationException {
        if (list == null) {
            MmTag nt = tag.copy();
            nt.isNull = true;
            enc.encodeTaggedPayload(new byte[0], nt.toBytes());
            return;
        }
        GrowableByteBuf body = new GrowableByteBuf();
        WireEncoder el = new WireEncoder();
        for (Object item : list) {
            el.reset();
            MmTag et = MmTag.empty();
            et.inheritFromArrayParent(tag);
            if (et.type == ValueType.UNKNOWN && item != null) {
                et.type = TypeInference.valueTypeForComponent(item.getClass());
                if (et.type == ValueType.UNKNOWN) {
                    et.type = ValueType.STRUCT;
                }
            }
            encodeValue(el, item, et);
            body.writeAll(el.toByteArray());
        }
        el.reset();
        el.encodeArrayPayload(body.copyRange(0, body.length()));
        enc.encodeTaggedPayload(el.toByteArray(), tag.toBytes());
    }
    
    private static void encodeObject(WireEncoder enc, java.util.Map<?, ?> map, MmTag tag) throws ReflectiveOperationException {
        if (map == null) {
            MmTag nt = tag.copy();
            nt.isNull = true;
            enc.encodeTaggedPayload(new byte[0], nt.toBytes());
            return;
        }
        GrowableByteBuf keysPacked = new GrowableByteBuf();
        GrowableByteBuf valsPacked = new GrowableByteBuf();
        WireEncoder tmp = new WireEncoder();
        for (java.util.Map.Entry<?, ?> entry : map.entrySet()) {
            Object key = entry.getKey();
            Object value = entry.getValue();
            
            // 编码键
            tmp.reset();
            MmTag keyTag = MmTag.empty();
            keyTag.type = ValueType.STRING; // 假设键是字符串类型
            encodeValue(tmp, key, keyTag);
            keysPacked.writeAll(tmp.toByteArray());
            
            // 编码值
            tmp.reset();
            MmTag valueTag = MmTag.empty();
            valueTag.inheritFromArrayParent(tag);
            if (valueTag.type == ValueType.UNKNOWN && value != null) {
                valueTag.type = TypeInference.valueTypeForComponent(value.getClass());
                if (valueTag.type == ValueType.UNKNOWN) {
                    valueTag.type = ValueType.STRUCT;
                }
            }
            encodeValue(tmp, value, valueTag);
            valsPacked.writeAll(tmp.toByteArray());
        }
        tmp.reset();
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()));
        GrowableByteBuf mapBody = new GrowableByteBuf();
        mapBody.writeAll(tmp.toByteArray());
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()));
        tmp.reset();
        tmp.encodeObjectPayload(mapBody.copyRange(0, mapBody.length()));
        enc.encodeTaggedPayload(tmp.toByteArray(), tag.toBytes());
    }

    private static byte[] uuidBytes(UUID u) {
        java.nio.ByteBuffer bb = java.nio.ByteBuffer.allocate(16);
        bb.putLong(u.getMostSignificantBits());
        bb.putLong(u.getLeastSignificantBits());
        return bb.array();
    }
}
