package io.metamessage.mm;

import java.util.ArrayList;
import java.util.List;

/**
 * 类型转换器注册表，用于管理所有的类型转换器。
 */
public class TypeConverterRegistry {
    
    private static final List<TypeConverter> converters = new ArrayList<>();
    
    static {
        // 注册默认的类型转换器
        // 这里可以添加默认的类型转换器
    }
    
    /**
     * 注册一个类型转换器。
     * 
     * @param converter 类型转换器
     */
    public static void register(TypeConverter converter) {
        converters.add(converter);
    }
    
    /**
     * 查找支持指定 Java 类型的类型转换器。
     * 
     * @param type Java 类型
     * @return 类型转换器，如果没有找到，返回 null
     */
    public static TypeConverter find(Class<?> type) {
        for (TypeConverter converter : converters) {
            if (converter.supports(type)) {
                return converter;
            }
        }
        return null;
    }
    
    /**
     * 将 Java 对象转换为 MM 类型。
     * 
     * @param value Java 对象
     * @param tag MM 标签
     * @return 转换后的 MM 类型对象
     */
    public static Object toMmType(Object value, MmTag tag) {
        if (value == null) {
            return null;
        }
        TypeConverter converter = find(value.getClass());
        if (converter != null) {
            return converter.toMmType(value, tag);
        }
        return value;
    }
    
    /**
     * 将 MM 类型对象转换为 Java 类型。
     * 
     * @param value MM 类型对象
     * @param tag MM 标签
     * @param targetType 目标 Java 类型
     * @return 转换后的 Java 对象
     */
    public static Object fromMmType(Object value, MmTag tag, Class<?> targetType) {
        if (value == null) {
            return null;
        }
        TypeConverter converter = find(targetType);
        if (converter != null) {
            return converter.fromMmType(value, tag, targetType);
        }
        return value;
    }
}
