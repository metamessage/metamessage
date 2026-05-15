package io.github.metamessage.jsonc

import io.github.metamessage.ir.Node

class JsoncException(message: String) : Exception(message)

object Jsonc {
    fun parseFromString(s: String): Node? {
        return parseJsonc(s)
    }

    fun parseFromBytes(b: ByteArray): Node? {
        return parseJsonc(String(b, Charsets.UTF_8))
    }

    fun toJsonc(n: Node): String {
        return JsoncPrinter.toString(n)
    }
}

fun toJsonc(node: Node): String {
    val sb = StringBuilder()
    if (node.tag != null && node.tag!!.toString().isNotEmpty()) {
        sb.append("// mm: ${node.tag!!.toString()}\n")
    }
    sb.append(JsoncPrinter.toString(node))
    return sb.toString()
}
