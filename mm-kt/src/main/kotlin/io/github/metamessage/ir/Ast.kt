package io.github.metamessage.ir

enum class NodeType(val label: String) {
    NodeTypeUnknown("unknown"),
    NodeTypeObject("object"),
    NodeTypeArray("array"),
    NodeTypeValue("value"),
    NodeTypeDoc("doc");

    override fun toString(): String {
        return label
    }
}

fun parseNodeType(s: String): NodeType {
    return when (s) {
        "object" -> NodeType.NodeTypeObject
        "array" -> NodeType.NodeTypeArray
        "value" -> NodeType.NodeTypeValue
        "doc" -> NodeType.NodeTypeDoc
        else -> NodeType.NodeTypeUnknown
    }
}

sealed class Node {
    abstract var tag: Tag?
    abstract var path: String

    abstract fun getType(): NodeType
}

data class Field(val key: String, val value: Node)

data class Object(
        var fields: MutableList<Field> = mutableListOf(),
        override var tag: Tag? = null,
        override var path: String = ""
) : Node() {
    override fun getType(): NodeType = NodeType.NodeTypeObject
}

data class Array(
        var items: MutableList<Node> = mutableListOf(),
        override var tag: Tag? = null,
        override var path: String = ""
) : Node() {
    override fun getType(): NodeType = NodeType.NodeTypeArray
}

data class Value(
        var data: Any? = null,
        var text: String = "",
        override var tag: Tag? = null,
        override var path: String = ""
) : Node() {
    override fun getType(): NodeType = NodeType.NodeTypeValue
}

class Doc(
        var fields: MutableList<Field> = mutableListOf(),
        override var tag: Tag? = null,
        override var path: String = ""
) : Node() {
    override fun getType(): NodeType = NodeType.NodeTypeDoc
}

fun dump(obj: Node): String {
    return dumpImpl(obj, 0)
}

private fun dumpImpl(obj: Any?, depth: Int): String {
    val indent = "  ".repeat(depth)
    val nextDepth = depth + 1

    if (obj == null) return "null"

    when (obj) {
        is String -> return "\"$obj\""
        is Number -> return obj.toString()
        is Boolean -> return obj.toString()
        is Enum<*> -> return obj.name
    }

    if (obj is Collection<*>) {
        val sb = StringBuilder()
        sb.append("[\n")
        obj.forEach {
            sb.append("  ".repeat(nextDepth))
            sb.append(dumpImpl(it, nextDepth))
            sb.append(",\n")
        }
        sb.append(indent)
        sb.append("]")
        return sb.toString()
    }

    if (obj is Map<*, *>) {
        val sb = StringBuilder()
        sb.append("{\n")
        obj.forEach { (k, v) ->
            sb.append("  ".repeat(nextDepth))
            sb.append("$k: ")
            sb.append(dumpImpl(v, nextDepth))
            sb.append(",\n")
        }
        sb.append(indent)
        sb.append("}")
        return sb.toString()
    }

    val className = obj::class.java.name
    if (!className.startsWith("io.github.metamessage.ir")) {
        return "${obj::class.simpleName}(${obj})"
    }

    val sb = StringBuilder()
    sb.append(obj::class.simpleName).append(" {\n")

    obj::class.java.declaredFields.forEach { field ->
        field.isAccessible = true
        val name = field.name
        val value =
                try {
                    field.get(obj)
                } catch (e: Exception) {
                    "err"
                }

        sb.append("  ".repeat(nextDepth))
        sb.append(name).append(": ")
        sb.append(dumpImpl(value, nextDepth))
        sb.append(",\n")
    }

    sb.append(indent).append("}")
    return sb.toString()
}
