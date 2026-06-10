package io.github.metamessage.jsonc

import io.github.metamessage.ir.Doc
import io.github.metamessage.ir.Node
import io.github.metamessage.ir.NodeArray
import io.github.metamessage.ir.NodeNull
import io.github.metamessage.ir.NodeObject as AstObject
import io.github.metamessage.ir.NodeScalar
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.ValueType

class JsoncPrinter {
    companion object {
        private const val INDENT = "\t"

        fun toString(node: Node, indentLevel: Int = 0): String {
            return when (node) {
                is AstObject -> objectToString(node, indentLevel)
                is NodeArray -> arrayToString(node, indentLevel)
                is NodeScalar -> valueToString(node)
                is NodeNull -> "null"
                is Doc -> docToString(node, indentLevel)
            }
        }

        private fun docToString(doc: Doc, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("{\n")

            for (field in doc.fields) {
                writeLeadingComments(sb, field.value.tag, indentLevel + 1)

                writeIndent(sb, indentLevel + 1)
                sb.append("\"${field.key}\": ")
                sb.append(toString(field.value, indentLevel + 1))
                sb.append(",\n")
            }

            writeIndent(sb, indentLevel)
            sb.append("}")
            return sb.toString()
        }

        private fun objectToString(obj: AstObject, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("{\n")

            for (field in obj.fields) {
                writeLeadingComments(sb, field.value.tag, indentLevel + 1)

                writeIndent(sb, indentLevel + 1)
                sb.append("\"${field.key}\": ")
                sb.append(toString(field.value, indentLevel + 1))
                sb.append(",\n")
            }

            writeIndent(sb, indentLevel)
            sb.append("}")
            return sb.toString()
        }

        private fun arrayToString(arr: NodeArray, indentLevel: Int): String {
            val sb = StringBuilder()
            sb.append("[\n")

            for (item in arr.items) {
                writeLeadingComments(sb, item.tag, indentLevel + 1)

                writeIndent(sb, indentLevel + 1)
                sb.append(toString(item, indentLevel + 1))
                sb.append(",\n")
            }

            writeIndent(sb, indentLevel)
            sb.append("]")
            return sb.toString()
        }

        private fun valueToString(value: NodeScalar): String {
            return valueToStringOnly(value)
        }

        private fun valueToStringOnly(value: NodeScalar): String {
            val tag = value.tag
            val type = tag?.type ?: ValueType.UNKNOWN

            val needsQuotes =
                    when (type) {
                        ValueType.STR,
                        ValueType.BYTES,
                        ValueType.DATETIME,
                        ValueType.DATE,
                        ValueType.TIME,
                        ValueType.UUID,
                        ValueType.IP,
                        ValueType.URL,
                        ValueType.EMAIL,
                        ValueType.ENUMS,
                        ValueType.MEDIA -> true
                        else -> false
                    }

            return if (needsQuotes) {
                "\"${value.text}\""
            } else {
                value.text
            }
        }

        private fun writeIndent(sb: StringBuilder, indentLevel: Int) {
            for (i in 0 until indentLevel) {
                sb.append(INDENT)
            }
        }

        private fun writeLeadingComments(sb: StringBuilder, tag: Tag?, indentLevel: Int) {
            val tagStr = tag?.toString() ?: ""
            if (tagStr.isNotEmpty()) {
                sb.append("\n")
                writeIndent(sb, indentLevel)
                sb.append("// mm: $tagStr\n")
            }
        }

        fun toCompactString(node: Node): String {
            return when (node) {
                is AstObject -> compactObject(node)
                is NodeArray -> compactArray(node)
                is NodeScalar -> compactValue(node)
                is NodeNull -> "null"
                is Doc -> compactDoc(node)
            }
        }

        private fun compactDoc(doc: Doc): String {
            val fields = doc.fields.map { "\"${it.key}\": ${toCompactString(it.value)}" }
            return "{${fields.joinToString(",")}}"
        }

        private fun compactObject(obj: AstObject): String {
            val fields = obj.fields.map { "\"${it.key}\": ${toCompactString(it.value)}" }
            return "{${fields.joinToString(",")}}"
        }

        private fun compactArray(arr: NodeArray): String {
            val items = arr.items.map { toCompactString(it) }
            return "[${items.joinToString(",")}]"
        }

        private fun compactValue(value: NodeScalar): String {
            return valueToStringOnly(value)
        }
    }
}
