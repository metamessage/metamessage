package io.github.metamessage.mm

import java.math.BigInteger
import java.time.LocalDateTime
import java.util.UUID
import io.github.metamessage.ast.Tag

import io.github.metamessage.ast.Node
import io.github.metamessage.ast.Object
import io.github.metamessage.ast.Array
import io.github.metamessage.ast.Value
import io.github.metamessage.MM

object Binder {
    fun <T> bind(node: Node, clazz: Class<T>): T {
        require(node is Object) { "root must be object" }
        val inst = clazz.getDeclaredConstructor().newInstance()
        val byKey = fieldsByJsonKey(clazz)
        for (e in node.fields) {
            val f = byKey[e.key] ?: continue
            f.isAccessible = true
            f.set(inst, materialize(f, e.value))
        }
        return inst
    }

    private fun fieldsByJsonKey(c: Class<*>): Map<String, java.lang.reflect.Field> {
        val m = LinkedHashMap<String, java.lang.reflect.Field>()
        for (f in c.declaredFields) {
            if (java.lang.reflect.Modifier.isStatic(f.modifiers)) continue
            val t = TypeInference.forField(f)
            val mm = f.getAnnotation(MM::class.java)
            if (mm != null && mm.name == "-") continue
            val key = fieldKey(f, t, mm)
            m[key] = f
        }
        return m
    }

    private fun fieldKey(f: java.lang.reflect.Field, ft: Tag, mm: MM?): String {
        if (mm != null && mm.name.isNotEmpty() && mm.name != "-") return mm.name
        if (ft.name.isNotEmpty()) return ft.name
        return CamelToSnake.convert(f.name)
    }

    private fun materialize(f: java.lang.reflect.Field, node: Node): Any? {
        val ft = f.type
        if (Map::class.java.isAssignableFrom(ft)) {
            if (node is Object) {
                return mapFrom(node, f)
            }
        }
        return when (node) {
            is Value -> scalarToField(ft, node)
            is Array -> listFrom(node, f)
            is Object -> bind(node, ft)
            is io.github.metamessage.ast.Doc -> bind(node, ft)
        }
    }

    @Suppress("UNCHECKED_CAST")
    private fun listFrom(arr: Array, f: java.lang.reflect.Field): List<*> {
        val gt = f.genericType
        var elemClass: Class<Any> = String::class.java as Class<Any>
        if (gt is java.lang.reflect.ParameterizedType) {
            val typeArgs = gt.actualTypeArguments
            if (typeArgs.isNotEmpty()) {
                val a0 = typeArgs[0]
                if (a0 is Class<*>) {
                    elemClass = a0 as Class<Any>
                }
            }
        }
        val out = mutableListOf<Any?>()
        for (ch in arr.items) {
            when (ch) {
                is Value -> out.add(scalarToField(elemClass, ch))
                is Object -> {
                    val sub = elemClass.getDeclaredConstructor().newInstance()
                    val byKey = fieldsByJsonKey(elemClass)
                    for (e in ch.fields) {
                        val sf = byKey[e.key]
                        if (sf != null) {
                            sf.isAccessible = true
                            sf.set(sub, materialize(sf, e.value))
                        }
                    }
                    out.add(sub)
                }
                else -> out.add(null)
            }
        }
        return out
    }

    @Suppress("UNCHECKED_CAST")
    private fun mapFrom(obj: Object, f: java.lang.reflect.Field): Map<*, *> {
        val gt = f.genericType
        var keyClass: Class<Any> = String::class.java as Class<Any>
        var valueClass: Class<Any> = Any::class.java
        
        if (gt is java.lang.reflect.ParameterizedType) {
            val typeArgs = gt.actualTypeArguments
            if (typeArgs.size >= 2) {
                val a0 = typeArgs[0]
                if (a0 is Class<*>) {
                    keyClass = a0 as Class<Any>
                }
                val a1 = typeArgs[1]
                if (a1 is Class<*>) {
                    valueClass = a1 as Class<Any>
                }
            }
        }
        
        val out = mutableMapOf<Any?, Any?>()
        for ((key, value) in obj.fields) {
            // Convert key to the appropriate type
            val convertedKey = when (keyClass) {
                String::class.java -> key
                Int::class.java, Int::class.javaPrimitiveType -> key.toInt()
                Long::class.java, Long::class.javaPrimitiveType -> key.toLong()
                else -> key
            }
            
            // Convert value to the appropriate type
            val convertedValue = when (value) {
                is Value -> {
                    val scalarValue = scalarToField(valueClass, value)
                    // Ensure the value is of the correct type
                    when {
                        valueClass == Int::class.java || valueClass == Int::class.javaPrimitiveType -> scalarValue as? Int ?: scalarValue.toString().toInt()
                        valueClass == Long::class.java || valueClass == Long::class.javaPrimitiveType -> scalarValue as? Long ?: scalarValue.toString().toLong()
                        valueClass == Double::class.java || valueClass == Double::class.javaPrimitiveType -> scalarValue as? Double ?: scalarValue.toString().toDouble()
                        valueClass == Float::class.java || valueClass == Float::class.javaPrimitiveType -> scalarValue as? Float ?: scalarValue.toString().toFloat()
                        else -> scalarValue
                    }
                }
                is Object -> bind(value, valueClass)
                is Array -> listFrom(value, f)
                is io.github.metamessage.ast.Doc -> bind(value, valueClass)
            }
            
            out[convertedKey] = convertedValue
        }
        return out
    }

    private fun scalarToField(ft: Class<*>, sc: Value): Any? {
        val d = sc.data
        return when {
            ft == Int::class.javaPrimitiveType || ft == Int::class.java -> if (d is Number) d.toInt() else sc.text.toInt()
            ft == Long::class.javaPrimitiveType || ft == Long::class.java -> if (d is Number) d.toLong() else sc.text.toLong()
            ft == Boolean::class.javaPrimitiveType || ft == Boolean::class.java -> d == true
            ft == String::class.java -> sc.text
            ft == Double::class.javaPrimitiveType || ft == Double::class.java -> if (d is Number) d.toDouble() else sc.text.toDouble()
            ft == Float::class.javaPrimitiveType || ft == Float::class.java -> if (d is Number) d.toFloat() else sc.text.toFloat()
            ft == LocalDateTime::class.java -> when (d) {
                is LocalDateTime -> d
                else -> LocalDateTime.parse(sc.text.replace(' ', 'T'))
            }
            ft == UUID::class.java && d is UUID -> d
            ft == BigInteger::class.java && d is BigInteger -> d
            else -> {
                // For other types, try to convert from text if possible
                when (ft) {
                    Int::class.java -> sc.text.toInt()
                    Long::class.java -> sc.text.toLong()
                    Double::class.java -> sc.text.toDouble()
                    Float::class.java -> sc.text.toFloat()
                    else -> d
                }
            }
        }
    }
}
