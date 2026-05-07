package io.github.metamessage.jsonc

import io.github.metamessage.mm.MmValidator
import io.github.metamessage.mm.MmTag
import io.github.metamessage.mm.ValueType

class JsoncParser(private val tokens: List<JsoncToken>) {
    private var pos: Int = 0
    private val pendingComments = mutableListOf<JsoncToken>()
    
    // 转换 JsoncTag 为 MmTag
    private fun convertJsoncTagToMmTag(jsoncTag: JsoncTag): MmTag {
        val mmTag = MmTag()
        mmTag.isNull = jsoncTag.isNull
        mmTag.desc = jsoncTag.desc
        mmTag.type = when (jsoncTag.type) {
            JsoncValueType.String -> ValueType.STRING
            JsoncValueType.Int -> ValueType.INT
            JsoncValueType.Int8 -> ValueType.INT8
            JsoncValueType.Int16 -> ValueType.INT16
            JsoncValueType.Int32 -> ValueType.INT32
            JsoncValueType.Int64 -> ValueType.INT64
            JsoncValueType.Uint -> ValueType.UINT
            JsoncValueType.Uint8 -> ValueType.UINT8
            JsoncValueType.Uint16 -> ValueType.UINT16
            JsoncValueType.Uint32 -> ValueType.UINT32
            JsoncValueType.Uint64 -> ValueType.UINT64
            JsoncValueType.Float32 -> ValueType.FLOAT32
            JsoncValueType.Float64 -> ValueType.FLOAT64
            JsoncValueType.Bool -> ValueType.BOOL
            JsoncValueType.Bytes -> ValueType.BYTES
            JsoncValueType.BigInt -> ValueType.BIGINT
            JsoncValueType.DateTime -> ValueType.DATETIME
            JsoncValueType.Date -> ValueType.DATE
            JsoncValueType.Time -> ValueType.TIME
            JsoncValueType.UUID -> ValueType.UUID
            JsoncValueType.Decimal -> ValueType.DECIMAL
            JsoncValueType.IP -> ValueType.IP
            JsoncValueType.URL -> ValueType.URL
            JsoncValueType.Email -> ValueType.EMAIL
            JsoncValueType.Enum -> ValueType.ENUM
            JsoncValueType.Array -> ValueType.ARRAY
            JsoncValueType.Struct -> ValueType.STRUCT
            else -> ValueType.STRING
        }
        mmTag.raw = jsoncTag.raw
        mmTag.nullable = jsoncTag.nullable
        mmTag.allowEmpty = jsoncTag.allowEmpty
        mmTag.unique = jsoncTag.unique
        mmTag.defaultValue = jsoncTag.defaultValue
        mmTag.min = jsoncTag.min
        mmTag.max = jsoncTag.max
        mmTag.size = jsoncTag.size
        mmTag.enumValues = jsoncTag.enum
        mmTag.pattern = jsoncTag.pattern
        mmTag.locationHours = jsoncTag.location.toIntOrNull() ?: 0
        mmTag.version = jsoncTag.version
        mmTag.mime = jsoncTag.mime
        return mmTag
    }

    private fun peek(): JsoncToken {
        return if (pos >= tokens.size) tokens.last() else tokens[pos]
    }

    private fun next(): JsoncToken {
        val t = peek()
        pos++
        return t
    }

    private fun consumeComments(): JsoncTag? {
        if (pendingComments.isEmpty()) return null

        val comments = pendingComments.toList()
        pendingComments.clear()

        if (comments.isEmpty()) return null

        var merged = JsoncTag()
        for (c in comments) {
            val t = tagFromComment(c.literal)
            if (t != null) {
                merged = mergeTag(merged, t)
            }
        }
        return merged
    }

    fun parse(): JsoncNode {
        val tok = peek()
        if (tok.type == JsoncTokenType.EOF) {
            return JsoncValue()
        }

        return when (tok.type) {
            JsoncTokenType.LBrace -> parseObject()
            JsoncTokenType.LBracket -> parseArray()
            else -> parseValue()
        }
    }

    private fun parseObject(): JsoncObject {
        next()
        val tag = consumeComments()
        val obj = JsoncObject(tag = tag)
        
        if (tag != null) {
            // 验证结构体 tag
            val mmTag = convertJsoncTagToMmTag(tag)
            val result = MmValidator.validateStruct(mmTag)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Struct validation failed")
            }
        }

        while (peek().type != JsoncTokenType.RBrace && peek().type != JsoncTokenType.EOF) {
            val tok = peek()
            if (tok.type == JsoncTokenType.LeadingComment || tok.type == JsoncTokenType.TrailingComment) {
                pendingComments.add(next())
                continue
            }

            val keyToken = next()
            if (keyToken.type != JsoncTokenType.String) break

            next()

            val fieldPath = "${obj.path}.${keyToken.literal}"
            val value = parseValue()
            obj.fields.add(JsoncField(keyToken.literal, value))

            if (peek().type == JsoncTokenType.Comma) {
                next()
            }
        }

        if (peek().type == JsoncTokenType.RBrace) {
            next()
        }

        return obj
    }

    private fun parseArray(): JsoncArray {
        next()
        val tag = consumeComments()
        val arr = JsoncArray(tag = tag)
        var index = 0
        
        if (tag != null) {
            // 验证数组 tag
            val mmTag = convertJsoncTagToMmTag(tag)
            val result = MmValidator.validateStruct(mmTag)
            if (!result.valid) {
                throw JsoncException(result.error ?: "Array validation failed")
            }
        }

        while (peek().type != JsoncTokenType.RBracket && peek().type != JsoncTokenType.EOF) {
            val tok = peek()
            if (tok.type == JsoncTokenType.LeadingComment || tok.type == JsoncTokenType.TrailingComment) {
                pendingComments.add(next())
                continue
            }

            val itemPath = "${arr.path}[$index]"
            val item = parseValue()
            arr.items.add(item)
            index++

            if (peek().type == JsoncTokenType.Comma) {
                next()
            }
        }

        if (peek().type == JsoncTokenType.RBracket) {
            next()
        }

        return arr
    }

    private fun parseValue(): JsoncNode {
        val tok = peek()

        when (tok.type) {
            JsoncTokenType.LBrace -> return parseObject()
            JsoncTokenType.LBracket -> return parseArray()
            else -> {}
        }

        val actualToken = next()
        val tag = consumeComments()

        when (actualToken.type) {
            JsoncTokenType.String -> {
                val effectiveTag = tag ?: JsoncTag().also { it.type = JsoncValueType.String }
                if (effectiveTag.type == JsoncValueType.Unknown) {
                    effectiveTag.type = JsoncValueType.String
                }
                
                // 验证值
                val mmTag = convertJsoncTagToMmTag(effectiveTag)
                val result = MmValidator.validate(actualToken.literal, mmTag)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "String validation failed")
                }
                
                return JsoncValue(data = actualToken.literal, text = actualToken.literal, tag = effectiveTag)
            }
            JsoncTokenType.Number -> {
                val effectiveTag = tag ?: JsoncTag()
                if (effectiveTag.type == JsoncValueType.Unknown) {
                    effectiveTag.type = if (actualToken.literal.contains(".")) JsoncValueType.Float64 else JsoncValueType.Int
                }
                val data = if (actualToken.literal.contains(".")) {
                    actualToken.literal.toDoubleOrNull()
                } else {
                    actualToken.literal.toLongOrNull()
                }
                
                // 验证值
                val mmTag = convertJsoncTagToMmTag(effectiveTag)
                val result = MmValidator.validate(data, mmTag)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Number validation failed")
                }
                
                return JsoncValue(data = data, text = actualToken.literal, tag = effectiveTag)
            }
            JsoncTokenType.True -> {
                val effectiveTag = tag ?: JsoncTag()
                effectiveTag.type = JsoncValueType.Bool
                
                // 验证值
                val mmTag = convertJsoncTagToMmTag(effectiveTag)
                val result = MmValidator.validate(true, mmTag)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Boolean validation failed")
                }
                
                return JsoncValue(data = true, text = "true", tag = effectiveTag)
            }
            JsoncTokenType.False -> {
                val effectiveTag = tag ?: JsoncTag()
                effectiveTag.type = JsoncValueType.Bool
                
                // 验证值
                val mmTag = convertJsoncTagToMmTag(effectiveTag)
                val result = MmValidator.validate(false, mmTag)
                if (!result.valid) {
                    throw JsoncException(result.error ?: "Boolean validation failed")
                }
                
                return JsoncValue(data = false, text = "false", tag = effectiveTag)
            }
            JsoncTokenType.Null -> {
                val effectiveTag = tag ?: JsoncTag()
                effectiveTag.isNull = true
                return JsoncValue(data = null, text = "null", tag = effectiveTag)
            }
            else -> {
                return JsoncValue(data = null, text = "", tag = tag)
            }
        }
    }

    private fun tagFromComment(comment: String): JsoncTag? {
        val trimmed = comment.trim()
        if (!trimmed.startsWith("mm:")) return null
        val tagStr = trimmed.removePrefix("mm:").trim()
        if (tagStr.isEmpty()) return null
        return parseTag(tagStr)
    }

    private fun parseTag(tagStr: String): JsoncTag {
        val tag = JsoncTag()
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
                    tag.type = when (value) {
                        "str" -> JsoncValueType.String
                        "i" -> JsoncValueType.Int
                        "i8" -> JsoncValueType.Int8
                        "i16" -> JsoncValueType.Int16
                        "i32" -> JsoncValueType.Int32
                        "i64" -> JsoncValueType.Int64
                        "u" -> JsoncValueType.Uint
                        "u8" -> JsoncValueType.Uint8
                        "u16" -> JsoncValueType.Uint16
                        "u32" -> JsoncValueType.Uint32
                        "u64" -> JsoncValueType.Uint64
                        "f32" -> JsoncValueType.Float32
                        "f64" -> JsoncValueType.Float64
                        "bool" -> JsoncValueType.Bool
                        "bytes" -> JsoncValueType.Bytes
                        "bi" -> JsoncValueType.BigInt
                        "datetime" -> JsoncValueType.DateTime
                        "date" -> JsoncValueType.Date
                        "time" -> JsoncValueType.Time
                        "uuid" -> JsoncValueType.UUID
                        "decimal" -> JsoncValueType.Decimal
                        "ip" -> JsoncValueType.IP
                        "url" -> JsoncValueType.URL
                        "email" -> JsoncValueType.Email
                        "enum" -> JsoncValueType.Enum
                        "arr" -> JsoncValueType.Array
                        "obj" -> JsoncValueType.Struct
                        else -> JsoncValueType.Unknown
                    }
                }
                "nullable" -> tag.nullable = true
                "raw" -> tag.raw = true
                "allow_empty" -> tag.allowEmpty = true
                "unique" -> tag.unique = true
                "default" -> tag.defaultValue = value
                "min" -> tag.min = value
                "max" -> tag.max = value
                "size" -> tag.size = value.toIntOrNull() ?: 0
                "enum" -> {
                    tag.type = JsoncValueType.Enum
                    tag.enum = value
                }
                "pattern" -> tag.pattern = value
                "location" -> tag.location = value
                "version" -> tag.version = value.toIntOrNull() ?: 0
                "mime" -> tag.mime = value
            }
        }
        return tag
    }

    private fun mergeTag(a: JsoncTag, b: JsoncTag): JsoncTag {
        if (a.type != JsoncValueType.Unknown) b.type = a.type
        if (a.desc.isNotEmpty()) b.desc = a.desc
        if (a.nullable) b.nullable = true
        if (a.isNull) b.isNull = true
        if (a.defaultValue.isNotEmpty()) b.defaultValue = a.defaultValue
        if (a.min.isNotEmpty()) b.min = a.min
        if (a.max.isNotEmpty()) b.max = a.max
        if (a.size != 0) b.size = a.size
        if (a.enum.isNotEmpty()) b.enum = a.enum
        if (a.pattern.isNotEmpty()) b.pattern = a.pattern
        if (a.location.isNotEmpty()) b.location = a.location
        if (a.version != 0) b.version = a.version
        if (a.mime.isNotEmpty()) b.mime = a.mime
        return b
    }
}

fun parseJsonc(source: String): JsoncNode {
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