package io.github.metamessage.jsonc

import io.github.metamessage.ast.Node
import io.github.metamessage.ast.Array
import io.github.metamessage.ast.Value
import io.github.metamessage.ast.Object
import io.github.metamessage.ast.ValueType
import io.github.metamessage.ast.Tag
import io.github.metamessage.ast.Doc



class JsoncPrinter {
    companion object {
        private const val INDENT = "    "

        fun toString(node: Node, indentLevel: Int = 0): String {
            return when (node) {
                is Object -> objectToString(node, indentLevel)
                is Array -> arrayToString(node, indentLevel)
                is Value -> valueToString(node)
                is Doc -> docToString(node, indentLevel)
            }
        }

        private fun docToString(doc: Doc, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("{\n")

            for (field in doc.fields) {
                for (i in 0 until indentLevel + 1) sb.append(INDENT)

                val fieldValue = field.value
                if (fieldValue is Value && fieldValue.tag != null && fieldValue.tag!!.toString().isNotEmpty()) {
                    sb.append("// mm: ${fieldValue.tag!!.toString()}\n")
                    for (i in 0 until indentLevel + 1) sb.append(INDENT)
                }

                sb.append("\"${field.key}\": ")
                sb.append(toString(fieldValue, indentLevel + 1))
                sb.append(",\n")
            }

            for (i in 0 until indentLevel) sb.append(INDENT)
            sb.append("}")
            return sb.toString()
        }

        private fun objectToString(obj: Object, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("{\n")

            for (field in obj.fields) {
                for (i in 0 until indentLevel + 1) sb.append(INDENT)

                val fieldValue = field.value
                if (fieldValue.tag != null && fieldValue.tag!!.toString().isNotEmpty()) {
                    sb.append("// mm: ${fieldValue.tag!!.toString()}\n")
                    for (i in 0 until indentLevel + 1) sb.append(INDENT)
                }
                sb.append("\"${field.key}\": ")
                sb.append(toString(fieldValue, indentLevel + 1))
                sb.append(",\n")
            }

            for (i in 0 until indentLevel) sb.append(INDENT)
            sb.append("}")
            return sb.toString()
        }

        private fun arrayToString(arr: Array, indentLevel: Int): String {
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

        private fun valueToString(value: Value): String {
            val sb = StringBuilder()
            sb.append(valueToStringOnly(value))
            return sb.toString()
        }

        private fun valueToStringOnly(value: Value): String {
            val tag = value.tag
            val type = tag?.type ?: ValueType.UNKNOWN

            val needsQuotes = when (type) {
                ValueType.STRING,
                ValueType.BYTES,
                ValueType.DATETIME,
                ValueType.DATE,
                ValueType.TIME,
                ValueType.UUID,
                ValueType.IP,
                ValueType.URL,
                ValueType.EMAIL,
                ValueType.ENUM -> true
                else -> false
            }

            return if (needsQuotes) {
                "\"${value.text}\""
            } else {
                value.text
            }
        }

        fun toCompactString(node: Node): String {
            return when (node) {
                is Object -> compactObject(node)
                is Array -> compactArray(node)
                is Value -> compactValue(node)
                is Doc -> compactDoc(node)
            }
        }

        private fun compactDoc(doc: Doc): String {
            val fields = doc.fields.map { "\"${it.key}\": ${toCompactString(it.value)}" }
            return "{${fields.joinToString(",")}}"
        }

        private fun compactObject(obj: Object): String {
            val fields = obj.fields.map { "\"${it.key}\": ${toCompactString(it.value)}" }
            return "{${fields.joinToString(",")}}"
        }

        private fun compactArray(arr: Array): String {
            val items = arr.items.map { toCompactString(it) }
            return "[${items.joinToString(",")}]"
        }

        private fun compactValue(value: Value): String {
            return valueToStringOnly(value)
        }
    }
}

