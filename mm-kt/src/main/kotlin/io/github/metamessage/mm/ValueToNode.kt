package io.github.metamessage.mm

import java.math.BigInteger
import java.net.InetAddress
import java.net.URI
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneOffset
import java.time.format.DateTimeFormatter
import java.util.UUID
import kotlin.reflect.KClass
import kotlin.reflect.KMutableProperty
import kotlin.reflect.KProperty
import kotlin.reflect.KVisibility
import kotlin.reflect.full.declaredMemberProperties
import kotlin.reflect.full.findAnnotation

import io.github.metamessage.ast.Tag
import io.github.metamessage.ast.ValueType
import io.github.metamessage.ast.Field
import io.github.metamessage.ast.Node
import io.github.metamessage.ast.Value as AstValue
import io.github.metamessage.ast.Array as AstArray
import io.github.metamessage.ast.Object as AstObject

private const val maxDepth = 32

fun valueToNode(v: Any?, tag: Tag, path: String): Node {
    return valueToNode(v, tag, 0, path)
}

private fun valueToNode(v: Any?, tag: Tag?, depth: Int, path: String): Node {
    var workTag = tag ?: Tag.empty()
    
    var data: Any? = null
    var text: String = ""

    when (v) {
        null -> {
            if (workTag.type == ValueType.UNKNOWN) {
                throw IllegalArgumentException("invalid input: v is untyped nil (no concrete type/value)")
            }
            workTag.isNull = true
        }

        is ByteArray -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.SLICE
            }
            when (workTag.type) {
                ValueType.BYTES -> {
                    val result = workTag.validateBytes(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.IMAGE -> {
                    val result = workTag.validateImage(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.SLICE -> {
                    return anyToJSONC(v, workTag, depth, path)
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Boolean -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.BOOL
            }
            when (workTag.type) {
                ValueType.BOOL -> {
                    val result = workTag.validateBool(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Byte -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT8
            }
            when (workTag.type) {
                ValueType.INT8 -> {
                    val result = workTag.validateInt8(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Short -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT16
            }
            when (workTag.type) {
                ValueType.INT16 -> {
                    val result = workTag.validateInt16(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Int -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT
            }
            when (workTag.type) {
                ValueType.INT -> {
                    val result = workTag.validateInt(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Long -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT64
            }
            when (workTag.type) {
                ValueType.INT64 -> {
                    val result = workTag.validateInt64(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.UINT64 -> {
                    val result = workTag.validateUint64(BigInteger.valueOf(v))
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Float -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.FLOAT32
            }
            when (workTag.type) {
                ValueType.FLOAT32 -> {
                    if (v.isInfinite() || v.isNaN()) {
                        throw IllegalArgumentException("${workTag.type} unsupported value: ${if (v.isInfinite()) "Inf" else "NaN"}")
                    }
                    val result = workTag.validateFloat32(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Double -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.FLOAT64
            }
            when (workTag.type) {
                ValueType.FLOAT64 -> {
                    if (v.isInfinite() || v.isNaN()) {
                        throw IllegalArgumentException("${workTag.type} unsupported value: ${if (v.isInfinite()) "Inf" else "NaN"}")
                    }
                    val result = workTag.validateFloat64(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.DECIMAL -> {
                    val result = workTag.validateDecimal(v.toString())
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is String -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.STRING
            }
            when (workTag.type) {
                ValueType.STRING -> {
                    val result = workTag.validateString(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.DECIMAL -> {
                    val result = workTag.validateDecimal(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.EMAIL -> {
                    val result = workTag.validateEmail(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.UUID -> {
                    val result = workTag.validateUUID(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.URL -> {
                    val result = workTag.validateURL(v)
                    data = result.data
                    text = result.text ?: ""    
                }
                ValueType.IP -> {
                    val result = workTag.validateIP(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.ENUM -> {
                    val result = workTag.validateEnum(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is BigInteger -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.BIGINT
            }
            when (workTag.type) {
                ValueType.BIGINT -> {
                    val result = workTag.validateBigInt(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is UUID -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.UUID
            }
            when (workTag.type) {
                ValueType.UUID -> {
                    val result = workTag.validateUUID(v.toString())
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is InetAddress -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.IP
            }
            when (workTag.type) {
                ValueType.IP -> {
                    val result = workTag.validateIP(v.hostAddress ?: "")
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is URI -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.URL
            }
            when (workTag.type) {
                ValueType.URL -> {
                    val result = workTag.validateURL(v.toString())
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is LocalDateTime -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.DATETIME
            }
            when (workTag.type) {
                ValueType.DATETIME -> {
                    val result = workTag.validateDateTime(v)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.DATE -> {
                    val result = workTag.validateDate(v.toLocalDate())
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.TIME -> {
                    val result = workTag.validateTime(v.toLocalTime())
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is LocalDate -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.DATE
            }
            when (workTag.type) {
                ValueType.DATE -> {
                    val result = workTag.validateDate(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is LocalTime -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.TIME
            }
            when (workTag.type) {
                ValueType.TIME -> {
                    val result = workTag.validateTime(v)
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        is Instant -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.DATETIME
            }
            val localDateTime = LocalDateTime.ofInstant(v, ZoneOffset.UTC)
            when (workTag.type) {
                ValueType.DATETIME -> {
                    val result = workTag.validateDateTime(localDateTime)
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.DATE -> {
                    val result = workTag.validateDate(localDateTime.toLocalDate())
                    data = result.data
                    text = result.text ?: ""
                }
                ValueType.TIME -> {
                    val result = workTag.validateTime(localDateTime.toLocalTime())
                    data = result.data
                    text = result.text ?: ""
                }
                else -> throw IllegalArgumentException("${workTag.type} unsupported type: ${v.javaClass.name}")
            }
        }

        else -> {
            return anyToJSONC(v, workTag, depth, path)
        }
    }

    return AstValue(data, text, workTag, path)
}

private fun anyToJSONC(obj: Any, tag: Tag, depth: Int, path: String): Node {
    val newDepth = depth + 1
    if (newDepth > maxDepth) {
        throw IllegalArgumentException("max depth: $maxDepth")
    }

    val kClass = obj::class
    val workTag = tag.copy()

    when {
        obj is List<*> -> {
            workTag.type = ValueType.SLICE
            return listToNode(obj, workTag, newDepth, path)
        }
        obj is Map<*, *> -> {
            workTag.type = ValueType.MAP
            return mapToNode(obj, workTag, newDepth, path)
        }
        kClass.isData || kClass.java.isAnnotationPresent(MM::class.java) -> {
            workTag.type = ValueType.STRUCT
            return structToNode(obj, workTag, newDepth, path)
        }
        else -> {
            workTag.type = ValueType.STRUCT
            return structToNode(obj, workTag, newDepth, path)
        }
    }
}

private fun structToNode(o: Any, objTag: Tag, depth: Int, path: String): Node {
    val objNode = AstObject(mutableListOf(), objTag, path)
    val kClass = o::class

    for (prop in kClass.declaredMemberProperties) {
        if (prop.visibility == KVisibility.PRIVATE) continue

        val fieldKey = getFieldKey(prop, objTag)
        if (fieldKey == "-") continue

        val mm = prop.findAnnotation<MM>()
        var ft = Tag.empty()

        if (mm != null) {
            val ann = Tag.fromAnnotation(mm)
            ft = ann.copy()
            if (ft.type == ValueType.UNKNOWN) {
                ft.type = TypeInference.forProperty(prop)
            }
        } else {
            ft.type = TypeInference.forProperty(prop)
        }

        val fieldPath = if (path.isEmpty()) fieldKey else "$path.$fieldKey"
        
        val fieldValue = try {
            (prop as? KMutableProperty<*>)?.getter?.call(o) ?: prop.getter.call(o)
        } catch (e: Exception) {
            null
        }

        val fieldNode = valueToNode(fieldValue, ft, depth, fieldPath)
        objNode.fields.add(Field(fieldKey, fieldNode))
    }

    return objNode
}

private fun listToNode(list: List<*>, tag: Tag, depth: Int, path: String): Node {
    val arrNode = AstArray(mutableListOf(), tag, path)
    var setTag = false

    for ((i, item) in list.withIndex()) {
        val et = Tag.empty()
        et.inheritFromArrayParent(tag)

        if (et.type == ValueType.UNKNOWN && item != null) {
            et.type = TypeInference.valueTypeForComponent(item.javaClass)
            if (et.type == ValueType.UNKNOWN) {
                et.type = ValueType.STRUCT
            }
        }

        val itemPath = "$path[$i]"
        val itemNode = valueToNode(item, et, depth, itemPath)

        updateChildTag(arrNode, et, !setTag)
        setTag = true

        arrNode.items.add(itemNode)
    }

    if (list.isEmpty()) {
        val exampleVal = createExampleValue(tag.childType)
        val et = Tag.empty()
        et.inheritFromArrayParent(tag)
        et.example = true

        val itemPath = "$path[0]"
        val itemNode = valueToNode(exampleVal, et, depth, itemPath)

        updateChildTag(arrNode, et, !setTag)
        arrNode.items.add(itemNode)
    }

    return arrNode
}

private fun mapToNode(map: Map<*, *>, tag: Tag, depth: Int, path: String): Node {
    val objNode = AstObject(mutableListOf(), tag, path)
    var setTag = false

    for ((key, value) in map) {
        val keyStr = key?.toString()?.let { CamelToSnake.convert(it) } ?: ""
        val valueTag = Tag.empty()
        
        if (tag.childType != ValueType.UNKNOWN) {
            valueTag.type = tag.childType
        } else if (value != null) {
            valueTag.type = TypeInference.valueTypeForComponent(value.javaClass)
        }

        val fieldPath = "$path[$keyStr]"
        val valueNode = valueToNode(value, valueTag, depth, fieldPath)

        updateChildTag(objNode, valueTag, !setTag)
        setTag = true

        objNode.fields.add(Field(keyStr, valueNode))
    }

    if (map.isEmpty()) {
        val exampleVal = createExampleValue(tag.childType)
        val valueTag = Tag.empty()
        valueTag.type = tag.childType
        valueTag.example = true

        val fieldPath = "$path[]"
        val valueNode = valueToNode(exampleVal, valueTag, depth, fieldPath)

        updateChildTag(objNode, valueTag, !setTag)
        objNode.fields.add(Field("", valueNode))
    }

    return objNode
}

private fun updateChildTag(parent: Node, childTag: Tag, update: Boolean) {
    if (!update) return

    val parentTag = parent.tag ?: return
    parentTag.childType = childTag.type
    parentTag.childDesc = childTag.desc
    parentTag.childNullable = childTag.nullable
    parentTag.childAllowEmpty = childTag.allowEmpty
    parentTag.childUnique = childTag.unique
    parentTag.childDefault = childTag.default
    parentTag.childMin = childTag.min
    parentTag.childMax = childTag.max
    parentTag.childSize = childTag.size
    parentTag.childEnum = childTag.enum
    parentTag.childPattern = childTag.pattern
    parentTag.childLocation = childTag.location
    parentTag.childVersion = childTag.version
    parentTag.childMime = childTag.mime
}

private fun getFieldKey(prop: KProperty<*>, tag: Tag): String {
    val mm = prop.findAnnotation<MM>()
    if (mm != null && mm.name.isNotEmpty()) {
        return mm.name
    }
    if (tag.name.isNotEmpty()) {
        return tag.name
    }
    return CamelToSnake.convert(prop.name)
}

private fun createExampleValue(type: ValueType): Any? {
    return when (type) {
        ValueType.INT8 -> 0.toByte()
        ValueType.INT16 -> 0.toShort()
        ValueType.INT32 -> 0
        ValueType.INT64 -> 0L
        ValueType.UINT, ValueType.UINT8, ValueType.UINT16, ValueType.UINT32 -> 0
        ValueType.UINT64 -> 0L
        ValueType.FLOAT32 -> 0.0f
        ValueType.FLOAT64 -> 0.0
        ValueType.STRING, ValueType.EMAIL, ValueType.URL, ValueType.IP -> ""
        ValueType.BOOL -> false
        ValueType.BYTES -> ByteArray(0)
        ValueType.BIGINT -> BigInteger.ZERO
        ValueType.UUID -> UUID.randomUUID()
        ValueType.DATETIME -> LocalDateTime.now()
        ValueType.DATE -> LocalDate.now()
        ValueType.TIME -> LocalTime.now()
        ValueType.ENUM -> 0
        else -> null
    }
}