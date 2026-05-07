package io.github.metamessage.mm

import java.math.BigInteger
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.util.UUID

object TypeInference {
    fun forField(f: java.lang.reflect.Field): MmTag {
        val ann = f.getAnnotation(MM::class.java)
        if (ann != null) {
            val t = MmTag.fromAnnotation(ann)
            if (List::class.java.isAssignableFrom(f.type) && t.type == ValueType.UNKNOWN) {
                t.type = ValueType.SLICE
            }
            return t
        }
        val t = MmTag.empty()
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
        val t = MmTag.empty()
        applyClass(t, c)
        return t.type
    }

    private fun applyClass(t: MmTag, ft: Class<*>) {
        if (ft.isPrimitive) {
            t.type = primitiveType(ft)
        } else if (String::class.java == ft) {
            t.type = ValueType.STRING
        } else if (Boolean::class.java == ft || Boolean::class.javaPrimitiveType == ft) {
            t.type = ValueType.BOOL
        } else if (Byte::class.java == ft || Byte::class.javaPrimitiveType == ft) {
            t.type = ValueType.INT8
        } else if (Short::class.java == ft || Short::class.javaPrimitiveType == ft) {
            t.type = ValueType.INT16
        } else if (Int::class.java == ft || Int::class.javaPrimitiveType == ft) {
            t.type = ValueType.INT
        } else if (Long::class.java == ft || Long::class.javaPrimitiveType == ft) {
            t.type = ValueType.INT64
        } else if (Float::class.java == ft || Float::class.javaPrimitiveType == ft) {
            t.type = ValueType.FLOAT32
        } else if (Double::class.java == ft || Double::class.javaPrimitiveType == ft) {
            t.type = ValueType.FLOAT64
        } else if (ByteArray::class.java == ft) {
            t.type = ValueType.BYTES
        } else if (BigInteger::class.java == ft) {
            t.type = ValueType.BIGINT
        } else if (UUID::class.java == ft) {
            t.type = ValueType.UUID
        } else if (LocalDateTime::class.java == ft || Instant::class.java == ft) {
            t.type = ValueType.DATETIME
        } else if (LocalDate::class.java == ft) {
            t.type = ValueType.DATE
        } else if (LocalTime::class.java == ft) {
            t.type = ValueType.TIME
        } else if (List::class.java.isAssignableFrom(ft)) {
            t.type = ValueType.SLICE
        } else if (Map::class.java.isAssignableFrom(ft)) {
            t.type = ValueType.MAP
        } else {
            t.type = ValueType.STRUCT
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
