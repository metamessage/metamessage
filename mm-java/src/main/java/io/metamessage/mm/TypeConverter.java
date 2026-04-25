package io.metamessage.mm;

/**
 * 自定义类型转换器接口，用于在 Java 类型和 MM 类型之间进行转换。
 */
public interface TypeConverter {
    
    /**
     * 将 Java 对象转换为 MM 类型。
     * 
     * @param value Java 对象
     * @param tag MM 标签
     * @return 转换后的 MM 类型对象
     */
    Object toMmType(Object value, MmTag tag);
    
    /**
     * 将 MM 类型对象转换为 Java 类型。
     * 
     * @param value MM 类型对象
     * @param tag MM 标签
     * @param targetType 目标 Java 类型
     * @return 转换后的 Java 对象
     */
    Object fromMmType(Object value, MmTag tag, Class<?> targetType);
    
    /**
     * 检查此转换器是否支持指定的 Java 类型。
     * 
     * @param type Java 类型
     * @return 如果支持，返回 true；否则返回 false
     */
    boolean supports(Class<?> type);
}
