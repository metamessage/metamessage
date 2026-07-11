package io.github.metamessage.core

import io.github.metamessage.ir.Node
import io.github.metamessage.ir.NodeArray
import io.github.metamessage.ir.NodeNull
import io.github.metamessage.ir.NodeObject
import io.github.metamessage.ir.NodeScalar
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.ValueType
import java.math.BigInteger
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.util.UUID

object Binder {
    @Suppress("UNCHECKED_CAST")
    fun <T> bind(node: Node, clazz: Class<T>): T {
        return bind(node, clazz, null)
    }

    @Suppress("UNCHECKED_CAST")
    fun <T> bind(node: Node, clazz: Class<T>, elementClass: Class<*>? = null): T {
        when (node) {
            is NodeNull -> {
                @Suppress("UNCHECKED_CAST") return null as T
            }
            is NodeObject -> {
                val tag = node.tag
                if (tag != null && tag.type == ValueType.OBJ) {
                    val inst = clazz.getDeclaredConstructor().newInstance()
                    convertObj(node, inst as Any)
                    return inst
                } else {
                    @Suppress("UNCHECKED_CAST") val inst = mutableMapOf<String, Any?>()
                    convertMap(node, inst)
                    return inst as T
                }
            }
            is NodeArray -> {
                val tag = node.tag
                if (tag != null && tag.size > 0 && tag.type == ValueType.ARR) {
                    @Suppress("UNCHECKED_CAST") return convertArr(node, elementClass ?: clazz) as T
                } else {
                    if (elementClass != null) {
                        @Suppress("UNCHECKED_CAST") return convertVec(node, elementClass) as T
                    }
                    @Suppress("UNCHECKED_CAST") return convertVecDynamic(node) as T
                }
            }
            is NodeScalar -> {
                return convertScalar(node, clazz)
            }
            else ->
                    throw IllegalArgumentException(
                            "unsupported node type: ${node::class.java.name}"
                    )
        }
    }

    private fun convertObj(obj: NodeObject, out: Any) {
        val outClazz = out.javaClass
        val nameToField = mutableMapOf<String, java.lang.reflect.Field>()
        for (f in outClazz.declaredFields) {
            if (java.lang.reflect.Modifier.isStatic(f.modifiers)) continue
            nameToField[f.name] = f
        }

        for (field in obj.fields) {
            val fieldKey = field.key
            val camelName = SnakeToCamel.convert(fieldKey)
            val structField = nameToField[camelName] ?: continue
            structField.isAccessible = true
            try {
                val fieldVal = materialize(structField, field.value)
                structField.set(out, fieldVal)
            } catch (e: Exception) {
                throw RuntimeException("failed to bind field ${field.key}: ${e.message}", e)
            }
        }
    }

    private fun convertMap(obj: NodeObject, out: MutableMap<String, Any?>) {
        for (field in obj.fields) {
            val key = field.key
            val value =
                    when (val v = field.value) {
                        is NodeNull -> null
                        is NodeScalar -> convertScalarToAny(v)
                        is NodeObject -> {
                            val map = mutableMapOf<String, Any?>()
                            convertMap(v, map)
                            map
                        }
                        is NodeArray -> {
                            convertVecDynamic(v)
                        }
                        else -> null
                    }
            out[key] = value
        }
    }

    private fun convertArr(arr: NodeArray, elementClass: Class<*>): Any {
        val list = mutableListOf<Any?>()
        for (item in arr.items) {
            when (item) {
                is NodeScalar -> list.add(convertScalarToAny(item))
                is NodeObject -> {
                    val inst = elementClass.getDeclaredConstructor().newInstance()
                    convertObj(item, inst as Any)
                    list.add(inst)
                }
                else -> list.add(null)
            }
        }
        return list
    }

    @Suppress("UNCHECKED_CAST")
    private fun convertVec(arr: NodeArray, field: java.lang.reflect.Field): Any {
        val elementClass = extractElementType(field)
        return convertVec(arr, elementClass)
    }

    @Suppress("UNCHECKED_CAST")
    private fun convertVec(arr: NodeArray, elementClass: Class<*>): Any {
        val list = mutableListOf<Any?>()
        for (item in arr.items) {
            when (item) {
                is NodeScalar -> list.add(convertScalarToAny(item))
                is NodeObject -> {
                    val inst = elementClass.getDeclaredConstructor().newInstance()
                    convertObj(item, inst as Any)
                    list.add(inst)
                }
                else -> list.add(null)
            }
        }
        return list
    }

    private fun extractElementType(field: java.lang.reflect.Field): Class<*> {
        val genericType = field.genericType
        if (genericType is java.lang.reflect.ParameterizedType) {
            val typeArgs = genericType.actualTypeArguments
            if (typeArgs.isNotEmpty()) {
                val arg = typeArgs[typeArgs.size - 1]
                if (arg is Class<*>) {
                    return arg
                }
                if (arg is java.lang.reflect.ParameterizedType) {
                    val rawType = arg.rawType
                    if (rawType is Class<*>) {
                        return rawType
                    }
                }
            }
        }
        return Any::class.java
    }

    private fun convertVecDynamic(arr: NodeArray): List<Any?> {
        val list = mutableListOf<Any?>()
        for (item in arr.items) {
            when (item) {
                is NodeScalar -> list.add(convertScalarToAny(item))
                is NodeObject -> {
                    val map = mutableMapOf<String, Any?>()
                    convertMap(item, map)
                    list.add(map)
                }
                else -> list.add(null)
            }
        }
        return list
    }

    @Suppress("UNCHECKED_CAST", "UNUSED_PARAMETER")
    private fun <T> convertScalar(value: NodeScalar, clazz: Class<T>): T {
        val tag = value.tag ?: Tag.empty()
        val data = value.data
        val text = value.text

        if (tag.isNull) {
            @Suppress("UNCHECKED_CAST") return null as T
        }

        return when (tag.type) {
            ValueType.DATETIME -> {
                (data as? LocalDateTime ?: LocalDateTime.of(1970, 1, 1, 0, 0, 0)) as T
            }
            ValueType.DATE -> {
                when (data) {
                    is LocalDate -> data as T
                    is LocalDateTime ->
                            (if (clazz == LocalDate::class.java) data.toLocalDate() else data) as T
                    else -> LocalDate.of(1970, 1, 1) as T
                }
            }
            ValueType.TIME -> {
                when (data) {
                    is LocalTime -> data as T
                    is LocalDateTime ->
                            (if (clazz == LocalTime::class.java) data.toLocalTime() else data) as T
                    else -> LocalTime.of(0, 0, 0) as T
                }
            }
            ValueType.BIGINT -> {
                (data as? BigInteger ?: BigInteger.ZERO) as T
            }
            ValueType.UUID -> {
                (data as? UUID ?: UUID(0, 0)).toString() as T
            }
            ValueType.DECIMAL, ValueType.EMAIL -> {
                (data as? String ?: text) as T
            }
            ValueType.IP -> {
                when (clazz) {
                    String::class.java ->
                            ((data as? java.net.InetAddress)?.hostAddress ?: text) as T
                    else -> (data as? java.net.InetAddress) as T
                }
            }
            ValueType.URL -> {
                when (clazz) {
                    String::class.java -> ((data as? java.net.URI)?.toString() ?: text) as T
                    else -> (data as? java.net.URI) as T
                }
            }
            ValueType.ENUMS -> {
                text as T
            }
            ValueType.I -> {
                val num = (data as? Number) ?: text.toInt()
                convertNumberToType(num, clazz) as T
            }
            ValueType.I8 -> {
                val num = (data as? Number) ?: text.toByte()
                convertNumberToType(num, clazz) as T
            }
            ValueType.I16 -> {
                val num = (data as? Number) ?: text.toShort()
                convertNumberToType(num, clazz) as T
            }
            ValueType.I32 -> {
                val num = (data as? Number) ?: text.toInt()
                convertNumberToType(num, clazz) as T
            }
            ValueType.I64 -> {
                val num = (data as? Number) ?: text.toLong()
                convertNumberToType(num, clazz) as T
            }
            ValueType.U -> {
                val num = (data as? Number) ?: text.toInt()
                convertNumberToType(num, clazz) as T
            }
            ValueType.U8 -> {
                val num = (data as? Number) ?: text.toShort()
                convertNumberToType(num, clazz) as T
            }
            ValueType.U16 -> {
                val num = (data as? Number) ?: text.toInt()
                convertNumberToType(num, clazz) as T
            }
            ValueType.U32 -> {
                val num = (data as? Number) ?: text.toInt()
                convertNumberToType(num, clazz) as T
            }
            ValueType.U64 -> {
                val num = (data as? Number) ?: text.toLong()
                convertNumberToType(num, clazz) as T
            }
            ValueType.F32 -> {
                val num = (data as? Number) ?: text.toFloat()
                convertNumberToType(num, clazz) as T
            }
            ValueType.F64 -> {
                val num = (data as? Number) ?: text.toDouble()
                convertNumberToType(num, clazz) as T
            }
            ValueType.STR -> {
                (data as? String ?: text) as T
            }
            ValueType.BOOL -> {
                (data as? Boolean ?: text.toBoolean()) as T
            }
            ValueType.BYTES -> {
                (data as? ByteArray) as T
            }
            else -> throw IllegalArgumentException("unsupported type: ${tag.type}")
        }
    }

    private fun convertScalarToAny(value: NodeScalar): Any? {
        return convertScalar(value, Any::class.java)
    }

    @Suppress("UNCHECKED_CAST")
    private fun <T> convertNumberToType(num: Number, clazz: Class<T>): Any {
        return when (clazz) {
            Byte::class.java, java.lang.Byte.TYPE -> num.toByte()
            Short::class.java, java.lang.Short.TYPE -> num.toShort()
            Int::class.java, java.lang.Integer.TYPE -> num.toInt()
            Long::class.java, java.lang.Long.TYPE -> num.toLong()
            Float::class.java, java.lang.Float.TYPE -> num.toFloat()
            Double::class.java, java.lang.Double.TYPE -> num.toDouble()
            BigInteger::class.java -> BigInteger.valueOf(num.toLong())
            else -> num
        } as
                Any
    }

    private fun materialize(f: java.lang.reflect.Field, node: Node): Any? {
        return when (node) {
            is NodeNull -> null
            is NodeScalar -> convertScalar(node, f.type)
            is NodeObject -> {
                val tag = node.tag
                if (tag != null && tag.type == ValueType.MAP) {
                    val map = mutableMapOf<String, Any?>()
                    convertMap(node, map)
                    map
                } else {
                    val inst = f.type.getDeclaredConstructor().newInstance()
                    convertObj(node, inst)
                    inst
                }
            }
            is NodeArray -> {
                convertVec(node, f)
            }
            else -> null
        }
    }
}
