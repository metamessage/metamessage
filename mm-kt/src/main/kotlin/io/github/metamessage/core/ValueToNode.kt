package io.github.metamessage.core

import io.github.metamessage.MM
import io.github.metamessage.ir.Array as AstArray
import io.github.metamessage.ir.Field
import io.github.metamessage.ir.Node
import io.github.metamessage.ir.Object as AstObject
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.Value as AstValue
import io.github.metamessage.ir.ValueType
import java.math.BigInteger
import java.net.InetAddress
import java.net.URI
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.reflect.KMutableProperty
import kotlin.reflect.KVisibility
import kotlin.reflect.full.declaredMemberProperties
import kotlin.reflect.full.findAnnotation

private const val maxDepth = 32

private const val Null = "null"
private const val True = "true"
private const val False = "false"

fun nilToNode(valueType: ValueType): Node {
    val tag = Tag.empty()
    tag.type = valueType
    return AstValue(null, Null, tag, "")
}

fun valueToNode(v: Any?, tag: Tag, path: String): Node {
    return valueToNode(v, tag, 0, path)
}

private fun valueToNode(v: Any?, tag: Tag?, depth: Int, path: String): Node {
    var workTag = tag ?: Tag.empty()

    var data: Any? = null
    var text: String = Null

    when (v) {
        null -> {
            if (workTag.type == ValueType.UNKNOWN) {
                throw IllegalArgumentException(
                        "invalid input: v is untyped nil (no concrete type/value)"
                )
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
                    text = result.text ?: Null
                }
                ValueType.IMAGE -> {
                    val result = workTag.validateImage(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.SLICE -> {
                    return anyToJSONC(v, workTag, depth, path)
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Byte -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT8
            }
            when (workTag.type) {
                ValueType.INT8 -> {
                    val result = workTag.validateI8(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Short -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT16
            }
            when (workTag.type) {
                ValueType.INT16 -> {
                    val result = workTag.validateI16(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Int -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT32
            }
            when (workTag.type) {
                ValueType.INT, ValueType.INT32 -> {
                    val result = workTag.validateI32(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Long -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.INT64
            }
            when (workTag.type) {
                ValueType.INT64 -> {
                    val result = workTag.validateI64(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.UINT64 -> {
                    val result = workTag.validateU64(BigInteger.valueOf(v))
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Float -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.FLOAT32
            }
            when (workTag.type) {
                ValueType.FLOAT32 -> {
                    if (v.isInfinite() || v.isNaN()) {
                        val desc =
                                when {
                                    v.isInfinite() && v > 0 -> "+Inf"
                                    v.isInfinite() && v < 0 -> "-Inf"
                                    else -> "NaN"
                                }
                        throw IllegalArgumentException("${workTag.type} unsupported value: $desc")
                    }
                    val result = workTag.validateF32(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is Double -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.FLOAT64
            }
            when (workTag.type) {
                ValueType.FLOAT64 -> {
                    if (v.isInfinite() || v.isNaN()) {
                        val desc =
                                when {
                                    v.isInfinite() && v > 0 -> "+Inf"
                                    v.isInfinite() && v < 0 -> "-Inf"
                                    else -> "NaN"
                                }
                        throw IllegalArgumentException("${workTag.type} unsupported value: $desc")
                    }
                    val result = workTag.validateF64(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.DECIMAL -> {
                    val result = workTag.validateDecimal(v.toString())
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is String -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.STRING
            }
            when (workTag.type) {
                ValueType.STRING -> {
                    val result = workTag.validateStr(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.DECIMAL -> {
                    val result = workTag.validateDecimal(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.EMAIL -> {
                    val result = workTag.validateEmail(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.ENUM -> {
                    val result = workTag.validateEnum(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.UUID -> {
                    val result = workTag.validateUUID(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.URL -> {
                    val result = workTag.validateURL(v)
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.IP -> {
                    val result = workTag.validateIP(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
            }
        }
        is BigInteger -> {
            if (workTag.type == ValueType.UNKNOWN) {
                workTag.type = ValueType.BIGINT
            }
            when (workTag.type) {
                ValueType.BIGINT -> {
                    val result = workTag.validateBigint(v)
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                ValueType.DATE -> {
                    val result = workTag.validateDate(v.toLocalDate())
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.TIME -> {
                    val result = workTag.validateTime(v.toLocalTime())
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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
                    text = result.text ?: Null
                }
                ValueType.DATE -> {
                    val result = workTag.validateDate(localDateTime.toLocalDate())
                    data = result.data
                    text = result.text ?: Null
                }
                ValueType.TIME -> {
                    val result = workTag.validateTime(localDateTime.toLocalTime())
                    data = result.data
                    text = result.text ?: Null
                }
                else ->
                        throw IllegalArgumentException(
                                "${workTag.type} unsupported type: ${v.javaClass.name}"
                        )
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

    if (workTag.toString().isEmpty()) {
        val mmAnnotation = kClass.java.getAnnotation(MM::class.java)
        if (mmAnnotation != null) {
            val mmTag = Tag.fromAnnotation(mmAnnotation)
            mergeTag(workTag, mmTag)
        }
    }

    workTag.type = ValueType.STRUCT

    var typeName = CamelToSnake.convert(kClass.simpleName ?: "")
    if (workTag.name.isNotEmpty()) {
        typeName = workTag.name
    }
    workTag.name = typeName

    val objPath =
            if (typeName.isNotEmpty()) {
                if (path.isEmpty()) typeName else "$path.$typeName"
            } else {
                path
            }

    when {
        obj is List<*> -> {
            workTag.type = ValueType.SLICE
            return convertVec(obj, workTag, newDepth, objPath)
        }
        obj is Map<*, *> -> {
            workTag.type = ValueType.MAP
            return convertMap(obj, workTag, newDepth, objPath)
        }
        obj is Array<*> -> {
            workTag.type = ValueType.SLICE
            return convertVec(obj.toList(), workTag, newDepth, objPath)
        }
        else -> {
            workTag.type = ValueType.STRUCT
            return convertObj(obj, workTag, newDepth, objPath)
        }
    }
}

private fun convertObj(obj: Any, objTag: Tag, depth: Int, path: String): Node {
    val kClass = obj::class
    val objNode = AstObject(mutableListOf(), objTag, path)

    for (prop in kClass.declaredMemberProperties) {
        if (prop.visibility == KVisibility.PRIVATE) continue

        val fieldKey = CamelToSnake.convert(prop.name)

        val mmAnnotation = prop.findAnnotation<MM>()
        var fieldTag: Tag? = null
        if (mmAnnotation != null) {
            fieldTag = Tag.fromAnnotation(mmAnnotation)
            if (fieldTag.name.isNotEmpty()) {
                if (fieldTag.name == "-") continue
            }
        }

        if (fieldTag == null) {
            fieldTag = Tag.empty()
        }

        val effectiveFieldKey =
                when {
                    fieldTag.name.isNotEmpty() -> fieldTag.name
                    else -> fieldKey
                }

        if (fieldTag.name.isEmpty()) {
            fieldTag.name = effectiveFieldKey
        }

        if (fieldTag.type == ValueType.UNKNOWN) {
            fieldTag.type = TypeInference.forProperty(prop)
        }

        val fieldPath = "$path.$effectiveFieldKey"

        val fieldValue =
                try {
                    (prop as? KMutableProperty<*>)?.getter?.call(obj) ?: prop.getter.call(obj)
                } catch (e: Exception) {
                    null
                }

        val fieldNode = valueToNode(fieldValue, fieldTag, depth, fieldPath)
        objNode.fields.add(Field(effectiveFieldKey, fieldNode))
    }

    val result = objTag.validateObj()
    if (!result.valid) {
        throw IllegalArgumentException("validate failed: ${result.error}")
    }

    return objNode
}

private fun convertVec(list: List<*>, tag: Tag, depth: Int, path: String): Node {
    val node = AstArray(mutableListOf(), tag, path)
    var setTag = false

    for ((i, item) in list.withIndex()) {
        val itemTag = Tag.empty()
        itemTag.inheritFromArrayParent(tag)

        val itemPath = "$path[$i]"
        val itemNode = valueToNode(item, itemTag, depth, itemPath)

        val resultTag = itemNode.tag
        if (!setTag && resultTag != null) {
            node.tag?.apply {
                childDesc = resultTag.desc
                childType = resultTag.type
                childRaw = resultTag.raw
                childNullable = resultTag.nullable
                childAllowEmpty = resultTag.allowEmpty
                childUnique = resultTag.unique
                childDefault = resultTag.default
                childMin = resultTag.min
                childMax = resultTag.max
                childSize = resultTag.size
                childEnum = resultTag.enum
                childPattern = resultTag.pattern
                childLocation = resultTag.location
                childVersion = resultTag.version
                childMime = resultTag.mime
            }
            setTag = true
        }

        node.items.add(itemNode)
    }

    if (list.isEmpty()) {
        val itemTag = Tag.empty()
        itemTag.inheritFromArrayParent(tag)
        itemTag.example = true

        val itemPath = "$path[0]"
        val exampleVal = createExampleValue(tag.childType)
        val itemNode = valueToNode(exampleVal, itemTag, depth, itemPath)

        val resultTag = itemNode.tag
        if (!setTag && resultTag != null) {
            node.tag?.apply {
                childDesc = resultTag.desc
                childType = resultTag.type
                childRaw = resultTag.raw
                childNullable = resultTag.nullable
                childAllowEmpty = resultTag.allowEmpty
                childUnique = resultTag.unique
                childDefault = resultTag.default
                childMin = resultTag.min
                childMax = resultTag.max
                childSize = resultTag.size
                childEnum = resultTag.enum
                childPattern = resultTag.pattern
                childLocation = resultTag.location
                childVersion = resultTag.version
                childMime = resultTag.mime
            }
        }

        node.items.add(itemNode)
    }

    val result = tag.validateVec(node.items)
    if (!result.valid) {
        throw IllegalArgumentException("validate failed: ${result.error}")
    }

    return node
}

private fun convertMap(map: Map<*, *>, tag: Tag, depth: Int, path: String): Node {
    val node = AstObject(mutableListOf(), tag, path)
    var setTag = false

    for ((key, value) in map) {
        val keyStr = CamelToSnake.convert(key?.toString() ?: "")

        val valueTag = Tag.empty()
        valueTag.inheritFromArrayParent(tag)
        valueTag.name = keyStr

        val fieldPath = "$path[$keyStr]"
        val valueNode = valueToNode(value, valueTag, depth, fieldPath)

        val resultTag = valueNode.tag
        if (!setTag && resultTag != null) {
            node.tag?.apply {
                childDesc = resultTag.desc
                childType = resultTag.type
                childRaw = resultTag.raw
                childNullable = resultTag.nullable
                childAllowEmpty = resultTag.allowEmpty
                childUnique = resultTag.unique
                childDefault = resultTag.default
                childMin = resultTag.min
                childMax = resultTag.max
                childSize = resultTag.size
                childEnum = resultTag.enum
                childPattern = resultTag.pattern
                childLocation = resultTag.location
                childVersion = resultTag.version
                childMime = resultTag.mime
            }
            setTag = true
        }

        node.fields.add(Field(keyStr, valueNode))
    }

    if (map.isEmpty()) {
        val valueTag = Tag.empty()
        valueTag.inheritFromArrayParent(tag)
        valueTag.name = ""
        valueTag.example = true

        val fieldPath = "$path[]"
        val exampleVal = createExampleValue(tag.childType)
        val valueNode = valueToNode(exampleVal, valueTag, depth, fieldPath)

        val resultTag = valueNode.tag
        if (!setTag && resultTag != null) {
            node.tag?.apply {
                childDesc = resultTag.desc
                childType = resultTag.type
                childRaw = resultTag.raw
                childNullable = resultTag.nullable
                childAllowEmpty = resultTag.allowEmpty
                childUnique = resultTag.unique
                childDefault = resultTag.default
                childMin = resultTag.min
                childMax = resultTag.max
                childSize = resultTag.size
                childEnum = resultTag.enum
                childPattern = resultTag.pattern
                childLocation = resultTag.location
                childVersion = resultTag.version
                childMime = resultTag.mime
            }
        }

        node.fields.add(Field("", valueNode))
    }

    val result = tag.validateMap()
    if (!result.valid) {
        throw IllegalArgumentException("validate failed: ${result.error}")
    }

    return node
}

private fun mergeTag(dst: Tag, src: Tag) {
    if (src.isNull) {
        dst.isNull = src.isNull
    }

    if (src.example) {
        dst.example = src.example
    }

    if (src.desc.isNotEmpty()) {
        dst.desc = src.desc
    }

    if (src.type != ValueType.UNKNOWN) {
        dst.type = src.type
    }

    if (src.raw) {
        dst.raw = true
    }

    if (src.nullable) {
        dst.nullable = true
    }

    if (src.allowEmpty) {
        dst.allowEmpty = true
    }

    if (src.unique) {
        dst.unique = true
    }

    if (src.default.isNotEmpty()) {
        dst.default = src.default
    }

    if (src.min.isNotEmpty()) {
        dst.min = src.min
    }

    if (src.max.isNotEmpty()) {
        dst.max = src.max
    }

    if (src.size != 0) {
        dst.size = src.size
    }

    if (src.enum.isNotEmpty()) {
        dst.enum = src.enum
    }

    if (src.pattern.isNotEmpty()) {
        dst.pattern = src.pattern
    }

    if (src.location != 0) {
        dst.location = src.location
    }

    if (src.version != Tag.DEFAULT_VERSION) {
        dst.version = src.version
    }

    if (src.mime.isNotEmpty()) {
        dst.mime = src.mime
    }

    if (src.childDesc.isNotEmpty()) {
        dst.childDesc = src.childDesc
    }

    if (src.childType != ValueType.UNKNOWN) {
        dst.childType = src.childType
    }

    if (src.childRaw) {
        dst.childRaw = true
    }

    if (src.childNullable) {
        dst.childNullable = true
    }

    if (src.childAllowEmpty) {
        dst.childAllowEmpty = true
    }

    if (src.childUnique) {
        dst.childUnique = true
    }

    if (src.childDefault.isNotEmpty()) {
        dst.childDefault = src.childDefault
    }

    if (src.childMin.isNotEmpty()) {
        dst.childMin = src.childMin
    }

    if (src.childMax.isNotEmpty()) {
        dst.childMax = src.childMax
    }

    if (src.childSize != 0) {
        dst.childSize = src.childSize
    }

    if (src.childEnum.isNotEmpty()) {
        dst.childEnum = src.childEnum
    }

    if (src.childPattern.isNotEmpty()) {
        dst.childPattern = src.childPattern
    }

    if (src.childLocation != 0) {
        dst.childLocation = src.childLocation
    }

    if (src.childVersion != Tag.DEFAULT_VERSION) {
        dst.childVersion = src.childVersion
    }

    if (src.childMime.isNotEmpty()) {
        dst.childMime = src.childMime
    }
}

private fun createExampleValue(type: ValueType): Any? {
    return when (type) {
        ValueType.INT8 -> 0.toByte()
        ValueType.INT16 -> 0.toShort()
        ValueType.INT, ValueType.INT32 -> 0
        ValueType.INT64 -> 0L
        ValueType.UINT, ValueType.UINT8, ValueType.UINT16, ValueType.UINT32 -> 0
        ValueType.UINT64 -> BigInteger.ZERO
        ValueType.FLOAT32 -> 0.0f
        ValueType.FLOAT64 -> 0.0
        ValueType.STRING, ValueType.DECIMAL, ValueType.EMAIL, ValueType.URL, ValueType.IP -> ""
        ValueType.BOOL -> false
        ValueType.BYTES, ValueType.IMAGE -> ByteArray(0)
        ValueType.BIGINT -> BigInteger.ZERO
        ValueType.UUID -> UUID(0, 0)
        ValueType.DATETIME -> LocalDateTime.of(1970, 1, 1, 0, 0, 0)
        ValueType.DATE -> LocalDate.of(1970, 1, 1)
        ValueType.TIME -> LocalTime.of(0, 0, 0)
        ValueType.ENUM -> 0
        ValueType.SLICE -> emptyList<Any>()
        ValueType.MAP -> emptyMap<String, Any>()
        ValueType.STRUCT -> emptyMap<String, Any>()
        else -> null
    }
}
