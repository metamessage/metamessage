package io.github.metamessage.jsonc

import io.github.metamessage.core.CamelToSnake
import io.github.metamessage.ir.Array as AstArray
import io.github.metamessage.ir.Field
import io.github.metamessage.ir.Node
import io.github.metamessage.ir.Object as AstObject
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.Value
import io.github.metamessage.ir.ValueType
import java.math.BigInteger
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.format.DateTimeFormatter
import java.util.Base64

class JsoncParser(private val tokens: List<JsoncToken>) {
    private var pos: Int = 0
    private val pendingComments = mutableListOf<JsoncToken>()
    private var depth = 0

    private fun peek(): JsoncToken {
        return if (pos >= tokens.size) tokens.last() else tokens[pos]
    }

    private fun next(): JsoncToken {
        val t = peek()
        pos++
        return t
    }

    private fun consumeCommentsFor(anchorLine: Int): Tag? {
        if (pendingComments.isEmpty()) return null

        val last = pendingComments.last()
        if (anchorLine - last.line > 1) {
            pendingComments.clear()
            return null
        }

        var merged: Tag? = null
        for (ct in pendingComments) {
            val parsed = tagFromComment(ct.literal)
            if (parsed != null) {
                merged = mergeTag(merged, parsed)
            }
        }

        pendingComments.clear()
        return merged
    }

    fun parse(): Node? {
        var result: Node? = null
        while (true) {
            val tok = peek()
            if (tok.type == JsoncTokenType.EOF) {
                break
            }

            if (tok.type == JsoncTokenType.LeadingComment) {
                if (pendingComments.isNotEmpty()) {
                    val last = pendingComments.last()
                    if (tok.line - last.line > 1) {
                        pendingComments.clear()
                    }
                }
                pendingComments.add(next())
                continue
            }

            if (tok.type == JsoncTokenType.TrailingComment) {
                if (result != null) {
                    val parsed = tagFromComment(tok.literal)
                    if (parsed != null) {
                        mergeNodeTag(result, parsed)
                    }
                }
                next()
                continue
            }

            result = parse("")
        }
        return result
    }

    private fun parse(path: String): Node? {
        while (true) {
            val tok = next()
            when (tok.type) {
                JsoncTokenType.EOF -> return null
                JsoncTokenType.LBrace -> return parseObject(tok.line, path)
                JsoncTokenType.LBracket -> return parseArray(tok.line, path)
                JsoncTokenType.String -> {
                    val tag = consumeCommentsFor(tok.line) ?: Tag()
                    val text = tok.literal

                    if (tag.type == ValueType.UNKNOWN) {
                        tag.type = ValueType.STR
                    }

                    val result =
                            when (tag.type) {
                                ValueType.STR -> parseStringValue(text, tag)
                                ValueType.BYTES -> parseBytesValue(text, tag)
                                ValueType.DATETIME -> parseDateTimeValue(text, tag)
                                ValueType.DATE -> parseDateValue(text, tag)
                                ValueType.TIME -> parseTimeValue(text, tag)
                                ValueType.UUID -> parseUUIDValue(text, tag)
                                ValueType.DECIMAL -> parseDecimalValue(text, tag)
                                ValueType.IP -> parseIPValue(text, tag)
                                ValueType.URL -> parseURLValue(text, tag)
                                ValueType.EMAIL -> parseEmailValue(text, tag)
                                ValueType.ENUM -> parseEnumValue(text, tag)
                                ValueType.IMAGE -> parseImageValue(text, tag)
                                else ->
                                        throw JsoncException(
                                                "unsupported type ${tag.type} for string literal"
                                        )
                            }

                    return Value(
                            data = result.first,
                            text = result.second ?: "",
                            tag = tag,
                            path = path
                    )
                }
                JsoncTokenType.Number -> {
                    val tag = consumeCommentsFor(tok.line) ?: Tag()
                    val text = tok.literal

                    val result =
                            when {
                                text.contains(".") -> parseFloatValue(text, tag)
                                text.startsWith("-") -> parseNegativeValue(text, tag)
                                else -> parsePositiveValue(text, tag)
                            }

                    val data =
                            if (result.first is Int) {
                                (result.first as Int).toLong()
                            } else {
                                result.first
                            }
                    return Value(data = data, text = result.second ?: "", tag = tag, path = path)
                }
                JsoncTokenType.True -> {
                    val tag = consumeCommentsFor(tok.line) ?: Tag()
                    if (tag.type == ValueType.UNKNOWN) {
                        tag.type = ValueType.BOOL
                    }

                    if (tag.type != ValueType.BOOL) {
                        throw JsoncException("unsupported type ${tag.type} for boolean literal")
                    }

                    if (tag.isNull) {
                        throw JsoncException("bool must false when bool is null")
                    }

                    val result = tag.validateBool(true)
                    if (!result.valid) {
                        throw JsoncException(result.error ?: "Boolean validation failed")
                    }

                    return Value(data = true, text = "true", tag = tag, path = path)
                }
                JsoncTokenType.False -> {
                    val tag = consumeCommentsFor(tok.line) ?: Tag()
                    if (tag.type == ValueType.UNKNOWN) {
                        tag.type = ValueType.BOOL
                    }

                    if (tag.type != ValueType.BOOL) {
                        throw JsoncException("unsupported type ${tag.type} for boolean literal")
                    }

                    if (!tag.isNull) {
                        val result = tag.validateBool(false)
                        if (!result.valid) {
                            throw JsoncException(result.error ?: "Boolean validation failed")
                        }
                    }

                    return Value(data = false, text = "false", tag = tag, path = path)
                }
                else -> throw JsoncException("unexpected token ${tok.type}")
            }
        }
    }

    private fun parseStringValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid string: \"$text\", valid: \"\"")
            }
            "" to ""
        } else {
            val result = tag.validateStr(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "String validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseBytesValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid bytes: \"$text\", valid: \"\"")
            }
            ByteArray(0) to ""
        } else {
            val decoded = Base64.getDecoder().decode(text)
            val result = tag.validateBytes(decoded)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Bytes validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseDateTimeValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            val defaultText = "2006-01-02 15:04:05"
            if (text != defaultText) {
                throw JsoncException("invalid datetime: \"$text\", valid: \"$defaultText\"")
            }
            LocalDateTime.parse(defaultText, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")) to
                    defaultText
        } else {
            val dateTime =
                    LocalDateTime.parse(text, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))
            val result = tag.validateDatetime(dateTime)
            if (!result.valid) {
                throw JsoncException(result.error ?: "DateTime validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseDateValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            val defaultText = "2006-01-02"
            if (text != defaultText) {
                throw JsoncException("invalid date: \"$text\", valid: \"$defaultText\"")
            }
            LocalDate.parse(defaultText) to defaultText
        } else {
            val date = LocalDate.parse(text)
            val result = tag.validateDate(date)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Date validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseTimeValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            val defaultText = "15:04:05"
            if (text != defaultText) {
                throw JsoncException("invalid time: \"$text\", valid: \"$defaultText\"")
            }
            LocalTime.parse(defaultText) to defaultText
        } else {
            val time = LocalTime.parse(text)
            val result = tag.validateTime(time)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Time validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseUUIDValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid uuid: \"$text\", valid: \"\"")
            }
            ByteArray(16) to ""
        } else {
            val result = tag.validateUUID(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "UUID validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseDecimalValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid decimal: \"$text\", valid: \"\"")
            }
            "" to ""
        } else {
            val result = tag.validateDecimal(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Decimal validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseIPValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid ip: \"$text\", valid: \"\"")
            }
            text to ""
        } else {
            val result = tag.validateIP(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "IP validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseURLValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid url: \"$text\", valid: \"\"")
            }
            text to ""
        } else {
            val result = tag.validateURL(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "URL validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseEmailValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid email: \"$text\", valid: \"\"")
            }
            "" to ""
        } else {
            val result = tag.validateEmail(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Email validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseEnumValue(text: String, tag: Tag): Pair<Any?, String?> {
        if (tag.enum.isEmpty()) {
            throw JsoncException("enum empty")
        }

        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid enum: \"$text\", valid: \"\"")
            }
            -1 to ""
        } else {
            val result = tag.validateEnum(text)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Enum validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseImageValue(text: String, tag: Tag): Pair<Any?, String?> {
        return if (tag.isNull) {
            if (text != "") {
                throw JsoncException("invalid image: \"$text\", valid: \"\"")
            }
            ByteArray(0) to ""
        } else {
            val decoded = Base64.getDecoder().decode(text)
            val result = tag.validateImage(decoded)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Image validation failed")
            }
            result.data to result.text
        }
    }

    private fun parseFloatValue(text: String, tag: Tag): Pair<Any?, String?> {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.F64
        }

        return when (tag.type) {
            ValueType.F32 -> {
                if (tag.isNull) {
                    if (text != "0.0") {
                        throw JsoncException("invalid float32: $text, valid: 0.0")
                    }
                    0.0f to "0.0"
                } else {
                    val value = text.toFloat()
                    val result = tag.validateF32(value)
                    if (!result.valid) {
                        throw JsoncException(result.error ?: "Float32 validation failed")
                    }
                    result.data to result.text
                }
            }
            ValueType.F64 -> {
                if (tag.isNull) {
                    if (text != "0.0") {
                        throw JsoncException("invalid float64: $text, valid: 0.0")
                    }
                    0.0 to "0.0"
                } else {
                    val value = text.toDouble()
                    val result = tag.validateF64(value)
                    if (!result.valid) {
                        throw JsoncException(result.error ?: "Float64 validation failed")
                    }
                    result.data to result.text
                }
            }
            else -> throw JsoncException("unsupported numeric type ${tag.type} for float literal")
        }
    }

    private fun parseNegativeValue(text: String, tag: Tag): Pair<Any?, String?> {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.I
        }

        return when (tag.type) {
            ValueType.I -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int: $text, valid: 0")
                    0 to "0"
                } else {
                    val value = text.toInt()
                    val result = tag.validateI(value)
                    if (!result.valid) throw JsoncException(result.error ?: "Int validation failed")
                    result.data to result.text
                }
            }
            ValueType.I8 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int8: $text, valid: 0")
                    0.toByte() to "0"
                } else {
                    val value = text.toByte()
                    val result = tag.validateI8(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int8 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I16 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int16: $text, valid: 0")
                    0.toShort() to "0"
                } else {
                    val value = text.toShort()
                    val result = tag.validateI16(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int16 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I32 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int32: $text, valid: 0")
                    0 to "0"
                } else {
                    val value = text.toInt()
                    val result = tag.validateI32(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int32 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I64 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int64: $text, valid: 0")
                    0L to "0"
                } else {
                    val value = text.toLong()
                    val result = tag.validateI64(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int64 validation failed")
                    result.data to result.text
                }
            }
            ValueType.BIGINT -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid bigint: $text, valid: 0")
                    BigInteger.ZERO to "0"
                } else {
                    val value = BigInteger(text)
                    val result = tag.validateBigint(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "BigInt validation failed")
                    result.data to result.text
                }
            }
            else ->
                    throw JsoncException(
                            "unsupported numeric type ${tag.type} for negative literal"
                    )
        }
    }

    private fun parsePositiveValue(text: String, tag: Tag): Pair<Any?, String?> {
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.I
        }

        return when (tag.type) {
            ValueType.I -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int: $text, valid: 0")
                    0 to "0"
                } else {
                    val value = text.toInt()
                    val result = tag.validateI(value)
                    if (!result.valid) throw JsoncException(result.error ?: "Int validation failed")
                    result.data to result.text
                }
            }
            ValueType.I8 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int8: $text, valid: 0")
                    0.toByte() to "0"
                } else {
                    val value = text.toByte()
                    val result = tag.validateI8(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int8 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I16 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int16: $text, valid: 0")
                    0.toShort() to "0"
                } else {
                    val value = text.toShort()
                    val result = tag.validateI16(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int16 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I32 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int32: $text, valid: 0")
                    0 to "0"
                } else {
                    val value = text.toInt()
                    val result = tag.validateI32(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int32 validation failed")
                    result.data to result.text
                }
            }
            ValueType.I64 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid int64: $text, valid: 0")
                    0L to "0"
                } else {
                    val value = text.toLong()
                    val result = tag.validateI64(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Int64 validation failed")
                    result.data to result.text
                }
            }
            ValueType.U -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid uint: $text, valid: 0")
                    0L to "0"
                } else {
                    val value = text.toLong()
                    val result = tag.validateU(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Uint validation failed")
                    result.data to result.text
                }
            }
            ValueType.U8 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid uint8: $text, valid: 0")
                    0.toShort() to "0"
                } else {
                    val value = text.toShort()
                    val result = tag.validateU8(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Uint8 validation failed")
                    result.data to result.text
                }
            }
            ValueType.U16 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid uint16: $text, valid: 0")
                    0 to "0"
                } else {
                    val value = text.toInt()
                    val result = tag.validateU16(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Uint16 validation failed")
                    result.data to result.text
                }
            }
            ValueType.U32 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid uint32: $text, valid: 0")
                    0L to "0"
                } else {
                    val value = text.toLong()
                    val result = tag.validateU32(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Uint32 validation failed")
                    result.data to result.text
                }
            }
            ValueType.U64 -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid uint64: $text, valid: 0")
                    BigInteger.ZERO to "0"
                } else {
                    val value = BigInteger(text)
                    val result = tag.validateU64(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "Uint64 validation failed")
                    result.data to result.text
                }
            }
            ValueType.BIGINT -> {
                if (tag.isNull) {
                    if (text != "0") throw JsoncException("invalid bigint: $text, valid: 0")
                    BigInteger.ZERO to "0"
                } else {
                    val value = BigInteger(text)
                    val result = tag.validateBigint(value)
                    if (!result.valid)
                            throw JsoncException(result.error ?: "BigInt validation failed")
                    result.data to result.text
                }
            }
            else -> throw JsoncException("unsupported numeric type ${tag.type}")
        }
    }

    private fun parseObject(openLine: Int, path: String): AstObject {
        depth++
        if (depth > 32) {
            throw JsoncException("max depth: 32")
        }

        var tag = consumeCommentsFor(openLine) ?: Tag()
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = ValueType.OBJ
        }

        var currentPath = path
        if (tag.name.isNotEmpty()) {
            currentPath = if (path.isEmpty()) tag.name else "$path.${tag.name}"
        }

        val obj = AstObject(tag = tag, path = currentPath)

        var valNode: Node? = null
        while (true) {
            val tok = peek()
            if (tok.type == JsoncTokenType.EOF) break
            if (tok.type == JsoncTokenType.RBrace) {
                next()
                break
            }

            if (tok.type == JsoncTokenType.LeadingComment) {
                if (pendingComments.isNotEmpty()) {
                    val last = pendingComments.last()
                    if (tok.line - last.line > 1) {
                        pendingComments.clear()
                    }
                }
                pendingComments.add(next())
                continue
            }

            if (tok.type == JsoncTokenType.TrailingComment) {
                if (valNode != null) {
                    val parsed = tagFromComment(tok.literal)
                    if (parsed != null) {
                        mergeNodeTag(valNode, parsed)
                    }
                }
                next()
                continue
            }

            val key = next()
            if (key.type != JsoncTokenType.String) {
                throw JsoncException("expect string key")
            }
            val keyStr = CamelToSnake.convert(key.literal)

            next()
            val fieldPath = "$currentPath.$keyStr"
            valNode = parse(fieldPath) ?: continue

            val childTag = valNode.tag
            if (childTag != null && childTag.type == ValueType.MAP) {
                childTag.inheritFromArrayParent(tag)
            }

            obj.fields.add(Field(keyStr, valNode))

            if (peek().type == JsoncTokenType.Comma) {
                next()
            }
        }

        when (tag.type) {
            ValueType.MAP -> {
                val result = tag.validateMap()
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Map validation failed")
                }
            }
            ValueType.OBJ -> {
                val result = tag.validateObj()
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Struct validation failed")
                }
            }
            else -> {}
        }

        depth--
        return obj
    }

    private fun parseArray(openLine: Int, path: String): AstArray {
        depth++
        if (depth > 32) {
            throw JsoncException("max depth: 32")
        }

        var tag = consumeCommentsFor(openLine) ?: Tag()
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = if (tag.size > 0) ValueType.ARRAY else ValueType.VEC
        }

        var currentPath = path
        if (tag.name.isNotEmpty()) {
            currentPath = "$path.${tag.name}"
        }

        val arr = AstArray(tag = tag, path = currentPath)

        var item: Node? = null
        var index = 0
        while (true) {
            val tok = peek()
            if (tok.type == JsoncTokenType.EOF) break
            if (tok.type == JsoncTokenType.RBracket) {
                next()
                break
            }

            if (tok.type == JsoncTokenType.LeadingComment) {
                if (pendingComments.isNotEmpty()) {
                    val last = pendingComments.last()
                    if (tok.line - last.line > 1) {
                        pendingComments.clear()
                    }
                }
                pendingComments.add(next())
                continue
            }

            if (tok.type == JsoncTokenType.TrailingComment) {
                if (item != null) {
                    val parsed = tagFromComment(tok.literal)
                    if (parsed != null) {
                        mergeNodeTag(item, parsed)
                    }
                }
                next()
                continue
            }

            val itemPath = "$currentPath[$index]"
            item = parse(itemPath) ?: continue

            val childTag = item.tag
            if (childTag != null) {
                childTag.inheritFromArrayParent(tag)
            }

            arr.items.add(item)
            index++

            if (peek().type == JsoncTokenType.Comma) {
                next()
            }
        }

        when (tag.type) {
            ValueType.ARRAY -> {
                val result = tag.validateArr(arr.items)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Array validation failed")
                }
            }
            ValueType.VEC -> {
                val result = tag.validateVec(arr.items)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Slice validation failed")
                }
            }
            else -> {}
        }

        depth--
        return arr
    }

    private fun tagFromComment(comment: String): Tag? {
        val trimmed = comment.trim()
        if (!trimmed.startsWith("mm:")) return null
        val tagStr = trimmed.removePrefix("mm:").trim()
        if (tagStr.isEmpty()) return null
        return parseTag(tagStr)
    }

    private fun parseTag(tagStr: String): Tag {
        val tag = Tag()
        val parts = tagStr.split(";").map { it.trim() }

        for (part in parts) {
            if (part.isEmpty()) continue
            val kv = part.split("=", limit = 2)
            val key = kv[0].lowercase()
            val value = if (kv.size > 1) kv[1].trim() else ""

            when (key) {
                "is_null" -> {
                    tag.isNull = true
                    tag.nullable = true
                }
                "example" -> {}
                "desc" -> tag.desc = value
                "type" -> {
                    tag.type =
                            when (value) {
                                "str" -> ValueType.STR
                                "i" -> ValueType.I
                                "i8" -> ValueType.I8
                                "i16" -> ValueType.I16
                                "i32" -> ValueType.I32
                                "i64" -> ValueType.I64
                                "u" -> ValueType.U
                                "u8" -> ValueType.U8
                                "u16" -> ValueType.U16
                                "u32" -> ValueType.U32
                                "u64" -> ValueType.U64
                                "f32" -> ValueType.F32
                                "f64" -> ValueType.F64
                                "bool" -> ValueType.BOOL
                                "bytes" -> ValueType.BYTES
                                "bigint" -> ValueType.BIGINT
                                "datetime" -> ValueType.DATETIME
                                "date" -> ValueType.DATE
                                "time" -> ValueType.TIME
                                "uuid" -> ValueType.UUID
                                "decimal" -> ValueType.DECIMAL
                                "ip" -> ValueType.IP
                                "url" -> ValueType.URL
                                "email" -> ValueType.EMAIL
                                "enum" -> ValueType.ENUM
                                "arr" -> ValueType.ARRAY
                                "vec" -> ValueType.VEC
                                "obj" -> ValueType.OBJ
                                "map" -> ValueType.MAP
                                else -> ValueType.UNKNOWN
                            }
                }
                "nullable" -> tag.nullable = true
                "raw" -> tag.raw = true
                "allow_empty" -> tag.allowEmpty = true
                "unique" -> tag.unique = true
                "default" -> tag.default = value
                "min" -> tag.min = value
                "max" -> tag.max = value
                "size" -> tag.size = value.toIntOrNull() ?: 0
                "enum" -> {
                    tag.type = ValueType.ENUM
                    tag.enum = value
                }
                "pattern" -> tag.pattern = value
                "location" -> tag.location = value.toIntOrNull() ?: 0
                "version" -> tag.version = value.toIntOrNull() ?: 0
                "mime" -> tag.mime = value
            }
        }
        return tag
    }

    private fun mergeTag(a: Tag?, b: Tag): Tag {
        if (a == null) return b
        if (a.type != ValueType.UNKNOWN) b.type = a.type
        if (a.desc.isNotEmpty()) b.desc = a.desc
        if (a.nullable) b.nullable = true
        if (a.isNull) b.isNull = true
        if (a.default.isNotEmpty()) b.default = a.default
        if (a.min.isNotEmpty()) b.min = a.min
        if (a.max.isNotEmpty()) b.max = a.max
        if (a.size != 0) b.size = a.size
        if (a.enum.isNotEmpty()) b.enum = a.enum
        if (a.pattern.isNotEmpty()) b.pattern = a.pattern
        if (a.location != 0) b.location = a.location
        if (a.version != 0) b.version = a.version
        if (a.mime.isNotEmpty()) b.mime = a.mime
        return b
    }

    private fun mergeNodeTag(n: Node, parsed: Tag) {
        val existing = n.tag
        val merged = mergeTag(existing, parsed)
        when (n) {
            is Value -> n.tag = merged
            is AstObject -> n.tag = merged
            is AstArray -> n.tag = merged
            is io.github.metamessage.ir.Doc -> n.tag = merged
        }
    }
}

fun parseJsonc(source: String): Node? {
    val scanner = JsoncScanner(source)
    val tokens = mutableListOf<JsoncToken>()
    while (true) {
        val tok = scanner.nextToken()
        tokens.add(tok)
        if (tok.type == JsoncTokenType.EOF) break
    }
    val parser = JsoncParser(tokens)
    return parser.parse()
}
