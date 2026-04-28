package io.metamessage.jsonc

class JsoncPrinter {
    companion object {
        private const val INDENT = "    "

        fun toString(node: JsoncNode, indentLevel: Int = 0): String {
            return when (node) {
                is JsoncObject -> objectToString(node, indentLevel)
                is JsoncArray -> arrayToString(node, indentLevel)
                is JsoncValue -> valueToString(node, indentLevel)
            }
        }

        private fun objectToString(obj: JsoncObject, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("{\n")

            for (field in obj.fields) {
                for (i in 0 until indentLevel + 1) sb.append(INDENT)

                val fieldValue = field.value
                if (fieldValue is JsoncValue && fieldValue.tag != null && fieldValue.tag!!.desc.isNotEmpty()) {
                    sb.append("// mm: ${tagToString(fieldValue.tag!!)}\n")
                    for (i in 0 until indentLevel + 1) sb.append(INDENT)
                }

                sb.append("\"${field.key}\": ")
                sb.append(toString(fieldValue ?: JsoncValue(), indentLevel + 1))
                sb.append(",\n")
            }

            for (i in 0 until indentLevel) sb.append(INDENT)
            sb.append("}")
            return sb.toString()
        }

        private fun arrayToString(arr: JsoncArray, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("[\n")

            for (item in arr.items) {
                for (i in 0 until indentLevel + 1) sb.append(INDENT)
                sb.append(toString(item, indentLevel + 1))
                sb.append(",\n")
            }

            for (i in 0 until indentLevel) sb.append(INDENT)
            sb.append("]")
            return sb.toString()
        }

        private fun valueToString(value: JsoncValue, indentLevel: Int): String {
            val sb = StringBuilder()
            if (value.tag != null && value.tag!!.desc.isNotEmpty()) {
                sb.append("// mm: ${tagToString(value.tag!!)}\n")
                for (i in 0 until indentLevel) sb.append(INDENT)
            }
            sb.append(valueToStringOnly(value))
            return sb.toString()
        }

        private fun valueToStringOnly(value: JsoncValue): String {
            val tag = value.tag
            val type = tag?.type ?: JsoncValueType.Unknown

            val needsQuotes = when (type) {
                JsoncValueType.String,
                JsoncValueType.Bytes,
                JsoncValueType.DateTime,
                JsoncValueType.Date,
                JsoncValueType.Time,
                JsoncValueType.UUID,
                JsoncValueType.IP,
                JsoncValueType.URL,
                JsoncValueType.Email,
                JsoncValueType.Enum -> true
                else -> false
            }

            return if (needsQuotes) {
                "\"${value.text}\""
            } else {
                value.text
            }
        }

        private fun tagToString(tag: JsoncTag): String {
            val parts = mutableListOf<String>()
            if (tag.type != JsoncValueType.Unknown) {
                parts.add("type=${typeToString(tag.type)}")
            }
            if (tag.desc.isNotEmpty()) {
                parts.add("desc=${tag.desc}")
            }
            if (tag.nullable) {
                parts.add("nullable")
            }
            if (tag.isNull) {
                parts.add("is_null")
            }
            if (tag.raw) {
                parts.add("raw")
            }
            if (tag.defaultValue.isNotEmpty()) {
                parts.add("default=${tag.defaultValue}")
            }
            if (tag.enum.isNotEmpty()) {
                parts.add("enum=${tag.enum}")
            }
            return parts.joinToString("; ")
        }

        private fun typeToString(type: JsoncValueType): String {
            return when (type) {
                JsoncValueType.String -> "str"
                JsoncValueType.Int -> "i"
                JsoncValueType.Int8 -> "i8"
                JsoncValueType.Int16 -> "i16"
                JsoncValueType.Int32 -> "i32"
                JsoncValueType.Int64 -> "i64"
                JsoncValueType.Uint -> "u"
                JsoncValueType.Uint8 -> "u8"
                JsoncValueType.Uint16 -> "u16"
                JsoncValueType.Uint32 -> "u32"
                JsoncValueType.Uint64 -> "u64"
                JsoncValueType.Float32 -> "f32"
                JsoncValueType.Float64 -> "f64"
                JsoncValueType.Bool -> "bool"
                JsoncValueType.Bytes -> "bytes"
                JsoncValueType.BigInt -> "bi"
                JsoncValueType.DateTime -> "datetime"
                JsoncValueType.Date -> "date"
                JsoncValueType.Time -> "time"
                JsoncValueType.UUID -> "uuid"
                JsoncValueType.Decimal -> "decimal"
                JsoncValueType.IP -> "ip"
                JsoncValueType.URL -> "url"
                JsoncValueType.Email -> "email"
                JsoncValueType.Enum -> "enum"
                JsoncValueType.Array -> "arr"
                JsoncValueType.Struct -> "obj"
                JsoncValueType.Null -> "null"
                JsoncValueType.Unknown -> "unknown"
                else -> "unknown"
            }
        }

        fun toCompactString(node: JsoncNode): String {
            return when (node) {
                is JsoncObject -> compactObject(node)
                is JsoncArray -> compactArray(node)
                is JsoncValue -> compactValue(node)
            }
        }

        private fun compactObject(obj: JsoncObject): String {
            val fields = obj.fields.map { "\"${it.key}\": ${toCompactString(it.value ?: JsoncValue())}" }
            return "{${fields.joinToString(",")}}"
        }

        private fun compactArray(arr: JsoncArray): String {
            val items = arr.items.map { toCompactString(it) }
            return "[${items.joinToString(",")}]"
        }

        private fun compactValue(value: JsoncValue): String {
            return valueToStringOnly(value)
        }
    }
}

fun printJsonc(node: JsoncNode): String = JsoncPrinter.toString(node)
fun printJsoncCompact(node: JsoncNode): String = JsoncPrinter.toCompactString(node)