package io.github.metamessage.core

import io.github.metamessage.ir.Array
import io.github.metamessage.ir.Node
import io.github.metamessage.ir.Object
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.Value
import io.github.metamessage.ir.ValueType
import java.math.BigInteger
import java.time.LocalDateTime
import java.util.UUID

object Binder {
    fun <T> bind(node: Node, clazz: Class<T>): T {
        when (node) {
            is Object -> {
                val tag = node.tag
                if (tag != null && tag.type == ValueType.STRUCT) {
                    val inst = clazz.getDeclaredConstructor().newInstance()
                    convertObj(node, inst as Any)
                    return inst
                } else {
                    @Suppress("UNCHECKED_CAST") val inst = mutableMapOf<String, Any?>()
                    convertMap(node, inst)
                    return inst as T
                }
            }
            is Array -> {
                val tag = node.tag
                if (tag != null && tag.size > 0 && tag.type == ValueType.ARRAY) {
                    @Suppress("UNCHECKED_CAST") return convertArr(node, clazz) as T
                } else {
                    return convertVec(node, clazz)
                }
            }
            is Value -> {
                return convertScalar(node, clazz)
            }
            else ->
                    throw IllegalArgumentException(
                            "unsupported node type: ${node::class.java.name}"
                    )
        }
    }

    private fun convertObj(obj: Object, out: Any) {
        val outClazz = out.javaClass
        val nameToField = mutableMapOf<String, java.lang.reflect.Field>()
        for (f in outClazz.declaredFields) {
            if (java.lang.reflect.Modifier.isStatic(f.modifiers)) continue
            nameToField[f.name] = f
        }

        for (field in obj.fields) {
            val fieldKey = field.key
            val runes = fieldKey.toCharArray()
            if (runes.isNotEmpty()) {
                runes[0] = runes[0].uppercaseChar()
            }
            val name = String(runes)
            val structField = nameToField[name] ?: continue
            structField.isAccessible = true
            try {
                val fieldVal = materialize(structField, field.value)
                structField.set(out, fieldVal)
            } catch (e: Exception) {
                throw RuntimeException("failed to bind field ${field.key}: ${e.message}", e)
            }
        }
    }

    private fun convertMap(obj: Object, out: MutableMap<String, Any?>) {
        for (field in obj.fields) {
            val key = field.key
            val value =
                    when (val v = field.value) {
                        is Value -> convertScalarToAny(v)
                        is Object -> {
                            val map = mutableMapOf<String, Any?>()
                            convertMap(v, map)
                            map
                        }
                        is Array -> {
                            convertVec(v, List::class.java).toList()
                        }
                        else -> null
                    }
            out[key] = value
        }
    }

    private fun convertArr(arr: Array, clazz: Class<*>): Any {
        val size = arr.tag?.size ?: arr.items.size
        val list = mutableListOf<Any?>()
        for (item in arr.items) {
            when (item) {
                is Value -> list.add(convertScalarToAny(item))
                is Object -> {
                    val inst = clazz.getDeclaredConstructor().newInstance()
                    convertObj(item, inst as Any)
                    list.add(inst)
                }
                else -> list.add(null)
            }
        }
        return list
    }

    @Suppress("UNCHECKED_CAST")
    private fun <T> convertVec(arr: Array, clazz: Class<T>): T {
        val list = mutableListOf<Any?>()
        var elemClass: Class<*> = Any::class.java
        val gt = clazz.typeParameters.firstOrNull()
        if (gt != null && gt is Class<*>) {
            elemClass = gt
        }
        for (item in arr.items) {
            when (item) {
                is Value -> list.add(convertScalarToAny(item))
                is Object -> {
                    val inst = clazz.getDeclaredConstructor().newInstance()
                    convertObj(item, inst as Any)
                    list.add(inst)
                }
                else -> list.add(null)
            }
        }
        return list as T
    }

    @Suppress("UNCHECKED_CAST")
    private fun <T> convertScalar(value: Value, clazz: Class<T>): T {
        val tag = value.tag ?: Tag.empty()
        val data = value.data
        val text = value.text

        return when (tag.type) {
            ValueType.DATETIME, ValueType.DATE, ValueType.TIME -> {
                (data as? LocalDateTime ?: LocalDateTime.of(1970, 1, 1, 0, 0, 0)) as T
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
                (data as? java.net.InetAddress) as T
            }
            ValueType.URL -> {
                (data as? java.net.URI) as T
            }
            ValueType.ENUM -> {
                text as T
            }
            ValueType.INT -> {
                ((data as? Number)?.toInt() ?: text.toInt()) as T
            }
            ValueType.INT8 -> {
                ((data as? Number)?.toByte() ?: text.toByte()) as T
            }
            ValueType.INT16 -> {
                ((data as? Number)?.toShort() ?: text.toShort()) as T
            }
            ValueType.INT32 -> {
                ((data as? Number)?.toInt() ?: text.toInt()) as T
            }
            ValueType.INT64 -> {
                ((data as? Number)?.toLong() ?: text.toLong()) as T
            }
            ValueType.UINT -> {
                ((data as? Number)?.toInt() ?: text.toInt()) as T
            }
            ValueType.UINT8 -> {
                ((data as? Number)?.toShort() ?: text.toShort()) as T
            }
            ValueType.UINT16 -> {
                ((data as? Number)?.toInt() ?: text.toInt()) as T
            }
            ValueType.UINT32 -> {
                ((data as? Number)?.toInt() ?: text.toInt()) as T
            }
            ValueType.UINT64 -> {
                ((data as? Number)?.toLong() ?: text.toLong()) as T
            }
            ValueType.FLOAT32 -> {
                ((data as? Number)?.toFloat() ?: text.toFloat()) as T
            }
            ValueType.FLOAT64 -> {
                ((data as? Number)?.toDouble() ?: text.toDouble()) as T
            }
            ValueType.STRING -> {
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

    private fun convertScalarToAny(value: Value): Any? {
        return convertScalar(value, Any::class.java) as? Any?
    }

    private fun materialize(f: java.lang.reflect.Field, node: Node): Any? {
        return when (node) {
            is Value -> convertScalar(node, f.type)
            is Object -> {
                val inst = f.type.getDeclaredConstructor().newInstance()
                convertObj(node, inst)
                inst
            }
            is Array -> {
                val list = convertVec(node, f.type)
                list
            }
            else -> null
        }
    }
}
