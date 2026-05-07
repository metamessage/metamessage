package io.github.metamessage.jsonc

class JsoncException(message: String) : Exception(message)

object Jsonc {
    fun parseFromString(s: String): JsoncNode {
        return parseJsonc(s)
    }

    fun parseFromBytes(b: ByteArray): JsoncNode {
        return parseJsonc(String(b, Charsets.UTF_8))
    }

    fun toString(n: JsoncNode): String {
        return JsoncPrinter.toString(n)
    }

    fun toCompactString(n: JsoncNode): String {
        return JsoncPrinter.toCompactString(n)
    }

    fun print(n: JsoncNode) {
        println(toString(n))
    }
}