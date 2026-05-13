package io.github.metamessage.mm

import java.math.BigInteger
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.util.UUID
import kotlin.reflect.KClass
import kotlin.reflect.KProperty
import io.github.metamessage.ast.Tag
import io.github.metamessage.ast.ValueType
import io.github.metamessage.MM


object TypeInference {
    fun forProperty(prop: KProperty<*>): ValueType {
        val returnType = prop.returnType
        val classifier = returnType.classifier
        if (classifier is KClass<*>) {
            return valueTypeForComponent(classifier.java)
        }
        return ValueType.UNKNOWN
    }

    fun forField(f: java.lang.reflect.Field): Tag {
        val ann = f.getAnnotation(MM::class.java)
        if (ann != null) {
            val t = Tag.fromAnnotation(ann)
            if (List::class.java.isAssignableFrom(f.type) && t.type == ValueType.UNKNOWN) {
                t.type = ValueType.SLICE
            }
            return t
        }
        val t = Tag.empty()
        applyClass(t, f.type)
        val gt = f.genericType
        if (gt is java.lang.reflect.ParameterizedType && gt.actualTypeArguments.isNotEmpty()) {
            val a0 = gt.actualTypeArguments[0]
            if (a0 is Class<*>) {
                if (t.type == ValueType.MAP) {
                    // For Map, first type arg is key, second is value
                    if (gt.actualTypeArguments.size > 1) {
                        val a1 = gt.actualTypeArguments[1]
                        if (a1 is Class<*>) {
                            t.childType = valueTypeForComponent(a1)
                        }
                    }
                } else {
                    t.childType = valueTypeForComponent(a0)
                }
            }
        }
        if (t.type == ValueType.SLICE && t.childType == ValueType.UNKNOWN) {
            t.childType = ValueType.STRING
        }
        if (t.type == ValueType.MAP && t.childType == ValueType.UNKNOWN) {
            t.childType = ValueType.STRING
        }
        return t
    }

    fun valueTypeForComponent(c: Class<*>): ValueType {
        val t = Tag.empty()
        applyClass(t, c)
        return t.type
    }

    private fun applyClass(t: Tag, ft: Class<*>) {
        t.type = when {
            // 原生类型
            ft.isPrimitive -> primitiveType(ft)

            // String
            ft == String::class.java -> ValueType.STRING

            // Boolean
            ft == Boolean::class.javaObjectType || ft == Boolean::class.javaPrimitiveType -> ValueType.BOOL

            // 整数
            ft == Byte::class.javaObjectType || ft == Byte::class.javaPrimitiveType -> ValueType.INT8
            ft == Short::class.javaObjectType || ft == Short::class.javaPrimitiveType -> ValueType.INT16
            ft == Int::class.javaObjectType || ft == Int::class.javaPrimitiveType -> ValueType.INT
            ft == Long::class.javaObjectType || ft == Long::class.javaPrimitiveType -> ValueType.INT64

            // 浮点数
            ft == Float::class.javaObjectType || ft == Float::class.javaPrimitiveType -> ValueType.FLOAT32
            ft == Double::class.javaObjectType || ft == Double::class.javaPrimitiveType -> ValueType.FLOAT64

            // 字节数组
            ft == ByteArray::class.java -> ValueType.BYTES

            // 大整数
            ft == BigInteger::class.java -> ValueType.BIGINT

            // UUID
            ft == UUID::class.java -> ValueType.UUID

            // 时间类型
            ft == LocalDateTime::class.java || ft == Instant::class.java -> ValueType.DATETIME
            ft == LocalDate::class.java -> ValueType.DATE
            ft == LocalTime::class.java -> ValueType.TIME

            // 集合
            List::class.java.isAssignableFrom(ft) -> ValueType.SLICE
            Map::class.java.isAssignableFrom(ft) -> ValueType.MAP

            // 默认结构体
            else -> ValueType.STRUCT
        }
    }

    private fun primitiveType(ft: Class<*>): ValueType {
        return when (ft) {
            Int::class.javaPrimitiveType -> ValueType.INT
            Long::class.javaPrimitiveType -> ValueType.INT64
            Boolean::class.javaPrimitiveType -> ValueType.BOOL
            Byte::class.javaPrimitiveType -> ValueType.INT8
            Short::class.javaPrimitiveType -> ValueType.INT16
            Float::class.javaPrimitiveType -> ValueType.FLOAT32
            Double::class.javaPrimitiveType -> ValueType.FLOAT64
            else -> ValueType.INT
        }
    }
}
